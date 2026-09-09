# Subtask 87.03: Fast Parallel Relative Paths Checker

## Objective
Upgrade `linter-scripts/check-relative-paths.py` with chunked parallel execution (8 files/chunk across 10+ workers), 5-second snapshot heartbeat, worker picked logs, and `--changed-only` Git commit history filtering.

## Proposed Changes
1. Arguments:
   - `--changed-only`, `-c`: Scan only files changed in last $N$ commits.
   - `--commits`, `-n`: Number of commits to check (default: 20).
   - `--workers`, `-w`: Worker thread count (default: 10 or os.cpu_count()).
   - `--chunk-size`: Files per chunk (default: 8).
   - `--all`, `-a`: Scan all tracked files (default if `--changed-only` not specified).
2. Dynamic Discovery & Changed-Files Resolution:
   - If `--changed-only` is passed, check if `.lovable/temp/git-changed-files.json` exists.
   - If missing or older than 60s, automatically invoke `03-ai-scripts/27-git-changed-files.py` to regenerate it.
   - Run dictionary deduplication on file list.
3. High-Speed Chunked Parallel Execution:
   - Chunk files into batches of 8.
   - Dispatch to `ThreadPoolExecutor(max_workers=workers)`.
   - Log picked chunk: `[Worker-X] Picked chunk (8 files)...` (in verbose mode or initial picks).
   - Periodic 5-second snapshot heartbeat printing: `[Snapshot 5s] Scanned X/Y files (Z%) | Throughput: N files/sec`.
4. Allowlist & Regex Matching:
   - Preserve existing `FORBIDDEN_PATTERNS` and `ALLOWLIST_FILES`.
   - Fast line scanning with memory safety and error resilience.

## Acceptance Criteria
- [ ] `python linter-scripts/check-relative-paths.py` runs parallelly and passes with 0 violations.
- [ ] `python linter-scripts/check-relative-paths.py --changed-only` finishes in < 0.5s.
- [ ] 5-second heartbeat telemetry works correctly.
- [ ] Coding guidelines met (<= 15 lines/func, zero nested ifs).
