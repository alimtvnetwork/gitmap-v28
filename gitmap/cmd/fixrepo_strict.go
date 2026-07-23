package cmd

// fixrepo_strict.go — orchestrates the optional post-rewrite `go test`
// step gated by --strict / -Strict. Mirrors the gofmt step's shape
// (dry-run skip → no-files skip → toolchain-missing warn → invoke)
// so the trailing summary block reads consistently across post-steps.
//
// Why this exists: the rewriter is byte-level and cannot detect
// SEMANTIC desyncs — e.g. the v9→v10/v12 width-crossing failure
// closed by v4.12.0, where a hard-coded sibling literal (`"9"`) drifted
// from its {base}-v9 neighbor. `go test` on the touched packages
// catches such drift immediately, BEFORE the bumped commit lands.
// Off by default so non-Go repos and machines without a Go toolchain
// stay unaffected; opt-in for Go repos in CI.
//
// Failure mode contract: prints the captured `go test` output to
// stderr, returns false → caller exits with FixRepoExitTestsFailed (9).

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runFixRepoStrict is the public entry called from runFixRepo after
// the gofmt step. Returns true on success (or skip), false on test
// failure. Never returns false for "go missing" or "no files" — those
// are considered safe skips so --strict can be left on across mixed
// environments (Linux CI with Go installed; macOS dev machine without).
func runFixRepoStrict(repoRoot string, goFiles []string, opts fixRepoOptions) bool {
	if !opts.isStrict {
		return true
	}
	if opts.isDryRun {
		fmt.Print(constants.FixRepoMsgStrictSkipDryRun)

		return true
	}
	if len(goFiles) == 0 {
		fmt.Print(constants.FixRepoMsgStrictNoGoFiles)

		return true
	}
	packages := derivePackagesFromGoFiles(repoRoot, goFiles)
	if len(packages) == 0 {
		fmt.Print(constants.FixRepoMsgStrictNoPackages)

		return true
	}
	if _, err := exec.LookPath("go"); err != nil {
		fmt.Fprint(os.Stderr, constants.FixRepoErrStrictMissing)

		return true
	}

	return invokeGoTest(repoRoot, packages)
}

// invokeGoTest runs `go test <pkg1> <pkg2> ...` from repoRoot and
// returns true on success. Output is captured (CombinedOutput) so the
// caller can include it verbatim in the failure message — `go test`
// writes the actual test failures to stdout, the compile errors to
// stderr, and CI users routinely lose one half when only one stream
// is captured. Batching all packages into one invocation is faster
// than per-package calls (one process startup, one module-graph load)
// AND keeps the output a single coherent block.
func invokeGoTest(repoRoot string, packages []string) bool {
	fmt.Printf(constants.FixRepoMsgStrictRunFmt,
		len(packages), strings.Join(packages, " "))
	args := append([]string{"test"}, packages...)
	cmd := exec.Command("go", args...)
	cmd.Dir = repoRoot
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.FixRepoErrStrictFailFmt, err, string(out))

		return false
	}
	fmt.Printf(constants.FixRepoMsgStrictPassFmt, len(packages))

	return true
}
