package cmdpull

import (
	"bytes"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestResolveConciseRepoColWidth(t *testing.T) {
	// 1. Empty states should fall back to minConciseRepoColWidth (26)
	if width := ResolveConciseRepoColWidth(nil); width < 26 {
		t.Fatalf("expected width >= 26 for nil states, got %d", width)
	}

	// 2. Short repos under 26 characters should remain at 26
	shortStates := []*PullRepoState{
		{RepoName: "gstack"},
		{RepoName: "alim-cv-v8"},
	}
	if width := ResolveConciseRepoColWidth(shortStates); width != 26 {
		t.Fatalf("expected width 26 for short repos, got %d", width)
	}

	// 3. Long repo names should dynamically expand colWidth
	longName := "kita-social-media-content-calender-v2"
	mixedStates := []*PullRepoState{
		{RepoName: "gstack"},
		{RepoName: "ai-empathy-prompt-tuner-v1"},
		{RepoName: longName},
	}
	expectedLen := len(longName) // 37
	if width := ResolveConciseRepoColWidth(mixedStates); width < expectedLen {
		t.Fatalf("expected width >= %d for long repo, got %d", expectedLen, width)
	}

	// 4. Stable column width calculated from allRecords even if active repos are short
	allRecords := []model.ScanRecord{
		{RepoName: "short"},
		{RepoName: longName},
	}
	if width := ResolveConciseRepoColWidth(shortStates, allRecords); width != expectedLen {
		t.Fatalf("expected width %d from allRecords, got %d", expectedLen, width)
	}
}

func TestResolveRepoStatusLabel(t *testing.T) {
	tests := []struct {
		input string
		want  string
	}{
		{"", "up-to-date"},
		{"synced", "up-to-date"},
		{"dirty", "dirty"},
		{"failed", "failed"},
		{"+4941/-484 (68)", "+4941/-484 (68)"},
	}

	for _, tt := range tests {
		got := ResolveRepoStatusLabel(tt.input)
		if got != tt.want {
			t.Errorf("ResolveRepoStatusLabel(%q) = %q, want %q", tt.input, got, tt.want)
		}
	}
}

func TestFormatConciseActiveResultLine_VerticalAlignment(t *testing.T) {
	colWidth := 37
	line1 := FormatConciseActiveResultLine(colWidth, "alim-cv-v8", "dirty")
	line2 := FormatConciseActiveResultLine(colWidth, "ai-empathy-prompt-tuner-v1", "up-to-date")
	line3 := FormatConciseActiveResultLine(colWidth, "kita-social-media-content-calender-v2", "up-to-date")

	idx1 := strings.Index(line1, "dirty")
	idx2 := strings.Index(line2, "up-to-date")
	idx3 := strings.Index(line3, "up-to-date")

	if idx1 != idx2 || idx2 != idx3 {
		t.Fatalf("column misalignment: idx1=%d, idx2=%d, idx3=%d\nline1: %s\nline2: %s\nline3: %s",
			idx1, idx2, idx3, line1, line2, line3)
	}
}

func TestRenderConciseActiveResultsTo(t *testing.T) {
	states := []*PullRepoState{
		{RepoName: "alim-cv-v8", Changes: "dirty"},
		{RepoName: "ai-empathy-prompt-tuner-v1", Changes: ""},
		{RepoName: "Antigravity-Manager", Changes: "+4941/-484 (68)"},
		{RepoName: "kita-social-media-content-calender-v2", Changes: "synced"},
	}

	var buf bytes.Buffer
	RenderConciseActiveResultsTo(&buf, states)

	rawLines := strings.Split(buf.String(), "\n")
	var bulletLines []string
	for _, l := range rawLines {
		trimmed := strings.TrimSpace(l)
		if strings.HasPrefix(trimmed, "•") {
			bulletLines = append(bulletLines, l)
		}
	}

	// Only 2 active repos (alim-cv-v8 and Antigravity-Manager) should be rendered; 2 up-to-date omitted!
	if len(bulletLines) != 2 {
		t.Fatalf("expected 2 active repo bullet lines rendered, got %d:\n%s", len(bulletLines), buf.String())
	}

	for i, line := range bulletLines {
		if !strings.HasPrefix(line, "    • ") {
			t.Errorf("line %d missing bullet prefix: %q", i, line)
		}
	}

	if !strings.Contains(buf.String(), "Reason:") {
		t.Errorf("expected Reason in output, got: %s", buf.String())
	}
	if !strings.Contains(buf.String(), "Next Step:") {
		t.Errorf("expected Next Step in output, got: %s", buf.String())
	}
}

func TestRenderConciseActiveResultsTo_AllUpToDate(t *testing.T) {
	states := []*PullRepoState{
		{RepoName: "repo-a", Changes: "up-to-date"},
		{RepoName: "repo-b", Changes: ""},
	}

	var buf bytes.Buffer
	RenderConciseActiveResultsTo(&buf, states)

	if !strings.Contains(buf.String(), "(all repositories are up-to-date)") {
		t.Fatalf("expected all repositories up to date notice, got: %s", buf.String())
	}
}

func TestStyleRepoStatusLabel(t *testing.T) {
	if s := StyleRepoStatusLabel("up-to-date"); !strings.Contains(s, "up-to-date") {
		t.Fatalf("expected up-to-date in styled output, got %s", s)
	}
	if s := StyleRepoStatusLabel("dirty"); !strings.Contains(s, "dirty") {
		t.Fatalf("expected dirty in styled output, got %s", s)
	}
	if s := StyleRepoStatusLabel("+576/-119 (32)"); !strings.Contains(s, "+576") || !strings.Contains(s, "/-119") {
		t.Fatalf("expected diff parts in styled output, got %s", s)
	}
}

func TestFormatWrappedInactiveList(t *testing.T) {
	names := []string{
		"repo-alpha", "repo-beta", "repo-gamma", "repo-delta",
		"repo-epsilon", "repo-zeta", "repo-eta", "repo-theta",
	}

	wrapped := FormatWrappedInactiveList(names, "      ", 40)
	lines := strings.Split(wrapped, "\n")
	if len(lines) <= 1 {
		t.Fatalf("expected multiple wrapped lines, got %d:\n%s", len(lines), wrapped)
	}

	for i, l := range lines {
		if !strings.HasPrefix(l, "      ") {
			t.Errorf("line %d missing 6-space indent: %q", i, l)
		}
		// Ensure line does not exceed maxLineLen (unless single token)
		if len(l) > 42 {
			t.Errorf("line %d exceeds expected width (len=%d): %q", i, len(l), l)
		}
	}
}

func TestResolveDualPullRemediationHints_Conflict(t *testing.T) {
	err := "Merge conflict detected during pull: automatic merge failed"
	l1, c1, l2, c2 := ResolveDualPullRemediationHints(err, "", "cat-my-v12")

	if l1 != "Stash & Re-pull" || c1 != "gitmap fix cat-my-v12 stash" {
		t.Fatalf("unexpected option 1: label=%q, cmd=%q", l1, c1)
	}
	if l2 != "Discard & Abort" || c2 != "gitmap fix cat-my-v12 discard" {
		t.Fatalf("unexpected option 2: label=%q, cmd=%q", l2, c2)
	}
}
