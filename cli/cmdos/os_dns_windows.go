//go:build windows

package cmdos

import (
	"net"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type windowsDNSEngine struct{}

func newPlatformDNSEngine() DNSOperator {
	return &windowsDNSEngine{}
}

func (w *windowsDNSEngine) GetDefaultInterface() (string, error) {
	ifaces, err := net.Interfaces()
	if err != nil {
		return "Ethernet", nil
	}

	for _, iface := range ifaces {
		isDown := iface.Flags&net.FlagUp == 0
		isLoopback := iface.Flags&net.FlagLoopback != 0
		isVirtual := strings.Contains(iface.Name, "vEthernet") || strings.Contains(iface.Name, "Loopback")
		if isDown || isLoopback || isVirtual {
			continue
		}

		return iface.Name, nil
	}

	return "Wi-Fi", nil
}

func (w *windowsDNSEngine) SetDNS(iface string, p DNSProvider) error {
	cmd1 := exec.Command("netsh", "interface", "ip", "set", "dns", "name="+iface, "static", p.Primary)
	if err := cmd1.Run(); err != nil {
		return apperror.WrapSimple(err, "failed to set primary DNS via netsh")
	}

	if p.Secondary != "" {
		cmd2 := exec.Command("netsh", "interface", "ip", "add", "dns", "name="+iface, p.Secondary, "index=2")
		_ = cmd2.Run()
	}

	return nil
}

func (w *windowsDNSEngine) SetDHCP(iface string) error {
	cmd := exec.Command("netsh", "interface", "ip", "set", "dns", "name="+iface, "dhcp")
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "failed to revert DNS to DHCP via netsh")
	}

	return nil
}

func (w *windowsDNSEngine) GetDNS(iface string) ([]string, error) {
	cmd := exec.Command("netsh", "interface", "ip", "show", "dns", "name="+iface)
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to query DNS via netsh")
	}

	lines := strings.Split(string(out), "\n")
	var servers []string
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if ip := net.ParseIP(trimmed); ip != nil {
			servers = append(servers, trimmed)
		}
	}

	return servers, nil
}
