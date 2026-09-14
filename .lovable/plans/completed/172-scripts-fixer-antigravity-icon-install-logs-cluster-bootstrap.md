# Plan 172: Scripts-Fixer Parity — Antigravity Desktop Icon & Duplicate Purge, Installation Telemetry SQLite DB, and Cluster SSH Node Bootstrap

> **Task Origin:** User request to align GitMap with the latest 20 commits of `D:\work\scripts-fixer` (up to `v1.47.0`), fixing Antigravity launcher duplicates and desktop icon trust, logging installation telemetry in SQLite, and integrating cluster SSH RSA bootstrap.  
> **Execution Strategy:** 2-Phase Continuous Self-Loop (Planning & Subtask Generation -> Parallel Execution -> Minor Release Orchestration).  
> **Execution Steps Completed:** 3 Subtasks across 4 parallel subagents; total loops: 2.

---

## 1. Subtask 01: Antigravity Desktop Icon & Duplicate Launcher Purge
* **Target Files:**
  - `cli/cmdinstall/installantigravity_launcher_linux.go` [NEW]
  - `cli/cmdinstall/installantigravity_deploy_linux.go` [MODIFIED]
  - `cli/cmdinstall/installantigravity_cleanup.go` [MODIFIED]
  - `cli/cmdinstall/installantigravity_launcher_test.go` [NEW]
* **Implementation Highlights:**
  - Purged 5 conflicting desktop entries: `antigravity-ide.desktop`, `Google Antigravity.desktop`, `Google-Antigravity.desktop`, `Antigravity.desktop`, `antigravity.desktop` across `~/.local/share/applications`, `~/Desktop`, and `/usr/share/applications`.
  - Deployed canonical `antigravity.desktop` with `StartupWMClass=Antigravity` and chmod 0755 to user applications and `~/Desktop`.
  - Marked desktop file as trusted via `gio set metadata::trusted true`.
  - Deployed multi-resolution icons (256x256, 512x512, pixmaps, `/usr/share/pixmaps`).
  - Refreshed GNOME icon cache and desktop database.

---

## 2. Subtask 02: Installation Telemetry SQLite Logging
* **Target Files:**
  - `cli/store/install_logs_types.go` [NEW]
  - `cli/store/installation_split_db.go` [MODIFIED]
  - `cli/store/installation_split_install_logs.go` [NEW]
  - `cli/store/installation_split_install_logs_test.go` [NEW]
  - `cli/cmdinstall/installantigravity_db.go` [MODIFIED]
* **Implementation Highlights:**
  - Added unified `install_logs` table matching `scripts-fixer`:
    ```sql
    CREATE TABLE IF NOT EXISTS install_logs (
        id TEXT PRIMARY KEY,
        target_type TEXT NOT NULL,
        target_name TEXT NOT NULL,
        action TEXT NOT NULL,
        status TEXT NOT NULL,
        exit_code INTEGER NOT NULL DEFAULT 0,
        error_message TEXT,
        log_path TEXT,
        started_at TEXT NOT NULL,
        ended_at TEXT
    );
    CREATE INDEX IF NOT EXISTS idx_install_logs_target ON install_logs(target_type, target_name);
    CREATE INDEX IF NOT EXISTS idx_install_logs_status ON install_logs(status);
    ```
  - Implemented `RecordInstallStart`, `RecordInstallSuccess`, `RecordInstallFailure`, `RecordInstallSkipped`, `GetInstallLog`, `ListInstallLogs`.
  - Hooked Antigravity and tool installations into `install_logs`.

---

## 3. Subtask 03: Cluster SSH Node Bootstrap & Sudoers NOPASSWD
* **Target Files:**
  - `cli/cmdssh/cluster_bootstrap_cmd.go` [NEW]
  - `cli/cmdssh/cluster_bootstrap_cmd_test.go` [NEW]
  - `cli/cmd/cluster.go` [MODIFIED]
  - `cli/cmdssh/sshjoin_cmd.go` [MODIFIED]
  - `cli/helptext/cluster.md` [MODIFIED]
* **Implementation Highlights:**
  - Implemented `gitmap cluster bootstrap <target> [password]` and `gitmap sj bootstrap <target> [password]`.
  - Discovers or creates local RSA cluster keypair (`~/.ssh/id_rsa`, 4096-bit).
  - Injects public key into remote `~/.ssh/authorized_keys` using `SSH_ASKPASS` (`attachAskPass`).
  - Injects `/etc/sudoers.d/<user>` with `<user> ALL=(ALL) NOPASSWD:ALL` and chmod 0440 (default `--sudo=true`).
  - Verifies key-based authentication with `BatchMode=yes` (`ssh -i <key> -o BatchMode=yes <user>@<ip> "echo ssh_ok"`).
  - Wipes plaintext password from memory immediately upon key deployment.
  - Enrolls host into SQLite `ssh_hosts` and `ssh_history`.
  - Displays formatted execution summary table (`NODE`, `IP`, `SUDO`, `KEY_AUTH`, `STATUS`, `DURATION`).

---

## 4. Release Status
- Staged and verified via code formatters, LF normalizers, and encoding checkers.
- Published minor release **v6.239.0** to GitHub releases via `29-release-orchestrator.py --tier minor --skip-tests`.
