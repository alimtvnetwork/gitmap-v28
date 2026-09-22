//go:build !linux

package cmdos

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type otherDMManager struct{}

func newLinuxDMManager() DMManager {
	return &otherDMManager{}
}

func (m *otherDMManager) GetStatus() (DMStatus, error) {
	return DMStatus{
		Name:             runtime.GOOS,
		ServiceStatus:    "unsupported",
		SessionType:      runtime.GOOS,
		IsWaylandEnabled: false,
		IsAutoLoginSet:   false,
		AutoLoginUser:    "",
		ConfigFile:       "N/A",
	}, nil
}

func (m *otherDMManager) SetWayland(_ bool) error {
	msg := "Display Manager (DM) configuration is only supported on Linux"

	return apperror.NewSimple(msg, "E_OS_UNSUPPORTED")
}

func (m *otherDMManager) RestartService() error {
	msg := "Display Manager restart is only supported on Linux"

	return apperror.NewSimple(msg, "E_OS_UNSUPPORTED")
}
