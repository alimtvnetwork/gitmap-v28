package cmdos

// TweakTargetType defines system tweak categories.
type TweakTargetType string

const (
	TweakContextMenuType TweakTargetType = "context-menu"
	TweakStartMenuType   TweakTargetType = "start-menu"
	TweakPowerSchemeType TweakTargetType = "power"
	TweakHibernateType   TweakTargetType = "hibernate"
)

// TweakStatus models the observed state of system tweaks.
type TweakStatus struct {
	IsClassicContextMenu bool
	IsClassicStartMenu   bool
	IsUltimatePower      bool
	IsHibernateEnabled   bool
}

// TweakEngine abstracts tweak application and inspection.
type TweakEngine interface {
	SetContextMenu(isClassic bool) error
	SetStartMenu(isClassic bool) error
	SetPowerScheme(isUltimate bool) error
	SetHibernate(isEnabled bool) error
	GetStatus() (TweakStatus, error)
}
