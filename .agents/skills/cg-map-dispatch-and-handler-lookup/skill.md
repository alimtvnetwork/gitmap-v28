---
name: cg-map-dispatch-and-handler-lookup
description: Autonomously audits, refactors, and replaces repetitive chained if-else if ladders and switch statements matching keys or tokens with declarative, high-efficiency map dispatch tables and discrete handler functions across polyglot codebases.
---

# Skill: Map Dispatch & Handler Table Patterns (`cg-map-dispatch-and-handler-lookup`)

> **Invocation:** `@[cg-map-dispatch-and-handler-lookup](file;.agents/skills/cg-map-dispatch-and-handler-lookup)`

This skill governs autonomous execution for auditing, refactoring, and standardizing repetitive chained conditional ladders (`if ... else if ... else if`) and verbose `switch-case` statements that match keys, event codes, or string tokens to invoke discrete functions, replacing them with declarative map lookup tables across TypeScript, Go, Python, Rust, and PHP.

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

## 2. Core Execution Instructions

1. Master Specification: `02-spec/02-coding-guidelines/01-cross-language/35-map-dispatch-and-handler-lookup.md`.
2. Source Prompt: `01-prompts/15-cg-execute/38-map-dispatch-and-handler-lookup.md` (and v4 mirror).
3. Search codebase exclusively via GitMap (`gitmap aum search "<pattern>" [dir] [-e <.ext>]`, `gitmap find`); TOTAL BAN on `rg`, `ripgrep`, `grep`, `git grep`, `Select-String`.
4. Turn 1 MUST showcase the Extreme Understood Turn-1 Showcase featuring all 5 Mandatory Pillars and 4-part verification schema.
5. Phase 1 Step 3 Brain & Comprehensive Roadmap Creation mandate: exhaustive inventory of target files, handler functions, and disjoint subtasks.
6. Immediate Dynamic Respawn & Dynamic Collaborative Join: 100% duty cycle, zero idle capacity.

---

## 3. Step 0 Preflight with Data Schema & Heatmap Verification

Before enqueuing subtasks or claiming work, initialize and verify the SQLite task database schema and companion heatmap telemetry:

```bash
# Initialize SQLite task engine
gitmap agent task init --name "<task name>" --budget 400
gitmap agent task status --ddl
gitmap agent task status --json

# Mandatory Heatmap & Companion RC Preflight
gitmap rc heatmap tables
gitmap heatmap status
gitmap agent heatmap
```

- **Contention Inspection:** Check active write leases and contention hotspots via `gitmap agent heatmap` and `gitmap heatmap status`.
- **Cognitive Action Logging:** Record file exploration touches and editing decisions with explicit reasoning:
  ```bash
  gitmap heatmap log -t <task-slug> -a "READ" -f "<path>" -r "Reasoning: <why>" -s "EXPLORATION"
  gitmap heatmap log -t <task-slug> -a "EDIT" -f "<path>" -r "Decision: <why>" -s "EXECUTION"
  gitmap rc heatmap log -t <task-slug> -a "<READ|EDIT>" -f "<path>" -r "<reasoning>"
  ```

---

## 4. The Map Dispatch Transformation Standard

### 4.1 Anti-Pattern: Chained Branching Dispatch
```typescript
// ❌ ANTI-PATTERN: High cyclomatic complexity, repetitive boilerplate
const handleKeyDown = (e: KeyboardEvent) => {
  if (e.key === 'Escape') {
    e.preventDefault();
    onClose();
  } else if (e.key === '+' || e.key === '=') {
    e.preventDefault();
    handleZoomIn();
  } else if (e.key === '-') {
    e.preventDefault();
    handleZoomOut();
  } else if (e.key === '0') {
    e.preventDefault();
    handleResetZoom();
  }
};
```

### 4.2 Canonical Solution: Declarative Map Lookup
```typescript
// ✅ CANONICAL PATTERN: Clean, single-line lookup and invocation
const KEY_ACTIONS: Record<string, () => void> = {
  Escape: onClose,
  '+': handleZoomIn,
  '=': handleZoomIn,
  '-': handleZoomOut,
  '0': handleResetZoom,
};

const handleKeyDown = (e: KeyboardEvent) => {
  const action = KEY_ACTIONS[e.key];
  if (action) {
    e.preventDefault();
    action();
  }
};
```

### 4.3 Multi-Language Dispatch Reference
- **Go:** `var commandRegistry = map[string]func(*Context) *appfault.AppError`
- **Python:** `EVENT_DISPATCH: dict[str, Callable[[dict], None]] = {"event": handler}`
- **Rust:** `const ACTION_MAP: &[(&str, ActionHandler)] = &[("key", handler)];`
- **PHP:** `$handlers = ['action' => fn() => $this->handleAction()];`

---

## 5. Routine Execution & GitMap Hyphen Commit Policy

- **NO FULL CI/CD RUNNER (Strict Policy):** DO NOT run `python 03-ai-scripts/06-cicd-local-runner.py` during routine turns.
- **NO TEST RUNNING (TOTAL BAN):** NEVER run test suites (`go test ./...`, `pytest`) during routine execution.
- **NO BUILD CHECKING (TOTAL BAN):** NEVER run build commands (`go build`, `npm run build`) after individual file edits.
- **GitMap Hyphen Commit Syntax:** All commits MUST use GitMap semantic commands with hyphen separator:
  - Feature: `gitmap cpf "<module> - <summary>"`
  - Bug Fix: `gitmap cpb "<module> - <summary>"`
  - TOTAL BAN on colons `:` inside the commit message argument.
  - TOTAL BAN on raw `git commit`, `git add -A`, or per-file commits.
- **Mandatory Completion Invariant ('No Push = Not Done'):** Tasks are only completed when changes are committed and pushed upstream to GitHub.
