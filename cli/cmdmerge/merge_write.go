package cmdmerge

// merge_write.go — `gitmap merge write` / `gitmap mr`: open a GitHub pull
// request from the CLI (program 273 WS3, spec 05). Flag parsing and flow
// orchestration live here; the GitHub host call lives in merge_write_gh.go;
// help text + the help-example registry live in merge_write_help.go.

import (
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// merge-write subcommand token, flag names, and user-facing messages.
const (
	mergeWriteName = "merge write"
	mergeWriteSub  = "write"
	mwFlagTitle    = "title"
	mwFlagBody     = "body"
	mwFlagSource   = "source"
	mwFlagTarget   = "target"
	mwFlagDraft    = "draft"
	mwFlagYes      = "yes"
	mwFlagYesShort = "y"

	mwUsageLine  = "Usage: gitmap merge write --title <title> [--body <body>] [--source <branch>] [--target <branch>] [--draft] [--yes]\n       gitmap mr --title <title> [same flags]\n"
	mwErrTitle   = "missing required flag: --title (every pull request needs a title)"
	mwSummaryFmt = "[merge write] opening pull request\n  repo:   %s\n  %s -> %s%s\n  title:  %s\n"
	mwConfirmCue = "Proceed? [y/N] "
	mwAborted    = "merge write: aborted by user"
	mwResultFmt  = "pr #%d opened: %s -> %s%s\n%s\n"
)

// mergeWriteOpts carries the parsed flags for `merge write`.
type mergeWriteOpts struct {
	title   string
	body    string
	source  string
	target  string
	isDraft bool
	isYes   bool
}

// isMergeWriteInvocation reports whether argv routes to `merge write`
// (the `mr` alias, or `merge` followed by the `write` subcommand).
func isMergeWriteInvocation(command string, argv []string) bool {
	if command == constants.CmdMergeWriteA {
		return true
	}
	if command != constants.CmdMerge {
		return false
	}
	return len(argv) > 2 && argv[2] == mergeWriteSub
}

// stripMergeWriteSubcommand drops the `write` token (`merge write ...`
// arrives as argv[2:]=["write", ...]; the `mr` alias has no subcommand).
func stripMergeWriteSubcommand(args []string) []string {
	if len(args) > 0 && args[0] == mergeWriteSub {
		return args[1:]
	}
	return args
}

// runMergeWrite executes `gitmap merge write` / `gitmap mr`.
func runMergeWrite(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(mwHelpText)
		cliexit.Exit(0)
	}
	opts, appErr := parseMergeWriteFlags(stripMergeWriteSubcommand(args))
	if appErr != nil {
		return appErr
	}
	return executeMergeWrite(opts)
}

// executeMergeWrite runs the flow: preflight, resolve, confirm, create.
func executeMergeWrite(opts mergeWriteOpts) *appfault.AppError {
	slug, resolved, appErr := prepareMergeWrite(opts)
	if appErr != nil {
		return appErr
	}
	if !isMergeWriteConfirmed(slug, resolved) {
		fmt.Println(mwAborted)
		return nil
	}
	return createAndReportMergeWrite(slug, resolved)
}

// prepareMergeWrite runs preflight checks and fills branch defaults.
func prepareMergeWrite(opts mergeWriteOpts) (string, mergeWriteOpts, *appfault.AppError) {
	if appErr := requireMergeWriteGh(); appErr != nil {
		return "", opts, appErr
	}
	slug, appErr := resolveGitHubSlug()
	if appErr != nil {
		return "", opts, appErr
	}
	resolved, appErr := resolveMergeWriteBranches(opts, slug)
	if appErr != nil {
		return "", opts, appErr
	}
	return slug, resolved, nil
}

// createAndReportMergeWrite opens the PR and prints the result lines.
func createAndReportMergeWrite(slug string, opts mergeWriteOpts) *appfault.AppError {
	result, appErr := openGitHubPullRequest(slug, opts)
	if appErr != nil {
		return appErr
	}
	printMergeWriteResult(result, opts)
	return nil
}

// parseMergeWriteFlags parses the `merge write` flag set.
func parseMergeWriteFlags(args []string) (mergeWriteOpts, *appfault.AppError) {
	opts := mergeWriteOpts{}
	fs := flag.NewFlagSet(mergeWriteName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	bindMergeWriteFlags(fs, &opts)
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(mwUsageLine + "error: " + err.Error())
	}
	if strings.TrimSpace(opts.title) == "" {
		return opts, appfault.NewValidationError(mwErrTitle)
	}
	return opts, nil
}

// bindMergeWriteFlags registers the `merge write` flag set.
func bindMergeWriteFlags(fs *flag.FlagSet, opts *mergeWriteOpts) {
	fs.StringVar(&opts.title, mwFlagTitle, "", "PR title (required)")
	fs.StringVar(&opts.body, mwFlagBody, "", "PR body (markdown allowed)")
	fs.StringVar(&opts.source, mwFlagSource, "", "head branch (default: current branch)")
	fs.StringVar(&opts.target, mwFlagTarget, "", "base branch (default: repo default branch)")
	fs.BoolVar(&opts.isDraft, mwFlagDraft, false, "open as a draft PR")
	fs.BoolVar(&opts.isYes, mwFlagYes, false, "skip the confirmation prompt")
	fs.BoolVar(&opts.isYes, mwFlagYesShort, false, "skip the confirmation prompt (shorthand)")
}

// isMergeWriteConfirmed prints the summary and prompts unless --yes.
func isMergeWriteConfirmed(slug string, opts mergeWriteOpts) bool {
	draft := ""
	if opts.isDraft {
		draft = " (draft)"
	}
	fmt.Printf(mwSummaryFmt, slug, opts.source, opts.target, draft, opts.title)
	if opts.isYes {
		return true
	}
	fmt.Print(mwConfirmCue)
	return readMergeWriteAnswer()
}

// readMergeWriteAnswer reads a y/n answer from stdin.
func readMergeWriteAnswer() bool {
	var line string
	_, err := fmt.Scanln(&line)
	if err != nil {
		return false
	}
	answer := strings.ToLower(strings.TrimSpace(line))
	return answer == "y" || answer == "yes"
}
