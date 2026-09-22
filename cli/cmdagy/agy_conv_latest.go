// Package cmdagy — agy_conv_latest.go fetches latest conversation title and ID per project.
package cmdagy

import (
	"database/sql"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AgyLatestConv holds the conversation ID and title for a project.
type AgyLatestConv struct {
	ConvID string
	Title  string
}

// loadLatestConversationsMap queries conversation_summaries.db for each project's latest conversation.
func loadLatestConversationsMap() map[string]AgyLatestConv {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return nil
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return nil
	}
	defer conn.Close()

	return queryLatestConversations(conn)
}

func queryLatestConversations(conn *sql.DB) map[string]AgyLatestConv {
	query := "SELECT project_id, conversation_id, title FROM conversation_summaries WHERE project_id != '' ORDER BY last_modified_time DESC"
	rows, err := conn.Query(query)
	if err != nil {
		return nil
	}
	defer rows.Close()

	result := make(map[string]AgyLatestConv)
	for rows.Next() {
		var pid, cid, title string
		if scanErr := rows.Scan(&pid, &cid, &title); scanErr != nil {
			continue
		}
		trimmed := strings.TrimSpace(title)
		existing, exists := result[pid]
		if !exists {
			result[pid] = AgyLatestConv{ConvID: cid, Title: trimmed}
		} else if existing.Title == "" && trimmed != "" {
			existing.Title = trimmed
			result[pid] = existing
		}
	}

	return result
}

func resolveProjectConvDetails(projectID string, convMap map[string]AgyLatestConv) (string, string) {
	if convMap == nil {
		return "—", "—"
	}

	conv, hasConv := convMap[projectID]
	if !hasConv {
		return "—", "—"
	}

	title := conv.Title
	if title == "" {
		title = "—"
	}

	cid := conv.ConvID
	if cid == "" {
		cid = "—"
	}

	return title, cid
}
