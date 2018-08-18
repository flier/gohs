package hyperscan

import (
	"errors"
	"strconv"
	"strings"
)

// The expression of pattern
type Expression string

func (e Expression) String() string { return string(e) }

type Pattern struct {
	Expression             // The expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	Id         int         // The ID number to be associated with the corresponding pattern
	info       *ExprInfo
}

func NewPattern(expr string, flags CompileFlag) *Pattern {
	return &Pattern{Expression: Expression(expr), Flags: flags}
}

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

func (p *Pattern) String() string {
	return "/" + string(p.Expression) + "/" + p.Flags.String()
}

/*

Parse pattern from a formated string

	/<expression>/[flags]

For example, the following pattern will match `test` in the caseless and multi-lines mode

	/test/im

*/
func ParsePattern(s string) (*Pattern, error) {
	var p Pattern

	if n := strings.LastIndex(s, "/"); n < 1 || !strings.HasPrefix(s, "/") {
		p.Expression = Expression(s)
	} else {
		p.Expression = Expression(s[1:n])

		flags, err := ParseCompileFlag(strings.ToLower(s[n+1:]))

		if err != nil {
			return nil, errors.New("invalid pattern, " + err.Error())
		}

		p.Flags = flags
	}

	info, err := hsExpressionInfo(string(p.Expression), p.Flags)

	if err != nil {
		return nil, errors.New("invalid pattern, " + err.Error())
	}

	p.info = info

	return &p, nil
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

	needSomLeftMost := false

	for _, pattern := range b.Patterns {
		if (pattern.Flags & SomLeftMost) == SomLeftMost {
			needSomLeftMost = true
		}
	}

	mode := b.Mode

	if mode == 0 {
		mode = BlockMode
	}

	if mode == StreamMode && needSomLeftMost {
		mode |= SomHorizonSmallMode
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
