# Subtask 02: Deep Survey & Aging Specs Inventory Specification

- **Parent Plan:** [236-deep-spec-consolidation-and-canonical-reduction.md](../../236-deep-spec-consolidation-and-canonical-reduction.md)
- **Subtask ID:** `02`
- **Status:** `IN_PROGRESS`
- **Assigned Worker:** Research 01 & Research 02 (Synthesis: Author 01)
- **Target Scope:** `02-spec/21-app/`, `.ai-memory/plans/completed/`, `.ai-memory/plans/pending/`, `.ai-memory/memory/`

---

## 1. Context & Objective

Per user instruction, perform a comprehensive deep survey across application specifications, completed plans, and AI memory to inventory aging, obsolete, and redundant files. This subtask documents the exact pre-compaction inventory, analyzes draft vs. final implementations (X, Y vs. A, B), audits the Section 3 raw text bloat across historical plans, and validates the strict isolation of all pending plans, active plans, and active subtasks.

---

## 2. Global Pre-Compaction Inventory Matrix

| Target Subsystem / Directory | Pre-Compaction Inventory | Primary Source of Bloat / Aging Data | Post-Compaction Target | Projected Reduction |
| :--- | :---: | :--- | :---: | :---: |
| **`02-spec/21-app/` Standalone Directories** | **72 directories** (309 internal files) | Fragmented across 3 historical waves (200-series, mid-tier, early legacy) | **0 standalone directories** | **-100.0%** |
| **`02-spec/21-app/` Standalone Root Files** | **181 loose files** | 178 loose `.md` files, `.gitkeep`, ERD diagram, legacy migration plans | **1 master `readme.md`** | **-99.4%** |
| **`02-spec/21-app/` Canonical Clusters** | **8 clusters** | Retain ratified domain specifications with enriched component specs | **8 clusters** | Authoritative Core |
| **`.ai-memory/plans/completed/` Milestones** | **43 files** (11,665 lines) | Milestones 01–10 hold 9,448 lines (~81% bloat) due to verbatim task dumps | **18 files** (~2,200 lines) | **-58.1% files / -81.1% lines** |
| **`.ai-memory/memory/` Knowledge Base** | **71 files** across 12 dirs | 35 fragmented micro-snippets + 36 dead-weight/stale legacy docs | **10 files** (8 clusters + 1 gov + readme) | **-85.9%** |
| **`.ai-memory/plans/pending/`** | **4 files** (273 lines) | Unfinished future tasks | **4 files** (Untouched) | **0% (Protected)** |
| **`.ai-memory/plans/` Active Plans** | **13 active plans** | Work-in-progress feature roadmaps | **13 active plans** (Untouched) | **0% (Protected)** |
| **`.ai-memory/plans/subtasks/` Active Dirs** | **12 directories** (51 files) | Subtasks belonging to open or pending plans | **12 directories** (Untouched) | **0% (Protected)** |

---

## 3. Deep Survey of `02-spec/21-app/`

### 3.1 Surviving Standalone Directories Breakdown (72 Folders)

The 72 surviving standalone directories represent three historical waves of feature implementation:

