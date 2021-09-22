package hyperscan

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

var (
	ErrNoFound    = errors.New("no found")
	ErrUnexpected = errors.New("unexpected")
)

// Expression of pattern.
type Expression string

func (e Expression) String() string { return string(e) }

// Patterns is a set of matching patterns.
type Patterns []*Pattern

// Pattern is a matching pattern.
type Pattern struct {
	Expression             // The expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	Id         int         // The ID number to be associated with the corresponding pattern
	info       *ExprInfo
	ext        *ExprExt
}

// NewPattern returns a new pattern base on expression and compile flags.
func NewPattern(expr string, flags CompileFlag) *Pattern {
	return &Pattern{Expression: Expression(expr), Flags: flags}
}

// IsValid validate the pattern contains a regular expression.
func (p *Pattern) IsValid() bool {
	_, err := p.Info()

	return err == nil
}

// Info provides information about a regular expression.
func (p *Pattern) Info() (*ExprInfo, error) {
	if p.info == nil {
		info, err := hsExpressionInfo(string(p.Expression), p.Flags)
		if err != nil {
			return nil, err
		}

		p.info = info
	}

	return p.info, nil
}

// WithExt is used to set the additional parameters related to an expression.
func (p *Pattern) WithExt(exts ...Ext) *Pattern {
	if p.ext == nil {
		p.ext = new(ExprExt)
	}

	p.ext.With(exts...)

	return p
}

// Ext provides additional parameters related to an expression.
func (p *Pattern) Ext() (*ExprExt, error) {
	if p.ext == nil {
		ext, info, err := hsExpressionExt(string(p.Expression), p.Flags)
		if err != nil {
			return nil, err
		}

		p.ext = ext
		p.info = info
	}

	return p.ext, nil
}

func (p *Pattern) String() string {
	var b strings.Builder

	if p.Id > 0 {
		fmt.Fprintf(&b, "%d:", p.Id)
	}

	fmt.Fprintf(&b, "/%s/%s", p.Expression, p.Flags)

	if p.ext != nil {
		b.WriteString(p.ext.String())
	}

	return b.String()
}

/*
ParsePattern parse pattern from a formated string.

	<integer id>:/<expression>/<flags>

For example, the following pattern will match `test` in the caseless and multi-lines mode

	/test/im

*/
func ParsePattern(s string) (*Pattern, error) {
	var p Pattern

	i := strings.Index(s, ":/")
	j := strings.LastIndex(s, "/")

	if i > 0 && j > i+1 {
		id, err := strconv.Atoi(s[:i])
		if err != nil {
			return nil, fmt.Errorf("invalid pattern id `%s`, %w", s[:i], ErrInvalid)
		}

		p.Id = id
		s = s[i+1:]
	}

	if n := strings.LastIndex(s, "/"); n > 1 && strings.HasPrefix(s, "/") {
		p.Expression = Expression(s[1:n])
		s = s[n+1:]

		if n = strings.Index(s, "{"); n > 0 && strings.HasSuffix(s, "}") {
			ext, err := ParseExprExt(s[n:])
			if err != nil {
				return nil, fmt.Errorf("invalid expression extensions `%s`, %w", s[n:], err)
			}

			p.ext = ext
			s = s[:n]
		}

		flags, err := ParseCompileFlag(s)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern flags `%s`, %w", s, err)
		}

		p.Flags = flags
	} else {
		p.Expression = Expression(s)
	}

	info, err := hsExpressionInfo(string(p.Expression), p.Flags)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern `%s`, %w", p.Expression, err)
	}

	p.info = info

	return &p, nil
}

// ParsePatterns parse lines as `Patterns`.
func ParsePatterns(r io.Reader) (patterns Patterns, err error) {
	s := bufio.NewScanner(r)

	for s.Scan() {
		line := strings.TrimSpace(s.Text())

		if line == "" {
			// skip empty line
			continue
		}

		if strings.HasPrefix(line, "#") {
			// skip comment
			continue
		}

		p, err := ParsePattern(line)
		if err != nil {
			return nil, err
		}

		patterns = append(patterns, p)
	}

	return
}

// Platform is a type containing information on the target platform.
type Platform interface {
	// Information about the target platform which may be used to guide the optimisation process of the compile.
	Tune() TuneFlag

	// Relevant CPU features available on the target platform
	CpuFeatures() CpuFeature
}

