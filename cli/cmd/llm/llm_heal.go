package llm

import (
	"fmt"
	"runtime"
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdautofix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfix"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
)

// healPhase is TrainPhase 6 ("Heal & Fix"): sub-step A heals workspace git-state
// via cmdfix, sub-step B audits file content via the `gitmap fix` engine
// (cli/cmdautofix). Report-only by default; --heal-apply applies git-state
// recipes only (stash|wip|discard; never under --text-only).
type healPhase struct{}

func (healPhase) Name() string { return "6. Heal & Fix" }

func (healPhase) Run(opts TrainOptions) *apperror.AppError {
	if err := runHealGitState(opts); err != nil {
		return err
	}
	return runFixContentAudit()
}

// healCategoryBuckets maps SummaryReason keywords to report buckets in
// priority order; unknown reasons fall back to "other".
var healCategoryBuckets = [][2]string{
	{"diverged", "diverged"}, {"behind", "diverged"}, {"ahead", "diverged"},
	{"untracked", "untracked"}, {"dirty", "dirty-worktree"},
	{"modified", "dirty-worktree"}, {"uncommitted", "dirty-worktree"},
	{"merge conflict", "dirty-worktree"},
}

func bucketHealReason(reason string) string {
	lower := strings.ToLower(reason)
	for _, b := range healCategoryBuckets {
		if strings.Contains(lower, b[0]) {
			return b[1]
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
	fmt.Println("Next actions: `gitmap stash ls`, `gitmap stash <repo> <action>` (stash|wip|discard)")
	if opts.HealApply == "" {
		return nil
	}
	if opts.IsTextOnly {
		fmt.Println("[heal] --heal-apply ignored under --text-only (no mutations)")
		return nil
	}
	action := strings.ToLower(opts.HealApply)
	if action != "stash" && action != "wip" && action != "discard" {
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
	for _, item := range items {
		bucket := bucketHealReason(item.SummaryReason)
		counts[bucket]++
		if example[bucket] == "" {
			example[bucket] = item.RepoName
		}
	}
	fmt.Println("CATEGORY | REPOS | EXAMPLE REPO")
	total := 0
	for _, bucket := range []string{"dirty-worktree", "diverged", "untracked", "other"} {
		total += counts[bucket]
		if example[bucket] == "" {
			example[bucket] = "-"
		}
		fmt.Printf("%s | %d | %s\n", bucket, counts[bucket], example[bucket])
	}
	if total == 0 {
		fmt.Println("Workspace git-state is clean — nothing to heal.")
		return
	}
	fmt.Printf("Total repos needing remediation: %d\n", total)
}

// runFixContentAudit is sub-step B: check-only scan via the fix engine's Scan
// entry (never prompts, never writes); the operator runs `gitmap fix all -y`.
func runFixContentAudit() *apperror.AppError {
	fmt.Println("STAGE 6b: FIX — file-content audit via fix engine (report-only)")
	result, err := cmdautofix.Scan(cmdautofix.Options{Workers: runtime.NumCPU()})
	if err != nil {
		return apperror.WrapSimple(err, "healPhase.scan")
	}
	printScanRollup(result)
	fmt.Println("----------------------------------------------------------------------\nHEAL & FIX ROLL-UP\n----------------------------------------------------------------------")
	fmt.Println("git-state remediation : report-only — run `gitmap stash ls` / `gitmap stash <repo> <action>`")
	fmt.Println("content-fix audit     : report-only — run `gitmap fix all -y` to write fixes")
	return nil
}

// printScanRollup prints the per-category roll-up from a ScanResult.
func printScanRollup(result *cmdautofix.ScanResult) {
	counts := map[string]int{}
	files := map[string]map[string]bool{}
	for _, v := range result.Violations {
		counts[v.Category]++
		if files[v.Category] == nil {
			files[v.Category] = map[string]bool{}
		}
		files[v.Category][v.Path] = true
	}
	cats := make([]string, 0, len(counts))
	for cat := range counts {
		cats = append(cats, cat)
	}
	sort.Strings(cats)
	fmt.Println("CATEGORY | VIOLATIONS | FILES FLAGGED")
	total := 0
	for _, cat := range cats {
		fmt.Printf("%s | %d | %d\n", cat, counts[cat], len(files[cat]))
		total += counts[cat]
	}
	fmt.Printf("Files scanned: %d | Violations: %d | Time: %s\n",
		result.FilesScanned, total, result.Elapsed)
	if total == 0 {
		fmt.Println("Content audit clean — no violations found.")
	}
}
