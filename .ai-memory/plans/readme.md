# Master Plans Registry (`.ai-memory/plans`)

This registry catalogs the active, completed, and pending execution plans for GitMap. In accordance with Plan 236 ([Spec 236](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md)), historical completed plans have been consolidated into **32 Authoritative Milestones** under `.ai-memory/plans/completed/`.

---

## 1. Active Plans

The following plans represent open and actively maintained development streams:

- [333-clone-and-fix-git-repositories.md](completed/333-clone-and-fix-git-repositories.md) — Plan 333: GitMap Clone Name Resolution, Suggestion Engine & Antigravity IDE Integration: Bare Repo Name & Slug Resolution, Catppuccin Fuzzy Suggestion Box, CFR/CRFB Routing Disambiguation, Dedicated CRF/CRFB Aliases, Automated Post-Clone Antigravity Workspace Registration, and Target Repos Cloning (`yt-manage-v1`, `go-email-sender`) (Spec: [333-clone-and-fix-git-repositories](../../02-spec/21-app/333-clone-and-fix-git-repositories/readme.md)) (Status: `completed`)

- [330-context-menu-and-gitmap-ui-enhancements.md](completed/330-context-menu-and-gitmap-ui-enhancements.md) — Plan 330: Context Menu and GitMap UI Enhancements: File Right-Click Integration (*\shell\gitmap), Explorer Cascade Flattening, PowerShell Duplicate Cleaner, Help Visibility, and Web UI Settings (Spec: [330-context-menu-and-gitmap-ui-enhancements](../../02-spec/21-app/330-context-menu-and-gitmap-ui-enhancements/00-master-audit-ledger.md)) (Status: `completed`)
- [329-gitmap-context-menu-and-version-fix.md](completed/329-gitmap-context-menu-and-version-fix.md) — Plan 329: GitMap Context Menu and Version Fix: Update Downgrade Guard, Sibling Probe Resolution, Cross-Platform Right-Click Context Menu, TortoiseGit-Inspired Show Log Dual-Pane UI, Beyond Compare Integration, and Background Install Wizard (Spec: [329-gitmap-context-menu-and-version-fix](../../02-spec/21-app/329-gitmap-context-menu-and-version-fix/00-master-audit-ledger.md)) (Status: `completed`)
- [316-repository-management-and-testing-instructions.md](316-repository-management-and-testing-instructions.md) — Repository Management and Testing Instructions across AGM, GitMap, img-pdf, and yt-manage (Spec: [316-repository-management-and-testing-instructions](../../02-spec/21-app/316-repository-management-and-testing-instructions/01-architecture-spec.md)) (Status: `active`)
- [gitmap-enhancement-and-issue-resolution.md](gitmap-enhancement-and-issue-resolution.md) — GitMap Enhancement and Issue Resolution: Database Read Introspection & 5-10s Metadata Cache, Visualizer UI Binary Cell Formatting & Drawer, Storytelling-v1 README Casing Remediation & Repo-Cache Recipe, Remote Node Clone GitHub Desktop & VS Code Registration, Remote SSH WinCredMan Resolution, and Machine Set Interactive Prompt (Spec: [gitmap-enhancement-and-issue-resolution](../../02-spec/21-app/gitmap-enhancement-and-issue-resolution/01-architecture-spec.md)) (Status: `active`)
- [goemailsender-connection-and-testing.md](goemailsender-connection-and-testing.md) — GoEmailSender Connection, Multi-Profile Fallback & Fleet Nodes Delegation (Spec: [goemailsender-connection-and-testing](../../02-spec/21-app/goemailsender-connection-and-testing/01-architecture-spec.md)) (Status: `in_progress`)
- [update-antigravity-ide.md](update-antigravity-ide.md) — Antigravity Desktop IDE Update Integration: GitMap CLI `update antigravity` / `update ide` commands, scripts fixture (`scripts-fixer`) `update` verb, Antigravity tools updater script, and live Ubuntu system updates (Spec: [update-antigravity-ide](../../02-spec/21-app/update-antigravity-ide/readme.md)) (Status: `active`)
- [git-prime-machine-os-version-and-ui-help-docs-improvement.md](git-prime-machine-os-version-and-ui-help-docs-improvement.md) — Git Prime Machine OS Version and UI Help Docs Improvement: Structured OS Version Object & JSON Breakdown, HD Docs Web UI Command Expansion, and Merge-AI Parameterized Examples (Spec: [git-prime-machine-os-version-and-ui-help-docs-improvement](../../02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/01-architecture-spec.md)) (Status: `active`)
- [219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md](219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md) — Ubuntu Pull-All Remediation Engine, Oh-My-Zsh & Default System Directory Scanner Exclusions, Cross-OS Path Healing, Failure Subtree Formatting (Spec: [03-git-operations-and-pull](../../02-spec/21-app/03-git-operations-and-pull/01-architecture-spec.md))

