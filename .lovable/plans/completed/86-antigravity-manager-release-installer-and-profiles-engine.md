# Plan 86: Antigravity Manager Release Installer, Dual CLI Parity, and Installation Profiles Engine

**Title:** Antigravity & Antigravity Manager Release Installer and Installation Profiles Engine  
**Status:** Completed  
**Parent Task Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (Prompt v2.1.0)  
**Budget (N):** 50 steps  
**Target Directory:** `.lovable/plans/completed/`  

---

## 1. Context & Objectives

1. **Antigravity Manager Latest Release Installer**:
   - Upstream repo: `https://github.com/lbjlaq/Antigravity-Manager`.
   - Automatically determine latest release via GitHub API with fallback to tag redirect parsing.
   - Match platform assets (`.exe`, `.msi`, `.deb`, `.AppImage`, `.dmg`).
   - Download and run silent/native installation with audit metrics recorded to `installation.db`.
2. **Antigravity CLI Live Installer**:
   - Upgrade `runInstallAntigravity()` from placeholder message to cross-platform installer (`get.antigravity.dev` or `npm install -g @google/antigravity`).
   - Probe binary and version post-install.
3. **Dual CLI Parity**:
   - Installable via `gitmap install ag-manager` and `gitmap install antigravity`.
   - Also installable via `gitmap agy install` (and `gitmap ag install`).
4. **Installation Profiles System (following `scripts-fixer`)**:
   - Introduce profiles: `minimal`, `dev`, `ubuntu` (alias `ubuntu-dev`), `ubuntu-dev-ai`, `ai`, `backend`, `fullstack`.
   - Display dedicated "Installation Profiles" table in `gitmap install ls`.
   - Support `gitmap install profile <name>` and `gitmap install <profile_name>` directly.

---

## 2. Subtasks Ledger

- **Subtask 86.01 (DONE):** Upgrade Antigravity Manager release resolution with GitHub API User-Agent and tag redirect fallback in `gitmap/cmd/installagmanager.go` and `gitmap/cmd/installagmanager_fetch.go`.
- **Subtask 86.02 (DONE):** Upgrade `runInstallAntigravity()` in `gitmap/cmd/installantigravity.go` to live cross-platform installer with npm fallback and DB audit tracking.
- **Subtask 86.03 (DONE):** Add `gitmap agy install` command in `gitmap/cmd/agy_install.go` and wire into `gitmap/cmd/agy_cmd.go`.
- **Subtask 86.04 (DONE):** Implement `Installation Profiles Engine` in `gitmap/cmd/installprofiles.go` and execution runner in `gitmap/cmd/installprofiles_exec.go`.
- **Subtask 86.05 (DONE):** Add Profiles table block to `printInstallListGrouped()` in `gitmap/cmd/installlist.go` and profile command routing in `gitmap/cmd/install.go`.
- **Subtask 86.06 (DONE):** Add unit tests in `gitmap/cmd/installprofiles_test.go`, verify linters (`check-nested-ifs.py`, `check-enum-and-boolean.py`), and test live binary.
