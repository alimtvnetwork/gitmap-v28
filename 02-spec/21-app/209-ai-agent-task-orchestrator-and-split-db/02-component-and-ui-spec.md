# AI Agent Task Orchestrator: Component & Web UI Specification

## 1. Interactive Web Visualizer Dashboard (`gitmap agent ui` / `gitmap ai-agents ui`)

### 1.1 Architecture & Endpoints
The visualizer dashboard runs an embedded HTTP server providing real-time visibility into the multi-agent task hierarchy, agent lifecycles, and crash autopsies:
- **Port Resolution:** Starts scanning from port 8095 upwards (up to 50 ports) to ensure collision-free local binding on `127.0.0.1`.
- **Browser Automation:** Automatically launches the default web browser via Windows `start`, macOS `open`, or Linux `xdg-open`.
- **REST Endpoints:**
  - `GET /api/agent/status`: Returns current orchestrator status, temp directory location, active parent task, and global metrics.
  - `GET /api/agent/tasks`: Returns list of parent tasks and decomposed subtasks with current statuses.
  - `GET /api/agent/logs?taskId=<id>&agent=<agent>`: Returns granular action logs (file touched, action type, line ranges, command query, duration, error messages).
  - `GET /api/agent/crashes?taskId=<id>`: Pinpoints autopsy report for any crashed or abandoned agents.
  - `POST /api/agent/clear`: Triggers task cleanup or temp-folder wipe.

### 1.2 Embedded Single-Page Dashboard (HTML/CSS/JS)
- **Theme:** High-contrast developer dark theme matching GitMap's aesthetic.
- **Header:** Live status pill (`RUNNING`, `IDLE`, `COMPLETED`, `CRASH_DETECTED`), active task badge, elapsed duration, steps budget meter (`CurrentStep / TotalBudget`).
- **Fleet & Task Tree View:** Interactive hierarchical tree displaying:
  - Repository Root (`ai_agents.db`)
    - Parent Tasks (`79-ai-agent-task-orchestrator-and-split`, etc.)
      - Subtasks (`Task-01` .. `Task-07`)
        - Agent Split DBs (`agents/worker-01.db`, `agents/worker-02.db`)
- **Action Log Explorer & Audit Timeline:** Filterable table by agent, action type (`READ`, `WRITE`, `SEARCH`, `EXEC`), target file, and execution status with drill-down modals.
- **Crash Autopsy Modal:** Visually surfaces the exact last action, target file, timestamp, and diagnosis for incomplete or crashed subtasks.

---

## 2. Configurable Temp Folder Architecture & Lifecycle Cleanup

### 2.1 Resolution Hierarchy
1. Explicit CLI Flag: `--dir <path>`
2. Environment Variable: `GITMAP_AGENT_TEMP_DIR`
3. GitMap SQLite Setting: `db.GetSetting("agent.temp_dir")` (configured via `gitmap config set agent.temp-dir <path>`)
4. Default Fallback: `.ai-memory/temp-agents/`

### 2.2 Lifecycle Commands
- `gitmap agent clear [--task-id <id>] [-y]`: Cleans completed or specified parent task directory and database entries.
- `gitmap agent reset [-y]`: Resets `ai_agents.db` and clears all agent split databases while preserving run directories.
- `gitmap agent temp-clear [-y]`: Completely purges the `.ai-memory/temp-agents/` directory, recreating clean schemas on the next run.

---

## 3. Migration Protocol: Updating Coding Guidelines & Execute Prompts

### 3.1 Legacy Python to Native GitMap Agent Mapping
| Legacy Python Invocation | Native GitMap Agent Command |
| :--- | :--- |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py init --name "<name>" --budget 300` | `gitmap agent task init --name "<name>" --budget 300` |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py add-subtasks --db <db> --tasks-file <file>` | `gitmap agent subtask add --parent <id> --file <file>` |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py claim --db <db> --agent "<role>"` | `gitmap agent subtask claim --agent "<role>"` |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py log-action --db <db> --subtask-id <id> ...` | `gitmap agent log --subtask <id> --agent "<role>" --action <type> ...` |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py complete --db <db> --subtask-id <id> ...` | `gitmap agent subtask complete <id> --agent "<role>" --evidence "<ev>"` |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py fail --db <db> --subtask-id <id> ...` | `gitmap agent subtask fail <id> --agent "<role>" --reason "<reason>"` |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py diagnose --db <db>` | `gitmap agent crashed` (or `gitmap agent diagnose`) |
| `python 03-ai-scripts/46-agent-sqlite-task-manager.py status --db <db>` | `gitmap agent task status` |

### 3.2 Targeted Skill Updates
1. `.agents/skills/execute-parent-task-with-n-steps-v6/skill.md`: Replace python script calls with native `gitmap agent` commands.
2. `01-prompts/14-execute/13-execute-parent-task-with-n-steps-v6.md`: Update step 0, worker briefs, and crash forensics sections.
3. Dedicated Antigravity Skill: Author `.agents/skills/gitmap-agent-orchestrator/skill.md`.
