package cmd

// Diff summary emitter for `gitmap regoldens`. Reads the git working
// tree to report which `testdata/` golden files were touched by the
// pass-1 regenerate step. Output is intentionally concise: status
// letter, path, and (+adds / -dels) line counts per file plus an
// aggregate totals line. This file is split out so regoldens.go
// stays well under the 200-line cap.

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// goldenDiffEntry captures one changed testdata/ file. Status uses
// the porcelain letter set: A (added/untracked), M (modified),
// D (deleted), R (renamed), ? (untracked — normalized to A).
// renamedFrom is populated only for R entries (full mode).
type goldenDiffEntry struct {
	status      string
	path        string
	renamedFrom string
	added       int
	deleted     int
}

// Path-scope predicates (goldenDiffPathFragment, goldenDiffBasenameFragment,
// isGoldenFixturePath) live in regoldens_diff_scope.go to keep this
// file under the 200-line cap.

// emitGoldenDiffSummary prints the post-pass-1 diff summary in the
// requested mode ("short" | "full"). Errors from git invocations are
// surfaced (zero-swallow policy) but never fatal — the diff is
// informational and must not block pass 2.
func emitGoldenDiffSummary(mode string) {
	if !isGitWorkingTree() {
		fmt.Fprint(os.Stderr, constants.MsgRegoldensDiffSkipped)
		return
	}
	entries, err := collectGoldenDiffEntries()
	if err != nil {
		fmt.Fprintf(os.Stderr, "regoldens: diff summary failed: %v\n", err)
		return
	}
	fmt.Fprintf(os.Stdout, constants.MsgRegoldensDiffHeader, mode)
	if len(entries) == 0 {
		fmt.Fprint(os.Stdout, constants.MsgRegoldensDiffNoChanges)
		return
	}
	printGoldenDiffEntries(entries, mode)
}

// isGitWorkingTree returns true when the current directory is inside
// a git repository. Used to gate the diff feature gracefully.
func isGitWorkingTree() bool {
	cmd := exec.Command("git", "rev-parse", "--is-inside-work-tree")
	cmd.Stderr = nil
	out, err := cmd.Output()
	if err != nil {
		return false
	}
	return strings.TrimSpace(string(out)) == "true"
}

// collectGoldenDiffEntries merges porcelain status (covers
// untracked/added/deleted) with numstat (covers +/- line counts for
// modified files) into a unified per-file record list.
func collectGoldenDiffEntries() ([]goldenDiffEntry, error) {
	statuses, err := readPorcelainStatuses()
	if err != nil {
		return nil, fmt.Errorf("git status: %w", err)
	}
	numstat, err := readNumstatCounts()
	if err != nil {
		return nil, fmt.Errorf("git diff numstat: %w", err)
	}
	return mergeStatusAndNumstat(statuses, numstat), nil
}

// readPorcelainStatuses returns a map of testdata/ path -> entry
// (status letter + optional renamedFrom) from `git status --porcelain`.
// Untracked entries (`??`) are normalized to "A" (added).
func readPorcelainStatuses() (map[string]goldenDiffEntry, error) {
	out, err := runGitCapture("status", "--porcelain", "--", "*"+goldenDiffPathFragment+"*")
	if err != nil {
		return nil, err
	}
	result := make(map[string]goldenDiffEntry)
	for _, line := range strings.Split(out, "\n") {
		if len(line) < 4 {
			continue
		}
		path, from := splitPorcelainPath(line[3:])
		if !isGoldenFixturePath(path) {
			continue
		}
		result[path] = goldenDiffEntry{
			status:      normalizePorcelainStatus(strings.TrimSpace(line[:2])),
			path:        path,
			renamedFrom: from,
		}
	}
	return result, nil
}

// splitPorcelainPath extracts (newPath, oldPath) from the path
// portion of a porcelain status line. Renames are formatted as
// "old -> new"; everything else is a plain path. Whitespace is
// trimmed in either branch.
func splitPorcelainPath(raw string) (newPath, oldPath string) {
	if idx := strings.Index(raw, " -> "); idx >= 0 {
		return strings.TrimSpace(raw[idx+4:]), strings.TrimSpace(raw[:idx])
	}
	return strings.TrimSpace(raw), ""
}

// normalizePorcelainStatus collapses git's two-letter porcelain
// codes to a single display letter. Order matters: deletion wins
// over modification when both are present (rare but possible for
// staged-then-deleted files).
func normalizePorcelainStatus(code string) string {
	if code == "??" {
		return "A"
	}
	if strings.Contains(code, "D") {
		return "D"
	}
	if strings.Contains(code, "R") {
		return "R"
	}
	if strings.Contains(code, "A") {
		return "A"
	}
	return "M"
}

// readNumstatCounts returns added/deleted line counts for tracked
// modifications. Untracked files do not appear here — that's why
// readPorcelainStatuses runs in parallel.
func readNumstatCounts() (map[string][2]int, error) {
	out, err := runGitCapture("diff", "--numstat", "--", "*"+goldenDiffPathFragment+"*")
	if err != nil {
		return nil, err
	}
	result := make(map[string][2]int)
	for _, line := range strings.Split(out, "\n") {
		fields := strings.Fields(line)
		if len(fields) < 3 {
			continue
		}
		if !isGoldenFixturePath(fields[2]) {
			continue
		}
		added, _ := strconv.Atoi(fields[0]) // "-" (binary) becomes 0
		deleted, _ := strconv.Atoi(fields[1])
		result[fields[2]] = [2]int{added, deleted}
	}
	return result, nil
}

// mergeStatusAndNumstat joins the two maps by path, defaulting line
// counts to 0 for added/deleted/untracked files (numstat omits them).
// renamedFrom is preserved from the porcelain entry.
func mergeStatusAndNumstat(statuses map[string]goldenDiffEntry, counts map[string][2]int) []goldenDiffEntry {
	entries := make([]goldenDiffEntry, 0, len(statuses))
	for _, e := range statuses {
		c := counts[e.path]
		e.added = c[0]
		e.deleted = c[1]
		entries = append(entries, e)
	}
	sortGoldenDiffEntries(entries)
	return entries
}

// runGitCapture executes a git subcommand and returns trimmed stdout.
// stderr is captured into the returned error so failures are visible.
func runGitCapture(args ...string) (string, error) {
	cmd := exec.Command("git", args...)
	var stderr strings.Builder
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		var exitErr *exec.ExitError
		if errors.As(err, &exitErr) {
			return "", fmt.Errorf("%w: %s", err, strings.TrimSpace(stderr.String()))
		}
		return "", err
	}
	return string(out), nil
}
