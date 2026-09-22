//go:build !windows && !linux

package cmdos

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type otherThemeEngine struct{}

func newPlatformThemeEngine() ThemeEngine {
	return &otherThemeEngine{}
}

func (o *otherThemeEngine) SetTheme(_ ThemeMode) error {
	msg := "desktop theme switching is not supported on " + runtime.GOOS

	return apperror.NewSimple(msg, "E_OS_UNSUPPORTED")
}

func (o *otherThemeEngine) GetTheme() (ThemeMode, error) {
	return ThemeModeDark, nil
}
