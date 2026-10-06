# Subtask 05: Synchronize Registries, Execute Quality Gates, and Perform Final Release v6.501.0

- **Parent Plan:** [236-deep-spec-consolidation-and-canonical-reduction.md](../../236-deep-spec-consolidation-and-canonical-reduction.md)
- **Specification:** [02-component-and-cli-spec.md](../../../02-spec/21-app/236-deep-spec-consolidation-and-canonical-reduction/02-component-and-cli-spec.md)
- **Subtask ID:** `05`
- **Status:** `COMPLETED`
- **Assigned Worker:** Lead Orchestrator
- **Target Scope:** `02-spec/21-app/readme.md`, `.ai-memory/plans/readme.md`, `.ai-memory/what-to-read.md`, quality linters, release manifests

---

## 1. Context & Objective

Following the execution of Subtask 03 (spec folder 21 deep compaction into 8 clusters) and Subtask 04 (completed plans compaction into 18 milestones and AI memory reduction to 10 files), the repository must undergo registry synchronization, quality verification gating, and the final minor version release ceremony for `v6.501.0`.

The primary objectives of Subtask 05 are to:
1. **Synchronize Navigational Registries:**
   - Rebuild `02-spec/21-app/readme.md` to index surviving specifications strictly within the 8 Canonical Domain Clusters, removing obsolete links and providing a legacy-to-cluster mapping table.
   - Update `.ai-memory/plans/readme.md` to index Milestones `01` through `18` in the Completed Plans section, preserving active and pending plans.
   - Update `.ai-memory/what-to-read.md` with Task 236 compaction milestones and updated developer orientation paths.
2. **Execute Mandatory Quality Verification Gates:**
   - Gate 1: Relative Path Integrity Linter (`check-relative-paths.py`).
   - Gate 2: Forbidden Strings & Secrets Linter (`check-forbidden-strings.py`).
   - Gate 3: Markdown Document Cross-Link Linter (`03-ai-scripts/22-doc-path-linter.py`).
   - Gate 4: Go Static Analysis Gate (`go vet ./...`).
3. **Execute Final Release Ceremony (v6.501.0):**
   - Execute minor version bump from `v6.500.0` to `v6.501.0`.
   - Update all version manifests and root documentation synchronously.
   - Tag `v6.501.0` and deploy release branch `release/v6.501.0`.
   - Publish the final before-and-after compaction scorecard.

---

## 2. Step-by-Step Operational Implementation Plan

### Step 5.1: Navigational Registry Updates

#### 5.1.1 Update `02-spec/21-app/readme.md`
Worker 01 will update `02-spec/21-app/readme.md`:
1. **Cluster Indexing:** Structure the catalog strictly around the 8 Canonical Domain Clusters:
   - `01-cli-architecture/`
   - `02-scanner-and-projects/`
   - `03-git-operations-and-pull/`
   - `04-fleet-nodes-and-ssh/`
   - `05-antigravity-and-ide/`
   - `06-database-and-split-db/`
   - `07-pipeline-and-diagnostics/`
   - `08-distribution-and-release/`
2. **Dead Link Pruning:** Ensure no links point to deleted Category A files or superseded directories.
3. **Legacy Mapping Matrix:** Provide a reference table mapping former spec slugs (e.g. `180`, `204`, `216`, `227`) to their corresponding canonical cluster locations.

#### 5.1.2 Update `.ai-memory/plans/readme.md`
Worker 01 will update `.ai-memory/plans/readme.md`:
1. **Completed Plans Section:** Replace the legacy 43 milestone entries with the 18 consolidated milestones:
   - `01-generation-1-core-foundation-milestones.md` (Milestones 01–15)
   - `02-generation-2-fleet-automation-and-modularization.md` (Milestones 16–27)
   - `03-agy-prompts-templates-and-rerun-suite.md` through `18-app-spec-and-completed-plans-consolidation-and-reduction.md`
2. **Active & Pending Sections:** Retain all active open plans and pending plans verbatim.
3. **Subtasks Directory Section:** Ensure only active WIP subtask directories are listed.

#### 5.1.3 Update `.ai-memory/what-to-read.md`
Worker 01 will update `.ai-memory/what-to-read.md`:
1. **Changelog Entry:** Append a record for Task 236 documenting the second-generation deep compaction metrics.
2. **Reading Guidance Paths:** Update the "Before any task" and "Before writing code" reading lists to point directly to:
   - `.ai-memory/memory/00-project-governance-and-invariants.md`
   - Canonical Domain Memory References in `.ai-memory/memory/`
   - Authoritative Cluster Specifications in `02-spec/21-app/`
   - Consolidated Milestones in `.ai-memory/plans/completed/`

