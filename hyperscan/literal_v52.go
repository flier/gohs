//go:build hyperscan_v52 || hyperscan_v54
// +build hyperscan_v52 hyperscan_v54

package hyperscan

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/flier/gohs/internal/hs"
)

type Literals []*Literal

// Pure literal is a special case of regular expression.
// A character sequence is regarded as a pure literal if and
// only if each character is read and interpreted independently.
// No syntax association happens between any adjacent characters.
// nolint: golint,revive,stylecheck
type Literal struct {
	Expression string      // The expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	Id         int         // The ID number to be associated with the corresponding pattern
	*ExprInfo
}

// NewLiteral returns a new Literal base on expression and compile flags.
func NewLiteral(expr string, flags ...CompileFlag) *Literal {
	var v CompileFlag
	for _, f := range flags {
		v |= f
	}
	return &Literal{Expression: expr, Flags: v}
}

// IsValid validate the literal contains a pure literal.
func (lit *Literal) IsValid() bool {
	_, err := lit.Info()

	return err == nil
}

// Provides information about a regular expression.
func (lit *Literal) Info() (*ExprInfo, error) {
	if lit.ExprInfo == nil {
		info, err := hs.ExpressionInfo(lit.Expression, lit.Flags)
		if err != nil {
			return nil, err //nolint: wrapcheck
		}

		lit.ExprInfo = info
	}

	return lit.ExprInfo, nil
}

func (lit *Literal) String() string {
	var b strings.Builder

	if lit.Id > 0 {
		fmt.Fprintf(&b, "%d:", lit.Id)
	}

	fmt.Fprintf(&b, "/%s/%s", lit.Expression, lit.Flags)

	return b.String()
}

/*
Parse literal from a formated string

	<integer id>:/<expression>/<flags>

For example, the following literal will match `test` in the caseless and multi-lines mode

	/test/im
*/
func ParseLiteral(s string) (*Literal, error) {
	var lit Literal

	i := strings.Index(s, ":/")
	j := strings.LastIndex(s, "/")
	if i > 0 && j > i+1 {
		id, err := strconv.Atoi(s[:i])
		if err != nil {
			return nil, fmt.Errorf("invalid pattern id `%s`, %w", s[:i], ErrInvalid)
		}
		lit.Id = id
		s = s[i+1:]
	}

	if n := strings.LastIndex(s, "/"); n > 1 && strings.HasPrefix(s, "/") {
		lit.Expression = s[1:n]
		s = s[n+1:]

		flags, err := ParseCompileFlag(s)
		if err != nil {
			return nil, fmt.Errorf("invalid pattern flags `%s`, %w", s, err)
		}
		lit.Flags = flags
	} else {
		lit.Expression = s
	}

	info, err := hs.ExpressionInfo(lit.Expression, lit.Flags)
	if err != nil {
		return nil, fmt.Errorf("invalid pattern `%s`, %w", lit.Expression, err)
	}
	lit.ExprInfo = info

	return &lit, nil
}

func (lit *Literal) Literal() *hs.Literal {
	return &hs.Literal{
		Expr:     lit.Expression,
		Flags:    lit.Flags,
		ID:       lit.Id,
		ExprInfo: lit.ExprInfo,
	}
}

func (lit *Literal) Literals() (r []*hs.Literal) {
	return []*hs.Literal{lit.Literal()}
}

func (lit *Literal) Build(mode ModeFlag) (Database, error) {
	return lit.ForPlatform(mode, nil)
}

func (lit *Literal) ForPlatform(mode ModeFlag, platform Platform) (Database, error) {
	if mode == 0 {
		mode = BlockMode
	} else if mode == StreamMode {
		som := (lit.Flags & SomLeftMost) == SomLeftMost

		if som && mode&(SomHorizonSmallMode|SomHorizonMediumMode|SomHorizonLargeMode) == 0 {
			mode |= SomHorizonSmallMode
		}
	}

	p, _ := platform.(*hs.PlatformInfo)

	db, err := hs.CompileLit(lit.Expression, lit.Flags, mode, p)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	switch mode & hs.ModeMask {
	case StreamMode:
		return newStreamDatabase(db), nil
	case VectoredMode:
		return newVectoredDatabase(db), nil
	case BlockMode:
		return newBlockDatabase(db), nil
	}

	return nil, fmt.Errorf("mode %d, %w", mode, ErrInvalid)
}

func (literals Literals) Literals() (r []*hs.Literal) {
	r = make([]*hs.Literal, len(literals))

	for i, lit := range literals {
		r[i] = lit.Literal()
	}

	return
}

func (literals Literals) Build(mode ModeFlag) (Database, error) {
	return literals.ForPlatform(mode, nil)
}

func (literals Literals) ForPlatform(mode ModeFlag, platform Platform) (Database, error) {
	if mode == 0 {
		mode = BlockMode
	} else if mode == StreamMode {
		som := false

		for _, lit := range literals {
			if (lit.Flags & SomLeftMost) == SomLeftMost {
				som = true
			}
		}

		if som && mode&(SomHorizonSmallMode|SomHorizonMediumMode|SomHorizonLargeMode) == 0 {
			mode |= SomHorizonSmallMode
		}
	}

	p, _ := platform.(*hs.PlatformInfo)

	db, err := hs.CompileLitMulti(literals, mode, p)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	switch mode & hs.ModeMask {
	case StreamMode:
		return newStreamDatabase(db), nil
	case VectoredMode:
		return newVectoredDatabase(db), nil
	case BlockMode:
		return newBlockDatabase(db), nil
	}

	return nil, fmt.Errorf("mode %d, %w", mode, ErrInvalid)
}
