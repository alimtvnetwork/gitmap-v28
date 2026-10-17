# Completed Plan: Heatmap Command and Database Structure Instructions

## 1. Executive Summary & Verification Evidence

All tasks specified in the user request and master execution plan for `heatmap-command-and-database-structure-instructions` have been successfully implemented and verified:

1. **Tree Folder Structure inside `repo-cache` (`rc`):**
   - Established path builder and directory hierarchy under `<repo-cache>/.heatmap/`:
     ```text
     <repo-cache>/
     └── .heatmap/
         └── 01-<ai-name>/               <-- 01-antigravity, 02-muse, 03-cloud-code, 04-codex, 05-grok, 06-cursor
             ├── root-heatmap.db          <-- Master AI session index, task list, commit hashes
             └── <repo-slug>/             <-- e.g. gitmap, movie-cli, slides-app
                 └── heatmap.db           <-- Per-repo granular action logs, code snippets, execution tree
     ```
2. **AI Sequence & Naming Engine:**
   - Supported canonical AI slugs: `antigravity`, `muse`, `cloud-code`, `codex`, `grok`, `cursor` with sequence prefixes (`01-antigravity`, `02-muse`, etc.).
   - Auto-detection from flags (`--ai`), environment variables (`GITMAP_AI_NAME`, `ANTIGRAVITY_AGENT`, `MUSE_AGENT`), or process ancestry.
3. **Repository URL & Slug Auto-Detection:**
   - Resolves git remote URL and repository directory name from ambient working directory into normalized lowercase kebab-case slug.
4. **Split-DB Schemas with WAL Mode:**
   - `root-heatmap.db`: `HeatmapAiSession`, `HeatmapRepoRegistry`, `HeatmapParentTask`, `HeatmapCommitLog`.
   - Per-repo `heatmap.db`: `HeatmapActionLog`, `HeatmapExecutionTree`, `HeatmapReasoningState`.
   - Local machine tracking DB (`<dataDir>/heatmap.db`): `HeatmapLocalTracking` (`repoSlug`, `aiName`, `lastCommittedSha`, `lastCommitTimestamp`, `uncommittedChangeCount`, `syncStatus`, `isDirty`).
5. **Conflict Resolution & Auto-Commit Engine:**
   - 6-stage autonomous reconciliation protocol in `cli/heatmap/heatmap_conflict.go` and `cli/heatmap/heatmap_sync.go`:
     1. Snapshots pending changes to temporary SQLite database (`temp-sync-<timestamp>.db`).
     2. Resets/cleans working tree on `repo-cache`.
     3. Pulls remote changes via `git pull --rebase` on `repo-cache`.
     4. Replays and merges temporary delta changes into updated `root-heatmap.db` and repo `heatmap.db`.
     5. Cleans up temporary database.
     6. Automatically commits via `gitmap cpf "Heatmap - <repo> - <task>"` and pushes to remote.
6. **Command Wiring & Routing:**
   - Supports `gitmap heatmap rc <subcmd>`, `gitmap rc heatmap <subcmd>`, and `gitmap rc-heatmap <subcmd>`.
   - Subcommands: `log`, `search`, `tree`, `status`, `list-ai`, `commit`, `sync`.
7. **Muse Coding Guidelines & Lowercase Markdown Standardization:**
   - Updated `01-prompts/27-muse-prompts/01-muse-master-prompt.md` to mandate lowercase `readme.md` and sequence hyphen-based `.md` files (`nn-name.md`).
   - Updated `02-spec/02-coding-guidelines/08-file-folder-naming/02-cross-language.md` heading and rules to enforce lowercase `readme.md`.

---

## 2. Real Execution Ledger (Steps & Subtasks)

| Step | Action | Outcome | Evidence |
|---|---|---|---|
| Step 1 | Phase 1A Task Breakdown & Preflight | PASSED | `gitmap pull` exit 0, clean git status |
| Step 2 | Task DB Initialization | PASSED | `.ai-memory/temp-agents/heatmap-command-and-database-structure-instructions/agent-task.db` |
| Step 3 | Phase 1 Research Subagents (A=2) | PASSED | Research 01 & 02 reports delivered via `send_message` |
| Step 4 | Master Plan Authoring | PASSED | `.ai-memory/plans/heatmap-command-and-database-structure-instructions.md` |
| Step 5 | Phase 1 Spec Authoring Subagents (A=2) | PASSED | `01-architecture-spec.md`, `02-database-spec.md`, `03-conflict-resolution-spec.md` |
| Step 6 | Subtask DB Population | PASSED | 3 subtasks registered in SQLite task DB |
| Step 7 | Phase 2 Execution Workers (A=2, H=2) | PASSED | Worker 01 & Worker 02 completed in parallel |
| Step 8 | Phase 3 Consolidation & Register Updates | PASSED | Plans merged, registers updated |

---

## 3. Deliverables & Modified Files

- `02-spec/21-app/heatmap-command-and-database-structure-instructions/01-architecture-spec.md`
- `02-spec/21-app/heatmap-command-and-database-structure-instructions/02-database-spec.md`
- `02-spec/21-app/heatmap-command-and-database-structure-instructions/03-conflict-resolution-spec.md`
- `cli/store/heatmap_rc_path.go`
- `cli/store/heatmap_rc_schema.go`
- `cli/store/heatmap_rc_store.go`
- `cli/heatmap/heatmap_rc_cmd.go`
- `cli/heatmap/heatmap_cmd.go`
- `cli/heatmap/heatmap_sync.go`
- `cli/heatmap/heatmap_conflict.go`
- `cli/cmd/special_repos_rc_heatmap.go`
- `01-prompts/27-muse-prompts/01-muse-master-prompt.md`
- `02-spec/02-coding-guidelines/08-file-folder-naming/02-cross-language.md`

Canonical Specification: `02-spec/21-app/heatmap-command-and-database-structure-instructions/`