---

## 2. Completed Plans (Consolidated Milestones `01`–`32`)

All completed work has been merged into 32 dense, authoritative milestone summaries preserving 100% of verified technical outcomes, Go type contracts (`*appfault.AppError`, `Result[T]`), and architectural invariants:

### Generational Milestones (Historical Foundations)
- [01-generation-1-core-foundation-milestones.md](completed/01-generation-1-core-foundation-milestones.md) — Core Architecture, Scanner, Types, Database, Git Operations, Terminal UI, Installers & Chrome Profile Vault (Synthesizing Historical Milestones 01–15)
- [02-generation-2-fleet-automation-and-modularization.md](completed/02-generation-2-fleet-automation-and-modularization.md) — Cluster Fleet, Macro Engine, Nuclear Modularization, Smart Test Runner, DB Locking, Antigravity Packaging & Quality Gates (Synthesizing Historical Milestones 16–27)

### Modern Milestone Sequence (Milestones 03–32)
- [03-agy-prompts-templates-and-rerun-suite.md](completed/03-agy-prompts-templates-and-rerun-suite.md) — AGY Prompt Templates, Decision Logs, Non-Destructive Rerun Commands & Prompt Sync
- [04-auto-aliasing-and-error-storage-reset.md](completed/04-auto-aliasing-and-error-storage-reset.md) — Auto-Aliasing Migration, Internal Errors Database & Safe DB Reset
- [05-ssh-exec-copy-mv-env-and-rm-sync-resilience.md](completed/05-ssh-exec-copy-mv-env-and-rm-sync-resilience.md) — Remote SSH Operations, Streaming Upload, Deploy Keys & Sync Modes
- [06-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md](completed/06-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md) — Pipeline Path Normalization, Cache Clearing & Accurate ETA Forecasting
- [07-pipeline-deep-eta-and-repo-folder-parity.md](completed/07-pipeline-deep-eta-and-repo-folder-parity.md) — Deep ETA Profiling, Heatmap Telemetry & Unit Test Traceback Extraction
- [08-ai-scripts-engine-ssh-authkey-and-feature-parity.md](completed/08-ai-scripts-engine-ssh-authkey-and-feature-parity.md) — Automation Scripts Engine, SSH Authorized Keys & Platform Parity
- [09-generic-ai-scripts-creator-and-scaffolding-engine.md](completed/09-generic-ai-scripts-creator-and-scaffolding-engine.md) — Scaffolding Engine, Toolchain Scripts & Automation Runners
- [10-native-automation-engine-lazy-regex-and-benchmarks.md](completed/10-native-automation-engine-lazy-regex-and-benchmarks.md) — Native Automation Engine, Lazy Regex Compilation & AUM Search Benchmarks
- [11-fleet-nodes-ping-envelope-and-remote-clone.md](completed/11-fleet-nodes-ping-envelope-and-remote-clone.md) — Unified Fleet Nodes Command, Dual-Stack Ping, Typed JSON Envelopes & Remote Clone
- [12-pas-worker-concurrency-pull-cache-and-split-db.md](completed/12-pas-worker-concurrency-pull-cache-and-split-db.md) — GitMap PAS Formula, Worker Concurrency, Ignore Engine & Split-DB Repo Cache
- [13-semantic-flat-commit-and-macro-execution-resilience.md](completed/13-semantic-flat-commit-and-macro-execution-resilience.md) — Semantic Flat Commit Suite (`gitmap c`), Auto-Staging & Macro Idempotency
- [14-ssh-password-interception-and-rsa-credential-vault.md](completed/14-ssh-password-interception-and-rsa-credential-vault.md) — Masked Password Interception, RSA-OAEP Salt Credential Vault & Push Self-Healing
- [15-ubuntu-fleet-migration-workstation-governance-and-customization.md](completed/15-ubuntu-fleet-migration-workstation-governance-and-customization.md) — Windows to Ubuntu Fleet Migration, GNOME 140% Scaling & Workstation Governance
- [16-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md](completed/16-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md) — AGY UI Plugins, Skills Sync, Terminal Freeze Remediation & Multi-IDE Scan Sync
- [17-codebase-review-remediation-and-preconsolidation-baseline.md](completed/17-codebase-review-remediation-and-preconsolidation-baseline.md) — Root Clutter Purge, README Index Reduction, Security Token Purge & Baseline Release
- [18-app-spec-and-completed-plans-consolidation-and-reduction.md](completed/18-app-spec-and-completed-plans-consolidation-and-reduction.md) — Application Specifications Consolidation, 8 Canonical Clusters & Memory Compaction
- [19-deep-spec-consolidation-and-canonical-reduction.md](completed/19-deep-spec-consolidation-and-canonical-reduction.md) — Deep Spec Consolidation, Canonical 8-Cluster Reduction & Memory Unification
- [20-git-history-purge-and-undo.md](completed/20-git-history-purge-and-undo.md) — Git History Purge & Undo Engine with Pre-Flight Graph Diff & SplitDB Journal
- [21-aria2c-download-and-multi-user-switching.md](completed/21-aria2c-download-and-multi-user-switching.md) — Aria2c Download Method, Smart Update Engine, Multi-User Project Switching & GitHub CLI Auth Sync
- [22-path-normalization-and-backup-guard-suite.md](completed/22-path-normalization-and-backup-guard-suite.md) — Forward Slash Path Normalization, Backup Branch Clean Tree Guard & E2E Verification
- [23-agm-update-tui-and-update-ui-alignment.md](completed/23-agm-update-tui-and-update-ui-alignment.md) — AGM Update TUI Fix, UI Color Alignment & Inline Release Notes
- [24-pull-remediation-automerge-and-large-files.md](completed/24-pull-remediation-automerge-and-large-files.md) — Pull Remediation Optimization, Auto-Merge Conflict Resolution & Large Files Split Engine
- [25-pipeline-parallel-download-and-heartbeat.md](completed/25-pipeline-parallel-download-and-heartbeat.md) — Pipeline Parallel Download, Heartbeat Telemetry & Caching
- [26-codebase-improvements-and-fleet-sync.md](completed/26-codebase-improvements-and-fleet-sync.md) — Codebase Improvement Suite, Prompt Catalog & Fleet Sync
- [27-result-wrapper-and-monadic-containers.md](completed/27-result-wrapper-and-monadic-containers.md) — Monadic Result Containers, AppError Architecture & Stutter Eradication
- [28-coding-guidelines-and-code-hygiene-suite.md](completed/28-coding-guidelines-and-code-hygiene-suite.md) — Coding Guidelines, Condition Flattening & Code Hygiene Suite
- [29-workspace-pull-all-and-ssh-fleet-acceleration.md](completed/29-workspace-pull-all-and-ssh-fleet-acceleration.md) — Workspace-Wide Default Pull & Accelerated SSH Fleet Remote Execution
- [30-commit-tracking-heatmap-url-and-database.md](completed/30-commit-tracking-heatmap-url-and-database.md) — Commit Tracking, Heatmap Telemetry, URL Correlation & Split-DB
- [31-historical-commands-stores-and-email-daemon.md](completed/31-historical-commands-stores-and-email-daemon.md) — Historical Commands, Stores Architecture & Email Automation
- [32-installer-discovery-warning-suppression-and-release.md](completed/32-installer-discovery-warning-suppression-and-release.md) — Installer Discovery Warning Suppression, Dual Installer Sync & Release v6.596.0

