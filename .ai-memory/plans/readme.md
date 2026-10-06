# Master Plans Registry (`.ai-memory/plans`)

This registry catalogs the active, completed, and pending execution plans for GitMap. In accordance with Plan 235 ([Spec 235](../../02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/01-architecture-spec.md)), historical completed plans have been consolidated into **42 Authoritative Milestones** under `.ai-memory/plans/completed/`.

---

## 1. Active Plans

The following plans represent open and actively maintained development streams:

- [235-app-spec-and-completed-plans-consolidation-and-reduction.md](235-app-spec-and-completed-plans-consolidation-and-reduction.md) — Application Specifications & Completed Plans Consolidation & Memory Reduction (Spec: [235](../../02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/01-architecture-spec.md))
- [219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md](219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md) — Ubuntu Pull-All Remediation Engine, Oh-My-Zsh & Default System Directory Scanner Exclusions, Cross-OS Path Healing, Failure Subtree Formatting (Spec: [219](../../02-spec/21-app/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/01-architecture-spec.md))
- [226-ubuntu-cursor-memories-conversations-and-projects-migration.md](226-ubuntu-cursor-memories-conversations-and-projects-migration.md) — Migrate Cursor IDE internal state, memories, conversations, and projects from local Windows workstation to Ubuntu node u1 (Spec: [226](../../02-spec/21-app/226-ubuntu-cursor-memories-conversations-and-projects-migration/01-architecture-spec.md))
- [227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md](227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md) — GitMap Cursor Fleet Delegation, Settings UI & Cross-Node Path Suggestions
- [225-ubuntu-cursor-start-menu-dock-position-and-profile-migration.md](225-ubuntu-cursor-start-menu-dock-position-and-profile-migration.md) — Ubuntu Cursor Start Menu Dock Position & Profile Migration
- [223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md](223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md) — Ubuntu Cursor GitMap Python Git Tracing & AUM Search
- [217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md](217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md) — Antigravity Fleet Parity, UI Theme Seeds, Turbo Permission Preset, 4 Plugins & 43 Skills Sync
- [216-gitmap-prompting-freeze-and-suggestion-engine-fix.md](216-gitmap-prompting-freeze-and-suggestion-engine-fix.md) — GitMap Terminal Prompt Input Freeze Remediation & Typo Suggestion Overhaul
- [214-ubuntu-fleet-automation-and-workstation-governance.md](214-ubuntu-fleet-automation-and-workstation-governance.md) — Ubuntu Fleet Automation, 74 Workspaces Clone & Parity Audit, GNOME Ergonomics

---

## 2. Completed Plans (Consolidated Milestones `01`–`42`)

All completed work has been merged into 42 dense, authoritative milestone summaries preserving 100% of verified technical outcomes, Go type contracts (`*appfault.AppError`, `Result[T]`), and architectural invariants:

