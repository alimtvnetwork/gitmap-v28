# Component and CLI Specification: Automatic Binary Swap & Process Hand-off

**Sequence ID:** 336  
**Document Index:** 02  
**Target Package:** `cli/cmdswap` (or subpackage `cli/cmdbinary`)  
**Status:** Approved  
**Author:** Spec Author Subagent 2  

---

## 1. Overview & Architectural Goals

The Automatic Binary Swap & Process Hand-off component provides a deterministic, zero-downtime, and cross-platform mechanism for swapping active or locked GitMap binaries on Windows, Linux, and macOS.

Because running executables on Windows are locked by the operating system (`ERROR_SHARING_VIOLATION`), a direct in-place file replacement fails while any instance of the binary or an IDE plugin is actively running. On Unix systems, file unlinking or atomic directory replacement (`mv`) is possible, but running daemon processes (e.g. background runners, lemon daemons, UI webview hosts) still execute stale memory copies.

The `cmdswap` package orchestrates:
1. **Lock Identification & Process Resolution**: Identifies locking PIDs and existing GitMap processes targeting the active executable.
2. **Deterministic Script Plan Generation**: Generates an external detached script (`.ps1` for Windows PowerShell, `.sh` for POSIX systems) capable of:
   - Waiting for parent process termination or terminating specified PIDs.
   - Forcibly killing locking processes if requested (`--kill` / `--force`).
   - Retrying binary replacement (`Move-Item -Force` / `Copy-Item -Force` or atomic `mv`).
   - Reopening or running post-swap health checks (`gitmap --version`) and resuming background services.
3. **Execution or Script Output**: Runs the script detached in a new process group or outputs the script payload directly when `--script-only` is provided.
4. **Structured Telemetry & Output**: Returns structured JSON (`--json`) or human-readable Lipgloss terminal UI status messages.

---

## 2. CLI Command & Flag Definitions

### 2.1 Command Invocation
```bash
gitmap binary swap [flags]
# Primary alias:
gitmap swap [flags]
```

### 2.2 Flag Specifications

| Flag | Type | Default Value | Description |
| :--- | :--- | :--- | :--- |
| `--from` | `string` | Auto-detected (`./cli/gitmap.exe`, `./bin/gitmap.exe`, `./gitmap.exe`) | Source new binary path to copy from. |
| `--to` | `string` | Active Deployed Binary (`cmdinstall.ResolveDeployedAndConfigPaths()`) | Destination target binary path to be replaced. |
| `--pid` | `int` | `0` (ignored) | Specific process PID to terminate prior to file swap. |
| `--kill` | `bool` | `false` | Terminate all active processes locking or executing the target executable. |
| `--restart`, `--reopen` | `bool` | `false` | Execute the newly installed binary (`gitmap --version`) or launch UI after swapping. |
| `--force`, `-f` | `bool` | `false` | Forcefully kill locking processes (`Stop-Process -Force` / `kill -9`). |
| `--script-only` | `bool` | `false` | Emit the generated detached handoff script without executing it immediately. |
| `--json` | `bool` | `false` | Emit output in structured JSON format for machine integration and IDEs. |

### 2.3 Flag Parsing & Resolution Rules
- If `--from` is omitted, the engine scans the current directory and parent build trees in order:
  1. `./cli/gitmap.exe` (or `./cli/gitmap` on Unix)
  2. `./bin/gitmap.exe` (or `./bin/gitmap` on Unix)
  3. `./gitmap.exe` (or `./gitmap` on Unix)
  If none exist, the command returns an `ErrSourceBinaryNotFound` validation error.
- If `--to` is omitted, the engine invokes `cmdinstall.ResolveDeployedAndConfigPaths()`. If the deployed path cannot be resolved, it falls back to `os.Executable()`.
- If `--script-only` is `true`, no background process is spawned, no processes are killed, and the script contents are written to standard output or a temporary script file path.

---

## 3. Data Structures & Contract Model

All data structures enforce positive boolean conventions (`is`/`has` prefixes), zero pointer overhead on values, and single reusable definitions.

```go
package cmdswap

import (
	"time"
)

// SwapOptions encapsulates all input parameters for the binary swap workflow.
type SwapOptions struct {
	FromPath       string `json:"fromPath"`
	ToPath         string `json:"toPath"`
	TargetPID      int    `json:"targetPid"`
	HasKill        bool   `json:"hasKill"`
	HasRestart     bool   `json:"hasRestart"`
	HasForce       bool   `json:"hasForce"`
	IsScriptOnly   bool   `json:"isScriptOnly"`
	IsJSONOutput   bool   `json:"isJsonOutput"`
	MaxWaitSeconds int    `json:"maxWaitSeconds"`
}

// ScriptPlan models the platform-specific detached execution plan.
type ScriptPlan struct {
	Platform       string        `json:"platform"`
	ScriptPath     string        `json:"scriptPath"`
	ScriptContent  string        `json:"scriptContent"`
	TargetPIDList  []int         `json:"targetPidList"`
	IsDetached     bool          `json:"isDetached"`
	ExecutionArgs  []string      `json:"executionArgs"`
	EstimatedWait  time.Duration `json:"estimatedWait"`
}

// SwapResult encapsulates the execution outcome returned to callers or emitted as JSON.
type SwapResult struct {
	IsSuccess      bool        `json:"isSuccess"`
	FromPath       string      `json:"fromPath"`
	ToPath         string      `json:"toPath"`
	TerminatedPIDs []int       `json:"terminatedPids"`
	ScriptPlan     *ScriptPlan `json:"scriptPlan,omitempty"`
	ErrorMessage   string      `json:"errorMessage,omitempty"`
	OutputMessage  string      `json:"outputMessage"`
	DurationMs     int64       `json:"durationMs"`
}

// ProcessLockInfo describes a detected process holding a file handle or executing the target binary.
type ProcessLockInfo struct {
	PID         int    `json:"pid"`
	ProcessName string `json:"processName"`
	Executable  string `json:"executable"`
	IsParent    bool   `json:"isParent"`
}
```

