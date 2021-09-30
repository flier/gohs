package ch

import (
	"fmt"
	"reflect"
	"runtime"
	"unsafe"

	"github.com/flier/gohs/internal/handle"
)

/*
#include <stdint.h>

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

static inline
ch_error_t HS_CDECL _ch_scan(const ch_database_t *db, const char *data,
                            unsigned int length, unsigned int flags,
                            ch_scratch_t *scratch,
                            ch_match_event_handler onEvent,
                            ch_error_event_handler onError,
                            uintptr_t context) {
	return ch_scan(db, data, length, flags, scratch, onEvent, onError, (void *)context);
}
*/
import "C"

// Callback return value used to tell the Chimera matcher what to do after processing this match.
type Callback C.ch_callback_t

const (
	Continue    Callback = C.CH_CALLBACK_CONTINUE     // Continue matching.
	Terminate   Callback = C.CH_CALLBACK_TERMINATE    // Terminate matching.
	SkipPattern Callback = C.CH_CALLBACK_SKIP_PATTERN // Skip remaining matches for this ID and continue.
)

// Capture representing a captured subexpression within a match.
type Capture struct {
	From  uint64 // offset at which this capture group begins.
	To    uint64 // offset at which this capture group ends.
	Bytes []byte // matches of the expression
}

// Definition of the match event callback function type.
type MatchEventHandler func(id uint, from, to uint64, flags uint, captured []*Capture, context interface{}) Callback

// Type used to differentiate the errors raised with the `ErrorEventHandler` callback.
type ErrorEvent C.ch_error_event_t

const (
	// PCRE hits its match limit and reports PCRE_ERROR_MATCHLIMIT.
	ErrMatchLimit ErrorEvent = C.CH_ERROR_MATCHLIMIT
	// PCRE hits its recursion limit and reports PCRE_ERROR_RECURSIONLIMIT.
	ErrRecursionLimit ErrorEvent = C.CH_ERROR_RECURSIONLIMIT
)

func (e ErrorEvent) Error() string {
	switch e {
	case ErrMatchLimit:
		return "PCRE hits its match limit."
	case ErrRecursionLimit:
		return "PCRE hits its recursion limit."
	}

	return fmt.Sprintf("ErrorEvent(%d)", int(C.ch_error_event_t(e)))
}

// Definition of the Chimera error event callback function type.
type ErrorEventHandler func(event ErrorEvent, id uint, info, context interface{}) Callback

type eventContext struct {
	data    []byte
	onEvent MatchEventHandler
	onError ErrorEventHandler
	context interface{}
}

//export matchEventCallback
func matchEventCallback(id C.uint, from, to C.ulonglong, flags, size C.uint,
	cap *C.capture_t, data unsafe.Pointer) C.ch_callback_t {
	ctx, ok := handle.Handle(data).Value().(eventContext)
	if !ok {
		return C.CH_CALLBACK_TERMINATE
	}

	captured := make([]*Capture, size)
	for i, c := range (*[1 << 30]C.capture_t)(unsafe.Pointer(cap))[:size:size] {
		if c.flags == C.CH_CAPTURE_FLAG_ACTIVE {
			captured[i] = &Capture{uint64(c.from), uint64(c.to), ctx.data[c.from:c.to]}
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
type ScanFlag uint

// The block regular expression scanner.
func Scan(db Database, data []byte, flags ScanFlag, scratch Scratch,
	onEvent MatchEventHandler, onError ErrorEventHandler, context interface{}) error {
	if data == nil {
		return ErrInvalid
	}

	h := handle.New(eventContext{data, onEvent, onError, context})
	defer h.Delete()

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data)) // FIXME: Zero-copy access to go data

	ret := C._ch_scan(db,
		(*C.char)(unsafe.Pointer(hdr.Data)),
		C.uint(hdr.Len),
		C.uint(flags),
		scratch,
		C.ch_match_event_handler(C.matchEventCallback),
		C.ch_error_event_handler(C.errorEventCallback),
		C.uintptr_t(h))

	// Ensure go data is alive before the C function returns
	runtime.KeepAlive(data)

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

type MatchEvent struct {
	ID       uint
	From, To uint64
	Flag     ScanFlag
	Captured []*Capture
}

type MatchRecorder struct {
	Events []MatchEvent
	Err    error
}

func (h *MatchRecorder) Matched() bool { return len(h.Events) > 0 }

func (h *MatchRecorder) OnMatch(id uint, from, to uint64, flags uint, captured []*Capture, ctx interface{}) Callback {
	h.Events = append(h.Events, MatchEvent{id, from, to, ScanFlag(flags), captured})

	return Continue
}

func (h *MatchRecorder) OnError(evt ErrorEvent, id uint, info, ctx interface{}) Callback {
	h.Err = evt

	return Terminate
}
