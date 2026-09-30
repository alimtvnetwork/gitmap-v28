package cmdos

import (
	"fmt"
	"strings"
)

func renderDoneView(m OSTUIModel) string {
	var b strings.Builder
	header := fmt.Sprintf("GITMAP OS DASHBOARD - EXECUTION RESULTS (%d executed)\n", len(m.Results))
	b.WriteString(tuiStyleTitle.Render(header))
	b.WriteString(strings.Repeat("─", 80) + "\n")

	for _, res := range m.Results {
		b.WriteString(renderResultLine(res))
		b.WriteByte('\n')
	}

	b.WriteString(strings.Repeat("─", 80) + "\n")
	b.WriteString(tuiStyleFooter.Render("[Enter / Esc / Backspace] Return to Dashboard | [q] Quit\n"))

	return b.String()
}

func renderResultLine(r OSTUIActionResult) string {
	glyph := lipglossGlyph(r.IsSuccess)
	return fmt.Sprintf("  %s %-32s: %s", glyph, r.Title, r.Message)
}

func lipglossGlyph(isSuccess bool) string {
	if isSuccess {
		return tuiStyleTitle.Foreground(tuiColorSuccess).Render("✔")
	}

	return tuiStyleTitle.Foreground(tuiColorDanger).Render("✖")
}
