package cmdscan

// scan_diff.go — `gitmap scan diff` / `gitmap scand`: compare the current
// scan index (the DB) against the last-snapshot state: added, removed, and
// moved repos (program 275, proposal P-17). Read-only: no dry-run needed
// (shared-patterns spec §6). One concern per file (≤300 lines).

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// scan diff messages.
const (
	scanDiffSnapshotName = "scan-diff-snapshot.json"
	scanDiffUsageLine    = "Usage: gitmap scan diff\n       gitmap scand   (alias; identical behavior)\n"
	scanDiffBaselineFmt  = "baseline snapshot saved (%d repos) · run again after changes to see a diff\n"
	scanDiffNoChangeFmt  = "no changes since %s (%d repos tracked)\n"
)

// scanDiffEntry is one repo's identity in the snapshot file.
type scanDiffEntry struct {
	Slug string `json:"slug"`
	Path string `json:"path"`
	URL  string `json:"url"`
}

// scanDiffSnapshot is the persisted last-scan state.
type scanDiffSnapshot struct {
	SavedAt string                   `json:"savedAt"`
	Repos   map[string]scanDiffEntry `json:"repos"`
}

// RunScanDiff implements `gitmap scan diff` / `gitmap scand` for rootdata wiring.
func RunScanDiff(args []string) error {
	if appErr := runScanDiff(args); appErr != nil {
		return appErr
	}
	return nil
}

// runScanDiff executes the diff flow: load index, load snapshot, compare, save.
func runScanDiff(args []string) *appfault.AppError {
	if hasHelpArg(args) {
		printScanDiffHelp()
		return nil
	}
	db, appErr := openScanDiffDB()
	if appErr != nil {
		return appErr
	}
	defer db.Close()
	records, appErr := listScanDiffRepos(db)
	if appErr != nil {
		return appErr
	}
	current := buildScanDiffEntries(records)
	snapPath := scanDiffSnapshotPath()
	previous, isFirst, appErr := loadScanDiffSnapshot(snapPath)
	if appErr != nil {
		return appErr
	}
	if isFirst {
		return saveScanDiffBaseline(snapPath, current)
	}
	printScanDiffReport(previous, current)
	return saveScanDiffSnapshot(snapPath, current)
}

// saveScanDiffBaseline saves the first snapshot and reports the baseline.
func saveScanDiffBaseline(snapPath string, current map[string]scanDiffEntry) *appfault.AppError {
	if appErr := saveScanDiffSnapshot(snapPath, current); appErr != nil {
		return appErr
	}
	fmt.Printf(scanDiffBaselineFmt, len(current))

	return nil
}

// openScanDiffDB opens the default gitmap database.
func openScanDiffDB() (*store.DB, *appfault.AppError) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, apperror.WrapSimple(err, "scan diff: could not open the gitmap database")
	}
	return db, nil
}

// listScanDiffRepos reads the current scan index.
func listScanDiffRepos(db *store.DB) ([]model.ScanRecord, *appfault.AppError) {
	records, err := db.ListRepos()
	if err != nil {
		return nil, apperror.WrapSimple(err, "scan diff: could not list repositories")
	}
	return records, nil
}

// scanDiffURL picks the identity URL for an entry (same priority as export).
func scanDiffURL(r model.ScanRecord) string {
	if r.HTTPSUrl != "" {
		return r.HTTPSUrl
	}
	if r.SSHUrl != "" {
		return r.SSHUrl
	}
	return r.DiscoveredURL
}

// buildScanDiffEntries indexes records by normalized absolute path.
func buildScanDiffEntries(records []model.ScanRecord) map[string]scanDiffEntry {
	entries := make(map[string]scanDiffEntry, len(records))
	for _, r := range records {
		key := normalizeScanDiffPath(r.AbsolutePath)
		if key == "" {
			continue
		}
		entries[key] = scanDiffEntry{Slug: r.Slug, Path: r.AbsolutePath, URL: scanDiffURL(r)}
	}
	return entries
}

// normalizeScanDiffPath lowercases + slash-normalizes a path for comparison.
func normalizeScanDiffPath(p string) string {
	return strings.ToLower(filepath.Clean(p))
}

// scanDiffSnapshotPath stores the snapshot next to the active database.
func scanDiffSnapshotPath() string {
	return filepath.Join(filepath.Dir(store.DefaultDBPath()), scanDiffSnapshotName)
}

// loadScanDiffSnapshot reads the snapshot; isFirst=true when none exists.
func loadScanDiffSnapshot(path string) (scanDiffSnapshot, bool, *appfault.AppError) {
	var snap scanDiffSnapshot
	data, err := os.ReadFile(path)
	if err != nil {
		return scanDiffReadErr(path, err, snap)
	}
	if err := json.Unmarshal(data, &snap); err != nil {
		return snap, false, appfault.NewValidationError(fmt.Sprintf("scan diff: snapshot %s is corrupt: %v (delete it to re-baseline)", path, err))
	}
	if snap.Repos == nil {
		snap.Repos = map[string]scanDiffEntry{}
	}
	return snap, false, nil
}

