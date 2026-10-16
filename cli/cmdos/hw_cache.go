package cmdos

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/charmbracelet/lipgloss"
)

// HardwareMetricsTTL defines the duration hardware telemetry is cached across processes.
const HardwareMetricsTTL = 5 * time.Second

var (
	hwCacheLock   sync.RWMutex
	hwMemoryCache *HardwareMetrics

	hwStyleCardTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	hwStyleHeader      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	hwStyleDim         = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	hwStyleSuccess     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	hwStyleWarn        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c"))
	hwStyleDanger      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
	hwStyleCachedBadge = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Italic(true)
	hwStyleLiveBadge   = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Italic(true)
)

func isMemoryCacheValid() bool {
	hasCache := hwMemoryCache != nil
	if hasCache {
		isFresh := time.Since(hwMemoryCache.Timestamp) < HardwareMetricsTTL
		return isFresh
	}
	return false
}

// GetHardwareMetrics returns cached hardware metrics or triggers a fresh platform probe.
func GetHardwareMetrics() (HardwareMetrics, error) {
	hwCacheLock.RLock()
	isMemValid := isMemoryCacheValid()
	if isMemValid {
		cached := *hwMemoryCache
		cached.IsCached = true
		hwCacheLock.RUnlock()
		return cached, nil
	}
	hwCacheLock.RUnlock()

	hwCacheLock.Lock()
	defer hwCacheLock.Unlock()

	isMemValidSecond := isMemoryCacheValid()
	if isMemValidSecond {
		cached := *hwMemoryCache
		cached.IsCached = true
		return cached, nil
	}

	diskMetrics, hasDiskCache := readDiskHWCache()
	if hasDiskCache {
		hwMemoryCache = &diskMetrics
		cached := diskMetrics
		cached.IsCached = true
		return cached, nil
	}

	start := time.Now()
	fresh, err := executePlatformHWProbe()
	hasProbeErr := err != nil
	if hasProbeErr {
		return HardwareMetrics{}, apperror.WrapSimple(err, "GetHardwareMetrics.executePlatformHWProbe")
	}
	fresh.Timestamp = time.Now()
	fresh.DurationMs = time.Since(start).Milliseconds()
	fresh.IsCached = false

	hwMemoryCache = &fresh
	_ = writeDiskHWCache(fresh)

	return fresh, nil
}

func resolveHWCacheFilePath() string {
	return filepath.Join(os.TempDir(), "gitmap_hw_metrics_cache.json")
}

func readDiskHWCache() (HardwareMetrics, bool) {
	cachePath := resolveHWCacheFilePath()
	data, err := os.ReadFile(cachePath)
	hasReadErr := err != nil
	if hasReadErr {
		return HardwareMetrics{}, false
	}

	var env HardwareCacheEnvelope
	unmarshalErr := json.Unmarshal(data, &env)
	hasUnmarshalErr := unmarshalErr != nil
	if hasUnmarshalErr {
		return HardwareMetrics{}, false
	}

	isExpired := time.Since(env.CachedAt) >= HardwareMetricsTTL
	if isExpired {
		return HardwareMetrics{}, false
	}

	metrics := env.Metrics
	metrics.Timestamp = env.CachedAt
	return metrics, true
}

func writeDiskHWCache(m HardwareMetrics) error {
	cachePath := resolveHWCacheFilePath()
	env := HardwareCacheEnvelope{
		CachedAt:   m.Timestamp,
		TTLSeconds: int(HardwareMetricsTTL.Seconds()),
		Metrics:    m,
	}
	data, err := json.Marshal(env)
	hasMarshalErr := err != nil
	if hasMarshalErr {
		return apperror.WrapSimple(err, "writeDiskHWCache.marshal")
	}

	tmpPath := fmt.Sprintf("%s.%d.tmp", cachePath, time.Now().UnixNano())
	writeErr := os.WriteFile(tmpPath, data, 0600)
	hasWriteErr := writeErr != nil
	if hasWriteErr {
		return apperror.WrapSimple(writeErr, "writeDiskHWCache.writeFile")
	}

	renameErr := os.Rename(tmpPath, cachePath)
	hasRenameErr := renameErr != nil
	if hasRenameErr {
		_ = os.Remove(tmpPath)
		return apperror.WrapSimple(renameErr, "writeDiskHWCache.rename")
	}
	return nil
}

