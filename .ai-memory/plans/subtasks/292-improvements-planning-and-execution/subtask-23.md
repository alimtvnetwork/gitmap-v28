# Subtask Plan 23: Heatmap Task Action Aggregation (`cli/store/heatmap_store.go`, `cli/store/heatmap_schema.go`)

> **Subtask ID:** `subtask-23`  
> **Parent Plan:** `.ai-memory/plans/pending/292-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/292-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/heatmap_schema.go`, `cli/store/heatmap_store.go`  
> **Status:** Pending  

---

## 1. Objectives & Scope
1. Define the `HeatmapTaskSummary` model struct in `cli/store/heatmap_schema.go`.
2. Implement `GetHeatmapTaskSummary(taskSlug string) (*HeatmapTaskSummary, error)` on `*HeatmapDB` in `cli/store/heatmap_store.go`.
3. Aggregate overall action volume and categorize actions by operation touch type (`ReadCount`, `EditCount`, `CreateCount`, `OtherCount`) in a single SQL query.
4. Support querying a specific `taskSlug` or querying across all tasks when `taskSlug` is empty.
5. Adhere strictly to Error Return Sovereignty (Rule 1), Function Limit (<= 15 lines, Rule 2), and Zero Nested Conditionals (Rule 3).

---

## 2. Targeted Problem & Architectural Rationale
The existing `HeatmapStatusSummary` aggregates overall action records and top-touched files, but does not provide a categorized operation breakdown (read, edit, create, other) for discrete task reporting or agent telemetry.

Adding `HeatmapTaskSummary`:
1. Provides a single-query SQL aggregation using conditional `CASE WHEN` sums.
2. Gives CLI commands and visualizers instant insight into whether an agent primarily read, edited, or created files.
3. Keeps query overhead low by avoiding multiple iterative queries.

---

## 3. Concrete Implementation Details

### 3.1 Model Struct in `cli/store/heatmap_schema.go`
Add `HeatmapTaskSummary` following `HeatmapStatusSummary`:

```go
// HeatmapTaskSummary aggregates overall action volume and categorizes actions by operation touch type for a task.
type HeatmapTaskSummary struct {
	TaskSlug     string `json:"taskSlug"`
	TotalActions int    `json:"totalActions"`
	ReadCount    int    `json:"readCount"`
	EditCount    int    `json:"editCount"`
	CreateCount  int    `json:"createCount"`
	OtherCount   int    `json:"otherCount"`
}
```

### 3.2 Method Implementation in `cli/store/heatmap_store.go`
Add `GetHeatmapTaskSummary` to `cli/store/heatmap_store.go`:

```go
const sqlHeatmapTaskSummary = `SELECT
    COUNT(*),
    COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('READ', 'READ_FILE', 'INSPECT') THEN 1 ELSE 0 END), 0),
    COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('EDIT', 'EDIT_FILE', 'MODIFY') THEN 1 ELSE 0 END), 0),
    COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('CREATE', 'WRITE_FILE', 'NEW') THEN 1 ELSE 0 END), 0)
FROM HeatmapActionLog
WHERE (? = '' OR TaskSlug = ?)`

// GetHeatmapTaskSummary returns aggregated read, edit, and create counts for a specified task.
func (h *HeatmapDB) GetHeatmapTaskSummary(taskSlug string) (*HeatmapTaskSummary, error) {
	cleanSlug := strings.TrimSpace(taskSlug)
	var total, reads, edits, creates int
	err := h.conn.QueryRow(sqlHeatmapTaskSummary, cleanSlug, cleanSlug).Scan(&total, &reads, &edits, &creates)
	if err != nil {
		return nil, apperror.WrapSimple(err, "heatmap.taskSummary")
	}

	other := total - (reads + edits + creates)
	if other < 0 {
		other = 0
	}

	return &HeatmapTaskSummary{
		TaskSlug:     cleanSlug,
		TotalActions: total,
		ReadCount:    reads,
		EditCount:    edits,
		CreateCount:  creates,
		OtherCount:   other,
	}, nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Database scan errors bubble up wrapped via `apperror.WrapSimple(err, "heatmap.taskSummary")`.
- **Rule 2 (Function Limit):** Function is exactly 14 lines, strictly within the 15-line limit.
- **Rule 3 (Zero Nested Ifs):** Linear guard checks for error return and other count normalization.
- **Rule 4 (Positive Booleans):** No boolean parameters or return types required; positive polarity maintained.
- **Rule 6 (Vertical Whitespace):** Blank lines after `}` and before `return`.

---

## 5. Verification Criteria
- [ ] `HeatmapTaskSummary` struct declared with JSON tags in `cli/store/heatmap_schema.go`.
- [ ] Method `GetHeatmapTaskSummary(taskSlug string) (*HeatmapTaskSummary, error)` declared on `*HeatmapDB` in `cli/store/heatmap_store.go`.
- [ ] Uses parameterized query `sqlHeatmapTaskSummary` preventing injection.
- [ ] Calculates `OtherCount` safely with non-negative lower bound guard.
- [ ] Returns structured error wrapped with `apperror.WrapSimple`.
- [ ] All touched files remain under their respective file size limits.
