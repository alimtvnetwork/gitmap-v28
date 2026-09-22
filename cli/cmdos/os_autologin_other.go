//go:build !windows && !linux

package cmdos

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type unsupportedAutoLoginEngine struct{}

func newPlatformAutoLoginEngine() AutoLoginEngine {
	return &unsupportedAutoLoginEngine{}
}

func (u *unsupportedAutoLoginEngine) Configure(cfg AutoLoginConfig) error {
	return apperror.NewSimple("os autologin is not supported on "+runtime.GOOS, "E_OS_AUTOLOGIN_UNSUPPORTED")
}

func (u *unsupportedAutoLoginEngine) Disable() error {
	return apperror.NewSimple("os autologin is not supported on "+runtime.GOOS, "E_OS_AUTOLOGIN_UNSUPPORTED")
}

func (u *unsupportedAutoLoginEngine) Status() (AutoLoginStatus, error) {
	return AutoLoginStatus{
		IsEnabled:      false,
		DisplayManager: runtime.GOOS,
	}, nil
}
