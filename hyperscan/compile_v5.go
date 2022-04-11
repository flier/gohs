//go:build !hyperscan_v4
// +build !hyperscan_v4

package hyperscan

import (
	"github.com/flier/gohs/internal/hs"
)

const (
	// Combination represents logical combination.
	Combination CompileFlag = hs.Combination
	// Quiet represents don't do any match reporting.
	Quiet CompileFlag = hs.Quiet
)
