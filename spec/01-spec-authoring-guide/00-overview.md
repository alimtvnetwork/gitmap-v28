# Spec Authoring Guide

**Version:** 3.2.0  
**Updated:** 2026-04-16  
**Status:** Active  
**AI Confidence:** Production-Ready  
**Ambiguity:** None

---

## Overview

This is the **definitive guide** for any AI agent or human contributor to understand, navigate, create, and maintain specifications in this repository. It covers folder structure, naming conventions, required files, templates, cross-referencing, scoring metrics, reliability validation, and the `.lovable/` institutional memory system — everything needed to produce spec-compliant documentation from scratch.

### How to Write an Overview (`00-overview.md`)

Every module's `00-overview.md` must follow this structure:

1. **Title** — H1 heading matching the module name
2. **Metadata block** — Version, Updated date, Status, AI Confidence Score (0–100%), Ambiguity Score (0–100%)
3. **Overview paragraph** — 1–2 paragraphs explaining what the module covers, its purpose, and who benefits
4. **Scoring section** — AI Confidence Score, Ambiguity Score, and Health Score with thresholds
5. **Keywords section** — Searchable tags for AI discovery
6. **Files table** — Numbered inventory of all files with links and descriptions
7. **Cross-References table** — Links to related modules, memories, and external docs

> **Rule:** The overview is the **entry point**. An AI agent reading only the overview should understand the module's scope, readiness, and how to find any file within it.

---

## Scoring Metrics

Every `00-overview.md` MUST include these three scores:

### AI Confidence

Measures how ready the specification is for an AI agent to implement the described feature.

| Tier | Icon | Meaning | When to Use |
|------|------|---------|-------------|
| **Production-Ready** | ✅ | Specs are complete, unambiguous, and fully implementable | All interfaces defined, acceptance criteria explicit, error codes mapped |
| **High** | 🟢 | Minor gaps exist but AI can proceed with reasonable assumptions | Most sections complete; a few edge cases undefined |
| **Medium** | 🟡 | Significant gaps; AI will need clarification or make risky assumptions | Missing types, partial acceptance criteria, unclear validation rules |
| **Low** | 🔴 | Major sections missing; do NOT attempt implementation | No interfaces, no acceptance criteria, vague requirements |

**Factors that increase confidence:** Complete interfaces, explicit acceptance criteria, error code mappings, clear data models, defined API contracts.

**Factors that decrease confidence:** Missing types, vague requirements ("should handle errors"), undefined edge cases, no acceptance criteria.

### Ambiguity

Measures how much interpretation is required. **Lower tiers are better.**

| Tier | Icon | Meaning | When to Use |
|------|------|---------|-------------|
| **None** | ✅ | No interpretation needed; every detail is explicit | All fields typed, all flows documented, all errors handled |
| **Low** | 🟢 | A few areas need assumptions; acceptable for implementation | Minor gaps in edge cases or optional features |
| **Medium** | 🟡 | Multiple areas require interpretation; review recommended | Several undefined behaviors, partial validation rules |
| **High** | 🟠 | Many areas open to interpretation; rewrite recommended | Missing data models, unclear permissions, vague UI specs |
| **Critical** | 🔴 | Spec is too vague to implement; MUST be rewritten | No clear structure, contradictory requirements, missing core definitions |

**Common ambiguity sources:** Undefined field types, unclear validation rules, missing error handling paths, unspecified permissions, vague UI requirements.

### Health Score (0–100)

Structural compliance score calculated by the dashboard scanner:

| Criterion | Weight |
|-----------|--------|
| `00-overview.md` present | 25% |
| `99-consistency-report.md` present | 25% |
| Lowercase kebab-case naming | 25% |
| Unique numeric sequence prefixes | 25% |

**100/100 = A+** — All four criteria met.

---

## Keywords

`spec-authoring` · `ai-guide` · `folder-structure` · `naming-conventions` · `required-files` · `module-template` · `cli-template` · `cross-references` · `health-score` · `ai-confidence` · `ambiguity-score` · `consistency-report` · `reliability-report` · `memory-folder` · `lovable-folder` · `kebab-case` · `numeric-prefix`

---

## File Categories

Every spec file falls into one of these categories. When creating a new file, assign it the correct category to help AI agents and contributors navigate:

