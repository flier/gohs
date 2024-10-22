package hs

import (
	"runtime"
	"runtime/cgo"
	"unsafe"
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

	var p runtime.Pinner
	defer p.Unpin()

	buf := unsafe.SliceData(data)
	p.Pin(buf)

	h := cgo.NewHandle(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_scan_stream(stream,
		(*C.char)(unsafe.Pointer(buf)),
		C.uint(len(data)),
		C.uint(flags),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func FreeStream(stream Stream) {
	C.hs_close_stream(stream, nil, nil, nil)
}

func CloseStream(stream Stream, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	h := cgo.NewHandle(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_close_stream(stream,
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ResetStream(stream Stream, flags ScanFlag, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	h := cgo.NewHandle(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_reset_stream(stream,
		C.uint(flags),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(h))

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
	h := cgo.NewHandle(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_reset_and_copy_stream(to,
		from,
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func CompressStream(stream Stream, b []byte) ([]byte, error) {
	var p runtime.Pinner
	defer p.Unpin()

	buf := unsafe.SliceData(b)
	p.Pin(buf)

	var size C.size_t

	ret := C.hs_compress_stream(stream, (*C.char)(unsafe.Pointer(buf)), C.size_t(len(b)), &size)
	if ret == C.HS_INSUFFICIENT_SPACE {
		b = make([]byte, size)

		buf = unsafe.SliceData(b)
		p.Pin(buf)

		ret = C.hs_compress_stream(stream, (*C.char)(unsafe.Pointer(buf)), C.size_t(len(b)), &size)
	}

	if ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return b[:size], nil
}

func ExpandStream(db Database, stream *Stream, b []byte) error {
	var p runtime.Pinner
	defer p.Unpin()

	buf := unsafe.SliceData(b)
	p.Pin(buf)

	ret := C.hs_expand_stream(db, (**C.hs_stream_t)(stream), (*C.char)(unsafe.Pointer(buf)), C.size_t(len(b)))
	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ResetAndExpandStream(stream Stream, b []byte, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	var p runtime.Pinner
	defer p.Unpin()

	buf := unsafe.SliceData(b)
	p.Pin(buf)

	h := cgo.NewHandle(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_reset_and_expand_stream(stream,
		(*C.char)(unsafe.Pointer(buf)),
		C.size_t(len(b)),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}
