//go:build !windows && !linux && !darwin

package power

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

type fallbackDriver struct{}

func init() {
	RegisterDriver(runtime.GOOS, func(r CmdRunner) Manager {
		return &fallbackDriver{}
	})
}

func (f *fallbackDriver) Platform() string {
	return runtime.GOOS
}

func (f *fallbackDriver) GetStatus() (Settings, error) {
	return Settings{
		Platform: runtime.GOOS,
		Source:   "unsupported",
	}, nil
}

func (f *fallbackDriver) SetNeverSleep() error {
	return apperror.NewSimple("power management unsupported on "+runtime.GOOS, "E_UNSUPPORTED")
}

func (f *fallbackDriver) SetTimeouts(displayMinutes, sleepMinutes int) error {
	return apperror.NewSimple("power management unsupported on "+runtime.GOOS, "E_UNSUPPORTED")
}

func (f *fallbackDriver) ApplySettings(s Settings) error {
	return apperror.NewSimple("power management unsupported on "+runtime.GOOS, "E_UNSUPPORTED")
}
