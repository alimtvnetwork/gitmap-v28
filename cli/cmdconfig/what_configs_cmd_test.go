package cmdconfig

import (
	"path/filepath"
	"strings"
	"testing"

	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
)

func TestWhatConfigs_Inspection(t *testing.T) {
	fixtureFiles := []string{
		filepath.Join("..", "jsonx", "fixtures", "commit-pull-config.json"),
		filepath.Join("..", "jsonx", "fixtures", "macro.json"),
		filepath.Join("..", "jsonx", "fixtures", "ssh-nodes.json"),
		filepath.Join("..", "jsonx", "fixtures", "ui-settings.json"),
		filepath.Join("..", "..", cfgconst.FileVersionJSON),
	}

	results := inspectConfigFiles(fixtureFiles)
	if len(results) != 5 {
		t.Fatalf("expected 5 inspection results, got %d", len(results))
	}

	for _, r := range results {
		if r.Category == "" || r.Category == cfgconst.CategoryUnreadable {
			t.Errorf("file %s has invalid category: %s", r.FilePath, r.Category)
		}

		if r.ManageCmd == "" || r.ManageCmd == cfgconst.DashPlaceholder {
			t.Errorf("file %s has invalid manage command: %s", r.FilePath, r.ManageCmd)
		}
	}

	// Verify version.json
	versionResult := results[4]
	if !strings.Contains(versionResult.Category, "Version") {
		t.Errorf("expected version.json to be categorized as Version site, got: %s", versionResult.Category)
	}
}
