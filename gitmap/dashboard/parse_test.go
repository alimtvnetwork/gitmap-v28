package dashboard

import (
	"testing"
)

func TestParseCommitLog(t *testing.T) {
	raw := "abcd1234ef5678|abcd123|Alice|alice@example.com|2026-08-20T10:00:00Z|feat: initial|parent1\n 1 2 file.txt"
	commits := parseCommitLog(raw)
	if len(commits) != 1 {
		t.Fatalf("expected 1 commit, got %d", len(commits))
	}
	c := commits[0]
	if c.SHA != "abcd1234ef5678" || c.ShortSHA != "abcd123" || c.Author != "Alice" {
		t.Errorf("unexpected commit info: %+v", c)
	}
	if c.FilesChanged != 1 || c.Insertions != 1 || c.Deletions != 2 {
		t.Errorf("unexpected numstat counts: %+v", c)
	}
}

func TestParseNumstat(t *testing.T) {
	lines := []string{"10\t5\ta.go", "2\t1\tb.go"}
	files, ins, del := parseNumstat(lines)
	if files != 2 || ins != 12 || del != 6 {
		t.Errorf("parseNumstat returned (%d, %d, %d)", files, ins, del)
	}
}

func TestParseBranchLines(t *testing.T) {
	lines := []string{"main|abcd123|2026-08-20", "origin/feat|efgh456|2026-08-19"}
	branches := parseBranchLines(lines)
	if len(branches) != 2 {
		t.Fatalf("expected 2 branches, got %d", len(branches))
	}
	if branches[0].IsRemote || !branches[1].IsRemote {
		t.Errorf("branch remote detection incorrect")
	}
}

func TestParseTagLines(t *testing.T) {
	lines := []string{"v1.0.0|abcd123|2026-08-18"}
	tags := parseTagLines(".", lines)
	if len(tags) != 1 || tags[0].Name != "v1.0.0" {
		t.Errorf("parseTagLines returned unexpected tags: %+v", tags)
	}
}
