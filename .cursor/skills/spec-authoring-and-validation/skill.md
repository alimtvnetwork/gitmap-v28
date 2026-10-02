---
name: spec-authoring-and-validation
description: Author, structure, sequence, and validate repository specifications adhering to 02-spec/01-spec-authoring-guide/.
---

# Specification Authoring & Validation Guide

This skill governs the creation, organization, and automated validation of architectural specifications in `02-spec/`, primarily under the canonical application specifications folder `02-spec/21-app/`.

## RULE 00: Pure Specification Authoring Only (Total Ban on Code Execution)

> [!CAUTION]
> **THIS SKILL IS STRICTLY FOR SPECIFICATION AUTHORING AND SUBTASK PLANNING.**
> - **NEVER write, edit, compile, or delete application source code files** (`.go`, `.ts`, `.tsx`, `.py`, `.php`, `.cs`, etc.).
> - **NEVER execute unit tests, integration tests, migrations, or application builds.**
> - **NEVER run modifying shell commands against application packages.**
> - All agent tool calls are strictly restricted to:
>   1. Reading existing codebase files and specifications for architectural context.
>   2. Decoding and saving base64 screenshots to `assets/screenshots/`.
>   3. Authoring canonical product specifications directly inside `02-spec/21-app/`.
>   4. Generating execution plans and subtasks inside `.ai-memory/plans/`.
>   5. Updating `02-spec/21-app/readme.md` and `.ai-memory/plans/readme.md`.

## Pre-Planning Step 0: Task Extraction & Chat Output Gate (Mandatory First Action)

When a large prompt or complex set of requirements is given, the AI cannot understand everything at once. Therefore, before doing any deep planning, codebase searches, or spec writing, the AI MUST first break down whatever requirements the user has given (regardless of formatting) into discrete, actionable items (`Task-01`, `Task-02`).
- Respect whatever requirements the user has given, parse every request completely, and format each task clearly with proper markdown indentation, vertical blank lines, task state (`State: [IN PROGRESS — EXECUTING IMMEDIATELY]`), and an explicit understanding indicator bracket (`Understood: [YES — ...]`).
- TOTAL BAN ON UNFORMATTED RUN-ON TEXT: Never concatenate tasks into a single unformatted line or paragraph block (e.g. NEVER `#1. Task-01: ... #2. Task-02: ...`). Every task must be its own clearly separated markdown item.
- Line-by-Line Output Format Structure:
  - Line 1: Header `### 📋 Confirmed Task Breakdown & Requirement Ingestion`
  - Line 2: Empty blank line
  - Line 3: Numbered task title `1. **Task-01: [Descriptive Task Title]**`
  - Line 4: Indented state bullet (3 spaces) `   - **State:** [IN PROGRESS — EXECUTING IMMEDIATELY]`
  - Line 5: Indented understanding check (3 spaces) `   - **Understood:** [YES] — [1-2 concise sentences proving understanding of intent, scope, and verified constraints]`
  - Line 6: Indented actionable scope bullet (3 spaces) `   - **Actionable Scope:** [Precise technical deliverable and specification scope]`
  - Line 7: Indented target files bullet (3 spaces) `   - **Target Files / Area:** [02-spec/21-app/relative/path]`
  - Line 8: Empty blank line (vertical gap before next task)
  - Concluding Line: `Proceeding directly to specification authoring in 02-spec/21-app/ and subtask planning (Active Tool Call Running Below).`
- MANDATORY SAME-TURN TOOL CHAIN: You MUST output this confirmed breakdown directly in chat, but you MUST NOT end your turn or pause after emitting it! In the EXACT SAME RESPONSE turn, you MUST immediately invoke your first tool call (e.g. `write_to_file` to initialize the spec in `02-spec/21-app/`). Never emit the breakdown with text alone.

```markdown
### 📋 Confirmed Task Breakdown & Requirement Ingestion

1. **Task-01: [Descriptive Task Title]**
   - **State:** `[IN PROGRESS — EXECUTING IMMEDIATELY]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and specification scope]
   - **Target Files / Area:** `[02-spec/21-app/relative/path]`

2. **Task-02: [Descriptive Task Title]**
   - **State:** `[QUEUED — EXECUTING NOW WITHOUT USER PROMPT]`
   - **Understood:** `[YES]` — [Concise 1-sentence verification of user requirement, intent, and verified constraints]
   - **Actionable Scope:** [Precise technical deliverable and specification scope]
   - **Target Files / Area:** `[02-spec/21-app/relative/path]`

Proceeding directly to specification authoring in 02-spec/21-app/ and subtask planning (Active Tool Call Running Below).
```

MANDATORY SAME-TURN TOOL CHAIN: The breakdown text above and your first tool call MUST be emitted in the EXACT SAME TURN. Never end the turn with text alone.

---

## Canonical Specification Destination in Folder 21 (`02-spec/21-app/`)

All architectural specifications, feature definitions, data contracts, and product-level documentation MUST be authored inside `02-spec/21-app/`:
- **Focused Feature (Single File):** If the deliverable is self-contained or under 100 lines, write:
  `02-spec/21-app/xx-<slug>.md`
- **Complex Feature (Segmented Subfolder):** If the deliverable exceeds 3 subtasks or spans UI + API + database, create a dedicated sequential directory:
  `02-spec/21-app/xx-<slug>/` containing:
  - `01-overview.md` — High-level architecture, module interactions, and the mandatory `## User Request (Verbatim)` section.
  - `02-data-contracts.md` — Types, schemas, API request/response contracts, and database models.
  - `03-workflow-and-state.md` — Control flows, state machine transitions, and business validation rules.
  - `04-ui-ux-spec.md` — Visual layout, design tokens, typography, and embedded relative screenshot links (`assets/screenshots/...`).
  - `05-acceptance-criteria.md` — Testable verification rules and quality gates.
