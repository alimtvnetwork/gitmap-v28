# Task 227 Master Audit Ledger: GitMap Cursor Fleet Delegation, Settings UI Modernization & Cross-Node Path Suggestions

**Slug:** `227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions`  
**Status:** In Progress (Phase 1 Spec & Planning)  
**Parent Orchestrator:** Antigravity v6 Protocol (`N=300, A=2, H=2`)  
**Created:** 2026-10-05  

---

## 1. Executive Summary & Intent

This task executes four core user requirements:
1. **Confirmation of Cursor IDE Migration on Ubuntu Node (`u1`)**:
   - Verification that all paths are mapped to `/home/a/git-work/*` (Linux POSIX convention).
   - SQLite databases (`state.vscdb`, `conversation-search.db`), 49 workspaces, 61 project manager projects, and 47 project directories restored with 100% integrity.
2. **Native GitMap Cursor Delegation Command (`gitmap cursor migrate` / `gitmap cur delegate`)**:
   - Integrate end-to-end delegation CLI into `cli/cmdcursor/`.
   - Support flags: `--node <alias>` (default: `u1`), `--exclude <p1,p2>`, `--skip-settings`, `--settings-only`, `--dry-run`, `--no-backup`, `--force`, `--output <path>`.
   - Wire into `cli/cmdcursor/cursor_cmd.go` with aliases `migrate`, `m`, `delegate`, `del`.
   - Record dual execution telemetry into SQLite Split-DB (`commands.db` and `ai-instruction/sql.db`).
3. **Settings UI Redesign & Embedded Interactive Terminal**:
   - Overhaul `#tab-settings` in `cli/cmdui/ui_assets.go` according to `02-spec/24-app-ui-design-system/`.
   - Resolve critical two-way binding data loss bug where Card 2 speed settings (`lap.default_hours`, `account_switch.threshold`, `telegram.*`, `email.*`, etc.) were ignored by `loadSettings` and `saveSettings`.
   - Organize Settings into clean sub-sections: General, Cursor Fleet, SSH Nodes, Telemetry, Cache & DB, Terminal.
   - Implement `/api/terminal/exec` in `cli/cmdui/ui_server.go` and embed interactive terminal widget into `#tab-settings` with real-time execution, command history, and quick chips.
4. **Cross-Node Path Suggestion Engine for Install, Update & CLI**:
   - Resolve missing suggestion options during install/update and CLI interactions.
   - Implement candidate directory discovery in `install-quick.ps1` and `install-quick.sh`.
   - Implement cross-node repo completion and suggestion aggregation querying local repos and cached remote fleet repos (e.g. `u1:<repo>`).
5. **Quality Verification, Minor Bump & Release**:
   - Enforce zero linter errors (`check-nested-ifs.py`, `check-enum-and-boolean.py`).
   - Run Go unit tests across `cmdcursor`, `cmdui`, and completion packages.
   - Perform minor SemVer bump ceremony and commit/push atomically via GitMap.

---

## 2. Phase 1 Discovery Subagent Findings

### Subagent 1: Cursor Delegation CLI Architecture
- **Router**: `cli/cmdcursor/cursor_cmd.go:routeCursorSubcommand` dispatches subcommands. Adding `case "migrate", "m", "delegate", "del": return RunCursorMigrate(args[1:])` immediately routes both `gitmap cursor migrate` and `gitmap cur delegate`.
- **Modularity**: Split Go implementation into `cursor_migrate.go` (options & flag parsing) and `cursor_migrate_exec.go` (script resolver, `cmdpy.RunPy` execution, Split-DB telemetry) to maintain strict < 100 line SRP boundaries.
- **Python Script Upgrades**: Enhance `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py` with `--exclude`, `--skip-settings`, and `--settings-only` support.

### Subagent 2: Settings UI & Path Suggestions Architecture
- **Settings Bug Discovered**: `saveSettings()` only saved Card 1 (8 fields) and dropped all 11 Card 2 speed settings.
- **API Endpoint**: Implement `/api/terminal/exec` in `cli/cmdui/ui_server.go` with 30s context timeout, separate stdout/stderr capture, exit code extraction, and remote SSH routing.
- **Path Suggestions**: Update `install-quick.ps1` and `install-quick.sh` with candidate directory probing; extend CLI completion to incorporate cached remote fleet paths.

---

## 3. Disjoint Subtask Allocation

| Subtask ID | Title | Primary Files | Status |
|:---|:---|:---|:---|
| `01` | Cursor Migrate CLI Command & Script Flags | `cli/cmdcursor/*`, `repo-secrets/05-scripts/*` | Ready |
| `02` | Settings UI Modernization & Embedded Terminal | `cli/cmdui/ui_server.go`, `ui_assets.go`, `ui_types.go` | Ready |
| `03` | Cross-Node Path Suggestions Engine | `cli/cmdprompt/*`, `cli/completion/*`, `install-quick.*` | Ready |
| `04` | Quality Verification, Minor Bump & Release | `version.json`, `changelog.md`, linters, tests | Queued |
