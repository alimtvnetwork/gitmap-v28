package cmd

import (
	"encoding/json"
	"fmt"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type ghJobsResponse struct {
	Jobs []ghJobItem `json:"jobs"`
}

type ghJobItem struct {
	DatabaseID  int64        `json:"databaseId"`
	Name        string       `json:"name"`
	Status      string       `json:"status"`
	Conclusion  string       `json:"conclusion"`
	StartedAt   string       `json:"startedAt"`
	CompletedAt string       `json:"completedAt"`
	Steps       []ghStepItem `json:"steps"`
}

type ghStepItem struct {
	Number      int    `json:"number"`
	Name        string `json:"name"`
	Status      string `json:"status"`
	Conclusion  string `json:"conclusion"`
	StartedAt   string `json:"startedAt"`
	CompletedAt string `json:"completedAt"`
}

// queryRunJobs fetches job and segment/step details for a GitHub Actions run.
func queryRunJobs(repo string, runID int64) []ghJobItem {
	if runID <= 0 {
		return nil
	}
	args := []string{"run", "view", fmt.Sprintf("%d", runID), "--json", "jobs"}
	if repo != "" {
		args = append(args, "--repo", repo)
	}
	cmd := exec.Command("gh", args...)
	out, err := cmd.Output()
	if err != nil || len(out) == 0 {
		return nil
	}
	var resp ghJobsResponse
	if err := json.Unmarshal(out, &resp); err != nil {
		return nil
	}
	return resp.Jobs
}

// renderSegmentBreakdown renders the progress of segments/steps for active runs.
func renderSegmentBreakdown(repo string, runID int64) {
	jobs := queryRunJobs(repo, runID)
	if len(jobs) == 0 {
		return
	}
	fmt.Printf("\n  %s● Pipeline Segments:%s\n", constants.ColorCyan, constants.ColorReset)
	for _, j := range jobs {
		if len(j.Steps) == 0 {
			printJobStatusLine(j.Name, j.Status, j.Conclusion)
			continue
		}
		fmt.Printf("    %sJob: %s%s\n", constants.ColorWhite, j.Name, constants.ColorReset)
		for _, s := range j.Steps {
			if strings.HasPrefix(s.Name, "Post ") {
				continue
			}
			renderStepProgressLine(s)
		}
	}
}

func printJobStatusLine(name, status, conclusion string) {
	badge := formatSegmentBadge(status, conclusion)
	fmt.Printf("    • %-32s %s\n", name, badge)
}

func renderStepProgressLine(s ghStepItem) {
	badge := formatSegmentBadge(s.Status, s.Conclusion)
	dur := computeStepDurationString(s.StartedAt, s.CompletedAt, s.Status)
	fmt.Printf("      • %-34s %s %s\n", s.Name, badge, dur)
}

func formatSegmentBadge(status, conclusion string) string {
	switch {
	case status == "completed" && conclusion == "success":
		return constants.ColorGreen + "✔ completed" + constants.ColorReset
	case status == "completed" && conclusion == "failure":
		return constants.ColorRed + "✖ failed   " + constants.ColorReset
	case status == "in_progress":
		return constants.ColorYellow + "● running  " + constants.ColorReset
	case status == "pending" || status == "queued":
		return constants.ColorDim + "○ queued   " + constants.ColorReset
	default:
		return status
	}
}

func computeStepDurationString(startedStr, completedStr, status string) string {
	if status == "in_progress" {
		return formatInProgressDuration(startedStr)
	}
	if startedStr == "" || completedStr == "" || startedStr == "0001-01-01T00:00:00Z" {
		return ""
	}
	t1, err1 := time.Parse(time.RFC3339, startedStr)
	t2, err2 := time.Parse(time.RFC3339, completedStr)
	if err1 != nil || err2 != nil || t2.Before(t1) {
		return ""
	}
	return fmt.Sprintf("(%ds)", int(t2.Sub(t1).Seconds()))
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

func formatInProgressDuration(startedStr string) string {
	t, err := time.Parse(time.RFC3339, startedStr)
	if err != nil {
		return ""
	}
	return fmt.Sprintf("(%ds elapsed)", int(time.Since(t).Seconds()))
}
