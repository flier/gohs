package hyperscan

import (
	"errors"
	"regexp"

	"github.com/hashicorp/go-multierror"
)

type Database interface {
	// Provides information about a database.
	Info() (DatabaseInfo, error)

	// Provides the size of the given database in bytes.
	Size() (int, error)

	// Free a compiled pattern database.
	Close() error

	// Serialize a pattern database to a stream of bytes.
	Marshal() ([]byte, error)

	// Reconstruct a pattern database from a stream of bytes at a given memory location.
	Unmarshal([]byte) error
}

type BlockDatabase interface {
	Database
	BlockScanner
	BlockMatcher
}

type StreamDatabase interface {
	Database
	StreamScanner
	StreamMatcher

	StreamSize() (int, error)
}

type VectoredDatabase interface {
	Database
	VectoredScanner
	VectoredMatcher
}

var (
	regexInfo = regexp.MustCompile(`^Version: (\d+\.\d+\.\d+) Features: ([\w\s]+) Mode: (\w+)$`)
)

// The version and platform information for the supplied database
type DatabaseInfo string

func (i DatabaseInfo) Version() (string, error) {
	matched := regexInfo.FindStringSubmatch(string(i))

	if len(matched) != 4 {
		return "", errors.New("invalid database info")
	}

	return matched[1], nil
}

func (i DatabaseInfo) Mode() (ModeFlag, error) {
	matched := regexInfo.FindStringSubmatch(string(i))

	if len(matched) != 4 {
		return 0, errors.New("invalid database info")
	}
	return ParseModeFlag(matched[3])
}

// Utility function for identifying this release version.
func Version() string { return hsVersion() }

type baseDatabase struct {
	db hsDatabase
}

// Utility function for reconstructing a pattern database from a stream of bytes.
func Unmarshal(data []byte) (Database, error) {
	db, err := hsDeserializeDatabase(data)

	if err != nil {
		return nil, err
	}

	return &baseDatabase{db}, nil
}

// Utility function for reporting the size that would be required by a database if it were deserialized.
func Size(data []byte) (int, error) { return hsSerializedDatabaseSize(data) }

// Utility function providing information about a serialized database.
func Info(data []byte) (DatabaseInfo, error) {
	i, err := hsSerializedDatabaseInfo(data)

	return DatabaseInfo(i), err
}

func (d *baseDatabase) Size() (int, error) { return hsDatabaseSize(d.db) }

func (d *baseDatabase) Info() (DatabaseInfo, error) {
	i, err := hsDatabaseInfo(d.db)

	return DatabaseInfo(i), err
}

func (d *baseDatabase) Close() error { return hsFreeDatabase(d.db) }

func (d *baseDatabase) Marshal() ([]byte, error) { return hsSerializeDatabase(d.db) }

func (d *baseDatabase) Unmarshal(data []byte) error { return hsDeserializeDatabaseAt(data, d.db) }

type blockDatabase struct {
	*baseDatabase
	*blockScanner
	*blockMatcher
}

func newBlockDatabase(db hsDatabase) (*blockDatabase, error) {
	bdb := &blockDatabase{baseDatabase: &baseDatabase{db}}

	bdb.blockScanner = newBlockScanner(bdb)
	bdb.blockMatcher = newBlockMatcher(bdb.blockScanner)

	return bdb, nil
}

func (d *blockDatabase) Close() error {
	var result *multierror.Error

	if d.blockMatcher != nil {
		if err := d.blockMatcher.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if err := d.baseDatabase.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

type streamDatabase struct {
	*baseDatabase
	*streamScanner
	*streamMatcher
}

func newStreamDatabase(db hsDatabase) (*streamDatabase, error) {
	sdb := &streamDatabase{baseDatabase: &baseDatabase{db}}

	sdb.streamScanner = newStreamScanner(sdb)
	sdb.streamMatcher = newStreamMatcher(sdb.streamScanner)

	return sdb, nil
}

func (d *streamDatabase) StreamSize() (int, error) { return hsStreamSize(d.db) }

func (d *streamDatabase) Close() error {
	var result *multierror.Error

	if d.streamMatcher != nil {
		if err := d.streamMatcher.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if err := d.baseDatabase.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}

type vectoredDatabase struct {
	*baseDatabase
	*vectoredScanner
	*vectoredMatcher
}

func newVectoredDatabase(db hsDatabase) (*vectoredDatabase, error) {
	vdb := &vectoredDatabase{baseDatabase: &baseDatabase{db}}

	vdb.vectoredScanner = newVectoredScanner(vdb)
	vdb.vectoredMatcher = newVectoredMatcher(vdb.vectoredScanner)

	return vdb, nil
}

func (d *vectoredDatabase) Close() error {
	var result *multierror.Error

	if d.vectoredMatcher != nil {
		if err := d.vectoredMatcher.Close(); err != nil {
			result = multierror.Append(result, err)
		}
	}

	if err := d.baseDatabase.Close(); err != nil {
		result = multierror.Append(result, err)
	}

	return result.ErrorOrNil()
}
