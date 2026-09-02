---
name: spec-authoring
description: >-
  Autonomously author, structure, and validate specification modules adhering to spec/01-spec-authoring-guide/
  with mandatory overview, acceptance criteria, and consistency reports.
---

# Specification Authoring & Quality Assurance Skill

Autonomously create, maintain, and audit specification files and folders adhering to `spec/01-spec-authoring-guide/` and `spec/17-consolidated-guidelines/01-spec-authoring.md`.

## Core Checkpoints & Mandatory Invariants

1. **Folder & File Layout:**
   - Folder naming: Lowercase kebab-case with 2-digit numeric prefixes `{NN}-{kebab-name}/`.
   - Core fundamentals: Modules `01` through `20`.
   - App-specific modules: Modules `21` and above.
   - File naming: Lowercase kebab-case `{NN}-{kebab-name}.md`. No spaces, underscores, or camelCase.

2. **Required Files Per Module:**
   - `00-overview.md`: Entry point with H1 title, metadata block (Version, Updated, Status, AI Confidence, Ambiguity), purpose, scoring, keywords, and file inventory table.
   - `97-acceptance-criteria.md`: Verification commands and Given/When/Then criteria.
   - `99-consistency-report.md`: Consistency audit across the module.

3. **Mandatory Module Scoring:**
   - **AI Confidence Score:** `Production-Ready` (all contracts defined), `High`, `Medium`, or `Low`.
   - **Ambiguity Score:** `None` (zero interpretation required), `Low`, `Medium`, `High`, or `Critical`.
   - **Health Score (0-100):** Calculated from 4 components (25% `00-overview.md` present + 25% `99-consistency-report.md` present + 25% lowercase kebab-case naming + 25% numeric prefixes).

4. **Cross-Reference Integrity:**
   - Use relative paths only (e.g. `../02-coding-guidelines/00-overview.md`). Never use root-relative or absolute paths.
   - Always include the `.md` file extension.
   - Validate using `python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .` and `node linter-scripts/generate-dashboard-data.cjs`.

5. **Institutional Memory Interaction:**
   - Specs document formal system architecture and contracts (`spec/`).
   - `.lovable/memory/` documents institutional knowledge, patterns, and decisions. Never duplicate full spec content in memory.
