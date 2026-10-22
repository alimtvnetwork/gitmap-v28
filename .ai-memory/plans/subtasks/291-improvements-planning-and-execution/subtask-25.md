# Subtask Plan 25: Log Pagination Query Offset (`cli/cmdagent/agent_ui_logs.go`)

> **Subtask ID:** `subtask-25`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdagent/agent_ui_logs.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `parseQueryOffset(val string) int` in `cli/cmdagent/agent_ui_logs.go`.
2. Add `Offset int` field to `ActionLogFilter` struct in `cli/cmdagent/agent_ui_server.go`.
3. Wire the `offset` URL query parameter in `parseLogFilter` to enable multi-page navigation in the agent visualizer.
4. Update `buildActionLogsSQL` to append `OFFSET ?` when `Offset > 0`.
5. Keep all functions strictly <= 15 lines.

---

## 2. Targeted Problem & Architectural Rationale
The agent visualizer logs endpoint (`/api/logs`) currently supports filtering by role, action type, touch type, subtask, search text, and `limit`. However, because it lacks an `offset` parameter, clients cannot paginate backwards to view older action logs.

Adding `Offset` with default fallback to 0 provides:
1. Standard SQL `LIMIT ? OFFSET ?` pagination.
2. Complete backwards compatibility with existing UI clients that omit the `offset` parameter.
3. Safe query validation preventing negative offset inputs.

---

## 3. Concrete Implementation Details

### 3.1 Field Addition in `cli/cmdagent/agent_ui_server.go`
Extend `ActionLogFilter`:

```go
type ActionLogFilter struct {
	TaskId     string
	AgentRole  string
	ActionType string
	TouchType  string
	SubtaskId  string
	Search     string
	Limit      int
	Offset     int
}
```

### 3.2 Helper & Query Additions in `cli/cmdagent/agent_ui_logs.go`

```go
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

func parseLogFilter(r *http.Request, defaultTaskId string) ActionLogFilter {
	return ActionLogFilter{
		TaskId:     resolveTaskId(r.URL.Query().Get("taskId"), defaultTaskId),
		AgentRole:  strings.TrimSpace(r.URL.Query().Get("agent")),
		ActionType: strings.TrimSpace(r.URL.Query().Get("actionType")),
		TouchType:  strings.TrimSpace(r.URL.Query().Get("touchType")),
		SubtaskId:  strings.TrimSpace(r.URL.Query().Get("subtask")),
		Search:     strings.TrimSpace(r.URL.Query().Get("q")),
		Limit:      parseQueryInt(r.URL.Query().Get("limit"), 100),
		Offset:     parseQueryOffset(r.URL.Query().Get("offset")),
	}
}

func appendPaginationLimitOffset(base string, args []any, limit, offset int) (string, []any) {
	base += " ORDER BY ActionLogId DESC LIMIT ?"
	args = append(args, limit)
	if offset > 0 {
		base += " OFFSET ?"
		args = append(args, offset)
	}

	return base, args
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Invalid integer inputs safely resolve to 0 rather than terminating or throwing.
- **Rule 2 (Function Limit):** `parseQueryOffset` (12 lines), `appendPaginationLimitOffset` (9 lines). All <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Guard clauses on empty string and negative values.
- **Rule 4 (Positive Booleans):** Positive conditions and clean flow.
- **Rule 6 (Vertical Whitespace):** Proper blank lines before returns and after `}`.

---

## 5. Verification Criteria
- [x] `parseQueryOffset(val string) int` parses non-negative integers and returns 0 on invalid values.
- [x] `ActionLogFilter` struct includes `Offset int`.
- [x] `OFFSET ?` SQL clause is appended only when `Offset > 0`.
- [x] Functions are <= 15 lines with zero nested conditionals.
- [x] All file paths are strictly relative lowercase Git paths.
