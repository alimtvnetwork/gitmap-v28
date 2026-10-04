package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// CrashAutopsy records forensic details of an incomplete or crashed subtask.
type CrashAutopsy struct {
	SubtaskId         int64  `json:"subtaskId"`
	TaskCode          string `json:"taskCode"`
	Title             string `json:"title"`
	AssignedAgentRole string `json:"assignedAgentRole"`
	LastActionType    string `json:"lastActionType"`
	LastTargetFile    string `json:"lastTargetFile"`
	LastActionDetails string `json:"lastActionDetails"`
	LastTimestamp     string `json:"lastTimestamp"`
	Diagnosis         string `json:"diagnosis"`
}

// SubtaskCounts holds aggregate counts across task lifecycle states.
type SubtaskCounts struct {
	Done       int `json:"done"`
	Pending    int `json:"pending"`
	InProgress int `json:"inProgress"`
	Failed     int `json:"failed"`
}

// DiagnosticReport contains comprehensive task forensic and autopsy data.
type DiagnosticReport struct {
	HasCrashesDetected bool           `json:"hasCrashesDetected"`
	CrashedAgents      []CrashAutopsy `json:"crashedAgents"`
	Counts             SubtaskCounts  `json:"counts"`
	TaskDirectory      string         `json:"taskDirectory,omitempty"`
}

// DiagnoseOptions holds flags for crashed and diagnose commands.
type DiagnoseOptions struct {
	TaskIdFlag string
	DirFlag    string
	DbFlag     string
	IsJson     bool
	IsAll      bool
}

// RunAgentCrashed checks for crashed or abandoned agents and outputs an autopsy report.
func RunAgentCrashed(args []string) *appfault.AppError {
	opts, parseErr := parseDiagnoseFlags(args)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return parseErr
	}

	taskDir, dirErr := resolveDiagnoseTaskDir(opts)
	hasDirErr := dirErr != nil
	if hasDirErr {
		return dirErr
	}

	report, diagErr := buildDiagnosticReport(taskDir)
	hasDiagErr := diagErr != nil
	if hasDiagErr {
		return diagErr
	}

	return renderCrashedOutput(report, opts.IsJson)
}

// RunAgentDiagnose performs deep forensic inspection of task status and crash autopsies.
func RunAgentDiagnose(args []string) *appfault.AppError {
	opts, parseErr := parseDiagnoseFlags(args)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return parseErr
	}

	taskDir, dirErr := resolveDiagnoseTaskDir(opts)
	hasDirErr := dirErr != nil
	if hasDirErr {
		return dirErr
	}

	report, diagErr := buildDiagnosticReport(taskDir)
	hasDiagErr := diagErr != nil
	if hasDiagErr {
		return diagErr
	}

	return renderDiagnoseOutput(report, opts.IsJson)
}

func parseDiagnoseFlags(args []string) (*DiagnoseOptions, *appfault.AppError) {
	opts := &DiagnoseOptions{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		hasEq := strings.HasPrefix(arg, "--") && strings.Contains(arg, "=")
		if hasEq {
			parts := strings.SplitN(arg, "=", 2)
			parseDiagnoseFlagValue(parts[0], parts[1], opts)
			continue
		}

		hasVal := i+1 < len(args) && !strings.HasPrefix(args[i+1], "-")
		if strings.HasPrefix(arg, "-") && hasVal {
			parseDiagnoseFlagValue(arg, args[i+1], opts)
			i++
			continue
		}

		parseDiagnoseFlagValue(arg, "", opts)
	}

	return opts, nil
}

func parseDiagnoseFlagValue(key, val string, opts *DiagnoseOptions) {
	normKey := strings.ToLower(strings.TrimLeft(key, "-"))
	switch normKey {
	case "task-id", "task_id", "task", "t":
		opts.TaskIdFlag = val
	case "dir", "d":
		opts.DirFlag = val
	case "db":
		opts.DbFlag = val
	case "json", "j":
		opts.IsJson = true
	case "all", "a":
		opts.IsAll = true
	}
}

