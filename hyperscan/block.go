package hyperscan

import (
	"errors"

	"github.com/flier/gohs/internal/hs"
)

// BlockScanner is the block (non-streaming) regular expression scanner.
type BlockScanner interface {
	// This is the function call in which the actual pattern matching takes place for block-mode pattern databases.
	Scan(data []byte, scratch *Scratch, handler MatchHandler, context interface{}) error
}

// BlockMatcher implements regular expression search.
type BlockMatcher interface {
	// Find returns a slice holding the text of the leftmost match in b of the regular expression.
	// A return value of nil indicates no match.
	Find(data []byte) []byte

	// FindIndex returns a two-element slice of integers defining
	// the location of the leftmost match in b of the regular expression.
	// The match itself is at b[loc[0]:loc[1]]. A return value of nil indicates no match.
	FindIndex(data []byte) []int

	// FindAll is the 'All' version of Find; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAll(data []byte, n int) [][]byte

	// FindAllIndex is the 'All' version of FindIndex; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAllIndex(data []byte, n int) [][]int

	// FindString returns a string holding the text of the leftmost match in s of the regular expression.
	// If there is no match, the return value is an empty string, but it will also be empty
	// if the regular expression successfully matches an empty string.
	// Use FindStringIndex if it is necessary to distinguish these cases.
	FindString(s string) string

	// FindStringIndex returns a two-element slice of integers defining
	// the location of the leftmost match in s of the regular expression.
	// The match itself is at s[loc[0]:loc[1]]. A return value of nil indicates no match.
	FindStringIndex(s string) []int

	// FindAllString is the 'All' version of FindString; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAllString(s string, n int) []string

	// FindAllStringIndex is the 'All' version of FindStringIndex;
	// it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAllStringIndex(s string, n int) [][]int

	// Match reports whether the pattern database matches the byte slice b.
	Match(b []byte) bool

	// MatchString reports whether the pattern database matches the string s.
	MatchString(s string) bool
}

// BlockDatabase scan the target data that is a discrete,
// contiguous block which can be scanned in one call and does not require state to be retained.
type BlockDatabase interface {
	Database
	BlockScanner
	BlockMatcher
}

type blockDatabase struct {
	*blockMatcher
}

func newBlockDatabase(db hs.Database) *blockDatabase {
	return &blockDatabase{newBlockMatcher(newBlockScanner(newBaseDatabase(db)))}
}

type blockScanner struct {
	*baseDatabase
}

func newBlockScanner(bdb *baseDatabase) *blockScanner {
	return &blockScanner{bdb}
}

func (bs *blockScanner) Scan(data []byte, s *Scratch, handler MatchHandler, context interface{}) (err error) {
	if s == nil {
		s, err = NewScratch(bs)

		if err != nil {
			return
		}

		defer func() {
			_ = s.Free()
		}()
	}

	return hs.Scan(bs.db, data, 0, s.s, handler, context) //nolint: wrapcheck
}

type blockMatcher struct {
	*blockScanner
	*hs.MatchRecorder
	n int
}

func newBlockMatcher(scanner *blockScanner) *blockMatcher {
	return &blockMatcher{blockScanner: scanner}
}

func (m *blockMatcher) Handle(id uint, from, to uint64, flags uint, context interface{}) error {
	err := m.MatchRecorder.Handle(id, from, to, flags, context)
	if err != nil {
		return err //nolint: wrapcheck
	}

	if m.n < 0 {
		return nil
	}

	if m.n < len(m.Events) {
		m.Events = m.Events[:m.n]

		return ErrTooManyMatches
	}

	return nil
}

func (m *blockMatcher) scan(data []byte) error {
	m.MatchRecorder = &hs.MatchRecorder{}

	return m.blockScanner.Scan(data, nil, m.Handle, nil)
}

const findIndexMatches = 2

func (m *blockMatcher) Find(data []byte) []byte {
	if loc := m.FindIndex(data); len(loc) == findIndexMatches {
		return data[loc[0]:loc[1]]
	}

	return nil
}

func (m *blockMatcher) FindIndex(data []byte) []int {
	if m.Match(data) && len(m.Events) == 1 {
		return []int{int(m.Events[0].From), int(m.Events[0].To)}
	}

	return nil
}

func (m *blockMatcher) FindAll(data []byte, n int) (matches [][]byte) {
	if locs := m.FindAllIndex(data, n); len(locs) > 0 {
		for _, loc := range locs {
			matches = append(matches, data[loc[0]:loc[1]])
		}
	}

	return
}

func (m *blockMatcher) FindAllIndex(data []byte, n int) (locs [][]int) {
	if n < 0 {
		n = len(data) + 1
	}

	m.n = n

	if err := m.scan(data); (err == nil || errors.Is(err, ErrScanTerminated)) && len(m.Events) > 0 {
		for _, e := range m.Events {
			locs = append(locs, []int{int(e.From), int(e.To)})
		}
	}

	return
}

func (m *blockMatcher) FindString(s string) string {
	return string(m.Find([]byte(s)))
}

func (m *blockMatcher) FindStringIndex(s string) (loc []int) {
	return m.FindIndex([]byte(s))
}

func (m *blockMatcher) FindAllString(s string, n int) (results []string) {
	for _, m := range m.FindAll([]byte(s), n) {
		results = append(results, string(m))
	}

	return
}

func (m *blockMatcher) FindAllStringIndex(s string, n int) [][]int {
	return m.FindAllIndex([]byte(s), n)
}

func (m *blockMatcher) Match(data []byte) bool {
	m.n = 1

	err := m.scan(data)

	return (err == nil || errors.Is(err, ErrScanTerminated)) && len(m.Events) == m.n
}

func (m *blockMatcher) MatchString(s string) bool {
	return m.Match([]byte(s))
}
