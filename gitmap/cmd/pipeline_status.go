package cmd

import (
	"fmt"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

func handlePipelineStatus(args []string) error {
	isJSON := hasArgFlag(args, "--json")
	repo := resolveCurrentRepoSlug()
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
	runs := queryWorkflowRuns(repo)
	etaSeconds := calculateETA(runs)

	isJSON := hasArgFlag(args, "--json")

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
	payload.LastStatus = latest.Status
	payload.LastConclusion = latest.Conclusion
	payload.LastRunURL = latest.URL

	runningCount := countRunningWorkflows(runs)
	payload.PendingPipelines = runningCount
	payload.IsRunning = runningCount > 0

	if payload.IsRunning {
		payload.ActiveWorkflow = latest.Name
		payload.EtaSeconds = calculateETA(runs)
	}

	return payload
}

func countRunningWorkflows(runs []ghRunItem) int {
	runningCount := 0

	for _, r := range runs {
		if r.Status == "in_progress" || r.Status == "queued" || r.Status == "waiting" {
			runningCount++
		}
	}

	return runningCount
}

func renderPipelineStatusTerminal(p PipelineStatusPayload) {
	fmt.Printf("  %s● Repo:%s             %s\n", constants.ColorCyan, constants.ColorReset, p.Repo)

	if p.IsRunning {
		fmt.Printf("  %s● Status:%s           %sRUNNING%s (%s, ETA: %ds)\n",
			constants.ColorCyan, constants.ColorReset,
			constants.ColorYellow, constants.ColorReset,
			p.ActiveWorkflow, p.EtaSeconds)
	} else {
		renderCompletedStatusLine(p)
	}

	fmt.Printf("  %s● Last Tag Release:%s %s\n", constants.ColorCyan, constants.ColorReset, p.LastTagRelease)
	fmt.Printf("  %s● Pending Pipelines:%s %d\n", constants.ColorCyan, constants.ColorReset, p.PendingPipelines)
	fmt.Printf("  %s● Pending PRs:%s       %d\n", constants.ColorCyan, constants.ColorReset, p.PendingPRs)

	if len(p.LastRunURL) > 0 {
		fmt.Printf("  %s● Run URL:%s           %s\n", constants.ColorCyan, constants.ColorReset, p.LastRunURL)
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
	if remainingSeconds < 20 {
		return 20
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
	totalDuration, completedCount := sumCompletedRunDurations(runs, workflowName)
	if completedCount > 0 {
		return totalDuration / completedCount
	}
	return fallbackWorkflowDuration(workflowName)
}

func sumCompletedRunDurations(runs []ghRunItem, workflowName string) (int, int) {
	totalDuration, completedCount := 0, 0
	for _, r := range runs {
		isMatchingCompleted := r.Status == "completed" && (r.Name == workflowName || workflowName == "")
		if !isMatchingCompleted {
			continue
		}
		dur := computeRunDuration(r.CreatedAt, r.UpdatedAt)
		if dur >= 10 {
			totalDuration += dur
			completedCount++
		}
	}
	return totalDuration, completedCount
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
		return 120
	}
	return 75
}
