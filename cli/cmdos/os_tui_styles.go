package cmdos

import "github.com/charmbracelet/lipgloss"

var (
	tuiColorPrimary   = lipgloss.Color("#3ddc84")
	tuiColorSecondary = lipgloss.Color("#888888")
	tuiColorSuccess   = lipgloss.Color("#50fa7b")
	tuiColorDanger    = lipgloss.Color("#ff5555")
	tuiColorAccent    = lipgloss.Color("#8be9fd")

	tuiStyleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(tuiColorPrimary)

	tuiStyleTabActive = lipgloss.NewStyle().
				Bold(true).
				Foreground(tuiColorPrimary).
				Underline(true)

	tuiStyleTabInactive = lipgloss.NewStyle().
				Foreground(tuiColorSecondary)

	tuiStyleCursor = lipgloss.NewStyle().
			Bold(true).
			Foreground(tuiColorAccent)

	tuiStyleNormal = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#f8f8f2"))

	tuiStyleBadge = lipgloss.NewStyle().
			Foreground(tuiColorSecondary)

	tuiStyleFooter = lipgloss.NewStyle().
			Foreground(lipgloss.Color("#a6adc8"))
)
