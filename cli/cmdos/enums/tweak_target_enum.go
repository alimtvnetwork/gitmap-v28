package enums

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// TweakTargetType defines system tweak categories.
type TweakTargetType string

const (
	TweakContextMenuType TweakTargetType = "context-menu"
	TweakStartMenuType   TweakTargetType = "start-menu"
	TweakPowerSchemeType TweakTargetType = "power"
	TweakHibernateType   TweakTargetType = "hibernate"
	TweakTelemetryType   TweakTargetType = "telemetry"
	TweakActivityType    TweakTargetType = "activity"
	TweakSearchType      TweakTargetType = "search"
)

// IsValid reports whether the tweak target matches a known variant.
func (t TweakTargetType) IsValid() bool {
	switch t {
	case TweakContextMenuType, TweakStartMenuType, TweakPowerSchemeType,
		TweakHibernateType, TweakTelemetryType, TweakActivityType,
		TweakSearchType:
		return true
	default:
		return false
	}
}

// String returns the string representation of the tweak target.
func (t TweakTargetType) String() string {
	return string(t)
}

// ParseTweakTargetType parses a raw string into a validated TweakTargetType.
func ParseTweakTargetType(s string) (TweakTargetType, *apperror.AppError) {
	val := TweakTargetType(strings.ToLower(strings.TrimSpace(s)))
	if val.IsValid() {
		return val, nil
	}

	return "", apperror.NewValidation("ParseTweakTargetType", "E1000", fmt.Sprintf("invalid tweak target: %q", s)).WithContext("value", s)
}

// AllTweakTargetTypes returns all valid TweakTargetType variants.
func AllTweakTargetTypes() []TweakTargetType {
	return []TweakTargetType{
		TweakContextMenuType,
		TweakStartMenuType,
		TweakPowerSchemeType,
		TweakHibernateType,
		TweakTelemetryType,
		TweakActivityType,
		TweakSearchType,
	}
}

// IsContextMenu reports whether the tweak target is context-menu.
func (t TweakTargetType) IsContextMenu() bool {
	return t == TweakContextMenuType
}

// IsStartMenu reports whether the tweak target is start-menu.
func (t TweakTargetType) IsStartMenu() bool {
	return t == TweakStartMenuType
}

// IsPowerScheme reports whether the tweak target is power.
func (t TweakTargetType) IsPowerScheme() bool {
	return t == TweakPowerSchemeType
}

// IsHibernate reports whether the tweak target is hibernate.
func (t TweakTargetType) IsHibernate() bool {
	return t == TweakHibernateType
}

// IsTelemetry reports whether the tweak target is telemetry.
func (t TweakTargetType) IsTelemetry() bool {
	return t == TweakTelemetryType
}

// IsActivity reports whether the tweak target is activity.
func (t TweakTargetType) IsActivity() bool {
	return t == TweakActivityType
}

// IsSearch reports whether the tweak target is search.
func (t TweakTargetType) IsSearch() bool {
	return t == TweakSearchType
}

// HasTweakTarget reports whether the tweak target is non-empty and valid.
func (t TweakTargetType) HasTweakTarget() bool {
	return t.IsValid()
}
