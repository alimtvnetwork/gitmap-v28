package cmdfix

import (
	"fmt"
	"strings"
	"sync"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagy"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdconflict"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdfixreleasetags"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdremediation"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/parallel"
)

func isFixAgyRequest(args []string) bool {
	for _, arg := range args {
		low := strings.ToLower(arg)
		if low == "agy" || low == "aef" || low == "agy-errors-fix" {
			return true
		}
	}
	return false
}

func isFixAllRequested(args []string) bool {
	for _, arg := range args {
		low := strings.ToLower(arg)
		if low == "all" || low == "--all" || low == "-all" || low == "-a" {
			return true
		}
	}
	return false
}

func isFixPromptRequested(args []string) bool {
	for _, arg := range args {
		low := strings.ToLower(arg)
		if low == "--prompt" || low == "-p" || low == "prompt" || low == "--interactive" || low == "-i" || low == "interactive" {
			return true
		}
	}
	return false
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
	hasPullFailures := cmdpull.HasPendingPullFailures()
	if len(items) == 0 && hasPullFailures {
		return cmdpull.RunBatchPullFix([]string{"--all"})
	}

	if len(items) == 0 {
		fmt.Printf("%s No pending repositories require remediation.\n", constants.ColorGreen+"✓"+constants.ColorReset)
		return nil
	}

	concurrency := parallel.ResolveConcurrency(0)
	fmt.Printf("%s Remediating %d repository(ies) in parallel (%d workers) with action: %s\n\n",
		constants.ColorCyan+"ℹ"+constants.ColorReset, len(items), concurrency, action)

	var (
		mu          sync.Mutex
		failedRepos []string
		completed   int
		total       = len(items)
	)

	_ = parallel.ExecuteParallel(items, concurrency, func(workerID int, item cmdremediation.RemediationItem) error {
		mu.Lock()
		fmt.Printf("%s [Worker %d] Starting remediation on: %s\n",
			constants.ColorCyan+"ℹ"+constants.ColorReset, workerID+1, item.RepoName)
		mu.Unlock()

		idx := parseRecipeIndex(action, item.Recipes)
		if idx < 0 || idx >= len(item.Recipes) {
			idx = 0
		}

		err := executeFixRecipe(&item, item.Recipes[idx])
		reportWorkerCompletion(&mu, &completed, total, workerID, item.RepoName, err, &failedRepos)
		return err
	})

	if hasPullFailures {
		fmt.Println()
		_ = cmdpull.RunBatchPullFix([]string{"--all"})
	}

	if len(failedRepos) > 0 {
		return fmt.Errorf("remediation failed on %d repository(ies): %s", len(failedRepos), strings.Join(failedRepos, ", "))
	}

	fmt.Printf("\n%s All %d repository(ies) remediated successfully.\n",
		constants.ColorGreen+"✓"+constants.ColorReset, len(items))
	return nil
}

func reportWorkerCompletion(mu *sync.Mutex, completed *int, total int, workerID int, repoName string, err error, failedRepos *[]string) {
	mu.Lock()
	defer mu.Unlock()
	*completed++
	if err != nil {
		*failedRepos = append(*failedRepos, repoName)
		fmt.Printf("%s [Worker %d] Failed (%d/%d): %s - %v\n",
			constants.ColorRed+"✗"+constants.ColorReset, workerID+1, *completed, total, repoName, err)
		return
	}
	fmt.Printf("%s [Worker %d] Completed (%d/%d): %s\n",
		constants.ColorGreen+"✓"+constants.ColorReset, workerID+1, *completed, total, repoName)
}

func isFixReleaseTagsRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))
	isReleaseTagsTwoWords := first == "release" && len(args) >= 2 && (strings.EqualFold(args[1], "tags") || strings.EqualFold(args[1], "tag"))
	if isReleaseTagsTwoWords {
		return true
	}

	return first == "release-tags" || first == "releasetags" || first == "release_tags" || first == "release-tag"
}

func extractFixReleaseTagsArgs(args []string) []string {
	if len(args) == 0 {
		return args
	}

	first := strings.ToLower(strings.TrimSpace(args[0]))
	isReleaseTagsTwoWords := first == "release" && len(args) >= 2 && (strings.EqualFold(args[1], "tags") || strings.EqualFold(args[1], "tag"))
	if isReleaseTagsTwoWords {
		remaining := make([]string, 0, len(args)-2)
		remaining = append(remaining, args[2:]...)
		return remaining
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
	fmt.Printf("%s No pending repositories require remediation.\n", constants.ColorGreen+"✓"+constants.ColorReset)
	fmt.Println("  Run 'gitmap fix ls' to inspect repositories with issues.")
	fmt.Println("  Run 'gitmap pull' to pull all tracked repositories.")
}
