# Subtask 01: Specification Authoring & Architecture Decomposition

## Status: COMPLETED
- **Subtask Slug:** `01-subtask-spec-and-planning`
- **Parent Plan:** `.ai-memory/plans/pending/327-heatmap-status-improvement-and-ai-confusion-handling.md`
- **Architecture Spec:** `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/01-architecture-spec.md`
- **Owner:** Subagent 1 (Worker 01)
- **Target Subsystem:** `02-spec/21-app/327-*/`, `.ai-memory/plans/subtasks/327-*/`, `cli/store/`
- **Exclusive Owned Files:**
  - `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/01-architecture-spec.md`
  - `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/01-subtask-spec-and-planning.md`
  - `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/02-subtask-sequence-folder-caching-and-1-to-n-schema.md`
  - `cli/store/ai_confusion_types.go`
  - `cli/store/ai_confusion_schema.go`

---

## 1. Context & Scope

This subtask establishes the comprehensive architectural and structural baseline for Task 327: *Heatmap Status Improvement and AI Confusion Handling*.

During sustained multi-agent swarm operations, seven architectural deficiencies were identified:
1. **O(N) Directory Traversal for Repo-Cache:** GitMap repeatedly scanned the filesystem to resolve `repo-cache` sequence directories (`XX-<repo>`), degrading CLI execution speed. A root SQLite `RepoCacheSequenceIndex` cache is required.
2. **Denormalized Confusion Files:** The AI confusion storage lacked a 1-to-N normalized table structure for referenced source files, line spans, why modified explanations, and code snippets.
3. **Multi-File Context Payloads:** Both terminal and JSON ingestion payloads needed extension to accept multiple files (`referencedFiles []AiConfusionFileInput`) and MD paths.
4. **SQLite Merging Hazards:** Binary SQLite files in `repo-cache` cannot merge via git; an autonomous export-revert-pull-replay synchronization flow is required on `gitmap rc ai-confusion status`.
5. **Relative File Context Resolution:** Agents need to view past cognitive confusion traces tied to specific relative file paths (`gitmap rc ai-confusion view <file>`).
6. **Remediation List Contamination:** Failed repository remediation outputs occasionally treated commit message strings as repository names.
7. **Minor Version Release Governance:** Packaging and releasing `v6.611.0` with semantic commit standards.

Subtask 01 authors the foundational architecture specification and technical implementation plans, and establishes the foundational 1-to-N Go data structures and SQLite DDL in `cli/store/`.

---

## 2. Deliverables

| Deliverable Artifact | Description | Target Path |
| :--- | :--- | :--- |
| **Comprehensive Architecture Spec** | Definitive architectural blueprint covering all 7 subsystems, data flows, schemas, and UI layouts. | `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/01-architecture-spec.md` |
| **Subtask 01 Specification** | Scope, boundaries, verification criteria, and execution tracking for specification authoring. | `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/01-subtask-spec-and-planning.md` |
| **Subtask 02 Technical Plan** | Implementation blueprint for root SQLite sequence folder cache (`RepoCacheSequenceIndex`) and 1-to-N `AiConfusionFile` schema & store. | `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/02-subtask-sequence-folder-caching-and-1-to-n-schema.md` |
| **Store Type Models** | `AiConfusionFileRecord`, `AiConfusionRecord.ReferencedFiles`, `AiConfusionFileInput`, and JSON mapping helpers. | `cli/store/ai_confusion_types.go` |
| **Store Schema DDL** | DDL and indexes for `AiConfusionFile` table with cascading foreign keys in `InitAiConfusionSchema`. | `cli/store/ai_confusion_schema.go` |

---

## 3. Strict Verification Criteria

To ensure architectural integrity, all authored specifications and code modifications must satisfy:

