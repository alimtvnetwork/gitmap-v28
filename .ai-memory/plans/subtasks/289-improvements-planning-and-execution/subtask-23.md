# Subtask 23: File Touch Counter Granularity

> **Subtask ID:** Subtask-23  
> **Parent Plan:** `.ai-memory/plans/pending/289-improvements-planning-and-execution.md`  
> **Component Spec:** `02-spec/21-app/289-improvements-planning-and-execution/03-migration-and-verification-spec.md`  
> **Target Path:** `cli/store/heatmap_schema.go`  
> **Status:** Completed  

---

## 1. Objectives & Scope
1. Extend `FileTouchCount` struct in `cli/store/heatmap_schema.go` with granular operational counters: `ReadCount`, `EditCount`, and `CreatedCount`.
2. Preserve existing `Count` field representing aggregate touches across all operations.
3. Ensure backwards compatibility with all JSON serialization consumers in web visualizers and CLI summaries.
4. Maintain default zero-value semantics for records where operation-specific counts are not yet populated.

---

## 2. Detailed Technical Plan

### 2.1 Struct Definition Update in `cli/store/heatmap_schema.go`
- Update `FileTouchCount`:
  ```go
  // FileTouchCount represents file touch frequency with granular operation metrics.
  type FileTouchCount struct {
  	FilePath     string `json:"filePath"`
  	Count        int    `json:"count"`
  	ReadCount    int    `json:"readCount"`
  	EditCount    int    `json:"editCount"`
  	CreatedCount int    `json:"createdCount"`
  }
  ```

### 2.2 Aggregation Query Integration
- Ensure subsequent heatmap aggregation queries that group by `FilePath` can populate `ReadCount`, `EditCount`, and `CreatedCount` using conditional aggregation:
  ```sql
  SELECT FilePath,
         COUNT(*) as Count,
         SUM(CASE WHEN ActionType IN ('READ', 'INSPECT') THEN 1 ELSE 0 END) as ReadCount,
         SUM(CASE WHEN ActionType IN ('EDIT', 'MODIFY') THEN 1 ELSE 0 END) as EditCount,
         SUM(CASE WHEN ActionType IN ('CREATE', 'NEW') THEN 1 ELSE 0 END) as CreatedCount
  FROM HeatmapActionLog
  GROUP BY FilePath
  ORDER BY Count DESC
  ```

---

## 3. Strict Coding Guidelines Compliance
- **Backwards Compatibility:** Additive fields only; no existing fields removed or renamed.
- **Naming Conventions:** All field names in PascalCase with matching camelCase JSON tags.
- **Zero Nested Ifs:** Straightforward struct field additions without branching.

---

## 4. Verification Criteria
- [x] `FileTouchCount` compiles cleanly with new fields.
- [x] JSON serialization produces `readCount`, `editCount`, and `createdCount` properties.
- [x] Unmarshaling legacy JSON payloads without the new fields populates default 0 values without error.
- [x] Relative paths verified across all spec references.
