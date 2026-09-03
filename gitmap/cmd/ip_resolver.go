package cmd

import (
	"context"
	"fmt"
	"net"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type IPResolver struct {
	Cache   map[string]string
	Timeout time.Duration
}

func (r *IPResolver) FetchLocalIP(ctx context.Context) (string, error) {
	return "", &apperror.AppError{
		Op:    "IPResolver.FetchLocalIP",
		Code:  "E_INTERNAL_ERROR",
		Ctx:   nil,
		Cause: nil,
	}
}

func GetLocalIP(ctx context.Context, isSkipLoopback bool, ifaceName string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", newLocalIPError("GetLocalIP", err, ifaceName)
	}

	ipStr := findMatchingInterfaceIP(ifaces, ifaceName, isSkipLoopback)
	if ipStr != "" {
		return ipStr, nil
	}

	return "", newLocalIPError("GetLocalIP", fmt.Errorf("no suitable IP found"), ifaceName)
}

func newLocalIPError(op string, cause error, ifaceName string) *apperror.AppError {
	return &apperror.AppError{
		Op:    op,
		Code:  "E_INTERNAL_ERROR",
		Cause: cause,
		Ctx:   map[string]any{"ifaceName": ifaceName},
	}
}

func findMatchingInterfaceIP(ifaces []net.Interface, ifaceName string, isSkipLoopback bool) string {
	for _, iface := range ifaces {
		if !isInterfaceEligible(iface, ifaceName, isSkipLoopback) {
			continue
		}

		ipStr := extractIPFromInterface(iface, isSkipLoopback)
		if ipStr != "" {
			return ipStr
		}
	}

	return ""
}

func isInterfaceEligible(iface net.Interface, ifaceName string, isSkipLoopback bool) bool {
	if ifaceName != "" && iface.Name != ifaceName {
		return false
	}

	if isSkipLoopback && (iface.Flags&net.FlagLoopback != 0) {
		return false
	}

	return iface.Flags&net.FlagUp != 0
}

func extractIPFromInterface(iface net.Interface, isSkipLoopback bool) string {
	addrs, err := iface.Addrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		ip := resolveIPv4FromAddr(addr, isSkipLoopback)
		if ip != nil {
			return ip.String()
		}
	}

	return ""
}

func resolveIPv4FromAddr(addr net.Addr, isSkipLoopback bool) net.IP {
	ip := extractRawIP(addr)
	if ip == nil {
		return nil
	}

	if isSkipLoopback && ip.IsLoopback() {
		return nil
	}

	return ip.To4()
}

func extractRawIP(addr net.Addr) net.IP {
	switch v := addr.(type) {
	case *net.IPNet:
		return v.IP
	case *net.IPAddr:
		return v.IP
	default:
		return nil
	}
}