// NewPlatform create a new platform information on the target platform.
func NewPlatform(tune TuneFlag, cpu CpuFeature) Platform { return newPlatformInfo(tune, cpu) }

// PopulatePlatform populates the platform information based on the current host.
func PopulatePlatform() Platform {
	platform, _ := hsPopulatePlatform()

	return platform
}

// DatabaseBuilder is a type to help to build up a database.
type DatabaseBuilder struct {
	// Array of patterns to compile.
	Patterns []*Pattern

	// Compiler mode flags that affect the database as a whole. (Default: block mode)
	Mode ModeFlag

	// If not nil, the platform structure is used to determine the target platform for the database.
	// If nil, a database suitable for running on the current host platform is produced.
	Platform Platform
}

// AddExpressions add more expressions to the database.
func (b *DatabaseBuilder) AddExpressions(exprs ...Expression) *DatabaseBuilder {
	for _, expr := range exprs {
		b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Id: len(b.Patterns) + 1})
	}

	return b
}

// AddExpressionWithFlags add more expressions with flags to the database.
func (b *DatabaseBuilder) AddExpressionWithFlags(expr Expression, flags CompileFlag) *DatabaseBuilder {
	b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Flags: flags, Id: len(b.Patterns) + 1})

	return b
}

// Build a database base on the expressions and platform.
func (b *DatabaseBuilder) Build() (Database, error) {
	if b.Patterns == nil {
		return nil, ErrNoFound
	}

	mode := b.Mode

	if mode == 0 {
		mode = BlockMode
	} else if mode == StreamMode {
		som := false

		for _, pattern := range b.Patterns {
			if (pattern.Flags & SomLeftMost) == SomLeftMost {
				som = true
			}
		}

		if som && mode&(SomHorizonSmallMode|SomHorizonMediumMode|SomHorizonLargeMode) == 0 {
			mode |= SomHorizonSmallMode
		}
	}

	platform, _ := b.Platform.(*hsPlatformInfo)

	db, err := hsCompileMulti(b.Patterns, mode, platform)
	if err != nil {
		return nil, err
	}

	switch mode & ModeMask {
	case StreamMode:
		return newStreamDatabase(db), nil
	case VectoredMode:
		return newVectoredDatabase(db), nil
	case BlockMode:
		return newBlockDatabase(db), nil
	default:
		return nil, fmt.Errorf("mode %d, %w", mode, ErrUnexpected)
	}
}

// NewBlockDatabase create a block database base on the patterns.
func NewBlockDatabase(patterns ...*Pattern) (BlockDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: BlockMode}

	db, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return db.(*blockDatabase), err
}

// NewStreamDatabase create a stream database base on the patterns.
func NewStreamDatabase(patterns ...*Pattern) (StreamDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: StreamMode}

	db, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return db.(*streamDatabase), err
}

// NewMediumStreamDatabase create a medium-sized stream database base on the patterns.
func NewMediumStreamDatabase(patterns ...*Pattern) (StreamDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: StreamMode | SomHorizonMediumMode}

	db, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return db.(*streamDatabase), err
}

// NewLargeStreamDatabase create a large-sized stream database base on the patterns.
func NewLargeStreamDatabase(patterns ...*Pattern) (StreamDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: StreamMode | SomHorizonLargeMode}

	db, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return db.(*streamDatabase), err
}

// NewVectoredDatabase create a vectored database base on the patterns.
func NewVectoredDatabase(patterns ...*Pattern) (VectoredDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: VectoredMode}

	db, err := builder.Build()
	if err != nil {
		return nil, err
	}

	return db.(*vectoredDatabase), err
}

// Compile a regular expression and returns, if successful,
// a pattern database in the block mode that can be used to match against text.
func Compile(expr string) (Database, error) {
	db, err := hsCompile(expr, SomLeftMost, BlockMode, nil)
	if err != nil {
		return nil, err
	}

	return newBlockDatabase(db), nil
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular expressions.
func MustCompile(expr string) Database {
	db, err := hsCompile(expr, SomLeftMost, BlockMode, nil)
	if err != nil {
		panic(`Compile(` + Quote(expr) + `): ` + err.Error())
	}

	return newBlockDatabase(db)
}

// Quote returns a quoted string literal representing s.
func Quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}

	return strconv.Quote(s)
}
