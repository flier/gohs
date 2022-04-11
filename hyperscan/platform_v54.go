//go:build hyperscan_v54
// +build hyperscan_v54

package hyperscan

import "github.com/flier/gohs/internal/hs"

const (
	// AVX512VBMI is a CPU features flag indicates that the target platform
	// supports Intel(R) Advanced Vector Extensions 512 Vector Byte Manipulation Instructions (Intel(R) AVX512VBMI)
	AVX512VBMI CpuFeature = hs.AVX512VBMI
)

const (
	// Icelake indicates that the compiled database should be tuned for the Icelake microarchitecture.
	Icelake TuneFlag = hs.Icelake
	// IcelakeServer indicates that the compiled database should be tuned for the Icelake Server microarchitecture.
	IcelakeServer TuneFlag = hs.IcelakeServer
)
