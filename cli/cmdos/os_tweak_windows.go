//go:build windows

package cmdos

import "os"

type windowsTweakEngine struct{}

func newPlatformTweakEngine() TweakEngine {
	return &windowsTweakEngine{}
}

func (w *windowsTweakEngine) SetContextMenu(isClassic bool) error {
	return applyContextMenuRegistry(isClassic)
}

func (w *windowsTweakEngine) SetStartMenu(isClassic bool) error {
	return applyStartMenuRegistry(isClassic)
}

func (w *windowsTweakEngine) SetPowerScheme(isUltimate bool) error {
	return applyPowerScheme(isUltimate)
}

func (w *windowsTweakEngine) SetHibernate(isEnabled bool) error {
	return applyHibernate(isEnabled)
}

func (w *windowsTweakEngine) SetTelemetry(isDisabled bool) error {
	return applyTelemetry(isDisabled)
}

func (w *windowsTweakEngine) SetActivityFeed(isDisabled bool) error {
	return applyActivityFeed(isDisabled)
}

func (w *windowsTweakEngine) SetBingSearch(isDisabled bool) error {
	return applyBingSearch(isDisabled)
}

func (w *windowsTweakEngine) GetStatus() (TweakStatus, error) {
	_, err := os.Stat(`C:\hiberfil.sys`)
	hasHiberfil := err == nil

	return TweakStatus{
		IsClassicContextMenu:   isClassicContextMenuActive(),
		IsClassicStartMenu:     isClassicStartMenuActive(),
		IsUltimatePower:        isUltimatePowerActive(),
		IsHibernateEnabled:     hasHiberfil,
		IsTelemetryDisabled:    isTelemetryDisabled(),
		IsActivityFeedDisabled: isActivityFeedDisabled(),
		IsBingSearchDisabled:   isBingSearchDisabled(),
	}, nil
}
