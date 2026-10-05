# Plan 227: GitMap Cursor Fleet Delegation, Settings UI Modernization & Cross-Node Path Suggestions

**Status:** IN PROGRESS  
**Created:** 2026-10-05  
**Spec Reference:** `02-spec/21-app/227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions/`  

---

## 1. Verbatim Requirements & User Context
- Confirm Cursor migration on Ubuntu node (`u1`): ensure all paths are respective to the Ubuntu path (`/home/a/git-work/*`).
- Integrate and implement native `gitmap cursor migrate` / `gitmap cur delegate` command in GitMap with options:
  - `--exclude <project1,project2>`
  - `--skip-settings`
  - `--settings-only`
  - `--dry-run`, `--no-backup`, `--force`
- Overhaul GitMap Settings UI (`gitmap settings` / `gitmap ui settings`) according to design system specifications:
  - Fix data-loss bug in settings persistence.
  - Expand configuration sections (General, Cursor Fleet, SSH Nodes, Telemetry, Cache).
  - Embed an interactive terminal inside the settings UI with `/api/terminal/exec` for command execution and live testing.
- Implement reliable cross-node path suggestions during install, update, and CLI prompt interactions:
  - Query suggestions from local and remote nodes (e.g. `u1`) via GitMap.
- Minor version bump and release.

---

## 2. Wave Breakdown & Disjoint Subtasks

### Wave 1: Implementation
- **Subtask 01 (`subtasks/227-.../01-cursor-migrate-cli-command.md`)**:
  - Update `repo-secrets/05-scripts/migrate-cursor-memories-conversations.py` to support `--exclude`, `--skip-settings`, and `--settings-only`.
  - Implement `cli/cmdcursor/cursor_migrate.go` and `cli/cmdcursor/cursor_migrate_exec.go`.
  - Wire into `cli/cmdcursor/cursor_cmd.go` and add unit tests in `cli/cmdcursor/cursor_migrate_test.go`.
- **Subtask 02 (`subtasks/227-.../02-settings-ui-design-and-embedded-terminal.md`)**:
  - Add `TerminalExecReq` and `TerminalExecResp` in `cli/cmdui/ui_types.go`.
  - Add `/api/terminal/exec` endpoint in `cli/cmdui/ui_server.go`.
  - Redesign `#tab-settings` in `cli/cmdui/ui_assets.go`: add sub-tabs, modern cards, full two-way data binding for all 19 settings, and embedded interactive terminal widget.
  - Add unit tests in `cli/cmdui/ui_test.go`.

### Wave 2: Suggestions, Linters & Release
- **Subtask 03 (`subtasks/227-.../03-cross-node-path-suggestions-engine.md`)**:
  - Enhance `install-quick.ps1` and `install-quick.sh` with candidate directory probing and suggestion menu.
  - Enhance repo completion in `cli/cmd/root_cobra_completion.go` / `cli/completion/dynamic.go` to support cross-node path suggestions.
- **Subtask 04 (`subtasks/227-.../04-testing-linters-and-release.md`)**:
  - Verify zero linter errors (`check-nested-ifs.py`, `check-enum-and-boolean.py`).
  - Verify Go tests compile and pass.
  - Bump minor version in `version.json`, update `changelog.md` and `readme.md`, and perform atomic commit.

---

## 3. Evidence Checklist
- [ ] Remote Ubuntu migration integrity report verified.
- [ ] `gitmap cursor migrate` flag parsing and execution verified with tests.
- [ ] `/api/terminal/exec` and Settings UI verified with tests.
- [ ] Linters pass with exit code 0.
- [ ] Atomic GitMap commit and push executed.
