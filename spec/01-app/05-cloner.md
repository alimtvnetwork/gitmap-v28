# Cloner

## Responsibility

Read a structured file (CSV, JSON, or text) and re-clone repositories,
preserving the original folder hierarchy. Also supports cloning a single
repository directly from a Git URL.

## Behavior

### File-based clone

1. Detect file format by extension (`.csv`, `.json`, `.txt`).
2. Parse records from the file.
3. For each record:
   a. Create the relative directory structure under `--target-dir`,
      preserving the **full nested folder hierarchy** captured during
      scan (no flattening, no path rewriting). A record with
      `RelativePath = group-a/sub/repo-x` always lands at
      `<target>/group-a/sub/repo-x`.
   b. Run `git clone -b <branch> <url> <target-path>`.
4. Log success or failure for each clone operation.
5. Print a summary: N succeeded, M failed.

### Parallel execution (`--max-concurrency N`)

When `N > 1` the per-record clone work is dispatched onto a bounded
worker pool (`gitmap/cloner/concurrent.go`):

- The default of `N = 1` keeps the legacy sequential runner so progress
  lines stay strictly ordered for users who script around them.
- The pool reads from a buffered job channel and writes to a single
  result channel; the collector goroutine serialises Progress + cache
  + summary updates so callers do not need their own locking.
- `Progress` (`progress.go`) and `CloneCache` (`cache.go`) carry their
  own mutexes — concurrent writes cannot interleave a half-written
  stderr line or corrupt the cache map.
- The on-disk **nested hierarchy is preserved at any N** because every
  worker still resolves `filepath.Join(targetDir, rec.RelativePath)`
  for its destination. Worker count never changes layout — only timing
  and progress-line ordering.
- Cache hits short-circuit before being dispatched to workers so a
  fully-cached run does no real work and prints the same `skipped
  (cached)` lines whether N is 1 or 100.

### Direct URL clone

1. Detect that the source is a URL (`https://`, `http://`, `git@`).
2. Derive the repo name from the URL (or use a custom folder name).
3. Run `git clone <url> <folder>`.
4. Upsert the repo record into the database.
5. Prompt to register with GitHub Desktop.

## Audit Mode

When `--audit` is passed, the cloner runs read-only:

1. Parse the source manifest as it would for a real clone run.
2. For each record, compute the exact `git clone` (or `git pull`) command
   that would be invoked, using the same branch-selection strategy as
   the live path (`pickCloneStrategy`).
3. Stat the destination to classify each record as one of:
   - `clone` (`+`) — target path does not exist yet
   - `pull` (`~`) — target is a git repository
   - `cached` (`=`) — clone-cache fingerprint matches local HEAD
   - `conflict` (`?`) — target exists but is not a git repo
   - `invalid` (`!`) — record has no `HTTPSUrl` / `SSHUrl`
4. Print a diff-style report to stdout and exit 0. Never invoke git,
   never write outside stdout, never touch the network.

Audit refuses direct-URL invocations (it only operates on manifests).

## Error Handling

- If a clone fails (network, auth, etc.), log the error and continue.
- Do not abort the entire run for a single failure.
- Summary at end lists all failures with reasons.
- For direct URL clone, if the target folder exists, exit with error.

## Input Formats

| Format | Structure                              |
|--------|----------------------------------------|
| CSV    | Standard CSV with headers              |
| JSON   | Array of `ScanRecord` objects          |
| Text   | One `git clone ...` command per line   |
| URL    | Direct HTTPS or SSH git URL            |
