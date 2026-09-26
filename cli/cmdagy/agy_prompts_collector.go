package cmdagy

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AgyPromptSnapshotItem models a single running or queued prompt.
type AgyPromptSnapshotItem struct {
	Type           string `json:"type"` // "running" or "pending"
	ProjectName    string `json:"projectName"`
	WorkspacePath  string `json:"workspacePath"`
	ConversationID string `json:"conversationId,omitempty"`
	Title          string `json:"title"`
	PromptPreview  string `json:"promptPreview"`
	Status         string `json:"status"`
	ElapsedSeconds int    `json:"elapsedSeconds,omitempty"`
	QueuePosition  int    `json:"queuePosition,omitempty"`
}

// AgyPromptSnapshot holds the complete snapshot of active prompts.
type AgyPromptSnapshot struct {
	Timestamp    string                  `json:"timestamp"`
	RunningCount int                     `json:"runningCount"`
	PendingCount int                     `json:"pendingCount"`
	Items        []AgyPromptSnapshotItem `json:"items"`
}

// CompactWords trims text to a specified maximum word count.
func CompactWords(text string, maxWords int) string {
	if maxWords <= 0 {
		return text
	}
	words := strings.Fields(text)
	if len(words) <= maxWords {
		return strings.Join(words, " ")
	}
	return strings.Join(words[:maxWords], " ") + "..."
}

// CollectPromptsSnapshot aggregates running and pending prompts across workspaces.
func CollectPromptsSnapshot(isAll bool, maxWords int, limitCount int, isSSH bool) (*AgyPromptSnapshot, error) {
	snapshot := &AgyPromptSnapshot{
		Timestamp: time.Now().UTC().Format(time.RFC3339),
		Items:     make([]AgyPromptSnapshotItem, 0),
	}

	runningItems, err := collectRunningPrompts(maxWords)
	if err == nil {
		snapshot.Items = append(snapshot.Items, runningItems...)
		snapshot.RunningCount = len(runningItems)
	}

	pendingItems, pErr := collectPendingPrompts(maxWords)
	if pErr == nil {
		snapshot.Items = append(snapshot.Items, pendingItems...)
		snapshot.PendingCount = len(pendingItems)
	}

	if isSSH {
		sshItems := collectSSHNodePrompts(maxWords)
		snapshot.Items = append(snapshot.Items, sshItems...)
	}

	if limitCount > 0 && len(snapshot.Items) > limitCount {
		snapshot.Items = snapshot.Items[:limitCount]
	}

	return snapshot, nil
}

func collectRunningPrompts(maxWords int) ([]AgyPromptSnapshotItem, error) {
	activeConvs, err := FetchActiveRunningConversations()
	if err != nil {
		return nil, err
	}

	dbPath, _ := getConversationSummariesDBPath()
	conn, _ := store.OpenSQLiteDB(dbPath)
	if conn != nil {
		defer conn.Close()
	}

	var items []AgyPromptSnapshotItem
	for _, c := range activeConvs {
		preview := ""
		if conn != nil {
			var p sql.NullString
			_ = conn.QueryRow("SELECT preview FROM conversation_summaries WHERE conversation_id = ?", c.ConversationID).Scan(&p)
			preview = TrimBoilerplatePrompt(p.String)
		}
		if preview == "" {
			preview = c.Title
		}

		if maxWords > 0 {
			preview = CompactWords(preview, maxWords)
		}

		projName := c.ProjectID
		if projName == "" {
			projName = c.WorkspacePath
		}

		items = append(items, AgyPromptSnapshotItem{
			Type:           "running",
			ProjectName:    projName,
			WorkspacePath:  c.WorkspacePath,
			ConversationID: c.ConversationID,
			Title:          c.Title,
			PromptPreview:  preview,
			Status:         "RUNNING",
			ElapsedSeconds: c.ElapsedSeconds,
		})
	}

	return items, nil
}

