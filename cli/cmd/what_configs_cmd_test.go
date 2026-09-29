package cmd

import (
	"path/filepath"
	"strings"
	"testing"
)

func TestWhatConfigs_Inspection(t *testing.T) {
	fixtureFiles := []string{
		filepath.Join("..", "jsonenvelope", "fixtures", "commit-pull-config.json"),
		filepath.Join("..", "jsonenvelope", "fixtures", "macro.json"),
		filepath.Join("..", "jsonenvelope", "fixtures", "ssh-nodes.json"),
		filepath.Join("..", "jsonenvelope", "fixtures", "ui-settings.json"),
		filepath.Join("..", "..", "version.json"),
	}

	results := inspectConfigFiles(fixtureFiles)
	if len(results) != 5 {
		t.Fatalf("expected 5 inspection results, got %d", len(results))
	}

	for _, r := range results {
		if r.Category == "" || r.Category == "Unreadable File" {
			t.Errorf("file %s has invalid category: %s", r.FilePath, r.Category)
		}
		if r.ManageCmd == "" || r.ManageCmd == "-" {
			t.Errorf("file %s has invalid manage command: %s", r.FilePath, r.ManageCmd)
		}
	}

	// Verify version.json
	versionResult := results[4]
	if !strings.Contains(versionResult.Category, "Version") {
		t.Errorf("expected version.json to be categorized as Version site, got: %s", versionResult.Category)
	}
}
