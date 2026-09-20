package store

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeSlugBasic(t *testing.T) {
	cases := []struct {
		input    string
		expected string
	}{
		{"my-repo", "my-repo"},
		{"My_Awesome_Repo!", "my_awesome_repo"},
		{"github.com/org/my-project.git", "org-my-project"},
		{"", "default"},
	}
	for _, tc := range cases {
		actual := SanitizeSlug(tc.input)
		if actual != tc.expected {
			t.Errorf("SanitizeSlug(%q) = %q; want %q", tc.input, actual, tc.expected)
		}
	}
}

func TestResolveSplitDbDirWithGitmap(t *testing.T) {
	tempDir := t.TempDir()
	gitmapDir := filepath.Join(tempDir, ".gitmap")
	if err := os.MkdirAll(gitmapDir, 0755); err != nil {
		t.Fatal(err)
	}

	dir := ResolveSplitDbDir(SectionAutomation, "my-repo", tempDir)
	expected := filepath.ToSlash(filepath.Join(gitmapDir, "data", "automation", "my-repo"))
	if dir != expected {
		t.Errorf("ResolveSplitDbDir = %q; want %q", dir, expected)
	}
}

func TestResolveSplitDbDirFallback(t *testing.T) {
	tempDir := t.TempDir()
	dir := ResolveSplitDbDir(SectionSites, "default", tempDir)
	expected := filepath.ToSlash(filepath.Join(BinaryDataDir(), "sites", "default"))
	if dir != expected {
		t.Errorf("ResolveSplitDbDir fallback = %q; want %q", dir, expected)
	}
}

func TestResolveSplitDbPathCanonical(t *testing.T) {
	tempDir := t.TempDir()
	gitmapDir := filepath.Join(tempDir, ".gitmap")
	_ = os.MkdirAll(gitmapDir, 0755)

	path := ResolveSplitDbPath(SectionPipeline, "test-pipe", tempDir)
	expectedSuffix := ".gitmap/data/pipeline/test-pipe/sql.db"
	if !strings.HasSuffix(path, expectedSuffix) {
		t.Errorf("ResolveSplitDbPath = %q; want suffix %q", path, expectedSuffix)
	}
}

func TestLegacyMigrationAutomation(t *testing.T) {
	tempDir := t.TempDir()
	gitmapData := filepath.Join(tempDir, ".gitmap", "data")
	legacyDir := filepath.Join(gitmapData, "repo1", "automation")
	_ = os.MkdirAll(legacyDir, 0755)
	legacyPath := filepath.Join(legacyDir, "sql.db")
	_ = os.WriteFile(legacyPath, []byte("legacy-data"), 0644)

	targetPath := ResolveSplitDbPath(SectionAutomation, "repo1", tempDir)
	assertMigratedTarget(t, targetPath, legacyPath)
}

func TestLegacyMigrationPipeline(t *testing.T) {
	tempDir := t.TempDir()
	pipeDir := filepath.Join(tempDir, ".gitmap", "pipeline", "repo2")
	_ = os.MkdirAll(pipeDir, 0755)
	legacyPath := filepath.Join(pipeDir, "pipeline.db")
	_ = os.WriteFile(legacyPath, []byte("pipe-data"), 0644)

	targetPath := ResolveSplitDbPath(SectionPipeline, "repo2", tempDir)
	assertMigratedTarget(t, targetPath, legacyPath)
}

func TestLegacyMigrationSchedule(t *testing.T) {
	tempDir := t.TempDir()
	gitmapData := filepath.Join(tempDir, ".gitmap", "data")
	schedDir := filepath.Join(gitmapData, "schedule")
	_ = os.MkdirAll(schedDir, 0755)
	legacyPath := filepath.Join(schedDir, "daily-sync.db")
	_ = os.WriteFile(legacyPath, []byte("sched-data"), 0644)

	targetPath := ResolveSplitDbPath(SectionSchedule, "daily-sync", tempDir)
	assertMigratedTarget(t, targetPath, legacyPath)
}

func assertMigratedTarget(t *testing.T, targetPath, legacyPath string) {
	data, err := os.ReadFile(targetPath)
	if err != nil {
		t.Fatalf("target file %s does not exist: %v", targetPath, err)
	}
	if len(data) == 0 {
		t.Errorf("target file %s is empty", targetPath)
	}
	if isFileExisting(legacyPath) {
		t.Errorf("legacy file %s still exists", legacyPath)
	}
}
