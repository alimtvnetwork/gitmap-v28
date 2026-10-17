package cmdmerge

// mr_view.go — `gitmap mr view <n>` / `gitmap mrv`: show PR detail — checks,
// reviews, mergeability — via the `gh` CLI (program 275, proposal P-02).
// One concern per file: flag parsing + orchestration + gh call + rendering +
// help text + help-example registry data (≤300 lines).

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

// mr view messages.
const (
	mrViewName      = "mr view"
	mrViewFlagRepo  = "repo"
	mrViewUsageLine = "Usage: gitmap mr view <n> [--repo <owner/repo>]\n       gitmap mrv <n> [same flags]\n"
	mrViewErrNum    = "missing required argument: <n> (pull request number, e.g. `gitmap mr view 42`)"
	mrViewErrNumFmt = "invalid pull request number %q: want a positive integer"
)

// mrViewOpts carries the parsed flags for `mr view`.
type mrViewOpts struct {
	number int
	repo   string
}

// mrReview decodes one entry of gh's reviews array.
type mrReview struct {
	Author mrAuthor `json:"author"`
	State  string   `json:"state"`
}

// mrCommits decodes gh's commits.totalCount object.
type mrCommits struct {
	TotalCount int `json:"totalCount"`
}

// mrViewPR decodes `gh pr view --json`.
type mrViewPR struct {
	Number            int          `json:"number"`
	Title             string       `json:"title"`
	State             string       `json:"state"`
	URL               string       `json:"url"`
	CreatedAt         string       `json:"createdAt"`
	HeadRefName       string       `json:"headRefName"`
	BaseRefName       string       `json:"baseRefName"`
	Mergeable         string       `json:"mergeable"`
	MergeStateStatus  string       `json:"mergeStateStatus"`
	Additions         int          `json:"additions"`
	Deletions         int          `json:"deletions"`
	ChangedFiles      int          `json:"changedFiles"`
	Commits           mrCommits    `json:"commits"`
	ReviewDecision    string       `json:"reviewDecision"`
	StatusCheckRollup []mrCheckRun `json:"statusCheckRollup"`
	Reviews           []mrReview   `json:"reviews"`
	Author            mrAuthor     `json:"author"`
}

// RunMrView implements `gitmap mr view` / `gitmap mrv` for rootdata wiring.
func RunMrView(args []string) error {
	if appErr := runMrView(args); appErr != nil {
		return appErr
	}
	return nil
}

// runMrView executes the view flow: preflight, resolve, fetch, render.
func runMrView(args []string) *appfault.AppError {
	if isHelpRequested(args) {
		fmt.Print(mrViewHelpText)
		cliexit.Exit(0)
	}
	opts, appErr := parseMrViewFlags(args)
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
	pr, appErr := fetchMrView(slug, opts.number)
	if appErr != nil {
		return appErr
	}
	printMrView(pr)
	return nil
}

// parseMrViewFlags parses the `mr view` flag set + the required <n>.
func parseMrViewFlags(args []string) (mrViewOpts, *appfault.AppError) {
	opts := mrViewOpts{}
	fs := flag.NewFlagSet(mrViewName, flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&opts.repo, mrViewFlagRepo, "", "owner/repo (default: the current repo's origin)")
	if err := fs.Parse(flagutil.ReorderFlagsBeforeArgs(args)); err != nil {
		return opts, appfault.NewValidationError(mrViewUsageLine + "error: " + err.Error())
	}
	rest := fs.Args()
	if len(rest) == 0 {
		fmt.Print(mrViewUsageLine)
		return opts, appfault.NewValidationError(mrViewErrNum)
	}
	n, err := strconv.Atoi(rest[0])
	if err != nil || n <= 0 {
		return opts, appfault.NewValidationError(fmt.Sprintf(mrViewErrNumFmt, rest[0]))
	}
	opts.number = n
	return opts, nil
}

// fetchMrView shells out to `gh pr view --json` and decodes the object.
func fetchMrView(slug string, number int) (mrViewPR, *appfault.AppError) {
	pr := mrViewPR{}
	out, appErr := runGhJSON([]string{
		"pr", "view", strconv.Itoa(number), "--repo", slug,
		"--json", "number,title,state,author,createdAt,url,headRefName,baseRefName,mergeable,mergeStateStatus,additions,deletions,changedFiles,commits,statusCheckRollup,reviews,reviewDecision",
	})
	if appErr != nil {
		return pr, appErr
	}
	if err := json.Unmarshal(bytes.TrimSpace(out), &pr); err != nil {
		return pr, appfault.NewValidationError("could not parse `gh pr view` output: " + string(bytes.TrimSpace(out)))
	}
	return pr, nil
}

