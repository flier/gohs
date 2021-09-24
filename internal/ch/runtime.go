package ch

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/flier/gohs/internal/handle"
)

/*
#include <ch.h>

typedef const ch_capture_t capture_t;

extern ch_callback_t matchEventCallback(unsigned int id,
										unsigned long long from,
										unsigned long long to,
										unsigned int flags,
										unsigned int size,
										const ch_capture_t *captured,
										void *ctx);

extern ch_callback_t errorEventCallback(ch_error_event_t error_type,
                                        unsigned int id,
										void *info,
                                        void *ctx);
*/
import "C"

// A Chimera scratch space.
type Scratch *C.ch_scratch_t

// Allocate a scratch space that is a clone of an existing scratch space.
func CloneScratch(scratch Scratch) (Scratch, error) {
	var clone *C.ch_scratch_t

	if ret := C.ch_clone_scratch(scratch, &clone); ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return clone, nil
}

// Provides the size of the given scratch space.
func ScratchSize(scratch Scratch) (int, error) {
	var size C.size_t

	if ret := C.ch_scratch_size(scratch, &size); ret != C.HS_SUCCESS {
		return 0, Error(ret)
	}

	return int(size), nil
}

// Free a scratch block previously allocated by @ref ch_alloc_scratch() or @ref ch_clone_scratch().
func FreeScratch(scratch Scratch) error {
	if ret := C.ch_free_scratch(scratch); ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

// Callback return value used to tell the Chimera matcher what to do after processing this match.
type Callback C.ch_callback_t

const (
	Continue    Callback = C.CH_CALLBACK_CONTINUE     // Continue matching.
	Terminate   Callback = C.CH_CALLBACK_TERMINATE    // Terminate matching.
	SkipPattern Callback = C.CH_CALLBACK_SKIP_PATTERN // Skip remaining matches for this ID and continue.
)

// Capture representing a captured subexpression within a match.
type Capture struct {
	From uint64 // offset at which this capture group begins.
	To   uint64 // offset at which this capture group ends.
}

// Definition of the match event callback function type.
type MatchEventHandler func(id uint, from, to uint64, flags uint, captured []*Capture, context interface{}) Callback

type ErrorEvent C.ch_error_event_t

const (
	// PCRE hits its match limit and reports PCRE_ERROR_MATCHLIMIT.
	MatchLimit ErrorEvent = C.CH_ERROR_MATCHLIMIT
	// PCRE hits its recursion limit and reports PCRE_ERROR_RECURSIONLIMIT.
	RecursionLimit ErrorEvent = C.CH_ERROR_RECURSIONLIMIT
)

// Definition of the Chimera error event callback function type.
type ErrorEventHandler func(event ErrorEvent, id uint, info, context interface{}) Callback

type eventContext struct {
	onEvent MatchEventHandler
	onError ErrorEventHandler
	context interface{}
}

//export matchEventCallback
func matchEventCallback(id C.uint, from, to C.ulonglong, flags, size C.uint, cap *C.capture_t, data unsafe.Pointer) C.ch_callback_t {
	ctx, ok := handle.Handle(data).Value().(eventContext)
	if !ok {
		return C.CH_CALLBACK_TERMINATE
	}

	captured := make([]*Capture, size)
	for i, c := range (*[1 << 30]C.capture_t)(unsafe.Pointer(cap))[:size:size] {
		if c.flags == C.CH_CAPTURE_FLAG_ACTIVE {
			captured[i] = &Capture{uint64(c.from), uint64(c.to)}
		}
	}

	return C.ch_callback_t(ctx.onEvent(uint(id), uint64(from), uint64(to), uint(flags), captured, ctx.context))
}

//export errorEventCallback
func errorEventCallback(evt C.ch_error_event_t, id C.uint, info, data unsafe.Pointer) C.ch_callback_t {
	ctx, ok := handle.Handle(data).Value().(eventContext)
	if !ok {
		return C.CH_CALLBACK_TERMINATE
	}

	return C.ch_callback_t(ctx.onError(ErrorEvent(evt), uint(id), nil, ctx.context))
}

// ScanFlag represents a scan flag.
type ScanFlag C.uint

// The block regular expression scanner.
func hsScan(db Database, data []byte, flags ScanFlag, scratch Scratch,
	onEvent MatchEventHandler, onError ErrorEventHandler, context interface{}) error {
	if data == nil {
		return ErrInvalid
	}

	h := handle.New(eventContext{onEvent, onError, context})
	defer h.Delete()

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data)) // FIXME: Zero-copy access to go data

	ret := C.ch_scan(db,
		(*C.char)(unsafe.Pointer(hdr.Data)),
		C.uint(hdr.Len),
		C.uint(flags),
		scratch,
		C.ch_match_event_handler(C.matchEventCallback),
		C.ch_error_event_handler(C.errorEventCallback),
		unsafe.Pointer(h))

	// Ensure go data is alive before the C function returns
	runtime.KeepAlive(data)

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}
