//go:build hyperscan_v54
// +build hyperscan_v54

package hs

/*
#include <hs.h>
*/
import "C"

const (
	// AVX512VBMI is a CPU features flag indicates that the target platform
	// supports Intel(R) Advanced Vector Extensions 512 Vector Byte Manipulation Instructions (Intel(R) AVX512VBMI)
	AVX512VBMI CpuFeature = C.HS_CPU_FEATURES_AVX512VBMI
)

const (
	// Icelake indicates that the compiled database should be tuned for the Icelake microarchitecture.
	Icelake TuneFlag = C.HS_TUNE_FAMILY_ICL
	// IcelakeServer indicates that the compiled database should be tuned for the Icelake Server microarchitecture.
	IcelakeServer TuneFlag = C.HS_TUNE_FAMILY_ICX
)
