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
	// Prepare cross-platform implementation structure
	return "", &apperror.AppError{
		Op:    "IPResolver.FetchLocalIP",
		Code:  "E_INTERNAL_ERROR",
		Ctx:   nil,
		Cause: nil,
	}
}

func GetLocalIP(ctx context.Context, skipLoopback bool, ifaceName string) (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "", &apperror.AppError{
			Op:    "GetLocalIP",
			Code:  "E_INTERNAL_ERROR",
			Cause: err,
			Ctx:   map[string]any{"ifaceName": ifaceName},
		}
	}

	for _, iface := range ifaces {
		if ifaceName != "" && iface.Name != ifaceName {
			continue
		}

		if skipLoopback && (iface.Flags&net.FlagLoopback != 0) {
			continue
		}

		if iface.Flags&net.FlagUp == 0 {
			continue
		}

		addrs, err := iface.Addrs()
		if err != nil {
			continue
		}

		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}

			if ip == nil || (ip.IsLoopback() && skipLoopback) {
				continue
			}

			ip = ip.To4()
			if ip == nil {
				continue
			}

			return ip.String(), nil
		}
	}

	return "", &apperror.AppError{
		Op:    "GetLocalIP",
		Code:  "E_INTERNAL_ERROR",
		Cause: fmt.Errorf("no suitable IP found"),
		Ctx:   map[string]any{"ifaceName": ifaceName},
	}
}
