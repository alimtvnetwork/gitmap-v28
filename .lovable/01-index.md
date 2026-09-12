# Repository Context & AI Architecture Index

> **Version:** 4.0.0
> **Target:** `.lovable/01-index.md`
> **Purpose:** Master repository index, directory router, and operational architecture guide for AI agents.

---

## 1. Repository Architecture Overview

This repository is the **GitMap multi-node, repository mapping and distributed workflow engine** written in Go with interactive terminal UI, SQLite split-databases, React documentation frontend, and cross-platform automation tooling. It consists of:

- **GitMap Core Engine (`gitmap/`):** Go CLI subcommands, database engines, macro engine, split databases, commitin workflow, and cluster orchestration.
- **Formal Specifications (`spec/`):** Layered hierarchy of specification modules covering coding guidelines, error management envelopes, database patterns, CLI contracts, and CI/CD pipelines.
- **Persistent AI Fix Toolchain (`03-ai-scripts/`):** High-speed, dual-platform Python automation engine providing deterministic file manipulation, path sanitization, naming validation, and local CI verification (`06-cicd-local-runner.py`).
- **Institutional Memory (`.lovable/memory/`):** Long-term cognitive storage tracking system architecture decisions, non-negotiable constraints, and root cause analyses.
- **Plan & Micro-Task Engine (`.lovable/plans/`):** Bounded task decomposition center (`pending/`, `subtasks/`, `completed/`).
- **Interactive Documentation Viewer (`src/`):** Modern React docs application with live search, tree navigation, and presentation capabilities.

---

## 2. Directory Navigation Router

| Location | Purpose | Key Entrypoint |
|---|---|---|
| `spec/` | Formal cross-language specifications & standards | `spec/01-spec-authoring-guide/00-overview.md` |
| `gitmap/` | Primary Go application codebase & CLI subcommands | `gitmap/main.go` |
| `03-ai-scripts/` | Reusable high-speed Python automation toolchain | `03-ai-scripts/01-index.md` |
| `.lovable/memory/` | Institutional knowledge base & CODE RED constraints | `.lovable/memory/01-index.md` |
| `.lovable/plans/` | Active roadmap, parent task specs & subtasks | `.lovable/plans/01-index.md` |
| `01-prompts/` | Prompt categories and compiled workflow prompts | `.lovable/prompts.md` |
| `.lovable/folder-structure.md` | Canonical folder structure & sequence ID rules | `.lovable/folder-structure.md` |
| `.lovable/strictly-avoid.md` | Universal hard prohibitions & negative constraints | `.lovable/strictly-avoid.md` |
| `.lovable/suggestions.md` | Architecture suggestions & enhancement proposals | `.lovable/suggestions.md` |

---

## 3. Core Universal Rules (CODE RED)

Violations of these principles are blocking:

1. **Never Swallow Errors:** Every error must be explicitly logged, handled, or wrapped in an application error envelope (`*apperror.AppError`).
2. **Zero Nesting (Flatten Conditionals):** Strictly avoid 4+ level nested `if` pyramids. Use shallow guard clauses and fast substring pre-filters.
3. **Implicit Positive Booleans:** NEVER compare booleans explicitly against `true` (`if is_ready:` mandatory; `if is_ready == True:` forbidden).
4. **Mandatory Boolean Prefixes:** All boolean variables, fields, and functions MUST use positive prefixes (`is_`, `has_`, `can_`, `should_`).
5. **Strict Relative Git Paths:** All links, markdown paths, and references must be relative to the repository root. Absolute paths (`C:\...`, `/home/...`, `file:///...`) are strictly prohibited.
6. **Strict Lowercase File Naming:** All files, scripts, and documentation must use lowercase naming (`01-index.md`, `02-shared-engine.py`).
7. **Quality Gate Verification:** Before completing any work session, run `python 03-ai-scripts/06-cicd-local-runner.py` to confirm all quality checks pass.

---

## 4. AI Operational Reading Sequence

1. `.lovable/01-index.md` → Understand repository scope and directory navigation.
2. `.lovable/strictly-avoid.md` → Load hard prohibitions and negative constraints.
3. `.lovable/memory/01-index.md` → Survey institutional knowledge and architectural decisions.
4. `.lovable/plans/01-index.md` → Review active plans and current task bounds.
5. `.lovable/prompts.md` → Access prompt catalog.

---

*Repository Architecture Index — v4.0.0 — 2026-09-12*