// Package cmdos — os_network_windows.go queries Windows network adapters via PowerShell CIM cmdlets.
//go:build windows

package cmdos

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const windowsNetAdapterPowerShellScript = `Get-NetAdapter | ForEach-Object { $a = $_; $stats = $a | Get-NetAdapterStatistics -ErrorAction SilentlyContinue; $ips = Get-NetIPAddress -InterfaceIndex $a.ifIndex -ErrorAction SilentlyContinue; [PSCustomObject]@{ name = $a.Name; description = $a.InterfaceDescription; status = $a.Status; isUp = ($a.Status -eq 'Up'); macAddress = $a.MacAddress; linkSpeed = $a.LinkSpeed; bytesSent = if ($stats) { [uint64]$stats.SentBytes } else { [uint64]0 }; bytesReceived = if ($stats) { [uint64]$stats.ReceivedBytes } else { [uint64]0 }; ipv4 = @($ips | Where-Object AddressFamily -eq 'IPv4' | Select-Object -ExpandProperty IPAddress); ipv6 = @($ips | Where-Object AddressFamily -eq 'IPv6' | Select-Object -ExpandProperty IPAddress) } } | ConvertTo-Json -Compress`

// collectNetworkReport gathers network adapter telemetry on Windows platforms.
func collectNetworkReport() (NetworkReport, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 8*time.Second)
	defer cancel()

	cmd := exec.CommandContext(ctx, "powershell", "-NoProfile", "-NonInteractive", "-Command", windowsNetAdapterPowerShellScript)
	out, err := cmd.Output()
	if err != nil {
		return collectFallbackNetworkReport()
	}

	trimmed := strings.TrimSpace(string(out))
	if len(trimmed) == 0 {
		return collectFallbackNetworkReport()
	}

	var adapters []NetworkAdapterInfo
	if strings.HasPrefix(trimmed, "[") {
		if unmarshalErr := json.Unmarshal([]byte(trimmed), &adapters); unmarshalErr != nil {
			return collectFallbackNetworkReport()
		}
	} else if strings.HasPrefix(trimmed, "{") {
		var single NetworkAdapterInfo
		if unmarshalErr := json.Unmarshal([]byte(trimmed), &single); unmarshalErr != nil {
			return collectFallbackNetworkReport()
		}
		adapters = []NetworkAdapterInfo{single}
	} else {
		return collectFallbackNetworkReport()
	}

	return assembleNetworkReport(adapters), nil
}

func assembleNetworkReport(adapters []NetworkAdapterInfo) NetworkReport {
	var totalSent, totalRecv uint64
	for i := range adapters {
		adapters[i].IPv4Addresses = adapters[i].IPv4
		adapters[i].IPv6Addresses = adapters[i].IPv6
		totalSent += adapters[i].BytesSent
		totalRecv += adapters[i].BytesReceived
	}

	hostname, _ := os.Hostname()
	return NetworkReport{
		Platform:       runtime.GOOS,
		Hostname:       hostname,
		Adapters:       adapters,
		TotalBytesSent: totalSent,
		TotalBytesRecv: totalRecv,
		CollectedAt:    time.Now().UTC().Format(time.RFC3339),
	}
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
			ip, _, cidrErr := net.ParseCIDR(addr.String())
			if cidrErr == nil && ip != nil {
				if ip.To4() != nil {
					ipv4List = append(ipv4List, ip.String())
				} else {
					ipv6List = append(ipv6List, ip.String())
				}
			}
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
