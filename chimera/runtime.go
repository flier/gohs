package chimera

import (
	"github.com/flier/gohs/hyperscan"
	"github.com/flier/gohs/internal/ch"
)

// Platform is a type containing information on the target platform.
type Platform = hyperscan.Platform

// ValidPlatform test the current system architecture.
var ValidPlatform = hyperscan.ValidPlatform

// Callback return value used to tell the Chimera matcher what to do after processing this match.
type Callback = ch.Callback

const (
	Continue    Callback = ch.Continue    // Continue matching.
	Terminate   Callback = ch.Terminate   // Terminate matching.
	SkipPattern Callback = ch.SkipPattern // Skip remaining matches for this ID and continue.
)

// Capture representing a captured subexpression within a match.
type Capture = ch.Capture

// Type used to differentiate the errors raised with the `ErrorEventHandler` callback.
type ErrorEvent = ch.ErrorEvent

const (
	// PCRE hits its match limit and reports PCRE_ERROR_MATCHLIMIT.
	ErrMatchLimit ErrorEvent = ch.ErrMatchLimit
	// PCRE hits its recursion limit and reports PCRE_ERROR_RECURSIONLIMIT.
	ErrRecursionLimit ErrorEvent = ch.ErrRecursionLimit
)

// Definition of the chimera event callback handler.
type Handler interface {
	// OnMatch will be invoked whenever a match is located in the target data during the execution of a scan.
	OnMatch(id uint, from, to uint64, flags uint, captured []*Capture, context interface{}) Callback

	// OnError will be invoked when an error event occurs during matching;
	// this indicates that some matches for a given expression may not be reported.
	OnError(event ErrorEvent, id uint, info, context interface{}) Callback
}

type MatchHandlerFunc func(id uint, from, to uint64, flags uint, captured []*Capture, context interface{}) Callback

func (f MatchHandlerFunc) OnMatch(id uint, from, to uint64, flags uint, captured []*Capture, ctx interface{}) Callback {
	return f(id, from, to, flags, captured, ctx)
}

func (f MatchHandlerFunc) OnError(event ErrorEvent, id uint, info, context interface{}) Callback {
	return Terminate
}
