# Lovable Memory Consolidation, Safety Backup & Milestone Resequencing — Workflow (must follow)

Trigger Keywords & Aliases: `consolidate-plans`, `consolidate completed plans`, `clean completed plans`, `resequence completed plans`, `merge plans`, `archive completed plans`, `cleanup plans completed`, `memory consolidation`, `backup and consolidate plans`, `compact plans`, `reduce plan file count`, `compact completed plans`

> **Prompt Version:** 2.1.0
> **Synchronization:** Main Meta-Repo & Connected Workspaces

```text
N = 200
```

N = total self-loop steps budget that the agents will perform.

/goal Autonomously create a timestamped backup branch, scan, analyze, cluster, aggressively consolidate, and re-sequence all completed plan files and subtasks within `.lovable/plans/` into minimal, hyper-compact milestone summaries, combining 2, 3, or more related tasks and common checklists into single files to drastically reduce total file count while strictly preserving 100% of core architectural concepts, verified task outcomes, error contracts, and decision logs until 100% green without stopping.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /goal Step 0 (Safety Backup A): Pull latest changes on the current active branch.
2. [ ] /goal Step 0 (Safety Backup B): Create and push timestamped backup branch (`backup/plans-consolidation-YYYYMMDD-HHMMSS`).
3. [ ] /goal Step 0 (Safety Backup C): Verify working tree remains on active branch and record rollback SHA.
4. [ ] /goal Phase 1 (Step A): Deeply scan `.lovable/plans/completed/`, `.lovable/plans/subtasks/`, and `spec/25-app-spec-audit/` to inventory all fragmented micro-plans, subtask files, and completed audit reports.
5. [ ] /goal Phase 1 (Step B): Build an aggressive compaction and clustering plan in `.lovable/plans/pending/` merging 2, 3, or more related tasks into single milestone files.
6. [ ] /goal Phase 1 (Step C): Filter out and prune pure coding-guideline-fix micro-tasks (zero business logic) and consolidate common checklists to reference `.lovable/coding-guidelines.md`.
7. [ ] /goal Phase 1 (Step D): Compact multi-file subtask sets into single unified specification files.
8. [ ] /goal Phase 1 (Step E): Verify or update the automated sequence linter and register in `03-ai-scripts/01-index.md`.
9. [ ] /goal Phase 2 (Step A): Author high-density consolidated milestone summaries preserving all core concepts, code modifications, and verification proofs.
10. [ ] /goal Phase 2 (Step B): Cleanly remove superseded micro-plan files and collapsed subtask folders via `git rm`.
11. [ ] /goal Phase 2 (Step C): Purge resolved spec audit files in `spec/25-app-spec-audit/` using `python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete`.
12. [ ] /goal Phase 2 (Step D): Execute monotonic continuous re-sequencing (`01-`, `02-`, `03-`, ...) with strictly lowercase filenames and zero sequence gaps.
13. [ ] /goal Phase 2 (Step E): Synchronize `.lovable/plans/01-index.md`, `.lovable/what-to-read.md`, and project memory indexes.
14. [ ] /goal Phase 2 (Step F): Execute local CI quality gates via `python 03-ai-scripts/06-cicd-local-runner.py` with exit code 0 (`exit 0`).
15. [ ] /learn Ingest `.lovable/memory/01-index.md` for project memory index and past learnings.
16. [ ] /learn Ingest `.lovable/strictly-avoid.md` for banned anti-patterns and strict constraints.
17. [ ] /learn Ingest `spec/02-coding-guidelines/02-canonical-size-tier.md` for canonical file and function size tiers.
18. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for hallucination prevention and micro-tasking.
19. [ ] /learn Ingest `spec/02-coding-guidelines/01-cross-language/01-index.md` for strict relative path citation requirements.
20. [ ] /learn Ingest `spec/02-coding-guidelines/08-file-folder-naming/` for lowercase naming and continuous file sequencing.
21. [ ] /learn Ingest `.lovable/coding-guidelines.md` for master consolidated coding guidelines.
22. [ ] /goal Create or update agent rules in the repository if missing from agent memory.