#### A. 200-Series Directories (34 Folders)
These directories record recent milestones (Milestones 36–42). They contain high-value ratified technical contracts (A, B) alongside intermediate experimental drafts:
- `202-gitmap-pull-errors-splitdb` (3 files)
- `204-ssh-password-interception-and-rsa-credential-vault` (2 files)
- `205-gitmap-u1-ubuntu-agm-fleet-integration` (1 file)
- `206-windows-to-ubuntu-fleet-migration-and-secrets-vault` (1 file)
- `207-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca` (1 file)
- `208-detectedproject-foreign-key-and-ssh-lifecycle-hardening` (1 file)
- `209-ai-agent-task-orchestrator-and-split-db` (2 files)
- `210-gitmap-push-fix-command-and-auth-recovery` (2 files)
- `211-ubuntu-fleet-git-clone-and-os-customization` (3 files)
- `212-ubuntu-fleet-full-customization-and-embedded-runner` (2 files)
- `213-antigravity-ubuntu-update-and-macro-automation` (3 files)
- `214-ubuntu-fleet-automation-and-workstation-governance` (2 files)
- `215-ubuntu-fleet-cleanup-antigravity-projects-and-app-manager` (2 files)
- `216-gitmap-prompting-freeze-and-suggestion-engine-fix` (1 file)
- `217-antigravity-fleet-parity-theme-preset-plugins-and-delegation` (2 files)
- `218-nodes-agy-ui-remote-settings-and-cursor-automation` (2 files)
- `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging` (1 file)
- `220-pipeline-pe-unit-test-traceback-and-heatmap` (3 files)
- `221-ci-cd-fix-nested-if-and-test-summary-remediation` (3 files)
- `222-ubuntu-pull-tree-omz-force-error-diagnostics-and-release` (2 files)
- `223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search` (3 files)
- `224-ubuntu-cursor-dock-pin-projects-os-update-and-remote-uninstall` (2 files)
- `225-ubuntu-cursor-start-menu-dock-position-and-profile-migration` (2 files)
- `226-ubuntu-cursor-memories-conversations-and-projects-migration` (3 files)
- `227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions` (3 files)
- `228-nodes-deploy-agm-accounts-and-fleet-sync` (1 file)
- `229-antigravity-ide-projects-and-settings-backup` (2 files)
- `229-nodes-deploy-repos-and-multi-ide-fleet-sync` (2 files)
- `230-antigravity-backup-e2e-and-release` (2 files)
- `230-token-purge-installer-workdir-pull-agm-and-ui-modernization` (4 files)
- `231-antigravity-ide-projects-and-repo-secrets-restore` (2 files)
- `232-pending-commits-sends-and-nodes-commit-suite` (2 files)
- `233-ubuntu-ide-and-github-desktop-scan-sync` (4 files)
- `234-codebase-review-remediation-and-consolidation` (4 files)
- `235-app-spec-and-completed-plans-consolidation-and-reduction` (3 files)

#### B. Mid-Tier Directories (10 Folders)
- `65-gitmap-update-all-zip-and-fixes` (2 files)
- `66-fix-gitmap-pa-duplicate-repos` (2 files)
- `67-nodes-cfr-remote-fleet-clone-enhancement` (2 files)
- `69-nodes-cfr-ui-async-delegation` (2 files)
- `71-pull-all-auto-ff-merge-cpu-scaling-and-pe-limit-fix` (3 files)
- `72-chrome-ext-test-and-vmware-spec` (2 files)
- `72-repo-dedup-os-aware-equalfold` (2 files)
- `145-ssh-fleet-deploy-remote-clone-api-ui` (4 files)
- `156-which-os-cross-platform-shell-and-node-profiling` (4 files)
- `181-gitmap-ignore-and-cache-engine` (3 files)

#### C. Early Legacy Directories (28 Folders)
- `01-vscode-project-manager-sync` (3 files)
- `03-commit-in` (10 files)
- `03-general` (20 files)
- `03-tasks` (3 files)
- `04-generic-cli` (37 files)
- `04-json-contract` (1 file)
- `07-error-and-logging` (3 files)
- `07-generic-release` (11 files)
- `08-generic-update` (11 files)
- `08-json-schemas` (29 files)
- `09-pipeline` (13 files)
- `09-pipeline-extend-v2` (6 files)
- `09-pipeline-historical-eta` (4 files)
- `10-pipeline-and-repo-split-db` (4 files)
- `11-pipeline-errorlogs-timeline-and-fix` (3 files)
- `12-repo-creation-profiles-and-cloud-backup` (4 files)
- `15-distribution-and-runner` (5 files)
- `16-scan-clone-cluster-interactive-macros` (9 files)
- `17-mv-rm-resolver-replace` (8 files)
- `18-install-cg-ssh` (1 file)
- `19-ssh-executor` (1 file)
- `23-app-db` (0 files — empty folder)
- `24-app-ui-design-system` (1 file)
- `25-app-spec-audit` (2 files)
- `26-coding-guideline-audit` (2 files)
- `37-project-detection` (13 files)
- `samples` (1 file)