// scanDiffReadErr maps a snapshot read error to (snapshot, isFirst, appErr).
func scanDiffReadErr(path string, err error, snap scanDiffSnapshot) (scanDiffSnapshot, bool, *appfault.AppError) {
	if os.IsNotExist(err) {
		return snap, true, nil
	}

	return snap, false, apperror.NewExecutionError(fmt.Sprintf("scan diff: could not read snapshot %s: %v", path, err))
}

// saveScanDiffSnapshot persists the current index state.
func saveScanDiffSnapshot(path string, entries map[string]scanDiffEntry) *appfault.AppError {
	snap := scanDiffSnapshot{SavedAt: nowUTC(), Repos: entries}
	data, err := json.MarshalIndent(snap, "", "  ")
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("scan diff: could not encode snapshot: %v", err))
	}
	if err := os.WriteFile(path, append(data, '\n'), 0o644); err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("scan diff: could not write snapshot %s: %v", path, err))
	}
	return nil
}

// scanDiffMove records one moved repo (old path -> new path).
type scanDiffMove struct {
	oldPath string
	newPath string
	slug    string
}

// diffScanEntries compares previous vs current snapshots.
func diffScanEntries(previous scanDiffSnapshot, current map[string]scanDiffEntry) (added, removed []scanDiffEntry, moved []scanDiffMove) {
	bySlug := map[string]string{}
	byURL := map[string]string{}
	for key, e := range previous.Repos {
		if e.Slug != "" {
			bySlug[e.Slug] = key
		}
		if e.URL != "" {
			byURL[e.URL] = key
		}
	}
	movedFrom := map[string]bool{}
	for key, e := range current {
		if _, kept := previous.Repos[key]; kept {
			continue
		}
		oldKey := ""
		if e.Slug != "" {
			oldKey = bySlug[e.Slug]
		}
		if oldKey == "" && e.URL != "" {
			oldKey = byURL[e.URL]
		}
		if oldKey != "" {
			moved = append(moved, scanDiffMove{oldPath: previous.Repos[oldKey].Path, newPath: e.Path, slug: e.Slug})
			movedFrom[oldKey] = true
			continue
		}
		added = append(added, e)
	}
	for key, e := range previous.Repos {
		if _, kept := current[key]; kept {
			continue
		}
		if movedFrom[key] {
			continue
		}
		removed = append(removed, e)
	}
	return added, removed, moved
}

// printScanDiffReport prints the proposal P-17 summary line.
func printScanDiffReport(previous scanDiffSnapshot, current map[string]scanDiffEntry) {
	added, removed, moved := diffScanEntries(previous, current)
	if len(added) == 0 && len(removed) == 0 && len(moved) == 0 {
		fmt.Printf(scanDiffNoChangeFmt, previous.SavedAt, len(current))
		return
	}
	parts := []string{}
	if len(added) > 0 {
		parts = append(parts, fmt.Sprintf("+ %d repo(s) (%s)", len(added), joinScanDiffPaths(added)))
	}
	if len(removed) > 0 {
		parts = append(parts, fmt.Sprintf("- %d repo(s) (%s)", len(removed), joinScanDiffPaths(removed)))
	}
	for _, m := range moved {
		parts = append(parts, fmt.Sprintf("~ 1 moved (%s -> %s)", shortenHome(m.oldPath), shortenHome(m.newPath)))
	}
	fmt.Println(strings.Join(parts, " · "))
}

// joinScanDiffPaths renders entry paths with ~ shortening.
func joinScanDiffPaths(entries []scanDiffEntry) string {
	paths := make([]string, 0, len(entries))
	for _, e := range entries {
		paths = append(paths, shortenHome(e.Path))
	}
	return strings.Join(paths, ", ")
}

// shortenHome replaces a $HOME prefix with ~ for display.
func shortenHome(p string) string {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return p
	}
	if p == home {
		return "~"
	}
	if strings.HasPrefix(p, home+string(os.PathSeparator)) {
		return "~" + p[len(home):]
	}
	return p
}

// printScanDiffHelp prints usage for `gitmap scan diff`.
func printScanDiffHelp() {
	fmt.Print(`gitmap scan diff — show what changed since the last scan snapshot

USAGE
  gitmap scan diff
  gitmap scand   (alias; identical behavior)

EXAMPLE
  gitmap scan diff

NOTES
  - Compares the current scan index against the last snapshot: added,
    removed, and moved repos (same slug/URL, different path).
  - First run saves a baseline snapshot next to the database; the snapshot
    refreshes after every successful diff.
  - Read-only: the index is never modified. For index-vs-disk truth, see
    'gitmap scan verify'.
`)
}

// Help-example data lives in cli/helpdoc/examples_275.go
// (single source of truth, program 275).