---

## 4. Execution Workflow & Step Decomposition

### Phase 1: Options Sanitization & Binary Verification
1. Resolve `--from` and `--to` absolute paths.
2. Verify `--from` exists and is a valid executable.
3. Compute and log SHA-256 checksums of source and destination binaries.
4. Verify write permissions on the target directory and destination binary.

### Phase 2: Process Lock Inspection
1. Query active system processes matching `ToPath` using native OS APIs or process table inspections:
   - Windows: Query `Get-Process` / WMI `Win32_Process` or handle enumeration.
   - Unix: Inspect `/proc` or `lsof`/`fuser` for open file descriptors.
2. If `--pid` is specified, verify process existence and add to `TargetPIDList`.
3. If `--kill` or `--force` is enabled, aggregate all detected PIDs into `TargetPIDList`.
4. If locking processes exist and neither `--kill` nor `--force` is enabled:
   - In non-script mode: return a descriptive error prompting the user to supply `--kill` or `--force`.

### Phase 3: Detached Script Plan Construction
1. Select template based on `runtime.GOOS`:
   - Windows: PowerShell `.ps1` template with retry loop, process termination, file copy, and restart check.
   - Unix: POSIX Shell `.sh` script with loop, `kill -9` fallback, atomic `mv`, and restart check.
2. Inject bounded retry intervals (e.g. 10 retries with 250ms exponential backoff).
3. Embed cleanup commands to remove the temporary handoff script after completion.

### Phase 4: Execution Dispatch
- If `IsScriptOnly`:
  - Print script body to standard output or return `ScriptPlan` in JSON mode.
- Else:
  - Write script to OS temporary directory (`os.TempDir()`).
  - Launch external process detached (`cmd.exe /c start /b powershell.exe ...` on Windows, `nohup sh ... &` on Unix).
  - Parent process exits cleanly with exit code 0 or emits success JSON.

---

## 5. Script Templates & Hand-off Mechanics

### 5.1 Windows PowerShell Hand-off Pattern (`swap_handoff.ps1`)
```powershell
param(
    [string]$SourcePath,
    [string]$TargetPath,
    [int[]]$PidsToKill,
    [switch]$ForceKill,
    [switch]$RestartAfter
)

# 1. Terminate locking processes
foreach ($p in $PidsToKill) {
    if ($p -gt 0) {
        Stop-Process -Id $p -Force -ErrorAction SilentlyContinue
    }
}

# 2. Wait for unlock and copy binary
$maxRetries = 20
$copied = $false
for ($i = 0; $i -lt $maxRetries; $i++) {
    try {
        Copy-Item -Path $SourcePath -Destination $TargetPath -Force -ErrorAction Stop
        $copied = $true
        break
    } catch {
        Start-Sleep -Milliseconds 250
    }
}

if (-not $copied) {
    Write-Error "Failed to swap binary after $maxRetries attempts."
    exit 1
}

# 3. Optional restart verification
if ($RestartAfter) {
    Start-Process -FilePath $TargetPath -ArgumentList "--version" -NoNewWindow
}
```

### 5.2 POSIX Shell Hand-off Pattern (`swap_handoff.sh`)
```sh
#!/bin/sh
SOURCE="$1"
TARGET="$2"
RESTART="$3"
shift 3
PIDS="$@"

for p in $PIDS; do
    kill -9 "$p" 2>/dev/null || true
done

sleep 0.2
cp -f "$SOURCE" "$TARGET"
chmod +x "$TARGET"

if [ "$RESTART" = "true" ]; then
    "$TARGET" --version
fi
```

---

## 6. Coding Guidelines & Invariant Enforcements

All implementations of this component must adhere to the following rules:

1. **Error Management**:
   - Leaf functions return `*apperror.AppError` directly.
   - Outer CLI dispatch handler (in `cli/cmdswap/swap_cmd.go` or `cli/cmd/dispatch.go`) captures errors and reports to `finishCommandAudit`.
   - Never call `os.Exit` or `cliexit.HandleError` inside business logic or leaf helpers.
2. **Function Sizing & Complexity**:
   - Each function must be <= 15 lines (preferably <= 8 lines).
   - Zero nested `if` statements. Use guard clauses, early returns, and discrete helper functions.
3. **Positive Boolean Naming**:
   - All booleans must begin with `is` or `has` (e.g. `isSuccess`, `hasKill`, `hasRestart`, `isDetached`).
   - Strictly no negative booleans (e.g. `noKill`, `isNotRestarting` are prohibited).
4. **Spacing & Hygiene**:
   - Ensure a blank line after `}` and a blank line before `return`.
   - All string literals for script commands and error messages must use centralized constants.

---

## 7. Verification & Acceptance Criteria

1. **CLI Flag Compliance**:
   - Invoking `gitmap swap --help` lists all specified flags (`--from`, `--to`, `--pid`, `--kill`, `--restart`, `--reopen`, `--force`, `--script-only`, `--json`).
2. **Source Detection**:
   - Successfully auto-detects `./cli/gitmap.exe` or returns `ErrSourceBinaryNotFound` if absent.
3. **Script-Only Mode**:
   - `gitmap swap --script-only` prints deterministic PowerShell/Bash content without executing or killing processes.
4. **JSON Mode**:
   - `gitmap swap --json` produces valid JSON complying with `SwapResult`.
5. **Clean Exit & Error Return**:
   - All leaf functions bubble errors without swallowing or exiting prematurely.
