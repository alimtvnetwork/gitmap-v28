---
name: cg-isolate-os-tests
description: Autonomously audits, refactors, and validates repository-wide unit tests to ensure that tests NEVER trigger real OS shutdown, reboot, power-off, system modifications, or heavy unmocked system calls, enforcing injectable executors and mock duration verification across polyglot stacks.
---

# Skill: Unit Test Isolation & Destructive OS Prevention (`cg-isolate-os-tests`)

> **Invocation:** `@[cg-isolate-os-tests](file;.cursor/skills/cg-isolate-os-tests)`

This skill governs autonomous auditing, refactoring, and validation of unit tests to strictly prevent execution of destructive host OS commands (shutdown, reboot, poweroff, system package modifications, user group alterations, unmocked filesystem purges) by enforcing injectable command executors and mock duration math.

## 1. Concurrency Invariants & Header Parameters

```text
N = 400 (Total self-loop steps budget — editable top-header parameter, default: 400)
A = 3   (MANDATORY number of spawned autonomous subagents running concurrently via invoke_subagent, default: 3)
H = 2   (Operational hands per agent: dual-task batch capacity & parallel tool dispatch, default: 2)
C = 30  (Tool calls per worker before it must report, default: 30)

System Concurrency Capacity = A × H = 3 agents × 2 hands = 6 concurrent subtask operations
Tree Hierarchy = Lead Orchestrator (Root Node) -> Branch Coordinators (Subtask Owners) -> Leaf Workers & Dynamic Joining Assistants
PHASE_1_BUDGET = 100     (Steps 1 .. 100: Planning, Discovery Subagents, Up-Front Brain & Roadmap Creation, Lean Subtasks)
PHASE_2_BUDGET = 150     (Steps 101 .. 250: Intermediate Archival, Sanitization, Prerequisite Refactoring, Intermediate Checkpoints)
PHASE_3_BUDGET = 150     (Steps 251 .. 400: Tree Swarm Parallel Execution, Dynamic Join Capacity Saturation, Targeted Quality Linting, Fleet Sync, Atomic Push)
WAVES = ceil(subtasks / (A x H))
```

---

## 2. Mandatory Architectural Invariants

### 1. Total Ban on Real OS Shutdown / Reboot in Unit Tests
The following command invocations MUST NEVER execute directly against the host OS during test runs:
- `shutdown` (`/s`, `/r`, `/a`, `-h`, `-r`, `-c`)
- `reboot`
- `poweroff` / `halt`
- `init 0` / `init 6`
- `systemctl poweroff` / `systemctl reboot`
- Direct package managers: `apt-get upgrade`, `Get-WindowsUpdate`, `softwareupdate`

### 2. Mandatory Injectable Executors
Every production subsystem that interacts with host power controls, system packages, or filesystem caches MUST expose an injectable executor or runner function:

```go
type OSActionExecutor func(params OSActionParams) *appfault.AppError
var DefaultOSActionExecutor OSActionExecutor = executeNativeOSAction
```

```go
type FileRemover func(filePath string) (int, int64)
var defaultFileRemover FileRemover = removeSingleFileSafely
```

```go
type TempDirResolver func() []string
var defaultTempDirResolver TempDirResolver = resolveTempDirectories
```

```go
type OSCommandRunner func(cmd *exec.Cmd) *appfault.AppError
var defaultOSCommandRunner OSCommandRunner = func(cmd *exec.Cmd) *appfault.AppError {
    err := cmd.Run()
    if err != nil {
        return appfault.Wrap(err, "defaultOSCommandRunner", appfault.ErrCodeCommandExecutionFailure)
    }
    return nil
}
```

### 3. Hermetic Mock Testing with `defer` Restoration
Unit tests must ALWAYS substitute the default executor inside a `defer` restoration block so that test execution remains 100% isolated:

```go
mock, cleanup := setupMockExecutor()
defer cleanup()

// Execute CLI or unit logic
err := RunSchedulePowerCLI(OSActionShutdown, []string{"1s"})
// Assert mock captured parameters without invoking host commands
```