func resolveDiagnoseTaskDir(opts *DiagnoseOptions) (string, *appfault.AppError) {
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

func buildDiagnosticReport(taskDir string) (*DiagnosticReport, *appfault.AppError) {
	db, err := OpenTier2DB(taskDir)
	hasErr := err != nil
	if hasErr {
		return nil, err
	}
	defer db.Close()

	crashes, crashErr := queryInProgressAutopsies(db, taskDir)
	hasCrashErr := crashErr != nil
	if hasCrashErr {
		return nil, crashErr
	}

	counts, countErr := querySubtaskCounts(db)
	hasCountErr := countErr != nil
	if hasCountErr {
		return nil, countErr
	}

	hasCrashes := len(crashes) > 0
	return &DiagnosticReport{
		HasCrashesDetected: hasCrashes,
		CrashedAgents:      crashes,
		Counts:             counts,
		TaskDirectory:      taskDir,
	}, nil
}

func queryInProgressAutopsies(db *sql.DB, taskDir string) ([]CrashAutopsy, *appfault.AppError) {
	query := `
	SELECT SubtaskId, TaskCode, Title, COALESCE(AssignedAgentRole, 'Worker')
	FROM Subtask
	WHERE Status = 'IN_PROGRESS'
	ORDER BY SubtaskId ASC;
	`
	rows, err := db.Query(query)
	hasErr := err != nil
	if hasErr {
		return nil, appfault.WrapExecution(err, "failed to query in-progress subtasks")
	}
	defer rows.Close()

	var autopsies []CrashAutopsy
	for rows.Next() {
		var subtaskId int64
		var taskCode, title, role string
		_ = rows.Scan(&subtaskId, &taskCode, &title, &role)

		actionType, targetFile, details, timestamp := fetchLastAction(db, taskDir, role, subtaskId)
		autopsy := constructAutopsy(subtaskId, taskCode, title, role, actionType, targetFile, details, timestamp)
		autopsies = append(autopsies, autopsy)
	}

	return autopsies, nil
}

func fetchLastAction(db *sql.DB, taskDir, role string, subtaskId int64) (string, string, string, string) {
	t3Type, t3File, t3Details, t3Time, hasT3 := fetchTier3LastAction(taskDir, role, subtaskId)
	if hasT3 {
		return t3Type, t3File, t3Details, t3Time
	}

	var actionType, targetFile, details, timestamp sql.NullString
	query := `
	SELECT ActionType, TargetFile, ActionDetails, CreatedAt
	FROM AgentActionLog
	WHERE SubtaskId = ?
	ORDER BY ActionLogId DESC LIMIT 1;
	`
	_ = db.QueryRow(query, subtaskId).Scan(&actionType, &targetFile, &details, &timestamp)

	return actionType.String, targetFile.String, details.String, timestamp.String
}

func fetchTier3LastAction(taskDir, role string, subtaskId int64) (string, string, string, string, bool) {
	agentSlug := Slugify(role)
	agentDbPath := filepath.Join(taskDir, AgentsSubdirName, fmt.Sprintf("%s.db", agentSlug))
	_, statErr := os.Stat(agentDbPath)
	hasFile := statErr == nil
	if !hasFile {
		return "", "", "", "", false
	}

	db, openErr := OpenSqliteDB(agentDbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return "", "", "", "", false
	}
	defer db.Close()

	var actionType, targetFile, details, timestamp sql.NullString
	query := `
	SELECT ActionType, TargetFile, ActionDetails, CreatedAt
	FROM AgentActionLog
	WHERE SubtaskId = ?
	ORDER BY ActionLogId DESC LIMIT 1;
	`
	err := db.QueryRow(query, subtaskId).Scan(&actionType, &targetFile, &details, &timestamp)
	hasResult := err == nil && actionType.Valid
	if !hasResult {
		return "", "", "", "", false
	}

	return actionType.String, targetFile.String, details.String, timestamp.String, true
}

func constructAutopsy(subtaskId int64, code, title, role, actionType, targetFile, details, timestamp string) CrashAutopsy {
	hasAction := strings.TrimSpace(actionType) != ""
	if !hasAction {
		actionType = "UNKNOWN"
	}

	hasTarget := strings.TrimSpace(targetFile) != ""
	if !hasTarget {
		targetFile = "workspace"
	}

	diag := fmt.Sprintf(
		"Agent '%s' was in progress on %s ('%s'). Last logged action: '%s' targeting '%s'. Agent terminated before completing. Likely caused by tool execution crash, OS file lock, or turn timeout.",
		role, code, title, actionType, targetFile,
	)

	return CrashAutopsy{
		SubtaskId:         subtaskId,
		TaskCode:          code,
		Title:             title,
		AssignedAgentRole: role,
		LastActionType:    actionType,
		LastTargetFile:    targetFile,
		LastActionDetails: details,
		LastTimestamp:     timestamp,
		Diagnosis:         diag,
	}
}

func querySubtaskCounts(db *sql.DB) (SubtaskCounts, *appfault.AppError) {
	var counts SubtaskCounts
	_ = db.QueryRow("SELECT COUNT(*) FROM Subtask WHERE Status = 'DONE';").Scan(&counts.Done)
	_ = db.QueryRow("SELECT COUNT(*) FROM Subtask WHERE Status = 'PENDING';").Scan(&counts.Pending)
	_ = db.QueryRow("SELECT COUNT(*) FROM Subtask WHERE Status = 'IN_PROGRESS';").Scan(&counts.InProgress)
	_ = db.QueryRow("SELECT COUNT(*) FROM Subtask WHERE Status = 'FAILED';").Scan(&counts.Failed)

	return counts, nil
}

func renderCrashedOutput(report *DiagnosticReport, isJson bool) *appfault.AppError {
	if isJson {
		resp := map[string]any{
			"hasCrashesDetected": report.HasCrashesDetected,
			"crashedAgents":      report.CrashedAgents,
		}
		data, err := json.MarshalIndent(resp, "", "  ")
		hasErr := err != nil
		if hasErr {
			return appfault.WrapExecution(err, "failed to format crashed json output")
		}
		fmt.Println(string(data))
		return nil
	}

	if !report.HasCrashesDetected {
		fmt.Printf("%s✓ No crashed or abandoned agents detected.%s All in-flight subtasks healthy.\n", constants.ColorGreen, constants.ColorReset)
		return nil
	}

	for _, crash := range report.CrashedAgents {
		renderTerminalAutopsyCard(crash)
	}

	return nil
}

func renderDiagnoseOutput(report *DiagnosticReport, isJson bool) *appfault.AppError {
	if isJson {
		data, err := json.MarshalIndent(report, "", "  ")
		hasErr := err != nil
		if hasErr {
			return appfault.WrapExecution(err, "failed to format diagnose json output")
		}
		fmt.Println(string(data))
		return nil
	}

	renderTerminalDiagnoseSummary(report)
	for _, crash := range report.CrashedAgents {
		renderTerminalAutopsyCard(crash)
	}

	return nil
}

func renderTerminalDiagnoseSummary(report *DiagnosticReport) {
	fmt.Printf("\n%s=== AGENT FLEET DIAGNOSTICS ===%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("Task Directory: %s%s%s\n", constants.ColorYellow, report.TaskDirectory, constants.ColorReset)
	fmt.Printf("Subtask Counts: [DONE: %s%d%s] [PENDING: %s%d%s] [IN_PROGRESS: %s%d%s] [FAILED: %s%d%s]\n",
		constants.ColorGreen, report.Counts.Done, constants.ColorReset,
		constants.ColorCyan, report.Counts.Pending, constants.ColorReset,
		constants.ColorYellow, report.Counts.InProgress, constants.ColorReset,
		constants.ColorRed, report.Counts.Failed, constants.ColorReset)

	hasCrashes := report.HasCrashesDetected
	if hasCrashes {
		fmt.Printf("Crash Status:   %s%d CRASH(ES) DETECTED%s\n\n", constants.ColorRed, len(report.CrashedAgents), constants.ColorReset)
	} else {
		fmt.Printf("Crash Status:   %sCLEAN - NO CRASHES DETECTED%s\n\n", constants.ColorGreen, constants.ColorReset)
	}
}

func renderTerminalAutopsyCard(crash CrashAutopsy) {
	fmt.Printf("%s╔══════════════════════════════════════════════════════════════════════════════╗%s\n", constants.ColorRed, constants.ColorReset)
	fmt.Printf("%s║ CRASH AUTOPSY: %s | %s (%s)%s\n", constants.ColorRed, crash.AssignedAgentRole, crash.TaskCode, crash.Title, constants.ColorReset)
	fmt.Printf("%s╠══════════════════════════════════════════════════════════════════════════════╣%s\n", constants.ColorRed, constants.ColorReset)
	fmt.Printf("  • Subtask ID:    %d\n", crash.SubtaskId)
	fmt.Printf("  • Assigned Role: %s%s%s\n", constants.ColorYellow, crash.AssignedAgentRole, constants.ColorReset)
	fmt.Printf("  • Last Action:   %s%s%s\n", constants.ColorCyan, crash.LastActionType, constants.ColorReset)
	fmt.Printf("  • Target File:   %s\n", crash.LastTargetFile)
	fmt.Printf("  • Action Detail: %s\n", crash.LastActionDetails)
	fmt.Printf("  • Timestamp:     %s\n", crash.LastTimestamp)
	fmt.Printf("  • Diagnosis:     %s%s%s\n", constants.ColorRed, crash.Diagnosis, constants.ColorReset)
	fmt.Printf("%s╚══════════════════════════════════════════════════════════════════════════════╝%s\n\n", constants.ColorRed, constants.ColorReset)
}
