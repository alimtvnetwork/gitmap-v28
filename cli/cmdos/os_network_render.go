// Package cmdos — os_network_render.go formats and renders network adapter telemetry tables and help.
package cmdos

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

var (
	netStyleCardTitle = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#bd93f9"))
	netStyleHeader    = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#8be9fd"))
	netStyleDim       = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	netStyleSuccess   = lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#50fa7b"))
	netStyleDown      = lipgloss.NewStyle().Foreground(lipgloss.Color("#6272a4"))
	netStyleTxRx      = lipgloss.NewStyle().Foreground(lipgloss.Color("#f1fa8c"))
)

func renderNetworkReportTable(report NetworkReport, opts NetworkOptions) {
	adapters := filterNetworkAdapters(report.Adapters, opts.IsAll)

	var sb strings.Builder
	sb.WriteString("\n")
	title := fmt.Sprintf("NETWORK ADAPTERS (Host: %s, Platform: %s)", report.Hostname, report.Platform)
	sb.WriteString(fmt.Sprintf("  %s\n", netStyleCardTitle.Render(title)))
	sb.WriteString(fmt.Sprintf("  %s\n", netStyleDim.Render(strings.Repeat("─", 102))))
	sb.WriteString(fmt.Sprintf("  %-22s %-8s %-18s %-20s %-12s %s\n",
		netStyleHeader.Render("ADAPTER"),
		netStyleHeader.Render("STATUS"),
		netStyleHeader.Render("MAC ADDRESS"),
		netStyleHeader.Render("IPV4 / IPV6"),
		netStyleHeader.Render("SPEED"),
		netStyleHeader.Render("TRAFFIC (TX / RX)"),
	))
	sb.WriteString(fmt.Sprintf("  %s\n", netStyleDim.Render(strings.Repeat("─", 102))))

	for _, a := range adapters {
		statusStyled := netStyleDown.Render("DOWN")
		if a.IsUp {
			statusStyled = netStyleSuccess.Render("UP")
		}

		ipStr := resolveAdapterIPString(a)
		speedStr := formatAdapterSpeed(a.LinkSpeed)
		tx := formatNetworkBytes(a.BytesSent)
		rx := formatNetworkBytes(a.BytesReceived)
		traffic := netStyleTxRx.Render(fmt.Sprintf("▲ %s", tx)) + " / " + netStyleTxRx.Render(fmt.Sprintf("▼ %s", rx))

		mac := fallbackDash(a.MACAddress)
		sb.WriteString(fmt.Sprintf("  %-22s %-8s %-18s %-20s %-12s %s\n",
			truncateText(a.Name, 22),
			statusStyled,
			truncateText(mac, 18),
			truncateText(ipStr, 20),
			speedStr,
			traffic,
		))
	}

	sb.WriteString(fmt.Sprintf("  %s\n", netStyleDim.Render(strings.Repeat("─", 102))))
	footer := fmt.Sprintf("Total Transmitted: %s  |  Total Received: %s  |  Timestamp: %s",
		formatNetworkBytes(report.TotalBytesSent),
		formatNetworkBytes(report.TotalBytesRecv),
		report.CollectedAt,
	)
	sb.WriteString(fmt.Sprintf("  %s\n\n", netStyleDim.Render(footer)))

	fmt.Print(sb.String())
}

func filterNetworkAdapters(adapters []NetworkAdapterInfo, isAll bool) []NetworkAdapterInfo {
	if isAll {
		return adapters
	}

	var filtered []NetworkAdapterInfo
	for _, a := range adapters {
		isIncluded := evaluateAdapterInclusion(a)
		if isIncluded {
			filtered = append(filtered, a)
		}
	}

	hasAdapters := len(filtered) > 0
	if hasAdapters {
		return filtered
	}
	return adapters
}

func evaluateAdapterInclusion(a NetworkAdapterInfo) bool {
	isLoopback := checkIsLoopbackAdapter(a)
	if isLoopback {
		return false
	}
	if a.IsUp {
		return true
	}
	hasIP := len(a.IPv4) > 0 || len(a.IPv4Addresses) > 0 || len(a.IPv6) > 0 || len(a.IPv6Addresses) > 0
	if hasIP {
		return true
	}
	hasTraffic := a.BytesSent > 0 || a.BytesReceived > 0
	return hasTraffic
}

func checkIsLoopbackAdapter(a NetworkAdapterInfo) bool {
	nameLower := strings.ToLower(a.Name)
	isNameLoopback := nameLower == "lo" || strings.HasPrefix(nameLower, "loopback")
	if isNameLoopback {
		return true
	}
	descLower := strings.ToLower(a.Description)
	isDescLoopback := strings.Contains(descLower, "loopback")
	if isDescLoopback {
		return true
	}
	isLoopbackIP := len(a.IPv4) == 1 && a.IPv4[0] == "127.0.0.1"
	return isLoopbackIP
}

func resolveAdapterIPString(a NetworkAdapterInfo) string {
	hasIPv4 := len(a.IPv4) > 0
	if hasIPv4 {
		return strings.Join(a.IPv4, ", ")
	}
	hasIPv4Addrs := len(a.IPv4Addresses) > 0
	if hasIPv4Addrs {
		return strings.Join(a.IPv4Addresses, ", ")
	}
	hasIPv6 := len(a.IPv6) > 0
	if hasIPv6 {
		return a.IPv6[0]
	}
	hasIPv6Addrs := len(a.IPv6Addresses) > 0
	if hasIPv6Addrs {
		return a.IPv6Addresses[0]
	}
	return "-"
}

func formatAdapterSpeed(speed string) string {
	trimmed := strings.TrimSpace(speed)
	isEmpty := len(trimmed) == 0 || trimmed == "-1" || trimmed == "0 bps"
	if isEmpty {
		return "-"
	}
	return trimmed
}

func renderNetworkHelp() {
	var sb strings.Builder
	sb.WriteString("\n")
	sb.WriteString(fmt.Sprintf("  %s\n", netStyleCardTitle.Render("NETWORK ADAPTER TELEMETRY (gitmap os network / gitmap network)")))
	sb.WriteString(fmt.Sprintf("  %s\n", netStyleDim.Render(strings.Repeat("─", 72))))
	sb.WriteString("  Usage:\n")
	sb.WriteString("    gitmap network [--json] [--all]\n")
	sb.WriteString("    gitmap os network [--json] [--all]\n\n")
	sb.WriteString("  Flags:\n")
	sb.WriteString(fmt.Sprintf("    %-16s %s\n", netStyleHeader.Render("-j, --json"), "Output network telemetry report as structured JSON"))
	sb.WriteString(fmt.Sprintf("    %-16s %s\n", netStyleHeader.Render("-a, --all"), "Display all interfaces including inactive and loopback"))
	sb.WriteString(fmt.Sprintf("    %-16s %s\n", netStyleHeader.Render("-h, --help"), "Display this help menu"))
	sb.WriteString(fmt.Sprintf("  %s\n\n", netStyleDim.Render(strings.Repeat("─", 72))))
	fmt.Print(sb.String())
}

func formatNetworkBytes(b uint64) string {
	const unit = 1024
	isSmall := b < unit
	if isSmall {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	isValidExp := exp < len(units)
	if isValidExp {
		return fmt.Sprintf("%.1f %s", float64(b)/float64(div), units[exp])
	}
	return fmt.Sprintf("%d B", b)
}

func truncateText(s string, maxLen int) string {
	isWithinLen := len(s) <= maxLen
	if isWithinLen {
		return s
	}
	isVeryShort := maxLen <= 3
	if isVeryShort {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
