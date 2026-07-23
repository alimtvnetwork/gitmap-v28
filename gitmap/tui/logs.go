package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type logsModel struct {
	db        *store.DB
	entries   []model.CommandHistoryRecord
	filtered  []model.CommandHistoryRecord
	cursor    int
	detail    bool
	searching bool
	query     string
}

func newLogsModel(db *store.DB) logsModel {
	var entries []model.CommandHistoryRecord
	if db != nil {
		entries, _ = db.ListHistory()
	}

	return logsModel{
		db:       db,
		entries:  entries,
		filtered: entries,
	}
}

func (m logsModel) Update(msg tea.Msg) (logsModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.searching {
		return m.updateSearch(keyMsg), nil
	}

	return m.handleKey(keyMsg), nil
}

func (m logsModel) updateSearch(msg tea.KeyMsg) logsModel {
	switch msg.String() {
	case "esc":
		m.searching = false
		m.query = ""
		m.filtered = m.entries
		m.cursor = 0
	case "enter":
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

	return m
}

func (m *logsModel) applyFilter() {
	if m.query == "" {
		m.filtered = m.entries
		m.cursor = 0

		return
	}

	q := strings.ToLower(m.query)
	var results []model.CommandHistoryRecord

	for _, e := range m.entries {
		if matchesLogQuery(e, q) {
			results = append(results, e)
		}
	}

	m.filtered = results
	m.cursor = 0
}

func matchesLogQuery(e model.CommandHistoryRecord, q string) bool {
	if strings.Contains(strings.ToLower(e.Command), q) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Alias), q) {
		return true
	}
	if strings.Contains(strings.ToLower(e.Args), q) {
		return true
	}
	if strings.Contains(fmt.Sprintf("%d", e.ExitCode), q) {
		return true
	}

	return false
}

func (m logsModel) handleKey(msg tea.KeyMsg) logsModel {
	max := len(m.filtered) - 1
	if max < 0 {
		if keys.search(msg) {
			m.searching = true
		}

		return m
	}

	switch {
	case keys.down(msg):
		if m.cursor < max {
			m.cursor++
		}
	case keys.up(msg):
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.enter(msg):
		m.detail = !m.detail
	case keys.refresh(msg):
		m.entries, _ = m.db.ListHistory()
		m.applyFilter()
	case keys.search(msg):
		m.searching = true
	}

	return m
}
