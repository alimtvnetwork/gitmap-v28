package tui

import (
	"strings"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// statusEntry holds computed git status for one repo.
type statusEntry struct {
	Slug      string
	Branch    string
	Status    string
	Ahead     int
	Behind    int
	Stash     int
	Untracked int
	Modified  int
	Staged    int
}

// refreshMsg carries freshly computed statuses.
type refreshMsg struct {
	entries []statusEntry
}

// tickMsg triggers a periodic auto-refresh.
type tickMsg struct{}

type dashboardModel struct {
	repos    []model.ScanRecord
	entries  []statusEntry
	cursor   int
	loading  bool
	interval time.Duration
}

func newDashboardModel(repos []model.ScanRecord, refreshSec int) dashboardModel {
	if refreshSec <= 0 {
		refreshSec = constants.DefaultDashboardRefresh
	}

	return dashboardModel{
		repos:    repos,
		loading:  true,
		interval: time.Duration(refreshSec) * time.Second,
	}
}

func (m dashboardModel) scheduleTick() tea.Cmd {
	return tea.Tick(m.interval, func(_ time.Time) tea.Msg {
		return tickMsg{}
	})
}

func makeStatusEntry(r model.ScanRecord) statusEntry {
	rs := gitutil.Status(r.AbsolutePath)

	return statusEntry{
		Slug:      r.Slug,
		Branch:    rs.Branch,
		Status:    statusLabel(rs.Dirty, rs.Unreachable),
		Ahead:     rs.Ahead,
		Behind:    rs.Behind,
		Stash:     rs.StashCount,
		Untracked: rs.Untracked,
		Modified:  rs.Modified,
		Staged:    rs.Staged,
	}
}

func collectStatusEntries(repos []model.ScanRecord) []statusEntry {
	entries := make([]statusEntry, 0, len(repos))
	for _, r := range repos {
		entries = append(entries, makeStatusEntry(r))
	}

	return entries
}

func refreshStatuses(repos []model.ScanRecord) tea.Cmd {
	return func() tea.Msg {
		return refreshMsg{entries: collectStatusEntries(repos)}
	}
}

func statusLabel(dirty, unreachable bool) string {
	if unreachable {
		return "error"
	}
	if dirty {
		return "dirty"
	}

	return "clean"
}

func (m dashboardModel) Init() tea.Cmd {
	return tea.Batch(refreshStatuses(m.repos), m.scheduleTick())
}

func (m dashboardModel) Update(msg tea.Msg) (dashboardModel, tea.Cmd) {
	switch msg := msg.(type) {
	case refreshMsg:
		m.entries, m.loading = msg.entries, false
		return m, m.scheduleTick()
	case tickMsg:
		m.loading = true
		return m, refreshStatuses(m.repos)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func maxIndex(a, b int) int {
	if a > 0 {
		return a - 1
	}

	return b - 1
}

func (m *dashboardModel) moveCursor(msg tea.KeyMsg, max int) {
	if keys.down(msg) && m.cursor < max {
		m.cursor++
	} else if keys.up(msg) && m.cursor > 0 {
		m.cursor--
	}
}

func (m dashboardModel) handleKey(msg tea.KeyMsg) (dashboardModel, tea.Cmd) {
	if keys.refresh(msg) {
		m.loading = true

		return m, refreshStatuses(m.repos)
	}
	m.moveCursor(msg, maxIndex(len(m.entries), len(m.repos)))

	return m, nil
}

func renderDashRows(b *strings.Builder, entries []statusEntry, cursor int) {
	for i, e := range entries {
		line := formatDashRow(e)
		if i == cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}
}

func (m dashboardModel) View() string {
	if len(m.repos) == 0 {
		return styleHint.Render(constants.TUINoRepos)
	}
	if m.loading {
		return styleHint.Render(constants.TUIRefreshing)
	}

	var b strings.Builder
	b.WriteString(styleHeader.Render(dashHeader()) + "\n")
	renderDashRows(&b, m.entries, m.cursor)
	b.WriteString("\n" + styleHint.Render(dashSummary(m.entries)))

	return b.String()
}
