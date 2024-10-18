//go:build go1.21
// +build go1.21

package handle

import "C"
import (
	"runtime"
	"unsafe"
)

// GoHandle provides a way to pass values that contain Go pointers (pointers to memory allocated by
// Go) between Go and C without breaking the cgo pointer passing rules. It uses the Pinner feature
// in the runtime instead of the cgo.Handle type in the cgo runtime to avoid lock contention due to
// the use of sync.Map in the cgo implementation.
type GoHandle struct {
	pinner *runtime.Pinner
	value  *any
}

// Value gets the handle value.
func (h *GoHandle) Value() any {
	return *h.value
}

// Set sets the handle value.
func (h *GoHandle) Set(v any) {
	h.Delete()
	pinner := &runtime.Pinner{}
	pinner.Pin(pinner)
	pinner.Pin(&v)
	h.pinner = pinner
	h.value = &v
}

// Delete stops the Go runtime tracking of the referred Go value.
//
// If the Go runtime does not manage the GoHandle storage, it cannot detect when the referred value
// becomes inaccessible. Before deallocating the GoHandle is necessary to close it so to prevent
// memory leaks.
func (h *GoHandle) Delete() {
	if h.pinner != nil {
		h.pinner.Unpin()
		h.pinner = nil
		h.value = nil
	}
}

// NewHandle creates a new handle and pins it to the given value.
func NewHandle(v any) *GoHandle {
	h := &GoHandle{}
	h.Set(v)
	return h
}

// Handle provides a way to pass values that contain Go pointers (pointers to memory allocated by Go)
// between Go and C without breaking the cgo pointer passing rules.
type Handle = GoHandle

// New returns a handle for a given value.
var New = NewHandle

// Pointer returns an unsafe.Pointer to the handle.
func Pointer(h *Handle) unsafe.Pointer {
	return unsafe.Pointer(h)
}
