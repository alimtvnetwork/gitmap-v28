---
name: cg-file-size-and-function-reduction
description: Autonomously audits, decomposes, and refactors files over 100 lines and functions over 8–15 lines across the codebase using a two-part decomposition strategy, wrapper objects, and strict whitespace preservation without intermediate test or build runs.
---

# Skill: File Size & Function Size Reduction (`cg-file-size-and-function-reduction`)

This skill governs autonomous auditing, decomposition, and refactoring of bloated source code files exceeding the **100-line standard file cap** (recommended <= 80 lines) and functions exceeding **8–15 lines** across polyglot codebases (Go, TypeScript, Python, PHP, Rust, C#).

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

### 1. The Two-Part Decomposition Strategy
Never attempt to split a file without first decomposing its functions:
- **Part 1 (Function Level):** Identify functions > 15 lines and decompose them into focused, single-responsibility helper functions targeting **<= 8 lines** (hard cap: 15 lines).
- **Part 2 (File Level):** If the file remains > 100 lines after function decomposition, group cohesive clusters of helper functions and extract them into dedicated sibling files (e.g. `<name>_validator.go`, `<name>_helpers.go`, `<name>_converter.go`).

### 2. Strict Allowed Exceptions (Exempt from 100-Line Limit)
The following file categories are explicitly permitted to exceed 100 lines:
- **Data & Schema Files:** JSON files (`*.json`), YAML/TOML configuration manifests.
- **Type Definition Files:** Dedicated types files (`types.go`, `*.d.ts`, `types.ts`, `models.go`, `schema.ts`) that centralize domain structs or discriminated unions.
- **Lookup & Map Files:** Static lookup tables, dictionary mappings, registries, and transition matrices (`map.go`, `registry.go`, `lookup.ts`, `routes.go`).
- **Variables & Constants Files:** Centralized constant packages, enum registrations, and configuration variable dictionaries (`vars.go`, `consts.go`, `constants.ts`, `vars.ts`).

All logic, services, controllers, handlers, utilities, hooks, CLI commands, and scripts MUST remain strictly under 100 lines.

### 3. Total Ban on Line-Compression Cheating
- **Zero Whitespace Stripping:** NEVER delete blank lines, squash vertical spacing, or compress `if/else` statements onto a single line to reduce line count.
- **No Inline Compound Condition Cramming:** NEVER cram variable declarations, type assertions, and compound conditions into the `if` header (e.g. `if v, isString := rawMap[key].(string); isString && len(v) > 0 {`). Put assignments on separate lines, evaluate booleans before `if`, and keep `if` conditions simple with one variable.
- **Zero Magic Strings & Constant Returns:** Never hardcode raw string literals (`"unknown"`, `"version"`) into function logic or return statements. Merge them into named constants/slices and return defined constants (`return VersionUnknown`).
- **Return New Line Concept (Mandatory):**
  - Exactly **ONE blank line BEFORE** every `return`, `throw`, or `break` statement.
  - Exactly **ONE blank line AFTER** every closing curly brace `}` of an `if`, `for`, `switch`, or helper block.
- Files must be reduced by modular decomposition into separate files, NOT by code compression.

### 4. Wrapper Objects & Parameter Reduction Pattern
- **Return Wrapper Struct:** If an extracted helper returns 2 or more related values (beyond standard error/Result), encapsulate them into a dedicated named struct or type.
- **Parameter Structs:** If an extracted function requires more than 2 parameters, combine them into an options or params struct (`*Params` / `*Options`).

### 5. Boolean Principles During Extraction
- **Implicit Positive Checks Only:** NEVER write `if isReady == true`. ALWAYS write `if isReady`.
- **No Negative Polarity:** NEVER combine a positive and negative condition (`if isA && !isB`). Split into separate guard clauses!
- **Positive Naming Only:** All booleans MUST use affirmative prefixes: `is*` or `has*` only (e.g., `isFail`, `isInvalid`, `isPaused`).
- **Zero Nested If Statements:** Every extracted function MUST maintain a nesting depth <= 1 using early guard returns.

---

## 3. Step 0 Preflight with Data Schema Verification

Before claiming subtasks or decomposing files, initialize and verify the SQLite task database schema:

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
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <decomposition-rationale>" -s "EXECUTION"

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

- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]`
- **List Indexed Files:** `gitmap list-files [pattern]` (alias `gitmap lf [pattern] [-ext <ext>]`)
- **Substring Match:** `gitmap find-files-any "<substring>"` (alias `gitmap ffa "<str>"`)
- **Stream File Content:** `gitmap cat <filepath>`
- **Check File Sizes:** `python 03-ai-scripts/13-file-size-guard.py`

---

## 6. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts, Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) after individual file edits. Build verification is checked ONLY at the final step.
- **NO LINE COMPRESSION (TOTAL BAN):** NEVER remove blank lines, merge statements, or compress `if/else` blocks to reduce line count.
- **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 7. Verification Checklist

- [ ] All functions reduced to <= 8 lines (hard cap: 15 lines).
- [ ] Files reduced to <= 100 lines (excluding allowed exceptions).
- [ ] Vertical line spacing preserved (Return New Line before return, after `}`).
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
