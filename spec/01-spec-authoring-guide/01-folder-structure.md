# Folder Structure

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The `spec/` directory is the canonical location for all project specifications. It uses a **numbered folder hierarchy** organized into functional layers. This document explains the complete tree layout, how modules are grouped, and how the numbering scheme works.

---

## Root Folder Structure — Mandatory Base

> **This file is the single source of truth for all folder structure rules.** The root-level `folder-structure-root.md` is a redirect that points here.

```
IMPORTANT — AI INSTRUCTION:

1. Read this file BEFORE creating, reorganizing, or validating any spec folder.
2. Verify that ALL required root folders listed below exist and are in the correct order.
3. If any required folder is missing, create it immediately with at least a 00-overview.md inside.
4. If folders are out of order or misnumbered, reorganize them to match this structure.
5. All additional folders and files MUST come AFTER the required base structure.
6. After any restructuring, update ALL cross-references that point to renamed or moved folders.
7. If ambiguity exists about folder placement, ASK before making structural assumptions.
```

### Numbering Policy

- **01–20**: Reserved exclusively for **core fundamentals** — foundational principles, standards, integrations, research, and consolidated summaries.
- **21+**: Reserved for **app-specific** content — application features, workflows, and issue tracking.

No app-specific content may appear in the 01–20 range. No foundational principles may appear at 21+.

### Required Root Folders

| # | Folder | Purpose |
|---|--------|---------|
| 01 | `01-spec-authoring-guide/` | Meta-guide — how to write and maintain specs |
| 02 | `02-coding-guidelines/` | Cross-language and language-specific coding standards |
| 03 | `03-error-manage/` | Error management, error codes, error architecture |
| 04 | `04-database-conventions/` | Database schema design, naming, ORM, views, testing, REST API format |
| 05 | `05-split-db-architecture/` | Split database pattern (operational + config DBs) |
| 06 | `06-seedable-config-architecture/` | Seedable configuration with versioning |
| 07 | `07-design-system/` | Design system tokens, components, theming |
| 08 | `08-docs-viewer-ui/` | Documentation viewer UI specifications |
| 09 | `09-code-block-system/` | Code block rendering pipeline, interactions, styling |
| 10 | `11-powershell-integration/` | PowerShell scripting conventions, cross-platform automation |
| 11 | `10-research/` | Comparative studies, technology evaluations, exploratory notes |
| 12 | `17-consolidated-guidelines/` | AI-readable summaries of every major spec module |
| 13 | `12-cicd-pipeline-workflows/` | CI/CD pipeline specs, deployment workflows, automation |
| 14 | `14-update/` | CLI self-update blueprints, release distribution, deploy strategies |
| 15–20 | _Reserved_ | Future core fundamental topics |
| 21 | `21-app/` | App-specific specs: features, workflows, architecture |
| 22 | `22-app-issues/` | App bug analysis, root cause analysis, fix documentation |

### Rules

1. **Required folders are fixed** — Their numbering and naming MUST NOT change unless this file is explicitly updated and all cross-references are audited.
2. **App content starts at 21** — App-specific folders MUST use numbers 21+. They MUST NOT be placed in the 01–20 range.
3. **Sequential numbering** — New folders use the next available number. Gaps are acceptable for historical reasons.
4. **Lowercase kebab-case** — All folder names: `{NN}-{kebab-case-name}/` (zero-padded, lowercase, hyphens only).
5. **Every folder needs `00-overview.md`** — Every root folder MUST contain a `00-overview.md` as its entry point.
6. **Cross-references must be updated after restructuring** — Grep `spec/` for old paths, update all references, run link scanner.

### Validation Checklist

Before any spec generation, restructuring, or audit, verify:

