package hyperscan

import (
	"errors"
	"fmt"
	"io"

	"github.com/flier/gohs/internal/hs"
)

// Stream exist in the Hyperscan library so that pattern matching state can be maintained
// across multiple blocks of target data.
type Stream interface {
	Scan(data []byte) error

	Close() error

	Reset() error

	Clone() (Stream, error)
}

// StreamScanner is the streaming regular expression scanner.
type StreamScanner interface {
	Open(flags ScanFlag, scratch *Scratch, handler MatchHandler, context interface{}) (Stream, error)

	Scan(reader io.Reader, scratch *Scratch, handler MatchHandler, context interface{}) error
}

// StreamMatcher implements regular expression search.
type StreamMatcher interface {
	// Find returns a slice holding the text of the leftmost match in b of the regular expression.
	// A return value of nil indicates no match.
	Find(reader io.ReadSeeker) []byte

	// FindIndex returns a two-element slice of integers defining
	// the location of the leftmost match in b of the regular expression.
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
	ResetAndExpand(stream Stream, buf []byte, flags ScanFlag, scratch *Scratch,
		handler MatchHandler, context interface{}) (Stream, error)
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

type streamDatabase struct {
	*streamMatcher
}

func newStreamDatabase(db hs.Database) *streamDatabase {
	return &streamDatabase{newStreamMatcher(newStreamScanner(newBaseDatabase(db)))}
}

func (db *streamDatabase) StreamSize() (int, error) { return hs.StreamSize(db.db) } //nolint: wrapcheck

const bufSize = 4096

type stream struct {
	stream       hs.Stream
	flags        ScanFlag
	scratch      hs.Scratch
	handler      hs.MatchEventHandler
	context      interface{}
	ownedScratch bool
}

func (s *stream) Scan(data []byte) error {
	return hs.ScanStream(s.stream, data, s.flags, s.scratch, s.handler, s.context) //nolint: wrapcheck
}

func (s *stream) Close() error {
	err := hs.CloseStream(s.stream, s.scratch, s.handler, s.context)

	if s.ownedScratch {
		_ = hs.FreeScratch(s.scratch)
	}

	return err //nolint: wrapcheck
}

func (s *stream) Reset() error {
	return hs.ResetStream(s.stream, s.flags, s.scratch, s.handler, s.context) //nolint: wrapcheck
}

func (s *stream) Clone() (Stream, error) {
	ss, err := hs.CopyStream(s.stream)
	if err != nil {
		return nil, fmt.Errorf("copy stream, %w", err)
	}

	scratch := s.scratch

	if s.ownedScratch {
		scratch, err = hs.CloneScratch(s.scratch)

		if err != nil {
			hs.FreeStream(ss)

			return nil, fmt.Errorf("clone scratch, %w", err)
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
	s, err := hs.OpenStream(ss.db, flags)
	if err != nil {
		return nil, fmt.Errorf("open stream, %w", err)
	}

	ownedScratch := false

	if sc == nil {
		sc, err = NewScratch(ss)
		if err != nil {
			return nil, fmt.Errorf("create scratch, %w", err)
		}

		ownedScratch = true
	}

	return &stream{s, flags, sc.s, handler, context, ownedScratch}, nil
}

func (ss *streamScanner) Scan(reader io.Reader, sc *Scratch, handler MatchHandler, context interface{}) error {
	stream, err := ss.Open(0, sc, handler, context)
	if err != nil {
		return err
	}
	defer stream.Close()

	buf := make([]byte, bufSize)

	for {
		n, err := reader.Read(buf)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read stream, %w", err)
		}

		if err = stream.Scan(buf[:n]); err != nil {
			return fmt.Errorf("scan stream, %w", err)
		}
	}
}

type streamMatcher struct {
	*streamScanner
	*hs.MatchRecorder
	n int
}

func newStreamMatcher(scanner *streamScanner) *streamMatcher {
	return &streamMatcher{streamScanner: scanner}
}

func (m *streamMatcher) Handle(id uint, from, to uint64, flags uint, context interface{}) error {
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

func (m *streamMatcher) scan(reader io.Reader) error {
	m.MatchRecorder = &hs.MatchRecorder{}

	stream, err := m.streamScanner.Open(0, nil, m.Handle, nil)
	if err != nil {
		return err
	}
	defer stream.Close()

	buf := make([]byte, bufSize)

	for {
		n, err := reader.Read(buf)

		if errors.Is(err, io.EOF) {
			return nil
		}

		if err != nil {
			return fmt.Errorf("read stream, %w", err)
		}

		if err = stream.Scan(buf[:n]); err != nil {
			return fmt.Errorf("scan stream, %w", err)
		}
	}
}

func (m *streamMatcher) read(reader io.ReadSeeker, loc []int) ([]byte, error) {
	if len(loc) != findIndexMatches {
		return nil, fmt.Errorf("location, %w", ErrInvalid)
	}

	offset := int64(loc[0])
	size := loc[1] - loc[0]

	_, err := reader.Seek(offset, io.SeekStart)
	if err != nil {
		return nil, fmt.Errorf("seek stream, %w", err)
	}

	buf := make([]byte, size)

	_, err = reader.Read(buf)

	if err != nil {
		return nil, fmt.Errorf("read data, %w", err)
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
	if m.Match(reader) && len(m.Events) == 1 {
		return []int{int(m.Events[0].From), int(m.Events[0].To)}
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

	if err := m.scan(reader); (err == nil || errors.Is(err, ErrScanTerminated)) && len(m.Events) > 0 {
		for _, e := range m.Events {
			locs = append(locs, []int{int(e.From), int(e.To)})
		}
	}

	return
}

func (m *streamMatcher) Match(reader io.Reader) bool {
	m.n = 1

	err := m.scan(reader)

	return (err == nil || errors.Is(err, ErrScanTerminated)) && len(m.Events) == m.n
}

var _ StreamCompressor = (*streamDatabase)(nil)

func (db *streamDatabase) Compress(s Stream) ([]byte, error) {
	size, err := db.StreamSize()
	if err != nil {
		return nil, fmt.Errorf("stream size, %w", err)
	}

	buf := make([]byte, size)

	buf, err = hs.CompressStream(s.(*stream).stream, buf)

	if err != nil {
		return nil, fmt.Errorf("compress stream, %w", err)
	}

	return buf, nil
}

func (db *streamDatabase) Expand(buf []byte, flags ScanFlag, sc *Scratch,
	handler MatchHandler, context interface{},
) (Stream, error) {
	var s hs.Stream

	err := hs.ExpandStream(db.db, &s, buf)
	if err != nil {
		return nil, fmt.Errorf("expand stream, %w", err)
	}

	ownedScratch := false

	if sc == nil {
		sc, err = NewScratch(db)
		if err != nil {
			return nil, fmt.Errorf("create scratch, %w", err)
		}

		ownedScratch = true
	}

	return &stream{s, flags, sc.s, handler, context, ownedScratch}, nil
}

func (db *streamDatabase) ResetAndExpand(s Stream, buf []byte, flags ScanFlag, sc *Scratch,
	handler MatchHandler, context interface{},
) (Stream, error) {
	ss, ok := s.(*stream)
	if !ok {
		return nil, fmt.Errorf("stream %v, %w", s, ErrInvalid)
	}

	ownedScratch := false

	if sc == nil {
		var err error

		sc, err = NewScratch(db)
		if err != nil {
			return nil, fmt.Errorf("create scratch, %w", err)
		}

		ownedScratch = true
	}

	err := hs.ResetAndExpandStream(ss.stream, buf, ss.scratch, ss.handler, ss.context)
	if err != nil {
		return nil, fmt.Errorf("reset and expand stream, %w", err)
	}

	return &stream{ss.stream, flags, sc.s, handler, context, ownedScratch}, nil
}
