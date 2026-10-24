package cloneconcurrency

import (
	"runtime"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cloneconcurrency/enums"

// CPUProfileType defines the CPU auto-scaling profile for pull operations.
// CPUProfileType aliases the strongly-typed enum from the enums package.
type CPUProfileType = enums.CPUProfileType

const (
	CPUProfileAuto     = enums.CPUProfileAuto
	CPUProfileHighPerf = enums.CPUProfileHighPerf
	CPUProfileLowCPU   = enums.CPUProfileLowCPU
)

// CPUProfile is an alias for CPUProfileType for backwards compatibility.
type CPUProfile = enums.CPUProfileType

// ConcurrencyPreset is an alias for CPUProfile for backwards compatibility with component specs.
type ConcurrencyPreset = enums.CPUProfileType

const (
	// PresetAutoScale is an alias for CPUProfileAuto.
	PresetAutoScale = CPUProfileAuto
	// PresetHighPerf is an alias for CPUProfileHighPerf.
	PresetHighPerf = CPUProfileHighPerf
	// PresetLowCPU is an alias for CPUProfileLowCPU.
	PresetLowCPU = CPUProfileLowCPU
)

// ResolveAdaptivePullConcurrency resolves the worker concurrency count based on the CPU profile
// and any user-specified overrides.
func ResolveAdaptivePullConcurrency(profile CPUProfile, userSpecified int) int {
	if userSpecified > 0 {
		return userSpecified
	}
	cores := runtime.NumCPU()
	if profile == CPUProfileHighPerf {
		return resolveHighPerfConcurrency(cores)
	}
	if profile == CPUProfileLowCPU {
		return resolveLowCPUConcurrency(cores)
	}
	return resolveAutoProfileConcurrency(cores)
}

// ResolveAdaptiveConcurrency delegates to ResolveAdaptivePullConcurrency.
func ResolveAdaptiveConcurrency(preset ConcurrencyPreset, userLimit int) int {
	return ResolveAdaptivePullConcurrency(preset, userLimit)
}

// ResolveBoostedPullConcurrency dynamically scales concurrency when running locally under idle conditions.
// If isLocalOnly is true and both CPU and RAM are under threshold, scales max workers up to max(runtime.NumCPU(), 6).
// Otherwise, falls back to standard adaptive auto profile.
func ResolveBoostedPullConcurrency(isLocalOnly bool, cpuPercent, ramPercent, threshold int) int {
	if !IsBoostEligible(isLocalOnly, cpuPercent, ramPercent, threshold) {
		return ResolveAdaptivePullConcurrency(CPUProfileAuto, 0)
	}
	cores := runtime.NumCPU()
	if cores < 6 {
		return 6
	}
	return cores
}

// ResolveBoostedWorkerHands returns the boosted worker count and hands per worker.
// When resource usage qualifies for boost, returns (max(runtime.NumCPU(), 6), 2).
func ResolveBoostedWorkerHands(isLocalOnly bool, cpuPercent, ramPercent, threshold int) (int, int) {
	if !IsBoostEligible(isLocalOnly, cpuPercent, ramPercent, threshold) {
		return ResolveAdaptivePullConcurrency(CPUProfileAuto, 0), 1
	}
	workers := ResolveBoostedPullConcurrency(isLocalOnly, cpuPercent, ramPercent, threshold)
	return workers, 2
}

// IsBoostEligible checks whether the resource conditions qualify for adaptive concurrency boost.
func IsBoostEligible(isLocalOnly bool, cpuPercent, ramPercent, threshold int) bool {
	if !isLocalOnly {
		return false
	}
	limit := threshold
	if limit <= 0 {
		limit = 30
	}
	return cpuPercent < limit && ramPercent < limit
}

func resolveHighPerfConcurrency(cores int) int {
	if cores < 1 {
		return 1
	}
	scaled := cores * 3
	if scaled < 12 {
		scaled = 12
	}
	if scaled > 24 {
		return 24
	}
	return scaled
}

func resolveLowCPUConcurrency(cores int) int {
	quarter := cores / 4
	if quarter < 1 {
		return 1
	}
	if quarter > 2 {
		return 2
	}
	return quarter
}

func resolveAutoProfileConcurrency(cores int) int {
	if cores < 1 {
		return 1
	}
	scaled := cores * 2
	if scaled < 8 {
		scaled = 8
	}
	if scaled > 20 {
		return 20
	}
	return scaled
}
