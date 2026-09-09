# Plan 87: Parallel CPU Chunking, 5-Second Heartbeat Telemetry & Git History Window Filter

## Overview
Transform repository linters and checkers from slow, single-threaded or unchunked execution into high-speed parallel engines that process files in 8-file chunks across 10+ workers, report 5-second snapshot heartbeats, and filter candidate files against the last 10–20 Git commits using deduplicated JSON/YAML/TXT exports in `.lovable/temp/`.

---

## Key Requirements & Scope
1. **Shared Engine Upgrades (`03-ai-scripts/02-shared-engine.py`)**:
   - `chunk_items(items, chunk_size=8)`: Partition file list into manageable 8-item chunks.
   - `run_chunked_worker_pool(...)`: Parallel chunk runner with worker identifier logs (`[Worker-X] Picked chunk (8 files)...`).
   - `WorkerHeartbeatMonitor`: Periodic 5-second snapshot reporter logging active worker activities, completed counts, and throughput (files/sec).
   - `extract_git_changed_files(repo_root, commit_count=20)`: Core extractor collecting files changed in last $N$ commits plus staged/untracked changes, deduplicated via dictionary lookup.
2. **Git Commit History Changed-Files Extractor (`03-ai-scripts/27-git-changed-files.py`)**:
   - Standalone CLI utility with `--commits <N>` / `-n <N>` (default: 20), `--output-dir` (default: `.lovable/temp/`), `--format` (`all`, `json`, `yaml`, `txt`).
   - Exports:
     - `.lovable/temp/git-changed-files.json`
     - `.lovable/temp/git-changed-files.yaml`
     - `.lovable/temp/git-changed-files.txt`
   - Deduplication dictionary lookup to guarantee zero duplicates.
   - Fast execution (< 0.2s).
3. **High-Speed Chunked Relative Paths Checker (`linter-scripts/check-relative-paths.py`)**:
   - Upfront file listing and extension/directory pre-filtering.
   - Support `--changed-only` / `-c` and `--recent [N]` reading from `.lovable/temp/git-changed-files.json`.
   - Auto-generate changed files artifact if missing when `--changed-only` is requested.
   - Parallel chunking (8 files/chunk, 10 workers or `os.cpu_count()`).
   - 5-second snapshot heartbeat and chunk pick logs.
   - Deduplication dictionary lookup.
4. **Nested Ifs & Enum/Boolean Checkers Upgrades**:
   - Enhance `linter-scripts/check-nested-ifs.py` and `linter-scripts/check-enum-and-boolean.py` to use chunked parallel execution (8 files/chunk) instead of individual future dispatch.
   - Support `--changed-only` reading from `.lovable/temp/git-changed-files.json`.
   - Periodic 5-second snapshot heartbeat.
5. **CI/CD Local Runner Integration (`03-ai-scripts/06-cicd-local-runner.py`)**:
   - Support `--changed-only` / `--recent` flag to forward changed-file scoping to linters.
   - Full repository verification ensuring 100% green exit code.

---

## Subtasks Breakdown
- `01-shared-engine-chunking-and-heartbeat.md`: Add chunking, heartbeat monitor, and git changed files extraction to `02-shared-engine.py`.
- `02-git-changed-files-extractor.md`: Implement standalone CLI tool `03-ai-scripts/27-git-changed-files.py` with multi-format export.
- `03-parallel-relative-paths-checker.md`: Refactor `check-relative-paths.py` with 8-file chunking, 10 workers, 5s heartbeat, and `--changed-only` filter.
- `04-parallel-nested-ifs-and-enum-upgrades.md`: Upgrade `check-nested-ifs.py` and `check-enum-and-boolean.py` to chunked execution and changed-files filter.
- `05-cicd-runner-integration-and-verification.md`: Integrate `--changed-only` into CI local runner and verify all gates.

---

## Success Criteria
- [ ] `03-ai-scripts/27-git-changed-files.py` extracts changed files from last 20 commits into JSON, YAML, and TXT in `.lovable/temp/` in < 0.2s.
- [ ] All exported files are deduplicated with zero duplicate entries.
- [ ] `linter-scripts/check-relative-paths.py` runs with 8 files/chunk, logs worker picked chunks, prints 5-second snapshot heartbeats, and supports `--changed-only`.
- [ ] `check-nested-ifs.py` and `check-enum-and-boolean.py` run chunked parallel execution with 5-second snapshot heartbeats.
- [ ] All linters exit 0 with 0 violations.
- [ ] Zero breaking changes to existing CI pipelines.
