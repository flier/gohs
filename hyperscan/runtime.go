package hyperscan

import (
	"errors"
	"fmt"
	"io"
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
	FindIndex(data []byte) []int

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
	FindStringIndex(s string) []int

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
	// Find returns a slice holding the text of the leftmost match in b of the regular expression.
	// A return value of nil indicates no match.
	Find(reader io.ReadSeeker) []byte

	// FindIndex returns a two-element slice of integers defining the location of the leftmost match in b of the regular expression.
	// The match itself is at b[loc[0]:loc[1]]. A return value of nil indicates no match.
	FindIndex(reader io.Reader) []int

	// FindAll is the 'All' version of Find; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAll(reader io.ReadSeeker, n int) [][]byte

	// FindAllIndex is the 'All' version of FindIndex; it returns a slice of all successive matches of the expression,
	// as defined by the 'All' description in the package comment. A return value of nil indicates no match.
	FindAllIndex(reader io.Reader, n int) [][]int

	// Match reports whether the pattern database matches the byte slice b.
	Match(reader io.Reader) bool
}

// StreamCompressor implements stream compressor.
type StreamCompressor interface {
	// Creates a compressed representation of the provided stream in the buffer provided.
	Compress(stream Stream) ([]byte, error)

	// Decompresses a compressed representation created by `CompressStream` into a new stream.
	Expand(buf []byte, flags ScanFlag, scratch *Scratch, handler MatchHandler, context interface{}) (Stream, error)

	// Decompresses a compressed representation created by `CompressStream` on top of the 'to' stream.
	ResetAndExpand(stream Stream, buf []byte, flags ScanFlag, scratch *Scratch, handler MatchHandler, context interface{}) (Stream, error)
}

// VectoredScanner is the vectored regular expression scanner.
type VectoredScanner interface {
	Scan(data [][]byte, scratch *Scratch, handler MatchHandler, context interface{}) error
}

// VectoredMatcher implements regular expression search.
type VectoredMatcher interface {
}

type stream struct {
	stream       hsStream
	flags        ScanFlag
	scratch      hsScratch
	handler      hsMatchEventHandler
	context      interface{}
	ownedScratch bool
}

func (s *stream) Scan(data []byte) error {
	return hsScanStream(s.stream, data, s.flags, s.scratch, s.handler, s.context)
}

func (s *stream) Close() error {
	err := hsCloseStream(s.stream, s.scratch, s.handler, s.context)

	if s.ownedScratch {
		hsFreeScratch(s.scratch)
	}

	return err
}

func (s *stream) Reset() error {
	return hsResetStream(s.stream, s.flags, s.scratch, s.handler, s.context)
}

func (s *stream) Clone() (Stream, error) {
	ss, err := hsCopyStream(s.stream)

	if err != nil {
		return nil, err
	}

	scratch := s.scratch

	if s.ownedScratch {
		scratch, err = hsCloneScratch(s.scratch)

		if err != nil {
			return nil, err
		}
	}

	return &stream{ss, s.flags, scratch, s.handler, s.context, s.ownedScratch}, nil
}

type streamScanner struct {
	*baseDatabase
}

func newStreamScanner(db *baseDatabase) *streamScanner {
	return &streamScanner{baseDatabase: db}
}

func (ss *streamScanner) Open(flags ScanFlag, sc *Scratch, handler MatchHandler, context interface{}) (Stream, error) {
	s, err := hsOpenStream(ss.db, flags)

	if err != nil {
		return nil, err
	}

	ownedScratch := false

	if sc == nil {
		sc, err = NewScratch(ss)

		if err != nil {
			return nil, err
		}

		ownedScratch = true
	}

	return &stream{s, flags, sc.s, hsMatchEventHandler(handler), context, ownedScratch}, nil
}

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
			return err
		}

		defer s.Free()
	}

	return hsScanVector(vs.db, data, 0, s.s, hsMatchEventHandler(handler), context)
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

	return hsScan(bs.db, data, 0, s.s, hsMatchEventHandler(handler), context)
}

type blockMatcher struct {
	*blockScanner
	*matchRecorder
	n int
}

func newBlockMatcher(scanner *blockScanner) *blockMatcher {
	return &blockMatcher{blockScanner: scanner}
}

func (m *blockMatcher) Handle(id uint, from, to uint64, flags uint, context interface{}) error {
	err := m.matchRecorder.Handle(id, from, to, flags, context)

	if err != nil {
		return err
	}

	if m.n < 0 {
		return nil
	}

	if m.n < len(m.matched) {
		m.matched = m.matched[:m.n]

		return errTooManyMatches
	}

	return nil
}

func (m *blockMatcher) scan(data []byte) error {
	m.matchRecorder = &matchRecorder{}

	return m.blockScanner.Scan(data, nil, m.Handle, nil)
}

func (m *blockMatcher) Find(data []byte) []byte {
	if loc := m.FindIndex(data); len(loc) == 2 {
		return data[loc[0]:loc[1]]
	}

	return nil
}

