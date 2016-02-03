package hyperscan

import (
	"fmt"
	"unsafe"
)

/*
#cgo LDFLAGS: -lstdc++
#cgo pkg-config: libhs

#include <hs.h>

#define DEFINE_ALLOCTOR(ID, TYPE) \
	extern void *hs ## ID ## Alloc(size_t size); \
	extern void hs ## ID ## Free(void *ptr); \
	static inline void *hs ## ID ## Alloc_cgo(size_t size) { return hs ## ID ## Alloc(size); } \
	static inline void hs ## ID ## Free_cgo(void *ptr) { hs ## ID ## Free(ptr); } \
	static inline hs_error_t hs_set_ ## TYPE ## _allocator_cgo(void *alloc, void *free) \
	{ return hs_set_ ## TYPE ## _allocator(alloc ? hs ## ID ## Alloc_cgo : NULL, free ? hs ## ID ## Free_cgo : NULL); }

DEFINE_ALLOCTOR(Db, database);
DEFINE_ALLOCTOR(Misc, misc);
DEFINE_ALLOCTOR(Scratch, scratch);
DEFINE_ALLOCTOR(Stream, stream);

extern int hsMatchEventCallback(unsigned int id, unsigned long long from, unsigned long long to, unsigned int flags, void *context);

static inline
hs_error_t hs_scan_cgo(const hs_database_t *db, const char *data, unsigned int length,
					   unsigned int flags, hs_scratch_t *scratch, void *context) {
	return hs_scan(db, data, length, flags, scratch, hsMatchEventCallback, context);
}

static inline
hs_error_t hs_scan_vector_cgo(const hs_database_t *db, const char *const *data, const unsigned int *length, unsigned int count,
					   		  unsigned int flags, hs_scratch_t *scratch, void *context) {
	return hs_scan_vector(db, data, length, count, flags, scratch, hsMatchEventCallback, context);
}

static inline
hs_error_t hs_scan_stream_cgo(hs_stream_t *id, const char *data, unsigned int length,
							  unsigned int flags, hs_scratch_t *scratch, void *context) {
	return hs_scan_stream(id, data, length, flags, scratch, hsMatchEventCallback, context);
}

static inline
hs_error_t hs_close_stream_cgo(hs_stream_t *id, hs_scratch_t *scratch, void *context) {
	return hs_close_stream(id, scratch, hsMatchEventCallback, context);
}

static inline
hs_error_t hs_reset_stream_cgo(hs_stream_t *id, unsigned int flags, hs_scratch_t *scratch, void *context) {
	return hs_reset_stream(id, flags, scratch, hsMatchEventCallback, context);
}

static inline
hs_error_t hs_reset_and_copy_stream_cgo(hs_stream_t *to_id, const hs_stream_t *from_id, hs_scratch_t *scratch, void *context) {
	return hs_reset_and_copy_stream(to_id, from_id, scratch, hsMatchEventCallback, context);
}
*/
import "C"

type CompileFlag uint

const (
	Caseless    CompileFlag = C.HS_FLAG_CASELESS     // Set case-insensitive matching.
	DotAll                  = C.HS_FLAG_DOTALL       // Matching a `.` will not exclude newlines.
	MultiLine               = C.HS_FLAG_MULTILINE    // Set multi-line anchoring.
	SingleMatch             = C.HS_FLAG_SINGLEMATCH  // Set single-match only mode.
	AllowEmpty              = C.HS_FLAG_ALLOWEMPTY   // Allow expressions that can match against empty buffers.
	Utf8                    = C.HS_FLAG_UTF8         // Enable UTF-8 mode for this expression.
	Ucp                     = C.HS_FLAG_UCP          // Enable Unicode property support for this expression.
	Prefilter               = C.HS_FLAG_PREFILTER    // Enable prefiltering mode for this expression.
	SomLeftMost             = C.HS_FLAG_SOM_LEFTMOST // Enable leftmost start of match reporting.
)

type CpuFeature int

const (
	AVX2 CpuFeature = C.HS_CPU_FEATURES_AVX2 // Intel(R) Advanced Vector Extensions 2 (Intel(R) AVX2)
)

type TuneFlag int

