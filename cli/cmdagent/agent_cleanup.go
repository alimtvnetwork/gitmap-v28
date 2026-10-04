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

// CleanupOptions captures parameters for clear, reset, and temp-clear commands.
type CleanupOptions struct {
	TaskIdFlag  string
	DirFlag     string
	IsConfirmed bool
	IsJson      bool
}

// RunAgentClear cleans a specified parent task or all completed tasks in temp dir.
func RunAgentClear(args []string) *appfault.AppError {
	opts, parseErr := parseCleanupFlags(args)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return parseErr
	}

	tempDir := ResolveAgentTempDir(opts.DirFlag)
	hasTaskId := strings.TrimSpace(opts.TaskIdFlag) != ""
	if hasTaskId {
		return executeClearSpecific(tempDir, opts.TaskIdFlag, opts.IsJson)
	}

	return executeClearCompleted(tempDir, opts.IsJson)
}

// RunAgentReset clears agent split databases and resets tables while preserving run dirs.
func RunAgentReset(args []string) *appfault.AppError {
	opts, parseErr := parseCleanupFlags(args)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return parseErr
	}

	tempDir := ResolveAgentTempDir(opts.DirFlag)

	return executeReset(tempDir, opts.IsJson)
}

// RunAgentTempClear completely purges the agent temp directory.
func RunAgentTempClear(args []string) *appfault.AppError {
	opts, parseErr := parseCleanupFlags(args)
	hasParseErr := parseErr != nil
	if hasParseErr {
		return parseErr
	}

	tempDir := ResolveAgentTempDir(opts.DirFlag)

	return executeTempClear(tempDir, opts.IsJson)
}

func parseCleanupFlags(args []string) (*CleanupOptions, *appfault.AppError) {
	opts := &CleanupOptions{}
	for i := 0; i < len(args); i++ {
		arg := args[i]
		hasEq := strings.HasPrefix(arg, "--") && strings.Contains(arg, "=")
		if hasEq {
			parts := strings.SplitN(arg, "=", 2)
			parseCleanupFlagValue(parts[0], parts[1], opts)
			continue
		}

		hasVal := i+1 < len(args) && !strings.HasPrefix(args[i+1], "-")
		if strings.HasPrefix(arg, "-") && hasVal {
			parseCleanupFlagValue(arg, args[i+1], opts)
			i++
			continue
		}

		parseCleanupFlagValue(arg, "", opts)
	}

	return opts, nil
}

func parseCleanupFlagValue(key, val string, opts *CleanupOptions) {
	normKey := strings.ToLower(strings.TrimLeft(key, "-"))
	switch normKey {
	case "task-id", "task_id", "task", "t":
		opts.TaskIdFlag = val
	case "dir", "d":
		opts.DirFlag = val
	case "y", "yes", "f", "force":
		opts.IsConfirmed = true
	case "json", "j":
		opts.IsJson = true
	}
}

func executeClearSpecific(tempDir, taskId string, isJson bool) *appfault.AppError {
	taskDir, findErr := ResolveTaskDir(tempDir, taskId)
	hasFindErr := findErr != nil
	if hasFindErr {
		return findErr
	}

	removeErr := os.RemoveAll(taskDir)
	hasRemoveErr := removeErr != nil
	if hasRemoveErr {
		return appfault.WrapExecution(removeErr, "failed to delete task run directory")
	}

	removeRegistryEntries(tempDir, taskId)
	payload := map[string]any{"status": "CLEARED", "taskId": taskId, "removedDir": filepath.ToSlash(taskDir)}

	return renderCleanupSuccess("CLEARED", payload, isJson)
}

func executeClearCompleted(tempDir string, isJson bool) *appfault.AppError {
	runDirs := scanRunDirs(tempDir)
	var cleared []string

	for _, dir := range runDirs {
		isDone := isTaskRunCompleted(dir)
		if isDone {
			_ = os.RemoveAll(dir)
			removeRegistryEntries(tempDir, filepath.Base(dir))
			cleared = append(cleared, filepath.ToSlash(dir))
		}
	}

	payload := map[string]any{"status": "CLEARED", "clearedTasksCount": len(cleared), "clearedDirs": cleared}

	return renderCleanupSuccess("CLEARED", payload, isJson)
}

