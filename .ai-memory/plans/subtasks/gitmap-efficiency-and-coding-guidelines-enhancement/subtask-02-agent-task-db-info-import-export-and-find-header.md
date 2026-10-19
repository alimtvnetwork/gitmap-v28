# Subtask 02: Agent Task DB Path & Query Commands, JSON/YAML Import/Export, Muse Skills & Find First Line

- **Subtask ID**: `subtask-02-agent-task-db-info-import-export-and-find-header`
- **Parent Task**: `gitmap-efficiency-and-coding-guidelines-enhancement`
- **Owned Files**:
  - `cli/cmdagent/agent_task.go`
  - `cli/cmdagent/agent_task_export.go`
  - `cli/cmdagent/agent_cmd.go`
  - `cli/cmd/llm/llm_train_muse_ui.go`
  - `cli/cmdfind/find_files.go`
- **Assigned Worker**: Worker 02

---

## 1. Objectives
1. In `cli/cmdagent/agent_task.go` (`renderTaskInitResult`):
   - Include `databasePath`: relative path to task DB (`.ai-memory/temp-agents/<taskSlug>/agent-task.db`).
   - Include `rootDbPath`: relative path to root DB (`gitmap.db` or task root).
   - Include `rootDbInfo`: metadata map with master DB path and schema.
   - Include `queryCommands`: list of actionable terminal commands to query tasks from sqlite3 and gitmap CLI:
     - `gitmap agent task ls`
     - `gitmap agent subtask list --task-id <id>`
     - `sqlite3 <relDbPath> "SELECT * FROM tasks;"`
     - `sqlite3 <relMasterDbPath> "SELECT * FROM parent_tasks WHERE parent_task_id = '<id>';"`
2. Implement `gitmap agent task export` and `gitmap agent task import`:
   - Support `--format json|yaml` (default json) and `--out <file>` for export.
   - Support `--format json|yaml` and `--file <file>` for import.
   - Use `gopkg.in/yaml.v3` for YAML serialization/deserialization.
   - Wire commands into `TaskCmd` under `initTaskCommands()`.
3. In `cli/cmd/llm/llm_train_muse_ui.go`:
   - Add step 5 in `printTraversalChecklist` with GitMap skills file location (`.agents/skills/gitmap/SKILL.md`) and command (`gitmap cat .agents/skills/gitmap/SKILL.md`).
   - Include skills path in `printDigestMissingGuidance` as well.
4. In `cli/cmdfind/find_files.go`:
   - Start stopwatch with `time.Now()`.
   - On line 1 before printing matches, output:
     `Search completed: %d hit(s) in %v`
     (matching `gitmap aum search` format so match count and duration are immediately visible).

---

## 2. Verification
- `go test ./cmdagent -v`
- `go test ./cmdfind -v`
- Execute `gitmap find "*install*"` to verify first line summary.