func printJSON(v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	hasErr := err != nil
	if hasErr {
		return apperror.WrapSimple(err, "printJSON")
	}
	fmt.Println(string(data))
	return nil
}

func runCPUCLI(args []string) error {
	checkHelp("cpu", args)
	isJSON := hasOSJSONArg(args)
	metrics, err := GetHardwareMetrics()
	hasErr := err != nil
	if hasErr {
		return err
	}
	if isJSON {
		return printJSON(metrics.CPU)
	}
	renderCPUView(metrics.CPU, metrics.IsCached, metrics.DurationMs)
	return nil
}

func runRAMCLI(args []string) error {
	checkHelp("ram", args)
	isJSON := hasOSJSONArg(args)
	metrics, err := GetHardwareMetrics()
	hasErr := err != nil
	if hasErr {
		return err
	}
	if isJSON {
		return printJSON(metrics.RAM)
	}
	renderRAMView(metrics.RAM, metrics.IsCached, metrics.DurationMs)
	return nil
}

func runDiskCLI(args []string) error {
	checkHelp("disk", args)
	isJSON := hasOSJSONArg(args)
	metrics, err := GetHardwareMetrics()
	hasErr := err != nil
	if hasErr {
		return err
	}
	if isJSON {
		return printJSON(metrics.Disks)
	}
	renderDiskView(metrics.Disks, metrics.IsCached, metrics.DurationMs)
	return nil
}

func renderProgressBar(pct int, width int) string {
	isUnder := pct < 0
	if isUnder {
		pct = 0
	}
	isOver := pct > 100
	if isOver {
		pct = 100
	}
	filled := (pct * width) / 100
	empty := width - filled
	fillChar := strings.Repeat("█", filled)
	emptyChar := strings.Repeat("░", empty)
	return fillChar + emptyChar
}

func resolvePercentStyle(pct int) lipgloss.Style {
	isLow := pct <= 50
	if isLow {
		return hwStyleSuccess
	}
	isMedium := pct <= 80
	if isMedium {
		return hwStyleWarn
	}
	return hwStyleDanger
}

func renderCPUView(cpu CPUMetrics, isCached bool, durationMs int64) {
	badge := hwStyleLiveBadge.Render(fmt.Sprintf("[live %dms]", durationMs))
	if isCached {
		badge = hwStyleCachedBadge.Render("[cached <1ms]")
	}

	pctStyle := resolvePercentStyle(cpu.LoadPercent)
	bar := renderProgressBar(cpu.LoadPercent, 20)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", hwStyleCardTitle.Render("CPU TELEMETRY"), badge))
	sb.WriteString(fmt.Sprintf("  %s\n", hwStyleDim.Render(strings.Repeat("─", 52))))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", hwStyleHeader.Render("Model:"), cpu.ModelName))
	sb.WriteString(fmt.Sprintf("  %-16s %s\n", hwStyleHeader.Render("Architecture:"), cpu.Architecture))
	sb.WriteString(fmt.Sprintf("  %-16s %d logical cores\n", hwStyleHeader.Render("Cores:"), cpu.Cores))
	sb.WriteString(fmt.Sprintf("  %-16s %s %s\n", hwStyleHeader.Render("Current Load:"),
		pctStyle.Render(fmt.Sprintf("%d%%", cpu.LoadPercent)),
		pctStyle.Render(bar)))
	sb.WriteString(fmt.Sprintf("  %s\n\n", hwStyleDim.Render(strings.Repeat("─", 52))))

	fmt.Print(sb.String())
}

