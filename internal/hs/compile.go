package hs

/*
#include <limits.h>

#include <hs.h>
*/
import "C"

import (
	"errors"
	"fmt"
	"runtime"
	"sort"
	"strings"
	"unsafe"
)

var (
	// ErrNoFound means patterns not found.
	ErrNoFound = errors.New("no found")
	// ErrUnexpected means item is unexpected.
	ErrUnexpected = errors.New("unexpected")
)

// A type containing error details that is returned by the compile calls on failure.
//
// The caller may inspect the values returned in this type to determine the cause of failure.
type CompileError struct {
	// A human-readable error message describing the error.
	Message string
	// The zero-based number of the expression that caused the error.
	Expression int
}

func (e *CompileError) Error() string { return e.Message }

// CpuFeature is the CPU feature support flags
type CpuFeature int //nolint: golint,stylecheck,revive

const (
	// AVX2 is a CPU features flag indicates that the target platform supports AVX2 instructions.
	AVX2 CpuFeature = C.HS_CPU_FEATURES_AVX2
	// AVX512 is a CPU features flag indicates that the target platform supports AVX512 instructions,
	// specifically AVX-512BW. Using AVX512 implies the use of AVX2.
	AVX512 CpuFeature = C.HS_CPU_FEATURES_AVX512
)

// TuneFlag is the tuning flags
type TuneFlag int

const (
	// Generic indicates that the compiled database should not be tuned for any particular target platform.
	Generic TuneFlag = C.HS_TUNE_FAMILY_GENERIC
	// SandyBridge indicates that the compiled database should be tuned for the Sandy Bridge microarchitecture.
	SandyBridge TuneFlag = C.HS_TUNE_FAMILY_SNB
	// IvyBridge indicates that the compiled database should be tuned for the Ivy Bridge microarchitecture.
	IvyBridge TuneFlag = C.HS_TUNE_FAMILY_IVB
	// Haswell indicates that the compiled database should be tuned for the Haswell microarchitecture.
	Haswell TuneFlag = C.HS_TUNE_FAMILY_HSW
	// Silvermont indicates that the compiled database should be tuned for the Silvermont microarchitecture.
	Silvermont TuneFlag = C.HS_TUNE_FAMILY_SLM
	// Broadwell indicates that the compiled database should be tuned for the Broadwell microarchitecture.
	Broadwell TuneFlag = C.HS_TUNE_FAMILY_BDW
	// Skylake indicates that the compiled database should be tuned for the Skylake microarchitecture.
	Skylake TuneFlag = C.HS_TUNE_FAMILY_SKL
	// SkylakeServer indicates that the compiled database should be tuned for the Skylake Server microarchitecture.
	SkylakeServer TuneFlag = C.HS_TUNE_FAMILY_SKX
	// Goldmont indicates that the compiled database should be tuned for the Goldmont microarchitecture.
	Goldmont TuneFlag = C.HS_TUNE_FAMILY_GLM
)

type PlatformInfo struct {
	Platform C.struct_hs_platform_info
}

// Tune returns the tuning flags of the platform.
func (i *PlatformInfo) Tune() TuneFlag { return TuneFlag(i.Platform.tune) }

// CpuFeatures returns the CPU features of the platform.
func (i *PlatformInfo) CpuFeatures() CpuFeature { //nolint: golint,stylecheck,revive
	return CpuFeature(i.Platform.cpu_features)
}

func NewPlatformInfo(tune TuneFlag, cpu CpuFeature) *PlatformInfo {
	var platform C.struct_hs_platform_info

	platform.tune = C.uint(tune)
	platform.cpu_features = C.ulonglong(cpu)

	return &PlatformInfo{platform}
}

func PopulatePlatform() (*PlatformInfo, error) {
	var platform C.struct_hs_platform_info

	if ret := C.hs_populate_platform(&platform); ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return &PlatformInfo{platform}, nil
}

