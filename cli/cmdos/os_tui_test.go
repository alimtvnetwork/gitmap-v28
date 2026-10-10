package cmdos

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestOSTUIModel_Init(t *testing.T) {
	m := NewOSTUIModel(true)
	if len(m.Tabs) != 6 {
		t.Fatalf("expected 6 tabs, got %d", len(m.Tabs))
	}
	if m.ActiveTab != 0 {
		t.Fatalf("expected active tab 0, got %d", m.ActiveTab)
	}
	if m.Mode != OSTUIViewModeNavType {
		t.Fatalf("expected nav mode, got %s", m.Mode)
	}
}

func TestOSTUIModel_Navigation(t *testing.T) {
	m := NewOSTUIModel(true)

	// Tab switch right
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyTab})
	m = newM.(OSTUIModel)
	if m.ActiveTab != 1 {
		t.Fatalf("expected active tab 1 after Tab, got %d", m.ActiveTab)
	}

	// Number key jump
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("4")})
	m = newM.(OSTUIModel)
	if m.ActiveTab != 3 {
		t.Fatalf("expected active tab 3 after pressing '4', got %d", m.ActiveTab)
	}

	// Down arrow moves cursor
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyDown})
	m = newM.(OSTUIModel)
	if m.Cursor != 1 {
		t.Fatalf("expected cursor 1 after Down, got %d", m.Cursor)
	}
}

func TestOSTUIModel_ToggleAndSelect(t *testing.T) {
	m := NewOSTUIModel(true)

	// Space toggles item 0
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeySpace})
	m = newM.(OSTUIModel)
	if !m.Tabs[0].Items[0].IsChecked {
		t.Fatalf("expected item 0 to be checked after space")
	}

	// 'n' clears all
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("n")})
	m = newM.(OSTUIModel)
	if countTabChecked(m.Tabs[0]) != 0 {
		t.Fatalf("expected 0 checked items after 'n'")
	}

	// 'a' checks all
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("a")})
	m = newM.(OSTUIModel)
	if countTabChecked(m.Tabs[0]) != len(m.Tabs[0].Items) {
		t.Fatalf("expected all items checked after 'a'")
	}
}

func TestOSTUIModel_DryRunExecution(t *testing.T) {
	m := NewOSTUIModel(true)
	// Check one item
	m.Tabs[0].Items[0].IsChecked = true

	// Press Enter to trigger execution
	newM, _ := m.Update(tea.KeyMsg{Type: tea.KeyEnter})
	m = newM.(OSTUIModel)

	if m.Mode != OSTUIViewModeDoneType {
		t.Fatalf("expected mode done after execution, got %s", m.Mode)
	}
	if len(m.Results) != 1 {
		t.Fatalf("expected 1 execution result, got %d", len(m.Results))
	}
	if m.Results[0].IsFail() {
		t.Fatalf("expected dry-run result to be success")
	}

	// Press Esc to return to Nav mode
	newM, _ = m.Update(tea.KeyMsg{Type: tea.KeyEsc})
	m = newM.(OSTUIModel)
	if m.Mode != OSTUIViewModeNavType {
		t.Fatalf("expected nav mode after esc, got %s", m.Mode)
	}
}

func TestOSTUIModel_ViewOutput(t *testing.T) {
	m := NewOSTUIModel(true)
	navView := m.View()

	if !strings.Contains(navView, "GITMAP OS DASHBOARD") {
		t.Fatalf("nav view missing header")
	}
	if !strings.Contains(navView, "Tweaks") {
		t.Fatalf("nav view missing Tweaks tab")
	}

	m.Mode = OSTUIViewModeDoneType
	m.Results = []OSTUIActionResult{
		{Title: "Sample Action", IsSuccess: true, Message: "done"},
	}
	doneView := m.View()
	if !strings.Contains(doneView, "EXECUTION RESULTS") {
		t.Fatalf("done view missing execution header")
	}
	if !strings.Contains(doneView, "Sample Action") {
		t.Fatalf("done view missing result title")
	}
}