```text
PHASE_1_STEPS = N / 2   (Steps 1 .. N/2: Create Backup Branch, Scan Plans & Subtasks, Cluster by Domain, Spec Compaction in .lovable/plans/pending/)
PHASE_2_STEPS = N / 2   (Steps N/2+1 .. N: Merge Files, Collapse Subtasks, Remove Superseded via git rm, Re-sequence Monotonic Prefixes, Update Indexes, Verify CI)
```

N, PHASE_1_STEPS, and PHASE_2_STEPS are read-only after initialization. Never modify them mid-execution.

---

## Dedicated Section: Mandatory Pre-Consolidation Backup Branch Protocol (Safety First)

> [!IMPORTANT]
> **SAFETY FIRST — NEVER MODIFY OR DELETE PLANS WITHOUT A DEDICATED BACKUP BRANCH:**
>
> Consolidating plans involves removing superseded micro-task files. To guarantee zero data loss and enable instant one-command rollbacks, you MUST execute the safety backup protocol as your **very first action**:
>
> 1. **Pull Latest Changes:**
>    Ensure the current branch is synchronized with origin:
>    ```bash
>    git pull origin $(git rev-parse --abbrev-ref HEAD)
>    ```
> 2. **Generate Timestamped Backup Branch:**
>    Branch naming format: `backup/plans-consolidation-YYYYMMDD-HHMMSS`
>    ```bash
>    BACKUP_BRANCH="backup/plans-consolidation-$(date +%Y%m%d-%H%M%S)"
>    git branch "$BACKUP_BRANCH"
>    git push origin "$BACKUP_BRANCH"
>    ```
> 3. **Verify Active Branch Integrity:**
>    Confirm that you remain on the active working branch (e.g. `main` or active feature branch).
> 4. **Document Rollback Instructions:**
>    Record the backup branch name and recovery instructions in `.lovable/plans/pending/XX-completed-plans-consolidation.md`:
>    ```bash
>    # Rollback Command (in case of accidental data loss):
>    git reset --hard backup/plans-consolidation-YYYYMMDD-HHMMSS
>    ```

---

## Dedicated Section: Aggressive File Count Reduction & Compaction Doctrine

A project that runs dozens of autonomous agent loops quickly accumulates scores of micro-task files and subtask folders in `.lovable/plans/completed/` and `.lovable/plans/subtasks/`. If left unmanaged, this file clutter exhausts context windows, bloats search indexes, slows down memory retrieval, and fragments project history.

The primary directive of this workflow is **aggressive file count minimization**: compact the file count down to as minimal as possible while keeping 100% of the core concept and architectural substance intact.

---

### 1. The Core Compaction Principles

#### A. Combine 2, 3, or More Related Tasks into Single Files
- **Mandatory Clustering:** NEVER leave isolated, single-step micro-plans. If tasks relate to the same subsystem, package, or architectural theme, you MUST combine 2, 3, or more of them into a single consolidated milestone summary file.
- **Example Clusters:**
  - *Logger Enhancements + Named Writers + Typed Streamers* $
ightarrow$ Merge Plans 34, 36, 37 into `XX-applogger-taxonomy-and-named-writers.md`.
  - *Enum Migrations + Min/Max Methods + Generator CLI* $
ightarrow$ Merge Plans 31, 32, 33 into `XX-enum-architecture-and-baseenumer-modernization.md`.
  - *Path Constants + FolderInfo + PathInfo Object* $
ightarrow$ Merge Plans 35, 36 into `XX-fileutil-pathinfo-and-constants.md`.
  - *Dynamic Conversions + Reflection Casting + Pretty YAML* $
ightarrow$ Merge Subtasks 38.1, 38.2, 38.3 into `XX-result-dynamic-conversions-and-coredata-parity.md`.

