//go:build linux

package cmdos

import (
	"net"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type linuxDNSEngine struct{}

func newPlatformDNSEngine() DNSOperator {
	return &linuxDNSEngine{}
}

func (l *linuxDNSEngine) GetDefaultInterface() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "eth0", nil
	}

	for _, iface := range ifaces {
		isDown := iface.Flags&net.FlagUp == 0
		isLoopback := iface.Flags&net.FlagLoopback != 0
		isVirtual := strings.HasPrefix(iface.Name, "docker") || strings.HasPrefix(iface.Name, "veth")
		if isDown || isLoopback || isVirtual {
			continue
		}

		return iface.Name, nil
	}

	return "eth0", nil
}

func (l *linuxDNSEngine) SetDNS(iface string, p DNSProvider) error {
	args := []string{"dns", iface, p.Primary}
	if p.Secondary != "" {
		args = append(args, p.Secondary)
	}

	cmd := exec.Command("resolvectl", args...)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "resolvectlSetDNS")
	}

	return nil
}

func (l *linuxDNSEngine) SetDHCP(iface string) error {
	cmd := exec.Command("resolvectl", "revert", iface)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "resolvectlRevertDNS")
	}

	return nil
}

func (l *linuxDNSEngine) GetDNS(iface string) ([]string, error) {
	cmd := exec.Command("resolvectl", "dns", iface)
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "resolvectlQueryDNS")
	}

	fields := strings.Fields(string(out))
	var servers []string
	for _, f := range fields {
		if ip := net.ParseIP(f); ip != nil {
			servers = append(servers, f)
		}
	}

	return servers, nil
}
