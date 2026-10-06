package cmdnodes

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"io"
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

func TestParseDeployRepoOptions(t *testing.T) {
	args := []string{
		"awesome-repo",
		"--target=node-alpha",
		"--dest=D:\\work\\awesome",
		"--with-ides",
		"--with-pinned",
		"--with-conversations",
		"--clean",
		"--clone",
		"--dry-run",
		"--json",
		"--include=node-alpha,node-beta",
		"--except=node-gamma",
		"--include-main",
		"--open-only",
	}

	opts, err := parseDeployRepoOptions(args)
	if err != nil {
		t.Fatalf("unexpected parse error: %v", err)
	}

	assertParsedOptions(t, opts)
}

func assertParsedOptions(t *testing.T, opts DeployRepoOptions) {
	if opts.Repo != "awesome-repo" {
		t.Errorf("expected repo awesome-repo, got %q", opts.Repo)
	}
	if opts.Target != "node-alpha" {
		t.Errorf("expected target node-alpha, got %q", opts.Target)
	}
	if opts.Dest != "D:\\work\\awesome" {
		t.Errorf("expected dest D:\\work\\awesome, got %q", opts.Dest)
	}
	if !opts.WithIDEs || !opts.WithPinned || !opts.WithConversations {
		t.Errorf("expected IDE and conversation flags set")
	}
	if !opts.Clean || !opts.Clone || !opts.DryRun || !opts.IsJSON {
		t.Errorf("expected execution modifier flags set")
	}
	if !opts.IncludeMain || !opts.OpenOnly {
		t.Errorf("expected include-main and open-only flags set")
	}
	if len(opts.Include) != 2 || len(opts.Except) != 1 {
		t.Errorf("expected 2 includes and 1 except, got %d and %d", len(opts.Include), len(opts.Except))
	}
}

func TestArchiveCreationFiltering(t *testing.T) {
	tmpDir := t.TempDir()
	setupTestRepoFiles(t, tmpDir)

	dataUnclean, _, errUnclean := PackageRepoArchive(tmpDir, DeployRepoOptions{Clean: false})
	if errUnclean != nil {
		t.Fatalf("unclean package failed: %v", errUnclean)
	}

	uncleanEntries := readTarEntries(t, dataUnclean)
	if !containsEntry(uncleanEntries, "node_modules/index.js") {
		t.Errorf("expected node_modules/index.js in unclean archive")
	}

	dataClean, _, errClean := PackageRepoArchive(tmpDir, DeployRepoOptions{Clean: true})
	if errClean != nil {
		t.Fatalf("clean package failed: %v", errClean)
	}

	cleanEntries := readTarEntries(t, dataClean)
	assertCleanArchiveEntries(t, cleanEntries)
}

func setupTestRepoFiles(t *testing.T, baseDir string) {
	files := []string{
		"main.go",
		"README.md",
		"node_modules/index.js",
		"target/debug/app",
		".venv/pyvenv.cfg",
		"dist/bundle.js",
		".git/config",
		".git/objects/pack/test.pack",
	}

	for _, f := range files {
		fullPath := filepath.Join(baseDir, filepath.FromSlash(f))
		dir := filepath.Dir(fullPath)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatalf("failed to create dir %s: %v", dir, err)
		}
		if err := os.WriteFile(fullPath, []byte("content: "+f), 0644); err != nil {
			t.Fatalf("failed to create file %s: %v", fullPath, err)
		}
	}
}

func readTarEntries(t *testing.T, data []byte) []string {
	gz, errGz := gzip.NewReader(bytes.NewReader(data))
	if errGz != nil {
		t.Fatalf("failed to create gzip reader: %v", errGz)
	}
	defer gz.Close()

	tr := tar.NewReader(gz)
	var entries []string
	for {
		hdr, err := tr.Next()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed reading tar header: %v", err)
		}
		entries = append(entries, hdr.Name)
	}
	return entries
}

func containsEntry(entries []string, target string) bool {
	for _, e := range entries {
		if e == target {
			return true
		}
	}
	return false
}

func assertCleanArchiveEntries(t *testing.T, entries []string) {
	if !containsEntry(entries, "main.go") {
		t.Errorf("expected main.go in clean archive")
	}
	if !containsEntry(entries, "README.md") {
		t.Errorf("expected README.md in clean archive")
	}
	if !containsEntry(entries, ".git/config") {
		t.Errorf("expected .git/config in clean archive")
	}

	excludedFiles := []string{
		"node_modules/index.js",
		"target/debug/app",
		".venv/pyvenv.cfg",
		"dist/bundle.js",
		".git/objects/pack/test.pack",
	}

	for _, ef := range excludedFiles {
		if containsEntry(entries, ef) {
			t.Errorf("unexpected excluded file in clean archive: %s", ef)
		}
	}
}

func TestDryRunSimulationBehavior(t *testing.T) {
	c1 := db.SSHConnection{Alias: "worker-1", IPAddress: "10.0.0.1", OS: "windows"}
	c2 := db.SSHConnection{Alias: "worker-2", IPAddress: "10.0.0.2", OS: "linux"}
	cMain := db.SSHConnection{Alias: "main", IPAddress: "10.0.0.3", OS: "linux"}
	cLocal := db.SSHConnection{Alias: "local", IPAddress: "127.0.0.1", OS: "windows"}

	allConns := []db.SSHConnection{c1, c2, cMain, cLocal}
	opts := DeployRepoOptions{
		NodeFilterOptions: NodeFilterOptions{IncludeMain: false},
		DryRun:            true,
		Repo:              "sample-repo",
	}

	filtered := FilterDeployRepoNodes(allConns, opts)
	if len(filtered) != 2 {
		t.Fatalf("expected 2 filtered nodes without main/local, got %d", len(filtered))
	}

	repoRec := &model.ScanRecord{
		Slug:     "sample-repo",
		RepoName: "sample-repo",
	}

	errDry := executeDryRunDeployRepo(filtered, opts, repoRec, "D:\\dummy\\path")
	if errDry != nil {
		t.Fatalf("unexpected error during dry run simulation: %v", errDry)
	}
}
