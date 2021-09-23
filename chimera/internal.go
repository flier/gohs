package chimera

/*
#cgo pkg-config: libch libhs
#cgo linux LDFLAGS: -lm -lpcre
#cgo darwin LDFLAGS: -lpcre

#include <stdlib.h>
#include <limits.h>
#include <stdint.h>

#include <ch.h>
*/
import "C"
import "fmt"

// ChError represents an error
type ChError C.ch_error_t

const (
	// ErrSuccess is the error returned if the engine completed normally.
	ErrSuccess ChError = C.CH_SUCCESS
	// ErrInvalid is the error returned if a parameter passed to this function was invalid.
	ErrInvalid ChError = C.CH_INVALID
	// ErrNoMemory is the error returned if a memory allocation failed.
	ErrNoMemory ChError = C.CH_NOMEM
	// ErrScanTerminated is the error returned if the engine was terminated by callback.
	ErrScanTerminated ChError = C.CH_SCAN_TERMINATED
	// ErrCompileError is the error returned if the pattern compiler failed.
	ErrCompileError ChError = C.CH_COMPILER_ERROR
	// ErrDatabaseVersionError is the error returned if the given database was built for a different version of Hyperscan.
	ErrDatabaseVersionError ChError = C.CH_DB_VERSION_ERROR
	// ErrDatabasePlatformError is the error returned if the given database was built for a different platform (i.e., CPU type).
	ErrDatabasePlatformError ChError = C.CH_DB_PLATFORM_ERROR
	// ErrDatabaseModeError is the error returned if the given database was built for a different mode of operation.
	ErrDatabaseModeError ChError = C.CH_DB_MODE_ERROR
	// ErrBadAlign is the error returned if a parameter passed to this function was not correctly aligned.
	ErrBadAlign ChError = C.CH_BAD_ALIGN
	// ErrBadAlloc is the error returned if the memory allocator did not correctly return memory suitably aligned.
	ErrBadAlloc ChError = C.CH_BAD_ALLOC
	// ErrScratchInUse is the error returned if the scratch region was already in use.
	ErrScratchInUse ChError = C.CH_SCRATCH_IN_USE
	// ErrUnknown is the unexpected internal error from Hyperscan.
	ErrUnknownHSError ChError = C.CH_UNKNOWN_HS_ERROR
)

var chErrorMessages = map[ChError]string{
	C.CH_SUCCESS:           "The engine completed normally.",
	C.CH_INVALID:           "A parameter passed to this function was invalid.",
	C.CH_NOMEM:             "A memory allocation failed.",
	C.CH_SCAN_TERMINATED:   "The engine was terminated by callback.",
	C.CH_COMPILER_ERROR:    "The pattern compiler failed.",
	C.CH_DB_VERSION_ERROR:  "The given database was built for a different version of Hyperscan.",
	C.CH_DB_PLATFORM_ERROR: "The given database was built for a different platform (i.e., CPU type).",
	C.CH_DB_MODE_ERROR:     "The given database was built for a different mode of operation.",
	C.CH_BAD_ALIGN:         "A parameter passed to this function was not correctly aligned.",
	C.CH_BAD_ALLOC:         "The memory allocator did not correctly return aligned memory.",
	C.CH_SCRATCH_IN_USE:    "The scratch region was already in use.",
	C.CH_UNKNOWN_HS_ERROR:  "Unexpected internal error from Hyperscan.",
}

func (e ChError) Error() string {
	if msg, exists := chErrorMessages[e]; exists {
		return msg
	}

	return fmt.Sprintf("unexpected error, %d", int(e))
}

type (
	chDatabase *C.ch_database_t
	chScratch  *C.ch_scratch_t
)

func chVersion() string {
	return C.GoString(C.ch_version())
}
