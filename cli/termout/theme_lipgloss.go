package termout

import (
	"fmt"
	"strings"

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

// TokyoBannerStyle renders title inside a bold styled banner.
func TokyoBannerStyle(title string) string {
	if !colorAllowed() {
		return title
	}

	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(TokyoPrimary)).
		Padding(0, 1)

	return style.Render(title)
}

// TokyoBoxStyle renders content within a rounded border box.
func TokyoBoxStyle(content string) string {
	if !colorAllowed() {
		return content
	}

	style := lipgloss.NewStyle().
		BorderStyle(lipgloss.RoundedBorder()).
		BorderForeground(lipgloss.Color(TokyoMuted)).
		Padding(0, 1)

	return style.Render(content)
}

// TokyoPill renders text inside a styled pill badge with safe plain fallback.
func TokyoPill(text string, style lipgloss.Style) string {
	if !colorAllowed() {
		return fmt.Sprintf("[%s]", text)
	}

	return style.Padding(0, 1).Render(text)
}

// TokyoOKPill renders a success status pill badge.
func TokyoOKPill(text string) string {
	style := SafeTokyoStyle(TokyoOK).
		Background(lipgloss.Color(TokyoSurface))

	return TokyoPill(text, style)
}

// TokyoWarnPill renders a warning status pill badge.
func TokyoWarnPill(text string) string {
	style := SafeTokyoStyle(TokyoWarn).
		Background(lipgloss.Color(TokyoSurface))

	return TokyoPill(text, style)
}

// TokyoErrorPill renders an error status pill badge.
func TokyoErrorPill(text string) string {
	style := SafeTokyoStyle(TokyoError).
		Background(lipgloss.Color(TokyoSurface))

	return TokyoPill(text, style)
}

// RenderTokyoProgress renders an inline Tokyo Night progress bar.
func RenderTokyoProgress(current, total, barWidth int) string {
	if total <= 0 || barWidth <= 0 {
		return "[░░░░] 0%"
	}
	if current < 0 {
		current = 0
	}
	if current > total {
		current = total
	}

	return assembleProgressBar(current, total, barWidth)
}

func assembleProgressBar(current, total, barWidth int) string {
	filled := (current * barWidth) / total
	empty := barWidth - filled
	pct := (current * 100) / total
	filledStr := strings.Repeat("█", filled)
	emptyStr := strings.Repeat("░", empty)
	if colorAllowed() {
		filledStr = SafeTokyoStyle(TokyoOK).Render(filledStr)
		emptyStr = SafeTokyoStyle(TokyoMuted).Render(emptyStr)
	}

	return fmt.Sprintf("[%s%s] %d%%", filledStr, emptyStr, pct)
}

