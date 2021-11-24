package ch

import "unsafe"

/*
#include <stdlib.h>
#include <stdint.h>

#include <ch.h>

static inline void* aligned64_malloc(size_t size) {
	void* result;
#ifdef _WIN32
	result = _aligned_malloc(size, 64);
#else
	if (posix_memalign(&result, 64, size)) {
		result = 0;
	}
#endif
	return result;
}

static inline void aligned64_free(void *ptr) {
#ifdef _WIN32
	_aligned_free(ptr);
#else
	free(ptr);
#endif
}

extern void *chDefaultAlloc(size_t size);
extern void chDefaultFree(void *ptr);
extern void *chDbAlloc(size_t size);
extern void chDbFree(void *ptr);
extern void *chMiscAlloc(size_t size);
extern void chMiscFree(void *ptr);
extern void *chScratchAlloc(size_t size);
extern void chScratchFree(void *ptr);
*/
import "C"

type (
	// The type of the callback function that will be used by Chimera to allocate more memory at runtime as required.
	AllocFunc func(uint) unsafe.Pointer
	// The type of the callback function that will be used by Chimera to free memory regions previously
	// allocated using the @ref ch_alloc_t function.
	FreeFunc func(unsafe.Pointer)
)

type Allocator struct {
	Alloc AllocFunc
	Free  FreeFunc
}

func DefaultAllocator() *Allocator  { return &allocator }
func DatabaseAllocator() *Allocator { return &dbAllocator }
func MiscAllocator() *Allocator     { return &miscAllocator }
func ScratchAllocator() *Allocator  { return &scratchAllocator }

var (
	defaultAllocator = Allocator{DefaultAlloc, DefaultFree}
	allocator        = defaultAllocator
	dbAllocator      = defaultAllocator
	miscAllocator    = defaultAllocator
	scratchAllocator = defaultAllocator
)

func DefaultAlloc(size uint) unsafe.Pointer {
	return C.malloc(C.size_t(size))
}

func DefaultFree(ptr unsafe.Pointer) {
	C.free(ptr)
}

func AlignedAlloc(size uint) unsafe.Pointer {
	return C.aligned64_malloc(C.size_t(size))
}

func AlignedFree(ptr unsafe.Pointer) {
	C.aligned64_free(ptr)
}

//export chDefaultAlloc
func chDefaultAlloc(size C.size_t) unsafe.Pointer {
	return allocator.Alloc(uint(size))
}

//export chDefaultFree
func chDefaultFree(ptr unsafe.Pointer) {
	allocator.Free(ptr)
}

// Set the allocate and free functions used by Chimera for allocating
// memory at runtime for stream state, scratch space, database bytecode,
// and various other data structure returned by the Chimera API.
func SetAllocator(allocFunc AllocFunc, freeFunc FreeFunc) (err error) {
	var ret C.ch_error_t

	if allocFunc == nil || freeFunc == nil {
		allocator = defaultAllocator
		ret = C.ch_set_allocator(nil, nil)
	} else {
		allocator = Allocator{allocFunc, freeFunc}
		ret = C.ch_set_allocator(C.ch_alloc_t(C.chDefaultAlloc), C.ch_free_t(C.chDefaultFree))
	}

	if ret != C.CH_SUCCESS {
		err = Error(ret)
	}

	dbAllocator = allocator
	miscAllocator = allocator
	scratchAllocator = allocator

	return
}

func ClearAllocator() error {
	return SetAllocator(nil, nil)
}

//export chDbAlloc
func chDbAlloc(size C.size_t) unsafe.Pointer {
	return dbAllocator.Alloc(uint(size))
}

//export chDbFree
func chDbFree(ptr unsafe.Pointer) {
	dbAllocator.Free(ptr)
}

// Set the allocate and free functions used by Chimera for allocating memory
// for database bytecode produced by the compile calls (@ref ch_compile() and
// @ref ch_compile_multi()).
func SetDatabaseAllocator(allocFunc AllocFunc, freeFunc FreeFunc) (err error) {
	var ret C.ch_error_t

	if allocFunc == nil || freeFunc == nil {
		dbAllocator = defaultAllocator
		ret = C.ch_set_database_allocator(nil, nil)
	} else {
		dbAllocator = Allocator{allocFunc, freeFunc}
		ret = C.ch_set_database_allocator(C.ch_alloc_t(C.chDbAlloc), C.ch_free_t(C.chDbFree))
	}

	if ret != C.CH_SUCCESS {
		err = Error(ret)
	}

	return
}

func ClearDatabaseAllocator() error {
	return SetDatabaseAllocator(nil, nil)
}

//export chMiscAlloc
func chMiscAlloc(size C.size_t) unsafe.Pointer {
	return miscAllocator.Alloc(uint(size))
}

//export chMiscFree
func chMiscFree(ptr unsafe.Pointer) {
	miscAllocator.Free(ptr)
}

// Set the allocate and free functions used by Chimera for allocating memory
// for items returned by the Chimera API such as @ref ch_compile_error_t.
func SetMiscAllocator(allocFunc AllocFunc, freeFunc FreeFunc) (err error) {
	var ret C.ch_error_t

	if allocFunc == nil || freeFunc == nil {
		miscAllocator = defaultAllocator
		ret = C.ch_set_misc_allocator(nil, nil)
	} else {
		miscAllocator = Allocator{allocFunc, freeFunc}
		ret = C.ch_set_misc_allocator(C.ch_alloc_t(C.chMiscAlloc), C.ch_free_t(C.chMiscFree))
	}

	if ret != C.CH_SUCCESS {
		err = Error(ret)
	}

	return
}

func ClearMiscAllocator() error {
	return SetMiscAllocator(nil, nil)
}

//export chScratchAlloc
func chScratchAlloc(size C.size_t) unsafe.Pointer {
	return scratchAllocator.Alloc(uint(size))
}

//export chScratchFree
func chScratchFree(ptr unsafe.Pointer) {
	scratchAllocator.Free(ptr)
}

// Set the allocate and free functions used by Chimera for allocating memory
// for scratch space by @ref ch_alloc_scratch() and @ref ch_clone_scratch().
func SetScratchAllocator(allocFunc AllocFunc, freeFunc FreeFunc) (err error) {
	var ret C.ch_error_t

	if allocFunc == nil || freeFunc == nil {
		scratchAllocator = defaultAllocator
		ret = C.ch_set_scratch_allocator(nil, nil)
	} else {
		scratchAllocator = Allocator{allocFunc, freeFunc}
		ret = C.ch_set_scratch_allocator(C.ch_alloc_t(C.chScratchAlloc), C.ch_free_t(C.chScratchFree))
	}

	if ret != C.CH_SUCCESS {
		err = Error(ret)
	}

	return
}

func ClearScratchAllocator() error {
	return SetScratchAllocator(nil, nil)
}
