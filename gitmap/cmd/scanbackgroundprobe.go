package cmd

// Background-probe glue for `gitmap scan` (v3.123.0+).
//
// scan.go calls exactly two functions from this file:
//
//	probeRunner := startBackgroundProbe(records, probeOpts, quiet)
//	defer drainBackgroundProbe(probeRunner, probeOpts, quiet)
//
// All pool / sink / wait mechanics live behind those two calls so the
// scan command stays free of goroutine + channel + mutex code. The
// runner itself is implemented in gitmap/probe/background.go and
// already handles nil-receivers, so callers can pretend the
// "disabled" path is just a no-op.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/errreport"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/probe"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

// startBackgroundProbe decides whether the background probe should
// run for this scan and, if so, spins up a runner pre-loaded with
// every record. Returns nil when the runner is disabled (--no-probe,
// concurrency<=0, empty record set, or auto-trigger ceiling exceeded
// without an explicit --probe-concurrency flag).
//
// The runner is intentionally created with its own *store.DB handle
// because workers persist results from goroutines and the SQLite
// connection-pool restriction (`SetMaxOpenConns(1)`) means we want
// the background pool to share a single dedicated connection rather
// than fight with the foreground for the main handle.
func startBackgroundProbe(records []model.ScanRecord, opts ScanProbeOptions, quiet bool, errCollector *errreport.Collector) *probe.BackgroundRunner {
	workers := resolveProbeWorkers(records, opts, quiet)
	if workers < 1 {
		return nil
	}

	db, err := store.OpenDefault()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProbeOpenDB, err)

		return nil
	}
	if migErr := db.Migrate(); migErr != nil {
		fmt.Fprintln(os.Stderr, migErr.Error())
		db.Close()

		return nil
	}

	runner := probe.NewBackgroundRunner(workers, len(records),
		pickProbeURL,
		func(rec model.ScanRecord, res probe.Result) {
			recordProbeResult(db, rec, res)
		})
	// Install the failure hook AND clone-depth BEFORE enqueueing
	// jobs so workers observe non-default values on every dequeue.
	// Setting either after Start would race with the first probe.
	installProbeFailureHook(runner, errCollector)
	runner.SetCloneDepth(opts.Depth)

	enqueueProbeJobs(runner, records)
	if !quiet {
		fmt.Printf(constants.MsgScanProbeStartFmt, len(records), workers)
	}

	return runner
}

// resolveProbeWorkers applies the documented rules:
//   - --no-probe wins, returns 0.
//   - explicit --probe-concurrency wins next, even past the ceiling.
//   - otherwise auto-trigger: only fire when len(records) is below
//     the ceiling so big scans don't accidentally hammer GitHub.
//
// Empty record set always returns 0 (nothing to probe).
func resolveProbeWorkers(records []model.ScanRecord, opts ScanProbeOptions, quiet bool) int {
	if opts.Disable || len(records) == 0 {
		return 0
	}
	if opts.ConcurrencySet {
		return opts.Concurrency
	}
	if len(records) >= constants.ScanProbeAutoTriggerCeiling {
		if !quiet {
			fmt.Printf(constants.MsgScanProbeSkippedAutoFmt,
				len(records), constants.ScanProbeAutoTriggerCeiling)
		}

		return 0
	}

	return constants.ScanProbeDefaultConcurrency
}

// enqueueProbeJobs hands every record to the runner. Split out so
// startBackgroundProbe stays under the function-length budget.
func enqueueProbeJobs(runner *probe.BackgroundRunner, records []model.ScanRecord) {
	for _, rec := range records {
		runner.Start(rec)
	}
}

// drainBackgroundProbe blocks on the runner until every queued probe
// has persisted, OR — when --no-probe-wait was passed — returns
// immediately and prints a single line so users know jobs are still
// running. Safe to call with a nil runner.
func drainBackgroundProbe(runner *probe.BackgroundRunner, opts ScanProbeOptions, quiet bool) {
	if runner == nil {
		return
	}
	if opts.NoWait {
		if !quiet {
			fmt.Print(constants.MsgScanProbeDetached)
		}

		return
	}
	if !quiet {
		fmt.Printf(constants.MsgScanProbeWaitingFmt, runner.Remaining())
	}
	stats := runner.Wait()
	if !quiet {
		fmt.Printf(constants.MsgScanProbeDoneFmt,
			stats.Available, stats.Unchanged, stats.Failed)
	}
}
