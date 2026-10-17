package cmdupdate

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ExtractReleaseNotesSummary parses a markdown release body into 1-2 clean bullet points.
func ExtractReleaseNotesSummary(body string, maxLines int) []string {
	if maxLines <= 0 {
		maxLines = 2
	}
	cleanBody := strings.ReplaceAll(body, "\r\n", "\n")
	lines := strings.Split(cleanBody, "\n")
	var results []string
	inCodeBlock := false

	for _, raw := range lines {
		if len(results) >= maxLines {
			break
		}
		trimmed := strings.TrimSpace(raw)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock || isMarkdownIgnoredLine(trimmed) {
			continue
		}
		text, isBullet := extractBulletText(trimmed)
		if !isBullet || len(text) == 0 {
			continue
		}
		results = append(results, clampNoteLine(text, 72))
	}
	if len(results) == 0 {
		results = fallbackExtractPlainLines(lines, maxLines)
	}
	return results
}

func fallbackExtractPlainLines(lines []string, maxLines int) []string {
	var results []string
	inCodeBlock := false
	for _, raw := range lines {
		if len(results) >= maxLines {
			break
		}
		trimmed := strings.TrimSpace(raw)
		if strings.HasPrefix(trimmed, "```") {
			inCodeBlock = !inCodeBlock
			continue
		}
		if inCodeBlock || isMarkdownIgnoredLine(trimmed) || len(trimmed) == 0 {
			continue
		}
		results = append(results, clampNoteLine(trimmed, 72))
	}
	return results
}

func isMarkdownIgnoredLine(line string) bool {
	if len(line) == 0 {
		return true
	}
	if strings.HasPrefix(line, "#") || strings.HasPrefix(line, ">") {
		return true
	}
	if strings.HasPrefix(line, "---") || strings.HasPrefix(line, "___") || strings.HasPrefix(line, "***") {
		return true
	}
	return false
}

func extractBulletText(line string) (string, bool) {
	if strings.HasPrefix(line, "- ") {
		return strings.TrimSpace(line[2:]), true
	}
	if strings.HasPrefix(line, "* ") {
		return strings.TrimSpace(line[2:]), true
	}
	if strings.HasPrefix(line, "• ") {
		return strings.TrimSpace(strings.TrimPrefix(line, "• ")), true
	}
	dotIdx := strings.Index(line, ". ")
	if dotIdx > 0 && dotIdx <= 3 && isAllDigits(line[:dotIdx]) {
		return strings.TrimSpace(line[dotIdx+2:]), true
	}
	return "", false
}

func isAllDigits(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < '0' || s[i] > '9' {
			return false
		}
	}
	return true
}

func clampNoteLine(s string, maxLen int) string {
	runes := []rune(s)
	if len(runes) <= maxLen {
		return s
	}
	if maxLen <= 3 {
		return string(runes[:maxLen])
	}
	return string(runes[:maxLen-3]) + "..."
}

// FormatReleaseNotesLines renders notes for styled lipgloss or plain output.
func FormatReleaseNotesLines(notes []string, isPlain bool) []string {
	if isPlain {
		return formatPlainReleaseNotes(notes)
	}
	return formatStyledReleaseNotes(notes)
}

func formatPlainReleaseNotes(notes []string) []string {
	var out []string
	for _, note := range notes {
		out = append(out, fmt.Sprintf("    - %s", note))
	}
	return out
}

func formatStyledReleaseNotes(notes []string) []string {
	var out []string
	bulletStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(DraculaGray))
	textStyle := lipgloss.NewStyle().Foreground(lipgloss.Color(DraculaFg))
	for _, note := range notes {
		line := fmt.Sprintf("%s%s", bulletStyle.Render("  • "), textStyle.Render(note))
		out = append(out, line)
	}
	return out
}

func appendSummaryNotes(lines []string, notes []string, lbl lipgloss.Style) []string {
	lines = append(lines, "")
	lines = append(lines, lbl.Render("Release Notes:"))
	formatted := FormatReleaseNotesLines(notes, false)
	return append(lines, formatted...)
}

func printPlainReleaseNotes(notes []string) {
	if len(notes) == 0 {
		return
	}
	fmt.Println("  Release Notes:")
	plainNotes := FormatReleaseNotesLines(notes, true)
	for _, note := range plainNotes {
		fmt.Println(note)
	}
}

func formatSummaryLines(result UpdateResult) []string {
	title := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DraculaGreen)).Render("✔ Gitmap CLI Updated Successfully")
	lbl := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color(DraculaCyan))
	prevVer := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f8f8f2")).Render(result.PreviousVersion)
	arrow := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#f1fa8c")).Render("→")
	currVer := lipgloss.NewStyle().Bold(true).Foreground(lipgloss.Color("#ffff55")).Render(result.CurrentVersion)
	delta := fmt.Sprintf("%s  %s  %s", prevVer, arrow, currVer)
	sec := float64(result.DurationMs) / 1000.0
	lines := []string{
		title,
		"",
		fmt.Sprintf("%s : %s", lbl.Render("Version   "), delta),
		fmt.Sprintf("%s : %s", lbl.Render("Binary    "), result.BinaryPath),
		fmt.Sprintf("%s : SHA256 verified against checksums.txt", lbl.Render("Integrity ")),
		fmt.Sprintf("%s : %s", lbl.Render("Engine    "), result.DownloaderEngine),
		fmt.Sprintf("%s : %.1fs", lbl.Render("Elapsed   "), sec),
	}
	if len(result.ReleaseNotes) > 0 {
		lines = appendSummaryNotes(lines, result.ReleaseNotes, lbl)
	}
	return lines
}

// MakeProgressCallback creates a throttled detail updater for UpdateUIRenderer.
func MakeProgressCallback(ui UpdateUIRenderer) func(written, total int64) {
	if ui == nil {
		return nil
	}
	return func(written, total int64) {
		if total > 0 {
			pct := float64(written) / float64(total) * 100
			ui.OnStageProgress(fmt.Sprintf("%.1f / %.1f MB (%.0f%%)", float64(written)/(1024*1024), float64(total)/(1024*1024), pct))
			return
		}
		ui.OnStageProgress(fmt.Sprintf("%.1f MB", float64(written)/(1024*1024)))
	}
}
