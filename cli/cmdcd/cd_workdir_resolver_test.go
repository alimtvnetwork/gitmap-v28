package cmdcd

import (
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestIsWorkDirKeyword(t *testing.T) {
	expectedKeywords := []string{
		"work",
		"$work",
		`\$work`,
		"def",
		"$def",
		`\$def`,
		"workdir",
		"$workdir",
		`\$workdir`,
		"default",
		"$default",
		`\$default`,
		"wd",
		"$wd",
		`\$wd`,
		// Uppercase/mixed-case variants
		"WORK",
		"$WORK",
		"DEF",
		"$DEF",
	}

	for _, kw := range expectedKeywords {
		if !isWorkDirKeyword(kw) {
			t.Errorf("expected isWorkDirKeyword(%q) to be true", kw)
		}
	}

	nonKeywords := []string{
		"unknown",
		"my-repo",
		"$foo",
		"workspace",
		"",
	}

	for _, nkw := range nonKeywords {
		if isWorkDirKeyword(nkw) {
			t.Errorf("expected isWorkDirKeyword(%q) to be false", nkw)
		}
	}
}

func TestResolveCDWorkDirPath_DefaultWorkDirVariants(t *testing.T) {
	tempDir := t.TempDir()

	db, errDB := store.OpenDefault()
	if errDB != nil {
		t.Skip("sqlite db unavailable in test environment")
	}
	defer db.Close()

	if errEnsureTable := db.EnsureWorkDirsTable(); errEnsureTable != nil {
		t.Fatalf("failed to ensure workdirs table: %v", errEnsureTable)
	}

	previousDefault, errGetPrev := db.GetDefaultWorkDir()

	t.Cleanup(func() {
		cleanupDB, errCleanup := store.OpenDefault()
		if errCleanup != nil {
			return
		}
		defer cleanupDB.Close()

		if errGetPrev == nil && previousDefault != nil {
			_ = cleanupDB.SetDefaultWorkDir(previousDefault.AbsolutePath)
		}
	})

	createdWd, errEnsure := db.EnsureWorkDir(tempDir, "test-cd-workdir", true)
	if errEnsure != nil {
		t.Fatalf("failed to register default workdir: %v", errEnsure)
	}

	if errSet := db.SetDefaultWorkDir(createdWd.AbsolutePath); errSet != nil {
		t.Fatalf("failed to set default workdir: %v", errSet)
	}

	// Verify tempDir is indeed an existing directory
	info, errStat := os.Stat(tempDir)
	if errStat != nil || !info.IsDir() {
		t.Fatalf("tempDir %s is not a valid directory", tempDir)
	}

	variants := []string{
		"$work",
		`\$work`,
		"$def",
		`\$def`,
		"def",
		"$workdir",
		"$default",
		"$wd",
		"work",
		"workdir",
		"default",
		"wd",
	}

	for _, variant := range variants {
		resolvedPath, hasWorkDir := resolveCDWorkDirPath(variant)
		if !hasWorkDir {
			t.Errorf("variant %q failed to resolve to default workdir", variant)
			continue
		}
		if resolvedPath != tempDir {
			t.Errorf("variant %q resolved to %q, want %q", variant, resolvedPath, tempDir)
		}
	}
}

func TestResolveDefaultWorkDirPath_FallbackEnv(t *testing.T) {
	tempDir := t.TempDir()
	t.Setenv("GITMAP_WORK_DIR", tempDir)

	resolved, ok := resolveEnvWorkDir()
	if !ok || resolved != tempDir {
		t.Fatalf("resolveEnvWorkDir() = %q, %v; want %q, true", resolved, ok, tempDir)
	}
}
