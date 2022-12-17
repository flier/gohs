package hyperscan

import (
	"fmt"
	"runtime"
	"strconv"
	"strings"

	"github.com/flier/gohs/internal/hs"
)

// A type containing error details that is returned by the compile calls on failure.
//
// The caller may inspect the values returned in this type to determine the cause of failure.
type CompileError = hs.CompileError

// CompileFlag represents a pattern flag.
type CompileFlag = hs.CompileFlag

const (
	// Caseless represents set case-insensitive matching.
	Caseless CompileFlag = hs.Caseless
	// DotAll represents matching a `.` will not exclude newlines.
	DotAll CompileFlag = hs.DotAll
	// MultiLine set multi-line anchoring.
	MultiLine CompileFlag = hs.MultiLine
	// SingleMatch set single-match only mode.
	SingleMatch CompileFlag = hs.SingleMatch
	// AllowEmpty allow expressions that can match against empty buffers.
	AllowEmpty CompileFlag = hs.AllowEmpty
	// Utf8Mode enable UTF-8 mode for this expression.
	Utf8Mode CompileFlag = hs.Utf8Mode
	// UnicodeProperty enable Unicode property support for this expression.
	UnicodeProperty CompileFlag = hs.UnicodeProperty
	// PrefilterMode enable prefiltering mode for this expression.
	PrefilterMode CompileFlag = hs.PrefilterMode
	// SomLeftMost enable leftmost start of match reporting.
	SomLeftMost CompileFlag = hs.SomLeftMost
)

/*
ParseCompileFlag parse the compile pattern flags from string

	i	Caseless 		Case-insensitive matching
	s	DotAll			Dot (.) will match newlines
	m	MultiLine		Multi-line anchoring
	H	SingleMatch		Report match ID at most once (`o` deprecated)
	V	AllowEmpty		Allow patterns that can match against empty buffers (`e` deprecated)
	8	Utf8Mode		UTF-8 mode (`u` deprecated)
	W	UnicodeProperty		Unicode property support (`p` deprecated)
	P	PrefilterMode		Prefiltering mode (`f` deprecated)
	L	SomLeftMost		Leftmost start of match reporting (`l` deprecated)
	C	Combination		Logical combination of patterns (Hyperscan 5.0)
	Q	Quiet			Quiet at matching (Hyperscan 5.0)
*/
func ParseCompileFlag(s string) (CompileFlag, error) {
	var flags CompileFlag

	for _, c := range s {
		if flag, exists := hs.CompileFlags[c]; exists {
			flags |= flag
		} else if flag, exists := hs.DeprecatedCompileFlags[c]; exists {
			flags |= flag
		} else {
			return 0, fmt.Errorf("flag `%c`, %w", c, ErrInvalid)
		}
	}

	return flags, nil
}

// ModeFlag represents the compile mode flags.
type ModeFlag = hs.ModeFlag

const (
	// BlockMode for the block scan (non-streaming) database.
	BlockMode ModeFlag = hs.BlockMode
	// NoStreamMode is alias for Block.
	NoStreamMode ModeFlag = hs.NoStreamMode
	// StreamMode for the streaming database.
	StreamMode ModeFlag = hs.StreamMode
	// VectoredMode for the vectored scanning database.
	VectoredMode ModeFlag = hs.VectoredMode
	// SomHorizonLargeMode use full precision to track start of match offsets in stream state.
	SomHorizonLargeMode ModeFlag = hs.SomHorizonLargeMode
	// SomHorizonMediumMode use medium precision to track start of match offsets in stream state (within 2^32 bytes).
	SomHorizonMediumMode ModeFlag = hs.SomHorizonMediumMode
	// SomHorizonSmallMode use limited precision to track start of match offsets in stream state (within 2^16 bytes).
	SomHorizonSmallMode ModeFlag = hs.SomHorizonSmallMode
)

// ParseModeFlag parse a database mode from string.
func ParseModeFlag(s string) (ModeFlag, error) {
	if mode, exists := hs.ModeFlags[strings.ToUpper(s)]; exists {
		return mode, nil
	}

	return BlockMode, fmt.Errorf("database mode %s, %w", s, ErrInvalid)
}

// Builder creates a database with the given mode and target platform.
type Builder interface {
	// Build the database with the given mode.
	Build(mode ModeFlag) (Database, error)

	// ForPlatform determine the target platform for the database
	ForPlatform(mode ModeFlag, platform Platform) (Database, error)
}

// Build the database with the given mode.
func (p *Pattern) Build(mode ModeFlag) (Database, error) {
	return p.ForPlatform(mode, nil)
}

// ForPlatform determine the target platform for the database.
func (p *Pattern) ForPlatform(mode ModeFlag, platform Platform) (Database, error) {
	b := DatabaseBuilder{Patterns: Patterns{p}, Mode: mode, Platform: platform}
	return b.Build()
}

