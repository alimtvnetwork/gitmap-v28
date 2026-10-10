# Engineering Subtask Plan: Distributed Fleet Nodes, Repo Feature & AI Merge Orchestrator

**Subtask ID:** `02-nodes-and-merge-ai`  
**Parent Plan:** `gitmap-command-enhancements`  
**Run Number:** 82  
**Spec Reference:**  
- [02-component-spec.md](../../../../02-spec/21-app/gitmap-command-enhancements/02-component-spec.md)  
- [repo-feature.md](../../../../02-spec/21-app/gitmap-command-enhancements/repo-feature.md)  
- [01-architecture-spec.md](../../../../02-spec/21-app/gitmap-command-enhancements/01-architecture-spec.md)  
**Status:** `READY FOR IMPLEMENTATION`  

---

## 1. Targeted File Inventory

| Component Area | Primary Code Files | Test Files |
| :--- | :--- | :--- |
| **Fleet Nodes Router** | `cli/cmdnodes/nodes_cmd.go` | `cli/cmdnodes/nodes_cmd_test.go` |
| **Fleet Nodes Protocol** | `cli/cmdnodes/nodes_summary.go` | `cli/cmdnodes/nodes_summary_test.go` |
| **Repo Feature Resolver** | `cli/cmdresolver/repo_feature.go` | `cli/cmdresolver/repo_feature_test.go` |
| **AI Merge Orchestrator** | `cli/cmdmergeai/merge_ai_core.go` | `cli/cmdmergeai/merge_ai_test.go` |
| **Merge AI Manifest & Instructions** | `cli/cmdmergeai/merge_ai_manifest.go` | `cli/cmdmergeai/merge_ai_test.go` |

---

## 2. Step-by-Step Implementation Tasks

### Step 1: Repair CLI Nodes Argument Slicing in `cli/cmdnodes/nodes_cmd.go`
* **Objective:** Ensure flags and subcommand variants (`+pe`, `fspe`, `pipe-error-all`, `--json`, `--force-all`) are not lost during argument dispatch.
* **Exact Modifications:**
  1. Inspect `isNodesFullSummary(args)` in `RunUnifiedNodesCLI`.
  2. Stop naive slicing (`args[1:]` / `args[2:]`) that strips `fs+pe` or `fspe` tokens.
  3. Reconstruct or pass intact slice to `RunNodesFullSummary(args)`.
  4. Ensure `isNodesPipelineAll(args)` routes `RunNodesPipelineErrorsAll(args)` preserving `--json` and `--force-all`.
* **Acceptance Criteria:**
  - `gitmap nodes fs+pe` passes `+pe` flag into `RunNodesFullSummary` so pipeline diagnostics are activated.
  - `gitmap nodes fspe --json` sets both `withPE = true` and `isJSON = true`.
  - `gitmap nodes pe all --force-all` correctly forwards `--force-all` to remote node execution commands.

---

### Step 2: Implement 4-Phase Distributed Fleet Nodes Protocol in `cli/cmdnodes/nodes_summary.go`
* **Objective:** Prevent redundant SSH invocations across fleet nodes through local-precedence deduplication.
* **Exact Modifications:**
  1. **Phase 1 (Local Catalog):** Interrogate local workspace via `cmdpending.ResolveWorkspaceRepositories(cwd)` and populate case-folded `localURLMap` and `localSlugMap`.
  2. **Phase 2 (Remote Discovery):** Concurrently query registered SSH nodes using `currentSSHExecutor.Execute` running `gitmap list --json` (or `gitmap.exe list --json` on Windows). Unmarshal output into `[]RemoteRepoManifestItem`.
  3. **Phase 3 (Local Precedence Deduplication):**
     - For each remote repository, verify if URL, slug, or folder name exists in local maps.
     - If repository exists locally, exclude from remote query list.
     - If repository is unique to the remote node, retain in `discoveredRemotes[node].uniqueURLs`.
  4. **Phase 4 (Execution & Aggregation):**
     - Execute local evaluation first: `cmdsummary.RunFullSummary` or `cmdsummary.RunPipelineErrorsAll`.
     - Concurrently dispatch targeted summaries/checks to remote nodes hosting unique repositories.
     - Skip remote nodes where `len(uniqueURLs) == 0`.
     - If `--json` is supplied, serialize into unified `FleetNodesEnvelope` containing `attributes` and `data`.
