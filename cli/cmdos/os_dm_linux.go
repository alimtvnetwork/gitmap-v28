//go:build linux

package cmdos

import (
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type linuxDMManager struct {
	customConfPath string
	defaultDMPath  string
}

func newLinuxDMManager() *linuxDMManager {
	return &linuxDMManager{
		customConfPath: "/etc/gdm3/custom.conf",
		defaultDMPath:  "/etc/X11/default-display-manager",
	}
}

func (m *linuxDMManager) GetStatus() (DMStatus, error) {
	dmName := m.detectActiveDM()
	sessType := os.Getenv("XDG_SESSION_TYPE")
	if sessType == "" {
		sessType = "x11"
	}

	isWayland := m.checkWaylandState()
	autoUser, isAuto := m.checkAutoLogin()

	return DMStatus{
		Name:             dmName,
		ServiceStatus:    m.getServiceState(),
		SessionType:      sessType,
		IsWaylandEnabled: isWayland,
		IsAutoLoginSet:   isAuto,
		AutoLoginUser:    autoUser,
		ConfigFile:       m.customConfPath,
	}, nil
}

func (m *linuxDMManager) detectActiveDM() string {
	data, err := os.ReadFile(m.defaultDMPath)
	if err != nil {
		return string(DMTypeGDM3)
	}

	lines := strings.TrimSpace(string(data))
	if strings.Contains(lines, "lightdm") {
		return string(DMTypeLightDM)
	}
	if strings.Contains(lines, "sddm") {
		return string(DMTypeSDDM)
	}

	return string(DMTypeGDM3)
}

func (m *linuxDMManager) getServiceState() string {
	cmd := exec.Command("systemctl", "is-active", "display-manager")
	out, err := cmd.Output()
	if err != nil {
		return "inactive"
	}

	return strings.TrimSpace(string(out))
}

func (m *linuxDMManager) RestartService() error {
	cmd := exec.Command("systemctl", "restart", "display-manager")
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "restartDisplayManager")
	}

	return nil
}
