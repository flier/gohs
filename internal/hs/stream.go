package hs

import (
	"reflect"
	"runtime"
	"unsafe"

	"github.com/flier/gohs/internal/handle"
)

/*
#include <stdint.h>

#include <hs.h>

extern int hsMatchEventCallback(unsigned int id,
								unsigned long long from,
								unsigned long long to,
								unsigned int flags,
								void *context);
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

func ScanStream(stream Stream, data []byte, flags ScanFlag, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	if data == nil {
		return Error(C.HS_INVALID)
	}

	h := handle.New(MatchEventContext{cb, ctx})
	defer h.Delete()

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data)) // FIXME: Zero-copy access to go data

	ret := C.hs_scan_stream(stream,
		(*C.char)(unsafe.Pointer(hdr.Data)),
		C.uint(hdr.Len),
		C.uint(flags),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(&h))

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

func CloseStream(stream Stream, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	h := handle.New(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_close_stream(stream,
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(&h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ResetStream(stream Stream, flags ScanFlag, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	h := handle.New(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_reset_stream(stream,
		C.uint(flags),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(&h))

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

func ResetAndCopyStream(to, from Stream, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	h := handle.New(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_reset_and_copy_stream(to,
		from,
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(&h))

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

func ResetAndExpandStream(stream Stream, buf []byte, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	h := handle.New(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_reset_and_expand_stream(stream,
		(*C.char)(unsafe.Pointer(&buf[0])),
		C.size_t(len(buf)),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(&h))

	runtime.KeepAlive(buf)

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}