// Build the database with the given mode.
func (p Patterns) Build(mode ModeFlag) (Database, error) {
	return p.ForPlatform(mode, nil)
}

// ForPlatform determine the target platform for the database.
func (p Patterns) ForPlatform(mode ModeFlag, platform Platform) (Database, error) {
	b := DatabaseBuilder{Patterns: p, Mode: mode, Platform: platform}
	return b.Build()
}

// DatabaseBuilder creates a database that will be used to matching the patterns.
type DatabaseBuilder struct {
	// Array of patterns to compile.
	Patterns

	// Compiler mode flags that affect the database as a whole. (Default: block mode)
	Mode ModeFlag

	// If not nil, the platform structure is used to determine the target platform for the database.
	// If nil, a database suitable for running on the current host platform is produced.
	Platform Platform
}

// AddExpressions add more expressions to the database.
func (b *DatabaseBuilder) AddExpressions(exprs ...string) *DatabaseBuilder {
	for _, expr := range exprs {
		b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Id: len(b.Patterns) + 1})
	}

	return b
}

// AddExpressionWithFlags add more expressions with flags to the database.
func (b *DatabaseBuilder) AddExpressionWithFlags(expr string, flags CompileFlag) *DatabaseBuilder {
	b.Patterns = append(b.Patterns, &Pattern{Expression: expr, Flags: flags, Id: len(b.Patterns) + 1})

	return b
}

// Build a database base on the expressions and platform.
func (b *DatabaseBuilder) Build() (Database, error) {
	if b.Patterns == nil {
		return nil, ErrInvalid
	}

	mode := b.Mode

	if mode == 0 {
		mode = BlockMode
	} else if mode == StreamMode {
		som := false

		for _, pattern := range b.Patterns {
			if (pattern.Flags & SomLeftMost) == SomLeftMost {
				som = true
			}
		}

		if som && mode&(SomHorizonSmallMode|SomHorizonMediumMode|SomHorizonLargeMode) == 0 {
			mode |= SomHorizonSmallMode
		}
	}

	platform, _ := b.Platform.(*hs.PlatformInfo)

	db, err := hs.CompileMulti(b.Patterns, mode, platform)
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
	default:
		return nil, fmt.Errorf("mode %d, %w", mode, ErrInvalid)
	}
}

// NewBlockDatabase create a block database base on the patterns.
func NewBlockDatabase(patterns ...*Pattern) (bdb BlockDatabase, err error) {
	var db Database
	db, err = Patterns(patterns).Build(BlockMode)
	if err != nil {
		return
	}

	bdb, _ = db.(*blockDatabase)
	return
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

// NewStreamDatabase create a stream database base on the patterns.
func NewStreamDatabase(patterns ...*Pattern) (sdb StreamDatabase, err error) {
	var db Database
	db, err = Patterns(patterns).Build(StreamMode)
	if err != nil {
		return
	}

	sdb, _ = db.(*streamDatabase)
	return
}

// NewManagedStreamDatabase is a wrapper for NewStreamDatabase that
// sets a finalizer on the Scratch instance so that memory is
// freed once the object is no longer in use.
func NewManagedStreamDatabase(patterns ...*Pattern) (StreamDatabase, error) {
	db, err := NewStreamDatabase(patterns...)
	if err != nil {
		return nil, err
	}

	runtime.SetFinalizer(db, func(obj StreamDatabase) {
		_ = obj.Close()
	})

	return db, nil
}

// NewMediumStreamDatabase create a medium-sized stream database base on the patterns.
func NewMediumStreamDatabase(patterns ...*Pattern) (sdb StreamDatabase, err error) {
	var db Database
	db, err = Patterns(patterns).Build(StreamMode | SomHorizonMediumMode)
	if err != nil {
		return
	}

	sdb, _ = db.(*streamDatabase)
	return
}

// NewLargeStreamDatabase create a large-sized stream database base on the patterns.
func NewLargeStreamDatabase(patterns ...*Pattern) (sdb StreamDatabase, err error) {
	var db Database
	db, err = Patterns(patterns).Build(StreamMode | SomHorizonLargeMode)
	if err != nil {
		return
	}

	sdb, _ = db.(*streamDatabase)
	return
}

// NewVectoredDatabase create a vectored database base on the patterns.
func NewVectoredDatabase(patterns ...*Pattern) (vdb VectoredDatabase, err error) {
	var db Database
	db, err = Patterns(patterns).Build(VectoredMode)
	if err != nil {
		return
	}

	vdb, _ = db.(*vectoredDatabase)
	return
}

// Compile a regular expression and returns, if successful,
// a pattern database in the block mode that can be used to match against text.
func Compile(expr string) (Database, error) {
	db, err := hs.Compile(expr, SomLeftMost, BlockMode, nil)
	if err != nil {
		return nil, err //nolint: wrapcheck
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
