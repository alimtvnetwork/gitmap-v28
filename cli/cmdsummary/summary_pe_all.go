// Package cmdsummary implements high-speed 48h-filtered CI/CD error diagnosis across repositories.
package cmdsummary

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
	"github.com/pterm/pterm"
)

// PipelineAllItem models pipeline health state for one repository in pe all.
type PipelineAllItem struct {
	RepoSlug     string             `json:"repoSlug"`
	RepoPath     string             `json:"repoPath"`
	IsActive     bool               `json:"isActive"`
	Status       string             `json:"status"` // "clean", "failed", "no-data"
	Pipeline     *RepoPipelineError `json:"pipeline,omitempty"`
	LastCommitAt string             `json:"lastCommitAt"`
}

// PipelineAllPayload models the JSON envelope for pe all.
type PipelineAllPayload struct {
	Attributes struct {
		GeneratedAt   string `json:"generatedAt"`
		TotalScanned  int    `json:"totalScanned"`
		ActiveScanned int    `json:"activeScanned"`
		ForceAll      bool   `json:"forceAll"`
	} `json:"attributes"`
	Data struct {
		CleanCount  int               `json:"cleanCount"`
		FailedCount int               `json:"failedCount"`
		Items       []PipelineAllItem `json:"items"`
	} `json:"data"`
}

// RunPipelineErrorsAll executes workspace-wide CI/CD pipeline error discovery with 48h filtering.
func RunPipelineErrorsAll(args []string) error {
	isJSON := false
	forceAll := false

	for _, arg := range args {
		lower := strings.ToLower(arg)
		if lower == "--json" || lower == "-j" {
			isJSON = true
		} else if lower == "--force-all" || lower == "--all" || lower == "-a" {
			forceAll = true
		}
	}

	cwd, _ := os.Getwd()
	records := resolveWorkspaceRepositories(cwd)
	if len(records) == 0 {
		fmt.Println("[GitMap] No repositories discovered in workspace.")
		return nil
	}

	cutoff := time.Now().Add(-48 * time.Hour)
	var items []PipelineAllItem
	cleanCount := 0
	failedCount := 0

	for _, rec := range records {
		dir := rec.AbsolutePath
		if dir == "" {
			dir = rec.RelativePath
		}

		lastCommitEpoch, _ := gitExec(dir, "log", "-1", "--format=%ct")
		var lastCommitAt time.Time
		if epochSec, err := strconv.ParseInt(lastCommitEpoch, 10, 64); err == nil && epochSec > 0 {
			lastCommitAt = time.Unix(epochSec, 0)
		} else {
			lastCommitAt = time.Now()
		}

		isActive := lastCommitAt.After(cutoff)
		if !forceAll && !isActive {
			continue
		}

		slug := pipelinedb.SanitizeRepoSlug(rec.Slug)
		if slug == "" {
			slug = filepath.Base(dir)
		}

		item := inspectPipelineForItem(slug, dir, isActive, lastCommitAt)
		if item.Status == "clean" {
			cleanCount++
		} else if item.Status == "failed" {
			failedCount++
		}

		items = append(items, item)
	}

	var payload PipelineAllPayload
	payload.Attributes.GeneratedAt = time.Now().UTC().Format(time.RFC3339)
	payload.Attributes.TotalScanned = len(records)
	payload.Attributes.ActiveScanned = len(items)
	payload.Attributes.ForceAll = forceAll
	payload.Data.CleanCount = cleanCount
	payload.Data.FailedCount = failedCount
	payload.Data.Items = items

	if isJSON {
		enc := json.NewEncoder(os.Stdout)
		enc.SetIndent("", "  ")
		return enc.Encode(payload)
	}

	renderPipelineAllTerminal(payload)
	return nil
}

