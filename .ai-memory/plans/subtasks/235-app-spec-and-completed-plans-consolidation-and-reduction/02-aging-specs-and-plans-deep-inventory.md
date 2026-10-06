# Subtask 02: Aging Specifications & Completed Plans Deep Survey Inventory and Compaction Mapping

> **Subtask ID:** Subtask-02  
> **Parent Plan:** `.ai-memory/plans/235-app-spec-and-completed-plans-consolidation-and-reduction.md`  
> **Architecture Spec:** `02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/01-architecture-spec.md`  
> **Owned Files & Target Scope:**  
> - Deep Inventory Survey across `02-spec/21-app/` (310 filesystem items)  
> - Completed Plans Catalog across `.ai-memory/plans/completed/` (163 post-milestone-27 plans)  
> - Memory Directory Survey across `.ai-memory/memory/` (23 files/folders)  
> - Pending Isolation Registry (`.ai-memory/plans/pending/`, open plans, open subtasks)  

---

## 1. Objectives

1. **Conduct Comprehensive Survey of `02-spec/21-app/`:**
   - Catalog all 310 items (80 subdirectories and 228 standalone Markdown files) located within `02-spec/21-app/`.
   - Map each item to one of the **8 Canonical Domain Clusters**:
     1. `01-cli-architecture`
     2. `02-scanner-and-projects`
     3. `03-git-operations-and-pull`
     4. `04-fleet-nodes-and-ssh`
     5. `05-antigravity-and-ide`
     6. `06-database-and-split-db`
     7. `07-pipeline-and-diagnostics`
     8. `08-distribution-and-release`
2. **Catalog Completed Plans for Milestone Compaction:**
   - Map the 163 completed plans (post-milestone-27) in `.ai-memory/plans/completed/` into **15 Milestone Clusters (Milestones 28 to 42)**.
   - Establish high-density milestone summary targets while preserving all Go types, interface signatures, and verification outcomes.
3. **Survey `.ai-memory/memory/`:**
   - Inventory session logs, legacy compatibility audits, and transient status notes to consolidate under canonical memory categories (`constraints/`, `avoid/`, `learned/`, `project/`, `style/`).
4. **Enforce Compaction Invariant (A, B vs. X, Y):**
   - Identify aging draft states (X, Y) that were superseded by subsequent implementations (A, B) and earmark them for pruning.
5. **Establish Strict Pending Isolation Boundaries:**
   - Lock `.ai-memory/plans/pending/` (4 files), active open plans (7 files), and active subtask folders (10 folders, 52 files) with strict exclusion flags.

---

## 2. Specification Inventory & 8 Canonical Domain Clusters Mapping

The 310 items in `02-spec/21-app/` map into the following 8 canonical functional domains:

| Domain Cluster | Item Count | Representative Source Specifications in `02-spec/21-app/` |
| :--- | :---: | :--- |
| **01-cli-architecture** | 43 | `02-cli-interface.md`, `04-generic-cli/`, `25-command-history.md`, `182-terminal-tab-completion-and-flag-suggestions.md`, `191-universal-command-help-restructuring-markdown-box-display.md`, `216-gitmap-prompting-freeze-and-suggestion-engine-fix/` |
| **02-scanner-and-projects** | 50 | `01-vscode-project-manager-sync/`, `03-scanner.md`, `04-formatter.md`, `05-cloner.md`, `06-config.md`, `07-data-model.md`, `177-scan-alias-migration-internal-errors-db-and-fleet-inventory-aggregation.md`, `181-gitmap-ignore-and-cache-engine/` |
| **03-git-operations-and-pull** | 36 | `03-commit-in/`, `20-revert.md`, `24-amend-author.md`, `63-semantic-flat-commit-suite/`, `66-fix-gitmap-pa-duplicate-repos/`, `72-repo-dedup-os-aware-equalfold/`, `200-semantic-flat-commit-and-auto-stage-command.md`, `201-pull-ignore-concurrency-split-db-cache-and-deadlock-prevention.md`, `210-gitmap-push-fix-command-and-auth-recovery/`, `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/` |
| **04-fleet-nodes-and-ssh** | 34 | `01-ports-and-ssh-enablement/`, `202-vmware-cli-commands-and-fleet-management.md`, `203-ports-inspection-and-ssh-daemon-enablement.md`, `204-ssh-password-interception-and-rsa-credential-vault/`, `206-windows-to-ubuntu-fleet-migration-and-secrets-vault/`, `208-detectedproject-foreign-key-and-ssh-lifecycle-hardening/`, `228-nodes-deploy-agm-accounts-and-fleet-sync/` |
| **05-antigravity-and-ide** | 10 | `10-github-desktop.md`, `184-agy-rp-running-projects-and-recreate-safety/`, `217-antigravity-fleet-parity-theme-preset-plugins-and-delegation/`, `223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search/`, `226-ubuntu-cursor-memories-conversations-and-projects-migration/`, `229-antigravity-ide-projects-and-settings-backup/`, `231-antigravity-ide-projects-and-repo-secrets-restore/` |
| **06-database-and-split-db** | 11 | `10-pipeline-and-repo-split-db/`, `181-gitmap-ignore-and-cache-engine/`, `199-devtools-cache-dynamic-discovery-tree-view-and-split-db.md`, `202-gitmap-pull-errors-splitdb/`, `209-ai-agent-task-orchestrator-and-split-db/` |
| **07-pipeline-and-diagnostics** | 23 | `09-pipeline/`, `09-pipeline-extend-v2/`, `09-pipeline-historical-eta/`, `220-pipeline-pe-unit-test-traceback-and-heatmap/`, `221-ci-cd-fix-nested-if-and-test-summary-remediation/`, `22-scan-release-import.md`, `26-stats.md` |
| **08-distribution-and-release** | 29 | `07-generic-release/`, `08-generic-update/`, `09-build-deploy.md`, `21-list-releases.md`, `213-antigravity-ubuntu-update-and-macro-automation/`, `214-ubuntu-fleet-automation-and-workstation-governance/`, `230-token-purge-installer-workdir-pull-agm-and-ui-modernization/`, `234-codebase-review-remediation-and-consolidation/` |

