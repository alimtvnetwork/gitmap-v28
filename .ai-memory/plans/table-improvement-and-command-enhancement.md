# Master Plan: Table Improvement and Command Enhancement

**Slug:** `table-improvement-and-command-enhancement`  
**Run Number:** 81  
**Target Specification Directory:** `02-spec/21-app/table-improvement-and-command-enhancement/`  
**Subtasks Directory:** `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/`  
**Task DB:** `.ai-memory/temp-agents/81-table-improvement-and-command-enhancement/agent-task.db`  
**Status:** ACTIVE (Phase 1 Synthesis)

---

## User Request (Verbatim)

```text
# Table Improvement and Command Enhancement: high priority instruction, non-negotiable task

I think you need to improve this table. So what I need to know from this table is that the repo name, the short version, and how to fix this. I don't need untracked, modified, staged. Probably I need to have uncommitted. Combining all these things. Each one of them should have a tree view, and the tree would have a solution, like I could run this and make a fix. Overall, there would be a command that would fix every one of them. Make sure that you write the spec, you implement that, you add the help everywhere, and also add the skills section with the help so that any LLM should run this to see any other changing, and from this, they could fix it. These statuses should be cached for one or two minutes. Remember that if the time is expired, the cache should be removed and do not trust it. It should have a cache DB, which actually deals with it. One DB would keep all these backups, so that it can be served again. Include the help text for the new commands that in the last 100 commits has been added. Try to observe that, add it to the help, and also add a new command like new commands. If I do that, then last, let's say, 100 new commands should show up as a filter. Try to include with the examples. Make sure you complete the coding, and at the end, you do the bump in the minor version and release.

## slug: table-improvement-and-command-enhancement

# Actionable Items Must Follow Non-Negotiable

1. Write spec under 02-spec/21-app/table-improvement-and-command-enhancement/ and enqueue plan task in .ai-memory/plans/table-improvement-and-command-enhancement.md (subtasks in .ai-memory/plans/subtasks/table-improvement-and-command-enhancement/) first
2. Search codebase exclusively via GitMap (gitmap aum search, gitmap find, gitmap cat, gitmap ps, gitmap py, gitmap llm train); TOTAL BAN on rg, ripgrep, grep, git grep, Select-String
3. Strictly use relative Git paths (02-spec/..., .ai-memory/..., cmd/...); only add the relative paths, never add the absolute path during your work, and ensure this is respected on the release page and in release notes as well
4. Use `gitmap` AI agents to enter data
5. Task completion includes committing and pushing to Git
6. Improve the table to include repo name, short version, and fix instructions, excluding untracked, modified, staged statuses, but including uncommitted.
7. Implement a tree view for each item with a solution command.
8. Create an overall command to fix all items.
9. Add help text and skills section for LLM compatibility.
10. Implement caching for statuses with a cache DB, expiring after 1-2 minutes.
11. Include help text for new commands from the last 100 commits.
12. Add a filter command to show the last 100 new commands with examples.
13. Complete coding, bump minor version, and release.

## Must follow and spawn agent using

@[.agents/skills/execute-parent-task-with-n-steps-v6]

## Additional Instructions

- /plan first before doing the work to reduce the credits.
- /learn from @[.agents/skills/gitmap] skill to leverage GitMap high-speed search, toolchain discovery, and caching.
- Only add the relative paths, never add the absolute path during your work; this should be respected on the release page and in release notes as well./learn
```

---

## 1. Architectural Synthesis & Objectives

1. **Table View Improvement:**
   - Display: Repository Name (`REPOSITORY`), Short Version & Branch (`VER/BRANCH`), Consolidated Uncommitted Files (`UNCOMMITTED`), Unpushed Commits (`UNPUSHED`), and Status (`STATUS`).
   - Remove individual columns for untracked, modified, staged; collapse into single `UNCOMMITTED` metric.
   - For dirty repositories: render tree-view with two solution commands (`Option 1: commit & push`, `Option 2: stash`).
   - Footer: Render overall fleet fix command (`gitmap cpar "wip: save changes"`).

2. **Dual-Database Status Caching & Backup System:**
   - Primary Cache DB (`pending_commits_cache.db`): strictly enforces 90-second TTL (1.5 min, within 1-2 min range). Stale entries are immediately purged from DB and never trusted. Bypass flags: `--no-cache`, `--refresh`.
   - Secondary Backup DB (`pending_commits_backup.db`): stores historical / snapshot backups of repository statuses with timestamps, branch, version, and payload JSON. Not auto-purged on 90s TTL. Allows serving backed-up statuses again via `--backup` flag.

3. **New Commands Discovery Engine (`gitmap new-commands` / `gitmap nc`):**
   - Catalog the last 100 new commands from recent GitMap history (e.g., `cpar`, `pe all`, `te all`, `prompt show --copy`, `scan export`, `scan merge`, `agent task`, `space backup-branch`, etc.).
   - Support `--limit` (default: 100), `--filter` / `-f` / `-q` search, `--category` filter, and `--json` format.
   - Every command includes syntax, version, category, description, and concrete copy-pasteable example.

4. **CLI Help Coverage & LLM Skills Integration:**
   - CLI help documents: `cli/helpdoc/pending-commits.md` and `cli/helpdoc/new-commands.md`.
   - CLI constants in `cli/constants/constants_cli.go`.
   - Embed in `.agents/skills/gitmap/SKILL.md` for LLM agents to detect uncommitted changes and run fixes.

5. **Minor Version Bump & Release Ceremony:**
   - Update `version.json` from `6.524.0` to `6.525.0` (minor bump).
   - Generate changelog and release notes adhering strictly to relative Git paths.
   - Atomic commit & push via `gitmap cpf "pending-commits - improve table status cache backup db and new commands"`.

---

## 2. Decomposed Subtask Plan

| Subtask Code | Title | Assigned Role | Target Files |
|---|---|---|---|
| `Task-01` | Architecture Spec & Backup DB Engine Architecture | Worker 01 | `02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md`, `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/01-table-and-backup-db-spec.md` |
| `Task-02` | Component Spec & New Commands Filter Architecture | Worker 02 | `02-spec/21-app/table-improvement-and-command-enhancement/02-component-spec.md`, `.ai-memory/plans/subtasks/table-improvement-and-command-enhancement/02-component-and-commands-spec.md` |
| `Task-03` | Backup DB Engine & Dual-DB Cache Integration | Worker 01 | `cli/cmdpending/pending_commits_backup.go`, `cli/cmdpending/pending_commits_backup_test.go`, `cli/cmdpending/pending_commits_cache.go`, `cli/cmdpending/pending_commits_cache_test.go` |
| `Task-04` | Table Rendering, Wire Cache/Backup into CLI, Help & LLM Skills Sync | Worker 02 | `cli/cmdpending/pending_commits_cmd.go`, `cli/cmdpending/pending_commits_types.go`, `cli/cmdpending/new_commands_cmd.go`, `cli/constants/constants_cli.go`, `cli/helpdoc/pending-commits.md`, `cli/helpdoc/new-commands.md`, `.agents/skills/gitmap/SKILL.md` |
