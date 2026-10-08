# Subtask 01: Architecture Specification, Master Ledger & Subtask Planning Scaffolding

- **Parent Task:** `240-mcp-server-ai-analysis-agm-nodes-and-git-cache`
- **Subtask ID:** `01-spec-and-plan`
- **Spec Reference:** `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md`
- **Status:** `IN_PROGRESS`
- **Assigned Subagent:** Spec Subagent 01 (Phase 1)

---

## 1. Objective & Scope

Author the canonical architecture specification, master execution ledger, and structured subtask execution plans for the **GitMap AI MCP Server, AI Analysis Engine, Split SQLite Storage & Zero-Loss File Removal Subsystem**.

This subtask lays the architectural foundation by formalizing:
1. **Core Invariants:**
   - Invariant 1: Split SQLite Isolation (`ai-analysis.db` under `.gitmap/data/ai-analysis/`).
   - Invariant 2: Zero-Loss File Removal (`<temp>/gitmap/removed/<taskId>/` staging before unlinking).
   - Invariant 3: Strict Relative Path Hygiene across all database tables, JSON payloads, and CLI output.
   - Invariant 4: Minimal Root Help Invariant (bare `gitmap` outputs only version block + location triplet + short footer; full catalog reserved for `gitmap help` / `-h`).
   - Invariant 5: Deterministic AI Reasoning Traceability (every task and line mutation records architectural reason).
2. **Database Schemas:**
   - DDL specifications and Go models for `AiTask`, `AiTaskFile`, and `AiTaskLine` adhering to PascalCase naming, primary key `{Table}Id`, positive boolean flags (`IsActive`, `HasFailed`, `IsModified`, `IsRemoved`, `IsApplied`), and standard documentation columns (`Description`, `Notes`, `Comments`).
3. **Subtask Decomposition:**
   - Scaffolding subtask plans 01 through 05 with explicit file targets, acceptance criteria, and verification commands.

---

## 2. Concrete Files Created & Managed

| File | Purpose |
| :--- | :--- |
| `02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md` | Authoritative architecture spec defining system topology, invariants, DDL, Go models, and CLI contracts. |
| `.ai-memory/plans/subtasks/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-spec-and-plan.md` | Subtask 01 plan: Spec authoring, invariants, and planning scaffolding. |
| `.ai-memory/plans/subtasks/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/02-compact-root-help.md` | Subtask 02 plan: Minimal Root Help Invariant implementation. |
| `.ai-memory/plans/subtasks/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/03-ai-analysis-split-db.md` | Subtask 03 plan: Split SQLite storage engine (`ai-analysis.db`) & CRUD layer. |
| `.ai-memory/plans/subtasks/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/04-task-file-removal.md` | Subtask 04 plan: Zero-loss file removal vault and SHA-256 revert flow. |
| `.ai-memory/plans/subtasks/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/05-llm-train-clear-transfer.md` | Subtask 05 plan: LLM training dataset generation, pruning, and cross-system export/import. |

---

## 3. Acceptance Criteria

- [x] Architecture specification `01-architecture-spec.md` is fully written with:
  - Executive Summary & Background on GitMap as an AI MCP Server.
  - The 5 Non-Negotiable Invariants.
  - SQLite DDL schemas and Go struct definitions for `AiTask`, `AiTaskFile`, and `AiTaskLine`.
  - Zero-Loss file removal staging and revert workflows.
  - LLM training dataset export (`gitmap ai-analysis llm-train`), pruning (`gitmap ai-analysis clear`), and cross-system transfer flows.
  - Mermaid and ASCII diagrams detailing system topology, sequence workflows, and architecture.
- [x] All 5 subtask plans (`01` through `05`) authored with concrete file assignments, acceptance criteria, and validation commands.
- [x] Strict Relative Path Hygiene maintained across all spec and plan documents.
- [x] Zero git commands executed and zero unintended files modified.

---

## 4. Verification Commands

```powershell
# 1. Verify existence of specification file
Test-Path "02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md"

# 2. Verify existence of all subtask plan files
Get-ChildItem ".ai-memory/plans/subtasks/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/*.md" | Select-Object -ExpandProperty Name
```
