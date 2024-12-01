package chimera

import (
	"github.com/flier/gohs/internal/ch"
)

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

// HandlerFunc type is an adapter to allow the use of ordinary functions as Chimera handlers.
// If f is a function with the appropriate signature, HandlerFunc(f) is a Handler that calls f.
type HandlerFunc func(id uint, from, to uint64, flags uint, captured []*Capture, context interface{}) Callback

// OnMatch will be invoked whenever a match is located in the target data during the execution of a scan.
func (f HandlerFunc) OnMatch(id uint, from, to uint64, flags uint, captured []*Capture, ctx interface{}) Callback {
	return f(id, from, to, flags, captured, ctx)
}

// OnError will be invoked when an error event occurs during matching;
// this indicates that some matches for a given expression may not be reported.
func (f HandlerFunc) OnError(event ErrorEvent, id uint, info, context interface{}) Callback {
	return Terminate
}