func isTaskRunCompleted(taskDir string) bool {
	dbPath := filepath.Join(taskDir, Tier2TaskDBName)
	_, statErr := os.Stat(dbPath)
	hasDb := statErr == nil
	if !hasDb {
		return false
	}

	db, err := OpenSqliteDB(dbPath)
	hasErr := err != nil
	if hasErr {
		return false
	}
	defer db.Close()

	var hasCompleted int
	_ = db.QueryRow("SELECT HasCompleted FROM ParentTask ORDER BY ParentTaskId DESC LIMIT 1;").Scan(&hasCompleted)
	isParentDone := hasCompleted == 1
	if isParentDone {
		return true
	}

	var pendingCount, inProgressCount int
	_ = db.QueryRow("SELECT COUNT(*) FROM Subtask WHERE Status IN ('PENDING', 'IN_PROGRESS');").Scan(&pendingCount)

	return pendingCount == 0 && inProgressCount == 0
}

func removeRegistryEntries(tempDir, taskId string) {
	masterPath := filepath.Join(tempDir, Tier1MasterDBName)
	_, statErr := os.Stat(masterPath)
	hasMaster := statErr == nil
	if !hasMaster {
		return
	}

	db, err := OpenSqliteDB(masterPath)
	hasErr := err != nil
	if hasErr {
		return
	}
	defer db.Close()

	_, _ = db.Exec("DELETE FROM ParentTaskRegistry WHERE TaskSlug LIKE ? OR RunDirectory LIKE ?;", "%"+taskId+"%", "%"+taskId+"%")
}

func executeReset(tempDir string, isJson bool) *appfault.AppError {
	resetAiAgentsMasterDB(tempDir)
	runDirs := scanRunDirs(tempDir)
	for _, dir := range runDirs {
		resetTaskRunDir(dir)
	}

	payload := map[string]any{"status": "RESET", "resetTasksCount": len(runDirs)}

	return renderCleanupSuccess("RESET", payload, isJson)
}

func resetAiAgentsMasterDB(tempDir string) {
	masterPath := filepath.Join(tempDir, Tier1MasterDBName)
	_, statErr := os.Stat(masterPath)
	hasMaster := statErr == nil
	if !hasMaster {
		return
	}

	db, err := OpenSqliteDB(masterPath)
	hasErr := err != nil
	if hasErr {
		return
	}
	defer db.Close()

	_, _ = db.Exec("DELETE FROM ParentTaskRegistry;")
	_, _ = db.Exec("DELETE FROM AgentRegistry;")
	_, _ = db.Exec("DELETE FROM GlobalLifecycleMetrics;")
}

func resetTaskRunDir(taskDir string) {
	agentsDir := filepath.Join(taskDir, AgentsSubdirName)
	_ = os.RemoveAll(agentsDir)
	_ = os.MkdirAll(agentsDir, 0755)

	dbPath := filepath.Join(taskDir, Tier2TaskDBName)
	_, statErr := os.Stat(dbPath)
	hasDb := statErr == nil
	if !hasDb {
		return
	}

	db, err := OpenSqliteDB(dbPath)
	hasErr := err != nil
	if hasErr {
		return
	}
	defer db.Close()

	_, _ = db.Exec("UPDATE Subtask SET Status = 'PENDING', HasCompleted = 0, IsBlocked = 0, AssignedAgentRole = NULL, Evidence = NULL;")
	_, _ = db.Exec("DELETE FROM AgentActionLog;")
}

func executeTempClear(tempDir string, isJson bool) *appfault.AppError {
	removeErr := os.RemoveAll(tempDir)
	hasRemoveErr := removeErr != nil
	if hasRemoveErr {
		return appfault.WrapExecution(removeErr, "failed to delete agent temp directory")
	}

	mkdirErr := os.MkdirAll(tempDir, 0755)
	hasMkdirErr := mkdirErr != nil
	if hasMkdirErr {
		return appfault.WrapExecution(mkdirErr, "failed to recreate empty agent temp directory")
	}

	payload := map[string]any{"status": "TEMP_CLEARED", "tempDir": filepath.ToSlash(tempDir)}

	return renderCleanupSuccess("TEMP_CLEARED", payload, isJson)
}

func renderCleanupSuccess(action string, payload any, isJson bool) *appfault.AppError {
	if isJson {
		data, err := json.MarshalIndent(payload, "", "  ")
		hasErr := err != nil
		if hasErr {
			return appfault.WrapExecution(err, "failed to format cleanup json response")
		}
		fmt.Println(string(data))
		return nil
	}

	fmt.Printf("%s✓ %s completed successfully.%s\n", constants.ColorGreen, action, constants.ColorReset)

	return nil
}
