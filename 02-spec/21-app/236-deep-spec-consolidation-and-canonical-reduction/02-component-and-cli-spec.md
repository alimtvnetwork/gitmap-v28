# 236: Deep Spec Consolidation & Canonical Reduction — Component & CLI Specification

- **Specification Slug:** `236-deep-spec-consolidation-and-canonical-reduction`
- **Spec Document:** `02-component-and-cli-spec.md`
- **Spec Status:** `APPROVED`
- **Architecture Reference:** [01-architecture-spec.md](01-architecture-spec.md)
- **Master Ledger:** [00-master-audit-ledger.md](00-master-audit-ledger.md)
- **Target Subsystems:** `02-spec/21-app/`, `.ai-memory/plans/completed/`, `.ai-memory/memory/`, `.ai-memory/plans/readme.md`, `.ai-memory/what-to-read.md`
- **Baseline Release:** `v6.500.0`
- **Target Version:** `v6.501.0`

---

## 1. Executive Summary & Operational Boundaries

This component specification establishes the exhaustive operational rules, consolidation maps, schema contracts, file reduction manifests, and release verification gates for executing the second-generation deep consolidation across the GitMap codebase.

### 1.1 Core Principles
1. **Final Version Invariant:** Only the final ratified architectural version (contracts A and B) is preserved. Intermediate exploratory notes, provisional scratchpads, and iterative refactoring checklists (drafts X and Y) are permanently eliminated.
2. **Zero Technical Loss:** 100% of verified technical outcomes, Go type contracts (such as `*appfault.AppError`, `Result[T]`, and `StreamWriter`), database table schemas, CLI parameter flags, and test proofs are preserved in synthesized reference documents.
3. **Bloat Pruning:** Raw verbatim copies of task prompts and repetitive subtask chronicles (Section 3 dumps in legacy milestone documents) are completely excised, shedding over 9,000 lines of noise.
4. **Strict Pending and Active Isolation:** All files located in `.ai-memory/plans/pending/` and active WIP plans/subtasks remain untouched.
5. **Strict Relative Path Governance:** All cross-document links and filesystem references must strictly adhere to relative Git paths without exception.

---

## 2. Compaction Plan for Completed Plans (43 Milestones down to 18)

The historical milestone ledger currently contains 43 completed milestone files under `.ai-memory/plans/completed/`. Under this compaction plan, the directory is streamlined into exactly 18 cohesive milestone files:
- **Milestones 01–27** are synthesized into **2 Generational Milestones** (`01` and `02`).
- **Milestones 28–43** are renumbered sequentially as **Milestones 03 through 18**.

```
.ai-memory/plans/completed/
├── 01-generation-1-core-foundation-milestones.md
├── 02-generation-2-fleet-automation-and-modularization.md
├── 03-agy-prompts-templates-and-rerun-suite.md
├── 04-auto-aliasing-and-error-storage-reset.md
├── 05-ssh-exec-copy-mv-env-and-rm-sync-resilience.md
├── 06-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md
├── 07-pipeline-deep-eta-and-repo-folder-parity.md
├── 08-ai-scripts-engine-ssh-authkey-and-feature-parity.md
├── 09-generic-ai-scripts-creator-and-scaffolding-engine.md
├── 10-native-automation-engine-lazy-regex-and-benchmarks.md
├── 11-fleet-nodes-ping-envelope-and-remote-clone.md
├── 12-pas-worker-concurrency-pull-cache-and-split-db.md
├── 13-semantic-flat-commit-and-macro-execution-resilience.md
├── 14-ssh-password-interception-and-rsa-credential-vault.md
├── 15-ubuntu-fleet-migration-workstation-governance-and-customization.md
├── 16-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md
├── 17-codebase-review-remediation-and-preconsolidation-baseline.md
└── 18-app-spec-and-completed-plans-consolidation-and-reduction.md
```

### 2.1 Synthesis of Historical Milestones 01–27 into 2 Generational Milestones

#### 2.1.1 Generational Milestone 1: Core Foundation (`01-generation-1-core-foundation-milestones.md`)
Generational Milestone 1 synthesizes the foundational 15 milestone files (`01` through `15`), consolidating 568 KB of legacy markdown into a dense, authoritative architectural record.

