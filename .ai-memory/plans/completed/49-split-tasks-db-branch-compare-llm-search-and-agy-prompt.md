# Plan 49: Split Tasks DB, Branch Compare (Beyond Compare), LLM AI Execution History, Search DB & AGY Prompt Enhancements

> **Execution Summary:**
> - **Origin:** Initiated from user request to separate tasks from the root DB into `.gitmap/data/tasks/sql.db` and section tasks DBs, implement task queue with undo/redo and history pagination (with negative offset), build branch compare with Beyond Compare installer/dev-profile integration, record AI execution history and search logs with normalized categories and views, and enhance GitMap EDI/AGY prompt commands.
> - **Workflow:** 3-Phase Parent Task N-Step Loop (2 Planning Subagents, 4 Execution Subagents).
> - **Total Loops/Steps:** 6 atomic subtasks executed across 2 phases.
> - **Outcome:** 100% completed with zero errors and zero linter violations.

---

## Consolidated Subtasks & Deliverables

### Subtask 01: Tasks Split DB Architecture
- **Files Modified/Created:**
  - `cli/store/split_db_path.go`
  - `cli/store/tasks_split_db.go`
  - `cli/store/split_database_registry_sync.go`
  - `cli/cmd/pendingtaskhelper.go`
  - `cli/cmd/pending.go`
  - `cli/cmd/dopending.go`
  - `cli/cmd/pendingclear.go`
  - `cli/cmd/tasks_list.go`
  - `cli/cmd/tasks_ops.go`
- **Accomplishments:**
  - Implemented `ResolveTasksRootDbPath` resolving to `.gitmap/data/tasks/sql.db`.
  - Implemented `TasksSplitDB` with tables: `TaskType`, `PendingTask`, `CompletedTask`, `TaskQueue`, and `TaskHistory`.
  - Registered Tasks Root DB in `SplitDatabaseRegistry` in `gitmap.db`.
  - Redirected all pending task queries, executions, and clears from `gitmap.db` to `.gitmap/data/tasks/sql.db`.

### Subtask 02: Section Tasks DBs, Task Queue, Undo/Redo & History
- **Files Modified/Created:**
  - `cli/store/section_tasks_db.go`
  - `cli/cmdssh/ssh_history_db.go`
  - `cli/cmdssh/ssh_history_types.go`
  - `cli/cmdssh/ssh_history_query.go`
  - `cli/cmdssh/ssh_undo_cmd.go`
  - `cli/cmdssh/ssh_keys_manage.go`
  - `cli/cmdssh/ssh.go`
  - `cli/cmdssh/sshdelete.go`
  - `cli/cmdssh/sshjoin_cmd.go`
  - `cli/cmdssh/sshjoin_enroll.go`
  - `cli/cmdssh/sshjoin_rm_cmd.go`
  - `cli/cmdtask/task_history_cmd.go`
  - `cli/cmd/tasks.go`
  - `cli/helptext/task.md`
  - `cli/helptext/tasks.md`
- **Accomplishments:**
  - Implemented `ResolveSectionTasksDbPath` resolving to `.gitmap/data/<section>/<section>-tasks.db` (e.g. `pipeline-tasks.db`, `ssh-tasks.db`).
  - Migrated SSH task history from legacy `history/task/sql.db` to `.gitmap/data/ssh/ssh-tasks.db`.
  - Enqueued `node remove`, `ssh add`, `ssh ip remove`, `ssh keys manage`, and `ssh keys remove` into TaskQueue with forward and inverse payloads for undo/redo.
  - Implemented `gitmap task history [limit] [offset]` supporting default 100 entries, custom limit, and negative offset (`-K`).

### Subtask 03: Branch Compare Tool
- **Files Modified/Created:**
  - `cli/constants/constants_compare.go`
  - `cli/cmd/dispatchcompare.go`
  - `cli/cmd/root.go`
  - `cli/cmdgit/compare_branches.go`
  - `cli/cmdgit/compare_worktree.go`
  - `cli/cmdgit/compare_launcher.go`