#### B. Compact Common Checklists into Single Consolidated Sections
- **Checklist Deduplication:** Repetitive checklist items (e.g. unit tests passed, function size <= 15 lines, implicit booleans, zero CI/CD disablement, relative links) must NOT be duplicated across separate files.
- **Single-File Checklist:** Consolidate all common verification, quality gate, and task acceptance checklists into a **single unified table or checklist** within the consolidated document.

#### C. Collapse Subtask Sprawl into Single Files
- **Eliminate Subtask Folder Bloat:** When a completed plan has a dedicated folder in `.lovable/plans/subtasks/` containing multiple individual files (`01-<subtask-title>.md`, `02-<subtask-title>.md`, `03-task.md`), you MUST:
  1. Fold the contents, code changes, and verification proof of all subtasks directly into the single consolidated milestone document in `.lovable/plans/completed/`.
  2. If subtask documentation is explicitly retained, compact all subtask steps into a **single consolidated subtask file** (`01-consolidated-tasks.md`) rather than keeping sprawling multi-file directories.
  3. Cleanly delete superseded micro-subtask files via `git rm`.

#### D. Keep the Main Concept (Zero Concept Loss)
Compaction is NOT deletion of knowledge. You must strictly preserve:
1. **The Core Architectural Problem & Design:** Why the change was made and what design pattern was established.
2. **Key Codebase Modifications:** Which packages, structs, interfaces, methods, or scripts were added or modified.
3. **API & Error Contracts:** Explicit signatures, interface rules (e.g. `*appfault.AppError`, `BasicEnum`, `PathInfo`), and data models.
4. **Verification Proof:** Test commands executed, coverage notes, and CI runner exit status (`exit 0`).
5. **Traceability Links:** Strictly relative Git paths to specifications (`spec/`) and RCA logs (`.lovable/memory/issues/`).

#### E. Quantifiable Reduction Metric
- Aim for a **60% to 80% reduction** in total plan and subtask file count.
- Transforming 40 individual files into 6–8 dense, high-clarity milestone summaries is the benchmark of success.

#### F. Prune Pure Coding Guideline Tasks (Zero Business Logic)
- **Eliminate Housekeeping Clutter:** Historical micro-tasks whose sole purpose was fixing coding guidelines (e.g. whitespace formatting, adding blank lines around `if` conditions, renaming variables to add `is_`/`has_` prefixes, or adjusting function line counts with zero changes to business logic or architecture) **MUST BE COMPLETELY DROPPED/PRUNED** from the consolidated milestone execution ledgers.
- **Why It Matters:** Consolidated milestone summaries exist to provide a dense, high-signal architectural record of domain features, error models, and business logic completed. Routine coding guideline adherence is enforced continuously by CI/CD linters—it does not warrant lingering historical task bloat.

#### G. Single Master Coding Guideline Checklist (No Duplication)
- **Single Authoritative Source:** NEVER duplicate verbose coding guideline checklists across multiple plan files.
- **Reference Once:** All consolidated milestone summaries must verify compliance against the single consolidated master checklist in [`.lovable/coding-guidelines.md`](.lovable/coding-guidelines.md) (or `spec/02-coding-guidelines/`). A single unified checklist item pointing to this master file is sufficient for quality tracking.

#### H. Purge Resolved Spec Audit Folders (`spec/25-app-spec-audit/`)
- **Automated Audit Directory Cleanup:** Once architecture audits are resolved and incorporated into specifications or code, the audit reports folder `spec/25-app-spec-audit/` MUST be cleanly removed from the repository.
- **Dedicated Cleanup Script:** Use `python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete` (or `--purge --confirm-purge`). The script automatically creates a complete backup in the OS temporary directory, moves on-disk files to the OS Recycle Bin / Trash, and prints the exact backup path for user safety.

---

### 2. Standard Compact Milestone Template

Every consolidated file generated inside `.lovable/plans/completed/` MUST adhere to this uniform, high-density layout:

