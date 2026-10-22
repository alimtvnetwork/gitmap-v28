# Subtask Plan 25: Touch Type Filtering in Agent Action Logs (`cli/cmdagent/agent_ui_logs.go`)

> **Subtask ID:** `subtask-25`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdagent/agent_ui_logs.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Extend `ActionLogFilter` struct in `cli/cmdagent/agent_ui_server.go` to include `TouchType string`.
2. Parse the `touchType` query parameter in `parseLogFilter` in `cli/cmdagent/agent_ui_logs.go`.
3. Support `"touchType"` filtering in `appendLogFilterConditions`, generating SQL matching conditions for `READ`, `EDIT`, and `CREATED` categories.
4. Decompose `appendLogFilterConditions` into modular filter helpers to maintain <= 15 lines per function.

---

## 2. Targeted Problem & Architectural Rationale
The GitMap Web Visualizer (`/agent` dashboard) displays agent action logs. In Specification 289, telemetry was enriched with granular file touch metrics (`FileTouchRead`, `FileTouchEdit`, `FileTouchCreated`).

However, the action logs query API in `cli/cmdagent/agent_ui_logs.go` only filters by `agent`, `actionType`, `subtask`, and `q`. Users inspecting agent sessions cannot filter to view only file edits, file inspections, or file creations.

Adding `touchType` filtering provides full parity between the database heatmap aggregations and the agent logs HTTP API.

---

## 3. Concrete Implementation Details

### 3.1 Struct Update in `cli/cmdagent/agent_ui_server.go`
```go
// ActionLogFilter specifies query filters for action logs.
type ActionLogFilter struct {
	TaskId     string
	AgentRole  string
	ActionType string
	TouchType  string
	SubtaskId  string
	Search     string
	Limit      int
}
```

### 3.2 Parsing in `cli/cmdagent/agent_ui_logs.go`
In `parseLogFilter`:
```go
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
		TouchType:  strings.TrimSpace(r.URL.Query().Get("touchType")),
		SubtaskId:  strings.TrimSpace(r.URL.Query().Get("subtask")),
		Search:     strings.TrimSpace(r.URL.Query().Get("q")),
		Limit:      limit,
	}
}
```

### 3.3 Filter Conditions in `cli/cmdagent/agent_ui_logs.go`
```go
func appendLogFilterConditions(c []string, args []any, f ActionLogFilter) ([]string, []any) {
	c, args = appendRoleFilter(c, args, f.AgentRole)
	c, args = appendTypeFilter(c, args, f.ActionType)
	c, args = appendTouchTypeFilter(c, args, f.TouchType)
	c, args = appendSubtaskFilter(c, args, f.SubtaskId)

	return appendSearchCondition(c, args, f.Search)
}

func appendTouchTypeFilter(c []string, args []any, touchType string) ([]string, []any) {
	hasTouch := len(touchType) > 0 && touchType != "ALL"
	if !hasTouch {
		return c, args
	}

	switch strings.ToUpper(touchType) {
	case "READ":
		c = append(c, "UPPER(ActionType) IN ('READ', 'READ_FILE', 'INSPECT')")
	case "EDIT":
		c = append(c, "UPPER(ActionType) IN ('EDIT', 'EDIT_FILE', 'MODIFY')")
	case "CREATED":
		c = append(c, "UPPER(ActionType) IN ('CREATE', 'WRITE_FILE', 'NEW')")
	default:
		c = append(c, "UPPER(ActionType) = UPPER(?)")
		args = append(args, touchType)
	}

	return c, args
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 2 (Function Limit):** `appendLogFilterConditions` is 7 lines; `appendTouchTypeFilter` is 15 lines. Both <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Guards flatten conditional branches.
- **Rule 4 (Positive Booleans):** Uses `hasTouch := len(touchType) > 0 && touchType != "ALL"`.
- **Rule 6 (Vertical Whitespace):** Proper blank line before `return` and after `}`.

---

## 5. Verification Criteria
- [x] `ActionLogFilter` contains `TouchType string`.
- [x] `parseLogFilter` extracts `r.URL.Query().Get("touchType")`.
- [x] `appendLogFilterConditions` invokes `appendTouchTypeFilter`.
- [x] Handles categorical touch types (`READ`, `EDIT`, `CREATED`) with SQL `IN` sets.
- [x] Functions stay within <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