| Source Milestone File | Original Scope | Synthesized Domain Focus | Preserved Type & Interface Contracts |
| :--- | :--- | :--- | :--- |
| `01-coding-guidelines-and-style-audits.md` | Coding guidelines, style linting | Static analysis & code style | Inverted guard clauses, positive booleans, blank line rules |
| `02-error-management-and-cliexit-architecture.md` | Error management, CLI exits | Centralized error hierarchy | `*appfault.AppError`, exit codes, error code enum registry |
| `03-type-safety-function-signatures-and-contracts.md` | Monadic Result types, signatures | Go type safety | `Result[T]`, `types.go` centralization, generic error wrappers |
| `04-cicd-pipelines-runners-and-streaming-telemetry.md` | CI/CD pipelines, local runners | Pipeline automation | `run.py`, GitHub Actions workflow telemetry, bounded test logs |
| `05-database-engine-sqlite-joins-and-scanners.md` | SQLite schema, multi-table joins | Database query engine | PascalCase tables, camelCase columns, SQLite foreign key enforcement |
| `06-git-operations-commit-engines-and-remediation.md` | Git commit engine, staging, push | Core Git operations | Flat commit semantics, auto-staging, push retry automation |
| `07-ssh-nodes-cluster-delegation-and-remote-exec.md` | Remote cluster delegation, SSH | Cluster execution | Node registry, SSH exec streaming, remote command timeouts |
| `08-terminal-ui-help-parity-and-cli-commands.md` | Terminal UI, help boxes, commands | CLI UI & formatting | Markdown box rendering, termpad, termtable styling |
| `09-chrome-profile-management-picker-and-token-vault.md` | Chrome profiles, token vault | Profile & auth storage | Profile discovery, encrypted tokens, profile switcher |
| `10-installers-multios-setup-and-web-stacks.md` | Cross-platform setup, web stack | Distribution & installers | NSIS installer script, Linux tarball generation, setup hooks |
| `11-completed-plans-consolidation.md` | Early plan consolidation | Milestone consolidation | Historical plan synthesis ledger |
| `12-git-deleted-files-tracer-and-purger.md` | Git deleted file tracer | Git hygiene & purging | `tracer.go`, dangling commit discovery, prune filters |
| `13-pipeline-sqlite-logs-and-history.md` | SQLite pipeline logging | Pipeline persistence | `pipeline.db` schema, stage duration history tables |
| `14-pipeline-agy-fix-injection-suite.md` | Pipeline AGY fix injection | Automated CI remediation | AGY error parser, failure injection hooks, patch verification |
| `15-ssh-multicommand-and-liveness-parity.md` | SSH multi-command execution | SSH reliability | Liveness ping, multi-command pipeline executor, error wrappers |

#### 2.1.2 Generational Milestone 2: Fleet Automation & Modularization (`02-generation-2-fleet-automation-and-modularization.md`)
Generational Milestone 2 synthesizes the intermediate 12 milestone files (`16` through `27`), consolidating architectural achievements in modularization, testing, and automation.