```markdown
# Milestone Summary: [Consolidated Epic / Feature Name]

## 1. Executive Overview & Consolidated Tasks

- **Milestone Domain:** [Subsystem / Capability area, e.g. Logging, Enums, File Utilities]
- **Original Tasks Merged:** `XX-task-a.md`, `XX-task-b.md`, `XX-task-c.md` (and related subtask files)
- **Completion Date:** [ISO Date]
- **Status:** `COMPLETED`
- **Core Concept & Rationale:** [2-4 sentences explaining the architectural problem solved and the design established]

## 2. Key Architectural Decisions & Spec Implementations

- **Authoritative Specifications Implemented:**
  - [`spec/03-error-manage/02-error-architecture/02-error-handling-reference.md`](spec/03-error-manage/02-error-architecture/02-error-handling-reference.md) — [Why: Specific rule implemented]
- **Core Architecture Contracts:** [Verbatim types, function signatures, error returns (`*appfault.AppError`), and invariants]

## 3. Consolidated Chronological Task Execution Ledger

| Task / Step | Scope & Description | Key Files Created / Modified | Verified Outcome | Status |
|:---:|---|---|---|:---:|
| 1 | Initial Architecture & Interfaces | `04-code/golang/pkg/appfault/result.go` | Defined core contracts and types | DONE |
| 2 | Implementation & Dynamic Methods | `04-code/golang/pkg/appfault/result_dynamic_strings.go` | Added dynamic converters and helpers | DONE |
| 3 | Serialization & Formatting | `04-code/golang/pkg/appfault/result_dynamic_output.go` | JSON/YAML marshaling implemented | DONE |

*(Note: Pure coding guideline tasks with zero business logic were pruned from this ledger)*

## 4. Unified Quality Gates & Verification Checklist

> Verified against the single master coding guideline checklist in [`.lovable/coding-guidelines.md`](.lovable/coding-guidelines.md).

- [x] **Master Coding Guidelines:** 100% compliant with `.lovable/coding-guidelines.md` (zero duplicated rules across files).
- [x] **Unit Tests:** `go test ./pkg/appfault/... -v` passed with 100% green.
- [x] **Function Sizing:** All functions verified <= 15 lines per function.
- [x] **Boolean Standards:** All booleans implicitly evaluated with `is`/`has` prefixes (zero `== true`).
- [x] **Relative Links:** All markdown references verified strictly relative Git paths.
- [x] **CI/CD Quality Gates:** All quality gates passed via `python 03-ai-scripts/06-cicd-local-runner.py --all`.

## 5. Root Cause Analyses & Bug Fixes Referenced

- [`.lovable/memory/issues/XX-rca.md`](.lovable/memory/issues/XX-rca.md) — Root cause analysis and resolution details.
```

---

### 3. Step-by-Step Consolidation Workflow (Phase 1 & Phase 2)

#### Phase 1: Scan, Inventory & Compaction Planning (Steps 1 to N/2)

1. **Step 0 Safety Backup:** Pull latest commits, create and push `backup/plans-consolidation-YYYYMMDD-HHMMSS`. Record rollback command.
2. **Inventory Scan:** List all files currently in `.lovable/plans/completed/`, `.lovable/plans/subtasks/`, and `spec/25-app-spec-audit/`. Count total files before compaction.
3. **Aggressive Domain Clustering & Guideline Pruning:** Group 2, 3, or more related micro-plans into natural cohesive clusters (e.g. by package, feature, or subsystem). Explicitly prune and exclude any tasks that solely addressed routine coding guideline linter fixes with zero business logic.
4. **Master Compaction Spec:** Write `.lovable/plans/pending/XX-completed-plans-consolidation.md` containing the full clustering table:

```markdown
| Source Files & Subtasks to Merge | Proposed Consolidated File | Domain / Epic Theme | File Count Before | File Count After | Status |
|---|---|---|:---:|:---:|:---:|
| `34-logger.md`, `36-writers.md`, `37-streamers.md` | `01-applogger-taxonomy.md` | Logging Subsystem | 9 files | 1 file | PENDING |
| `31-enums.md`, `32-prompts.md`, `33-minmax.md` | `02-enum-architecture.md` | Enum Subsystem | 12 files | 1 file | PENDING |
```