### 3.2 Standalone Root Files Breakdown (181 Items)
- **178 Loose Markdown Files:** Dating from early CLI commands (`02-cli-interface.md`) up to specialized utilities (`200-semantic-flat-commit-and-auto-stage-command.md`).
- **1 Database Diagram:** `gitmap-database-erd.mmd` (Authoritative Mermaid ERD to be folded into `06-database-and-split-db/`).
- **1 Empty Git Placeholder:** `.gitkeep` (obsolete root marker).
- **1 Migration Plan:** `xx-migration-plan.md` (superseded by modern Split-DB engine).

---

## 4. Deep Survey of `.ai-memory/plans/completed/`

### 4.1 Historical Milestones Line Count & Bloat Analysis
The completed plans directory contains 43 milestone files totaling **11,665 lines**:
- **Milestones 01–10 (Early Pass):** **9,448 lines** (~81.0% of total lines across all 43 files).
  - Average lines per file: **945 lines**.
  - Cause of bloat: Section 3 contains verbatim raw markdown dumps of 18+ individual subtask files rather than concise structured outcome ledgers.
  - Top bloated files:
    - `10-installers-multios-setup-and-web-stacks.md`: 2,849 lines (Section 3 raw dumps: 2,796 lines)
    - `04-cicd-pipelines-runners-and-streaming-telemetry.md`: 1,455 lines (Section 3 raw dumps: 1,418 lines)
    - `06-git-operations-commit-engines-and-remediation.md`: 1,273 lines (Section 3 raw dumps: 1,236 lines)
    - `08-terminal-ui-help-parity-and-cli-commands.md`: 982 lines (Section 3 raw dumps: 945 lines)
    - `05-database-engine-sqlite-joins-and-scanners.md`: 785 lines (Section 3 raw dumps: 748 lines)
    - `01-coding-guidelines-and-style-audits.md`: 746 lines (Section 3 raw dumps: 709 lines)
    - `03-type-safety-function-signatures-and-contracts.md`: 496 lines (Section 3 raw dumps: 459 lines)
    - `09-chrome-profile-management-picker-and-token-vault.md`: 477 lines (Section 3 raw dumps: 440 lines)
    - `02-error-management-and-cliexit-architecture.md`: 276 lines (Section 3 raw dumps: 239 lines)
- **Milestones 11–27 (Second Pass):** **1,000 lines** (Average: 58 lines/file). Highly structured but heavily duplicate domain topics covered in Milestones 01–10.
- **Milestones 28–43 (Modern Pass):** **1,217 lines** (Average: 76 lines/file). Crisp, authoritative records of recent development.

### 4.2 Generational Compaction Architecture
Historical Milestones 01–27 are consolidated into **2 dense Generational Milestones**:
1. **`01-generation-1-core-foundation-milestones.md`** (Synthesizing Historical Milestones 01–15):
   - Merges core Go architecture, scanner, project classification, error handling (`*appfault.AppError`), monadic results (`Result[T]`), Split-DB foundation, and flat commit suite into a dense execution ledger.
2. **`02-generation-2-fleet-automation-and-modularization.md`** (Synthesizing Historical Milestones 16–27):
   - Merges cluster orchestration, macro streaming, nuclear package modularization, and smart test inventory into a dense execution ledger.
3. **Renumbering Surviving Milestones 28–43:**
   - Surviving Milestones 28–43 are cleanly renumbered as Milestones `03` through `18`.
   - Result: 43 files reduced to **18 authoritative milestone files** (-58.1% file reduction, -81.1% line reduction).

---

## 5. Deep Survey of `.ai-memory/memory/`