// ExprInfo containing information related to an expression
type ExprInfo struct {
	// The minimum length in bytes of a match for the pattern.
	MinWidth uint
	// The maximum length in bytes of a match for the pattern.
	MaxWidth uint
	// Whether this expression can produce matches that are not returned in order,
	// such as those produced by assertions.
	ReturnUnordered bool
	// Whether this expression can produce matches at end of data (EOD).
	AtEndOfData bool
	// Whether this expression can *only* produce matches at end of data (EOD).
	OnlyAtEndOfData bool
}

// UnboundedMaxWidth represents the pattern expression has an unbounded maximum width
const UnboundedMaxWidth = C.UINT_MAX

func NewExprInfo(info *C.hs_expr_info_t) *ExprInfo {
	return &ExprInfo{
		MinWidth:        uint(info.min_width),
		MaxWidth:        uint(info.max_width),
		ReturnUnordered: info.unordered_matches != 0,
		AtEndOfData:     info.matches_at_eod != 0,
		OnlyAtEndOfData: info.matches_only_at_eod != 0,
	}
}

// ExtFlag are used in ExprExt.Flags to indicate which fields are used.
type ExtFlag uint64

const (
	// ExtMinOffset is a flag indicating that the ExprExt.MinOffset field is used.
	ExtMinOffset ExtFlag = C.HS_EXT_FLAG_MIN_OFFSET
	// ExtMaxOffset is a flag indicating that the ExprExt.MaxOffset field is used.
	ExtMaxOffset ExtFlag = C.HS_EXT_FLAG_MAX_OFFSET
	// ExtMinLength is a flag indicating that the ExprExt.MinLength field is used.
	ExtMinLength ExtFlag = C.HS_EXT_FLAG_MIN_LENGTH
	// ExtEditDistance is a flag indicating that the ExprExt.EditDistance field is used.
	ExtEditDistance ExtFlag = C.HS_EXT_FLAG_EDIT_DISTANCE
	// ExtHammingDistance is a flag indicating that the ExprExt.HammingDistance field is used.
	ExtHammingDistance ExtFlag = C.HS_EXT_FLAG_HAMMING_DISTANCE
)

// ExprExt is a structure containing additional parameters related to an expression.
type ExprExt struct {
	// Flags governing which parts of this structure are to be used by the compiler.
	Flags ExtFlag
	// The minimum end offset in the data stream at which this expression should match successfully.
	MinOffset uint64
	// The maximum end offset in the data stream at which this expression should match successfully.
	MaxOffset uint64
	// The minimum match length (from start to end) required to successfully match this expression.
	MinLength uint64
	// Allow patterns to approximately match within this edit distance.
	EditDistance uint32
	// Allow patterns to approximately match within this Hamming distance.
	HammingDistance uint32
}

func (e *ExprExt) c() *C.hs_expr_ext_t {
	if e == nil {
		return nil
	}

	var ext C.hs_expr_ext_t

	if e.Flags&ExtMinOffset != 0 {
		ext.flags |= C.HS_EXT_FLAG_MIN_OFFSET
		ext.min_offset = C.ulonglong(e.MinOffset)
	}
	if e.Flags&ExtMaxOffset != 0 {
		ext.flags |= C.HS_EXT_FLAG_MAX_OFFSET
		ext.max_offset = C.ulonglong(e.MaxOffset)
	}
	if e.Flags&ExtMinLength != 0 {
		ext.flags |= C.HS_EXT_FLAG_MIN_LENGTH
		ext.min_length = C.ulonglong(e.MinLength)
	}
	if e.Flags&ExtEditDistance != 0 {
		ext.flags |= C.HS_EXT_FLAG_EDIT_DISTANCE
		ext.edit_distance = C.uint(e.EditDistance)
	}
	if e.Flags&ExtHammingDistance != 0 {
		ext.flags |= C.HS_EXT_FLAG_HAMMING_DISTANCE
		ext.hamming_distance = C.uint(e.HammingDistance)
	}

	return &ext
}

