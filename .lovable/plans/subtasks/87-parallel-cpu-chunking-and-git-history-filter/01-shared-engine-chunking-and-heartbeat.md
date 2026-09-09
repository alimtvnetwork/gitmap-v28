# Subtask 87.01: Shared Engine Chunking & Heartbeat Monitor

## Objective
Implement core reusable chunking, worker identifier logging, 5-second snapshot heartbeat monitor, and git changed files extraction in `03-ai-scripts/02-shared-engine.py`.

## Proposed Changes
1. `chunk_items(items: list[T], chunk_size: int = 8) -> list[list[T]]`:
   - Splits a flat list into chunks of `chunk_size` elements.
2. `class WorkerHeartbeatMonitor`:
   - Runs a background daemon thread reporting snapshots every 5.0 seconds.
   - Tracks active worker statuses, chunk indices, processed item counts, elapsed time, and items/sec.
   - Formats clean terminal output: `[Snapshot 5s] Processed 120/800 (15.0%) | 10 workers | 24.0 items/sec`.
   - Stops cleanly upon pool completion.
3. `run_chunked_worker_pool(...)`:
   - Takes list of items, batches them into 8-item chunks (configurable).
   - Submits chunks to `ThreadPoolExecutor(max_workers=workers)`.
   - Emits pick logs: `[Worker-X] Picked chunk Y (8 files)...`.
   - Supports 5-second snapshot heartbeat.
   - Gathers results and handles errors without unhandled exceptions.
4. `extract_git_changed_files(repo_root: Path, commit_count: int = 20) -> list[dict[str, Any]]`:
   - Runs `git diff --name-only HEAD~N..HEAD` and `git status --porcelain`.
   - Deduplicates files using canonical relative path dictionary lookup.
   - Returns sorted list of file records: `{"path": rel_path, "status": status, "extension": ext}`.

## Acceptance Criteria
- [ ] `chunk_items` partitions correctly for exact, remainder, and empty lists.
- [ ] `WorkerHeartbeatMonitor` triggers cleanly every 5.0 seconds without deadlocks.
- [ ] `run_chunked_worker_pool` aggregates results across workers.
- [ ] `extract_git_changed_files` produces deduplicated results.
- [ ] Conforms to coding guidelines (<= 15 lines/func, blank lines before return, zero nested ifs).
