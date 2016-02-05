package hyperscan

import (
	"regexp"
)

type Database interface {
	Version() string

	Mode() ModeFlag

	Info() string

	DatabaseSize() int

	StreamSize() int

	Close() error

	Marshal() ([]byte, error)

	Unmarshal([]byte) error
}

var (
	regexInfo = regexp.MustCompile(`^Version: (\d+\.\d+\.\d+) Features: ([\w\s]+) Mode: (\w+)$`)
)

func EngineVersion() string { return hsVersion() }

type database struct {
	db hsDatabase
}

func UnmarshalDatabase(data []byte) (Database, error) {
	db, err := hsDeserializeDatabase(data)

	if err != nil {
		return nil, err
	}

	return &database{db}, nil
}

func DatabaseSize(data []byte) (int, error) { return hsSerializedDatabaseSize(data) }

func DatabaseInfo(data []byte) (string, error) { return hsSerializedDatabaseInfo(data) }

func (d *database) DatabaseSize() int {
	size, err := hsDatabaseSize(d.db)

	if err != nil {
		panic(err)
	}

	return size
}

func (d *database) StreamSize() int {
	size, err := hsStreamSize(d.db)

	if err != nil {
		panic(err)
	}

	return size
}

func (d *database) Version() string {
	return regexInfo.FindStringSubmatch(d.Info())[1]
}

func (d *database) Mode() ModeFlag {
	mode, _ := ParseModeFlag(regexInfo.FindStringSubmatch(d.Info())[3])

	return mode
}

func (d *database) Info() string {
	info, err := hsDatabaseInfo(d.db)

	if err != nil {
		panic(err)
	}

	return info
}

func (d *database) Close() error { return hsFreeDatabase(d.db) }

func (d *database) Marshal() ([]byte, error) { return hsSerializeDatabase(d.db) }

func (d *database) Unmarshal(data []byte) error { return hsDeserializeDatabaseAt(data, d.db) }
