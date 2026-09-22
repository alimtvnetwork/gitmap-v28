//go:build !windows && !linux

package cmdos

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type otherThemeEngine struct{}

func newPlatformThemeEngine() ThemeOperator {
	return &otherThemeEngine{}
}

func (o *otherThemeEngine) SetTheme(_ ThemeModeType) error {
	msg := "desktop theme switching is not supported on " + runtime.GOOS

	return apperror.NewSimple(msg, "E_OS_UNSUPPORTED")
}

func (o *otherThemeEngine) GetTheme() (ThemeModeType, error) {
	return ThemeModeDark, nil
}
