//go:build linux

package cmdos

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	lightdmConfigDir = "/etc/lightdm/lightdm.conf.d"
	lightdmFile      = "/etc/lightdm/lightdm.conf.d/50-autologin.conf"
)

func configureLightDM(user string) error {
	_ = os.MkdirAll(filepath.Dir(lightdmFile), 0755)
	content := fmt.Sprintf("[Seat:*]\nautologin-user=%s\nautologin-user-timeout=0\n", user)
	if err := os.WriteFile(lightdmFile, []byte(content), 0644); err != nil {
		return apperror.WrapSimple(err, "failed writing LightDM config (run with sudo)")
	}
	return nil
}

func readLightDMStatus() (AutoLoginStatus, error) {
	content, err := os.ReadFile(lightdmFile)
	if err != nil {
		return AutoLoginStatus{DisplayManager: "LightDM"}, nil
	}
	user := parseLightDMUser(strings.Split(string(content), "\n"))
	return AutoLoginStatus{
		IsEnabled:      len(user) > 0,
		Username:       user,
		DisplayManager: "LightDM",
	}, nil
}

func parseLightDMUser(lines []string) string {
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "autologin-user=") {
			return strings.TrimPrefix(trimmed, "autologin-user=")
		}
	}
	return ""
}
