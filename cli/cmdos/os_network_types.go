// Package cmdos — os_network_types.go provides network adapter metrics models, CLI options, and reporting dispatch.
package cmdos

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// NetworkAdapterInfo defines the cross-platform telemetry collected per adapter.
type NetworkAdapterInfo struct {
	Name            string   `json:"name"`
	Description     string   `json:"description,omitempty"`
	Status          string   `json:"status"`
	IsUp            bool     `json:"isUp"`
	IPv4            []string `json:"ipv4"`
	IPv4Addresses   []string `json:"ipv4Addresses,omitempty"`
	IPv6            []string `json:"ipv6,omitempty"`
	IPv6Addresses   []string `json:"ipv6Addresses,omitempty"`
	MACAddress      string   `json:"macAddress"`
	LinkSpeed       string   `json:"linkSpeed,omitempty"`
	BytesSent       uint64   `json:"bytesSent"`
	BytesReceived   uint64   `json:"bytesReceived"`
	PacketsSent     uint64   `json:"packetsSent,omitempty"`
	PacketsReceived uint64   `json:"packetsReceived,omitempty"`
}

// NetworkReport encapsulates the aggregate telemetry report for the active host.
type NetworkReport struct {
	Platform       string               `json:"platform"`
	Hostname       string               `json:"hostname"`
	Adapters       []NetworkAdapterInfo `json:"adapters"`
	TotalBytesSent uint64               `json:"totalBytesSent"`
	TotalBytesRecv uint64               `json:"totalBytesRecv"`
	CollectedAt    string               `json:"collectedAt"`
}

// NetworkAdapterMetrics is a type alias for spec compatibility.
type NetworkAdapterMetrics = NetworkAdapterInfo

// NetworkTelemetryPayload is a type alias for spec compatibility.
type NetworkTelemetryPayload = NetworkReport

// NetworkOptions defines parsed CLI execution flags.
type NetworkOptions struct {
	IsJSON   bool
	IsAll    bool
	IsBypass bool
}

// parseNetworkOptions extracts flags from CLI arguments using positive booleans.
func parseNetworkOptions(args []string) NetworkOptions {
	opts := NetworkOptions{}
	for _, arg := range args {
		normalized := strings.ToLower(strings.TrimSpace(arg))
		if normalized == "--json" || normalized == "-j" || normalized == "json" {
			opts.IsJSON = true
		}
		if normalized == "--all" || normalized == "-a" || normalized == "all" {
			opts.IsAll = true
		}
		if normalized == "--bypass" || normalized == "--no-cache" || normalized == "-b" {
			opts.IsBypass = true
		}
	}
	return opts
}

func isNetworkHelpArg(arg string) bool {
	normalized := strings.ToLower(strings.TrimSpace(arg))
	return normalized == "-h" || normalized == "--help" || normalized == "help"
}

// RunOSNetworkCLI executes the OS network adapter inspection command.
func RunOSNetworkCLI(args []string) error {
	if len(args) > 0 && isNetworkHelpArg(args[0]) {
		renderNetworkHelp()
		return nil
	}

	opts := parseNetworkOptions(args)
	report, err := collectNetworkReport()
	if err != nil {
		return apperror.WrapSimple(err, "collect network telemetry report")
	}

	if opts.IsJSON {
		data, marshalErr := json.MarshalIndent(report, "", "  ")
		if marshalErr != nil {
			return apperror.WrapSimple(marshalErr, "marshal network telemetry json")
		}
		fmt.Println(string(data))
		return nil
	}

	renderNetworkReportTable(report, opts)
	return nil
}

// RunNetworkCLI is an alias for RunOSNetworkCLI for top-level routing.
func RunNetworkCLI(args []string) error {
	return RunOSNetworkCLI(args)
}

func renderNetworkReportTable(report NetworkReport, opts NetworkOptions) {
	fmt.Printf("\n%s● Network Adapters & Interfaces (Host: %s)%s\n", constants.ColorCyan, report.Hostname, constants.ColorReset)
	fmt.Printf("  %-22s %-8s %-18s %-18s %-12s %s\n", "ADAPTER", "STATUS", "MAC ADDRESS", "IPV4", "SPEED", "TRAFFIC (TX / RX)")
	fmt.Printf("  %s\n", strings.Repeat("─", 102))

	for _, a := range report.Adapters {
		statusColor := constants.ColorDim
		statusText := "DOWN"
		if a.IsUp {
			statusColor = constants.ColorGreen
			statusText = "UP"
		}

		ipStr := "-"
		if len(a.IPv4) > 0 {
			ipStr = strings.Join(a.IPv4, ", ")
		}

		speedStr := a.LinkSpeed
		if len(speedStr) == 0 {
			speedStr = "-"
		}

		tx := formatNetworkBytes(a.BytesSent)
		rx := formatNetworkBytes(a.BytesReceived)
		traffic := fmt.Sprintf("▲ %s / ▼ %s", tx, rx)

		fmt.Printf("  %-22s %s%-8s%s %-18s %-18s %-12s %s\n",
			truncateText(a.Name, 22),
			statusColor, statusText, constants.ColorReset,
			fallbackDash(a.MACAddress),
			truncateText(ipStr, 18),
			speedStr,
			traffic,
		)
	}

	fmt.Printf("  %s\n", strings.Repeat("─", 102))
	fmt.Printf("  Total Transmitted: %s  |  Total Received: %s  |  Timestamp: %s\n\n",
		formatNetworkBytes(report.TotalBytesSent),
		formatNetworkBytes(report.TotalBytesRecv),
		report.CollectedAt,
	)
}

func renderNetworkHelp() {
	fmt.Printf("\n%s● Network Adapter Telemetry (gitmap os network / gitmap network)%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println("  Usage:")
	fmt.Println("    gitmap network [--json] [--all]")
	fmt.Println("    gitmap os network [--json] [--all]")
	fmt.Println()
	fmt.Println("  Flags:")
	fmt.Println("    -j, --json     Output network telemetry report as structured JSON")
	fmt.Println("    -a, --all      Display all interfaces including inactive and loopback")
	fmt.Println("    -h, --help     Display this help menu")
	fmt.Println()
}

func formatNetworkBytes(b uint64) string {
	const unit = 1024
	if b < unit {
		return fmt.Sprintf("%d B", b)
	}
	div, exp := int64(unit), 0
	for n := b / unit; n >= unit; n /= unit {
		div *= unit
		exp++
	}
	units := []string{"KB", "MB", "GB", "TB", "PB"}
	if exp < len(units) {
		return fmt.Sprintf("%.1f %s", float64(b)/float64(div), units[exp])
	}
	return fmt.Sprintf("%d B", b)
}

func truncateText(s string, maxLen int) string {
	if len(s) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return s[:maxLen]
	}
	return s[:maxLen-3] + "..."
}
