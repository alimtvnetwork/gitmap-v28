package cmdos

import "time"

// HardwareMetrics encapsulates consolidated hardware status telemetry.
type HardwareMetrics struct {
	CPU        CPUMetrics  `json:"cpu"`
	RAM        RAMMetrics  `json:"ram"`
	Disks      []DriveInfo `json:"disks"`
	Timestamp  time.Time   `json:"timestamp"`
	DurationMs int64       `json:"durationMs"`
	IsCached   bool        `json:"isCached"`
}

// CPUMetrics contains processor load and core information.
type CPUMetrics struct {
	ModelName    string `json:"modelName"`
	Cores        int    `json:"cores"`
	LoadPercent  int    `json:"loadPercent"`
	Architecture string `json:"architecture"`
}

// RAMMetrics contains physical memory capacity and utilization.
type RAMMetrics struct {
	TotalMB     int `json:"totalMb"`
	UsedMB      int `json:"usedMb"`
	FreeMB      int `json:"freeMb"`
	UsedPercent int `json:"usedPercent"`
}

// DriveInfo represents volume storage capacity and usage.
type DriveInfo struct {
	MountPoint  string  `json:"mountPoint"`
	TotalGB     float64 `json:"totalGb"`
	UsedGB      float64 `json:"usedGb"`
	FreeGB      float64 `json:"freeGb"`
	UsedPercent int     `json:"usedPercent"`
}

// HardwareCacheEnvelope represents the serialized JSON cache envelope.
type HardwareCacheEnvelope struct {
	CachedAt   time.Time       `json:"cachedAt"`
	TTLSeconds int             `json:"ttlSeconds"`
	Metrics    HardwareMetrics `json:"metrics"`
}

// HardwareProbeResult provides a flattened hardware probe telemetry structure.
type HardwareProbeResult struct {
	CPUPercent      int     `json:"cpuPercent"`
	RAMPercent      int     `json:"ramPercent"`
	TotalMB         int     `json:"totalMb"`
	UsedMB          int     `json:"usedMb"`
	FreeMB          int     `json:"freeMb"`
	DiskTotalGB     float64 `json:"diskTotalGb"`
	DiskFreeGB      float64 `json:"diskFreeGb"`
	DiskUsedGB      float64 `json:"diskUsedGb"`
	CollectedAtUnix int64   `json:"collectedAtUnix"`
	IsCached        bool    `json:"isCached"`
}

// ToProbeResult converts HardwareMetrics to a flattened HardwareProbeResult.
func (m HardwareMetrics) ToProbeResult() HardwareProbeResult {
	var totalDiskGB, freeDiskGB, usedDiskGB float64
	for _, d := range m.Disks {
		totalDiskGB += d.TotalGB
		freeDiskGB += d.FreeGB
		usedDiskGB += d.UsedGB
	}
	return HardwareProbeResult{
		CPUPercent:      m.CPU.LoadPercent,
		RAMPercent:      m.RAM.UsedPercent,
		TotalMB:         m.RAM.TotalMB,
		UsedMB:          m.RAM.UsedMB,
		FreeMB:          m.RAM.FreeMB,
		DiskTotalGB:     totalDiskGB,
		DiskFreeGB:      freeDiskGB,
		DiskUsedGB:      usedDiskGB,
		CollectedAtUnix: m.Timestamp.Unix(),
		IsCached:        m.IsCached,
	}
}
