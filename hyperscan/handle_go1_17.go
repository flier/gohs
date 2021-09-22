//go:build go1.17
// +build go1.17

package hyperscan

import "C"
import "runtime/cgo"

type Handle = cgo.Handle

var NewHandle = cgo.NewHandle
