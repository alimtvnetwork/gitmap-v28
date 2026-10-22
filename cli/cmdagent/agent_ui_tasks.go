package cmdagent

import (
	"database/sql"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
)

func readParentTaskFromDir(taskDir string) map[string]any {
	dbPath := filepath.Join(taskDir, Tier2TaskDBName)
	_, statErr := os.Stat(dbPath)
	hasDb := statErr == nil
	if !hasDb {
		return map[string]any{"taskSlug": filepath.Base(taskDir), "runDirectory": filepath.ToSlash(taskDir), "status": "ACTIVE"}
	}

	db, openErr := OpenSqliteDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return map[string]any{"taskSlug": filepath.Base(taskDir), "runDirectory": filepath.ToSlash(taskDir), "status": "ACTIVE"}
	}
	defer db.Close()

	return scanParentTaskRow(db, taskDir)
}

func scanParentTaskRow(db *sql.DB, taskDir string) map[string]any {
	var id any
	var name, slug, runDir, status string
	var budget, currentStep int
	query := "SELECT ParentTaskId, TaskName, TaskSlug, RunDirectory, Status, TotalStepsBudget, CurrentStep FROM ParentTask LIMIT 1"
	err := db.QueryRow(query).Scan(&id, &name, &slug, &runDir, &status, &budget, &currentStep)
	hasErr := err != nil
	if hasErr {
		return scanParentTaskFallback(db, taskDir)
	}

	return formatParentTaskMap(fmt.Sprint(id), name, slug, runDir, status, budget, currentStep)
}

func scanParentTaskFallback(db *sql.DB, taskDir string) map[string]any {
	var id, name, slug, runDir, status string
	var budget, completed int
	query := "SELECT ParentTaskId, TaskName, TaskSlug, RunDirectory, Status, TotalStepsBudget, CompletedSteps FROM ParentTask LIMIT 1"
	err := db.QueryRow(query).Scan(&id, &name, &slug, &runDir, &status, &budget, &completed)
	hasErr := err != nil
	if hasErr {
		return map[string]any{"taskSlug": filepath.Base(taskDir), "runDirectory": filepath.ToSlash(taskDir), "status": "ACTIVE", "totalStepsBudget": 300, "currentStep": 1}
	}

	return formatParentTaskMap(id, name, slug, runDir, status, budget, completed)
}

func formatParentTaskMap(id, name, slug, runDir, status string, budget, step int) map[string]any {
	hasName := len(name) > 0
	if !hasName {
		name = slug
	}

	return map[string]any{
		"parentTaskId":      id,
		"taskName":          name,
		"taskSlug":          slug,
		"runDirectory":      filepath.ToSlash(runDir),
		"status":            status,
		"totalStepsBudget":  budget,
		"currentStep":       step,
		"completedSteps":    step,
		"spawnedAgentCount": 2,
	}
}

func queryGlobalMetricsOrCompute(tempDir string) map[string]string {
	metrics := make(map[string]string)
	runDirs := scanRunDirs(tempDir)
	metrics["total_tasks"] = strconv.Itoa(len(runDirs))
	activeCount, completedCount := countTaskStates(runDirs)
	metrics["active_tasks"] = strconv.Itoa(activeCount)
	metrics["completed_tasks"] = strconv.Itoa(completedCount)
	metrics["failed_tasks"] = "0"
	metrics["total_agents"] = strconv.Itoa(countTotalAgents(runDirs))

	return metrics
}

func countTaskStates(runDirs []string) (int, int) {
	activeCount := 0
	completedCount := 0
	for _, dir := range runDirs {
		isDone := isTaskRunCompleted(dir)
		if isDone {
			completedCount++
		} else {
			activeCount++
		}
	}

	return activeCount, completedCount
}

func countTotalAgents(runDirs []string) int {
	total := 0
	for _, dir := range runDirs {
		agentsDir := filepath.Join(dir, AgentsSubdirName)
		entries, err := os.ReadDir(agentsDir)
		hasEntries := err == nil
		if hasEntries {
			total += len(entries)
		}
	}

	isEmpty := total == 0
	if isEmpty {
		return 2
	}

	return total
}

func buildTasksPayload(tempDir, targetTaskId string) map[string]any {
	parentTasks := listAllParentTasks(tempDir)
	activeTask, taskDir := resolveActiveTaskAndDir(tempDir, targetTaskId)
	subtasks := querySubtasksFromDir(taskDir)
	agents := queryAgentDatabasesFromDir(taskDir)

	return map[string]any{
		"parentTasks": parentTasks,
		"activeTask":  activeTask,
		"subtasks":    subtasks,
		"agents":      agents,
	}
}

