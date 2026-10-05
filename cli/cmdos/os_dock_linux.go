//go:build !windows && !darwin

package cmdos

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type linuxDockEngine struct{}

func newPlatformDockEngine() DockOperator {
	return &linuxDockEngine{}
}

func (l *linuxDockEngine) GetDockConfig() (DockConfig, error) {
	cfg := DockConfig{OS: "linux", IsPanelMode: true}
	out, err := execLinuxGSettings("get", "org.gnome.shell.extensions.dash-to-dock", "dock-position")
	if err != nil {
		cfg.Position = string(DockPositionBottom)
		cfg.RawPosition = "BOTTOM"
		return cfg, nil
	}
	raw := strings.Trim(strings.TrimSpace(string(out)), "'\"")
	cfg.RawPosition = raw
	cfg.Position = strings.ToLower(raw)
	return cfg, nil
}

func (l *linuxDockEngine) SetDockPosition(pos DockPosition) error {
	upper := strings.ToUpper(string(pos))
	val := fmt.Sprintf("'%s'", upper)
	_, err := execLinuxGSettings("set", "org.gnome.shell.extensions.dash-to-dock", "dock-position", val)
	if err != nil {
		return apperror.WrapSimple(err, "gsettings set dock-position failed")
	}
	return nil
}

func execLinuxGSettings(args ...string) ([]byte, error) {
	cmd := exec.Command("gsettings", args...)
	configureLinuxDBusEnv(cmd)
	return cmd.Output()
}

func configureLinuxDBusEnv(cmd *exec.Cmd) {
	if os.Getenv("DBUS_SESSION_BUS_ADDRESS") != "" {
		return
	}
	busPath := fmt.Sprintf("/run/user/%d/bus", os.Getuid())
	if _, err := os.Stat(busPath); err == nil {
		cmd.Env = append(os.Environ(), "DBUS_SESSION_BUS_ADDRESS=unix:path="+busPath)
	}
}
