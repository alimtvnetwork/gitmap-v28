//go:build tempe2e

package e2e_test

import (
	"bytes"
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
)

func TestTempE2E_PAEColumnWidthAlignment(t *testing.T) {
	if os.Getenv("RUN_TEMP_E2E") != "1" {
		t.Skip("skipping temporary E2E test: RUN_TEMP_E2E=1 not set")
	}

	states := []*cmdpull.PullRepoState{
		{RepoName: "ai-empathy-prompt-tuner-v1", Changes: ""},
		{RepoName: "alim-cv-v8", Changes: "dirty"},
		{RepoName: "alim-karim-profile-v2", Changes: "synced"},
		{RepoName: "Antigravity-Manager", Changes: "+4941/-484 (68)"},
		{RepoName: "bsrm-presentation-hiltrax-v4", Changes: "up-to-date"},
		{RepoName: "gstack", Changes: "+13217/-2450 (177)"},
		{RepoName: "kita-social-media-content-calender-v2", Changes: "up-to-date"},
	}

	var buf bytes.Buffer
	cmdpull.RenderConciseActiveResultsTo(&buf, states)

	rawLines := strings.Split(buf.String(), "\n")
	var lines []string
	for _, l := range rawLines {
		if strings.TrimSpace(l) != "" {
			lines = append(lines, l)
		}
	}
	if len(lines) != len(states) {
		t.Fatalf("expected %d rendered lines, got %d", len(states), len(lines))
	}

	var expectedStatusIndex int
	for i, line := range lines {
		if !strings.HasPrefix(line, "    • ") {
			t.Fatalf("line %d missing bullet prefix: %s", i, line)
		}

		expectedLabel := cmdpull.ResolveRepoStatusLabel(states[i].Changes)
		statusIdx := strings.LastIndex(line, expectedLabel)
		if statusIdx == -1 {
			t.Fatalf("line %d does not contain expected status %q: %s", i, expectedLabel, line)
		}

		if i == 0 {
			expectedStatusIndex = statusIdx
		} else if statusIdx != expectedStatusIndex {
			t.Fatalf("line %d status column index (%d) does not match expected (%d):\nLine 0: %s\nLine %d: %s",
				i, statusIdx, expectedStatusIndex, lines[0], i, line)
		}
	}
}