### 5.1 Current Inventory Analysis (71 Files)
The `.ai-memory/memory/` folder contains 71 files across 12 subdirectories + root:
- **Active Structured Categories (35 files, 535 lines):**
  - `avoid/` (5 files): anti-patterns for tests, macros, git, booleans, paths.
  - `features/` (5 files): CLI, Git, SSH, DB, Antigravity catalogs.
  - `issues/` (10 files): RCA matrices across major subsystems.
  - `learned/` (15 files): short 12–20 line notes on patterns and decisions.
- **Stale, Dead-Weight, and Duplicate Files (36 files, ~1,850 lines):**
  - `01-index.md` in root memory: exact byte-for-byte duplicate of `readme.md`.
  - `project/what-to-read.md`: ancient April 2026 onboarding map referencing deleted directories.
  - `tech/` (7 files), `workflow/` (4 files), `constraints/` (4 files), `style/` (3 files), `specs/` (2 files), `suggestions/` (3 files): legacy v2.x/v3.x session files.

### 5.2 Compaction into 8 Canonical Domain References
All memory is unified into **8 Canonical Domain Memory References** matching the 8 Canonical Spec Clusters, plus 1 Project Governance reference:
1. `00-project-governance-and-invariants.md`
2. `01-cli-architecture.md`
3. `02-scanner-and-projects.md`
4. `03-git-operations-and-pull.md`
5. `04-fleet-nodes-and-ssh.md`
6. `05-antigravity-and-ide.md`
7. `06-database-and-split-db.md`
8. `07-pipeline-and-diagnostics.md`
9. `08-distribution-and-release.md`
10. `readme.md` (Master catalog linking only to the 9 references)

Result: 71 files reduced to **10 authoritative memory files** (-85.9% reduction).

---

## 6. Strict Pending Isolation Verification

The survey verified that active tasks, pending roadmaps, and ongoing development pipelines are completely quarantined from modification or deletion:

### 6.1 Pending Plans in `.ai-memory/plans/pending/` (4 Files)
1. `01-ports-and-ssh-enablement.md` — Ports inspection & cross-platform OpenSSH daemon governance.
2. `56-vmware-hardware-batch-and-macro-orchestration.md` — VMware multi-VM batch ops & macro execution.
3. `67-nodes-cfr-remote-fleet-clone-enhancement.md` — Nodes CFR remote fleet clone and async probing.
4. `75-gitmap-u1-ubuntu-agm-fleet-integration.md` — GitMap U1 Ubuntu fleet integration & AGM migration.

### 6.2 Active Open Plans in `.ai-memory/plans/` (13 Feature Plans)
- `82-ubuntu-fleet-full-customization-and-embedded-runner.md`
- `214-ubuntu-fleet-automation-and-workstation-governance.md`
- `216-gitmap-prompting-freeze-and-suggestion-engine-fix.md`
- `217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md`
- `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md`
- `223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md`
- `225-ubuntu-cursor-start-menu-dock-position-and-profile-migration.md`
- `226-ubuntu-cursor-memories-conversations-and-projects-migration.md`
- `227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md`
- `229-nodes-deploy-repos-and-multi-ide-fleet-sync.md`
- `231-antigravity-ide-projects-and-repo-secrets-restore.md`
- `235-app-spec-and-completed-plans-consolidation-and-reduction.md`
- `236-deep-spec-consolidation-and-canonical-reduction.md`

### 6.3 Active Subtask Directories in `.ai-memory/plans/subtasks/` (12 Folders)
All 12 subtask directories belong to active or pending parent plans and are strictly quarantined.

### 6.4 Verification Properties
- `isPendingIsolated: true`
- `isActivePlansProtected: true`
- `isActiveSubtasksProtected: true`
- `isZeroMutationEnforced: true`

---

## 7. Next Steps & Worker Handoff

This inventory serves as the authoritative blueprint for:
1. **Subtask 03:** Operational implementation plan for folding the 72 legacy folders and 181 loose files in `02-spec/21-app/` into the 8 Canonical Domain Clusters.
2. **Subtask 04:** Operational compaction of completed plans and memory into generational milestones and canonical memory references.
3. **Subtask 05:** Registry synchronization and quality gate verification.
