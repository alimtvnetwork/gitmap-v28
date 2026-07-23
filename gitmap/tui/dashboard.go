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

func refreshStatuses(repos []model.ScanRecord) tea.Cmd {
	return func() tea.Msg {
		entries := make([]statusEntry, 0, len(repos))

		for _, r := range repos {
			rs := gitutil.Status(r.AbsolutePath)
			entries = append(entries, statusEntry{
				Slug:      r.Slug,
				Branch:    rs.Branch,
				Status:    statusLabel(rs.Dirty, rs.Unreachable),
				Ahead:     rs.Ahead,
				Behind:    rs.Behind,
				Stash:     rs.StashCount,
				Untracked: rs.Untracked,
				Modified:  rs.Modified,
				Staged:    rs.Staged,
			})
		}

		return refreshMsg{entries: entries}
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
		m.entries = msg.entries
		m.loading = false

		return m, m.scheduleTick()
	case tickMsg:
		m.loading = true

		return m, refreshStatuses(m.repos)
	case tea.KeyMsg:
		return m.handleKey(msg)
	}

	return m, nil
}

func (m dashboardModel) handleKey(msg tea.KeyMsg) (dashboardModel, tea.Cmd) {
	max := len(m.entries) - 1
	if max < 0 {
		max = len(m.repos) - 1
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
	case keys.refresh(msg):
		m.loading = true

		return m, refreshStatuses(m.repos)
	}

	return m, nil
}

func (m dashboardModel) View() string {
	if len(m.repos) == 0 {
		return styleHint.Render(constants.TUINoRepos)
	}
	if m.loading {
		return styleHint.Render(constants.TUIRefreshing)
	}

	var b strings.Builder

	b.WriteString(styleHeader.Render(dashHeader()))
	b.WriteString("\n")

	for i, e := range m.entries {
		line := formatDashRow(e)
		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(dashSummary(m.entries)))

	return b.String()
}
