# Engineering Subtask Plan: Split-DB Summary Engine, Workspace Heatmap & Pipeline Telemetry

**Subtask ID:** `01-summary-engine`  
**Subtask Code:** `Task-01-Summary`  
**Parent Plan:** `gitmap-command-enhancements` ([gitmap-command-enhancements.md](../../gitmap-command-enhancements.md))  
**Run Number:** 82  
**Assigned Agent Role:** Spec Author 01 / Implementation Worker 01  
**Spec Reference:**  
- [01-architecture-spec.md](../../../../02-spec/21-app/gitmap-command-enhancements/01-architecture-spec.md)  
- [02-component-spec.md](../../../../02-spec/21-app/gitmap-command-enhancements/02-component-spec.md)  
- [repo-feature.md](../../../../02-spec/21-app/gitmap-command-enhancements/repo-feature.md)  
- [readme.md](../../../../02-spec/21-app/gitmap-command-enhancements/readme.md)  
- [Companion Subtask: Nodes & Merge-AI](02-nodes-and-merge-ai.md)  
**Status:** `READY FOR IMPLEMENTATION`  

---

## 1. Targeted File Inventory

| Component Area | Primary Implementation Files | Test & Documentation Files |
| :--- | :--- | :--- |
| **Split-DB Storage Engine** | `cli/cmdsummary/summary_cache.go` | `cli/cmdsummary/summary_test.go` |
| **Single-Repo Summary Core** | `cli/cmdsummary/summary_core.go`, `cli/cmdsummary/summary_types.go` | `cli/cmdsummary/summary_test.go` |
| **Workspace Heatmap & TreeView** | `cli/cmdsummary/full_summary_core.go`, `cli/cmdsummary/summary_helpers.go` | `cli/cmdsummary/summary_test.go` |
| **CI/CD Telemetry Fusion** | `cli/cmdsummary/full_summary_pe.go`, `cli/cmdsummary/summary_pe_all.go` | `cli/cmdsummary/summary_test.go` |
| **CLI Root Wiring & Aliases** | `cli/constants/constants_cli.go`, `cli/cmd/roottooling.go`, `cli/cmd/rootutility.go` | `cli/cmd/root_no_args_test.go` |

---

## 2. Disjoint File Ownership & Strict Boundaries

To guarantee multi-agent safety and eliminate concurrency conflicts:
- **Owned Subtask Files:**
  * `02-spec/21-app/gitmap-command-enhancements/01-architecture-spec.md`
  * `.ai-memory/plans/subtasks/gitmap-command-enhancements/01-summary-engine.md`
- **Downstream Owned Implementation Files (for Worker 01):**
  * `cli/cmdsummary/summary_types.go`
  * `cli/cmdsummary/summary_cache.go`
  * `cli/cmdsummary/summary_core.go`
  * `cli/cmdsummary/full_summary_core.go`
  * `cli/cmdsummary/full_summary_pe.go`
  * `cli/cmdsummary/summary_pe_all.go`
  * `cli/cmdsummary/summary_helpers.go`
  * `cli/cmdsummary/summary_test.go`
- **Strictly Prohibited Actions:**
  * **TOTAL BAN ON GIT COMMANDS:** Subagents MUST NOT execute `git add`, `git commit`, `git status`, `git push`, or `git diff`.
  * **Strictly Relative Paths Only:** All file references across documentation, code comments, and test fixtures must use forward-slash paths relative to repository root (`02-spec/...`, `.ai-memory/...`, `cli/...`). Zero absolute filesystem paths or `file:///` URIs.
  * **Search Exclusively via GitMap:** Use `gitmap aum search`, `gitmap find`, and `gitmap cat`. Strict ban on `grep`, `ripgrep`, `rg`, `Select-String`.
  * **No Build / No Test in Spec Phase:** Spec authoring requires zero build or test execution.

---

## 3. Step-by-Step Implementation Tasks

