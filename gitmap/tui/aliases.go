package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type aliasesModel struct {
	db      *store.DB
	aliases []store.AliasWithRepo
	cursor  int
	message string
}

func newAliasesModel(db *store.DB) aliasesModel {
	var aliases []store.AliasWithRepo
	if db != nil {
		aliases, _ = db.ListAliasesWithRepo()
	}

	return aliasesModel{db: db, aliases: aliases}
}

func (m aliasesModel) Update(msg tea.Msg) (aliasesModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case keys.down(keyMsg):
		if m.cursor < len(m.aliases)-1 {
			m.cursor++
		}
	case keys.up(keyMsg):
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.delete(keyMsg):
		if len(m.aliases) > 0 {
			name := m.aliases[m.cursor].Alias.Alias
			m.message = fmt.Sprintf(constants.TUIAliasDeleteHint, name, name)
		}
	case keys.refresh(keyMsg):
		m.aliases, _ = m.db.ListAliasesWithRepo()
		m.message = constants.TUIAliasRefreshed
	case keys.create(keyMsg):
		m.message = constants.TUIAliasCreateHint
	}

	return m, nil
}

func (m aliasesModel) View() string {
	if len(m.aliases) == 0 {
		return styleHint.Render(constants.TUIAliasEmpty)
	}

	var b strings.Builder

	header := fmt.Sprintf("  %-4s %-15s %-25s %s",
		"", "Alias", "Slug", "Path")
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, a := range m.aliases {
		path := truncatePath(a.AbsolutePath, 35)
		line := fmt.Sprintf("%-15s %-25s %s", a.Alias.Alias, a.Slug, path)

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
