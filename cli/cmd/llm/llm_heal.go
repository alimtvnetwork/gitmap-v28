package llm

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautofix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
)

// healPhase is TrainPhase 6 ("Heal & Fix"): sub-step A heals workspace
// git-state via the cmdfix engine, sub-step B audits file content via the
// cmdautofix engine. Report-only by default; git-state recipes apply only
// with --heal-apply stash|wip|discard (ignored under --text-only).
type healPhase struct{}

func (healPhase) Name() string { return "6. Heal & Fix" }

func (healPhase) Run(opts TrainOptions) *apperror.AppError {
	if err := runHealGitState(opts); err != nil {
		return err
	}
	return runFixContentAudit()
}

// healCategoryBuckets maps SummaryReason keywords to report buckets, in
// priority order. Buckets are derived from the reason string; unknown
// reasons fall back to "other".
var healCategoryBuckets = []struct {
	keyword string
	bucket  string
}{
	{"diverged", "diverged"},
	{"behind", "diverged"},
	{"ahead", "diverged"},
	{"untracked", "untracked"},
	{"dirty", "dirty-worktree"},
	{"modified", "dirty-worktree"},
	{"uncommitted", "dirty-worktree"},
	{"merge conflict", "dirty-worktree"},
}

func bucketHealReason(reason string) string {
	lower := strings.ToLower(reason)
	for _, b := range healCategoryBuckets {
		if strings.Contains(lower, b.keyword) {
			return b.bucket
		}
	}
	return "other"
}

// runHealGitState is sub-step A: report pending git-state remediation per
// category, then optionally apply recipes via --heal-apply.
func runHealGitState(opts TrainOptions) *apperror.AppError {
	fmt.Println("STAGE 6a: HEAL — workspace git-state remediation (report-only)")
	items := cmdremediation.LoadRemediationState()
	if len(items) == 0 {
		if err := cmdfix.RunFix([]string{}, ""); err != nil {
			return apperror.WrapSimple(err, "healPhase.discover")
		}
		items = cmdremediation.LoadRemediationState()
	}
	printHealCategoryTable(items)
	fmt.Println("Next actions: `gitmap fix ls`, `gitmap fix <repo> <action>` (stash|wip|discard)")
	if opts.HealApply == "" {
		return nil
	}
	if opts.IsTextOnly {
		fmt.Println("[heal] --heal-apply ignored under --text-only (no mutations)")
		return nil
	}
	action := strings.ToLower(opts.HealApply)
	switch action {
	case "stash", "wip", "discard":
	default:
		return apperror.WrapSimple(
			fmt.Errorf("invalid --heal-apply %q: want stash|wip|discard", opts.HealApply),
			"healPhase.apply")
	}
	fmt.Printf("[heal] applying %q recipes across all pending repos...\n", action)
	if err := cmdfix.RunFix([]string{"all", action}, ""); err != nil {
		return apperror.WrapSimple(err, "healPhase.apply")
	}
	return nil
}

func printHealCategoryTable(items []cmdremediation.RemediationItem) {
	counts := map[string]int{}
	example := map[string]string{}
	buckets := []string{"dirty-worktree", "diverged", "untracked", "other"}
	for _, item := range items {
		bucket := bucketHealReason(item.SummaryReason)
		counts[bucket]++
		if example[bucket] == "" {
			example[bucket] = item.RepoName
		}
	}
	fmt.Println("CATEGORY | REPOS | EXAMPLE REPO")
	total := 0
	for _, bucket := range buckets {
		total += counts[bucket]
		name := example[bucket]
		if name == "" {
			name = "-"
		}
		fmt.Printf("%s | %d | %s\n", bucket, counts[bucket], name)
	}
	if total == 0 {
		fmt.Println("Workspace git-state is clean — nothing to heal.")
		return
	}
	fmt.Printf("Total repos needing remediation: %d\n", total)
}

// runFixContentAudit is sub-step B: dry-run the autofix engine over the
// workspace and re-print its per-category roll-up. The train never passes
// --apply; the operator runs `gitmap autofix --apply` to write fixes.
func runFixContentAudit() *apperror.AppError {
	fmt.Println("STAGE 6b: FIX — file-content audit via autofix engine (report-only)")
	if err := cmdautofix.RunAutofixCmd([]string{"."}); err != nil {
		return apperror.WrapSimple(err, "healPhase.autofix")
	}
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println("HEAL & FIX ROLL-UP")
	fmt.Println("----------------------------------------------------------------------")
	fmt.Println("git-state remediation : report-only — run `gitmap fix ls` / `gitmap fix <repo> <action>`")
	fmt.Println("content-fix audit     : report-only — run `gitmap autofix --apply` to write fixes")
	return nil
}
