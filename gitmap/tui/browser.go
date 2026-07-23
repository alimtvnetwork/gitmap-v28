package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/sahilm/fuzzy"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

type browserModel struct {
	repos      []model.ScanRecord
	filtered   []model.ScanRecord
	cursor     int
	selections map[int]bool
	searching  bool
	query      string
}

func newBrowserModel(repos []model.ScanRecord) browserModel {
	return browserModel{
		repos:      repos,
		filtered:   repos,
		selections: make(map[int]bool),
	}
}

func (m browserModel) Update(msg tea.Msg) (browserModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.searching {
		return m.updateSearch(keyMsg)
	}

	return m.updateBrowse(keyMsg), nil
}

func (m browserModel) updateSearch(msg tea.KeyMsg) (browserModel, tea.Cmd) {
	switch msg.String() {
	case "esc", "enter":
		m.searching = false
	case "backspace":
		if len(m.query) > 0 {
			m.query = m.query[:len(m.query)-1]
			m.applyFilter()
		}
	default:
		if len(msg.String()) == 1 {
			m.query += msg.String()
			m.applyFilter()
		}
	}

	return m, nil
}

func (m browserModel) updateBrowse(msg tea.KeyMsg) browserModel {
	switch {
	case keys.down(msg):
		if m.cursor < len(m.filtered)-1 {
			m.cursor++
		}
	case keys.up(msg):
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.selectItem(msg):
		m.selections[m.cursor] = !m.selections[m.cursor]
	case keys.selectAll(msg):
		allSelected := len(m.selections) == len(m.filtered)
		m.selections = make(map[int]bool)
		if !allSelected {
			for i := range m.filtered {
				m.selections[i] = true
			}
		}
	case keys.search(msg):
		m.searching = true
	}

	return m
}

func (m *browserModel) applyFilter() {
	if m.query == "" {
		m.filtered = m.repos
		m.cursor = 0

		return
	}

	slugs := make([]string, len(m.repos))
	for i, r := range m.repos {
		slugs[i] = r.Slug
	}

	matches := fuzzy.Find(m.query, slugs)
	m.filtered = make([]model.ScanRecord, len(matches))
	for i, match := range matches {
		m.filtered[i] = m.repos[match.Index]
	}

	m.cursor = 0
	m.selections = make(map[int]bool)
}

func (m browserModel) selected() []model.ScanRecord {
	var result []model.ScanRecord
	for i, sel := range m.selections {
		if sel && i < len(m.filtered) {
			result = append(result, m.filtered[i])
		}
	}

	return result
}

func (m browserModel) View() string {
	if len(m.repos) == 0 {
		return styleHint.Render(constants.TUINoRepos)
	}

	var b strings.Builder

	if m.searching {
		b.WriteString(styleSearch.Render(constants.TUISearchPrompt + m.query + "█"))
		b.WriteString("\n\n")
	} else if m.query != "" {
		b.WriteString(styleSearch.Render(constants.TUISearchPrompt + m.query))
		b.WriteString("\n\n")
	}

	header := fmt.Sprintf("  %-4s %-20s %-12s %s",
		"", constants.TUIColSlug, constants.TUIColBranch, constants.TUIColPath)
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, r := range m.filtered {
		marker := "  "
		if m.selections[i] {
			marker = "● "
		}

		line := fmt.Sprintf("%s%-20s %-12s %s", marker, r.Slug, r.Branch, truncatePath(r.AbsolutePath, 40))

		switch {
		case i == m.cursor:
			b.WriteString(styleCursorRow.Render("> " + line))
		case m.selections[i]:
			b.WriteString(styleSelectedRow.Render("  " + line))
		default:
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	sel := len(m.selected())
	if sel > 0 {
		b.WriteString(fmt.Sprintf("\n%d selected — press Tab for actions", sel))
	}

	return b.String()
}

func truncatePath(path string, max int) string {
	if len(path) <= max {
		return path
	}

	return "…" + path[len(path)-max+1:]
}
