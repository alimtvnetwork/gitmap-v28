# Subtask 01: Universal Polyglot File Runner & Task Audit DB

> **Subtask ID:** Subtask-01  
> **Parent Plan:** `.ai-memory/plans/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills.md`  
> **Target Subsystems:** `cli/cmdrun/`, `cli/cmd/macro_root_dispatch.go`, `cli/cmd/roottooling.go`  
> **Owned Files:**  
> - `cli/cmdrun/run_types.go`  
> - `cli/cmdrun/run_resolve.go`  
> - `cli/cmdrun/run_exec.go`  
> - `cli/cmdrun/run_audit_db.go`  
> - `cli/cmdrun/run_errors_db.go`  
> - `cli/cmdrun/run_errors_cmd.go`  
> - `cli/cmdrun/run_cmd.go`  
> - `cli/cmdrun/run_cmd_test.go`  
> - `cli/cmd/macro_root_dispatch.go`  
> - `cli/cmd/roottooling.go`  

---

## 1. Concrete Objectives

1. **Polyglot File Target & Extension Resolution:**
   - Implement `ResolveRunTarget(target string)` in `cli/cmdrun/run_resolve.go`.
   - If `target` exists directly on disk, inspect its extension.
   - If `target` does not exist directly on disk, probe candidate extensions in deterministic priority: `.py`, `.ps1`, `.sh`, `.js`, `.ts`, `.go`.
   - Reject ambiguous or non-executable targets with structured errors (`apperror.NewValidationError`).
2. **Interpreter Resolution & Execution Pipeline:**
   - Implement `ExecuteTarget(target RunTarget, args []string)` in `cli/cmdrun/run_exec.go`.
   - Map file extensions to system runtimes:
     - `.py`: Auto-resolve Python executable from `installation.db` cache or system PATH (`python3`, `python`, `py -3`).
     - `.ps1`: Execute with `pwsh` or `powershell -NoProfile -ExecutionPolicy Bypass -File`.
     - `.sh`: Execute with `bash` or `sh`.
     - `.js`: Execute with `node`, `bun`, or `deno`.
     - `.ts`: Execute with `bun`, `tsx`, `ts-node`, or `node`.
     - `.go`: Execute via `go run`.
   - Capture execution stdout, stderr, execution duration, and exit code.
3. **SQLite Task DB Audit Lifecycle:**
   - Implement task lifecycle tracking in `cli/cmdrun/run_audit_db.go` utilizing `store.OpenTasksRootSplitDB()`.
   - Before launching the child process, insert record into `TaskQueue` with `Action: "file-run"`, `Target: target.ResolvedPath`, `Status: "running"`.
   - Upon process termination with exit code `0`, commit audit record to `TaskHistory` with `Status: "completed"`.
4. **Run Errors Database (`run_errors`) & Error Capture:**
   - In `cli/cmdrun/run_errors_db.go`, ensure the `run_errors` table exists in SQLite split DB.
   - When execution produces exit code `!= 0` or fails to launch, insert structured `RunErrorRecord` into `run_errors` table containing:
     - `ErrorId`: Unique identifier (`runerr-<timestamp>`).
     - `FilePath`: Relative workspace path to the executed file.
     - `FileExtension`: File extension (`.py`, `.ps1`, etc.).
     - `Interpreter`: Invoked interpreter executable.
     - `ExitCode`: Process exit code.
     - `DurationMs`: Total elapsed runtime in milliseconds.
     - `ErrorMessage`: Summary of failure or exit code error.
     - `StdoutSnippet`: Bounded tail of stdout (last 2048 bytes).
     - `StderrSnippet`: Bounded tail of stderr (last 2048 bytes).
     - `ExecutedArgs`: Joined command arguments.
5. **Run Inspection CLI Commands:**
   - Implement `gitmap run errors` (aliases: `gitmap run-errors`, `gitmap run err`) displaying recent execution failures in an ANSI table.
   - Implement `gitmap run history` showing task audit records from `TaskHistory`.
   - Implement `gitmap run clear-errors` to purge `run_errors`.
6. **Command Dispatcher Disambiguation:**
   - Modify `runMacroRootRun` in `cli/cmd/macro_root_dispatch.go`:
     - If argument is `errors`, `history`, or `clear-errors`, dispatch to `cmdrun`.
     - If argument resolves to an existing file on disk or matches an auto-extension candidate, dispatch to `cmdrun.RunFile()`.
     - Fall back to macro execution (`cmdmacro.ExecuteDynamicMacroWithArgs`) only if no file target matches.
   - Register `run-errors` and `run-history` in `cli/cmd/roottooling.go`.

---

## 2. Core Domain Types & Structs

