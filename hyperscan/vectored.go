package hyperscan

import "github.com/flier/gohs/internal/hs"

// VectoredScanner is the vectored regular expression scanner.
type VectoredScanner interface {
	Scan(data [][]byte, scratch *Scratch, handler MatchHandler, context interface{}) error
}

// VectoredMatcher implements regular expression search.
type VectoredMatcher interface{}

type vectoredScanner struct {
	*baseDatabase
}

func newVectoredScanner(vdb *baseDatabase) *vectoredScanner {
	return &vectoredScanner{vdb}
}

func (vs *vectoredScanner) Scan(data [][]byte, s *Scratch, handler MatchHandler, context interface{}) (err error) {
	if s == nil {
		s, err = NewScratch(vs)

		if err != nil {
			return
		}

		defer func() {
			_ = s.Free()
		}()
	}

	return hs.ScanVector(vs.db, data, 0, s.s, handler, context) // nolint: wrapcheck
}

type vectoredMatcher struct {
	*vectoredScanner
}

func newVectoredMatcher(scanner *vectoredScanner) *vectoredMatcher {
	return &vectoredMatcher{vectoredScanner: scanner}
}
