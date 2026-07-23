// Package cmd — visibilityparallel.go: bounded worker pool that
// applies the per-repo visibility flip concurrently. Per-repo stdout
// is captured into a bytes.Buffer per worker and flushed atomically
// under a mutex so the interleaved output stays line-coherent.
//
// Audit writes (a.updateResult) are serialized through the same
// mutex because SQLite connections aren't safe for concurrent writers.
//
// Spec: spec/01-app/116-bulk-visibility-mapub-mapri.md §parallel.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/visibility"
)

// applyBulkLoopParallel fans the matched set across `flags.Parallel`
// workers and tallies the results. When Parallel<=1 it degrades to
// the sequential path so the per-line live output is preserved.
func applyBulkLoopParallel(ctx ownerContext, target string, matches []visibility.MatchedRepo, flags bulkFlags, audit *runAudit) (int, int, int) {
	workers := flags.Parallel
	if workers <= 1 {
		return applyBulkLoopSeq(ctx, target, matches, flags.Verbose, audit)
	}
	fmt.Fprintf(os.Stdout, constants.MsgBulkParallelFmt, workers)

	total := len(matches)
	type job struct {
		i int
		m visibility.MatchedRepo
	}
	jobs := make(chan job, total)
	for i, m := range matches {
		jobs <- job{i: i, m: m}
	}
	close(jobs)

	var (
		mu                       sync.Mutex
		changed, skipped, failed int
		wg                       sync.WaitGroup
	)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := range jobs {
				buf := &bytes.Buffer{}
				fmt.Fprintf(buf, constants.MsgBulkApplyItemFmt, j.i+1, total, j.m.RepoName)
				start := time.Now()
				status := applyOneRepoTo(buf, ctx, j.m.RepoName, target, flags.Verbose)
				mu.Lock()
				os.Stdout.Write(buf.Bytes())
				audit.updateResult(j.m.RepoName, status, status.prev, status.next, start)
				switch status.outcome {
				case "skip":
					skipped++
				case "ok":
					changed++
				default:
					failed++
				}
				mu.Unlock()
			}
		}()
	}
	wg.Wait()

	return changed, skipped, failed
}

// applyBulkLoopSeq is the original sequential implementation kept for
// --parallel=1 and tests that assert deterministic per-line order.
func applyBulkLoopSeq(ctx ownerContext, target string, matches []visibility.MatchedRepo, verbose bool, audit *runAudit) (int, int, int) {
	changed, skipped, failed := 0, 0, 0
	total := len(matches)
	for i, m := range matches {
		fmt.Fprintf(os.Stdout, constants.MsgBulkApplyItemFmt, i+1, total, m.RepoName)
		start := time.Now()
		status := applyOneRepo(ctx, m.RepoName, target, verbose)
		audit.updateResult(m.RepoName, status, status.prev, status.next, start)
		switch status.outcome {
		case "skip":
			skipped++
		case "ok":
			changed++
		default:
			failed++
		}
	}

	return changed, skipped, failed
}