func (m *blockMatcher) FindIndex(data []byte) []int {
	if m.Match(data) && len(m.matched) == 1 {
		return []int{int(m.matched[0].from), int(m.matched[0].to)}
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

	if err := m.scan(data); (err == nil || err.(HsError) == ErrScanTerminated) && len(m.matched) > 0 {
		for _, e := range m.matched {
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
	*matchRecorder
	n int
}

func newStreamMatcher(scanner *streamScanner) *streamMatcher {
	return &streamMatcher{streamScanner: scanner}
}

func (m *streamMatcher) Handle(id uint, from, to uint64, flags uint, context interface{}) error {
	err := m.matchRecorder.Handle(id, from, to, flags, context)

	if err != nil {
		return err
	}

	if m.n < 0 {
		return nil
	}

	if m.n < len(m.matched) {
		m.matched = m.matched[:m.n]

		return errTooManyMatches
	}

	return nil
}

func (m *streamMatcher) scan(reader io.Reader) error {
	m.matchRecorder = &matchRecorder{}

	stream, err := m.streamScanner.Open(0, nil, m.Handle, nil)

	if err != nil {
		return err
	}

	buf := make([]byte, 4096)

	for {
		n, err := reader.Read(buf)

		if err == io.EOF {
			return stream.Close()
		}

		if err != nil {
			return err
		}

		if err = stream.Scan(buf[:n]); err != nil {
			return err
		}
	}
}

func (m *streamMatcher) read(reader io.ReadSeeker, loc []int) ([]byte, error) {
	if len(loc) != 2 {
		return nil, fmt.Errorf("invalid location")
	}

	offset := int64(loc[0])
	size := loc[1] - loc[0]

	_, err := reader.Seek(offset, io.SeekStart)

	if err != nil {
		return nil, err
	}

	buf := make([]byte, size)

	_, err = reader.Read(buf)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (m *streamMatcher) Find(reader io.ReadSeeker) []byte {
	loc := m.FindIndex(reader)
	buf, err := m.read(reader, loc)

	if err != nil {
		return nil
	}

	return buf
}

func (m *streamMatcher) FindIndex(reader io.Reader) []int {
	if m.Match(reader) && len(m.matched) == 1 {
		return []int{int(m.matched[0].from), int(m.matched[0].to)}
	}

	return nil
}

func (m *streamMatcher) FindAll(reader io.ReadSeeker, n int) (result [][]byte) {
	for _, loc := range m.FindAllIndex(reader, n) {
		if buf, err := m.read(reader, loc); err == nil {
			result = append(result, buf)
		}
	}

	return
}

func (m *streamMatcher) FindAllIndex(reader io.Reader, n int) (locs [][]int) {
	m.n = n

	if err := m.scan(reader); (err == nil || err.(HsError) == ErrScanTerminated) && len(m.matched) > 0 {
		for _, e := range m.matched {
			locs = append(locs, []int{int(e.from), int(e.to)})
		}
	}

	return
}

func (m *streamMatcher) Match(reader io.Reader) bool {
	m.n = 1

	err := m.scan(reader)

	return err != nil && err.(HsError) == ErrScanTerminated
}

type vectoredMatcher struct {
	*vectoredScanner
}

func newVectoredMatcher(scanner *vectoredScanner) *vectoredMatcher {
	return &vectoredMatcher{vectoredScanner: scanner}
}

var _ StreamCompressor = (*streamDatabase)(nil)

func (db *streamDatabase) Compress(s Stream) ([]byte, error) {
	size, err := db.StreamSize()

	if err != nil {
		return nil, err
	}

	buf := make([]byte, size)

	buf, err = hsCompressStream(s.(*stream).stream, buf)

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (db *streamDatabase) Expand(buf []byte, flags ScanFlag, sc *Scratch, handler MatchHandler, context interface{}) (Stream, error) {
	var s hsStream

	err := hsExpandStream(db.db, &s, buf)

	if err != nil {
		return nil, err
	}

	ownedScratch := false

	if sc == nil {
		sc, err = NewScratch(db)

		if err != nil {
			return nil, err
		}

		ownedScratch = true
	}

	return &stream{s, flags, sc.s, hsMatchEventHandler(handler), context, ownedScratch}, nil
}

func (db *streamDatabase) ResetAndExpand(s Stream, buf []byte, flags ScanFlag, sc *Scratch, handler MatchHandler, context interface{}) (Stream, error) {
	ss := s.(*stream)

	ownedScratch := false

	if sc == nil {
		var err error

		sc, err = NewScratch(db)

		if err != nil {
			return nil, err
		}

		ownedScratch = true
	}

	err := hsResetAndExpandStream(ss.stream, buf, ss.scratch, ss.handler, ss.context)

	if err != nil {
		return nil, err
	}

	return &stream{ss.stream, flags, sc.s, hsMatchEventHandler(handler), context, ownedScratch}, nil
}
