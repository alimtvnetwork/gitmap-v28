---
name: pipeline-commit-history-errors-storage
description: >-
  Autonomously implement commit-based pipeline error inspection with negative offsets (-1, -2, -3),
  multi-run error aggregation, noise filtering, pipeline history tree, commit-scoped pipeline logs,
  repository SQLite telemetry, and root/OS storage command parity across GitMap.
---

# Pipeline Commit-Based Errors, History, Logs, SQLite Telemetry & Storage Parity

> **Prompt Version:** 2.1.0
> **Architecture:** GitMap Pipeline Engine & CLI Subcommand Suite

## Core Responsibilities

1. **Commit-Based Offset Inspection (`-1`, `-2`, `-3`, etc.)**:
   - Translate negative integer flags (`-1`, `-2`, `-3`, `-N`) in `gitmap pipeline errors` (and aliases `errorlogs`, `errors`, `err`) to mean **N commits before latest** (or the Nth previous distinct commit in pipeline run history), NOT raw individual run indices.
   - Aggregate all pipeline workflows (CI, Release, Lint, etc.) for that specific commit, combining all section failures into a single cohesive report.
   - Print a summary of previous commits with failed pipelines at the end.
   - Eliminate runaway blank newlines in output.

2. **Noise / Ok Line Filtering**:
   - Filter out passing test suites, ok lines, and runtime noise while preserving historical failure lines, tracebacks, and contextual snippets.

3. **Pipeline History Subcommand (`gitmap pipeline history`)**:
   - Display a tree/table of the last 5 (configurable `-n N`) commits with their pipeline run statuses (succeeded, failed, in-progress).
   - Display workflow breakdown and allow selecting/identifying which commit to inspect.

4. **Commit-Scoped Pipeline Logs (`gitmap pipeline logs <commit-or-offset>`)**:
   - Retrieve and display/export the consolidated logs for all workflows of a specific commit or offset.

5. **Repository SQLite Database Telemetry**:
   - Persist runs, commit metadata, section errors, and full logs in the repository-scoped SQLite database (`.gitmap/data/pipeline.db`).
   - Display SQLite database information: exact location and formatted database file size.

6. **Storage Command Parity (`gitmap storage` & `gitmap os storage`)**:
   - Implement `gitmap storage` and `gitmap os storage` commands to report storage breakdown: repository size, Git objects, GitHub Actions cache/artifacts, SQLite databases, temp files, and disk usage.
   - Maintain complete CLI help text parity in `cli/helptext/`.