- **Accomplishments:**
  - Implemented `gitmap compare (cmp) <branch1> <branch2> [branch3...]`.
  - Detached worktree isolation under `.gitmap/tmp/compare-...` for non-checked-out branches with defer cleanup.
  - External compare tool discovery: Beyond Compare (`bcomp`/`bcompare` on Windows and Linux), with fallback to `code`, `meld`, `winmerge`.

### Subtask 04: Beyond Compare Installer & Dev Profile Integration
- **Files Modified/Created:**
  - `cli/constants/constants_tools.go`
  - `cli/cmdinstall/install_beyondcompare.go`
  - `cli/cmdinstall/install_handlers.go`
  - `cli/cmdinstall/install_packages.go`
  - `cli/cmdinstall/installprobe.go`
  - `cli/cmdinstall/install_uninstall_custom.go`
  - `cli/cmdinstall/installprofiles.go`
- **Accomplishments:**
  - Added tool identifiers `ToolBeyondCompare4` and `ToolBeyondCompare5`.
  - Implemented silent installers for Beyond Compare 4 & 5 on Windows (Inno Setup `/VERYSILENT`) and Linux (DEB/tar.gz) with side-by-side support.
  - Implemented uninstaller in `install_uninstall_custom.go`.
  - Added Beyond Compare 5 to default dev profiles across Windows and Linux.

### Subtask 05: AI Instruction History & Search DB
- **Files Modified/Created:**
  - `cli/store/ai_instruction_db.go`
  - `cli/store/search_split_db.go`
  - `cli/searcher/search_logger.go`
  - `cli/cmd/search.go`
  - `cli/store/split_database_registry_sync.go`
- **Accomplishments:**
  - Implemented `.gitmap/data/ai-instruction/sql.db` with `AiCommandCategory`, `AiExecutionHistory`, and `AiExecutionView`.
  - Implemented `.gitmap/data/search/sql.db` with `SearchCategory` (`user`, `ai`, `automated_audit`), `SearchRecord`, and `SearchLogView`.
  - Supported `--ai` and `-ai` flags on search commands to classify caller as AI-initiated.
  - Registered `ai-instruction` and `search` in `SplitDatabaseRegistry`.

### Subtask 06: AGY Prompt Commands & Help
- **Files Modified/Created:**
  - `cli/cmdantigravity/prompt_cmd.go`
  - `cli/cmdantigravity/prompt_queue.go`
  - `cli/cmdantigravity/prompt_project.go`
  - `cli/cmdprompt/prompt_list.go`
  - `cli/cmdagy/agy_cmd.go`
  - `cli/cmdagy/agy_prompt_subcmds.go`
  - `cli/helptext/antigravity.md`
- **Accomplishments:**
  - Implemented `gitmap agy prompt-project` (`p`, `prompt-p`) targeting project by prefix with `-name(n)`, `-txt(t)`, and `--prefix(pf)/suffix(sf)`.
  - Enhanced `gitmap agy prompt` defaulting to current repo project: auto-registers repository in Antigravity if missing, and queues default "read all first" prompt before the user prompt.
  - Implemented `gitmap agy prompt-with-name` and `gitmap agy prompt-txt`.
  - Implemented `gitmap agy prompt ls` / `list prompts` rendering prompt templates in a table.
  - Updated help documentation in `cli/helptext/antigravity.md`.

---

## Verification & Compliance
- **Linters:**
  - `python linter-scripts/check-boolean-guidelines.py`: PASS (0 violations)
  - `python linter-scripts/check-nested-ifs.py --all`: PASS (0 violations across 3,423 files)
  - `python .github/scripts/go-format-check.py --check-only`: PASS (0 unformatted across 3,205 files)
- **Coding Guidelines:** All functions <= 15 lines, affirmative booleans only, zero nested ifs, universal `AppError`/`AppFault` wrapping.
