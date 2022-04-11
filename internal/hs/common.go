package hs

import (
	"runtime"
	"unsafe"
)

// #include <hs.h>
import "C"

type Database *C.hs_database_t

func Version() string {
	return C.GoString(C.hs_version())
}

func ValidPlatform() error {
	if ret := C.hs_valid_platform(); ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func FreeDatabase(db Database) (err error) {
	if ret := C.hs_free_database(db); ret != C.HS_SUCCESS {
		err = Error(ret)
	}

	return
}

func SerializeDatabase(db Database) (b []byte, err error) {
	var data *C.char
	var length C.size_t

	ret := C.hs_serialize_database(db, &data, &length)
	if ret != C.HS_SUCCESS {
		err = Error(ret)
	} else {
		defer C.free(unsafe.Pointer(data))

		b = C.GoBytes(unsafe.Pointer(data), C.int(length))
	}

	return
}

func DeserializeDatabase(data []byte) (Database, error) {
	var db *C.hs_database_t

	ret := C.hs_deserialize_database((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &db)

	runtime.KeepAlive(data)

	if ret != C.HS_SUCCESS {
		return nil, Error(ret)
	}

	return db, nil
}

func DeserializeDatabaseAt(data []byte, db Database) error {
	ret := C.hs_deserialize_database_at((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), db)

	runtime.KeepAlive(data)

	if ret != C.HS_SUCCESS {
		return Error(ret)
	}

	return nil
}

func StreamSize(db Database) (int, error) {
	var size C.size_t

	if ret := C.hs_stream_size(db, &size); ret != C.HS_SUCCESS {
		return 0, Error(ret)
	}

	return int(size), nil
}

func DatabaseSize(db Database) (int, error) {
	var size C.size_t

	if ret := C.hs_database_size(db, &size); ret != C.HS_SUCCESS {
		return -1, Error(ret)
	}

	return int(size), nil
}

func SerializedDatabaseSize(data []byte) (int, error) {
	var size C.size_t

	ret := C.hs_serialized_database_size((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &size)

	runtime.KeepAlive(data)

	if ret != C.HS_SUCCESS {
		return 0, Error(ret)
	}

	return int(size), nil
}

func DatabaseInfo(db Database) (string, error) {
	var info *C.char

	if ret := C.hs_database_info(db, &info); ret != C.HS_SUCCESS {
		return "", Error(ret)
	}

	defer C.free(unsafe.Pointer(info))

	return C.GoString(info), nil
}

func SerializedDatabaseInfo(data []byte) (string, error) {
	var info *C.char

	ret := C.hs_serialized_database_info((*C.char)(unsafe.Pointer(&data[0])), C.size_t(len(data)), &info)

	runtime.KeepAlive(data)

	if ret != C.HS_SUCCESS {
		return "", Error(ret)
	}

	defer C.free(unsafe.Pointer(info))

	return C.GoString(info), nil
}
