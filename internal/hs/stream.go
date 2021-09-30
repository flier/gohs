package hs

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/flier/gohs/internal/handle"
)

/*
#include <hs.h>

extern int hsMatchEventCallback(unsigned int id,
								unsigned long long from,
								unsigned long long to,
								unsigned int flags,
								void *context);

static inline
hs_error_t HS_CDECL _hs_scan_stream(hs_stream_t *id, const char *data,
                                   unsigned int length, unsigned int flags,
                                   hs_scratch_t *scratch,
                                   match_event_handler onEvent,
								   uintptr_t context) {
	return hs_scan_stream(id, data, length, flags, scratch, onEvent, (void *)context);
}

static inline
hs_error_t HS_CDECL _hs_close_stream(hs_stream_t *id, hs_scratch_t *scratch,
                                    match_event_handler onEvent,
									uintptr_t context) {
	return hs_close_stream(id, scratch, onEvent, (void *)context);
}

static inline
hs_error_t HS_CDECL _hs_reset_stream(hs_stream_t *id, unsigned int flags,
                                    hs_scratch_t *scratch,
                                    match_event_handler onEvent,
									uintptr_t context) {
	return hs_reset_stream(id, flags, scratch, onEvent, (void *)context);
}

static inline
hs_error_t HS_CDECL _hs_reset_and_copy_stream(hs_stream_t *to_id,
                                             const hs_stream_t *from_id,
                                             hs_scratch_t *scratch,
                                             match_event_handler onEvent,
                                             uintptr_t context) {
	return hs_reset_and_copy_stream(to_id, from_id, scratch, onEvent, (void *)context);
}

static inline
hs_error_t HS_CDECL _hs_reset_and_expand_stream(hs_stream_t *to_stream,
                                               const char *buf, size_t buf_size,
                                               hs_scratch_t *scratch,
                                               match_event_handler onEvent,
                                               uintptr_t context) {
	return hs_reset_and_expand_stream(to_stream, buf, buf_size, scratch, onEvent, (void *)context);
}
*/
import "C"

type Stream *C.hs_stream_t

func OpenStream(db Database, flags ScanFlag) (Stream, error) {
	var stream *C.hs_stream_t

	if ret := C.hs_open_stream(db, C.uint(flags), &stream); ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return stream, nil
}

func ScanStream(stream Stream, data []byte, flags ScanFlag, scratch Scratch, onEvent MatchEventHandler, context interface{}) error {
	if data == nil {
		return Error(C.HS_INVALID)
	}

	h := handle.New(MatchEventContext{onEvent, context})
	defer h.Delete()

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data)) // FIXME: Zero-copy access to go data

	ret := C._hs_scan_stream(stream,
		(*C.char)(unsafe.Pointer(hdr.Data)),
		C.uint(hdr.Len),
		C.uint(flags),
		scratch,
		C.match_event_handler(C.hsMatchEventCallback),
		C.uintptr_t(h))

	// Ensure go data is alive before the C function returns
	runtime.KeepAlive(data)

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func FreeStream(stream Stream) {
	C.hs_close_stream(stream, nil, nil, nil)
}

func CloseStream(stream Stream, scratch Scratch, onEvent MatchEventHandler, context interface{}) error {
	h := handle.New(MatchEventContext{onEvent, context})
	defer h.Delete()

	ret := C._hs_close_stream(stream,
		scratch,
		C.match_event_handler(C.hsMatchEventCallback),
		C.uintptr_t(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ResetStream(stream Stream, flags ScanFlag, scratch Scratch, onEvent MatchEventHandler, context interface{}) error {
	h := handle.New(MatchEventContext{onEvent, context})
	defer h.Delete()

	ret := C._hs_reset_stream(stream,
		C.uint(flags),
		scratch,
		C.match_event_handler(C.hsMatchEventCallback),
		C.uintptr_t(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func CopyStream(stream Stream) (Stream, error) {
	var copied *C.hs_stream_t

	if ret := C.hs_copy_stream(&copied, stream); ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return copied, nil
}

func ResetAndCopyStream(to, from Stream, scratch Scratch, onEvent MatchEventHandler, context interface{}) error {
	h := handle.New(MatchEventContext{onEvent, context})
	defer h.Delete()

	ret := C._hs_reset_and_copy_stream(to,
		from,
		scratch,
		C.match_event_handler(C.hsMatchEventCallback),
		C.uintptr_t(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func CompressStream(stream Stream, buf []byte) ([]byte, error) {
	var size C.size_t

	ret := C.hs_compress_stream(stream, (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)), &size)

	if ret == C.HS_INSUFFICIENT_SPACE {
		buf = make([]byte, size)

		ret = C.hs_compress_stream(stream, (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)), &size)
	}

	if ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return buf[:size], nil
}

func ExpandStream(db Database, stream *Stream, buf []byte) error {
	ret := C.hs_expand_stream(db, (**C.hs_stream_t)(stream), (*C.char)(unsafe.Pointer(&buf[0])), C.size_t(len(buf)))

	runtime.KeepAlive(buf)

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ResetAndExpandStream(stream Stream, buf []byte, scratch Scratch, onEvent MatchEventHandler, context interface{}) error {
	h := handle.New(MatchEventContext{onEvent, context})
	defer h.Delete()

	ret := C._hs_reset_and_expand_stream(stream,
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.size_t(len(buf)),
		scratch,
		C.match_event_handler(C.hsMatchEventCallback),
		C.uintptr_t(h))

	runtime.KeepAlive(buf)

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}
