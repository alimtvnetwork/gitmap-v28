package cmd

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// runReleasePull implements `gitmap release-pull` (alias `relp`).
//
// It is sugar for `release` that first runs `git pull` in the CURRENT
// repo (cwd) using the chosen mode, then delegates to runRelease with
// the remaining args. Modes (mutually exclusive):
//
//	--ff-only (default)  fast-forward only; hard-fail on divergence.
//	--rebase             rebase local commits onto upstream; abort on conflict.
//	--merge              classic merge (passes --no-rebase); creates merge commit.
func runReleasePull(args []string) {
	checkHelp(constants.CmdReleasePull, args)
	printCanonicalCmdBanner(constants.CmdReleasePull, constants.CmdReleasePullAlias)

	mode, dryRun, verbose, rest := parseReleasePullFlags(args)
	dir := requireReleasePullCwd()

	pullCurrentRepo(dir, mode, dryRun, verbose)
	forceYesOverride = true
	runRelease(ensureYesForward(rest))
}

// ensureYesForward guarantees `-y` is present in the args forwarded to
// `runRelease` when invoked via `pull-release` / `pr`. Rationale: the
// user already performed an explicit pull as part of this single
// command, so the post-release auto-commit prompt is just a stutter in
// the pipeline — we treat the whole `pr <version>` invocation as an
// implicit consent to commit + push any release-adjacent changes. If
// the caller already passed `-y` or `--yes` we leave args untouched so
// downstream flag parsing stays idempotent.
func ensureYesForward(args []string) []string {
	for _, a := range args {
		if a == "-y" || a == "--yes" || strings.HasPrefix(a, "--yes=") {
			return args
		}
	}

	return append(args, "-y")
}

// requireReleasePullCwd validates we are inside a repo and returns cwd.
func requireReleasePullCwd() string {
	if !release.IsInsideGitRepo() {
		fmt.Fprint(os.Stderr, constants.ErrRPNotInRepo)
		os.Exit(1)
	}

	cwd, err := os.Getwd()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRPCwdFailedFmt, err)
		os.Exit(1)
	}

	return cwd
}

// parseReleasePullFlags extracts release-pull-specific flags and returns
// the mode plus the remaining args to forward to `runRelease`.
func parseReleasePullFlags(args []string) (mode string, dryRun, verbose bool, rest []string) {
	fs := flag.NewFlagSet(constants.CmdReleasePull, flag.ExitOnError)
	ffOnly := fs.Bool(constants.FlagRPFFOnly, false, constants.FlagDescRPFFOnly)
	rebase := fs.Bool(constants.FlagRPRebase, false, constants.FlagDescRPRebase)
	merge := fs.Bool(constants.FlagRPMerge, false, constants.FlagDescRPMerge)
	dry := fs.Bool(constants.FlagRPDryRun, false, constants.FlagDescRPDryRun)
	vrb := fs.Bool(constants.FlagRPVerbose, false, constants.FlagDescRPVerbose)

	relevant, forwarded := splitReleasePullArgs(args)
	if err := fs.Parse(relevant); err != nil {
		os.Exit(2)
	}

	mode = resolvePullMode(*ffOnly, *rebase, *merge)
	rest = append(fs.Args(), forwarded...)

	return mode, *dry, *vrb, rest
}

// splitReleasePullArgs separates release-pull's own flags from args
// destined for the embedded `release` call. Anything matching one of
// our flag tokens goes to the local FlagSet; the rest is forwarded.
func splitReleasePullArgs(args []string) (own, forwarded []string) {
	ownFlags := map[string]bool{
		"--" + constants.FlagRPFFOnly:  true,
		"--" + constants.FlagRPRebase:  true,
		"--" + constants.FlagRPMerge:   true,
		"--" + constants.FlagRPDryRun:  true,
		"--" + constants.FlagRPVerbose: true,
	}

	for _, a := range args {
		token := a
		if eq := strings.IndexByte(a, '='); eq > 0 {
			token = a[:eq]
		}
		if ownFlags[token] {
			own = append(own, a)
		} else {
			forwarded = append(forwarded, a)
		}
	}

	return own, forwarded
}

// resolvePullMode enforces mutual exclusion and applies the default.
func resolvePullMode(ffOnly, rebase, merge bool) string {
	count := 0
	for _, b := range []bool{ffOnly, rebase, merge} {
		if b {
			count++
		}
	}

	if count > 1 {
		fmt.Fprintf(os.Stderr, constants.ErrRPModeConflictFmt, describePickedModes(ffOnly, rebase, merge))
		os.Exit(2)
	}

	if rebase {
		return constants.RPModeRebase
	}
	if merge {
		return constants.RPModeMerge
	}

	return constants.RPModeFFOnly
}

// describePickedModes builds a "--ff-only,--rebase" style summary of
// every flag the user actually set, for the conflict error message.
func describePickedModes(ffOnly, rebase, merge bool) string {
	picked := []string{}
	if ffOnly {
		picked = append(picked, "--"+constants.FlagRPFFOnly)
	}
	if rebase {
		picked = append(picked, "--"+constants.FlagRPRebase)
	}
	if merge {
		picked = append(picked, "--"+constants.FlagRPMerge)
	}

	return strings.Join(picked, ",")
}

// pullCurrentRepo runs `git pull <modeFlag>` in dir, exiting on error.
// On rebase failure it attempts `git rebase --abort` to leave the tree
// in a recoverable state before exiting.
func pullCurrentRepo(dir, mode string, dryRun, verbose bool) {
	flagArg := pullFlagForMode(mode)

	if dryRun {
		fmt.Printf(constants.MsgRPDryRunFmt, flagArg, dir)

		return
	}

	fmt.Printf(constants.MsgRPPullingFmt, flagArg, dir)

	if verbose {
		fmt.Fprintf(os.Stderr, "+ %s %s %s (cwd=%s)\n", constants.GitBin, constants.GitPull, flagArg, dir)
	}

	cmd := exec.Command(constants.GitBin, constants.GitPull, flagArg)
	cmd.Dir = dir
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		handleReleasePullFailure(dir, mode, flagArg, err)
	}
}

// pullFlagForMode maps an internal mode constant to the git pull flag.
func pullFlagForMode(mode string) string {
	switch mode {
	case constants.RPModeRebase:
		return constants.GitPullRebaseFlag
	case constants.RPModeMerge:
		return constants.GitNoRebaseFlag
	default:
		return constants.GitFFOnlyFlag
	}
}

// handleReleasePullFailure logs the failure, runs rebase --abort when
// applicable, and exits non-zero so we never tag on top of a broken
// tree.
func handleReleasePullFailure(dir, mode, flagArg string, err error) {
	if mode == constants.RPModeRebase {
		abort := exec.Command(constants.GitBin, constants.GitRebase, constants.GitRebaseAbortFlag)
		abort.Dir = dir
		_ = abort.Run()
		fmt.Fprintf(os.Stderr, constants.ErrRPRebaseAbortFmt, dir, err)
		os.Exit(1)
	}

	fmt.Fprintf(os.Stderr, constants.ErrRPPullFailedFmt, flagArg, dir, err)
	os.Exit(1)
}
