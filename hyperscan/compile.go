package hyperscan

type Expression string

type Pattern interface {
	Expression() Expression

	Flags() CompileFlag

	Id() uint
}

type Database interface {
	Len() int

	Info() string

	Close() error
}

type DatabaseBuilder struct {
	Patterns []Pattern

	Mode ModeFlag

	Platform Platform
}

func (b *DatabaseBuilder) Build() Database {
	return nil
}

type database struct {
	db hsDatabase
}
