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
// under FixRepoGofmtMaxCmdLen, and prefers repo-relative paths to
// shrink each arg.

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runFixRepoGofmt formats every modified .go file in goFiles via
// `gofmt -w`. Returns true on success (or no-op), false if gofmt
// itself failed; the caller treats false as a write failure.
func runFixRepoGofmt(goFiles []string, opts fixRepoOptions) bool {
	if opts.isDryRun {
		fmt.Print(constants.FixRepoMsgGofmtSkip)

		return true
	}
	if len(goFiles) == 0 {
		fmt.Print(constants.FixRepoMsgGofmtNoneFmt)

		return true
	}
	if _, err := exec.LookPath("gofmt"); err != nil {
		fmt.Fprint(os.Stderr, constants.FixRepoErrGofmtMissing)

		return true
	}

	return invokeGofmt(shortenGofmtPaths(goFiles))
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
// exec's command line under the Windows CreateProcess cap.
func invokeGofmt(goFiles []string) bool {
	batches := chunkPathsForGofmt(goFiles, constants.FixRepoGofmtMaxCmdLen)
	for _, batch := range batches {
		args := append([]string{"-w"}, batch...)
		cmd := exec.Command("gofmt", args...)
		out, err := cmd.CombinedOutput()
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.FixRepoErrGofmtFmt, err, string(out))

			return false
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
