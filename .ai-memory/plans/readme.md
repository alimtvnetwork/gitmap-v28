# Master Plans Registry (`.ai-memory/plans`)

This registry catalogs the active, completed, and pending execution plans for GitMap. In accordance with Plan 236 ([Spec 236](../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/01-architecture-spec.md)), historical completed plans have been consolidated into **18 Authoritative Milestones** under `.ai-memory/plans/completed/`.

---

## 1. Active Plans

The following plans represent open and actively maintained development streams:

- [gitmap-and-antigravity-manager-agm-package-enhancements.md](gitmap-and-antigravity-manager-agm-package-enhancements.md) — GitMap Single-Repo Pull UX Simplification, Progress Bar Column Width Wrap Fix, Shell Wrapping RCA & AGM Subsystem Alignment (Spec: [gitmap-and-antigravity-manager-agm-package-enhancements](../../02-spec/21-app/gitmap-and-antigravity-manager-agm-package-enhancements/01-architecture-spec.md)) (Status: `in_progress`)

- [heatmap-commands-testing-and-rc-admin-instructions.md](completed/heatmap-commands-testing-and-rc-admin-instructions.md) — Heatmap Commands Testing & RC Admin Instructions: Table inspection (`tables`), record viewing with limits & JSON/YAML formatting (`view`), direct read-only SQL query engine (`query`), and on-machine E2E test suite (Spec: [heatmap-commands-testing-and-rc-admin-instructions](../../02-spec/21-app/heatmap-commands-testing-and-rc-admin-instructions/01-architecture-spec.md)) (Status: `completed`)

- [heatmap-command-and-database-structure-instructions.md](completed/heatmap-command-and-database-structure-instructions.md) — Heatmap Command and Database Structure Instructions: Repo-cache (`rc`) tree folder structure (`.heatmap/01-<ai-name>/`), Split-DB architecture (`root-heatmap.db` and per-repo `heatmap.db`), local CLI tracking database, autonomous conflict resolution via temporary SQLite delta replay, and Muse lowercase guideline enforcement (Spec: [heatmap-command-and-database-structure-instructions](../../02-spec/21-app/heatmap-command-and-database-structure-instructions/01-architecture-spec.md)) (Status: `completed`)

- [muse-memory-io-and-gitmap-integration.md](muse-memory-io-and-gitmap-integration.md) — Muse Memory IO & GitMap Integration: 7-day digest into repo-cache `12-muse-ai/` (README.md + sequenced readers + chat seeds + 00-new-muse-onboarding.md), "say READ" re-sync convention, Muse prompt/skill repo-cache updates, go-email-reader email docs (config from repo-secrets `11-go-email-reader/`), `gitmap llm train muse [--days N]` command (repo-cache vault, gitmap-tracked last-read watermark, DECISION TRACE), new-command docs + changelog (Spec: [277](../../02-spec/21-app/277-muse-memory-io-and-gitmap-integration/readme.md)) (Status: `active`)

