# `.lovable/` Folder Structure Guide

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The `.lovable/` directory is the **canonical location** for all AI-readable project knowledge — memory, tasks, suggestions, constraints, and onboarding context. Any AI model reading this folder should be able to fully understand the project's conventions, active work, and hard constraints.

> **CRITICAL:** There is exactly ONE memory folder: `.lovable/memory/`. The variant `.lovable/memories/` is **prohibited** and must be deleted if found. Never maintain two memory folders.

---

## Canonical Folder Structure

```
.lovable/
├── overview.md                     # AI onboarding — read this FIRST
├── user-preferences                # User communication preferences
│
├── memory/                         # Institutional knowledge (patterns, decisions, rules)
│   ├── index.md                    # Canonical index of all memory files
│   ├── architecture/               # System architecture decisions
│   ├── constraints/                # Hard constraints and rules
│   ├── features/                   # Feature-specific knowledge
│   ├── issues/                     # Issue-specific knowledge
│   ├── patterns/                   # Reusable patterns/templates
│   ├── processes/                  # Workflow processes
│   ├── project/                    # Project-level status/decisions
│   ├── standards/                  # Technical standards
│   └── style/                      # Code style rules
│
├── suggestions/                    # Suggestion details (one .md per suggestion)
│   └── completed/                  # Completed suggestions archive
│
├── pending-tasks/                  # Active work-in-progress tasks (one .md per task)
├── completed-tasks/                # Finished tasks archive
│
├── pending-issues/                 # Open issues awaiting resolution
├── solved-issues/                  # Resolved issues with root cause analysis
│
├── strictly-avoid/                 # ⛔ Things the AI must NEVER do (one .md per rule)
│
├── plan.md                         # Current roadmap / active plan
├── suggestions.md                  # Pending suggestions summary (bullet points)
└── strictly-avoid.md               # ⛔ Quick-read summary of all strict avoidance rules
```

---

## File Descriptions

### Root Files

| File | Purpose | Read Priority |
|------|---------|---------------|
| `overview.md` | AI onboarding document — project summary, conventions, how to navigate | **READ FIRST** |
| `user-preferences` | User's communication style, timezone, response format preferences | High |
| `plan.md` | Current active plan — what's being worked on now, next steps | High |
| `suggestions.md` | Bullet-point summary of all pending suggestions (details in `suggestions/`) | Medium |
| `strictly-avoid.md` | Quick-reference list of all strict avoidance rules (details in `strictly-avoid/`) | **READ FIRST** |

### `memory/` — Institutional Knowledge

Stores patterns, architectural decisions, naming conventions, and workflow rules that persist across AI sessions.

- **`index.md`** — Canonical index listing all memory files with paths and one-line descriptions
- Subfolders use **kebab-case WITHOUT numeric prefixes**
- Files use **kebab-case** and MAY have numeric prefixes for ordering

| Subfolder | Purpose | Examples |
|-----------|---------|----------|
| `architecture/` | System design decisions, data pipelines | `database-schema.md`, `split-database.md` |
| `constraints/` | Non-negotiable rules and boundaries | `axios-version-pinning.md` |
| `features/` | Feature-specific implementation knowledge | `self-update-architecture.md` |
| `issues/` | Issue-specific knowledge | `nested-code-fence-data-corruption.md` |
| `patterns/` | Reusable templates and patterns | `spec-template.md` |
| `processes/` | How work gets done — workflows, conventions | `development-workflow.md` |
| `project/` | Project-level status and decisions | `documentation-standards.md` |
| `standards/` | Technical standards and enforcement rules | `code-red-guidelines.md`, `enum-standards.md` |
| `style/` | Code style and naming rules | `naming-conventions.md`, `powershell-naming.md` |

### `suggestions/` — Suggestion Tracking

- Each suggestion gets its own `.md` file with full detail
- `completed/` subfolder archives implemented suggestions
- `suggestions.md` (root) provides a bullet-point summary of **pending** suggestions only

### `pending-tasks/` and `completed-tasks/` — Task Lifecycle

- Each task is a single `.md` file with description, status, and notes
- When a task is finished, move it from `pending-tasks/` to `completed-tasks/`
- `plan.md` (root) provides the high-level roadmap linking to active tasks

### `pending-issues/` and `solved-issues/` — Issue Lifecycle

- Each issue is a single `.md` file with root cause, steps to reproduce, and resolution
- When resolved, move from `pending-issues/` to `solved-issues/`

### `strictly-avoid/` — Hard Prohibitions ⛔

This folder contains things the AI must **NEVER** do. Each rule gets its own `.md` file with:

- **What** is prohibited
- **Why** it's prohibited
- **What to do instead**

Example files:
- `no-memories-folder.md` — Never create `.lovable/memories/` (use `memory/` only)
- `no-error-swallowing.md` — Never catch and ignore errors
- `no-uuid-primary-keys.md` — Never use UUIDs for primary keys

The root `strictly-avoid.md` file provides a **quick-read summary** of all rules so the AI can internalize them without reading every individual file.

---

## Naming Conventions

### Folders

- **kebab-case WITHOUT numeric prefixes**
- Examples: `memory/`, `pending-tasks/`, `strictly-avoid/`

### Files

- **kebab-case**, numeric prefix optional for ordering
- Examples: `index.md`, `01-plan-tracker.md`, `database-schema.md`

---

## AI Reading Order

When an AI model first encounters this project, it should read `.lovable/` files in this order:

1. **`overview.md`** — Understand the project
2. **`strictly-avoid.md`** — Know what NOT to do
3. **`user-preferences`** — Adapt communication style
4. **`memory/index.md`** — Survey all institutional knowledge
5. **`plan.md`** — Understand current work context
6. **`suggestions.md`** — See pending improvement ideas
7. **Individual memory files** — Deep-dive as needed per task

---

## Task & Suggestion Workflow

### Tasks

```
plan.md (roadmap) → pending-tasks/ → completed-tasks/
```

1. Define task in `plan.md`
2. Create detailed `.md` file in `pending-tasks/`
3. Move to `completed-tasks/` when done

### Suggestions

```
suggestions.md (summary) → suggestions/ (details) → suggestions/completed/
```

1. Add bullet point to `suggestions.md`
2. Create detailed `.md` file in `suggestions/`
3. Move to `suggestions/completed/` when implemented
4. Remove from `suggestions.md`

---

## Relationship to `spec/`

| Aspect | `spec/` | `.lovable/` |
|--------|---------|-------------|
| Purpose | Formal specifications | Institutional knowledge & AI context |
| Naming | Numeric prefix **required** | Numeric prefix **optional** |
| Required files | `00-overview.md`, `99-consistency-report.md` | `overview.md`, `memory/index.md` |
| Health scoring | ✅ Tracked by dashboard | ❌ Not scored |
| AI consumption | Referenced during implementation | Referenced during **all** interactions |
| Depth | Up to 3 levels | Up to 2 levels (category/file) |

---

## Consolidation Rule

> **There is only ONE memory folder: `.lovable/memory/`.** The legacy `.lovable/memories/` variant is prohibited. If found, its contents must be migrated to `.lovable/memory/` and the folder deleted.

---

*`.lovable/` folder structure guide — v3.2.0 — 2026-04-16*
