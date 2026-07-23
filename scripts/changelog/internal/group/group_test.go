package group

import (
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/scripts/changelog/internal/gitlog"
)

func TestByPrefixGroupsConventionalCommits(t *testing.T) {
	commits := []gitlog.Commit{
		{Hash: "a1", Subject: "feat(cli): add --dry-run"},
		{Hash: "b2", Subject: "fix: handle empty config"},
		{Hash: "c3", Subject: "docs: clarify install flow"},
		{Hash: "d4", Subject: "feat: second feature"},
		{Hash: "e5", Subject: "Changes"},                    // skipped
		{Hash: "f6", Subject: "refactor!: drop legacy api"}, // breaking marker
	}

	sections, skipped := ByPrefix(commits)

	if len(skipped) != 1 || skipped[0].Hash != "e5" {
		t.Fatalf("expected one skipped commit (e5), got %#v", skipped)
	}

	bySection := map[string][]string{}
	for _, s := range sections {
		bySection[s.Name] = s.Items
	}

	if got := bySection["Added"]; len(got) != 2 || got[0] != "add --dry-run" || got[1] != "second feature" {
		t.Fatalf("Added section mismatch: %#v", got)
	}
	if got := bySection["Fixed"]; len(got) != 1 || got[0] != "handle empty config" {
		t.Fatalf("Fixed section mismatch: %#v", got)
	}
	if got := bySection["Refactor"]; len(got) != 1 || got[0] != "drop legacy api" {
		t.Fatalf("Refactor section mismatch: %#v", got)
	}
}

func TestByPrefixPreservesSectionOrder(t *testing.T) {
	commits := []gitlog.Commit{
		{Hash: "a", Subject: "chore: tidy"},
		{Hash: "b", Subject: "feat: alpha"},
		{Hash: "c", Subject: "fix: bravo"},
	}

	sections, _ := ByPrefix(commits)

	if len(sections) != 3 {
		t.Fatalf("want 3 sections, got %d", len(sections))
	}
	if sections[0].Name != "Added" || sections[1].Name != "Fixed" || sections[2].Name != "Chore" {
		t.Fatalf("order mismatch: %v %v %v", sections[0].Name, sections[1].Name, sections[2].Name)
	}
}
