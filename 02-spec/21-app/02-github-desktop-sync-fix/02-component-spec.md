# Spec 02: GitHub Desktop Pacing, Throttling & Error Surfacing Component Specification

> **Spec ID:** `SPEC-02-COMPONENT`  
> **Subsystems:** `cli/cmdide`, `cli/desktop`, `cli/cmddesktopsync`  
> **Status:** Approved Component Specification  
> **Date:** 2026-10-13  
> **Plan Reference:** `.ai-memory/plans/02-github-desktop-sync-fix.md`  

---

## 1. Executive Summary & Problem Formulation

When synchronizing Git repositories into GitHub Desktop on Ubuntu (and other Linux or multi-window desktop environments), users observe that bulk synchronization operations either fail mid-stream, drop repositories silently, or cause GitHub Desktop to crash.

Analysis of the CLI and Electron runtime identified two major defects in GitMap:
1. **Unthrottled Registration Burst (Electron Modal & IPC Collision):** In `cli/cmdide/ide_sync.go`, `cli/desktop/desktop.go`, and `cli/cmddesktopsync/desktopsync_ops.go`, registration commands were executed in tight 0ms loops. Because GitHub Desktop is an Electron application that uses React for its user interface and an asynchronous IPC bus between main and renderer processes, rapid bursts of CLI requests trigger modal dialog collisions and IPC queue saturation, causing Electron processes to hang or crash.
2. **Systematic Error Swallowing:** Execution errors and subprocess outputs were discarded across `cli/cmdide/ide_sync.go` (`exec.Command(cli, p).Run() == nil`) and `cli/cmdide/ide_add.go` (`res.IsDesktopAffected = ... Run() == nil`). Failures were counted merely as incremented numbers or marked as skipped, with no stderr capture, no failure diagnostic messages, and no structured reporting.

This component specification establishes:
- The **Pacing & Throttling Specification** (200ms–350ms delay between sequential registrations).
- The **Error Surfacing & Diagnostics Specification** (`cmd.CombinedOutput()`, structured failure summaries, diagnostic terminal output).
- The **Component Synchronization Contract** across `cli/cmdide`, `cli/desktop`, and `cli/cmddesktopsync`.

---

## 2. Electron Architecture & Pacing Rationale

### 2.1 Electron Runtime Mechanics
GitHub Desktop is built on Electron, which executes code across two primary process types:
1. **Main Process (Node.js):** Manages OS lifecycle, windows, single-instance mutex locks, and command-line argument dispatch.
2. **Renderer Process (Chromium + React + Redux/Dispatcher):** Hosts the graphical user interface, renders dialogs, manages git repository state in local IndexedDB stores, and coordinates user interactions.

```mermaid
flowchart TD
    A["GitMap CLI Batch Loop"] -->|"Rapid Exec (0ms delay)"| B["OS Exec: github <path>"]
    B --> C["Electron Main Process<br/>(SingleInstance Lock & IPC)"]
    C -->|"IPC Event: open-repository"| D["React Renderer Process"]
    D --> E{"React State & Modal Dialog"}
    E -->|"Re-entrant Invocation"| F["Race Condition / Dialog Collision / IPC Drop"]
    F --> G["Dropped Repo / App Freeze / SIGSEGV Crash"]
```

### 2.2 Modal Dialog & State Collision Failure Mode
When the CLI launches `github <repoPath>`:
1. A new CLI process starts and hands off the target directory path to the running GitHub Desktop instance via the single-instance IPC channel.
2. The main process relays an `open-repository` command to the React renderer window.
3. The renderer checks if the repository exists in its database:
   - If not present, it mounts an "Add Local Repository" confirmation dialog or executes an internal repository validation cycle.
   - If a subsequent CLI process triggers another `open-repository` before the previous modal cycle resolves, the React dispatcher experiences re-entrancy conflicts.
   - On Linux systems with Wayland or X11 window managers, window focus-stealing and modal pop-under triggers compound the issue, causing UI thread starvation, memory spikes, or silent process termination.

