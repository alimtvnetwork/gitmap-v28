---
name: pipeline-compact-error-logs
description: Autonomously filter ok lines from pipeline error logs by default, persist dual detailed and compact logs in repository-scoped database across 3 connected tables, and support detailed/verbose flags.
---

# Pipeline Compact Error Logs & 3-Table Repository Database Architecture

## Overview
Filters passing `ok` lines, test status lines, and non-failure noise from `gitmap pipeline error-logs` (and aliases `errorlogs`, `errors`, `err`, `last-failed-logs`) by default after reading from the server.
Supports `--detailed`, `--verbose`, `--v`, and `-v` flags to display full, uncompressed logs with all lines preserved.

## Repository DB Isolation (Strict Policy)
Pipeline execution records and error logs must be kept strictly in the repository-scoped split database (`data/pipeline_db/pipeline_<repoSlug>.db`), **never** in the main global application database (`gitmap.db`).

## 3 Connected Tables Architecture
The repository-scoped database organizes pipeline diagnostics into 3 connected tables linked via `RunId`:

1. **Master Table (`PipelineRun`)**:
   - Contains all pipeline run details, execution status, metadata, commit SHA, branch, duration, ETA, conclusion, and URLs.
   - Master record for every CI/CD workflow run.
2. **Detailed Error Logs Sub-Table (`PipelineDetailErrorLog`)**:
   - Connected to master table via `RunId REFERENCES PipelineRun(RunId)`.
   - Stores full, uncompressed raw logs fetched from the server, preserving all diagnostic lines verbatim.
3. **Compact Error Logs Sub-Table (`PipelineCompactErrorLog`)**:
   - Connected to master table via `RunId REFERENCES PipelineRun(RunId)`.
   - Stores compact error logs with all passing `ok lines` removed (`ok\t`, `ok `, `?\t`, `? `, `--- PASS:`, `=== RUN`, `PASS`, `✔ ok`, `✔ Macro`).
   - Tracks `FilteredOkCount` indicating how many noise lines were eliminated.

## CLI Flag Matrix
- **Default (Compact)**:
  `gitmap pipeline error-logs`, `errors`, `errorlogs`, `err`, `last-failed-logs`
  Displays compact version without `ok lines`.
- **Non-Compact (Detailed)**:
  `gitmap pipeline error-logs --detailed` (or `--verbose`, `--v`, `-v`, `-V`)
  Displays full detailed logs with all lines included.