func renderRAMView(ram RAMMetrics, isCached bool, durationMs int64) {
	badge := hwStyleLiveBadge.Render(fmt.Sprintf("[live %dms]", durationMs))
	if isCached {
		badge = hwStyleCachedBadge.Render("[cached <1ms]")
	}

	pctStyle := resolvePercentStyle(ram.UsedPercent)
	bar := renderProgressBar(ram.UsedPercent, 20)

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", hwStyleCardTitle.Render("RAM TELEMETRY"), badge))
	sb.WriteString(fmt.Sprintf("  %s\n", hwStyleDim.Render(strings.Repeat("─", 52))))
	sb.WriteString(fmt.Sprintf("  %-16s %d MB (%.1f GB)\n", hwStyleHeader.Render("Total Capacity:"), ram.TotalMB, float64(ram.TotalMB)/1024.0))
	sb.WriteString(fmt.Sprintf("  %-16s %d MB (%.1f GB)\n", hwStyleHeader.Render("Used Memory:"), ram.UsedMB, float64(ram.UsedMB)/1024.0))
	sb.WriteString(fmt.Sprintf("  %-16s %d MB (%.1f GB)\n", hwStyleHeader.Render("Available/Free:"), ram.FreeMB, float64(ram.FreeMB)/1024.0))
	sb.WriteString(fmt.Sprintf("  %-16s %s %s\n", hwStyleHeader.Render("Utilization:"),
		pctStyle.Render(fmt.Sprintf("%d%%", ram.UsedPercent)),
		pctStyle.Render(bar)))
	sb.WriteString(fmt.Sprintf("  %s\n\n", hwStyleDim.Render(strings.Repeat("─", 52))))

	fmt.Print(sb.String())
}

func renderDiskRows(disks []DriveInfo, sb *strings.Builder) {
	sb.WriteString(fmt.Sprintf("  %-12s %-12s %-12s %-12s %-16s\n",
		hwStyleHeader.Render("MOUNT"),
		hwStyleHeader.Render("TOTAL"),
		hwStyleHeader.Render("USED"),
		hwStyleHeader.Render("FREE"),
		hwStyleHeader.Render("USAGE")))
	sb.WriteString(fmt.Sprintf("  %s\n", hwStyleDim.Render(strings.Repeat("─", 68))))

	for _, d := range disks {
		pctStyle := resolvePercentStyle(d.UsedPercent)
		bar := renderProgressBar(d.UsedPercent, 12)
		sb.WriteString(fmt.Sprintf("  %-12s %-12s %-12s %-12s %s %s\n",
			d.MountPoint,
			fmt.Sprintf("%.1f GB", d.TotalGB),
			fmt.Sprintf("%.1f GB", d.UsedGB),
			fmt.Sprintf("%.1f GB", d.FreeGB),
			pctStyle.Render(fmt.Sprintf("%3d%%", d.UsedPercent)),
			pctStyle.Render(bar)))
	}
}

func renderDiskView(disks []DriveInfo, isCached bool, durationMs int64) {
	badge := hwStyleLiveBadge.Render(fmt.Sprintf("[live %dms]", durationMs))
	if isCached {
		badge = hwStyleCachedBadge.Render("[cached <1ms]")
	}

	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s %s\n", hwStyleCardTitle.Render("DISK STORAGE TELEMETRY"), badge))
	sb.WriteString(fmt.Sprintf("  %s\n", hwStyleDim.Render(strings.Repeat("─", 68))))

	hasDrives := len(disks) > 0
	if hasDrives {
		renderDiskRows(disks, &sb)
	}
	isEmptyDrives := len(disks) == 0
	if isEmptyDrives {
		sb.WriteString(fmt.Sprintf("  %s\n", hwStyleDim.Render("No active drive volumes discovered.")))
	}

	sb.WriteString(fmt.Sprintf("  %s\n\n", hwStyleDim.Render(strings.Repeat("─", 68))))
	fmt.Print(sb.String())
}
