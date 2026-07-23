package tui

import "github.com/charmbracelet/lipgloss"

// Color palette derived from terminal green accent.
var (
	colorPrimary   = lipgloss.Color("#3ddc84")
	colorSecondary = lipgloss.Color("#888888")
	colorMuted     = lipgloss.Color("#555555")
	colorDanger    = lipgloss.Color("#ff5555")
	colorSuccess   = lipgloss.Color("#50fa7b")
	colorSelected  = lipgloss.Color("#2a2a4e")
)

// Layout styles.
var (
	styleTitle = lipgloss.NewStyle().
			Bold(true).
			Foreground(colorPrimary).
			MarginBottom(1)

	styleTab = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(colorSecondary)

	styleActiveTab = lipgloss.NewStyle().
			Padding(0, 2).
			Foreground(colorPrimary).
			Bold(true).
			Underline(true)

	styleStatusBar = lipgloss.NewStyle().
			Foreground(colorMuted).
			MarginTop(1)

	styleSelectedRow = lipgloss.NewStyle().
				Background(colorSelected).
				Foreground(colorPrimary)

	styleCursorRow = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	styleNormalRow = lipgloss.NewStyle().
			Foreground(colorSecondary)

	styleHeader = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true).
			Underline(true)

	styleSearch = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	styleDirty = lipgloss.NewStyle().
			Foreground(colorDanger).
			Bold(true)

	styleClean = lipgloss.NewStyle().
			Foreground(colorSuccess)

	styleGroupName = lipgloss.NewStyle().
			Foreground(colorPrimary).
			Bold(true)

	styleHint = lipgloss.NewStyle().
			Foreground(colorMuted).
			Italic(true)
)
