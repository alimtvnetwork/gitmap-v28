package cmd

// Post-rewrite gofmt step for `gitmap fix-repo`. The token rewriter is
// byte-level and does not understand Go's column-aligned map literals
// or const blocks, so a width-crossing bump (e.g. v9 -> v12) silently
// produces gofmt-dirty files and trips the CI gofmt gate. Running
// `gofmt -w` on every modified .go file restores deterministic
// formatting before the command exits. See
// .lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md.
//
// Windows CreateProcess caps the assembled command line at 32,767
// characters. On a large repo where hundreds of touched .go files sit
// under a long absolute path, a single-batch `gofmt -w <p1> <p2> ...`
// exec overflows that cap and Go surfaces the failure as
// "fork/exec ...gofmt.exe: The filename or extension is too long."
// invokeGofmt therefore chunks the argument list so each exec stays
// under effectiveGofmtBudget(opts), and prefers repo-relative paths
// to shrink each arg. The per-run budget is tunable via
// --gofmt-max-cmd-len (v6.80.1+).

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// gofmtArgvOverhead approximates the fixed cost of the argv prefix
// ("gofmt -w " plus the executable's resolved path) that Windows
// counts against the 32,767-char CreateProcess cap in addition to the
// file-path arguments themselves. Kept conservative on purpose.
const gofmtArgvOverhead = 32

// effectiveGofmtBudget returns the per-batch argv budget honouring
// the CLI override when set. Zero / negative user input falls back to
// the compiled-in default; the parser rejects sub-floor values before
// we ever reach here, but we guard anyway.
func effectiveGofmtBudget(opts fixRepoOptions) int {
	if opts.gofmtMaxCmdLen > 0 {
		return opts.gofmtMaxCmdLen
	}

	return constants.FixRepoGofmtMaxCmdLen
}

// runFixRepoGofmt formats every modified .go file in goFiles via
// `gofmt -w`. Returns true on success (or no-op), false if gofmt
// itself failed; the caller treats false as a write failure.
func runFixRepoGofmt(goFiles []string, opts fixRepoOptions) bool {
	if opts.isDryRun {
		return emitGofmtDryRunPreview(goFiles, opts)
	}
	if len(goFiles) == 0 {
		fmt.Print(constants.FixRepoMsgGofmtNoneFmt)

		return true
	}
	if _, err := exec.LookPath("gofmt"); err != nil {
		fmt.Fprint(os.Stderr, constants.FixRepoErrGofmtMissing)

		return true
	}

	return invokeGofmt(shortenGofmtPaths(goFiles), opts)
}

// emitGofmtDryRunPreview replaces the old single-line "skipped
// (dry-run)" message with a per-batch table so users can see cmd-line
// sizes before running the real rewrite. The tag column flags
// batches at or above 90% of the budget (NEAR-LIMIT) and any batch
// that would still overflow (OVER-LIMIT, only reachable for the
// pathological single-huge-path edge case documented on
// chunkPathsForGofmt).
func emitGofmtDryRunPreview(goFiles []string, opts fixRepoOptions) bool {
	if len(goFiles) == 0 {
		fmt.Print(constants.FixRepoMsgGofmtSkip)

		return true
	}
	budget := effectiveGofmtBudget(opts)
	paths := shortenGofmtPaths(goFiles)
	batches := chunkPathsForGofmt(paths, budget)
	fmt.Printf(constants.FixRepoMsgGofmtDryFmt, len(batches), len(paths), budget)
	for i, batch := range batches {
		cmdLen := batchCmdLen(batch)
		pct := 0
		if budget > 0 {
			pct = cmdLen * 100 / budget
		}
		tag := ""
		if pct >= 100 {
			tag = constants.FixRepoMsgGofmtDryOverTag
		} else if pct >= constants.FixRepoGofmtNearLimitPct {
			tag = constants.FixRepoMsgGofmtDryNearTag
		}
		fmt.Printf(constants.FixRepoMsgGofmtDryBatchFmt,
			i+1, len(batches), len(batch), cmdLen, pct, tag)
	}

	return true
}