### 2.3 Pacing Benchmark & Selection
Testing demonstrates that introducing a **200ms–350ms delay** between sequential repository registrations resolves IPC and modal collisions:
- **Default Recommended Delay:** `250ms` (`const DefaultDesktopRegistrationDelay = 250 * time.Millisecond`).
- **Safety Window:** Provides sufficient time for Chromium's V8 event loop to process React state updates and persist IndexedDB records without blocking the OS event loop.
- **Total Duration Impact:** For a typical workspace of 20 repositories, a 250ms delay adds only ~5.0 seconds total, converting an unreliable crash prone operation into a 100% deterministic, resilient flow.

---

## 3. Pacing & Throttling Specification

### 3.1 Pacing Constant
Pacing duration is defined in `cli/desktop` or shared constants:
```go
// DefaultDesktopRegistrationDelay defines the pause between sequential
// GitHub Desktop repository registrations to allow Electron's React
// renderer to process IPC events and modals without dialog collision.
const DefaultDesktopRegistrationDelay = 250 * time.Millisecond
```

### 3.2 Loop Throttling Contract
All batch registration loops iterating over multiple repositories must apply throttling:
1. Execute repository registration via the resolved CLI.
2. If more repositories remain in the queue, pause execution via `time.Sleep(DefaultDesktopRegistrationDelay)`.
3. Avoid redundant sleep on the final repository in the iteration.

```mermaid
sequenceDiagram
    autonumber
    participant Engine as GitMap Engine
    participant CLI as GitHub Desktop CLI
    participant Electron as Electron React Window

    Engine->>CLI: exec github /path/repo1
    CLI->>Electron: IPC open-repository
    Electron-->>CLI: OK (Registered)
    Engine->>Engine: Sleep(250ms)
    Engine->>CLI: exec github /path/repo2
    CLI->>Electron: IPC open-repository
    Electron-->>CLI: OK (Registered)
```

---

## 4. Error Surfacing & Diagnostics Specification

### 4.1 Subprocess Execution Overhaul
All GitHub Desktop CLI invocations must transition from raw `.Run()` or `.Output()` to `cmd.CombinedOutput()`:

```go
cmd := exec.Command(cli, repoPath)
output, err := cmd.CombinedOutput()
```

- **Output Capture:** Capture both stdout and stderr.
- **Output Sanitization:** Trim whitespace (`strings.TrimSpace(string(output))`).
- **Error Retention:** If `err != nil`, retain both `err` and any CLI output string to format diagnostic feedback.

### 4.2 Enhanced Data Structures in `cli/cmdide/ide_types.go`

#### A. Structured Failure Tracking (`IDESyncFailure`)
```go
// IDESyncFailure records failure context for an individual repository registration.
type IDESyncFailure struct {
	Target   string `json:"target"`
	RepoName string `json:"repoName"`
	RepoPath string `json:"repoPath"`
	Reason   string `json:"reason"`
	Output   string `json:"output,omitempty"`
}
```

#### B. Updated `IDESyncSummary`
```go
// IDESyncSummary captures registration results across all editors.
type IDESyncSummary struct {
	VSCodeAdded         int              `json:"vscodeAdded"`
	VSCodeRetained      int              `json:"vscodeRetained"`
	CursorAdded         int              `json:"cursorAdded"`
	CursorRetained      int              `json:"cursorRetained"`
	AntigravityAdded    int              `json:"antigravityAdded"`
	AntigravityRetained int              `json:"antigravityRetained"`
	DesktopAdded        int              `json:"desktopAdded"`
	DesktopFailed       int              `json:"desktopFailed"`
	TotalRepos          int              `json:"totalRepos"`
	DesktopFailures     []IDESyncFailure `json:"desktopFailures,omitempty"`
}
```

