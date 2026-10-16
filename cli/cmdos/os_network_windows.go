// Package cmdos — os_network_windows.go queries Windows network adapters via PowerShell CIM cmdlets.
//go:build windows

package cmdos

import (
	"context"
	"encoding/json"
	"net"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const windowsNetAdapterPowerShellScript = `Get-NetAdapter | ForEach-Object { $a = $_; $stats = $a | Get-NetAdapterStatistics -ErrorAction SilentlyContinue; $ips = Get-NetIPAddress -InterfaceIndex $a.ifIndex -ErrorAction SilentlyContinue; $rxPkt = if ($stats) { [uint64]($stats.ReceivedUnicastPackets + $stats.ReceivedMulticastPackets + $stats.ReceivedBroadcastPackets) } else { [uint64]0 }; $txPkt = if ($stats) { [uint64]($stats.SentUnicastPackets + $stats.SentMulticastPackets + $stats.SentBroadcastPackets) } else { [uint64]0 }; [PSCustomObject]@{ name = $a.Name; description = $a.InterfaceDescription; status = $a.Status; isUp = ($a.Status -eq 'Up'); macAddress = $a.MacAddress; linkSpeed = $a.LinkSpeed; bytesSent = if ($stats) { [uint64]$stats.SentBytes } else { [uint64]0 }; bytesReceived = if ($stats) { [uint64]$stats.ReceivedBytes } else { [uint64]0 }; packetsSent = $txPkt; packetsReceived = $rxPkt; ipv4 = @($ips | Where-Object AddressFamily -eq 'IPv4' | Select-Object -ExpandProperty IPAddress); ipv6 = @($ips | Where-Object AddressFamily -eq 'IPv6' | Select-Object -ExpandProperty IPAddress) } } | ConvertTo-Json -Compress`

// collectNetworkReport gathers network adapter telemetry on Windows platforms.
func collectNetworkReport() (NetworkReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", windowsNetAdapterPowerShellScript)
	out, err := cmd.Output()
	if err != nil {
		return collectFallbackNetworkReport()
	}

	adapters, parseErr := parseWindowsAdaptersJSON(string(out))
	if parseErr != nil {
		return collectFallbackNetworkReport()
	}

	return assembleNetworkReport(adapters), nil
}

func parseWindowsAdaptersJSON(raw string) ([]NetworkAdapterInfo, error) {
	trimmed := strings.TrimSpace(raw)
	isEmpty := len(trimmed) == 0
	if isEmpty {
		return nil, apperror.NewSimple("empty powershell json output")
	}

	isArray := strings.HasPrefix(trimmed, "[")
	if isArray {
		return unmarshalAdapterList([]byte(trimmed))
	}

	isObject := strings.HasPrefix(trimmed, "{")
	if isObject {
		return unmarshalAdapterSingle([]byte(trimmed))
	}

	return nil, apperror.NewSimple("unexpected json format from powershell")
}

func unmarshalAdapterList(data []byte) ([]NetworkAdapterInfo, error) {
	var list []NetworkAdapterInfo
	err := json.Unmarshal(data, &list)
	if err != nil {
		return nil, apperror.WrapSimple(err, "unmarshal adapter array")
	}
	return list, nil
}

func unmarshalAdapterSingle(data []byte) ([]NetworkAdapterInfo, error) {
	var single NetworkAdapterInfo
	err := json.Unmarshal(data, &single)
	if err != nil {
		return nil, apperror.WrapSimple(err, "unmarshal adapter object")
	}
	return []NetworkAdapterInfo{single}, nil
}

// collectFallbackNetworkReport provides resilience by querying standard Go net.Interfaces if PowerShell fails.
func collectFallbackNetworkReport() (NetworkReport, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return NetworkReport{}, apperror.WrapSimple(err, "enumerate fallback net interfaces")
	}

	var adapters []NetworkAdapterInfo
	for _, iface := range ifaces {
		isUp := (iface.Flags & net.FlagUp) != 0
		status := "Down"
		if isUp {
			status = "Up"
		}

		var ipv4List, ipv6List []string
		addrs, _ := iface.Addrs()
		for _, addr := range addrs {
			appendParsedIP(addr, &ipv4List, &ipv6List)
		}

		adapters = append(adapters, NetworkAdapterInfo{
			Name:          iface.Name,
			Status:        status,
			IsUp:          isUp,
			IPv4:          ipv4List,
			IPv4Addresses: ipv4List,
			IPv6:          ipv6List,
			IPv6Addresses: ipv6List,
			MACAddress:    iface.HardwareAddr.String(),
		})
	}

	return assembleNetworkReport(adapters), nil
}
