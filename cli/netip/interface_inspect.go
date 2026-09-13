package netip

import (
	"fmt"
	"net"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ListSystemInterfaces returns all network interfaces found on the machine.
func ListSystemInterfaces() InterfaceSliceResult {
	ifaces, err := net.Interfaces()
	if err != nil {
		return result.FailSlice[InterfaceInfo](apperror.WrapSimple(err, "ListSystemInterfaces"))
	}

	infos := make([]InterfaceInfo, 0, len(ifaces))
	for _, iface := range ifaces {
		infos = append(infos, buildInterfaceInfo(iface))
	}

	return result.OkSlice(infos)
}

// GetSystemInterface returns details for a named interface.
func GetSystemInterface(name string) InterfaceResult {
	iface, err := net.InterfaceByName(name)
	if err != nil {
		return result.Fail[InterfaceInfo](apperror.WrapSimple(err, "GetSystemInterface"))
	}

	return result.Ok(buildInterfaceInfo(*iface))
}

func buildInterfaceInfo(iface net.Interface) InterfaceInfo {
	ip, mask := extractInterfaceIPv4(iface)
	isUp := iface.Flags&net.FlagUp != 0
	isLoopback := iface.Flags&net.FlagLoopback != 0

	return InterfaceInfo{
		Name:       iface.Name,
		IP:         ip,
		Netmask:    mask,
		MAC:        iface.HardwareAddr.String(),
		IsUp:       isUp,
		IsLoopback: isLoopback,
		IsValid:    ip != "",
	}
}

func extractInterfaceIPv4(iface net.Interface) (string, string) {
	addrs, err := iface.Addrs()
	if err != nil {
		return "", ""
	}

	return findFirstIPv4(addrs)
}

func findFirstIPv4(addrs []net.Addr) (string, string) {
	for _, addr := range addrs {
		ipStr, maskStr := inspectAddr(addr)
		if ipStr != "" {
			return ipStr, maskStr
		}
	}

	return "", ""
}

func inspectAddr(addr net.Addr) (string, string) {
	ipNet, isIPNet := addr.(*net.IPNet)
	if !isIPNet {
		return "", ""
	}

	ip4 := ipNet.IP.To4()
	if ip4 == nil {
		return "", ""
	}

	return ip4.String(), formatMask(ipNet.Mask)
}

func formatMask(mask net.IPMask) string {
	if len(mask) == 4 {
		return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
	}

	return "255.255.255.0"
}
