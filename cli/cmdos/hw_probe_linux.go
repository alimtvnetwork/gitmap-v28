//go:build linux

package cmdos

import (
	"bufio"
	"os"
	"runtime"
	"strconv"
	"strings"
	"syscall"
	"time"
)

func executePlatformHWProbe() (HardwareMetrics, error) {
	cpuMetrics := probeLinuxCPU()
	ramMetrics := probeLinuxRAM()
	diskMetrics := probeLinuxDisks()

	return HardwareMetrics{
		CPU:   cpuMetrics,
		RAM:   ramMetrics,
		Disks: diskMetrics,
	}, nil
}

func readProcStatTicks() (idle uint64, total uint64, hasRead bool) {
	file, err := os.Open("/proc/stat")
	hasErr := err != nil
	if hasErr {
		return 0, 0, false
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		isCPU := strings.HasPrefix(line, "cpu ")
		if isCPU {
			fields := strings.Fields(line)
			hasEnoughFields := len(fields) >= 5
			if hasEnoughFields {
				var sum uint64
				for i := 1; i < len(fields); i++ {
					val, _ := strconv.ParseUint(fields[i], 10, 64)
					sum += val
				}
				idleVal, _ := strconv.ParseUint(fields[4], 10, 64)
				return idleVal, sum, true
			}
		}
	}
	return 0, 0, false
}

func probeLinuxCPU() CPUMetrics {
	idle1, total1, hasFirst := readProcStatTicks()
	var loadPercent int
	if hasFirst {
		time.Sleep(100 * time.Millisecond)
		idle2, total2, hasSecond := readProcStatTicks()
		if hasSecond {
			deltaTotal := total2 - total1
			deltaIdle := idle2 - idle1
			hasDelta := deltaTotal > 0 && deltaTotal >= deltaIdle
			if hasDelta {
				deltaBusy := deltaTotal - deltaIdle
				loadPercent = int((deltaBusy * 100) / deltaTotal)
			}
		}
	}

	modelName := probeLinuxCPUModel()
	isEmptyModel := len(modelName) == 0
	if isEmptyModel {
		modelName = "Linux Host Processor"
	}

	return CPUMetrics{
		ModelName:    modelName,
		Cores:        runtime.NumCPU(),
		LoadPercent:  loadPercent,
		Architecture: runtime.GOARCH,
	}
}

func probeLinuxCPUModel() string {
	file, err := os.Open("/proc/cpuinfo")
	hasErr := err != nil
	if hasErr {
		return ""
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		isModel := strings.HasPrefix(line, "model name")
		if isModel {
			parts := strings.Split(line, ":")
			hasParts := len(parts) >= 2
			if hasParts {
				return strings.TrimSpace(parts[1])
			}
		}
	}
	return ""
}

func probeLinuxRAM() RAMMetrics {
	file, err := os.Open("/proc/meminfo")
	hasErr := err != nil
	if hasErr {
		return RAMMetrics{}
	}
	defer file.Close()

	var totalKB, availableKB, freeKB int
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fields := strings.Fields(line)
		hasFields := len(fields) >= 2
		if hasFields {
			val, _ := strconv.Atoi(fields[1])
			switch fields[0] {
			case "MemTotal:":
				totalKB = val
			case "MemAvailable:":
				availableKB = val
			case "MemFree:":
				freeKB = val
			}
		}
	}

	totalMB := totalKB / 1024
	availMB := availableKB / 1024
	isAvailZero := availMB == 0
	if isAvailZero {
		availMB = freeKB / 1024
	}
	usedMB := totalMB - availMB
	usedPercent := 0
	hasTotal := totalMB > 0
	if hasTotal {
		usedPercent = int(float64(usedMB) / float64(totalMB) * 100.0)
	}

	return RAMMetrics{
		TotalMB:     totalMB,
		UsedMB:      usedMB,
		FreeMB:      availMB,
		UsedPercent: usedPercent,
	}
}

func probeLinuxDisks() []DriveInfo {
	var stat syscall.Statfs_t
	err := syscall.Statfs("/", &stat)
	hasErr := err != nil
	if hasErr {
		return []DriveInfo{}
	}

	totalBytes := stat.Blocks * uint64(stat.Bsize)
	freeBytes := stat.Bavail * uint64(stat.Bsize)
	usedBytes := totalBytes - freeBytes

	totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
	usedGB := float64(usedBytes) / (1024 * 1024 * 1024)
	freeGB := float64(freeBytes) / (1024 * 1024 * 1024)
	usedPercent := 0
	hasBytes := totalBytes > 0
	if hasBytes {
		usedPercent = int(float64(usedBytes) / float64(totalBytes) * 100.0)
	}

	return []DriveInfo{
		{
			MountPoint:  "/",
			TotalGB:     totalGB,
			UsedGB:      usedGB,
			FreeGB:      freeGB,
			UsedPercent: usedPercent,
		},
	}
}
