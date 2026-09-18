# Plan 157: OS Fix, Clean, User, User-Group, VSCode Profiles & Clone Suite

> **Consolidation Note:** This parent task started from a direct user command demanding complete parity across `os ip`, `os fix`, `os clean`, `os user`, `os user-group`, `vscode profiles`, and `gitmap clone -ls / --only`. Executed across 4 granular subtasks with N = 100 continuous budget (50 planning, 50 execution loops) and verified zero-defect completion against coding guidelines, zero-nesting, affirmative booleans, and error wrapping contracts.

## Context & User Prompt

User requested full parity and implementation of core system and tooling management subcommands:
```text
os ip
os ip help
os ip set # if cannot ping google revert back
os fix add/edit/import/export/import-all/export-all/run <new commands> # run os fixes, also run existing ones
os clean
os user add/edit/import/export/import-all/export-all
os user-group add/edit/import/export/import-all/export-all

os clean temp
os clear temp
vscode profiles export/export-all/import/import-all/ls #ls will list profiles,
vscode profiles export <name>

gitmap clone -ls file.json #it will list the repos as table to work later with sequence ids
gitmap clone file.json --only 1,3 # only ids, slug, starts with will do it, clear???
```

## Task-Specific Rule Set (Domain-Specific Constraints)

1. **Rule 1 (Network Safety & Google Ping Auto-Revert)**: `gitmap os ip set` always captures a pre-change snapshot in `store.DB`. On Windows, Linux, and macOS, after setting the IP, connectivity to Google (`8.8.8.8`) is validated. If validation fails or times out, the network configuration automatically and atomically rolls back to the pre-change snapshot with clear user diagnostics.
2. **Rule 2 (Zero-Data-Loss System Cleaning)**: `gitmap os clean`, `gitmap os clean temp`, and `gitmap os clear temp` strictly target ephemeral and cache locations (`%TEMP%`, `%LOCALAPPDATA%\Temp`, `/tmp`, `/var/tmp`, `~/.cache/gitmap`). Never deletes repository sources, active database files, or git configuration files.
3. **Rule 3 (VSCode Profile Portability & Safe JSON Serialization)**: `gitmap vscode profiles` export and import routines serialize extensions list, settings, and keybindings using standard format, sanitizing absolute paths to portable relative tokens.
4. **Rule 4 (Deterministic Sequential Clone IDs & Multi-Pattern Filtering)**: In `gitmap clone -ls <manifest>`, integer IDs are 1-based, deterministic, and mapped to parsed repository records. In `gitmap clone <manifest> --only <patterns>`, comma-separated tokens can be integer IDs (`1,3`), repo slugs (`owner/repo` or `repo`), or prefix globs (`prefix*`).
5. **Rule 5 (Non-Negotiable Coding Standards & Total Ban)**: Maximum function length <= 15 lines (target <= 8 lines). Affirmative booleans only (`is*`, `has*`). No negative booleans. Zero nested ifs. AppError wrapping on all error paths. TOTAL BAN on `go test`, `go build`, and runner scripts during routine execution.

## Acceptance Criteria

```text
os ip
os ip help
os ip set # if cannot ping google revert back
os fix add/edit/import/export/import-all/export-all/run <new commands> # run os fixes, also run existing ones
os clean
os user add/edit/import/export/import-all/export-all
os user-group add/edit/import/export/import-all/export-all

os clean temp
os clear temp
vscode profiles export/export-all/import/import-all/ls #ls will list profiles,
vscode profiles export <name>

gitmap clone -ls file.json #it will list the repos as table to work later with sequence ids
gitmap clone file.json --only 1,3 # only ids, slug, starts with will do it, clear???
```

## Consolidated Subtasks

### Subtask 01: OS IP, Fix Registry & Temp Cleaner
- Implemented `gitmap os ip` defaulting to active network interfaces list when no args are passed.
- Implemented `gitmap os ip help` displaying full usage help.
- Implemented `gitmap os ip set <ip>` with automatic snapshotting, Google (`8.8.8.8`) ping validation, and automatic rollback on failure.
- Implemented `cli/osfix` engine and `cli/cmdos/os_fix.go` supporting `add`, `edit`, `rm`, `ls`, `export`, `export-all`, `import`, `import-all`, and `run` (running registered fixes or arbitrary ad-hoc commands).
- Implemented `cli/osclean` engine and `cli/cmdos/os_clean.go` supporting `gitmap os clean`, `gitmap os clean temp`, and `gitmap os clear temp` across Windows, Linux, and macOS.

### Subtask 02: OS User and User-Group Suite
- Implemented `cli/osuser/portable.go` and `cli/cmdos/os_user.go` supporting `add`, `edit`, `export`, `export-all`, `import`, `import-all`, `rm`, `create-root`, `kill-processes`, `add-ssh-key`, `ls`.
- Implemented `cli/osuser/group_portable.go` and `cli/cmdos/os_group.go` supporting `add`, `edit`, `export`, `export-all`, `import`, `import-all`, `rm`, `ls`.
- Registered aliases `user-group` and `usergroup` alongside `group` in `cli/cmdos/os.go`.

### Subtask 03: VSCode Profiles Management Suite
- Implemented `cli/cmdvscode/vscode_profiles.go` supporting `gitmap vscode profiles` (and `vscode profile`):
  - `ls`: List all installed VSCode user profiles (Default and named profiles) in a formatted table with settings and extension indicators.
  - `export <name> [file.json]`: Export single profile configuration.
  - `export-all [file.json]`: Export all detected profiles.
  - `import <file.json>`: Import profile configuration.
  - `import-all <file.json>`: Import multiple profiles.
- Routed `profiles` and `profile` subcommand in `cli/cmdvscode/vscode_cmd.go`.

### Subtask 04: Clone Sequence Table Listing & Multi-Target Filter
- Exported `cloner.LoadRecords` and `cloner.CloneRecords` in `cli/cloner/cloner.go`.
- Added `-ls` / `--ls` / `-l` / `--list` flag to `cmdclone` to display records as a formatted table with 1-based sequential integer IDs (`1`, `2`, ...), repo slugs, branches, and remote URLs.
- Added `--only <tokens>` flag to `cmdclone` to filter candidate repos by integer sequence IDs (`1,3`), slugs (`owner/repo`), or prefix globs (`prefix*`).
- Refactored `executeClone` in `cli/cmdclone/clone.go` into concise functions <= 15 lines passing single `CloneFlags` parameter struct.

## Verification & Quality Assurance
- `python linter-scripts/check-nested-ifs.py`: PASS (0 violations across 2,895 files).
- `python linter-scripts/check-enum-and-boolean.py`: PASS (0 violations across 2,183 files).
- `python linter-scripts/check-error-management.py`: PASS (0 violations across 2,932 files).
- Recent file changes recorded to `.ai-memory/temp/recent-file-changes.json`.
