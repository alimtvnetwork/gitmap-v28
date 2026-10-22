# Architecture Specification: Plan 293 — GitMap Plans Command and Fleet Skills Synchronization

## Status: APPROVED
- **Spec ID:** 293-01
- **Task Slug:** `293-gitmap-plans-command-and-fleet-sync`
- **Component Spec:** [02-component-spec.md](02-component-spec.md)
- **Migration & Verification Spec:** [03-migration-and-verification-spec.md](03-migration-and-verification-spec.md)
- **Master Plan:** [.ai-memory/plans/pending/293-gitmap-plans-command-and-fleet-sync.md](../../../.ai-memory/plans/pending/293-gitmap-plans-command-and-fleet-sync.md)
- **Execution Mandate:** Zero Builds, Zero Tests, Strictly Relative Lowercase Git Paths

---

## 1. Problem Statement & Architectural Vision

### 1.1 The Problem
Software engineers and autonomous AI agents working in polyglot multi-repository workspaces need visibility into pending and completed architectural plans stored in `.ai-memory/plans/`. Currently, inspecting plans requires manual filesystem commands such as:
```powershell
Get-ChildItem -Path ".ai-memory/plans" -Name
```
This produces an unformatted, unparsed list of filenames without indicating completion status, chronological order, active tasks, or pending backlogs.

### 1.2 The Solution
Introduce the first-class `gitmap plans` command (with alias `gitmap plan`). The command autonomously scans:
1. `.ai-memory/plans/` (root directory containing completed and top-level master plans)
2. `.ai-memory/plans/pending/` (directory containing pending, queued, and in-progress plans)
3. `.ai-memory/plans/completed/` (directory containing archived completed plans, if present)

It parses markdown headers, extracts sequence IDs, titles, normalized completion states, and renders an Tokyo Night styled terminal presentation or machine-readable JSON/YAML output.

```
┌─────────────────────────────────────────────────────────────────────────────┐
│ gitmap plans [N] [--json] [--yaml] [--all] [--pending] [--completed]        │
├─────────────────────────────────────────────────────────────────────────────┤
│ • Parses sequence ID, Title, Status, Relative File Path                     │
│ • Displays last N completed plans (default 10) in Tokyo Night Lipgloss table │
│ • Displays top N pending plans (default 10) with status pills               │
│ • Renders pending backlog overflow notification with exact CLI command      │
│ • Supports pure machine-readable --json and --yaml streaming to os.Stdout    │
│ • Integrates into GitMap skill and synchronizes across all fleet repos      │
└─────────────────────────────────────────────────────────────────────────────┘
```

---

## 2. Architectural Principles & Coding Invariants

### 2.1 Error Sovereignty
- All leaf functions in `cli/cmdplans/` return `error` or `*apperror.AppError`.
- No leaf function invokes `os.Exit` or `cliexit.HandleError` directly. Errors bubble up to `cli/cmd/root.go` dispatchers where `finishCommandAudit` captures execution metrics and failures into the SQLite audit store.

### 2.2 Function Sizing & Cyclomatic Simplicity
- All functions strictly adhere to the $\le 15$ line rule ($\le 8$ lines preferred).
- Zero nested `if` statements. Complex branches are decomposed into linear guard clauses with early returns.

### 2.3 Positive Boolean Conventions
- All boolean variables and struct fields must begin with `is` or `has` (`isJSON`, `isYAML`, `isAll`, `hasMorePending`, `hasMoreCompleted`).
- Banned prefixes: `can`, `should`, `was`, `will`, and negative boolean names.

### 2.4 Modularity & File Size Cap Ratchet
- Total lines per file strictly bounded by the 500-line cap ratchet (all modules targeting $\le 120$ lines).
- Module decomposition:
  - `cli/cmdplans/plans_types.go`: Data structures, enums, options, and JSON/YAML models.
  - `cli/cmdplans/plans_parser.go`: Filesystem scanner and markdown metadata extractor.
  - `cli/cmdplans/plans_render.go`: Tokyo Night Lipgloss tables and metric pill renderers.
  - `cli/cmdplans/plans_format.go`: Machine-readable JSON and YAML serializers.
  - `cli/cmdplans/plans.go`: Command entry point, argument parsing, and workflow coordinator.

### 2.5 Terminal Stream Separation
- When `--json` or `--yaml` is specified, all output is written exclusively to `os.Stdout` without ANSI escape sequences, diagnostic banners, or extraneous text to allow pipe composition (`gitmap plans --json | jq .`).
- Interactive notices and table framing are restricted to human TTY invocations.

---

## 3. Fleet Synchronization Strategy

1. **Local Skill Enrichment**: Update `.agents/skills/gitmap/SKILL.md` to document the new `gitmap plans` command, flag taxonomy, and replacement for manual directory inspection.
2. **Upstream Canonical Enrichment**: Mirror skill changes to `d:/work/coding-guidelines/.agents/skills/gitmap/SKILL.md`.
3. **Fleet Broadcast**: Execute `python 03-ai-scripts/38-sync-prompts-skills-scripts.py` to propagate skills, rules, and scripts across connected fleet repositories.
