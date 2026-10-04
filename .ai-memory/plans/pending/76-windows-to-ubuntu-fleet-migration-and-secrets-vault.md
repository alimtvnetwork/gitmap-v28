# Plan 76: Windows to Ubuntu Fleet Migration, Zero-to-End Provisioning & Secrets Vault

## Status: Active
- **Plan ID:** 76
- **Spec Reference:** [02-spec/21-app/206-windows-to-ubuntu-fleet-migration-and-secrets-vault/01-architecture-spec.md](../../02-spec/21-app/206-windows-to-ubuntu-fleet-migration-and-secrets-vault/01-architecture-spec.md)
- **Scope:** Cross-Platform Fleet Provisioning, Secrets Isolation, Multi-Agent Verification
- **Created At:** 2026-10-04
- **Parent Goal:** Build a repeatable Windows PowerShell migration runner, create a target node zero-to-end Linux installer injecting Supabase/Telegram/accounts/licenses, enforce complete secrets isolation in external `repo-secrets`, and prove 100% security hygiene.

---

## 1. Executive Summary

This plan orchestrates the creation and verification of automated provisioning scripts connecting the Windows workstation to target machine `u1`. It guarantees that all private tokens, accounts, and licenses reside strictly in `d:\work\repo-secrets\`, confirms zero secrets in the public repository, and verifies complete headless account switching on `u1`.

---

## 2. Subtasks Breakdown

### Subtask 76.1: Repository Secrets & Git History Audit Gate
- **Spec Reference:** Spec 206 §2
- **File:** `linter-scripts/check-forbidden-strings.py`, git commits
- **Actions:**
  - Execute exhaustive regex scans across all tracked files in GitMap.
  - Verify recent commits (`fa8460fd`, `2de54811`, `84732dc7`, `2c617e26`, `e6a607ff`).
  - Confirm zero exposure of passwords, tokens, API keys, or private accounts.

### Subtask 76.2: External Repo Secrets Vault & Dual Licenses Storage Setup
- **Spec Reference:** Spec 206 §5
- **File:** `d:\work\repo-secrets\08-licenses\`
- **Actions:**
  - Initialize `d:\work\repo-secrets\08-licenses\` directory.
  - Create `01-primary-enterprise.json` and `02-cluster-fleet-node.json` using typed JSON envelopes.
  - Verify integration with `supabase_config.json` and `telegram_config.json`.

### Subtask 76.3: Windows-to-Target One-Click PowerShell Migration Script
- **Spec Reference:** Spec 206 §3
- **File:** `d:\work\repo-secrets\migrate-to-u1.ps1`
- **Actions:**
  - Author robust PowerShell script that packages vault, credentials, and licenses from Windows.
  - Automatically SCPs archive and Linux installer to node `u1`.
  - Dispatches remote execution and streams verification logs.

### Subtask 76.4: Target Machine Zero-to-End Autonomous Installation Script
- **Spec Reference:** Spec 206 §4
- **File:** `d:\work\repo-secrets\install-zero-to-end.sh`
- **Actions:**
  - Create idempotent Bash installer for Ubuntu target node.
  - Installs/verifies `gitmap` and `agm` CLI binaries.
  - Resolves launcher recursion bug (`antigravity-launcher.sh`).
  - Injects Supabase endpoints, Telegram bot config, 35 Google AI accounts, and licenses.

### Subtask 76.5: End-to-End Remote Verification & Proof Capture
- **Spec Reference:** Spec 206 §6
- **Actions:**
  - Execute test run against target node `u1` (`192.168.1.22`).
  - Verify `agm switch` activates accounts and writes tokens to `~/.gemini/oauth_creds.json`.
  - Document complete evidence trail in audit ledger.
