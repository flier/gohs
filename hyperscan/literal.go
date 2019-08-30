// +build !hyperscan_v4,!hyperscan_v5_1

package hyperscan

import (
	"errors"
	"fmt"
	"strings"
)

type Literal struct {
	Expression []byte      // The expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	Id         int         // The ID number to be associated with the corresponding pattern
}

func NewLiteral(expr []byte, flags CompileFlag) *Literal {
	return &Literal{Expression: expr, Flags: flags}
}

/*

Parse literal from a formated string

	/<expression>/[flags]

For example, the following pattern will match `test` in the caseless and multi-lines mode

	/test/im

*/
func ParseLiteral(s string) (*Literal, error) {
	var lit Literal

	if n := strings.LastIndex(s, "/"); n < 1 || !strings.HasPrefix(s, "/") {
		lit.Expression = []byte(s)
	} else {
		lit.Expression = []byte(s[1:n])

		flags, err := ParseCompileFlag(s[n+1:])

		if err != nil {
			return nil, fmt.Errorf("invalid pattern, %v", err)
		}

		if f := flags &^ (Caseless | MultiLine | SingleMatch | SomLeftMost); f != 0 {
			return nil, fmt.Errorf("invalid flags, %s of %s", f, flags)
		}

		lit.Flags = flags
	}

	return &lit, nil
}

func (lit *Literal) String() string {
	expr := fmt.Sprintf("%q", lit.Expression)
	return fmt.Sprintf("/%s/%s", expr[1:len(expr)-1], lit.Flags)
}

func (lit *Literal) Compile(mode ModeFlag, platform Platform) (Database, error) {
	p, _ := platform.(*hsPlatformInfo)
	db, err := hsCompileLit(lit.Expression, lit.Flags, mode, p)

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

// A type to help to build up a literal database
type LiteralDatabaseBuilder struct {
	// Array of literals to compile.
	Literals []*Literal

	// Compiler mode flags that affect the database as a whole. (Default: block mode)
	Mode ModeFlag

	// If not nil, the platform structure is used to determine the target platform for the database.
	// If nil, a database suitable for running on the current host platform is produced.
	Platform Platform
}

func (b *LiteralDatabaseBuilder) Build() (Database, error) {
	mode := b.Mode
	if mode == 0 {
		mode = BlockMode
	}
	platform, _ := b.Platform.(*hsPlatformInfo)
	db, err := hsCompileLitMulti(b.Literals, mode, platform)

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