### Step 1: Implement Normalized SQLite Split-DB Cache Engine in `cli/cmdsummary/summary_cache.go`
* **Objective:** Establish high-speed, persistent caching of immutable release summaries, heated file churn metrics, and mutable repository head states to achieve $< 15\text{ms}$ query latency.
* **Exact Modifications:**
  1. Implement `resolveSummaryDBPath()`:
     - Check for local workspace directory: `filepath.Join(cwd, ".gitmap", "summary.db")`.
     - Fall back to global user profile: `filepath.Join(userHome, ".gitmap", "summary.db")`.
  2. Implement `OpenSummaryDB() (*sql.DB, error)` with robust connection pooling and SQLite pragmas:
     - `PRAGMA journal_mode = WAL;`
     - `PRAGMA busy_timeout = 5000;`
     - `PRAGMA synchronous = NORMAL;`
     - `PRAGMA foreign_keys = ON;`
     - Constrain connections: `SetMaxOpenConns(1)`, `SetMaxIdleConns(1)`, `SetConnMaxLifetime(10 * time.Minute)`.
  3. Implement `initSummaryDBSchema(db *sql.DB) error`:
     - Create tables: `Repositories`, `ReleaseSummaries`, `RepoHeadStates`.
     - Enforce `UNIQUE(repo_id, tag_name, tag_commit_hash)` on `ReleaseSummaries`.
     - Create view: `ViewRepoReleaseSummaries`.
  4. Implement `GetCachedRelease(repoURL, tag, tagCommitHash string) (*ReleaseSummaryRecord, bool)`:
     - Execute single-row index scan matching `(repo_url, tag_name, tag_commit_hash)`.
     - If matched, unmarshal `heated_files_json`, set `IsCacheHit = true`, and return in $< 15\text{ms}$.
  5. Implement `SaveCachedRelease(repoURL, slug, localPath string, rec ReleaseSummaryRecord) error`:
     - Upsert repository record in `Repositories` (`ON CONFLICT(repo_url) DO UPDATE`).
     - Upsert release record in `ReleaseSummaries` (`ON CONFLICT(repo_id, tag_name, tag_commit_hash) DO UPDATE`).
  6. Implement `UpdateRepoHeadState(repoURL, slug, localPath, headHash string, isDirty bool, dirtyCount int, lastActivity time.Time) error`:
     - Update mutable head commit hash, dirty flag, dirty files count, and activity epoch in `RepoHeadStates`.
* **Acceptance Criteria:**
  - Database schema initializes idempotently on cold start without errors.
  - Second invocation on an unchanged repository returns 100% cache hits in $< 15\text{ms}$.
  - Concurrent readers never receive table lock errors (`SQLITE_BUSY`) due to WAL mode and busy timeout.

---

### Step 2: Implement Single-Repo Summary Engine in `cli/cmdsummary/summary_core.go`
* **Objective:** Summarize the last $N$ releases (default $N = 8$) for a single repository, enforcing $\le 200$ words gist, file list suppression, and top 5 heated file churn metrics.
* **Exact Modifications:**
  1. Implement argument parser `parseSummaryArgs(args []string) (string, int)`:
     - Defaults target directory to `.`.
     - Defaults releases limit $N$ to `8`.
     - Extracts custom integer limit and target directory path or URL.
  2. Implement tag and commit extraction `extractRepoReleaseSummaries(dir, remoteURL string, limit int)`:
     - Query tags sorted chronologically: `git tag --sort=-creatordate`.
     - If tags are absent, invoke fallback `extractFallbackCommitSummaries(dir, limit)` using recent git commit logs.
     - For each tag, retrieve commit hash (`git rev-list -n 1 <tag>`).
     - Check cache via `GetCachedRelease`. On hit, append cached record and increment cache hit tally.
     - On cache miss, compute release delta: `computeReleaseSummary(dir, tag, prevTag, tagCommitHash)`.
  3. Implement semantic gist synthesizer `synthesizeReleaseGist(commitLog string) string`:
     - Parse commit subjects across the release boundary.
     - Construct high-level executive summary prose.
     - Enforce hard ceiling of $\le 200$ words (gracefully truncate at 180 words with `...`).
     - Strictly suppress raw file paths from prose gist.
  4. Implement heated files extractor `extractHeatedFilesFromDiff(numstat string) []HeatedFileMetric`:
     - Parse `git diff --numstat` lines into insertions, deletions, and relative path.
     - Calculate $\text{ChangesCount} = \text{Insertions} + \text{Deletions}$.
     - Sort descending by $\text{ChangesCount}$ and select top 5 files.
     - Assign functional rationale description via `summarizeFileChurn(path, ins, del)` based on extension (`.go`, `.ts`, `.md`, `.json`, etc.).
  5. Implement terminal renderer `renderSingleRepoSummaryTerminal(...)`:
     - Print repository header (name, URL, branch, short HEAD).
     - Print release cards with green checkmark (`✓`), release date, bold tag name, `(cached)` badge if applicable, summary gist, and indented heated files.
     - Print completion footer with elapsed seconds and cache hit ratio (`Cache: M/N hits`).
