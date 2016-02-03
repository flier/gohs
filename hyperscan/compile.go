package hyperscan

type Platform interface {
	Tune() TuneFlag

	CpuFeatures() CpuFeature
}

type Expression string

type Pattern struct {
	Expression Expression  // The NULL-terminated expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	Id         uint
}

func NewPlatform(tune TuneFlag, cpu CpuFeature) Platform { return newPlatformInfo(tune, cpu) }

func CurrentPlatform() Platform {
	platform, _ := hsPopulatePlatform()

	return platform
}

type DatabaseBuilder struct {
	Patterns []Pattern

	Mode ModeFlag

	Platform Platform
}

func (b *DatabaseBuilder) Build() (Database, error) {
	expressions := make([]string, len(b.Patterns))
	flags := make([]CompileFlag, len(b.Patterns))
	ids := make([]uint, len(b.Patterns))

	for i, pattern := range b.Patterns {
		expressions[i] = string(pattern.Expression)
		flags[i] = pattern.Flags
		ids[i] = pattern.Id
	}

	platform, _ := b.Platform.(*hsPlatformInfo)

	db, err := hsCompileMulti(expressions, flags, ids, b.Mode, platform)

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
