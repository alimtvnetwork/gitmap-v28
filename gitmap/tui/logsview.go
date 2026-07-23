package tui

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
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

func (m logsModel) viewList() string {
	var b strings.Builder

	if m.searching {
		b.WriteString(styleSearch.Render(constants.TUISearchPrompt + m.query + "█"))
		b.WriteString("\n")
	} else if len(m.query) > 0 {
		b.WriteString(styleHint.Render(fmt.Sprintf(constants.TUILogFilterActive, m.query, len(m.filtered))))
		b.WriteString("\n")
	}

	header := fmt.Sprintf("  %-4s %-16s %-10s %-30s %-10s %-6s %s",
		"", constants.TUIColCommand, constants.TUIColAlias,
		constants.TUIColArgs, constants.TUIColDuration,
		constants.TUIColExit, constants.TUIColDate)
	b.WriteString(styleHeader.Render(header))
	b.WriteString("\n")

	if len(m.filtered) == 0 {
		b.WriteString(styleHint.Render(constants.TUILogNoMatch))
		b.WriteString("\n")
	}

	for i, e := range m.filtered {
		line := formatLogRow(e)
		if i == m.cursor {
			b.WriteString(styleCursorRow.Render("> " + line))
		} else {
			b.WriteString(styleNormalRow.Render("  " + line))
		}
		b.WriteString("\n")
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render(fmt.Sprintf("  %d log(s)  •  enter: detail  •  r: refresh  •  /: filter", len(m.filtered))))

	return b.String()
}

func (m logsModel) viewDetail() string {
	e := m.filtered[m.cursor]

	var b strings.Builder
	b.WriteString(styleGroupName.Render(fmt.Sprintf("  Command: %s", e.Command)))
	b.WriteString("\n\n")

	writeField(&b, "Alias", e.Alias)
	writeField(&b, "Args", e.Args)
	writeField(&b, "Flags", e.Flags)
	writeField(&b, "Started", e.StartedAt)
	writeField(&b, "Finished", e.FinishedAt)
	writeField(&b, "Duration", formatDurationMs(e.DurationMs))
	writeField(&b, "Exit Code", fmt.Sprintf("%d", e.ExitCode))
	writeField(&b, "Repo Count", fmt.Sprintf("%d", e.RepoCount))

	if len(e.Summary) > 0 {
		b.WriteString("\n")
		writeField(&b, "Summary", e.Summary)
	}

	b.WriteString("\n")
	b.WriteString(styleHint.Render("  enter: back to list"))

	return b.String()
}
