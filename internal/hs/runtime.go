package hs

import (
	"errors"
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

// ScanFlag represents a scan flag
type ScanFlag uint

type MatchEventHandler func(id uint, from, to uint64, flags uint, context interface{}) error

type MatchEventContext struct {
	handler MatchEventHandler
	context interface{}
}

//export hsMatchEventCallback
func hsMatchEventCallback(id C.uint, from, to C.ulonglong, flags C.uint, data unsafe.Pointer) C.int {
	h := cgo.Handle(data)
	ctx, ok := h.Value().(MatchEventContext)
	if !ok {
		return C.HS_INVALID
	}

	err := ctx.handler(uint(id), uint64(from), uint64(to), uint(flags), ctx.context)
	if err != nil {
		var hsErr Error
		if errors.As(err, &hsErr) {
			return C.int(hsErr)
		}

		return C.HS_SCAN_TERMINATED
	}

	return C.HS_SUCCESS
}

func Scan(db Database, data []byte, flags ScanFlag, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	if data == nil {
		return Error(C.HS_INVALID)
	}

	var p runtime.Pinner
	defer p.Unpin()

	buf := unsafe.SliceData(data)
	p.Pin(buf)

	h := cgo.NewHandle(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_scan(db,
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

func ScanVector(db Database, data [][]byte, flags ScanFlag, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	if data == nil {
		return Error(C.HS_INVALID)
	}

	var p runtime.Pinner
	defer p.Unpin()

	cdata := make([]uintptr, len(data))
	clength := make([]C.uint, len(data))

	for i, d := range data {
		if d == nil {
			return Error(C.HS_INVALID)
		}

		buf := unsafe.SliceData(d)
		p.Pin(buf)

		cdata[i] = uintptr(unsafe.Pointer(buf))
		clength[i] = C.uint(len(d))
	}

	cdataBuf := unsafe.SliceData(cdata)
	p.Pin(cdataBuf)

	clengthBuf := unsafe.SliceData(clength)
	p.Pin(clengthBuf)

	h := cgo.NewHandle(MatchEventContext{cb, ctx})
	defer h.Delete()

	ret := C.hs_scan_vector(db,
		(**C.char)(unsafe.Pointer(cdataBuf)),
		(*C.uint)(unsafe.Pointer(clengthBuf)),
		C.uint(len(cdata)),
		C.uint(flags),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(h))

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

type MatchEvent struct {
	ID       uint
	From, To uint64
	ScanFlag
}

type MatchRecorder struct {
	Events []MatchEvent
	Err    error
}

func (h *MatchRecorder) Matched() bool { return len(h.Events) > 0 }

func (h *MatchRecorder) Handle(id uint, from, to uint64, flags uint, _ interface{}) error {
	if len(h.Events) > 0 {
		tail := &h.Events[len(h.Events)-1]

		if tail.ID == id && tail.From == from && tail.To < to {
			tail.To = to

			return h.Err
		}
	}

	h.Events = append(h.Events, MatchEvent{id, from, to, ScanFlag(flags)})

	return h.Err
}