---

## 2b. Recently Completed (awaiting milestone consolidation)

- [349-gitmap-timezone-and-commit-scheduling.md](completed/349-gitmap-timezone-and-commit-scheduling.md) — GitMap Timezone and Commit Scheduling: 3-tier timezone resolution (--tz > config > local), weekday < 08:00 AM cutoff policy, weekend 24h unconstrained scheduling, interactive & one-liner modes, Settings UI, terminal rich help, and historical prompt/skill renewal (Spec: [349-gitmap-timezone-and-commit-scheduling](../../02-spec/21-app/349-gitmap-timezone-and-commit-scheduling/readme.md)) (Status: `completed`)
- [348-recommit-and-skill-creation-instructions.md](completed/348-recommit-and-skill-creation-instructions.md) — Recommit & Skill Creation Instructions: GitMap time-restriction on-call skill, AI automation bot on-call skill, Letterly prompt 12, and V7 execution engine (Spec: [348-recommit-and-skill-creation-instructions](../../02-spec/21-app/348-recommit-and-skill-creation-instructions/readme.md)) (Status: `completed`)
- [319-package-constants-and-enums-dry-restructuring.md](completed/319-package-constants-and-enums-dry-restructuring.md) — Package Constants & Enums DRY Restructuring: Modular enum types (*Type suffix), centralized constants (pkg/constants/constants.go), and 5-package rollout (cmdpull, cloner, cmdclone, store, cmdcommit) (Spec: [319-package-constants-and-enums-dry-restructuring](../../02-spec/21-app/319-package-constants-and-enums-dry-restructuring/readme.md)) (Status: `completed`)
- [320-package-constants-and-enums-dry-restructuring.md](completed/320-package-constants-and-enums-dry-restructuring.md) — Package Constants & Enums DRY Restructuring Phase 2: Modular enum types (*Type suffix), leaf enums/constants packages, error return sovereignty, and 5-package rollout (cmdsync, cmdplans, cmdstatus, cmdlogin, cmdcargo) (Status: `completed`)
- [321-package-constants-and-enums-dry-restructuring.md](completed/321-package-constants-and-enums-dry-restructuring.md) — Package Constants & Enums DRY Restructuring Phase 3: Modular enum types (*Type suffix), leaf enums/constants packages, error return sovereignty, and 5-package rollout (cmdconfig, cmddiff, cmdsearch, cmdbranch, cmdhealth) (Status: `completed`)
- [322-package-constants-and-enums-dry-restructuring.md](completed/322-package-constants-and-enums-dry-restructuring.md) — Package Constants & Enums DRY Restructuring Phase 4: Modular enum types (*Type suffix), leaf enums/constants packages, error return sovereignty, and 5-package rollout (cmdstats, cmdhooks, cmdprune, cmdadd, cmddocs) (Status: `completed`)
- [323-package-constants-and-enums-dry-restructuring.md](completed/323-package-constants-and-enums-dry-restructuring.md) — Package Constants & Enums DRY Restructuring Phase 5: Modular enum types (*Type suffix), leaf enums/constants packages, error return sovereignty, and 5-package rollout (cmdundo, cmdrevert, cmdstale, cmdconflict, cmdcopy) (Status: `completed`)

