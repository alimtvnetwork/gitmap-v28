package clonepick

// picker_view.go: View() implementation for the bubbletea picker.
// Kept in its own file so picker.go can stay under 200 lines and the
// rendering loop has room to breathe.

import (
	"fmt"
	"strings"
)

// View renders the picker as a single string per tea-model contract.
// Only the visible window (viewportHeight rows starting at
// scrollOffset) is rendered so large repos don't blow past the
// terminal. Header shows "selected/total" plus an indicator of
// where in the list the window sits.
func (m pickerModel) View() string {
	if len(m.paths) == 0 {
		return "clone-pick: repository has no tracked files\n"
	}
	var b strings.Builder
	fmt.Fprintf(&b, "gitmap clone-pick --ask  (%d/%d selected, rows %d-%d)\n",
		m.countPicked(), len(m.paths),
		m.firstVisible()+1, m.lastVisible()+1)
	m.renderRows(&b)
	b.WriteString("\nspace toggle | a all | n none | pgup/pgdn page | g/G home/end | s save | q quit\n")

	return b.String()
}

// firstVisible / lastVisible report the inclusive 0-based bounds of
// the row window currently rendered. Used in the header so the user
// can see "rows 41-60 of 487" without counting.
func (m pickerModel) firstVisible() int { return m.scrollOffset }

func (m pickerModel) lastVisible() int {
	last := m.scrollOffset + m.viewportHeight - 1
	if last >= len(m.paths) {
		last = len(m.paths) - 1
	}

	return last
}

// renderRows writes one line per visible path. The window is
// [scrollOffset, scrollOffset+viewportHeight) clamped to len(paths)
// so we never read past the slice.
func (m pickerModel) renderRows(b *strings.Builder) {
	end := m.scrollOffset + m.viewportHeight
	if end > len(m.paths) {
		end = len(m.paths)
	}
	for i := m.scrollOffset; i < end; i++ {
		b.WriteString(formatRow(i == m.cursor, m.picked[i],
			IsAutoExcluded(m.paths[i]), m.paths[i]))
		b.WriteByte('\n')
	}
}

// formatRow returns the single-line representation of one picker
// entry. Cursor row gets a leading ">", everything else gets two
// spaces so columns line up.
func formatRow(isCursor, isPicked, isGreyed bool, path string) string {
	prefix := "  "
	if isCursor {
		prefix = "> "
	}
	mark := pickMark(isPicked, isGreyed)
	suffix := ""
	if isGreyed {
		suffix = "  (auto-greyed)"
	}

	return prefix + mark + " " + path + suffix
}

// pickMark returns the bracketed checkbox glyph for the row state.
// Greyed rows use "[-]" so the user can still see they're toggleable
// (versus the disabled-looking "[ ]").
func pickMark(isPicked, isGreyed bool) string {
	switch {
	case isPicked:
		return "[x]"
	case isGreyed:
		return "[-]"
	default:
		return "[ ]"
	}
}

// countPicked is the header counter. O(rows) once per render, fine
// for the row counts we expect (<10k entries before the windowed
// scroller lands).
func (m pickerModel) countPicked() int {
	n := 0
	for _, on := range m.picked {
		if on {
			n++
		}
	}

	return n
}
