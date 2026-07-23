package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// relTriggerModel handles release initiation from the TUI.
type relTriggerModel struct {
	active    bool
	cursor    int
	version   string
	typing    bool
	confirmed bool
	options   []relOption
}

type relOption struct {
	label string
	flag  string
}

func newRelTriggerModel() relTriggerModel {
	return relTriggerModel{
		options: []relOption{
			{label: "Patch (0.0.x)", flag: "--bump patch"},
			{label: "Minor (0.x.0)", flag: "--bump minor"},
			{label: "Major (x.0.0)", flag: "--bump major"},
			{label: "Custom version", flag: ""},
		},
	}
}

func (m relTriggerModel) Update(msg tea.Msg) (relTriggerModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.typing {
		return m.handleTyping(keyMsg), nil
	}

	return m.handleNav(keyMsg), nil
}

func (m relTriggerModel) handleNav(msg tea.KeyMsg) relTriggerModel {
	switch {
	case keys.quit(msg):
		m.active = false
		m.confirmed = false
		m.version = ""
	case keys.up(msg):
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.down(msg):
		if m.cursor < len(m.options)-1 {
			m.cursor++
		}
	case keys.enter(msg):
		if m.cursor == len(m.options)-1 {
			m.typing = true
			m.version = ""
		} else {
			m.confirmed = true
		}
	}

	return m
}

func (m relTriggerModel) handleTyping(msg tea.KeyMsg) relTriggerModel {
	switch msg.String() {
	case "enter":
		if len(m.version) > 0 {
			m.confirmed = true
			m.typing = false
		}
	case "esc":
		m.typing = false
		m.version = ""
	case "backspace":
		if len(m.version) > 0 {
			m.version = m.version[:len(m.version)-1]
		}
	default:
		ch := msg.String()
		if len(ch) == 1 && isVersionChar(ch[0]) {
			m.version += ch
		}
	}

	return m
}

func isVersionChar(c byte) bool {
	return (c >= '0' && c <= '9') || c == '.' || c == 'v' || c == '-'
}

func (m relTriggerModel) buildCommand() string {
	if len(m.version) > 0 {
		return fmt.Sprintf(constants.TUIRelTriggerCmd, m.version)
	}

	if m.cursor < len(m.options) {
		return fmt.Sprintf(constants.TUIRelTriggerBumpCmd, m.options[m.cursor].flag)
	}

	return ""
}

func (m relTriggerModel) View() string {
	if m.confirmed {
		return m.viewConfirmed()
	}

	if m.typing {
		return m.viewTyping()
	}

	return m.viewMenu()
}

func (m relTriggerModel) viewMenu() string {
	var b strings.Builder
	b.WriteString(styleGroupName.Render(constants.TUIRelTriggerTitle))
	b.WriteString("\n\n")

	for i, opt := range m.options {
		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + opt.label))
		} else {
			b.WriteString(styleNormalRow.Render("  " + opt.label))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(constants.TUIRelTriggerNavHint))

	return b.String()
}

func (m relTriggerModel) viewTyping() string {
	var b strings.Builder
	b.WriteString(styleGroupName.Render(constants.TUIRelTriggerTitle))
	b.WriteString("\n\n")
	b.WriteString(styleNormalRow.Render(constants.TUIRelTriggerVerPrompt + m.version + "█"))
	b.WriteString("\n\n")
	b.WriteString(styleHint.Render(constants.TUIRelTriggerTypeHint))

	return b.String()
}

func (m relTriggerModel) viewConfirmed() string {
	cmd := m.buildCommand()

	var b strings.Builder
	b.WriteString(styleGroupName.Render(constants.TUIRelTriggerReady))
	b.WriteString("\n\n")
	b.WriteString(styleHeader.Render("  " + cmd))
	b.WriteString("\n\n")
	b.WriteString(styleHint.Render(constants.TUIRelTriggerRunHint))

	return b.String()
}
