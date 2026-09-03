# Consolidated: `.lovable/` Folder Structure

**Version:** 4.0.0
**Updated:** 2026-08-31
**Source:** [`02-spec/01-spec-authoring-guide/09-memory-folder-guide.md`](../01-spec-authoring-guide/09-memory-folder-guide.md)

---

## Purpose

This is the **standalone consolidated reference** for the `.lovable/` folder structure — the AI context layer. An AI reading only this file must be able to create, maintain, and navigate the `.lovable/` directory correctly.

---

## Canonical Structure

```
.lovable/
├── 01-overview.md                   # AI onboarding — read FIRST
├── 02-user-preferences              # User communication preferences
├── 03-strictly-avoid.md             # ⛔ Quick-read prohibition summary
├── 04-suggestions.md                # Pending suggestions (bullet points)
├── 05-plan.md                       # Current roadmap / active plan
├── 06-what-to-read.md               # Router and reading sequence
│
├── ai-fix-scripts/                  # Persistent AI automation toolchain
│   ├── 01-index.md                  # Master catalog & search tag registry
│   ├── 02-shared-engine.py          # Shared engine: constants, lazy regex, cache
│   └── 03..20-*.py                  # Specialized linters and runners
│
├── plans/                           # Execution hub
│   ├── 01-index.md                  # Master plan index
│   ├── pending/                     # Active parent task specs
│   ├── subtasks/                    # Bounded micro-tasks (XX-<slug>/)
│   └── completed/                   # Archived completed plans
│
├── memory/                          # Institutional knowledge (SINGULAR)
│   ├── 01-index.md                  # Canonical index of all memory files
│   ├── architecture/                # System architecture decisions
│   ├── constraints/                 # Hard constraints and rules
│   ├── done/                        # Completed tasks archive
│   ├── features/                    # Feature-specific knowledge
│   ├── issues/                      # Issue-specific knowledge
│   ├── patterns/                    # Reusable patterns/templates
│   ├── processes/                   # Workflow processes
│   ├── project/                     # Project-level status/decisions
│   ├── standards/                   # Technical standards
│   ├── style/                       # Code style rules
│   ├── suggestions/                 # Suggestion tracker
│   └── workflow/                    # Workflow trackers
│
├── prompts/                         # AI Prompt Repository
│   ├── 01-prompts-category/         # Categorized source prompt modules (01-22)
│   └── *.md                         # Flat synced prompts
│
├── release/                         # Release automation & version bumping
│   ├── release-method.md            # Version bump specification
│   ├── bump_versions.py             # Version bumper
│   └── issues/                      # Release diagnostics
│
├── question-and-ambiguity/          # Ambiguity logs & iteration counters
├── suggestions/                     # Granular suggestion proposals
├── cicd-issues/                     # CI pipeline diagnostics & RCAs
└── assets/                          # Mockups, diagrams, media
```

---

## Critical Rules

> **There is exactly ONE memory folder: `.lovable/memory/` (singular).** The variant `.lovable/memories/` (plural) is **prohibited**. If found, migrate contents and delete it.

> **`memory/01-index.md` is the single source of truth** for all memory files. Every memory file must be listed there. Orphaned files (in `memory/` but not in `index.md`) must be indexed or removed.

---

## AI Reading Order

1. `01-overview.md` → understand the project
2. `03-strictly-avoid.md` → know what NOT to do
3. `02-user-preferences` → adapt communication style
4. `memory/01-index.md` → survey all institutional knowledge
5. `05-plan.md` → understand current work context
6. `04-suggestions.md` → see pending ideas

---

## Naming Conventions

- **Folders:** kebab-case, 2-digit zero-padded prefix when sequenced (`01-prompts-category/`, `ai-fix-scripts/`)
- **Files:** strictly lowercase, kebab-case, numeric prefix where sequenced (`01-index.md`, `02-shared-engine.py`)
- **No spaces**, no uppercase letters, no camelCase in filenames.

---

## Workflows

### Tasks: `05-plan.md` → `plans/pending/` → `plans/completed/`

1. High-level items tracked in `05-plan.md` as a roadmap.
2. When work begins, create a detailed file in `plans/pending/` and decompose into `plans/subtasks/XX-<slug>/`.
3. On completion, move to `plans/completed/` with results noted.

---

*Consolidated .lovable folder structure — v4.0.0 — 2026-08-31*
