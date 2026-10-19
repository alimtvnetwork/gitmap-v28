package cmdagent

import (
	"encoding/json"
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"github.com/alimtvnetwork/gitmap-v28/cli/types"
)

func RunTaskStatus(opts TaskStatusOptions) *appfault.AppError {
	trimmed := strings.TrimSpace(opts.TaskId)
	hasRef := len(trimmed) > 0
	if hasRef {
		return runTaskStatusResolved(trimmed, opts.IsJson)
	}
	dbPath, findErr := store.FindTaskDbPath("", "")
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}
	summary, statusErr := store.GetTaskStatus(dbPath, "")
	hasStatusErr := statusErr != nil
	if hasStatusErr {
		return statusErr
	}
	if opts.IsJson {
		return renderStatusJson(summary)
	}
	renderStatusText(summary)

	return nil
}

// runTaskStatusResolved resolves an ID-or-slug reference: slug input goes
// through the Tier 1 lookup and merges subtasks stored under either the
// canonical ID or the slug form; raw IDs keep the legacy direct path.
func runTaskStatusResolved(taskIdOrSlug string, isJson bool) *appfault.AppError {
	masterDb := store.ResolveMasterAgentDbPath("")
	task, findErr := store.FindParentTaskBySlug(masterDb, taskIdOrSlug)
	if findErr != nil {
		return findErr
	}
	if task == nil {
		return runTaskStatusByPath(taskIdOrSlug, isJson)
	}

	summary, sumErr := loadTaskSummary(task)
	if sumErr != nil {
		return sumErr
	}
	if isJson {
		return renderStatusJson(summary)
	}
	renderStatusText(summary)

	return nil
}

func runTaskStatusByPath(taskIdOrSlug string, isJson bool) *appfault.AppError {
	dbPath, dbErr := store.FindTaskDbPath("", taskIdOrSlug)
	if dbErr != nil {
		return dbErr
	}
	summary, statusErr := store.GetTaskStatus(dbPath, taskIdOrSlug)
	if statusErr != nil {
		return statusErr
	}
	if isJson {
		return renderStatusJson(summary)
	}
	renderStatusText(summary)

	return nil
}

func renderStatusJson(summary *types.TaskStatusSummary) *appfault.AppError {
	b, err := json.MarshalIndent(summary, "", "  ")
	if err != nil {
		return appfault.WrapSimple(err, "renderStatusJson")
	}
	fmt.Println(string(b))

	return nil
}

func renderStatusText(s *types.TaskStatusSummary) {
	hasParent := len(s.ParentTask.ParentTaskId) > 0
	if !hasParent {
		fmt.Println("Parent task not found or not initialized.")
		return
	}
	printStatusHeader(s)
	printSubtaskRows(s.Subtasks)
}

func printStatusHeader(s *types.TaskStatusSummary) {
	p := s.ParentTask
	fmt.Printf("Task ID:    %s\n", p.ParentTaskId)
	fmt.Printf("Slug:       %s\n", p.TaskSlug)
	fmt.Printf("Status:     %s\n", p.Status)
	fmt.Printf("Progress:   %d/%d steps (Completed: %t)\n", p.CompletedSteps, p.TotalStepsBudget, s.IsCompleted)
	fmt.Printf("Subtasks:   %d total (%d pending, %d in progress, %d done, %d failed)\n\n",
		len(s.Subtasks), s.PendingCount, s.InProgressCount, s.CompletedCount, s.FailedCount)
}

func printSubtaskRows(subtasks []types.Subtask) {
	hasSub := len(subtasks) > 0
	if !hasSub {
		fmt.Println("No subtasks registered.")
		return
	}
	fmt.Printf("%-14s  %-12s  %-15s  %s\n", "CODE", "STATUS", "AGENT", "TITLE")
	for _, sub := range subtasks {
		agent := sub.AssignedAgentRole
		hasAgent := len(agent) > 0
		if !hasAgent {
			agent = "-"
		}
		fmt.Printf("%-14s  %-12s  %-15s  %s\n", sub.TaskCode, sub.Status, agent, sub.Title)
	}
}
