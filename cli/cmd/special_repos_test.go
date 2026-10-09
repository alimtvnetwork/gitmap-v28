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

	folder1Repeat, err := db.ResolveOrCreateRepoFolder("rs", repoSecretsRoot, "gitmap")
	if err != nil || folder1Repeat != folder1 {
		t.Fatalf("expected idempotent folder resolution %s, got %s, err: %v", folder1, folder1Repeat, err)
	}
}

func TestStoreSpecialRepoStandaloneHelpers(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-helpers.db")
	db, err := store.OpenSpecialReposSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenSpecialReposSplitDBAt failed: %v", err)
	}
	defer db.Close()

	verifyStandaloneRepoQueries(t, db)
	verifyStandaloneSeqHelpers(t, db, tempDir)
}

func verifyStandaloneRepoQueries(t *testing.T, db *store.SpecialReposSplitDB) {
	rec, err := store.GetSpecialRepoByShortKey(db.Conn(), "rs")
	if err != nil || rec.RepoKey != "repo-secrets" {
		t.Fatalf("GetSpecialRepoByShortKey failed: %+v, err: %v", rec, err)
	}
	recKey, err := store.GetSpecialRepoByKey(db.Conn(), "repo-cache")
	if err != nil || recKey.ShortKey != "rc" {
		t.Fatalf("GetSpecialRepoByKey failed: %+v, err: %v", recKey, err)
	}
	all, err := store.GetAllSpecialRepos(db.Conn())
	if err != nil || len(all) < 2 {
		t.Fatalf("GetAllSpecialRepos failed: %d items, err: %v", len(all), err)
	}
}

func verifyStandaloneSeqHelpers(t *testing.T, db *store.SpecialReposSplitDB, tempDir string) {
	prefix1, err := store.GetNextRepoFolderSeq(db.Conn(), "rs", "app-one")
	if err != nil || prefix1 != "01-app-one" {
		t.Fatalf("expected 01-app-one, got %s, err: %v", prefix1, err)
	}
	prefix2, err := store.GetNextRepoFolderSeq(db.Conn(), "rs", "app-two")
	if err != nil || prefix2 != "02-app-two" {
		t.Fatalf("expected 02-app-two, got %s, err: %v", prefix2, err)
	}
	fileSeq, err := store.GetNextFileSeq(tempDir, "")
	if err != nil || fileSeq != 1 {
		t.Fatalf("expected fileSeq 1, got %d, err: %v", fileSeq, err)
	}
}

