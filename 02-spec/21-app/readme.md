# Application Specifications Index (`02-spec/21-app`)

This directory houses the authoritative application specifications for GitMap. In accordance with the compaction and consolidation architecture ([Spec 235](235-app-spec-and-completed-plans-consolidation-and-reduction/01-architecture-spec.md)), core systems are organized into **Eight Canonical Domain Clusters**, representing the ratified architecture across all subsystems.

---

## 1. Canonical Architecture Domain Clusters

Each domain cluster represents the authoritative, consolidated specification for a major GitMap subsystem, embodying the Final Version Invariant (retaining ratified contracts and pruning intermediate drafts):

| Cluster | Subsystem Scope | Specifications | Status |
| :--- | :--- | :--- | :---: |
| **`01-cli-architecture/`** | Cobra CLI hierarchy, middle-ellipsized tables, shell tab completion, Win32 CP handoff, typed JSON Envelope V2 | [Architecture](01-cli-architecture/01-architecture-spec.md) • [Component](01-cli-architecture/02-component-spec.md) | `ratified` |
| **`02-scanner-and-projects/`** | High-speed filesystem discovery, project classification heuristics, `.gitmapignore` engine, AUM indexing, dedup | [Architecture](02-scanner-and-projects/01-architecture-spec.md) • [Component](02-scanner-and-projects/02-component-spec.md) | `ratified` |
| **`03-git-operations-and-pull/`** | 8-worker concurrency pull pool, semantic flat commit suite (`gitmap c`), auto-staging, push self-healing, Oh-My-Zsh ignore rules | [Architecture](03-git-operations-and-pull/01-architecture-spec.md) • [Component](03-git-operations-and-pull/02-component-spec.md) | `ratified` |
| **`04-fleet-nodes-and-ssh/`** | Unified `gitmap nodes` CLI, dual-stack ICMP/TCP ping probing, RSA-OAEP salt credential vault, except-self remote clone | [Architecture](04-fleet-nodes-and-ssh/01-architecture-spec.md) • [Component](04-fleet-nodes-and-ssh/02-component-spec.md) | `ratified` |
| **`05-antigravity-and-ide/`** | Google Antigravity SDK workflows, multi-conversation prompt dispatch, theme parity, plugins/skills sync, Cursor IDE integration | [Architecture](05-antigravity-and-ide/01-architecture-spec.md) • [Component](05-antigravity-and-ide/02-component-spec.md) | `ratified` |
| **`06-database-and-split-db/`** | Three-tier SQLite Split-DB engine (`gitmap.db`, `installation.db`, `repodb/pipeline.db`), WAL mode, `SetMaxOpenConns(1)`, PascalCase schema | [Architecture](06-database-and-split-db/01-architecture-spec.md) • [Component](06-database-and-split-db/02-component-spec.md) | `ratified` |
| **`07-pipeline-and-diagnostics/`** | CI/CD pipeline telemetry (`pe`/`pea`), traceback extractors, heatmap failure summaries, dynamic ETA forecaster, 4-part RCA engine | [Architecture](07-pipeline-and-diagnostics/01-architecture-spec.md) • [Component](07-pipeline-and-diagnostics/02-component-spec.md) | `ratified` |
| **`08-distribution-and-release/`** | NSIS Windows installer payload auto-detection, Linux archives, SemVer release ceremony, cross-platform runners (`run.ps1`, `run.sh`) | [Architecture](08-distribution-and-release/01-architecture-spec.md) • [Component](08-distribution-and-release/02-component-spec.md) | `ratified` |

---

## 2. Active Application Specifications

The following application specifications record active feature developments and recent ratified milestones:

