package ch

import (
	"fmt"
	"runtime"
	"sort"
	"strings"
	"unsafe"

	"github.com/flier/gohs/internal/hs"
)

// #include <ch.h>
import "C"

// Pattern is a matching pattern.
type Pattern struct {
	Expression string      // The expression to parse.
	Flags      CompileFlag // Flags which modify the behaviour of the expression.
	ID         int         // The ID number to be associated with the corresponding pattern
}

// NewPattern returns a new pattern base on expression and compile flags.
func NewPattern(expr string, flags CompileFlag) *Pattern {
	return &Pattern{
		Expression: expr,
		Flags:      flags,
	}
}

// A type containing error details that is returned by the compile calls on failure.
//
// The caller may inspect the values returned in this type to determine the cause of failure.
type CompileError struct {
	Message    string // A human-readable error message describing the error.
	Expression int    // The zero-based number of the expression that caused the error.
}

func (e *CompileError) Error() string { return e.Message }

type CompileFlag uint

const (
	// Caseless represents set case-insensitive matching.
	Caseless CompileFlag = C.CH_FLAG_CASELESS
	// DotAll represents matching a `.` will not exclude newlines.
	DotAll CompileFlag = C.CH_FLAG_DOTALL
	// MultiLine set multi-line anchoring.
	MultiLine CompileFlag = C.CH_FLAG_MULTILINE
	// SingleMatch set single-match only mode.
	SingleMatch CompileFlag = C.CH_FLAG_SINGLEMATCH
	// Utf8Mode enable UTF-8 mode for this expression.
	Utf8Mode CompileFlag = C.CH_FLAG_UTF8
	// UnicodeProperty enable Unicode property support for this expression.
	UnicodeProperty CompileFlag = C.HS_FLAG_UCP
)

var CompileFlags = map[rune]CompileFlag{
	'i': Caseless,
	's': DotAll,
	'm': MultiLine,
	'H': SingleMatch,
	'8': Utf8Mode,
	'W': UnicodeProperty,
}

func (flags CompileFlag) String() string {
	var values []string

	for c, flag := range CompileFlags {
		if (flags & flag) == flag {
			values = append(values, string(c))
		}
	}

	sort.Strings(values)

	return strings.Join(values, "")
}

// CompileMode flags.
type CompileMode int

const (
	// Disable capturing groups.
	NoGroups CompileMode = C.CH_MODE_NOGROUPS

	// Enable capturing groups.
	Groups CompileMode = C.CH_MODE_GROUPS
)

// The basic regular expression compiler.
func Compile(expression string, flags CompileFlag, mode CompileMode, info *hs.PlatformInfo) (Database, error) {
	var db *C.ch_database_t
	var err *C.ch_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = (*C.struct_hs_platform_info)(unsafe.Pointer(&info.Platform))
	}

	expr := C.CString(expression)

	defer C.free(unsafe.Pointer(expr))

	ret := C.ch_compile(expr, C.uint(flags), C.uint(mode), platform, &db, &err)

	if err != nil {
		defer C.ch_free_compile_error(err)
	}

	if ret == C.CH_SUCCESS {
		return db, nil
	}

	if ret == C.CH_COMPILER_ERROR && err != nil {
		return nil, &CompileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}

type Patterns interface {
	Patterns() []*Pattern
}

// The multiple regular expression compiler.
func CompileMulti(p Patterns, mode CompileMode, info *hs.PlatformInfo) (Database, error) {
	return CompileExtMulti(p, mode, info, 0, 0)
}

// The multiple regular expression compiler.
func CompileExtMulti(p Patterns, mode CompileMode, info *hs.PlatformInfo,
	matchLimit, matchLimitRecursion uint) (Database, error) {
	var db *C.ch_database_t
	var err *C.ch_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = (*C.struct_hs_platform_info)(unsafe.Pointer(&info.Platform))
	}

	patterns := p.Patterns()

	cexprs := (**C.char)(C.calloc(C.size_t(len(patterns)), C.size_t(unsafe.Sizeof(uintptr(0)))))
	exprs := (*[1 << 30]*C.char)(unsafe.Pointer(cexprs))[:len(patterns):len(patterns)]

	cflags := (*C.uint)(C.calloc(C.size_t(len(patterns)), C.size_t(unsafe.Sizeof(C.uint(0)))))
	flags := (*[1 << 30]C.uint)(unsafe.Pointer(cflags))[:len(patterns):len(patterns)]

	cids := (*C.uint)(C.calloc(C.size_t(len(patterns)), C.size_t(unsafe.Sizeof(C.uint(0)))))
	ids := (*[1 << 30]C.uint)(unsafe.Pointer(cids))[:len(patterns):len(patterns)]

	for i, pattern := range patterns {
		exprs[i] = C.CString(pattern.Expression)
		flags[i] = C.uint(pattern.Flags)
		ids[i] = C.uint(pattern.ID)
	}

	ret := C.ch_compile_ext_multi(cexprs, cflags, cids, C.uint(len(patterns)), C.uint(mode),
		C.ulong(matchLimit), C.ulong(matchLimitRecursion), platform, &db, &err)

	for _, expr := range exprs {
		C.free(unsafe.Pointer(expr))
	}

	C.free(unsafe.Pointer(cexprs))
	C.free(unsafe.Pointer(cflags))
	C.free(unsafe.Pointer(cids))

	runtime.KeepAlive(patterns)

	if err != nil {
		defer C.ch_free_compile_error(err)
	}

	if ret == C.CH_SUCCESS {
		return db, nil
	}

	if ret == C.CH_COMPILER_ERROR && err != nil {
		return nil, &CompileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}
