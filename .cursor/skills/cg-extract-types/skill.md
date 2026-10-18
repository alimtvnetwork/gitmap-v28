---
name: cg-extract-types
description: Autonomously audits, extracts, and centralizes domain structs, raw generic instantiations, and Result wrappers into dedicated types.go files as single reusable named types across the codebase.
---

# Skill: Extracting Generic Types & Models to `types.go` (`cg-extract-types`)

This skill governs autonomous auditing, extraction, and centralization of scattered domain structs, raw generic Result envelopes (`ResultSlice[T]`, `ResultMap[K, V]`, `Result[T]`, `Wrap[T]`), and repeated generic type parameters into dedicated, package-level `types.go` files (or language-equivalent leaf `types/` modules) as single reusable named types.

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

1. **Dedicated `types.go` Single Reusable Type Definition:**
   - Every package managing domain models, payloads, or Result envelopes MUST define them inside a dedicated `types.go` file within the package directory as a single reusable named type.
   - Never declare ad-hoc unexported structs or raw generic Result envelopes inline in implementation files (`importer.go`, `store.go`, `service.go`).
   - Sibling implementation files in the package import or consume the single named type directly without re-specifying generic parameters (`ResultSlice[...]`).

2. **Single Reusable Type Alias for Repeated Generics:**
   - When a generic type is used 3 or more times with the same type argument, or wraps a domain payload struct, create a canonical type alias in `types.go`.
   - Pattern:
     ```go
     type (
         ScheduleExportBundle struct { ... }
         ScheduleExportBundleResult = result.ResultSlice[ScheduleExportBundle]
         ScheduleExportBundleSingleResult = result.Wrap[ScheduleExportBundle]
     )
     ```

3. **Exported Models & Affirmative Field Naming:**
   - All domain payload structs used across function boundaries MUST be exported with PascalCase names.
   - All boolean fields in models, structs, and envelopes MUST carry affirmative prefixes (`is*` or `has*`):
     - `isDefined bool` (TOTAL BAN on bare `defined bool`)
     - `hasRecords bool` / `hasRecord bool`
     - `isStopped bool` (TOTAL BAN on bare `stop bool`)
     - `isForced bool` / `isPaused bool` / `isEnabled bool`

4. **Pointer-Attached Null Safety & Fluent Predicates:**
   - Result containers returned by functions (`Result[T]`, `ResultSlice[T]`, `ResultMap[K, V]`) must have inspection methods attached to pointer receivers.
   - Code consuming `types.go` aliases must utilize the 4 fluent pointer-safe predicates:
     - `res.IsCountOtherThan(expectedCount int) bool`: Returns `true` if operation failed or count != expected.
     - `res.IsEmpty() bool`: Returns `true` if empty collection, nil/zero payload, or nil receiver.
     - `res.HasRecord() bool` (or `res.HasRecords()`): Returns `true` if succeeded and count > 0.
     - `res.IsDefined() bool`: Returns `true` if succeeded and payload is non-null/non-empty with count > 0.
   - **Mandatory Replacement for `!isEmpty`:** NEVER use inverted negative empty checks (`!isEmpty`, `!res.IsEmpty()`). Always use affirmative `isDefined` / `res.IsDefined()` when asserting that data or records are present. Map lookups use `val, isFound := userMap[id]` or `val, isUserExist := userMap[id]` (never bare `ok`, and never use `isDefined` for map lookups).

5. **Clean Package Boundaries & Zero Dependency Cycles:**
   - `types.go` files are leaf definitions within their package.
   - `types.go` MUST only import foundation utility packages (e.g. `pkg/result`, `pkg/appfault`) and standard library packages.
   - `types.go` MUST NEVER import heavy business logic, services, or circular sibling packages.

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
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <type-extraction-rationale>" -s "EXECUTION"

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

## 5. High-Speed Discovery via GitMap Primacy

All type searches MUST use GitMap. Shell search tools (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`) are strictly banned:

```bash
# Inventory unexported structs in Go implementation files
gitmap aum search -r "type [a-z][A-Za-z0-9]+ struct \{" [dir] -e .go

# Discover raw generic ResultSlice returns outside types.go
gitmap aum search "result.ResultSlice[" [dir] -e .go

# Discover raw generic Wrap returns outside types.go
gitmap aum search "result.Wrap[" [dir] -e .go

# Discover raw generic ResultMap returns outside types.go
gitmap aum search "appfault.ResultMap[" [dir] -e .go

# Universal File Search
gitmap find "types.go"
```

---

## 6. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts, Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) to verify compilation. Build verification is checked later on in CI/CD.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 7. Verification Checklist

- [ ] All package domain structs and Result wrappers centralized in `types.go`.
- [ ] Single reusable named type aliases exported and consumed across implementation files.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
