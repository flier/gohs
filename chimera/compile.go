//go:build chimera
// +build chimera

package chimera

import (
	"fmt"
	"runtime"
	"strconv"

	"github.com/flier/gohs/hyperscan"
	"github.com/flier/gohs/internal/ch"
	"github.com/flier/gohs/internal/hs"
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

// Builder creates a database with the given mode and target platform.
type Builder interface {
	// Build the database with the given mode.
	Build(mode CompileMode) (Database, error)

	// ForPlatform determine the target platform for the database
	ForPlatform(mode CompileMode, platform hyperscan.Platform) (Database, error)
}

// Build the database with the given mode.
func (p *Pattern) Build(mode CompileMode) (Database, error) {
	return p.ForPlatform(mode, nil)
}

// ForPlatform determine the target platform for the database.
func (p *Pattern) ForPlatform(mode CompileMode, platform hyperscan.Platform) (Database, error) {
	b := DatabaseBuilder{Patterns: Patterns{p}, Mode: mode, Platform: platform}
	return b.Build()
}

// Build the database with the given mode.
func (p Patterns) Build(mode CompileMode) (Database, error) {
	return p.ForPlatform(mode, nil)
}

// ForPlatform determine the target platform for the database.
func (p Patterns) ForPlatform(mode CompileMode, platform hyperscan.Platform) (Database, error) {
	b := DatabaseBuilder{Patterns: p, Mode: mode, Platform: platform}
	return b.Build()
}

// DatabaseBuilder creates a database that will be used to matching the patterns.
type DatabaseBuilder struct {
	// Array of patterns to compile.
	Patterns

	// Compiler mode flags that affect the database as a whole. (Default: capturing groups mode)
	Mode CompileMode

	// If not nil, the platform structure is used to determine the target platform for the database.
	// If nil, a database suitable for running on the current host platform is produced.
	hyperscan.Platform

	// A limit from pcre_extra on the amount of match function called in PCRE to limit backtracking that can take place.
	MatchLimit uint

	// A limit from pcre_extra on the recursion depth of match function in PCRE.
	MatchLimitRecursion uint
}

// AddExpressions add more expressions to the database.
func (b *DatabaseBuilder) AddExpressions(exprs ...string) *DatabaseBuilder {
	for _, expr := range exprs {
		b.Patterns = append(b.Patterns, &Pattern{Expression: expr, ID: len(b.Patterns) + 1})
	}

	return b
}

// AddExpressionWithFlags add more expressions with flags to the database.
func (b *DatabaseBuilder) AddExpressionWithFlags(expr string, flags CompileFlag) *DatabaseBuilder {
	b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Flags: flags, ID: len(b.Patterns) + 1})

	return b
}

// Build a database base on the expressions and platform.
func (b *DatabaseBuilder) Build() (Database, error) {
	if b.Patterns == nil {
		return nil, ErrInvalid
	}

	platform, _ := b.Platform.(*hs.PlatformInfo)

	db, err := ch.CompileExtMulti(b.Patterns, b.Mode, platform, b.MatchLimit, b.MatchLimitRecursion)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	return newBlockDatabase(db), nil
}

// NewBlockDatabase compile expressions into a pattern database.
func NewBlockDatabase(patterns ...*Pattern) (BlockDatabase, error) {
	db, err := Patterns(patterns).Build(Groups)
	if err != nil {
		return nil, err
	}

	return db.(BlockDatabase), err
}

// NewManagedBlockDatabase is a wrapper for NewBlockDatabase that
// sets a finalizer on the Scratch instance so that memory is
// freed once the object is no longer in use.
func NewManagedBlockDatabase(patterns ...*Pattern) (BlockDatabase, error) {
	db, err := NewBlockDatabase(patterns...)
	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(db, func(obj BlockDatabase) {
		_ = obj.Close()
	})

	return db, nil
}

// Compile a regular expression and returns, if successful,
// a pattern database in the block mode that can be used to match against text.
func Compile(expr string) (Database, error) {
	db, err := ch.Compile(expr, 0, ch.Groups, nil)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	return newBlockDatabase(db), nil
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
