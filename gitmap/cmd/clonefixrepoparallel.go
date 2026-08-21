// Package cmd — clonefixrepoparallel.go: comma-separated URL
// fan-out for `gitmap cfr` / `gitmap cfrp`. Mirrors the
// `visibilityparallel.go` pattern (mapub/mapri): bounded worker
// pool, per-worker bytes.Buffer captured stdout, mutex-guarded
// atomic flush so interleaved lines stay coherent.
//
// Each worker re-execs the current binary with a single URL so the
// existing single-URL pipeline (chdir, fix-repo chaining, exit
// codes, transport persistence) stays the source of truth. Trade:
// each clone forks one extra process — the network/IO dwarfs it,
// and isolation is worth it.
//
// Result ORDER matches input order: per-URL output buffers are
// flushed under the mutex in the same goroutine that finished work,
// so terminal interleaving is line-coherent but URL order is
// effectively "completion order". A pre-flight banner lists the URL
// set and worker count so the user can correlate finished blocks
// back to inputs.
package cmd

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// splitCommaURLs splits a positional URL token on commas. Whitespace
// is trimmed; empty fragments are dropped. Returns a single-element
// slice for non-comma input so callers can treat the result
// uniformly.
func splitCommaURLs(raw string) []string {
	if !strings.Contains(raw, ",") {
		return handleSingleURLSplit(raw)
	}
	parts := strings.Split(raw, ",")
	out := make([]string, 0, len(parts))
	for _, p := range parts {
		t := strings.TrimSpace(p)
		if t != "" {
			out = append(out, t)
		}
	}
	return out
}

type cfrJob struct {
	i   int
	url string
}

func initCFRJobQueue(urls []string) chan cfrJob {
	jobs := make(chan cfrJob, len(urls))
	for i, u := range urls {
		jobs <- cfrJob{i: i, url: u}
	}
	close(jobs)
	return jobs
}

func resolveCFRWorkers(workers, total int) int {
	if workers <= 0 {
		workers = constants.CloneFixRepoDefaultParallel
	}
	if workers > total {
		return total
	}
	return workers
}

func startCFRWorker(bin, subcmd string, total int, leadingMods, passthroughFlags []string, jobs <-chan cfrJob, mu *sync.Mutex, failed *int, wg *sync.WaitGroup) {
	defer wg.Done()
	for j := range jobs {
		if !runOneCFRJob(bin, subcmd, j.url, j.i+1, total, leadingMods, passthroughFlags, mu) {
			mu.Lock()
			*failed++
			mu.Unlock()
		}
	}
}

func runCFRWorkerPool(bin, subcmd string, urls []string, leadingMods, passthroughFlags []string, workers int) int {
	jobs := initCFRJobQueue(urls)
	var (
		mu     sync.Mutex
		wg     sync.WaitGroup
		failed int
	)
	for w := 0; w < workers; w++ {
		wg.Add(1)
		go startCFRWorker(bin, subcmd, len(urls), leadingMods, passthroughFlags, jobs, &mu, &failed, &wg)
	}
	wg.Wait()
	return failed
}

// runCloneFixRepoParallel fans `urls` across `workers` goroutines.
// Each worker re-execs `gitmap <subcmd> <url> [passthroughFlags...]`.
// passthroughFlags carries the user's original non-positional flags
// (--ssh, --no-vscode-sync, --require-version, --dry-run, ...) so
// each worker observes the same semantics.
//
// Returns the count of failed URLs (exit-non-zero from the re-exec).
func runCloneFixRepoParallel(urls []string, subcmd string, leadingMods, passthroughFlags []string, workers int) int {
	workers = resolveCFRWorkers(workers, len(urls))
	bin, err := os.Executable()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrCloneFixRepoExecFmt, err)
		return len(urls)
	}
	fmt.Fprintf(os.Stdout, constants.MsgCloneFixRepoParallelHeader, len(urls), workers, subcmd)
	failed := runCFRWorkerPool(bin, subcmd, urls, leadingMods, passthroughFlags, workers)
	printCloneFixRepoParallelSummary(len(urls), failed)
	return failed
}

func buildCFRJobArgs(subcmd, url string, leadingMods, passthroughFlags []string) []string {
	args := make([]string, 0, 2+len(leadingMods)+len(passthroughFlags))
	args = append(args, subcmd)
	args = append(args, leadingMods...)
	args = append(args, url)
	return append(args, passthroughFlags...)
}

func logCFRJobResult(buf *bytes.Buffer, idx, total int, url string, elapsed time.Duration, runErr error) {
	if runErr == nil {
		fmt.Fprintf(buf, constants.MsgCloneFixRepoParallelItemOk, idx, total, url, elapsed)
		return
	}
	fmt.Fprintf(buf, constants.MsgCloneFixRepoParallelItemFail, idx, total, url, elapsed, runErr)
}

