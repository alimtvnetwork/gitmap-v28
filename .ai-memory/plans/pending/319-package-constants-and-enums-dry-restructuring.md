# Master Plan: 319 — Package Constants & Enums DRY Restructuring

> **Status:** Pending Execution (Ready for 5-Package Rollout Turn)  
> **Sequence ID:** 319  
> **Slug:** `package-constants-and-enums-dry-restructuring`  
> **Canonical Spec:** `02-spec/21-app/319-package-constants-and-enums-dry-restructuring/`  
> **Subtask Root:** `.ai-memory/plans/subtasks/319-package-constants-and-enums-dry-restructuring/`  

---

## 1. User Request (Verbatim)

```text
@[01-prompts/15-cg-execute/06-constants-and-enums.md] first make a plan for each 4 packages with a pkg/enums/<enum>, pkg/constants/constants.go and other go files to organize constants to make the code DRY and next when I share @[d:\work\gitmap\01-prompts\15-cg-execute\06-constants-and-enums.md] 

You should apply on the 5 packages, clear???
```

---

## 2. 5-Pillar Architectural Mandate

### Pillar 1: Strategic Intent & Systems Reasoning
Provide a comprehensive architecture and execution blueprint for eliminating raw string literals, magic numbers, untyped strings, and duplicate literals across GitMap's core Go packages.
1. **Strongly-Typed Enum Types:** Group related values into dedicated `enums/` packages using the mandatory `*Type` suffix (e.g., `PullModeType`, `CloneProtocolType`, `StoreTableType`, `CommitVerbType`), complete with `IsValid() bool` and `String() string` methods.
2. **Centralized Constants:** Organize all command flags, telemetry tags, table names, format strings, and process environment variables into dedicated `constants/constants.go` files (kept $\le 300$ lines).
3. **Decoupled Architecture:** Maintain strict acyclic package dependencies. Enums and constants reside in leaf packages to ensure zero circular imports.
4. **Phased Rollout:** Phase 1 (current turn) creates the complete architecture, discovery audits, and subtask ledger. Phase 2 & 3 executes the 5-package refactor immediately when the user shares `01-prompts/15-cg-execute/06-constants-and-enums.md`.

### Pillar 2: Explicit Scope Boundaries
- **In-Scope (Current Turn):**
  - Canonical spec in `02-spec/21-app/319-package-constants-and-enums-dry-restructuring/` (01 to 03 + readme).
  - Parent plan in `.ai-memory/plans/pending/319-package-constants-and-enums-dry-restructuring.md`.
  - 5 subtask plans in `.ai-memory/plans/subtasks/319-package-constants-and-enums-dry-restructuring/`.
  - Discovery audits across `cli/cmdpull`, `cli/cloner`, `cli/cmdclone`, `cli/store`, and `cli/cmdcommit`.
  - Registry synchronization in `.ai-memory/plans/readme.md` and `02-spec/21-app/readme.md`.
- **In-Scope (Next Turn Rollout across 5 Packages):**
  - Package 1: `cli/cmdpull` (`PullModeType`, `PullStatusType`, flags, summary tokens).
  - Package 2: `cli/cloner` (`CloneProtocolType`, `CloneDepthType`, git args, hooks).
  - Package 3: `cli/cmdclone` (`CloneStrategyType`, `CloneProfileType`, telemetry, prompts).
  - Package 4: `cli/store` (`StoreTableType`, `DbOpType`, pragmas, column names).
  - Package 5: `cli/cmdcommit` (`CommitVerbType`, `CommitFormatType`, module templates).
- **Excluded / Deferred:**
  - Production code edits to unrelated packages (`cli/cmdlemon`, `cli/cmdssh`, etc.).
  - Running full local test suites or builds (FastGate AST and targeted linters only).

### Pillar 3: Non-Negotiable Constraints & Negative Guardrails
- **Function Limit:** Maximum 15 lines per function; $\le 8$ lines preferred.
- **File Limit:** Maximum 300 lines per `.go` file.
- **Positive Booleans Only:** Booleans must begin with `is*` or `has*`. Zero negative booleans.
- **Zero Nested Ifs:** All branches flattened with early returns and guard clauses.
- **Return Formatting:** Blank line after `}` and blank line before `return`.
- **Search Exclusivity:** Exclusively use GitMap search tools (`gitmap aum search`, `gitmap find`); ban on raw `rg`, `ripgrep`, `grep`, `Select-String`.
- **Git Safety:** Zero git commands in subagents. Lead commits with `gitmap cpf "<module> - <summary>"`.

### Pillar 4: Blast Radius & Dependency Topology
- **Isolated Leaf Packages:** Enums live in subpackages that import only standard library packages (`fmt`, `strings`).
- **Zero Cycle Invariant:** Parent packages import `enums` and `constants`; `enums` and `constants` NEVER import the parent package.
- **Backwards Compatibility:** Exported legacy identifiers in `cli/constants/` remain as type aliases during transitional phases.

### Pillar 5: Multi-Agent Concurrency Topology & Wave Orchestration
- **Swarm Capacity:** $A = 3$ concurrent subagents $\times$ $H = 2$ operational hands $= 6$ concurrent operations.
- **Planning Wave (Current):** 3 research subagents auditing packages in parallel.
- **Execution Waves (Next Turn):** Wave 1 refactors `cmdpull`, `cloner`, `cmdclone`, `store`; Wave 2 refactors `cmdcommit` and runs AST verification.

---

## 3. Subtask Catalog & Execution Status

| Subtask ID | Target Package | Status | Assigned Primary Owner | Disjoint Boundaries |
| :--- | :--- | :--- | :--- | :--- |
| `Subtask-01` | `cli/cmdpull` | Ready for Execution | Worker 1 | `cli/cmdpull/enums/*`, `cli/cmdpull/constants/*` |
| `Subtask-02` | `cli/cloner` | Ready for Execution | Worker 2 | `cli/cloner/enums/*`, `cli/cloner/constants/*` |
| `Subtask-03` | `cli/cmdclone` | Ready for Execution | Worker 2 | `cli/cmdclone/enums/*`, `cli/cmdclone/constants/*` |
| `Subtask-04` | `cli/store` | Ready for Execution | Worker 3 | `cli/store/enums/*`, `cli/store/constants/*` |
| `Subtask-05` | `cli/cmdcommit` (Rollout) | Ready for Execution | Worker 1 (Wave 2) | `cli/cmdcommit/enums/*`, `cli/cmdcommit/constants/*` |

---

## 4. Verification Evidence Contract
- FastGate AST validation (`python 03-ai-scripts/50-fastgate.py`) returns exit code 0.
- Enum coding guideline linter (`python linter-scripts/check-enum-guidelines.py`) returns exit code 0.
- Forbidden strings scanner (`python linter-scripts/check-forbidden-strings.py`) returns exit code 0.
- `gofmt -w` ensures 100% formatted Go files across all 5 packages.
