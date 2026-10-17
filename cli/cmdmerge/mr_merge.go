package cmdmerge

// mr_merge.go — `gitmap mr merge <n>` / `gitmap mrm`: merge a pull request
// via the `gh` CLI (program 275, proposal P-03). Side-effect command: it
// requires --yes or an interactive confirmation, and supports --dry-run
// (shared-patterns spec §5). One concern per file (≤300 lines).

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// mr merge flag names, messages, and merge strategies.
const (
	mrMergeName          = "mr merge"
	mrMergeFlagStrategy  = "strategy"
	mrMergeFlagSquash    = "squash"
	mrMergeFlagMerge     = "merge"
	mrMergeFlagRebase    = "rebase"
	mrMergeFlagDelBranch = "delete-branch"
	mrMergeFlagRepo      = "repo"
	mrMergeFlagYes       = "yes"
	mrMergeFlagYesShort  = "y"
	mrMergeFlagDryRun    = "dry-run"
	mrMergeDefStrategy   = "squash"

	mrMergeUsageLine = "Usage: gitmap mr merge <n> [--strategy squash|merge|rebase] [--delete-branch] [--repo <owner/repo>] [--yes] [--dry-run]\n       gitmap mrm <n> [same flags]\n"
	mrMergeErrNum    = "missing required argument: <n> (pull request number, e.g. `gitmap mr merge 42`)"
	mrMergeErrNumFmt = "invalid pull request number %q: want a positive integer"
	mrMergeErrStrat  = "invalid --strategy %q: want one of squash|merge|rebase"

	mrMergeSummaryFmt = "[mr merge] merging pull request #%d\n  repo:     %s\n  strategy: %s%s\n"
	mrMergeConfirmCue = "Proceed? [y/N] "
	mrMergeAborted    = "mr merge: aborted by user"
	mrMergeDryRunFmt  = "[mr merge] dry-run — no changes made\n  would run: gh pr merge %d --repo %s --%s%s\n"
	mrMergeResultFmt  = "merged #%d (%s) -> %s%s\n"
)

// mrMergeOpts carries the parsed flags for `mr merge`.
type mrMergeOpts struct {
	number         int
	strategy       string
	isDeleteBranch bool
	isYes          bool
	isDryRun       bool
	repo           string
}

// mrMergeCommit decodes gh's mergeCommit object.
type mrMergeCommit struct {
	Oid string `json:"oid"`
}

// RunMrMerge implements `gitmap mr merge` / `gitmap mrm` for rootdata wiring.
func RunMrMerge(args []string) error {
	if appErr := runMrMerge(args); appErr != nil {
		return appErr
	}
	return nil
}

// runMrMerge executes the merge flow: preflight, confirm, merge, report.
func runMrMerge(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(mrMergeHelpText)
		cliexit.Exit(0)
	}
	opts, appErr := parseMrMergeFlags(args)
	if appErr != nil {
		return appErr
	}
	if appErr := requireMrGh(); appErr != nil {
		return appErr
	}
	slug, appErr := resolveMrRepo(opts.repo)
	if appErr != nil {
		return appErr
	}
	if opts.isDryRun {
		printMrMergeDryRun(slug, opts)
		return nil
	}
	if !isMrMergeConfirmed(slug, opts) {
		fmt.Println(mrMergeAborted)
		return nil
	}
	return executeMrMerge(slug, opts)
}

// parseMrMergeFlags parses the `mr merge` flag set + the required <n>.
func parseMrMergeFlags(args []string) (mrMergeOpts, *appfault.AppError) {
	opts := mrMergeOpts{strategy: mrMergeDefStrategy}
	var isMergeFlag, isRebaseFlag, isSquashFlag bool
	fs := flag.NewFlagSet(mrMergeName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.strategy, mrMergeFlagStrategy, mrMergeDefStrategy, "merge strategy: squash|merge|rebase")
	fs.BoolVar(&isSquashFlag, mrMergeFlagSquash, false, "squash-merge (default strategy)")
	fs.BoolVar(&isMergeFlag, mrMergeFlagMerge, false, "regular merge commit")
	fs.BoolVar(&isRebaseFlag, mrMergeFlagRebase, false, "rebase-merge")
	fs.BoolVar(&opts.isDeleteBranch, mrMergeFlagDelBranch, false, "delete the head branch after merging")
	fs.BoolVar(&opts.isYes, mrMergeFlagYes, false, "skip the confirmation prompt")
	fs.BoolVar(&opts.isYes, mrMergeFlagYesShort, false, "skip the confirmation prompt (shorthand)")
	fs.BoolVar(&opts.isDryRun, mrMergeFlagDryRun, false, "print what would happen, change nothing")
	fs.StringVar(&opts.repo, mrMergeFlagRepo, "", "owner/repo (default: the current repo's origin)")
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(mrMergeUsageLine + "error: " + err.Error())
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Print(mrMergeUsageLine)
		return opts, appfault.NewValidationError(mrMergeErrNum)
	}
	n, err := strconv.Atoi(rest[0])
	if err != nil || n <= 0 {
		return opts, appfault.NewValidationError(fmt.Sprintf(mrMergeErrNumFmt, rest[0]))
	}
	opts.number = n
	switch {
	case isMergeFlag:
		opts.strategy = "merge"
	case isRebaseFlag:
		opts.strategy = "rebase"
	case isSquashFlag:
		opts.strategy = "squash"
	}
	switch opts.strategy {
	case "squash", "merge", "rebase":
	default:
		return opts, appfault.NewValidationError(fmt.Sprintf(mrMergeErrStrat, opts.strategy))
	}
	return opts, nil
}

