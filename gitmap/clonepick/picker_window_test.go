package clonepick

// picker_window_test.go: coverage for the windowed scroller added
// in v4.22.0 (cursor stays in view, page keys jump by viewport,
// View() only renders the visible window).

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func makeRows(n int) []string {
	out := make([]string, n)
	for i := 0; i < n; i++ {
		out[i] = stringFromInt(i)
	}

	return out
}

// stringFromInt avoids strconv import noise in test paths -- the
// values are only used as unique row labels.
func stringFromInt(i int) string {
	const digits = "0123456789"
	if i < 10 {
		return string(digits[i])
	}
	if i < 100 {
		return string(digits[i/10]) + string(digits[i%10])
	}

	return string(digits[i/100]) + string(digits[(i/10)%10]) + string(digits[i%10])
}

func TestClampScrollKeepsCursorInWindow(t *testing.T) {
	cases := []struct {
		cursor, offset, height, total, want int
	}{
		{cursor: 0, offset: 0, height: 5, total: 100, want: 0},
		{cursor: 4, offset: 0, height: 5, total: 100, want: 0},
		{cursor: 5, offset: 0, height: 5, total: 100, want: 1}, // scrolled down
		{cursor: 99, offset: 0, height: 5, total: 100, want: 95},
		{cursor: 2, offset: 10, height: 5, total: 100, want: 2}, // jumped up
		{cursor: 0, offset: 0, height: 5, total: 0, want: 0},    // empty
	}
	for _, tc := range cases {
		got := clampScroll(tc.cursor, tc.offset, tc.height, tc.total)
		if got != tc.want {
			t.Errorf("clampScroll(c=%d,o=%d,h=%d,t=%d) = %d, want %d",
				tc.cursor, tc.offset, tc.height, tc.total, got, tc.want)
		}
	}
}

func TestWindowSizeMsgUpdatesViewportHeight(t *testing.T) {
	m := newTestModel(makeRows(50), nil)
	next, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 12})
	got := next.(pickerModel).viewportHeight
	want := 12 - chromeRows
	if got != want {
		t.Fatalf("viewportHeight = %d, want %d", got, want)
	}
}

func TestPageDownAdvancesByViewportHeight(t *testing.T) {
	m := newTestModel(makeRows(100), nil)
	resized, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 13})
	mr := resized.(pickerModel)
	next, _ := mr.handleKey(tea.KeyMsg{Type: tea.KeyPgDown})
	got := next.(pickerModel).cursor
	if got != mr.viewportHeight {
		t.Fatalf("pgdn cursor = %d, want %d", got, mr.viewportHeight)
	}
}

func TestEndKeyJumpsToLastRow(t *testing.T) {
	m := newTestModel(makeRows(40), nil)
	next, _ := m.handleKey(keyMsg("G"))
	if got := next.(pickerModel).cursor; got != 39 {
		t.Fatalf("G cursor = %d, want 39", got)
	}
}

func TestViewRendersOnlyVisibleRows(t *testing.T) {
	m := newTestModel(makeRows(500), nil)
	resized, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 23})
	out := resized.(pickerModel).View()
	// Output shape: header\n + N row\n + blank\n + footer\n
	// => total newlines = N + 3.
	rowLines := strings.Count(out, "\n") - 3
	want := 23 - chromeRows
	if rowLines != want {
		t.Fatalf("rendered %d row lines, want %d (out=%q)",
			rowLines, want, out)
	}
}

func TestViewHeaderShowsRowRange(t *testing.T) {
	m := newTestModel(makeRows(50), nil)
	resized, _ := m.Update(tea.WindowSizeMsg{Width: 80, Height: 13})
	out := resized.(pickerModel).View()
	if !strings.Contains(out, "rows 1-10") {
		t.Fatalf("view missing row-range indicator: %q", out)
	}
}
