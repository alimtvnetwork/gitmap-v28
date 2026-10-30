package cmdfix

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixreleasetags"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/parallel"
)

func hasAnyMatch(args []string, candidates ...string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		for _, c := range candidates {
			if low == c {
				return true
			}
		}
	}

	return false
}

func isFixAgyRequest(args []string) bool {
	return hasAnyMatch(args, "agy", "aef", "agy-errors-fix")
}

func isFixCasingRequest(args []string) bool {
	return hasAnyMatch(args, "casing", "case", "readme-casing")
}

func isFixAllRequested(args []string) bool {
	return hasAnyMatch(args, "all", "--all", "-all", "-a")
}

func isFixPromptRequested(args []string) bool {
	return hasAnyMatch(args, "--prompt", "-p", "prompt", "--interactive", "-i", "interactive")
}

func executeTargetedFixAction(repoPath string, action string) (bool, error) {
	return cmdconflict.ResolveConflictAction(repoPath, action)
}

func resolveFixAllAction(args []string, aliasOverride string) string {
	if aliasOverride != "" {
		return aliasOverride
	}

	for _, a := range args {
		low := strings.ToLower(a)
		if isFixActionWord(low) {
			return low
		}
	}

	return "stash"
}

func runFixAll(action string, items []cmdremediation.RemediationItem) error {
	reconcileAllCasingDiscrepancies()
	if len(items) == 0 && cmdpull.HasPendingPullFailures() {
		return cmdpull.RunBatchPullFix([]string{"--all"})
	}
	if len(items) == 0 {
		fmt.Printf("%s No pending repositories require remediation.\n", constants.ColorGreen+"✓"+constants.ColorReset)

		return nil
	}

	return executeParallelFixAll(action, items, cmdpull.HasPendingPullFailures())
}

func reconcileAllCasingDiscrepancies() {
	for _, r := range discoverReposForFixLs() {
		discrepancies, err := DetectCasingDiscrepancies(r.AbsolutePath)
		if err == nil && len(discrepancies) > 0 {
			_ = ReconcileDiscrepancyList(discrepancies)
		}
	}
}

func executeParallelFixAll(action string, items []cmdremediation.RemediationItem, hasPullFailures bool) error {
	printParallelFixHeader(len(items), action)
	concurrency := parallel.ResolveConcurrency(0)
	failedRepos := runParallelFixWorkers(action, items, concurrency)
	handlePostFixPull(hasPullFailures)

	return evaluateParallelFixResult(len(items), failedRepos)
}

func printParallelFixHeader(itemCount int, action string) {
	concurrency := parallel.ResolveConcurrency(0)
	effectiveWorkers := concurrency
	if effectiveWorkers > itemCount {
		effectiveWorkers = itemCount
	}
	fmt.Printf("\n  %s→%s %sgitmap fix --all%s %s(action: %s)%s\n",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, constants.ColorReset,
		constants.ColorDim, action, constants.ColorReset)
	fmt.Printf("  Remediating %d repository(ies) in parallel (%d workers)\n\n",
		itemCount, effectiveWorkers)
}

func handlePostFixPull(hasPullFailures bool) {
	if hasPullFailures {
		fmt.Println()
		_ = cmdpull.RunBatchPullFix([]string{"--all"})
	}
}

func evaluateParallelFixResult(total int, failedRepos []string) error {
	if len(failedRepos) > 0 {
		fmt.Printf("\n  %s Remediation completed with %d failure(s):\n",
			constants.ColorYellow+"⚠"+constants.ColorReset, len(failedRepos))
		for _, f := range failedRepos {
			fmt.Printf("    • %s: remediation failed (run 'gitmap fix %s' to inspect)\n", f, f)
		}
		fmt.Println()

		return apperror.NewExecutionError(fmt.Sprintf("remediation failed on %d repository(ies): %s",
			len(failedRepos), strings.Join(failedRepos, ", ")))
	}
	fmt.Printf("\n  %s All %d repository(ies) remediated successfully.\n\n",
		constants.ColorGreen+"✔"+constants.ColorReset, total)

	return nil
}

func runParallelFixWorkers(action string, items []cmdremediation.RemediationItem, concurrency int) []string {
	var mu sync.Mutex
	var failedRepos []string
	completed := 0
	total := len(items)

	_ = parallel.ExecuteParallel(items, concurrency, func(workerID int, item cmdremediation.RemediationItem) error {
		return executeFixWorkerStep(&mu, &completed, total, action, item, &failedRepos)
	})

	return failedRepos
}

