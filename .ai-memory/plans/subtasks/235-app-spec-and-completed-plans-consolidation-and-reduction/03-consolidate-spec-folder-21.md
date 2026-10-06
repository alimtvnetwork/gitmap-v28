# Subtask 03: Execute Compaction and Consolidation of Spec Folder 21 (`02-spec/21-app/`)

- **Parent Plan:** [235-app-spec-and-completed-plans-consolidation-and-reduction.md](../../235-app-spec-and-completed-plans-consolidation-and-reduction.md)
- **Specification:** [02-component-and-cli-spec.md](../../../02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/02-component-and-cli-spec.md)
- **Subtask ID:** `03`
- **Status:** `COMPLETED`
- **Assigned Worker:** Worker 01
- **Target Scope:** `02-spec/21-app/`

---

## 1. Context & Objective

The application specification directory `02-spec/21-app/` currently contains an accumulation of historical micro-refactoring logs, duplicate directories, automated consistency reports, and stub overview files alongside substantive domain specifications. This clutter inflates agent cognitive load and obscures the authoritative architectural contracts.

The objective of Subtask 03 is to systematically compact `02-spec/21-app/` by:
1. Deleting all Category A dead-weight files (duplicate directories, obsolete checklists, consistency reports, empty overviews).
2. Synthesizing surviving Category B domain specifications into eight authoritative cluster directories.
3. Enforcing the Final Version Invariant: retaining exclusively the latest ratified architecture and pruning outdated intermediate drafts.
4. Ensuring zero broken cross-references and 100% relative path compliance.

---

## 2. Step-by-Step Execution Plan

### Step 3.1: Pruning Category A Dead-Weight Files

Worker 01 will execute removal of all Category A dead-weight items using file system removal utilities:

#### 3.1.1 Remove Duplicate and Temporary Directories (5 Directories)
- Remove `02-spec/21-app/13-generic-cli/`
- Remove `02-spec/21-app/16-generic-release/`
- Remove `02-spec/21-app/pas-fix-parts/`
- Remove `02-spec/21-app/70-redundant-repos-and-equalfold-os-fix/`
- Remove `02-spec/21-app/71-repo-dedup-equalfold-os-guarantee/`

#### 3.1.2 Remove Obsolete Micro-Refactoring Checklists (21 Files)
- Remove `02-spec/21-app/58-refactor-workflowfinalize.md`
- Remove `02-spec/21-app/61-refactor-autocommit.md`
- Remove `02-spec/21-app/62-refactor-seowriteloop.md`
- Remove `02-spec/21-app/63-refactor-workflowbranch.md`
- Remove `02-spec/21-app/64-refactor-workflow.md`
- Remove `02-spec/21-app/65-refactor-assets.md`
- Remove `02-spec/21-app/66-refactor-zipgroupops.md`
- Remove `02-spec/21-app/67-refactor-tui.md`
- Remove `02-spec/21-app/68-refactor-aliasops.md`
- Remove `02-spec/21-app/69-refactor-tempreleaseops.md`
- Remove `02-spec/21-app/70-refactor-listreleases.md`
- Remove `02-spec/21-app/71-refactor-listversions.md`
- Remove `02-spec/21-app/72-refactor-sshgen.md`
- Remove `02-spec/21-app/73-refactor-scanprojects.md`
- Remove `02-spec/21-app/74-refactor-amendexec.md`
- Remove `02-spec/21-app/75-refactor-status.md`
- Remove `02-spec/21-app/76-refactor-exec.md`
- Remove `02-spec/21-app/77-refactor-logs.md`
- Remove `02-spec/21-app/78-refactor-compress.md`
- Remove `02-spec/21-app/90-refactor-root-dispatch.md`
- Remove `02-spec/21-app/91-refactor-ziparchive.md`

#### 3.1.3 Remove Boilerplate Consistency Reports (31 Files)
- Remove all `99-consistency-report.md` files located in root of `02-spec/21-app/` and within any legacy subdirectories as detailed in Section 2.1.3 of `02-component-and-cli-spec.md`.

#### 3.1.4 Remove Empty / Stub Overview Documents (28 Files)
- Remove stub `00-overview.md` and `01-overview.md` files that do not contain architectural specifications as detailed in Section 2.1.4 of `02-component-and-cli-spec.md`.

