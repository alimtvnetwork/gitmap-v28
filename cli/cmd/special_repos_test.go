package cmd

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestSpecialReposNormalizeAndSchema(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-special-repos.db")
	db, err := store.OpenSpecialReposSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenSpecialReposSplitDBAt failed: %v", err)
	}
	defer db.Close()

	verifyDefaultSpecialRepos(t, db)
	verifySpecialRepoFolderAllocation(t, db, tempDir)
}

func verifyDefaultSpecialRepos(t *testing.T, db *store.SpecialReposSplitDB) {
	repos, err := db.ListSpecialRepos()
	if err != nil {
		t.Fatalf("ListSpecialRepos failed: %v", err)
	}
	if len(repos) < 2 {
		t.Fatalf("expected at least 2 default special repos, got %d", len(repos))
	}
	secretRec, err := db.GetSpecialRepo("rs")
	if err != nil || secretRec.RepoKey != "repo-secrets" {
		t.Fatalf("unexpected secret rec: %+v, err: %v", secretRec, err)
	}
	cacheRec, err := db.GetSpecialRepo("rc")
	if err != nil || cacheRec.RepoKey != "repo-cache" {
		t.Fatalf("unexpected cache rec: %+v, err: %v", cacheRec, err)
	}
}

func verifySpecialRepoFolderAllocation(t *testing.T, db *store.SpecialReposSplitDB, tempDir string) {
	repoSecretsRoot := filepath.Join(tempDir, "repo-secrets")
	folder1, err := db.ResolveOrCreateRepoFolder("rs", repoSecretsRoot, "gitmap")
	if err != nil {
		t.Fatalf("ResolveOrCreateRepoFolder failed: %v", err)
	}
	if filepath.Base(folder1) != "01-gitmap" {
		t.Fatalf("expected folder 01-gitmap, got %s", filepath.Base(folder1))
	}

	folder2, err := db.ResolveOrCreateRepoFolder("rs", repoSecretsRoot, "coding-guidelines")
	if err != nil {
		t.Fatalf("ResolveOrCreateRepoFolder 2 failed: %v", err)
	}
	if filepath.Base(folder2) != "02-coding-guidelines" {
		t.Fatalf("expected folder 02-coding-guidelines, got %s", filepath.Base(folder2))
	}

	// Idempotent resolution
	folder1Repeat, err := db.ResolveOrCreateRepoFolder("rs", repoSecretsRoot, "gitmap")
	if err != nil || folder1Repeat != folder1 {
		t.Fatalf("expected idempotent folder resolution %s, got %s, err: %v", folder1, folder1Repeat, err)
	}
}

func TestDeriveTextSlugAndFilename(t *testing.T) {
	res1 := deriveTextSlugAndFilename("hello world secret token", "", ".env")
	if res1 != "hello-world-secret-token.env" {
		t.Fatalf("expected hello-world-secret-token.env, got %s", res1)
	}

	res2 := deriveTextSlugAndFilename("Get-Process | Select-Object -First 5", "procs", ".ps1")
	if res2 != "procs.ps1" {
		t.Fatalf("expected procs.ps1, got %s", res2)
	}

	res3 := deriveTextSlugAndFilename("", "", "")
	if res3 != "note.txt" {
		t.Fatalf("expected note.txt, got %s", res3)
	}
}

func TestIsSpecialRepoCDAlias(t *testing.T) {
	if !isSpecialRepoCDAlias("rs") {
		t.Fatalf("expected rs to be special repo CD alias")
	}
	if !isSpecialRepoCDAlias("repo-secrets") {
		t.Fatalf("expected repo-secrets to be special repo CD alias")
	}
	if !isSpecialRepoCDAlias("rc") {
		t.Fatalf("expected rc to be special repo CD alias")
	}
	if !isSpecialRepoCDAlias("repo-cache") {
		t.Fatalf("expected repo-cache to be special repo CD alias")
	}
	if isSpecialRepoCDAlias("other-repo") {
		t.Fatalf("expected other-repo not to be special repo CD alias")
	}
}

func TestAllocateSequencedItemPath(t *testing.T) {
	tempDir := t.TempDir()
	path1 := allocateSequencedItemPath(tempDir, "token.txt")
	if filepath.Base(path1) != "01-token.txt" {
		t.Fatalf("expected 01-token.txt, got %s", filepath.Base(path1))
	}
	_ = os.WriteFile(path1, []byte("test"), 0644)

	path2 := allocateSequencedItemPath(tempDir, "key.txt")
	if filepath.Base(path2) != "02-key.txt" {
		t.Fatalf("expected 02-key.txt, got %s", filepath.Base(path2))
	}
}

func TestSpecialRepoFlagShorthands(t *testing.T) {
	if !isRepoFlag("--repo") || !isRepoFlag("-r") {
		t.Fatalf("expected --repo and -r to be valid repo flags")
	}
	if !isSlugFlag("--slug") || !isSlugFlag("-s") || !isSlugFlag("--name") || !isSlugFlag("-n") {
		t.Fatalf("expected slug and name shorthands to be valid")
	}
	if !isExtFlag("--ext") || !isExtFlag("-e") {
		t.Fatalf("expected --ext and -e to be valid ext flags")
	}
}

func TestResolveAutoPutActionFunc(t *testing.T) {
	tempDir := t.TempDir()
	testFile := filepath.Join(tempDir, "sample.txt")
	_ = os.WriteFile(testFile, []byte("sample"), 0644)

	fnFile := resolveAutoPutActionFunc(testFile)
	if fnFile == nil {
		t.Fatalf("expected valid file put func")
	}

	fnDir := resolveAutoPutActionFunc(tempDir)
	if fnDir == nil {
		t.Fatalf("expected valid dir put func")
	}

	fnText := resolveAutoPutActionFunc("raw inline text note")
	if fnText == nil {
		t.Fatalf("expected valid text put func")
	}
}

