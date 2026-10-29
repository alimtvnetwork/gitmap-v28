---
name: cg-variadic-and-spread-parameters
description: Autonomously audits, refactors, and enforces variadic and spread parameters (...T, ...items, &[T]) across polyglot codebases to eliminate artificial slice wrappers.
---

# Skill: Variadic & Spread Parameters, Rest Elements & Slices (`cg-variadic-and-spread-parameters`)

> **Invocation:** `@[cg-variadic-and-spread-parameters](file;.agents/skills/cg-variadic-and-spread-parameters)`

This skill governs autonomous execution for auditing, refactoring, and enforcing variadic parameters, rest elements, and slice reference patterns across Go, TypeScript, and Rust codebases.

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

1. **Variadic Parameter Primacy:**
   - Functions accepting homogeneous collections of items MUST use variadic signatures (`...T` in Go, `...items: readonly T[]` in TypeScript, `&[T]` in Rust) instead of rigid slice or array types whenever single-item or variable-length calls occur.
   - Total ban on forcing callers to construct artificial wrapper slices/arrays for single-element invocations.

2. **Zero-Wrapper Single Invocations:**
   - Single-item callers MUST pass the argument directly without wrapper syntax (`DeleteUsers(id)` instead of `DeleteUsers([]string{id})`).

3. **Seamless Spread Forwarding:**
   - Callers possessing an existing slice or array must be able to unpack it directly using native spread syntax (`DeleteUsers(allIDs...)` in Go, `deleteUsers(...allIds)` in TypeScript, `process_records(&records)` in Rust).

4. **Trailing Position & Parameter Limits:**
   - The variadic parameter MUST be the final parameter in the signature.
   - Total function parameters must not exceed 2–3 loose arguments. Pack complex configuration options into a dedicated `*Params` struct placed before the variadic parameter.

5. **Safe Empty & Nil Handling:**
   - Functions MUST handle zero arguments cleanly without error.
   - Affirmative length checks must guard processing (`hasItems := len(items) > 0`).
   - If not `hasItems`, return early cleanly.

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
gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <variadic-refactor-rationale>" -s "EXECUTION"

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

## 5. Polyglot Code Transformations

### 5.1 Go (`...T` Variadic Slices)

```go
// ✅ REQUIRED — variadic slice parameter with clean affirmative check
func PurgeTags(resourceID string, tags ...string) *appfault.AppError {
    hasTags := len(tags) > 0
    if !hasTags {
        return nil
    }

    for _, tag := range tags {
        appErr := removeTag(resourceID, tag)
        if appErr != nil {
            return appErr
        }
    }

    return nil
}
```

### 5.2 TypeScript (`...items: readonly T[]`)

```typescript
// ✅ REQUIRED — Rest parameter with readonly safety
export async function trackEvents(
  ...events: readonly string[]
): Promise<Result<void>> {
  const hasEvents = events.length > 0;
  if (!hasEvents) {
    return Result.ok(undefined);
  }

  for (const event of events) {
    const result = await dispatchEvent(event);
    if (!result.isSuccess) {
      return result;
    }
  }

  return Result.ok(undefined);
}
```

---

## 6. High-Speed AST Search via GitMap Primacy

All searches MUST use GitMap. Shell search tools (`rg`, `ripgrep`, `grep`, `git grep`, `Select-String`, `findstr`) are strictly banned:

```bash
# Go: Find function definitions with slice parameters
gitmap aum search -r "func [A-Za-z0-9_]+\([^)]*\[\][A-Za-z0-9_.*]+(\s*\))?" -e .go

# Go: Find call sites wrapping single items in slice literals
gitmap aum search -r "\b[A-Za-z0-9_]+\(\[\]string\{" -e .go

# TypeScript: Find function definitions with array parameters
gitmap aum search -r "\bfunction [A-Za-z0-9_]+\([^)]*:\s*[A-Za-z0-9_]+\[\]" -e .ts

# Rust: Find functions taking owned Vec<T>
gitmap aum search -r "fn [A-Za-z0-9_]+\([^)]*Vec<" -e .rs
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

- [ ] Variadic parameters placed in trailing position.
- [ ] No more than 2–3 loose parameters (use `*Params` struct for configuration).
- [ ] Call sites refactored to eliminate artificial slice wrappers (`[]string{x}`).
- [ ] Preflight schema verified via `gitmap agent task status --ddl` and `gitmap agent task status --json`.
- [ ] 3-Tier Database Location Transparency contract fulfilled.
- [ ] HEATMAP TELEMETRY & CONTENTION INSPECTION: Verified that multi-agent contention was inspected via gitmap heatmap discuss (or gitmap agent heatmap) before editing shared files, and that all file exploration and modifications were logged with reasoning via gitmap heatmap log (or gitmap agent log) into .ai-memory/heatmap.db for crash forensics and repo-cache synchronization.
