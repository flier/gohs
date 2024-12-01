package hyperscan

import (
	"fmt"
	"regexp"

	"github.com/flier/gohs/internal/hs"
)

// HsError is the type type for errors returned by Hyperscan functions.
type HsError = Error

// Error is the type type for errors returned by Hyperscan functions.
type Error = hs.Error

const (
	// ErrSuccess is the error returned if the engine completed normally.
	ErrSuccess Error = hs.ErrSuccess
	// ErrInvalid is the error returned if a parameter passed to this function was invalid.
	ErrInvalid Error = hs.ErrInvalid
	// ErrNoMemory is the error returned if a memory allocation failed.
	ErrNoMemory Error = hs.ErrNoMemory
	// ErrScanTerminated is the error returned if the engine was terminated by callback.
	ErrScanTerminated Error = hs.ErrScanTerminated
	// ErrCompileError is the error returned if the pattern compiler failed.
	ErrCompileError Error = hs.ErrCompileError
	// ErrDatabaseVersionError is the error returned if the given database was built for a different version of Hyperscan.
	ErrDatabaseVersionError Error = hs.ErrDatabaseVersionError
	// ErrDatabasePlatformError is the error returned if the given database was built for a different platform.
	ErrDatabasePlatformError Error = hs.ErrDatabasePlatformError
	// ErrDatabaseModeError is the error returned if the given database was built for a different mode of operation.
	ErrDatabaseModeError Error = hs.ErrDatabaseModeError
	// ErrBadAlign is the error returned if a parameter passed to this function was not correctly aligned.
	ErrBadAlign Error = hs.ErrBadAlign
	// ErrBadAlloc is the error returned if the memory allocator did not correctly return memory suitably aligned.
	ErrBadAlloc Error = hs.ErrBadAlloc
	// ErrScratchInUse is the error returned if the scratch region was already in use.
	ErrScratchInUse Error = hs.ErrScratchInUse
	// ErrArchError is the error returned if unsupported CPU architecture.
	ErrArchError Error = hs.ErrArchError
	// ErrInsufficientSpace is the error returned if provided buffer was too small.
	ErrInsufficientSpace Error = hs.ErrInsufficientSpace
)

// Database is an immutable database that can be used by the Hyperscan scanning API.
type Database interface {
	// Provides information about a database.
	Info() (DbInfo, error)

	// Provides the size of the given database in bytes.
	Size() (int, error)

	// Free a compiled pattern database.
	Close() error

	// Serialize a pattern database to a stream of bytes.
	Marshal() ([]byte, error)

	// Reconstruct a pattern database from a stream of bytes at a given memory location.
	Unmarshal(b []byte) error
}

// DbInfo identify the version and platform information for the supplied database.
type DbInfo string //nolint: stylecheck

func (i DbInfo) String() string { return string(i) }

var regexDBInfo = regexp.MustCompile(`^Version: (\d+\.\d+\.\d+) Features: ([\w\s]+)? Mode: (\w+)$`)

const dbInfoMatches = 4

// Parse the version and platform information.
func (i DbInfo) Parse() (version, features, mode string, err error) {
	matched := regexDBInfo.FindStringSubmatch(string(i))

	if len(matched) != dbInfoMatches {
		err = fmt.Errorf("database info `%s`, %w", i, ErrInvalid)
	} else {
		version = matched[1]
		features = matched[2]
		mode = matched[3]
	}

	return
}

// Version is the version for the supplied database.
func (i DbInfo) Version() (string, error) {
	version, _, _, err := i.Parse()

	return version, err
}

// Mode is the scanning mode for the supplied database.
func (i DbInfo) Mode() (ModeFlag, error) {
	_, _, mode, err := i.Parse()
	if err != nil {
		return 0, err
	}

	return ParseModeFlag(mode)
}

// Version identify this release version. The return version is a string
// containing the version number of this release build and the date of the build.
func Version() string { return hs.Version() }

// ValidPlatform test the current system architecture.
func ValidPlatform() error { return hs.ValidPlatform() } //nolint: wrapcheck

type database interface {
	c() hs.Database
}

type baseDatabase struct {
	db hs.Database
}

func newBaseDatabase(db hs.Database) *baseDatabase {
	return &baseDatabase{db}
}

// UnmarshalDatabase reconstruct a pattern database from a stream of bytes.
func UnmarshalDatabase(data []byte) (Database, error) {
	db, err := hs.DeserializeDatabase(data)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return &baseDatabase{db}, nil
}

// UnmarshalBlockDatabase reconstruct a block database from a stream of bytes.
func UnmarshalBlockDatabase(data []byte) (BlockDatabase, error) {
	db, err := hs.DeserializeDatabase(data)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return newBlockDatabase(db), nil
}

// UnmarshalStreamDatabase reconstruct a stream database from a stream of bytes.
func UnmarshalStreamDatabase(data []byte) (StreamDatabase, error) {
	db, err := hs.DeserializeDatabase(data)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return newStreamDatabase(db), nil
}

// UnmarshalVectoredDatabase reconstruct a vectored database from a stream of bytes.
func UnmarshalVectoredDatabase(data []byte) (VectoredDatabase, error) {
	db, err := hs.DeserializeDatabase(data)
	if err != nil {
		return nil, err //nolint: wrapcheck
	}

	return newVectoredDatabase(db), nil
}

// SerializedDatabaseSize reports the size that would be required by a database if it were deserialized.
func SerializedDatabaseSize(data []byte) (int, error) { return hs.SerializedDatabaseSize(data) } //nolint: wrapcheck

// SerializedDatabaseInfo provides information about a serialized database.
func SerializedDatabaseInfo(data []byte) (DbInfo, error) {
	i, err := hs.SerializedDatabaseInfo(data)

	return DbInfo(i), err
}

func (d *baseDatabase) c() hs.Database { return d.db }

func (d *baseDatabase) Size() (int, error) { return hs.DatabaseSize(d.db) } //nolint: wrapcheck

func (d *baseDatabase) Info() (DbInfo, error) {
	i, err := hs.DatabaseInfo(d.db)
	if err != nil {
		return "", err //nolint: wrapcheck
	}

	return DbInfo(i), nil
}

func (d *baseDatabase) Close() error { return hs.FreeDatabase(d.db) } //nolint: wrapcheck

func (d *baseDatabase) Marshal() ([]byte, error) { return hs.SerializeDatabase(d.db) } //nolint: wrapcheck

func (d *baseDatabase) Unmarshal(data []byte) error { return hs.DeserializeDatabaseAt(data, d.db) } //nolint: wrapcheck
