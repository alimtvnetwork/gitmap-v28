package cmdos

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestDiscoverAICleanTargetsIsolated(t *testing.T) {
	geminiDir, wsDir, tempDir, gitmapDir := setupAllMockDirectories(t)
	applyTestOverrides(geminiDir, wsDir, tempDir, gitmapDir)
	defer resetTestOverrides()

	assertBrainCategory(t)
	assertTasksCategory(t)
	assertTempCategory(t)
	assertGitmapCategory(t)
	assertAllCategoriesCount(t)
}

func setupAllMockDirectories(t *testing.T) (string, string, string, string) {
	geminiDir := t.TempDir()
	wsDir := t.TempDir()
	tempDir := t.TempDir()
	gitmapDir := t.TempDir()

	setupMockBrainCaches(t, geminiDir)
	setupMockTasks(t, wsDir, geminiDir)
	setupMockTempDumps(t, tempDir)
	setupMockGitmapCaches(t, gitmapDir)

	return geminiDir, wsDir, tempDir, gitmapDir
}

func setupMockBrainCaches(t *testing.T, baseDir string) {
	createTestFile(t, filepath.Join(baseDir, "crashes", "crash1.log"), "crash-data")
	createTestFile(t, filepath.Join(baseDir, "scratch", "scratch1.tmp"), "scratch-data")
	createTestFile(t, filepath.Join(baseDir, "brain", "tempmediaStorage", "media1.bin"), "media-data")
	createTestFile(t, filepath.Join(baseDir, "brain", "conv-001", "scratch", "run.log"), "log-data")
}

func setupMockTasks(t *testing.T, wsDir, geminiDir string) {
	createTestFile(t, filepath.Join(wsDir, ".system_generated", "tasks", "task1.json"), "task-data")
	createTestFile(t, filepath.Join(geminiDir, "brain", "conv-001", ".system_generated", "tasks", "sub.json"), "sub-data")
}

func setupMockTempDumps(t *testing.T, tempDir string) {
	createTestFile(t, filepath.Join(tempDir, "antigravity-dump.log"), "dump1")
	createTestFile(t, filepath.Join(tempDir, "gemini-cache.tmp"), "dump2")
	createTestFile(t, filepath.Join(tempDir, "agent-task.log"), "dump3")
}

func setupMockGitmapCaches(t *testing.T, gitmapDir string) {
	createTestFile(t, filepath.Join(gitmapDir, "downloads", "installer.tar.gz"), "bin-data")
	createTestFile(t, filepath.Join(gitmapDir, "build", "cache.o"), "build-data")
}

func applyTestOverrides(gemini, ws, temp, gitmap string) {
	OverrideAntigravityDir = gemini
	OverrideWorkspaceDir = ws
	OverrideTempDir = temp
	OverrideGitmapDir = gitmap
}

func resetTestOverrides() {
	OverrideAntigravityDir = ""
	OverrideWorkspaceDir = ""
	OverrideTempDir = ""
	OverrideGitmapDir = ""
}

func assertBrainCategory(t *testing.T) {
	cat := DiscoverAntigravityBrainCaches()
	if cat.FileCount != 4 {
		t.Errorf("DiscoverAntigravityBrainCaches FileCount = %d, want 4", cat.FileCount)
	}
}

func assertTasksCategory(t *testing.T) {
	cat := DiscoverSystemGeneratedTasks()
	if cat.FileCount != 2 {
		t.Errorf("DiscoverSystemGeneratedTasks FileCount = %d, want 2", cat.FileCount)
	}
}

func assertTempCategory(t *testing.T) {
	cat := DiscoverOSTempAIDumps()
	if cat.FileCount < 3 {
		t.Errorf("DiscoverOSTempAIDumps FileCount = %d, want at least 3", cat.FileCount)
	}
}

func assertGitmapCategory(t *testing.T) {
	cat := DiscoverGitmapInstallerCaches()
	if cat.FileCount != 2 {
		t.Errorf("DiscoverGitmapInstallerCaches FileCount = %d, want 2", cat.FileCount)
	}
}

func assertAllCategoriesCount(t *testing.T) {
	all := DiscoverAllAICleanTargets()
	if len(all) != 4 {
		t.Errorf("DiscoverAllAICleanTargets count = %d, want 4", len(all))
	}
}

func TestProtectedFilesPreservation(t *testing.T) {
	assertStandardProtectedFlags(t)
	testProtectedFileRemovalFailure(t)
}

