# Master Plans Registry (`.ai-memory/plans`)

This registry catalogs the active, completed, and pending execution plans for GitMap. In accordance with Plan 236 ([Spec 236](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md)), historical completed plans have been consolidated into **18 Authoritative Milestones** under `.ai-memory/plans/completed/`.

---

## 1. Active Plans

The following plans represent open and actively maintained development streams:

- [246-parallel-fix-command.md](246-parallel-fix-command.md) — Native parallel `gitmap autofix` command (8 AI fixer scripts ported to Go), LLM train Heal & Fix phase, skills docs, settings UI/UX overhaul (Spec: [246](../../02-spec/21-app/246-parallel-fix-command/01-overview.md))
- [240-mcp-server-ai-analysis-agm-nodes-and-git-cache.md](240-mcp-server-ai-analysis-agm-nodes-and-git-cache.md) — GitMap AI MCP Server, AI Analysis Engine, Split SQLite Storage, AGM Node Deployment & Git Cache (Spec: [240](../../02-spec/21-app/240-mcp-server-ai-analysis-agm-nodes-and-git-cache/01-architecture-spec.md))
- [239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills.md](239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills.md) — GitMap Runner, Audit DB, Supabase Vault, PE-AI & Muse Skills (Spec: [239](../../02-spec/21-app/239-gitmap-runner-audit-supabase-vault-pe-ai-and-muse-skills/01-architecture-spec.md))
- [236-deep-spec-consolidation-and-canonical-reduction.md](236-deep-spec-consolidation-and-canonical-reduction.md) — Deep Spec Consolidation & Canonical Reduction (Spec: [236](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md))
- [219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md](219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md) — Ubuntu Pull-All Remediation Engine, Oh-My-Zsh & Default System Directory Scanner Exclusions, Cross-OS Path Healing, Failure Subtree Formatting (Spec: [03-git-operations-and-pull](../../02-spec/21-app/03-git-operations-and-pull/01-architecture-spec.md))
- [226-ubuntu-cursor-memories-conversations-and-projects-migration.md](226-ubuntu-cursor-memories-conversations-and-projects-migration.md) — Migrate Cursor IDE internal state, memories, conversations, and projects from local Windows workstation to Ubuntu node u1 (Spec: [05-antigravity-and-ide](../../02-spec/21-app/05-antigravity-and-ide/01-architecture-spec.md))
- [227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md](227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md) — GitMap Cursor Fleet Delegation, Settings UI & Cross-Node Path Suggestions
- [225-ubuntu-cursor-start-menu-dock-position-and-profile-migration.md](225-ubuntu-cursor-start-menu-dock-position-and-profile-migration.md) — Ubuntu Cursor Start Menu Dock Position & Profile Migration
- [223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md](223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md) — Ubuntu Cursor GitMap Python Git Tracing & AUM Search
- [217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md](217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md) — Antigravity Fleet Parity, UI Theme Seeds, Turbo Permission Preset, 4 Plugins & 43 Skills Sync
- [216-gitmap-prompting-freeze-and-suggestion-engine-fix.md](216-gitmap-prompting-freeze-and-suggestion-engine-fix.md) — GitMap Terminal Prompt Input Freeze Remediation & Typo Suggestion Overhaul
- [214-ubuntu-fleet-automation-and-workstation-governance.md](214-ubuntu-fleet-automation-and-workstation-governance.md) — Ubuntu Fleet Automation, 74 Workspaces Clone & Parity Audit, GNOME Ergonomics

---

## 2. Completed Plans (Consolidated Milestones `01`–`20`)

All completed work has been merged into 20 dense, authoritative milestone summaries preserving 100% of verified technical outcomes, Go type contracts (`*appfault.AppError`, `Result[T]`), and architectural invariants:

### Generational Milestones (Historical Foundations)
- [01-generation-1-core-foundation-milestones.md](completed/01-generation-1-core-foundation-milestones.md) — Core Architecture, Scanner, Types, Database, Git Operations, Terminal UI, Installers & Chrome Profile Vault (Synthesizing Historical Milestones 01–15)
- [02-generation-2-fleet-automation-and-modularization.md](completed/02-generation-2-fleet-automation-and-modularization.md) — Cluster Fleet, Macro Engine, Nuclear Modularization, Smart Test Runner, DB Locking, Antigravity Packaging & Quality Gates (Synthesizing Historical Milestones 16–27)

### Modern Milestone Sequence (Milestones 03–20)
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
- [21-aria2c-download-smart-update-stats-subnode-remediation.md](completed/242-aria2c-download-smart-update-stats-subnode-remediation.md) — aria2c Download Method, Smart Update Engine, 5-Tag Fallback, Diffstat Isolation & Remediation Tree UI (Spec: [242](../../02-spec/21-app/242-aria2c-download-smart-update-stats-subnode-remediation/01-architecture-spec.md))
- [22-e2e-verify-program-243-changes.md](completed/244-e2e-verify-program-243-changes.md) — End-to-end verification of program 243 (backup-branch, help displayer, suggestions, enforcement): 26/27 assertions pass, 1 real defect found with RCA (Spec: [244](../../02-spec/21-app/244-e2e-verify-program-243-changes/01-test-scope-and-plan.md))
- [23-fix-backup-branch-lasterror-guard.md](completed/245-fix-backup-branch-lasterror-guard.md) — Fix backup-branch's clean-tree guard to ignore untracked gitmap state (.gitmap/): the last_error.log side effect no longer defeats --force; F01–F06 all pass (Spec: [245](../../02-spec/21-app/245-fix-backup-branch-lasterror-guard/01-fix-scope-and-plan.md))

---

## 3. Pending Plans

The following plans are protected by the Strict Pending Isolation Boundary and remain untouched:

- [01-ports-and-ssh-enablement.md](pending/01-ports-and-ssh-enablement.md) — Ports Inspection, Firewall Audit & Cross-Platform OpenSSH Daemon Management
- [56-vmware-hardware-batch-and-macro-orchestration.md](pending/56-vmware-hardware-batch-and-macro-orchestration.md) — VMware Hardware Customization, Multi-VM Batch Operations & Macro Orchestration
- [67-nodes-cfr-remote-fleet-clone-enhancement.md](pending/67-nodes-cfr-remote-fleet-clone-enhancement.md) — Nodes CFR Remote Fleet Clone Enhancement & Async Delegation
- [75-gitmap-u1-ubuntu-agm-fleet-integration.md](pending/75-gitmap-u1-ubuntu-agm-fleet-integration.md) — GitMap U1 Ubuntu Fleet Integration, AGM Migration & Cross-OS Automation

---

## 4. Active Subtasks Directories

Only active subtask directories are maintained under `.ai-memory/plans/subtasks/` (12 directories):

- `subtasks/236-deep-spec-consolidation-and-canonical-reduction/` (Active current parent plan)
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
