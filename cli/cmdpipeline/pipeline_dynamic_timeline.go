package cmdpipeline

import (
	"flag"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runPipelineDynamicTimeline runs an adaptive polling watch loop driven by ETA.
func runPipelineDynamicTimeline(repo string, isJSON bool) error {
	runs := queryWorkflowRuns(repo)
	if len(runs) == 0 {
		fmt.Printf("No recent pipeline runs found for %s.\n", repo)

		return nil
	}

	active := findActiveWorkflowRun(runs)
	if active == nil {
		return reportCompletedTimeline(runs[0], repo, isJSON)
	}

	return watchDynamicTimeline(repo, active.Name, isJSON)
}

func isTimelineTestMode() bool {
	if os.Getenv("CI") != "" {
		return true
	}

	if os.Getenv("GITMAP_TEST") != "" {
		return true
	}

	return flag.Lookup("test.v") != nil
}

func watchDynamicTimeline(repo, workflowName string, isJSON bool) error {
	fmt.Printf("%s● Watching pipeline [%s] dynamic timeline...%s\n", constants.ColorCyan, workflowName, constants.ColorReset)
	startTime := time.Now()
	for {
		isDone, err := pollTimelineStep(repo, isJSON, startTime)
		if err != nil || isDone {
			return err
		}

		if isTimelineTestMode() {
			break
		}
	}

	return nil
}

func pollTimelineStep(repo string, isJSON bool, startTime time.Time) (bool, error) {
	runs := queryWorkflowRuns(repo)
	active := findActiveWorkflowRun(runs)
	if active == nil && len(runs) > 0 {
		return true, reportCompletedTimeline(runs[0], repo, isJSON)
	}

	if active == nil {
		return true, nil
	}

	eta := calculateETA(runs)
	elapsed := int(time.Since(startTime).Seconds())
	printTimelineProgress(active.Name, eta, elapsed)
	if !isTimelineTestMode() {
		time.Sleep(time.Duration(computeAdaptiveInterval(eta)) * time.Second)
	}

	return false, nil
}

func computeAdaptiveInterval(eta int) int {
	if eta <= 0 {
		return 15
	}
	if eta > 120 {
		return 15
	}
	if eta > 60 {
		return 10
	}

	return 5
}

func printTimelineProgress(name string, eta, elapsed int) {
	fmt.Printf("  %s⏳ [%s] in progress:%s ETA %s (elapsed: %s)\n",
		constants.ColorYellow, name, constants.ColorReset, formatEtaDisplay(eta), formatDurationSeconds(elapsed))
}

func reportCompletedTimeline(latest ghRunItem, repo string, isJSON bool) error {
	if latest.Conclusion == "success" {
		fmt.Printf("\n%s✓ Pipeline [%s] completed successfully!%s\n", constants.ColorGreen, latest.Name, constants.ColorReset)

		return nil
	}

	fmt.Printf("\n%s✖ Pipeline [%s] finished with conclusion: %s%s\n", constants.ColorRed, latest.Name, latest.Conclusion, constants.ColorReset)
	renderFailureErrorSummary(repo, latest.DatabaseId)
	runs := queryWorkflowRuns(repo)
	if len(runs) > 0 {
		rerunETA := calculateAverageDuration(runs, latest.Name)
		printRerunETA(rerunETA)
	}

	return nil
}

func locateFailedRunId(runs []ghRunItem) uint64 {
	for _, r := range runs {
		if r.Conclusion == "failure" {
			return r.DatabaseId
		}
	}

	return 0
}

func renderFailureErrorSummary(repo string, runId uint64) {
	if runId == 0 {
		runId = locateFailedRunId(queryWorkflowRuns(repo))
	}

	if runId == 0 {
		return
	}

	rawLogs := queryFailedRunLogs(repo, runId)
	cleanErrors := extractCleanErrorLines(rawLogs)
	if cleanErrors == "" {
		return
	}

	fmt.Printf("\n  %s● CI/CD Error Diagnostics:%s\n", constants.ColorRed, constants.ColorReset)
	fmt.Printf("    %s\n", strings.Repeat("─", 78))
	for _, l := range strings.Split(cleanErrors, "\n") {
		fmt.Printf("    %s\n", l)
	}

	fmt.Printf("    %s\n\n", strings.Repeat("─", 78))
}

func computeRunDuration(createdStr, updatedStr string) int {
	createdAt, err1 := time.Parse(time.RFC3339, createdStr)
	updatedAt, err2 := time.Parse(time.RFC3339, updatedStr)
	if err1 != nil || err2 != nil || updatedAt.Before(createdAt) {
		return 0
	}

	return int(updatedAt.Sub(createdAt).Seconds())
}

func fallbackWorkflowDuration(workflowName string) int {
	repo := resolveCurrentRepoSlug()
	cached := GetCachedWorkflowETA(repo, workflowName)
	if cached > 0 {
		return cached
	}

	return resolveStaticFallbackDuration(workflowName)
}

func resolveStaticFallbackDuration(workflowName string) int {
	lowerName := strings.ToLower(workflowName)
	dur := resolveSpecializedWorkflowDuration(lowerName)
	if dur > 0 {
		return dur
	}

	return resolveStandardWorkflowDuration(lowerName)
}

func resolveSpecializedWorkflowDuration(lower string) int {
	if strings.Contains(lower, "beacon") {
		return 55
	}
	if strings.Contains(lower, "rewrite") || strings.Contains(lower, "smoke") {
		return 75
	}
	if strings.Contains(lower, "race") {
		return 130
	}

	return 0
}

func resolveStandardWorkflowDuration(lower string) int {
	if strings.Contains(lower, "release") {
		return 480
	}
	if strings.Contains(lower, "cross-platform") || strings.Contains(lower, "build") {
		return 420
	}
	if strings.Contains(lower, "ci") || strings.Contains(lower, "test") {
		return 450
	}

	return 300
}
