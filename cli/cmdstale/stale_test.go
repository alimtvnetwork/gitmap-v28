package cmdstale

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstale/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdstale/enums"
)

func TestParseStaleFlagsDefaults(t *testing.T) {
	opts, err := parseStaleFlags([]string{})
	if err != nil {
		t.Fatalf("unexpected error parsing default flags: %v", err)
	}
	if opts.staleDays != constants.DefaultStaleDays {
		t.Errorf("expected default stale days %d, got %d", constants.DefaultStaleDays, opts.staleDays)
	}
	if opts.scanRoot != constants.DefaultScanRoot {
		t.Errorf("expected default root %q, got %q", constants.DefaultScanRoot, opts.scanRoot)
	}
	if opts.isArchive || opts.isDryRun {
		t.Error("expected isArchive and isDryRun to be false by default")
	}
	if !opts.outputFormat.IsTable() {
		t.Errorf("expected default format table, got %s", opts.outputFormat)
	}
}

func TestParseStaleFlagsCustom(t *testing.T) {
	args := []string{
		constants.FlagArgDays, "45",
		constants.FlagArgRoot, "custom_root",
		constants.FlagArgArchive,
		constants.FlagArgDryRun,
		constants.FlagArgFormat, "json",
	}
	opts, err := parseStaleFlags(args)
	if err != nil {
		t.Fatalf("unexpected error parsing custom flags: %v", err)
	}
	if opts.staleDays != 45 || opts.scanRoot != "custom_root" {
		t.Errorf("mismatch in custom flags: %+v", opts)
	}
	if !opts.isArchive || !opts.isDryRun || !opts.outputFormat.IsJSON() {
		t.Errorf("flags not parsed properly: %+v", opts)
	}
}

func TestParseStaleFlagsErrors(t *testing.T) {
	testInvalidFlag(t)
	testInvalidFormat(t)
}

func testInvalidFlag(t *testing.T) {
	_, err := parseStaleFlags([]string{"--unsupported-flag"})
	if err == nil {
		t.Error("expected error for unsupported flag")
	}
}

func testInvalidFormat(t *testing.T) {
	_, err := parseStaleFlags([]string{constants.FlagArgFormat, "yaml"})
	if err == nil {
		t.Error("expected error for invalid format")
	}
}

func TestRunStaleValidationErrors(t *testing.T) {
	if err := RunStale([]string{"--bad-arg"}); err == nil {
		t.Error("expected RunStale to return error on bad argument")
	}
	if err := RunStale([]string{constants.FlagArgFormat, "xml"}); err == nil {
		t.Error("expected RunStale to return error on bad format")
	}
}

func TestParseEpochOutput(t *testing.T) {
	testEpochValid(t)
	testEpochEmpty(t)
	testEpochInvalid(t)
}

func testEpochValid(t *testing.T) {
	parsed, hasCommit := parseEpochOutput("1600000000\n")
	if !hasCommit {
		t.Error("expected hasCommit to be true for valid epoch")
	}
	if parsed.Unix() != 1600000000 {
		t.Errorf("expected unix timestamp 1600000000, got %d", parsed.Unix())
	}
}

func testEpochEmpty(t *testing.T) {
	_, hasCommit := parseEpochOutput("   ")
	if hasCommit {
		t.Error("expected hasCommit to be false for empty epoch")
	}
}

func testEpochInvalid(t *testing.T) {
	_, hasCommit := parseEpochOutput("not-a-number")
	if hasCommit {
		t.Error("expected hasCommit to be false for non-numeric epoch")
	}
}

func TestEmitStaleOutputs(t *testing.T) {
	repos := []staleRepo{
		{path: "repo1", last: time.Now().Add(-100 * 24 * time.Hour)},
		{path: "repo2", last: time.Now().Add(-120 * 24 * time.Hour)},
	}
	emitStale(repos, 90, enums.StaleOutputFormatTable)
	emitStale(repos, 90, enums.StaleOutputFormatJSON)
	emitStale(repos, 90, enums.StaleOutputFormatCSV)
	emitStale([]staleRepo{}, 90, enums.StaleOutputFormatTable)
}

func TestArchiveStaleReposDryRun(t *testing.T) {
	repos := []staleRepo{
		{path: "fake/path/repo", last: time.Now()},
	}
	archiveStaleRepos([]staleRepo{}, true)
	archiveStaleRepos(repos, true)
}

func TestArchiveStaleReposExecution(t *testing.T) {
	tmpDir := t.TempDir()
	repoPath := filepath.Join(tmpDir, "stale-repo")
	if err := os.Mkdir(repoPath, 0o755); err != nil {
		t.Fatalf("failed to create temp dir: %v", err)
	}

	archiveDir := filepath.Join(tmpDir, "archive")
	dest := filepath.Join(archiveDir, "stale-repo")
	executeArchiveMove(repoPath, dest, archiveDir)
	if _, err := os.Stat(dest); err != nil {
		t.Errorf("expected archived repo at dest, err: %v", err)
	}
}

func TestHygieneWorkersBoundaries(t *testing.T) {
	workers := hygieneWorkers()
	if workers < constants.MinHygieneWorkers || workers > constants.MaxHygieneWorkers {
		t.Errorf("workers %d outside bounds [%d, %d]", workers, constants.MinHygieneWorkers, constants.MaxHygieneWorkers)
	}
}
