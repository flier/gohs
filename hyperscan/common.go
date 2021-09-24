package hyperscan

import (
	"fmt"
	"regexp"

	"github.com/flier/gohs/internal/hs"
)

type HsError = hs.Error

const (
	// ErrSuccess is the error returned if the engine completed normally.
	ErrSuccess HsError = hs.ErrSuccess
	// ErrInvalid is the error returned if a parameter passed to this function was invalid.
	ErrInvalid HsError = hs.ErrInvalid
	// ErrNoMemory is the error returned if a memory allocation failed.
	ErrNoMemory HsError = hs.ErrNoMemory
	// ErrScanTerminated is the error returned if the engine was terminated by callback.
	ErrScanTerminated HsError = hs.ErrScanTerminated
	// ErrCompileError is the error returned if the pattern compiler failed.
	ErrCompileError HsError = hs.ErrCompileError
	// ErrDatabaseVersionError is the error returned if the given database was built for a different version of Hyperscan.
	ErrDatabaseVersionError HsError = hs.ErrDatabaseVersionError
	// ErrDatabasePlatformError is the error returned if the given database was built for a different platform.
	ErrDatabasePlatformError HsError = hs.ErrDatabasePlatformError
	// ErrDatabaseModeError is the error returned if the given database was built for a different mode of operation.
	ErrDatabaseModeError HsError = hs.ErrDatabaseModeError
	// ErrBadAlign is the error returned if a parameter passed to this function was not correctly aligned.
	ErrBadAlign HsError = hs.ErrBadAlign
	// ErrBadAlloc is the error returned if the memory allocator did not correctly return memory suitably aligned.
	ErrBadAlloc HsError = hs.ErrBadAlloc
	// ErrScratchInUse is the error returned if the scratch region was already in use.
	ErrScratchInUse HsError = hs.ErrScratchInUse
	// ErrArchError is the error returned if unsupported CPU architecture.
	ErrArchError HsError = hs.ErrArchError
	// ErrInsufficientSpace is the error returned if provided buffer was too small.
	ErrInsufficientSpace HsError = hs.ErrInsufficientSpace
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
	Unmarshal([]byte) error
}

// BlockDatabase scan the target data that is a discrete,
// contiguous block which can be scanned in one call and does not require state to be retained.
type BlockDatabase interface {
	Database
	BlockScanner
	BlockMatcher
}

// StreamDatabase scan the target data to be scanned is a continuous stream,
// not all of which is available at once;
// blocks of data are scanned in sequence and matches may span multiple blocks in a stream.
type StreamDatabase interface {
	Database
	StreamScanner
	StreamMatcher
	StreamCompressor

	StreamSize() (int, error)
}

// VectoredDatabase scan the target data that consists of a list of non-contiguous blocks
// that are available all at once.
type VectoredDatabase interface {
	Database
	VectoredScanner
	VectoredMatcher
}

const infoMatches = 4

var regexInfo = regexp.MustCompile(`^Version: (\d+\.\d+\.\d+) Features: ([\w\s]+)? Mode: (\w+)$`)

// DbInfo identify the version and platform information for the supplied database.
type DbInfo string // nolint: stylecheck

func (i DbInfo) String() string { return string(i) }

// Version is the version for the supplied database.
func (i DbInfo) Version() (string, error) {
	matched := regexInfo.FindStringSubmatch(string(i))

	if len(matched) != infoMatches {
		return "", fmt.Errorf("database info, %w", ErrInvalid)
	}

	return matched[1], nil
}

// Mode is the scanning mode for the supplied database.
func (i DbInfo) Mode() (ModeFlag, error) {
	matched := regexInfo.FindStringSubmatch(string(i))

	if len(matched) != infoMatches {
		return 0, fmt.Errorf("database info, %w", ErrInvalid)
	}

	return ParseModeFlag(matched[3])
}

// Version identify this release version. The return version is a string
// containing the version number of this release build and the date of the build.
func Version() string { return hs.Version() }

// ValidPlatform test the current system architecture.
func ValidPlatform() error { return hs.ValidPlatform() } // nolint: wrapcheck

type database interface {
	Db() hs.Database
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
		return nil, err // nolint: wrapcheck
	}

	return &baseDatabase{db}, nil
}

// UnmarshalBlockDatabase reconstruct a block database from a stream of bytes.
func UnmarshalBlockDatabase(data []byte) (BlockDatabase, error) {
	db, err := hs.DeserializeDatabase(data)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	return newBlockDatabase(db), nil
}

// UnmarshalStreamDatabase reconstruct a stream database from a stream of bytes.
func UnmarshalStreamDatabase(data []byte) (StreamDatabase, error) {
	db, err := hs.DeserializeDatabase(data)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	return newStreamDatabase(db), nil
}

// UnmarshalVectoredDatabase reconstruct a vectored database from a stream of bytes.
func UnmarshalVectoredDatabase(data []byte) (VectoredDatabase, error) {
	db, err := hs.DeserializeDatabase(data)
	if err != nil {
		return nil, err // nolint: wrapcheck
	}

	return newVectoredDatabase(db), nil
}

// SerializedDatabaseSize reports the size that would be required by a database if it were deserialized.
func SerializedDatabaseSize(data []byte) (int, error) { return hs.SerializedDatabaseSize(data) } // nolint: wrapcheck

// SerializedDatabaseInfo provides information about a serialized database.
func SerializedDatabaseInfo(data []byte) (DbInfo, error) {
	i, err := hs.SerializedDatabaseInfo(data)

	return DbInfo(i), err
}

func (d *baseDatabase) Db() hs.Database { return d.db } // nolint: stylecheck

func (d *baseDatabase) Size() (int, error) { return hs.DatabaseSize(d.db) } // nolint: wrapcheck

func (d *baseDatabase) Info() (DbInfo, error) {
	i, err := hs.DatabaseInfo(d.db)
	if err != nil {
		return "", err //nolint: wrapcheck
	}

	return DbInfo(i), nil
}

func (d *baseDatabase) Close() error { return hs.FreeDatabase(d.db) } // nolint: wrapcheck

func (d *baseDatabase) Marshal() ([]byte, error) { return hs.SerializeDatabase(d.db) } // nolint: wrapcheck

func (d *baseDatabase) Unmarshal(data []byte) error { return hs.DeserializeDatabaseAt(data, d.db) } // nolint: wrapcheck

type blockDatabase struct {
	*blockMatcher
}

func newBlockDatabase(db hs.Database) *blockDatabase {
	return &blockDatabase{newBlockMatcher(newBlockScanner(newBaseDatabase(db)))}
}

type streamDatabase struct {
	*streamMatcher
}

func newStreamDatabase(db hs.Database) *streamDatabase {
	return &streamDatabase{newStreamMatcher(newStreamScanner(newBaseDatabase(db)))}
}

func (db *streamDatabase) StreamSize() (int, error) { return hs.StreamSize(db.db) } // nolint: wrapcheck

type vectoredDatabase struct {
	*vectoredMatcher
}

func newVectoredDatabase(db hs.Database) *vectoredDatabase {
	return &vectoredDatabase{newVectoredMatcher(newVectoredScanner(newBaseDatabase(db)))}
}