```go
package cmdrun

import (
	"time"
)

// RunTarget represents a resolved executable script file.
type RunTarget struct {
	RawInput      string
	ResolvedPath  string
	Extension     string
	Interpreter   string
	IsDirectMatch bool
}

// RunOptions configures process execution parameters.
type RunOptions struct {
	Args       []string
	Timeout    time.Duration
	WorkingDir string
}

// RunResult captures process termination output and telemetry.
type RunResult struct {
	Target     RunTarget
	ExitCode   int
	DurationMs int64
	Stdout     string
	Stderr     string
	Err        error
}

// RunErrorRecord represents a persisted execution failure in SQLite.
type RunErrorRecord struct {
	ErrorID       string    `json:"error_id"`
	FilePath      string    `json:"file_path"`
	FileExtension string    `json:"file_extension"`
	Interpreter   string    `json:"interpreter"`
	ExitCode      int       `json:"exit_code"`
	DurationMs    int64     `json:"duration_ms"`
	ErrorMessage  string    `json:"error_message"`
	StdoutSnippet string    `json:"stdout_snippet"`
	StderrSnippet string    `json:"stderr_snippet"`
	ExecutedArgs  string    `json:"executed_args"`
	CreatedAt     time.Time `json:"created_at"`
}
```

---

## 3. Implementation Checklist

- [ ] **1. Data Models & Constants (`cli/cmdrun/run_types.go`):**
  - Define `RunTarget`, `RunOptions`, `RunResult`, and `RunErrorRecord`.
  - Define supported extensions slice: `[]string{".py", ".ps1", ".sh", ".js", ".ts", ".go"}`.
- [ ] **2. File Target Resolution (`cli/cmdrun/run_resolve.go`):**
  - Implement `ResolveRunTarget(input string) (*RunTarget, error)`.
  - Check `os.Stat(input)`. If file exists, return directly.
  - If not found, iterate supported extensions, checking `input + ext`.
  - Convert absolute paths to workspace-relative paths.
  - Return `RunTarget` with identified interpreter.
- [ ] **3. Process Execution Engine (`cli/cmdrun/run_exec.go`):**
  - Implement `ExecuteTarget(target *RunTarget, opts RunOptions) (*RunResult, error)`.
  - Assemble command arguments per interpreter:
    - `.py` -> `[]string{pythonBin, target.ResolvedPath, ...args}`
    - `.ps1` -> `[]string{pwshBin, "-NoProfile", "-File", target.ResolvedPath, ...args}`
    - `.sh` -> `[]string{bashBin, target.ResolvedPath, ...args}`
    - `.js` -> `[]string{nodeBin, target.ResolvedPath, ...args}`
    - `.ts` -> `[]string{tsBin, target.ResolvedPath, ...args}`
    - `.go` -> `[]string{"go", "run", target.ResolvedPath, ...args}`
  - Stream process output while buffering tail bytes for error telemetry.
  - Calculate elapsed execution time in milliseconds.
- [ ] **4. SQLite Task Audit & Errors DB (`cli/cmdrun/run_audit_db.go`, `cli/cmdrun/run_errors_db.go`):**
  - Implement `EnsureRunErrorsTable(conn *sql.DB) error`.
  - Implement `EnqueueRunTaskAudit(target *RunTarget, args []string) (string, error)`.
  - Implement `CompleteRunTaskAudit(taskId string, durationMs int64) error`.
  - Implement `RecordRunError(record RunErrorRecord) error`.
  - Implement `QueryRecentRunErrors(limit int) ([]RunErrorRecord, error)`.
  - Implement `ClearRunErrors() error`.
- [ ] **5. CLI Handlers & Table Rendering (`cli/cmdrun/run_cmd.go`, `cli/cmdrun/run_errors_cmd.go`):**
  - Implement `Run(args []string) error` routing between file run and subcommands.
  - Implement `RunErrorsCmd(args []string) error` rendering ANSI table with columns: `ID`, `FILE`, `INTERPRETER`, `EXIT`, `DURATION`, `ERROR`.
  - Implement `RunHistoryCmd(args []string) error` querying task audit.
  - Implement `RunClearErrorsCmd(args []string) error`.
- [ ] **6. Dispatcher Wiring (`cli/cmd/macro_root_dispatch.go`, `cli/cmd/roottooling.go`):**
  - Update `runMacroRootRun`: Check `args[0]` against file resolver and subcommands before macro dispatch.
  - Register `run-errors` and `run-history` in `toolingAuditEntries()` in `cli/cmd/roottooling.go`.
- [ ] **7. Unit & Integration Testing (`cli/cmdrun/run_cmd_test.go`):**
  - Test resolution of existing and extensionless files.
  - Test failure recording into SQLite DB.
  - Test ANSI table rendering for errors.

---

## 4. Acceptance Criteria

1. Running `gitmap run test_script` automatically discovers `test_script.py` (or `.ps1`/`.sh`/`.js`/`.ts`/`.go`) and executes it.
2. Every invocation enqueues an audit entry in the SQLite tasks database.
3. If the script exits with non-zero status, a failure record is persisted to `run_errors`.
4. Running `gitmap run errors` lists the failure record with exit code, duration, and error snippet.
5. All file paths emitted or stored are strictly relative Git paths.
6. Existing macro functionality remains fully backwards-compatible.

---

## 5. Verification Commands

```bash
# Run unit test suite
go test -v ./cli/cmdrun/...

# Verify build compilation
go build -v ./cli/...

# Manual verification tests
gitmap run errors
gitmap run history
```