func assertStandardProtectedFlags(t *testing.T) {
	assertProtectedFlag(t, "antigravity_state.pbtxt", true)
	assertProtectedFlag(t, "installation_id", true)
	assertProtectedFlag(t, "config.json", true)
	assertProtectedFlag(t, ".gitkeep", true)
	assertProtectedFlag(t, "random_cache.log", false)
}

func testProtectedFileRemovalFailure(t *testing.T) {
	tempDir := t.TempDir()
	protectedPath := filepath.Join(tempDir, "config.json")
	createTestFile(t, protectedPath, `{"key": "value"}`)

	count, bytesFreed := removeSingleFileSafely(protectedPath)
	if count != 0 || bytesFreed != 0 {
		t.Errorf("removeSingleFileSafely removed protected file: count=%d, bytes=%d", count, bytesFreed)
	}
	if !pathExists(protectedPath) {
		t.Error("protected file was deleted from disk")
	}
}

func assertProtectedFlag(t *testing.T, filename string, isExpected bool) {
	isProtected := isProtectedFile(filename)
	if isProtected != isExpected {
		t.Errorf("isProtectedFile(%q) = %v, want %v", filename, isProtected, isExpected)
	}
}

func TestRenderAICleanPreflightTable(t *testing.T) {
	cats := []AICleanCategory{
		{Name: "Antigravity Brain Caches", FileCount: 10, TotalBytes: 1024 * 1024},
		{Name: "System Generated Tasks", FileCount: 5, TotalBytes: 512 * 1024},
	}

	tbl := renderAICleanPreflightTable(cats, 15, 1536*1024)
	assertTableContains(t, tbl, "Category")
	assertTableContains(t, tbl, "Files")
	assertTableContains(t, tbl, "Size")
	assertTableContains(t, tbl, "Antigravity Brain Caches")
	assertTableContains(t, tbl, "Total")
}

func assertTableContains(t *testing.T, tbl, substr string) {
	hasSub := strings.Contains(tbl, substr)
	if !hasSub {
		t.Errorf("table missing expected string %q", substr)
	}
}

func TestRunOSAICleanCLIHelp(t *testing.T) {
	if err := RunOSAICleanCLI([]string{"--help"}); err != nil {
		t.Errorf("RunOSAICleanCLI --help failed: %v", err)
	}

	if err := RunOSAICleanCLI([]string{"-h"}); err != nil {
		t.Errorf("RunOSAICleanCLI -h failed: %v", err)
	}

	if err := RunOSAICleanCLI([]string{"help"}); err != nil {
		t.Errorf("RunOSAICleanCLI help failed: %v", err)
	}
}

func TestRunOSAICleanCLIDryRunAndJSON(t *testing.T) {
	emptyDir := t.TempDir()
	applyTestOverrides(emptyDir, emptyDir, emptyDir, emptyDir)
	defer resetTestOverrides()

	if err := RunOSAICleanCLI([]string{"--dry-run"}); err != nil {
		t.Errorf("RunOSAICleanCLI --dry-run failed: %v", err)
	}

	if err := RunOSAICleanCLI([]string{"--json", "--dry-run"}); err != nil {
		t.Errorf("RunOSAICleanCLI --json --dry-run failed: %v", err)
	}
}

func TestRunOSAICleanCLIPurgeExecution(t *testing.T) {
	tempDir := t.TempDir()
	cacheFile := filepath.Join(tempDir, "purge.tmp")
	createTestFile(t, cacheFile, "purge-me")

	count, bytes := purgeCategoryFiles([]string{cacheFile})
	if count != 1 || bytes == 0 {
		t.Errorf("purgeCategoryFiles count=%d bytes=%d, want count=1 bytes>0", count, bytes)
	}

	if pathExists(cacheFile) {
		t.Error("purged file still exists on disk")
	}
}

func TestRunOSAICleanCLIConfirmationAbort(t *testing.T) {
	emptyDir := t.TempDir()
	applyTestOverrides(emptyDir, emptyDir, emptyDir, emptyDir)
	defer resetTestOverrides()

	origPrompt := promptAICleanConfirmFn
	defer func() { promptAICleanConfirmFn = origPrompt }()

	promptAICleanConfirmFn = mockRejectPrompt

	if err := RunOSAICleanCLI([]string{}); err != nil {
		t.Errorf("RunOSAICleanCLI aborted confirmation failed: %v", err)
	}
}

func mockRejectPrompt(msg string) (bool, error) {
	return false, nil
}

func createTestFile(t *testing.T, fullPath, content string) {
	dir := filepath.Dir(fullPath)
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("MkdirAll failed: %v", err)
	}

	if err := os.WriteFile(fullPath, []byte(content), 0o644); err != nil {
		t.Fatalf("WriteFile failed: %v", err)
	}
}