---

## 3. Pending Plans

- [332-gitmap-ui-and-functionality-improvements.md](pending/332-gitmap-ui-and-functionality-improvements.md) — GitMap UI and Functionality Improvements: Show Log folder inspection, modern scrollable responsive UI, server exit button, web context menu, working tree diff with stack traces, context menu Install & Repos submenus, and terminal help expansion (Spec: [332-gitmap-ui-and-functionality-improvements](../../02-spec/21-app/332-gitmap-ui-and-functionality-improvements/00-master-audit-ledger.md)) (Status: `pending`)
- [331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation.md](pending/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation.md) — Git Commit Prefix Inconsistency & Raw Git Fallback Remediation: Formatter engine release module canonicalization, 47-reconcile script commit standardization, and rules/skills sync (Spec: [331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation](../../02-spec/21-app/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation/00-master-audit-ledger.md)) (Status: `pending`)
- [328-flexible-path-parsing-and-asset-image-reading.md](pending/328-flexible-path-parsing-and-asset-image-reading.md) — Flexible Path Parsing and Asset Image Reading: Support for @[...] bracket paths, markdown image/link URLs, file:/// schemes, list/comment stripping, cross-repo relative paths, safe image inspection in cat, and minor release v6.612.0 (Spec: [328-flexible-path-parsing-and-asset-image-reading](../../02-spec/21-app/328-flexible-path-parsing-and-asset-image-reading/00-master-audit-ledger.md)) (Status: `completed`)
- [327-heatmap-status-improvement-and-ai-confusion-handling.md](pending/327-heatmap-status-improvement-and-ai-confusion-handling.md) — Heatmap Status Improvement and AI Confusion Handling: Root SQLite sequence folder caching, 1-to-N normalized AiConfusionFile schema, auto-syncing status workflow, relative file context querying, ghost repo name candidate sanitization, and minor release ceremony (Spec: [327-heatmap-status-improvement-and-ai-confusion-handling](../../02-spec/21-app/327-heatmap-status-improvement-and-ai-confusion-handling/00-master-audit-ledger.md)) (Status: `pending`)
- [326-heatmap-status-improvement-and-ai-confusion-database.md](pending/326-heatmap-status-improvement-and-ai-confusion-database.md) — Heatmap Status Improvement and AI Confusion Database: Visual padding overhaul, task integration, JSON telemetry, safe stashing, repo-cache AI confusion database, and commit prefix remediation (Spec: [326-heatmap-status-improvement-and-ai-confusion-database](../../02-spec/21-app/326-heatmap-status-improvement-and-ai-confusion-database/00-master-audit-ledger.md)) (Status: `completed`)
- [318-gitmap-command-and-coding-guideline-synchronization.md](pending/318-gitmap-command-and-coding-guideline-synchronization.md) — GitMap Command and Coding Guideline Synchronization: Alternative Commands in GitMap Skill, Third-Bracket Commit Elimination, Magic Strings Refactoring to Domain Enums, and V6 Execute Prompt Error Management Embedding (Spec: [318-gitmap-command-and-coding-guideline-synchronization](../../02-spec/21-app/318-gitmap-command-and-coding-guideline-synchronization/00-overview.md)) (Status: `pending`)

