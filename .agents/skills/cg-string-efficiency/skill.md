---
name: cg-string-efficiency
description: Autonomously audits, optimizes, and refactors string operations and comparisons across polyglot codebases for zero-allocation case folding, short-circuiting lazy evaluation, and loop hoisting while strictly preserving business logic.
---

# Skill: String Operations & Memory Efficiency (`cg-string-efficiency`)

> **Invocation:** `@[cg-string-efficiency](file;.agents/skills/cg-string-efficiency)`

This skill governs autonomous auditing, optimization, and refactoring of inefficient string operations, comparisons, case conversions, and loop allocations across polyglot codebases (Go, TypeScript, Python, PHP, Rust, C#).

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

### 1. Prime Directive: Business Logic Safety First (Zero-Regression Mandate)
Code efficiency optimizations MUST NEVER alter program behavior. Before modifying any string check, the agent MUST analyze whether the original code required:
- **Exact Case-Insensitive Equality:** Use zero-allocation folding (`strings.EqualFold` in Go).
- **Case-Insensitive Substring Containment:** DO NOT change to `EqualFold`! Keep substring search semantics (`Contains`), but optimize via lazy short-circuiting and loop hoisting.
- **Prefix / Suffix Matching:** Use `HasPrefix` / `HasSuffix` or `startsWith` / `endsWith`.

### 2. Zero-Allocation Case-Folding
- **Go:** Replace `strings.ToLower(a) == strings.ToLower(b)` or `strings.ToLower(s) == "val"` with `strings.EqualFold(a, b)` (0 allocations, ~10x faster).
- **TypeScript/JS:** Avoid lowering both sides when comparing with a constant literal (`status.toLowerCase() === "active"`).
- **Python:** Use `a.casefold() == b.casefold()` for full Unicode case-folding.
- **Rust:** Replace `a.to_lowercase() == b.to_lowercase()` with `a.eq_ignore_ascii_case(&b)` (0 heap allocations).
- **C#:** Replace `a.ToLower() == b.ToLower()` with `string.Equals(a, b, StringComparison.OrdinalIgnoreCase)`.
- **PHP:** Replace `strtolower($a) === strtolower($b)` with `strcasecmp($a, $b) === 0`.

### 3. Short-Circuiting Lazy Evaluation
In multi-field matching predicates, never compute all lowerings or searches eagerly:
- Return immediately on the first true branch using early guard returns.
- Avoid executing secondary case conversions when the first condition already matches.

### 4. Loop Hoisting of Case Conversions
- Never repeatedly convert the same search pattern or filter string inside an iteration.
- Lower the pattern once prior to the loop and pass the pre-lowered term to item matchers.

### 5. String Builders for Cumulative Concatenation
- In Go, replace `s += chunk` inside loops with `strings.Builder`.
- In C#, use `StringBuilder`. In Python, use `''.join(parts)`. In TS, use array accumulator `parts.join('')`.

### 6. Function Sizing & File Hygiene
- Functions target **<= 8 lines** of body logic (hard cap: 15 lines).
- Files target **<= 80 lines** (hard cap: 100 lines, excluding allowed exceptions).
- Whitespace preservation: exactly ONE blank line before return, exactly ONE blank line after closing brace `}`.

---

## 3. Step 0 Preflight with Data Schema Verification

Before claiming subtasks or refactoring string operations, initialize and verify the SQLite task database schema:

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
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <string-optimization-rationale>" -s "EXECUTION"

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

All string operation searches MUST use GitMap. Shell search tools (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`) are strictly banned:

- **Discover Case Lowerings in Go:** `gitmap aum search "strings.ToLower" [dir] -e .go`
- **Discover Case Lowerings in TypeScript:** `gitmap aum search "toLowerCase()" [dir] -e .ts`
- **Discover EqualFold Candidates:** `gitmap aum search -r "strings\.ToLower\(.*==.*strings\.ToLower" [dir] -e .go`
- **Universal File Search:** `gitmap find "<pattern>" [-ext <ext>]`
- **Stream File Content:** `gitmap cat <filepath>`

---

## 6. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine coding guideline execution turns or micro-batch loops. Verify code strictly using targeted file-level linters / autofixers on the specific modified files.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run any tests using Python scripts, Go (`go test ./...`), or any test runner during routine execution turns. Testing is strictly checked later on in CI/CD.
- **NO INTERMEDIATE BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`, compiler checks) after individual file edits. Build verification is checked ONLY at the final step.
- **NO AUTOMATIC RELEASES (TOTAL BAN):** NEVER bump versions, update changelogs, or trigger releases unless explicitly commanded by the user.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - Chore: `gitmap cpc "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.

---

## 7. Verification Checklist

- [ ] Zero-allocation `strings.EqualFold` used for case-insensitive equality; business logic semantics preserved.
- [ ] Loop hoisting of invariant string normalizations applied.
- [ ] String concatenations in loops converted to `strings.Builder`.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
