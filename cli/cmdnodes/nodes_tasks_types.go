package cmdnodes

import (
	"strconv"
	"strings"
)

// TaskItemSummary represents an enqueued or running agent subtask.
type TaskItemSummary struct {
	TaskID       int64  `json:"taskId"`
	TaskSlug     string `json:"taskSlug"`
	TaskCode     string `json:"taskCode"`
	Title        string `json:"title"`
	AssignedRole string `json:"assignedRole"`
	Status       string `json:"status"`
	IsHeavyDuty  bool   `json:"isHeavyDuty"`
	CurrentStep  int    `json:"currentStep"`
	TotalSteps   int    `json:"totalSteps"`
	CreatedAt    string `json:"createdAt"`
}

// NodeTaskReport summarizes active tasks and queue depth for a single fleet node.
type NodeTaskReport struct {
	NodeAlias        string            `json:"nodeAlias"`
	NodeHost         string            `json:"nodeHost"`
	NodePort         int               `json:"nodePort"`
	IsLocal          bool              `json:"isLocal"`
	Status           string            `json:"status"`
	ActiveTaskCount  int               `json:"activeTaskCount"`
	PendingTaskCount int               `json:"pendingTaskCount"`
	HeavyTaskCount   int               `json:"heavyTaskCount"`
	ActiveTasks      []TaskItemSummary `json:"activeTasks"`
	PendingTasks     []TaskItemSummary `json:"pendingTasks"`
	ErrorMessage     string            `json:"errorMessage,omitempty"`
}

// FleetTasksPayload encapsulates the complete response across all cluster nodes.
type FleetTasksPayload struct {
	Timestamp      string           `json:"timestamp"`
	TotalNodes     int              `json:"totalNodes"`
	TotalActive    int              `json:"totalActive"`
	TotalPending   int              `json:"totalPending"`
	TotalHeavyDuty int              `json:"totalHeavyDuty"`
	NodeReports    []NodeTaskReport `json:"nodeReports"`
}

type nodesTasksOptions struct {
	isJSON       bool
	limit        int
	targetFilter string
}

func parseNodesTasksOptions(args []string) nodesTasksOptions {
	opts := nodesTasksOptions{
		limit: 10,
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		if a == "--json" || a == "-j" {
			opts.isJSON = true
		} else if (a == "--limit" || a == "-l") && i+1 < len(args) {
			i++
			if n, err := strconv.Atoi(args[i]); err == nil && n > 0 {
				opts.limit = n
			}
		} else if isTaskFilterCandidate(a) && opts.targetFilter == "" {
			opts.targetFilter = a
		}
	}
	return opts
}

func isTaskFilterCandidate(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return false
	}
	tok := strings.ToLower(arg)
	isSubcmd := tok == "tasks" || tok == "task" || tok == "task-queue" || tok == "tasks-queue" || tok == "queue" || tok == "jobs" || tok == "delegation"
	if isSubcmd {
		return false
	}
	return true
}

func detectHeavyDutyTask(totalStepsBudget int, title, taskCode string) bool {
	isBudgetHeavy := totalStepsBudget >= 100
	if isBudgetHeavy {
		return true
	}
	text := strings.ToLower(title + " " + taskCode)
	keywords := []string{"heavy", "deep-scan", "migration", "full-audit", "rebuild-all"}
	for _, kw := range keywords {
		hasKw := strings.Contains(text, kw)
		if hasKw {
			return true
		}
	}
	return false
}
