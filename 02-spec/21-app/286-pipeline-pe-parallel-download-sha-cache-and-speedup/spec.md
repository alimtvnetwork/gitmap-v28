# Specification: 286-pipeline-pe-parallel-download-sha-cache-and-speedup

## 1. Overview & Objective
Optimize `gitmap pe` (pipeline error logs) execution performance, eliminate sequential network bottlenecks, accelerate log line extraction, prevent stale TTL cache leaks across commit boundaries, and introduce a persistent commit SHA cache (`PipelineShaCache`) in `repodb/pipeline.db` for instant (<10ms) responses on repeated calls.

---

## 2. Architecture & Bottlenecks Addressed

### Bottleneck 1: Sequential Fallback Job Log Fetching
- **Current Behavior:** In `cli/cmdpipeline/pipeline_query_logs.go:124-146`, when top-level `--log-failed` fails or returns partial logs, `fetchJobLevelFailedLogs` iterates through failing jobs with `for _, j := range jobs` and calls `fetchFailedJobLogs(repo, j.DatabaseId)` one-by-one. In runs with 4 to 8 failing jobs, this incurs 15 to 60 seconds of sequential network blocking.
- **Remediation:** Refactor `fetchJobLevelFailedLogs` to download failing job logs concurrently using a bounded worker pool (concurrency = 4) and wait group. Maintain exact deterministic ordering of output logs.

### Bottleneck 2: Redundant Remote Metadata Probing on Cache Hits
- **Current Behavior:** In `cli/cmdpipeline/pipeline_logs_payload.go:167`, `enrichErrorLogsMetadata` unconditionally invokes `queryReleaseAndPRsConcurrently(repo)` which executes `gh release list` and `gh pr list`. Even when telemetry is retrieved from SQLite cache, this stalls execution for 1.5s to 3.0s.
- **Remediation:** Introduce a 10-minute in-memory/TTL cache for repository release versions and open PR counts in `cli/cmdpipeline/pipeline_logs_repo.go`. Bypass remote queries on cache hits when fresh metadata is already recorded.

### Bottleneck 3: Single-Threaded Regex & Line Processing
- **Current Behavior:** In `cli/cmdpipeline/pipeline_extract_clean.go:10`, every line unconditionally triggers `ansiRegex.ReplaceAllString(raw, "")`, even though >95% of log lines contain no ANSI escape sequences (`0x1b`).
- **Remediation:** Introduce `stripAnsiIfPresent(raw)` checking `strings.IndexByte(raw, 0x1b) != -1` before regex substitution. Accelerate line scanning and failure detection.

### Bottleneck 4: Stale TTL Cache Leak Across Commits
- **Current Behavior:** In `cli/cmdpipeline/pipeline_cache_eval.go:67-69`, `evaluateDecisionFromRuns` checks `checkCompletedCommitCacheHit` (which verifies `headSha == latest.Sha`). If `false` (e.g., after pushing a new commit), it falls through to `checkTtlCacheHit(db.Path)`. If the DB was modified within TTL (15 min), it returns `true` and incorrectly serves the previous commit's cached failures!
- **Remediation:** Enforce that `within_ttl` cache hit is strictly invalidated whenever local `headSha != ""` and `!matchCommitSha(headSha, latest.Sha)`.

### Bottleneck 5: Lack of Persistent SHA-Keyed Extracted Cache
- **Current Behavior:** `gitmap pe` repeatedly re-downloads or re-parses logs because `repodb/pipeline.db` only stores raw error rows, not the assembled structured payload.
- **Remediation:** Create `PipelineShaCache` table in `repodb/pipeline.db` (`CommitSha`, `RunId`, `RepoSlug`, `PayloadJSON`, `CreatedAt`, `UpdatedAt`). On completed runs, store the serialized payload. On invocation of `gitmap pe`, check `PipelineShaCache` for target commit SHA; if found, deserialize and return immediately in <10ms.

---

## 3. Database Schema Changes

### `PipelineShaCache` Table in `repodb/pipeline.db`
```sql
CREATE TABLE IF NOT EXISTS PipelineShaCache (
    PipelineShaCacheId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommitSha TEXT NOT NULL,
    RunId INTEGER NOT NULL,
    RepoSlug TEXT NOT NULL,
    PayloadJSON TEXT NOT NULL,
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP,
    UpdatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE UNIQUE INDEX IF NOT EXISTS IdxPipelineShaCache_Sha_RepoSlug ON PipelineShaCache (CommitSha, RepoSlug);
CREATE INDEX IF NOT EXISTS IdxPipelineShaCache_RunId ON PipelineShaCache (RunId);
```

---

## 4. File Ownership & Boundaries

### Worker 01 Files:
- `cli/cmdpipeline/pipeline_query_logs.go`
- `cli/cmdpipeline/pipeline_extract_scan.go`
- `cli/cmdpipeline/pipeline_extract_clean.go`

### Worker 02 Files:
- `cli/pipelinedb/pipeline_split_schema.go`
- `cli/pipelinedb/pipeline_split_schema_exec.go`
- `cli/pipelinedb/pipeline_split_queries.go`
- `cli/cmdpipeline/pipeline_sha_cache.go` (new)
- `cli/cmdpipeline/pipeline_logs_entry.go`
- `cli/cmdpipeline/pipeline_cache_eval.go`
- `cli/cmdpipeline/pipeline_logs_payload.go`
- `cli/cmdpipeline/pipeline_logs_repo.go`

---

## 5. Coding & Non-Negotiable Invariants
- All modified and newly created Go files must stay strictly <= 300 lines with nesting depth <= 1.
- No raw `git add`, `git commit`, `git push` by subagents.
- No heavy builds (`go build`, `npm run build`, `go test`, `pytest`).
- Zero absolute paths and zero `file:///` URIs.
