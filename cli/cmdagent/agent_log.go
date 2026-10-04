package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type ActionType string

const (
	ActionTypeSearch   ActionType = "SEARCH"
	ActionTypeRead     ActionType = "READ"
	ActionTypeWrite    ActionType = "WRITE"
	ActionTypeExec     ActionType = "EXEC"
	ActionTypeLint     ActionType = "LINT"
	ActionTypeCheck    ActionType = "CHECK"
	ActionTypeClaim    ActionType = "CLAIM"
	ActionTypeStart    ActionType = "START"
	ActionTypeComplete ActionType = "COMPLETE"
	ActionTypeFail     ActionType = "FAIL"
	ActionTypeCrash    ActionType = "CRASH"
)

// LogOptions encapsulates command line parameters for agent telemetry logging.
type LogOptions struct {
	AgentRole    string
	SubtaskId    int64
	Action       ActionType
	TargetFile   string
	StartLine    int
	EndLine      int
	Query        string
	Details      string
	DurationMs   int64
	Status       string
	ErrorMessage string
	DirFlag      string
	TaskIdFlag   string
	DbFlag       string
}

// RunAgentLog executes gitmap agent log command.
func RunAgentLog(args []string) *appfault.AppError {
	opts, parseErr := parseLogFlags(args)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return parseErr
	}

	valErr := validateLogOptions(opts)
	hasValErr := valErr != nil
	if hasValErr {
		return valErr
	}

	taskDir, dirErr := resolveLogTargetDir(opts)
	hasDirErr := dirErr != nil
	if hasDirErr {
		return dirErr
	}

	actionLogId, execErr := executeLogRecord(taskDir, opts)
	hasExecErr := execErr != nil
	if hasExecErr {
		return execErr
	}

	syncToTaskDB(taskDir, opts)

	return renderLogResponse(actionLogId)
}

func parseLogFlags(args []string) (*LogOptions, *appfault.AppError) {
	opts := &LogOptions{
		AgentRole: "Worker",
		Action:    ActionTypeExec,
		Status:    "IN_PROGRESS",
	}

	for i := 0; i < len(args); i++ {
		arg := args[i]
		hasEq := strings.HasPrefix(arg, "--") && strings.Contains(arg, "=")
		if hasEq {
			parts := strings.SplitN(arg, "=", 2)
			parseFlagKeyValue(parts[0], parts[1], opts)
			continue
		}

		hasVal := i+1 < len(args) && !strings.HasPrefix(args[i+1], "-")
		if strings.HasPrefix(arg, "-") && hasVal {
			parseFlagKeyValue(arg, args[i+1], opts)
			i++
		}
	}

	return opts, nil
}

func parseFlagKeyValue(key, val string, opts *LogOptions) {
	normKey := strings.ToLower(strings.TrimLeft(key, "-"))
	switch normKey {
	case "agent", "a", "role":
		opts.AgentRole = val
	case "subtask", "subtask-id", "subtask_id", "s":
		opts.SubtaskId = parseSubtaskId(val)
	case "action", "act":
		opts.Action = NormalizeActionType(val)
	case "file", "target", "f":
		opts.TargetFile = filepath.ToSlash(val)
	case "start-line", "start_line", "start":
		opts.StartLine = parsePositiveInt(val)
	case "end-line", "end_line", "end":
		opts.EndLine = parsePositiveInt(val)
	case "query", "command", "cmd", "q":
		opts.Query = val
	case "details", "detail", "d":
		opts.Details = val
	case "duration-ms", "duration_ms", "duration":
		opts.DurationMs = parsePositiveInt64(val)
	case "status":
		opts.Status = strings.ToUpper(val)
	case "error", "err":
		opts.ErrorMessage = val
	case "dir":
		opts.DirFlag = val
	case "task-id", "task_id", "task":
		opts.TaskIdFlag = val
	case "db":
		opts.DbFlag = val
	}
}

func parseSubtaskId(val string) int64 {
	parsed, err := strconv.ParseInt(val, 10, 64)
	hasParsed := err == nil && parsed > 0
	if hasParsed {
		return parsed
	}

	var digits strings.Builder
	for _, r := range val {
		isDigit := r >= '0' && r <= '9'
		if isDigit {
			digits.WriteRune(r)
		}
	}

	extracted, extErr := strconv.ParseInt(digits.String(), 10, 64)
	hasExt := extErr == nil && extracted > 0
	if hasExt {
		return extracted
	}

	return 1
}

func parsePositiveInt(val string) int {
	num, err := strconv.Atoi(val)
	hasNum := err == nil && num >= 0
	if hasNum {
		return num
	}

	return 0
}

func parsePositiveInt64(val string) int64 {
	num, err := strconv.ParseInt(val, 10, 64)
	hasNum := err == nil && num >= 0
	if hasNum {
		return num
	}

	return 0
}

