# Folder Structure

**Version:** 5.0.0
**Updated:** 2026-10-27

---

## Overview

The `02-spec/` directory is the canonical location for all project specifications. It uses a **numbered folder hierarchy** organized into functional layers. This document explains the complete tree layout, how modules are grouped, and how the numbering scheme works.

---

## Root Folder Structure — Mandatory Base

> **This file is the single source of truth for all folder structure rules.** The root-level `folder-structure-root.md` is a redirect that points here.

```
IMPORTANT — AI INSTRUCTION:

1. Read this file BEFORE creating, reorganizing, or validating any spec folder.
2. Verify that ALL required root folders listed below exist and are in the correct order.
3. If any required folder is missing, create it immediately with at least a readme.md inside.
4. If folders are out of order or misnumbered, reorganize them to match this structure.
5. All additional folders and files MUST come AFTER the required base structure.
6. After any restructuring, update ALL cross-references that point to renamed or moved folders.
7. If ambiguity exists about folder placement, ASK before making structural assumptions.
```

### Numbering Policy

- **01–20**: Reserved exclusively for **core fundamentals** — foundational principles, standards, integrations, research, and consolidated summaries.
- **21+**: Reserved for **app-specific modules and audits** — application features, workflows, issue tracking, domain databases, UI tokens, and specification audits (`21-app/` through `25-spec-audits/`).

No app-specific content may appear in the 01–20 range. No foundational principles may appear at 21+. All 21+ directories reside directly under `02-spec/` as top-level containers.

### Required Root Folders

| # | Folder | Purpose |
|---|--------|---------|
| 01 | `01-spec-authoring-guide/` | Meta-guide — how to write and maintain specs |
| 02 | `02-coding-guidelines/` | Cross-language and language-specific coding standards |
| 03 | `03-error-manage/` | Error management, application errors, and error envelopes |
| 04 | `04-database-conventions/` | Database schema design, naming, ORM, views, testing, REST API format |
| 05 | `05-split-db-architecture/` | Split database pattern (operational + config DBs) |
| 06 | `06-seedable-config-architecture/` | Seedable configuration with versioning and changelogs |
| 07 | `07-design-system/` | Design system tokens, components, theming, visual rules |
| 08 | `08-docs-viewer-ui/` | Documentation viewer UI specifications |
| 09 | `09-code-block-system/` | Code block rendering pipeline, interactions, styling |
| 10 | `10-research/` | Comparative studies, technology evaluations, exploratory notes |
| 11 | `11-powershell-integration/` | PowerShell scripting conventions, cross-platform automation |
| 12 | `12-cicd-pipeline-workflows/` | CI/CD pipeline specs, deployment workflows, quality gates |
| 13 | `13-generic-cli/` | Generic CLI terminal UX, colors, flags, and help contracts |
| 14 | `14-update/` | CLI self-update blueprints, release distribution, deploy strategies |
| 15 | `15-distribution-and-runner/` | Cross-platform installers, runners, and repository forwarding |
| 16 | `16-generic-release/` | Universal release pipeline blueprint and asset matrices |
| 17 | `17-consolidated-guidelines/` | Single-file AI summaries of every major spec module |
| 18 | `18-wp-plugin-how-to/` | WordPress plugin architecture, admin UI, and REST API |
| 19 | `19-main-worker-service/` | Split-tier architecture, credential-blind proxy, and backup nodes |
| 20 | `20-*` *(Reserved)* | Reserved for future core fundamentals |
| 21 | `21-app/` | App-specific specs: features, workflows, architecture |
| 22 | `22-app-issues/` | App bug analysis, root cause analysis, fix documentation |
| 23 | `23-app-db/` | App-specific data model, table designs, migration strategies |
| 24 | `24-app-ui-design-system/` | App-specific UI, design system, theming, component patterns |
| 25 | `25-spec-audits/` | Specification audits, blind-AI readiness scorecards, and matrices |

### Rules