| Source Milestone File | Original Scope | Synthesized Domain Focus | Preserved Type & Interface Contracts |
| :--- | :--- | :--- | :--- |
| `16-cluster-sc-and-kubernetes-suite.md` | Cluster discovery, k8s commands | Cluster tooling | Node topology structures, cluster status commands |
| `17-macro-streaming-export-and-schedule.md` | Macro recording, replay streaming | Macro engine | Event streaming format, macro action replay loop, cron scheduler |
| `18-nuclear-package-modularization.md` | Modularizing monolith packages | Package architecture | Package dependency decoupling, DAG enforcement, acyclic imports |
| `19-smart-test-runner-and-inventory.md` | Incremental testing, ETA sleep | Smart test runner | Central test inventory manifest, duration profiling, test batches |
| `20-error-management-and-errorwrapper.md` | AppError wrappers across packages | Standardized error wrapping | Wrap with caller context, root cause chaining, user message mapping |
| `21-type-safety-result-and-types-go.md` | Result wrapper migration | Domain type extraction | Centralized `types.go`, removal of multi-value error maps |
| `22-database-transactions-and-schemas.md` | Database transactions, locks | DB concurrency | Reentrant lock pooling, transaction rollbacks, schema migrations |
| `23-installers-antigravity-and-archives.md` | Antigravity installers, packaging | Installer archives | Embedded assets extraction, self-extracting zip structures |
| `24-os-management-and-power-lifecycle.md` | OS power and reboot isolation | Test execution isolation | Mock system calls, safe test harnesses, OS shutdown guards |
| `25-terminal-ui-help-and-agy-prompts.md` | Terminal UI help, AGY prompts | Help & prompt rendering | Help topic hierarchy, command syntax flags, prompt templates |
| `26-coding-guidelines-and-linter-audits.md` | Coding guideline audits | Linter compliance | Zero-storage Actions, regex optimization, variable naming standards |
| `27-completed-plans-consolidation.md` | Plan consolidation milestone | Milestone accounting | Legacy milestone aggregation and cross-linking |

### 2.2 Renumbering and Modernization of Recent Milestones (28–43 to 03–18)

The 16 recent milestone documents are renumbered sequentially into `03` through `18`. Section 3 raw verbatim dumps are removed while preserving all technical outcome ledgers and Go type definitions.

| New Milestone ID & Slug | Former Milestone ID | Primary Architectural Domain | Consolidated Source Plans |
| :--- | :--- | :--- | :--- |
| `03-agy-prompts-templates-and-rerun-suite.md` | 28 | AGY prompt templates, decision logs, non-destructive rerun | Plans 112, 161, 164, 167 |
| `04-auto-aliasing-and-error-storage-reset.md` | 29 | Dynamic command auto-aliasing, internal error database, reset safety | Plans 174, 177, 202 |
| `05-ssh-exec-copy-mv-env-and-rm-sync-resilience.md` | 30 | SSH remote operations, streaming file copy, deploy key isolation | Plans 100, 101, 106, 163, 178 |
| `06-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md` | 31 | Pipeline forward slash normalization, cache clearing, ETA forecast | Plans 162, 175, 176 |
| `07-pipeline-deep-eta-and-repo-folder-parity.md` | 32 | Deep ETA profiling, heatmap telemetry, unit test traceback extraction | Plans 220, 221 |
| `08-ai-scripts-engine-ssh-authkey-and-feature-parity.md` | 33 | OS profiling, SSH installer payload detection, hermetic isolation | Plans 108, 109, 110, 171, 172 |
| `09-generic-ai-scripts-creator-and-scaffolding-engine.md` | 34 | Scaffolding scripts engine, guideline auto-sync, fleet update polish | Plans 170, 173, 179 |
| `10-native-automation-engine-lazy-regex-and-benchmarks.md` | 35 | Help modernization, markdown box displays, version pinning flags | Plans 180, 181, 182 |
| `11-fleet-nodes-ping-envelope-and-remote-clone.md` | 36 | Unified fleet nodes, ICMP/TCP ping, typed JSON envelopes, clone | Plans 183, 184, 185, 186, 188, 191, 192, 194, 196 |
| `12-pas-worker-concurrency-pull-cache-and-split-db.md` | 37 | GitMap PAS formula, worker concurrency, ignore engine, CPAR | Plans 57, 58, 60, 61, 62, 64, 66, 193, 195, 197, 198, 199, 201 |
| `13-semantic-flat-commit-and-macro-execution-resilience.md` | 38 | Semantic flat commit (`cm`), auto-staging, macro idempotency | Plans 54, 63, 65, 187, 189, 200 |
| `14-ssh-password-interception-and-rsa-credential-vault.md` | 39 | Masked password capture, user consent, RSA-OAEP credentials vault | Plans 67, 68, 78, 80, 203, 204, 208, 210 |
| `15-ubuntu-fleet-migration-workstation-governance-and-customization.md` | 40 | Ubuntu migration, GNOME scaling, automount, workstation governance | Plans 52, 76, 81, 206, 211, 212, 214, 215, 222 |
| `16-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md` | 41 | AGY UI plugins, skills sync, terminal freeze fix, multi-IDE sync | Plans 83, 213, 216, 217, 218, 223, 224, 225, 227, 228, 229, 233 |
| `17-codebase-review-remediation-and-preconsolidation-baseline.md` | 42 | Root clutter purge, readme index reduction, token purge | Plans 230, 231, 232, 234 |
| `18-app-spec-and-completed-plans-consolidation-and-reduction.md` | 43 | App spec consolidation into 8 clusters, completed subtask folding | Plans 235 |