---

## 4. Active Subtasks Directories

Only active subtask directories are maintained under `.ai-memory/plans/subtasks/`:

- [332-gitmap-ui-and-functionality-improvements](subtasks/332-gitmap-ui-and-functionality-improvements/) — Show Log folder inspection, modern scrollable responsive UI, server exit button, web context menu, working tree diff with stack traces, context menu Install & Repos submenus, and terminal help expansion
- [331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation](subtasks/331-git-commit-prefix-inconsistency-and-raw-git-fallback-remediation/) — Formatter engine release module canonicalization, 47-reconcile script commit standardization, and rules/skills sync
- [328-flexible-path-parsing-and-asset-image-reading](subtasks/328-flexible-path-parsing-and-asset-image-reading/) — Support for @[...] bracket paths, markdown image/link URLs, file:/// schemes, list/comment stripping, cross-repo relative paths, safe image inspection in cat, and minor release v6.612.0
- [327-heatmap-status-improvement-and-ai-confusion-handling](subtasks/327-heatmap-status-improvement-and-ai-confusion-handling/) — Root SQLite sequence folder caching, 1-to-N normalized AiConfusionFile schema, auto-sync status workflow, relative file context querying, ghost repo candidate sanitization, and release ceremony
- [326-heatmap-status-improvement-and-ai-confusion-database](subtasks/326-heatmap-status-improvement-and-ai-confusion-database/) — Heatmap status redesign, completed/pending task display, work directory guard, safe stashing, repo-cache AI confusion SQLite DB, and commit prefix remediation
- [319-package-constants-and-enums-dry-restructuring](subtasks/319-package-constants-and-enums-dry-restructuring/) — Enums and constants DRY restructuring across cmdpull, cloner, cmdclone, store, and cmdcommit
- [318-gitmap-command-and-coding-guideline-synchronization](subtasks/318-gitmap-command-and-coding-guideline-synchronization/) — Alternative commands in GitMap skill, third bracket elimination in commit tools, magic strings replacement, V6 execute prompt error management, and multi-repo sync

- [316-repository-management-and-testing-instructions](subtasks/316-repository-management-and-testing-instructions/) — AGM autoswitch, fast-forward rerun, gitmap commands, img-pdf, yt-manage
- [219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging](subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/) — Ubuntu Pull-All Remediation Engine & Ignore Rules
- [226-ubuntu-cursor-memories-conversations-and-projects-migration](subtasks/226-ubuntu-cursor-memories-conversations-and-projects-migration/) — Cursor IDE Internal State, Memories & Conversations
- [56-vmware-hardware-batch-and-macro-orchestration](subtasks/56-vmware-hardware-batch-and-macro-orchestration/) — VMware Hardware Customization & Multi-VM Batch Operations
- [75-gitmap-u1-ubuntu-agm-fleet-integration](subtasks/75-gitmap-u1-ubuntu-agm-fleet-integration/) — GitMap U1 Ubuntu Fleet Integration & AGM Migration
- [git-prime-machine-os-version-and-ui-help-docs-improvement](subtasks/git-prime-machine-os-version-and-ui-help-docs-improvement/) — Prime Machine OS Version Breakdown & Help Docs Expansion
- [gitmap-enhancement-and-issue-resolution](subtasks/gitmap-enhancement-and-issue-resolution/) — Database Read Introspection, Visualizer Drawer & Storytelling-v1 Casing
- [gitmap-node-and-prompt-system-issues](subtasks/gitmap-node-and-prompt-system-issues/) — GitHub Node Aliasing, Prompt Tree Suite & Telegram Sorting
- [goemailsender-connection-and-testing](subtasks/goemailsender-connection-and-testing/) — GoEmailSender Connection & Multi-Profile Fallback
- [muse-memory-io-and-gitmap-integration](subtasks/muse-memory-io-and-gitmap-integration/) — Muse Memory IO, 7-Day Digest & Training Vault
- [update-antigravity-ide](subtasks/update-antigravity-ide/) — Antigravity Desktop IDE Update CLI Commands

## Recent Completed Tasks Register

- 2026-10-29 — Recorded the warm-up gate session. Sequences 08 through 13 remain pending catalogs. The live gate engine is in go-email-reader commit `f3864ee`.
