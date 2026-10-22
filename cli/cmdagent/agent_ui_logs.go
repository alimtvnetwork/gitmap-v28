package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func parseLogFilter(r *http.Request, defaultTaskId string) ActionLogFilter {
	limit := parseQueryInt(r.URL.Query().Get("limit"), 100)
	taskId := r.URL.Query().Get("taskId")
	hasTaskId := len(strings.TrimSpace(taskId)) > 0
	if !hasTaskId {
		taskId = defaultTaskId
	}

	return ActionLogFilter{
		TaskId:     taskId,
		AgentRole:  strings.TrimSpace(r.URL.Query().Get("agent")),
		ActionType: strings.TrimSpace(r.URL.Query().Get("actionType")),
		SubtaskId:  strings.TrimSpace(r.URL.Query().Get("subtask")),
		Search:     strings.TrimSpace(r.URL.Query().Get("q")),
		Limit:      limit,
	}
}

func queryFilteredActionLogs(tempDir string, f ActionLogFilter) []map[string]any {
	taskDir, err := ResolveTaskDir(tempDir, f.TaskId)
	hasDir := err == nil && len(taskDir) > 0
	if !hasDir {
		return []map[string]any{}
	}

	dbPath := filepath.Join(taskDir, Tier2TaskDBName)
	db, openErr := OpenSqliteDB(dbPath)
	hasOpenErr := openErr != nil
	if hasOpenErr {
		return []map[string]any{}
	}
	defer db.Close()

	return executeLogsQuery(db, f)
}

func executeLogsQuery(db *sql.DB, f ActionLogFilter) []map[string]any {
	query, args := buildActionLogsSQL(f)
	rows, err := db.Query(query, args...)
	hasErr := err != nil
	if hasErr {
		return []map[string]any{}
	}
	defer rows.Close()

	return scanActionLogsRows(rows)
}

func buildActionLogsSQL(f ActionLogFilter) (string, []any) {
	var conditions []string
	var args []any
	conditions, args = appendLogFilterConditions(conditions, args, f)

	base := "SELECT ActionLogId, SubtaskId, AgentRole, ActionType, COALESCE(TargetFile, ''), ActionDetails, Status, CreatedAt FROM AgentActionLog"
	hasCond := len(conditions) > 0
	if hasCond {
		base += " WHERE " + strings.Join(conditions, " AND ")
	}
	base += " ORDER BY ActionLogId DESC LIMIT ?"
	args = append(args, f.Limit)

	return base, args
}

func appendLogFilterConditions(c []string, args []any, f ActionLogFilter) ([]string, []any) {
	hasRole := len(f.AgentRole) > 0 && f.AgentRole != "ALL"
	if hasRole {
		c = append(c, "LOWER(AgentRole) = LOWER(?)")
		args = append(args, f.AgentRole)
	}

	hasType := len(f.ActionType) > 0 && f.ActionType != "ALL"
	if hasType {
		c = append(c, "ActionType = ?")
		args = append(args, f.ActionType)
	}

	hasSub := len(f.SubtaskId) > 0 && f.SubtaskId != "ALL"
	if hasSub {
		c = append(c, "(SubtaskId = ? OR CAST(SubtaskId AS TEXT) = ?)")
		args = append(args, f.SubtaskId, f.SubtaskId)
	}

	return appendSearchCondition(c, args, f.Search)
}

func appendSearchCondition(c []string, args []any, search string) ([]string, []any) {
	hasSearch := len(search) > 0
	if hasSearch {
		pattern := "%" + search + "%"
		c = append(c, "(TargetFile LIKE ? OR ActionDetails LIKE ?)")
		args = append(args, pattern, pattern)
	}

	return c, args
}

func scanActionLogsRows(rows *sql.Rows) []map[string]any {
	var list []map[string]any
	for rows.Next() {
		var id int64
		var subId any
		var role, aType, file, details, status, cAt string
		scanErr := rows.Scan(&id, &subId, &role, &aType, &file, &details, &status, &cAt)
		hasScan := scanErr == nil
		if hasScan {
			item := formatActionLogMap(id, fmt.Sprint(subId), role, aType, file, details, status, cAt)
			list = append(list, item)
		}
	}

	return list
}

func deriveTouchType(actionType string) string {
	upper := strings.ToUpper(strings.TrimSpace(actionType))
	switch {
	case strings.Contains(upper, "READ") || strings.Contains(upper, "INSPECT"):
		return "READ"
	case strings.Contains(upper, "EDIT") || strings.Contains(upper, "MODIFY"):
		return "EDIT"
	case strings.Contains(upper, "CREATE") || strings.Contains(upper, "WRITE"):
		return "CREATED"
	default:
		return "OTHER"
	}
}

func formatActionLogMap(id int64, subId, role, aType, file, details, status, cAt string) map[string]any {
	return map[string]any{
		"actionLogId":    id,
		"subtaskId":      subId,
		"agentRole":      role,
		"actionType":     aType,
		"touchType":      deriveTouchType(aType),
		"targetFile":     file,
		"actionDetails":  details,
		"status":         status,
		"createdAt":      cAt,
		"queryOrCommand": "",
		"durationMs":     0,
	}
}

func writeJSON(w http.ResponseWriter, statusCode int, data any) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Cache-Control", "no-cache")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func parseQueryInt(raw string, defaultVal int) int {
	v, err := strconv.Atoi(raw)
	hasParsed := err == nil && v > 0
	if hasParsed {
		return v
	}

	return defaultVal
}