- [ ] All required root folders (01–14, 21–23) exist
- [ ] Core fundamental folders are numbered within 01–20
- [ ] App-specific folders start at 21+
- [ ] All folders use lowercase kebab-case naming
- [ ] Every folder contains `00-overview.md`
- [ ] Root-level `00-overview.md` and `99-consistency-report.md` exist
- [ ] All cross-references resolve correctly
- [ ] Zero broken links reported by `node linter-scripts/generate-dashboard-data.cjs`

---

## Root-Level Tree

```
spec/
├── 00-overview.md                          # Master index — links every module
├── folder-structure-root.md                # Redirect → this file is the canonical source
├── 99-consistency-report.md               # Root-level consistency report
│
│ ── CORE FUNDAMENTALS (01–20) ──
├── 01-spec-authoring-guide/                 # THIS GUIDE — how to write specs
├── 02-coding-guidelines/                   # Consolidated coding standards
├── 03-error-manage/                   # Error management, codes, architecture
├── 04-database-conventions/                # Database schema design, naming, ORM, views
├── 05-split-db-architecture/               # Split database pattern
├── 06-seedable-config-architecture/        # Seedable configuration with versioning
├── 07-design-system/                       # Design system tokens, components, theming
├── 08-docs-viewer-ui/                      # Documentation viewer UI specifications
├── 09-code-block-system/                   # Code block rendering pipeline
├── 11-powershell-integration/              # PowerShell scripting & automation
├── 10-research/                            # Research: studies, evaluations, explorations
├── 17-consolidated-guidelines/             # AI-readable summaries of all modules
├── 12-cicd-pipeline-workflows/             # CI/CD pipelines, deployment workflows
├── 14-update/              # CLI self-update, release distribution
├── 15–20 (reserved)                        # Future core fundamentals
│
│ ── APP-SPECIFIC (21+) ──
├── 21-app/                                  # App feature specs, workflows
├── 22-app-issues/                          # App bug analysis, root cause, fixes
│
│ ── ARCHIVE & GOVERNANCE ──
├── 99-archive/                             # Deprecated specifications
└── validation-reports/                     # Audit and validation artifacts
```

---

## Functional Layers

Modules are organized into three conceptual layers:

### Layer 1: Core Fundamentals (01–20)

Foundational rules, standards, and integrations that all other work depends on. This range is **reserved** — no app-specific content may appear here.

| # | Module | Purpose |
|---|--------|---------|
| 01 | spec-authoring-guide | This guide — how to write and maintain specs |
| 02 | coding-guidelines | Consolidated language standards (cross-language, TS, Go, PHP, Rust) |
| 03 | error-manage | Error management, error codes, error architecture |
| 04 | database-conventions | Database schema design, naming, ORM, views, testing, REST API |
| 05 | split-db-architecture | Split database pattern (operational + config DBs) |
| 06 | seedable-config-architecture | Seedable configuration with changelog versioning |
| 07 | design-system | Design system tokens, components, theming |
| 08 | docs-viewer-ui | Documentation viewer UI specifications |
| 09 | code-block-system | Code block rendering pipeline, interactions, styling |
| 10 | powershell-integration | PowerShell scripting conventions, cross-platform automation |
| 11 | research | Comparative studies, technology evaluations, exploratory notes |
| 12 | consolidated-guidelines | Single-file AI-readable summaries of all major modules |
| 13 | cicd-pipeline-workflows | CI/CD pipeline specs, deployment workflows, automation |
| 14 | self-update-app-update | CLI self-update blueprints, release distribution |

### Layer 2: App-Specific (21+)

Application-level content — features, workflows, and issue tracking.

| # | Module | Purpose |
|---|--------|---------|
| 21 | app | App feature specs, workflows, architecture decisions |
| 22 | app-issues | App bug analysis, root cause analysis, fix documentation |

### Layer 3: Archive & Governance (99, validation-reports)

Historical and governance artifacts.

| # | Module | Purpose |
|---|--------|---------|
| 99 | archive | Deprecated and superseded specs |
| — | validation-reports | Audit certificates and validation logs |

---

