package termout

// DefaultSubcommandIndicator is the arrow symbol shown for subcommands.
const DefaultSubcommandIndicator = "➔"

// Theme holds color constants for help rendering. Every slot is filled
// from a semantic Role in the active design-system palette (themes.go) —
// never from literal color codes — so UI help screens follow --theme.
type Theme struct {
	BannerColor  string
	HeaderColor  string
	CommandColor string
	DescColor    string
	HintColor    string
	TipColor     string
}

// DefaultTheme returns the help-menu palette resolved from the active
// design-system theme via semantic roles:
//
//	BannerColor  → accent    the single Von Restorff accent per region
//	HeaderColor  → primary   section headings (bold)
//	CommandColor → ok        command tokens read as "do this"
//	DescColor    → fg        body text
//	HintColor    → fg-muted  subcommand hints, secondary cues
//	TipColor     → accent    tips are highlighted callouts
//
// On a non-TTY or under the mono theme every slot is "" (RoleSGR() emits
// nothing), so piped help stays byte-clean without a strip pass.
func DefaultTheme() Theme {
	return Theme{
		BannerColor:  RoleSGR(RoleAccent),
		HeaderColor:  RoleSGRBold(RolePrimary),
		CommandColor: RoleSGR(RoleOK),
		DescColor:    RoleSGR(RoleFG),
		HintColor:    RoleSGR(RoleFGMuted),
		TipColor:     RoleSGR(RoleAccent),
	}
}
