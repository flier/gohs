package hyperscan

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/flier/gohs/internal/hs"
)

// ExprInfo containing information related to an expression.
type ExprInfo = hs.ExprInfo

// ExtFlag are used in ExprExt.Flags to indicate which fields are used.
type ExtFlag = hs.ExtFlag

const (
	// ExtMinOffset is a flag indicating that the ExprExt.MinOffset field is used.
	ExtMinOffset ExtFlag = hs.ExtMinOffset
	// ExtMaxOffset is a flag indicating that the ExprExt.MaxOffset field is used.
	ExtMaxOffset ExtFlag = hs.ExtMaxOffset
	// ExtMinLength is a flag indicating that the ExprExt.MinLength field is used.
	ExtMinLength ExtFlag = hs.ExtMinLength
	// ExtEditDistance is a flag indicating that the ExprExt.EditDistance field is used.
	ExtEditDistance ExtFlag = hs.ExtEditDistance
	// ExtHammingDistance is a flag indicating that the ExprExt.HammingDistance field is used.
	ExtHammingDistance ExtFlag = hs.ExtHammingDistance
)

// Ext is a option containing additional parameters related to an expression.
type Ext func(ext *ExprExt)

// MinOffset given the minimum end offset in the data stream at which this expression should match successfully.
func MinOffset(n uint64) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtMinOffset
		ext.MinOffset = n
	}
}

// MaxOffset given the maximum end offset in the data stream at which this expression should match successfully.
func MaxOffset(n uint64) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtMaxOffset
		ext.MaxOffset = n
	}
}

// MinLength given the minimum match length (from start to end) required to successfully match this expression.
func MinLength(n uint64) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtMinLength
		ext.MinLength = n
	}
}

// EditDistance allow patterns to approximately match within this edit distance.
func EditDistance(n uint32) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtEditDistance
		ext.EditDistance = n
	}
}

// HammingDistance allow patterns to approximately match within this Hamming distance.
func HammingDistance(n uint32) Ext {
	return func(ext *ExprExt) {
		ext.Flags |= ExtHammingDistance
		ext.HammingDistance = n
	}
}

// ExprExt is a structure containing additional parameters related to an expression.
type ExprExt hs.ExprExt

func NewExprExt(exts ...Ext) (ext *ExprExt) {
	if len(exts) == 0 {
		return
	}

	ext = new(ExprExt)

	for _, f := range exts {
		f(ext)
	}

	return ext
}

// With specifies the additional parameters related to an expression.
func (ext *ExprExt) With(exts ...Ext) *ExprExt {
	for _, f := range exts {
		f(ext)
	}

	return ext
}

func (ext *ExprExt) String() string {
	var values []string

	if (ext.Flags & ExtMinOffset) == ExtMinOffset {
		values = append(values, fmt.Sprintf("min_offset=%d", ext.MinOffset))
	}

	if (ext.Flags & ExtMaxOffset) == ExtMaxOffset {
		values = append(values, fmt.Sprintf("max_offset=%d", ext.MaxOffset))
	}

	if (ext.Flags & ExtMinLength) == ExtMinLength {
		values = append(values, fmt.Sprintf("min_length=%d", ext.MinLength))
	}

	if (ext.Flags & ExtEditDistance) == ExtEditDistance {
		values = append(values, fmt.Sprintf("edit_distance=%d", ext.EditDistance))
	}

	if (ext.Flags & ExtHammingDistance) == ExtHammingDistance {
		values = append(values, fmt.Sprintf("hamming_distance=%d", ext.HammingDistance))
	}

	return "{" + strings.Join(values, ",") + "}"
}

const keyValuePair = 2

