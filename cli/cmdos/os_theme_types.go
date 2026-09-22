package cmdos

// ThemeModeType specifies desktop visual appearance.
type ThemeModeType string

const (
	ThemeModeDark  ThemeModeType = "dark"
	ThemeModeLight ThemeModeType = "light"
)

// ThemeOperator sets desktop theme.
type ThemeOperator interface {
	SetTheme(mode ThemeModeType) error
	GetTheme() (ThemeModeType, error)
}
