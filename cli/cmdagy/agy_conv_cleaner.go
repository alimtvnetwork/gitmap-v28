// Package cmdagy — agy_conv_cleaner.go purges project conversation state from Gemini/Antigravity storage.
package cmdagy

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func purgeProjectConversations(p AgyProject) (int, error) {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return 0, nil
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return 0, apperror.WrapSimple(openErr, "open summaries db")
	}
	defer conn.Close()

	convIDs, findErr := findConversationIDs(conn, p)
	if findErr != nil {
		return 0, findErr
	}

	return deleteConversations(conn, convIDs)
}

func findConversationIDs(conn *sql.DB, p AgyProject) ([]string, error) {
	query := "SELECT conversation_id FROM conversation_summaries WHERE project_id = ? OR workspace_uris LIKE ? OR LOWER(title) = LOWER(?)"
	pattern := "%" + p.GetPath() + "%"
	rows, err := conn.Query(query, p.ID, pattern, p.Name)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query conversation ids")
	}
	defer rows.Close()

	var ids []string
	for rows.Next() {
		var id string
		if scanErr := rows.Scan(&id); scanErr == nil && id != "" {
			ids = append(ids, id)
		}
	}
	return ids, nil
}

func deleteConversations(conn *sql.DB, ids []string) (int, error) {
	home, homeErr := os.UserHomeDir()
	if homeErr != nil {
		return 0, apperror.WrapSimple(homeErr, "user home dir")
	}

	convDir := filepath.Join(home, ".gemini", "antigravity", "conversations")
	brainDir := filepath.Join(home, ".gemini", "antigravity", "brain")
	count := 0

	for _, id := range ids {
		removeSingleConvFiles(convDir, brainDir, id)
		if delErr := deleteSummaryRecord(conn, id); delErr != nil {
			return count, delErr
		}
		count++
	}

	return count, nil
}

func removeSingleConvFiles(convDir, brainDir, id string) {
	_ = os.Remove(filepath.Join(convDir, id+".db"))
	_ = os.Remove(filepath.Join(convDir, id+".db-wal"))
	_ = os.Remove(filepath.Join(convDir, id+".db-shm"))
	_ = os.RemoveAll(filepath.Join(brainDir, id))
}

func deleteSummaryRecord(conn *sql.DB, id string) error {
	res := store.ExecWrapper(conn, "DELETE FROM conversation_summaries WHERE conversation_id = ?", id)
	if res.IsFailure {
		return fmt.Errorf("delete conversation %s from summaries: %w", id, res.Error)
	}
	return nil
}
