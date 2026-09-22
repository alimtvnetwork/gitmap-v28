//go:build linux

package cmdos

import (
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type linuxThemeEngine struct{}

func newPlatformThemeEngine() ThemeOperator {
	return &linuxThemeEngine{}
}

func (l *linuxThemeEngine) SetTheme(mode ThemeModeType) error {
	val := "'prefer-dark'"
	if mode == ThemeModeLight {
		val = "'default'"
	}

	cmd := exec.Command("gsettings", "set", "org.gnome.desktop.interface", "color-scheme", val)
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "gsettingsSetColorScheme")
	}

	return nil
}

func (l *linuxThemeEngine) GetTheme() (ThemeModeType, error) {
	cmd := exec.Command("gsettings", "get", "org.gnome.desktop.interface", "color-scheme")
	out, err := cmd.Output()
	if err != nil {
		return ThemeModeDark, nil
	}

	val := strings.TrimSpace(string(out))
	if strings.Contains(val, "dark") {
		return ThemeModeDark, nil
	}

	return ThemeModeLight, nil
}