func ExpressionInfo(expression string, flags CompileFlag) (*ExprInfo, error) {
	var info *C.hs_expr_info_t
	var err *C.hs_compile_error_t

	expr := C.CString(expression)

	defer C.free(unsafe.Pointer(expr))

	ret := C.hs_expression_info(expr, C.uint(flags), &info, &err)

	if ret == C.HS_SUCCESS && info != nil {
		defer MiscAllocator().Free(unsafe.Pointer(info))

		return NewExprInfo(info), nil
	}

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &CompileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}

func ExpressionExt(expression string, flags CompileFlag) (ext *ExprExt, info *ExprInfo, err error) {
	var exprInfo *C.hs_expr_info_t
	var compileErr *C.hs_compile_error_t

	ext = new(ExprExt)
	expr := C.CString(expression)

	defer C.free(unsafe.Pointer(expr))

	ret := C.hs_expression_ext_info(expr, C.uint(flags), (*C.hs_expr_ext_t)(unsafe.Pointer(ext)), &exprInfo, &compileErr)

	if exprInfo != nil {
		defer MiscAllocator().Free(unsafe.Pointer(exprInfo))

		info = NewExprInfo(exprInfo)
	}

	if compileErr != nil {
		defer C.hs_free_compile_error(compileErr)
	}

	if ret == C.HS_COMPILER_ERROR && compileErr != nil {
		err = &CompileError{C.GoString(compileErr.message), int(compileErr.expression)}
	} else {
		err = fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
	}

	return
}

// CompileFlag represents a pattern flag
type CompileFlag uint

const (
	// Caseless represents set case-insensitive matching.
	Caseless CompileFlag = C.HS_FLAG_CASELESS
	// DotAll represents matching a `.` will not exclude newlines.
	DotAll CompileFlag = C.HS_FLAG_DOTALL
	// MultiLine set multi-line anchoring.
	MultiLine CompileFlag = C.HS_FLAG_MULTILINE
	// SingleMatch set single-match only mode.
	SingleMatch CompileFlag = C.HS_FLAG_SINGLEMATCH
	// AllowEmpty allow expressions that can match against empty buffers.
	AllowEmpty CompileFlag = C.HS_FLAG_ALLOWEMPTY
	// Utf8Mode enable UTF-8 mode for this expression.
	Utf8Mode CompileFlag = C.HS_FLAG_UTF8
	// UnicodeProperty enable Unicode property support for this expression.
	UnicodeProperty CompileFlag = C.HS_FLAG_UCP
	// PrefilterMode enable prefiltering mode for this expression.
	PrefilterMode CompileFlag = C.HS_FLAG_PREFILTER
	// SomLeftMost enable leftmost start of match reporting.
	SomLeftMost CompileFlag = C.HS_FLAG_SOM_LEFTMOST
)

var CompileFlags = map[rune]CompileFlag{
	'i': Caseless,
	's': DotAll,
	'm': MultiLine,
	'H': SingleMatch,
	'V': AllowEmpty,
	'8': Utf8Mode,
	'W': UnicodeProperty,
	'P': PrefilterMode,
	'L': SomLeftMost,
}

