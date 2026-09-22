//go:build windows

package cmdos

import (
	"golang.org/x/sys/windows/registry"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const themePersonalizeKey = `Software\Microsoft\Windows\CurrentVersion\Themes\Personalize`

type windowsThemeEngine struct{}

func newPlatformThemeEngine() ThemeEngine {
	return &windowsThemeEngine{}
}

func (w *windowsThemeEngine) SetTheme(mode ThemeMode) error {
	k, _, err := registry.CreateKey(registry.CURRENT_USER, themePersonalizeKey, registry.SET_VALUE)
	if err != nil {
		return apperror.WrapSimple(err, "failed to open Themes Personalize key")
	}
	defer k.Close()

	val := uint32(1)
	if mode == ThemeModeDark {
		val = 0
	}

	_ = k.SetDWordValue("AppsUseLightTheme", val)
	_ = k.SetDWordValue("SystemUsesLightTheme", val)

	return nil
}

func (w *windowsThemeEngine) GetTheme() (ThemeMode, error) {
	k, err := registry.OpenKey(registry.CURRENT_USER, themePersonalizeKey, registry.QUERY_VALUE)
	if err != nil {
		return ThemeModeDark, nil
	}
	defer k.Close()

	val, _, err := k.GetIntegerValue("AppsUseLightTheme")
	if err != nil || val == 0 {
		return ThemeModeDark, nil
	}

	return ThemeModeLight, nil
}
