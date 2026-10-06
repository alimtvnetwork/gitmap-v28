# Subtask 04: Consolidate Completed Plans into Milestones 28–42 and Compact AI Memory

- **Parent Plan:** [235-app-spec-and-completed-plans-consolidation-and-reduction.md](../../235-app-spec-and-completed-plans-consolidation-and-reduction.md)
- **Specification:** [02-component-and-cli-spec.md](../../../02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/02-component-and-cli-spec.md)
- **Subtask ID:** `04`
- **Status:** `COMPLETED`
- **Assigned Worker:** Worker 02
- **Target Scope:** `.ai-memory/plans/completed/`, `.ai-memory/plans/subtasks/`, `.ai-memory/memory/`

---

## 1. Context & Objective

The repository currently maintains 190 completed plan files in `.ai-memory/plans/completed/` (with 163 completed plans created since Milestone 27), 229 subtask files across 42 directories in `.ai-memory/plans/subtasks/` (including 177 completed subtasks across 32 completed folders), and over 85 fragmented memory notes in `.ai-memory/memory/`.

The objective of Subtask 04 is to:
1. Merge the 163 post-Milestone-27 completed plans into 15 dense, authoritative milestone summaries (Milestones 28 to 42) in `.ai-memory/plans/completed/`.
2. Fold the 177 completed subtasks across 32 folders into their respective milestone documents, preserving 100% of verified architectural outcomes and error contracts.
3. Remove the 32 completed subtask folders once synthesis is verified.
4. Compact `.ai-memory/memory/` down to approximately 35 clean, cohesive reference summaries.
5. Strictly isolate all pending plans, active open plans, and active subtask folders.

---

## 2. Step-by-Step Execution Plan

### Step 4.1: Pre-Flight Strict Isolation Verification

Prior to any file manipulation, Worker 02 will verify the isolation boundary protecting open work:
1. Confirm that `.ai-memory/plans/pending/` contains 4 pending files that remain untouched:
   - `56-vmware-hardware-batch-and-macro-orchestration.md`
   - `75-gitmap-u1-ubuntu-agm-fleet-integration.md`
   - `77-gitmap-db-reset-ssh-key-management-firewall-and-project-detection-rca.md`
   - And any pending plan restored via V6 protocol.
2. Confirm that active open plans in `.ai-memory/plans/` remain untouched:
   - `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging.md`
   - `226-ubuntu-cursor-memories-conversations-and-projects-migration.md`
   - `235-app-spec-and-completed-plans-consolidation-and-reduction.md`
3. Confirm that the 10 active subtask directories (52 files total) under `.ai-memory/plans/subtasks/` remain untouched:
   - `219-ubuntu-pull-all-remediation-omz-ignore-and-error-logging/`
   - `226-ubuntu-cursor-memories-conversations-and-projects-migration/`
   - `235-app-spec-and-completed-plans-consolidation-and-reduction/`
   - `56-vmware-hardware-batch-and-macro-orchestration/`
   - `75-gitmap-u1-ubuntu-agm-fleet-integration/`
   - And all other active work directories.

---

### Step 4.2: Merging 163 Completed Plans into Milestones 28–42

Worker 02 will author 15 dense milestone summaries in `.ai-memory/plans/completed/`, consolidating 163 source plans while preserving all technical decisions, Go type contracts (`*appfault.AppError`, `Result[T]`), and verified testing outcomes:

1. `28-agy-prompts-templates-and-rerun-suite.md`
   - Synthesizes Plans 112, 161, 164, 167: Antigravity prompt templates, decision logs, non-destructive rerun commands, and green automation.
2. `29-auto-aliasing-and-error-storage-reset.md`
   - Synthesizes Plans 174, 177, 202: Auto-aliasing migration, internal errors database, and DB reset safety mechanisms.
3. `30-ssh-exec-copy-mv-env-and-rm-sync-resilience.md`
   - Synthesizes Plans 100, 101, 106, 163, 178: Remote SSH execution, streaming uploads, deploy keys, and smart file deployment sync modes.
4. `31-pipeline-repo-folder-forward-slash-clear-and-accurate-eta.md`
   - Synthesizes Plans 162, 175, 176: Pipeline stage timings, forward slash path normalization, cache clearing, and accurate ETA forecasting.
5. `32-pipeline-deep-eta-and-repo-folder-parity.md`
   - Synthesizes Plans 220, 221: Deep ETA profiling, heatmap telemetry, unit test traceback extraction, and nested if linter fixes.
6. `33-ai-scripts-engine-ssh-authkey-and-feature-parity.md`
   - Synthesizes Plans 108, 109, 110, 171, 172: First-time login OS profiling, SSH installer payload detection, hermetic test isolation, and auto-update guards.
7. `34-generic-ai-scripts-creator-and-scaffolding-engine.md`
   - Synthesizes Plans 170, 173, 179: Coding guideline auto-sync guards, fleet liveness error remediation, and fleet update JSON polish.
8. `35-native-automation-engine-lazy-regex-and-benchmarks.md`
   - Synthesizes Plans 180, 181, 182: Universal command help modernization, Markdown box displays, and version pinning flags.