* **Acceptance Criteria:**
  - `gitmap summary 8` generates summaries for exactly 8 releases (or all available if fewer).
  - Word count of every generated release gist is $\le 200$ words.
  - Heated files section contains at most 5 files, reporting additions, deletions, and functional rationale.
  - Output displays green checkmarks (`✓`) and execution duration.

---

### Step 3: Implement Workspace Activity Heatmap TreeView in `cli/cmdsummary/full_summary_core.go`
* **Objective:** Scan workspace repositories, filter by the 48-hour activity window or dirty state, render an ANSI box-drawing tree view with $N=3$ releases, and provide a single master sanitize command footer.
* **Exact Modifications:**
  1. Implement options parser `parseFullSummaryOptions(args []string) FullSummaryOptions`:
     - Default `ReleasesCount = 3`, `ActivityHours = 48`.
     - Parse flags: `--json`, `+pe`, `fspe`, `--force-all`, `--verbose`.
  2. Implement repository activity evaluator `evaluateRepoActivity(rec, cutoff, opts)`:
     - Run `git status --porcelain` to extract untracked, modified, staged counts and pending file paths.
     - Inspect last commit epoch timestamp: `git log -1 --format=%ct`.
     - Calculate `isRecentlyActive = lastActivity.After(cutoff)`.
     - **48-Hour Filter Rule:** If `!opts.ForceAll && !isDirty && !isRecentlyActive`, return `nil` (suppress dormant clean repository).
     - If active or dirty, extract last $N=3$ release summaries using `summary.db` cache.
     - Generate per-repo suggested commit hint: `gitmap -C <dir> cpf "wip: save active progress"`.
     - If `opts.WithPE` is true, invoke `attachPipelineErrorTelemetry(record)`.
     - Update SQLite head state via `UpdateRepoHeadState`.
  3. Implement ANSI TreeView renderer `renderFullSummaryTreeView(payload, opts)`:
     - Render header with activity window hours and discovered active vs. total repository counts.
     - Render parent nodes using `├── 📁` (or `└── 📁` for last repository) with repo name, URL, and status badges (`[CLEAN]` in green, `[DIRTY: N uncommitted]` in yellow).
     - Render indented children (`│   ├──`, `│   └──`):
       * Pending changes list (if dirty).
       * Releases subtree with green checkmarks (`✓`), tags, release dates, and gists.
       * CI/CD status (if `+pe` enabled).
       * Suggested commit hint (if dirty).
  4. Implement Master Sanitize Command Footer:
     - If `payload.Data.DirtyReposCount > 0`, render prominent bounding box containing:
       `gitmap sanitize-all --message "wip: save active progress across dirty repos"`
  5. Implement JSON output serialization when `--json` flag is provided:
     - Emit typed JSON envelope conforming to `attributes` and `data.repositories` schema.
* **Acceptance Criteria:**
  - Repositories with no activity in 48 hours and a clean tree are excluded from output.
  - Repositories with uncommitted files are always included, displaying dirty file count and pending paths.
  - Passing repositories display `[CLEAN]` in green; modified repositories display `[DIRTY]` in yellow.
  - When dirty repositories exist, the master sanitize command banner is rendered at the bottom.
  - Supplying `--json` emits valid JSON with zero ANSI formatting codes.

---