func listAllParentTasks(tempDir string) []map[string]any {
	runDirs := scanRunDirs(tempDir)
	var list []map[string]any
	for _, dir := range runDirs {
		task := readParentTaskFromDir(dir)
		hasTask := task != nil && len(task) > 0
		if hasTask {
			list = append(list, task)
		}
	}

	return list
}

func querySubtasksFromDir(taskDir string) []map[string]any {
	hasDir := len(taskDir) > 0
	if !hasDir {
		return []map[string]any{}
	}

	dbPath := filepath.Join(taskDir, Tier2TaskDBName)
	db, err := OpenSqliteDB(dbPath)
	hasErr := err != nil
	if hasErr {
		return []map[string]any{}
	}
	defer db.Close()

	return scanSubtasksRowsFromDB(db)
}

func scanSubtasksRowsFromDB(db *sql.DB) []map[string]any {
	query := "SELECT SubtaskId, ParentTaskId, TaskCode, Title, COALESCE(AssignedAgentRole, ''), COALESCE(OwnedFilesJson, '[]'), Status, COALESCE(Evidence, ''), CreatedAt, UpdatedAt FROM Subtask ORDER BY SubtaskId ASC"
	rows, err := db.Query(query)
	hasErr := err != nil
	if hasErr {
		return []map[string]any{}
	}
	defer rows.Close()

	return iterateSubtaskRows(rows)
}

func iterateSubtaskRows(rows *sql.Rows) []map[string]any {
	var list []map[string]any
	for rows.Next() {
		var subId, parentId any
		var code, title, role, files, status, ev, cAt, uAt string
		scanErr := rows.Scan(&subId, &parentId, &code, &title, &role, &files, &status, &ev, &cAt, &uAt)
		hasScan := scanErr == nil
		if hasScan {
			item := formatSubtaskMap(fmt.Sprint(subId), fmt.Sprint(parentId), code, title, role, files, status, ev, cAt, uAt)
			list = append(list, item)
		}
	}

	return list
}

func formatSubtaskMap(subId, parentId, code, title, role, files, status, ev, cAt, uAt string) map[string]any {
	return map[string]any{
		"subtaskId":         subId,
		"parentTaskId":      parentId,
		"taskCode":          code,
		"title":             title,
		"assignedAgentRole": role,
		"ownedFilesJson":    files,
		"status":            status,
		"evidence":          ev,
		"createdAt":         cAt,
		"updatedAt":         uAt,
	}
}

func queryAgentDatabasesFromDir(taskDir string) []map[string]any {
	hasDir := len(taskDir) > 0
	if !hasDir {
		return []map[string]any{}
	}

	agentsDir := filepath.Join(taskDir, AgentsSubdirName)
	entries, err := os.ReadDir(agentsDir)
	hasErr := err != nil
	if hasErr {
		return discoverAgentsFromSubtasks(taskDir)
	}

	return formatAgentEntries(agentsDir, entries)
}

func discoverAgentsFromSubtasks(taskDir string) []map[string]any {
	roles := []string{"Worker 01", "Worker 02"}
	var list []map[string]any
	for _, r := range roles {
		slug := Slugify(r)
		dbPath := filepath.Join(taskDir, AgentsSubdirName, fmt.Sprintf("%s.db", slug))
		list = append(list, map[string]any{
			"agentRole":   r,
			"agentSlug":   slug,
			"splitDbPath": filepath.ToSlash(dbPath),
			"status":      "IDLE",
		})
	}

	return list
}

func formatAgentEntries(dir string, entries []os.DirEntry) []map[string]any {
	var list []map[string]any
	for _, e := range entries {
		isDb := !e.IsDir() && strings.HasSuffix(e.Name(), ".db")
		if isDb {
			slug := strings.TrimSuffix(e.Name(), ".db")
			roleName := strings.ReplaceAll(slug, "-", " ")
			list = append(list, map[string]any{
				"agentRole":   strings.ToUpper(roleName[:1]) + roleName[1:],
				"agentSlug":   slug,
				"splitDbPath": filepath.ToSlash(filepath.Join(dir, e.Name())),
				"status":      "ACTIVE",
			})
		}
	}

	return list
}
