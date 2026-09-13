package netip

import (
	"fmt"
	"net"
	"strconv"
	"strings"
)

// NetmaskToCIDR converts an IPv4 subnet mask to CIDR prefix length (default 24).
func NetmaskToCIDR(netmask string) int {
	ip := net.ParseIP(strings.TrimSpace(netmask)).To4()
	if ip == nil {
		return 24
	}

	return calcMaskOnes(ip)
}

func calcMaskOnes(ip net.IP) int {
	mask := net.IPv4Mask(ip[0], ip[1], ip[2], ip[3])
	ones, _ := mask.Size()
	if ones == 0 {
		return 24
	}

	return ones
}

// FormatCIDR returns the IP with CIDR suffix, e.g. "192.168.1.50/24".
func FormatCIDR(ip, netmask string) string {
	cidr := NetmaskToCIDR(netmask)

	return fmt.Sprintf("%s/%d", strings.TrimSpace(ip), cidr)
}

// CIDRToNetmask converts a prefix integer into a dotted quad netmask string.
func CIDRToNetmask(prefix int) string {
	if prefix < 0 || prefix > 32 {
		return "255.255.255.0"
	}

	mask := net.CIDRMask(prefix, 32)

	return fmt.Sprintf("%d.%d.%d.%d", mask[0], mask[1], mask[2], mask[3])
}

// ParseCIDRPrefix extracts the prefix int from an IP/CIDR string or returns 24.
func ParseCIDRPrefix(ipWithCIDR string) int {
	parts := strings.Split(ipWithCIDR, "/")
	if len(parts) < 2 {
		return 24
	}

	prefix, err := strconv.Atoi(parts[1])
	if err != nil {
		return 24
	}

	return prefix
}
