package cmdos

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdos/enums"
)

// ThemeModeType specifies desktop visual appearance.
type ThemeModeType = enums.ThemeModeType

const (
	ThemeModeDark  = enums.ThemeModeDark
	ThemeModeLight = enums.ThemeModeLight
)

// ThemeOperator sets desktop termout.
type ThemeOperator interface {
	SetTheme(mode ThemeModeType) error
	GetTheme() (ThemeModeType, error)
}
