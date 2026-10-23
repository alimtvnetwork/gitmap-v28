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

// TokyoTag renders a key-value pair formatted with Tokyo Night tokens.
// In non-TTY or no-color environments, returns "[label: value]" or "[label]".
func TokyoTag(label, value string) string {
	if !colorAllowed() {
		return formatPlainTag(label, value)
	}

	return formatColoredTag(label, value)
}

func formatPlainTag(label, value string) string {
	if value == "" {
		return fmt.Sprintf("[%s]", label)
	}

	return fmt.Sprintf("[%s: %s]", label, value)
}

func formatColoredTag(label, value string) string {
	lbl := renderTagLabel(label)
	if value == "" {
		return lbl
	}

	val := renderTagValue(value)

	return lbl + val
}

func renderTagLabel(label string) string {
	style := SafeTokyoStyle(TokyoFG).
		Bold(true).
		Background(lipgloss.Color(TokyoSurface)).
		Padding(0, 1)

	return style.Render(label)
}

func renderTagValue(value string) string {
	style := SafeTokyoStyle(TokyoPrimary).
		Background(lipgloss.Color(TokyoBase)).
		Padding(0, 1)

	return style.Render(value)
}

func resolveMetricColor(label string, count int) string {
	lower := strings.ToLower(label)
	if lower == "healthy" || lower == "pass" || lower == "ok" {
		return TokyoOK
	}
	if count > 0 && (lower == "issues" || lower == "errors" || lower == "fail") {
		return TokyoWarn
	}

	return TokyoPrimary
}

// TokyoMetricPill renders a labeled count metric with Tokyo Night styling.
// In non-TTY or no-color environments, returns "[label: count]".
func TokyoMetricPill(label string, count int) string {
	if !colorAllowed() {
		return fmt.Sprintf("[%s: %d]", label, count)
	}

	color := resolveMetricColor(label, count)
	style := SafeTokyoStyle(color).
		Bold(true).
		Background(lipgloss.Color(TokyoSurface)).
		Padding(0, 1)

	return style.Render(fmt.Sprintf("%s: %d", label, count))
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

// TokyoBadge renders text inside a styled badge with custom foreground and background colors.
// When color is disallowed or stdout is not a TTY, returns plain "[label]".
func TokyoBadge(label, bgHex, fgHex string) string {
	if label == "" {
		return ""
	}
	if !colorAllowed() {
		return fmt.Sprintf("[%s]", label)
	}

	return renderStyledBadge(label, bgHex, fgHex)
}

func renderStyledBadge(label, bgHex, fgHex string) string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(fgHex)).
		Background(lipgloss.Color(bgHex)).
		Padding(0, 1)

	return style.Render(label)
}

// TokyoHeader renders title inside a bold Tokyo Night section header banner.
// When color is disallowed, returns uppercase title without ANSI styling.
func TokyoHeader(title string) string {
	trimmed := strings.TrimSpace(title)
	if trimmed == "" {
		return ""
	}
	if !colorAllowed() {
		return strings.ToUpper(trimmed)
	}

	return renderStyledHeader(strings.ToUpper(trimmed))
}

func renderStyledHeader(upperTitle string) string {
	style := lipgloss.NewStyle().
		Bold(true).
		Foreground(lipgloss.Color(TokyoPrimary)).
		MarginBottom(1)

	return style.Render(upperTitle)
}
