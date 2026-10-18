// Package cmdagent — agent_task_sqlite_subtask.go implements subtask registration and claiming.
package cmdagent

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type subtaskInput struct {
	Code       string   `json:"code"`
	TaskCode   string   `json:"task_code"`
	Title      string   `json:"title"`
	AgentRole  string   `json:"agent_role"`
	OwnedFiles []string `json:"owned_files"`
}

func runTaskSQLiteAdd(args []string) error {
	fs := flag.NewFlagSet("task add", flag.ContinueOnError)
	dbFlag := fs.String("db", "", "Path to SQLite database")
	jsonStr := fs.String("tasks-json", "", "JSON string representing list of subtasks")
	jsonFile := fs.String("tasks-file", "", "Path to JSON file representing list of subtasks")
	parentIdFlag := fs.Int64("parent-id", 0, "Optional ParentTaskId")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	rawJson, readErr := readTasksRawJson(*jsonStr, *jsonFile)
	if readErr != nil {
		return readErr
	}
	if rawJson == "" {
		return appfault.NewValidationError("missing required flag: --tasks-json or --tasks-file")
	}

	subtasks, parseErr := parseTaskManagerSubtasksJson(rawJson)
	if parseErr != nil {
		return parseErr
	}

	if len(subtasks) == 0 {
		return appfault.NewValidationError("subtask list is empty")
	}

	dbPath := resolveDbPath(*dbFlag)
	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	parentId, resolveErr := resolveParentTaskId(db, *parentIdFlag)
	if resolveErr != nil {
		return resolveErr
	}

	now := NowISO()
	var createdIds []int64

	for idx, st := range subtasks {
		code := st.Code
		if code == "" {
			code = st.TaskCode
		}
		if code == "" {
			code = fmt.Sprintf("Task-%02d", idx+1)
		}

		ownedFilesBytes, _ := json.Marshal(st.OwnedFiles)
		ownedJson := string(ownedFilesBytes)

		insertQuery := `
			INSERT INTO Subtask (
				ParentTaskId, TaskCode, Title, AssignedAgentRole, OwnedFilesJson,
				Status, IsBlocked, HasCompleted, Evidence, CreatedAt, UpdatedAt
			) VALUES (?, ?, ?, ?, ?, 'PENDING', 0, 0, NULL, ?, ?);
		`
		res, execErr := db.Exec(insertQuery, parentId, code, st.Title, st.AgentRole, ownedJson, now, now)
		if execErr != nil {
			return appfault.WrapExecution(execErr, "failed to insert subtask")
		}
		newId, _ := res.LastInsertId()
		createdIds = append(createdIds, newId)
	}

	out := map[string]any{
		"status":          "SUCCESS",
		"subtasksCreated": len(createdIds),
		"subtaskIds":      createdIds,
	}

	return printJson(out)
}

func parseTaskManagerSubtasksJson(raw string) ([]subtaskInput, error) {
	var directList []subtaskInput
	err := json.Unmarshal([]byte(raw), &directList)
	if err == nil {
		return directList, nil
	}

	var wrapper map[string]any
	wrapErr := json.Unmarshal([]byte(raw), &wrapper)
	if wrapErr != nil {
		return nil, appfault.WrapValidation(err, "invalid JSON for subtasks")
	}

	var itemsRaw any
	if list, ok := wrapper["subtasks"]; ok {
		itemsRaw = list
	} else if list, ok := wrapper["tasks"]; ok {
		itemsRaw = list
	}

	if itemsRaw == nil {
		return nil, appfault.NewValidationError("dict input must contain 'tasks' or 'subtasks' list")
	}

	bytes, _ := json.Marshal(itemsRaw)
	var parsed []subtaskInput
	if err := json.Unmarshal(bytes, &parsed); err != nil {
		return nil, appfault.WrapValidation(err, "invalid subtasks array in dict")
	}

	return parsed, nil
}

func runTaskSQLiteClaim(args []string) error {
	fs := flag.NewFlagSet("task claim", flag.ContinueOnError)
	dbFlag := fs.String("db", "", "Path to SQLite database")
	agent := fs.String("agent", "", "Agent role name (e.g. Worker 01)")
	subtaskId := fs.Int64("subtask-id", 0, "Optional specific SubtaskId to claim")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	if strings.TrimSpace(*agent) == "" {
		return appfault.NewValidationError("missing required flag: --agent")
	}

	dbPath := resolveDbPath(*dbFlag)
	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	var row *sql.Row
	if *subtaskId > 0 {
		row = db.QueryRow(`
			SELECT SubtaskId, TaskCode, Title, OwnedFilesJson
			FROM Subtask
			WHERE SubtaskId = ? AND Status = 'PENDING'
			LIMIT 1;
		`, *subtaskId)
	} else {
		row = db.QueryRow(`
			SELECT SubtaskId, TaskCode, Title, OwnedFilesJson
			FROM Subtask
			WHERE Status = 'PENDING'
			ORDER BY SubtaskId ASC
			LIMIT 1;
		`)
	}

	var id int64
	var code, title, ownedJson string
	scanErr := row.Scan(&id, &code, &title, &ownedJson)
	if scanErr == sql.ErrNoRows {
		return printJson(map[string]any{"status": "NO_TASKS_AVAILABLE"})
	}
	if scanErr != nil {
		return appfault.WrapExecution(scanErr, "failed to query pending subtask")
	}

	now := NowISO()
	_, updateErr := db.Exec(`
		UPDATE Subtask
		SET Status = 'IN_PROGRESS', AssignedAgentRole = ?, UpdatedAt = ?
		WHERE SubtaskId = ?;
	`, *agent, now, id)
	if updateErr != nil {
		return appfault.WrapExecution(updateErr, "failed to update subtask to IN_PROGRESS")
	}

	_, logErr := db.Exec(`
		INSERT INTO AgentActionLog (
			SubtaskId, AgentRole, ActionType, ActionDetails, Status, CreatedAt
		) VALUES (?, ?, 'CLAIM', 'Claimed subtask for execution', 'IN_PROGRESS', ?);
	`, id, *agent, now)
	if logErr != nil {
		return appfault.WrapExecution(logErr, "failed to record claim in action log")
	}

	var ownedFiles []string
	_ = json.Unmarshal([]byte(ownedJson), &ownedFiles)

	out := map[string]any{
		"status":            "CLAIMED",
		"subtaskId":         id,
		"taskCode":          code,
		"title":             title,
		"assignedAgentRole": *agent,
		"ownedFiles":        ownedFiles,
	}

	return printJson(out)
}
