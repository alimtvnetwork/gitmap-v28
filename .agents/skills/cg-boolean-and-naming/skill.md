---
name: cg-boolean-and-naming
description: Autonomously audits, refactors, and validates repository-wide boolean conventions, positive prefixes, implicit checks, enum Type suffixes, and nested if flattening against 02-spec/02-coding-guidelines/.
---

# Skill: Coding Guidelines — Booleans, Naming & Enums (`cg-boolean-and-naming`)

> **Invocation:** `@[cg-boolean-and-naming](file;.agents/skills/cg-boolean-and-naming)`

This skill governs autonomous execution for boolean conventions, semantic naming, enum standardization, and conditional flattening across all source files in alignment with `02-spec/02-coding-guidelines/`.

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

## 2. Mandatory Architectural Rules

1. **Implicit Boolean Checks Only:**
   - NEVER write `if isReady == true` or `if (isValid === true)`.
   - Positive booleans MUST ALWAYS be evaluated implicitly: `if isReady { ... }` or `if !isReady { ... }`.
   - Never compare against boolean literals (`== false`, `!= true`).

2. **Boolean Prefixes (`is`, `has`) & Affirmative Naming:**
   - `is`, `has` as prefix is only acceptable and nothing else acceptable including but not limited to `can`, `should`, `was`, `will`, `did`, `must`, etc.
   - No negative boolean identifiers (`isNotValid`, `isUndefined`, `isNotDefined`, `isNotSet`, `hasNoData` are banned).
   - **Try `IsDefined` instead of negatives:** When verifying presence, definition, or initialization, always use affirmative `isDefined` / `IsDefined` (or `isValid`, `hasValue`, `isReady`, `isFound`). Invert only once at the callsite guard clause (`if !isDefined { ... }`) if handling the missing case.
   - **Mandatory Replacement for `!isEmpty`:** NEVER use inverted negative empty checks (`!isEmpty`, `!res.IsEmpty()`). Always use affirmative `isDefined` / `res.IsDefined()` when asserting that data or records are present.
   - **Map Lookups (Canonical):** For map lookups, always use `val, isFound := map[k]` (or `val, isUserExist := map[k]`). Revert any improper usage of `isDefined` for map lookups.
   - **Total Ban on Single-Letter Parameters:** NEVER use single-letter boolean parameters (`v bool`, `b bool`, `val bool`, `flag bool`) in method and function signatures (e.g. setters).
   - **Total Ban on Bare Unprefixed Names:** NEVER use bare verbs, nouns, or adjectives (`stop bool`, `pause bool`, `force bool`, `dryRun bool`, `header bool`).
   - **Mandatory Affirmative Prefixes:** Every boolean parameter, struct field, property, and variable MUST carry an affirmative prefix (`is*` or `has*`):
     - `stop` -> `isStopped`
     - `stopOnFail` -> `isStopOnFail` (e.g. `SetStopOnFail(isStopOnFail bool)`)
     - `defined` -> `isDefined` (e.g. struct field `isDefined bool`, method `IsDefined() bool`)
     - `pause` / `paused` -> `isPaused`
     - `force` -> `isForced` or `isForce`
     - `enable` / `enabled` -> `isEnabled`
     - `dryRun` -> `isDryRun`
     - `debug` -> `isDebug`
     - `verbose` -> `isVerbose`
     - `header` -> `hasHeader`
   - **Map Lookups vs `isDefined`:** For map lookups, always use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]`. NEVER use `isDefined` for map lookups; `isDefined` / `res.IsDefined()` is strictly reserved for replacing inverted `!isEmpty` / `!res.IsEmpty()`.
   - **Total Ban on Compound Negative Chains (`!a || !b || c`):** Chaining inverted negative checks (such as `!state.IsDefined || !state.IsEmpty || state.IsRepo`) violates both discrete assertion rules and positive logic standards. In tests, write discrete assertions; in app code, extract an affirmative composite predicate.

### Generic Code Patterns (Affirmative Naming)

#### Pattern A: Setter Method Parameter & Field Assignment (`v bool` -> `isStopOnFail bool`)

```go
// ❌ ANTI-PATTERN: Single-letter parameter `v bool` and un-prefixed field
func (p *BatchProgress) SetStopOnFail(v bool) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.stopOnFail = v
}

