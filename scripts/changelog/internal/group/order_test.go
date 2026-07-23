package group

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v27/scripts/changelog/internal/gitlog"
)

// TestByPrefixSortsItemsLexicographically pins the deterministic
// item-ordering contract: items inside each section come out
// alphabetically regardless of the order git happened to emit them.
func TestByPrefixSortsItemsLexicographically(t *testing.T) {
	commits := []gitlog.Commit{
		{Hash: "aaa", Subject: "feat: zeta feature"},
		{Hash: "bbb", Subject: "feat: alpha feature"},
		{Hash: "ccc", Subject: "feat: mu feature"},
		{Hash: "ddd", Subject: "fix: zulu bug"},
		{Hash: "eee", Subject: "fix: alpha bug"},
	}

	sections, _ := ByPrefix(commits)
	if len(sections) != 2 {
		t.Fatalf("want 2 sections, got %d", len(sections))
	}

	wantAdded := []string{"alpha feature", "mu feature", "zeta feature"}
	wantFixed := []string{"alpha bug", "zulu bug"}

	for i, want := range wantAdded {
		if sections[0].Items[i] != want {
			t.Fatalf("Added[%d]: want %q got %q", i, want, sections[0].Items[i])
		}
	}

	for i, want := range wantFixed {
		if sections[1].Items[i] != want {
			t.Fatalf("Fixed[%d]: want %q got %q", i, want, sections[1].Items[i])
		}
	}
}

// TestByPrefixIsStableAcrossInputPermutations regenerates the section
// list from two different input orders and asserts byte-identical
// output. This is the full determinism guarantee for the renderer.
func TestByPrefixIsStableAcrossInputPermutations(t *testing.T) {
	a := []gitlog.Commit{
		{Hash: "1", Subject: "feat: one"},
		{Hash: "2", Subject: "fix: two"},
		{Hash: "3", Subject: "feat: three"},
	}

	b := []gitlog.Commit{
		{Hash: "3", Subject: "feat: three"},
		{Hash: "1", Subject: "feat: one"},
		{Hash: "2", Subject: "fix: two"},
	}

	sa, _ := ByPrefix(a)
	sb, _ := ByPrefix(b)

	if len(sa) != len(sb) {
		t.Fatalf("section count differs: %d vs %d", len(sa), len(sb))
	}

	for i := range sa {
		if sa[i].Name != sb[i].Name {
			t.Fatalf("section[%d] name differs", i)
		}

		if len(sa[i].Items) != len(sb[i].Items) {
			t.Fatalf("section[%d] item count differs", i)
		}

		for j := range sa[i].Items {
			if sa[i].Items[j] != sb[i].Items[j] {
				t.Fatalf("section[%d].Items[%d]: %q vs %q", i, j, sa[i].Items[j], sb[i].Items[j])
			}
		}
	}
}