### 2.3 Elimination of Section 3 Raw Verbatim Dumps

In legacy milestones, Section 3 contained multi-thousand-line dumps of raw subtask instructions and verbatim planning logs that created severe navigational bloat:
- **Milestone 01:** 690 lines of raw task copies in Section 3.
- **Milestone 04:** 1,406 lines of raw task copies in Section 3.
- **Milestone 06:** 1,120 lines of raw task copies in Section 3.
- **Milestone 10:** 2,850 lines of raw task copies in Section 3.
- **Milestone 05 & 08:** ~1,800 lines of raw task copies in Section 3.

#### The Streamlining Contract
1. Section 3 raw text dumps are replaced with a dense, structured **Architecture & Type Contract Summary Table**.
2. All technical outcomes, Go structs, error enums, and database migrations are preserved in concise bulleted technical ledgers.
3. This elimination saves over **9,000 lines of markdown clutter** across the completed plans directory.

---

## 3. Compaction Plan for AI Memory (71 Files down to 10 Files = 85.9% Reduction)

The directory `.ai-memory/memory/` contains 71 files distributed across 15 subdirectories and root markdown files. This sprawling layout is consolidated into **8 Canonical Domain Memory References**, plus `00-project-governance-and-invariants.md` and a clean `readme.md`, achieving an 85.9% file count reduction.

### 3.1 Target Canonical Structure

```
.ai-memory/memory/
├── 00-project-governance-and-invariants.md
├── 01-cli-architecture-and-contracts.md
├── 02-scanner-projects-and-deduplication.md
├── 03-git-operations-commit-and-pull.md
├── 04-fleet-nodes-ssh-and-credentials.md
├── 05-antigravity-and-ide-ecosystem.md
├── 06-database-engine-and-split-storage.md
├── 07-pipeline-diagnostics-and-telemetry.md
├── 08-distribution-installers-and-release.md
└── readme.md
```

### 3.2 Purging the 36 Obsolete Legacy Files

The following 36 files contain superseded, duplicate, or point-in-time notes. Their valid invariants are incorporated into the 10 canonical files, and the obsolete files are permanently removed:

#### 3.2.1 Obsolete Files in `tech/` (7 Files)
1. `tech/01-version-json-architecture.md` (superseded by `08-distribution-installers-and-release.md`)
2. `tech/ci-pipeline-architecture.md` (superseded by `07-pipeline-diagnostics-and-telemetry.md`)
3. `tech/static-analysis-security.md` (superseded by `00-project-governance-and-invariants.md`)
4. `tech/release-pipeline.md` (superseded by `08-distribution-installers-and-release.md`)
5. `tech/ci-release-automation.md` (superseded by `08-distribution-installers-and-release.md`)
6. `tech/dependency-management.md` (superseded by `00-project-governance-and-invariants.md`)
7. `tech/versioning-strategy.md` (superseded by `08-distribution-installers-and-release.md`)

#### 3.2.2 Obsolete Files in `workflow/` (4 Files)
8. `workflow/04-ci-hardening-session.md` (historical session log; superseded by pipeline telemetry)
9. `workflow/01-plan.md` (superseded by `.ai-memory/plans/readme.md`)
10. `workflow/02-ssh-plan.md` (superseded by `04-fleet-nodes-ssh-and-credentials.md`)
11. `workflow/03-gitmap-dir-plan.md` (superseded by `02-scanner-projects-and-deduplication.md`)