1. **Required folders are fixed** — Their numbering and naming MUST NOT change unless this file is explicitly updated and all cross-references are audited.
2. **App content starts at 21** — App-specific folders MUST use numbers 21+. They MUST NOT be placed in the 01–20 range.
3. **Top-level direct placement** — Folders 21+ (`21-app/` through `25-spec-audits/`) are top-level containers located directly under `02-spec/`. They must NEVER be placed inside `02-spec/02-coding-guidelines/`.
4. **Sequential numbering** — New folders use the next available number. Gaps are acceptable for historical reasons.
5. **Lowercase kebab-case** — All folder names: `{NN}-{kebab-case-name}/` (zero-padded, lowercase, hyphens only).
6. **Every folder needs `readme.md`** — Every root folder MUST contain a `readme.md` as its entry point.
7. **Cross-references must be updated after restructuring** — Search `02-spec/` for old paths, update all references, and run link validation.

### Validation Checklist

Before any spec generation, restructuring, or audit, verify:

- [ ] All required root folders (01–19, 21–25) exist directly under `02-spec/`
- [ ] Core fundamental folders are numbered strictly within 01–20
- [ ] App-specific and audit folders start at 21+
- [ ] All folders use lowercase kebab-case naming
- [ ] Every folder contains `readme.md`
- [ ] Root-level `readme.md` and `99-consistency-report.md` exist
- [ ] All cross-references resolve correctly
- [ ] Zero broken links reported across spec documentation

---

## Root-Level Tree

