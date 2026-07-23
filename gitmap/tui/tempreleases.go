package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// tempReleasesModel manages the TUI view for temp-release branches.
type tempReleasesModel struct {
	db       *store.DB
	records  []model.TempRelease
	cursor   int
	detail   bool
	groupBy  bool
	filter   string
	filtered []model.TempRelease
}

func newTempReleasesModel(db *store.DB) tempReleasesModel {
	var records []model.TempRelease
	if db != nil {
		records, _ = db.ListTempReleases()
	}

	return tempReleasesModel{
		db:       db,
		records:  records,
		filtered: records,
	}
}

func (m tempReleasesModel) Update(msg tea.Msg) (tempReleasesModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	return m.handleKey(keyMsg), nil
}

func (m tempReleasesModel) handleKey(msg tea.KeyMsg) tempReleasesModel {
	max := len(m.filtered) - 1

	switch {
	case keys.down(msg):
		if max >= 0 && m.cursor < max {
			m.cursor++
		}
	case keys.up(msg):
		if m.cursor > 0 {
			m.cursor--
		}
	case keys.enter(msg):
		if max >= 0 {
			m.detail = !m.detail
		}
	case keys.refresh(msg):
		m.records, _ = m.db.ListTempReleases()
		m.filtered = filterTRByPrefix(m.records, m.filter)
		m.cursor = 0
	case msg.String() == "g":
		m.groupBy = !m.groupBy
	}

	return m
}

func (m tempReleasesModel) View() string {
	if len(m.records) == 0 {
		return styleHint.Render(constants.TUITREmpty)
	}

	if m.detail && m.cursor < len(m.filtered) {
		return m.viewDetail()
	}

	if m.groupBy {
		return m.viewGrouped()
	}

	return m.viewList()
}

func (m tempReleasesModel) viewList() string {
	var b strings.Builder

	header := fmt.Sprintf("  %-4s %-28s %-12s %-5s %-10s %s",
		"", constants.TUIColTRBranch, constants.TUIColTRPrefix,
		constants.TUIColTRSeq, constants.TUIColTRCommit, constants.TUIColDate)
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, r := range m.filtered {
		line := formatTRRow(r)
		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(fmt.Sprintf(
		"  %d branch(es)  •  enter: detail  •  g: group by prefix  •  r: refresh",
		len(m.filtered))))

	return b.String()
}

func (m tempReleasesModel) viewGrouped() string {
	groups := groupTRByPrefix(m.filtered)

	var b strings.Builder
	b.WriteString(styleHeader.Render(fmt.Sprintf("  %-20s %-8s %-12s %s",
		"Prefix", "Count", "Seq Range", "Latest Commit")))
	b.WriteString("\n")

	for _, g := range groups {
		line := fmt.Sprintf("  %-20s %-8d %-12s %s",
			g.prefix, g.count, g.seqRange, shortSHA(g.latestCommit))
		b.WriteString(styleNormalRow.Render(line))
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(fmt.Sprintf(
		"  %d prefix group(s)  •  g: flat list  •  r: refresh", len(groups))))

	return b.String()
}

func (m tempReleasesModel) viewDetail() string {
	r := m.filtered[m.cursor]

	var b strings.Builder
	b.WriteString(styleGroupName.Render(fmt.Sprintf("  Temp Release: %s", r.Branch)))
	b.WriteString("\n\n")

	writeField(&b, "Branch", r.Branch)
	writeField(&b, "Version Prefix", r.VersionPrefix)
	writeField(&b, "Sequence", fmt.Sprintf("%d", r.SequenceNumber))
	writeField(&b, "Commit", shortSHA(r.CommitSha))
	writeField(&b, "Message", truncateStr(r.CommitMessage, 60))
	writeField(&b, "Created", r.CreatedAt)

	b.WriteString("\n")
	b.WriteString(styleHint.Render("  enter: back to list"))

	return b.String()
}
