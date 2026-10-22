# Subtask 09: Split-DB Table Name Types & Telemetry Entity Constants

> **Subtask ID:** Subtask-09  
> **Parent Plan:** `.ai-memory/plans/291-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/291-improvements-planning-and-execution/subtask-09.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Modernize SQLite table identifiers in `cli/constants/constants_store.go` by introducing strongly-typed `DatabaseTableType` definitions.
2. Add typed and untyped constants for telemetry and pipeline Split-DB tables:
   - `TablePipelineRun = "PipelineRun"` / `TablePipelineRun DatabaseTableType = "PipelineRun"`
   - `TableHeatmapActionLog = "HeatmapActionLog"` / `TableHeatmapActionLog DatabaseTableType = "HeatmapActionLog"`
3. Enforce the PascalCase and singular entity naming conventions across all Split-DB tiers (master, installation, repodb, pipeline, heatmap).
4. Maintain clean type definition standards ending in `Type` suffix.

---

## 2. Detailed Technical Plan

### 2.1 Define `DatabaseTableType` in `cli/constants/constants_store.go`
- Add `DatabaseTableType` (or alias with existing `TableNameType`):
  ```go
  // DatabaseTableType represents strongly-typed database table names across all Split-DB tiers.
  type DatabaseTableType string

  const (
      TableRepo               DatabaseTableType = "Repo"
      TableGroup              DatabaseTableType = "Group"
      TableGroupRepo          DatabaseTableType = "GroupRepo"
      TableRelease            DatabaseTableType = "Release"
      TableAlias              DatabaseTableType = "Alias"
      TableRepoVersionHistory DatabaseTableType = "RepoVersionHistory"
      TableDetectedProject    DatabaseTableType = "DetectedProject"
      TablePipelineRun        DatabaseTableType = "PipelineRun"
      TableHeatmapActionLog   DatabaseTableType = "HeatmapActionLog"
  )
  ```

### 2.2 Add String Constants for Migration & Query Compatibility
- Ensure untyped / string constants exist for query builders and DDL strings:
  ```go
  const (
      TablePipelineRunStr     = "PipelineRun"
      TableHeatmapActionLogStr = "HeatmapActionLog"
  )
  ```

### 2.3 Verification of PascalCase + Positive Boolean Conventions
- Verify that table definitions enforce:
  - Table names: PascalCase and singular (`PipelineRun`, `HeatmapActionLog`).
  - Column names: camelCase or PascalCase per schema spec, with PKs as `{Table}Id`.
  - Boolean columns: Strictly positive booleans (`isActive`, `hasPassed`, `isSuccess`).

---

## 3. Verification Criteria

- [ ] `DatabaseTableType` is defined ending in `Type` suffix in `cli/constants/constants_store.go`.
- [ ] `TablePipelineRun` and `TableHeatmapActionLog` are defined as strongly-typed constants.
- [ ] Table names strictly follow PascalCase and singular naming rules.
- [ ] `cli/constants/constants_store.go` complies with spacing rules and the 500-line limit.
