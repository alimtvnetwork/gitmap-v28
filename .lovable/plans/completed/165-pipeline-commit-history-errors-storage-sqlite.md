# Plan 165: Pipeline Commit-Based Errors, History, Logs, SQLite Telemetry & Storage Suite

## 1. Executive Summary & Accomplishments

This plan resolved the positional offset misalignment in `gitmap pipeline errors` where negative offsets (`-1`, `-2`, `-3`) previously targeted raw workflow run indexes instead of distinct Git commits, resulting in duplicate inspections of the same commit and false "PASSING (clean)" statuses on in-progress runs. It implemented true commit-scoped grouping, multi-workflow error aggregation, live in-progress state tracking, runaway newline suppression, new `gitmap pipeline history` and `gitmap pipeline logs` subcommands, repository-scoped SQLite database telemetry (`.gitmap/data/pipeline.db`), and complete `gitmap storage` and `gitmap os storage` command suites with CLI help text parity.

### Core Achievements
1. **Commit-Scoped Offset Resolution (`-1`, `-2`, `-3`)**:
   - Implemented `GroupRunsByCommit` grouping runs by `HeadSha`.
   - Translated negative offsets `-1`, `-2`, `-3` to target previous distinct commits across historical pipeline runs.
   - For each commit, aggregated all workflow runs (CI, Release, Lint), combined failed sections, and reported accurate `IN_PROGRESS` when workflows are running.
   - Rendered a compact recent commits pipeline summary table at the end of inspection.
   - Eliminated runaway blank lines via `CollapseConsecutiveEmptyLines`.
2. **`gitmap pipeline history` Subcommand**:
   - Displays a formatted execution tree across the last 5 (or `-n N`) commits with workflow branches, durations, and status badges.
3. **`gitmap pipeline logs [commit-or-offset]` Subcommand**:
   - Inspects and exports full or noise-filtered error logs for any target commit or offset.
   - Supports `--clip` (clipboard copy), `--tempfile`, `--file`, `--failed-only`, and `--json`.
4. **Repository SQLite Database Telemetry**:
   - Stored pipeline runs, commit SHAs, section errors, and logs in `.gitmap/data/pipeline.db`.
   - Formatted database telemetry (exact path, file size in human-readable units, total runs, failure counts).
   - Added SQLite vacuuming support (`db.Vacuum()`).
5. **`gitmap storage` & `gitmap os storage` Command Suites**:
   - Implemented `gitmap storage` / `status` reporting filesystem capacity and SQLite database inventory.
   - Implemented `gitmap storage ls` listing all repository-scoped and global databases.
   - Implemented `gitmap storage clean` purging ephemeral logs and temporary files (`--dry-run`, `--force`, `--vacuum`).
   - Integrated `gitmap os storage` via cross-package function delegation avoiding circular imports.
   - Updated help text across `storage.md`, `os.md`, and `pipeline.md`.

---

## 2. Verification Results
- **Cross-OS Type Checking**: `go vet ./...` passed with 0 errors across Windows, Linux (`GOOS=linux`), and macOS (`GOOS=darwin`).
- **Boolean & Enum Linters**: `check-boolean-guidelines.py` and `check-enum-and-boolean.py` passed with 0 violations across 2,918 files.
- **Control Flow Flattening**: `check-nested-ifs.py` passed with 0 nested-if or single-line compression violations across 2,918 files.
- **CI/CD Quality Gates**: `06-cicd-local-runner.py --no-tests` passed 38/38 gates 100% green in 35.07s.
- **Test Inventory**: Recorded 24 modified files, tracking 323 associated test files.
