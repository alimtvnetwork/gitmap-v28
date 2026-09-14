---
name: profile-vmware-installer-sqlite
description: Autonomously implement and verify VMware installation, profile installation idempotency, shell & SQLite database tracking, tree display for already-installed profiles, startup mounting persistence, error logging, and AppError envelopes across GitMap and common-linux-installer.
---

# Profile & VMware Installer with SQLite Database Tracking

Autonomously implement and verify idempotency and persistent tracking for profiles and VMware tooling across Linux (Shell) and Windows.

## Core Mandates

1. **Idempotency & Tree Display:** When installing a profile or package that is already installed, detect this state, report "already installed", and display an installation summary tree with status, paths, and timestamps.
2. **SQLite Database & Log Tracking:** Record all installation attempts, package metadata, execution steps, durations, statuses, and failure error logs into an untracked system SQLite database (e.g. `~/.gitmap/installer.db` or `/var/lib/gitmap/installer.db`) and log files. The database must never be committed to Git (enforce in `.gitignore`).
3. **Common Database Helper for Shell:** Provide reusable shell library functions (e.g. `shared/db-helper.sh`) to migrate schema, insert records, query status, update failure logs, and check installed state.
4. **VMware Cross-Platform CLI & Shell Parity:**
   - Linux Shell (`vmware/`): automate `open-vm-tools` / `open-vm-tools-desktop` installation, `/mnt/hgfs` mounting, desktop symlink creation (run once, idempotent), and cron/systemd startup job creation.
   - Windows & Unix CLI: provide `gitmap vmware` commands for install, shared mount, status, and desktop link.
5. **Universal Error Handling & Stack Traces:** All Go CLI operations must return structured `*AppError` envelopes with stack trace suppression where appropriate and full error logging to the SQLite database.
