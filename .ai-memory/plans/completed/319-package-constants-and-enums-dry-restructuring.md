# Master Plan: 319 — Package Constants & Enums DRY Restructuring

> **Status:** Completed (100% Verified Across 5 Packages)  
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
Successfully executed package-level constants and enum deduplication across 5 Go packages:
1. **Strongly-Typed Enum Types:** Grouped related values into dedicated `enums/` packages using the mandatory `*Type` suffix (`PullStepType`, `PullOutcomeStatusType`, `ConflictSideType`, `AuditActionType`, `GitUnmergedStatusType`, `CloneDepthType`, `AuthMethodType`, `CloneStrategyType`, `BatchStatusType`, `StoreTableType`, `SQLitePragmaType`, `SplitDatabaseType`, `CommitVerbType`, `ConventionalVerbType`, `CommitFormatType`), complete with `IsValid() bool` and `String() string` methods.
2. **Centralized Constants:** Organized all command flags, table headers, tree glyphs, static git argument slices, and process environment variables into dedicated `constants/constants.go` files (all kept $\le 300$ lines).
3. **Decoupled Architecture:** Strict acyclic package dependencies maintained. `enums` and `constants` reside in leaf packages importing only standard library packages.
4. **Phased Rollout:** Phase 1 authored specifications and plans; Phase 2 & 3 executed complete 5-package implementation, refactoring, and FastGate AST verification.

### Pillar 2: Explicit Scope Boundaries
- **Delivered Packages (5 packages):**
  - Package 1: `cli/cmdpull` (`PullStepType`, `PullOutcomeStatusType`, `ConflictSideType`, flags, tree glyphs, table headers).
  - Package 2: `cli/cloner` (`AuditActionType`, `GitUnmergedStatusType`, `CloneDepthType`, immutable static CLI argument slices).
  - Package 3: `cli/cmdclone` (`AuthMethodType`, `CloneStrategyType`, `BatchStatusType`, clone flags, banners).
  - Package 4: `cli/store` (`StoreTableType`, `SQLitePragmaType`, `SplitDatabaseType`, default pragma runner).
  - Package 5: `cli/cmdcommit` (`CommitVerbType`, `ConventionalVerbType`, `CommitFormatType`, flags, format patterns).

### Pillar 3: Non-Negotiable Constraints & Negative Guardrails
- **Function Limit:** Maximum 15 lines per function; $\le 8$ lines preferred.
- **File Limit:** Maximum 300 lines per `.go` file.
- **Positive Booleans Only:** Booleans must begin with `is*` or `has*`. Zero negative booleans.
- **Zero Nested Ifs:** All branches flattened with early returns and guard clauses.
- **Return Formatting:** Blank line after `}` and blank line before `return`.
- **Git Safety:** Zero git commands in subagents. Lead commits with `gitmap cpf "<module> - <summary>"`.

### Pillar 4: Blast Radius & Dependency Topology
- **Isolated Leaf Packages:** Subpackages only import standard library packages (`fmt`, `strings`).
- **Zero Cycle Invariant:** Parent packages import `enums` and `constants`; `enums` and `constants` NEVER import the parent package.
- **Backwards Compatibility:** Exported legacy identifiers in `cli/constants/` and package roots remain as type aliases.

### Pillar 5: Multi-Agent Concurrency Topology & Wave Orchestration
- **Swarm Execution:** $A = 3$ concurrent self subagents with disjoint package ownership:
  - Worker 1: `cli/cmdpull` (Enums + Constants + Call Sites)
  - Worker 2: `cli/cloner` and `cli/cmdclone` (Enums + Constants + Call Sites + Error Sovereignty)
  - Worker 3: `cli/store` and `cli/cmdcommit` (Enums + Constants + Call Sites + Format Helpers)

---

## 3. Subtask Catalog & Execution Status

| Subtask ID | Target Package | Status | Assigned Primary Owner | Disjoint Boundaries |
| :--- | :--- | :--- | :--- | :--- |
| `Subtask-01` | `cli/cmdpull` | Complete (Verified) | Worker 1 | `cli/cmdpull/enums/*`, `cli/cmdpull/constants/*` |
| `Subtask-02` | `cli/cloner` | Complete (Verified) | Worker 2 | `cli/cloner/enums/*`, `cli/cloner/constants/*` |
| `Subtask-03` | `cli/cmdclone` | Complete (Verified) | Worker 2 | `cli/cmdclone/enums/*`, `cli/cmdclone/constants/*` |
| `Subtask-04` | `cli/store` | Complete (Verified) | Worker 3 | `cli/store/enums/*`, `cli/store/constants/*` |
| `Subtask-05` | `cli/cmdcommit` | Complete (Verified) | Worker 3 | `cli/cmdcommit/enums/*`, `cli/cmdcommit/constants/*` |

---

## 4. Verification Evidence Contract
- FastGate AST validation (`python 03-ai-scripts/50-fastgate.py`) passed in 0.044s with 0 violations.
- Enum coding guideline linter (`python linter-scripts/check-enum-guidelines.py`) passed with exit code 0.
- All unit tests across all 5 packages and 10 new subpackages passed 100% green (`go test -count=1 ./...`).
- `gofmt -w` ensured 100% formatted Go files across all 5 packages.