const (
	Generic     TuneFlag = C.HS_TUNE_FAMILY_GENERIC // Generic
	SandyBridge          = C.HS_TUNE_FAMILY_SNB     // Intel(R) microarchitecture code name Sandy Bridge
	IvyBridge            = C.HS_TUNE_FAMILY_IVB     // Intel(R) microarchitecture code name Ivy Bridge
	Haswell              = C.HS_TUNE_FAMILY_HSW     // Intel(R) microarchitecture code name Haswell
	Silvermont           = C.HS_TUNE_FAMILY_SLM     // Intel(R) microarchitecture code name Silvermont
	Broadwell            = C.HS_TUNE_FAMILY_BDW     // Intel(R) microarchitecture code name Broadwell
)

type ModeFlag uint

const (
	Block    ModeFlag = C.HS_MODE_BLOCK    // Block scan (non-streaming) database.
	NoStream          = C.HS_MODE_NOSTREAM // Alias for Block.
	Stream            = C.HS_MODE_STREAM   // Streaming database.
	Vectored          = C.HS_MODE_VECTORED // Vectored scanning database.
)

type ScanFlag uint

type hsError int

const (
	Success               hsError = C.HS_SUCCESS
	Invalid                       = C.HS_INVALID
	NoMemory                      = C.HS_NOMEM
	ScanTerminated                = C.HS_SCAN_TERMINATED
	CompileError                  = C.HS_COMPILER_ERROR
	DatabaseVersionError          = C.HS_DB_VERSION_ERROR
	DatabasePlatformError         = C.HS_DB_PLATFORM_ERROR
	DatabaseModeError             = C.HS_DB_MODE_ERROR
	BadAlign                      = C.HS_BAD_ALIGN
	BadAlloc                      = C.HS_BAD_ALLOC
)

var (
	hsErrorMessages = map[hsError]string{
		C.HS_SUCCESS:           "The engine completed normally.",
		C.HS_INVALID:           "A parameter passed to this function was invalid.",
		C.HS_NOMEM:             "A memory allocation failed.",
		C.HS_SCAN_TERMINATED:   "The engine was terminated by callback.",
		C.HS_COMPILER_ERROR:    "The pattern compiler failed.",
		C.HS_DB_VERSION_ERROR:  "The given database was built for a different version of Hyperscan.",
		C.HS_DB_PLATFORM_ERROR: "The given database was built for a different platform (i.e., CPU type).",
		C.HS_DB_MODE_ERROR:     "The given database was built for a different mode of operation.",
		C.HS_BAD_ALIGN:         "A parameter passed to this function was not correctly aligned.",
		C.HS_BAD_ALLOC:         "The memory allocator did not correctly return aligned memory.",
	}
)

func (e hsError) Error() string {
	if msg, exists := hsErrorMessages[e]; exists {
		return msg
	}

	return fmt.Sprintf("unexpected error, %d", int(e))
}

type compileError struct {
	msg  string
	expr int
}

// A human-readable error message describing the error.
func (e *compileError) Error() string { return e.msg }

// The zero-based number of the expression that caused the error (if this can be determined).
// If the error is not specific to an expression, then this value will be less than zero.
func (e *compileError) Expression() int { return e.expr }

type hsPlatformInfo struct {
	info C.struct_hs_platform_info
}

func hsPopulatePlatform() (*hsPlatformInfo, error) {
	var platform hsPlatformInfo

	if ret := C.hs_populate_platform(&platform.info); ret != C.HS_SUCCESS {
		return &platform, hsError(ret)
	}

	return &platform, nil
}

type hsDatabase *C.hs_database_t
type hsScratch *C.hs_scratch_t
type hsStream *C.hs_stream_t

type hsExprInfo struct {
	MinWidth, MaxWidth          uint
	Unordered, AtEod, OnlyAtEod bool
}

type hsAllocFunc func(uint) unsafe.Pointer
type hsFreeFunc func(unsafe.Pointer)

type hsAllocator struct {
	Alloc hsAllocFunc
	Free  hsFreeFunc
}

var (
	defaultAllocator hsAllocator
	dbAllocator      hsAllocator
	miscAllocator    hsAllocator
	scratchAllocator hsAllocator
	streamAllocator  hsAllocator
)

