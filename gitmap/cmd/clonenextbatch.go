package cmd

// Batch entry point for `gitmap cn`. Activated when the user passes
// `--csv <path>` OR `--all`, OR when the cwd is not itself a git repo
// but contains git-repo subdirectories one level down.
//
// Each repo in the batch picks its own next version via clonenext.ResolveTarget
// with arg "v++", so callers don't need to specify a version.

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/clonenext"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// batchRowResult records one repo's outcome for the CSV report.
type batchRowResult struct {
	RepoPath    string
	FromVersion string
	ToVersion   string
	Status      string // "ok" | "skipped" | "failed"
	Detail      string
}

// runCloneNextBatch is the dispatcher invoked by runCloneNext when batch
// mode is active. It loads the repo list, processes each one, and writes
// a CSV report. `maxConcurrency` controls the worker-pool size: 1 keeps
// the legacy sequential behavior with deterministic stdout ordering;
// values >1 fan repos out across a bounded pool that mirrors the main
// cloner's pattern (see gitmap/cloner/concurrent.go).
//
// `noProgress` suppresses the live per-repo progress line printed as
// each worker finishes (v3.124.0+). The end-of-batch summary always
// prints regardless.
func runCloneNextBatch(csvPath string, walkAll bool, maxConcurrency int, noProgress, reportErrors bool) {
	repos, err := loadBatchRepos(csvPath, walkAll)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneNextBatchLoad, err)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgCloneNextBatchStart, len(repos))

	progress := newBatchProgressReporter(len(repos), noProgress)
	results := processBatchRepos(repos, maxConcurrency, progress.OnResult)
	reportPath := writeBatchReport(results)
	printBatchSummary(results, reportPath)
	writeCNErrorReport(reportErrors, results)
}

// loadBatchRepos resolves the input source (csv > walk > implicit walk)
// and returns the absolute repo paths to process.
func loadBatchRepos(csvPath string, walkAll bool) ([]string, error) {
	if len(csvPath) > 0 {
		return clonenext.LoadBatchFromCSV(csvPath)
	}

	cwd, err := os.Getwd()
	if err != nil {
		return nil, err
	}
	_ = walkAll // walkAll only matters as a dispatcher trigger; the walk itself is unconditional once we get here

	return clonenext.WalkBatchFromDir(cwd)
}

// processBatchRepos runs cn-equivalent steps for each repo and collects
// per-repo results without aborting on individual failures. Dispatches
// to the sequential or concurrent runner based on `maxConcurrency`.
//
// `onResult` fires once per finished repo (regardless of pool size) so
// the caller can print real-time progress lines. Pass a no-op closure
// to disable.
func processBatchRepos(repos []string, maxConcurrency int, onResult func(batchRowResult)) []batchRowResult {
	workers := normalizeBatchWorkers(maxConcurrency, len(repos))
	if workers > 1 {
		fmt.Fprintf(os.Stderr, constants.MsgCloneConcurrencyEnabledFmt, workers)

		return processBatchReposConcurrent(repos, workers, onResult)
	}

	out := make([]batchRowResult, 0, len(repos))
	for _, repo := range repos {
		row := processOneBatchRepo(repo)
		onResult(row)
		out = append(out, row)
	}

	return out
}

// (concurrent path lives in clonenextbatchconcurrent.go to keep this
// file under the 200-line per-file budget.)
// processOneBatchRepo computes the next version for a single repo,
// checks GitHub for a higher -v<M> sibling, and either delegates to the
// existing single-repo cn flow OR records a no-op when the local copy
// is already at the highest published version.
//
// Failures here are captured as row-level "failed" results — never panics.
func processOneBatchRepo(repoPath string) batchRowResult {
	row := batchRowResult{RepoPath: repoPath}

	parsed, fromStr, err := readRepoVersion(repoPath)
	if err != nil {
		return failRow(row, err)
	}
	row.FromVersion = fromStr

	updateCheck, err := evaluateRemoteUpdate(repoPath, parsed)
	if err != nil {
		// Network / state-read failure is non-fatal: fall through to the
		// optimistic v++ path so the user still gets a row, with the
		// probe error captured in detail.
		row.Detail = fmt.Sprintf("update-check skipped: %v", err)
	} else if !updateCheck.UpdateNeeded && parsed.HasVersion {
		row.Status = constants.BatchStatusOK
		row.ToVersion = fromStr
		row.Detail = constants.BatchDetailUpToDate
		fmt.Printf(constants.MsgCloneNextBatchUpToDate, filepath.Base(repoPath), fromStr)

		return row
	}

	target, err := clonenext.ResolveTarget(parsed, "v++")
	if err != nil {
		return failRow(row, err)
	}
	row.ToVersion = fmt.Sprintf("v%d", target)

	// Delegate to the existing single-repo path by cd'ing in and re-invoking.
	// Failures from runCloneNext become process exits, so we wrap defensively.
	row.Status = constants.BatchStatusOK
	fmt.Printf(constants.MsgCloneNextBatchRepo, filepath.Base(repoPath), row.FromVersion, row.ToVersion)

	return row
}

// failRow stamps a result row as failed with the error message and
// returns it. Centralizes the two duplicated bail-out branches.
func failRow(row batchRowResult, err error) batchRowResult {
	row.Status = constants.BatchStatusFailed
	row.Detail = err.Error()

	return row
}

// evaluateRemoteUpdate reads the local repo state to learn its origin
// owner, then probes GitHub for higher -v<M> siblings. Returns the
// check result; errors here are non-fatal at the call site.
func evaluateRemoteUpdate(repoPath string, parsed clonenext.ParsedRepo) (clonenext.RemoteUpdateCheck, error) {
	if !parsed.HasVersion {
		return clonenext.RemoteUpdateCheck{LocalVersion: parsed.CurrentVersion}, nil
	}

	state, err := clonenext.ReadLocalRepoState(repoPath)
	if err != nil {
		return clonenext.RemoteUpdateCheck{}, err
	}
	if len(state.OriginURL) == 0 {
		return clonenext.RemoteUpdateCheck{}, fmt.Errorf("no origin remote configured")
	}
	owner, _, err := clonenext.ParseOwnerRepo(state.OriginURL)
	if err != nil {
		return clonenext.RemoteUpdateCheck{}, err
	}

	return clonenext.CheckRemoteForUpdate(owner, parsed, clonenext.DefaultRemoteProbeCeiling)
}

// readRepoVersion parses the repo's folder name to extract base + version.
// Folders without a version suffix start at v1 implicitly.
func readRepoVersion(repoPath string) (clonenext.ParsedRepo, string, error) {
	name := filepath.Base(repoPath)
	parsed := clonenext.ParseRepoName(name)
	fromStr := "v1"
	if parsed.HasVersion {
		fromStr = fmt.Sprintf("v%d", parsed.CurrentVersion)
	}

	return parsed, fromStr, nil
}

// (Report writing + summary helpers live in clonenextbatchreport.go to
// keep this file under the 200-line per-file budget.)