func TestOneTimePromptPersistence(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-prompt.db")
	db, err := store.OpenSpecialReposSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenSpecialReposSplitDBAt failed: %v", err)
	}
	defer db.Close()

	if err := store.MarkSpecialRepoPromptAnswered(db.Conn(), "rs", "approved"); err != nil {
		t.Fatalf("MarkSpecialRepoPromptAnswered failed: %v", err)
	}
	rec, err := store.GetSpecialRepoByShortKey(db.Conn(), "rs")
	if err != nil || !rec.HasAnsweredPrompt() {
		t.Fatalf("expected prompt to be answered, got %+v, err: %v", rec, err)
	}

	records, err := CheckSpecialReposOnScanWithDB(db, tempDir, true)
	if err != nil || len(records) < 2 {
		t.Fatalf("CheckSpecialReposOnScanWithDB failed: %v", err)
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

func TestDeriveFileAndFolderDestinationName(t *testing.T) {
	name1 := deriveFileDestinationName("/path/to/my-secret.env", "token")
	if name1 != "token.env" {
		t.Fatalf("expected token.env, got %s", name1)
	}
	name2 := deriveFileDestinationName("/path/to/my-secret.env", "")
	if name2 != "my-secret.env" {
		t.Fatalf("expected my-secret.env, got %s", name2)
	}
	dir1 := deriveFolderDestinationName("/path/to/my_certs", "certs")
	if dir1 != "certs" {
		t.Fatalf("expected certs, got %s", dir1)
	}
}

func TestResolveSpecialRepoCDPathWithDB(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-cd.db")
	db, err := store.OpenSpecialReposSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenSpecialReposSplitDBAt failed: %v", err)
	}
	defer db.Close()

	pathRS, err := resolveSpecialRepoCDPathWithDB(db, "rs", tempDir)
	if err != nil || pathRS != filepath.Join(tempDir, "repo-secrets") {
		t.Fatalf("expected fallback %s, got %s, err: %v", filepath.Join(tempDir, "repo-secrets"), pathRS, err)
	}

	customDir := filepath.Join(tempDir, "custom-cache")
	_ = os.MkdirAll(customDir, 0755)
	_ = store.UpdateSpecialRepoPath(db.Conn(), "rc", customDir, "")

	pathRC, err := resolveSpecialRepoCDPathWithDB(db, "rc", tempDir)
	if err != nil || pathRC != customDir {
		t.Fatalf("expected customDir %s, got %s, err: %v", customDir, pathRC, err)
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

func TestSpecialRepoOptionFlags(t *testing.T) {
	opts := parseSpecialCLIOptions([]string{"file", "path/secret.env", "--repo", "gitmap", "--name", "token", "--msg", "init", "--no-push", "--json"})
	if opts.subcmd != "file" || opts.primaryArg != "path/secret.env" {
		t.Fatalf("unexpected subcmd or primaryArg: %+v", opts)
	}
	if opts.repoName != "gitmap" || opts.slug != "token" || opts.commitMsg != "init" {
		t.Fatalf("unexpected parsed flags: %+v", opts)
	}
	if !opts.isNoPush || !opts.isJSON {
		t.Fatalf("expected isNoPush and isJSON to be true")
	}
}

func TestExecuteSpecialRepoTextAndFileSequencing(t *testing.T) {
	tempDir := t.TempDir()
	dbPath := filepath.Join(tempDir, "test-exec.db")
	db, err := store.OpenSpecialReposSplitDBAt(dbPath)
	if err != nil {
		t.Fatalf("OpenSpecialReposSplitDBAt failed: %v", err)
	}
	defer db.Close()

	verifySpecialRepoTextExecution(t, db, tempDir)
	verifySpecialRepoFileExecution(t, db, tempDir)
}

func verifySpecialRepoTextExecution(t *testing.T, db *store.SpecialReposSplitDB, tempDir string) {
	specialRoot := filepath.Join(tempDir, "repo-secrets")
	textOpts := specialCLIOptions{primaryArg: "API_KEY=12345", slug: "api-key", ext: ".env", repoName: "demo", isNoPush: true}
	resText, err := executeSpecialRepoText(db, "rs", specialRoot, textOpts)
	if err != nil {
		t.Fatalf("executeSpecialRepoText failed: %v", err)
	}
	if filepath.Base(resText.TargetPath) != "01-api-key.env" {
		t.Fatalf("expected 01-api-key.env, got %s", filepath.Base(resText.TargetPath))
	}
}

func verifySpecialRepoFileExecution(t *testing.T, db *store.SpecialReposSplitDB, tempDir string) {
	specialRoot := filepath.Join(tempDir, "repo-secrets")
	srcFile := filepath.Join(tempDir, "source-token.txt")
	_ = os.WriteFile(srcFile, []byte("token-content"), 0644)
	fileOpts := specialCLIOptions{primaryArg: srcFile, slug: "token", repoName: "demo", isNoPush: true}
	resFile, err := executeSpecialRepoFile(db, "rs", specialRoot, fileOpts)
	if err != nil {
		t.Fatalf("executeSpecialRepoFile failed: %v", err)
	}
	if filepath.Base(resFile.TargetPath) != "02-token.txt" {
		t.Fatalf("expected 02-token.txt, got %s", filepath.Base(resFile.TargetPath))
	}
}