- **Mandatory Registry Update:** Register newly created specs in `02-spec/21-app/readme.md` under `## Contents`.

---

## Lossless Verbatim Ingestion (Zero-Loss Requirement Capture)

> [!IMPORTANT]
> Every specification authored in `02-spec/21-app/` MUST include a top-level section:
> `## User Request (Verbatim)`
> Copy and paste the user's complete prompt text, instructions, parameters, edge cases, and examples character-for-character, word-for-word, without alteration or omission.

---

## Screenshot & Print Screen Base64 Ingestion Protocol (Mandatory in Specs)

If the user request or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
1. Convert & Save Locally: Immediately decode the base64 encoding or download the image from the URL to the local filesystem under `assets/screenshots/<slug>-<NN>.png` or `assets/ui/<slug>-<NN>.png`.
2. Never Embed Raw Base64 or Remote URLs: Never leave raw base64 strings or ephemeral external URLs inside specification files or plans.
3. Strict Relative Path Referencing: In the master spec, domain documentation, and subtasks, refer back to the saved image file strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
4. Visual Ground Truth: Use the saved screenshot as the visual ground truth for layout, colors, component hierarchy, spacing, and state transitions during spec authoring and UI task execution.

---

## Structure & File Naming Conventions

1. Folder Naming:
   - Folders follow the hyphenated two-digit sequence pattern: `02-spec/<NN>-<slug>/` (e.g. `02-spec/02-coding-guidelines/`, `02-spec/21-app/`).

2. Mandatory Files per Spec Folder:
   - `readme.md`: Primary entry point explaining scope, version, goal, and learn checklists.
   - Numbered markdown files: Detailed topic-specific policies.
   - `97-acceptance-criteria.md`: Verification commands and criteria.
   - `98-changelog.md`: Evolution history of the specification.
   - `99-consistency-report.md`: Audit log verifying alignment with global rules.

3. Strict Relative Path Rules:
   - All internal links must use relative paths starting from the repository root or relative markdown paths.
   - Never write absolute filesystem paths or `file:///` URIs.

4. Actionable Subtask Generation in `.ai-memory/plans/` (Decoupled Planning):
   - Parent Plan: `.ai-memory/plans/pending/xx-<slug>.md`
   - Subtasks: `.ai-memory/plans/subtasks/xx-<slug>/001-<task>.md`, `002-<task>.md`, etc.
   - **MANDATORY SPEC LINK:** Every subtask MUST include a prominent relative markdown link back to the canonical spec in `02-spec/21-app/`:
     ```markdown
     Spec Reference: [02-spec/21-app/xx-<slug>.md](../../../02-spec/21-app/xx-<slug>.md)
     ```
   - **Lean Subtasks Mandate:** Subtasks in `.ai-memory/plans/subtasks/<plan-slug>/` MUST NOT repeat common repository boilerplate, universal coding rules, banned operations, or generic guidelines. Universal rules exist in root guidelines and the canonical spec. Subtasks must contain strictly the unique, task-specific details, exact file paths, symbol modifications, and runnable verification checks.
   - Register plan in `.ai-memory/plans/readme.md`.

5. Validation Checklist:
   - Run spec cross-link validation:
     ```bash
     python linter-scripts/check-spec-cross-links.py --root 02-spec --repo-root .
     ```
   - Ensure header spacing and markdown gap linters pass.

---

## End-of-Run Comprehensive Summary & Traceability Report (Mandatory Output)

At the conclusion of the specification turn, you MUST output the following structured summary in chat:

```markdown
### 📑 Specifications Created in Folder 21 (`02-spec/21-app/`)

- [02-spec/21-app/xx-<slug>.md](02-spec/21-app/xx-<slug>.md) (Lines: <count>) — [Short description of scope]
  *(or list files in 02-spec/21-app/xx-<slug>/ if segmented)*

### 📋 Actionable Tasks & Subtasks Created (`.ai-memory/plans/`)

- Parent Plan: [.ai-memory/plans/pending/xx-<slug>.md](.ai-memory/plans/pending/xx-<slug>.md)
- Subtask Count: <count> subtasks in `.ai-memory/plans/subtasks/xx-<slug>/`
  1. `001-<task>.md`: [Title] -> Target: `[relative/path]` | State: `[QUEUED]`
  2. `002-<task>.md`: [Title] -> Target: `[relative/path]` | State: `[QUEUED]`

### 🔗 Requirements Traceability Matrix

| Requirement / Prompt Item | Canonical Spec File (`02-spec/21-app/`) | Subtask File (`.ai-memory/plans/subtasks/`) | Target Code Files |
|:---|:---|:---|:---|
| [User Requirement 1] | `02-spec/21-app/xx-<slug>.md` | `.ai-memory/plans/subtasks/xx-<slug>/001-<task>.md` | `path/to/file.go` |
| [User Requirement 2] | `02-spec/21-app/xx-<slug>.md` | `.ai-memory/plans/subtasks/xx-<slug>/002-<task>.md` | `path/to/file.ts` |

### 🚀 Next Steps: Execution Command

To execute these generated subtasks in continuous sequence with 2-agent concurrency:
> Run: `06-execute-parent-task-with-n-steps-v2.md` with plan slug `xx-<slug>`
```