- [update-antigravity-ide.md](update-antigravity-ide.md) — Antigravity Desktop IDE Update Integration: GitMap CLI `update antigravity` / `update ide` commands, scripts fixture (`scripts-fixer`) `update` verb, Antigravity tools updater script, and live Ubuntu system updates (Spec: [update-antigravity-ide](../../02-spec/21-app/update-antigravity-ide/readme.md)) (Status: `active`)
- [agent-ai-synchronization-and-task-management.md](completed/agent-ai-synchronization-and-task-management.md) — Agent AI Synchronization & Task Management: MachineId Tracking, ActionType Enum Lifecycle, Real-Time Fast Heatmap Consensus & Multi-Agent Telemetry (Spec: [agent-ai-synchronization-and-task-management](../../02-spec/21-app/agent-ai-synchronization-and-task-management/01-architecture-spec.md)) (Status: `completed`)
- [git-prime-machine-os-version-and-ui-help-docs-improvement.md](git-prime-machine-os-version-and-ui-help-docs-improvement.md) — Git Prime Machine OS Version and UI Help Docs Improvement: Structured OS Version Object & JSON Breakdown, HD Docs Web UI Command Expansion, and Merge-AI Parameterized Examples (Spec: [git-prime-machine-os-version-and-ui-help-docs-improvement](../../02-spec/21-app/git-prime-machine-os-version-and-ui-help-docs-improvement/01-architecture-spec.md)) (Status: `active`)
- [249-ui-accessibility.md](249-ui-accessibility.md) — Dashboard keyboard accessibility + compact professional polish: full keyboard operability, shortcuts + `?` overlay, ARIA tablist/tab/tabpanel, compact density on all 10 tabs (Spec: [249](../../02-spec/21-app/249-ui-accessibility/01-overview.md))
- [gitmap-repo-cache-scan-clone-and-ls-ecosystem.md](gitmap-repo-cache-scan-clone-and-ls-ecosystem.md) — GitMap Repo-Cache Scan (`--rc`, `--sep`), Clone/CFR/CFRP (`--rc`, Split-DB Cache, TreeView Selector), and `ls rc` Format Inspector (Spec: [gitmap-repo-cache-scan-clone-and-ls-ecosystem](../../02-spec/21-app/gitmap-repo-cache-scan-clone-and-ls-ecosystem/01-architecture-spec.md))
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
- [24-agm-update-tui-fix-and-release.md](completed/247-agm-update-tui-fix-and-release.md) — `gitmap agm update` TUI fix (coordinated pterm UI, captured output) + `agm` alias misrouting fix + failure-path trace removal; U01–U06 all pass; released as 6.517.0 (Spec: [247](../../02-spec/21-app/247-agm-update-tui-fix-and-release/01-rca-scope-and-plan.md))
- [25-agent-collision-and-dispatch-registry-plan.md](250-gitmap-agent-collision-and-dispatch-registry-plan.md) — Dispatch registry + agent collision feature (task-DB file claims, look-ahead, scoped staging) + A08 swallowed-error fix + version truth unification; release 6.518.0 (Spec: [250](../../02-spec/21-app/250-gitmap-agent-collision-and-dispatch-registry-plan/01-design.md)) (Status: `awaiting-review`)
- [26-table-improvement-and-command-enhancement.md](completed/81-table-improvement-and-command-enhancement.md) — Pending commits table overhaul (UNCOMMITTED column, tree remediation, fleet batch fix footer), dual SQLite status cache and backup DBs, new-commands discovery engine, helpdoc and LLM skills sync (Spec: [01-architecture-spec.md](../../02-spec/21-app/table-improvement-and-command-enhancement/01-architecture-spec.md))
- [27-gitmap-command-implementation.md](completed/83-gitmap-command-implementation.md) — GitMap Tree Export Flags (`--file`/`-f`), Positional Glob Crash Remediation, Tree Search Command Family & Tree Learn Split SQLite DB Engine (Spec: [01-architecture-spec.md](../../02-spec/21-app/gitmap-command-implementation/01-architecture-spec.md))
- [28-repository-fix-and-solution-implementation.md](completed/84-repository-fix-and-solution-implementation.md) — Pull Failure Subtree Yellow Solutions, Non-Repo Folder Smart Clone Resolution & Unified Batch Fix (Spec: [01-architecture-spec.md](../../02-spec/21-app/repository-fix-and-solution-implementation/01-architecture-spec.md))
- [29-gitmap-command-enhancements.md](completed/gitmap-command-enhancements.md) — OS JSON Output, Hardware Telemetry (CPU, RAM, Disk) with 5-Second Cache, Network Adapter Metrics, and Fleet Node Task Queue Delegation (Spec: [01-architecture-spec.md](../../02-spec/21-app/gitmap-command-enhancements/01-architecture-spec.md))
- [29-repository-creation-and-folder-structure-instructions.md](completed/repository-creation-and-folder-structure-instructions.md) — Multi-Repo Non-Destructive Amalgamation (Icon Studio 2026 V1), Zero Replacement Invariant, Commit Metadata Tracking, AI Folder Structure Consolidation Guide & Multi-IDE Ecosystem Registration (Spec: [01-architecture-spec.md](../../02-spec/21-app/repository-creation-and-folder-structure-instructions/01-architecture-spec.md))
- [30-gitmap-repository-management-and-release-process.md](completed/gitmap-repository-management-and-release-process.md) — Repo-Cache Split-DB JSON Injection, Two-Tier Subtree Solutions Formatting, Native GitMap Remediation Commands, Antivirus False-Positive Mitigation & Go Direct Path Query CLI Acceleration (Spec: [01-architecture-spec.md](../../02-spec/21-app/gitmap-repository-management-and-release-process/01-architecture-spec.md))
- [31-github-desktop-sync-fix.md](completed/02-github-desktop-sync-fix.md) — GitHub Desktop Linux Registration Crash Remediation, 250ms IPC Pacing, Canonical 'open' Subcommand Dispatch, Candidate Resolution & Non-Swallowed Error Surfacing (Spec: [01-architecture-spec.md](../../02-spec/21-app/02-github-desktop-sync-fix/01-architecture-spec.md))
- [32-gitmap-enhancement-and-optimization.md](completed/gitmap-enhancement-and-optimization.md) — Resilient Autostash Pull, Interactive Post-Pull Remediation Prompting, ANSI Progress Bar Color Fix, 48h Commit Recency for PAE, Multi-Node Resource Monitoring & Adaptive Concurrency Boost, Semantic Sends CP All Grouping & Fleet SSH Parity (Spec: [01-architecture-spec.md](../../02-spec/21-app/gitmap-enhancement-and-optimization/01-architecture-spec.md))
- [33-bug-fix-and-terminal-ui-improvement.md](completed/bug-fix-and-terminal-ui-improvement.md) — Pure Binary Updater, Source-Build Inversion, Truecolor ANSI Lipgloss Terminal UI Installers & Standalone Local Installers (Spec: [readme.md](../../02-spec/21-app/bug-fix-and-terminal-ui-improvement/readme.md))
- [34-gitmap-pull-fix-remediation-optimization.md](completed/278-gitmap-pull-fix-remediation-optimization.md) — Windows GCM Credential Store Repair (`wincredman`), Single-Repo Reclone URL Handler, Generic Typed Parallel Fix-All, Dynamic 20-Core Pull Concurrency, Redesigned Post-Pull Menu, 24h Update Probe Throttle & Repo-Cache Multi-AI Memory Store (Spec: [278](../../02-spec/21-app/278-gitmap-pull-fix-remediation-optimization/readme.md))
- [35-automerge-conflict-resolution-and-ide-guidance.md](completed/279-automerge-conflict-resolution-and-ide-guidance.md) — Automerge Histogram Diff Optimization, Pre-Flight Divergence Checking, Git Rerere Memory Engine, Multi-IDE Launcher (`code`, `cursor`, `windsurf`, `codium`), Step-by-Step 3-Way Merge Editor Guide, `gitmap resolve` CLI, and Lipgloss Post-Pull Conflict UI (Spec: [279](../../02-spec/21-app/279-automerge-conflict-resolution-and-ide-guidance/readme.md))
- [36-large-files-split-flatten-and-muse-sync.md](completed/280-large-files-split-flatten-and-muse-sync.md) — Large Files Splitting Engine (`gitmap split large files` $> 10\text{ MB}$, multipart zip, STORE/DEFLATE heuristics, Git history purge), Git LFS Conversion (`gitmap lfs convert` $> 5\text{ MB}$), Repository Flattening (`gitmap repo flatten` squashed flat commit, archive, undo), and Muse LLM Training Source Prioritization (Spec: [280](../../02-spec/21-app/280-large-files-split-flatten-and-muse-sync/readme.md))
- [37-pull-summary-remediation-help-and-alim-seo-fix.md](completed/281-pull-summary-remediation-help-and-alim-seo-fix.md) — Pull Summary Transparency (Real-Time Worker Visibility & Up-to-Date Repositories Grouping), Remediation Menu Interactive Explanatory Guide ([h/help]), Safe-Pull Conflict Detection & Keyword Parsing Hardening, and Alim SEO Writing LFS Pointer Normalization (Spec: [281](../../02-spec/21-app/281-pull-summary-remediation-help-and-alim-seo-fix/readme.md))



