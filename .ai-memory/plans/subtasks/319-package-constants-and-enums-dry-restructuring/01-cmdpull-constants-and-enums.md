# Subtask 01: `cli/cmdpull` Constants & Enums Refactoring

> **Parent Plan:** [.ai-memory/plans/pending/319-package-constants-and-enums-dry-restructuring.md](../../pending/319-package-constants-and-enums-dry-restructuring.md)  
> **Target Package:** `cli/cmdpull`  
> **Status:** Ready for Execution (Next Turn)  

---

## 1. Scope & Technical Deliverables
1. **Enum: `PullModeType`** (`cli/cmdpull/enums/pull_mode_enum.go`):
   - Values: `PullModeInvalid` (0), `PullModeRebase`, `PullModeFfOnly`, `PullModeMerge`, `PullModeAuto`.
   - Methods: `IsValid() bool`, `String() string`, `IsRebase() bool`, `IsFastForward() bool`, `IsMerge() bool`, `IsAuto() bool`.
2. **Enum: `PullStatusType`** (`cli/cmdpull/enums/pull_status_enum.go`):
   - Values: `PullStatusInvalid` (0), `PullStatusUpToDate`, `PullStatusFastForward`, `PullStatusConflict`, `PullStatusDiverged`, `PullStatusError`.
   - Methods: `IsValid() bool`, `String() string`, `IsClean() bool`, `HasConflict() bool`.
3. **Constants: `cli/cmdpull/constants/constants.go`**:
   - CLI flags: `FlagRebase`, `FlagFfOnly`, `FlagAutoFix`, `FlagAll`, `FlagVerbose`.
   - Output tokens: `HeaderRepoArrow`, `SummaryStatHeader`, `StatusDirtyMarker`.
   - Error codes: `ErrPullConflict`, `ErrNetworkTimeout`.
4. **Call-Site Refactoring:**
   - Update `cli/cmdpull/pull_dispatch.go`, `cli/cmdpull/pull_cwd.go`, `cli/cmdpull/pull_summary_single.go`, `cli/cmdpull/pull_summary_changes.go` to consume `enums.*` and `constants.*`.

---

## 2. Coding Guidelines Invariants
- Files $\le 300$ lines.
- Functions $\le 15$ lines.
- Positive booleans only (`is*`, `has*`).
- Zero nested `if` statements.
- Blank line after `}` and before `return`.

---

## 3. Verification Evidence Contract
- Run `gofmt -w cli/cmdpull/`.
- AST pass via FastGate.
