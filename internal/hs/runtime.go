package hs

import (
	"errors"
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

// ScanFlag represents a scan flag
type ScanFlag uint

type MatchEventHandler func(id uint, from, to uint64, flags uint, context interface{}) error

type MatchEventContext struct {
	handler MatchEventHandler
	context interface{}
}

//export hsMatchEventCallback
func hsMatchEventCallback(id C.uint, from, to C.ulonglong, flags C.uint, data unsafe.Pointer) C.int {
	h := (*handle.Handle)(data)
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

	h := handle.New(MatchEventContext{cb, ctx})
	defer h.Delete()

	hdr := (*reflect.SliceHeader)(unsafe.Pointer(&data)) // FIXME: Zero-copy access to go data

	ret := C.hs_scan(db,
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

func ScanVector(db Database, data [][]byte, flags ScanFlag, s Scratch, cb MatchEventHandler, ctx interface{}) error {
	if data == nil {
		return Error(C.HS_INVALID)
	}

	cdata := make([]uintptr, len(data))
	clength := make([]C.uint, len(data))

	for i, d := range data {
		if d == nil {
			return Error(C.HS_INVALID)
		}

		// FIXME: Zero-copy access to go data
		hdr := (*reflect.SliceHeader)(unsafe.Pointer(&d)) //nolint: scopelint
		cdata[i] = uintptr(unsafe.Pointer(hdr.Data))
		clength[i] = C.uint(hdr.Len)
	}

	h := handle.New(MatchEventContext{cb, ctx})
	defer h.Delete()

	cdataHdr := (*reflect.SliceHeader)(unsafe.Pointer(&cdata))     // FIXME: Zero-copy access to go data
	clengthHdr := (*reflect.SliceHeader)(unsafe.Pointer(&clength)) // FIXME: Zero-copy access to go data

	ret := C.hs_scan_vector(db,
		(**C.char)(unsafe.Pointer(cdataHdr.Data)),
		(*C.uint)(unsafe.Pointer(clengthHdr.Data)),
		C.uint(cdataHdr.Len),
		C.uint(flags),
		s,
		C.match_event_handler(C.hsMatchEventCallback),
		unsafe.Pointer(&h))

	// Ensure go data is alive before the C function returns
	runtime.KeepAlive(data)
	runtime.KeepAlive(cdata)
	runtime.KeepAlive(clength)

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
