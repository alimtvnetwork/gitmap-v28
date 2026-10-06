# Subtask 05: Synchronize Registries, Execute Quality Gates, and Perform Final Release v6.499.0

- **Parent Plan:** [235-app-spec-and-completed-plans-consolidation-and-reduction.md](../../235-app-spec-and-completed-plans-consolidation-and-reduction.md)
- **Specification:** [02-component-and-cli-spec.md](../../../02-spec/21-app/235-app-spec-and-completed-plans-consolidation-and-reduction/02-component-and-cli-spec.md)
- **Subtask ID:** `05`
- **Status:** `COMPLETED`
- **Assigned Worker:** Lead Orchestrator
- **Target Scope:** Documentation registries, linting scripts, release manifests

---

## 1. Context & Objective

Following the execution of Subtask 03 (spec folder 21 compaction) and Subtask 04 (completed plans and memory consolidation), the repository navigation indices must be updated to reflect the new canonical structure, eliminate dead links, verify repository health through rigorous quality gates, and perform the final release ceremony bumping the version to `v6.499.0`.

The objectives of Subtask 05 are to:
1. Rebuild and synchronize `02-spec/21-app/readme.md`, `.ai-memory/plans/readme.md`, and `.ai-memory/what-to-read.md`.
2. Execute quality verification protocols (`check-relative-paths.py`, `check-forbidden-strings.py`, `doc-path-linter.py`, and `go vet`).
3. Execute the post-consolidation final release ceremony bumping the version from `v6.498.0` to `v6.499.0`.
4. Publish comprehensive release notes and a showcase metrics scorecard.

---

## 2. Step-by-Step Execution Plan

### Step 5.1: Navigational Registry Updates

#### 5.1.1 Update `02-spec/21-app/readme.md`
Worker 01 will rewrite `02-spec/21-app/readme.md` to:
1. Index the eight canonical cluster directories:
   - `01-cli-architecture`
   - `02-scanner-and-projects`
   - `03-git-operations-and-pull`
   - `04-fleet-nodes-and-ssh`
   - `05-antigravity-and-ide`
   - `06-database-and-split-db`
   - `07-pipeline-and-diagnostics`
   - `08-distribution-and-release`
2. Remove all deleted Category A file links (the 21 `refactor-*.md` files, consistency reports, and stub overviews).
3. Provide a legacy-to-cluster mapping reference table for quick symbol lookup.

#### 5.1.2 Update `.ai-memory/plans/readme.md`
Worker 01 will update `.ai-memory/plans/readme.md` to:
1. Replace the long list of individual completed plans with links to Milestones 01 through 42 in `.ai-memory/plans/completed/`.
2. Retain Active Plans and Pending Plans sections with 100% fidelity.
3. Update the Subtasks Directory section to list only surviving active subtask directories (Plans 219, 226, 235, etc.), pruning folded subtask directories.

#### 5.1.3 Update `.ai-memory/what-to-read.md`
Worker 01 will update `.ai-memory/what-to-read.md` to:
1. Record a new Changelog entry for Plan 235 consolidation.
2. Update the "Before any task" and "Before writing code" reading lists to point to the consolidated memory reference summaries and cluster specifications.

---

### Step 5.2: Quality Verification Protocols

Worker 02 will execute four mandatory quality verification gates:

#### Gate 1: Relative Path Integrity Linter
- **Verification:** Run `python3 linter-scripts/check-relative-paths.py`
- **Criteria:** 0 absolute filesystem paths allowed anywhere in modified files or documentation.

#### Gate 2: Forbidden Strings & Secrets Linter
- **Verification:** Run `python3 linter-scripts/check-forbidden-strings.py`
- **Criteria:** 0 forbidden strings, 0 negative boolean patterns, and 0 leaked credentials.

#### Gate 3: Markdown Document Path Integrity
- **Verification:** Run `python3 03-ai-scripts/22-doc-path-linter.py`
- **Criteria:** 0 broken markdown links across `02-spec/`, `.ai-memory/`, and `01-prompts/`.

#### Gate 4: Go Static Analysis Gate
- **Verification:** Run `go vet ./...`
- **Criteria:** Clean compiler and static analysis report across all Go packages in `cli/`.

---

### Step 5.3: Final Release Ceremony (v6.499.0)

Upon passing all quality gates, the Lead Orchestrator executes the release ceremony:
1. Bump minor version to `v6.499.0` via `python3 03-ai-scripts/37-bump-version.py --minor` (or release orchestrator).
2. Validate that `version.json` records `6.499.0` as the authoritative version.
3. Update root `readme.md` version badges and release documentation headers.
4. Generate release notes detailing the consolidation metrics, preserved architectural invariants, and safety backups.
5. Create annotated release tag `v6.499.0` and publish to release branch `release/v6.499.0`.

---

### Step 5.4: Final Showcase Scorecard

The final task report will generate the compaction scorecard comparing baseline `v6.498.0` with release `v6.499.0`:

| Dimension | Baseline (`v6.498.0`) | Final (`v6.499.0`) | Reduction | Reduction % |
| :--- | :--- | :--- | :--- | :--- |
| `02-spec/21-app/` Files & Folders | ~120 items | ~35-45 items | ~75 items | **~62.5%** |
| `.ai-memory/plans/completed/` | 190 files | 42 milestone files | 148 files | **77.9%** |
| `.ai-memory/plans/subtasks/` | 229 files (42 dirs) | 52 files (10 dirs) | 177 files | **77.3%** |
| `.ai-memory/memory/` Files | ~85 files | ~35 files | ~50 files | **~58.8%** |
| **Total Artifact Count** | **~624 items** | **~174 items** | **~450 items** | **~72.1% net reduction** |

---

## 3. Acceptance Criteria & Quality Gates

- [ ] `02-spec/21-app/readme.md` synchronized with eight canonical clusters.
- [ ] `.ai-memory/plans/readme.md` updated with Milestones 01–42.
- [ ] `.ai-memory/what-to-read.md` updated with Plan 235 changelog and pointers.
- [ ] Relative paths verified (zero absolute paths).
- [ ] Forbidden strings check passed.
- [ ] Markdown link linter passed (zero broken links).
- [ ] Go vet static analysis passed.
- [ ] Final release `v6.499.0` completed and recorded.
- [ ] Final showcase metrics scorecard published.
