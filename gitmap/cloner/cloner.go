// Package cloner re-clones repos from structured files.
package cloner

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/formatter"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// CloneOptions tunes a CloneFromFileWithOptions run. The zero value
// keeps the historical behavior: sequential, progress-on-stderr.
//
// MaxConcurrency:
//   - <= 1 → sequential (one repo at a time, original ordering).
//   - >  1 → bounded worker pool (see concurrent.go). Per-repo paths
//     come from each ScanRecord.RelativePath unchanged, so the on-disk
//     nested folder hierarchy is preserved regardless of worker count.
//
// Quiet suppresses per-repo progress lines but keeps the final summary
// (matches the legacy CloneFromFileQuiet behavior).
type CloneOptions struct {
	SafePull       bool
	Quiet          bool
	MaxConcurrency int
	// DefaultBranch is the fallback branch name handed to `git clone -b`
	// for any record whose recorded (Branch, BranchSource) would
	// otherwise leave the cloner with no usable branch (empty Branch,
	// detached HEAD, unknown source, etc.). Empty preserves the legacy
	// "let the remote's default HEAD decide" behavior. Plumbed in by
	// the CLI from `--default-branch` (constants.FlagScanDefaultBranch),
	// so the wording and semantics match `gitmap scan --default-branch`.
	DefaultBranch string
}

// CloneFromFile reads a source file and clones all repos under targetDir.
func CloneFromFile(sourcePath, targetDir string, safePull bool) (model.CloneSummary, error) {
	return CloneFromFileWithOptions(sourcePath, targetDir, CloneOptions{SafePull: safePull})
}

// CloneFromFileQuiet reads a source file and clones with suppressed progress.
func CloneFromFileQuiet(sourcePath, targetDir string, safePull bool) (model.CloneSummary, error) {
	return CloneFromFileWithOptions(sourcePath, targetDir, CloneOptions{SafePull: safePull, Quiet: true})
}

// CloneFromFileWithOptions is the full-control entry point. The legacy
// helpers above are thin wrappers that fill in CloneOptions defaults.
func CloneFromFileWithOptions(sourcePath, targetDir string, opts CloneOptions) (model.CloneSummary, error) {
	records, err := loadRecords(sourcePath)
	if err != nil {
		return model.CloneSummary{}, err
	}

	return cloneAll(records, targetDir, opts), nil
}

// loadRecords detects file format and parses records.
//
// Errors are wrapped with the source path so the CLI can surface
// "which file failed" without callers needing to know the original
// argument. Parser errors additionally carry their own line context
// (see parseTextFile) so users can jump straight to the offending row.
func loadRecords(path string) ([]model.ScanRecord, error) {
	ext := strings.ToLower(filepath.Ext(path))
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("open clone source %q: %w", path, err)
	}
	defer file.Close()

	records, err := parseByExtension(ext, file)
	if err != nil {
		return nil, fmt.Errorf("parse clone source %q: %w", path, err)
	}

	return records, nil
}

// parseByExtension dispatches to the correct parser.
func parseByExtension(ext string, r io.Reader) ([]model.ScanRecord, error) {
	if ext == constants.ExtCSV {
		return formatter.ParseCSV(r)
	}
	if ext == constants.ExtJSON {
		return formatter.ParseJSON(r)
	}

	return parseTextFile(r)
}

// parseTextFile reads one git clone command per line. Scanner errors
// are wrapped with the line number of the last successfully read line
// so users can locate malformed input in long clone manifests.
func parseTextFile(r io.Reader) ([]model.ScanRecord, error) {
	var records []model.ScanRecord
	sc := bufio.NewScanner(r)
	lineNum := 0
	for sc.Scan() {
		lineNum++
		line := strings.TrimSpace(sc.Text())
		if len(line) > 0 {
			records = append(records, parseCloneLine(line))
		}
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("read line %d: %w", lineNum+1, err)
	}

	return records, nil
}

// parseCloneLine extracts url, branch, path from a git clone command.
func parseCloneLine(line string) model.ScanRecord {
	parts := strings.Fields(line)
	rec := model.ScanRecord{CloneInstruction: line}
	if len(parts) >= 5 {
		rec.Branch = parts[3]
		rec.HTTPSUrl = parts[4]
	}
	if len(parts) >= 6 {
		rec.RelativePath = parts[5]
	}

	return rec
}

// (Dispatcher + sequential runner moved to runners.go to keep this file
// focused on entry points + parsing. The parallel runner lives in
// concurrent.go.)

// hasExistingRepos checks if any target repo directories already exist.
func hasExistingRepos(records []model.ScanRecord, targetDir string) bool {
	for _, rec := range records {
		dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
		if isGitRepo(dest) {
			return true
		}
	}

	return false
}

// cloneOne clones a single repository. Errors include the destination
// path and the record's RelativePath/RepoName so failures point straight
// at the offending row in the source manifest.
func cloneOne(rec model.ScanRecord, targetDir string) model.CloneResult {
	dest := filepath.Join(targetDir, model.CleanRelativePath(rec.RelativePath))
	err := os.MkdirAll(filepath.Dir(dest), constants.DirPermission)
	if err != nil {
		msg := fmt.Sprintf("mkdir %q for %s: %v", filepath.Dir(dest), recordTag(rec), err)

		return model.CloneResult{Record: rec, Success: false, Error: msg}
	}

	return runClone(rec, dest)
}

// runClone executes the git clone command.
//
// The branch-selection strategy is driven by ScanRecord.BranchSource so
// that records captured in a detached or unknown state never produce
// "Remote branch not found" errors. When the source is trusted (HEAD,
// remote-tracking, default) the recorded branch is passed via -b; when it
// is untrusted (detached, unknown) git clone is invoked without -b and
// the remote's default HEAD decides the checkout.
//
// Failures are formatted with the URL, branch, destination, and record
// tag so a single error line is enough to identify which manifest row
// failed and why — no cross-referencing the source file required.
func runClone(rec model.ScanRecord, dest string) model.CloneResult {
	url := pickURL(rec)
	strategy := pickCloneStrategy(rec)

	args := []string{constants.GitClone}
	if strategy.useBranch {
		args = append(args, constants.GitBranchFlag, strategy.branch)
	}
	args = append(args, url, dest)

	cmd := exec.Command(constants.GitBin, args...)
	if isSSHCloneURL(url) {
		return runInteractiveClone(cmd, rec, url, dest, strategy)
	}
	out, err := cmd.CombinedOutput()
	notes := strategy.reason
	if err != nil && isLFSSmudgeFailure(string(out)) {
		// Broken/missing LFS object on remote: retry with smudge disabled
		// so pointer files stay as pointers and checkout completes.
		retryOut, retryErr := retryCloneSkipSmudge(constants.GitBin, args, dest)
		if retryErr == nil {
			if notes == "" {
				notes = LFSRetryNote
			} else {
				notes = notes + "; " + LFSRetryNote
			}
			return model.CloneResult{Record: rec, Success: true, Notes: notes}
		}
		out, err = retryOut, retryErr
	}
	if err != nil {
		msg := fmt.Sprintf(
			"git clone failed for %s: url=%q branch=%q dest=%q: %v: %s",
			recordTag(rec), url, strategy.branch, dest, err, strings.TrimSpace(string(out)),
		)

		return model.CloneResult{Record: rec, Success: false, Error: msg, Notes: strategy.reason}
	}

	return model.CloneResult{Record: rec, Success: true, Notes: notes}
}

// (recordTag, pickURL, updateSummary, and updateSummarySkipped moved
// to summary.go so this file stays focused on entry points + parsing.)
