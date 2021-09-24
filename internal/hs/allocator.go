package hs

import "unsafe"

/*
#include <stdlib.h>
#include <stdint.h>

#include <hs.h>

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

extern void *hsAlloc(size_t size);
extern void hsFree(void *ptr);
extern void *hsDbAlloc(size_t size);
extern void hsDbFree(void *ptr);
extern void *hsMiscAlloc(size_t size);
extern void hsMiscFree(void *ptr);
extern void *hsScratchAlloc(size_t size);
extern void hsScratchFree(void *ptr);
extern void *hsStreamAlloc(size_t size);
extern void hsStreamFree(void *ptr);
*/
import "C"

type (
	AllocFunc func(uint) unsafe.Pointer
	FreeFunc  func(unsafe.Pointer)
)

type Allocator struct {
	Alloc AllocFunc
	Free  FreeFunc
}

func DefaultAllocator() *Allocator  { return &allocator }
func DatabaseAllocator() *Allocator { return &dbAllocator }
func MiscAllocator() *Allocator     { return &miscAllocator }
func ScratchAllocator() *Allocator  { return &scratchAllocator }
func StreamAllocator() *Allocator   { return &streamAllocator }

var (
	defaultAllocator = Allocator{DefaultAlloc, DefaultFree}
	allocator        = defaultAllocator
	dbAllocator      = defaultAllocator
	miscAllocator    = defaultAllocator
	scratchAllocator = defaultAllocator
	streamAllocator  = defaultAllocator
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

//export hsAlloc
func hsAlloc(size C.size_t) unsafe.Pointer {
	return allocator.Alloc(uint(size))
}

//export hsFree
func hsFree(ptr unsafe.Pointer) {
	allocator.Free(ptr)
}

func SetAllocator(allocFunc AllocFunc, freeFunc FreeFunc) error {
	var ret C.hs_error_t

	if allocFunc == nil || freeFunc == nil {
		allocator = defaultAllocator
		ret = C.hs_set_allocator(nil, nil)
	} else {
		allocator = Allocator{allocFunc, freeFunc}
		ret = C.hs_set_allocator(C.hs_alloc_t(C.hsAlloc), C.hs_free_t(C.hsFree))
	}

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	dbAllocator = allocator
	miscAllocator = allocator
	scratchAllocator = allocator
	streamAllocator = allocator

	return nil
}

func ClearAllocator() error {
	return SetAllocator(nil, nil)
}

//export hsDbAlloc
func hsDbAlloc(size C.size_t) unsafe.Pointer {
	return dbAllocator.Alloc(uint(size))
}

//export hsDbFree
func hsDbFree(ptr unsafe.Pointer) {
	dbAllocator.Free(ptr)
}

func SetDatabaseAllocator(allocFunc AllocFunc, freeFunc FreeFunc) error {
	var ret C.hs_error_t

	if allocFunc == nil || freeFunc == nil {
		dbAllocator = defaultAllocator
		ret = C.hs_set_database_allocator(nil, nil)
	} else {
		dbAllocator = Allocator{allocFunc, freeFunc}
		ret = C.hs_set_database_allocator(C.hs_alloc_t(C.hsDbAlloc), C.hs_free_t(C.hsDbFree))
	}

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ClearDatabaseAllocator() error {
	return SetDatabaseAllocator(nil, nil)
}

//export hsMiscAlloc
func hsMiscAlloc(size C.size_t) unsafe.Pointer {
	return miscAllocator.Alloc(uint(size))
}

//export hsMiscFree
func hsMiscFree(ptr unsafe.Pointer) {
	miscAllocator.Free(ptr)
}

func SetMiscAllocator(allocFunc AllocFunc, freeFunc FreeFunc) error {
	var ret C.hs_error_t

	if allocFunc == nil || freeFunc == nil {
		miscAllocator = defaultAllocator
		ret = C.hs_set_misc_allocator(nil, nil)
	} else {
		miscAllocator = Allocator{allocFunc, freeFunc}
		ret = C.hs_set_misc_allocator(C.hs_alloc_t(C.hsMiscAlloc), C.hs_free_t(C.hsMiscFree))
	}

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ClearMiscAllocator() error {
	return SetMiscAllocator(nil, nil)
}

//export hsScratchAlloc
func hsScratchAlloc(size C.size_t) unsafe.Pointer {
	return scratchAllocator.Alloc(uint(size))
}

//export hsScratchFree
func hsScratchFree(ptr unsafe.Pointer) {
	scratchAllocator.Free(ptr)
}

func SetScratchAllocator(allocFunc AllocFunc, freeFunc FreeFunc) error {
	var ret C.hs_error_t

	if allocFunc == nil || freeFunc == nil {
		scratchAllocator = defaultAllocator
		ret = C.hs_set_scratch_allocator(nil, nil)
	} else {
		scratchAllocator = Allocator{allocFunc, freeFunc}
		ret = C.hs_set_scratch_allocator(C.hs_alloc_t(C.hsScratchAlloc), C.hs_free_t(C.hsScratchFree))
	}

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ClearScratchAllocator() error {
	return SetScratchAllocator(nil, nil)
}

//export hsStreamAlloc
func hsStreamAlloc(size C.size_t) unsafe.Pointer {
	return streamAllocator.Alloc(uint(size))
}

//export hsStreamFree
func hsStreamFree(ptr unsafe.Pointer) {
	streamAllocator.Free(ptr)
}

func SetStreamAllocator(allocFunc AllocFunc, freeFunc FreeFunc) error {
	var ret C.hs_error_t

	if allocFunc == nil || freeFunc == nil {
		streamAllocator = defaultAllocator
		ret = C.hs_set_stream_allocator(nil, nil)
	} else {
		streamAllocator = Allocator{allocFunc, freeFunc}
		ret = C.hs_set_stream_allocator(C.hs_alloc_t(C.hsStreamAlloc), C.hs_free_t(C.hsStreamFree))
	}

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func ClearStreamAllocator() error {
	return SetStreamAllocator(nil, nil)
}
