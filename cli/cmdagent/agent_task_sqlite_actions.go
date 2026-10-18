// Package cmdagent — agent_task_sqlite_actions.go implements complete, fail, and log-action operations.
package cmdagent

import (
	"flag"
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func runTaskSQLiteComplete(args []string) error {
	fs := flag.NewFlagSet("task complete", flag.ContinueOnError)
	dbFlag := fs.String("db", "", "Path to SQLite database")
	subtaskId := fs.Int64("subtask-id", 0, "SubtaskId to mark complete")
	agent := fs.String("agent", "Worker", "Agent role name")
	evidence := fs.String("evidence", "PASS exit 0", "Verification evidence")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if *subtaskId <= 0 {
		return appfault.NewValidationError("missing required flag: --subtask-id")
	}

	dbPath := resolveDbPath(*dbFlag)
	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	now := NowISO()
	_, updateErr := db.Exec(`
		UPDATE Subtask
		SET Status = 'DONE', HasCompleted = 1, IsBlocked = 0, Evidence = ?, UpdatedAt = ?
		WHERE SubtaskId = ?;
	`, *evidence, now, *subtaskId)
	if updateErr != nil {
		return appfault.WrapExecution(updateErr, "failed to update subtask to DONE")
	}

	_, logErr := db.Exec(`
		INSERT INTO AgentActionLog (
			SubtaskId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
		) VALUES (?, ?, 'COMPLETE', 'Subtask verified and completed', 'DONE', ?);
	`, *subtaskId, *agent, now)
	if logErr != nil {
		return appfault.WrapExecution(logErr, "failed to log complete action")
	}

	out := map[string]any{
		"status":    "COMPLETED",
		"subtaskId": *subtaskId,
	}

	return printJson(out)
}

func runTaskSQLiteFail(args []string) error {
	fs := flag.NewFlagSet("task fail", flag.ContinueOnError)
	dbFlag := fs.String("db", "", "Path to SQLite database")
	subtaskId := fs.Int64("subtask-id", 0, "SubtaskId to mark failed")
	agent := fs.String("agent", "Worker", "Agent role name")
	reason := fs.String("reason", "Failed", "Failure reason")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if *subtaskId <= 0 {
		return appfault.NewValidationError("missing required flag: --subtask-id")
	}

	dbPath := resolveDbPath(*dbFlag)
	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	now := NowISO()
	_, updateErr := db.Exec(`
		UPDATE Subtask
		SET Status = 'FAILED', IsBlocked = 1, Evidence = ?, UpdatedAt = ?
		WHERE SubtaskId = ?;
	`, *reason, now, *subtaskId)
	if updateErr != nil {
		return appfault.WrapExecution(updateErr, "failed to update subtask to FAILED")
	}

	_, logErr := db.Exec(`
		INSERT INTO AgentActionLog (
			SubtaskId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
		) VALUES (?, ?, 'FAIL', ?, 'FAILED', ?);
	`, *subtaskId, *agent, *reason, now)
	if logErr != nil {
		return appfault.WrapExecution(logErr, "failed to log fail action")
	}

	out := map[string]any{
		"status":    "FAILED",
		"subtaskId": *subtaskId,
	}

	return printJson(out)
}

func runTaskSQLiteLogAction(args []string) error {
	fs := flag.NewFlagSet("task log-action", flag.ContinueOnError)
	dbFlag := fs.String("db", "", "Path to SQLite database")
	subtaskId := fs.Int64("subtask-id", 0, "SubtaskId")
	agent := fs.String("agent", "", "Agent role name")
	action := fs.String("action", "", "Action type")
	file := fs.String("file", "", "Target relative file path")
	details := fs.String("details", "", "Action details")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if *subtaskId <= 0 || strings.TrimSpace(*agent) == "" || strings.TrimSpace(*action) == "" {
		return appfault.NewValidationError("missing required flags: --subtask-id, --agent, --action")
	}

	dbPath := resolveDbPath(*dbFlag)
	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	det := *details
	if det == "" {
		det = fmt.Sprintf("Executing %s on %s", *action, *file)
	}

	now := NowISO()
	res, execErr := db.Exec(`
		INSERT INTO AgentActionLog (
			SubtaskId, AgentRole, ActionType, TargetFile, ActionDetails, Status, CreatedAt
		) VALUES (?, ?, ?, ?, ?, 'IN_PROGRESS', ?);
	`, *subtaskId, *agent, *action, *file, det, now)
	if execErr != nil {
		return appfault.WrapExecution(execErr, "failed to log action")
	}

	logId, _ := res.LastInsertId()
	out := map[string]any{
		"status":      "LOGGED",
		"actionLogId": logId,
	}

	return printJson(out)
}
