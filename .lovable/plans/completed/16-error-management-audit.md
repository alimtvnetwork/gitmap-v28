# Master Plan: Repository-Wide Error Management Audit & Compliance (16-error-management-audit.md)

> **Status:** COMPLETED  
> **Budget:** N = 200 (PHASE_1_STEPS = 100, PHASE_2_STEPS = 100)  
> **Target:** Replace all bare panics (`panic("error")`), bare exits (`os.Exit`), and swallowed errors with structured `apperror.AppError` and centralized `cliexit.HandleError`.

---

## 1. Task-Specific Rule Set

1. **Zero Bare Panics**: Never call `panic("error")` or `panic(err)` in command logic or business operations. Wrap in `apperror.NewWithDetails` or `apperror.WrapWithDetails` with full context and pass to `cliexit.HandleError`.
2. **Zero Bare Exits**: Direct `os.Exit(...)` calls in command handlers are strictly forbidden. All command terminations must route through `cliexit.HandleError(appErr, exitCode)`.
3. **Preserve Exit Codes & Attribution**: For domain-specific CLI exit codes (e.g. `constants.ExitVisNotARepo`, `ExitVisAuthFailed`), pass the explicit code to `cliexit.HandleError(appErr, code)`.
4. **Automated CI/CD Linter Enforcement**: `linter-scripts/check-error-management.py` must run as part of `.lovable/ai-fix-scripts/03-cicd-local-runner.py` and pass with 0 errors.

---

## 2. Violation Inventory

### A. Bare Panics (`panic("error")`)

- `gitmap/cmd/vscodeworkspace.go` (Lines 48, 162, 178)
- `gitmap/cmd/vscode_cmd.go` (Lines 19, 39, 46, 58)
- `gitmap/cmd/zipgroup.go` (Line 57)
- `gitmap/cmd/zipgroupcreate.go` (Lines 18, 45, 51, 79, 87)
- `gitmap/cmd/zipgroupops.go` (Lines 15, 23, 32, 43, 50, 56, 68, 71, 94, 100)
- `gitmap/cmd/zipgroupshow.go` (Lines 17, 23, 54, 66, 72)

### B. Bare Exits (`os.Exit`)

- `gitmap/cmd/visibilitymakelast.go` (Lines 42, 52, 60, 68)
- `gitmap/cmd/visibilityredo.go` (Line 23)
- `gitmap/cmd/visibilityresolve.go` (Lines 35, 41, 47, 60, 66, 72, 143, 165, 170)
- `gitmap/cmd/visibilityundo.go` (Lines 45, 85, 107, 118, 143)
- `gitmap/cmd/visibilityundoflags.go` (Lines 45, 51)
- `gitmap/cmd/workflow_open_pr.go` (Line 49)
- `gitmap/cmd/workflow_recent_todo.go` (Line 22)
- `gitmap/cmd/zip.go` (Lines 30, 34, 38, 196)
- `gitmap/helptext/print.go` (Lines 34, 55)
- `gitmap-updater/cmd/*.go` (`check.go`, `run.go`, `worker.go`)

---

## 3. Subtask Breakdown

- **Subtask 01**: [`01-vscode-commands-error-refactor.md`](../subtasks/16-error-management/01-vscode-commands-error-refactor.md)
- **Subtask 02**: [`02-zip-commands-error-refactor.md`](../subtasks/16-error-management/02-zip-commands-error-refactor.md)
- **Subtask 03**: [`03-visibility-and-workflow-error-refactor.md`](../subtasks/16-error-management/03-visibility-and-workflow-error-refactor.md)
- **Subtask 04**: [`04-helptext-and-updater-error-refactor.md`](../subtasks/16-error-management/04-helptext-and-updater-error-refactor.md)
- **Subtask 05**: [`05-ci-linter-integration-and-verification.md`](../subtasks/16-error-management/05-ci-linter-integration-and-verification.md)
