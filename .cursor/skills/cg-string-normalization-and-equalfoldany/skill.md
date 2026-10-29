---
name: cg-string-normalization-and-equalfoldany
description: Autonomously audits, refactors, and standardizes string normalization, case-insensitive comparisons, and chained equality into reusable strutil helpers (EqualFoldAny, EqualFoldAnyTrim) across Go, TypeScript, and Rust.
---

# Skill: String Normalization, `EqualFoldAny` & Centralized Utility Discovery (`cg-string-normalization-and-equalfoldany`)

> **Invocation:** `@[cg-string-normalization-and-equalfoldany](file;.cursor/skills/cg-string-normalization-and-equalfoldany)`

This skill governs autonomous execution for auditing, refactoring, and standardizing repetitive string normalization (`ToLower`, `TrimSpace`) and chained equality checks (`EqualFold` OR-chains) across Go, TypeScript, and Rust codebases.

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

## 2. Search First Protocol for `<repo>/pkg/strutil/strutil.go`

Before modifying call sites, creating local helpers, or introducing duplicate logic, agents must strictly follow the **Search First Protocol**:

### 2.1 Canonical Package Locations

| Language | Primary Canonical Path | Permitted Internal Path | Prohibited Ad-Hoc Paths |
| :--- | :--- | :--- | :--- |
| **Go** | `<repo>/pkg/strutil/strutil.go` | `<repo>/internal/strutil/strutil.go` | `cmd/*.go`, `<repo>/pkg/util/helpers.go`, inline file helpers |
| **TypeScript** | `<repo>/src/lib/strutil.ts` | `<repo>/src/utils/strutil.ts` | `src/components/*.tsx`, `src/pages/*.ts`, local helpers |
| **Rust** | `<repo>/src/util/strutil.rs` | `<repo>/src/strutil/mod.rs` | `<repo>/src/main.rs`, inline sub-modules in `<repo>/src/cmd/` |
| **Python** | `<repo>/pkg/strutil/strutil.py` | `<repo>/src/strutil.py` | `scripts/*.py`, inline functions in individual script files |

### 2.2 Protocol Execution Flow

1. **Locate Canonical String Utility:**
   - Scan for string utilities using `gitmap find-files-any "strutil"`.
   - In Go: Check `<repo>/pkg/strutil/strutil.go` or `<repo>/04-code/golang/pkg/strutil/strutil.go`.
   - In TypeScript: Check `<repo>/src/lib/strutil.ts` or `<repo>/src/utils/strutil.ts`.
   - In Rust: Check `<repo>/src/util/strutil.rs` or `crate::strutil`.
2. **Re-use Mandate:**
   - All refactored call sites MUST import and use the canonical utility functions (`strutil.EqualFoldAny`, `strutil.EqualFoldAnyTrim`, `strutil.NormalizeLowerTrim`).
   - Total ban on authoring duplicate unshared functions (`isYes`, `checkConfirm`, `matchesAny`, `cleanEqual`) inside consumer packages.
3. **Absence Protocol:**
   - If the canonical utility is missing in the target repository, agents must create or register it following `02-spec/02-coding-guidelines/01-cross-language/34-string-normalization-and-equalfoldany.md` BEFORE modifying call sites.

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

## 5. Polyglot Code Transformation Catalog

### 5.1 Go (`<repo>/pkg/strutil/strutil.go`)

```go
// ✅ REQUIRED — single-line atomic check via canonical strutil helper
isElevated := strutil.EqualFoldAnyTrim(userRole, "admin", "owner", "moderator")
if isElevated {
    grantPermission()
}
```

### 5.2 TypeScript (`<repo>/src/lib/strutil.ts`)

```typescript
// ✅ REQUIRED — concise variadic helper invocation
const isAffirmative = strutil.equalFoldAnyTrim(input, 'yes', 'y', 'true', '1');
if (isAffirmative) {
  enableFeature();
}
```

### 5.3 Rust (`<repo>/src/util/strutil.rs`)

```rust
// ✅ REQUIRED — slice-based case-folding helper
let is_pending = strutil::equal_fold_any_trim(input, &["draft", "review"]);
if is_pending {
    process_document();
}
```

---

## 6. High-Speed AST Search & Discovery via GitMap

All searches MUST use GitMap. Shell search tools (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`) are strictly banned:

```bash
# Go: Find chained EqualFold checks
gitmap aum search -r "EqualFold\(.*\|\|.*EqualFold\(" [dir] -e .go

# Go: Find TrimSpace followed by EqualFold
gitmap aum search -r "TrimSpace\(.*EqualFold" [dir] -e .go

# Go: Find ToLower equality OR-chains
gitmap aum search -r "strings\.ToLower\(.*==.*\|\|" [dir] -e .go

# TypeScript: Find chained .toLowerCase() checks
gitmap aum search -r "\.toLowerCase\(\)\s*===\s*.*\|\|" [dir] -e .ts

# Rust: Find chained eq_ignore_ascii_case checks
gitmap aum search -r "eq_ignore_ascii_case\(.*\|\|" [dir] -e .rs
```

---

## 7. Routine Execution & GitMap Hyphen Commit Policy

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

## 8. Verification Checklist

- [ ] Canonical `<repo>/pkg/strutil/strutil.go` verified before refactoring call sites.
- [ ] Chained `EqualFold` and `ToLower` OR-chains converted to `strutil.EqualFoldAny` or `strutil.EqualFoldAnyTrim`.
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
