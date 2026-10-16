// Package cmdos — os_network_types.go provides network adapter metrics models, CLI options, and reporting dispatch.
package cmdos

import (
	"encoding/json"
	"fmt"
	"net"
	"os"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// NetworkAdapterInfo defines the cross-platform telemetry collected per network interface.
type NetworkAdapterInfo struct {
	Name                string   `json:"name"`
	Description         string   `json:"description,omitempty"`
	InterfaceAlias      string   `json:"interfaceAlias,omitempty"`
	Status              string   `json:"status"`
	IsUp                bool     `json:"isUp"`
	IPv4                []string `json:"ipv4"`
	IPv4Addresses       []string `json:"ipv4Addresses,omitempty"`
	IPv6                []string `json:"ipv6,omitempty"`
	IPv6Addresses       []string `json:"ipv6Addresses,omitempty"`
	MACAddress          string   `json:"macAddress"`
	LinkSpeed           string   `json:"linkSpeed,omitempty"`
	LinkSpeedDisplay    string   `json:"linkSpeedDisplay,omitempty"`
	LinkSpeedBitsPerSec uint64   `json:"linkSpeedBitsPerSec,omitempty"`
	BytesSent           uint64   `json:"bytesSent"`
	BytesReceived       uint64   `json:"bytesReceived"`
	PacketsSent         uint64   `json:"packetsSent,omitempty"`
	PacketsReceived     uint64   `json:"packetsReceived,omitempty"`
}

// NetworkReport encapsulates the aggregate telemetry report for the host.
type NetworkReport struct {
	Platform       string               `json:"platform"`
	Hostname       string               `json:"hostname"`
	MachineName    string               `json:"machineName,omitempty"`
	Adapters       []NetworkAdapterInfo `json:"adapters"`
	TotalBytesSent uint64               `json:"totalBytesSent"`
	TotalBytesRecv uint64               `json:"totalBytesRecv"`
	CollectedAt    string               `json:"collectedAt"`
	Timestamp      string               `json:"timestamp,omitempty"`
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
	hasHelp := len(args) > 0 && isNetworkHelpArg(args[0])
	if hasHelp {
		renderNetworkHelp()
		return nil
	}

	opts := parseNetworkOptions(args)
	report, err := collectNetworkReport()
	if err != nil {
		return apperror.WrapSimple(err, "collect network telemetry report")
	}

	if opts.IsJSON {
		filtered := filterNetworkAdapters(report.Adapters, opts.IsAll)
		report.Adapters = filtered
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

func assembleNetworkReport(adapters []NetworkAdapterInfo) NetworkReport {
	var totalSent, totalRecv uint64
	for i := range adapters {
		normalizeAdapterFields(&adapters[i])
		totalSent += adapters[i].BytesSent
		totalRecv += adapters[i].BytesReceived
	}

	hostname, _ := os.Hostname()
	timestamp := time.Now().UTC().Format(time.RFC3339)
	return NetworkReport{
		Platform:       runtime.GOOS,
		Hostname:       hostname,
		MachineName:    hostname,
		Adapters:       adapters,
		TotalBytesSent: totalSent,
		TotalBytesRecv: totalRecv,
		CollectedAt:    timestamp,
		Timestamp:      timestamp,
	}
}

func normalizeAdapterFields(a *NetworkAdapterInfo) {
	hasIPv4 := len(a.IPv4) > 0
	if hasIPv4 {
		a.IPv4Addresses = a.IPv4
	}
	hasIPv4Addrs := len(a.IPv4Addresses) > 0
	if hasIPv4Addrs {
		a.IPv4 = a.IPv4Addresses
	}
	hasIPv6 := len(a.IPv6) > 0
	if hasIPv6 {
		a.IPv6Addresses = a.IPv6
	}
	hasIPv6Addrs := len(a.IPv6Addresses) > 0
	if hasIPv6Addrs {
		a.IPv6 = a.IPv6Addresses
	}
	hasSpeed := len(a.LinkSpeed) > 0
	if hasSpeed {
		a.LinkSpeedDisplay = a.LinkSpeed
	}
	hasSpeedDisplay := len(a.LinkSpeedDisplay) > 0
	if hasSpeedDisplay {
		a.LinkSpeed = a.LinkSpeedDisplay
	}
	hasAlias := len(a.InterfaceAlias) > 0
	if !hasAlias {
		a.InterfaceAlias = a.Name
	}
}

func appendParsedIP(addr net.Addr, ipv4List *[]string, ipv6List *[]string) {
	ip, _, cidrErr := net.ParseCIDR(addr.String())
	hasErr := cidrErr != nil || ip == nil
	if hasErr {
		return
	}
	isV4 := ip.To4() != nil
	if isV4 {
		*ipv4List = append(*ipv4List, ip.String())
		return
	}
	*ipv6List = append(*ipv6List, ip.String())
}
