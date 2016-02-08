package hyperscan

type Scratch interface {
	Size() (int, error)

	Realloc(db Database) error

	Clone() (Scratch, error)

	Close() error
}

type scratch struct {
	s hsScratch
}

func Alloc(db Database) (Scratch, error) {
	s, err := hsAllocScratch(db.(*baseDatabase).db)

	if err != nil {
		return nil, err
	}

	return &scratch{s}, nil
}

func (s *scratch) Size() (int, error) { return hsScratchSize(s.s) }

func (s *scratch) Realloc(db Database) error {
	if err := hsReallocScratch(db.(*baseDatabase).db, &s.s); err != nil {
		return err
	}

	return nil
}

func (s *scratch) Clone() (Scratch, error) {
	cloned, err := hsCloneScratch(s.s)

	if err != nil {
		return nil, err
	}

	return &scratch{cloned}, nil
}

func (s *scratch) Close() error { return hsFreeScratch(s.s) }
