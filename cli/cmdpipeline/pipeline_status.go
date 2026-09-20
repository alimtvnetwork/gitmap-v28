package cmdpipeline

import (
	"fmt"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

func handlePipelineStatus(args []string) error {
	isJSON := hasArgFlag(args, "--json")
	hasTimeout := hasArgFlag(args, "-t") || hasArgFlag(args, "--timeout") || hasArgFlag(args, "--timeline")
	repo := resolveCurrentRepoSlug()

	if hasTimeout {
		return runPipelineDynamicTimeline(repo, isJSON)
	}

	runs := queryWorkflowRuns(repo)
	pendingPRs := queryPendingPRs(repo)
	lastTag := queryLatestTagRelease(repo)

	payload := buildStatusPayload(repo, lastTag, pendingPRs, runs)
	recordPipelineInDB(payload, runs)

	if isJSON {
		return printJSON(payload)
	}

	renderPipelineStatusTerminal(payload)

	return nil
}

func handlePipelineWaitTime(args []string) error {
	repo := resolveCurrentRepoSlug()
	hasTimeout := hasArgFlag(args, "-t") || hasArgFlag(args, "--timeout") || hasArgFlag(args, "--timeline")
	isJSON := hasArgFlag(args, "--json")

	if hasTimeout {
		return runPipelineDynamicTimeline(repo, isJSON)
	}

	runs := queryWorkflowRuns(repo)
	etaSeconds := calculateETA(runs)

	if isJSON {
		out := map[string]any{"etaSeconds": etaSeconds, "repo": repo}

		return printJSON(out)
	}

	if etaSeconds > 0 {
		fmt.Printf("%d\n", etaSeconds)
	} else {
		fmt.Println("0")
	}

	return nil
}

func buildStatusPayload(repo, lastTag string, pendingPRs int, runs []ghRunItem) PipelineStatusPayload {
	nowStr := time.Now().UTC().Format(time.RFC3339)
	payload := PipelineStatusPayload{
		Repo:           repo,
		LastTagRelease: lastTag,
		PendingPRs:     pendingPRs,
		UpdatedAt:      nowStr,
	}

	if len(runs) == 0 {
		return payload
	}

	latest := runs[0]
	payload.LastRunId = latest.DatabaseId
	payload.LastStatus = latest.Status
	payload.LastConclusion = latest.Conclusion
	payload.LastRunUrl = latest.Url

	runningCount := countRunningWorkflows(runs)
	payload.PendingPipelines = runningCount
	payload.IsRunning = runningCount > 0

	if payload.IsRunning {
		populateActiveRunPayload(&payload, runs, latest)
	}

	AttachLivePipelineErrors(&payload)

	return payload
}

func populateActiveRunPayload(payload *PipelineStatusPayload, runs []ghRunItem, latest ghRunItem) {
	activeRun := findActiveWorkflowRun(runs)
	if activeRun != nil {
		payload.ActiveWorkflow = activeRun.Name
		payload.LastRunId = activeRun.DatabaseId
		payload.LastRunUrl = activeRun.Url
		payload.EtaSeconds = calculateETA(runs)

		return
	}

	payload.ActiveWorkflow = latest.Name
	payload.EtaSeconds = calculateETA(runs)
}

func renderPipelineStatusTerminal(p PipelineStatusPayload) {
	fmt.Printf("  %s● Repo:%s             %s\n", constants.ColorCyan, constants.ColorReset, p.Repo)

	if p.IsRunning {
		fmt.Printf("  %s● Status:%s           %sRUNNING%s (%s, ETA: %s)\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorYellow, constants.ColorReset,
			p.ActiveWorkflow, formatEtaDisplay(p.EtaSeconds))
	} else {
		renderCompletedStatusLine(p)
	}

	fmt.Printf("  %s● Last Tag Release:%s %s\n", constants.ColorCyan, constants.ColorReset, p.LastTagRelease)
	fmt.Printf("  %s● Pending Pipelines:%s %d\n", constants.ColorCyan, constants.ColorReset, p.PendingPipelines)
	fmt.Printf("  %s● Pending PRs:%s       %d\n", constants.ColorCyan, constants.ColorReset, p.PendingPRs)

	if len(p.LastRunUrl) > 0 {
		fmt.Printf("  %s● Run URL:%s           %s\n", constants.ColorCyan, constants.ColorReset, p.LastRunUrl)
	}

	if p.IsRunning && p.LastRunId > 0 {
		renderSegmentBreakdown(p.Repo, p.LastRunId)
	}
}

func renderCompletedStatusLine(p PipelineStatusPayload) {
	statusColor := constants.ColorGreen

	if p.LastConclusion == "failure" {
		statusColor = constants.ColorRed
	}

	fmt.Printf("  %s● Status:%s           %s%s%s (conclusion: %s)\n",
		constants.ColorCyan, constants.ColorReset,
		statusColor, p.LastStatus, constants.ColorReset,
		p.LastConclusion)

	if p.LastConclusion == "failure" && p.Repo != "" {
		renderFailureErrorSummary(p.Repo, 0)
	}
}

func calculateETA(runs []ghRunItem) int {
	activeRun := findActiveWorkflowRun(runs)
	if activeRun == nil {
		return 0
	}

	createdAt, parseErr := time.Parse(time.RFC3339, activeRun.CreatedAt)
	if parseErr != nil {
		return 30
	}

	avgDuration := calculateAverageDuration(runs, activeRun.Name)
	elapsedSeconds := int(time.Since(createdAt).Seconds())
	remainingSeconds := avgDuration - elapsedSeconds
	if remainingSeconds < 15 {
		return 15
	}

	return remainingSeconds
}

func findActiveWorkflowRun(runs []ghRunItem) *ghRunItem {
	for i := range runs {
		status := runs[i].Status
		if status == "in_progress" || status == "queued" || status == "waiting" {
			return &runs[i]
		}
	}

	return nil
}

func calculateAverageDuration(runs []ghRunItem, workflowName string) int {
	durs := collectValidRunDurations(runs, workflowName)
	if len(durs) == 0 {
		durs = queryDbHistoricalDurations(workflowName)
	}
	if len(durs) > 0 {
		return computeBaselineDuration(durs)
	}

	return fallbackWorkflowDuration(workflowName)
}

func collectValidRunDurations(runs []ghRunItem, workflowName string) []int {
	var durs []int
	for _, r := range runs {
		if !isRunValidSuccess(r, workflowName) {
			continue
		}
		dur := computeRunDuration(r.CreatedAt, r.UpdatedAt)
		if dur >= getWorkflowMinDurationFloor(workflowName) {
			durs = append(durs, dur)
		}
	}

	return durs
}

func isRunValidSuccess(r ghRunItem, workflowName string) bool {
	if r.Status != "completed" || r.Conclusion != "success" {
		return false
	}
	if workflowName == "" || r.Name == workflowName {
		return true
	}

	return strings.EqualFold(r.Name, workflowName)
}

func getWorkflowMinDurationFloor(workflowName string) int {
	lower := strings.ToLower(workflowName)
	if strings.Contains(lower, "lint") || strings.Contains(lower, "format") {
		return 15
	}
	if strings.Contains(lower, "release") {
		return 45
	}

	return 45
}

func computeBaselineDuration(durs []int) int {
	if len(durs) == 1 {
		return durs[0]
	}
	sort.Ints(durs)
	mid := len(durs) / 2
	if len(durs)%2 == 1 {
		return durs[mid]
	}

	return (durs[mid-1] + durs[mid]) / 2
}

func queryDbHistoricalDurations(workflowName string) []int {
	slug := resolveCurrentRepoSlug()
	db, err := pipelinedb.OpenPipelineSplitDb(slug)
	if err != nil {
		return nil
	}
	defer db.Close()

	return db.QuerySuccessfulRunDurations(workflowName, 20)
}
