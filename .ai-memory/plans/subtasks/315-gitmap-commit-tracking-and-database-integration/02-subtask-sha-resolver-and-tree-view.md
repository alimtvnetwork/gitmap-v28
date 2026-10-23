# Subtask 02: Universal SHA Resolution, Multi-Database Telemetry & ASCII Tree View

## Status: READY
- **Task ID:** `315-subtask-02`
- **Parent Spec:** [02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/01-architecture-spec.md](../../../../02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/01-architecture-spec.md)
- **Module:** `cli/cmdcommit`, `cli/store`, `cli/repodb`
- **Owned Files:**
  - `cli/cmdcommit/commit_info_resolve.go`
  - `cli/cmdcommit/commit_info_query.go`
  - `cli/cmdcommit/commit_info_tree.go`
  - `cli/cmdcommit/commit_info_attribution.go`
  - `cli/store/commit_history_store.go`

---

## 1. Objectives

1. **Implement Universal 4-Stage SHA Resolution in `cli/cmdcommit/commit_info_resolve.go`:**
   - Resolve any input reference (40-char full SHA, 7-char short SHA, prefix, suffix, branch name, tag, or arbitrary substring) to an authoritative 40-character commit SHA.
   - Stage 1: Fast-path execution of `git rev-parse --verify <query>^{commit}` in the local repository.
   - Stage 2: Indexed query against Commit Split-DB (`repodb/commit_history.db`) matching `CommitSha`, `ShortSha`, prefix, and suffix.
   - Stage 3: Heatmap telemetry query against `HeatmapCommitLog` across discovered local and shared heatmap databases.
   - Stage 4: Streaming substring scan over `git rev-list --all` with bounded buffer and early exit.

2. **Implement Cross-Database Metadata Aggregation in `cli/cmdcommit/commit_info_query.go`:**
   - Query native Git commit metadata (`git log -1` and `git show --numstat`) to retrieve author, email, timestamp, message, parent SHAs, and file mutation deltas.
   - Query Commit Split-DB (`repodb/commit_history.db`) for indexed `CommitHistoryRecord`.
   - Query Heatmap databases for `HeatmapCommitLog`, `HeatmapAiSession`, and action log counts.
   - Query Multi-Agent databases (`ai_agents.db` and `agent-task.db`) for matching task plans, subtasks, and agent roles.

3. **Implement Origin Attribution Classifier in `cli/cmdcommit/commit_info_attribution.go`:**
   - Determine whether a commit was authored by:
     - Local Machine AI Agent (`OriginLocalAI`) when `IsAiCommitted` is true and `MachineHostname` matches the local host.
     - Remote Cluster Worker AI (`OriginForeignAI`) when `IsAiCommitted` is true and `MachineHostname` belongs to another machine.
     - External or Human Developer (`OriginHumanOrExternal`) when no AI session traces exist, emitting the standard status:
       `The current machine's AI has not committed this (likely another worker/node/human)`.

4. **Implement Rich ASCII Tree & JSON Formatter in `cli/cmdcommit/commit_info_tree.go`:**
   - Format hierarchical terminal output using standard box-drawing glyphs (`constants.TreeBranch`, `constants.TreeCorner`, `constants.TreePipe`, and `constants.TreeSpace`).
   - Colorize output with high-contrast palette (`ColorCyan`, `ColorNeonGreen`, `ColorNeonGold`, `ColorMagenta`, `ColorYellow`, `ColorDim`).
   - Support machine-readable JSON output mode (`--json`).

5. **Implement Split-DB Store Driver in `cli/store/commit_history_store.go`:**
   - Provide `OpenCommitHistoryDB(customPath string) (*CommitHistoryDB, error)`.
   - Implement `CommitHistoryRecord` persistence and indexed lookup methods.

---

## 2. Step-by-Step Implementation Work Breakdown

### Step 1: Implement Commit Split-DB Driver in `cli/store/commit_history_store.go`
- Define `CommitHistoryDB` struct wrapping `*DB` with mutex protection.
- Implement path resolution: `<effectiveDataDir>/repodb/commit_history.db`.
- Implement `InitSchema()` creating `CommitHistoryRecord` table and indexes.
- Implement `FindCommitBySha(repoSlug, query string) (*CommitHistoryRecord, error)`.
- Implement `UpsertCommitHistoryRecord(rec CommitHistoryRecord) error`.

