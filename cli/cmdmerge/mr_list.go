package cmdmerge

// mr_list.go — `gitmap mr list` / `gitmap mrl`: list pull requests via the
// `gh` CLI (program 275, proposal P-01). One concern per file: flag parsing
// + orchestration + gh call + table rendering + help text + help-example
// registry data all live here (≤300 lines).

import (
	"bytes"
	"encoding/json"
	"flag"
	"fmt"
	"io"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/flagutil"
)

// mr list flag names, defaults, and messages.
const (
	mrListName       = "mr list"
	mrListFlagState  = "state"
	mrListFlagLimit  = "limit"
	mrListFlagRepo   = "repo"
	mrListDefState   = "open"
	mrListDefLimit   = 30
	mrListUsageLine  = "Usage: gitmap mr list [--state open|closed|merged|all] [--limit <n>] [--repo <owner/repo>]\n       gitmap mrl [same flags]\n"
	mrListErrState   = "invalid --state %q: want one of open|closed|merged|all"
	mrListErrLimit   = "invalid --limit %d: want a positive number"
	mrListTitleWidth = 48
)

// mrListOpts carries the parsed flags for `mr list`.
type mrListOpts struct {
	state string
	limit int
	repo  string
}

// mrPRListItem decodes one entry of `gh pr list --json`.
type mrPRListItem struct {
	Number            int          `json:"number"`
	Title             string       `json:"title"`
	URL               string       `json:"url"`
	CreatedAt         string       `json:"createdAt"`
	HeadRefName       string       `json:"headRefName"`
	BaseRefName       string       `json:"baseRefName"`
	ReviewDecision    string       `json:"reviewDecision"`
	StatusCheckRollup []mrCheckRun `json:"statusCheckRollup"`
	Author            mrAuthor     `json:"author"`
}

// RunMrList implements `gitmap mr list` / `gitmap mrl` for rootdata wiring.
func RunMrList(args []string) error {
	if appErr := runMrList(args); appErr != nil {
		return appErr
	}
	return nil
}

// runMrList executes the list flow: preflight, resolve, fetch, render.
func runMrList(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(mrListHelpText)
		cliexit.Exit(0)
	}
	opts, appErr := parseMrListFlags(args)
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
	items, appErr := fetchMrList(slug, opts)
	if appErr != nil {
		return appErr
	}
	printMrListTable(items, slug, opts.state)
	return nil
}

// parseMrListFlags parses the `mr list` flag set.
func parseMrListFlags(args []string) (mrListOpts, *appfault.AppError) {
	opts := mrListOpts{state: mrListDefState, limit: mrListDefLimit}
	fs := flag.NewFlagSet(mrListName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.state, mrListFlagState, mrListDefState, "PR state: open|closed|merged|all")
	fs.IntVar(&opts.limit, mrListFlagLimit, mrListDefLimit, "max PRs to list")
	fs.StringVar(&opts.repo, mrListFlagRepo, "", "owner/repo (default: current repo)")
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(mrListUsageLine + "error: " + err.Error())
	}
	switch opts.state {
	case "open", "closed", "merged", "all":
	default:
		return opts, appfault.NewValidationError(fmt.Sprintf(mrListErrState, opts.state))
	}
	if opts.limit <= 0 {
		return opts, appfault.NewValidationError(fmt.Sprintf(mrListErrLimit, opts.limit))
	}
	return opts, nil
}

// fetchMrList shells out to `gh pr list --json` and decodes the array.
func fetchMrList(slug string, opts mrListOpts) ([]mrPRListItem, *appfault.AppError) {
	out, appErr := runGhJSON([]string{
		"pr", "list", "--repo", slug,
		"--state", opts.state, "--limit", fmt.Sprintf("%d", opts.limit),
		"--json", "number,title,author,createdAt,headRefName,baseRefName,reviewDecision,url,statusCheckRollup",
	})
	if appErr != nil {
		return nil, appErr
	}
	var items []mrPRListItem
	if err := json.Unmarshal(bytes.TrimSpace(out), &items); err != nil {
		return nil, appfault.NewValidationError("could not parse `gh pr list` output: " + string(bytes.TrimSpace(out)))
	}
	return items, nil
}

// printMrListTable renders the proposal P-01 table (#, title, age, checks,
// reviews) — plain text, no ANSI (theme layer handles color).
func printMrListTable(items []mrPRListItem, slug, state string) {
	if len(items) == 0 {
		fmt.Printf("no %s pull requests in %s\n", state, slug)
		return
	}
	fmt.Printf("%-4s %-48s %-5s %-14s %s\n", "#", "TITLE", "AGE", "CHECKS", "REVIEWS")
	for _, item := range items {
		title := item.Title
		if len(title) > mrListTitleWidth {
			title = title[:mrListTitleWidth-1] + "…"
		}
		fmt.Printf("%-4d %-48s %-5s %-14s %s\n",
			item.Number, title, mrRelativeAge(item.CreatedAt),
			mrChecksSummary(item.StatusCheckRollup), mrReviewDecisionText(item.ReviewDecision))
	}
	fmt.Printf("%d %s pull request(s) in %s\n", len(items), state, slug)
}

// mrListHelpText is the full `--help` text for `mr list`.
const mrListHelpText = `gitmap mr list — list pull requests for a repo

USAGE
  gitmap mr list [--state open|closed|merged|all] [--limit <n>] [--repo <owner/repo>]
  gitmap mrl     (alias; identical behavior)

FLAGS
  --state <s>      PR state: open (default), closed, merged, all
  --limit <n>      Max PRs to list (default: 30)
  --repo <slug>    owner/repo (default: the current repo's origin)

EXAMPLE
  gitmap mr list --state open --limit 10

NOTES
  - Lists via the 'gh' CLI: 'gh auth login', or GITHUB_TOKEN / GH_TOKEN.
  - CHECKS summarizes the check-rollup (passed/total, failures, pending).
`

// MrListHelpExamples registers the 3-part help example (spec 03 §2) for the
// lead's helpdoc registry wiring.
const (
	mrListExampleWhat = "Lists pull requests for a repo with age, check-rollup, and review state in one table — no browser round-trip."
	mrListExampleInv  = `gitmap mr list --state open --limit 10`
	mrListExampleOut  = "#    TITLE                                            AGE   CHECKS         REVIEWS\n42   release - bump to v1.1.0                     2d    3/3 ok         approved\n41   fix: queue dispatch double-send              5h    2/3 (1 failed) pending\n2 open pull request(s) in alimtvnetwork/go-email-reader"
)

var MrListHelpExamples = map[string]HelpExample{
	"mr list": {WhatItDoes: mrListExampleWhat, Examples: []Example{{Invocation: mrListExampleInv, ExpectedOutput: mrListExampleOut}}},
	"mrl":     {WhatItDoes: mrListExampleWhat, Examples: []Example{{Invocation: strings.Replace(mrListExampleInv, "gitmap mr list", "gitmap mrl", 1), ExpectedOutput: mrListExampleOut}}},
}
