package cmd

import (
	"testing"
)

func TestInstallerExportGitCmd(t *testing.T) {
	if err := executeExportGit([]string{}, true); err == nil {
		t.Fatal("expected error on empty args for export-all-git")
	}

	if err := executeExportGit([]string{"slug"}, false); err == nil {
		t.Fatal("expected error on single arg for export-git")
	}
}
