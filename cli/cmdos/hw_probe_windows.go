//go:build windows

package cmdos

import (
	"context"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"golang.org/x/sys/windows"
)

func executePlatformHWProbe() (HardwareMetrics, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	cmdScript := `$cpu = (Get-CimInstance Win32_Processor | Measure-Object -Property LoadPercentage -Average).Average; $os = Get-CimInstance Win32_OperatingSystem; $tot = [math]::Round($os.TotalVisibleMemorySize/1024); $free = [math]::Round($os.FreePhysicalMemory/1024); $used = $tot - $free; $ram = [math]::Round(($used/$tot)*100); Write-Output ("{0},{1},{2},{3}" -f [int]$cpu, [int]$ram, [int]$tot, [int]$used)`

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", cmdScript)
	outBytes, err := cmd.Output()
	if err != nil {
		return HardwareMetrics{}, apperror.WrapSimple(err, "executePlatformHWProbe.powershell")
	}

	lines := strings.Split(strings.TrimSpace(string(outBytes)), "\n")
	line := strings.TrimSpace(lines[len(lines)-1])
	parts := strings.Split(line, ",")
	if len(parts) < 4 {
		return HardwareMetrics{}, apperror.NewValidation("executePlatformHWProbe", "E_HW_CIM_PARSE_FAILED", "invalid CIM probe output: "+line)
	}

	cpuPercent, _ := strconv.Atoi(strings.TrimSpace(parts[0]))
	ramPercent, _ := strconv.Atoi(strings.TrimSpace(parts[1]))
	totalRAM, _ := strconv.Atoi(strings.TrimSpace(parts[2]))
	usedRAM, _ := strconv.Atoi(strings.TrimSpace(parts[3]))

	modelName := os.Getenv("PROCESSOR_IDENTIFIER")
	isEmptyModel := len(modelName) == 0
	if isEmptyModel {
		modelName = "Windows Host Processor"
	}

	disks := probeWindowsDisks()

	return HardwareMetrics{
		CPU: CPUMetrics{
			ModelName:    modelName,
			Cores:        runtime.NumCPU(),
			LoadPercent:  cpuPercent,
			Architecture: runtime.GOARCH,
		},
		RAM: RAMMetrics{
			TotalMB:     totalRAM,
			UsedMB:      usedRAM,
			FreeMB:      totalRAM - usedRAM,
			UsedPercent: ramPercent,
		},
		Disks: disks,
	}, nil
}

func probeWindowsDisks() []DriveInfo {
	var drives []DriveInfo
	bitmask, err := windows.GetLogicalDrives()
	if err != nil {
		return probeFallbackDrive()
	}

	for i := 0; i < 26; i++ {
		hasDrive := (bitmask & (1 << uint(i))) != 0
		if hasDrive {
			driveRoot := fmt.Sprintf("%c:\\", 'A'+i)
			volPtr, err := windows.UTF16PtrFromString(driveRoot)
			if err != nil {
				continue
			}
			var freeBytes, totalBytes, totalFreeBytes uint64
			errFree := windows.GetDiskFreeSpaceEx(volPtr, &freeBytes, &totalBytes, &totalFreeBytes)
			hasValidSpace := errFree == nil && totalBytes > 0
			if hasValidSpace {
				usedBytes := totalBytes - freeBytes
				totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
				usedGB := float64(usedBytes) / (1024 * 1024 * 1024)
				freeGB := float64(freeBytes) / (1024 * 1024 * 1024)
				usedPct := int((float64(usedBytes) / float64(totalBytes)) * 100.0)

				drives = append(drives, DriveInfo{
					MountPoint:  driveRoot,
					TotalGB:     totalGB,
					UsedGB:      usedGB,
					FreeGB:      freeGB,
					UsedPercent: usedPct,
				})
			}
		}
	}

	hasAny := len(drives) > 0
	if hasAny {
		return drives
	}
	return probeFallbackDrive()
}

func probeFallbackDrive() []DriveInfo {
	driveRoot := "C:\\"
	volPtr, err := windows.UTF16PtrFromString(driveRoot)
	if err != nil {
		return []DriveInfo{}
	}
	var freeBytes, totalBytes, totalFreeBytes uint64
	errFree := windows.GetDiskFreeSpaceEx(volPtr, &freeBytes, &totalBytes, &totalFreeBytes)
	hasValidSpace := errFree == nil && totalBytes > 0
	if hasValidSpace {
		usedBytes := totalBytes - freeBytes
		totalGB := float64(totalBytes) / (1024 * 1024 * 1024)
		usedGB := float64(usedBytes) / (1024 * 1024 * 1024)
		freeGB := float64(freeBytes) / (1024 * 1024 * 1024)
		usedPct := int((float64(usedBytes) / float64(totalBytes)) * 100.0)

		return []DriveInfo{
			{
				MountPoint:  driveRoot,
				TotalGB:     totalGB,
				UsedGB:      usedGB,
				FreeGB:      freeGB,
				UsedPercent: usedPct,
			},
		}
	}
	return []DriveInfo{}
}
