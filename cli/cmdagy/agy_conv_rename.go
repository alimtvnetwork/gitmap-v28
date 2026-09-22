// Package cmdagy — agy_conv_rename.go updates conversation titles in Antigravity conversation_summaries.db.
package cmdagy

import (
	"database/sql"
	"fmt"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ConvRenameCmd renames an Antigravity conversation title.
var ConvRenameCmd = &cobra.Command{
	Use:     "conv-rename <id|seq|startswith|slug|path> <new-name>",
	Aliases: []string{"cr", "rename-conv"},
	Short:   "Rename conversation title in Antigravity",
	Long: `Rename a conversation title in Antigravity conversation_summaries.db.
Target can be:
  • <seq>: Persistent sequence number (e.g. 1, 001)
  • <id>: Full conversation UUID or prefix (e.g. 88823909)
  • <startswith>: Conversation title prefix
  • <slug>: Project name or slug
  • <path>: Project directory path`,
	Example: `  gitmap agy conv-rename 1 "Refactor database architecture"
  gitmap agy cr 002 "Verify CI/CD pipeline fix"
  gitmap agy cr 88823909 "Implement AGY prompt injection"
  gitmap agy cr gitmap "Active session"`,
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(args) < 2 {
			return apperror.NewValidationError("usage: gitmap agy conv-rename <id|seq|startswith|slug|path> \"new name\"")
		}
		return RunConvRename(args[0], args[1])
	},
}

func renameProjectInitialConversation(p AgyProject, newTitle string) error {
	dbPath, err := getConversationSummariesDBPath()
	if err != nil {
		return apperror.WrapSimple(err, "get conversation summaries db path")
	}

	conn, openErr := store.OpenSQLiteDB(dbPath)
	if openErr != nil {
		return apperror.WrapSimple(openErr, "open summaries db for rename")
	}
	defer conn.Close()

	convID, findErr := findLatestConvIDForProject(conn, p)
	if findErr != nil {
		return findErr
	}
	if convID == "" {
		return nil
	}

	return executeConversationTitleUpdate(conn, convID, newTitle)
}

func findLatestConvIDForProject(conn *sql.DB, p AgyProject) (string, error) {
	query := "SELECT conversation_id FROM conversation_summaries WHERE project_id = ? OR workspace_uris LIKE ? ORDER BY last_modified_time DESC LIMIT 1"
	pattern := "%" + p.GetPath() + "%"
	row := conn.QueryRow(query, p.ID, pattern)

	var convID string
	err := row.Scan(&convID)
	if err == sql.ErrNoRows {
		return "", nil
	}
	if err != nil {
		return "", apperror.WrapSimple(err, "scan latest conv id")
	}

	return convID, nil
}

func executeConversationTitleUpdate(conn *sql.DB, convID, newTitle string) error {
	query := "UPDATE conversation_summaries SET title = ? WHERE conversation_id = ?"
	res := store.ExecWrapper(conn, query, newTitle, convID)
	if res.IsFailure {
		return fmt.Errorf("update conversation %s title to %s: %w", convID, newTitle, res.Error)
	}

	return nil
}
