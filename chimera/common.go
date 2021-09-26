package chimera

import (
	"github.com/flier/gohs/internal/ch"
)

// Error is the type for errors returned by Chimera functions.
type Error = ch.Error

const (
	// ErrSuccess is the error returned if the engine completed normally.
	ErrSuccess Error = ch.ErrSuccess
	// ErrInvalid is the error returned if a parameter passed to this function was invalid.
	ErrInvalid Error = ch.ErrInvalid
	// ErrNoMemory is the error returned if a memory allocation failed.
	ErrNoMemory Error = ch.ErrNoMemory
	// ErrScanTerminated is the error returned if the engine was terminated by callback.
	ErrScanTerminated Error = ch.ErrScanTerminated
	// ErrCompileError is the error returned if the pattern compiler failed.
	ErrCompileError Error = ch.ErrCompileError
	// ErrDatabaseVersionError is the error returned if the given database was built
	// for a different version of the Chimera matcher.
	ErrDatabaseVersionError Error = ch.ErrDatabaseVersionError
	// ErrDatabasePlatformError is the error returned if the given database was built for a different platform.
	ErrDatabasePlatformError Error = ch.ErrDatabasePlatformError
	// ErrDatabaseModeError is the error returned if the given database was built for a different mode of operation.
	ErrDatabaseModeError Error = ch.ErrDatabaseModeError
	// ErrBadAlign is the error returned if a parameter passed to this function was not correctly aligned.
	ErrBadAlign Error = ch.ErrBadAlign
	// ErrBadAlloc is the error returned if the memory allocator did not correctly return memory suitably aligned.
	ErrBadAlloc Error = ch.ErrBadAlloc
	// ErrScratchInUse is the error returned if the scratch region was already in use.
	ErrScratchInUse Error = ch.ErrScratchInUse
)

// DbInfo identify the version and platform information for the supplied database.
type DbInfo string // nolint: stylecheck

func (i DbInfo) String() string { return string(i) }

// Database is an immutable database that can be used by the Hyperscan scanning API.
type Database interface {
	// Provides information about a database.
	Info() (DbInfo, error)

	// Provides the size of the given database in bytes.
	Size() (int, error)

	// Free a compiled pattern database.
	Close() error
}

type database struct {
	db ch.Database
}

func (d *database) Size() (int, error) { return ch.DatabaseSize(d.db) } // nolint: wrapcheck

func (d *database) Info() (DbInfo, error) {
	i, err := ch.DatabaseInfo(d.db)
	if err != nil {
		return "", err //nolint: wrapcheck
	}

	return DbInfo(i), nil
}

func (d *database) Close() error { return ch.FreeDatabase(d.db) } // nolint: wrapcheck

// Version identify this release version.
//
// The return version is a string containing the version number of this release
// build and the date of the build.
func Version() string { return ch.Version() }