// printMrMergeDryRun prints exactly what would happen (spec §5).
func printMrMergeDryRun(slug string, opts mrMergeOpts) {
	del := ""
	if opts.isDeleteBranch {
		del = " --delete-branch"
	}
	fmt.Printf(mrMergeDryRunFmt, opts.number, slug, opts.strategy, del)
}

// isMrMergeConfirmed prints the summary and prompts unless --yes.
func isMrMergeConfirmed(slug string, opts mrMergeOpts) bool {
	del := ""
	if opts.isDeleteBranch {
		del = " · delete head branch"
	}
	fmt.Printf(mrMergeSummaryFmt, opts.number, slug, opts.strategy, del)
	if opts.isYes {
		return true
	}
	fmt.Print(mrMergeConfirmCue)
	return readMergeWriteAnswer()
}

// executeMrMerge runs `gh pr merge`, then reports number/strategy/sha.
func executeMrMerge(slug string, opts mrMergeOpts) *appfault.AppError {
	ghArgs := []string{"pr", "merge", strconv.Itoa(opts.number), "--repo", slug, "--" + opts.strategy}
	if opts.isDeleteBranch {
		ghArgs = append(ghArgs, "--delete-branch")
	}
	if _, appErr := runGhJSON(ghArgs); appErr != nil {
		return appErr
	}
	sha, head := fetchMrMergeResult(slug, opts.number)
	deleted := ""
	if opts.isDeleteBranch && head != "" {
		deleted = fmt.Sprintf(" · branch %s deleted", head)
	}
	fmt.Printf(mrMergeResultFmt, opts.number, opts.strategy, sha, deleted)
	return nil
}

// fetchMrMergeResult reads the merge commit sha + head branch after a merge.
// Best effort: an empty sha still means the merge itself succeeded.
func fetchMrMergeResult(slug string, number int) (sha, head string) {
	out, appErr := runGhJSON([]string{
		"pr", "view", strconv.Itoa(number), "--repo", slug,
		"--json", "mergeCommit,headRefName",
	})
	if appErr != nil {
		return "?", ""
	}
	var decoded struct {
		MergeCommit mrMergeCommit `json:"mergeCommit"`
		HeadRefName string        `json:"headRefName"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &decoded); err != nil {
		return "?", ""
	}
	sha = decoded.MergeCommit.Oid
	if len(sha) > 7 {
		sha = sha[:7]
	}
	if sha == "" {
		sha = "?"
	}
	return sha, decoded.HeadRefName
}

// mrMergeHelpText is the full `--help` text for `mr merge`.
const mrMergeHelpText = `gitmap mr merge — merge a pull request

USAGE
  gitmap mr merge <n> [--strategy squash|merge|rebase] [--delete-branch] [--repo <owner/repo>] [--yes] [--dry-run]
  gitmap mrm <n>     (alias; identical behavior)

FLAGS
  --strategy <s>   squash (default), merge, or rebase; --merge / --rebase / --squash are shorthands
  --delete-branch  Delete the head branch after merging
  --repo <slug>    owner/repo (default: the current repo's origin)
  --yes, -y        Skip the confirmation prompt
  --dry-run        Print what would happen, change nothing

EXAMPLE
  gitmap mr merge 42 --strategy squash --delete-branch --yes

NOTES
  - Requires --yes or an interactive confirmation; nothing merges silently.
  - Merges via the 'gh' CLI: 'gh auth login', or GITHUB_TOKEN / GH_TOKEN.
  - Check triage first: gitmap mr view <n>.
`

// MrMergeHelpExamples registers the 3-part help example (spec 03 §2) for the
// lead's helpdoc registry wiring.
const (
	mrMergeExampleWhat = "Merges an approved pull request with the chosen strategy, optionally deleting the head branch — the PR loop closed entirely inside gitmap."
	mrMergeExampleInv  = `gitmap mr merge 42 --strategy squash --delete-branch --yes`
	mrMergeExampleOut  = "merged #42 (squash) -> abc1234 · branch fix/foo deleted"
)

var MrMergeHelpExamples = map[string]HelpExample{
	"mr merge": {WhatItDoes: mrMergeExampleWhat, Examples: []Example{{Invocation: mrMergeExampleInv, ExpectedOutput: mrMergeExampleOut}}},
	"mrm":      {WhatItDoes: mrMergeExampleWhat, Examples: []Example{{Invocation: strings.Replace(mrMergeExampleInv, "gitmap mr merge", "gitmap mrm", 1), ExpectedOutput: mrMergeExampleOut}}},
}