- [235-app-spec-and-completed-plans-consolidation-and-reduction](235-app-spec-and-completed-plans-consolidation-and-reduction/01-architecture-spec.md) — Application Specifications & Completed Plans Consolidation & Memory Reduction (Specs: [Master Ledger](235-app-spec-and-completed-plans-consolidation-and-reduction/00-master-audit-ledger.md), [Architecture](235-app-spec-and-completed-plans-consolidation-and-reduction/01-architecture-spec.md), [Component](235-app-spec-and-completed-plans-consolidation-and-reduction/02-component-and-cli-spec.md)) (Status: `completed`)
- [234-codebase-review-remediation-and-consolidation](234-codebase-review-remediation-and-consolidation/01-architecture-spec.md) — Codebase Review Remediation, Root Clutter Purge, README Index Reduction, Version Identity & Clone Spec (Status: `completed`)
- [233-ubuntu-ide-and-github-desktop-scan-sync](233-ubuntu-ide-and-github-desktop-scan-sync/01-architecture-spec.md) — Ubuntu Multi-IDE & GitHub Desktop Scan Integration, CLI Management Suite & Storage Sync (Status: `completed`)
- [232-pending-commits-sends-and-nodes-commit-suite](232-pending-commits-sends-and-nodes-commit-suite/01-architecture-spec.md) — Pending Commits Scanner, Sends Command Dispatcher & Remote Nodes Cluster Delegation Suite (Status: `completed`)
- [231-antigravity-ide-projects-and-repo-secrets-restore](231-antigravity-ide-projects-and-repo-secrets-restore/01-architecture-spec.md) — Antigravity IDE Projects Ingestion (78 repos), Pinned Projects Suite & Repo-Secrets Migration Documentation (Status: `completed`)
- [230-token-purge-installer-workdir-pull-agm-and-ui-modernization](230-token-purge-installer-workdir-pull-agm-and-ui-modernization/readme.md) — Security Token Purge, Installer Navigation (`$work`, `$def`), Pull Remediation, AGM Linux Update, Multi-Instance API & UI Modernization (Status: `completed`)
- [230-antigravity-backup-e2e-and-release](230-antigravity-backup-e2e-and-release/01-architecture-spec.md) — Antigravity IDE Backup E2E Validation, Logging Verification & Minor Release v6.494.0 (Status: `completed`)
- [229-antigravity-ide-projects-and-settings-backup](229-antigravity-ide-projects-and-settings-backup/01-architecture-spec.md) — Antigravity IDE Projects Ingestion (78 repos), Pinned Projects Suite & Repo-Secrets Migration Documentation (Status: `completed`)
- [229-nodes-deploy-repos-and-multi-ide-fleet-sync](229-nodes-deploy-repos-and-multi-ide-fleet-sync/01-architecture-spec.md) — Fleet Repository Deployment, Remote Scanner Delegation & Multi-IDE Sync Engine (Status: `completed`)
- [228-nodes-deploy-agm-accounts-and-fleet-sync](228-nodes-deploy-agm-accounts-and-fleet-sync/01-architecture-spec.md) — Native Multi-Node AGM Accounts Deployment, Fleet Synchronization & AGM Integration (Status: `completed`)
- [221-ci-cd-fix-nested-if-and-test-summary-remediation](221-ci-cd-fix-nested-if-and-test-summary-remediation/01-architecture-spec.md) — CI/CD Fix: Nested If Linter & Test Failure Summary Remediation (Status: `completed`)
- [220-pipeline-pe-unit-test-traceback-and-heatmap](220-pipeline-pe-unit-test-traceback-and-heatmap/01-architecture-spec.md) — Pipeline PE Unit Test Traceback Extraction, Heatmap Modernization & CI Test Remediation (Status: `completed`)
- [219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging](219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/01-architecture-spec.md) — Ubuntu Pull-All Remediation Engine, Oh-My-Zsh & Default System Directory Scanner Exclusions, Cross-OS Path Healing, Failure Subtree Formatting (Status: `active`)
- [218-nodes-agy-ui-remote-settings-and-cursor-automation](218-nodes-agy-ui-remote-settings-and-cursor-automation/01-architecture-spec.md) — Fleet Nodes Settings Sync, Remote Project-Scoped Prompting, Nodes AGY UI Dashboard, Privacy IP Scrub, and Cursor Subsystem Automation (Status: `completed`)
- [217-antigravity-fleet-parity-theme-preset-plugins-and-delegation](217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/01-architecture-spec.md) — Antigravity Fleet Parity, UI Theme Seeds, Turbo Permission Preset, 4 Plugins & 43 Skills Sync, SUID Sandbox Hardening (Status: `completed`)
- [216-gitmap-prompting-freeze-and-suggestion-engine-fix](216-gitmap-prompting-freeze-and-suggestion-engine-fix/01-architecture-spec.md) — GitMap Terminal Prompt Input Freeze Remediation, Typo Suggestion Engine Overhaul, and Cobra Shell Tab Completion Restoration (Status: `completed`)
- [215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager](215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager/01-architecture-spec.md) — Ubuntu Fleet Filesystem Hygiene, Antigravity Projects Registry Engine, App Uninstaller Subsystem (Status: `completed`)
- [214-ubuntu-fleet-automation-and-workstation-governance](214-ubuntu-fleet-automation-and-workstation-governance/01-architecture-spec.md) — Ubuntu Fleet Automation, 74 Workspaces Clone & Parity Audit, GNOME Ergonomics, VMware Automount, Antigravity Update (Status: `completed`)
- [213-antigravity-ubuntu-update-and-macro-automation](213-antigravity-ubuntu-update-and-macro-automation/01-architecture-spec.md) — Antigravity Ubuntu Update, In-App Updater RCA, SUID Sandbox Hardening, GitMap Macro Automation (Status: `completed`)
- [211-ubuntu-fleet-git-clone-and-os-customization](211-ubuntu-fleet-git-clone-and-os-customization/01-architecture-spec.md) — Ubuntu Fleet Git Clone, Desktop Ergonomics, VMware Automount & Antigravity Cross-OS Sync (Status: `completed`)
- [210-gitmap-push-fix-command-and-auth-recovery](210-gitmap-push-fix-command-and-auth-recovery/01-architecture-spec.md) — GitMap Push-Fix Command Suite, SSH Auth Self-Healing, Non-Fast-Forward Auto-Rebase & Stack Trace Suppression (Status: `completed`)
- [209-ai-agent-task-orchestrator-and-split-db](209-ai-agent-task-orchestrator-and-split-db/01-architecture-spec.md) — AI Agent Task Orchestrator, 3-Tier Multi-Agent SQLite Split-DB Hierarchy, Native `gitmap agent` CLI, Telemetry Logging, Crash Forensics & Interactive Web UI (Status: `active`)
- [208-detectedproject-foreign-key-and-ssh-lifecycle-hardening](208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/01-architecture-spec.md) — DetectedProject Foreign Key RCA, Database Reset Engine, SSH Terminal Layout Reorder, Dedicated Key Viewer, and Cross-OS Port Management (Status: `active`)
- [206-windows-to-ubuntu-fleet-migration-and-secrets-vault](206-windows-to-ubuntu-fleet-migration-and-secrets-vault/01-architecture-spec.md) — Windows to Ubuntu Fleet Migration, Zero-to-End Autonomous Provisioning, and Repo Secrets Vault Isolation (Status: `completed`)
- [205-gitmap-u1-ubuntu-agm-fleet-integration](205-gitmap-u1-ubuntu-agm-fleet-integration/01-architecture-spec.md) — GitMap U1 Ubuntu Fleet Integration, AGM Migration & Cross-OS Automation (Status: `active`)
- [204-ssh-password-interception-and-rsa-credential-vault](204-ssh-password-interception-and-rsa-credential-vault/01-architecture-spec.md) — SSH Password Interception, Masked Terminal Prompt, User RSA Consent, and RSA-OAEP Salt Credential Vault (Status: `completed`)
- [202-vmware-cli-commands-and-fleet-management.md](202-vmware-cli-commands-and-fleet-management.md) — VMware Workstation CLI Command Suite, Fleet Management, Dual-Engine Hypervisor Driver, and Coordinated Scheduler Shutdown (Status: `active`)
