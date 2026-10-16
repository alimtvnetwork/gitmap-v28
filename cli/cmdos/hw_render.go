package cmdos

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	hwStyleCardTitle   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	hwStyleHeader      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	hwStyleDim         = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	hwStyleSuccess     = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	hwStyleWarn        = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c"))
	hwStyleDanger      = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ff5555"))
	hwStyleCachedBadge = lipgloss.NewStyle().Foreground(lipgloss.Color("#bd93f9")).Italic(true)
	hwStyleLiveBadge   = lipgloss.NewStyle().Foreground(lipgloss.Color("#50fa7b")).Italic(true)
)

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
