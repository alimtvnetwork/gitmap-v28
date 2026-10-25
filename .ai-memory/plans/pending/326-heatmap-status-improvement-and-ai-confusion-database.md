# Plan: Heatmap Status Improvement and AI Confusion Database

- **Task Sequence:** 326
- **Task Slug:** `heatmap-status-improvement-and-ai-confusion-database`
- **Spec Path:** `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/`
- **Status:** `COMPLETED`
- **Created At:** 2026-10-25
- **Completed At:** 2026-10-25

## Executive Summary

This plan modernizes GitMap's Heatmap telemetry, fixes directory-boundary commit errors in `gitmap fix all`, humanizes auto-stashing with transparent Git status flags, establishes the `repo-cache` AI Confusion SQLite database with bidirectional JSON synchronization, and purges malformed commit prefixes across AI execution prompts and skills.

## Subtasks

- [x] `01-subtask-spec-and-planning.md` — Author comprehensive architecture, component, and verification specifications.
- [x] `02-subtask-heatmap-status-improvement.md` — Redesign `gitmap heatmap status` with Lipgloss/termpad padding, colors, completed/pending tasks, tree-view command examples, and `--json` format.
- [x] `03-subtask-work-directory-guard-and-safe-stashing.md` — Prevent `gitmap fix all` from committing in non-git parent directories; humanize auto-stashing with `[M]`, `[D]`, `[?]` flags, collapse unchanged repos, and suggest one-liner status inspection.
- [x] `04-subtask-ai-confusion-database-and-rc-commands.md` — Implement normalized SQLite `AiConfusion` table in `repo-cache`, CLI commands (`gitmap rc ai-confusion` `add`, `view`, `complete`, `sync`), and JSON instruction ingestion.
- [x] `05-subtask-commit-format-and-skills-sync.md` — Eradicate malformed commit prefixes in V6 prompt, GitMap skills, and Muse prompts; add defensive deduplication in `cli/cmdcommit/commit_formatter.go`.
- [x] `06-subtask-verification-and-consolidation.md` — Retrospective verification, documentation index updates, and final GitMap atomic commit.

## Completed Deliverables

1. **Heatmap Visual & Task Overhaul:**
   - Lipgloss card layout with top/bottom padding, subtle cyan/purple accents, and active status indicators in `cli/heatmap/heatmap_cmd_render.go`.
   - Task inventory discovery in `cli/heatmap/heatmap_tasks.go` querying completed/pending plan files and `.ai-memory/temp-agents/ai_agents.db`.
   - `--json` telemetry output flag supporting `HeatmapStatusPayload` schema in `cli/heatmap/heatmap_cmd.go`.
   - Copy-paste command tree view for drilling into task details and confusion alerts.

2. **Work Directory Guard & Safe Stashing:**
   - Boundary checks in `cli/cmdfix/fix_nonrepo.go` verifying `.git` presence before staging or committing in workspace folders.
   - Transparent pre-stash file inventory displaying `[M]`, `[D]`, and `[?]` flags in `cli/cmdfix/fix_stashing.go`.
   - Collapsed summary for unchanged repositories (`Updated all other N repositories (no changes)`) with suggested one-liner inspection commands in `cli/cmdpull/pull_efficient_render.go`.

3. **AI Confusion SQLite Engine & CLI Commands:**
   - Normalized `AiConfusion` SQLite schema and persistence store in `cli/store/ai_confusion_schema.go` and `cli/store/ai_confusion_store.go`.
   - CLI dispatch in `cli/cmd/rc_ai_confusion.go` supporting `add`, `view`, `complete`, and `sync` subcommands with JSON ingestion (`--file`).
   - 6-stage safe synchronization and rebase replay engine in `cli/store/ai_confusion_sync.go`.

4. **Skill & Formatter Remediation:**
   - Updated `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`, `.agents/skills/gitmap/skill.md`, and `.agents/skills/muse-master-prompt/skill.md` with true architectural subsystem modules and Title Case prefixes.
   - Defensive commit formatter updates in `cli/cmdcommit/commit_formatter.go` cleansing redundant prefixes and pseudo-module collisions.

5. **Registry Synchronization & Quality Gates:**
   - Spec 326 registered in `02-spec/21-app/readme.md`.
   - Plan 326 registered in `.ai-memory/plans/readme.md`.
   - Traceability matrix reconciled in `02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md`.
   - Relative paths, forbidden strings, and secrets checks verified 100% clean.

## Disjoint File Ownership Matrix

| Worker | Subtasks | Exclusive Owned Files |
| :--- | :--- | :--- |
| **Worker 01** | Subtask 01, Subtask 02 | `02-spec/21-app/326-*/01-architecture-spec.md`, `cli/heatmap/heatmap_cmd_render.go`, `cli/heatmap/heatmap_tasks.go`, `cli/heatmap/heatmap_cmd.go`, `cli/heatmap/types.go` |
| **Worker 02** | Subtask 03, Subtask 04 | `02-spec/21-app/326-*/02-component-and-cli-spec.md`, `cli/cmdfix/fix_nonrepo.go`, `cli/cmdpull/pull_cwd.go`, `cli/cmdfix/fix_casing.go`, `cli/cmdpull/pull_efficient_render.go`, `cli/store/ai_confusion_*.go`, `cli/cmd/rc_ai_confusion.go`, `cli/cmd/special_repos_cli.go` |
| **Worker 03** | Subtask 05, Subtask 06 | `02-spec/21-app/326-*/03-migration-and-verification-spec.md`, `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`, `.agents/skills/gitmap/skill.md`, `.agents/skills/muse-master-prompt/skill.md`, `cli/cmdcommit/commit_formatter.go` |
