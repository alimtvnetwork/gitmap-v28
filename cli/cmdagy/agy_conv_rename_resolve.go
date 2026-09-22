package cmdagy

import (
	"database/sql"
	"fmt"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RunConvRename resolves and renames a conversation in Antigravity summaries DB.
func RunConvRename(target, newTitle string) error {
	newTitle = strings.TrimSpace(newTitle)
	if len(newTitle) == 0 {
		return apperror.NewValidationError("new conversation title cannot be empty")
	}

	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return apperror.WrapSimple(err, "get conversation summaries db path")
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "open summaries db for rename")
	}
	defer conn.Close()

	convID, oldTitle, projectID, seq, findErr := resolveTargetConversation(conn, target)
	if findErr != nil {
		return findErr
	}

	updateErr := executeConversationTitleUpdate(conn, convID, newTitle)
	if updateErr != nil {
		return updateErr
	}

	renderRenameSuccess(convID, projectID, oldTitle, newTitle, seq)

	return nil
}

func resolveTargetConversation(conn *sql.DB, target string) (string, string, string, int, error) {
	if convID, oldTitle, proj, seq, ok := tryResolveBySequence(conn, target); ok {
		return convID, oldTitle, proj, seq, nil
	}
	if convID, oldTitle, proj, seq, ok := tryResolveByIDOrPrefix(conn, target); ok {
		return convID, oldTitle, proj, seq, nil
	}
	if convID, oldTitle, proj, seq, ok := tryResolveByPattern(conn, target); ok {
		return convID, oldTitle, proj, seq, nil
	}

	return "", "", "", 0, apperror.NewNotFound("conv-rename", "E9035", fmt.Sprintf("no conversation found matching target: %s", target))
}

func tryResolveBySequence(conn *sql.DB, target string) (string, string, string, int, bool) {
	seq, parseErr := strconv.Atoi(target)
	if parseErr != nil || seq <= 0 {
		return "", "", "", 0, false
	}
	rec, err := GetProjectBySequence(seq)
	if err != nil || rec == nil {
		return "", "", "", 0, false
	}
	convID, oldTitle, findErr := findLatestConvForProjectData(conn, rec.ProjectID, rec.WorkspacePath)
	if findErr != nil || convID == "" {
		return "", "", "", 0, false
	}

	return convID, oldTitle, rec.ProjectName, rec.SequenceNum, true
}

func findLatestConvForProjectData(conn *sql.DB, projectID, workspacePath string) (string, string, error) {
	query := "SELECT conversation_id, title FROM conversation_summaries WHERE project_id = ? OR workspace_uris LIKE ? ORDER BY last_modified_time DESC LIMIT 1"
	pattern := "%" + workspacePath + "%"
	var convID, title string
	err := conn.QueryRow(query, projectID, pattern).Scan(&convID, &title)
	if err != nil {
		return "", "", err
	}

	return convID, title, nil
}

func tryResolveByIDOrPrefix(conn *sql.DB, target string) (string, string, string, int, bool) {
	query := "SELECT conversation_id, title, project_id FROM conversation_summaries WHERE conversation_id = ? OR conversation_id LIKE ? LIMIT 1"
	var convID, title, projID string
	err := conn.QueryRow(query, target, target+"%").Scan(&convID, &title, &projID)
	if err != nil {
		return "", "", "", 0, false
	}
	seqMap, _ := GetAllProjectSequences()
	seq := seqMap[projID]

	return convID, title, projID, seq, true
}

func tryResolveByPattern(conn *sql.DB, target string) (string, string, string, int, bool) {
	query := "SELECT conversation_id, title, project_id FROM conversation_summaries WHERE title LIKE ? OR project_id LIKE ? OR workspace_uris LIKE ? ORDER BY last_modified_time DESC LIMIT 1"
	pattern := "%" + target + "%"
	var convID, title, projID string
	err := conn.QueryRow(query, pattern, pattern, pattern).Scan(&convID, &title, &projID)
	if err != nil {
		return "", "", "", 0, false
	}
	seqMap, _ := GetAllProjectSequences()
	seq := seqMap[projID]

	return convID, title, projID, seq, true
}

func renderRenameSuccess(convID, projectID, oldTitle, newTitle string, seq int) {
	seqStr := ""
	if seq > 0 {
		seqStr = fmt.Sprintf(", SEQ: %03d", seq)
	}
	projStr := ""
	if len(projectID) > 0 {
		projStr = fmt.Sprintf(" (Project: '%s'%s)", projectID, seqStr)
	}
	fmt.Printf("  %s✔ Renamed conversation '%s'%s to: \"%s\"%s\n",
		constants.ColorGreen, convID, projStr, newTitle, constants.ColorReset)
}
