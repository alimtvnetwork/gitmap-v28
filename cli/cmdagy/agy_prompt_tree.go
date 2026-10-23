package cmdagy

import (
	"fmt"
	"os"
	"path/filepath"
	"time"
)

func normalizeWorkspacePath(path string) string {
	if path == "" {
		return "$WORKSPACE_ROOT"
	}
	clean := filepath.ToSlash(path)
	return clean
}

// BuildPromptTree aggregates running and queued prompts grouped by node and project.
func BuildPromptTree() (*PromptTreeResponse, error) {
	nodeAlias, _ := os.Hostname()
	running, _ := DiscoverRunningProjects()
	queuePath := getTempAgyFilePath("agy-prompt-queue.json")
	queuedRecords, _ := readRecordsFromFile(queuePath)

	projectMap := make(map[string]*ProjectPromptSummary)
	totalActive := 0

	for _, p := range running {
		projName := p.ProjectName
		if projName == "" {
			projName = filepath.Base(p.ProjectPath)
		}
		summary, exists := projectMap[projName]
		if !exists {
			summary = &ProjectPromptSummary{
				ProjectName:   projName,
				WorkspacePath: normalizeWorkspacePath(p.ProjectPath),
				ActivePrompts: []PromptRecord{},
				QueuedPrompts: []PromptRecord{},
			}
			projectMap[projName] = summary
		}
		if p.HasActivePrompt {
			totalActive++
			summary.ActivePrompts = append(summary.ActivePrompts, PromptRecord{
				ID:          fmt.Sprintf("p-%s", p.ConversationId),
				ProjectName: projName,
				Title:       p.PromptPreview,
				PromptText:  p.FullPrompt,
				Status:      p.Status,
				CreatedAt:   time.Now().Format(time.RFC3339),
			})
		}
	}

	totalQueued := len(queuedRecords)
	for _, q := range queuedRecords {
		projName := q.ProjectName
		if projName == "" {
			projName = "default"
		}
		summary, exists := projectMap[projName]
		if !exists {
			summary = &ProjectPromptSummary{
				ProjectName:   projName,
				WorkspacePath: "$WORKSPACE_ROOT",
				ActivePrompts: []PromptRecord{},
				QueuedPrompts: []PromptRecord{},
			}
			projectMap[projName] = summary
		}
		summary.QueuedPrompts = append(summary.QueuedPrompts, q)
	}

	var projects []ProjectPromptSummary
	for _, p := range projectMap {
		projects = append(projects, *p)
	}

	return &PromptTreeResponse{
		TotalActive: totalActive,
		TotalQueued: totalQueued,
		Nodes: []NodePromptTree{
			{
				NodeAlias: nodeAlias,
				IsOnline:  true,
				Projects:  projects,
			},
		},
	}, nil
}

func buildActiveSummaryFromProject(p RunningProjectRecord) ActivePromptSummary {
	snip := p.PromptPreview
	if snip == "" {
		snip = p.FullPrompt
	}
	if len(snip) > 80 {
		snip = snip[:77] + "..."
	}
	return ActivePromptSummary{
		WorkspacePath:  p.ProjectPath,
		ProjectName:    p.ProjectName,
		ConversationID: p.ConversationId,
		PromptTitle:    p.PromptPreview,
		PromptSnippet:  snip,
		Status:         p.Status,
	}
}

// CollectActivePrompts extracts summaries from running projects.
func CollectActivePrompts(running []RunningProjectRecord) []ActivePromptSummary {
	var summaries []ActivePromptSummary
	for _, p := range running {
		if p.HasActivePrompt {
			summaries = append(summaries, buildActiveSummaryFromProject(p))
		}
	}
	return summaries
}

// GetAGYUIStatus collects complete telemetry for the AGY web studio.
func GetAGYUIStatus() (*AGYUIStatusPayload, error) {
	node, _ := os.Hostname()
	running, _ := DiscoverRunningProjects()
	activePrompts := CollectActivePrompts(running)
	qStatus, _ := GetQueueStatus()
	saved, _ := LoadSavedPrompts()

	return &AGYUIStatusPayload{
		IsServerRunning: true,
		NodeAlias:       node,
		RunningProjects: running,
		ActivePrompts:   activePrompts,
		QueuedPrompts:   qStatus.Queued,
		SavedPrompts:    saved,
	}, nil
}
