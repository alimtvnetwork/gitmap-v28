package termout

import (
	"github.com/charmbracelet/lipgloss"
)

var (
	TokyoHeaderStyle       = SafeTokyoStyle(TokyoPrimary).Bold(true)
	TokyoAccentStyle       = SafeTokyoStyle(TokyoAccent)
	TokyoActiveBadgeStyle  = SafeTokyoStyle(TokyoOK).Bold(true)
	TokyoMissingBadgeStyle = SafeTokyoStyle(TokyoError).Bold(true)
	TokyoWarnBadgeStyle    = SafeTokyoStyle(TokyoWarn).Bold(true)
	TokyoMutedStyle        = SafeTokyoStyle(TokyoMuted)
)

// SafeTokyoStyle returns a Lipgloss style configured with foreground color fgHex,
// falling back to an unstyled empty style when colors are disallowed or stdout is not a TTY.
func SafeTokyoStyle(fgHex string) lipgloss.Style {
	if !colorAllowed() {
		return lipgloss.NewStyle()
	}

	return lipgloss.NewStyle().Foreground(lipgloss.Color(fgHex))
}

// SafeAdaptiveStyle constructs an adaptive Lipgloss style supporting light and dark backgrounds,
// with strict TTY validation.
func SafeAdaptiveStyle(light, dark string) lipgloss.Style {
	if !colorAllowed() {
		return lipgloss.NewStyle()
	}

	return lipgloss.NewStyle().Foreground(lipgloss.AdaptiveColor{
		Dark:  dark,
		Light: light,
	})
}
