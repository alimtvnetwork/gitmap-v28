package cmdpull

import (
	"bytes"
	"strings"
	"testing"
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
	var lines []string
	for _, l := range rawLines {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}

	if len(lines) != 4 {
		t.Fatalf("expected 4 lines rendered, got %d:\n%s", len(lines), buf.String())
	}

	// Verify that each status starts at the exact same column index
	var commonStatusCol int
	for i, line := range lines {
		// Bullet prefix must be present
		if !strings.HasPrefix(line, "    • ") {
			t.Errorf("line %d missing bullet prefix: %q", i, line)
		}

		expectedLabel := ResolveRepoStatusLabel(states[i].Changes)
		statusCol := strings.LastIndex(line, expectedLabel)
		if statusCol == -1 {
			t.Errorf("line %d does not contain expected status %q: %s", i, expectedLabel, line)
		}

		if i == 0 {
			commonStatusCol = statusCol
		} else if statusCol != commonStatusCol {
			t.Errorf("line %d status column (%d) != line 0 status column (%d):\nLine 0: %s\nLine %d: %s",
				i, statusCol, commonStatusCol, lines[0], i, line)
		}
	}
}
