# Subtask 294.5: Decompose Agent Subtask Subsystem & Implement Structured JSON Envelopes

- **Parent Plan:** [.ai-memory/plans/pending/294-code-refactoring-and-coding-guideline-update.md](../../pending/294-code-refactoring-and-coding-guideline-update.md)
- **Status:** Completed
- **Spec Reference:** [02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md](../../../02-spec/21-app/294-code-refactoring-and-coding-guideline-update/02-component-spec.md)
- **File Boundaries:**
  - `cli/cmdagent/agent_subtask_types.go`
  - `cli/cmdagent/agent_subtask_cli.go`
  - `cli/cmdagent/agent_subtask.go`
  - `cli/cmdagent/agent_subtask_resolve.go`

---

## 1. Objectives & Scope

1. **Modular Subsystem Decomposition ($\le 300$ Lines Cap):**
   - Currently `agent_subtask.go` (297 lines) and `agent_subtask_resolve.go` (309 lines) either border or violate the 300-line repository limit.
   - Decompose into four cohesive modules:
     - `cli/cmdagent/agent_subtask_types.go`: All command options, summary models, and JSON response envelopes.
     - `cli/cmdagent/agent_subtask_cli.go`: Cobra commands, flag bindings, and terminal table output renderers.
     - `cli/cmdagent/agent_subtask.go`: Core execution functions (`RunSubtaskAdd`, `RunSubtaskClaim`, `RunSubtaskStart`, `RunSubtaskComplete`, `RunSubtaskFail`, `RunSubtaskLs`).
     - `cli/cmdagent/agent_subtask_resolve.go`: Subtask identifier and parent task DB resolution, slug generation, and JSON payload normalization.
2. **Standardized Response Envelope for `gitmap agent subtask ls --json`:**
   - Replace unstructured flat array JSON output with a comprehensive envelope containing:
     - Top-level status (`status: "success"`).
     - Relative database path (`databasePath`).
     - Canonical parent task identifier (`parentTaskId`).
     - Aggregated status metrics rollup (`summary`: total, pending, in-progress, completed, failed).
     - Classified subtasks array with resolved file targets.
3. **Grouped Commonalities for `gitmap agent subtask add`:**
   - Group common metadata (`databasePath`, `parentTaskId`, `addedCount`) at the top level of the JSON response instead of repeating per subtask item.
   - Support optional `--json` flag in `gitmap agent subtask add`.
4. **Style & Guideline Invariants:**
   - Functions strictly $\le 15$ lines.
   - All booleans prefixed with `is` or `has` (`isJson`, `hasDb`, `hasFiles`).
   - Clean spacing with blank lines before `if`, after `}`, and before `return`.

---

## 2. Step-by-Step Execution Plan

### Step 1: Create `cli/cmdagent/agent_subtask_types.go`
- Extract all option types from `agent_subtask.go`:
  - `SubtaskAddOptions` (including `IsJson bool`)
  - `SubtaskClaimOptions`
  - `SubtaskStartOptions`
  - `SubtaskCompleteOptions`
  - `SubtaskFailOptions`
  - `SubtaskLsOptions`
- Define structured envelope types:
  - `SubtaskStatusSummary` (fields: `TotalCount`, `PendingCount`, `InProgressCount`, `CompletedCount`, `FailedCount`)
  - `SubtaskLsEnvelope` (fields: `Status`, `DatabasePath`, `ParentTaskId`, `Summary`, `Subtasks`)
  - `SubtaskAddEnvelope` (fields: `Status`, `DatabasePath`, `ParentTaskId`, `AddedCount`, `Subtasks`)

### Step 2: Create `cli/cmdagent/agent_subtask_cli.go`
- Extract Cobra command definitions from `agent_subtask.go`:
  - `subtaskAddCmd`, `subtaskClaimCmd`, `subtaskStartCmd`, `subtaskCompleteCmd`, `subtaskFailCmd`, `subtaskLsCmd`.
- Extract flag registration helpers:
  - `initSubtaskCommands()`, `registerTargetFlags()`.
  - Add `--json` flag to `subtaskAddCmd`.
- Extract and refactor output formatters:
  - `renderSubtasksTable([]types.Subtask)` (moved from `agent_subtask_resolve.go`).
  - `renderSubtasksJsonEnvelope(subtasks []types.Subtask, dbPath, parentId string) *appfault.AppError`.
  - `computeSubtaskSummary(subtasks []types.Subtask) SubtaskStatusSummary`.
  - `renderSubtaskAddEnvelope(subtasks []types.Subtask, dbPath, parentId string) *appfault.AppError`.

### Step 3: Refactor `cli/cmdagent/agent_subtask.go`
- Streamline to contain only core execution routines:
  - `RunSubtaskAdd`: support both standard table/text output and `renderSubtaskAddEnvelope` when `opts.IsJson` is true.
  - `RunSubtaskClaim`: invoke `renderClaimResult`.
  - `RunSubtaskStart`, `RunSubtaskComplete`, `RunSubtaskFail`.
  - `RunSubtaskLs`: pass `dbPath` and `canonicalId` into `renderSubtasksJsonEnvelope` on `opts.IsJson`.
- Confirm final line count of `agent_subtask.go` is $\le 200$ lines.

### Step 4: Refactor `cli/cmdagent/agent_subtask_resolve.go`
- Remove table rendering function (`renderSubtasksTable`) and JSON dumping function (`renderSubtasksJson`) which now live in `agent_subtask_cli.go`.
- Retain only resolution, slug, and parsing functions.
- Confirm final line count of `agent_subtask_resolve.go` is $\le 220$ lines.

---

## 3. Verification & Acceptance Criteria

- [x] All 4 subtask files in `cli/cmdagent/` are strictly $\le 300$ lines.
- [x] `agent_subtask_types.go` centralizes all options and envelope definitions.
- [x] `agent_subtask_cli.go` handles Cobra wiring, flag binding, and UI rendering.
- [x] `gitmap agent subtask ls --json` outputs `SubtaskLsEnvelope` with summary rollup and relative database path.
- [x] `gitmap agent subtask add` supports grouped commonalities in JSON mode.
- [x] Zero compiler build or test commands executed during planning or authoring.
- [x] Relative Git paths strictly used across all specifications and plans.
