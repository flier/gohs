//go:build hyperscan_v54
// +build hyperscan_v54

package ch

// #include <ch.h>
import "C"

// ErrUnknown is the unexpected internal error from Hyperscan.
const ErrUnknownHSError Error = C.CH_UNKNOWN_HS_ERROR

func init() {
	ErrorMessages[C.CH_UNKNOWN_HS_ERROR] = "Unexpected internal error from Hyperscan."
}
