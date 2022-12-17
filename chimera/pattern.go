//go:build chimera
// +build chimera

package chimera

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/flier/gohs/internal/ch"
)

// Pattern is a matching pattern.
type Pattern ch.Pattern

// NewPattern returns a new pattern base on expression and compile flags.
func NewPattern(expr string, flags CompileFlag) *Pattern {
	return &Pattern{Expression: expr, Flags: flags}
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
			return nil, fmt.Errorf("pattern id `%s`, %w", s[:i], ErrInvalid)
		}

		p.ID = id
		s = s[i+1:]
	}

	if n := strings.LastIndex(s, "/"); n > 1 && strings.HasPrefix(s, "/") {
		p.Expression = s[1:n]
		s = s[n+1:]

		flags, err := ParseCompileFlag(s)
		if err != nil {
			return nil, fmt.Errorf("pattern flags `%s`, %w", s, err)
		}

		p.Flags = flags
	} else {
		p.Expression = s
	}

	return &p, nil
}

func (p *Pattern) String() string {
	var b strings.Builder

	if p.ID > 0 {
		fmt.Fprintf(&b, "%d:", p.ID)
	}

	fmt.Fprintf(&b, "/%s/%s", p.Expression, p.Flags)

	return b.String()
}

// Patterns is a set of matching patterns.
type Patterns []*Pattern

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

func (p Patterns) Patterns() (r []*ch.Pattern) {
	r = make([]*ch.Pattern, len(p))
	for i, pat := range p {
		r[i] = (*ch.Pattern)(pat)
	}
	return
}
