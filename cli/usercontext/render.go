package usercontext

import (
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

const cardBoxWidth = 68

// RenderUserInfoCard renders an elegant Catppuccin Macchiato rounded box.
func RenderUserInfoCard(w io.Writer, summary UserSummary) error {
	writeCardTop(w, "Git & GitHub Identity Context", cardBoxWidth)
	writeCardRow(w, formatGhAuthLine(summary.GhAuth), cardBoxWidth)
	writeCardRow(w, formatActiveGitLine(summary), cardBoxWidth)
	writeCardRow(w, formatGlobalGitLine(summary.GlobalGitUser), cardBoxWidth)
	writeCardRow(w, formatProfileLine(summary.ActiveProfile), cardBoxWidth)
	writeCardRow(w, formatBindingLine(summary.BoundProfile), cardBoxWidth)
	writeCardBottom(w, cardBoxWidth)

	return nil
}

func writeCardTop(w io.Writer, title string, width int) {
	header := fmt.Sprintf("╭─ %s%s%s ", constants.ColorPastelMagenta, title, constants.ColorPastelCyan)
	plain := fmt.Sprintf("╭─ %s ", title)
	rem := width - len(plain) - 1
	if rem < 2 {
		rem = 2
	}
	fmt.Fprintf(w, "%s%s%s╮%s\n", constants.ColorPastelCyan, header, strings.Repeat("─", rem), constants.ColorReset)
}

func writeCardBottom(w io.Writer, width int) {
	fmt.Fprintf(w, "%s╰%s╯%s\n", constants.ColorPastelCyan, strings.Repeat("─", width-2), constants.ColorReset)
}

func writeCardRow(w io.Writer, content string, width int) {
	cleanLen := len([]rune(stripAnsi(content)))
	pad := width - 4 - cleanLen
	if pad < 0 {
		pad = 0
	}
	fmt.Fprintf(w, "%s│%s %s%s %s│%s\n", constants.ColorPastelCyan, constants.ColorReset, content, strings.Repeat(" ", pad), constants.ColorPastelCyan, constants.ColorReset)
}

func stripAnsi(s string) string {
	var b strings.Builder
	inSeq := false
	for _, r := range s {
		if r == '\033' {
			inSeq = true
			continue
		}
		if inSeq && r == 'm' {
			inSeq = false
			continue
		}
		if !inSeq {
			b.WriteRune(r)
		}
	}

	return b.String()
}
