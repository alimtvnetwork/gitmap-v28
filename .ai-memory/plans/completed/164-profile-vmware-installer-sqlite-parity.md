# 164-profile-vmware-installer-sqlite-parity.md: Profile VMware & Shell Installer SQLite Parity Suite

**Status: completed**
**Execution Loop Count: 8 steps completed in continuous N-step self-loop (N=200 budget)**
**Task Initiation: Initiated via user request to implement profile installation idempotency (announce "already installed", show component tree), installation logging and untracked SQLite database persistence, VMware commands across Windows and Linux Shell, automatic `/mnt/hgfs` mount, one-time desktop symlink creation, startup persistence via crontab and systemd, and universal AppError envelopes with stack trace auditing.**

---

## 1. Executive Summary

This suite established complete parity for profile installations, VMware tooling, and SQLite telemetry across both the Shell environment (`common-linux-installer`) and the Go CLI (`gitmap`):

1. **Standalone Shell SQLite Engine (`common-linux-installer/shared/db-helper.sh`)**:
   - Reusable SQLite helper library for shell scripts providing schema initialization, idempotency checks (`db_check_installed`), run recording (`db_record_start`, `db_record_success`, `db_record_failure`), execution telemetry (`db_record_log`), and ASCII component tree rendering (`db_show_tree`).
   - Root untracked system database location (`~/.local/share/linux-installer/installer.db` or `/var/lib/linux-installer/installer.db`), ignored by `.gitignore`.
2. **Modernized VMware Shell Installer & Startup Mount (`common-linux-installer/vmware/`)**:
   - `install-tools.sh`: Idempotent installer for `open-vm-tools` and `open-vm-tools-desktop`. If already installed, announces already installed, renders the component tree, and exits 0 immediately.
   - `vmware-mount-shared.sh`: Mounts `/mnt/hgfs` using `vmhgfs-fuse`, creates the Desktop symlink (`~/Desktop/SharedDirectories`) safely for `$SUDO_USER` only once, and persists the `@reboot` mount job into crontab and systemd.
   - `vm-mount.sh`: Lightweight boot mount runner.
3. **GitMap SQLite Database Schema & Audit Layer (`cli/store/`)**:
   - Expanded `installation.db` with `ProfileInstallation` and `PackageInstallation` tables adhering to `db-sqlite-conventions`.
   - Comprehensive tracking of execution durations, exit codes, output logs, and full AppError stack traces.
4. **GitMap Cross-Platform VMware Engine (`cli/cmdvmware/`)**:
   - Linux: Open-vm-tools detection, `/mnt/hgfs` mount management, crontab inspection, and desktop symlinking.
   - Windows: Service Control Manager (`sc query VMTools`), Windows Registry (`HKLM\SOFTWARE\VMware, Inc.\VMware Tools`), `vmrun.exe` locator, and UNC shared folder access (`\\vmware-host\Shared Folders`).
5. **Profile Idempotency & Tree Display (`cli/cmdinstall/` & `cli/cmd/profile.go`)**:
   - Routed `gitmap profile install <name>`.
   - Idempotency check before execution: if already installed and `--force` is omitted, prints `[INFO] Profile '<name>' is already installed (installed at: <time>)` and displays the component tree without re-executing.
6. **Universal AppError Envelope**:
   - All Go errors are wrapped in strongly typed `*apperror.AppError` and persisted into SQLite for root cause analysis and reversing.

---

## 2. Mandatory Rules & Invariants Followed

1. **Rule 1 (Strict Relative Git Paths)**: All paths are referenced relatively within each repo.
2. **Rule 2 (Strict Sizing & Style)**: Every Go function $\le 15$ lines (target $\le 8$). Mandatory blank line before return statements.
3. **Rule 3 (Affirmative Booleans)**: Zero negative booleans (`is*`, `has*` only). Resolved inverted success check (`!rec.IsSuccess`) to positive control flow.
4. **Rule 4 (Zero Swallowed Errors)**: Zero swallow error policy with full stack traces persisted to database.
5. **Rule 5 (Untracked Database Root)**: Database roots (`~/.local/share/linux-installer/` and `~/.gitmap/data/`) are never committed to Git and covered by `.gitignore`.
6. **Rule 6 (Quality Gates)**: Verified 100% clean across all 38 CI/CD local runner quality gates.

---

## 3. Subtasks Consolidated

### Subtask 01: Shell SQLite Database Engine & Repository Hygiene
- **Created**: `common-linux-installer/.gitignore`
- **Created**: `common-linux-installer/shared/db-helper.sh`
- **Created**: `common-linux-installer/shared/test-db-helper.sh`
- Implemented `db_init`, `db_get_path`, `db_check_installed`, `db_record_start`, `db_record_success`, `db_record_failure`, `db_record_log`, and `db_show_tree`. Verified all 7 unit tests green.

### Subtask 02: VMware Shell Installer & Startup Mount Automation
- **Refactored**: `common-linux-installer/vmware/install-tools.sh`
- **Refactored**: `common-linux-installer/vmware/vmware-mount-shared.sh`
- **Refactored**: `common-linux-installer/vmware/vm-mount.sh`
- **Refactored**: `common-linux-installer/vmware/vmware-mount-shared-first.sh`
- **Updated**: `common-linux-installer/vmware/makefile`
- Idempotent tool installation, one-time desktop symlinking for `$SUDO_USER`, crontab `@reboot` and systemd persistence.

### Subtask 03: GitMap SQLite Schema & Audit Layer
- **Modified**: `cli/store/installation_split_db.go`
- **Created**: `cli/store/installation_split_profile.go`
- **Created**: `cli/store/installation_split_profile_test.go`
- Added `ProfileInstallation` and `PackageInstallation` DDL and CRUD methods with AppError stack trace preservation.

### Subtask 04: GitMap VMware CLI Cross-Platform Engine & Profile Routing
- **Modified**: `cli/cmdvmware/vmware.go`
- **Modified**: `cli/cmdvmware/vmware_install.go`
- **Modified**: `cli/cmdvmware/vmware_shared.go`
- **Modified**: `cli/cmdvmware/vmware_status.go`
- **Created**: `cli/cmdvmware/vmware_windows.go`
- **Modified**: `cli/cmd/profile.go`
- **Modified**: `cli/cmdinstall/installprofiles_exec.go`
- **Updated**: `cli/helptext/vmware.md`
- **Updated**: `cli/helptext/profile.md`
- Supported `gitmap vmware install`, `gitmap vmware shared enable`, `gitmap vmware shared status`, `gitmap vmware status`, and `gitmap profile install <name>`.

---

## 4. Verification & Quality Gates

- **Unit Tests**:
  - `shared/test-db-helper.sh`: 7/7 tests passed.
  - `bash -n` syntax validation passed across all shell scripts in `common-linux-installer`.
  - `go vet ./...`: Passed on Windows, Linux, and macOS (GOOS=darwin, GOOS=linux).
- **CI/CD Local Quality Gates**:
  - `python 03-ai-scripts/06-cicd-local-runner.py --no-tests`: 38/38 quality gates passed in 38.63s.
- **Inventory Tracking**:
  - All 12 modified files recorded in `.ai-memory/test-inventory.json`.
