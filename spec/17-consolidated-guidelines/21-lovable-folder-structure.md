# Consolidated: `.lovable/` Folder Structure

**Version:** 3.3.0  
**Updated:** 2026-04-16  
**Source:** [`spec/01-spec-authoring-guide/07-memory-folder-guide.md`](../01-spec-authoring-guide/07-memory-folder-guide.md)

---

## Purpose

This is the **standalone consolidated reference** for the `.lovable/` folder structure — the AI context layer. An AI reading only this file must be able to create, maintain, and navigate the `.lovable/` directory correctly.

---

## Canonical Structure

```
.lovable/
├── overview.md                     # AI onboarding — read FIRST
├── user-preferences                # User communication preferences
├── plan.md                         # Current roadmap / active plan
├── suggestions.md                  # Pending suggestions (bullet points)
├── strictly-avoid.md               # ⛔ Quick-read prohibition summary
│
├── memory/                         # Institutional knowledge
│   ├── index.md                    # Canonical index of all memory files
│   ├── architecture/               # System architecture decisions
│   ├── constraints/                # Hard constraints and rules
│   ├── done/                       # Completed tasks archive
│   ├── features/                   # Feature-specific knowledge
│   ├── issues/                     # Issue-specific knowledge
│   ├── patterns/                   # Reusable patterns/templates
│   ├── processes/                  # Workflow processes
│   ├── project/                    # Project-level status/decisions
│   ├── standards/                  # Technical standards
│   ├── style/                      # Code style rules
│   ├── suggestions/                # Suggestion tracker
│   └── workflow/                   # Workflow trackers
│
├── suggestions/                    # Suggestion details (one .md per suggestion)
│   └── completed/                  # Completed suggestions archive
│
├── pending-tasks/                  # Active work-in-progress tasks
├── completed-tasks/                # Finished tasks archive
│
├── pending-issues/                 # Open issues
├── solved-issues/                  # Resolved issues with root cause
│
└── strictly-avoid/                 # ⛔ Hard prohibitions (one .md per rule)
```

---

## Critical Rules

> **There is exactly ONE memory folder: `.lovable/memory/` (singular).** The variant `.lovable/memories/` (plural) is **prohibited**. If found, migrate contents and delete it.

> **`memory/index.md` is the single source of truth** for all memory files. Every memory file must be listed there. Orphaned files (in `memory/` but not in `index.md`) must be indexed or removed.

---

## AI Reading Order

1. `overview.md` → understand the project
2. `strictly-avoid.md` → know what NOT to do
3. `user-preferences` → adapt communication style
4. `memory/index.md` → survey all institutional knowledge
5. `plan.md` → understand current work context
6. `suggestions.md` → see pending ideas

---

## Naming Conventions

- **Folders:** kebab-case, no numeric prefixes (`memory/`, `pending-tasks/`)
- **Files:** kebab-case, optional numeric prefix (`index.md`, `01-plan-tracker.md`)
- **No spaces**, no underscores, no PascalCase in filenames

---

## Memory File Conventions

### File Structure

Every memory file follows this template:

```markdown

# Memory: [descriptive-title]

**Category:** architecture | constraints | features | issues | patterns | processes | project | standards | style | workflow  
**Created:** YYYY-MM-DD  
**Updated:** YYYY-MM-DD

---

## Summary

One-paragraph description of what this memory captures.

## Details

Full content — rules, decisions, patterns, constraints.

## Related

- Links to spec files, other memories, or code files
```

### Categories and When to Use

| Category | Use When | Examples |
|----------|----------|---------|
| `architecture/` | System-level structural decisions | Database schema, caching policy, split-DB pattern |
| `constraints/` | Hard rules that cannot be violated | Version pinning, forbidden patterns |
| `features/` | Feature-specific knowledge and requirements | Self-update architecture, visual rendering |
| `issues/` | Specific bugs and their root cause analysis | Nested code fence corruption |
| `patterns/` | Reusable code/design patterns | Template patterns, composition patterns |
| `processes/` | Workflow and operational procedures | Development workflow, automated enforcement |
| `project/` | Project-level status and decisions | Documentation standards, author attribution |
| `standards/` | Technical standards and coding rules | Code Red guidelines, TypeScript patterns, enum standards |
| `style/` | Code style and naming conventions | Naming conventions, PowerShell naming |
| `workflow/` | Workflow tracking state | Sprint trackers, migration status |

### index.md Structure

The `memory/index.md` has two sections:

```markdown

# Project Memory

## Core

- 🔴 CODE RED: [critical rule — one line, <150 chars]
- [universal rule that applies to every action]

## Memories

- [Descriptive title](mem://category/filename) — One-line description
```

**Core section rules:**
- Only rules that apply to **every** action across the entire project
- Max ~150 characters per entry
- Prefixed with 🔴 for CODE RED severity items