* **Acceptance Criteria:**
  - In a cluster where 65 repos are shared between local and remote worker-1, and 2 repos exist solely on worker-1, remote execution on worker-1 queries strictly the 2 unique repos.
  - If all repos exist locally, zero remote execution calls are made.
  - Terminal output prints clear section headers with green checkmarks for completed checks.

---

### Step 3: Implement & Validate Universal Repo Feature Resolver in `cli/cmdresolver/repo_feature.go`
* **Objective:** Enable seamless resolution of any destination target passed to GitMap commands.
* **Exact Modifications:**
  1. Implement `ResolveRepoFeature(target string) (*ResolvedRepoFeature, error)`:
     - **Branch A (Remote URL):** Match HTTPS/SSH URLs. Query `store.OpenDefault().ListRepos()` to adopt existing local clone if present. If remote exists and is not local, execute clone to `<workspace>/<slug>`. If remote does not exist, provision repository via authenticated GitHub API and initialize local directory.
     - **Branch B (Local folder with `.git`):** Verify `.git` exists, inspect health, extract origin URL, adopt directory directly without modifying history.
     - **Branch C (Local folder without `.git`):** Execute `git init -b main`, query GitMap credential vault, create remote GitHub repository, link origin.
     - **Branch D (Bare slug):** Normalize slug, compute path `<workspace-root>/<slug>`, create directory, create GitHub repository via authenticated API, initialize git and link remote origin.
  2. Enforce zero interactive prompts; return structured `apperror.AppError` (`E9085`) on missing credentials or scopes.
  3. Enforce non-destructive file preservation: pre-existing files in target directories are never overwritten or deleted.
* **Acceptance Criteria:**
  - `ResolveRepoFeature("https://github.com/org/existing-local.git")` adopts local directory without re-cloning.
  - `ResolveRepoFeature("new-slug")` creates directory in workspace root, initializes git, and creates GitHub remote.
  - `ResolveRepoFeature("./existing-repo")` adopts directory preserving all uncommitted files.

---

### Step 4: Implement AI Merge Orchestrator Input Parsing & Chronological Sorting in `cli/cmdmergeai/merge_ai_core.go`
* **Objective:** Ingest source repositories from diverse input formats and sort them chronologically to determine the baseline repository.
* **Exact Modifications:**
  1. Implement `parseMergeAIArgs(args []string) (string, []string, error)`:
     - Handle `config.json` containing `{"destination": "...", "sources": [...]}`.
     - Handle `.txt` file containing line-delimited sources.
     - Handle comma-separated and space-separated argument tokens.
  2. Implement `prepareAndInspectSources(inputs []string, stagingBase string)`:
     - Inspect each source (clone with `--depth 50` into staging workspace if remote URL).
     - Extract `HEAD` commit date (`git log -1 --format=%ct`), commit message, active branch, and tracked file list (`git ls-files`).
  3. Sort sources chronologically from earliest to latest commit timestamp. The earliest source becomes $S_1$ (base repository).
* **Acceptance Criteria:**
  - `gitmap merge-ai <dest> file.txt` parses URLs ignoring comment lines and blank lines.
  - `gitmap merge-ai <dest> url1, url2, url3` correctly tokenizes into 3 discrete sources.
  - Earliest committed repository is consistently assigned sequence order 1 ($S_1$).

---

### Step 5: Implement Single-Commit Staging & `01_`/`02_` Collision Sequencer in `cli/cmdmergeai/merge_ai_core.go`
* **Objective:** Assemble multi-repo file trees into a unified staging tree for a single commit without file loss.
* **Exact Modifications:**
  1. Enforce **Single-Commit Staging Invariant**: No synthetic branch rebasing or stacking.
  2. Implement `stageAndSequenceFiles(destDir string, sources []SourceRepoConfig)`:
     - **Base Repo ($S_1$):** Copy all tracked files directly into `<destDir>/<relPath>`.
     - **Subsequent Repos ($S_i$, $i \ge 2$):**
       - If `<destDir>/<relPath>` does not exist: copy directly into destination path.
       - If `<destDir>/<relPath>` exists (first collision):
         - Rename existing file to `01_<filename>`.
         - Copy incoming file to `%02d_<filename>` (e.g. `02_<filename>`).
         - Instantiate entry in `collisionMap`.
       - If `<destDir>/<relPath>` already collided: copy incoming file as `%02d_<filename>` (e.g. `03_<filename>`).
