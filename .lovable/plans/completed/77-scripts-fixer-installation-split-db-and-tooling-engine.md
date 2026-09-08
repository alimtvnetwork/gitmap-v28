# 77-scripts-fixer-installation-split-db-and-tooling-engine.md

## 1. Executive Summary

This plan fulfills the integration of cross-platform developer tools, package management workflows, and installation scripts from `D:\work\scripts-fixer` directly into GitMap:

1. **Dedicated Installation Split Database (`installation.db`)**:
   - Strictly adheres to `spec/04-database-conventions/07-split-db-pattern.md` and repository database conventions.
   - Isolates package installation records and voluminous execution logs from the root SQLite database into `<BinaryDataDir>/installation.db` with `conn.SetMaxOpenConns(1)`.
   - Defines PascalCase schema: `InstalledTool` and `InstallationLog`.

2. **Master SQLite Database Registry (`SplitDatabaseRegistry`)**:
   - The main root database (`gitmap.db`) knows about and tracks `installation.db` (and all split DBs) via the new `SplitDatabaseRegistry` table.
   - Automatically synchronizes split database paths, table counts, record counts, file sizes, and health status during migration and runtime operations.
   - Exposes typed query APIs: `RegisterSplitDB`, `GetSplitDB`, `ListSplitDBs`, and `SyncKnownSplitDatabases`.

3. **Universal Execution & Failure Telemetry**:
   - Upgrades `gitmap/cmd/installtools.go` with a dual-stream audit command runner capturing duration (`DurationMs`), exit code (`ExitCode`), stdout, stderr, and the exact command line.
   - Captures and records all installation actions, intermediate phases, user cancellations, and execution failures into `InstallationLog` in `installation.db`.
   - Adds CLI log inspection: `gitmap install logs` / `gitmap in --logs`.

4. **Robust Ubuntu Google Chrome Installer Upgrade**:
   - Upgrades `gitmap/cmd/install_chrome_deb.go` by adopting the proven 6-step pipeline from `scripts-fixer` (`scripts/os/ubuntu/install-chrome.sh`):
     1. `sudo apt-get update`
     2. `sudo apt-get install -y wget curl`
     3. Download official `.deb` directly to `/tmp/google-chrome-stable_current_amd64.deb`
     4. `sudo apt-get install -y /tmp/google-chrome-stable_current_amd64.deb` (activates APT's local-deb dependency resolver)
     5. Immediate deletion of scratch `.deb`
     6. Post-install verification via `google-chrome --version`
   - Fully instruments all phases with duration, exit codes, and stdout/stderr logged to `installation.db`.

5. **Expanded Cross-Platform Tooling from `scripts-fixer`**:
   - Integrates tool identifiers, aliases, descriptions, and package manager mappings from `scripts-fixer` for Rust, .NET SDK, Java OpenJDK, Ollama, Docker, Kubernetes, Zsh, Flameshot, and others.

---

## 2. Mandatory Rules & Invariants

1. **Rule 1 (Strict Relative Git Paths)**: All paths in markdown files, artifacts, and documentation must be strictly relative to the repository root. Zero `file:///` URIs.
2. **Rule 2 (Strict File & Function Sizing)**: Every function must be <= 15 lines (preferred <= 8 lines), with a mandatory blank line before every return statement. All new and modified Go files must remain <= 200 lines.
3. **Rule 3 (Database Conventions Compliance)**: All tables singular PascalCase (`SplitDatabaseRegistry`, `InstalledTool`, `InstallationLog`), integer auto-increment PKs (`{TableName}Id`), affirmative booleans (`IsActive`, `IsSuccess`), and Rule 10/11 context columns (`Description`, `Notes`, `Comments`).
4. **Rule 4 (SQLite Concurrency Standard)**: All split databases must call `conn.SetMaxOpenConns(1)` and anchor paths via `store.BinaryDataDir()` / `filepath.EvalSymlinks(os.Executable())`.
5. **Rule 5 (CI/CD Execution Timing)**: CI/CD runner is only executed at the end of the final creation step, as requested by the user.

---

## 3. Subtasks Breakdown

1. [01-main-db-split-database-registry.md](../subtasks/77-scripts-fixer-installation-split-db-and-tooling-engine/01-main-db-split-database-registry.md): Implement `SplitDatabaseRegistry` schema, models, and sync methods in `gitmap/store/split_database_registry.go` and `gitmap/store/store.go`.
2. [02-installation-split-db-telemetry-and-runner.md](../subtasks/77-scripts-fixer-installation-split-db-and-tooling-engine/02-installation-split-db-telemetry-and-runner.md): Implement comprehensive execution/failure logging and dual-stream runner in `gitmap/store/installation_split_log.go` and `gitmap/cmd/installtools.go`.
3. [03-ubuntu-chrome-installer-upgrade.md](../subtasks/77-scripts-fixer-installation-split-db-and-tooling-engine/03-ubuntu-chrome-installer-upgrade.md): Upgrade Ubuntu Chrome installer in `gitmap/cmd/install_chrome_deb.go` with `scripts-fixer` 6-step pipeline and multi-phase audit logging.
4. [04-expand-scripts-fixer-tools-and-commands.md](../subtasks/77-scripts-fixer-installation-split-db-and-tooling-engine/04-expand-scripts-fixer-tools-and-commands.md): Expand tool definitions and package managers in `gitmap/constants/constants_install.go` and `gitmap/cmd/install.go`.
5. [05-install-logs-cli-and-unit-tests.md](../subtasks/77-scripts-fixer-installation-split-db-and-tooling-engine/05-install-logs-cli-and-unit-tests.md): Add CLI `gitmap install logs` command and comprehensive unit tests for Split Database Registry and Installation Logging.
6. [06-quality-gate-verification-and-cicd.md](../subtasks/77-scripts-fixer-installation-split-db-and-tooling-engine/06-quality-gate-verification-and-cicd.md): Execute all quality gates and local CI runner `python 03-ai-scripts/06-cicd-local-runner.py`.

---

## 4. Acceptance Criteria

- [x] `SplitDatabaseRegistry` table created in main root SQLite DB with PascalCase columns, integer PK, and affirmative booleans.
- [x] `SyncKnownSplitDatabases()` discovers and records `installation.db` (and existing split DBs) into `SplitDatabaseRegistry`.
- [x] `InstallationLog` in `installation.db` records duration, exit code, stdout, stderr, command line, and success/failure status for all installs.
- [x] Ubuntu Google Chrome installation follows the proven 6-step pipeline from `scripts-fixer` (`apt-get update`, install `wget curl`, download deb to `/tmp`, `apt-get install /tmp/*.deb`, cleanup, verify).
- [x] `gitmap install logs` prints formatted execution logs from `installation.db`.
- [x] All functions <= 15 lines with blank lines before return statements.
- [x] CI/CD local runner passes all 33 quality gates with `exit 0`.
