# Subtask 08: Targeted Interim Tests and Verification

> **Parent Plan:** [230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md](../../pending/230-token-purge-installer-workdir-pull-agm-and-ui-modernization.md)  
> **Spec Reference:** [02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md](../../../../02-spec/21-app/230-token-purge-installer-workdir-pull-agm-and-ui-modernization/02-component-and-cli-spec.md)  
> **Status:** `QUEUED`  
> **Target Subsystems:**  
> - `cli/cmdpull/`  
> - `cli/cmdpipeline/`  
> - `cli/scanner/`  

---

## 1. Technical Objective

Execute targeted interim verification and regression testing across subsystems modified in recent commits (`cli/cmdpull/`, `cli/cmdpipeline/`, and `cli/scanner/`). Validate that newly added features—such as `executeAllPipelineErrorLogs` for `gitmap pe all`, `--is-all` CLI flags, Ubuntu `pull-all` failure trees, and `DefaultScanExcludeDirs`—operate reliably without breaking existing behaviors. Enforce **Rule R1**: full repository test runs (`go test ./...`) and global CI/CD runners remain strictly banned. Only file-scoped and package-scoped tests are permitted.

---

## 2. Test Execution Inventory

| Subsystem | Targeted Test Files | Verifications Performed | Permitted Command Scope |
|:---|:---|:---|:---|
| **Pull & Sync** | `cli/cmdpull/pull_fallback_test.go`<br>`cli/cmdpull/pull_flags_test.go`<br>`cli/cmdpull/pull_progress_bar_test.go` | Fast-forward retry logic, auto-stash fallback, progress indicators, non-blocking network error handling | `go test -v -run TestPull ./cli/cmdpull/` |
| **Pipeline & Telemetry** | `cli/cmdpipeline/pipeline_errorlogs_test.go`<br>`cli/cmdpipeline/pipeline_all_errors_test.go`<br>`cli/cmdpipeline/pipeline_failure_tree_test.go` | Multi-commit error aggregation, `executeAllPipelineErrorLogs`, `--is-all` flag propagation, failure tree hierarchy | `go test -v -run TestPipeline ./cli/cmdpipeline/` |
| **Project Scanner** | `cli/scanner/scanner_test.go`<br>`cli/scanner/scanner_worktree_test.go`<br>`cli/scanner/scanner_depth_test.go` | `DefaultScanExcludeDirs` enforcement, nested git worktree detection, case-insensitive path uniqueness | `go test -v -run TestScanner ./cli/scanner/` |

---

## 3. Step-by-Step Verification Plan

### Step 3.1: Verify `cli/cmdpull/` Package
1. Run file-scoped unit tests:
   ```bash
   go test -v ./cli/cmdpull/pull_fallback_test.go ./cli/cmdpull/pull_fallback.go ./cli/cmdpull/pull_types.go
   ```
2. Validate that `executePullAllWithFallback` retries fast-forward cleanly and does not abort early on individual remote failures.
3. Validate that progress bar calculations do not divide by zero on empty queues.

### Step 3.2: Verify `cli/cmdpipeline/` Package
1. Run package-scoped unit tests:
   ```bash
   go test -v -run TestPipelineAllErrors ./cli/cmdpipeline/
   go test -v -run TestPipelineErrorLogs ./cli/cmdpipeline/
   ```
2. Confirm that `executeAllPipelineErrorLogs` aggregates errors across all tracked repositories without deadlock or SQLite concurrency lockouts.
3. Verify that `--is-all` flag produces valid JSON envelopes matching the positive boolean schema.

### Step 3.3: Verify `cli/scanner/` Package
1. Run scanner unit tests:
   ```bash
   go test -v -run TestScannerWorktree ./cli/scanner/
   go test -v -run TestScannerExcludeDirs ./cli/scanner/
   ```
2. Confirm that directories specified in `DefaultScanExcludeDirs` (`.cache`, `.oh-my-zsh`, `node_modules`, `.venv`) are strictly ignored during directory traversal.

### Step 3.4: Telemetry Schema & Positive Boolean Check
1. Execute CLI JSON telemetry commands in dry-run mode:
   ```bash
   ./gitmap pe --json | jq .
   ./gitmap rp ls --json | jq .
   ```
2. Validate that all JSON response fields adhere to positive polarity:
   - `isSuccess`, `hasErrors`, `isRunning`, `isDefault` (zero negative booleans like `isNotFailed` or `disabled`).

---

## 4. Verification Commands & Expected Output

```bash
# 1. Targeted execution for cmdpull
go test -v -run "TestPullFallback|TestPullFlags" ./cli/cmdpull/

# 2. Targeted execution for cmdpipeline
go test -v -run "TestPipelineErrorLogs|TestPipelineAllErrors" ./cli/cmdpipeline/

# 3. Targeted execution for scanner
go test -v -run "TestScannerWorktree|TestScannerExclude" ./cli/scanner/

# 4. Confirm zero full test suite execution
git status --porcelain
```

---

## 5. Acceptance Criteria

- [ ] All targeted unit tests in `cli/cmdpull/` pass (exit 0).
- [ ] All targeted unit tests in `cli/cmdpipeline/` pass (exit 0).
- [ ] All targeted unit tests in `cli/scanner/` pass (exit 0).
- [ ] Telemetry JSON outputs conform to positive boolean attributes.
- [ ] Zero invocations of `go test ./...` or global CI test runners.
