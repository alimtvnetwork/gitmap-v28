// Package cmdrelease implements the CLI commands for gitmap.
package cmdrelease

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/lazyregex"
	"github.com/alimtvnetwork/gitmap-v28/cli/release"
)

// versionLikeArgPattern matches positional args the user probably meant
// as a version (`v1`, `v3.1`, `v2.233.0`, `1.4.0`). Used by
// runReleasePending to detect the classic `rp` (release-pending) vs `pr`
// (pull-release) confusion — see v5.19.0 changelog.
var versionLikeArgPattern = lazyregex.New(`^v?\d+(\.\d+){0,2}(-[A-Za-z0-9.]+)?$`)

// RunReleasePending handles the 'release-pending' command.
func RunReleasePending(args []string) error {
	checkHelp("release-pending", args)
	if isRunningProjectsAliasArg(args) {
		return cmdagy.RunRunningProjectsCLI(args)
	}

	printCanonicalCmdBanner(constants.CmdReleasePending, constants.CmdReleasePendingAlias)

	if err := rejectVersionArgOnPending(args); err != nil {
		return err
	}

	return executeReleasePendingRun(args)
}

func executeReleasePendingRun(args []string) error {
	assets, notes, draft, dryRun, verbose, noCommit, yes := parseReleasePendingFlags(args)
	_ = verbose

	err := release.ExecutePending(assets, notes, draft, dryRun, noCommit, yes)
	if err != nil {
		return apperror.WrapSimple(err, constants.ErrBareFmt)
	}

	return nil
}

func isRunningProjectsAliasArg(args []string) bool {
	if len(args) == 0 {
		return false
	}

	first := strings.ToLower(args[0])

	return first == "ls" || first == "list"
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
// (pull-release). release-pending takes NO positional version.
func rejectVersionArgOnPending(args []string) error {
	for _, a := range args {
		if strings.HasPrefix(a, "-") || !versionLikeArgPattern.MatchString(a) {
			continue
		}

		return apperror.NewValidationError(fmt.Sprintf(
			"gitmap release-pending takes no version argument (got %q); did you mean: gitmap pr %s or gitmap release %s",
			a, a, a,
		))
	}

	return nil
}

// parseReleasePendingFlags parses flags for the release-pending command.
func parseReleasePendingFlags(args []string) (string, string, bool, bool, bool, bool, bool) {
	fs, assetsFlag, notesFlag, draftFlag, dryRunFlag, verboseFlag, noCommitFlag, yesFlag := initReleasePendingFlags()
	fs.Parse(args)

	return *assetsFlag, *notesFlag, *draftFlag, *dryRunFlag, *verboseFlag, *noCommitFlag, *yesFlag
}

func initReleasePendingFlags() (*flag.FlagSet, *string, *string, *bool, *bool, *bool, *bool, *bool) {
	fs := flag.NewFlagSet(constants.CmdReleasePending, flag.ExitOnError)
	assetsFlag := fs.String("assets", "", constants.FlagDescAssets)
	notesFlag := fs.String("notes", "", constants.FlagDescNotes)
	draftFlag := fs.Bool("draft", false, constants.FlagDescDraft)
	dryRunFlag := fs.Bool("dry-run", false, constants.FlagDescDryRun)
	verboseFlag := fs.Bool("verbose", false, constants.FlagDescVerbose)
	noCommitFlag := fs.Bool("no-commit", false, constants.FlagDescNoCommit)
	yesFlag := fs.Bool("yes", false, constants.FlagDescYes)
	fs.StringVar(notesFlag, "N", "", constants.FlagDescNotes)
	fs.BoolVar(yesFlag, "y", false, constants.FlagDescYes)

	return fs, assetsFlag, notesFlag, draftFlag, dryRunFlag, verboseFlag, noCommitFlag, yesFlag
}
