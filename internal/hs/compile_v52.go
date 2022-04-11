//go:build hyperscan_v52 || hyperscan_v54
// +build hyperscan_v52 hyperscan_v54

package hs

import (
	"fmt"
	"unsafe"
)

/*
#include <hs.h>
*/
import "C"

// Pure literal is a special case of regular expression.
// A character sequence is regarded as a pure literal if and
// only if each character is read and interpreted independently.
// No syntax association happens between any adjacent characters.
type Literal struct {
	Expr  string      // The expression to parse.
	Flags CompileFlag // Flags which modify the behaviour of the expression.
	ID    int         // The ID number to be associated with the corresponding pattern
	*ExprInfo
}

func CompileLit(expression string, flags CompileFlag, mode ModeFlag, info *PlatformInfo) (Database, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = (*C.struct_hs_platform_info)(unsafe.Pointer(&info.Platform))
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
		return nil, &CompileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}

type Literals interface {
	Literals() []*Literal
}

func CompileLitMulti(input Literals, mode ModeFlag, info *PlatformInfo) (Database, error) {
	var db *C.hs_database_t
	var err *C.hs_compile_error_t
	var platform *C.hs_platform_info_t

	if info != nil {
		platform = (*C.struct_hs_platform_info)(unsafe.Pointer(&info.Platform))
	}

	literals := input.Literals()
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
		exprs[i] = C.CString(lit.Expr)
		lens[i] = C.size_t(len(lit.Expr))
		flags[i] = C.uint(lit.Flags)
		ids[i] = C.uint(lit.ID)
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
		return nil, &CompileError{C.GoString(err.message), int(err.expression)}
	}

	return nil, fmt.Errorf("compile error %d, %w", int(ret), ErrCompileError)
}
