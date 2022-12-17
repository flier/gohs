package hyperscan

import "github.com/flier/gohs/internal/hs"

type TuneFlag = hs.TuneFlag

const (
	// Generic indicates that the compiled database should not be tuned for any particular target platform.
	Generic TuneFlag = hs.Generic
	// SandyBridge indicates that the compiled database should be tuned for the Sandy Bridge microarchitecture.
	SandyBridge TuneFlag = hs.SandyBridge
	// IvyBridge indicates that the compiled database should be tuned for the Ivy Bridge microarchitecture.
	IvyBridge TuneFlag = hs.IvyBridge
	// Haswell indicates that the compiled database should be tuned for the Haswell microarchitecture.
	Haswell TuneFlag = hs.Haswell
	// Silvermont indicates that the compiled database should be tuned for the Silvermont microarchitecture.
	Silvermont TuneFlag = hs.Silvermont
	// Broadwell indicates that the compiled database should be tuned for the Broadwell microarchitecture.
	Broadwell TuneFlag = hs.Broadwell
	// Skylake indicates that the compiled database should be tuned for the Skylake microarchitecture.
	Skylake TuneFlag = hs.Skylake
	// SkylakeServer indicates that the compiled database should be tuned for the Skylake Server microarchitecture.
	SkylakeServer TuneFlag = hs.SkylakeServer
	// Goldmont indicates that the compiled database should be tuned for the Goldmont microarchitecture.
	Goldmont TuneFlag = hs.Goldmont
)

// CpuFeature is the CPU feature support flags.
type CpuFeature = hs.CpuFeature //nolint: golint,stylecheck,revive

const (
	// AVX2 is a CPU features flag indicates that the target platform supports AVX2 instructions.
	AVX2 CpuFeature = hs.AVX2
	// AVX512 is a CPU features flag indicates that the target platform supports AVX512 instructions,
	// specifically AVX-512BW. Using AVX512 implies the use of AVX2.
	AVX512 CpuFeature = hs.AVX512
)

// Platform is a type containing information on the target platform.
type Platform interface {
	// Information about the target platform which may be used to guide the optimisation process of the compile.
	Tune() TuneFlag

	// Relevant CPU features available on the target platform
	CpuFeatures() CpuFeature
}

// NewPlatform create a new platform information on the target platform.
func NewPlatform(tune TuneFlag, cpu CpuFeature) Platform { return hs.NewPlatformInfo(tune, cpu) }

// PopulatePlatform populates the platform information based on the current host.
func PopulatePlatform() Platform {
	platform, _ := hs.PopulatePlatform()

	return platform
}
