//go:build hyperscan_v54
// +build hyperscan_v54

package hyperscan

/*
#include <hs.h>
*/
import "C"

const (
	AVX512VBMI CpuFeature = C.HS_CPU_FEATURES_AVX512VBMI // AVX512VBMI is a CPU features flag indicates that the target platform supports Intel(R) Advanced Vector Extensions 512 Vector Byte Manipulation Instructions (Intel(R) AVX512VBMI)
)

const (
	Icelake       TuneFlag = C.HS_TUNE_FAMILY_ICL // Icelake indicates that the compiled database should be tuned for the Icelake microarchitecture.
	IcelakeServer TuneFlag = C.HS_TUNE_FAMILY_ICX // IcelakeServer indicates that the compiled database should be tuned for the Icelake Server microarchitecture.
)
