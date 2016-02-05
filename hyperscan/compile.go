package hyperscan

import (
	"errors"
	"strconv"
	"strings"
)

type Expression string

type Pattern struct {
	Expression             // The NULL-terminated expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	Id         int
	info       *ExprInfo
}

func (p *Pattern) IsValid() bool {
	_, err := p.Info()

	return err == nil
}

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

func ParsePattern(s string) (*Pattern, error) {
	var p Pattern

	if n := strings.LastIndex(s, "/"); n < 1 || !strings.HasPrefix(s, "/") {
		p.Expression = Expression(s)
	} else {
		p.Expression = Expression(s[1:n])

		if flags, err := ParseCompileFlag(strings.ToLower(s[n+1:])); err != nil {
			return nil, errors.New("invalid pattern, " + err.Error())
		} else {
			p.Flags = flags
		}
	}

	info, err := hsExpressionInfo(string(p.Expression), p.Flags)

	if err != nil {
		return nil, errors.New("invalid pattern, " + err.Error())
	}

	p.info = info

	return &p, nil
}

type Platform interface {
	Tune() TuneFlag

	CpuFeatures() CpuFeature
}

func NewPlatform(tune TuneFlag, cpu CpuFeature) Platform { return newPlatformInfo(tune, cpu) }

func CurrentPlatform() Platform {
	platform, _ := hsPopulatePlatform()

	return platform
}

type DatabaseBuilder struct {
	Patterns []*Pattern

	Mode ModeFlag

	Platform Platform
}

func (b *DatabaseBuilder) AddPatterns(exprs ...Expression) *DatabaseBuilder {
	for _, expr := range exprs {
		b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Id: len(b.Patterns) + 1})
	}

	return b
}

func (b *DatabaseBuilder) AddPatternWithFlags(expr Expression, flags CompileFlag) *DatabaseBuilder {
	b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Flags: flags, Id: len(b.Patterns) + 1})

	return b
}

func (b *DatabaseBuilder) Build() (Database, error) {
	if b.Patterns == nil {
		return nil, errors.New("no patterns")
	}

	expressions := make([]string, len(b.Patterns))
	flags := make([]CompileFlag, len(b.Patterns))
	ids := make([]uint, len(b.Patterns))

	for i, pattern := range b.Patterns {
		expressions[i] = string(pattern.Expression)
		flags[i] = pattern.Flags
		ids[i] = uint(pattern.Id)
	}

	mode := b.Mode

	if mode == 0 {
		mode = Block
	}

	platform, _ := b.Platform.(*hsPlatformInfo)

	db, err := hsCompileMulti(expressions, flags, ids, mode, platform)

	if err != nil {
		return nil, err
	}

	return &database{db}, nil
}

func Compile(expr string) (Database, error) {
	db, err := hsCompile(expr, 0, Block, nil)

	if err != nil {
		return nil, err
	}

	return &database{db}, nil
}

func MustCompile(expr string) Database {
	db, err := hsCompile(expr, 0, Block, nil)

	if err != nil {
		panic(`Compile(` + quote(expr) + `): ` + err.Error())
	}

	return &database{db}
}

func quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}
	return strconv.Quote(s)
}
