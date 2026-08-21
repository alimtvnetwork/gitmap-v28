package tui

import (
	"fmt"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

type releasesModel struct {
	db       *store.DB
	releases []model.ReleaseRecord
	cursor   int
	detail   bool
	trigger  relTriggerModel
}

func newReleasesModel(db *store.DB) releasesModel {
	var releases []model.ReleaseRecord
	if db != nil {
		releases, _ = db.ListReleases()
	}

	return releasesModel{
		db:       db,
		releases: releases,
		trigger:  newRelTriggerModel(),
	}
}

func (m *releasesModel) updateTrigger(msg tea.Msg) bool {
	if !m.trigger.active {
		return false
	}
	m.trigger, _ = m.trigger.Update(msg)
	if !m.trigger.active {
		m.trigger = newRelTriggerModel()
	}

	return true
}

func (m releasesModel) Update(msg tea.Msg) (releasesModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}
	if m.updateTrigger(msg) {
		return m, nil
	}

	return m.handleKey(keyMsg), nil
}

func (m *releasesModel) handleNavKey(msg tea.KeyMsg, max int) {
	switch {
	case keys.down(msg) && max >= 0 && m.cursor < max:
		m.cursor++
	case keys.up(msg) && m.cursor > 0:
		m.cursor--
	case keys.enter(msg) && max >= 0:
		m.detail = !m.detail
	}
}

func (m releasesModel) handleKey(msg tea.KeyMsg) releasesModel {
	if msg.String() == "n" {
		m.trigger.active = true
	} else if keys.refresh(msg) && m.db != nil {
		m.releases, _ = m.db.ListReleases()
	} else {
		m.handleNavKey(msg, len(m.releases)-1)
	}

	return m
}

func (m releasesModel) View() string {
	if m.trigger.active {
		return m.trigger.View()
	}
	if len(m.releases) == 0 {
		return styleHint.Render(constants.TUIRelEmpty)
	}
	if m.detail && m.cursor < len(m.releases) {
		return m.viewDetail()
	}

	return m.viewList()
}

func relListHeader() string {
	return fmt.Sprintf("  %-4s %-12s %-14s %-20s %-8s %-8s %-8s %s",
		"", constants.TUIColVersion, constants.TUIColTag,
		constants.TUIColBranch, constants.TUIColDraft,
		constants.TUIColLatest, constants.TUIColSource, constants.TUIColDate)
}

func renderRelRows(b *strings.Builder, releases []model.ReleaseRecord, cursor int) {
	for i, r := range releases {
		line := formatRelRow(r)
		if i == cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}
}

func (m releasesModel) viewList() string {
	var b strings.Builder
	b.WriteString(styleHeader.Render(relListHeader()) + "\n")
	renderRelRows(&b, m.releases, m.cursor)
	b.WriteString("\n" + styleHint.Render(fmt.Sprintf("  %d release(s)  •  enter: detail  •  r: refresh", len(m.releases))))

	return b.String()
}

func writeRelDetailFields(b *strings.Builder, r model.ReleaseRecord) {
	writeField(b, "Tag", r.Tag)
	writeField(b, "Branch", r.Branch)
	writeField(b, "Source Branch", r.SourceBranch)
	writeField(b, "Commit", shortSHA(r.CommitSha))
	writeField(b, "Source", r.Source)
	writeField(b, "Date", r.CreatedAt)
	writeField(b, "Draft", boolLabel(r.IsDraft))
	writeField(b, "Pre-release", boolLabel(r.IsPreRelease))
	writeField(b, "Latest", boolLabel(r.IsLatest))
	if len(r.Notes) > 0 {
		b.WriteString("\n")
		writeField(b, "Notes", r.Notes)
	}
}

func writeRelChangelog(b *strings.Builder, changelog string) {
	if len(changelog) == 0 {
		return
	}
	b.WriteString("\n  Changelog:\n")
	for _, line := range strings.Split(changelog, "\n") {
		b.WriteString(styleHint.Render("    "+line) + "\n")
	}
}

func (m releasesModel) viewDetail() string {
	r := m.releases[m.cursor]
	var b strings.Builder
	b.WriteString(styleGroupName.Render(fmt.Sprintf("  Release %s", r.Version)) + "\n\n")
	writeRelDetailFields(&b, r)
	writeRelChangelog(&b, r.Changelog)
	b.WriteString("\n" + styleHint.Render("  enter: back to list"))

	return b.String()
}
