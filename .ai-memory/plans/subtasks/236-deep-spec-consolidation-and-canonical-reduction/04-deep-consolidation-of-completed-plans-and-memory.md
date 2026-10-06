# Subtask 04: Deep Consolidation of Completed Plans (43 to 18 Milestones) and AI Memory (71 to 10 Files)

- **Parent Plan:** [236-deep-spec-consolidation-and-canonical-reduction.md](../../236-deep-spec-consolidation-and-canonical-reduction.md)
- **Specification:** [02-component-and-cli-spec.md](../../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/02-component-and-cli-spec.md)
- **Subtask ID:** `04`
- **Status:** `COMPLETED`
- **Assigned Worker:** Worker 02
- **Target Scope:** `.ai-memory/plans/completed/`, `.ai-memory/memory/`

---

## 1. Context & Objective

The objective of Subtask 04 is to execute the second-generation deep compaction of historical completed plans and AI memory files:
1. **Completed Plans Compaction (43 to 18 Milestones):**
   - Synthesize historical Milestones 01–27 into 2 dense Generational Milestones:
     - `01-generation-1-core-foundation-milestones.md` (consolidating Milestones 01–15)
     - `02-generation-2-fleet-automation-and-modularization.md` (consolidating Milestones 16–27)
   - Renumber the 16 recent milestones (former 28–43) as Milestones `03` through `18`.
   - Excise all Section 3 raw verbatim dumps, saving over 9,000 lines of noise while preserving 100% of verified technical outcomes and Go type contracts (`Result[T]`, `*appfault.AppError`, `StreamWriter`).
2. **AI Memory Compaction (71 to 10 Files = 85.9% Reduction):**
   - Purge the 36 obsolete legacy files across `tech/`, `workflow/`, `constraints/`, `style/`, `project/`, root memory, and scratch folders.
   - Unify surviving architectural knowledge into 8 Canonical Domain Memory References + `00-project-governance-and-invariants.md` + clean `readme.md` (10 files total).
3. **Strict Isolation Boundary:**
   - Enforce an untouchable boundary protecting `.ai-memory/plans/pending/`, active open plans (`219`, `226`, `236`), and active subtasks.

---

## 2. Step-by-Step Operational Implementation Plan

### Step 4.1: Pre-Flight Strict Isolation Verification

Before executing any file modifications, Worker 02 must verify the isolation boundary protecting open work:
1. **Pending Plans Verification:**
   - Verify `.ai-memory/plans/pending/` exists and contains all enqueued pending files.
   - Assert `isPendingIsolated == true` (no modifications permitted under `.ai-memory/plans/pending/`).
2. **Active Plans Verification:**
   - Confirm active plans remain untouched in `.ai-memory/plans/`:
     - `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md`
     - `226-ubuntu-cursor-memories-conversations-and-projects-migration.md`
     - `236-deep-spec-consolidation-and-canonical-reduction.md` (active self-plan)
3. **Active Subtasks Verification:**
   - Confirm active subtask folders under `.ai-memory/plans/subtasks/` remain untouched:
     - `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/`
     - `226-ubuntu-cursor-memories-conversations-and-projects-migration/`
     - `236-deep-spec-consolidation-and-canonical-reduction/`
     - `56-vmware-hardware-batch-and-macro-orchestration/`
     - `75-gitmap-u1-ubuntu-agm-fleet-integration/`

---

### Step 4.2: Operational Execution of Completed Plans Compaction (43 -> 18)

#### 4.2.1 Author `01-generation-1-core-foundation-milestones.md`
Worker 02 will author the first generational milestone consolidating historical milestones 01 through 15:
1. **Consolidated Scope:**
   - Milestones 01–03: Coding guidelines, inverted guard clauses, `*appfault.AppError`, and `Result[T]` contracts.
   - Milestones 04–06: CI/CD runner architecture, SQLite joins, and flat commit staging.
   - Milestones 07–09: SSH cluster delegation, terminal UI markdown boxes, and Chrome profile token vaults.
   - Milestones 10–12: Multi-OS installers, plan consolidations, and deleted files tracer.
   - Milestones 13–15: Pipeline history logging, AGY fix injections, and SSH multi-command execution.
