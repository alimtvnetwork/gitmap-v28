package cmdagent

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"net/http"
	"path/filepath"
	"strconv"
	"strings"
)

func resolveTaskId(raw, defaultId string) string {
	clean := strings.TrimSpace(raw)
	hasId := len(clean) > 0
	if hasId {
		return clean
	}

	return defaultId
}

func parseQueryOffset(val string) int {
	clean := strings.TrimSpace(val)
	if clean == "" {
		return 0
	}

	offset, err := strconv.Atoi(clean)
	if err != nil || offset < 0 {
		return 0
	}

	return offset
}

func parseSortColumn(raw string) string {
	switch strings.ToLower(strings.TrimSpace(raw)) {
	case "createdat", "time", "date":
		return "CreatedAt"
	case "actiontype", "type":
		return "ActionType"
	case "agentrole", "role":
		return "AgentRole"
	case "subtaskid", "subtask":
		return "SubtaskId"
	default:
		return "ActionLogId"
	}
}

func parseSortOrder(raw string) string {
	if strings.EqualFold(strings.TrimSpace(raw), "asc") {
		return "ASC"
	}

	return "DESC"
}

func parseLogFilter(r *http.Request, defaultTaskId string) ActionLogFilter {
	return ActionLogFilter{
		TaskId:     resolveTaskId(r.URL.Query().Get("taskId"), defaultTaskId),
		AgentRole:  strings.TrimSpace(r.URL.Query().Get("agent")),
		ActionType: strings.TrimSpace(r.URL.Query().Get("actionType")),
		TouchType:  strings.TrimSpace(r.URL.Query().Get("touchType")),
		SubtaskId:  strings.TrimSpace(r.URL.Query().Get("subtask")),
		Search:     strings.TrimSpace(r.URL.Query().Get("q")),
		SortBy:     parseSortColumn(r.URL.Query().Get("sort")),
		Order:      parseSortOrder(r.URL.Query().Get("order")),
		Limit:      parseQueryInt(r.URL.Query().Get("limit"), 100),
		Offset:     parseQueryOffset(r.URL.Query().Get("offset")),
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

	return appendPaginationLimitOffset(base, args, f)
}

func appendPaginationLimitOffset(base string, args []any, f ActionLogFilter) (string, []any) {
	sortCol := parseSortColumn(f.SortBy)
	sortDir := parseSortOrder(f.Order)
	base += fmt.Sprintf(" ORDER BY %s %s LIMIT ?", sortCol, sortDir)
	args = append(args, f.Limit)
	if f.Offset > 0 {
		base += " OFFSET ?"
		args = append(args, f.Offset)
	}

	return base, args
}

func appendLogFilterConditions(c []string, args []any, f ActionLogFilter) ([]string, []any) {
	c, args = appendRoleCondition(c, args, f.AgentRole)
	c, args = appendActionTypeCondition(c, args, f.ActionType)
	c, args = appendTouchTypeCondition(c, args, f.TouchType)
	c, args = appendSubtaskCondition(c, args, f.SubtaskId)

	return appendSearchCondition(c, args, f.Search)
}

func appendRoleCondition(c []string, args []any, role string) ([]string, []any) {
	hasRole := len(role) > 0 && role != "ALL"
	if !hasRole {
		return c, args
	}

	c = append(c, "LOWER(AgentRole) = LOWER(?)")

	return c, append(args, role)
}

func appendActionTypeCondition(c []string, args []any, actionType string) ([]string, []any) {
	hasType := len(actionType) > 0 && actionType != "ALL"
	if !hasType {
		return c, args
	}

	c = append(c, "ActionType = ?")

	return c, append(args, actionType)
}

func appendSubtaskCondition(c []string, args []any, subtaskId string) ([]string, []any) {
	hasSub := len(subtaskId) > 0 && subtaskId != "ALL"
	if !hasSub {
		return c, args
	}

	c = append(c, "(SubtaskId = ? OR CAST(SubtaskId AS TEXT) = ?)")

	return c, append(args, subtaskId, subtaskId)
}

func appendTouchTypeCondition(c []string, args []any, touchType string) ([]string, []any) {
	hasTouch := len(touchType) > 0 && touchType != "ALL"
	if !hasTouch {
		return c, args
	}

	clause, param := resolveTouchTypeClause(strings.ToUpper(touchType))
	c = append(c, clause)
	if len(param) > 0 {
		args = append(args, param)
	}

	return c, args
}

func resolveTouchTypeClause(upper string) (string, string) {
	switch upper {
	case "READ":
		return "UPPER(ActionType) IN ('READ', 'READ_FILE', 'INSPECT')", ""
	case "EDIT":
		return "UPPER(ActionType) IN ('EDIT', 'EDIT_FILE', 'MODIFY')", ""
	case "CREATED":
		return "UPPER(ActionType) IN ('CREATE', 'WRITE_FILE', 'NEW')", ""
	default:
		return "UPPER(ActionType) = UPPER(?)", upper
	}
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
	w.Header().Set(constants.HeaderContentType, constants.MediaTypeJSON)
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
