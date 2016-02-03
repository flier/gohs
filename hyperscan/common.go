package hyperscan

func Version() string { return hsVersion() }

type Database interface {
	Size() int

	StreamSize() int

	Info() string

	Close() error

	Marshal() ([]byte, error)

	Unmarshal([]byte) error
}

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

func (d *database) Size() int {
	size, _ := hsDatabaseSize(d.db)

	return size
}

func (d *database) StreamSize() int {
	size, _ := hsStreamSize(d.db)

	return size
}

func (d *database) Info() string {
	info, _ := hsDatabaseInfo(d.db)

	return info
}

func (d *database) Close() error { return hsFreeDatabase(d.db) }

func (d *database) Marshal() ([]byte, error) { return hsSerializeDatabase(d.db) }

func (d *database) Unmarshal(data []byte) error { return hsDeserializeDatabaseAt(data, d.db) }