2. **Excision of Verbatim Section 3 Dumps:**
   - Exclude raw prompt dumps and task chronicles.
   - Retain dense Architectural Decisions tables, verified test commands, and Go type definitions.
3. **Preserved Go Type Contracts:**
   ```go
   // Preserved core error envelope
   type AppError struct {
       Code       ErrorCode
       Message    string
       Err        error
       Caller     string
       StatusCode int
   }

   // Preserved monadic result wrapper
   type Result[T any] struct {
       Value T
       Error *appfault.AppError
   }
   ```

#### 4.2.2 Author `02-generation-2-fleet-automation-and-modularization.md`
Worker 02 will author the second generational milestone consolidating historical milestones 16 through 27:
1. **Consolidated Scope:**
   - Milestones 16–18: Kubernetes/cluster CLI commands, macro streaming export, and nuclear package modularization.
   - Milestones 19–21: Smart test runner inventory manifests, standardized AppError wrapping, and `types.go` centralization.
   - Milestones 22–24: SQLite transaction reentrant locks, Antigravity packaging archives, and OS mock isolation test harnesses.
   - Milestones 25–27: Terminal UI prompt rendering, zero-storage Actions linting, and legacy milestone deduplication.
2. **Excision of Verbatim Section 3 Dumps:**
   - Exclude raw prompt dumps.
   - Retain verified test results, package dependency graphs, and lock acquisition protocols.

#### 4.2.3 Renumber and Streamline Recent Milestones (28–43 -> 03–18)
Worker 02 will renumber the 16 recent milestone documents:
- Former `28` -> `03-agy-prompts-templates-and-rerun-suite.md`
- Former `29` -> `04-auto-aliasing-and-error-storage-reset.md`
- Former `30` -> `05-ssh-exec-copy-mv-env-and-rm-sync-resilience.md`
- Former `31` -> `06-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md`
- Former `32` -> `07-pipeline-deep-eta-and-repo-folder-parity.md`
- Former `33` -> `08-ai-scripts-engine-ssh-authkey-and-feature-parity.md`
- Former `34` -> `09-generic-ai-scripts-creator-and-scaffolding-engine.md`
- Former `35` -> `10-native-automation-engine-lazy-regex-and-benchmarks.md`
- Former `36` -> `11-fleet-nodes-ping-envelope-and-remote-clone.md`
- Former `37` -> `12-pas-worker-concurrency-pull-cache-and-split-db.md`
- Former `38` -> `13-semantic-flat-commit-and-macro-execution-resilience.md`
- Former `39` -> `14-ssh-password-interception-and-rsa-credential-vault.md`
- Former `40` -> `15-ubuntu-fleet-migration-workstation-governance-and-customization.md`
- Former `41` -> `16-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md`
- Former `42` -> `17-codebase-review-remediation-and-preconsolidation-baseline.md`
- Former `43` -> `18-app-spec-and-completed-plans-consolidation-and-reduction.md`

#### 4.2.4 Removal of Historical Milestone Files 01–27
Once synthesis into `01` and `02` is validated, remove the 27 original historical files (`01-coding-guidelines-...md` through `27-completed-plans-...md`).

---

### Step 4.3: Operational Execution of AI Memory Compaction (71 -> 10)

#### 4.3.1 Purge Manifest for 36 Obsolete Legacy Memory Files
Worker 02 will execute the deletion of the 36 obsolete legacy files across `.ai-memory/memory/`:
- **`tech/` (7 files):** `01-version-json-architecture.md`, `ci-pipeline-architecture.md`, `static-analysis-security.md`, `release-pipeline.md`, `ci-release-automation.md`, `dependency-management.md`, `versioning-strategy.md`
- **`workflow/` (4 files):** `04-ci-hardening-session.md`, `01-plan.md`, `02-ssh-plan.md`, `03-gitmap-dir-plan.md`
- **`constraints/` (4 files):** `constants-ownership.md`, `ci-release-pipeline-untouchable.md`, `clone-preserves-version-folder.md`, `strictly-prohibited.md`
- **`style/` (3 files):** `02-enums-audit-results.md`, `01-ts-enums-and-query-wrappers.md`, `code-quality-improvement.md`
- **`project/` (5 files):** `01-overview.md`, `version-bump-procedure.md`, `unified-directory-structure.md`, `what-to-read.md`, `release-keyword.md`
- **Root Memory (10 files):** `last-failure.md`, `ssh-public-key-display-and-clipboard.md`, `01-index.md`, `01-replace-command-plan.md`, `index.md`, `03-v3.12.1-session.md`, `release-architecture-map.md`, `learned.md`, `02-v15-legacy-compat-audit.md`, `readme.md` (legacy version)
- **Scratch Folders (3+ files):** `reports/20260723-rejog-reliability.md`, `specs/01-version-and-code-quality-rules.md`, `specs/02-lfs-smudge-rca.md`, `suggestions/` stub files