### Step 4: Implement CI/CD Pipeline Telemetry Fusion in `cli/cmdsummary/full_summary_pe.go` & `summary_pe_all.go`
* **Objective:** Interrogate repository Split-DB pipeline stores (`repodb/pipeline.db`), enforce the Concise Green Rule, and extract bounded 25-line stack traces with line numbers.
* **Exact Modifications:**
  1. Implement telemetry attachment `attachPipelineErrorTelemetry(record *RepoSummaryRecord)`:
     - Resolve pipeline DB path via `pipelinedb.ResolvePipelineDbPath(record.CanonicalSlug)`.
     - Open pipeline Split-DB via `pipelinedb.OpenPipelineSplitDb(slug)`.
     - Query latest run: `db.QueryRunByNegativeOffset(-1)`.
     - **Concise Green Rule:** If run succeeded, set `record.IsPipelineClean = true` and return without reading log traces.
     - **Failing Run Extraction:** If run failed, extract `WorkflowName`, `JobName`, `StepName`, `ExitCode`, `RunUrl`.
     - Retrieve compact error logs: `db.QueryCompactErrorLogsByRunId(run.RunId)`.
     - **Bounded 25-Line Ceiling:** Truncate trace to the last 25 lines:
       ```go
       if len(traceLines) > 25 {
           traceLines = traceLines[len(traceLines)-25:]
       }
       ```
     - Assemble `RepoPipelineError` model and attach to record.
  2. Implement workspace-wide aggregator `RunPipelineErrorsAll(args []string)`:
     - Parse `--json` and `--force-all` flags.
     - Resolve workspace repositories; filter by 48-hour activity window unless `--force-all` is set.
     - Concurrently inspect pipeline status for each repository.
     - Tally `cleanCount` and `failedCount`.
     - Render terminal diagnostics: passing repositories output a single concise green line:
       `✓ [<repoSlug>] green (all checks passed)`
     - Failing repositories expand with workflow name, step name, exit code, and 25-line stack trace.
     - Conclude with summary tallies.
* **Acceptance Criteria:**
  - Passing CI/CD pipelines produce exactly one green indicator line with zero verbose log dumps.
  - Failing CI/CD pipelines display failure summary and at most 25 lines of error stack traces.
  - Stack traces preserve file names and line numbers for automated AI diagnosis.
  - `gitmap pe all` runs within 48h active filter by default and respects `--force-all`.

---

### Step 5: Implement CLI Root Dispatch Wiring & Constants
* **Objective:** Register all summary and telemetry command variants, short aliases, and multi-word rewrites in the root CLI router.
* **Exact Modifications:**
  1. In `cli/constants/constants_cli.go`:
     - Define `CmdSummary = "summary"`.
     - Define `CmdFullSummary = "full-summary"`, `CmdFullStatus = "full-status"`, `CmdFs = "fs"`.
     - Define `CmdFullSummaryPE = "full-summary+pe"`, `CmdFullStatusPE = "full-status+pe"`, `CmdFsPE = "fs+pe"`, `CmdFspe = "fspe"`.
     - Define `CmdNewCommands = "new-commands"`, `CmdNewCommandsAlias = "nc"`.
  2. In `cli/cmd/roottooling.go`:
     - Register `summary` routing to `cmdsummary.RunSummary(argsTail())`.
     - Register `full-summary`, `full-status`, `fs` routing to `cmdsummary.RunFullSummary(argsTail())`.
     - Register `full-summary+pe`, `full-status+pe`, `fs+pe`, `fspe` routing to `cmdsummary.RunFullSummary(append([]string{"+pe"}, argsTail()...))`.
  3. In `cli/cmd/rootutility.go`:
     - Intercept `all` argument under `pe`: `pe all` routes to `cmdsummary.RunPipelineErrorsAll(args[1:])`.
     - Register `pipe-error-all`, `pipe-errors-all` routing to `cmdsummary.RunPipelineErrorsAll(argsTail())`.
* **Acceptance Criteria:**
  - Executing `gitmap fs` triggers full summary without unknown command errors.
  - Executing `gitmap fspe` triggers full summary with pipeline error telemetry.
  - Executing `gitmap pe all` triggers workspace-wide pipeline error diagnostics.

---