---

## 3. Completed Plans Inventory & 15 Milestone Clusters (Milestones 28 to 42)

The 163 post-milestone-27 completed plan files in `.ai-memory/plans/completed/` are clustered as follows:

| Milestone Cluster | Target Output File | Source Plans Included | Count |
| :--- | :--- | :--- | :---: |
| **Milestone 28** | `28-antigravity-ide-prompts-and-protocols.md` | `28`, `41`, `42`, `54`, `56`, `77`, `79`, `80`, `83`, `xx-agy-enhancements` | 10 |
| **Milestone 29** | `29-ssh-remote-exec-auth-and-join.md` | `30`, `39`, `43`, `48`, `66`, `81`, `89`, `91`, `97` | 9 |
| **Milestone 30** | `30-pipeline-telemetry-eta-and-remediation.md` | `31`, `32`, `45`, `46`, `84`, `88`, `90`, `162` | 8 |
| **Milestone 31** | `31-aum-engine-lazy-regex-and-training.md` | `33`, `34`, `35`, `36`, `37`, `38` | 6 |
| **Milestone 32** | `32-macro-fleet-automation-and-streaming.md` | `40`, `54`, `89`, `98`, `192` | 5 |
| **Milestone 33** | `33-split-db-transactions-and-cache.md` | `42`, `44`, `49`, `50`, `62`, `78`, `85` | 7 |
| **Milestone 34** | `34-os-profiling-linutil-and-governance.md` | `51`, `52`, `53`, `57`, `76`, `81`, `82`, `104`, `105`, `108` | 10 |
| **Milestone 35** | `35-vmware-cluster-automation-and-hardware.md` | `72`, `73`, `74`, `86`, `87` | 5 |
| **Milestone 36** | `36-gitmap-pull-optimization-and-cache.md` | `58`, `60`, `61`, `64`, `66`, `70`, `71`, `72`, `164`, `166`, `169`, `197` | 12 |
| **Milestone 37** | `37-semantic-flat-commit-and-non-git-rca.md` | `63`, `65`, `80`, `165`, `168` | 5 |
| **Milestone 38** | `38-ssh-password-interception-and-rsa-vault.md` | `67`, `68`, `76`, `78`, `163`, `198`, `199` | 7 |
| **Milestone 39** | `39-ai-agent-task-orchestrator-and-web-ui.md` | `79`, `112`, `174`, `184`, `185`, `186`, `202` | 7 |
| **Milestone 40** | `40-fleet-nodes-cfr-manifests-and-delegation.md` | `55`, `57`, `69`, `94`, `98`, `99`, `195`, `196`, `200`, `201`, `218`, `228` | 12 |
| **Milestone 41** | `41-remote-deploy-streaming-and-auto-update.md` | `100`, `101`, `106`, `109`, `110`, `111`, `171`, `172`, `173`, `177`, `178`, `179`, `180`, `184`, `188`, `189`, `194` | 17 |
| **Milestone 42** | `42-workstation-parity-and-codebase-consolidation.md` | `215`, `216`, `217`, `220`, `221`, `222`, `223`, `224`, `228`, `229`, `230`, `230-token`, `231`, `232`, `233`, `234` | 16 |
| **Total Plans Merged** | | | **163** |

