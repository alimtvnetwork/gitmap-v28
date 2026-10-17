package cmdmerge

// mr_checkout.go — `gitmap mr checkout <n>` / `gitmap mrco`: check out a
// pull request locally via the `gh` CLI (program 275, proposal P-04).
// One concern per file: flag parsing + orchestration + gh call + rendering +
// help text + help-example registry data (≤300 lines).

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// mr checkout messages.
const (
	mrCheckoutName      = "mr checkout"
	mrCheckoutFlagRepo  = "repo"
	mrCheckoutUsageLine = "Usage: gitmap mr checkout <n> [--repo <owner/repo>]\n       gitmap mrco <n> [same flags]\n"
	mrCheckoutErrNum    = "missing required argument: <n> (pull request number, e.g. `gitmap mr checkout 42`)"
	mrCheckoutErrNumFmt = "invalid pull request number %q: want a positive integer"
	mrCheckoutResultFmt = "fetched PR #%d -> branch %s · checked out%s\n"
)

// mrCheckoutOpts carries the parsed flags for `mr checkout`.
type mrCheckoutOpts struct {
	number int
	repo   string
}

// RunMrCheckout implements `gitmap mr checkout` / `gitmap mrco` for rootdata wiring.
func RunMrCheckout(args []string) error {
	if appErr := runMrCheckout(args); appErr != nil {
		return appErr
	}
	return nil
}

// runMrCheckout executes the checkout flow: preflight, resolve, checkout.
func runMrCheckout(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(mrCheckoutHelpText)
		cliexit.Exit(0)
	}
	opts, appErr := parseMrCheckoutFlags(args)
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
	head, base, appErr := fetchMrCheckoutRefs(slug, opts.number)
	if appErr != nil {
		return appErr
	}
	if _, appErr := runGhJSON([]string{"pr", "checkout", strconv.Itoa(opts.number), "--repo", slug}); appErr != nil {
		return appErr
	}
	behind := mrCommitsBehind(base)
	fmt.Printf(mrCheckoutResultFmt, opts.number, head, behind)
	return nil
}

// parseMrCheckoutFlags parses the `mr checkout` flag set + the required <n>.
func parseMrCheckoutFlags(args []string) (mrCheckoutOpts, *appfault.AppError) {
	opts := mrCheckoutOpts{}
	fs := flag.NewFlagSet(mrCheckoutName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.repo, mrCheckoutFlagRepo, "", "owner/repo (default: the current repo's origin)")
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(mrCheckoutUsageLine + "error: " + err.Error())
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Print(mrCheckoutUsageLine)
		return opts, appfault.NewValidationError(mrCheckoutErrNum)
	}
	n, err := strconv.Atoi(rest[0])
	if err != nil || n <= 0 {
		return opts, appfault.NewValidationError(fmt.Sprintf(mrCheckoutErrNumFmt, rest[0]))
	}
	opts.number = n
	return opts, nil
}

// fetchMrCheckoutRefs reads the PR's head/base branch names via gh.
func fetchMrCheckoutRefs(slug string, number int) (head, base string, appErr *appfault.AppError) {
	out, appErr := runGhJSON([]string{
		"pr", "view", strconv.Itoa(number), "--repo", slug,
		"--json", "headRefName,baseRefName",
	})
	if appErr != nil {
		return "", "", appErr
	}
	var decoded struct {
		HeadRefName string `json:"headRefName"`
		BaseRefName string `json:"baseRefName"`
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &decoded); err != nil {
		return "", "", appfault.NewValidationError("could not parse `gh pr view` output: " + string(bytes.TrimSpace(out)))
	}
	return decoded.HeadRefName, decoded.BaseRefName, nil
}

// mrCommitsBehind counts how far HEAD is behind the base branch.
// Best effort: an unresolvable base simply omits the suffix.
func mrCommitsBehind(base string) string {
	if strings.TrimSpace(base) == "" {
		return ""
	}
	out, err := exec.Command("git", "rev-list", "--count", "HEAD.."+base).Output()
	if err != nil {
		return ""
	}
	n, err := strconv.Atoi(strings.TrimSpace(string(out)))
	if err != nil {
		return ""
	}
	return fmt.Sprintf(" · %d commits behind %s", n, base)
}

// mrCheckoutHelpText is the full `--help` text for `mr checkout`.
const mrCheckoutHelpText = `gitmap mr checkout — check out a pull request locally

USAGE
  gitmap mr checkout <n> [--repo <owner/repo>]
  gitmap mrco <n>        (alias; identical behavior)

FLAGS
  --repo <slug>    owner/repo (default: the current repo's origin)

EXAMPLE
  gitmap mr checkout 42

NOTES
  - Runs 'gh pr checkout' in the current directory (a git repo).
  - Reports how many commits the PR branch is behind the base branch.
  - Via the 'gh' CLI: 'gh auth login', or GITHUB_TOKEN / GH_TOKEN.
`

// MrCheckoutHelpExamples registers the 3-part help example (spec 03 §2) for
// the lead's helpdoc registry wiring.
const (
	mrCheckoutExampleWhat = "Fetches a pull request's head branch and checks it out locally — no hand-typed fetch refspecs, no typos."
	mrCheckoutExampleInv  = `gitmap mr checkout 42`
	mrCheckoutExampleOut  = "fetched PR #42 -> branch fix/queue-double-send · checked out · 3 commits behind main"
)

var MrCheckoutHelpExamples = map[string]HelpExample{
	"mr checkout": {WhatItDoes: mrCheckoutExampleWhat, Examples: []Example{{Invocation: mrCheckoutExampleInv, ExpectedOutput: mrCheckoutExampleOut}}},
	"mrco":        {WhatItDoes: mrCheckoutExampleWhat, Examples: []Example{{Invocation: strings.Replace(mrCheckoutExampleInv, "gitmap mr checkout", "gitmap mrco", 1), ExpectedOutput: mrCheckoutExampleOut}}},
}
