package cmdmerge

// mr_family.go — shared plumbing for the `mr` PR family (program 275):
// `mr list|view|merge|checkout` (+ single-word aliases mrl/mrv/mrm/mrco).
// The two-word form is routed here from DispatchMoveMerge before the
// `merge write` branch; the aliases are wired by the lead in rootdata.go
// to the exported RunMr* entry points in mr_list.go / mr_view.go /
// mr_merge.go / mr_checkout.go. Help text + help-example registry data
// live at the bottom of each command file (merge_write_help.go pattern).

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// mr family subcommand tokens.
const (
	mrSubList     = "list"
	mrSubView     = "view"
	mrSubMerge    = "merge"
	mrSubCheckout = "checkout"
)

// mrErrGhMissing names the missing tool, per shared-patterns spec §7.
const mrErrGhMissing = "the `gh` CLI is required for the `mr` commands but was not found on PATH.\nFix: install gh (https://cli.github.com), then `gh auth login`, or set GITHUB_TOKEN / GH_TOKEN."

// mrAuthor decodes gh's author object.
type mrAuthor struct {
	Login string `json:"login"`
	Name  string `json:"name"`
}

// mrCheckRun decodes one entry of gh's statusCheckRollup array.
type mrCheckRun struct {
	Name       string `json:"name"`
	Status     string `json:"status"`
	Conclusion string `json:"conclusion"`
}

// requireMrGh ensures the `gh` CLI is available.
func requireMrGh() *appfault.AppError {
	if _, err := exec.LookPath("gh"); err != nil {
		return appfault.NewNotFoundError(mrErrGhMissing)
	}
	return nil
}

// mrFamilySubcommand splits a leading mr-family verb from its args.
func mrFamilySubcommand(args []string) (string, []string, bool) {
	if len(args) == 0 {
		return "", args, false
	}
	switch args[0] {
	case mrSubList, mrSubView, mrSubMerge, mrSubCheckout:
		return args[0], args[1:], true
	}
	return "", args, false
}

// dispatchMrFamily routes `gitmap mr <list|view|merge|checkout>`; it
// returns handled=false for anything else (including bare `mr`, which
// stays on the merge-write branch).
func dispatchMrFamily(command string, args []string) (bool, error) {
	if command != constants.CmdMergeWriteA {
		return false, nil
	}
	sub, subArgs, isFamily := mrFamilySubcommand(args)
	if !isFamily {
		return false, nil
	}
	switch sub {
	case mrSubList:
		return true, RunMrList(subArgs)
	case mrSubView:
		return true, RunMrView(subArgs)
	case mrSubMerge:
		return true, RunMrMerge(subArgs)
	case mrSubCheckout:
		return true, RunMrCheckout(subArgs)
	}
	return false, nil
}

// runGhJSON runs `gh` with the given args and returns stdout; gh's stderr
// becomes the error detail — never swallowed.
func runGhJSON(ghArgs []string) ([]byte, *appfault.AppError) {
	cmd := exec.Command("gh", ghArgs...)
	var stderr bytes.Buffer
	cmd.Stderr = &stderr
	out, err := cmd.Output()
	if err != nil {
		return nil, apperror.NewExecutionError(fmt.Sprintf("gh %s failed: %s", strings.Join(ghArgs, " "), ghFailureDetail(stderr.String(), err)))
	}
	return out, nil
}

// resolveMrRepo returns the --repo flag value or the current repo's slug.
func resolveMrRepo(repoFlag string) (string, *appfault.AppError) {
	if strings.TrimSpace(repoFlag) != "" {
		return strings.TrimSpace(repoFlag), nil
	}
	return resolveGitHubSlug()
}

// mrRelativeAge renders an RFC3339 timestamp as "5m" / "5h" / "2d".
func mrRelativeAge(rfc3339 string) string {
	t, err := time.Parse(time.RFC3339, rfc3339)
	if err != nil {
		return "?"
	}
	d := time.Since(t)
	if d < 0 {
		d = 0
	}
	switch {
	case d < time.Hour:
		return fmt.Sprintf("%dm", int(d.Minutes()))
	case d < 24*time.Hour:
		return fmt.Sprintf("%dh", int(d.Hours()))
	default:
		return fmt.Sprintf("%dd", int(d.Hours()/24))
	}
}

// mrChecksSummary summarizes a statusCheckRollup array as "3/3 ok",
// "2/3 (1 failed)", "1/2 (1 pending)", or "none".
func mrChecksSummary(rollup []mrCheckRun) string {
	total := len(rollup)
	if total == 0 {
		return "none"
	}
	passed := 0
	failed := 0
	for _, c := range rollup {
		switch c.Conclusion {
		case "SUCCESS", "SKIPPED", "NEUTRAL":
			passed++
		case "FAILURE", "TIMED_OUT", "ACTION_REQUIRED", "CANCELED", "STALE":
			failed++
		}
	}
	pending := total - passed - failed
	switch {
	case failed > 0:
		return fmt.Sprintf("%d/%d (%d failed)", passed, total, failed)
	case pending > 0:
		return fmt.Sprintf("%d/%d (%d pending)", passed, total, pending)
	default:
		return fmt.Sprintf("%d/%d ok", passed, total)
	}
}

// mrReviewDecisionText renders gh's reviewDecision enum for the list table.
func mrReviewDecisionText(decision string) string {
	switch decision {
	case "APPROVED":
		return "approved"
	case "CHANGES_REQUESTED":
		return "changes requested"
	case "REVIEW_REQUIRED":
		return "pending"
	default:
		return "none"
	}
}
