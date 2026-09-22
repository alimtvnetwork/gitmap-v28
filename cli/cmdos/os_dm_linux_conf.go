//go:build linux

package cmdos

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func (m *linuxDMManager) checkWaylandState() bool {
	data, err := os.ReadFile(m.customConfPath)
	if err != nil {
		return true
	}

	lines := strings.Split(string(data), "\n")
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "WaylandEnable=false") {
			return false
		}
	}

	return true
}

func (m *linuxDMManager) checkAutoLogin() (string, bool) {
	data, err := os.ReadFile(m.customConfPath)
	if err != nil {
		return "", false
	}

	lines := strings.Split(string(data), "\n")
	user := ""
	isAuto := false
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "AutomaticLoginEnable=true") {
			isAuto = true
		}

		if strings.HasPrefix(trimmed, "AutomaticLogin=") && !strings.HasPrefix(trimmed, "AutomaticLoginEnable=") {
			user = strings.TrimPrefix(trimmed, "AutomaticLogin=")
		}
	}

	return user, isAuto
}

func (m *linuxDMManager) SetWayland(isEnabled bool) error {
	data, err := os.ReadFile(m.customConfPath)
	if err != nil {
		return apperror.WrapSimple(err, "readGDM3CustomConf")
	}

	updated := updateWaylandInConfig(string(data), isEnabled)

	return os.WriteFile(m.customConfPath, []byte(updated), 0644)
}
