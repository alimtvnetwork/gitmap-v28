package cmdagy

import (
	"encoding/json"
	"os"
	"path/filepath"
)

// AgyWorkspaceQueueSummary encapsulates prompt queue status for a specific workspace.
type AgyWorkspaceQueueSummary struct {
	ProjectName string                `json:"projectName"`
	Workspace   string                `json:"workspace"`
	QueueFile   string                `json:"queueFile"`
	ActiveItem  *AgyPromptQueueEntry  `json:"activeItem,omitempty"`
	QueuedItems []AgyPromptQueueEntry `json:"queuedItems"`
	TotalQueued int                   `json:"totalQueued"`
}

func handleAgyQueuesInspector(args []string) error {
	summaries, err := DiscoverAllWorkspaceQueues()
	if err != nil {
		return err
	}
	if hasArgFlag(args, "--json") || hasArgFlag(os.Args, "--json") {
		return renderQueuesJSON(summaries)
	}
	renderQueuesTerminal(summaries)
	return nil
}

// DiscoverAllWorkspaceQueues scans projects and CWD for pending prompt queues.
func DiscoverAllWorkspaceQueues() ([]AgyWorkspaceQueueSummary, error) {
	wsMap := collectCandidateWorkspaces()
	var list []AgyWorkspaceQueueSummary

	for ws, projName := range wsMap {
		summary, isFound := inspectSingleWorkspaceQueue(ws, projName)
		if isFound {
			list = append(list, summary)
		}
	}
	return list, nil
}

func inspectSingleWorkspaceQueue(ws, projName string) (AgyWorkspaceQueueSummary, bool) {
	qPath := filepath.Join(ws, ".ai-memory", "temp", "agy-prompt-queue.json")
	data, err := os.ReadFile(qPath)
	if err != nil {
		return AgyWorkspaceQueueSummary{}, false
	}
	var q AgyPromptQueueFile
	if jsonErr := json.Unmarshal(data, &q); jsonErr != nil {
		return AgyWorkspaceQueueSummary{}, false
	}
	queuedList := filterPendingQueuedItems(q)
	total := len(queuedList)
	if total == 0 {
		return AgyWorkspaceQueueSummary{}, false
	}
	return AgyWorkspaceQueueSummary{
		ProjectName: projName,
		Workspace:   ws,
		QueueFile:   qPath,
		ActiveItem:  q.Active,
		QueuedItems: queuedList,
		TotalQueued: total,
	}, true
}

func filterPendingQueuedItems(q AgyPromptQueueFile) []AgyPromptQueueEntry {
	var items []AgyPromptQueueEntry
	if q.Active != nil && (q.Active.Status == "queued" || q.Active.Status == "pending") {
		items = append(items, *q.Active)
	}
	for _, item := range q.Queued {
		if isQueuePendingStatus(item.Status) {
			items = append(items, item)
		}
	}
	return items
}

func isQueuePendingStatus(status string) bool {
	return status == "queued" || status == "requeued_with_check_prefix" || status == "pending"
}
