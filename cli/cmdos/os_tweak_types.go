package cmdos

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos/enums"
)

// TweakTargetType defines system tweak categories.
type TweakTargetType = enums.TweakTargetType

const (
	TweakContextMenuType = enums.TweakContextMenuType
	TweakStartMenuType   = enums.TweakStartMenuType
	TweakPowerSchemeType = enums.TweakPowerSchemeType
	TweakHibernateType   = enums.TweakHibernateType
	TweakTelemetryType   = enums.TweakTelemetryType
	TweakActivityType    = enums.TweakActivityType
	TweakSearchType      = enums.TweakSearchType
)

// TweakStatus models the observed state of system tweaks.
type TweakStatus struct {
	IsClassicContextMenu   bool
	IsClassicStartMenu     bool
	IsUltimatePower        bool
	IsHibernateEnabled     bool
	IsTelemetryDisabled    bool
	IsActivityFeedDisabled bool
	IsBingSearchDisabled   bool
}

// TweakOperator abstracts tweak application and inspection.
type TweakOperator interface {
	SetContextMenu(isClassic bool) error
	SetStartMenu(isClassic bool) error
	SetPowerScheme(isUltimate bool) error
	SetHibernate(isEnabled bool) error
	SetTelemetry(isDisabled bool) error
	SetActivityFeed(isDisabled bool) error
	SetBingSearch(isDisabled bool) error
	GetStatus() (TweakStatus, error)
}
