//go:build go1.17
// +build go1.17

package handle

import "C"
import "runtime/cgo"

// Handle provides a way to pass values that contain Go pointers (pointers to memory allocated by Go)
// between Go and C without breaking the cgo pointer passing rules.
type Handle = cgo.Handle

// New returns a handle for a given value.
var New = cgo.NewHandle
