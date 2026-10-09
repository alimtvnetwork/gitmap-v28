package cmdimport

import (
	"path/filepath"
	"testing"
)

func TestImportAllJSON_DryRun(t *testing.T) {
	fixtureFiles := []string{
		filepath.Join("..", "jsonenvelope", "fixtures", "commit-pull-config.json"),
		filepath.Join("..", "jsonenvelope", "fixtures", "macro.json"),
		filepath.Join("..", "jsonenvelope", "fixtures", "ssh-nodes.json"),
		filepath.Join("..", "jsonenvelope", "fixtures", "ui-settings.json"),
		filepath.Join("..", "jsonenvelope", "fixtures", "unmatched.json"),
	}

	results := executeImportAllJSONBatch(fixtureFiles, true)
	if len(results) != 5 {
		t.Fatalf("expected 5 results, got %d", len(results))
	}

	succeeded, skipped, failed := tallyImportResults(results)
	if succeeded != 4 {
		t.Errorf("expected 4 succeeded (previewed) items, got %d", succeeded)
	}
	if skipped != 1 {
		t.Errorf("expected 1 skipped item (unmatched.json), got %d", skipped)
	}
	if failed != 0 {
		t.Errorf("expected 0 failed items, got %d", failed)
	}

	// Verify unmatched result
	unmatched := results[4]
	if unmatched.Status != "SKIPPED" {
		t.Errorf("expected unmatched.json to have status SKIPPED, got %s", unmatched.Status)
	}
}
