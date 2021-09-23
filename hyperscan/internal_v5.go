//go:build !hyperscan_v4
// +build !hyperscan_v4

package hyperscan

/*
#include <hs.h>
*/
import "C"

import (
	"fmt"
	"unsafe"
)

const (
	// Combination represents logical combination.
	Combination CompileFlag = C.HS_FLAG_COMBINATION
	// Quiet represents don't do any match reporting.
	Quiet CompileFlag = C.HS_FLAG_QUIET
)

func init() {
	compileFlags['C'] = Combination
	compileFlags['Q'] = Quiet
}

func hsCompileLit(expression string, flags CompileFlag, mode ModeFlag, info *hsPlatformInfo) (hsDatabase, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = &info.platform
	}

	expr := C.CString(expression)

	defer C.free(unsafe.Pointer(expr))

	ret := C.hs_compile_lit(expr, C.uint(flags), C.ulong(len(expression)), C.uint(mode), platform, &db, &err)

	if err != nil {
		defer C.hs_free_compile_error(err)
	}

	if ret == C.HS_SUCCESS {
		return db, nil
	}

	if ret == C.HS_COMPILER_ERROR && err != nil {
		return nil, &compileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}

func hsCompileLitMulti(literals []*Literal, mode ModeFlag, info *hsPlatformInfo) (hsDatabase, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = &info.platform
	}

	count := len(literals)

	cexprs := (**C.char)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(uintptr(0)))))
	exprs := (*[1 << 30]*C.char)(unsafe.Pointer(cexprs))[:len(literals):len(literals)]

	clens := (*C.size_t)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(uintptr(0)))))
	lens := (*[1 << 30]C.size_t)(unsafe.Pointer(clens))[:len(literals):len(literals)]

	cflags := (*C.uint)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(C.uint(0)))))
	flags := (*[1 << 30]C.uint)(unsafe.Pointer(cflags))[:len(literals):len(literals)]

	cids := (*C.uint)(C.calloc(C.size_t(len(literals)), C.size_t(unsafe.Sizeof(C.uint(0)))))
	ids := (*[1 << 30]C.uint)(unsafe.Pointer(cids))[:len(literals):len(literals)]

	for i, lit := range literals {
		exprs[i] = C.CString(string(lit.Expression))
		lens[i] = C.size_t(len(lit.Expression))
		flags[i] = C.uint(lit.Flags)
		ids[i] = C.uint(lit.Id)
	}

	ret := C.hs_compile_lit_multi(cexprs, cflags, cids, clens, C.uint(count), C.uint(mode), platform, &db, &err)

	for _, expr := range exprs {
		C.free(unsafe.Pointer(expr))
	}

	C.free(unsafe.Pointer(cexprs))
	C.free(unsafe.Pointer(clens))
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

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}
