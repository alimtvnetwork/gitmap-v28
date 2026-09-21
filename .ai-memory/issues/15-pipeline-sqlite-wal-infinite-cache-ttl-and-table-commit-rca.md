# Root Cause Analysis (RCA-15): SQLite WAL Connection-Open ModTime Invalidation and Table Inversion

## 1. Symptom
When executing `gitmap pipeline` (`gitmap pl`) or `gitmap pipeline-errors` (`gitmap pe`), gitmap persistently served cached pipeline run data from old commits (e.g., `76db95c`) for minutes or hours instead of respecting the 5-second cache TTL window (`GITMAP_PIPELINE_CACHE_TTL_SEC` / `pipelineCacheTtlSec`). The terminal output continuously reported:
```text
  Cache: ⚡ Served from local SQLITE DB cache (commit 76db95c)
```
Consequently:
1. `gitmap` did not query GitHub Actions for fresh runs.
2. Even after hours and repeated runs, fresh runs triggered by new commits (e.g., `7c1d6f9a`) were never fetched.
3. The recent commits pipeline summary table retained outdated commits and never showed the active running pipeline or updated statuses.

---

## 2. Root Cause
The root cause was that `resolveDbModTime` in `cli/cmdpipeline/pipeline_cache_eval.go` inspected the filesystem modification time of the SQLite Write-Ahead Log (`sql.db-wal`):
```go
func resolveDbModTime(dbPath string, defaultTime time.Time) time.Time {
    walPath := dbPath + "-wal"
    walInfo, err := os.Stat(walPath)
    if err == nil && walInfo.ModTime().After(defaultTime) {
        return walInfo.ModTime()
    }
    return defaultTime
}
```
Whenever `pipelinedb.OpenPipelineSplitDb` was called to evaluate the cache, modernc SQLite opened the connection and executed `PRAGMA journal_mode = WAL`. On opening a connection in WAL mode, SQLite touched or created `sql.db-wal` at the current instant (`time.Now()`).

When `checkTtlCacheHit(db.Path)` evaluated `time.Since(modTime) < resolvePipelineCacheTTL()`, `modTime` was the timestamp of `sql.db-wal` that was modified merely ~14 milliseconds earlier during the connection open call. Thus:
1. `time.Since(modTime)` was always approximately 10–20ms.
2. `time.Since(modTime) < 5*time.Second` was **always true** forever.
3. The cache hit condition never failed, permanently locking gitmap into serving whatever outdated runs existed in the local SQLite database.

---

## 3. Resolution
1. **Decoupled Cache Sync Metadata File:**
   Introduced a dedicated synchronization manifest `cache_sync.json` stored in the repository-scoped pipeline data folder (`.gitmap/data/pipeline/<repoSlug>/cache_sync.json`).
   - Fields: `fetchedAt` (RFC3339 UTC timestamp), `fetchedSha` (head SHA of primary run), and `runCount`.
   - Written exclusively upon querying fresh runs from GitHub Actions and persisting them via `RecordFetchedRunsToSplitDb` and `recordInPipelineSplitDb`.
2. **Deterministic Cache Evaluation:**
   Updated `checkTtlCacheHit(dbPath)` in `cli/cmdpipeline/pipeline_cache_eval.go`:
   - Inspects `cache_sync.json` directly.
   - If `cache_sync.json` is present, checks `time.Since(meta.FetchedAt) < resolvePipelineCacheTTL()`.
   - If older than 5s (or configured TTL), returns `false` (cache miss), forcing fresh fetch from GitHub Actions.
   - Completely eliminated `resolveDbModTime` and all inspections of `sql.db-wal`.
   - Fallback when `cache_sync.json` is absent checks `os.Stat(dbPath)` on `sql.db` alone.
3. **Quality Gates & Comprehensive Unit Testing:**
   Added tests in `cli/cmdpipeline/pipeline_cache_eval_test.go` verifying:
   - Hit within TTL with `cache_sync.json`.
   - Miss/expiration when TTL window has elapsed.
   - Respect for `GITMAP_PIPELINE_CACHE_TTL_SEC` env overrides.
   - Safe fallback when sync manifest is missing.

---

## 4. Prevention & Learnings
1. **Never Stat Ephemeral Database Internals:** Active database engines (SQLite WAL, shm, lockfiles) mutate auxiliary files during standard read transactions and connection handshakes. Application-level cache TTLs must never depend on database engine internal journal timestamps.
2. **Dedicated Decoupled Sync Manifests:** Always use dedicated metadata manifests (`cache_sync.json`) with explicit UTC timestamps to govern caching windows.
3. **Verify TTL Invalidation Live:** In cache evaluation test suites, always test both sides of the window: `time.Since(t) < TTL` (hit) and `time.Since(t) >= TTL` (miss).
