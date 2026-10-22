# Subtask Plan 23: Heatmap SQL SUM Aggregation for Touch Operations (`cli/store/heatmap_store.go`)

> **Subtask ID:** `subtask-23`  
> **Parent Plan:** `.ai-memory/plans/pending/290-improvements-planning-and-execution.md`  
> **Parent Spec:** `02-spec/21-app/290-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/heatmap_store.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Update `populateTopTouchedFiles` in `cli/store/heatmap_store.go` to compute operational counts (`ReadCount`, `EditCount`, `CreatedCount`) via SQL `SUM(CASE WHEN ... THEN 1 ELSE 0 END)`.
2. Populate all fields of `FileTouchCount` struct (added in Specification 289).
3. Eliminate silent error swallowing in the rows scanning loop.
4. Decompose the scanning logic into modular helpers to maintain functions <= 15 lines.

---

## 2. Targeted Problem & Architectural Rationale
In `cli/store/heatmap_store.go`, `populateTopTouchedFiles` currently executes:
```go
const sqlTop = `SELECT FilePath, COUNT(*) AS c FROM HeatmapActionLog
	WHERE (? = '' OR TaskSlug = ?) AND FilePath != ''
	GROUP BY FilePath ORDER BY c DESC LIMIT 10`
```
And scans only `FilePath` and `Count`:
```go
if err := rows.Scan(&item.FilePath, &item.Count); err == nil {
	s.TopTouchedFiles = append(s.TopTouchedFiles, item)
}
```
This causes:
1. `ReadCount`, `EditCount`, and `CreatedCount` remain unpopulated (0) in status summaries and agent dashboard visualizations.
2. Scan errors are silently ignored via `err == nil`, hiding database column mismatch bugs or schema corruption.

---

## 3. Concrete Implementation Details

### 3.1 Refactor `populateTopTouchedFiles` in `cli/store/heatmap_store.go`

```go
func (h *HeatmapDB) populateTopTouchedFiles(s *HeatmapStatusSummary, taskSlug string) error {
	const sqlTop = `SELECT FilePath, COUNT(*) AS c,
		COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('READ', 'READ_FILE', 'INSPECT') THEN 1 ELSE 0 END), 0) AS reads,
		COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('EDIT', 'EDIT_FILE', 'MODIFY') THEN 1 ELSE 0 END), 0) AS edits,
		COALESCE(SUM(CASE WHEN UPPER(TRIM(ActionType)) IN ('CREATE', 'WRITE_FILE', 'NEW') THEN 1 ELSE 0 END), 0) AS creates
		FROM HeatmapActionLog
		WHERE (? = '' OR TaskSlug = ?) AND FilePath != ''
		GROUP BY FilePath ORDER BY c DESC LIMIT 10`
	rows, err := h.conn.Query(sqlTop, taskSlug, taskSlug)
	if err != nil {
		return apperror.WrapSimple(err, "heatmap.statusTopFiles")
	}

	defer rows.Close()

	return scanTopTouchedRows(rows, s)
}

func scanTopTouchedRows(rows *sql.Rows, s *HeatmapStatusSummary) error {
	for rows.Next() {
		var item FileTouchCount
		err := rows.Scan(&item.FilePath, &item.Count, &item.ReadCount, &item.EditCount, &item.CreatedCount)
		if err != nil {
			return apperror.WrapSimple(err, "heatmap.scanTopFiles")
		}

		s.TopTouchedFiles = append(s.TopTouchedFiles, item)
	}

	return nil
}
```

---

## 4. Strict Coding Guidelines Compliance
- **Rule 1 (Error Return Sovereignty):** Scan errors return `apperror.WrapSimple(err, ...)` directly. No errors swallowed.
- **Rule 2 (Function Limit):** Both `populateTopTouchedFiles` (15 lines) and `scanTopTouchedRows` (14 lines) strictly satisfy the <= 15 line limit.
- **Rule 3 (Zero Nested Ifs):** Error checking uses early guard returns; no nested branching.
- **Rule 6 (Vertical Whitespace):** One blank line after `}` and before `return`.

---

## 5. Verification Criteria
- [x] `populateTopTouchedFiles` executes SQL SUM aggregation for `reads`, `edits`, `creates`.
- [x] Populates `ReadCount`, `EditCount`, and `CreatedCount` on each `FileTouchCount`.
- [x] Scan errors are propagated via `scanTopTouchedRows`.
- [x] Both functions <= 15 lines.
- [x] All file paths are strictly relative lowercase Git paths.
