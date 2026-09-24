package cmdvscode

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/vscodepm"
)

func TestParseRepairOptions(t *testing.T) {
	opts := parseRepairOptions([]string{"--force", "--kill", "--dry-run"})
	if !opts.Force {
		t.Errorf("expected Force to be true")
	}
	if !opts.Kill {
		t.Errorf("expected Kill to be true")
	}
	if !opts.DryRun {
		t.Errorf("expected DryRun to be true")
	}

	shortOpts := parseRepairOptions([]string{"-f", "-k", "-d"})
	if !shortOpts.Force || !shortOpts.Kill || !shortOpts.DryRun {
		t.Errorf("expected short flags to be parsed correctly")
	}
}

func TestSanitizeEntries(t *testing.T) {
	tmpDir := t.TempDir()
	existingPath := filepath.Join(tmpDir, "repo-a")
	if err := os.Mkdir(existingPath, 0755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	entries := []vscodepm.Entry{
		{
			Name:     "Existing",
			RootPath: existingPath,
			Paths:    nil,
			Tags:     nil,
		},
		{
			Name:     "Missing",
			RootPath: filepath.Join(tmpDir, "nonexistent-repo"),
			Paths:    nil,
			Tags:     nil,
		},
	}

	missing := sanitizeEntries(entries)
	if missing != 1 {
		t.Errorf("expected 1 missing path, got %d", missing)
	}
	if entries[0].Paths == nil || len(entries[0].Tags) == 0 {
		t.Errorf("expected Paths and Tags to be initialized")
	}
}

func TestCommitDirRegex(t *testing.T) {
	valid := []string{"2242ebbb54", "7debcd0e2a", "abcdef1234"}
	for _, v := range valid {
		if !commitDirRegex.MatchString(v) {
			t.Errorf("expected %s to match commitDirRegex", v)
		}
	}

	invalid := []string{"2242ebbb5", "2242ebbb541", "not-a-hash!", "2242EBBB54"}
	for _, inv := range invalid {
		if commitDirRegex.MatchString(inv) {
			t.Errorf("expected %s not to match commitDirRegex", inv)
		}
	}
}