* **Acceptance Criteria:**
  - Non-colliding files from all sources are placed directly into destination tree.
  - Colliding files are renamed with sequential prefixes (`01_main.go`, `02_main.go`) preserving relative subdirectories.
  - Zero files are overwritten or silently discarded.

---

### Step 6: Implement Manifest & Instruction Generator in `cli/cmdmergeai/merge_ai_manifest.go`
* **Objective:** Produce root `merge-ai-manifest.json` and `instruction.md` for AI agent consolidation.
* **Exact Modifications:**
  1. Implement `WriteMergeAIManifest(destDir string, manifest MergeAIManifest)`:
     - Serializes destination info, `repoSequence` array (order, URL, branch, commit range, description, release info, folder tree), `fileCollisions` array, and `uniqueFilesDirectlyPlaced`.
  2. Implement `WriteMergeAIInstruction(destDir string, collisions []FileCollisionRecord)`:
     - Generates Markdown checklist detailing all detected collisions, instructions to preserve business logic, synthesize unified implementations, flatten sequence prefixes back to canonical file names, and verify build.
* **Acceptance Criteria:**
  - Destination root contains valid, indented `merge-ai-manifest.json`.
  - Destination root contains actionable `instruction.md` detailing every collision variant and exact consolidation checklist.

---

### Step 7: Unit & Integration Test Suites
* **Objective:** Ensure automated regression coverage across all modified packages.
* **Target Tests:**
  1. `cli/cmdnodes/nodes_cmd_test.go`:
     - Test argument preservation for `fs+pe`, `fspe`, `full-summary+pe`, `pipe-error-all`.
  2. `cli/cmdnodes/nodes_summary_test.go`:
     - Test local precedence deduplication logic with mock fleet connections.
  3. `cli/cmdresolver/repo_feature_test.go`:
     - Test classification of Branch A (URL), Branch B (local git), Branch C (unversioned folder), Branch D (bare slug).
  4. `cli/cmdmergeai/merge_ai_test.go`:
     - Test argument parsing (`config.json`, `.txt`, comma-separated).
     - Test chronological sorting of source repos.
     - Test collision sequencing (`01_`, `02_`) on synthetic test directories.
     - Test manifest and instruction generation.
* **Acceptance Criteria:**
  - All test suites pass cleanly without test regressions.

---

## 3. Verification & Acceptance Checklist

- [ ] **Fleet Nodes CLI Router:**
  - Verify `gitmap nodes fs+pe` does NOT drop `+pe`.
  - Verify `gitmap nodes fspe --json` emits JSON with pipeline error diagnostics.
  - Verify `gitmap nodes pe all` executes without hanging.
- [ ] **Fleet Deduplication Protocol:**
  - Verify Phase 1 local cataloging registers all local repos.
  - Verify Phase 2 discovery queries remote nodes concurrently.
  - Verify Phase 3 skips remote SSH calls for repositories already present locally.
  - Verify Phase 4 outputs aggregated tree view with node headers.
- [ ] **Universal Repo Feature Resolver:**
  - Test remote URL: clones if missing, adopts if local.
  - Test folder with `.git`: adopts directly.
  - Test folder without `.git`: initializes git and links origin.
  - Test bare slug: creates folder in workspace root, initializes git, and links origin.
  - Verify zero interactive prompts.
- [ ] **AI Merge Orchestrator:**
  - Test `gitmap merge-ai <dest> <src1> <src2>` with overlapping files.
  - Verify base repo files placed directly.
  - Verify conflicting files sequenced to `01_<filename>` and `02_<filename>`.
  - Verify root `merge-ai-manifest.json` conforms to schema.
  - Verify root `instruction.md` contains step-by-step AI checklist.
  - Verify single-commit staging invariant.
