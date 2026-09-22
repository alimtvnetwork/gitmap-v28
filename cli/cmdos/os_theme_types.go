package cmdos

// ThemeMode specifies desktop visual appearance.
type ThemeMode string

const (
	ThemeModeDark  ThemeMode = "dark"
	ThemeModeLight ThemeMode = "light"
)

// ThemeEngine sets desktop theme.
type ThemeEngine interface {
	SetTheme(mode ThemeMode) error
	GetTheme() (ThemeMode, error)
}
