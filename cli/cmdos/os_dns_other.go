//go:build !windows && !linux

package cmdos

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type otherDNSEngine struct{}

func newPlatformDNSEngine() DNSEngine {
	return &otherDNSEngine{}
}

func (o *otherDNSEngine) GetDefaultInterface() (string, error) {
	return "en0", nil
}

func (o *otherDNSEngine) SetDNS(_ string, _ DNSProvider) error {
	msg := "DNS configuration is not supported on " + runtime.GOOS

	return apperror.NewSimple(msg, "E_OS_UNSUPPORTED")
}

func (o *otherDNSEngine) SetDHCP(_ string) error {
	msg := "DNS DHCP revert is not supported on " + runtime.GOOS

	return apperror.NewSimple(msg, "E_OS_UNSUPPORTED")
}

func (o *otherDNSEngine) GetDNS(_ string) ([]string, error) {
	return []string{}, nil
}