1. **Strictly Relative Paths:** Every Markdown link, code reference, and file pointer must be relative to the repository root (e.g., `02-spec/...`, `cli/...`). Zero absolute paths (`C:\...`, `D:\...`, `file:///...`) are permitted.
2. **Positive Boolean Naming:** All schema fields, data structures, and condition references must use positive prefixes (`is*`, `has*`, e.g., `isCompleted`, `isFailure`, `hasFiles`). Banned prefixes: `un*`, `non*`, `not*`, `dis*`, `can*`, `should*`.
3. **Database Convention Compliance:** All SQL tables must adhere to PascalCase names (`AiConfusion`, `AiConfusionFile`, `RepoCacheSequenceIndex`), camelCase column names (`aiConfusionId`, `filePath`), and strongly typed integer enums.
4. **Function Line Limits:** Every Go function must be $\le 15$ lines maximum ($\le 8$ lines preferred).
5. **Zero Nested If Statements:** All conditionals must be flattened with early returns and guard clauses.
6. **Whitespace Discipline:** Blank line before `if`, blank line after `}`, and blank line before `return`.
7. **Strict No-Build, No-Test Invariant:** During specification and type authoring, zero build tools (`go build`, `cargo build`) and zero test runners (`go test`, `pytest`) may be executed. Verification is purely structural, syntactic, and cognitive.

---

## 4. Boundaries & Constraints

- **File Ownership Isolation:** Subagent 1 writes exclusively to its assigned 5 files:
  - `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/01-architecture-spec.md`
  - `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/01-subtask-spec-and-planning.md`
  - `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/02-subtask-sequence-folder-caching-and-1-to-n-schema.md`
  - `cli/store/ai_confusion_types.go`
  - `cli/store/ai_confusion_schema.go`
- Files assigned to peer subagents (`cli/store/ai_confusion_store.go`, `cli/store/repo_cache_sequence_store.go`, `cli/cmd/rc_ai_confusion.go`, `cli/cmdpull/`, `version.json`, etc.) must not be modified by Subagent 1.
- **Traceability Guarantee:** Every requirement in the user request is mapped 1:1 to an implementation subtask and acceptance criterion in the master audit ledger.

---

## 5. Execution Steps

- [x] **Step 1: Ingestion & Baseline Audit**
  - Read `00-master-audit-ledger.md` and user verbatim requirements for Task 327.
  - Inspect `cli/store/ai_confusion_types.go` and `cli/store/ai_confusion_schema.go`.
- [x] **Step 2: Update Store Types (`cli/store/ai_confusion_types.go`)**
  - Define `AiConfusionFileRecord` struct with all 9 fields.
  - Add `ReferencedFiles []AiConfusionFileRecord` to `AiConfusionRecord`.
  - Define `AiConfusionFileInput` struct with `ToRecord(uuid string)` helper.
  - Update `AiConfusionJSONInput` with `ReferencedFiles []AiConfusionFileInput`.
  - Implement decomposed conversion helpers strictly $\le 15$ lines each.
- [x] **Step 3: Update Store Schema DDL (`cli/store/ai_confusion_schema.go`)**
  - Define `SQLCreateAiConfusionFileTable` with `FOREIGN KEY ... ON DELETE CASCADE`.
  - Define `SQLCreateIdxAiConfusionFileConfusionId` and `SQLCreateIdxAiConfusionFileFilePath`.
  - Add statements to `aiConfusionSchemaStatements` executed by `InitAiConfusionSchema`.
  - Verify vertical whitespace rules (blank line before `if`, blank line after `}`).
- [x] **Step 4: Author Master Architecture Specification**
  - Author `02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/01-architecture-spec.md` covering all 7 subsystems, Mermaid diagrams, DDLs, and traceability matrix.
- [x] **Step 5: Author Subtask 01 Planning Document**
  - Author `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/01-subtask-spec-and-planning.md`.
- [x] **Step 6: Author Subtask 02 Technical Implementation Plan**
  - Author `.ai-memory/plans/subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/02-subtask-sequence-folder-caching-and-1-to-n-schema.md`.
- [x] **Step 7: Verification & Notification**
  - Confirm all 5 owned files are synchronized and compliant.
  - Notify parent coordinator with structured diffstat.

---

## 6. Acceptance Criteria

- [x] `01-architecture-spec.md` exists and details all 7 core functional requirements.
- [x] `01-subtask-spec-and-planning.md` exists with clear boundaries and positive boolean standards.
- [x] `02-subtask-sequence-folder-caching-and-1-to-n-schema.md` exists with concrete implementation details.
- [x] `cli/store/ai_confusion_types.go` contains `AiConfusionFileRecord`, `ReferencedFiles`, and input adapters.
- [x] `cli/store/ai_confusion_schema.go` contains table and index DDL executed by `InitAiConfusionSchema`.
- [x] All file paths are strictly relative to repository root.
- [x] Zero compilation or test commands executed during this turn.