---

## 2b. Recently Completed (awaiting milestone consolidation)

- [271-improvement-program.md](271-improvement-program.md) — Program 271: workflow bug fixes, dead code 36→0, deprecation hook, v6.551.0; WS2 nc auto-discovery reverted as broken (Status: `completed` 2026-10-15)
- [272-decision-record-improvement-scope.md](272-decision-record-improvement-scope.md) — Owner decision: items 1–3 approved, 4–6 explicitly NOT approved (Status: `completed` 2026-10-15)
- [272-deprecation-coverage-nc-codegen.md](272-deprecation-coverage-nc-codegen.md) — Deprecation policy removed, coverage floor gate, nc codegen, v6.553.0 (Status: `completed` 2026-10-15)
- [273-gitmap-help-text-and-ui-enhancement.md](273-gitmap-help-text-and-ui-enhancement.md) — 4 design-system themes, `help example`, `mr` command, 20 proposals, virus RCA (Status: `completed` 2026-10-16)
- [274-cicd-verification-and-enhancement.md](274-cicd-verification-and-enhancement.md) — CI 44/44 green, v6.558.0, 10 proposals + 20 improvements, 2 E2E bugs documented (Status: `completed` 2026-10-17)
- [275-implement-30-new-commands.md](275-implement-30-new-commands.md) — All 30 proposals implemented, E2E 30/30, v6.559.0 (Status: `completed` 2026-10-17)
- [282-help-doc-ui-hd-theme-and-terminal-runner.md](282-help-doc-ui-hd-theme-and-terminal-runner.md) — Help Doc UI HD Theme & Live Terminal Runner: Registered help text across terminal and UI, added prompt finder and categories, updated Help Dashboard themes, built interactive TerminalDrawer.tsx, and promoted terminal in cli/cmdui/ (Status: `completed` 2026-10-17)
- [283-cli-update-ui-color-alignment-and-release-notes.md](283-cli-update-ui-color-alignment-and-release-notes.md) — CLI Update UI Color Alignment, Spacing, Padding & Release Notes: Clean white-to-yellow version transition, 4-space indented SUCCESS badges matching step indicators, top/bottom padding, duplicate output prevention, and 1–2 line release notes snippet (Status: `completed` 2026-10-17)
- [284-agent-heatmap-db-introspection-and-subtask-code-completion.md](284-agent-heatmap-db-introspection-and-subtask-code-completion.md) — Agent Heatmap Split-DB Architecture, Native DB Introspection CLI (`gitmap db ...`), Subtask `--code` Resolution & Clean Error Handling (Status: `completed` 2026-10-18)
- [285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix.md](285-gitmap-pae-infinite-hang-and-heartbeat-rca-fix.md) — GitMap PAE Post-Pull Infinite Hang & Heartbeat Teardown RCA Fix: Immediate teardown of bar/heartbeat, 3s bounded ignore collect timeout, heartbeat self-termination guard (Status: `completed` 2026-10-18)

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
