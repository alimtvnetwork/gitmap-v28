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

func (m releasesModel) Update(msg tea.Msg) (releasesModel, tea.Cmd) {
	keyMsg, ok := msg.(tea.KeyMsg)
	if !ok {
		return m, nil
	}

	if m.trigger.active {
		m.trigger, _ = m.trigger.Update(msg)
		if !m.trigger.active {
			m.trigger = newRelTriggerModel()
		}

		return m, nil
	}

	return m.handleKey(keyMsg), nil
}

func (m releasesModel) handleKey(msg tea.KeyMsg) releasesModel {
	max := len(m.releases) - 1

	switch {
	case msg.String() == "n":
		m.trigger.active = true
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
		m.releases, _ = m.db.ListReleases()
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

func (m releasesModel) viewList() string {
	var b strings.Builder

	header := fmt.Sprintf("  %-4s %-12s %-14s %-20s %-8s %-8s %-8s %s",
		"", constants.TUIColVersion, constants.TUIColTag,
		constants.TUIColBranch, constants.TUIColDraft,
		constants.TUIColLatest, constants.TUIColSource, constants.TUIColDate)
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	for i, r := range m.releases {
		line := formatRelRow(r)
		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(fmt.Sprintf("  %d release(s)  •  enter: detail  •  r: refresh", len(m.releases))))

	return b.String()
}

func (m releasesModel) viewDetail() string {
	r := m.releases[m.cursor]

	var b strings.Builder
	b.WriteString(styleGroupName.Render(fmt.Sprintf("  Release %s", r.Version)))
	b.WriteString("\n\n")

	writeField(&b, "Tag", r.Tag)
	writeField(&b, "Branch", r.Branch)
	writeField(&b, "Source Branch", r.SourceBranch)
	writeField(&b, "Commit", shortSHA(r.CommitSha))
	writeField(&b, "Source", r.Source)
	writeField(&b, "Date", r.CreatedAt)
	writeField(&b, "Draft", boolLabel(r.IsDraft))
	writeField(&b, "Pre-release", boolLabel(r.IsPreRelease))
	writeField(&b, "Latest", boolLabel(r.IsLatest))

	if len(r.Notes) > 0 {
		b.WriteString("\n")
		writeField(&b, "Notes", r.Notes)
	}

	if len(r.Changelog) > 0 {
		b.WriteString("\n  Changelog:\n")
		for _, line := range strings.Split(r.Changelog, "\n") {
			b.WriteString(styleHint.Render("    " + line))
			b.WriteString("\n")
		}
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render("  enter: back to list"))

	return b.String()
}
