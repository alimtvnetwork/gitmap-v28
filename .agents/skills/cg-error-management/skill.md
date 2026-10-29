---
name: cg-error-management
description: >-
  Autonomously audits, refactors, and validates repository-wide error management against 02-spec/03-error-manage/ using AppError wrappers, universal response envelopes, and CI linters.
---

# Error Management & Architecture Coding Guidelines (`cg-error-management`)

> **Invocation:** `@[cg-error-management](file;.agents/skills/cg-error-management)`

This skill provides autonomous audit, refactoring, and validation of repository-wide error handling based on `02-spec/03-error-manage/` and `.ai-memory/coding-guidelines.md`.

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

## 2. Core Architectural Invariants

1. **Top-Instruction Priority Mandate (Preamble Precedence):** Any directive, constraint, or checklist placed at the top of an incoming prompt, header alert block, or user request represents an absolute MUST FOLLOW mandate that takes highest priority and strictly overrides generic guidelines below it.
2. **Strict Golang Error Wrapping Mandate (`*appfault.AppError` & `appfault.Fault`):**
   - Whenever ANY Go function encounters, intercepts, or receives an error (from stdlib `os`, `io`, `json`, `sql`, `net`, or downstream packages), it MUST be immediately embedded and wrapped using `appfault.Wrap(errType, err, "opName")` or `result.WrapFailure[T]`.
   - Raw standard library `error` returns (`return err`) in domain packages are strictly prohibited.
   - Functions with side-effects only MUST return `*appfault.AppError`. Functions returning data and possible error MUST return `result.Wrap[T]`.
3. **Zero Swallowed Errors Policy (TOTAL BAN):**
   - NEVER swallow, suppress, or silently ignore errors under any circumstances.
   - NO empty `catch` or `except:` blocks.
   - NO blank identifier error discards (`_ = err` or `val, _ := fn()`).
   - NO returning fallback default values (`return nil`, `return ""`, `return false`) to mask an underlying error without caller notification.
   - Every caught or received error MUST either be completely resolved with structured context logging (operation name, input parameters) OR embedded/wrapped in `*appfault.AppError` and returned to the caller.
4. **No Bare Panics or Bare Exits**: Zero calls to `panic("...")`, `panic(err)`, or `os.Exit(...)` outside the central dispatcher (`cliexit.HandleError`).
5. **Context-Rich Metadata**: All errors MUST use structured `*appfault.AppError` from `04-code/golang/pkg/appfault` and be wrapped with `Op`, `Code`, `Type`, `Severity`, `Creator`, `Message`, `Ctx`, and `Cause`.
6. **Single Result Container Return Types (`pkg/appfault`) & `types.go` Mandate**:
   - Multi-value returns returning errors (`(map[K]V, error)`, `([]T, error)`, `(T, error)`) are strictly banned.
   - Functions returning maps MUST return `appfault.ResultMap[K, V]` (or domain alias).
   - Functions returning slices MUST return `appfault.ResultSlice[T]` (or domain alias).
   - Functions returning single values MUST return `appfault.Result[T]` (or domain alias).
   - **Mandatory `types.go` Single Reusable Type Definition:** All domain payload structs and repeated generic Result envelopes MUST be defined in a dedicated `types.go` file within the package as a single reusable named type everywhere.
7. **Universal Response Envelope**: All API endpoints return `{ "data": ..., "errors": [...], "meta": ... }`.
8. **Targeted Verification**: Continuous verification via `python linter-scripts/check-error-management.py <files>`. DO NOT run the full CI/CD pipeline runner (`06-cicd-local-runner.py`) during routine fixes.
9. **No Releases**: Strictly forbidden from bumping versions or cutting releases at the end of this task.
10. **Atomic Change Tracking**: Append all modified files to `.ai-memory/temp/recent-file-changes.json` under lock (`python 03-ai-scripts/33-test-inventory-generator.py --record <files...>`), mapping to associated tests in `.ai-memory/test-inventory.json`.
11. **Typed Errors for Unknown Enum Values (canonical example)**: Unknown enum value, failed `Parse`/`TryParse`, or unhandled switch case over an action enum MUST return a typed error (the repo's `*appfault.AppError` / result convention) with the offending value in the message — never a panic, never a silent default, never a swallowed error. Canonical pattern: `action, err := ActionTypeParse(raw)` → on failure wrap and return; exhaustive switch `default:` → same typed error.
12. **Repository Error Wrapping (canonical example)**: Errors crossing a repository boundary MUST be wrapped at that boundary with context naming the repository and the operation (e.g. `appfault.Wrap(errType, err, "UserRepository.FindByID")`) before reaching services. Unwrapped repository/driver errors in the service layer are strictly prohibited.

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

- [ ] Structured Go errors return `*appfault.AppError`, never bare `error`.
- [ ] Functions return monadic `Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`; zero multi-value error tuples.
- [ ] Dedicated `types.go` centralizes domain structs and Result wrappers.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
- [ ] Unknown enum value / failed parse / unhandled switch case returns a typed `AppError` carrying the offending value; zero panics, zero silent defaults.
- [ ] Repository errors wrapped with repository + operation context at the repository boundary; zero unwrapped driver errors in services.