func inspectPipelineForItem(slug, dir string, isActive bool, lastCommit time.Time) PipelineAllItem {
	item := PipelineAllItem{
		RepoSlug:     slug,
		RepoPath:     dir,
		IsActive:     isActive,
		Status:       "no-data",
		LastCommitAt: lastCommit.Format("2006-01-02 15:04"),
	}

	dbPath := pipelinedb.ResolvePipelineDbPath(slug)
	if _, statErr := os.Stat(dbPath); statErr != nil {
		return item
	}

	db, errOpen := pipelinedb.OpenPipelineSplitDb(slug)
	if errOpen != nil {
		return item
	}
	defer db.Close()

	run, runErr := db.QueryRunByNegativeOffset(-1)
	if runErr != nil || run == nil {
		return item
	}

	if run.IsSuccess || strings.EqualFold(run.Conclusion, "success") {
		item.Status = "clean"
		return item
	}

	item.Status = "failed"
	item.Pipeline = &RepoPipelineError{
		RepoSlug:     slug,
		WorkflowName: run.WorkflowName,
		JobName:      "pipeline-job",
		StepName:     "failing-step",
		ExitCode:     1,
		ErrorSummary: "Pipeline failure detected in latest run",
		RunURL:       run.RunUrl,
	}

	compactRes := db.QueryCompactErrorLogsByRunId(run.RunId)
	if compactRes.IsSuccess() && len(compactRes.Data) > 0 {
		first := compactRes.Data[0]
		item.Pipeline.ErrorSummary = first.ErrorText
		item.Pipeline.StepName = first.StepName
		item.Pipeline.JobName = "workflow-job"
		item.Pipeline.ExitCode = 1

		var traceLines []string
		for _, logRec := range compactRes.Data {
			if logRec.CompactLogs != "" {
				for _, l := range strings.Split(logRec.CompactLogs, "\n") {
					trimmed := strings.TrimRight(l, "\r")
					if trimmed != "" {
						traceLines = append(traceLines, trimmed)
					}
				}
			} else if logRec.ErrorText != "" {
				traceLines = append(traceLines, strings.TrimRight(logRec.ErrorText, "\r"))
			}
		}
		if len(traceLines) > 25 {
			traceLines = traceLines[len(traceLines)-25:]
		}
		item.Pipeline.StackTrace = strings.Join(traceLines, "\n")
	}

	return item
}

func renderPipelineAllTerminal(payload PipelineAllPayload) {
	fmt.Printf("[GitMap] CI/CD Pipeline Health Diagnostics (Active Filter: %d Repositories)\n", payload.Attributes.ActiveScanned)
	if !payload.Attributes.ForceAll {
		fmt.Println("Filtered by last 48 hours activity window. Use --force-all to inspect all repositories.")
	}
	fmt.Println()

	for _, item := range payload.Data.Items {
		switch item.Status {
		case "clean":
			fmt.Printf("%s [%s] green (all checks passed)\n", pterm.Green("✓"), item.RepoSlug)
		case "failed":
			fmt.Printf("%s [%s] FAILED\n", pterm.Red("❌"), pterm.Bold.Sprint(item.RepoSlug))
			if item.Pipeline != nil {
				fmt.Printf("   Workflow: %s | Step: %s | Exit Code: %d\n", item.Pipeline.WorkflowName, item.Pipeline.StepName, item.Pipeline.ExitCode)
				if item.Pipeline.ErrorSummary != "" {
					fmt.Printf("   Error: %s\n", item.Pipeline.ErrorSummary)
				}
				if item.Pipeline.StackTrace != "" {
					fmt.Println("   Stack Trace (last 25 lines):")
					for _, line := range strings.Split(item.Pipeline.StackTrace, "\n") {
						if strings.TrimSpace(line) != "" {
							fmt.Printf("     %s\n", line)
						}
					}
				}
				if item.Pipeline.RunURL != "" {
					fmt.Printf("   Run URL: %s\n", item.Pipeline.RunURL)
				}
			}
			fmt.Println()
		}
	}

	fmt.Printf("\nSummary: %d clean, %d failed out of %d active repos.\n",
		payload.Data.CleanCount, payload.Data.FailedCount, payload.Attributes.ActiveScanned)
}
