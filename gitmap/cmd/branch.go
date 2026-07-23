// Package cmd — `gitmap branch <subcommand>` dispatcher.
package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
)

// runBranch dispatches `gitmap branch <subcommand>` (alias `b`).
//
// Today the only subcommand is `default` (alias `def`). Future
// subcommands (e.g. `current`, `list`) plug in by adding cases to the
// switch — the dispatcher pattern is intentionally explicit instead of
// table-driven because we expect <5 subcommands long-term and want
// per-subcommand help text to stay obvious.
//
// We don't gate on `--help` here because each subcommand handler does
// its own help check via checkHelp at the top.
func runBranch(args []string) {
	if len(args) == 0 {
		fmt.Fprint(os.Stderr, constants.ErrBranchMissingSubcommand)
		os.Exit(1)
	}
	sub, rest := args[0], args[1:]
	switch sub {
	case constants.CmdBranchSubDefault, constants.CmdBranchSubDefaultAlias:
		runBranchDefault(rest)

		return
	}
	fmt.Fprintf(os.Stderr, constants.ErrBranchUnknownSubcommand, sub)
	os.Exit(1)
}

// runBranchDefault implements `gitmap branch default` / `b def`.
//
// Resolution flow:
//  1. Verify cwd is inside a git work tree (reuses gitutil.IsInsideWorkTree
//     for parity with `lb`).
//  2. Resolve the default branch name via gitutil.ResolveDefaultBranchName,
//     which prefers `git symbolic-ref refs/remotes/origin/HEAD` and falls
//     back to constants.DefaultBranch ("main") so even brand-new repos
//     without an origin still get a sane target.
//  3. Run `git checkout` and surface git's status line verbatim — same
//     UX as `lb --switch` so users get one mental model for both flows.
//
// Errors are printed with the standardized "  ✗ ..." prefix and exit 1
// to short-circuit shell pipelines.
func runBranchDefault(args []string) {
	checkHelp("branch", args)
	if !gitutil.IsInsideWorkTree() {
		fmt.Fprint(os.Stderr, constants.ErrBranchNotRepo)
		os.Exit(1)
	}
	target := gitutil.ResolveDefaultBranchName(".")
	fmt.Printf(constants.MsgBranchDefaultSwitching, target)
	out, err := gitutil.CheckoutBranch(".", target)
	if len(out) > 0 {
		fmt.Println(out)
	}
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBranchDefaultFailed, target, err)
		os.Exit(1)
	}
}