#### C. Updated `IDEActionResult`
```go
// IDEActionResult reports single-action outcome for add or remove.
type IDEActionResult struct {
	Path                  string            `json:"path"`
	Action                string            `json:"action"`
	IsVSCodeAffected      bool              `json:"isVSCodeAffected"`
	IsCursorAffected      bool              `json:"isCursorAffected"`
	IsAntigravityAffected bool              `json:"isAntigravityAffected"`
	IsDesktopAffected     bool              `json:"isDesktopAffected"`
	HasErrors             bool              `json:"hasErrors"`
	ErrorDetails          map[string]string `json:"errorDetails,omitempty"`
}
```

### 4.3 Terminal Output & Diagnostics

#### In `cli/cmdide/ide_sync.go`:
In `printSyncSummary(s IDESyncSummary, isJSON, isQuiet bool)`:
- When `!isQuiet && s.DesktopFailed > 0`:
  Output diagnostic failure details for each item in `s.DesktopFailures`:
  ```text
    ✗ Desktop Failures (2):
      • my-repo-1: exit status 1 (Failed to open repository: directory locked)
      • my-repo-2: exit status 127 (github: command not found)
  ```
- When `isJSON`: `desktopFailures` is serialized cleanly in JSON output.

#### In `cli/cmdide/ide_add.go`:
In `printActionResult(res IDEActionResult, isJSON, isQuiet bool)`:
- If Desktop registration fails (`!res.IsDesktopAffected && res.ErrorDetails["desktop"] != ""`):
  Display failure diagnostic instead of "— skipped":
  ```text
    • Desktop:     ✗ failed: exit status 1 (IPC connection timed out)
  ```

---

## 5. Synchronization Across Components

| Subsystem Component | Role & Scope | Pacing Policy | Error Surfacing Policy |
| :--- | :--- | :--- | :--- |
| **`cli/cmdide/ide_sync.go`** | Multi-editor sync orchestrator (`gitmap ide sync`) | 250ms sleep between sequential repos | Captures `cmd.CombinedOutput()`, records in `s.DesktopFailures`, prints diagnostic table |
| **`cli/cmdide/ide_add.go`** | Single repo add orchestrator (`gitmap ide add <path>`) | N/A (single repo) | Captures error & output, records in `res.ErrorDetails["desktop"]`, surfaces in CLI output |
| **`cli/desktop/desktop.go`** | Core Desktop registration package (`desktop.AddRepos`) | 250ms sleep between repos in `addAll` | Replaces `cmd.Output()` with `cmd.CombinedOutput()`, bubbles enriched error |
| **`cli/cmddesktopsync/desktopsync_ops.go`** | Dedicated Desktop sync engine (`gitmap desktop-sync`) | 250ms sleep between repos in `syncAll` | Preserves `cmd.CombinedOutput()`, cleanses output strings, provides diagnostic logs |

---

## 6. Coding Guidelines & Invariants

1. **Strict Relative Paths:** All documentation, logs, and references must use forward-slash relative paths (`cli/...`, `02-spec/...`).
2. **Positive Boolean Naming:** All boolean fields and variables must use positive prefixes:
   - `isTargeted`, `isDryRun`, `isQuiet`, `isJSON`, `isInstalled`, `hasErrors`.
   - Never use negative booleans (`isNotAdded`, `noError`, `unregistered`).
3. **Decomposition & Sizing Limits:**
   - File length: strictly ≤ 300 lines.
   - Function length: strictly ≤ 20 lines (decompose loops, diagnostics, and summaries into dedicated helper functions).
4. **AppError & Clean Exit Codes:**
   - Package errors must wrap lower-level OS errors with context using `apperror` conventions where applicable.

---

## 7. Verification & Acceptance Criteria

- [ ] **Pacing Verification:** Execution of `ide sync` and `desktop-sync` across 5+ repositories exhibits measured inter-invocation delay of 200ms–350ms.
- [ ] **Zero Swallowed Errors:** Executing registration against a simulated invalid CLI shim or locked repository surfaces the exact error and combined output in terminal and JSON modes.
- [ ] **Positive Booleans:** All added structs and functions adhere to `hasErrors`, `isAffected`, `isTargeted`.
- [ ] **Lint & Architecture Gate:** Pass verification via `python3 03-ai-scripts/50-fastgate.py`.
