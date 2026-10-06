# Master Audit Ledger: 232-ubuntu-ide-and-github-desktop-scan-sync

- **Request Slug:** `232-ubuntu-ide-and-github-desktop-scan-sync`
- **Request Verbatim First Line:** `Serious issue in the gitmap. When we do the scan, it should actually add the repos to the GitHub Desktop and VS Code.`
- **Status:** `ACTIVE`
- **Phase:** `1 (Planning & Discovery)`
- **Wave:** `0 / 3`
- **Step:** `1 / 300`
- **Last Completed Action:** `Task DB Initialized and Ledger Created`
- **Next Action:** `Phase 1 Planning Subagents Discovery Dispatch`
- **Workers In Flight:** `none`
- **Commits:** `none`
- **Pushed:** `no`
- **Branch:** `main` | **Tree at Start:** `clean`
- **Tools:** `invoke_subagent=yes send_message=yes ask_question=yes gitmap=yes sqlite_db=.ai-memory/temp-agents/ubuntu-ide-and-github-desktop-scan-sync/agent-task.db`

---

## Task Matrix

| Task-ID | Subtask | Owner | Owned Files | Status | Evidence |
| :--- | :--- | :--- | :--- | :--- | :--- |
| **Task-01** | `01-e2e-investigation-and-grounded-rca` | Worker 01 | `02-spec/22-app-issues/232-ubuntu-ide-github-desktop-scan-omission.md`, `.ai-memory/issues/` | PENDING | - |
| **Task-02** | `02-standalone-ide-sync-scripts` | Worker 01 | `03-ai-scripts/40-ubuntu-ide-desktop-sync.py`, `03-ai-scripts/` | PENDING | - |
| **Task-03** | `03-gitmap-ide-command-suite` | Worker 01 | `cli/cmdide/`, `cli/cmd/` | PENDING | - |
| **Task-04** | `04-scan-integration-and-dedup-check` | Worker 02 | `cli/scanner/`, `cli/cmdscan/` | PENDING | - |
| **Task-05** | `05-scan-sync-and-skip-flags` | Worker 02 | `cli/scanner/`, `cli/cmd/` | PENDING | - |
| **Task-06** | `06-targeted-verification-and-release-ceremony` | Worker 02 | `version.json`, `package.json`, `readme.md` | PENDING | - |

---

## Assumptions & Decision Boundaries
- Target platforms: Ubuntu Linux (with cross-platform hooks for macOS and Windows).
- GitHub Desktop storage on Linux: Flatpak / Native deb packages store state in `~/.config/GitHub Desktop/IndexedDB` / SQLite LevelDB or registry config files.
- VS Code on Linux: Project Manager (`~/.config/Code/User/globalStorage/alefragnani.project-manager/projects.json`) and VS Code Recent Workspaces (`~/.config/Code/User/globalStorage/state.vscdb`).
- Antigravity on Linux: `~/.gemini/config/projects/` and `~/.gemini/antigravity/`.
- Cursor on Linux: `~/.config/Cursor/User/globalStorage/`.
- Skip flags: `--skip-sync`, `--exclude-sync`, `--sync-ide`.
- CLI commands: `gitmap ide add <path>`, `gitmap ide sync`, `gitmap ide remove <path>`, `gitmap ide ls`.
