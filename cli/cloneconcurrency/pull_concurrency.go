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

// ResolveBoostedPullConcurrency dynamically scales concurrency when running locally under idle conditions.
// If isLocalOnly is true and both CPU and RAM are under threshold, scales max workers up to max(runtime.NumCPU(), 6).
// Otherwise, falls back to standard adaptive auto profile.
func ResolveBoostedPullConcurrency(isLocalOnly bool, cpuPercent, ramPercent, threshold int) int {
	if isLocalOnly {
		if threshold <= 0 {
			threshold = 30
		}
		isCPUIdle := cpuPercent < threshold
		isRAMIdle := ramPercent < threshold
		if isCPUIdle && isRAMIdle {
			cores := runtime.NumCPU()
			if cores < 6 {
				return 6
			}
			return cores
		}
	}
	return ResolveAdaptivePullConcurrency(CPUProfileAuto, 0)
}

// ResolveBoostedWorkerHands returns the boosted worker count and hands per worker.
// When resource usage qualifies for boost, returns (max(runtime.NumCPU(), 6), 2).
func ResolveBoostedWorkerHands(isLocalOnly bool, cpuPercent, ramPercent, threshold int) (int, int) {
	if isLocalOnly {
		if threshold <= 0 {
			threshold = 30
		}
		isCPUIdle := cpuPercent < threshold
		isRAMIdle := ramPercent < threshold
		if isCPUIdle && isRAMIdle {
			workers := ResolveBoostedPullConcurrency(isLocalOnly, cpuPercent, ramPercent, threshold)
			return workers, 2
		}
	}
	return ResolveAdaptivePullConcurrency(CPUProfileAuto, 0), 1
}

// IsBoostEligible checks whether the resource conditions qualify for adaptive concurrency boost.
func IsBoostEligible(isLocalOnly bool, cpuPercent, ramPercent, threshold int) bool {
	if isLocalOnly {
		if threshold <= 0 {
			threshold = 30
		}
		isCPUIdle := cpuPercent < threshold
		isRAMIdle := ramPercent < threshold
		return isCPUIdle && isRAMIdle
	}
	return false
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
