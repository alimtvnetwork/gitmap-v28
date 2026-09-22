//go:build !windows

package cmdos

import (
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type unsupportedTweakEngine struct{}

func newPlatformTweakEngine() TweakOperator {
	return &unsupportedTweakEngine{}
}

func (u *unsupportedTweakEngine) SetContextMenu(isClassic bool) error {
	return apperror.NewSimple("context-menu tweak is Windows-only, not supported on "+runtime.GOOS, "E_TWEAK_UNSUPPORTED")
}

func (u *unsupportedTweakEngine) SetStartMenu(isClassic bool) error {
	return apperror.NewSimple("start-menu tweak is Windows-only, not supported on "+runtime.GOOS, "E_TWEAK_UNSUPPORTED")
}

func (u *unsupportedTweakEngine) SetPowerScheme(isUltimate bool) error {
	return apperror.NewSimple("power scheme tweak is Windows-only, not supported on "+runtime.GOOS, "E_TWEAK_UNSUPPORTED")
}

func (u *unsupportedTweakEngine) SetHibernate(isEnabled bool) error {
	return apperror.NewSimple("hibernate tweak is Windows-only, not supported on "+runtime.GOOS, "E_TWEAK_UNSUPPORTED")
}

func (u *unsupportedTweakEngine) SetTelemetry(isDisabled bool) error {
	return apperror.NewSimple("telemetry tweak is Windows-only, not supported on "+runtime.GOOS, "E_TWEAK_UNSUPPORTED")
}

func (u *unsupportedTweakEngine) SetActivityFeed(isDisabled bool) error {
	return apperror.NewSimple("activity feed tweak is Windows-only, not supported on "+runtime.GOOS, "E_TWEAK_UNSUPPORTED")
}

func (u *unsupportedTweakEngine) SetBingSearch(isDisabled bool) error {
	return apperror.NewSimple("search tweak is Windows-only, not supported on "+runtime.GOOS, "E_TWEAK_UNSUPPORTED")
}

func (u *unsupportedTweakEngine) GetStatus() (TweakStatus, error) {
	return TweakStatus{}, nil
}
