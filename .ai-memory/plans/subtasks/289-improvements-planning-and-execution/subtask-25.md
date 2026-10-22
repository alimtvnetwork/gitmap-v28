# Subtask 25: Touch Type Telemetry in Agent Action Logs API

> **Subtask ID:** Subtask-25  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/cmdagent/agent_ui_server.go`  
> **Status:** PENDING  

---

## 1. Objectives & Scope
1. Expose granular touch type telemetry (`touchType`) in the Agent UI Server REST API in `cli/cmdagent/agent_ui_server.go`.
2. Map raw action log action types to standard `FileTouchType` string representations (`READ`, `EDIT`, `CREATED`, or `OTHER`).
3. Enrich the action log response payload in `formatActionLogMap` without breaking existing UI dashboard clients.
4. Keep all helper functions <= 15 lines and enforce zero nested conditionals.

---

## 2. Detailed Technical Plan

### 2.1 Action Type to Touch Type Classifier
- Implement `deriveTouchType` helper in `cli/cmdagent/agent_ui_server.go`:
  ```go
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
  ```

### 2.2 Enrichment of Action Log Map
- Update `formatActionLogMap` to inject `touchType`:
  ```go
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
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Rule 2 (Function Limit):** `deriveTouchType` is 13 lines, `formatActionLogMap` is 14 lines. Both strictly under the 15-line limit.
- **Rule 3 (Zero Nested Ifs):** Flat switch-case matching; zero nested if statements.
- **Rule 6 (Vertical Whitespace):** Single blank line before return statements.

---

## 4. Verification Criteria
- [ ] Action logs endpoint `/api/action-logs` includes `touchType` in JSON objects.
- [ ] Actions containing "READ" yield `"touchType": "READ"`.
- [ ] Actions containing "EDIT" yield `"touchType": "EDIT"`.
- [ ] Actions containing "CREATE" or "WRITE" yield `"touchType": "CREATED"`.
- [ ] Existing client consumers continue functioning without breaking changes.
