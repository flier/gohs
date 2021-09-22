//go:build go1.17
// +build go1.17

package hyperscan

import "C"
import "runtime/cgo"

// Handle provides a way to pass values that contain Go pointers (pointers to memory allocated by Go)
// between Go and C without breaking the cgo pointer passing rules.
type Handle = cgo.Handle

// NewHandle returns a handle for a given value.
var NewHandle = cgo.NewHandle
