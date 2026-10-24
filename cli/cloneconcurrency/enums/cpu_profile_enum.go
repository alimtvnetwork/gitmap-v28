package enums

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"strings"
)

// CPUProfileType defines the CPU auto-scaling profile for pull operations.
type CPUProfileType string

const (
	// CPUProfileAuto dynamically scales concurrency based on available CPU logical cores.
	CPUProfileAuto CPUProfileType = "auto"
	// CPUProfileHighPerf scales concurrency aggressively for high-performance low-pressure environments.
	CPUProfileHighPerf CPUProfileType = "high-perf"
	// CPUProfileLowCPU throttles concurrency conservatively for resource-constrained systems.
	CPUProfileLowCPU CPUProfileType = "low-cpu"
)

// IsValid reports whether the value matches a known variant.
func (c CPUProfileType) IsValid() bool {
	switch c {
	case CPUProfileAuto, CPUProfileHighPerf, CPUProfileLowCPU:
		return true
	default:
		return false
	}
}

// String returns the string representation of the value.
func (c CPUProfileType) String() string {
	return string(c)
}

// ParseCPUProfileType parses a raw string into a validated CPUProfileType.
func ParseCPUProfileType(s string) (CPUProfileType, *apperror.AppError) {
	trimmed := strings.TrimSpace(s)
	for _, val := range AllCPUProfileTypes() {
		if trimmed == val.String() || strings.EqualFold(trimmed, val.String()) {
			return val, nil
		}
	}
	switch strings.ToLower(trimmed) {
	case "cpuprofileauto":
		return CPUProfileAuto, nil
	case "auto":
		return CPUProfileAuto, nil
	case "cpuprofilehighperf":
		return CPUProfileHighPerf, nil
	case "highperf":
		return CPUProfileHighPerf, nil
	case "cpuprofilelowcpu":
		return CPUProfileLowCPU, nil
	case "lowcpu":
		return CPUProfileLowCPU, nil
	}
	return "", apperror.NewValidation("ParseCPUProfileType", "E1000", fmt.Sprintf("invalid cpuprofiletype: %q", s)).WithContext("value", s)
}

// AllCPUProfileTypes returns all valid variants.
func AllCPUProfileTypes() []CPUProfileType {
	return []CPUProfileType{
		CPUProfileAuto,
		CPUProfileHighPerf,
		CPUProfileLowCPU,
	}
}

// IsAuto reports whether the value is CPUProfileAuto.
func (c CPUProfileType) IsAuto() bool {
	return c == CPUProfileAuto
}

// IsHighPerf reports whether the value is CPUProfileHighPerf.
func (c CPUProfileType) IsHighPerf() bool {
	return c == CPUProfileHighPerf
}

// IsLowCPU reports whether the value is CPUProfileLowCPU.
func (c CPUProfileType) IsLowCPU() bool {
	return c == CPUProfileLowCPU
}

// HasProfile reports whether the value is non-empty and valid.
func (c CPUProfileType) HasProfile() bool {
	return c.IsValid()
}
