# Subtask 03: `cli/cmdclone` Constants & Enums Refactoring

> **Parent Plan:** [.ai-memory/plans/pending/319-package-constants-and-enums-dry-restructuring.md](../../pending/319-package-constants-and-enums-dry-restructuring.md)  
> **Target Package:** `cli/cmdclone`  
> **Status:** Ready for Execution (Next Turn)  

---

## 1. Scope & Technical Deliverables
1. **Enum: `CloneStrategyType`** (`cli/cmdclone/enums/clone_strategy_enum.go`):
   - Values: `CloneStrategyInvalid` (0), `CloneStrategyBatch`, `CloneStrategySequential`, `CloneStrategyInteractive`, `CloneStrategyParallel`.
   - Methods: `IsValid() bool`, `String() string`, `IsParallel() bool`, `IsInteractive() bool`.
2. **Enum: `CloneProfileType`** (`cli/cmdclone/enums/clone_profile_enum.go`):
   - Values: `CloneProfileInvalid` (0), `CloneProfileDefault`, `CloneProfileWork`, `CloneProfilePersonal`, `CloneProfileFleet`.
   - Methods: `IsValid() bool`, `String() string`, `IsFleet() bool`.
3. **Constants: `cli/cmdclone/constants/constants.go`**:
   - Flags: `FlagStrategy`, `FlagProfile`, `FlagConcurrency`, `FlagDryRun`.
   - Prompts: `PromptSelectTargetDir`, `BannerCloneStarting`.
   - Telemetry tags: `TagCloneBatchSize`, `TagCloneTotalBytes`.
4. **Call-Site Refactoring:**
   - Update `cli/cmdclone/cmd_clone.go`, `cli/cmdclone/clone_batch.go`, `cli/cmdclone/clone_interactive.go`.

---

## 2. Coding Guidelines Invariants
- Files $\le 300$ lines.
- Functions $\le 15$ lines.
- Positive booleans only (`is*`, `has*`).
- Zero nested `if` statements.

---

## 3. Verification Evidence Contract
- Run `gofmt -w cli/cmdclone/`.
- AST pass via FastGate.
