package cmdagy

import (
	"database/sql"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AgyActiveConversation represents an Antigravity conversation with an active prompt running.
type AgyActiveConversation struct {
	ConversationID   string `json:"conversationId"`
	Title            string `json:"title"`
	WorkspacePath    string `json:"workspacePath"`
	ProjectID        string `json:"projectId"`
	StepCount        int    `json:"stepCount"`
	LastModifiedTime string `json:"lastModifiedTime"`
	ElapsedSeconds   int    `json:"elapsedSeconds"`
	IsRunning        bool   `json:"isRunning"`
}

func handleAgyActivePrompts(args []string) error {
	convs, err := FetchActiveRunningConversations()
	if err != nil {
		return err
	}
	if hasArgFlag(args, "--json") || hasArgFlag(os.Args, "--json") {
		return renderActivePromptsJSON(convs)
	}
	renderActivePromptsTerminal(convs)
	return nil
}

// FetchActiveRunningConversations queries conversation_summaries.db for not_fully_idle != 0.
func FetchActiveRunningConversations() ([]AgyActiveConversation, error) {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return nil, err
	}
	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return nil, openErr
	}
	defer conn.Close()

	return queryRunningConversations(conn)
}

func queryRunningConversations(conn *sql.DB) ([]AgyActiveConversation, error) {
	query := `SELECT conversation_id, title, workspace_uris, project_id, step_count, last_modified_time
		FROM conversation_summaries
		WHERE (killed IS NULL OR killed = 0) AND not_fully_idle != 0
		ORDER BY last_modified_time DESC;`
	rows, err := conn.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query running conversations")
	}
	defer rows.Close()

	return scanRunningConversations(rows)
}

func scanRunningConversations(rows *sql.Rows) ([]AgyActiveConversation, error) {
	var list []AgyActiveConversation
	for rows.Next() {
		var c AgyActiveConversation
		var wsRaw, pid sql.NullString
		if err := rows.Scan(&c.ConversationID, &c.Title, &wsRaw, &pid, &c.StepCount, &c.LastModifiedTime); err != nil {
			continue
		}
		c.WorkspacePath = extractCleanWorkspaceFromURIs(wsRaw.String)
		c.ProjectID = pid.String
		c.IsRunning = true
		c.ElapsedSeconds = calculateElapsedSeconds(c.LastModifiedTime)
		list = append(list, c)
	}
	return list, nil
}

func calculateElapsedSeconds(lastMod string) int {
	if lastMod == "" {
		return 0
	}
	t, err := time.Parse(time.RFC3339, lastMod)
	if err != nil {
		return 0
	}
	elapsed := int(time.Since(t).Seconds())
	if elapsed < 0 {
		return 0
	}
	return elapsed
}
