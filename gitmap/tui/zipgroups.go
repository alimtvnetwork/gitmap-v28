package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type zipGroupsModel struct {
	db      *store.DB
	groups  []store.ZipGroupWithCount
	cursor  int
	detail  string
	message string
}

func newZipGroupsModel(db *store.DB) zipGroupsModel {
	var groups []store.ZipGroupWithCount
	if db != nil {
		groups, _ = db.ListZipGroupsWithCount()
	}

	return zipGroupsModel{db: db, groups: groups}
}

func (m zipGroupsModel) Update(msg tea.Msg) (zipGroupsModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	switch {
	case keys.down(keyMsg):
		if m.cursor < len(m.groups)-1 {
			m.cursor++
		}
		m.detail = ""
	case keys.up(keyMsg):
		if m.cursor > 0 {
			m.cursor--
		}
		m.detail = ""
	case keys.enter(keyMsg):
		m.detail = m.loadDetail()
	case keys.refresh(keyMsg):
		m.groups, _ = m.db.ListZipGroupsWithCount()
		m.message = constants.TUIZGRefreshed
	case keys.create(keyMsg):
		m.message = constants.TUIZGCreateHint
	case keys.delete(keyMsg):
		if len(m.groups) > 0 {
			name := m.groups[m.cursor].Name
			m.message = fmt.Sprintf(constants.TUIConfirmDelete, name)
		}
	}

	return m, nil
}

func (m zipGroupsModel) loadDetail() string {
	if len(m.groups) == 0 {
		return ""
	}

	g := m.groups[m.cursor]
	items, err := m.db.ListZipGroupItems(g.Name)
	if err != nil {
		return fmt.Sprintf("  Error: %v", err)
	}

	var b strings.Builder
	b.WriteString(fmt.Sprintf("  %s (%d items):\n", g.Name, len(items)))

	for _, item := range items {
		if item.IsFolder {
			b.WriteString(fmt.Sprintf("    📁 %s\n", item.RelativePath))
		} else {
			b.WriteString(fmt.Sprintf("    📄 %s\n", item.RelativePath))
		}
	}

	if len(g.ArchiveName) > 0 {
		b.WriteString(fmt.Sprintf("  Archive: %s\n", g.ArchiveName))
	}

	return b.String()
}

func (m zipGroupsModel) View() string {
	if len(m.groups) == 0 {
		return styleHint.Render(constants.TUIZGEmpty)
	}

	var b strings.Builder

	header := fmt.Sprintf("  %-4s %-20s %-8s %s",
		"", "Name", "Items", "Archive")
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, g := range m.groups {
		archive := g.ArchiveName
		if len(archive) == 0 {
			archive = g.Name + ".zip"
		}

		line := fmt.Sprintf("%-20s %-8d %s", g.Name, g.ItemCount, archive)

		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	if m.detail != "" {
		b.WriteString("\n")
		b.WriteString(styleNormalRow.Render(m.detail))
	}

	if m.message != "" {
		b.WriteString("\n")
		b.WriteString(styleSearch.Render(m.message))
	}

	return b.String()
}