// ✅ REQUIRED: Meaningful, affirmative boolean parameter and property
func (p *BatchProgress) SetStopOnFail(isStopOnFail bool) {
    p.mu.Lock()
    defer p.mu.Unlock()
    p.stopOnFail = isStopOnFail
}
```

#### Pattern B: Generic State Flag & Struct Worker (`stop` -> `isStopped`)

```go
// ❌ ANTI-PATTERN: Bare verb `stop` and lazy `b bool` in stateful worker
type TaskWorker struct {
    stop bool
}

func (w *TaskWorker) SetStop(b bool) {
    w.stop = b
}

func (w *TaskWorker) Run() {
    for {
        if w.stop {
            break
        }
        processTask()
    }
}

// ✅ REQUIRED: Generic affirmative `isStopped` state and parameter
type TaskWorker struct {
    isStopped bool
}

func (w *TaskWorker) SetStopped(isStopped bool) {
    w.isStopped = isStopped
}

func (w *TaskWorker) Run() {
    for {
        if w.isStopped {
            break
        }
        processTask()
    }
}
```

#### Pattern C: Struct Field Definition State (`defined bool` -> `isDefined bool`)

```go
// ❌ ANTI-PATTERN: Bare field name `defined bool` in wrapper struct
type Result[T any] struct {
    value   T
    err     *appfault.AppError
    defined bool // VIOLATION: Bare boolean without is/has prefix
}

// ✅ REQUIRED: Meaningful, affirmative boolean struct field `isDefined bool`
type Result[T any] struct {
    value     T
    err       *appfault.AppError
    isDefined bool // REQUIRED: Explicit affirmative boolean prefix
}
```

#### Pattern D: Generic Transformation Reference Table

| Target Category | ❌ Anti-Pattern (Lazy / Bare) | ✅ Required Affirmative Identifier | Context / Description |
|---|---|---|---|
| Setter Parameter | `SetStopOnFail(v bool)` | `SetStopOnFail(isStopOnFail bool)` | Early termination flag parameter |
| State Variable | `stop := false` | `isStopped := false` | Process / loop cancellation state |
| Method Parameter | `Stop(stop bool)` | `SetStopped(isStopped bool)` | State toggle parameter |
| Struct Field | `defined bool` | `isDefined bool` | Value/record definition presence indicator |
| Method Name | `Defined() bool` | `IsDefined() bool` | Definition verification predicate |
| Struct Field | `pause bool` | `isPaused bool` | Pause / suspend indicator |
| CLI / Config Flag | `force bool` | `isForced bool` | Force override flag |
| Struct Field | `dryRun bool` | `isDryRun bool` | Dry run simulation flag |
| Option Parameter | `debug bool` | `isDebug bool` | Debug mode toggle |
| Struct Field | `header bool` | `hasHeader bool` | Header presence indicator |
| Option Parameter | `records bool` | `hasRecords bool` | Records presence requirement |
| Struct Field | `exists bool` / `isExists bool` | `isDefined bool` | Presence/definition indicator (ban `isExists`) |
| Missing Check | `isUndefined` / `isNotDefined` | `isDefined` (invert with `!isDefined`) | Try IsDefined instead of negatives |
| Missing Value | `hasNoValue` / `isMissing` | `hasValue` / `isDefined` | Affirmative presence check |
| Negative State | `isNotValid` / `isInvalid` | `isValid` (invert with `!isValid`) | Check positive validity |
| Map Comma-Ok | `val, ok` / `val, isExists` | `val, isFound` / `val, isUserExist` | Map lookup presence boolean (revert original name) |
| Non-Empty / Data Present | `!isEmpty` / `!res.IsEmpty()` | `isDefined` / `res.IsDefined()` | Mandatory: Affirmative IsDefined instead of inverted !isEmpty |

#### Pattern E: `IsDefined` vs `IsExists` & Compound Negative Decomposition

```go
// ❌ FORBIDDEN: Compound negative chain and awkward isExists in test assertions
if !state.IsExists || !state.IsEmpty || state.IsRepo {
    t.Errorf("expected empty non-repo directory: %+v", state)
}

// ✅ REQUIRED: Affirmative IsDefined field + discrete individual assertions
if !state.IsDefined {
    t.Errorf("expected directory to be defined: %+v", state)
}

