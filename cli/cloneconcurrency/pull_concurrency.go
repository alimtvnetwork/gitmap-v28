package cloneconcurrency

import (
	"runtime"
)

// CPUProfile defines the CPU auto-scaling profile for pull operations.
type CPUProfile string

const (
	// CPUProfileAuto dynamically scales concurrency based on available CPU logical cores.
	CPUProfileAuto CPUProfile = "auto"
	// CPUProfileHighPerf scales concurrency aggressively for high-performance low-pressure environments.
	CPUProfileHighPerf CPUProfile = "high-perf"
	// CPUProfileLowCPU throttles concurrency conservatively for resource-constrained systems.
	CPUProfileLowCPU CPUProfile = "low-cpu"
)

// ConcurrencyPreset is an alias for CPUProfile for backwards compatibility with component specs.
type ConcurrencyPreset = CPUProfile

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

func resolveHighPerfConcurrency(cores int) int {
	if cores < 1 {
		return 1
	}
	if cores > 12 {
		return 12
	}
	return cores
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
	if cores <= 2 {
		return 1
	}
	if cores <= 4 {
		return 2
	}
	if cores <= 8 {
		return 3
	}
	if cores <= 16 {
		return 4
	}
	return 6
}