---

### Step 5.2: Quality Verification Protocols

Worker 02 will execute four mandatory quality verification gates:

#### Gate 1: Relative Path Integrity Linter
- **Command:** `python3 linter-scripts/check-relative-paths.py`
- **Pass Requirement:** 0 absolute filesystem paths allowed anywhere in modified files, specifications, memory plans, or release notes.
- **Assertion:** `isPathRelative == true`.

#### Gate 2: Forbidden Strings & Secrets Linter
- **Command:** `python3 linter-scripts/check-forbidden-strings.py`
- **Pass Requirement:** 0 forbidden strings, 0 negative boolean patterns, and 0 leaked secrets or tokens.
- **Assertion:** `hasCleanTokens == true` and `isPositiveBooleanCompliant == true`.

#### Gate 3: Markdown Document Path Integrity
- **Command:** `python3 03-ai-scripts/22-doc-path-linter.py`
- **Pass Requirement:** 0 broken markdown links across `02-spec/`, `.ai-memory/`, and `01-prompts/`.
- **Assertion:** `hasValidDocLinks == true`.

#### Gate 4: Go Static Analysis Gate
- **Command:** `go vet ./...`
- **Pass Requirement:** Clean compiler and static analysis report across all Go packages in `cli/`.
- **Assertion:** `isGoVetClean == true`.

---

### Step 5.3: Final Release Ceremony (v6.501.0)

Upon successful verification of all quality gates, the Lead Orchestrator executes the release ceremony:
1. **Minor Version Bump:**
   - Execute minor version bump from `v6.500.0` to `v6.501.0` via `python3 03-ai-scripts/37-bump-version.py --minor` (or release orchestrator).
2. **Synchronous Manifest Verification:**
   - `version.json`: Authoritative root version updated to `6.501.0`.
   - `package.json`: Version field synchronized to `6.501.0`.
   - `.gitmap/release/latest.json`: Latest release metadata updated.
   - `readme.md`: Badges and version headers updated to `v6.501.0`.
   - `what-to-read.md`: Release notes and version references synchronized.
   - `cli/constants/constants.go`: Version constant updated to `6.501.0`.
   - `changelog.md`: Detailed entry for `v6.501.0` documenting the deep spec consolidation, completed plans compaction, and memory reduction.
3. **Git Release Artifacts:**
   - Annotated release tag `v6.501.0`.
   - Dedicated release branch `release/v6.501.0`.
   - Push release branch and tag to remote origin.

---

### Step 5.4: Final Compaction Showcase Scorecard

The final completion summary will publish the before-and-after compaction scorecard:

| Target Subsystem | Baseline Count (`v6.500.0`) | Consolidated Count (`v6.501.0`) | Net File Reduction | Net Reduction % |
| :--- | :--- | :--- | :--- | :--- |
| `02-spec/21-app/` Files & Dirs | ~45 items | ~16 items (8 clusters + docs) | ~29 items | **~64.4%** |
| `.ai-memory/plans/completed/` | 43 milestone files | 18 milestone files | 25 files | **58.1%** |
| `.ai-memory/plans/completed/` Lines | ~35,000 lines | ~12,000 lines | ~23,000 lines | **~65.7%** |
| `.ai-memory/memory/` Files | 71 files | 10 files | 61 files | **85.9%** |
| **Total Target Files** | **159 items** | **44 items** | **115 items** | **72.3% net reduction** |

---

## 3. Acceptance Criteria & Quality Gates

- [ ] `02-spec/21-app/readme.md` updated and indexed strictly to the 8 Canonical Domain Clusters.
- [ ] `.ai-memory/plans/readme.md` updated with Milestones `01` through `18` in Completed Plans.
- [ ] `.ai-memory/what-to-read.md` updated with Task 236 changelog and canonical reading paths.
- [ ] Gate 1: Relative Path Integrity Linter passed with 0 violations.
- [ ] Gate 2: Forbidden Strings & Secrets Linter passed with 0 violations.
- [ ] Gate 3: Markdown Document Path Integrity passed with 0 broken links.
- [ ] Gate 4: Go static analysis (`go vet ./...`) passed cleanly.
- [ ] Minor version bumped to `v6.501.0` across all manifests and documentation.
- [ ] Release branch `release/v6.501.0` and tag `v6.501.0` pushed to remote origin.
- [ ] Compaction scorecard published in final task report.