### Milestones 01–27 (Historical Core)
- [01-core-architecture-and-scanner.md](completed/01-core-architecture-and-scanner.md) — Core GitMap Architecture & Scanner Engine
- [02-project-detection-and-formatting.md](completed/02-project-detection-and-formatting.md) — Project Language Detection & Multi-Format Exporters
- [03-git-operations-and-cloning.md](completed/03-git-operations-and-cloning.md) — Parallel Cloning & High-Speed Manifest Ingestion
- [04-ignore-rules-and-cache.md](completed/04-ignore-rules-and-cache.md) — GitIgnore Patterns & Cache Storage
- [05-pull-and-branch-management.md](completed/05-pull-and-branch-management.md) — Pull Synchronization & Branch Operations
- [06-split-database-engine.md](completed/06-split-database-engine.md) — Three-Tier SQLite Split-DB Storage Architecture
- [07-ssh-delegation-and-remote-exec.md](completed/07-ssh-delegation-and-remote-exec.md) — SSH Delegation & Cross-Node Remote Execution
- [08-fleet-nodes-and-discovery.md](completed/08-fleet-nodes-and-discovery.md) — Fleet Discovery & Multi-Target Clustering
- [09-credential-vault-and-security.md](completed/09-credential-vault-and-security.md) — Credential Vault & RSA Encryption
- [10-antigravity-integration-and-prompts.md](completed/10-antigravity-integration-and-prompts.md) — Antigravity SDK & Prompt Automation
- [11-pipeline-diagnostics-and-rca.md](completed/11-pipeline-diagnostics-and-rca.md) — CI/CD Pipeline Telemetry & Automated RCA
- [12-terminal-ui-and-cli-styling.md](completed/12-terminal-ui-and-cli-styling.md) — Terminal UI Layouts, ANSI Styling & ANSI Tables
- [13-cross-platform-runners-and-packaging.md](completed/13-cross-platform-runners-and-packaging.md) — Cross-Platform Shell Runners & Binary Packaging
- [14-ide-synchronization-and-projects.md](completed/14-ide-synchronization-and-projects.md) — VS Code & Multi-IDE Project Registration
- [15-coding-guidelines-and-quality-gates.md](completed/15-coding-guidelines-and-quality-gates.md) — Coding Standards & Quality Verification Gates
- [16-database-migrations-and-schemas.md](completed/16-database-migrations-and-schemas.md) — Database Schema Evolution & PascalCase Tables
- [17-macro-automation-and-replay.md](completed/17-macro-automation-and-replay.md) — Macro Recording, Idempotency & Replay Engine
- [18-performance-benchmarks-and-aum.md](completed/18-performance-benchmarks-and-aum.md) — AUM Sub-Millisecond Search & Micro-Benchmarks
- [19-fleet-sync-and-distribution.md](completed/19-fleet-sync-and-distribution.md) — Fleet Distribution & Remote Archive Deployment
- [20-error-handling-and-appfault.md](completed/20-error-handling-and-appfault.md) — Structured Errors (`*appfault.AppError`) & Resilient Exits
- [21-semver-release-ceremony.md](completed/21-semver-release-ceremony.md) — Automated Release Orchestration & SemVer Governance
- [22-fast-file-reader-and-caching.md](completed/22-fast-file-reader-and-caching.md) — Fast File Reader Engine & LRU Cache Hierarchy
- [23-vmware-and-hypervisor-controls.md](completed/23-vmware-and-hypervisor-controls.md) — VMware Workstation CLI & Coordinated Lifecycle
- [24-firewall-and-port-management.md](completed/24-firewall-and-port-management.md) — Cross-OS Firewall Management & Port Inspection
- [25-cross-platform-path-normalization.md](completed/25-cross-platform-path-normalization.md) — Cross-OS Path Normalization & Collation
- [26-ai-agent-task-database.md](completed/26-ai-agent-task-database.md) — AI Agent Task Orchestrator & Task Database
- [27-secrets-governance-and-vault.md](completed/27-secrets-governance-and-vault.md) — Secrets Governance, Repo-Secrets & Token Purge

