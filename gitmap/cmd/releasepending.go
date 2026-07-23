// Package cmd implements the CLI commands for gitmap.
package cmd

import (
	"flag"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/release"
)

// versionLikeArgPattern matches positional args the user probably meant
// as a version (`v1`, `v3.1`, `v2.233.0`, `1.4.0`). Used by
// runReleasePending to detect the classic `rp` (release-pending) vs `pr`
// (pull-release) confusion — see v5.19.0 changelog.
var versionLikeArgPattern = regexp.MustCompile(`^v?\d+(\.\d+){0,2}(-[A-Za-z0-9.]+)?$`)

// runReleasePending handles the 'release-pending' command.
func runReleasePending(args []string) {
	checkHelp("release-pending", args)
	printCanonicalCmdBanner(constants.CmdReleasePending, constants.CmdReleasePendingAlias)
	rejectVersionArgOnPending(args)
	assets, notes, draft, dryRun, verbose, noCommit, yes := parseReleasePendingFlags(args)
	_ = verbose

	err := release.ExecutePending(assets, notes, draft, dryRun, noCommit, yes)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrBareFmt, err)
		os.Exit(1)
	}
}

// printCanonicalCmdBanner prints "→ Running: gitmap <canonical> (alias: <alias>)"
// to stderr so users running an alias (e.g. `rp`) see the resolved
// canonical command name and stop confusing `rp` (release-pending) with
// `pr` (pull-release). v5.19.0+.
func printCanonicalCmdBanner(canonical, alias string) {
	fmt.Fprintf(os.Stderr, "  → Running: gitmap %s  (alias: %s)\n\n", canonical, alias)
}

// rejectVersionArgOnPending catches the classic mis-invocation
// `gitmap rp v3.1` where the user meant `gitmap pr v3.1`
// (pull-release). release-pending takes NO positional version — it
// scans for already-pending branches/metadata and releases all of
// them. Silently ignoring the version arg (the pre-v5.19.0 behavior)
// caused users to release unrelated versions (e.g. v2.233.0) when they
// asked for v3.1. Exit 2 with a precise suggestion.
func rejectVersionArgOnPending(args []string) {
	for _, a := range args {
		if strings.HasPrefix(a, "-") {
			continue
		}
		if !versionLikeArgPattern.MatchString(a) {
			continue
		}
		fmt.Fprintf(os.Stderr,
			"  ✗ gitmap release-pending (rp) takes no version argument (got %q).\n"+
				"    It releases EVERY pending branch + orphan metadata file — not a specific version.\n\n"+
				"    Did you mean:  gitmap pr %s        # pull-release: pull, then release %s\n"+
				"               or  gitmap release %s   # release %s directly\n\n",
			a, a, a, a, a)
		os.Exit(2)
	}
}

// parseReleasePendingFlags parses flags for the release-pending command.
func parseReleasePendingFlags(args []string) (assets, notes string, draft, dryRun, verbose, noCommit, yes bool) {
	fs := flag.NewFlagSet(constants.CmdReleasePending, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	notesFlag := fs.String("notes", "", constants.FlagDescNotes)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	noCommitFlag := fs.Bool("no-commit", false, constants.FlagDescNoCommit)
	yesFlag := fs.Bool("yes", false, constants.FlagDescYes)

	// Register -N as shorthand for --notes, -y as shorthand for --yes.
	fs.StringVar(notesFlag, "N", "", constants.FlagDescNotes)
	fs.BoolVar(yesFlag, "y", false, constants.FlagDescYes)

	fs.Parse(args)

	return *assetsFlag, *notesFlag, *draftFlag, *dryRunFlag, *verboseFlag, *noCommitFlag, *yesFlag
}