// batchCmdLen approximates the exec.Command line length for a single
// gofmt batch: sum(len(path)+1) for the argv separators plus the
// fixed argv overhead ("gofmt -w " + resolved binary path headroom).
func batchCmdLen(batch []string) int {
	n := gofmtArgvOverhead
	for _, p := range batch {
		n += len(p) + 1
	}

	return n
}

// shortenGofmtPaths converts absolute paths to repo-relative form
// when possible. Shorter args reduce pressure on the Windows argv
// budget and keep the on-screen error output readable. Falls back to
// the original path when filepath.Rel fails (cross-volume on Windows,
// unreadable cwd, etc.).
func shortenGofmtPaths(paths []string) []string {
	cwd, err := os.Getwd()
	if err != nil {
		return paths
	}
	out := make([]string, len(paths))
	for i, p := range paths {
		if rel, err := filepath.Rel(cwd, p); err == nil && !strings.HasPrefix(rel, "..") {
			out[i] = rel
		} else {
			out[i] = p
		}
	}

	return out
}

// invokeGofmt is the actual `gofmt -w` call, chunked to keep each
// exec's command line under the Windows CreateProcess cap. When
// opts.isVerbose is set, prints per-batch progress with cmd-length
// and a rolling ETA computed from average per-batch wall time.
func invokeGofmt(goFiles []string, opts fixRepoOptions) bool {
	budget := effectiveGofmtBudget(opts)
	batches := chunkPathsForGofmt(goFiles, budget)
	if opts.isVerbose && len(batches) > 0 {
		fmt.Printf(constants.FixRepoMsgGofmtVerbHeaderFmt, len(batches), len(goFiles), budget)
	}
	start := time.Now()
	for i, batch := range batches {
		if opts.isVerbose {
			fmt.Printf(constants.FixRepoMsgGofmtVerbBatchStartFmt,
				i+1, len(batches), len(batch), batchCmdLen(batch))
		}
		args := append([]string{"-w"}, batch...)
		cmd := exec.Command("gofmt", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.FixRepoErrGofmtFmt, err, string(out))

			return false
		}
		if opts.isVerbose {
			elapsed := time.Since(start).Truncate(time.Millisecond)
			done := i + 1
			eta := time.Duration(int64(elapsed) / int64(done) * int64(len(batches)-done)).Truncate(time.Millisecond)
			fmt.Printf(constants.FixRepoMsgGofmtVerbBatchDoneFmt, elapsed, eta)
		}
	}
	if len(batches) > 1 {
		fmt.Printf(constants.FixRepoMsgGofmtBatchFmt, len(goFiles), len(batches))
	} else {
		fmt.Printf(constants.FixRepoMsgGofmtFmt, len(goFiles))
	}

	return true
}

// chunkPathsForGofmt splits paths into groups whose joined length
// (path bytes plus one separator per path) stays at or below budget.
// A path longer than budget is emitted in its own single-element
// chunk: gofmt itself has no argv limit, only Windows CreateProcess
// does, and a single argument that large is a pathological case the
// caller can't do better on.
func chunkPathsForGofmt(paths []string, budget int) [][]string {
	if len(paths) == 0 {
		return nil
	}
	if budget <= 0 {
		budget = constants.FixRepoGofmtMaxCmdLen
	}
	var batches [][]string
	var cur []string
	curLen := 0
	for _, p := range paths {
		// +1 accounts for the space separator between argv entries.
		cost := len(p) + 1
		if len(cur) > 0 && curLen+cost > budget {
			batches = append(batches, cur)
			cur = nil
			curLen = 0
		}
		cur = append(cur, p)
		curLen += cost
	}
	if len(cur) > 0 {
		batches = append(batches, cur)
	}

	return batches
}

// isGoSourceFile returns true when rel ends with `.go`. Kept tiny so
// the sweep loop can call it inline without obscuring the per-file
// branching in processFixRepoFile.
func isGoSourceFile(rel string) bool {
	return strings.EqualFold(filepath.Ext(rel), ".go")
}
