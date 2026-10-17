package cmdscan

// scan_verify.go — `gitmap scan verify` / `gitmap scanverify`: two-way
// index↔filesystem consistency check (program 275, proposal C-08). It finds
// DB entries pointing at deleted/moved dirs (stale) AND on-disk repos
// missing from the index (unindexed). `--repair` prunes stale entries.
// Read-only by default (shared-patterns spec §6). One concern per file.

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// scan verify messages.
const (
	scanVerifyFlagRepair = "--repair"
	scanVerifyReportFmt  = "index %d repos · %d on disk ok · %d stale entries (%s) · %d unindexed (%s)%s\n"
	scanVerifyPrunedFmt  = "pruned %d stale entr(ies) from the index\n"
	scanVerifyRepairHint = " · run with --repair to prune"
)

// nowUTC returns the current UTC time in RFC3339 (shared scan helper).
func nowUTC() string {
	return time.Now().UTC().Format(time.RFC3339)
}

// RunScanVerify implements `gitmap scan verify` / `gitmap scanverify` for rootdata wiring.
func RunScanVerify(args []string) error {
	if appErr := runScanVerify(args); appErr != nil {
		return appErr
	}
	return nil
}

// runScanVerify executes the verify flow: load index, check both directions,
// optionally repair.
func runScanVerify(args []string) *appfault.AppError {
	if hasHelpArg(args) {
		printScanVerifyHelp()
		return nil
	}
	isRepair := hasScanVerifyRepairFlag(args)
	db, appErr := openScanVerifyDB()
	if appErr != nil {
		return appErr
	}
	defer db.Close()
	records, appErr := listScanVerifyRepos(db)
	if appErr != nil {
		return appErr
	}
	index := indexScanVerifyPaths(records)
	stale := findScanVerifyStale(index)
	unindexed := findScanVerifyUnindexed(index)
	pruned := 0
	if isRepair {
		var repairErr *appfault.AppError
		pruned, repairErr = pruneScanVerifyStale(db, stale)
		if repairErr != nil {
			return repairErr
		}
		stale = stale[pruned:]
	}
	printScanVerifyReport(len(records), stale, unindexed, isRepair, pruned)
	return nil
}

// hasScanVerifyRepairFlag reports whether --repair was passed.
func hasScanVerifyRepairFlag(args []string) bool {
	for _, a := range args {
		if a == scanVerifyFlagRepair {
			return true
		}
	}
	return false
}

// openScanVerifyDB opens the default gitmap database.
func openScanVerifyDB() (*store.DB, *appfault.AppError) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "scan verify: could not open the gitmap database")
	}
	return db, nil
}

// listScanVerifyRepos reads the current scan index.
func listScanVerifyRepos(db *store.DB) ([]model.ScanRecord, *appfault.AppError) {
	records, err := db.ListRepos()
	if err != nil {
		return nil, apperror.WrapSimple(err, "scan verify: could not list repositories")
	}
	return records, nil
}

// indexScanVerifyPaths indexes records by normalized absolute path.
func indexScanVerifyPaths(records []model.ScanRecord) map[string]model.ScanRecord {
	index := make(map[string]model.ScanRecord, len(records))
	for _, r := range records {
		if key := normalizeScanDiffPath(r.AbsolutePath); key != "" {
			index[key] = r
		}
	}
	return index
}

// findScanVerifyStale returns index entries whose directory is gone.
func findScanVerifyStale(index map[string]model.ScanRecord) []model.ScanRecord {
	var stale []model.ScanRecord
	for _, r := range index {
		if !isScanVerifyRepoDir(r.AbsolutePath) {
			stale = append(stale, r)
		}
	}
	return stale
}

// isScanVerifyRepoDir reports whether path is an existing directory.
func isScanVerifyRepoDir(path string) bool {
	info, err := os.Stat(path)
	return err == nil && info.IsDir()
}

// findScanVerifyUnindexed walks the parent dirs of indexed repos one level
// deep and returns .git dirs missing from the index.
func findScanVerifyUnindexed(index map[string]model.ScanRecord) []string {
	var unindexed []string
	seenParents := map[string]bool{}
	for key := range index {
		parent := filepath.Dir(key)
		if seenParents[parent] {
			continue
		}
		seenParents[parent] = true
		entries, err := os.ReadDir(parent)
		if err != nil {
			continue
		}
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			full := filepath.Join(parent, e.Name())
			if _, known := index[normalizeScanDiffPath(full)]; known {
				continue
			}
			if dirHasGitDir(full) {
				unindexed = append(unindexed, full)
			}
		}
	}
	return unindexed
}

// dirHasGitDir reports whether dir contains a .git dir or worktree file.
func dirHasGitDir(dir string) bool {
	info, err := os.Stat(filepath.Join(dir, ".git"))
	return err == nil && (info.IsDir() || info.Mode().IsRegular())
}

// pruneScanVerifyStale deletes stale index rows; returns the pruned count.
func pruneScanVerifyStale(db *store.DB, stale []model.ScanRecord) (int, *appfault.AppError) {
	pruned := 0
	for _, r := range stale {
		n, err := db.DeleteByPath(r.AbsolutePath)
		if err != nil {
			return pruned, apperror.WrapSimple(err, fmt.Sprintf("scan verify: could not prune %s", r.AbsolutePath))
		}
		if n > 0 {
			pruned++
		}
	}
	return pruned, nil
}

// printScanVerifyReport prints the proposal C-08 summary line.
func printScanVerifyReport(total int, stale []model.ScanRecord, unindexed []string, isRepair bool, pruned int) {
	onDisk := total - len(stale)
	stalePaths := make([]string, 0, len(stale))
	for _, r := range stale {
		stalePaths = append(stalePaths, shortenHome(r.AbsolutePath))
	}
	unindexedPaths := make([]string, 0, len(unindexed))
	for _, p := range unindexed {
		unindexedPaths = append(unindexedPaths, shortenHome(p))
	}
	hint := ""
	if len(stale) > 0 && !isRepair {
		hint = scanVerifyRepairHint
	}
	fmt.Printf(scanVerifyReportFmt, total, onDisk, len(stale),
		strings.Join(stalePaths, ", "), len(unindexed), strings.Join(unindexedPaths, ", "), hint)
	if isRepair {
		fmt.Printf(scanVerifyPrunedFmt, pruned)
	}
}

// printScanVerifyHelp prints usage for `gitmap scan verify`.
func printScanVerifyHelp() {
	fmt.Print(`gitmap scan verify — two-way index/filesystem consistency check

USAGE
  gitmap scan verify [--repair]
  gitmap scanverify [same flags]   (alias; identical behavior)

FLAGS
  --repair   Prune stale index entries (dirs that no longer exist)

EXAMPLE
  gitmap scan verify --repair

NOTES
  - Finds BOTH directions: index entries pointing at deleted/moved dirs
    (stale) and on-disk repos missing from the index (unindexed).
  - Read-only by default; --repair only deletes stale index rows, never
    touches the filesystem.
  - For "what changed since last scan", see 'gitmap scan diff'.
`)
}

// ScanVerify help-example data (3-part, spec 03 §2) for the lead's helpdoc wiring.
const (
	scanVerifyExampleWhat       = "Two-way index↔filesystem consistency check: stale DB entries for deleted dirs and on-disk repos missing from the index, with --repair to prune the stale ones."
	scanVerifyExampleInvocation = `gitmap scan verify`
	scanVerifyExampleOutput     = `index 43 repos · 41 on disk ok · 2 stale entries (~/tmp/scratch, ~/work/old) · 0 unindexed () · run with --repair to prune`
)