func validateLogOptions(opts *LogOptions) *appfault.AppError {
	hasRole := strings.TrimSpace(opts.AgentRole) != ""
	if !hasRole {
		opts.AgentRole = "Worker"
	}

	hasDetails := strings.TrimSpace(opts.Details) != ""
	if !hasDetails {
		opts.Details = fmt.Sprintf("Executing %s on %s", opts.Action, defaultTarget(opts.TargetFile))
	}

	hasSubtask := opts.SubtaskId > 0
	if !hasSubtask {
		return appfault.NewValidationError("missing required flag --subtask <id>")
	}

	return nil
}

func defaultTarget(file string) string {
	hasFile := strings.TrimSpace(file) != ""
	if hasFile {
		return file
	}

	return "workspace"
}

func resolveLogTargetDir(opts *LogOptions) (string, *appfault.AppError) {
	hasDbFlag := strings.TrimSpace(opts.DbFlag) != ""
	if hasDbFlag {
		cleanDb := filepath.Clean(opts.DbFlag)
		dir := filepath.Dir(cleanDb)
		isAgentsSubdir := filepath.Base(dir) == AgentsSubdirName
		if isAgentsSubdir {
			return filepath.Dir(dir), nil
		}
		return dir, nil
	}

	tempDir := ResolveAgentTempDir(opts.DirFlag)

	return ResolveTaskDir(tempDir, opts.TaskIdFlag)
}

func executeLogRecord(taskDir string, opts *LogOptions) (int64, *appfault.AppError) {
	agentDb, err := OpenTier3DB(taskDir, opts.AgentRole)
	hasErr := err != nil
	if hasErr {
		return 0, err
	}
	defer agentDb.Close()

	created := NowISO()
	query := `
	INSERT INTO AgentActionLog (
		SubtaskId, AgentRole, ActionType, TargetFile, StartLine, EndLine,
		QueryOrCommand, ActionDetails, DurationMs, Status, ErrorMessage, CreatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`
	res, execErr := agentDb.Exec(query,
		opts.SubtaskId, opts.AgentRole, string(opts.Action), opts.TargetFile,
		opts.StartLine, opts.EndLine, opts.Query, opts.Details,
		opts.DurationMs, opts.Status, opts.ErrorMessage, created,
	)
	hasExecErr := execErr != nil
	if hasExecErr {
		return 0, appfault.WrapExecution(execErr, "failed to record action log in agent split db")
	}

	rowId, idErr := res.LastInsertId()
	hasIdErr := idErr != nil
	if hasIdErr {
		return 1, nil
	}

	return rowId, nil
}

func syncToTaskDB(taskDir string, opts *LogOptions) {
	taskDbPath := filepath.Join(taskDir, Tier2TaskDBName)
	_, statErr := os.Stat(taskDbPath)
	hasTaskDb := statErr == nil
	if !hasTaskDb {
		return
	}

	db, openErr := OpenTier2DB(taskDir)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return
	}
	defer db.Close()

	now := NowISO()
	_, _ = db.Exec("UPDATE Subtask SET UpdatedAt = ? WHERE SubtaskId = ?;", now, opts.SubtaskId)
	_, _ = db.Exec(`
		INSERT INTO AgentActionLog (
			SubtaskId, AgentRole, ActionType, TargetFile, StartLine, EndLine,
			QueryOrCommand, ActionDetails, DurationMs, Status, ErrorMessage, CreatedAt
		) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, opts.SubtaskId, opts.AgentRole, string(opts.Action), opts.TargetFile,
		opts.StartLine, opts.EndLine, opts.Query, opts.Details,
		opts.DurationMs, opts.Status, opts.ErrorMessage, now)
}

func renderLogResponse(actionLogId int64) *appfault.AppError {
	resp := map[string]any{
		"status":      "LOGGED",
		"actionLogId": actionLogId,
	}

	payload, err := json.MarshalIndent(resp, "", "  ")
	hasErr := err != nil
	if hasErr {
		return appfault.WrapExecution(err, "failed to serialize log response")
	}

	fmt.Println(string(payload))

	return nil
}

// NormalizeActionType normalizes tool execution and raw names to standard ActionType.
func NormalizeActionType(raw string) ActionType {
	cleaned := strings.ToUpper(strings.TrimSpace(raw))
	switch cleaned {
	case "SEARCH", "SEARCH_WEB", "SEARCHER", "FIND", "GREP":
		return ActionTypeSearch
	case "READ", "VIEW_FILE", "READ_URL_CONTENT", "READ_MEMORY":
		return ActionTypeRead
	case "WRITE", "WRITE_TO_FILE", "REPLACE_FILE_CONTENT", "EDIT":
		return ActionTypeWrite
	case "EXEC", "RUN_COMMAND", "EXECUTE", "BASH", "PWSH":
		return ActionTypeExec
	case "LINT", "LINTER":
		return ActionTypeLint
	case "CHECK", "TEST", "VALIDATE":
		return ActionTypeCheck
	case "CLAIM":
		return ActionTypeClaim
	case "START":
		return ActionTypeStart
	case "COMPLETE", "DONE":
		return ActionTypeComplete
	case "FAIL", "FAILED", "BLOCKED":
		return ActionTypeFail
	case "CRASH", "CRASHED":
		return ActionTypeCrash
	default:
		return ActionType(cleaned)
	}
}