```
02-spec/
├── readme.md                               # Master index — links every module
├── folder-structure-root.md                # Redirect → this file is the canonical source
├── 99-consistency-report.md                # Root-level consistency report
│
│ ── CORE FUNDAMENTALS (01–20) ──
├── 01-spec-authoring-guide/                # THIS GUIDE — how to write specs
├── 02-coding-guidelines/                   # Consolidated coding standards
├── 03-error-manage/                        # Error management, codes, architecture
├── 04-database-conventions/                # Database schema design, naming, ORM, views
├── 05-split-db-architecture/               # Split database pattern
├── 06-seedable-config-architecture/        # Seedable configuration with versioning
├── 07-design-system/                       # Design system tokens, components, theming
├── 08-docs-viewer-ui/                      # Documentation viewer UI specifications
├── 09-code-block-system/                   # Code block rendering pipeline
├── 10-research/                            # Research: studies, evaluations, explorations
├── 11-powershell-integration/              # PowerShell scripting & automation
├── 12-cicd-pipeline-workflows/             # CI/CD pipelines, deployment workflows
├── 13-generic-cli/                         # Generic CLI UX, exit codes, and flags
├── 14-update/                              # CLI self-update, release distribution
├── 15-distribution-and-runner/             # Installers, runners, repository forwarding
├── 16-generic-release/                     # Universal release pipeline blueprint
├── 17-consolidated-guidelines/             # AI-readable summaries of all modules
├── 18-wp-plugin-how-to/                    # WordPress plugin architecture
├── 19-main-worker-service/                 # Split-tier service and proxy nodes
├── 20-reserved/                            # Reserved for future core fundamentals
│
│ ── APP-SPECIFIC & AUDITS (21+) ──
├── 21-app/                                 # App feature specs, workflows
├── 22-app-issues/                          # App bug analysis, root cause, fixes
├── 23-app-db/                              # App data models, migrations
├── 24-app-ui-design-system/                # App UI components, tokens
├── 25-spec-audits/                         # Specification audits and scorecards
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
| 02 | coding-guidelines | Consolidated language standards (cross-language, TS, Go, PHP, Rust, C#) |
| 03 | error-manage | Error management, application errors, and error envelopes |
| 04 | database-conventions | Database schema design, naming, ORM, views, testing, REST API |
| 05 | split-db-architecture | Split database pattern (operational + config DBs) |
| 06 | seedable-config-architecture | Seedable configuration with changelog versioning |
| 07 | design-system | Design system tokens, components, theming |
| 08 | docs-viewer-ui | Documentation viewer UI specifications |
| 09 | code-block-system | Code block rendering pipeline, interactions, styling |
| 10 | research | Comparative studies, technology evaluations, exploratory notes |
| 11 | powershell-integration | PowerShell scripting conventions, cross-platform automation |
| 12 | cicd-pipeline-workflows | CI/CD pipeline specs, deployment workflows, automation |
| 13 | generic-cli | Generic CLI UX, ANSI colors, flags, and help contracts |
| 14 | update | CLI self-update blueprints, release distribution, deploy strategies |
| 15 | distribution-and-runner | Cross-platform installers, runners, and repository forwarding |
| 16 | generic-release | Universal release pipeline blueprint and asset matrices |
| 17 | consolidated-guidelines | Single-file AI-readable summaries of all major modules |
| 18 | wp-plugin-how-to | WordPress plugin architecture, admin UI, and REST API |
| 19 | main-worker-service | Split-tier architecture, credential-blind proxy, and backup nodes |

### Layer 2: App-Specific & Audits (21+)

Application-level content — features, workflows, issue tracking, databases, UI design systems, and specification audits.

| # | Module | Purpose |
|---|--------|---------|
| 21 | app | App feature specs, workflows, architecture decisions |
| 22 | app-issues | App bug analysis, root cause analysis, fix documentation |
| 23 | app-db | App-specific data models, table designs, and migration strategies |
| 24 | app-ui-design-system | App-specific UI components, tokens, and design patterns |
| 25 | spec-audits | Specification audits, blind-AI readiness scorecards, and matrices |

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
| 00 | Root files (overview, folder guideline) and `readme.md` within folders |
| 01–20 | **Core fundamentals** — standards, principles, integrations, research, consolidated summaries |
| 21+ | **App-specific & audits** — application features, workflows, issue analysis, UI, databases, audits |
| 97 | Acceptance criteria |
| 98 | Changelogs |
| 99 | Consistency reports and archive |

---

## Module Templates

Two templates exist depending on project type:

| Type | Template | Pattern |
|------|----------|---------|
| **CLI tool** | [CLI Module Template](./06-cli-module-template.md) | `01-backend/`, `02-frontend/`, `03-deploy/` |
| **App / WordPress** | [App Project Template](./07-app-project-template.md) | `02-fundamentals.md`, `02-features/`, `03-issues/` |

- **CLI tools** — Go/Rust command-line tools with backend, optional frontend, and deploy specs
- **App/WordPress projects** — Plugins, themes, or web apps with features and issue tracking. Features use `01-backend.md`, `02-frontend.md`, `03-wp-admin.md` sub-structure.

---

## Subfolder Depth Rules

- **Maximum depth:** 3 levels (e.g., `02-spec/13-generic-cli/03-subcommand-architecture.md`)
- **Each level** follows the same `{NN}-{name}` convention
- **Every folder** at any depth must have `readme.md`
- **Only top-level** and major subfolder boundaries require `99-consistency-report.md`

---

## Non-Contiguous Numbering

Historical gaps exist in the numbering sequence (e.g., numbers previously used for modules now consolidated elsewhere). These gaps are preserved to avoid mass-renaming existing modules. New modules should use the **next available number** after the highest existing module within their range (01–20 for core, 21+ for app).

---

## Direct Top-Level Placement for 21+ Containers

A critical architectural invariant is that **all 21+ containers reside directly under `02-spec/`**:

```
02-spec/
├── 21-app/
├── 22-app-issues/
├── 23-app-db/
├── 24-app-ui-design-system/
└── 25-spec-audits/
```

- **NEVER nest under `02-coding-guidelines/`**: `02-spec/02-coding-guidelines/` is strictly reserved for language-specific standards (`01-cross-language/`, `02-golang/`, `03-typescript/`, etc.). App-specific specifications or audit documents must never be created inside `02-coding-guidelines/`.
- **Isolation Boundary**: Cross-repository synchronization tools strictly exclude `02-spec/21-*` through `02-spec/25-*` from syncing into external repositories to prevent domain-specific leaks.

### Decision Guide — Where to Place New Content

```
AI INSTRUCTION — Spec Placement Decision Guide:

1. Is the content reusable, foundational, or principle-driven across all codebases?
   → Place in 01–20 (core fundamentals range under 02-spec/)

2. Is the content exploratory, comparative, or evaluative?
   → Place in 02-spec/10-research/

3. Does the content define a specific application feature or workflow?
   → Place in 02-spec/21-app/

4. Does the content analyze application bugs, failures, root causes, or fixes?
   → Place in 02-spec/22-app-issues/

5. Does the content define application data models, tables, or database migrations?
   → Place in 02-spec/23-app-db/

6. Does the content define application-specific UI tokens, screens, or layouts?
   → Place in 02-spec/24-app-ui-design-system/

7. Does the content contain specification audits, readiness scorecards, or gap analysis?
   → Place in 02-spec/25-spec-audits/

8. Never place app-specific content in the 01–20 range.
9. Never place foundational principles in the 21+ range.
```