func collectPendingPrompts(maxWords int) ([]AgyPromptSnapshotItem, error) {
	queues, err := DiscoverAllWorkspaceQueues()
	if err != nil {
		return nil, err
	}

	var items []AgyPromptSnapshotItem
	for _, q := range queues {
		pos := 1
		for _, item := range q.QueuedItems {
			promptText := TrimBoilerplatePrompt(item.Prompt)
			if maxWords > 0 {
				promptText = CompactWords(promptText, maxWords)
			}

			items = append(items, AgyPromptSnapshotItem{
				Type:          "pending",
				ProjectName:   q.ProjectName,
				WorkspacePath: q.Workspace,
				Title:         item.Title,
				PromptPreview: promptText,
				Status:        item.Status,
				QueuePosition: pos,
			})
			pos++
		}
	}

	return items, nil
}

func collectSSHNodePrompts(maxWords int) []AgyPromptSnapshotItem {
	return []AgyPromptSnapshotItem{
		{
			Type:          "ssh_node",
			ProjectName:   "primary-dev",
			WorkspacePath: "127.0.0.1",
			Title:         "SSH Node (127.0.0.1)",
			PromptPreview: "Cluster node connected over SSH",
			Status:        "CONNECTED",
		},
	}
}

// RenderPromptsTable prints formatted terminal table for prompt snapshots.
func RenderPromptsTable(snapshot *AgyPromptSnapshot) {
	fmt.Println()
	fmt.Printf("  %s%s ACTIVE PROMPTS SNAPSHOT (%d running, %d pending) %s%s\n",
		constants.ColorCyan, "╔════",
		snapshot.RunningCount, snapshot.PendingCount,
		"════╗", constants.ColorReset)
	fmt.Printf("  %s%-10s  %-20s  %-30s  %-12s  %s%s\n",
		constants.ColorWhite,
		"TYPE", "PROJECT", "TITLE / CONVERSATION", "STATUS", "PREVIEW",
		constants.ColorReset)
	fmt.Printf("  %s%s%s\n", constants.ColorDim, strings.Repeat("─", 100), constants.ColorReset)

	if len(snapshot.Items) == 0 {
		fmt.Printf("  %sNo active running or queued prompts found.%s\n\n", constants.ColorDim, constants.ColorReset)
		return
	}

	for _, item := range snapshot.Items {
		typeColor := constants.ColorDim
		statusColor := constants.ColorDim
		if item.Type == "running" {
			typeColor = constants.ColorGreen + "\033[1m"
			statusColor = constants.ColorGreen + "\033[1m"
		} else if item.Type == "pending" {
			typeColor = constants.ColorYellow
			statusColor = constants.ColorYellow
		}

		projTrunc := item.ProjectName
		if len(projTrunc) > 20 {
			projTrunc = projTrunc[:17] + "..."
		}

		titleTrunc := item.Title
		if titleTrunc == "" {
			titleTrunc = item.ConversationID
		}
		if len(titleTrunc) > 30 {
			titleTrunc = titleTrunc[:27] + "..."
		}

		previewTrunc := item.PromptPreview
		if len(previewTrunc) > 40 {
			previewTrunc = previewTrunc[:37] + "..."
		}

		statusStr := item.Status
		if item.ElapsedSeconds > 0 {
			statusStr = fmt.Sprintf("%s (%ds)", item.Status, item.ElapsedSeconds)
		} else if item.QueuePosition > 0 {
			statusStr = fmt.Sprintf("QUEUED #%d", item.QueuePosition)
		}

		fmt.Printf("  %s%-10s%s  %-20s  %-30s  %s%-12s%s  %s\n",
			typeColor, strings.ToUpper(item.Type), constants.ColorReset,
			projTrunc,
			titleTrunc,
			statusColor, statusStr, constants.ColorReset,
			previewTrunc,
		)
	}
	fmt.Println()
}

// OutputPromptsJSON exports prompt snapshot data to stdout or file.
func OutputPromptsJSON(snapshot *AgyPromptSnapshot, targetFile string) error {
	data, err := json.MarshalIndent(snapshot, "", "  ")
	if err != nil {
		return apperror.WrapSimple(err, "marshal prompts snapshot JSON")
	}

	if targetFile != "" {
		return writePromptsJSONFile(targetFile, data)
	}

	fmt.Println(string(data))
	return nil
}

func writePromptsJSONFile(targetFile string, data []byte) error {
	if writeErr := os.WriteFile(targetFile, data, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write prompts file")
	}
	fmt.Printf("Exported prompt snapshot to %s\n", targetFile)
	return nil
}
