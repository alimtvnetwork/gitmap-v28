// Package cmd — visibility.go: entry points for `gitmap make-public`
// and `gitmap make-private`.
//
// These commands toggle the current repository's visibility on the
// remote provider (GitHub or GitLab). They wrap the host CLI (`gh`
// or `glab`) so we don't have to ship OAuth tokens — if the CLI is
// authenticated, so are we.
//
// Spec parity: spec-authoring/23-visibility-change/01-spec.md.
//
// Forms:
//
//	gitmap make-public  [--yes] [--dry-run] [--verbose]
//	gitmap make-private        [--dry-run] [--verbose]
//
// `--yes` is a no-op for `make-private` (no confirmation is shown
// when going public → private; the asymmetry matches the PowerShell
// reference and is intentional — exposing a private repo is the
// risky direction, hiding a public one is reversible).
package cmdvisibility

import (
	"flag"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runMakePublic implements `gitmap make-public`.
func RunMakePublic(args []string) error {
	checkHelp(constants.CmdMakePublic, args)
	runVisibility(args, constants.VisibilityPublic)

	return nil
}

// runMakePrivate implements `gitmap make-private`.
func RunMakePrivate(args []string) error {
	checkHelp(constants.CmdMakePrivate, args)
	runVisibility(args, constants.VisibilityPrivate)

	return nil
}

// VisibilityFlags captures parsed flag state. Kept as a struct so
// runVisibility stays under the 15-line limit.
type VisibilityFlags struct {
	yes     bool
	dryRun  bool
	verbose bool
}

// runVisibility is the shared core for both commands. Steps mirror
// the PowerShell reference verbatim so behavior parity is auditable.
func runVisibility(args []string, target string) error {
	opts, positional := parseVisibilityFlags(args, target)

	// Spec 113 §2.2 — bulk form takes over when positional args
	// describe a `<repo> <count>` or `<count>` request.
	if maybeRunBulkVisibility(positional, target, opts) {
		return nil
	}

	ctx := mustResolveVisibilityContext()
	mustEnsureProviderCLI(ctx.Provider, opts.verbose)

	current := mustReadCurrentVisibility(ctx, opts.verbose)
	if current == target {
		fmt.Printf(constants.MsgVisAlreadyFmt, current, ctx.Provider)
		cliexit.HandleError(nil, constants.ExitVisOK)
	}

	if target == constants.VisibilityPublic && !opts.yes {
		confirmPublicOrExit(ctx)
	}

	if opts.dryRun {
		fmt.Printf(constants.MsgVisDryRunFmt, current, target, ctx.Slug, ctx.Provider)
		cliexit.HandleError(nil, constants.ExitVisOK)
	}

	applyVisibilityOrExit(ctx, target, opts.verbose)
	verifyVisibilityOrExit(ctx, target, opts.verbose)
	fmt.Printf(constants.MsgVisChangedFmt, current, target, ctx.Slug, ctx.Provider)

	return nil
}

// parseVisibilityFlags reads the supported flags and returns the
// leftover positional args (spec 113 §2.2 bulk form). cmdName is
// passed in only so FlagSet error output names the right command.
func parseVisibilityFlags(args []string, target string) (VisibilityFlags, []string) {
	cmdName := constants.CmdMakePublic
	if target == constants.VisibilityPrivate {
		cmdName = constants.CmdMakePrivate
	}

	fs := flag.NewFlagSet(cmdName, flag.ExitOnError)
	yesLong := fs.Bool(constants.FlagVisYes, false, constants.FlagDescVisYes)
	yesShort := fs.Bool(constants.FlagVisYesAlt, false, constants.FlagDescVisYes)
	dry := fs.Bool(constants.FlagVisDryRun, false, constants.FlagDescVisDryRun)
	vrb := fs.Bool(constants.FlagVisVerbose, false, constants.FlagDescVisVerbose)

	if err := fs.Parse(args); err != nil {
		cliexit.HandleError(err, constants.ExitVisBadFlag)
	}

	out := VisibilityFlags{yes: *yesLong || *yesShort, dryRun: *dry, verbose: *vrb}

	return out, fs.Args()
}