func execCFRJobCommand(bin string, args []string, buf *bytes.Buffer) (time.Duration, error) {
	start := time.Now()
	cmd := exec.Command(bin, args...)
	cmd.Stdout = buf
	cmd.Stderr = buf
	runErr := cmd.Run()
	return time.Since(start).Round(time.Millisecond), runErr
}

// runOneCFRJob re-execs the binary with a single URL, captures both
// streams into one buffer, and flushes under the shared mutex so
// per-URL blocks stay contiguous. Returns true on exit code 0.
func runOneCFRJob(bin, subcmd, url string, idx, total int, leadingMods, passthroughFlags []string, mu *sync.Mutex) bool {
	args := buildCFRJobArgs(subcmd, url, leadingMods, passthroughFlags)
	buf := &bytes.Buffer{}
	fmt.Fprintf(buf, constants.MsgCloneFixRepoParallelItem, idx, total, url)
	elapsed, runErr := execCFRJobCommand(bin, args, buf)
	logCFRJobResult(buf, idx, total, url, elapsed, runErr)
	mu.Lock()
	_, _ = os.Stdout.Write(buf.Bytes())
	mu.Unlock()
	return runErr == nil
}

func printCloneFixRepoParallelSummary(total, failed int) {
	if failed == 0 {
		fmt.Fprintf(os.Stdout, constants.MsgCloneFixRepoParallelDoneOk, total)
		return
	}
	fmt.Fprintf(os.Stdout, constants.MsgCloneFixRepoParallelDoneFail, total-failed, failed)
}

func isParallelFlagName(name string) bool {
	return name == "--parallel" || name == "-parallel" || name == "--p" || name == "-p"
}

func parseParallelArgStep(args []string, i int) (int, int, bool) {
	name, val, hasVal := splitFlagEq(args[i])
	if !isParallelFlagName(name) {
		return 0, 0, false
	}
	if hasVal {
		return atoiSafe(val), 0, true
	}
	if i+1 < len(args) {
		return atoiSafe(args[i+1]), 1, true
	}
	return 0, 0, true
}

// extractParallelFlag scans args for --parallel=N / --parallel N /
// -p=N / -p N. Returns the resolved value (0 = unset → caller picks
// default) and the args slice with the flag removed so the residual
// flags can be forwarded to single-URL workers without duplication.
func extractParallelFlag(args []string) (int, []string) {
	out := make([]string, 0, len(args))
	parallel := 0
	for i := 0; i < len(args); i++ {
		if pVal, advance, ok := parseParallelArgStep(args, i); ok {
			if pVal != 0 {
				parallel = pVal
			}
			i += advance
			continue
		}
		out = append(out, args[i])
	}
	return parallel, out
}

// splitFlagEq splits "--name=value" into ("--name", "value", true).
// Returns (raw, "", false) when no `=` is present.
func splitFlagEq(a string) (string, string, bool) {
	idx := strings.IndexByte(a, '=')
	if idx < 0 {
		return a, "", false
	}
	return a[:idx], a[idx+1:], true
}

func atoiSafe(s string) int {
	n := 0
	for _, c := range s {
		if c < '0' || c > '9' {
			return 0
		}
		n = n*10 + int(c-'0')
	}
	return n
}

func buildCFRTransportAndSyncFlags(noVSCodeSync, requireVersion, useSSH, useHTTPS bool) []string {
	out := make([]string, 0, 4)
	if noVSCodeSync {
		out = append(out, "--"+constants.FlagNoVSCodeSync)
	}
	if requireVersion {
		out = append(out, "--"+constants.FlagRequireVersion)
	}
	if useSSH {
		out = append(out, "--ssh")
	}
	if useHTTPS {
		out = append(out, "--https")
	}
	return out
}

func buildCFRActionFlags(autoYes, dryRun, noCommit, noPush bool) []string {
	out := make([]string, 0, 4)
	if autoYes {
		out = append(out, "--yes")
	}
	if dryRun {
		out = append(out, "--"+constants.FlagCloneDryRun)
	}
	if noCommit {
		out = append(out, "--"+constants.FlagCGNoCommit)
	}
	if noPush {
		out = append(out, "--"+constants.FlagCGNoPush)
	}
	return out
}

// buildCFRPassthroughFlags reconstructs the non-positional flag set
// for re-exec workers. Mirrors the flags recognised by
// parseCloneFixRepoArgs so each worker sees the same semantics. The
// --parallel flag is intentionally NOT forwarded — workers run a
// single URL and must not recurse into another fan-out.
func buildCFRPassthroughFlags(noVSCodeSync, requireVersion, useSSH, useHTTPS, autoYes, dryRun, noCommit, noPush bool) []string {
	flags := buildCFRTransportAndSyncFlags(noVSCodeSync, requireVersion, useSSH, useHTTPS)
	return append(flags, buildCFRActionFlags(autoYes, dryRun, noCommit, noPush)...)
}

func handleSingleURLSplit(raw string) []string {
	trimmed := strings.TrimSpace(raw)
	if trimmed == "" {
		return nil
	}
	return []string{trimmed}
}