#### 4.3.2 Author 8 Canonical Domain Memory References + Governance + Readme
Worker 02 will author the 10 consolidated reference files in `.ai-memory/memory/`:

1. `00-project-governance-and-invariants.md`:
   - Non-negotiable architectural invariants: positive booleans, 100% relative paths, zero-build execution, branch immutability.
   - Codebase structure and quality guidelines.
2. `01-cli-architecture-and-contracts.md`:
   - Cobra CLI setup, command dispatching, terminal boxes, typo suggestions, JSON envelope schema.
3. `02-scanner-projects-and-deduplication.md`:
   - Scanner traversal algorithm, project detection heuristics, OS path collation, `.gitmapignore`, split-DB repo cache.
4. `03-git-operations-commit-and-pull.md`:
   - Flat commit mechanics (`gitmap commit`, `cm`), auto-stage, PAS worker concurrency formula, OMZ ignore handling.
5. `04-fleet-nodes-ssh-and-credentials.md`:
   - Fleet nodes topology, ICMP/TCP ping probing, RSA-OAEP credentials vault, masked password capture, remote clone.
6. `05-antigravity-and-ide-ecosystem.md`:
   - AGY agent system, prompt manager, decision logs, multi-IDE sync, desktop workspace configurations.
7. `06-database-engine-and-split-storage.md`:
   - Three-tier SQLite architecture (`gitmap.db`, `installation.db`, `pipeline.db`), PascalCase tables, locking, transaction isolation.
8. `07-pipeline-diagnostics-and-telemetry.md`:
   - CI/CD pipeline telemetry, PE failure extraction, heatmap telemetry, historical stage timings, 4-part RCA.
9. `08-distribution-installers-and-release.md`:
   - Cross-platform runners (`run.ps1`, `run.sh`), NSIS/archive installers, SemVer bump mechanics, release ceremonies.
10. `readme.md`:
    - Clean master navigation catalog linking all 8 domain references and the governance invariant document.

---

### Step 4.4: Validation & Quality Checks for Preserved Contracts

Worker 02 will execute validation checks across synthesized files:
1. **Contract Integrity Check:**
   - Confirm `*appfault.AppError`, `Result[T]`, and database schema models are documented with full type fidelity.
2. **Positive Booleans Check:**
   - Verify all boolean identifiers across new documents follow positive naming conventions (`isConsolidated`, `isAuthoritative`, `isPreserved`, `hasVerifiedOutcome`, `isPendingIsolated`).
3. **Relative Path Verification:**
   - Confirm all markdown cross-links use strictly relative paths (`../../02-spec/...`, `../../.ai-memory/...`).
4. **File Count Verification:**
   - `.ai-memory/plans/completed/`: Exactly 18 files.
   - `.ai-memory/memory/`: Exactly 10 files.

---

## 3. Acceptance Criteria & Quality Gates

- [x] All pending plans in `.ai-memory/plans/pending/` verified untouched (`isPendingIsolated == true`).
- [x] Active open plans and active subtask folders verified untouched.
- [x] Historical Milestones 01–27 synthesized into `01-generation-1-core-foundation-milestones.md` and `02-generation-2-fleet-automation-and-modularization.md`.
- [x] Recent Milestones 28–43 renumbered sequentially to `03` through `18`.
- [x] Section 3 verbatim task dumps excised, eliminating 9,000+ lines of raw bloat.
- [x] 36 obsolete legacy files across `.ai-memory/memory/` permanently purged.
- [x] 8 Canonical Domain Memory References + `00-project-governance-and-invariants.md` + clean `readme.md` authored (10 files total).
- [x] 100% of verified technical outcomes and Go type contracts preserved.
- [x] All markdown links adhere strictly to relative Git paths.
