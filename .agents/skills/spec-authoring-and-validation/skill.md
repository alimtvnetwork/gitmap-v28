---
name: spec-authoring-and-validation
description: Author, structure, sequence, and validate repository specifications adhering to 02-spec/01-spec-authoring-guide/.
---

# Specification Authoring & Validation Guide

This skill governs the creation, organization, and automated validation of architectural specifications in `02-spec/`.

## Pre-Planning Step 0: Task Extraction & Chat Confirmation Gate (Mandatory First Action)

When a large prompt or complex set of requirements is given, the AI cannot understand everything at once. Therefore, before doing any deep planning, codebase searches, or spec writing, the AI MUST first break down the requirements into smaller tasks and a discrete list of ordered items (`#1. Task-01:`, `#2. Task-02:`) without hard brackets, and output this confirmed task breakdown directly in chat:

```markdown
### Confirmed Task Breakdown
#1. Task-01: [Actionable deliverable description]
#2. Task-02: [Actionable deliverable description]
Proceeding directly to detailed specification and subtask planning.
```

After outputting this confirmed breakdown in chat, provide each task into the spec in a very detailed manner:
- Architectural context, domain logic, and module interactions.
- Input and output data contracts.
- Exact symbol signatures and target files.
- Visual specification references: If screenshot URLs or base64 data URIs are provided, decode and save them locally under `assets/screenshots/<slug>-<NN>.png` and refer back to them via relative paths (e.g. `![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
- Acceptance criteria and verification checks.

## Screenshot & Print Screen Base64 Ingestion Protocol (Mandatory in Specs)

If the user request or prompt contains a screenshot URL, print screen link, or base64 data URI (e.g. `data:image/png;base64,...`):
1. Convert & Save Locally: Immediately decode the base64 encoding or download the image from the URL to the local filesystem under `assets/screenshots/<slug>-<NN>.png` or `assets/ui/<slug>-<NN>.png`.
2. Never Embed Raw Base64 or Remote URLs: Never leave raw base64 strings or ephemeral external URLs inside specification files or plans.
3. Strict Relative Path Referencing: In the master spec, domain documentation, and subtasks, refer back to the saved image file strictly as a relative markdown link (e.g. `![Screenshot](assets/screenshots/<slug>-<NN>.png)`).
4. Visual Ground Truth: Use the saved screenshot as the visual ground truth for layout, colors, component hierarchy, spacing, and state transitions during spec authoring and UI task execution.

## Structure & File Naming Conventions

1. Folder Naming:
   - Folders follow the hyphenated two-digit sequence pattern: `02-spec/<NN>-<slug>/` (e.g. `02-spec/02-coding-guidelines/`, `02-spec/21-app/`).

2. Mandatory Files per Spec Folder:
   - `01-index.md`: Primary entry point explaining scope, version, goal, and learn checklists.
   - Numbered markdown files: Detailed topic-specific policies.
   - `97-acceptance-criteria.md`: Verification commands and criteria.
   - `98-changelog.md`: Evolution history of the specification.
   - `99-consistency-report.md`: Audit log verifying alignment with global rules.

3. Strict Relative Path Rules:
   - All internal links must use relative paths starting from the repository root or relative markdown paths.
   - Never write absolute filesystem paths or `file:///` URIs.

4. Lean Subtasks Mandate (No Common Boilerplate in Subtasks):
   - Subtasks in `.ai-memory/plans/subtasks/<plan-slug>/` MUST NOT repeat common repository boilerplate, universal coding rules, banned operations, or generic guidelines.
   - Universal rules exist in root guidelines and parent spec. Subtasks must contain strictly the unique, task-specific details, exact file paths, symbol modifications, and runnable verification checks.

5. Validation Checklist:
   - Run spec cross-link validation:
     ```bash
     python linter-scripts/check-spec-cross-links.py --root 02-spec --repo-root .
     ```
   - Ensure header spacing and markdown gap linters pass.

## End-of-Turn Verification & Confidence Reporting (Mandatory Output)

When all tasks or planning phases are completed (or if the run concludes), you MUST output:

```markdown
### Task Completion Summary
✅ #1. Task-01: [Task description] — Completed
✅ #2. Task-02: [Task description] — Completed
(If any task failed or was deferred, mark with ❌ or ⏳ and explain why)

### Modified Files Summary
- [relative path to modified file 1]
- [relative path to modified file 2]

### Implementation Confidence Score
- Confidence: [e.g. 98% or 100%]
- Rationale: [Detailed explanation of verified quality gates, passing linters, contract adherence, and zero regressions]
```
