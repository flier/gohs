package chimera

import (
	"fmt"
	"strconv"

	"github.com/flier/gohs/internal/ch"
)

// A type containing error details that is returned by the compile calls on failure.
//
// The caller may inspect the values returned in this type to determine the cause of failure.
type CompileError = ch.CompileError

// CompileFlag represents a pattern flag.
type CompileFlag = ch.CompileFlag

const (
	// Caseless represents set case-insensitive matching.
	Caseless CompileFlag = ch.Caseless
	// DotAll represents matching a `.` will not exclude newlines.
	DotAll CompileFlag = ch.DotAll
	// MultiLine set multi-line anchoring.
	MultiLine CompileFlag = ch.MultiLine
	// SingleMatch set single-match only mode.
	SingleMatch CompileFlag = ch.SingleMatch
	// Utf8Mode enable UTF-8 mode for this expression.
	Utf8Mode CompileFlag = ch.Utf8Mode
	// UnicodeProperty enable Unicode property support for this expression.
	UnicodeProperty CompileFlag = ch.UnicodeProperty
)

/*
ParseCompileFlag parse the compile pattern flags from string

	i	Caseless 		Case-insensitive matching
	s	DotAll			Dot (.) will match newlines
	m	MultiLine		Multi-line anchoring
	H	SingleMatch		Report match ID at most once (`o` deprecated)
	8	Utf8Mode		UTF-8 mode (`u` deprecated)
	W	UnicodeProperty		Unicode property support (`p` deprecated)
*/
func ParseCompileFlag(s string) (CompileFlag, error) {
	var flags CompileFlag

	for _, c := range s {
		if flag, exists := ch.CompileFlags[c]; exists {
			flags |= flag
		} else {
			return 0, fmt.Errorf("flag `%c`, %w", c, ErrInvalid)
		}
	}

	return flags, nil
}

// CompileMode flags.
type CompileMode = ch.CompileMode

const (
	// Disable capturing groups.
	NoGroups CompileMode = ch.NoGroups

	// Enable capturing groups.
	Groups CompileMode = ch.Groups
)

// Compile a regular expression and returns, if successful,
// a pattern database in the block mode that can be used to match against text.
func Compile(expr string) (Database, error) {
	db, err := ch.Compile(expr, 0, ch.Groups, nil)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	return &database{db}, nil
}

// MustCompile is like Compile but panics if the expression cannot be parsed.
// It simplifies safe initialization of global variables holding compiled regular expressions.
func MustCompile(expr string) Database {
	db, err := Compile(expr)
	if err != nil {
		panic(`Compile(` + Quote(expr) + `): ` + err.Error())
	}

	return db
}

// Quote returns a quoted string literal representing s.
func Quote(s string) string {
	if strconv.CanBackquote(s) {
		return "`" + s + "`"
	}

	return strconv.Quote(s)
}
