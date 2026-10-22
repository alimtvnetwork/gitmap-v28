# Subtask 09: Store Table Name Enum Expansion for Heatmap & Pipeline Logs

> **Subtask ID:** Subtask-09  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-09.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Expand strongly-typed database table names in `cli/constants/constants_store.go` to cover `HeatmapLocalTracking` and `PipelineErrorLog`.
2. Add typed enum constants `TableNameHeatmapLocalTracking` and `TableNamePipelineErrorLog` to `TableNameType` (and its alias `DatabaseTableType`).
3. Add corresponding untyped string constants `TableHeatmapLocalTracking = "HeatmapLocalTracking"` and `TablePipelineErrorLog = "PipelineErrorLog"` for SQL query templating and backward compatibility.
4. Enforce strict coding guidelines:
   - PascalCase table names: Singular/PascalCase table naming standards (`HeatmapLocalTracking`, `PipelineErrorLog`).
   - Strong typing: All enum types must end with `Type`.
   - File length ratchet: Maintain `cli/constants/constants_store.go` well under the 500-line ceiling (current length ~354 lines).

---

## 2. Detailed Technical Plan

### 2.1 Expand `TableNameType` in `cli/constants/constants_store.go`
- Locate the `TableNameType` constant block:
  ```go
  const (
      TableNameRepo                 TableNameType = "Repo"
      TableNameGroup                TableNameType = "Group"
      TableNameGroupRepo            TableNameType = "GroupRepo"
      TableNameRelease              TableNameType = "Release"
      TableNameAlias                TableNameType = "Alias"
      TableNameRepoVersionHistory   TableNameType = "RepoVersionHistory"
      TableNameDetectedProject      TableNameType = "DetectedProject"
      TableNamePipelineRun          TableNameType = "PipelineRun"
      TableNameHeatmapActionLog     TableNameType = "HeatmapActionLog"
      TableNameHeatmapLocalTracking TableNameType = "HeatmapLocalTracking"
      TableNamePipelineErrorLog     TableNameType = "PipelineErrorLog"
  )
  ```

### 2.2 Expand Untyped Table Name Constants
- Locate the `Table names (v15: PascalCase + singular + {Table}Id PK)` constant block:
  ```go
  const (
      TableRepo                 = "Repo"
      TableGroup                = "Group"
      TableGroupRepo            = "GroupRepo"
      TableRelease              = "Release"
      TableAlias                = "Alias"
      TableRepoVersionHistory   = "RepoVersionHistory"
      TableDetectedProject      = "DetectedProject"
      TablePipelineRun          = "PipelineRun"
      TableHeatmapActionLog     = "HeatmapActionLog"
      TableHeatmapLocalTracking = "HeatmapLocalTracking"
      TablePipelineErrorLog     = "PipelineErrorLog"
  )
  ```

### 2.3 File Sizing & Hygiene Verification
- Verify that `cli/constants/constants_store.go` remains formatted with clean indentation and comment annotations.
- Verify file line count remains comfortably below 370 lines, adhering to the 500-line ratchet.
- Ensure no circular dependencies are introduced.

---

## 3. Verification Criteria

- [ ] `TableNameHeatmapLocalTracking` and `TableNamePipelineErrorLog` are added to `TableNameType`.
- [ ] `TableHeatmapLocalTracking` and `TablePipelineErrorLog` string constants are defined.
- [ ] Table names follow strict PascalCase formatting rules.
- [ ] File size of `cli/constants/constants_store.go` remains safely under 500 lines.
- [ ] Strictly relative Git paths used throughout all references.
