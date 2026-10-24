package cmdagy

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func fetchInstanceRunningPrompts(inst AgyInstanceInfo, maxWords int) []AgyPromptSnapshotItem {
	dbPath := inst.SummariesDBPath
	if dbPath == "" {
		return nil
	}

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil
	}
	conn := openRes.Value
	defer conn.Close()

	activeConvs, qErr := queryRunningConversations(conn)
	if qErr != nil {
		return nil
	}

	var items []AgyPromptSnapshotItem
	for _, c := range activeConvs {
		preview := queryConversationPreview(conn, c.ConversationID)
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

	return items
}

func fetchInstanceQueuedPrompts(inst AgyInstanceInfo) []AgyPromptQueueEntry {
	var queued []AgyPromptQueueEntry

	for _, ws := range inst.ActiveWorkspaces {
		queueFile := filepath.Join(ws, ".ai-memory", "temp", "agy-prompt-queue.json")
		data, err := os.ReadFile(queueFile)
		if err != nil {
			continue
		}

		var q AgyPromptQueueFile
		if jsonErr := json.Unmarshal(data, &q); jsonErr == nil {
			queued = append(queued, q.Queued...)
		}
	}

	return queued
}

func fetchInstanceRecentConversations(inst AgyInstanceInfo, limit, maxWords int) []AgyConversationPreview {
	dbPath := inst.SummariesDBPath
	if dbPath == "" {
		return nil
	}

	openRes := store.OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil
	}
	conn := openRes.Value
	defer conn.Close()

	if limit <= 0 {
		limit = 10
	}

	query := `SELECT conversation_id, title, workspace_uris, preview, step_count, last_modified_time, not_fully_idle
		FROM conversation_summaries
		WHERE (killed IS NULL OR killed = 0)
		ORDER BY last_modified_time DESC
		LIMIT ?;`

	rows, qErr := conn.Query(query, limit)
	if qErr != nil {
		return nil
	}
	defer rows.Close()

	var convs []AgyConversationPreview
	for rows.Next() {
		var cid, title, wsRaw, prev, lastMod sql.NullString
		var stepCount, idle int
		if scanErr := rows.Scan(&cid, &title, &wsRaw, &prev, &stepCount, &lastMod, &idle); scanErr != nil {
			continue
		}

		previewText := prev.String
		if maxWords > 0 && previewText != "" {
			previewText = CompactWords(previewText, maxWords)
		}

		cleanWs := extractCleanWorkspaceFromURIs(wsRaw.String)

		convs = append(convs, AgyConversationPreview{
			ConversationID: cid.String,
			Title:          title.String,
			WorkspacePath:  cleanWs,
			LastPromptText: previewText,
			StepCount:      stepCount,
			UpdatedAt:      lastMod.String,
			IsActive:       idle != 0,
		})
	}

	return convs
}
