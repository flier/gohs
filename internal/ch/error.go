package ch

// #include <ch.h>
import "C"

import "fmt"

// Error represents an error.
type Error C.ch_error_t

const (
	// ErrSuccess is the error returned if the engine completed normally.
	ErrSuccess Error = C.CH_SUCCESS
	// ErrInvalid is the error returned if a parameter passed to this function was invalid.
	ErrInvalid Error = C.CH_INVALID
	// ErrNoMemory is the error returned if a memory allocation failed.
	ErrNoMemory Error = C.CH_NOMEM
	// ErrScanTerminated is the error returned if the engine was terminated by callback.
	ErrScanTerminated Error = C.CH_SCAN_TERMINATED
	// ErrCompileError is the error returned if the pattern compiler failed.
	ErrCompileError Error = C.CH_COMPILER_ERROR
	// ErrDatabaseVersionError is the error returned if the given database was built for a different version of Hyperscan.
	ErrDatabaseVersionError Error = C.CH_DB_VERSION_ERROR
	// ErrDatabasePlatformError is the error returned if the given database was built for a different platform (i.e., CPU type).
	ErrDatabasePlatformError Error = C.CH_DB_PLATFORM_ERROR
	// ErrDatabaseModeError is the error returned if the given database was built for a different mode of operation.
	ErrDatabaseModeError Error = C.CH_DB_MODE_ERROR
	// ErrBadAlign is the error returned if a parameter passed to this function was not correctly aligned.
	ErrBadAlign Error = C.CH_BAD_ALIGN
	// ErrBadAlloc is the error returned if the memory allocator did not correctly return memory suitably aligned.
	ErrBadAlloc Error = C.CH_BAD_ALLOC
	// ErrScratchInUse is the error returned if the scratch region was already in use.
	ErrScratchInUse Error = C.CH_SCRATCH_IN_USE
)

var ErrorMessages = map[Error]string{
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
}

func (e Error) Error() string {
	if msg, exists := ErrorMessages[e]; exists {
		return msg
	}

	return fmt.Sprintf("unexpected error, %d", int(C.ch_error_t(e)))
}