### 4. Fast Duration Math (1s/2s Checks)
Tests for countdown timers, delayed power commands, or recurring schedules must test:
- Duration string parsing (`ParseScheduleDuration`)
- Arithmetic target calculation (`TargetTime = ScheduledAt + Duration`)
- Fast simulated delays (`1s`, `2s`) with mock tick callbacks rather than sleeping or arming host OS timers.

### 5. Function Sizing & Structured Errors
- Functions target **<= 8 lines** of body logic (hard cap: 15 lines).
- Affirmative booleans only (`is*`, `has*`). No negative booleans.
- Universal `*appfault.AppError` return wrapping (standardized per coding guidelines).

---

## 3. Step 0 Preflight with Data Schema Verification

Before enqueuing subtasks or claiming work, initialize and verify the SQLite task database schema:

```bash
gitmap agent task init --name "<task name>" --budget 400
gitmap agent task status --ddl
gitmap agent task status --json
```

- **DDL & Schema Verification Protocol:** Inspect tables (`ParentTask`, `Subtask`, `AgentActionLog`), affirmative boolean flags (`IsActive`, `HasCompleted`, `IsBlocked`), and foreign key linkages.
- **Forensic Autopsy Check:** Check `diagnostics.hasCrashesDetected`. If true, inspect `diagnostics.crashedAgents` to identify failed subtasks, targeting files, and last logged actions to facilitate clean resume.

---

## 4. 3-Tier Database Location Transparency & Diagnostic Suite

Whenever executing or logging subtask telemetry, fulfill the Database Location Transparency Contract:

| Database Tier | Physical File Location | Access Mode & Concurrency Protocol | Data Retained & Ownership |
| :--- | :--- | :--- | :--- |
| **Tier 1: Machine-Local Root Database** | `~/.gitmap/agent-root.db` | Machine sessions, cross-repo leases | Machine Sessions & Cross-Repo Leases |
| **Tier 2: Repo-Cache Companion Storage** | `repo-cache/.heatmap/<ai-slug>/<repo-slug>/heatmap.db` | Isolated Split-DB with periodic sync (`gitmap heatmap rc sync`) | Cross-repository multi-agent actions, tool usage traces, commit checkpoints |
| **Tier 3: In-Flight Local Run Database** | `.ai-memory/temp-agents/<slug>/agent-task.db` | WAL Mode (`PRAGMA journal_mode=WAL; busy_timeout=5000`) | Task status, worker claims, subtasks, in-flight action logs (<1ms Go engine) |
| **Tier 3: Local Heatmap Cache** | `.ai-memory/heatmap.db` | WAL Mode + Read-Only Queries (`mode=ro`) | Cognitive action telemetry, file exploration reads/edits, contention leases |

### Executable Diagnostic Commands:

```bash
# Pre-edit contention and lease check
gitmap heatmap discuss [--window 15m] [--threshold 2] [--json]

# Cognitive file exploration and editing telemetry
gitmap heatmap log -t <task-slug> -a "READ" -f "<path>" -r "Reasoning: <explanation>" -s "EXPLORATION"
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <isolation-refactor-rationale>" -s "EXECUTION"

# Real-time continuous activity stream & tree
gitmap heatmap tree [-t <task-slug>]
gitmap heatmap live

# Read-only companion query and schema inspection
gitmap heatmap rc query "SELECT ActionType, FilePath, ReasoningText FROM AgentActionLog ORDER BY AgentActionLogId DESC LIMIT 5" --format table
gitmap heatmap rc tables
gitmap heatmap rc view AgentActionLog --limit 10
gitmap heatmap rc sync
```

---

## 5. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine turns. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run test suites (`go test ./...`, `pytest`) during routine execution.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`) after individual file edits.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 6. Verification Checklist

- [ ] All destructive host OS commands routed through injectable executors.
- [ ] Unit tests hermetically isolated with `defer` restoration.
- [ ] Errors wrapped using `*appfault.AppError`.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