#### 3.2.3 Obsolete Files in `constraints/` (4 Files)
12. `constraints/constants-ownership.md` (superseded by `00-project-governance-and-invariants.md`)
13. `constraints/ci-release-pipeline-untouchable.md` (superseded by `00-project-governance-and-invariants.md`)
14. `constraints/clone-preserves-version-folder.md` (superseded by `02-scanner-projects-and-deduplication.md`)
15. `constraints/strictly-prohibited.md` (superseded by `.ai-memory/strictly-avoid.md`)

#### 3.2.4 Obsolete Files in `style/` (3 Files)
16. `style/02-enums-audit-results.md` (point-in-time audit report)
17. `style/01-ts-enums-and-query-wrappers.md` (superseded by `02-spec/02-coding-guidelines/`)
18. `style/code-quality-improvement.md` (superseded by `00-project-governance-and-invariants.md`)

#### 3.2.5 Obsolete Files in `project/` (5 Files)
19. `project/01-overview.md` (superseded by `00-project-governance-and-invariants.md`)
20. `project/version-bump-procedure.md` (superseded by `08-distribution-installers-and-release.md`)
21. `project/unified-directory-structure.md` (superseded by `00-project-governance-and-invariants.md`)
22. `project/what-to-read.md` (duplicate of root `.ai-memory/what-to-read.md`)
23. `project/release-keyword.md` (superseded by `08-distribution-installers-and-release.md`)

#### 3.2.6 Obsolete Files in Root Memory & Misc Folders (13 Files)
24. `last-failure.md` (transient failure state log)
25. `ssh-public-key-display-and-clipboard.md` (superseded by `04-fleet-nodes-ssh-and-credentials.md`)
26. `01-index.md` (redundant index)
27. `01-replace-command-plan.md` (one-off completed command plan)
28. `index.md` (redundant legacy index)
29. `03-v3.12.1-session.md` (obsolete historical session log)
30. `release-architecture-map.md` (superseded by `08-distribution-installers-and-release.md`)
31. `learned.md` (redundant stub file)
32. `02-v15-legacy-compat-audit.md` (obsolete v15 compatibility audit)
33. `reports/20260723-rejog-reliability.md` (point-in-time historical report)
34. `specs/02-lfs-smudge-rca.md` (historical RCA note)
35. `specs/01-version-and-code-quality-rules.md` (superseded by coding guidelines)
36. `suggestions/01-suggestions.md`, `suggestions/readme.md`, `suggestions/_template.md` (obsolete scratch stubs)

### 3.3 Domain Scope of the 8 Canonical Memory References

#### Reference 00: `00-project-governance-and-invariants.md`
- Core architecture invariants: Strict relative paths, positive booleans, no-build execution rule, branch immutability.
- Codebase guidelines reference links (`02-spec/02-coding-guidelines/`).
- Isolation boundaries for pending and active plans.

#### Reference 01: `01-cli-architecture-and-contracts.md`
- Cobra command initialization and registration patterns.
- Universal command help renderer and markdown box display format.
- Typed JSON envelope specification (`v1` and `v2`) with status enums.
- Typo suggestion engine, predictive suggestions, and terminal clear primitives.

