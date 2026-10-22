# Subtask 03: Repo Manage UI Error Return Sovereignty & Fallback Bubble-Up

> **Subtask ID:** Subtask-03  
> **Parent Plan:** `.ai-memory/plans/292-improvements-planning-and-execution.md`  
> **Target Path:** `.ai-memory/plans/subtasks/292-improvements-planning-and-execution/subtask-03.md`  
> **Status:** Completed  

---

## 1. Objectives & Scope

1. Modernize `RunRepoManageUI` in `cli/cmd/cmd_see.go` to eliminate swallowed execution errors and silent discards (`_ = cmdWeb.Start()`).
2. Ensure that terminal interactive UI failures and web dashboard fallback errors are properly captured and bubbled as `*apperror.AppError` to the root command dispatcher.
3. Clean up vertical whitespace formatting across `cli/cmd/cmd_see.go`, ensuring blank lines after closing braces `}` and before all `return` statements.
4. Enforce strict coding guidelines:
   - Function limit: $\le 15$ lines per function ($\le 8$ lines preferred).
   - Conditionals: Zero nested `if` statements with guard clause decomposition.
   - Positive booleans: Prefix boolean checks strictly with `is` or `has` (`hasErr`, `isFallbackFailed`).
   - File length: Keep `cli/cmd/cmd_see.go` comfortably under the 500-line ceiling.

---

## 2. Detailed Technical Plan

### 2.1 Refactor `RunRepoManageUI` in `cli/cmd/cmd_see.go`
- Locate `RunRepoManageUI()`:
  - Decompose into primary launcher and web fallback launcher to respect the $\le 15$-line limit.
  - Implement `launchInteractiveDashboard() error`:
    ```go
    func launchInteractiveDashboard() error {
        fmt.Println("Launching GitMap Repository Management Dashboard...")
        cmd := exec.Command("gitmap", "interactive")
        cmd.Stdin = os.Stdin
        cmd.Stdout = os.Stdout
        cmd.Stderr = os.Stderr

        return cmd.Run()
    }
    ```
  - Implement `launchWebDashboardFallback() *apperror.AppError`:
    ```go
    func launchWebDashboardFallback() *apperror.AppError {
        fmt.Println("Interactive dashboard unavailable; launching web UI fallback...")
        cmdWeb := exec.Command("gitmap", "ui")
        if err := cmdWeb.Start(); err != nil {
            return apperror.WrapSimple(err, "launch web UI fallback")
        }

        return nil
    }
    ```
  - Assemble `RunRepoManageUI() *apperror.AppError`:
    ```go
    func RunRepoManageUI() *apperror.AppError {
        if err := launchInteractiveDashboard(); err != nil {
            return launchWebDashboardFallback()
        }

        return nil
    }
    ```

### 2.2 Whitespace & Coding Standards Cleanup in `cli/cmd/cmd_see.go`
- Audit helper functions in `cli/cmd/cmd_see.go` (e.g., `resolveErrorsSSHArgs`):
  - Ensure mandatory blank line after closing brace `}`.
  - Ensure mandatory blank line before every `return` statement.
  - Enforce positive boolean naming (`isLonger`, `hasArgs`).

### 2.3 Dispatcher Audit Verification
- Verify `cli/cmd/rootcore.go` mapping:
  - `{[]string{"repo-manage", "repo-manage-ui"}, func() error { return RunRepoManageUI() }}`
  - Confirms that any returned `*apperror.AppError` propagates directly into `finishCommandAudit`, guaranteeing accurate audit telemetry in `gitmap.db`.

---

## 3. Verification Criteria

- [x] Zero swallowed errors in `RunRepoManageUI` and `_ = cmdWeb.Start()` is completely removed.
- [x] Fallback failure is captured and returned as a typed `*apperror.AppError`.
- [x] All functions in `cli/cmd/cmd_see.go` adhere to the $\le 15$ lines limit.
- [x] Mandatory blank lines present after `}` and before every `return`.
- [x] All booleans adhere to positive `is*` / `has*` naming conventions.
- [x] Strictly relative Git paths used throughout all references.
