package hyperscan

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
	"strings"
)

// The expression of pattern
type Expression string

func (e Expression) String() string { return string(e) }

type Patterns []*Pattern

type Pattern struct {
	Expression             // The expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	Id         int         // The ID number to be associated with the corresponding pattern
	Ext        *ExprExt    // The matching behaviour of a pattern
	info       *ExprInfo
	ext        *ExprExt
}

/// NewPattern returns a new pattern base on expression and compile flags.
func NewPattern(expr string, flags CompileFlag) *Pattern {
	return &Pattern{Expression: Expression(expr), Flags: flags}
}

/// IsValid validate the pattern contains a regular expression.
func (p *Pattern) IsValid() bool {
	_, err := p.Info()

	return err == nil
}

// Provides information about a regular expression.
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

type Ext func(ext *ExprExt)

func NewExprExt(exts ...Ext) (ext *ExprExt) {
	ext = &ExprExt{}
	for _, f := range exts {
		f(ext)
	}
	return
}

// The minimum end offset in the data stream at which this expression should match successfully.
func MinOffset(n uint64) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtMinOffset
		ext.MinOffset = n
	}
}

// The maximum end offset in the data stream at which this expression should match successfully.
func MaxOffset(n uint64) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtMaxOffset
		ext.MaxOffset = n
	}
}

// The minimum match length (from start to end) required to successfully match this expression.
func MinLength(n uint64) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtMinLength
		ext.MinLength = n
	}
}

// Allow patterns to approximately match within this edit distance.
func EditDistance(n uint) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtEditDistance
		ext.EditDistance = n
	}
}

// Allow patterns to approximately match within this Hamming distance.
func HammingDistance(n uint) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtHammingDistance
		ext.HammingDistance = n
	}
}

func (p *Pattern) WithExt(exts ...Ext) *Pattern {
	if exts != nil {
		p.Ext = NewExprExt(exts...)
	}
	return p
}

// Provides additional parameters related to an expression.
func (p *Pattern) Exts() (*ExprExt, error) {
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

Parse pattern from a formated string

	<integer id>:/<expression>/<flags>

For example, the following pattern will match `test` in the caseless and multi-lines mode

	/test/im

*/
func ParsePattern(s string) (*Pattern, error) {
	var p Pattern

	i := strings.Index(s, ":/")
	j := strings.LastIndex(s, "/")
	if i > 0 && j > i+1 {
		id, err := strconv.ParseInt(s[:i], 10, 32)
		if err != nil {
			return nil, errors.New("invalid pattern id: " + s[:i])
		}
		p.Id = int(id)
		s = s[i+1:]
	}

	if n := strings.LastIndex(s, "/"); n > 1 && strings.HasPrefix(s, "/") {
		p.Expression = Expression(s[1:n])
		s = s[n+1:]

		if n = strings.Index(s, "{"); n > 0 && strings.HasSuffix(s, "}") {
			ext, err := ParseExprExt(s[n:])
			if err != nil {
				return nil, fmt.Errorf("invalid expression extensions: %s, %w", s[n:], err)
			}
			p.ext = ext
			s = s[:n]
		}

		flags, err := ParseCompileFlag(s)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern flags: %s, %w", s, err)
		}
		p.Flags = flags
	} else {
		p.Expression = Expression(s)
	}

	info, err := hsExpressionInfo(string(p.Expression), p.Flags)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern: %s, %w", p.Expression, err)
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

// A type containing information on the target platform.
type Platform interface {
	// Information about the target platform which may be used to guide the optimisation process of the compile.
	Tune() TuneFlag

	// Relevant CPU features available on the target platform
	CpuFeatures() CpuFeature
}

func NewPlatform(tune TuneFlag, cpu CpuFeature) Platform { return newPlatformInfo(tune, cpu) }

// Populates the platform information based on the current host.
func PopulatePlatform() Platform {
	platform, _ := hsPopulatePlatform()

	return platform
}

// A type to help to build up a database
type DatabaseBuilder struct {
	// Array of patterns to compile.
	Patterns []*Pattern

	// Compiler mode flags that affect the database as a whole. (Default: block mode)
	Mode ModeFlag

	// If not nil, the platform structure is used to determine the target platform for the database.
	// If nil, a database suitable for running on the current host platform is produced.
	Platform Platform
}

func (b *DatabaseBuilder) AddExpressions(exprs ...Expression) *DatabaseBuilder {
	for _, expr := range exprs {
		b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Id: len(b.Patterns) + 1})
	}

	return b
}

func (b *DatabaseBuilder) AddExpressionWithFlags(expr Expression, flags CompileFlag) *DatabaseBuilder {
	b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Flags: flags, Id: len(b.Patterns) + 1})

	return b
}

func (b *DatabaseBuilder) Build() (Database, error) {
	if b.Patterns == nil {
		return nil, errors.New("no patterns")
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
		return newStreamDatabase(db)
	case VectoredMode:
		return newVectoredDatabase(db)
	case BlockMode:
		return newBlockDatabase(db)
	}

	return nil, errors.New("unknown mode")
}

func NewBlockDatabase(patterns ...*Pattern) (BlockDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: BlockMode}

	db, err := builder.Build()

	if err != nil {
		return nil, err
	}

	return db.(*blockDatabase), err
}

func NewStreamDatabase(patterns ...*Pattern) (StreamDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: StreamMode}

	db, err := builder.Build()

	if err != nil {
		return nil, err
	}

	return db.(*streamDatabase), err
}

func NewMediumStreamDatabase(patterns ...*Pattern) (StreamDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: StreamMode | SomHorizonMediumMode}

	db, err := builder.Build()

	if err != nil {
		return nil, err
	}

	return db.(*streamDatabase), err
}

func NewLargeStreamDatabase(patterns ...*Pattern) (StreamDatabase, error) {
	builder := &DatabaseBuilder{Patterns: patterns, Mode: StreamMode | SomHorizonLargeMode}

	db, err := builder.Build()

	if err != nil {
		return nil, err
	}

	return db.(*streamDatabase), err
}

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

	return newBlockDatabase(db)
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular expressions.
func MustCompile(expr string) Database {
	db, err := hsCompile(expr, SomLeftMost, BlockMode, nil)

	if err != nil {
		panic(`Compile(` + Quote(expr) + `): ` + err.Error())
	}

	bdb, err := newBlockDatabase(db)

	if err != nil {
		panic(`Compile(` + Quote(expr) + `): ` + err.Error())
	}

	return bdb
}

func Quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}
