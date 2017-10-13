package hyperscan

import (
	"errors"
)

var (
	errTooManyMatches = errors.New("too many matches")
)

// Scratch is a Hyperscan scratch space.
type Scratch struct {
	s hsScratch
}

// NewScratch allocate a "scratch" space for use by Hyperscan.
// This is required for runtime use, and one scratch space per thread,
// or concurrent caller, is required.
func NewScratch(db Database) (*Scratch, error) {
	s, err := hsAllocScratch(db.(database).Db())

	if err != nil {
		return nil, err
	}

	return &Scratch{s}, nil
}

// Size provides the size of the given scratch space.
func (s *Scratch) Size() (int, error) { return hsScratchSize(s.s) }

// Realloc reallocate the scratch for another database.
func (s *Scratch) Realloc(db Database) error {
	return hsReallocScratch(db.(database).Db(), &s.s)
}

// Clone allocate a scratch space that is a clone of an existing scratch space.
func (s *Scratch) Clone() (*Scratch, error) {
	cloned, err := hsCloneScratch(s.s)

	if err != nil {
		return nil, err
	}

	return &Scratch{cloned}, nil
}

// Free a scratch block previously allocated
func (s *Scratch) Free() error { return hsFreeScratch(s.s) }

type MatchContext interface {
	Database() Database

	Scratch() Scratch

	UserData() interface{}
}

type MatchEvent interface {
	Id() uint

	From() uint64

	To() uint64

	Flags() ScanFlag
}

type MatchHandler hsMatchEventHandler

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

	// FindIndex returns a two-element slice of integers defining the location of the leftmost match in b of the regular expression.
	// The match itself is at b[loc[0]:loc[1]]. A return value of nil indicates no match.
	FindIndex(data []byte) (loc []int)

	// FindAll is the 'All' version of Find; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAll(data []byte, n int) [][]byte

	// FindAllIndex is the 'All' version of FindIndex; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAllIndex(data []byte, n int) [][]int

	// FindString returns a string holding the text of the leftmost match in s of the regular expression.
	// If there is no match, the return value is an empty string, but it will also be empty
	// if the regular expression successfully matches an empty string. Use FindStringIndex if it is necessary to distinguish these cases.
	FindString(s string) string

	// FindStringIndex returns a two-element slice of integers defining the location of the leftmost match in s of the regular expression.
	// The match itself is at s[loc[0]:loc[1]]. A return value of nil indicates no match.
	FindStringIndex(s string) (loc []int)

	// FindAllString is the 'All' version of FindString; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAllString(s string, n int) []string

	// FindAllStringIndex is the 'All' version of FindStringIndex; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAllStringIndex(s string, n int) [][]int

	// Match reports whether the pattern database matches the byte slice b.
	Match(b []byte) bool

	// MatchString reports whether the pattern database matches the string s.
	MatchString(s string) bool
}

// Stream exist in the Hyperscan library so that pattern matching state can be maintained across multiple blocks of target data
type Stream interface {
	Scan(data []byte) error

	Close() error

	Reset() error

	Clone() (Stream, error)
}

// StreamScanner is the streaming regular expression scanner.
type StreamScanner interface {
	Open(flags ScanFlag, scratch *Scratch, handler MatchHandler, context interface{}) (Stream, error)
}

// StreamMatcher implements regular expression search.
type StreamMatcher interface {
}

// VectoredScanner is the vectored regular expression scanner.
type VectoredScanner interface {
	Scan(data [][]byte, scratch *Scratch, handler MatchHandler, context interface{}) error
}

// VectoredMatcher implements regular expression search.
type VectoredMatcher interface {
}

type stream struct {
	stream  hsStream
	flags   ScanFlag
	scratch hsScratch
	handler hsMatchEventHandler
	context interface{}
}

func (s *stream) Scan(data []byte) error {
	return hsScanStream(s.stream, data, s.flags, s.scratch, s.handler, s.context)
}