### Step 2: Implement 4-Stage SHA Resolver in `cli/cmdcommit/commit_info_resolve.go`
- Implement `ResolveUniversalCommitSha(repoPath, query string) (string, error)`:
  - Call `resolveStage1GitRevParse(repoPath, query)`.
  - If resolved, return full SHA.
  - Call `resolveStage2SplitDb(repoPath, query)`.
  - If resolved, return full SHA.
  - Call `resolveStage3HeatmapLog(repoPath, query)`.
  - If resolved, return full SHA.
  - Call `resolveStage4GitRevListStream(repoPath, query)`.
  - If resolved, return full SHA.
  - Return `apperror.WrapSimple(ErrCommitNotFound, "universal sha resolution failed")`.
- Keep each stage resolver strictly $\le 15$ lines.

### Step 3: Implement Metadata Aggregator in `cli/cmdcommit/commit_info_query.go`
- Implement `BuildCommitInfoReport(repoPath, fullSha string) (*CommitInfoReport, error)`:
  - Call `queryGitCommitMetadata(repoPath, fullSha)`.
  - Call `queryGitFileDeltas(repoPath, fullSha)`.
  - Call `querySplitDbRecord(repoPath, fullSha)`.
  - Call `queryHeatmapTelemetry(repoPath, fullSha)`.
  - Call `queryAgentTaskDatabases(repoPath, fullSha)`.
  - Call `DetectCommitOrigin(report)`.
  - Assemble and return unified `*CommitInfoReport`.

### Step 4: Implement Origin Attribution in `cli/cmdcommit/commit_info_attribution.go`
- Implement `DetectCommitOrigin(rep *CommitInfoReport) CommitProvenance`:
  - Retrieve local host name via `os.Hostname()`.
  - If `!rep.CrossReferences.HasSplitDbRecord && !rep.CrossReferences.HasHeatmapRecord`:
    - Return `OriginHumanOrExternal` with message:
      `The current machine's AI has not committed this (likely another worker/node/human)`.
  - If commit has AI session records:
    - Compare `MachineHostname` to local hostname.
    - If matches: return `OriginLocalAI` with AI name and agent role.
    - If differs: return `OriginForeignAI` with worker hostname.

### Step 5: Implement Tree Visualizer & JSON Output in `cli/cmdcommit/commit_info_tree.go`
- Implement `RenderCommitInfoReport(w io.Writer, rep *CommitInfoReport, isJson bool) error`:
  - If `isJson`: call `renderCommitInfoJSON(w, rep)`.
  - Otherwise: call `renderCommitInfoTree(w, rep)`.
- Implement `renderCommitInfoTree(w io.Writer, rep *CommitInfoReport)`:
  - Render root commit header (Full SHA and short SHA).
  - Render metadata branch (`TreeBranch Author`, `TreeBranch Date`, `TreeBranch Message`).
  - Render provenance branch (`TreeBranch Provenance` with nested `TreePipe` + `TreeBranch`/`TreeCorner`).
  - Render database cross-references branch (`Split-DB`, `Heatmap`, `Agent Tasks`).
  - Render file mutations branch (`TreeBranch File Mutations` listing individual files with `+`/`-`).
  - Render summary footer (`TreeCorner Summary`).

---

## 3. Strict Coding Guidelines & Constraints

- **15-Line Function Limit:** Every function must be $\le 15$ lines (preferred $\le 8$ lines).
- **Positive Booleans Only:** Strict positive naming (`isResolved`, `isExactMatch`, `hasFoundRecord`, `isLocalHost`, `isBinaryFile`, `isLastEntry`).
- **Zero Nested Conditionals:** Guard clauses and early returns only; zero nested `if` statements.
- **Error Wrapping:** Wrap all underlying SQLite and subprocess errors with `*apperror.AppError` via `apperror.WrapSimple(err, "context")`.
- **Relative Path Hygiene:** Zero hardcoded absolute drive letters or URI schemes.
- **Zero Builds / Zero Tests:** Do not execute `go build` or `go test` during authoring.

---

## 4. Verification Evidence & Acceptance Criteria

1. Resolution succeeds for exact 40-character commit hashes via Stage 1.
2. Resolution succeeds for short 7-character commit hashes via Stage 1 or Stage 2.
3. Resolution succeeds for arbitrary substrings via Stage 3 or Stage 4.
4. Non-matching query returns `ErrCommitNotFound` cleanly without panic.
5. Commits without local AI traces render the standard text:
   `The current machine's AI has not committed this (likely another worker/node/human)`.
6. Terminal output accurately renders `constants.TreeBranch`, `TreeCorner`, `TreePipe`, and `TreeSpace` connectors.
7. `--json` output returns valid JSON conforming to `CommitInfoReport`.
8. All Go files satisfy line limits and boolean naming invariants.