if !state.IsEmpty {
    t.Errorf("expected directory to be empty: %+v", state)
}

if state.IsRepo {
    t.Errorf("expected non-repo directory: %+v", state)
}

// ❌ FORBIDDEN: Compound negative in application logic
if !params.State.IsExists || params.State.IsEmpty {
    performFreshClone(params)
}

// ✅ REQUIRED: Extract affirmative composite predicate
isCloneTargetFresh := !params.State.IsDefined || params.State.IsEmpty

if isCloneTargetFresh {
    performFreshClone(params)
}
```

#### Pattern F: Boolean Evaluation BEFORE the `if` Statement (Simple One-Variable `if` Checking)

```go
// ❌ BANNED ANTI-PATTERN:
// 1. Cramming type assertion assignment and compound condition into one line.
// 2. Hardcoding magic strings ("Version", "version", "unknown") inline.
// 3. Returning raw fallback literal instead of a defined constant.
func extractVersionValue(rawMap map[string]interface{}) string {
    for _, key := range []string{"Version", "version"} {
        if v, isString := rawMap[key].(string); isString && len(v) > 0 {
            return v
        }
    }

    return "unknown"
}

// ✅ MANDATORY CLEAN PATTERN:
// 1. Zero magic strings: extract lookup keys and defaults into constants.
// 2. Merge repeated/related strings into reusable collections (versionKeys).
// 3. Assignment on its own dedicated line.
// 4. Affirmative boolean (hasContent) pre-evaluated BEFORE the if statement.
// 5. Clean vertical breathing room (blank line before if).
// 6. Dead-simple if statement evaluating exactly ONE variable.
// 7. Return defined constant (VersionUnknown) instead of raw magic string literal.
const (
    VersionUnknown  = "unknown"
    versionKeyUpper = "Version"
    versionKeyLower = "version"
)

var versionKeys = []string{versionKeyUpper, versionKeyLower}

func extractVersionValue(rawMap map[string]interface{}) string {
    for _, key := range versionKeys {
        v, isString := rawMap[key].(string)
        hasContent := isString && len(v) > 0

        if hasContent {
            return v
        }
    }

    return VersionUnknown
}
```

3. **No Inverted Success Checks:**
   - Never invert positive success checks (e.g. `!response.isSuccess`).
   - Use explicit failure states (e.g. `response.isFail`, `isError`).

4. **Zero Tolerance for Nested `if` (Nesting Depth <= 1):**
   - No `if` statements inside another `if` block.
   - Flatten all conditionals using guard clauses and early returns.
   - Never combine mixed polarity (`if isA && !isB` -> split into separate guard clauses).

5. **Enum Suffix `Type`:**
   - All enum declarations across TypeScript, Go, PHP must end with `Type` (e.g. `UserRoleType`, `CommandStatusType`).

6. **Function and File Size Caps:**
   - Functions: <= 8 lines preferred, <= 15 lines maximum.
   - Files: <= 100 lines coding maximum (recommended <= 80 lines).
   - Zero line compression (no single-line `if/else`, no deleted blank lines).

---

## 3. Step 0 Preflight with Data Schema Verification

Before claiming subtasks or refactoring booleans, inspect and verify the SQLite schema:

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
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <refactor-rationale>" -s "EXECUTION"

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

## 5. Fast File Discovery & Reading Toolchain (GitMap Primacy)

All code searches MUST use GitMap. Shell search tools (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`) are strictly banned:

- **Live Streaming Search:** `gitmap aum search "<pattern>" [dir] [-e <.ext>] [-r] [-i]`
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]`
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>` (streams to stdout with zero disk writes)
- **Instant Code Search:** `gitmap search "<term>"` (immediate multi-core filesystem walk)

---

## 6. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts, Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **NO RUNNER SCRIPTS (TOTAL BAN):** NEVER launch background test runners, worker pools, or test inventory loops during routine execution.
- **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 7. Verification Checklist

- [ ] All booleans use positive prefixes (`is`, `has`) only; zero negative boolean identifiers.
- [ ] Implicit boolean checks enforced (`if isReady`, never `if isReady == true`).
- [ ] No mixed polarity conditions (`if isA && !isB` is banned).
- [ ] Enums end with `Type` suffix.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