func hsDefaultAlloc(size uint) unsafe.Pointer {
	var ptr unsafe.Pointer

	if ret := C.posix_memalign(&ptr, 64, C.size_t(size)); ret == 0 {
		return ptr
	}

	return nil
}

func hsDefaultFree(ptr unsafe.Pointer) {
	C.free(ptr)
}

//export hsDbAlloc
func hsDbAlloc(size C.size_t) unsafe.Pointer {
	if dbAllocator.Alloc != nil {
		return dbAllocator.Alloc(uint(size))
	}

	return hsDefaultAlloc(uint(size))
}

//export hsDbFree
func hsDbFree(ptr unsafe.Pointer) {
	if dbAllocator.Free != nil {
		dbAllocator.Free(ptr)
	} else {
		hsDefaultFree(ptr)
	}
}

func hsSetDatabaseAllocator(allocFunc hsAllocFunc, freeFunc hsFreeFunc) error {
	dbAllocator = hsAllocator{allocFunc, freeFunc}

	if ret := C.hs_set_database_allocator_cgo(unsafe.Pointer(&dbAllocator.Alloc), unsafe.Pointer(&dbAllocator.Free)); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

//export hsMiscAlloc
func hsMiscAlloc(size C.size_t) unsafe.Pointer {
	if miscAllocator.Alloc != nil {
		return miscAllocator.Alloc(uint(size))
	}

	return hsDefaultAlloc(uint(size))
}

//export hsMiscFree
func hsMiscFree(ptr unsafe.Pointer) {
	if miscAllocator.Free != nil {
		miscAllocator.Free(ptr)
	} else {
		hsDefaultFree(ptr)
	}
}

func hsSetMiscAllocator(allocFunc hsAllocFunc, freeFunc hsFreeFunc) error {
	miscAllocator = hsAllocator{allocFunc, freeFunc}

	if ret := C.hs_set_misc_allocator_cgo(unsafe.Pointer(&miscAllocator.Alloc), unsafe.Pointer(&miscAllocator.Free)); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

//export hsScratchAlloc
func hsScratchAlloc(size C.size_t) unsafe.Pointer {
	if scratchAllocator.Alloc != nil {
		return scratchAllocator.Alloc(uint(size))
	}

	return hsDefaultAlloc(uint(size))
}

//export hsScratchFree
func hsScratchFree(ptr unsafe.Pointer) {
	if scratchAllocator.Free != nil {
		scratchAllocator.Free(ptr)
	} else {
		hsDefaultFree(ptr)
	}
}

func hsSetScratchAllocator(allocFunc hsAllocFunc, freeFunc hsFreeFunc) error {
	scratchAllocator = hsAllocator{allocFunc, freeFunc}

	if ret := C.hs_set_scratch_allocator_cgo(unsafe.Pointer(&scratchAllocator.Alloc), unsafe.Pointer(&scratchAllocator.Free)); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

//export hsStreamAlloc
func hsStreamAlloc(size C.size_t) unsafe.Pointer {
	if streamAllocator.Alloc != nil {
		return streamAllocator.Alloc(uint(size))
	}

	return hsDefaultAlloc(uint(size))
}

//export hsStreamFree
func hsStreamFree(ptr unsafe.Pointer) {
	if streamAllocator.Free != nil {
		streamAllocator.Free(ptr)
	} else {
		hsDefaultFree(ptr)
	}
}

func hsSetStreamAllocator(allocFunc hsAllocFunc, freeFunc hsFreeFunc) error {
	streamAllocator = hsAllocator{allocFunc, freeFunc}

	if ret := C.hs_set_stream_allocator_cgo(unsafe.Pointer(&streamAllocator.Alloc), unsafe.Pointer(&streamAllocator.Free)); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsVersion() string {
	return C.GoString(C.hs_version())
}

func hsFreeDatabase(db hsDatabase) error {
	if ret := C.hs_free_database(db); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsSerializeDatabase(db hsDatabase) ([]byte, error) {
	var data *C.char
	var length C.size_t

	if ret := C.hs_serialize_database(db, &data, &length); ret != C.HS_SUCCESS {
		return nil, hsError(ret)
	}

	return C.GoBytes(unsafe.Pointer(data), C.int(length)), nil
}

func hsDeserializeDatabase(data []byte) (hsDatabase, error) {
	var db *C.hs_database_t

	if ret := C.hs_deserialize_database((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &db); ret != C.HS_SUCCESS {
		return nil, hsError(ret)
	}

	return db, nil
}

func hsDeserializeDatabaseAt(data []byte, db hsDatabase) error {
	if ret := C.hs_deserialize_database_at((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), db); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsStreamSize(db hsDatabase) (int, error) {
	var size C.size_t

	if ret := C.hs_stream_size(db, &size); ret != C.HS_SUCCESS {
		return 0, hsError(ret)
	}

	return int(size), nil
}

func hsDatabaseSize(db hsDatabase) (int, error) {
	var size C.size_t

	if ret := C.hs_database_size(db, &size); ret != C.HS_SUCCESS {
		return -1, hsError(ret)
	}

	return int(size), nil
}

func hsSerializedDatabaseSize(data []byte) (int, error) {
	var size C.size_t

	if ret := C.hs_serialized_database_size((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &size); ret != C.HS_SUCCESS {
		return 0, hsError(ret)
	}

	return int(size), nil
}

func hsDatabaseInfo(db hsDatabase) (string, error) {
	var info *C.char

	if ret := C.hs_database_info(db, &info); ret != C.HS_SUCCESS {
		return "", hsError(ret)
	}

	return C.GoString(info), nil
}

func hsSerializedDatabaseInfo(data []byte) (string, error) {
	var info *C.char

	if ret := C.hs_serialized_database_info((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &info); ret != C.HS_SUCCESS {
		return "", hsError(ret)
	}

	return C.GoString(info), nil
}

func hsCompile(expression string, flags CompileFlag, mode ModeFlag, platform *hsPlatformInfo) (hsDatabase, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t

	expr := C.CString(expression)

	ret := C.hs_compile(expr, C.uint(flags), C.uint(mode), &platform.info, &db, &err)

	C.free(unsafe.Pointer(expr))

	if ret == C.HS_SUCCESS {
		return db, nil
	}

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &compileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error, %d", int(ret))
}

func hsCompileMulti(expressions []string, flags []CompileFlag, ids []uint, mode ModeFlag, platform *hsPlatformInfo) (hsDatabase, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t

	cexprs := make([]*C.char, len(expressions))

	for i, expr := range expressions {
		cexprs[i] = C.CString(expr)
	}

	var cflags, cids *C.uint

	if flags != nil {
		values := make([]C.uint, len(flags))

		for i, flag := range flags {
			values[i] = C.uint(flag)
		}

		cflags = &values[0]
	}

	if ids != nil {
		values := make([]C.uint, len(ids))

		for i, id := range ids {
			values[i] = C.uint(id)
		}

		cids = &values[0]
	}

	ret := C.hs_compile_multi(&cexprs[0], cflags, cids, C.uint(len(cexprs)), C.uint(mode), &platform.info, &db, &err)

	for _, expr := range cexprs {
		C.free(unsafe.Pointer(expr))
	}

	if ret == C.HS_SUCCESS {
		return db, nil
	}

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &compileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error, %d", int(ret))
}

func hsExpressionInfo(expression string, flags CompileFlag) (*hsExprInfo, error) {
	var info *C.hs_expr_info_t
	var err *C.hs_compile_error_t

	expr := C.CString(expression)

	ret := C.hs_expression_info(expr, C.uint(flags), &info, &err)

	C.free(unsafe.Pointer(expr))

	if ret == C.HS_SUCCESS && info != nil {
		defer hsMiscFree(unsafe.Pointer(info))

		return &hsExprInfo{
			MinWidth:  uint(info.min_width),
			MaxWidth:  uint(info.max_width),
			Unordered: info.unordered_matches != 0,
			AtEod:     info.matches_at_eod != 0,
			OnlyAtEod: info.matches_only_at_eod != 0,
		}, nil
	}

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &compileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error, %d", int(ret))
}

func hsAllocScratch(db hsDatabase) (hsScratch, error) {
	var scratch *C.hs_scratch_t

	if ret := C.hs_alloc_scratch(db, &scratch); ret != C.HS_SUCCESS {
		return nil, hsError(ret)
	}

	return scratch, nil
}

func hsCloneScratch(scratch hsScratch) (hsScratch, error) {
	var clone *C.hs_scratch_t

	if ret := C.hs_clone_scratch(scratch, &clone); ret != C.HS_SUCCESS {
		return nil, hsError(ret)
	}

	return clone, nil
}

func hsScratchSize(scratch hsScratch) (int, error) {
	var size C.size_t

	if ret := C.hs_scratch_size(scratch, &size); ret != C.HS_SUCCESS {
		return 0, hsError(ret)
	}

	return int(size), nil
}

func hsFreeScratch(scratch hsScratch) error {
	if ret := C.hs_free_scratch(scratch); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

type hsMatchEventHandler func(id uint, from, to uint64, flags uint, context interface{}) error

type hsMatchEventContext struct {
	handler hsMatchEventHandler
	context interface{}
}

//export hsMatchEventCallback
func hsMatchEventCallback(id C.uint, from, to C.ulonglong, flags C.uint, context unsafe.Pointer) C.int {
	ctxt := (*hsMatchEventContext)(context)

	if err := ctxt.handler(uint(id), uint64(from), uint64(to), uint(flags), ctxt.context); err != nil {
		return -1
	}

	return 0
}

func hsScan(db hsDatabase, data []byte, flags ScanFlag, scratch hsScratch, onEvent hsMatchEventHandler, context interface{}) error {
	ctxt := &hsMatchEventContext{onEvent, context}

	if ret := C.hs_scan_cgo(db, (*C.char)(unsafe.Pointer(&data[0])), C.uint(len(data)), C.uint(flags), scratch, unsafe.Pointer(ctxt)); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsScanVector(db hsDatabase, data [][]byte, flags ScanFlag, scratch hsScratch, onEvent hsMatchEventHandler, context interface{}) error {
	cdata := make([]*C.char, len(data))
	clength := make([]C.uint, len(data))

	for i, d := range data {
		cdata[i] = (*C.char)(unsafe.Pointer(&d[0]))
		clength[i] = C.uint(len(d))
	}

	if ret := C.hs_scan_vector_cgo(db, (**C.char)(unsafe.Pointer(&cdata[0])), &clength[0], C.uint(len(data)), C.uint(flags),
		scratch, unsafe.Pointer(&hsMatchEventContext{onEvent, context})); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsOpenStream(db hsDatabase, flags ScanFlag) (hsStream, error) {
	var stream *C.hs_stream_t

	if ret := C.hs_open_stream(db, C.uint(flags), &stream); ret != C.HS_SUCCESS {
		return nil, hsError(ret)
	}

	return stream, nil
}

func hsScanStream(stream hsStream, data []byte, flags ScanFlag, scratch hsScratch, onEvent hsMatchEventHandler, context interface{}) error {
	if ret := C.hs_scan_stream_cgo(stream, (*C.char)(unsafe.Pointer(&data[0])), C.uint(len(data)), C.uint(flags),
		scratch, unsafe.Pointer(&hsMatchEventContext{onEvent, context})); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsCloseStream(stream hsStream, scratch hsScratch, onEvent hsMatchEventHandler, context interface{}) error {
	if ret := C.hs_close_stream_cgo(stream, scratch, unsafe.Pointer(&hsMatchEventContext{onEvent, context})); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsResetStream(stream hsStream, flags ScanFlag, scratch hsScratch, onEvent hsMatchEventHandler, context interface{}) error {
	if ret := C.hs_reset_stream_cgo(stream, C.uint(flags), scratch, unsafe.Pointer(&hsMatchEventContext{onEvent, context})); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}

func hsCopyStream(stream hsStream) (hsStream, error) {
	var copied *C.hs_stream_t

	if ret := C.hs_copy_stream(&copied, stream); ret != C.HS_SUCCESS {
		return nil, hsError(ret)
	}

	return copied, nil
}

func hsResetAndCopyStream(to, from hsStream, scratch hsScratch, onEvent hsMatchEventHandler, context interface{}) error {
	if ret := C.hs_reset_and_copy_stream_cgo(to, from, scratch, unsafe.Pointer(&hsMatchEventContext{onEvent, context})); ret != C.HS_SUCCESS {
		return hsError(ret)
	}

	return nil
}
