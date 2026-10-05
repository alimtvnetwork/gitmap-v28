//go:build darwin

package cmdos

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type darwinDockEngine struct{}

func newPlatformDockEngine() DockOperator {
	return &darwinDockEngine{}
}

func (d *darwinDockEngine) GetDockConfig() (DockConfig, error) {
	cfg := DockConfig{OS: "darwin", IsPanelMode: false}
	out, err := exec.Command("defaults", "read", "com.apple.dock", "orientation").Output()
	if err != nil {
		cfg.Position = string(DockPositionBottom)
		cfg.RawPosition = "bottom"
		return cfg, nil
	}
	raw := strings.TrimSpace(string(out))
	cfg.RawPosition = raw
	cfg.Position = strings.ToLower(raw)
	return cfg, nil
}

func (d *darwinDockEngine) SetDockPosition(pos DockPosition) error {
	val := string(pos)
	if err := exec.Command("defaults", "write", "com.apple.dock", "orientation", val).Run(); err != nil {
		return apperror.WrapSimple(err, "defaults write orientation failed")
	}
	_ = exec.Command("killall", "Dock").Run()
	return nil
}
