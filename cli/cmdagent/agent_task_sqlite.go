package cmdagent

import (
	"database/sql"
	"encoding/json"
	"flag"
	"fmt"
	"math"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// DispatchTaskSQLite handles gitmap task agent SQLite operations.
func DispatchTaskSQLite(args []string) error {
	hasArgs := len(args) > 0
	if !hasArgs {
		return printTaskSQLiteHelp()
	}

	sub := strings.ToLower(args[0])
	rest := args[1:]

	switch sub {
	case "init":
		return runTaskSQLiteInit(rest)
	case "add", "add-subtasks":
		return runTaskSQLiteAdd(rest)
	case "claim":
		return runTaskSQLiteClaim(rest)
	case "complete":
		return runTaskSQLiteComplete(rest)
	case "fail":
		return runTaskSQLiteFail(rest)
	case "log-action":
		return runTaskSQLiteLogAction(rest)
	case "status":
		return runTaskSQLiteStatus(rest)
	case "schema":
		return runTaskSQLiteSchema(rest)
	case "help", "-h", "--help":
		return printTaskSQLiteHelp()
	default:
		return fmt.Errorf("unknown task subcommand: %s", sub)
	}
}

func printTaskSQLiteHelp() error {
	fmt.Println("GitMap Native SQLite Agent Task Engine")
	fmt.Println()
	fmt.Println("Usage:")
	fmt.Println("  gitmap task init --name <title> [--budget <N>] [--db <path>]")
	fmt.Println("  gitmap task add --db <path> --tasks-json '<json>'")
	fmt.Println("  gitmap task claim --db <path> --agent <agent-role> [--subtask-id <id>]")
	fmt.Println("  gitmap task complete --db <path> --subtask-id <id> [--evidence <summary>]")
	fmt.Println("  gitmap task fail --db <path> --subtask-id <id> [--reason <summary>]")
	fmt.Println("  gitmap task log-action --db <path> --subtask-id <id> --agent <role> --action <act>")
	fmt.Println("  gitmap task status --db <path> [--json]")
	fmt.Println("  gitmap task schema")
	fmt.Println()
	return nil
}

func runTaskSQLiteInit(args []string) error {
	fs := flag.NewFlagSet("task init", flag.ContinueOnError)
	name := fs.String("name", "", "Task name or prompt description")
	budget := fs.Int("budget", 300, "Total steps budget")
	dbFlag := fs.String("db", "", "Path to SQLite database")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	hasName := strings.TrimSpace(*name) != ""
	if !hasName {
		return appfault.NewValidationError("missing required flag: --name")
	}

	repoRoot := FindRepoRoot()
	tempDir := ResolveAgentTempDir("")
	slug := Slugify(*name)

	var dbPath string
	var runDir string

	hasDb := strings.TrimSpace(*dbFlag) != ""
	if hasDb {
		dbPath = resolveTaskDbPath(repoRoot, *dbFlag)
		runDir = filepath.Dir(dbPath)
	} else {
		nextPrefix := resolveNextRunPrefix(tempDir)
		runDirName := fmt.Sprintf("%02d-%s", nextPrefix, slug)
		runDir = filepath.Join(tempDir, runDirName)
		dbPath = filepath.Join(runDir, Tier2TaskDBName)
	}

	mkErr := os.MkdirAll(runDir, 0755)
	if mkErr != nil {
		return appfault.WrapExecution(mkErr, "failed to create run directory")
	}

	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	initErr := InitTier2Schema(db)
	if initErr != nil {
		return initErr
	}

	now := NowISO()
	insertQuery := `
		INSERT INTO ParentTask (
			TaskName, TaskSlug, RunDirectory, Status, IsActive, HasCompleted,
			TotalStepsBudget, CurrentStep, CreatedAt, UpdatedAt
		) VALUES (?, ?, ?, 'ACTIVE', 1, 0, ?, 1, ?, ?);
	`
	relRunDir := toRelativePath(repoRoot, runDir)
	res, execErr := db.Exec(insertQuery, *name, slug, relRunDir, *budget, now, now)
	if execErr != nil {
		return appfault.WrapExecution(execErr, "failed to insert parent task")
	}

	parentId, _ := res.LastInsertId()
	relDbPath := toRelativePath(repoRoot, dbPath)

	out := map[string]any{
		"action":           "INITIALIZED",
		"parentTaskId":     parentId,
		"slug":             slug,
		"runDirectory":     relRunDir,
		"databasePath":     relDbPath,
		"totalStepsBudget": *budget,
		"status":           "ACTIVE",
	}

	return printJson(out)
}

func resolveNextRunPrefix(tempDir string) int {
	entries, err := os.ReadDir(tempDir)
	if err != nil {
		return 1
	}

	re := regexp.MustCompile(`^(\d+)-`)
	maxPrefix := 0

	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		matches := re.FindStringSubmatch(entry.Name())
		if len(matches) <= 1 {
			continue
		}
		val, convErr := strconv.Atoi(matches[1])
		if convErr == nil && val > maxPrefix {
			maxPrefix = val
		}
	}

	return maxPrefix + 1
}

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

