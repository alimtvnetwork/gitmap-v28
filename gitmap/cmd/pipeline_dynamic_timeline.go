package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
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

func watchDynamicTimeline(repo, workflowName string, isJSON bool) error {
	fmt.Printf("%s● Watching pipeline [%s] dynamic timeline...%s\n", constants.ColorCyan, workflowName, constants.ColorReset)
	startTime := time.Now()
	for {
		runs := queryWorkflowRuns(repo)
		active := findActiveWorkflowRun(runs)
		if active == nil && len(runs) > 0 {
			return reportCompletedTimeline(runs[0], repo, isJSON)
		}
		if active == nil {
			break
		}
		eta := calculateETA(runs)
		elapsed := int(time.Since(startTime).Seconds())
		printTimelineProgress(active.Name, eta, elapsed)
		interval := computeAdaptiveInterval(eta)
		time.Sleep(time.Duration(interval) * time.Second)
	}
	return nil
}

func computeAdaptiveInterval(eta int) int {
	if eta > 120 {
		return 15
	}
	if eta > 60 {
		return 10
	}
	return 5
}

func printTimelineProgress(name string, eta, elapsed int) {
	fmt.Printf("  %s⏳ [%s] in progress:%s ETA ~%ds (elapsed: %ds)\n",
		constants.ColorYellow, name, constants.ColorReset, eta, elapsed)
}

// ErrorLogsTimelineParams specifies parameters for running errorlogs with dynamic timeline watching.
type ErrorLogsTimelineParams struct {
	Repo         string
	IsJSON       bool
	WantFix      bool
	WantCheck    bool
	FilePath     string
	TempFileName string
	Args         []string
}

// runPipelineErrorLogsDynamicTimeline watches the pipeline dynamic timeline and surfaces error logs upon completion.
func runPipelineErrorLogsDynamicTimeline(params ErrorLogsTimelineParams) error {
	runs := queryWorkflowRuns(params.Repo)
	active := findActiveWorkflowRun(runs)
	if active != nil {
		watchDynamicTimeline(params.Repo, active.Name, params.IsJSON)
		runs = queryWorkflowRuns(params.Repo)
	}

	payload := buildErrorLogsPayload(params.Repo, runs)
	if len(runs) > 0 {
		payload.RerunEtaSeconds = calculateAverageDuration(runs, payload.WorkflowName)
	}

	if params.WantFix || params.WantCheck {
		payload.CICDChecks = runInternalCICDChecks(params.WantFix)
	}

	err := writeOrRenderErrorLogs(ErrorLogOutputParams{
		Payload:  payload,
		IsJSON:   params.IsJSON,
		FilePath: params.FilePath,
		TempFile: params.TempFileName,
	})
	if err != nil {
		return err
	}

	maybeOfferAutoFix(payload, params.IsJSON, params.WantFix, params.WantCheck, params.FilePath, params.TempFileName, params.Args)

	return nil
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

func locateFailedRunID(runs []ghRunItem) int64 {
	for _, r := range runs {
		if r.Conclusion == "failure" {
			return r.DatabaseId
		}
	}
	return 0
}

func renderFailureErrorSummary(repo string, runID int64) {
	if runID <= 0 {
		runID = locateFailedRunID(queryWorkflowRuns(repo))
	}
	if runID <= 0 {
		return
	}
	rawLogs := queryFailedRunLogs(repo, runID)
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
	lowerName := strings.ToLower(workflowName)
	if strings.Contains(lowerName, "release") {
		return 95
	}
	if strings.Contains(lowerName, "ci") {
		return 180
	}
	return 90
}