func (s *stream) Close() error {
	return hsCloseStream(s.stream, s.scratch, s.handler, s.context)
}

func (s *stream) Reset() error {
	return hsResetStream(s.stream, s.flags, s.scratch, s.handler, s.context)
}

func (s *stream) Clone() (Stream, error) {
	ss, err := hsCopyStream(s.stream)

	if err != nil {
		return nil, err
	}

	return &stream{ss, s.flags, s.scratch, s.handler, s.context}, nil
}

type streamScanner struct {
	*baseDatabase
}

func newStreamScanner(db *baseDatabase) *streamScanner {
	return &streamScanner{baseDatabase: db}
}

func (ss *streamScanner) Close() error {
	return nil
}

func (ss *streamScanner) Open(flags ScanFlag, sc *Scratch, handler MatchHandler, context interface{}) (Stream, error) {
	s, err := hsOpenStream(ss.db, flags)

	if err != nil {
		return nil, err
	}

	if sc == nil {
		sc, err = NewScratch(ss)

		if err != nil {
			return nil, err
		}

		defer sc.Free()
	}

	return &stream{s, flags, sc.s, hsMatchEventHandler(handler), context}, nil
}

type vectoredScanner struct {
	*baseDatabase
}

func newVectoredScanner(vdb *baseDatabase) *vectoredScanner {
	return &vectoredScanner{vdb}
}

func (vs *vectoredScanner) Close() error { return nil }

func (vs *vectoredScanner) Scan(data [][]byte, s *Scratch, handler MatchHandler, context interface{}) (err error) {
	if s == nil {
		s, err = NewScratch(vs)

		if err != nil {
			return err
		}

		defer s.Free()
	}

	err = hsScanVector(vs.db, data, 0, s.s, hsMatchEventHandler(handler), context)

	if err != nil {
		return err
	}

	return nil
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
			return err
		}

		defer s.Free()
	}

	err = hsScan(bs.db, data, 0, s.s, hsMatchEventHandler(handler), context)

	if err != nil {
		return err
	}

	return nil
}

type blockMatcher struct {
	*blockScanner
	handler *matchRecorder
	n       int
}

func newBlockMatcher(scanner *blockScanner) *blockMatcher {
	return &blockMatcher{blockScanner: scanner}
}

func (m *blockMatcher) Handle(id uint, from, to uint64, flags uint, context interface{}) error {
	m.n--

	if m.n == 0 {
		m.handler.err = errTooManyMatches
	}

	return m.handler.Handle(id, from, to, flags, context)
}

func (m *blockMatcher) scan(data []byte) error {
	m.handler = &matchRecorder{}

	return m.blockScanner.Scan(data, nil, m.Handle, nil)
}

func (m *blockMatcher) Find(data []byte) []byte {
	if loc := m.FindIndex(data); loc != nil && len(loc) == 2 {
		return data[loc[0]:loc[1]]
	}

	return nil
}

func (m *blockMatcher) FindIndex(data []byte) []int {
	if m.Match(data) && len(m.handler.matched) == 1 {
		return []int{int(m.handler.matched[0].from), int(m.handler.matched[0].to)}
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
	m.n = n

	if err := m.scan(data); (err == nil || err.(HsError) == ErrScanTerminated) && len(m.handler.matched) > 0 {
		for _, e := range m.handler.matched {
			locs = append(locs, []int{int(e.from), int(e.to)})
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

	return err != nil && err.(HsError) == ErrScanTerminated
}

func (m *blockMatcher) MatchString(s string) bool {
	return m.Match([]byte(s))
}

type streamMatcher struct {
	*streamScanner
}

func newStreamMatcher(scanner *streamScanner) *streamMatcher {
	return &streamMatcher{streamScanner: scanner}
}

type vectoredMatcher struct {
	*vectoredScanner
}

func newVectoredMatcher(scanner *vectoredScanner) *vectoredMatcher {
	return &vectoredMatcher{vectoredScanner: scanner}
}