### Step 6: Implement Comprehensive Unit & Benchmark Tests in `cli/cmdsummary/summary_test.go`
* **Objective:** Ensure 100% test pass rate across argument parsing, gist word limits, heated file churn analysis, Split-DB round-trip caching, and sub-15ms cache retrieval latency.
* **Exact Modifications:**
  1. `TestParseSummaryArgs`: Validate default directory (`.`) and limit ($N=8$), custom targets, and custom limits.
  2. `TestParseFullSummaryOptions`: Validate flag parsing (`--json`, `+pe`, `--force-all`, custom limits).
  3. `TestSynthesizeReleaseGist`: Validate semantic distillation and strict word count ceiling ($\le 200$ words).
  4. `TestExtractHeatedFilesFromDiff`: Validate diffstat parsing, change counts calculation, and top 5 selection.
  5. `TestSplitDBCacheRoundTrip`: Validate inserting and retrieving `ReleaseSummaryRecord` into `summary.db`, asserting `IsCacheHit == true`.
  6. `TestPipelineErrorTruncation`: Validate that stack traces exceeding 25 lines are strictly truncated to the last 25 lines.
  7. Benchmark `BenchmarkSplitDBCacheHit`: Assert cache hit query execution completes in $< 15\text{ms}$.
* **Acceptance Criteria:**
  - All unit tests pass with zero failures.
  - Cache hit retrieval benchmark verifies latency $< 15\text{ms}$.

---

## 4. End-to-End Test Suite & Verification Matrix

| Test ID | Command Under Test | Key Validation Assertions | Success Evidence |
| :--- | :--- | :--- | :--- |
| **TEST-01** | `gitmap summary 8` | Returns up to 8 release summaries; word count $\le 200$; top 5 heated files surfaced. | Output contains `✓ Release`, heated files list with churn counts. |
| **TEST-02** | `gitmap summary 8` (2nd run) | Split-DB cache hit; zero git subprocess invocations; latency $< 15\text{ms}$. | Header reports `(cached)`, latency footer reports `Cache: 8/8 hits`. |
| **TEST-03** | `gitmap fs` | Filters repositories by 48-hour activity window or dirty state. Dormant clean repos omitted. | TreeView rendered; active/dirty repos displayed. |
| **TEST-04** | `gitmap fs` (dirty workspace) | Identifies dirty repos, lists uncommitted files, renders suggested commit hints and master sanitize command. | Footer contains `gitmap sanitize-all --message ...`. |
| **TEST-05** | `gitmap fs+pe` / `fspe` | Passing repos render `✓ CI/CD Pipeline PASSING`; failing repos expand with 25-line stack trace. | Failure box contains workflow, job, step, exit code, bounded trace. |
| **TEST-06** | `gitmap pe all` | Fleet diagnostics respecting 48h active filter. Passing builds render concise green line. | Outputs `✓ [<slug>] green (all checks passed)` for clean repos. |
| **TEST-07** | `gitmap pe all --force-all` | Inspects every repository in the workspace regardless of activity window. | Scanned count matches total workspace repositories. |
| **TEST-08** | `gitmap fs --json` | Serializes complete payload to stdout matching JSON Envelope specification. | Valid JSON with `attributes` and `data.repositories`. |

---

## 5. Definition of Done & Quality Gates

- [ ] **Split-DB Persistence:** `summary.db` initialized in WAL mode with `Repositories`, `ReleaseSummaries`, `RepoHeadStates`, and `ViewRepoReleaseSummaries`.
- [ ] **Cache Latency SLA:** Second invocation of `gitmap summary` returns in $< 15\text{ms}$ with `(cached)` badge.
- [ ] **Gist Word Ceiling:** Release summaries strictly capped at $\le 200$ words with zero raw file lists in prose.
- [ ] **Heated Files Ceiling:** Churn analysis surfaces at most 5 files, reporting additions, deletions, and functional rationale.
- [ ] **48-Hour Activity Filter:** `gitmap fs` suppresses dormant clean repositories by default.
- [ ] **Master Sanitize Banner:** Prominently rendered in table footer when dirty repositories exist.
- [ ] **Concise Green Rule:** Passing pipelines output exactly one green line with zero verbose log dumps.
- [ ] **Bounded Stack Traces:** Failing pipelines truncated to exactly the last 25 lines with line numbers and file paths.
- [ ] **CLI Grammar & Aliases:** `summary`, `full-summary`, `full-status`, `fs`, `fspe`, `fs+pe`, `pe all`, and `pipe-error-all` fully routed.
- [ ] **Strict Relative Paths:** Zero absolute filesystem paths or `file:///` URIs across all spec, code, and doc files.
- [ ] **Zero Git Subprocesses:** No git commands executed by subagents.