## Reserved Number Ranges

| Range | Purpose |
|-------|---------|
| 00 | Root files (overview, folder guideline) and `00-overview.md` within folders |
| 01–20 | **Core fundamentals** — standards, principles, integrations, research, consolidated summaries |
| 21+ | **App-specific** — application features, workflows, issue analysis |
| 97 | Acceptance criteria |
| 98 | Changelogs |
| 99 | Consistency reports and archive |

---

## Module Templates

Two templates exist depending on project type:

| Type | Template | Pattern |
|------|----------|---------|
| **CLI tool** | [CLI Module Template](./04-cli-module-template.md) | `01-backend/`, `02-frontend/`, `03-deploy/` |
| **App / WordPress** | [App Project Template](./05-app-project-template.md) | `01-fundamentals.md`, `02-features/`, `03-issues/` |

- **CLI tools** — Go/Rust command-line tools with backend, optional frontend, and deploy specs
- **App/WordPress projects** — Plugins, themes, or web apps with features and issue tracking. Features use `01-backend.md`, `02-frontend.md`, `03-wp-admin.md` sub-structure.

---

## Subfolder Depth Rules

- **Maximum depth:** 3 levels (e.g., `spec/13-wp-plugin/03-exam-manager/02-features/01-exam-builder/`)
- **Each level** follows the same `{NN}-{name}` convention
- **Every folder** at any depth must have `00-overview.md`
- **Only top-level** and major subfolder boundaries require `99-consistency-report.md`

---

## Non-Contiguous Numbering

Historical gaps exist in the numbering sequence (e.g., numbers previously used for modules now consolidated elsewhere). These gaps are preserved to avoid mass-renaming existing modules. New modules should use the **next available number** after the highest existing module within their range (01–20 for core, 21+ for app).

---

## Coding Guidelines Subfolder Numbering Convention

The `spec/02-coding-guidelines/` directory uses a **split numbering scheme** that separates core fundamentals from app-specific content:

### Core Fundamentals Range (01–20)

Folders 01–20 are **reserved exclusively** for core, foundational coding standards. This includes:

- Language-specific standards (TypeScript, Go, PHP, Rust, C#)
- Cross-language principles (naming, DRY, typing, booleans)
- AI optimization rules
- File & folder naming conventions
- PowerShell integration (`09`)
- Research — comparative studies, technology evaluations (`10`)
- Security policies (`11`)
- Database conventions (`04` — root level)

**No app-specific content may appear in this range.** New core fundamentals use the next available number within 01–20.

### App-Specific Range (21+)

Folders 21 and above are for **application-level** content:

| # | Folder | Purpose |
|---|--------|---------|
| 21 | `21-app/` | App feature specs, workflows, architecture decisions |
| 22 | `22-app-issues/` | App bug analysis, root cause analysis, fix documentation |

New app-related folders use the next available number after 22.

### Decision Guide — Where to Place New Content

```
AI INSTRUCTION — Coding Guidelines Placement:

1. Is the content reusable, foundational, or principle-driven?
   → Place in 01–20 (core fundamentals range)

2. Is the content exploratory, comparative, or evaluative?
   → Place in 10-research

3. Does the content define a specific application feature or workflow?
   → Place in 21-app

4. Does the content analyze bugs, failures, root causes, or fixes?
   → Place in 22-app-issues

5. Never place app-specific content in the 01–20 range.
6. Never place foundational principles in the 21+ range.
```

### Fixed Folder Assignments

| # | Folder | Status |
|---|--------|--------|
| 09 | `09-powershell-integration/` | Fixed — all PowerShell content goes here |
| 10 | `10-research/` | Fixed — all research content goes here |
| 21 | `21-app/` | Fixed — all app specs go here |
| 22 | `22-app-issues/` | Fixed — all app issue analysis goes here |

See [`00-overview.md`](../02-coding-guidelines/00-overview.md) for the full category listing.