| Category | Purpose | Examples |
|----------|---------|----------|
| **Overview** | Module entry point with metadata and file index | `00-overview.md` |
| **Architecture** | System design, data flow, component structure | `01-architecture.md`, `02-data-model.md` |
| **API / Interface** | Endpoints, contracts, request/response schemas | `03-api-design.md`, `04-rest-endpoints.md` |
| **Logic / Rules** | Business logic, validation rules, algorithms | `05-validation-rules.md`, `06-scoring-logic.md` |
| **UI / Frontend** | Component specs, layouts, user flows | `07-ui-components.md`, `08-user-flows.md` |
| **Backend** | Server-side implementation, database, services | `09-database-schema.md`, `10-service-layer.md` |
| **Diagrams** | Mermaid diagrams, architecture visuals, flow charts | `11-diagrams.md`, `12-sequence-diagrams.md` |
| **Testing** | Test plans, acceptance criteria, QA standards | `97-acceptance-criteria.md` |
| **Meta / Reports** | Consistency reports, changelogs, migration notes | `99-consistency-report.md`, `98-changelog.md` |

> **Rule:** Assign a category in the file's metadata block so AI agents can filter and prioritize what to read.

---

## Folder Structure Examples

### Standard Module (Flat)

```
spec/17-research-queries/
├── 00-overview.md
├── 01-query-types.md
├── 02-search-algorithms.md
├── 97-acceptance-criteria.md
└── 99-consistency-report.md
```

### CLI Tool Module (3-Folder Pattern)

```
spec/09-gsearch-cli/
├── 00-overview.md
├── 01-backend/
│   ├── 01-architecture.md
│   ├── 02-commands.md
│   └── 03-api-design.md
├── 02-frontend/
│   ├── 01-ui-components.md
│   └── 02-user-flows.md
├── 03-diagrams/
│   └── 01-architecture-diagram.md
├── 97-acceptance-criteria.md
└── 99-consistency-report.md
```

### WordPress / App Module (Features + Issues Pattern)

App and WordPress projects use `01-fundamentals.md` as the first content file, then `02-features/` and `03-issues/` folders:

```
spec/13-wp-plugin/03-exam-manager/
├── 00-overview.md
├── 01-fundamentals.md                    # Core architecture, schema, lifecycle
│
├── 02-features/                          # Feature specifications
│   ├── 00-overview.md                   # Feature index with status table
│   ├── 01-exam-builder/
│   │   ├── 00-overview.md
│   │   ├── 01-backend.md
│   │   ├── 02-frontend.md
│   │   └── 03-wp-admin.md
│   └── 02-question-bank/
│       ├── 00-overview.md
│       ├── 01-backend.md
│       └── 02-frontend.md
│
├── 03-issues/                            # Tracked issues and investigations
│   ├── 00-overview.md                   # Issue index with status/severity
│   ├── 01-score-rounding/
│   │   ├── 00-overview.md
│   │   ├── 01-investigation.md
│   │   └── 02-resolution.md
│   └── 02-wp-65-compat.md              # Simple issues can be single files
│
├── 97-acceptance-criteria.md
└── 99-consistency-report.md
```

> **Key insight:** App/WP projects split features into `01-backend.md`, `02-frontend.md`, `03-wp-admin.md` inside each feature folder. Issues follow the same `{NN}-{kebab-name}` convention with `00-overview.md` required for multi-file issues. See [05-app-project-template.md](./05-app-project-template.md) for the full template.

---

## Files

| # | File | Category | Description |
|---|------|----------|-------------|
| 01 | [01-folder-structure.md](./01-folder-structure.md) | Architecture | Complete spec tree layout, layer grouping, and numbering ranges |
| 02 | [02-naming-conventions.md](./02-naming-conventions.md) | Rules | File and folder naming rules (kebab-case, numeric prefixes, reserved ranges) |
| 03 | [03-required-files.md](./03-required-files.md) | Rules | Mandatory files every module must contain (overview, consistency report, etc.) |
| 04 | [04-cli-module-template.md](./04-cli-module-template.md) | Template | Step-by-step template for CLI tool spec modules (3-folder pattern) |
| 05 | [05-app-project-template.md](./05-app-project-template.md) | Template | Template for app/WordPress projects (fundamentals + features + issues) |
| 06 | [06-non-cli-module-template.md](./06-non-cli-module-template.md) | Template | Template for flat/non-CLI modules (research, utilities, standards) |
| 07 | [07-memory-folder-guide.md](./07-memory-folder-guide.md) | Guide | Structure and conventions for the `.lovable/memories/` tree |
| 08 | [08-cross-references.md](./08-cross-references.md) | Rules | How to write cross-references, relative paths, and link integrity rules |
| 09 | [09-exceptions.md](./09-exceptions.md) | Rules | All known exception cases with folder structure examples |
| 10 | [10-mandatory-linter-infrastructure.md](./10-mandatory-linter-infrastructure.md) | Rules | Mandatory linter scripts — AI must verify presence before validation |
| 11 | [11-root-readme-conventions.md](./11-root-readme-conventions.md) | Rules | **MANDATORY** root `readme.md` format — centered icon, hero block, author/company template, badges, §9 release-blocker checklist |