### Milestones 28–42 (Consolidated Feature Generations)
- [28-agy-prompts-templates-and-rerun-suite.md](completed/28-agy-prompts-templates-and-rerun-suite.md) — AGY Prompt Templates, Decision Logs, Non-Destructive Rerun Commands & Prompt Sync
- [29-auto-aliasing-and-error-storage-reset.md](completed/29-auto-aliasing-and-error-storage-reset.md) — Auto-Aliasing Migration, Internal Errors Database & Safe DB Reset
- [30-ssh-exec-copy-mv-env-and-rm-sync-resilience.md](completed/30-ssh-exec-copy-mv-env-and-rm-sync-resilience.md) — Remote SSH Operations, Streaming Upload, Deploy Keys & Sync Modes
- [31-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md](completed/31-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md) — Pipeline Path Normalization, Cache Clearing & Accurate ETA Forecasting
- [32-pipeline-deep-eta-and-repo-folder-parity.md](completed/32-pipeline-deep-eta-and-repo-folder-parity.md) — Deep ETA Profiling, Heatmap Telemetry & Unit Test Traceback Extraction
- [33-ai-scripts-engine-ssh-authkey-and-feature-parity.md](completed/33-ai-scripts-engine-ssh-authkey-and-feature-parity.md) — Automation Scripts Engine, SSH Authorized Keys & Platform Parity
- [34-generic-ai-scripts-creator-and-scaffolding-engine.md](completed/34-generic-ai-scripts-creator-and-scaffolding-engine.md) — Scaffolding Engine, Toolchain Scripts & Automation Runners
- [35-native-automation-engine-lazy-regex-and-benchmarks.md](completed/35-native-automation-engine-lazy-regex-and-benchmarks.md) — Native Automation Engine, Lazy Regex Compilation & AUM Search Benchmarks
- [36-fleet-nodes-ping-envelope-and-remote-clone.md](completed/36-fleet-nodes-ping-envelope-and-remote-clone.md) — Unified Fleet Nodes Command, Dual-Stack Ping, Typed JSON Envelopes & Remote Clone
- [37-pas-worker-concurrency-pull-cache-and-split-db.md](completed/37-pas-worker-concurrency-pull-cache-and-split-db.md) — GitMap PAS Formula, Worker Concurrency, Ignore Engine & Split-DB Repo Cache
- [38-semantic-flat-commit-and-macro-execution-resilience.md](completed/38-semantic-flat-commit-and-macro-execution-resilience.md) — Semantic Flat Commit Suite (`gitmap c`), Auto-Staging & Macro Idempotency
- [39-ssh-password-interception-and-rsa-credential-vault.md](completed/39-ssh-password-interception-and-rsa-credential-vault.md) — Masked Password Interception, RSA-OAEP Salt Credential Vault & Push Self-Healing
- [40-ubuntu-fleet-migration-workstation-governance-and-customization.md](completed/40-ubuntu-fleet-migration-workstation-governance-and-customization.md) — Windows to Ubuntu Fleet Migration, GNOME 140% Scaling & Workstation Governance
- [41-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md](completed/41-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md) — AGY UI Plugins, Skills Sync, Terminal Freeze Remediation & Multi-IDE Scan Sync
- [42-codebase-review-remediation-and-preconsolidation-baseline.md](completed/42-codebase-review-remediation-and-preconsolidation-baseline.md) — Root Clutter Purge, README Index Reduction, Security Token Purge & Baseline Release

---

## 3. Pending Plans

The following plans are protected by the Strict Pending Isolation Boundary and remain untouched:

- [01-ports-and-ssh-enablement.md](pending/01-ports-and-ssh-enablement.md) — Ports Inspection, Firewall Audit & Cross-Platform OpenSSH Daemon Management
- [56-vmware-hardware-batch-and-macro-orchestration.md](pending/56-vmware-hardware-batch-and-macro-orchestration.md) — VMware Hardware Customization, Multi-VM Batch Operations & Macro Orchestration
- [67-nodes-cfr-remote-fleet-clone-enhancement.md](pending/67-nodes-cfr-remote-fleet-clone-enhancement.md) — Nodes CFR Remote Fleet Clone Enhancement & Async Delegation
- [75-gitmap-u1-ubuntu-agm-fleet-integration.md](pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md) — GitMap U1 Ubuntu Fleet Integration, AGM Migration & Cross-OS Automation

---

## 4. Active Subtasks Directories

Only active subtask directories are maintained under `.ai-memory/plans/subtasks/` (11 directories, 55 files total):

- `subtasks/235-app-spec-and-completed-plans-consolidation-and-reduction/` (5 subtasks)
- `subtasks/219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/` (6 subtasks)
- `subtasks/226-ubuntu-cursor-memories-conversations-and-projects-migration/` (4 subtasks)
- `subtasks/56-vmware-hardware-batch-and-macro-orchestration/` (6 subtasks)
- `subtasks/60-gitmap-pas-command-fix/` (14 subtasks)
- `subtasks/67-ssh-password-interception-and-rsa-vault/` (2 subtasks)
- `subtasks/68-stats-float-scan-and-ssh/` (3 subtasks)
- `subtasks/71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix/` (2 subtasks)
- `subtasks/74-refresh-token-and-vmware-spec-validation/` (2 subtasks)
- `subtasks/75-gitmap-u1-ubuntu-agm-fleet-integration/` (6 subtasks)
- `subtasks/91-ssh-join-common-scan-and-agy-rop/` (5 subtasks)
