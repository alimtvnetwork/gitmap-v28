package constants

// Global `--theme` flag — selects which terminal color palette
// gitmap renders. Resolved once at startup (see gitmap/theme) and
// applied to all stdout / stderr output via an ANSI rewrite filter,
// so no individual Msg* constant has to know which palette is live.
const (
	FlagTheme = "theme"

	EnvTheme = "GITMAP_THEME"

	ThemeBright     = "bright"
	ThemeStandard   = "standard"
	ThemeMonochrome = "monochrome"
	ThemeMono       = "mono" // short alias for monochrome
	ThemeDefault    = ThemeBright
)
