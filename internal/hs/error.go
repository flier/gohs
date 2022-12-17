package hs

// #include <hs.h>
import "C"
import "fmt"

// Error represents an error
type Error int

const (
	// ErrSuccess is the error returned if the engine completed normally.
	ErrSuccess Error = C.HS_SUCCESS
	// ErrInvalid is the error returned if a parameter passed to this function was invalid.
	ErrInvalid Error = C.HS_INVALID
	// ErrNoMemory is the error returned if a memory allocation failed.
	ErrNoMemory Error = C.HS_NOMEM
	// ErrScanTerminated is the error returned if the engine was terminated by callback.
	ErrScanTerminated Error = C.HS_SCAN_TERMINATED
	// ErrCompileError is the error returned if the pattern compiler failed.
	ErrCompileError Error = C.HS_COMPILER_ERROR
	// ErrDatabaseVersionError is the error returned if the given database was built for a different version of Hyperscan.
	ErrDatabaseVersionError Error = C.HS_DB_VERSION_ERROR
	// ErrDatabasePlatformError is the error returned if the given database was built for a different platform.
	ErrDatabasePlatformError Error = C.HS_DB_PLATFORM_ERROR
	// ErrDatabaseModeError is the error returned if the given database was built for a different mode of operation.
	ErrDatabaseModeError Error = C.HS_DB_MODE_ERROR
	// ErrBadAlign is the error returned if a parameter passed to this function was not correctly aligned.
	ErrBadAlign Error = C.HS_BAD_ALIGN
	// ErrBadAlloc is the error returned if the memory allocator did not correctly return memory suitably aligned.
	ErrBadAlloc Error = C.HS_BAD_ALLOC
	// ErrScratchInUse is the error returned if the scratch region was already in use.
	ErrScratchInUse Error = C.HS_SCRATCH_IN_USE
	// ErrArchError is the error returned if unsupported CPU architecture.
	ErrArchError Error = C.HS_ARCH_ERROR
	// ErrInsufficientSpace is the error returned if provided buffer was too small.
	ErrInsufficientSpace Error = C.HS_INSUFFICIENT_SPACE
)

var errorMessages = map[Error]string{
	C.HS_SUCCESS:            "The engine completed normally.",
	C.HS_INVALID:            "A parameter passed to this function was invalid.",
	C.HS_NOMEM:              "A memory allocation failed.",
	C.HS_SCAN_TERMINATED:    "The engine was terminated by callback.",
	C.HS_COMPILER_ERROR:     "The pattern compiler failed.",
	C.HS_DB_VERSION_ERROR:   "The given database was built for a different version of Hyperscan.",
	C.HS_DB_PLATFORM_ERROR:  "The given database was built for a different platform (i.e., CPU type).",
	C.HS_DB_MODE_ERROR:      "The given database was built for a different mode of operation.",
	C.HS_BAD_ALIGN:          "A parameter passed to this function was not correctly aligned.",
	C.HS_BAD_ALLOC:          "The memory allocator did not correctly return aligned memory.",
	C.HS_SCRATCH_IN_USE:     "The scratch region was already in use.",
	C.HS_ARCH_ERROR:         "Unsupported CPU architecture.",
	C.HS_INSUFFICIENT_SPACE: "Provided buffer was too small.",
}

func (e Error) Error() string {
	if msg, exists := errorMessages[e]; exists {
		return msg
	}

	return fmt.Sprintf("unexpected error, %d", int(e))
}