## File Naming Convention (Quick Reference)

All files and folders in `spec/` and `.lovable/` MUST use **lowercase kebab-case**:

```
✅ 01-backend/                   ✅ 00-overview.md
✅ 09-gsearch-cli/               ✅ 03-api-design.md
✅ .lovable/memories/workflow/   ✅ file-naming-conventions.md

❌ 01-Backend/                   ❌ ApiDesign.md
❌ 09_gsearch_cli/               ❌ file_naming.md
```

**Numeric prefixes** are mandatory for spec files/folders and optional for memory files. See [02-naming-conventions.md](./02-naming-conventions.md) for full rules.

---

## Cross-Reference Validation Checklist

Every spec file must pass these cross-reference checks:

| # | Check | Rule |
|---|-------|------|
| 1 | **Relative paths only** | Never use root-relative (`/spec/...`) or absolute filesystem paths |
| 2 | **File extension included** | Always end with `.md` — never bare paths |
| 3 | **Target exists** | Every linked file must exist on disk |
| 4 | **Lowercase paths** | Path segments must be entirely lowercase kebab-case |
| 5 | **Depth correctness** | Count `../` levels carefully from source to target |
| 6 | **Bidirectional linking** | If module A references module B, module B should reference A |
| 7 | **Post-rename audit** | After any module renumbering, grep and update ALL references |

**Automated validation:** Run `node linter-scripts/generate-dashboard-data.cjs` — the output JSON reports all broken links. Zero broken links = passing.

See [08-cross-references.md](./08-cross-references.md) for full syntax and examples.

---

## Reliability Check Report

Every module MUST include a **reliability risk assessment** to evaluate implementation feasibility before coding begins. Reports are stored in `spec/validation-reports/` or inline within the module. The only allowed exception: pure-documentation modules with no implementable surface (e.g. `spec/_template.md`, `spec/folder-structure-root.md`) MAY omit the assessment when an `<!-- AUTHORING-WAIVER: documentation-only -->` comment is present near the title.

### What It Covers

| Section | Content |
|---------|---------|
| **Complexity Tier** | Simple / Medium / Complex Agentic / End-to-End |
| **Success Probability** | Estimated % chance of first-pass implementation success |
| **Failure Modes** | Where, why, and how failures can manifest |
| **Risk Mitigations** | Specific actions to reduce failure likelihood |
| **Dependency Risks** | External APIs, shared modules, or integrations that add risk |

### When to Create

- Before implementing any **Complex Agentic** or **End-to-End** module
- When a module's AI Confidence Score is below 70%
- When multiple modules have interdependencies
- After a major spec rewrite or architectural change

---

## The `.lovable/` Folder

The `.lovable/` directory is the **institutional knowledge hub** for the project. It persists AI-learned patterns, decisions, and workflows across sessions.

### Canonical Structure

```
.lovable/
├── memories/                    # ← CANONICAL memory folder (single source of truth)
│   ├── 00-memory-index.md       # Complete inventory of all memory files
│   ├── readme.md                # Simplified high-level overview
│   ├── architecture/            # System design decisions
│   ├── constraints/             # Hard rules (e.g., no-code policy, coding standards)
│   ├── features/                # Feature-specific knowledge
│   ├── guidelines/              # Development guidelines
│   ├── logic/                   # Business logic, formulas, algorithms
│   ├── patterns/                # Reusable templates
│   ├── pending/                 # Work-in-progress / pending tasks
│   ├── planned/                 # Planned tasks (queued for future work)
│   ├── done/                    # Completed tasks archive
│   ├── completed-issues/        # Resolved issues archive
│   ├── project/                 # Project status and tracking
│   ├── qa/                      # Quality standards
│   ├── reports/                 # Reliability reports, audit reports
│   ├── spec-management/         # Spec management conventions
│   ├── suggestions/             # Suggestion tracking
│   │   └── completed/           # Archived completed suggestions
│   ├── training/                # AI training materials
│   ├── ui/                      # UI component patterns
│   ├── workflow/                # Process conventions
│   └── wp-plugins/              # WordPress plugin knowledge
├── plan.md                      # Current execution plan
├── reliability-risk-report.md   # Project-level reliability assessment
└── [other root files]           # Standards archive, audit history, etc.
```

