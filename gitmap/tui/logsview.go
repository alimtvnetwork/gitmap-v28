package tui

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func (m logsModel) View() string {
	if len(m.entries) == 0 {
		return styleHint.Render(constants.TUILogEmpty)
	}

	if m.detail && m.cursor < len(m.filtered) {
		return m.viewDetail()
	}

	return m.viewList()
}

func (m logsModel) writeSearchOrFilterHeader(b *strings.Builder) {
	if m.searching {
		b.WriteString(styleSearch.Render(constants.TUISearchPrompt+m.query+"█") + "\n")
	} else if len(m.query) > 0 {
		b.WriteString(styleHint.Render(fmt.Sprintf(constants.TUILogFilterActive, m.query, len(m.filtered))) + "\n")
	}
}

func logListHeader() string {
	return fmt.Sprintf("  %-4s %-16s %-10s %-30s %-10s %-6s %s",
		"", constants.TUIColCommand, constants.TUIColAlias,
		constants.TUIColArgs, constants.TUIColDuration,
		constants.TUIColExit, constants.TUIColDate)
}

func renderLogRows(b *strings.Builder, filtered []model.CommandHistoryRecord, cursor int) {
	if len(filtered) == 0 {
		b.WriteString(styleHint.Render(constants.TUILogNoMatch) + "\n")
		return
	}

	for i, e := range filtered {
		line := formatLogRow(e)
		if i == cursor {
			b.WriteString(styleCursorRow.Render("> " + line) + "\n")
		} else {
			b.WriteString(styleNormalRow.Render("  " + line) + "\n")
		}
	}
}

func (m logsModel) viewList() string {
	var b strings.Builder
	m.writeSearchOrFilterHeader(&b)
	b.WriteString(styleHeader.Render(logListHeader()) + "\n")
	renderLogRows(&b, m.filtered, m.cursor)
	b.WriteString("\n" + styleHint.Render(fmt.Sprintf("  %d log(s)  •  enter: detail  •  r: refresh  •  /: filter", len(m.filtered))))

	return b.String()
}

func writeLogDetailFields(b *strings.Builder, e model.CommandHistoryRecord) {
	writeField(b, "Alias", e.Alias)
	writeField(b, "Args", e.Args)
	writeField(b, "Flags", e.Flags)
	writeField(b, "Started", e.StartedAt)
	writeField(b, "Finished", e.FinishedAt)
	writeField(b, "Duration", formatDurationMs(e.DurationMs))
	writeField(b, "Exit Code", fmt.Sprintf("%d", e.ExitCode))
	writeField(b, "Repo Count", fmt.Sprintf("%d", e.RepoCount))
	if len(e.Summary) > 0 {
		b.WriteString("\n")
		writeField(b, "Summary", e.Summary)
	}
}

func (m logsModel) viewDetail() string {
	e := m.filtered[m.cursor]
	var b strings.Builder
	b.WriteString(styleGroupName.Render(fmt.Sprintf("  Command: %s", e.Command)) + "\n\n")
	writeLogDetailFields(&b, e)
	b.WriteString("\n" + styleHint.Render("  enter: back to list"))

	return b.String()
}
