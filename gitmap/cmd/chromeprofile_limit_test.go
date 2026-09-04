package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestParseChromeTransferOptions(t *testing.T) {
	args := []string{"export", "--limit=1", "--profile=Work", "out.json", "--format=json"}
	opts := parseChromeTransferOptions(args)

	if opts.Limit != 1 {
		t.Errorf("expected limit 1, got %d", opts.Limit)
	}
	if opts.Profile != "Work" {
		t.Errorf("expected profile Work, got %s", opts.Profile)
	}
	if opts.Format != "json" {
		t.Errorf("expected format json, got %s", opts.Format)
	}
	if len(opts.Positional) != 2 || opts.Positional[0] != "export" || opts.Positional[1] != "out.json" {
		t.Errorf("unexpected positional: %+v", opts.Positional)
	}

	args2 := []string{"--limit", "2", "-p", "Default", "-n", "3"}
	opts2 := parseChromeTransferOptions(args2)
	if opts2.Limit != 3 {
		t.Errorf("expected last limit 3, got %d", opts2.Limit)
	}
	if opts2.Profile != "Default" {
		t.Errorf("expected profile Default, got %s", opts2.Profile)
	}
}

func TestChromeExportAndImportWithLimitAndProfile(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_CHROME_USER_DATA", tempDir)
	seedChromeProfileTree(t, tempDir)

	workDir := t.TempDir()

	// 1. Export with --limit 1 to JSON
	outJSON := filepath.Join(workDir, "limited.json")
	if err := runChromeProfileExport([]string{outJSON, "--format=json", "--limit=1"}); err != nil {
		t.Fatalf("export with --limit=1 failed: %v", err)
	}
	if _, err := os.Stat(outJSON); err != nil {
		t.Fatalf("expected limited.json to exist: %v", err)
	}

	// 2. Export single profile with --profile
	outSingle := filepath.Join(workDir, "single_flag.json")
	if err := runChromeProfileExport([]string{outSingle, "--profile=Default"}); err != nil {
		t.Fatalf("export with --profile=Default failed: %v", err)
	}
	if _, err := os.Stat(outSingle); err != nil {
		t.Fatalf("expected single_flag.json to exist: %v", err)
	}

	// 3. Import with --limit 1
	if err := runChromeProfileImport([]string{outJSON, "--limit=1"}); err != nil {
		t.Fatalf("import with --limit=1 failed: %v", err)
	}

	// 4. Import single profile with --profile
	if err := runChromeProfileImport([]string{outJSON, "--profile=Default"}); err != nil {
		t.Fatalf("import with --profile=Default failed: %v", err)
	}
}

func TestAgyClearProtectsPinnedProjects(t *testing.T) {
	tmp := t.TempDir()
	t.Setenv("USERPROFILE", tmp)
	t.Setenv("HOME", tmp)

	pinnedPath := filepath.Join(tmp, "pinned-repo")
	_ = os.MkdirAll(pinnedPath, 0755)

	pinned, err := addPinnedProjectTarget(pinnedPath)
	if err != nil {
		t.Fatalf("failed to pin project: %v", err)
	}

	// Now pretend the directory is missing (deleted)
	_ = os.RemoveAll(pinnedPath)

	projects := []AgyProject{
		{
			ID:   pinned.ID,
			Name: "pinned-repo",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{
						GitFolder: &AgyGitFolder{
							FolderURI: "file:///" + filepath.ToSlash(pinnedPath),
						},
					},
				},
			},
		},
		{
			ID:   "unpinned-missing",
			Name: "unpinned-missing",
			ProjectResources: &AgyProjectResources{
				Resources: []AgyResource{
					{
						GitFolder: &AgyGitFolder{
							FolderURI: "file:///" + filepath.ToSlash(filepath.Join(tmp, "unpinned-missing")),
						},
					},
				},
			},
		},
	}

	targets := selectClearTargets(projects)
	for _, target := range targets {
		if target.ID == pinned.ID {
			t.Errorf("agy clear targeted pinned project %s for removal!", pinned.ID)
		}
	}
	if len(targets) != 1 || targets[0].ID != "unpinned-missing" {
		t.Errorf("expected only unpinned-missing in targets, got: %+v", targets)
	}
}
