package cmdagy

import (
	"strings"
	"testing"
	"time"
)

func TestAgyActiveAndQueueFormatting(t *testing.T) {
	secs := 125
	formatted := formatElapsedDuration(secs)
	if formatted != "2m 5s" {
		t.Errorf("expected 2m 5s, got %s", formatted)
	}

	badge := formatQueueStatusBadge("queued")
	if !strings.Contains(badge, "queued") {
		t.Errorf("expected badge to contain queued, got %s", badge)
	}

	snippet := sanitizePromptSnippet("First line of a very long prompt that goes on and on and on and exceeds the limit\nSecond line", 30)
	if len(snippet) > 30 {
		t.Errorf("expected snippet len <= 30, got %d (%s)", len(snippet), snippet)
	}

	summaries := []AgyWorkspaceQueueSummary{
		{TotalQueued: 2},
		{TotalQueued: 3},
	}
	tot := countTotalQueuedPrompts(summaries)
	if tot != 5 {
		t.Errorf("expected 5 queued, got %d", tot)
	}

	nowStr := time.Now().UTC().Format(time.RFC3339)
	elapsed := calculateElapsedSeconds(nowStr)
	if elapsed < 0 || elapsed > 5 {
		t.Errorf("expected elapsed between 0 and 5s, got %d", elapsed)
	}
}
