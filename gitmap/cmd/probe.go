package cmd

// probe.go — `gitmap probe` dispatcher and capped worker pool.
//
// JSON shaping, persistence, URL preference, and the per-repo summary
// line live in probereport.go. Flag parsing lives in probeflags.go.
// This file owns just two responsibilities:
//
//  1. Top-level dispatch: parse flags, resolve targets, hand off.
//  2. Fan-out: stand up an opts.workers-sized pool, slot results back
//     into input order, and serialize counter updates through one mutex.

import (
	"fmt"
	"os"
	"path/filepath"
	"sync"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// runProbe dispatches `gitmap probe [<repo-path>|--all] [--json] [--workers N]`.
// The probe pool is capped at constants.ProbeMaxWorkers (default 2) to
// stay under provider rate limits.
func runProbe(args []string) {
	checkHelp("probe", args)
	opts := mustParseProbeArgs(args)

	db := openSfDB()
	defer db.Close()

	targets := mustResolveProbeTargets(db, opts.rest)
	if len(targets) == 0 {
		emitProbeEmpty(opts.jsonOut)
		return
	}
	probeAndReport(db, targets, opts)
}

// mustParseProbeArgs is a fatal wrapper around parseProbeArgs.
func mustParseProbeArgs(args []string) probeOptions {
	opts, err := parseProbeArgs(args)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return opts
}

// mustResolveProbeTargets is a fatal wrapper around resolveProbeTargets.
func mustResolveProbeTargets(db *store.DB, rest []string) []model.ScanRecord {
	targets, err := resolveProbeTargets(db, rest)
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}

	return targets
}

// emitProbeEmpty handles the "no targets" case in either output mode.
func emitProbeEmpty(jsonOut bool) {
	if jsonOut {
		fmt.Println("[]")
		return
	}
	fmt.Print(constants.MsgProbeNoTargets)
}

// resolveProbeTargets converts CLI args into a list of repos to probe.
func resolveProbeTargets(db *store.DB, args []string) ([]model.ScanRecord, error) {
	if len(args) == 0 || args[0] == constants.ProbeFlagAll {
		return db.ListRepos()
	}

	absPath, err := filepath.Abs(args[0])
	if err != nil {
		return nil, fmt.Errorf(constants.ErrSFAbsResolve, args[0], err)
	}

	matches, err := db.FindByPath(absPath)
	if err != nil {
		return nil, err
	}
	if len(matches) == 0 {
		return nil, fmt.Errorf(constants.ErrProbeNoRepo, absPath)
	}

	return matches, nil
}

// probeAndReport executes the probe across opts.workers goroutines,
// then emits either the human summary or a JSON array.
func probeAndReport(db *store.DB, targets []model.ScanRecord, opts probeOptions) {
	if !opts.jsonOut {
		fmt.Printf(constants.MsgProbeStartFmt, len(targets))
	}

	entries, available, unchanged, failed := runProbePool(db, targets, opts)

	// JSON wins only when the user did NOT also pass --output terminal;
	// explicit terminal opts in to the human format and trumps json.
	if opts.jsonOut && !opts.termOut {
		emitProbeJSON(entries)
		return
	}
	fmt.Printf(constants.MsgProbeDoneFmt, available, unchanged, failed)
}

// probeJob is a single unit of work for the probe pool.
type probeJob struct {
	idx  int
	repo model.ScanRecord
}

// runProbePool fans the targets across opts.workers goroutines. Entries
// are slotted back into input order so the JSON output is deterministic
// regardless of completion order. Per-repo human progress lines print as
// workers finish (matches the cloner pattern); the trailing summary
// totals are guarded by counterMu so concurrent tallies cannot lose
// updates.
func runProbePool(db *store.DB, targets []model.ScanRecord, opts probeOptions) ([]probeJSONEntry, int, int, int) {
	jobs := make(chan probeJob, len(targets))
	entries := make([]probeJSONEntry, len(targets))
	var counterMu sync.Mutex
	available, unchanged, failed := 0, 0, 0

	var wg sync.WaitGroup
	for w := 0; w < opts.workers; w++ {
		wg.Add(1)
		go probeWorker(db, jobs, entries, &counterMu, &available, &unchanged, &failed, opts, &wg)
	}
	for i, repo := range targets {
		jobs <- probeJob{idx: i, repo: repo}
	}
	close(jobs)
	wg.Wait()

	return entries, available, unchanged, failed
}

// probeWorker drains the job channel, writes its result to its own
// index in entries (no contention — each index is owned by exactly one
// job), and serializes counter/print updates through counterMu so the
// final tallies and the human progress lines stay coherent. opts is
// passed by value so each worker reads jsonOut/depth without locking.
func probeWorker(db *store.DB, jobs <-chan probeJob, entries []probeJSONEntry,
	counterMu *sync.Mutex, available, unchanged, failed *int, opts probeOptions, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		result := executeOneProbe(db, j.repo, opts.depth)
		entries[j.idx] = makeProbeEntry(j.repo, result)
		counterMu.Lock()
		// In --output terminal mode we suppress the legacy 1-line
		// "ok/none/fail" tally print so the standardized block is
		// the only per-repo output. Counters still update.
		quietTally := opts.jsonOut || opts.termOut
		*available, *unchanged, *failed = tallyProbe(j.repo, result, *available, *unchanged, *failed, quietTally)
		if opts.termOut {
			emitProbeTermBlock(j.idx+1, j.repo, result)
		}
		counterMu.Unlock()
	}
}
