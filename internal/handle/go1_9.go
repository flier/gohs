//go:build go1.9 && !go1.15
// +build go1.9,!go1.15

package handle

import "unsafe"

// Delete invalidates a handle. This method should only be called once
// the program no longer needs to pass the handle to C and the C code
// no longer has a copy of the handle value.
//
// The method panics if the handle is invalid.
func (h Handle) Delete() {
	_, ok := handles.Load(uintptr(h))
	if ok {
		handles.Delete(uintptr(h))
	} else {
		panic("runtime/cgo: misuse of an invalid Handle")
	}
}

// Pointer returns an unsafe.Pointer to the handle.
func Pointer(h Handle) unsafe.Pointer {
	return unsafe.Pointer(&h)
}
