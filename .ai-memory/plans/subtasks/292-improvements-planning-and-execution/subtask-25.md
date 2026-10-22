# Subtask Plan 25: Sortable Agent Action Log Query Engine (`cli/cmdagent/agent_ui_logs.go`, `cli/cmdagent/agent_ui_server.go`)

> **Subtask ID:** `subtask-25`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdagent/agent_ui_server.go`, `cli/cmdagent/agent_ui_logs.go`  
> **Status:** Pending  

---

## 1. Objectives & Scope
1. Extend `ActionLogFilter` in `cli/cmdagent/agent_ui_server.go` with `SortBy string` and `Order string`.
2. Implement strict whitelist sanitization functions `parseSortColumn` and `parseSortOrder` in `cli/cmdagent/agent_ui_logs.go` to prevent SQL injection.
3. Wire query parameters `sort` and `order` in `parseLogFilter`.
4. Update `appendPaginationLimitOffset` to dynamically apply the whitelisted sort column and order.
5. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
The HTTP API `/api/agent/logs` serves action logs to the Agent Visualizer UI. Currently, the SQL generator hardcodes `" ORDER BY ActionLogId DESC LIMIT ?"`, preventing front-end clients from ordering records by timestamp (`CreatedAt`), action type (`ActionType`), agent role (`AgentRole`), or subtask ID (`SubtaskId`).

By introducing whitelisted sorting:
1. UI clients can sort chronologically or reverse-chronologically across multiple fields.
2. Strict whitelist switch-cases prevent any possibility of SQL injection through query parameters.
3. Fallbacks ensure unprovided or unrecognized sort columns safely default to `"ActionLogId"` and `"DESC"`.

---

## 3. Concrete Implementation Details

### 3.1 Struct Update in `cli/cmdagent/agent_ui_server.go`
Update `ActionLogFilter` to include `SortBy` and `Order`:

```go
type ActionLogFilter struct {
	TaskId     string
	AgentRole  string
	ActionType string
	TouchType  string
	SubtaskId  string
	Search     string
	SortBy     string
	Order      string
	Limit      int
	Offset     int
}
```

### 3.2 Whitelist Parsing in `cli/cmdagent/agent_ui_logs.go`
Add helper functions to parse and sanitize sort parameters:

```go
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
```

### 3.3 Filter Formulation in `cli/cmdagent/agent_ui_logs.go`
In `parseLogFilter`, populate `SortBy` and `Order`:

```go
SortBy: parseSortColumn(r.URL.Query().Get("sort")),
Order:  parseSortOrder(r.URL.Query().Get("order")),
```

Update `appendPaginationLimitOffset`:

```go
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
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** No exit handlers. Parsing and formatting return sanitized primitives and slices directly.
- **Rule 2 (Function Limit):** `parseSortColumn` is 14 lines; `parseSortOrder` is 7 lines; `appendPaginationLimitOffset` is 12 lines. All <= 15 lines.
- **Rule 3 (Zero Nested Ifs):** Linear early returns and switch statement guards.
- **Rule 4 (Positive Booleans):** Case-insensitive string equality with positive condition.
- **Rule 6 (Vertical Whitespace):** Blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [ ] `ActionLogFilter` struct updated with `SortBy` and `Order` fields.
- [ ] `parseSortColumn` strictly maps input to valid columns (`CreatedAt`, `ActionType`, `AgentRole`, `SubtaskId`, `ActionLogId`).
- [ ] `parseSortOrder` strictly maps input to `"ASC"` or `"DESC"`.
- [ ] `appendPaginationLimitOffset` formats dynamic `ORDER BY <column> <order> LIMIT ?`.
- [ ] Maintains backward compatibility when sort parameters are omitted.
