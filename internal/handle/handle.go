// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package handle

import (
	"runtime"
	"strconv"
	"sync"
	"sync/atomic"
)

// Handle provides a way to pass values that contain Go pointers
// (pointers to memory allocated by Go) between Go and C without
// breaking the cgo pointer passing rules. A Handle is an integer
// value that can represent any Go value. A Handle can be passed
// through C and back to Go, and Go code can use the Handle to
// retrieve the original Go value.
//
// The underlying type of Handle is guaranteed to fit in an integer type
// that is large enough to hold the bit pattern of any pointer. The zero
// value of a Handle is not valid, and thus is safe to use as a sentinel
// in C APIs.
//
// For instance, on the Go side:
//
//	package main
//
//	/*
//	#include <stdint.h> // for uintptr_t
//
//	extern void MyGoPrint(uintptr_t handle);
//	void myprint(uintptr_t handle);
//	*/
//	import "C"
//	import "runtime/cgo"
//
//	//export MyGoPrint
//	func MyGoPrint(handle C.uintptr_t) {
//		h := cgo.Handle(handle)
//		val := h.Value().(string)
//		println(val)
//		h.Delete()
//	}
//
//	func main() {
//		val := "hello Go"
//		C.myprint(C.uintptr_t(cgo.NewHandle(val)))
//		// Output: hello Go
//	}
//
// and on the C side:
//
//	#include <stdint.h> // for uintptr_t
//
//	// A Go function
//	extern void MyGoPrint(uintptr_t handle);
//
//	// A C function
//	void myprint(uintptr_t handle) {
//	    MyGoPrint(handle);
//	}
type Handle uintptr

// New returns a handle for a given value.
//
// The handle is valid until the program calls Delete on it. The handle
// uses resources, and this package assumes that C code may hold on to
// the handle, so a program must explicitly call Delete when the handle
// is no longer needed.
//
// The intended use is to pass the returned handle to C code, which
// passes it back to Go, which calls Value.
func New(v interface{}) Handle {
	h := atomic.AddUintptr(&handleIdx, 1)
	if h == 0 {
		panic("runtime/cgo: ran out of handle space")
	}

	handles(Handle(h)).Store(h, v)

	return Handle(h)
}

// Value returns the associated Go value for a valid handle.
//
// The method panics if the handle is invalid.
func (h Handle) Value() interface{} {
	v, ok := handles(h).Load(uintptr(h))
	if !ok {
		panic("runtime/cgo: misuse of an invalid Handle")
	}
	return v
}

func (h Handle) Delete() {
	_, ok := handles(h).LoadAndDelete(uintptr(h))
	if !ok {
		panic("runtime/cgo: misuse of an invalid Handle")
	}
}

func handles(h Handle) *sync.Map {
	return &handleMap[uintptr(h)%uintptr(len(handleMap))]
}

var (
	mapSize   = ""             // --ldflags="-X 'github.com/flier/gohs/internal/handle.mapSize=64'"
	handleMap = newHandleMap() // []map[Handle]interface{}
	handleIdx uintptr          // atomic
)

func newHandleMap() []sync.Map {
	n, err := strconv.Atoi(mapSize)
	if err != nil {
		n = runtime.NumCPU() * 2
	}

	return make([]sync.Map, n)
}