### Task & Issue Tracking Folders

The memory folder includes dedicated folders for tracking work items:

| Folder | Purpose | When to Use |
|--------|---------|-------------|
| `pending/` | Tasks currently in progress or awaiting action | Active work items, blocked tasks |
| `planned/` | Tasks queued for future execution | Upcoming batches, prioritized backlog |
| `done/` | Completed tasks archive | Finished work items (move from pending/planned) |
| `completed-issues/` | Resolved issues archive | Bug fixes, resolved problems, closed issues |

> **Workflow:** Create task files in `planned/` → move to `pending/` when work starts → move to `done/` when complete. Issues follow the same flow but end in `completed-issues/`.

### Consolidation Rule

> **There is only ONE memory folder: `.lovable/memories/`.** The legacy `.lovable/memory/` variant is prohibited. If found during audits, migrate all contents to `.lovable/memories/` and delete the legacy folder.

### Memory File Conventions

| Rule | Detail |
|------|--------|
| **Naming** | Lowercase kebab-case; numeric prefixes optional (unlike spec files) |
| **Format** | Every file: H1 title → metadata (Updated, Version, Status) → Overview → Content → Cross-References |
| **Index** | Always update `00-memory-index.md` when adding/removing files |
| **Depth** | Maximum 2 levels: `memories/{category}/{file}.md` |
| **No duplication** | Don't duplicate information already in spec files |

### What Goes in Memories vs. Specs

| Content | Location |
|---------|----------|
| Formal specifications, APIs, data models | `spec/` |
| Architectural decisions, conventions, patterns | `.lovable/memories/` |
| Execution plans and batch tracking | `.lovable/plan.md` |
| Suggestion tracking | `.lovable/memories/suggestions/` |
| Pending / planned / done tasks | `.lovable/memories/pending/`, `planned/`, `done/` |
| Completed issues | `.lovable/memories/completed-issues/` |
| Reliability assessments | `.lovable/memories/reports/` or `spec/validation-reports/` |

### Where AI Should Write Updates

When an AI agent learns something new or the user provides instructions:

| User Says | AI Writes To |
|-----------|-------------|
| "Remember this pattern" | `.lovable/memories/patterns/` or relevant category |
| "Add this to the plan" | `.lovable/plan.md` |
| "Track this task" | `.lovable/memories/planned/` or `pending/` |
| "This issue is resolved" | Move to `.lovable/memories/completed-issues/` |
| "Update coding guidelines" | `.lovable/memories/constraints/` |
| "Add WP plugin spec" | `spec/XX-wp-plugin-name/` (spec tree) |
| "Remember WP plugin convention" | `.lovable/memories/wp-plugins/` |

See [07-memory-folder-guide.md](./07-memory-folder-guide.md) for the complete memory folder guide.

---

## Quick Start for AI Agents

> 🔴 **MANDATORY — Commit to Memory Before Implementation**
>
> After reading this spec authoring guide and the linked coding guidelines, you **MUST** internalize and retain the following rules in your working memory for the entire session. These are non-negotiable and must be applied to every code change, every file you create, and every review you perform:
>
> 1. **Error Management is the highest priority** — error handling patterns from [03-error-manage/](../03-error-manage/00-overview.md) must be implemented from the very first line of code. Never defer error handling to "later."
> 2. **Boolean and if/else naming** — follow the strict boolean naming conventions (`is`/`has`/`should` prefixes, positive-only names, no negatives). Extract complex conditions into named variables. See [Coding Guidelines](../02-coding-guidelines/00-overview.md).
> 3. **Database conventions** — singular table names, PascalCase, `{TableName}Id` integer PKs, exact FK name matching. See [Database Conventions](../04-database-conventions/00-overview.md).
> 4. **Never hallucinate** — if any requirement is unclear, ambiguous, or missing from the spec, **stop and ask clarifying questions** rather than guessing or making assumptions. Incorrect assumptions waste more time than asking.
> 5. **Zero-nesting rule** — no nested `if` blocks. Use early returns and guard clauses.
>
> Failure to follow these rules will produce code that fails review and must be rewritten.