// ParseExprExt parse containing additional parameters from string.
func ParseExprExt(s string) (ext *ExprExt, err error) {
	ext = new(ExprExt)

	if strings.HasPrefix(s, "{") && strings.HasSuffix(s, "}") {
		s = strings.TrimSuffix(strings.TrimPrefix(s, "{"), "}")
	}

	for _, s := range strings.Split(s, ",") {
		parts := strings.SplitN(s, "=", keyValuePair)

		if len(parts) != keyValuePair {
			continue
		}

		key := strings.ToLower(parts[0])
		value := parts[1]

		var n int

		if n, err = strconv.Atoi(value); err != nil {
			return ext, fmt.Errorf("parse value, %w", err)
		}

		switch key {
		case "min_offset":
			ext.Flags |= ExtMinOffset
			ext.MinOffset = uint64(n)

		case "max_offset":
			ext.Flags |= ExtMaxOffset
			ext.MaxOffset = uint64(n)

		case "min_length":
			ext.Flags |= ExtMinLength
			ext.MinLength = uint64(n)

		case "edit_distance":
			ext.Flags |= ExtEditDistance
			ext.EditDistance = uint32(n)

		case "hamming_distance":
			ext.Flags |= ExtHammingDistance
			ext.HammingDistance = uint32(n)
		}
	}

	return //nolint: nakedret
}

// Pattern is a matching pattern.
type Pattern struct {
	Expression string      // The expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	// The ID number to be associated with the corresponding pattern
	Id   int //nolint: revive,stylecheck
	info *ExprInfo
	ext  *ExprExt
}

// NewPattern returns a new pattern base on expression and compile flags.
func NewPattern(expr string, flags CompileFlag, exts ...Ext) *Pattern {
	return &Pattern{
		Expression: expr,
		Flags:      flags,
		ext:        NewExprExt(exts...),
	}
}

func (p *Pattern) Pattern() *hs.Pattern {
	return &hs.Pattern{
		Expr:  p.Expression,
		Flags: p.Flags,
		ID:    p.Id,
		Ext:   (*hs.ExprExt)(p.ext),
	}
}

func (p *Pattern) Patterns() []*hs.Pattern {
	return []*hs.Pattern{p.Pattern()}
}

// IsValid validate the pattern contains a regular expression.
func (p *Pattern) IsValid() bool {
	_, err := p.Info()

	return err == nil
}

// Info provides information about a regular expression.
func (p *Pattern) Info() (*ExprInfo, error) {
	if p.info == nil {
		info, err := hs.ExpressionInfo(p.Expression, p.Flags)
		if err != nil {
			return nil, err //nolint: wrapcheck
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
		ext, info, err := hs.ExpressionExt(p.Expression, p.Flags)
		if err != nil {
			return nil, err //nolint: wrapcheck
		}

		p.ext = (*ExprExt)(ext)
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
			return nil, fmt.Errorf("pattern id `%s`, %w", s[:i], ErrInvalid)
		}

		p.Id = id
		s = s[i+1:]
	}

	if n := strings.LastIndex(s, "/"); n > 1 && strings.HasPrefix(s, "/") {
		p.Expression = s[1:n]
		s = s[n+1:]

		if n = strings.Index(s, "{"); n > 0 && strings.HasSuffix(s, "}") {
			ext, err := ParseExprExt(s[n:])
			if err != nil {
				return nil, fmt.Errorf("expression extensions `%s`, %w", s[n:], err)
			}

			p.ext = ext
			s = s[:n]
		}

		flags, err := ParseCompileFlag(s)
		if err != nil {
			return nil, fmt.Errorf("pattern flags `%s`, %w", s, err)
		}

		p.Flags = flags
	} else {
		p.Expression = s
	}

	info, err := hs.ExpressionInfo(p.Expression, p.Flags)
	if err != nil {
		return nil, fmt.Errorf("pattern `%s`, %w", p.Expression, err)
	}

	p.info = info

	return &p, nil
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

func (p Patterns) Patterns() (r []*hs.Pattern) {
	r = make([]*hs.Pattern, len(p))

	for i, pat := range p {
		r[i] = pat.Pattern()
	}

	return
}
