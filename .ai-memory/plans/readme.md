# Master Plans Registry (`.ai-memory/plans`)

This registry catalogs the active, completed, and pending execution plans for GitMap. In accordance with Plan 236 ([Spec 236](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md)), historical completed plans have been consolidated into **32 Authoritative Milestones** under `.ai-memory/plans/completed/`.

---

## 1. Active Plans

The following plans represent open and actively maintained development streams:

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

*(All historical completed plans have been consolidated into Milestones 01–32 under `completed/`.)*

---

## 3. Pending Plans

None. All active execution tracks are maintained directly under `.ai-memory/plans/subtasks/`.

---

## 4. Active Subtasks Directories

Only active subtask directories are maintained under `.ai-memory/plans/subtasks/`:

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