func runTaskSQLiteStatus(args []string) error {
	fs := flag.NewFlagSet("task status", flag.ContinueOnError)
	dbFlag := fs.String("db", "", "Path to SQLite database")
	_ = fs.Bool("json", true, "Output results in JSON format")

	err := fs.Parse(args)
	if err != nil {
		return err
	}

	dbPath := resolveDbPath(*dbFlag)
	db, openErr := OpenSqliteDB(dbPath)
	if openErr != nil {
		return openErr
	}
	defer db.Close()

	var parentSlug, parentName, parentStatus string
	row := db.QueryRow("SELECT TaskSlug, TaskName, Status FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
	scanErr := row.Scan(&parentSlug, &parentName, &parentStatus)
	if scanErr != nil && scanErr != sql.ErrNoRows {
		return appfault.WrapExecution(scanErr, "failed to scan parent task record")
	}

	rows, queryErr := db.Query(`
		SELECT SubtaskId, TaskCode, Title, COALESCE(AssignedAgentRole, ''), Status, COALESCE(Evidence, '')
		FROM Subtask
		ORDER BY SubtaskId ASC;
	`)
	if queryErr != nil {
		return appfault.WrapExecution(queryErr, "failed to query subtasks")
	}
	defer rows.Close()

	type subtaskOut struct {
		SubtaskId         int64  `json:"SubtaskId"`
		TaskCode          string `json:"TaskCode"`
		Title             string `json:"Title"`
		AssignedAgentRole string `json:"AssignedAgentRole"`
		Status            string `json:"Status"`
		Evidence          string `json:"Evidence"`
	}

	var subtasks []subtaskOut
	counts := make(map[string]int)

	for rows.Next() {
		var s subtaskOut
		if err := rows.Scan(&s.SubtaskId, &s.TaskCode, &s.Title, &s.AssignedAgentRole, &s.Status, &s.Evidence); err == nil {
			subtasks = append(subtasks, s)
			counts[s.Status]++
		}
	}

	total := len(subtasks)
	done := counts["DONE"]
	var percent float64
	if total > 0 {
		percent = math.Round((float64(done)/float64(total)*100)*10) / 10
	}

	out := map[string]any{
		"taskSlug":        parentSlug,
		"taskName":        parentName,
		"status":          parentStatus,
		"totalSubtasks":   total,
		"completed":       done,
		"pending":         counts["PENDING"],
		"inProgress":      counts["IN_PROGRESS"],
		"failed":          counts["FAILED"],
		"percentComplete": percent,
		"subtasks":        subtasks,
	}

	return printJson(out)
}

func runTaskSQLiteSchema(args []string) error {
	schemaText := `
Schema Specification: agent-task.db (Tier 2 Run-Scoped SQLite Database)

1. ParentTask Table:
   - ParentTaskId: INTEGER PRIMARY KEY AUTOINCREMENT
   - TaskName: TEXT NOT NULL
   - TaskSlug: TEXT NOT NULL
   - RunDirectory: TEXT NOT NULL
   - Status: TEXT NOT NULL DEFAULT 'ACTIVE'
   - IsActive: INTEGER NOT NULL DEFAULT 1
   - HasCompleted: INTEGER NOT NULL DEFAULT 0
   - TotalStepsBudget: INTEGER NOT NULL DEFAULT 300
   - CurrentStep: INTEGER NOT NULL DEFAULT 1
   - Notes: TEXT NULL
   - CreatedAt: TEXT NOT NULL (ISO 8601 UTC)
   - UpdatedAt: TEXT NOT NULL (ISO 8601 UTC)

2. Subtask Table:
   - SubtaskId: INTEGER PRIMARY KEY AUTOINCREMENT
   - ParentTaskId: INTEGER NOT NULL (FK -> ParentTask.ParentTaskId)
   - TaskCode: TEXT NOT NULL (e.g. Task-01)
   - Title: TEXT NOT NULL
   - AssignedAgentRole: TEXT NULL (e.g. Worker 01)
   - OwnedFilesJson: TEXT NOT NULL DEFAULT '[]'
   - Status: TEXT NOT NULL DEFAULT 'PENDING' (PENDING, IN_PROGRESS, DONE, FAILED)
   - IsBlocked: INTEGER NOT NULL DEFAULT 0
   - HasCompleted: INTEGER NOT NULL DEFAULT 0
   - Evidence: TEXT NULL
   - CreatedAt: TEXT NOT NULL (ISO 8601 UTC)
   - UpdatedAt: TEXT NOT NULL (ISO 8601 UTC)

3. AgentActionLog Table:
   - ActionLogId: INTEGER PRIMARY KEY AUTOINCREMENT
   - SubtaskId: INTEGER NOT NULL (FK -> Subtask.SubtaskId)
   - AgentRole: TEXT NOT NULL
   - ActionType: TEXT NOT NULL (CLAIM, COMPLETE, FAIL, EXEC, etc.)
   - TargetFile: TEXT NULL
   - ActionDetails: TEXT NOT NULL
   - Status: TEXT NOT NULL DEFAULT 'IN_PROGRESS'
   - CreatedAt: TEXT NOT NULL (ISO 8601 UTC)
`
	fmt.Println(schemaText)
	return nil
}

func resolveTaskDbPath(repoRoot, raw string) string {
	cleaned := filepath.Clean(strings.TrimSpace(raw))
	if filepath.IsAbs(cleaned) {
		return cleaned
	}
	return filepath.Join(repoRoot, cleaned)
}

func resolveDbPath(flagVal string) string {
	clean := strings.TrimSpace(flagVal)
	repoRoot := FindRepoRoot()

	if clean != "" {
		return resolveTaskDbPath(repoRoot, clean)
	}

	tempDir := ResolveAgentTempDir("")
	latestRun, err := findLatestRunDir(tempDir)
	if err == nil && latestRun != "" {
		return filepath.Join(latestRun, Tier2TaskDBName)
	}

	return filepath.Join(tempDir, "01-agent", Tier2TaskDBName)
}

func toRelativePath(root, target string) string {
	rel, err := filepath.Rel(root, target)
	if err != nil {
		return filepath.ToSlash(target)
	}
	return filepath.ToSlash(rel)
}

func printJson(data any) error {
	bytes, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	fmt.Println(string(bytes))
	return nil
}

func readTasksRawJson(jsonStr, jsonFile string) (string, error) {
	raw := strings.TrimSpace(jsonStr)
	if raw != "" {
		return raw, nil
	}
	trimmedFile := strings.TrimSpace(jsonFile)
	if trimmedFile == "" {
		return "", nil
	}
	fileBytes, err := os.ReadFile(trimmedFile)
	if err != nil {
		return "", appfault.WrapExecution(err, "failed to read tasks-file")
	}
	return string(fileBytes), nil
}

func resolveParentTaskId(db *sql.DB, explicitId int64) (int64, error) {
	if explicitId > 0 {
		return explicitId, nil
	}
	return queryLatestParentTaskId(db)
}

func queryLatestParentTaskId(db *sql.DB) (int64, error) {
	var parentId int64
	row := db.QueryRow("SELECT ParentTaskId FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;")
	scanErr := row.Scan(&parentId)
	if scanErr != nil {
		return 0, appfault.WrapExecution(scanErr, "ParentTask record missing in database")
	}
	return parentId, nil
}
