# Subtask 01: Specification Authoring & Architecture Decomposition

## Status: COMPLETED
- **Subtask Slug:** `01-subtask-spec-and-planning`
- **Parent Plan:** `.ai-memory/plans/pending/328-flexible-path-parsing-and-asset-image-reading.md`
- **Architecture Spec:** `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/01-architecture-spec.md`
- **Owner:** Subagent 1 (Worker 01)
- **Target Subsystem:** `02-spec/21-app/328-*/`, `.ai-memory/plans/subtasks/328-*/`
- **Exclusive Owned Files:**
  - `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/00-master-audit-ledger.md`
  - `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/01-architecture-spec.md`
  - `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/02-component-and-cli-spec.md`
  - `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/03-migration-and-verification-spec.md`
  - `.ai-memory/plans/pending/328-flexible-path-parsing-and-asset-image-reading.md`
  - `.ai-memory/plans/subtasks/328-flexible-path-parsing-and-asset-image-reading/01-subtask-spec-and-planning.md`
  - `.ai-memory/plans/subtasks/328-flexible-path-parsing-and-asset-image-reading/02-subtask-flexible-path-parser.md`
  - `.ai-memory/plans/subtasks/328-flexible-path-parsing-and-asset-image-reading/03-subtask-safe-image-cat-and-skills-sync.md`

---

## 1. Context & Scope

This subtask establishes the complete specification suite, architectural contracts, and actionable implementation plans for Task 328: *Flexible Path Parsing and Asset Image Reading*.

In multi-agent collaborative engineering, automated agents and users paste file references across chat transcripts, markdown documents, and task lists. These inputs commonly take complex forms:
1. `@[C:\Users\Administrator\...\media_1792894009084_45b4893f.png]` (Windows absolute inside bracket wrapper)
2. `![Attached Image 1](file:///C:/Users/Administrator/.../media_1792894009084_45b4893f.png)` (Markdown image with `file:///`)
3. `1. @[assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo` (numbered list, repo-relative path, trailing comment)
4. `1. @[../repo-cache/assets/media_1792894009084_45b4893f.png] # consider the file is relative from the repo outside of the repo` (parent directory relative path)

Additionally, executing `gitmap cat` on raw binary image files historically flooded the terminal stdout stream with binary control characters, causing terminal corruption and freezing agent processes.

Subtask 01 authors the master audit ledger, architectural specifications, component interfaces, verification matrices, and subtask execution plans to govern implementation by Subagents 2 and 3.

---

## 2. Deliverables

| Deliverable Artifact | Description | Target Path |
| :--- | :--- | :--- |
| **Master Audit Ledger** | Verbatim intent, traceability matrix, risk registry, and subtask ownership mapping. | `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/00-master-audit-ledger.md` |
| **Comprehensive Architecture Spec** | Architectural blueprint covering the 4-stage path parser and safe image cat guard. | `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/01-architecture-spec.md` |
| **Component & CLI Spec** | Detailed Go types, function definitions, CLI flags, usage patterns, and terminal mockups. | `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/02-component-and-cli-spec.md` |
| **Migration & Verification Spec** | Test case matrices for all 4 path formats, image types, edge cases, and skills sync guidance. | `02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/03-migration-and-verification-spec.md` |
| **Master Pending Plan** | Tracking roadmap and subtask status registry. | `.ai-memory/plans/pending/328-flexible-path-parsing-and-asset-image-reading.md` |
| **Subtask 01 Plan** | Ownership boundaries and verification checklist for spec authoring. | `.ai-memory/plans/subtasks/328-flexible-path-parsing-and-asset-image-reading/01-subtask-spec-and-planning.md` |
| **Subtask 02 Plan** | Concrete implementation blueprint for `cli/pathparse` package. | `.ai-memory/plans/subtasks/328-flexible-path-parsing-and-asset-image-reading/02-subtask-flexible-path-parser.md` |
| **Subtask 03 Plan** | Concrete implementation blueprint for safe image cat and skills synchronization. | `.ai-memory/plans/subtasks/328-flexible-path-parsing-and-asset-image-reading/03-subtask-safe-image-cat-and-skills-sync.md` |

---

## 3. Strict Verification Criteria

All authored specifications satisfy the following project rules:
1. **Strictly Relative Paths:** Zero absolute file system paths are used in Markdown links or documentation references.
2. **Positive Boolean Naming:** All booleans exclusively use positive prefixes (`is*`, `has*`, e.g., `isImage`, `isAbsolute`, `hasRawFlag`, `hasTrailingComment`).
3. **15-Line Function Limit:** All proposed Go functions adhere to $\le 15$ lines maximum.
4. **Zero Nested Ifs:** All conditionals are structured with guard clauses and early returns.
5. **Whitespace Discipline:** Blank line before `if`, blank line after `}`, and blank line before `return`.
6. **Strict No-Build, No-Test Invariant:** Zero build or test tools are run. All validation is structural and cognitive.

---

## 4. Execution Steps

- [x] **Step 1: Ingestion & Baseline Audit**
  - Read user requirements and system instructions for Task 328.
  - Review existing path utilities in `cli/fsutil/` and command dispatcher in `cli/cmdmacro/cat.go`.
- [x] **Step 2: Author Master Audit Ledger (`00-master-audit-ledger.md`)**
  - Document verbatim user intent, 4 path format specifications, and traceability matrix.
- [x] **Step 3: Author Master Architecture Spec (`01-architecture-spec.md`)**
  - Detail 4-stage parsing pipeline, Mermaid sequence diagrams, and dual-layer image detection architecture.
- [x] **Step 4: Author Component & CLI Spec (`02-component-and-cli-spec.md`)**
  - Define Go types (`ParsedPathResult`, `ImageMetadata`), function signatures, CLI flags, and terminal outputs.
- [x] **Step 5: Author Migration & Verification Spec (`03-migration-and-verification-spec.md`)**
  - Formulate verification test cases across all 4 path formats, image types, and skills sync protocol.
- [x] **Step 6: Author Master Plan & Subtask Plans**
  - Author `.ai-memory/plans/pending/328-flexible-path-parsing-and-asset-image-reading.md`.
  - Author `01-subtask-spec-and-planning.md`.
  - Author `02-subtask-flexible-path-parser.md`.
  - Author `03-subtask-safe-image-cat-and-skills-sync.md`.
- [x] **Step 7: Verification & Notification**
  - Confirm all 8 owned files are synchronized and compliant.
  - Notify parent coordinator with structured diffstat.