var DeprecatedCompileFlags = map[rune]CompileFlag{
	'o': SingleMatch,
	'e': AllowEmpty,
	'u': Utf8Mode,
	'p': UnicodeProperty,
	'f': PrefilterMode,
	'l': SomLeftMost,
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

// ModeFlag represents the compile mode flags
type ModeFlag uint

const (
	// BlockMode for the block scan (non-streaming) database.
	BlockMode ModeFlag = C.HS_MODE_BLOCK
	// NoStreamMode is alias for Block.
	NoStreamMode ModeFlag = C.HS_MODE_NOSTREAM
	// StreamMode for the streaming database.
	StreamMode ModeFlag = C.HS_MODE_STREAM
	// VectoredMode for the vectored scanning database.
	VectoredMode ModeFlag = C.HS_MODE_VECTORED
	// SomHorizonLargeMode use full precision to track start of match offsets in stream state.
	SomHorizonLargeMode ModeFlag = C.HS_MODE_SOM_HORIZON_LARGE
	// SomHorizonMediumMode use medium precision to track start of match offsets in stream state. (within 2^32 bytes)
	SomHorizonMediumMode ModeFlag = C.HS_MODE_SOM_HORIZON_MEDIUM
	// SomHorizonSmallMode use limited precision to track start of match offsets in stream state. (within 2^16 bytes)
	SomHorizonSmallMode ModeFlag = C.HS_MODE_SOM_HORIZON_SMALL
	// ModeMask represents the mask of database mode
	ModeMask ModeFlag = 0xFF
)

var ModeFlags = map[string]ModeFlag{
	"STREAM":   StreamMode,
	"NOSTREAM": BlockMode,
	"VECTORED": VectoredMode,
	"BLOCK":    BlockMode,
}

func (m ModeFlag) String() string {
	switch m & 0xF {
	case BlockMode:
		return "BLOCK"
	case StreamMode:
		return "STREAM"
	case VectoredMode:
		return "VECTORED"
	default:
		panic(fmt.Sprintf("unknown mode: %d", m))
	}
}

func Compile(expression string, flags CompileFlag, mode ModeFlag, info *PlatformInfo) (Database, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = (*C.struct_hs_platform_info)(unsafe.Pointer(&info.Platform))
	}

	expr := C.CString(expression)

	defer C.free(unsafe.Pointer(expr))

	ret := C.hs_compile(expr, C.uint(flags), C.uint(mode), platform, &db, &err)

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_SUCCESS {
		return db, nil
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &CompileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}

type Pattern struct {
	Expr  string
	Flags CompileFlag
	ID    int
	Ext   *ExprExt
}

type Patterns interface {
	Patterns() []*Pattern
}

func CompileMulti(input Patterns, mode ModeFlag, info *PlatformInfo) (Database, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = (*C.struct_hs_platform_info)(unsafe.Pointer(&info.Platform))
	}

	patterns := input.Patterns()
	count := len(patterns)

	cexprs := (**C.char)(C.calloc(C.size_t(count), C.size_t(unsafe.Sizeof(uintptr(0)))))
	exprs := (*[1 << 30]*C.char)(unsafe.Pointer(cexprs))[:count:count]

	cflags := (*C.uint)(C.calloc(C.size_t(count), C.size_t(unsafe.Sizeof(C.uint(0)))))
	flags := (*[1 << 30]C.uint)(unsafe.Pointer(cflags))[:count:count]

	cids := (*C.uint)(C.calloc(C.size_t(count), C.size_t(unsafe.Sizeof(C.uint(0)))))
	ids := (*[1 << 30]C.uint)(unsafe.Pointer(cids))[:count:count]

	cexts := (**C.hs_expr_ext_t)(C.calloc(C.size_t(count), C.size_t(unsafe.Sizeof(uintptr(0)))))
	exts := (*[1 << 30]*C.hs_expr_ext_t)(unsafe.Pointer(cexts))[:count:count]

	for i, pattern := range patterns {
		exprs[i] = C.CString(pattern.Expr)
		flags[i] = C.uint(pattern.Flags)
		ids[i] = C.uint(pattern.ID)
		exts[i] = pattern.Ext.c()
	}

	ret := C.hs_compile_ext_multi(cexprs, cflags, cids, cexts, C.uint(count), C.uint(mode), platform, &db, &err)

	for _, expr := range exprs {
		C.free(unsafe.Pointer(expr))
	}

	C.free(unsafe.Pointer(cexprs))
	C.free(unsafe.Pointer(cflags))
	C.free(unsafe.Pointer(cexts))
	C.free(unsafe.Pointer(cids))

	runtime.KeepAlive(patterns)

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_SUCCESS {
		return db, nil
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &CompileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}
