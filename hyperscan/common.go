package hyperscan

import (
	"errors"
	"regexp"
)

var (
	regexInfo = regexp.MustCompile(`^Version: (\d+\.\d+\.\d+) Features: ([\w\s]+) Mode: (\w+)$`)
)

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

type Database interface {
	Info() (DatabaseInfo, error)

	Size() (int, error)

	Close() error

	Marshal() ([]byte, error)

	Unmarshal([]byte) error
}

type BlockDatabase interface {
	Database
}

type StreamDatabase interface {
	Database

	StreamSize() (int, error)
}

type VectoredDatabase interface {
	Database
}

func EngineVersion() string { return hsVersion() }

type baseDatabase struct {
	db hsDatabase
}

func UnmarshalDatabase(data []byte) (Database, error) {
	db, err := hsDeserializeDatabase(data)

	if err != nil {
		return nil, err
	}

	return &baseDatabase{db}, nil
}

func Size(data []byte) (int, error) { return hsSerializedDatabaseSize(data) }

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
}

type streamDatabase struct {
	*baseDatabase
}

func (d *streamDatabase) StreamSize() (int, error) { return hsStreamSize(d.db) }

type vectoredDatabase struct {
	*baseDatabase
}