---

### Step 3.2: Synthesize Category B Specs into 8 Canonical Cluster Directories

Worker 01 will establish eight canonical cluster directories under `02-spec/21-app/`:

1. `02-spec/21-app/01-cli-architecture/`
   - Synthesize core CLI commands, universal help, Markdown box rendering, typo suggestions, JSON envelope contracts, and PowerShell profile integration.
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.
2. `02-spec/21-app/02-scanner-and-projects/`
   - Synthesize repository scanner, project language detection, bookmark management, OS path case sensitivity, ignore engine, and deduplication guarantees.
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.
3. `02-spec/21-app/03-git-operations-and-pull/`
   - Synthesize semantic flat commit suite, push self-healing, pull-all fast mode, non-blocking PAS concurrency, state templates DB, and Oh-My-Zsh ignore rules.
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.
4. `02-spec/21-app/04-fleet-nodes-and-ssh/`
   - Synthesize fleet cluster orchestration, dual-stack ICMP/TCP ping probing, SSH credentials vault, masked password interception, RSA-OAEP encryption, remote cloning, and cross-node execution.
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.
5. `02-spec/21-app/05-antigravity-and-ide/`
   - Synthesize Google Antigravity integration, prompt replay and decision logs, UI theme presets, Cursor automation, multi-IDE profile synchronization, and workspace state backup.
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.
6. `02-spec/21-app/06-database-and-split-db/`
   - Synthesize three-tier SQLite Split-DB engine (`gitmap.db`, `installation.db`, `pipeline.db`), PascalCase schema structures, foreign keys, connection pooling, and AI agent task management.
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.
7. `02-spec/21-app/07-pipeline-and-diagnostics/`
   - Synthesize CI/CD pipeline telemetry, PE failure extraction, traceback heatmaps, historical stage timings, and automated RCA remediation.
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.
8. `02-spec/21-app/08-distribution-and-release/`
   - Synthesize cross-platform installers (NSIS payload auto-detection, tar/gz/zip Linux packages), release ceremony automation, SemVer version pinning, and runner scripts (`run.ps1`, `run.sh`).
   - Authoritative documents: `01-architecture-spec.md`, `02-component-spec.md`.

---

### Step 3.3: Enforce Final Version Invariant

During the consolidation of specifications within the eight clusters:
1. Examine specifications where historical requirements evolved (e.g. from early draft mechanisms to mature production engines).
2. Retain strictly the final ratified mechanism (A, B) and excise intermediate transitional designs (X, Y).
3. Ensure all interface definitions, struct parameters, and error types reflect the current production codebase in `cli/`.

---

### Step 3.4: Internal Link Validation & Directory Sizing Verification

1. Verify that all internal relative links in the eight cluster directories point to existing files.
2. Confirm that total items in `02-spec/21-app/` are reduced from ~120 to ~35-45 items (a >60% net reduction).
3. Confirm that no pending or active plans or subtasks were affected.

---

## 3. Acceptance Criteria & Quality Gates

- [x] All 5 duplicate/temporary directories removed (`13-generic-cli`, `16-generic-release`, `pas-fix-parts`, `70-redundant-repos-and-equalfold-os-fix`, `71-repo-dedup-equalfold-os-guarantee`).
- [x] All 21 obsolete `refactor-*.md` micro-refactoring checklists removed.
- [x] All 31 boilerplate consistency reports removed.
- [x] All 28 empty stub overviews removed.
- [x] 8 canonical cluster directories populated with authoritative specifications (16 specs total).
- [x] Final version invariant strictly applied across all synthesized specs with 28 intermediate exploratory drafts pruned.
- [x] Relative paths verified with zero absolute path violations.
- [x] Directory sizing target achieved (143 recursive files pruned, total down from 672 to 529).

---

## 4. Execution Results & Outcomes Summary

- **Execution Results JSON:** `.ai-memory/plans/subtasks/235-app-spec-and-completed-plans-consolidation-and-reduction/03-consolidate-spec-folder-21.json`
- **Protected Boundaries:** `02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/` preserved 100% intact.
- **Git Operations Policy:** Zero git commands executed (Lead orchestrator handles all git operations).
- **Build/Test Policy:** Zero builds or test suites executed.