9. `36-fleet-nodes-ping-envelope-and-remote-clone.md`
   - Synthesizes Plans 183, 184, 185, 186, 188, 191, 192, 194, 196: Unified `gitmap nodes` command, ICMP/TCP ping, typed JSON envelopes, and except-self clone.
10. `37-pas-worker-concurrency-pull-cache-and-split-db.md`
    - Synthesizes Plans 57, 58, 60, 61, 62, 64, 66, 193, 195, 197, 198, 199, 201: GitMap PAS formula, worker concurrency, ignore engine, CPAR, and split-DB repo cache.
11. `38-semantic-flat-commit-and-macro-execution-resilience.md`
    - Synthesizes Plans 54, 63, 65, 187, 189, 200: Semantic flat commit suite (`gitmap commit`, `cm`), auto-staging, and macro idempotency.
12. `39-ssh-password-interception-and-rsa-credential-vault.md`
    - Synthesizes Plans 67, 68, 78, 80, 203, 204, 208, 210: Masked password capture, user consent, RSA-OAEP salt credential encryption vault, and push self-healing.
13. `40-ubuntu-fleet-migration-workstation-governance-and-customization.md`
    - Synthesizes Plans 52, 76, 81, 206, 211, 212, 214, 215, 222: Windows to Ubuntu fleet migration, GNOME 140% scaling, VMware automount, and workstation governance.
14. `41-antigravity-fleet-parity-prompting-freeze-and-ide-sync.md`
    - Synthesizes Plans 83, 213, 216, 217, 218, 223, 224, 225, 227, 228, 229, 233: AGY UI plugins, skills sync, Win32 terminal freeze remediation, typo suggestions, and multi-IDE scan sync.
15. `42-codebase-review-remediation-and-preconsolidation-baseline.md`
    - Synthesizes Plans 230, 231, 232, 234: Root clutter purge, readme index reduction, security token purge, and baseline v6.498.0 release.

---

### Step 4.3: Folding and Clean Removal of 177 Completed Subtasks

1. In each synthesized milestone document, include a dedicated **Subtask Verification Ledger** mapping the folded subtask files, their verified criteria, and their code changes.
2. Once the synthesis is complete and validated, remove the 32 completed subtask directories under `.ai-memory/plans/subtasks/` listed in Section 3.2 of `02-component-and-cli-spec.md`.
3. Preserve the 10 active subtask directories completely untouched.

---

### Step 4.4: Compacting `.ai-memory/memory/` into ~35 Reference Summaries

Worker 02 will synthesize fragmented notes across `.ai-memory/memory/`:
1. **Learned Memory (`.ai-memory/memory/learned/`):** Consolidate 40+ granular files into 15 domain-focused summaries:
   - `01-project-context-and-standards.md`
   - `02-cli-contracts-and-help-architecture.md`
   - `03-scanner-deduplication-and-path-collation.md`
   - `04-git-operations-flat-commit-and-pull.md`
   - `05-ssh-fleet-ping-and-remote-exec.md`
   - `06-rsa-credential-vault-and-auth.md`
   - `07-antigravity-integration-and-ide-sync.md`
   - `08-sqlite-split-db-and-cache-lifecycle.md`
   - `09-pipeline-telemetry-and-rca-remediation.md`
   - `10-cross-platform-installers-and-runners.md`
   - `11-macro-automation-and-idempotency.md`
   - `12-ubuntu-workstation-governance-and-customization.md`
   - `13-streamwriter-contracts-and-concurrency.md`
   - `14-coding-guidelines-and-positive-booleans.md`
   - `15-fast-file-reader-and-testing-hygiene.md`
2. **Issues & RCA Memory (`.ai-memory/memory/issues/`):** Consolidate micro-RCA files into 10 grouped failure analyses.
3. **Avoid Rules (`.ai-memory/memory/avoid/`):** Consolidate into 5 core anti-pattern guides.
4. **Features Memory (`.ai-memory/memory/features/`):** Consolidate into 5 capability catalogs.

---

## 3. Acceptance Criteria & Quality Gates

- [x] All 163 completed plans clustered into Milestones 28–42.
- [x] All 177 completed subtask files folded with zero loss of verified outcomes.
- [x] 32 completed subtask directories removed.
- [x] 4 pending plans and active open plans (219, 226, 235) remain 100% untouched.
- [x] 10 active subtask directories remain 100% untouched.
- [x] `.ai-memory/memory/` consolidated into ~35 clean reference files.
- [x] Net reduction in `.ai-memory/plans/` and subtasks exceeds 75%.
- [x] All references use strictly relative Git paths.

---

## 4. Execution Summary & Metrics

- **Completed Plans:** Reduced from 190 to exactly 42 milestones (148 files removed, **77.89% reduction**).
- **Subtasks:** 32 completed subtask directories (174 files) removed, leaving 11 active subtask directories (55 files) 100% untouched (**75.98% reduction**).
- **Pending Isolation:** 4 pending plans in `.ai-memory/plans/pending/` and 18 active root plans preserved verbatim.
- **AI Memory Consolidation:** Consolidated down to 36 authoritative domain reference summaries across `learned/` (15), `issues/` (10), `avoid/` (5), `features/` (5), and `plans/` (1).
- **Path Portability:** All file references across markdown documents adhere strictly to relative Git paths.

