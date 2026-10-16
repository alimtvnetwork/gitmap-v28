//go:build !windows && !linux

package cmdos

import "runtime"

func executePlatformHWProbe() (HardwareMetrics, error) {
	return HardwareMetrics{
		CPU: CPUMetrics{
			ModelName:    "Generic POSIX Host Processor",
			Cores:        runtime.NumCPU(),
			LoadPercent:  0,
			Architecture: runtime.GOARCH,
		},
		RAM: RAMMetrics{
			TotalMB:     0,
			UsedMB:      0,
			FreeMB:      0,
			UsedPercent: 0,
		},
		Disks: []DriveInfo{},
	}, nil
}
