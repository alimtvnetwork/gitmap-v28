package tui

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

func testRepos() []model.ScanRecord {
	return []model.ScanRecord{
		{Slug: "alpha", Branch: "main", AbsolutePath: "/repos/alpha"},
		{Slug: "beta", Branch: "dev", AbsolutePath: "/repos/beta"},
		{Slug: "gamma", Branch: "main", AbsolutePath: "/repos/gamma"},
	}
}

func testGroups() []model.Group {
	return []model.Group{
		{Name: "frontend", Description: "UI repos"},
		{Name: "backend", Description: "API repos"},
	}
}

func keyMsg(key string) tea.KeyMsg {
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(key)}
}

func specialKey(k tea.KeyType) tea.KeyMsg {
	return tea.KeyMsg{Type: k}
}

// --- Root model: tab switching ---

func TestTabSwitching(t *testing.T) {
	m := newRootModel(nil, testRepos(), testGroups(), model.DefaultConfig())

	if m.activeTab != viewBrowser {
		t.Fatalf("expected initial tab %d, got %d", viewBrowser, m.activeTab)
	}

	result, _ := m.Update(specialKey(tea.KeyTab))
	m = result.(rootModel)

	if m.activeTab != viewActions {
		t.Fatalf("expected tab %d after first Tab, got %d", viewActions, m.activeTab)
	}

	// Tab wraps around
	for i := 0; i < viewCount-1; i++ {
		result, _ = m.Update(specialKey(tea.KeyTab))
		m = result.(rootModel)
	}

	if m.activeTab != viewBrowser {
		t.Fatalf("expected tab to wrap to %d, got %d", viewBrowser, m.activeTab)
	}
}

// --- Browser: cursor navigation ---

func TestBrowserCursorDown(t *testing.T) {
	b := newBrowserModel(testRepos())

	if b.cursor != 0 {
		t.Fatalf("expected cursor 0, got %d", b.cursor)
	}

	b, _ = b.Update(keyMsg("j"))
	if b.cursor != 1 {
		t.Fatalf("expected cursor 1 after j, got %d", b.cursor)
	}

	b, _ = b.Update(keyMsg("j"))
	b, _ = b.Update(keyMsg("j")) // past end
	if b.cursor != 2 {
		t.Fatalf("expected cursor clamped at 2, got %d", b.cursor)
	}
}

func TestBrowserCursorUp(t *testing.T) {
	b := newBrowserModel(testRepos())
	b, _ = b.Update(keyMsg("j"))
	b, _ = b.Update(keyMsg("k"))

	if b.cursor != 0 {
		t.Fatalf("expected cursor 0 after k, got %d", b.cursor)
	}

	b, _ = b.Update(keyMsg("k")) // past start
	if b.cursor != 0 {
		t.Fatalf("expected cursor clamped at 0, got %d", b.cursor)
	}
}

// --- Browser: selection ---

func TestBrowserSelection(t *testing.T) {
	b := newBrowserModel(testRepos())

	b, _ = b.Update(keyMsg(" "))
	sel := b.selected()

	if len(sel) != 1 {
		t.Fatalf("expected 1 selected, got %d", len(sel))
	}
	if sel[0].Slug != "alpha" {
		t.Fatalf("expected alpha selected, got %s", sel[0].Slug)
	}

	// Toggle off
	b, _ = b.Update(keyMsg(" "))
	if len(b.selected()) != 0 {
		t.Fatalf("expected 0 selected after toggle, got %d", len(b.selected()))
	}
}

func TestBrowserSelectAll(t *testing.T) {
	b := newBrowserModel(testRepos())

	b, _ = b.Update(keyMsg("a"))

	if len(b.selected()) != 3 {
		t.Fatalf("expected 3 selected, got %d", len(b.selected()))
	}

	// Deselect all
	b, _ = b.Update(keyMsg("a"))
	if len(b.selected()) != 0 {
		t.Fatalf("expected 0 after deselect all, got %d", len(b.selected()))
	}
}

// --- Browser: search ---

func TestBrowserSearch(t *testing.T) {
	b := newBrowserModel(testRepos())

	b, _ = b.Update(keyMsg("/"))
	if !b.searching {
		t.Fatal("expected searching mode")
	}

	b, _ = b.Update(keyMsg("a"))
	b, _ = b.Update(keyMsg("l"))
	b, _ = b.Update(keyMsg("p"))

	if len(b.filtered) != 1 {
		t.Fatalf("expected 1 match for 'alp', got %d", len(b.filtered))
	}
	if b.filtered[0].Slug != "alpha" {
		t.Fatalf("expected alpha, got %s", b.filtered[0].Slug)
	}

	// Exit search
	b, _ = b.Update(specialKey(tea.KeyEsc))
	if b.searching {
		t.Fatal("expected search mode off after Esc")
	}
}

// --- Groups: navigation ---

func TestGroupsCursorNav(t *testing.T) {
	g := newGroupsModel(testGroups())

	if g.cursor != 0 {
		t.Fatalf("expected cursor 0, got %d", g.cursor)
	}

	g, _ = g.Update(keyMsg("j"))
	if g.cursor != 1 {
		t.Fatalf("expected cursor 1, got %d", g.cursor)
	}

	g, _ = g.Update(keyMsg("j")) // clamp
	if g.cursor != 1 {
		t.Fatalf("expected cursor clamped at 1, got %d", g.cursor)
	}

	g, _ = g.Update(keyMsg("k"))
	if g.cursor != 0 {
		t.Fatalf("expected cursor 0 after k, got %d", g.cursor)
	}
}

// --- Releases: navigation and detail toggle ---

func TestReleasesDetailToggle(t *testing.T) {
	r := releasesModel{
		releases: []model.ReleaseRecord{
			{Version: "v1.0.0", Tag: "v1.0.0"},
		},
	}

	r, _ = r.Update(specialKey(tea.KeyEnter))
	if !r.detail {
		t.Fatal("expected detail on after Enter")
	}

	r, _ = r.Update(specialKey(tea.KeyEnter))
	if r.detail {
		t.Fatal("expected detail off after second Enter")
	}
}

// --- Releases: trigger activation ---

func TestReleaseTrigger(t *testing.T) {
	r := releasesModel{
		releases: []model.ReleaseRecord{
			{Version: "v1.0.0"},
		},
		trigger: newRelTriggerModel(),
	}

	r, _ = r.Update(keyMsg("n"))
	if !r.trigger.active {
		t.Fatal("expected trigger active after n")
	}

	// Navigate to patch and confirm
	r.trigger, _ = r.trigger.Update(specialKey(tea.KeyEnter))
	if !r.trigger.confirmed {
		t.Fatal("expected trigger confirmed")
	}

	cmd := r.trigger.buildCommand()
	if cmd == "" {
		t.Fatal("expected non-empty command")
	}
}

// --- Empty state rendering ---

func TestBrowserEmptyView(t *testing.T) {
	b := newBrowserModel(nil)
	view := b.View()

	if view == "" {
		t.Fatal("expected non-empty view for empty browser")
	}
}

func TestGroupsEmptyView(t *testing.T) {
	g := newGroupsModel(nil)
	view := g.View()

	if view == "" {
		t.Fatal("expected non-empty view for empty groups")
	}
}
