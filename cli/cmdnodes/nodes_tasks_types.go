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

// NodeTaskReport summarizes active tasks, hardware pressure, and queue depth for a single fleet node.
type NodeTaskReport struct {
	NodeAlias        string            `json:"nodeAlias"`
	NodeHost         string            `json:"nodeHost"`
	NodePort         int               `json:"nodePort"`
	IsLocal          bool              `json:"isLocal"`
	Status           string            `json:"status"`
	CPUPercent       int               `json:"cpuPercent"`
	RAMPercent       int               `json:"ramPercent"`
	TotalRAMMB       int               `json:"totalRamMb"`
	UsedRAMMB        int               `json:"usedRamMb"`
	RolePlaceholder  string            `json:"rolePlaceholder"`
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
	isOpenOnly   bool
	isHeavyOnly  bool
	limit        int
	targetFilter string
}

func parseNodesTasksOptions(args []string) nodesTasksOptions {
	opts := nodesTasksOptions{
		limit: 10,
	}
	for i := 0; i < len(args); i++ {
		a := args[i]
		switch {
		case a == "--json" || a == "-j":
			opts.isJSON = true
		case a == "--open-only" || a == "--online-only" || a == "--running-only" || a == "-ro":
			opts.isOpenOnly = true
		case a == "--heavy-only" || a == "--heavy" || a == "-ho":
			opts.isHeavyOnly = true
		case (a == "--limit" || a == "-l") && i+1 < len(args):
			i++
			opts.limit = parsePositiveIntLimit(args[i], opts.limit)
		case isTaskFilterCandidate(a) && opts.targetFilter == "":
			opts.targetFilter = a
		}
	}
	return opts
}

func parsePositiveIntLimit(val string, fallback int) int {
	n, err := strconv.Atoi(val)
	if err == nil && n > 0 {
		return n
	}
	return fallback
}

func isTaskFilterCandidate(arg string) bool {
	if strings.HasPrefix(arg, "-") {
		return false
	}
	tok := strings.ToLower(arg)
	isSubcmd := tok == "tasks" || tok == "task" || tok == "task-queue" || tok == "tasks-queue" || tok == "queue" || tok == "jobs" || tok == "delegation" || tok == "nodes"
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
	keywords := []string{"heavy", "deep-scan", "migration", "full-audit", "rebuild-all", "build-all", "e2e", "coverage"}
	for _, kw := range keywords {
		hasKw := strings.Contains(text, kw)
		if hasKw {
			return true
		}
	}
	return false
}
