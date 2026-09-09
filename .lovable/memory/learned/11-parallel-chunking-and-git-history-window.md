# Learning 11: Parallel Chunking, 5-Second Heartbeats & Git History Filter Window

## Context & Motivation
Large repository verification of 2,500 to 6,700+ files often suffers from two bottlenecks:
1. Sequential single-threaded loops or per-file future dispatch causing Python GIL thread contention and high scheduling overhead.
2. Checking unchanged files when only a small subset of files changed in recent commits.

## Solution Architecture
1. **8-File Chunking Across Workers**:
   - Files are partitioned using `chunk_items(items, chunk_size=8)` into 8-file chunks.
   - Workers pick chunks (`[Worker-X] Picked chunk Y (8 files)...`) and scan them in-memory, avoiding thread submission per file.
2. **5-Second Snapshot Heartbeat**:
   - `WorkerHeartbeatMonitor` background daemon reports live throughput and worker activity every 5 seconds.
3. **Git History Changed Files Extractor (`27-git-changed-files.py`)**:
   - Extracts files changed in the last $N$ commits (`git diff --name-only HEAD~N..HEAD`) + working tree changes (`git status --porcelain`).
   - Deduplicates using a canonical relative path dictionary.
   - Exports to `.lovable/temp/git-changed-files.json`, `.yaml`, and `.txt` in < 0.1s.
4. **Linter Scoping**:
   - Linters (`check-relative-paths.py`, `check-nested-ifs.py`, `check-enum-and-boolean.py`) support `--changed-only`.
   - Running across 100-200 changed files takes < 0.2s instead of scanning the entire repo.
5. **Reentrant Lock in CI Runner**:
   - `DISK_WRITE_LOCK` in `03-ai-scripts/06-cicd-local-runner.py` must use `threading.RLock()` to prevent reentrant deadlocks when telemetry handlers trigger direct append writes.