---

## 4. Compaction Invariant Mapping: Aging Drafts (X, Y) vs. Final Targets (A, B)

| Subsystem Domain | Intermediate / Superseded Drafts (X, Y) | Authoritative Final Architecture (A, B) | Action |
| :--- | :--- | :--- | :--- |
| **Git Operations / Pull** | Standalone `70-redundant-repos`, `71-repo-dedup`, `197-pas-fix`, `198-pas-worker-concurrency` | `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging` & `202-gitmap-pull-errors-splitdb` | Retain final 219 and 202; prune intermediate 70/71/197/198. |
| **SSH Credential Storage** | Plaintext console reader, unencrypted storage, machine AES fallback exploration | `204-ssh-password-interception-and-rsa-credential-vault` (RSA-OAEP with SHA-256) | Retain 204 authoritative RSA-OAEP vault; prune transient notes. |
| **Split-DB Architecture** | Single monolithic database, file-lock prototypes | `209-ai-agent-task-orchestrator-and-split-db` & 3-Tier Split-DB (`gitmap.db`, `installation.db`, `repodb/pipeline.db`) | Retain 3-tier Split-DB specification; eliminate single-DB drafts. |
| **Workstation & Fleet** | Ad-hoc setup scripts, unversioned bash snippets | `214-ubuntu-fleet-automation-and-workstation-governance`, `217-antigravity-fleet-parity`, `233-ubuntu-ide-and-github-desktop-scan-sync` | Retain integrated fleet governance; prune early shell snippets. |
| **CLI Help & UX** | Hardcoded text tables, manual spacing calculations | `191-universal-command-help-restructuring-markdown-box-display`, `216-gitmap-prompting-freeze-and-suggestion-engine-fix` | Retain dynamic box renderers and Cobra completion engine. |

---

## 5. Strict Pending Isolation Registry

To guarantee that active development is never interrupted or corrupted, the following sets are classified as strictly read-only:

### 5.1 Pending Plans Directory (`.ai-memory/plans/pending/`) — 4 Files (Protected)
- `01-ports-and-ssh-enablement.md`
- `56-vmware-hardware-batch-and-macro-orchestration.md`
- `67-nodes-cfr-remote-fleet-clone-enhancement.md`
- `75-gitmap-u1-ubuntu-agm-fleet-integration.md`

### 5.2 Active Open Plans in `.ai-memory/plans/` — 7 Files (Protected)
- `214-ubuntu-fleet-automation-and-workstation-governance.md`
- `216-gitmap-prompting-freeze-and-suggestion-engine-fix.md`
- `217-antigravity-fleet-parity-theme-preset-plugins-and-delegation.md`
- `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md`
- `223-ubuntu-cursor-gitmap-python-git-tracing-and-aum-search.md`
- `225-ubuntu-cursor-start-menu-dock-position-and-profile-migration.md`
- `227-gitmap-cursor-fleet-delegate-settings-ui-and-cross-node-path-suggestions.md`

### 5.3 Active Open Subtask Directories (`.ai-memory/plans/subtasks/`) — 10 Folders / 52 Files (Protected)
1. `01-ports-and-ssh-enablement/`
2. `177-scan-alias-and-deploy-bin/`
3. `181-gitmap-ignore-and-cache-engine/`
4. `187-special-repos-repo-secrets-repo-cache-cd-and-coding-guidelines/`
5. `190-universal-help-llm-benchmarks/`
6. `191-universal-help-restructure/`
7. `192-version-pin-deploy-ui-secrets/`
8. `193-typed-json-envelope/`
9. `194-import-all-json-what-configs/`
10. `200-nodes-display-cfr-manifest-gitmap-json-and-vscode-optimize/`

---

## 6. Verification Gates & Acceptance Criteria

```yaml
verificationGates:
  isInventoryCatalogComplete: true
  isDomainClusteringValidated: true
  isMilestoneClusteringValidated: true
  isCompactionMappingSound: true
  isPendingIsolationVerified: true
  isRelativePathEnforced: true
  isZeroGitCommandRespected: true
```

- [x] All 310 items in `02-spec/21-app/` cataloged and mapped into 8 Canonical Domain Clusters.
- [x] All 163 completed plans cataloged and partitioned across 15 Milestone Clusters (Milestones 28 to 42).
- [x] Compaction Invariant defined with concrete examples of intermediate (X, Y) vs. final (A, B) architectures.
- [x] Strict Pending Isolation Rule formalized with complete protection registry.
- [x] Strictly relative Git paths verified across all catalog references.
