package cmdagy

import (
	"testing"
	"time"
)

func TestTrimBoilerplatePrompt(t *testing.T) {
	raw := `<system_prompt>You are Antigravity coding assistant</system_prompt>
<context>Some context here</context>
Please implement the new feature for gitmap CLI.
Ensure all tests pass and coding guidelines are followed.`

	trimmed := TrimBoilerplatePrompt(raw)
	if trimmed == "" {
		t.Errorf("expected non-empty trimmed prompt")
	}
	if len(trimmed) > 205 {
		t.Errorf("expected trimmed length <= 203, got %d", len(trimmed))
	}
	if trimmed != "Please implement the new feature for gitmap CLI. Ensure all tests pass and coding guidelines are followed." {
		t.Errorf("unexpected trimmed text: %q", trimmed)
	}
}

func TestSaveAndLoadSequenceCache(t *testing.T) {
	entries := []CachedSequenceEntry{
		{Seq: 1, ID: "conv-1", Name: "First Task", Messages: 10, Status: "RUNNING", Queued: 2},
		{Seq: 2, ID: "conv-2", Name: "Second Task", Messages: 5, Status: "IDLE", Queued: 0},
	}

	err := SaveSequenceCache(entries)
	if err != nil {
		t.Fatalf("unexpected error saving cache: %v", err)
	}

	manifest, isValid := LoadSequenceCache()
	if !isValid || manifest == nil {
		t.Fatalf("expected valid cache manifest")
	}
	if len(manifest.Entries) != 2 {
		t.Errorf("expected 2 entries, got %d", len(manifest.Entries))
	}

	resolved, found := ResolveCachedSequence(1)
	if !found || resolved == nil {
		t.Fatalf("expected to resolve sequence 1")
	}
	if resolved.ID != "conv-1" {
		t.Errorf("expected ID conv-1, got %s", resolved.ID)
	}
	if resolved.Status != "RUNNING" {
		t.Errorf("expected status RUNNING, got %s", resolved.Status)
	}
}

func TestSortInspectRowsRunningFirst(t *testing.T) {
	now := time.Now().UTC().Format(time.RFC3339)
	earlier := time.Now().Add(-10 * time.Minute).UTC().Format(time.RFC3339)

	rows := []AgyConvInspectRow{
		{ID: "idle-new", Status: "IDLE", IsRunning: false, LastMod: now},
		{ID: "run-old", Status: "RUNNING", IsRunning: true, LastMod: earlier},
		{ID: "idle-old", Status: "IDLE", IsRunning: false, LastMod: earlier},
	}

	sortInspectRows(rows)
	assignInspectRowSequences(rows)

	if rows[0].ID != "run-old" {
		t.Errorf("expected RUNNING row first, got %s", rows[0].ID)
	}
	if rows[0].Seq != 1 {
		t.Errorf("expected Seq 1, got %d", rows[0].Seq)
	}
}