1. **Read [01-folder-structure.md](./01-folder-structure.md)** — the single source of truth for all folder structure rules
2. **Read this overview** to understand the file inventory, scoring, and conventions
3. **Read [01-folder-structure.md](./01-folder-structure.md)** for the tree layout
4. **Read [02-naming-conventions.md](./02-naming-conventions.md)** for naming rules
5. **Read [03-required-files.md](./03-required-files.md)** for mandatory file checklist
6. **Choose a template**: [04-cli-module-template.md](./04-cli-module-template.md), [05-app-project-template.md](./05-app-project-template.md), or [06-non-cli-module-template.md](./06-non-cli-module-template.md)
7. **Check [09-exceptions.md](./09-exceptions.md)** for edge cases before creating files
8. **Verify linter infrastructure** — read [10-mandatory-linter-infrastructure.md](./10-mandatory-linter-infrastructure.md) and confirm `linter-scripts/` exists
9. **Score the module** — set AI Confidence and Ambiguity percentages in `00-overview.md`
10. **Validate cross-references** — run the link scanner and fix any broken links
11. **Check `.lovable/memories/`** — read relevant memories before writing new specs

---

## Folder Structure Enforcement

When asked to "follow the spec authoring guideline and fix the folder structure," an AI agent MUST perform these steps in order:

### Step 1 — Verify Root Structure

1. Read [`01-folder-structure.md`](./01-folder-structure.md) — the single source of truth for all folder structure rules
2. Confirm all required root folders (01–11, 21–22) exist
3. Confirm they are correctly numbered and named
4. If any are missing, create them with at least `00-overview.md`

### Step 2 — Verify Folder Ordering

1. Core fundamentals use the 01–20 range; app-specific content uses 21+
2. No app-specific folders may appear in the 01–20 range
3. Numbering must be sequential (gaps are acceptable for historical reasons)

### Step 3 — Verify Naming Conventions

1. All folders use lowercase kebab-case: `{NN}-{kebab-case-name}/`
2. All files use lowercase kebab-case: `{NN}-{kebab-case-name}.md`
3. No spaces, underscores, or camelCase anywhere

### Step 4 — Verify Required Files

1. Every folder has `00-overview.md`
2. Top-level folders have `99-consistency-report.md`
3. Root `spec/` has `00-overview.md` and `99-consistency-report.md` (plus `folder-structure-root.md` as a redirect)

### Step 5 — Update Cross-References

1. If any folder was renamed, renumbered, or moved, grep for all old references
2. Update every broken reference to the new path
3. Run the link scanner: `node linter-scripts/generate-dashboard-data.cjs`
4. Confirm zero broken links

> **This process is NOT optional.** If inconsistencies exist, they MUST be fixed before any new spec work begins.

---

## Document Inventory

| File |
|------|
| 97-acceptance-criteria.md |
| 98-changelog.md |
| 99-consistency-report.md |


## Author & Attribution

All specifications in this repository are authored by **Md. Alim Ul Karim** — Chief Software Engineer, **Top 1% Crossover** status, [Stack Overflow](https://stackoverflow.com/users/513511/md-alim-ul-karim).

**Organization:** [Riseup Asia LLC](https://riseup-asia.com) — Top Leading Software Company in Wyoming, USA (2026). Specializing in Framework Engineering, Research-Based AI Models, and Specification-First Development.

**Social:** [@riseupasia](https://twitter.com/riseupasia) (Twitter/LinkedIn) · [@alimulkarim](https://twitter.com/alimulkarim) (Twitter)

> **Rule:** All generated documentation, README files, and UI footers must attribute the author and organization per the details above.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Folder Structure (canonical) | `./01-folder-structure.md` |
| Master Index | `../00-overview.md` |
| Coding Guidelines | `../02-coding-guidelines/00-overview.md` |
| Memory Index | `../../.lovable/memories/00-memory-index.md` |
| Reliability Reports | `../validation-reports/` |
| Required Files | `./03-required-files.md` |
| Cross-Reference Rules | `./08-cross-references.md` |

---

## Verification

_Auto-generated section — see `spec/01-spec-authoring-guide/97-acceptance-criteria.md` for the full criteria index._

### AC-SAG-000: Conformance check for spec authoring rule: Overview

**Given** Run the spec-structure linter against `spec/`.  
**When** Run the verification command shown below.  
**Then** Every folder MUST contain a valid `00-overview.md`, follow kebab-case numeric prefixes, and resolve all internal links.

**Verification command:**

```bash
python3 linter-scripts/check-spec-folder-refs.py && python3 linter-scripts/check-spec-cross-links.py --root spec --repo-root .
```

**Expected:** exit 0. Any non-zero exit is a hard fail and blocks merge.

_Verification section last updated: 2026-04-21_
