package cmd

// Post-rewrite gofmt step for `gitmap fix-repo`. The token rewriter is
// byte-level and does not understand Go's column-aligned map literals
// or const blocks, so a width-crossing bump (e.g. v9 → v12) silently
// produces gofmt-dirty files and trips the CI gofmt gate. Running
// `gofmt -w` on every modified .go file restores deterministic
// formatting before the command exits. See
// .lovable/memory/issues/2026-05-01-fixrepo-no-gofmt.md.

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

	return invokeGofmt(goFiles)
}

// invokeGofmt is the actual `gofmt -w` call, batched across all paths
// in one process to keep latency low even on repos with hundreds of
// rewritten Go files.
func invokeGofmt(goFiles []string) bool {
	args := append([]string{"-w"}, goFiles...)
	cmd := exec.Command("gofmt", args...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrGofmtFmt, err, string(out))

		return false
	}
	fmt.Printf(constants.FixRepoMsgGofmtFmt, len(goFiles))

	return true
}

// isGoSourceFile returns true when rel ends with `.go`. Kept tiny so
// the sweep loop can call it inline without obscuring the per-file
// branching in processFixRepoFile.
func isGoSourceFile(rel string) bool {
	return strings.EqualFold(filepath.Ext(rel), ".go")
}
