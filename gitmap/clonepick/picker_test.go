package clonepick

// picker_test.go: keyboard + state coverage for the bubbletea picker
// without launching a real terminal. Each test drives pickerModel
// directly via handleKey -- the same entry point Update() uses for
// tea.KeyMsg, so the assertions exercise the production code path.

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func newTestModel(paths, preselected []string) pickerModel {
	return newPickerModel(paths, preselected)
}

func keyMsg(s string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)}
}

func TestPreselectedPathsStartChecked(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, []string{"b"})
	if !m.picked[1] || m.picked[0] || m.picked[2] {
		t.Fatalf("preselected map wrong: %v", m.picked)
	}
}

func TestSpaceTogglesCurrentRow(t *testing.T) {
	m := newTestModel([]string{"a", "b"}, nil)
	next, _ := m.handleKey(keyMsg(" "))
	out := next.(pickerModel)
	if !out.picked[0] {
		t.Fatal("space did not toggle row 0 on")
	}
	again, _ := out.handleKey(keyMsg(" "))
	if again.(pickerModel).picked[0] {
		t.Fatal("second space did not toggle row 0 off")
	}
}

func TestDownKeyAdvancesCursorAndClampsAtEnd(t *testing.T) {
	m := newTestModel([]string{"a", "b"}, nil)
	step1, _ := m.handleKey(keyMsg("j"))
	step2, _ := step1.(pickerModel).handleKey(keyMsg("j"))
	if step2.(pickerModel).cursor != 1 {
		t.Fatalf("cursor = %d, want clamped at 1", step2.(pickerModel).cursor)
	}
}

func TestSelectAllSkipsAutoExcluded(t *testing.T) {
	paths := []string{"src/main.go", "node_modules/lib.js", "docs/README.md"}
	m := newTestModel(paths, nil)
	next, _ := m.handleKey(keyMsg("a"))
	out := next.(pickerModel)
	if !out.picked[0] || !out.picked[2] {
		t.Fatal("'a' should pick non-greyed rows")
	}
	if out.picked[1] {
		t.Fatal("'a' must not pick auto-greyed node_modules row")
	}
}

func TestSelectNoneClearsEverything(t *testing.T) {
	m := newTestModel([]string{"a", "b"}, []string{"a", "b"})
	next, _ := m.handleKey(keyMsg("n"))
	if len(next.(pickerModel).picked) != 0 {
		t.Fatalf("'n' should clear picks, got %v", next.(pickerModel).picked)
	}
}

func TestQuitFlagsCancelled(t *testing.T) {
	m := newTestModel([]string{"a"}, nil)
	next, cmd := m.handleKey(keyMsg("q"))
	if !next.(pickerModel).cancelled { //nolint:misspell // matches pickerModel.cancelled field spelling.
		t.Fatal("'q' should set the cancel flag")
	}
	if cmd == nil {
		t.Fatal("'q' should return tea.Quit cmd (non-nil)")
	}
}

func TestSaveFlagsDoneAndReturnsSelection(t *testing.T) {
	m := newTestModel([]string{"a", "b", "c"}, []string{"a", "c"})
	next, cmd := m.handleKey(keyMsg("s"))
	out := next.(pickerModel)
	if !out.done || cmd == nil {
		t.Fatalf("'s' should set done + return tea.Quit, got done=%v cmd=%v", out.done, cmd)
	}
	got := out.selected()
	if len(got) != 2 || got[0] != "a" || got[1] != "c" {
		t.Fatalf("selected() = %v, want [a c]", got)
	}
}

func TestIsAutoExcludedMatchesPrefixesNotSubstrings(t *testing.T) {
	cases := map[string]bool{
		"node_modules":            true,
		"node_modules/foo/bar":    true,
		"src/node_modules_helper": false, // substring, not a prefix
		"vendor/x.go":             true,
		"docs/README.md":          false,
	}
	for path, want := range cases {
		if got := IsAutoExcluded(path); got != want {
			t.Errorf("IsAutoExcluded(%q) = %v, want %v", path, got, want)
		}
	}
}

func TestViewIncludesCounterAndKeyHints(t *testing.T) {
	m := newTestModel([]string{"a", "b"}, []string{"a"})
	out := m.View()
	if !strings.Contains(out, "1/2 selected") {
		t.Errorf("view missing counter: %q", out)
	}
	if !strings.Contains(out, "space toggle") {
		t.Errorf("view missing key hints: %q", out)
	}
}

func TestViewHandlesEmptyRepoGracefully(t *testing.T) {
	m := newTestModel(nil, nil)
	out := m.View()
	if !strings.Contains(out, "no tracked files") {
		t.Errorf("empty-repo view wrong: %q", out)
	}
}
