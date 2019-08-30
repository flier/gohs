// +build !hyperscan_v4,!hyperscan_v5_1

package hyperscan

import (
	"fmt"
	"unsafe"
	"reflect"
)

/*
#include <hs.h>
*/
import "C"

func hsCompileLit(expression []byte, flags CompileFlag, mode ModeFlag, info *hsPlatformInfo) (hsDatabase, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = &info.platform
	}

	expr := C.CBytes(expression)

	ret := C.hs_compile_lit((*C.char)(expr), C.uint(flags), C.ulong(len(expression)), C.uint(mode), platform, &db, &err)

	C.free(expr)

	if ret == C.HS_SUCCESS {
		return db, nil
	}

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &compileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error, %d", int(ret))
}

func hsCompileLitMulti(literals []*Literal, mode ModeFlag, info *hsPlatformInfo) (hsDatabase, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = &info.platform
	}

	cexprs := (**C.char)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(uintptr(0)))))
	exprs := *(*[]*C.char)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cexprs)),
		Len:  len(literals),
		Cap:  len(literals),
	}))

	cflags := (*C.uint)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(C.uint(0)))))
	flags := *(*[]C.uint)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cflags)),
		Len:  len(literals),
		Cap:  len(literals),
	}))

	cids := (*C.uint)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(C.uint(0)))))
	ids := *(*[]C.uint)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(cids)),
		Len:  len(literals),
		Cap:  len(literals),
	}))

	clens := (*C.size_t)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(C.size_t(0)))))
	lens := *(*[]C.size_t)(unsafe.Pointer(&reflect.SliceHeader{
		Data: uintptr(unsafe.Pointer(clens)),
		Len:  len(literals),
		Cap:  len(literals),
	}))

	for i, literal := range literals {
		exprs[i] = (*C.char)(C.CBytes(literal.Expression))
		flags[i] = C.uint(literal.Flags)
		ids[i] = C.uint(literal.Id)
		lens[i] = C.size_t(len(literal.Expression))
	}

	ret := C.hs_compile_lit_multi(cexprs, cflags, cids, clens, C.uint(len(literals)), C.uint(mode), platform, &db, &err)

	for _, expr := range exprs {
		C.free(unsafe.Pointer(expr))
	}

	C.free(unsafe.Pointer(cexprs))
	C.free(unsafe.Pointer(cflags))
	C.free(unsafe.Pointer(cids))

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_SUCCESS {
		return db, nil
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &compileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error, %d", int(ret))
}