**Memories section rules:**
- Every memory file in `memory/` must have a corresponding entry
- Format: `[Title](mem://category/filename) — Description`
- Description must be specific enough to judge relevance without opening the file
- Sort by category, then alphabetically within category

### Creating a New Memory

1. Write the file to `memory/{category}/{kebab-case-name}.md`
2. Add entry to `memory/index.md` under `## Memories`
3. Verify no duplicate or overlapping memory already exists

### Updating a Memory

1. Edit the file directly
2. Update the `Updated:` date in the frontmatter
3. If the description changed materially, update the index entry too

### Deleting a Memory

1. Remove the file from `memory/{category}/`
2. Remove its entry from `memory/index.md`

---

## Workflows

### Tasks: `plan.md` → `pending-tasks/` → `completed-tasks/`

1. High-level items tracked in `plan.md` as a roadmap
2. When work begins, create a detailed file in `pending-tasks/`
3. On completion, move to `completed-tasks/` with results noted

### Task File Template

```markdown

# Pending Task: [Title]

**Priority:** High | Medium | Low  
**Status:** Not Started | In Progress | Complete ✅  
**Created:** YYYY-MM-DD

---

## Description

What needs to be done.

## Items

- [ ] Subtask 1
- [ ] Subtask 2
- [x] Completed subtask

---

*Pending task — [context] — vX.Y.Z — YYYY-MM-DD*
```

### Suggestions: `suggestions.md` → `suggestions/` → `suggestions/completed/`

1. Quick bullet points in `suggestions.md` for pending ideas
2. Detailed analysis in `suggestions/{suggestion-name}.md`
3. On implementation, move to `suggestions/completed/`

### Issues: `pending-issues/` → `solved-issues/`

1. Document the issue with symptom, diagnosis, and attempted fixes
2. On resolution, move to `solved-issues/` with root cause and fix documented

### Issue File Template

```markdown

# Issue: [Title]

**Severity:** Critical | High | Medium | Low  
**Status:** Open | Investigating | Resolved  
**Created:** YYYY-MM-DD

---

## Symptom

What was observed.

## Root Cause

Why it happened (filled after diagnosis).

## Fix

What was changed and why.

## Prevention

How to prevent recurrence.
```

---

## Strictly-Avoid Rules

### Summary File (`strictly-avoid.md`)

Quick-reference list of all prohibitions. One line per rule, linking to the detailed file.

### Detail Files (`strictly-avoid/{rule-name}.md`)

```markdown

# Strictly Avoid: [Rule Title]

**Rule:** One-sentence prohibition.

---

## What Is Prohibited

Specific behavior or pattern that must not occur.

## Why

Rationale — what goes wrong if violated.

## What To Do Instead

The correct alternative approach.
```

---

## Root-Level Files

### `overview.md`

Project onboarding for AI. Must contain:
- Project name and purpose
- Tech stack summary
- Key architecture decisions
- Links to critical spec modules

### `user-preferences`

User communication preferences (plain text, no frontmatter). Applied to every AI response. Examples:
- Preferred language/tone
- Timezone
- Response format preferences
- Things to always/never do

### `plan.md`

Current roadmap with active items. Updated as priorities shift.

### `suggestions.md`

Bullet-point list of pending suggestions with one-line descriptions.

---

*Consolidated .lovable folder structure — v3.3.0 — 2026-04-16*

---

## §X Project Memory — Active Core Rules (Mirror)

This section **mirrors** the operational rules stored in `.lovable/memory/index.md` Core section. A blind AI receiving only the consolidated folder would otherwise miss these — and violate at least three on its first PR. This mirror is **read-only documentation** of the rules; the canonical source remains `mem://index.md`.

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
| Spec changes | Spec-First — edit `spec/` then implement |
| Bug fixes | Issue-First — create `03-issues/<issue>.md` then fix |
| `.lovable/` structure | Single-file convention — `plan.md`, `suggestions.md`, `strictly-avoid.md` each hold their full history. **Never** create per-task folders |
| Multi-step requests | Break into discrete tasks. Wait for "next" prompt to continue |

### X.7 Dependency Pinning

| Package | Allowed Versions | Blocked |
|---------|------------------|---------|
| `axios` | `1.14.0`, `0.30.3` | `1.14.1`, `0.30.4` (security; never bump) |

Enforced by `linter-scripts/check-axios-version.sh`.

### X.8 Why This Mirror Exists

The canonical memory lives at `mem://index.md` and is automatically loaded into every Lovable AI session. **External AIs** (Claude, GPT, Gemini handed only this folder) have no access to that memory — this section is the only way they will learn these rules.

When updating: edit `mem://index.md` first, then sync this section. The mirror is allowed to lag by at most one minor version.

---

*Project Memory Core mirror added — v3.4.0 — 2026-04-22*
