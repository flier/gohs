package hyperscan

import (
	"errors"

	"github.com/flier/gohs/internal/hs"
)

// ErrTooManyMatches means too many matches.
var ErrTooManyMatches = errors.New("too many matches")

// MatchContext represents a match context.
type MatchContext interface {
	Database() Database

	Scratch() Scratch

	UserData() interface{}
}

// MatchEvent indicates a match event.
type MatchEvent interface {
	Id() uint

	From() uint64

	To() uint64

	Flags() ScanFlag
}

type ScanFlag = hs.ScanFlag

// MatchHandler handles match events.
type MatchHandler = hs.MatchEventHandler