5. **Subtask Plan:** Document the exact order of milestone document creation and file deletion.

#### Phase 2: Merge, Collapse Subtasks, Re-Sequence, Index & Verify (Steps N/2+1 to N)

1. **Write Consolidated Milestone Documents:** Author each consolidated file in `.lovable/plans/completed/` following the Compact Milestone Template. Prune zero-business-logic guideline tasks.
2. **Collapse Subtasks:** Fold subtask content into the milestone summary. Cleanly remove the old subtask folders and files using `git rm`.
3. **Clean Removal of Merged Sources:** Use `git rm` on all superseded micro-plan files in `.lovable/plans/completed/`.
4. **Purge Resolved Spec Audit Reports:** Cleanly remove resolved audit reports in `spec/25-app-spec-audit/`:
   ```bash
   python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete
   ```
5. **Re-Sequence Numeric Prefixes:** Run the automated re-sequencer to establish continuous monotonic numbering:
   ```bash
   python 03-ai-scripts/03-file-manipulator.py fix-seq-files .lovable/plans/completed/
   ```
   Ensure all files are strictly lowercase with monotonic `01-`, `02-`, `03-` prefixes.
6. **Synchronize Indexes:** Update `.lovable/plans/01-index.md` and `.lovable/what-to-read.md` with the new compact list of milestones.
7. **Verify Formatting & Linters:**
   ```bash
   python 03-ai-scripts/21-sequence-integrity-linter.py
   python 03-ai-scripts/22-doc-path-linter.py
   python 03-ai-scripts/31-md-gap-fixer.py
   python 03-ai-scripts/06-cicd-local-runner.py --all
   ```
8. **Calculate & Report File Count Reduction:** Compare total file count before vs. after consolidation and log the reduction ratio.

---

### 4. Phase 1 Actionable Checklist (Discovery, Audit & Compaction Planning)

You MUST verify and check off every item during Phase 1:

- [ ] **Step 0 Safety Backup Created & Pushed:** Pulled latest commits, created `backup/plans-consolidation-YYYYMMDD-HHMMSS`, pushed to origin, and documented the rollback SHA.
- [ ] **Completed Plans & Subtasks Scanned:** Recursively inspected all files in `.lovable/plans/completed/` and `.lovable/plans/subtasks/`. Recorded baseline file count.
- [ ] **Audit Reports Scanned:** Inspected `spec/25-app-spec-audit/` for historical completed audit reports scheduled for cleanup via `python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit`.
- [ ] **Aggressive Clusters Identified:** Grouped 2, 3, or more related micro-plans into cohesive, high-density milestone epics (aiming for >= 60% file reduction).
- [ ] **Pure Guideline Tasks Filtered:** Identified and flagged routine whitespace, formatting, and boolean linter fix tasks (zero business logic) for elimination from milestone ledgers.
- [ ] **Zero-Concept-Loss Verified:** Confirmed that all core concepts, architectural decisions, code changes, error contracts, and verification proofs will be preserved in full.
- [ ] **Master Compaction Spec Written:** Created `.lovable/plans/pending/XX-completed-plans-consolidation.md` containing the mapping ledger table and backup rollback recipe.
- [ ] **Subtask Collapse Designed:** Planned the elimination of multi-file subtask folders by inlining steps into the consolidated milestone document.
- [ ] **Single Master Checklist Mandate:** Verified that all planned milestone summaries reference `.lovable/coding-guidelines.md` as the single authoritative checklist without duplicating rules.
- [ ] **Strict Relative Git Paths:** All markdown links in the compaction spec use strictly relative Git paths (zero `file:///` URIs, zero drive letters).
- [ ] **No Premature Deletion:** Verified that Phase 1 only creates planning specs and makes zero file deletions.