func executeFixWorkerStep(mu *sync.Mutex, completed *int, total int, action string, item cmdremediation.RemediationItem, failedRepos *[]string) error {
	recipe := resolveFixStepRecipe(action, &item)
	res := executeFixRecipeBatch(&item, recipe)
	reportFixBatchResult(mu, completed, total, res, failedRepos)

	return res.Err
}

func resolveFixStepRecipe(action string, item *cmdremediation.RemediationItem) gitutil.RemediationRecipe {
	if len(item.Recipes) > 0 {
		idx := resolveRecipeIdx(action, item.Recipes)

		return item.Recipes[idx]
	}

	switch action {
	case "commit", "wip":
		return gitutil.GenerateCommitRecipe(item.RepoPath)
	case "discard", "clean":
		return gitutil.GenerateDiscardRecipe(item.RepoPath)
	default:
		return gitutil.GenerateStashRecipe(item.RepoPath)
	}
}

func resolveRecipeIdx(action string, recipes []gitutil.RemediationRecipe) int {
	idx := parseRecipeIndex(action, recipes)
	if idx < 0 || idx >= len(recipes) {
		return 0
	}

	return idx
}

func isFixReleaseTagsRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))
	if first == "release" && len(args) >= 2 && (strings.EqualFold(args[1], "tags") || strings.EqualFold(args[1], "tag")) {
		return true
	}

	return first == "release-tags" || first == "releasetags" || first == "release_tags" || first == "release-tag"
}

func extractFixReleaseTagsArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))
	if first == "release" && len(args) >= 2 && (strings.EqualFold(args[1], "tags") || strings.EqualFold(args[1], "tag")) {
		return args[2:]
	}

	if first == "release-tags" || first == "releasetags" || first == "release_tags" || first == "release-tag" {
		return args[1:]
	}

	return args
}

// RunFix coordinates targeted and batch fix actions.
func RunFix(args []string, aliasOverride string) error {
	if isFixReleaseTagsRequest(args) {
		return cmdfixreleasetags.RunFixReleaseTags(extractFixReleaseTagsArgs(args))
	}

	cmdName := "fix"
	if aliasOverride != "" {
		cmdName = aliasOverride
	}
	checkHelp(cmdName, args)

	if isFixAgyRequest(args) {
		return cmdagy.RunPipelineFixAgyCLI(args)
	}

	if isFixIgnoreRequest(args) {
		return runFixIgnoreDispatch(args)
	}

	if isFixLsRequest(args) {
		return runFixLs(args)
	}

	if isFixCasingRequest(args) {
		return RunFixCasing(args)
	}

	items := cmdremediation.LoadRemediationState()
	if isFixAllRequested(args) {
		action := resolveFixAllAction(args, aliasOverride)
		allItems := collectAllRemediationItems(items)
		return runFixAll(action, allItems)
	}

	if len(items) == 0 {
		return handleEmptyRemediationState(args, aliasOverride)
	}

	if isFixPromptRequested(args) {
		return cmdremediation.RunInteractiveRemediation(items)
	}

	if len(args) == 0 && aliasOverride == "" {
		cmdremediation.PrintRemediationSummary(items)
		return nil
	}

	return executeFixTarget(args, aliasOverride, items)
}

func executeFixTarget(args []string, aliasOverride string, items []cmdremediation.RemediationItem) error {
	item, action, err := resolveFixTarget(args, aliasOverride, items)
	if err != nil {
		return err
	}
	if item == nil {
		return nil
	}

	return applyFixRecipe(item, action)
}

func handleEmptyRemediationState(args []string, aliasOverride string) error {
	if len(args) > 0 && !isFixAllRequested(args) && !isFixPromptRequested(args) {
		return runFixDirect(args, aliasOverride)
	}

	if handleLiveDiscoveredIssues() {
		return nil
	}

	printNoPendingRemediationHelp()
	return nil
}

func handleLiveDiscoveredIssues() bool {
	liveItems := scanTrackedReposForIssues()
	if len(liveItems) == 0 {
		return false
	}

	_ = cmdremediation.SaveRemediationState(liveItems)
	cmdremediation.PrintRemediationSummary(liveItems)
	return true
}

func printNoPendingRemediationHelp() {
	fmt.Printf("%s No pending repositories require remediation.\n  Run 'gitmap fix ls' to inspect repositories with issues.\n  Run 'gitmap pull' to pull all tracked repositories.\n", constants.ColorGreen+"✓"+constants.ColorReset)
}