#### Reference 02: `02-scanner-projects-and-deduplication.md`
- High-speed repository discovery scanner algorithm and recursion limits.
- Project type detection heuristics (Go, Node, Python, Rust, PHP, C#).
- Case-insensitive OS path normalization and `EqualFoldAnyTrim` usage.
- Ignore rules (`.gitmapignore`) and split-DB repo cache lifecycle.

#### Reference 03: `03-git-operations-commit-and-pull.md`
- Semantic flat commit suite (`gitmap commit`, `cm`) and auto-staging behavior.
- Pull-all fast mode (`pa`, `pat`) and PAS worker concurrency formula.
- Oh-My-Zsh ignore rules and git status error recovery.
- Revert and amend command safety invariants.

#### Reference 04: `04-fleet-nodes-ssh-and-credentials.md`
- Fleet nodes command architecture (`gitmap nodes`, `gitmap cluster`).
- Dual-stack reachability probing (ICMP and TCP ping).
- RSA-OAEP credentials vault, master password hashing, and user consent gates.
- Remote clone, except-self clone filters, and remote execution streaming.

#### Reference 05: `05-antigravity-and-ide-ecosystem.md`
- Google Antigravity (AGY) agent architecture, toolchain, and prompt manager.
- Non-destructive rerun commands and decision logs database.
- Multi-IDE scan sync (VS Code, Cursor, Windsurf, JetBrains).
- Desktop dock positioning, start menu configurations, and settings sync.

#### Reference 06: `06-database-engine-and-split-storage.md`
- Three-tier SQLite database architecture (`gitmap.db`, `installation.db`, `pipeline.db`).
- PascalCase table schema design and camelCase column naming rules.
- Connection pooling, reentrant locking, and busy timeout configurations.
- Foreign key constraints, cascade triggers, and schema migration protocols.

#### Reference 07: `07-pipeline-diagnostics-and-telemetry.md`
- Pipeline Error (`pe`) extraction, failure traceback capture, and regex parsers.
- Heatmap telemetry data structures and historical stage duration forecasting.
- 4-part Root Cause Analysis (RCA) methodology.
- Automated pipeline remediation injection pipelines.

#### Reference 08: `08-distribution-installers-and-release.md`
- NSIS installer payload auto-detection and Linux tar/gz/zip installer automation.
- Cross-platform runner architecture (`run.ps1`, `run.sh`, `run.config.json`).
- Semantic Versioning automation (`version.json`, `package.json`, subpackages).
- Release ceremony steps: version bump, changelog, tag push, and backup branch.

#### Root Index: `readme.md`
- Master navigation catalog providing quick links and summary abstracts for all 8 domain references and the governance invariant document.

---

## 4. Strict Isolation Boundary for Pending & Active Work

To prevent disruption to ongoing features, a strict isolation boundary is enforced.

### 4.1 Untouchable Work-in-Progress Artifacts

1. **Pending Plans (`.ai-memory/plans/pending/`):**
   - `56-vmware-hardware-batch-and-macro-orchestration.md`
   - `75-gitmap-u1-ubuntu-agm-fleet-integration.md`
   - `77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md`
   - Any plan enqueued under `.ai-memory/plans/pending/`.
2. **Active Open Plans (`.ai-memory/plans/`):**
   - `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md`
   - `226-ubuntu-cursor-memories-conversations-and-projects-migration.md`
   - `236-deep-spec-consolidation-and-canonical-reduction.md` (current plan)
3. **Active Subtask Directories (`.ai-memory/plans/subtasks/`):**
   - `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/`
   - `226-ubuntu-cursor-memories-conversations-and-projects-migration/`
   - `236-deep-spec-consolidation-and-canonical-reduction/`
   - `56-vmware-hardware-batch-and-macro-orchestration/`
   - `75-gitmap-u1-ubuntu-agm-fleet-integration/`
   - All subtask files corresponding to active or pending tasks.

---

## 5. Navigational Registry & Index Synchronization Plan

### 5.1 Application Specification Index (`02-spec/21-app/readme.md`)
The application specification index will be rewritten to:
1. Index only the active specifications contained within the **8 Canonical Domain Clusters**:
   - `01-cli-architecture/`
   - `02-scanner-and-projects/`
   - `03-git-operations-and-pull/`
   - `04-fleet-nodes-and-ssh/`
   - `05-antigravity-and-ide/`
   - `06-database-and-split-db/`
   - `07-pipeline-and-diagnostics/`
   - `08-distribution-and-release/`
2. Remove dead markdown links pointing to folded or deleted specifications.
3. Provide a historical spec number mapping table to enable fast lookup of old spec numbers.

### 5.2 Plans Registry (`.ai-memory/plans/readme.md`)
The plans index will be updated to:
1. Present the **Completed Plans** section referencing the 18 consolidated milestones (`01` through `18`).
2. Retain Active Plans and Pending Plans sections with complete fidelity.
3. Update the **Subtasks Directory** section to reference active subtask directories only.

### 5.3 Cognitive Map (`.ai-memory/what-to-read.md`)
The primary developer orientation document will be updated to:
1. Add a changelog entry detailing Task 236 achievements and file reduction metrics.
2. Update the "Before any task" and "Before writing code" reading paths to point directly to:
   - `.ai-memory/memory/00-project-governance-and-invariants.md`
   - The relevant Domain Memory Reference file in `.ai-memory/memory/`
   - The relevant Canonical Cluster Specification in `02-spec/21-app/`
   - `.ai-memory/plans/completed/` milestones

---

## 6. Quality Verification Protocols

Before committing and releasing, all files must satisfy four strict verification gates:

### Gate 1: Relative Path Integrity Linter
- **Command:** `python3 linter-scripts/check-relative-paths.py`
- **Requirement:** Zero absolute filesystem paths anywhere in modified files, specifications, memory documents, or release notes.

### Gate 2: Forbidden Strings & Secrets Linter
- **Command:** `python3 linter-scripts/check-forbidden-strings.py`
- **Requirement:** Zero occurrences of banned patterns, private tokens, passwords, or deprecated negative booleans across the repository.

### Gate 3: Markdown Document Cross-Link Linter
- **Command:** `python3 03-ai-scripts/22-doc-path-linter.py`
- **Requirement:** Zero broken markdown relative links across `02-spec/`, `.ai-memory/`, and `01-prompts/`.

### Gate 4: Go Static Analysis Gate
- **Command:** `go vet ./...`
- **Requirement:** Zero compiler errors, static analysis diagnostics, or malformed package structures.

---

## 7. Final Release Ceremony & Compaction Scorecard (v6.501.0)

### 7.1 Minor Version Bump Protocol
Following successful verification:
1. Bump minor version from `v6.500.0` to `v6.501.0` via `python3 03-ai-scripts/37-bump-version.py --minor` (or release orchestrator).
2. Ensure version manifests are updated synchronously:
   - `version.json` (root version `6.501.0`)
   - `package.json`
   - `.gitmap/release/latest.json`
   - `readme.md` (badges and version header)
   - `what-to-read.md`
   - `cli/constants/constants.go`
   - `changelog.md`
3. Push tag `v6.501.0` and branch `release/v6.501.0` to remote origin.

### 7.2 Compaction Metrics Scorecard

| Target Subsystem | Baseline Count (`v6.500.0`) | Consolidated Count (`v6.501.0`) | Net File Reduction | Net Reduction % |
| :--- | :--- | :--- | :--- | :--- |
| `02-spec/21-app/` Files & Dirs | ~45 items | ~16 items (8 clusters + docs) | ~29 items | **~64.4%** |
| `.ai-memory/plans/completed/` | 43 milestone files | 18 milestone files | 25 files | **58.1%** |
| `.ai-memory/plans/completed/` Lines | ~35,000 lines | ~12,000 lines | ~23,000 lines | **~65.7%** |
| `.ai-memory/memory/` Files | 71 files | 10 files | 61 files | **85.9%** |
| **Total Target Files** | **159 items** | **44 items** | **115 items** | **72.3% net reduction** |

---

## 8. Positive Boolean & Naming Standards

In strict compliance with repository coding guidelines, all new components, configurations, and verification scripts enforce positive boolean naming conventions:

| Positive Identifier | Type | Operational Meaning | Deprecated Negative Equivalent |
| :--- | :--- | :--- | :--- |
| `isConsolidated` | `bool` | True if milestone or spec has been consolidated | `isNotConsolidated` |
| `isAuthoritative` | `bool` | True if spec represents final ratified version | `isNotObsolete` |
| `isPreserved` | `bool` | True if contract or proof is preserved in archive | `isNotPruned` |
| `hasVerifiedOutcome` | `bool` | True if subtask verification proof passed | `hasNoFailure` |
| `isPendingIsolated` | `bool` | True if pending plan is protected from modification | `isNonPendingModified` |
| `isPathRelative` | `bool` | True if filesystem path adheres to relative Git standard | `isNotAbsolutePath` |
| `hasActiveWork` | `bool` | True if directory contains active WIP subtasks | `isNotIdle` |
| `canRelease` | `bool` | True if all verification gates passed | `cannotRelease` |