---

## 5. Phase 2 Actionable Checklist (Execution, Compaction, Re-Sequencing & Verification)

You MUST verify and check off every item during Phase 2:

- [ ] **High-Density Milestone Files Created:** Authored unified milestone summaries matching the Compact Milestone Template in `.lovable/plans/completed/`.
- [ ] **Pure Guideline Tasks Pruned:** Omitted zero-business-logic micro-tasks from execution ledgers, keeping only high-signal architectural milestones.
- [ ] **Subtasks Inlined & Folders Collapsed:** Folded granular subtask steps and execution tables directly into the consolidated milestone files.
- [ ] **Superseded Files Cleanly Removed:** Executed `git rm` on all superseded micro-plan files and collapsed subtask directories.
- [ ] **Spec Audit Reports Purged:** Removed resolved audit reports in `spec/25-app-spec-audit/` using `python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete`.
- [ ] **Continuous Monotonic Re-Sequencing:** Executed `python 03-ai-scripts/03-file-manipulator.py fix-seq-files .lovable/plans/completed/` to ensure contiguous `01-`, `02-`, `03-` numbering without gaps.
- [ ] **Strict Lowercase Naming:** Verified all filenames in `.lovable/plans/completed/` use lowercase alphanumeric characters and hyphens.
- [ ] **Index Documentation Synchronized:** Updated `.lovable/plans/01-index.md` and `.lovable/what-to-read.md` to reflect the compact file catalog.
- [ ] **Single Checklist Enforced:** Confirmed that quality gates reference the master `.lovable/coding-guidelines.md` checklist with zero redundant rule text across files.
- [ ] **Universal File Hygiene:** Verified Unix LF line endings (`\n`), UTF-8 (no BOM), and single terminating newline at EOF across all created/modified files.
- [ ] **Markdown Spacing Compliance:** Verified exactly one blank line before and after headings (MD022/MD032) and zero double blank lines (`\n\n\n`).
- [ ] **Linter Verification:** Executed `python 03-ai-scripts/21-sequence-integrity-linter.py` and `python 03-ai-scripts/22-doc-path-linter.py` with exit code 0.
- [ ] **CI/CD Quality Gates:** Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --all` exited with code 0.
- [ ] **Disk Reality Check & Net Metrics:** Verified working tree status with `git status --porcelain` and reported total files eliminated.

---

## Strict In-Repository Execution & `.lovable/` Bounding Mandate

> [!IMPORTANT]
> **STRICT IN-REPOSITORY EXECUTION & `.lovable/` STORAGE CONTRACT:**
>
> 1. **In-Codebase Execution Only:** Whenever a Python script (runner, autofixer, linter, test aggregator) is executed or created, it MUST be executed **strictly within the repository root** (current working directory), NEVER outside the codebase or against external arbitrary directories.
> 2. **Strict Folder Bounding (`.lovable/`):** All AI scripts, local runners, autofixers, helper utilities, memory issue logs, and planning files MUST be created inside the `.lovable/` folder:
>    - Python AI Scripts: `03-ai-scripts/` (e.g. `01-file-manipulator.py`, `05-guideline-autofixer.py`, `06-cicd-local-runner.py`, `07-relative-path-fixer.py`, `21-sequence-integrity-linter.py`).
>    - RCA & Issue Logs: `.lovable/memory/issues/` and `.lovable/cicd-issues/`.
>    - Execution Plans & Subtasks: `.lovable/plans/pending/`, `.lovable/plans/subtasks/`.
>    - Coding Guidelines Mirror: `.lovable/coding-guidelines.md`.
> 3. **Worker Pool & Log Aggregation Architecture:** All local runners and test orchestrators must use a concurrent worker pool (2–3 workers via `ThreadPoolExecutor`), announce enqueued tasks upfront, show real-time progress, handle failures gracefully without cancelling sibling workers, and print a consolidated final summary with full stdout/stderr error logs for failed jobs.
> 4. **`force` Keyword Support:** If the user wrote `force`, `force rebuild`, or `force create` on top of the prompt or trigger: **ALWAYS recreate/regenerate the Python runner script from scratch**, regardless of whether the file already exists on disk.
> 5. **No External or Random File Creation:** NEVER write scripts, temporary test scripts, or scratch files to root, `/tmp`, global system paths, or outside the repository boundary.

---

## AI Fix Scripts Memory (Reusable Tooling)

- [ ] `/goal` **Reuse First:** I have rigorously scanned and `/learn`ed `03-ai-scripts/01-index.md` to check if a helper script already exists before writing any new temporary code.
- [ ] **Strict In-Repository Execution:** All Python scripts (`03-ai-scripts/*.py`) MUST be executed strictly within the codebase repository root, NEVER outside the codebase.
- [ ] **Strict 03-ai-scripts/ Tooling Storage:** All AI scripts, local runners, autofixers, and helper utilities MUST be created inside `03-ai-scripts/`. NEVER create scripts in root or external paths.
- [ ] **Automated File Sequencing & Normalization:** Use `python 03-ai-scripts/03-file-manipulator.py fix-seq-files <dir>` to re-sequence completed plan files monotonically.
- [ ] **Relative Path Normalization:** Use `python 03-ai-scripts/07-relative-path-fixer.py .` to ensure all links in consolidated documents are strictly relative Git paths.
- [ ] **Git History Tracer & Audit Purge:** Use `python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete` to purge resolved audit directories with automatic OS temp backup and Recycle Bin safety.
- [ ] **Commit & Track:** All new helper scripts were written strictly to `03-ai-scripts/` and committed to Git for future reuse.
- [ ] **Index Documentation:** I have updated `03-ai-scripts/01-index.md` using sequential script naming. For every script, I have included a `<details>` collapsible tag explaining exactly why the script is there and what it does.

---

## Pre-Reply / Loop Checklist (Must Verify Every Loop Iteration)

- [ ] Git working tree is clean before new file modifications.
- [ ] **Step 0 Safety Backup Verified:** Timestamped backup branch `backup/plans-consolidation-YYYYMMDD-HHMMSS` exists on origin.
- [ ] **Spec & Concept Protection:** I have manually verified that NONE of the merged files lost critical architectural concepts, domain models, contracts, or non-negotiable rules.
- [ ] **Aggressive File Reduction:** Total plan and subtask file count has been reduced significantly (combining 2, 3, or more related tasks into single files).
- [ ] **Zero Pure-Guideline Housekeeping Noise:** All historical micro-tasks solely addressing routine coding guidelines (whitespace, formatting, boolean conventions with zero business logic) were pruned.
- [ ] **Checklist Compaction:** Common checklists are consolidated into a single unified reference to `.lovable/coding-guidelines.md` per milestone file, avoiding multi-file duplication.
- [ ] **Spec Audit Reports Purged:** Verified `spec/25-app-spec-audit/` reports were cleaned using `python 03-ai-scripts/32-git-history-file-tracer.py --preset-audit --delete`.
- [ ] **Subtask Folder Collapse:** Fragmented subtask directories have been folded into single files or inlined into milestone summaries.
- [ ] **Strict Relative Git Paths:** All file paths, markdown links, citations, and subtask references in consolidated files are strictly relative to the git repository root. Zero absolute paths (`/absolute/path/to/...`) or `file:///` URIs.
- [ ] **Strict Lowercase Naming:** Every file in `.lovable/plans/completed/` uses strictly lowercase letters (e.g. `01-auth-system.md`).
- [ ] **Monotonic Sequencing:** File prefixes in `.lovable/plans/completed/` are continuous and monotonic (`01-`, `02-`, `03-`, ...) without gaps or duplicates.
- [ ] **Index Synchronization:** Both `.lovable/plans/01-index.md` and `.lovable/what-to-read.md` reflect the consolidated files and remove deleted entries.
- [ ] **LF Line Endings (`\n`):** All files use Unix LF line endings. Zero CRLF (`\r\n`).
- [ ] **UTF-8 Encoding (No BOM):** All files encoded in UTF-8 without BOM.
- [ ] **Single Trailing Newline:** Every file ends with exactly one terminating newline (`\n`).
- [ ] **Markdown Heading Spacing:** Exactly one blank line before and after headings (no leading blank line on line 1).
- [ ] **Zero Double Blank Lines:** No `\n\n\n` anywhere in markdown.
- [ ] `python 03-ai-scripts/21-sequence-integrity-linter.py` and `python 03-ai-scripts/22-doc-path-linter.py` exited with code 0.
- [ ] Local CI runner `python 03-ai-scripts/06-cicd-local-runner.py --all` exited with code 0.

---

## Non-Negotiable Coding Guidelines Checklist (Auto-Reject on Violation)

/goal You MUST verify every item on this checklist before committing any code. If a subagent violated one of these rules, you must reject their work.

- [ ] Strict Relative Git Paths: All file paths, markdown links, citations, and subtask references in plans, specs, and memory logs are strictly relative to the git repository root. Zero absolute paths or `file:///` URIs.
- [ ] Master Guidelines: I have fully read and strictly enforced `spec/02-coding-guidelines/08-file-folder-naming/` and `.lovable/coding-guidelines.md`.
- [ ] Backup Integrity: Verified `backup/plans-consolidation-YYYYMMDD-HHMMSS` branch exists on origin before touching plans.
- [ ] Concept & Spec Preservation: Zero concept loss, zero truncation, zero placeholder stubs (`TODO`, `[N]`, `// ...`).
- [ ] Monotonic Sequence: Verified sequential `01-`, `02-`, `03-` numbering across `.lovable/plans/completed/`.
- [ ] Compaction Proven: Proved that multiple micro-tasks were combined into high-density files with measurable file count reduction.
- [ ] Spec Audit Cleanup: Verified `spec/25-app-spec-audit/` was purged via `32-git-history-file-tracer.py`.
- [ ] High-Signal Task Ledger: Verified that zero-business-logic guideline tasks were pruned from milestone summaries.

### Master Task Checklist (Atomic Numbered Steps)

1. [ ] /learn the section as a /goal [AI Fix Scripts Memory](#ai-fix-scripts-memory)

- [ ] Action Summary: I have output a detailed `- [x]` checklist summarizing exactly what I accomplished this turn to prove I did not hallucinate.

---

## STRICT AVOIDANCE: Never Disable CI/CD

> [!CAUTION]
> **NEVER disable any CI/CD checks, GitHub Actions, or validation workflows.**
> Strictly avoid commenting out, bypassing, or deleting CI/CD steps to force a pipeline to pass. Your job is to fix the underlying code so that the CI/CD pipeline passes legitimately. Disabling CI/CD is an auto-reject failure.

---

## MUST FOLLOW NON-NEGOTIABLE

Listen, past runs of these turns have been sloppy and careless: wrong step counts, partial task lists dumped into chat instead of files, plans and session summaries half-filled with "[N]" placeholders, folders skimmed, open ambiguities ignored, CI/CD issues and `plans/subtasks/` forgotten, user commands dropped, coding guidelines bypassed, detailed specs chopped and summarized into useless junk, uppercase README files left uncorrected, `.lovable/memory/` created by accident, `strictly-avoid.md` overwritten, and explicit user instructions softened after being told not to. Read the whole codebase, read every folder in `spec/` and `.lovable/`, confirm root `readme.md` is strictly lowercase, find the root cause in one sentence, capture commands, issues, and pending tasks without omitting a single item, write the spec files and memory files in the right paths, update every index in the same turn, sync `readme.md` with `what-to-read.md`, preserve detailed specs verbatim with zero truncation, run builds and full unit tests, group commits with clear messages, and push everything to git before ending. Going deep IS the job. If you are not going deep, you are not doing the job.
