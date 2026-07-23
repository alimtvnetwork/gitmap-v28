package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

type groupsModel struct {
	groups  []model.Group
	cursor  int
	message string
}

func newGroupsModel(groups []model.Group) groupsModel {
	return groupsModel{groups: groups}
}

func (m groupsModel) Update(msg tea.Msg) (groupsModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case keys.down(keyMsg):
		if m.cursor < len(m.groups)-1 {
			m.cursor++
		}
	case keys.up(keyMsg):
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.create(keyMsg):
		m.message = "Group creation — use CLI: gitmap group create <name>"
	case keys.delete(keyMsg):
		if len(m.groups) > 0 {
			name := m.groups[m.cursor].Name
			m.message = fmt.Sprintf(constants.TUIConfirmDelete, name)
		}
	}

	return m, nil
}

func (m groupsModel) View() string {
	if len(m.groups) == 0 {
		return styleHint.Render(constants.TUINoGroups)
	}

	var b strings.Builder

	header := fmt.Sprintf("  %-4s %-20s %-30s",
		"", constants.TUIColGroup, "Description")
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, g := range m.groups {
		desc := g.Description
		if desc == "" {
			desc = "—"
		}
		line := fmt.Sprintf("%-20s %s", g.Name, desc)

		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(styleSearch.Render(m.message))
	}

	return b.String()
}
