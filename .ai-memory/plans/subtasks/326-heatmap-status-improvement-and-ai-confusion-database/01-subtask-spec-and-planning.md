# Subtask 01: Specification Authoring & Architecture Decomposition

## Status: IN_PROGRESS
- **Subtask Slug:** `01-subtask-spec-and-planning`
- **Parent Plan:** `.ai-memory/plans/pending/326-heatmap-status-improvement-and-ai-confusion-database.md`
- **Architecture Spec:** `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/01-architecture-spec.md`
- **Owner:** Spec Author 01 (Worker 01)
- **Target Subsystem:** `02-spec/21-app/326-*/`, `.ai-memory/plans/subtasks/326-*/`
- **Exclusive Owned Files:**
  - `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/01-architecture-spec.md`
  - `.ai-memory/plans/subtasks/326-heatmap-status-improvement-and-ai-confusion-database/01-subtask-spec-and-planning.md`
  - `.ai-memory/plans/subtasks/326-heatmap-status-improvement-and-ai-confusion-database/02-subtask-heatmap-status-improvement.md`

---

## 1. Context & Scope

This subtask establishes the architectural, structural, and execution baseline for Task 326: *Heatmap Status Improvement and AI Confusion Database*.

During previous autonomous runs, developers and AI swarms encountered usability and operational friction:
1. `gitmap heatmap status` presented unformatted terminal output lacking task context and machine-readable JSON flags.
2. `gitmap fix all` triggered invalid `git commit` commands when run from parent directories outside individual git roots.
3. Automated stashing during pull operations caused developer alarm due to scary diagnostic walls without explicit file flags.
4. AI cognitive confusions, questions, and ambiguous diffs were lost upon session termination without persistent storage.
5. Inconsistent commit prefixes (`feature:`, `[FIX]`, `FIX:`) infiltrated swarm commits, violating GitMap commit standards.

Subtask 01 authors the foundational specifications and subtask plans required to coordinate parallel worker execution across these domains.

---

## 2. Deliverables

| Deliverable Artifact | Description | Target Path |
| :--- | :--- | :--- |
| **Comprehensive Architecture Spec** | Definitive architectural blueprint covering all 7 subsystems, data flows, schemas, and UI layouts. | `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/01-architecture-spec.md` |
| **Subtask 01 Specification** | Scope, boundaries, verification criteria, and execution tracking for specification authoring. | `.ai-memory/plans/subtasks/326-heatmap-status-improvement-and-ai-confusion-database/01-subtask-spec-and-planning.md` |
| **Subtask 02 Technical Plan** | Implementation blueprint for `cli/heatmap/` (Lipgloss TokyoNight styling, task integration, `--json`, tree examples). | `.ai-memory/plans/subtasks/326-heatmap-status-improvement-and-ai-confusion-database/02-subtask-heatmap-status-improvement.md` |

---

## 3. Strict Verification Criteria

To ensure architectural integrity, all authored specifications must pass these non-negotiable checks:

1. **Strictly Relative Paths:** Every Markdown link, code reference, and file pointer must be relative to the repository root (e.g., `02-spec/...`, `cli/...`). Zero absolute paths (`C:\...`, `D:\...`, `file:///...`) are permitted.
2. **Positive Boolean Naming:** All schema fields, data structures, and condition references must use positive prefixes (`is*`, `has*`, e.g., `isCompleted`, `isResolved`, `hasMoreTasks`). Banned prefixes: `un*`, `non*`, `not*`, `dis*`, `can*`, `should*`.
3. **Database Convention Compliance:** All SQL tables must adhere to PascalCase names (e.g., `AiConfusion`, `HeatmapParentTask`), camelCase column names (`aiConfusionId`, `taskSlug`), and strongly typed integer enums (`02-spec/04-database-conventions/`).
4. **Strict No-Build, No-Test Invariant:** During specification authoring, zero build tools (`go build`, `cargo build`) and zero test runners (`go test`, `pytest`) may be executed. Verification is purely structural, syntactic, and cognitive.
5. **TokyoNight Design System Alignment:** Visual CLI specifications must explicitly bind to tokens from `cli/termout/theme_lipgloss.go` and `cli/termout/themes.go` (`TokyoPrimary`, `TokyoSurface`, `TokyoOK`, `TokyoWarn`, `TokyoError`, `TokyoMuted`).

---

## 4. Boundaries & Constraints

- **Documentation Sovereignty:** Subtask 01 owns documentation files only. It must NOT modify any Go source files in `cli/`, shell scripts, or agent prompts.
- **Worker File Isolation:** Subtask 01 writes exclusively to its assigned paths in `02-spec/21-app/326-*/` and `.ai-memory/plans/subtasks/326-*/`. Files assigned to Worker 02 (`cli/cmdfix/`, `cli/store/`) and Worker 03 (`.agents/skills/`, `cli/cmdcommit/`) must not be touched.
- **Traceability Guarantee:** Every requirement articulated in the user request must map 1:1 to an implementation subtask and acceptance criterion in the master audit ledger.

---

## 5. Execution Steps

- [x] **Step 1: Context Ingestion & Baseline Audit**
  - Read `00-master-audit-ledger.md` and user verbatim requirements.
  - Inspect `cli/heatmap/`, `cli/cmdplans/`, and `cli/termout/` to assess current implementations.
- [x] **Step 2: Author Master Architecture Specification**
  - Draft `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/01-architecture-spec.md`.
  - Document all 11 core sections including TokyoNight styling, task integration, `--json` schema, work directory guards, humanized stashing, `repo-cache` AI Confusion database, and semantic commit standards.
- [x] **Step 3: Author Subtask 01 Planning Document**
  - Draft `.ai-memory/plans/subtasks/326-heatmap-status-improvement-and-ai-confusion-database/01-subtask-spec-and-planning.md`.
  - Define scope, deliverables, verification criteria, and boundaries.
- [ ] **Step 4: Author Subtask 02 Technical Implementation Plan**
  - Draft `.ai-memory/plans/subtasks/326-heatmap-status-improvement-and-ai-confusion-database/02-subtask-heatmap-status-improvement.md`.
  - Specify micro-batching, data structures, and refactoring strategy for `cli/heatmap/`.
- [ ] **Step 5: Hand-off & Progress Notification**
  - Send structured completion message to parent coordinator.

---

## 6. Acceptance Criteria

- [x] `01-architecture-spec.md` exists and details all 7 core functional requirements.
- [x] `01-subtask-spec-and-planning.md` exists with clear boundaries and positive boolean standards.
- [x] `02-subtask-heatmap-status-improvement.md` exists with concrete refactoring plans for `cli/heatmap/`.
- [x] All file paths are strictly relative to repository root.
- [x] Zero compilation or test commands executed during this turn.
