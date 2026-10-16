// Package cmdos — os_network_linux.go gathers network metrics from /proc/net/dev, /sys/class/net, and net.Interfaces().
//go:build !windows

package cmdos

import (
	"bufio"
	"net"
	"os"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type procNetDevStats struct {
	bytesRecv   uint64
	packetsRecv uint64
	bytesSent   uint64
	packetsSent uint64
}

// collectNetworkReport gathers network adapter telemetry on Linux, macOS, and POSIX platforms.
func collectNetworkReport() (NetworkReport, error) {
	devStats := parseProcNetDev()
	ifaces, err := net.Interfaces()
	if err != nil {
		return NetworkReport{}, apperror.WrapSimple(err, "enumerate posix net interfaces")
	}

	var adapters []NetworkAdapterInfo
	var totalSent, totalRecv uint64

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

		linkSpeed := readLinuxInterfaceSpeed(iface.Name)
		stat := devStats[iface.Name]

		totalSent += stat.bytesSent
		totalRecv += stat.bytesRecv

		adapters = append(adapters, NetworkAdapterInfo{
			Name:            iface.Name,
			Status:          status,
			IsUp:            isUp,
			IPv4:            ipv4List,
			IPv4Addresses:   ipv4List,
			IPv6:            ipv6List,
			IPv6Addresses:   ipv6List,
			MACAddress:      iface.HardwareAddr.String(),
			LinkSpeed:       linkSpeed,
			BytesSent:       stat.bytesSent,
			BytesReceived:   stat.bytesRecv,
			PacketsSent:     stat.packetsSent,
			PacketsReceived: stat.packetsRecv,
		})
	}

	hostname, _ := os.Hostname()
	return NetworkReport{
		Platform:       runtime.GOOS,
		Hostname:       hostname,
		Adapters:       adapters,
		TotalBytesSent: totalSent,
		TotalBytesRecv: totalRecv,
		CollectedAt:    time.Now().UTC().Format(time.RFC3339),
	}, nil
}

func parseProcNetDev() map[string]procNetDevStats {
	stats := make(map[string]procNetDevStats)
	file, err := os.Open("/proc/net/dev")
	if err != nil {
		return stats
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	lineCount := 0
	for scanner.Scan() {
		lineCount++
		if lineCount <= 2 {
			continue
		}
		line := strings.TrimSpace(scanner.Text())
		colonIdx := strings.Index(line, ":")
		if colonIdx == -1 {
			continue
		}
		ifaceName := strings.TrimSpace(line[:colonIdx])
		fields := strings.Fields(line[colonIdx+1:])
		if len(fields) < 16 {
			continue
		}

		rxBytes, _ := strconv.ParseUint(fields[0], 10, 64)
		rxPackets, _ := strconv.ParseUint(fields[1], 10, 64)
		txBytes, _ := strconv.ParseUint(fields[8], 10, 64)
		txPackets, _ := strconv.ParseUint(fields[9], 10, 64)

		stats[ifaceName] = procNetDevStats{
			bytesRecv:   rxBytes,
			packetsRecv: rxPackets,
			bytesSent:   txBytes,
			packetsSent: txPackets,
		}
	}
	return stats
}

func readLinuxInterfaceSpeed(ifaceName string) string {
	speedPath := filepath.Join("/sys/class/net", ifaceName, "speed")
	content, err := os.ReadFile(speedPath)
	if err != nil {
		return ""
	}
	raw := strings.TrimSpace(string(content))
	if len(raw) == 0 || raw == "-1" {
		return ""
	}
	return raw + " Mbps"
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
