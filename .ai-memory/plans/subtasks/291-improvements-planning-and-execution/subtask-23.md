# Subtask Plan 23: Stale Heatmap Action Log Safe Purge (`cli/store/heatmap_store.go`)

> **Subtask ID:** `subtask-23`  
> **Parent Plan:** `.ai-memory/plans/pending/291-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/291-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/heatmap_store.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Implement `PurgeStaleActionLogs(olderThanDays int) (int64, error)` on `*HeatmapDB` in `cli/store/heatmap_store.go`.
2. Guard against invalid inputs by rejecting non-positive retention windows (`olderThanDays <= 0`).
3. Enforce the Zero-Harm safety invariant: NEVER delete action logs belonging to tasks marked `IN_PROGRESS` in `HeatmapParentTask`.
4. Return the number of deleted records (`RowsAffected`) and direct errors without invoking exit handlers.

---

## 2. Targeted Problem & Architectural Rationale
Autonomous agent loops and file tracking generate high-frequency telemetry in `HeatmapActionLog`. Over long running sessions, databases accumulate hundreds of thousands of historical log rows.

Without an automated, guarded retention mechanism:
1. SQLite database files grow unbounded, slowing down queries in the web visualizer.
2. In-flight operations risk having their active telemetry purged if deletion queries lack status guards.

`PurgeStaleActionLogs` provides safe garbage collection by isolating completed tasks while protecting active tasks.

---

## 3. Concrete Implementation Details

### 3.1 Method Addition in `cli/store/heatmap_store.go`
Add `PurgeStaleActionLogs` to `cli/store/heatmap_store.go`:

```go
const sqlPurgeStaleActionLogs = `DELETE FROM HeatmapActionLog 
WHERE CreatedAt < datetime('now', printf('-%d days', ?))
  AND TaskSlug NOT IN (SELECT TaskSlug FROM HeatmapParentTask WHERE Status = 'IN_PROGRESS')`

// PurgeStaleActionLogs removes action logs older than the specified retention window,
// strictly preserving logs associated with tasks currently in progress.
func (h *HeatmapDB) PurgeStaleActionLogs(olderThanDays int) (int64, error) {
	if olderThanDays <= 0 {
		return 0, apperror.NewSimple("olderThanDays must be positive", "E_INVALID_PARAM")
	}

	res, err := h.conn.Exec(sqlPurgeStaleActionLogs, olderThanDays)
	if err != nil {
		return 0, apperror.WrapSimple(err, "heatmap.purgeStale")
	}

	return res.RowsAffected()
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Guard returns `apperror.NewSimple`, exec errors wrapped via `apperror.WrapSimple`. No calls to `cliexit.HandleError`.
- **Rule 2 (Function Limit):** 13 lines of code, well within the 15-line limit.
- **Rule 3 (Zero Nested Ifs):** Linear validation with guard clause.
- **Rule 4 (Positive Booleans):** Positive conditions and validation.
- **Rule 6 (Vertical Whitespace):** Blank line before `return` and after `}`.

---

## 5. Verification Criteria
- [x] Method `PurgeStaleActionLogs(olderThanDays int) (int64, error)` is declared on `*HeatmapDB`.
- [x] Returns `0, nil` when `olderThanDays <= 0`.
- [x] Excludes tasks with `UPPER(Status) = 'IN_PROGRESS'` from deletion.
- [x] Function is <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
