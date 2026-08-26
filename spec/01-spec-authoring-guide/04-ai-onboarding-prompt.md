# AI Project Onboarding Protocol

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

> **Purpose:** This document is a mandatory onboarding sequence for any AI assistant joining this project. It ensures you internalize all specifications, rules, and conventions before writing a single line of code.

> **Rule #0:** Follow every phase sequentially. Do not skip, summarize prematurely, or assume knowledge from training data. The specs are the single source of truth.

---

## Table of Contents

1. [Phase 1 — AI Context Layer](#phase-1--ai-context-layer)
2. [Phase 2 — Consolidated Guidelines](#phase-2--consolidated-guidelines)
3. [Phase 3 — Spec Authoring Rules](#phase-3--spec-authoring-rules)
4. [Phase 4 — Deep-Dive Source Specs](#phase-4--deep-dive-source-specs-task-driven)
5. [Anti-Hallucination Contract](#anti-hallucination-contract)
6. [Memory Update Protocol](#memory-update-protocol)
7. [Completion Confirmation](#completion-confirmation)

---

## Phase 1 — AI Context Layer

**Goal:** Load the project's identity, hard rules, and institutional memory into your working context.

### Step 1.1 — Read core files in EXACT order

| Order | File | What You Learn |
|-------|------|----------------|
| 1 | `.lovable/overview.md` | Project summary, tech stack, navigation map |
| 2 | `.lovable/strictly-avoid.md` | **Hard prohibitions** — violating ANY of these is a critical failure |
| 3 | `.lovable/user-preferences` | How the human expects you to communicate and behave |
| 4 | `.lovable/memory/index.md` | Index of all institutional knowledge files |
| 5 | `.lovable/plan.md` | Current active roadmap and priorities |
| 6 | `.lovable/suggestions.md` | Pending improvement ideas (not yet approved) |

### Step 1.2 — Read EVERY file referenced in `.lovable/memory/index.md`

- If the index lists 12 files, you read 12 files. No exceptions.
- If there are subfolders, traverse them recursively.
- If a file is missing or empty, note it — do not silently skip.

### Step 1.3 — Self-check (answer these internally before continuing)

- [ ] What are the project's **CODE RED** rules?
- [ ] What naming conventions are enforced (files, folders, DB columns, variables)?
- [ ] What is the error handling philosophy?
- [ ] What is the current plan and what tasks are in progress?
- [ ] What patterns/tools/approaches are **strictly forbidden**?

> ⛔ **DO NOT proceed to Phase 2 until every file above has been read and internalized.**

---

## Phase 2 — Consolidated Guidelines

**Goal:** Absorb the project's unified rulebook — 22 self-contained guideline documents.

### Instructions

1. Navigate to `spec/17-consolidated-guidelines/`.
2. Read files in **numeric order**: `01-*.md` through `22-*.md`.
3. Each file is self-contained. Treat each as a standalone policy document.

### After reading, confirm internally

- [ ] Total number of guideline files read.
- [ ] One-sentence summary of the key rule from each file.
- [ ] Any rules that contradict your default training (these are intentional — the spec wins).

> ⛔ **DO NOT proceed to Phase 3 until all guideline files have been read.**

---

## Phase 3 — Spec Authoring Rules

**Goal:** Understand how specifications themselves are structured, so you can read them correctly and author new ones if asked.

### Instructions

1. Navigate to `spec/01-spec-authoring-guide/`.
2. Read all files in numeric order.

### After reading, confirm you understand

| Concept | Where It's Defined |
|---------|-------------------|
| File and folder naming conventions | Spec authoring guide |
| Required files in every spec folder (`00-overview.md`, `99-consistency-report.md`) | Spec authoring guide |
| The `.lovable/` folder structure and its purpose | `07-memory-folder-guide.md` |
| Linter infrastructure requirements | Spec authoring guide |

> ⛔ **DO NOT begin any task until Phases 1–3 are complete.**

---

## Phase 4 — Deep-Dive Source Specs (Task-Driven)

**Goal:** Before performing any task, read the relevant source spec(s) so your work is compliant.

### Lookup Table

| If your task involves... | Read this spec folder |
|--------------------------|----------------------|
| Writing or reviewing code | `spec/02-coding-guidelines/` |
| Error handling | `spec/03-error-manage/` |
| Database schema or queries | `spec/04-database-conventions/` |
| SQLite or multi-database architecture | `spec/05-split-db-architecture/` |
| Configuration systems | `spec/06-seedable-config-architecture/` |
| UI theming, CSS variables, design tokens | `spec/07-design-system/` |
| Documentation viewer features | `spec/08-docs-viewer-ui/` |
| Code block rendering | `spec/09-code-block-system/` |
| PowerShell scripts | `spec/11-powershell-integration/` |
| CI/CD pipelines | `spec/12-cicd-pipeline-workflows/` |
| CLI self-update system | `spec/14-update/` |
| WordPress plugins | `spec/18-wp-plugin-how-to/` |
| App-specific features | `spec/21-app/` |
| Known app bugs/issues | `spec/22-app-issues/` |
| App-specific database schema | `spec/23-app-db/` |
| App-specific UI and design system | `spec/24-app-ui-design-system/` |

### Reading order within each folder

1. `00-overview.md` — always first
2. All numbered files in order
3. `99-consistency-report.md` — always last (if present)

---

## Anti-Hallucination Contract

These rules are **absolute and non-negotiable**. Violating any of them is a critical failure.

### 1. Never Invent Rules

If a spec does not mention a rule, that rule does not exist. Do not fill gaps with assumptions from your training data.

### 2. Specs Override Training Data

If your pre-trained knowledge conflicts with a spec, **the spec wins**. Every time. No exceptions.

### 3. Cite Your Sources

When enforcing a rule, reference the **specific file and section**. Example:

> Per `spec/02-coding-guidelines/03-naming.md` § "Database Columns": all column names use PascalCase.

### 4. Ask When Uncertain

If a spec is ambiguous or silent on a topic, **ask the human**. Do not guess, infer, or "use best judgment."

### 5. Never Merge Conventions

This project has its own conventions (e.g., PascalCase DB columns). Do not blend them with conventions from other projects, languages, or frameworks you've seen in training.

### 6. Namespace Awareness

The project namespace is `<owner>/<repo>`. Any `v1` reference is a bug. All specs are at v3.2.0.

### 7. No Filler

Never append boilerplate like "Let me know if you have questions!" or "Hope this helps!" Just deliver the work.

---

## Memory Update Protocol

When you learn something new during a session, follow this decision tree:

```
New information discovered
│
├─ Is it institutional knowledge (pattern, convention, decision)?
│  └─ YES → Write to `.lovable/memory/` and update `.lovable/memory/index.md`
│
├─ Is it something that must NEVER be done?
│  └─ YES → Add to `.lovable/strictly-avoid.md`
│
├─ Is it a suggestion or improvement idea (not yet approved)?
│  └─ YES → Add to `.lovable/suggestions.md`
│
└─ None of the above → Do not persist it
```

### Critical Rules

- The memory folder is `.lovable/memory/` — **never** `.lovable/memories/` (no trailing `s`).
- When adding a new memory file, **always** update the index at `.lovable/memory/index.md`.
- When modifying an existing memory, preserve all other content — do not truncate or overwrite unrelated entries.

---

## Completion Confirmation

After completing **Phases 1 through 3**, respond with exactly this format:

```
✅ Onboarding complete.
- Memory files read: [X]
- Consolidated guidelines read: [Y]
- Spec authoring files read: [Z]

I understand:
- CODE RED rules: [list the top 3–5]
- Naming conventions: [brief summary]
- Error handling approach: [one sentence]
- Active plan: [current milestone or focus]
- Strict avoidances: [top 3–5 forbidden patterns]

Ready for tasks.
```

Then **stop and wait** for instructions. Do not suggest next steps. Do not ask exploratory questions. Just wait.

---

## Usage Notes

- **Full onboarding** — Give the entire prompt above for a fresh AI session
- **Quick refresh** — For an AI that has already been onboarded, ask it to re-read only `.lovable/strictly-avoid.md` and `.lovable/memory/index.md`
- **Task-specific** — For a focused task, give Phase 1 + the relevant Phase 4 folder only
- **"Read memory"** — Refers to executing this full onboarding protocol (see `.lovable/prompts/01-read-prompt.md`)

---

## Why Multi-Phase?

| Phase | Purpose |
|-------|---------|
| 1 | Establish hard constraints before reading any implementation details |
| 2 | Build a comprehensive mental model from pre-digested summaries |
| 3 | Understand how specs are structured so the AI can create/edit them |
| 4 | Deep-dive only when needed — avoids context window overflow |

---

## Cross-References

- `.lovable/prompts/01-read-prompt.md` — Quick-access copy of this protocol
- `.lovable/prompt.md` — Prompt index referencing all available prompts
- `.lovable/overview.md` — Shorter always-present onboarding doc

---

*AI onboarding prompt — v3.2.0 — 2026-04-16*