// mrReviewTally dedups reviews by author (keeping the latest) and counts
// approvals vs change-requests.
func mrReviewTally(reviews []mrReview) (approved, changes, total int) {
	latest := map[string]string{}
	order := []string{}
	for _, r := range reviews {
		login := r.Author.Login
		if login == "" {
			continue
		}
		if _, seen := latest[login]; !seen {
			order = append(order, login)
		}
		latest[login] = r.State
	}
	total = len(order)
	for _, login := range order {
		switch latest[login] {
		case "APPROVED":
			approved++
		case "CHANGES_REQUESTED":
			changes++
		}
	}
	return approved, changes, total
}

// mrMergeableText renders gh's mergeable enum.
func mrMergeableText(mergeable string) string {
	switch mergeable {
	case "MERGEABLE":
		return "YES"
	case "CONFLICTING":
		return "NO (conflicts)"
	case "UNKNOWN":
		return "UNKNOWN (checks pending)"
	default:
		return "?"
	}
}

// printMrView renders the proposal P-02 one-liner plus detail lines.
func printMrView(pr mrViewPR) {
	approved, _, reviewers := mrReviewTally(pr.Reviews)
	checks := mrChecksSummary(pr.StatusCheckRollup)
	reviews := fmt.Sprintf("%d/%d %s", approved, reviewers, mrReviewDecisionText(pr.ReviewDecision))
	fmt.Printf("PR #%d %q · checks: %s · reviews: %s · +%d/-%d · mergeable: %s\n",
		pr.Number, pr.Title, checks, reviews, pr.Additions, pr.Deletions, mrMergeableText(pr.Mergeable))
	fmt.Printf("  url:      %s\n", pr.URL)
	fmt.Printf("  branches: %s -> %s (%s)\n", pr.HeadRefName, pr.BaseRefName, strings.ToLower(pr.State))
	fmt.Printf("  diff:     +%d/-%d across %d files · %d commits · opened %s ago by %s\n",
		pr.Additions, pr.Deletions, pr.ChangedFiles, pr.Commits.TotalCount, mrRelativeAge(pr.CreatedAt), pr.Author.Login)
	if pr.MergeStateStatus != "" {
		fmt.Printf("  state:    %s\n", pr.MergeStateStatus)
	}
}

// mrViewHelpText is the full `--help` text for `mr view`.
const mrViewHelpText = `gitmap mr view — show pull request detail

USAGE
  gitmap mr view <n> [--repo <owner/repo>]
  gitmap mrv <n>     (alias; identical behavior)

FLAGS
  --repo <slug>    owner/repo (default: the current repo's origin)

EXAMPLE
  gitmap mr view 42

NOTES
  - Answers "can I merge this" from the terminal: check-rollup, per-author
    review tally, diff stats, and the mergeable flag.
  - Detail via the 'gh' CLI: 'gh auth login', or GITHUB_TOKEN / GH_TOKEN.
`

// MrViewHelpExamples registers the 3-part help example (spec 03 §2) for the
// lead's helpdoc registry wiring.
const (
	mrViewExampleWhat = "Shows one pull request's checks, per-author review tally, diff stats, and mergeability — the full triage picture without opening a browser."
	mrViewExampleInv  = `gitmap mr view 42`
	mrViewExampleOut  = "PR #42 \"release - bump to v1.1.0\" · checks: 3/3 ok · reviews: 1/1 approved · +180/-24 · mergeable: YES\n  url:      https://github.com/alimtvnetwork/go-email-reader/pull/42\n  branches: release/v1.1.0 -> main (open)\n  diff:     +180/-24 across 6 files · 3 commits · opened 2d ago by alimtvnetwork"
)

var MrViewHelpExamples = map[string]HelpExample{
	"mr view": {WhatItDoes: mrViewExampleWhat, Examples: []Example{{Invocation: mrViewExampleInv, ExpectedOutput: mrViewExampleOut}}},
	"mrv":     {WhatItDoes: mrViewExampleWhat, Examples: []Example{{Invocation: strings.Replace(mrViewExampleInv, "gitmap mr view", "gitmap mrv", 1), ExpectedOutput: mrViewExampleOut}}},
}
