package ch

import (
	"unsafe"
)

// #include <ch.h>
import "C"

type Database *C.ch_database_t

func FreeDatabase(db Database) (err error) {
	if ret := C.ch_free_database(db); ret != C.CH_SUCCESS {
		err = Error(ret)
	}

	return
}

func Version() string {
	return C.GoString(C.ch_version())
}

func DatabaseSize(db Database) (n int, err error) {
	var size C.size_t

	if ret := C.ch_database_size(db, &size); ret != C.CH_SUCCESS {
		err = Error(ret)
	} else {
		n = int(size)
	}

	return
}

func DatabaseInfo(db Database) (s string, err error) {
	var info *C.char

	if ret := C.ch_database_info(db, &info); ret != C.HS_SUCCESS {
		err = Error(ret)
	}

	s = C.GoString(info)
	C.free(unsafe.Pointer(info))

	return
}
