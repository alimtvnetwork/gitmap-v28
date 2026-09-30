# Consolidated: `.ai-memory/` Folder Structure

**Version:** 4.0.0
**Updated:** 2026-08-31
**Source:** [`02-spec/01-spec-authoring-guide/09-memory-folder-guide.md`](../01-spec-authoring-guide/09-memory-folder-guide.md)

---

## Purpose

This is the **standalone consolidated reference** for the `.ai-memory/` folder structure — the AI context layer. An AI reading only this file must be able to create, maintain, and navigate the `.ai-memory/` directory correctly.

---

## Canonical Structure

```
.ai-memory/
├── 01-overview.md                   # AI onboarding — read FIRST
├── 02-user-preferences              # User communication preferences
├── 03-strictly-avoid.md             # ⛔ Quick-read prohibition summary
├── 04-suggestions.md                # Pending suggestions (bullet points)
├── 05-plan.md                       # Current roadmap / active plan
├── 06-what-to-read.md               # Router and reading sequence
│
├── ai-fix-scripts/                  # Persistent AI automation toolchain
│   ├── readme.md                  # Master catalog & search tag registry
│   ├── 02-shared-engine.py          # Shared engine: constants, lazy regex, cache
│   └── 03..20-*.py                  # Specialized linters and runners
│
├── plans/                           # Execution hub
│   ├── readme.md                  # Master plan index
│   ├── pending/                     # Active parent task specs
│   ├── subtasks/                    # Bounded micro-tasks (XX-<slug>/)
│   └── completed/                   # Archived completed plans
│
├── memory/                          # Institutional knowledge (SINGULAR)
│   ├── readme.md                  # Canonical index of all memory files
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

> **There is exactly ONE memory folder: `.ai-memory/memory/` (singular).** The variant `.ai-memory/memories/` (plural) is **prohibited**. If found, migrate contents and delete it.

> **`memory/readme.md` is the single source of truth** for all memory files. Every memory file must be listed there. Orphaned files (in `memory/` but not in `index.md`) must be indexed or removed.

---

## AI Reading Order

1. `01-overview.md` → understand the project
2. `03-strictly-avoid.md` → know what NOT to do
3. `02-user-preferences` → adapt communication style
4. `memory/readme.md` → survey all institutional knowledge
5. `05-plan.md` → understand current work context
6. `04-suggestions.md` → see pending ideas

---

## Naming Conventions

- **Folders:** kebab-case, 2-digit zero-padded prefix when sequenced (`01-prompts-category/`, `ai-fix-scripts/`)
- **Files:** strictly lowercase, kebab-case, numeric prefix where sequenced (`readme.md`, `02-shared-engine.py`)
- **No spaces**, no uppercase letters, no camelCase in filenames.

---

## Workflows

### Tasks: `05-plan.md` → `plans/pending/` → `plans/completed/`

1. High-level items tracked in `05-plan.md` as a roadmap.
2. When work begins, create a detailed file in `plans/pending/` and decompose into `plans/subtasks/XX-<slug>/`.
3. On completion, move to `plans/completed/` with results noted.

---

*Consolidated .ai-memory folder structure — v4.0.0 — 2026-08-31*

---

## §X Project Memory — Active Core Rules (Mirror)

This section **mirrors** the operational rules stored in `.ai-memory/memory/readme.md` Core section. A blind AI receiving only the consolidated folder would otherwise miss these — and violate at least three on its first PR. This mirror is **read-only documentation** of the rules; the canonical source remains `mem://index.md`.

### X.1 Code-Red Quality Rules

| Rule | Enforced By |
|------|-------------|
| Never swallow errors. Zero-nesting (no nested `if`). Max 2 operands per condition. Positively named guard functions. | `linter-scripts/validate-guidelines.py` |
| Functions: 8–15 lines. Files: < 300 lines. React components: < 100 lines. | `linter-scripts/validate-guidelines.py` |

### X.2 Sync & Repo Rules

| Rule | Notes |
|------|-------|
| **Never** sync `01-app`, `02-app-issues`, `03-general`, `03-tasks`, or `12-consolidated-guidelines` from upstream sibling repos | All maintained locally |
| **Skip** from spec audits: `21-app`, `22-app-issues`, `23-app-db`, `24-app-ui-design-system` are intentional stubs | Never write 97/99 files for them; never demote to `_drafts/`; exclude from corpus averages |
| Repo identity: `alimtvnetwork/coding-guidelines-v24` | Install scripts live at repo root (`install.ps1` / `install.sh`) |

### X.3 Naming Rules

| Domain | Convention | Exception |
|--------|------------|-----------|
| Internal IDs, DB, JSON, Types | PascalCase | Rust uses `snake_case` identifiers |
| DB tables | PascalCase, **singular** | — |
| DB primary keys | `{TableName}Id` (INTEGER PRIMARY KEY AUTOINCREMENT) | No UUIDs |

### X.4 DB Boolean Rules

- **Forbidden** prefixes: `Not`, `No`
- **Approved Inverses** (allowed despite negative semantics): `IsDisabled`, `IsInvalid`, `IsIncomplete`, `IsUnavailable`, `IsUnread`, `IsHidden`, `IsBroken`, `IsLocked`, `IsUnpublished`, `IsUnverified`
- Inverses are derived in code via Rule 9 codegen (never stored as separate columns)

### X.5 DB Descriptive Column Rules (Rules 10/11/12)

| Table Type | Required Columns |
|------------|------------------|
| Entity tables | `Description TEXT NULL` |
| Transactional tables | `Notes TEXT NULL` + `Comments TEXT NULL` |

Enforcement: see `18-database-conventions.md` §18 (rule presence) and §19 (waiver syntax).

### X.6 Workflow Rules

| Rule | Pattern |
|------|---------|
| Spec changes | Spec-First — edit `02-spec/` then implement |
| Bug fixes | Issue-First — create `03-issues/<issue>.md` then fix |
| `.ai-memory/` structure | Single-file convention — `plan.md`, `suggestions.md`, `strictly-avoid.md` each hold their full history. **Never** create per-task folders |
| Multi-step requests | Break into discrete tasks. Wait for "next" prompt to continue |

### X.7 Dependency Pinning

| Package | Allowed Versions | Blocked |
|---------|------------------|---------|
| `axios` | `1.14.0`, `0.30.3` | `1.14.1`, `0.30.4` (security; never bump) |

Enforced by `linter-scripts/check-axios-version.sh`.

### X.8 Why This Mirror Exists

The canonical memory lives at `mem://index.md` and is automatically loaded into every Lovable AI session. External AIs have no access to that memory — this section is the only way they will learn these rules.

When updating: edit `mem://index.md` first, then sync this section. The mirror is allowed to lag by at most one minor version.

