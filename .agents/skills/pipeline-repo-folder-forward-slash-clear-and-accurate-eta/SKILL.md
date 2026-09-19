---
name: pipeline-repo-folder-forward-slash-clear-and-accurate-eta
description: >-
  Autonomously implement repository-scoped pipeline data directories, uniform forward-slash
  path formatting, repository pipeline clear commands (both within repo and via arguments),
  and accurate historical pipeline ETA duration baseline calculations across GitMap.
---

# Pipeline Repo Folder Architecture, Forward Slash Formatting, Clear Commands & Accurate ETA

> **Prompt Version:** 2.2.0
> **Architecture:** GitMap Pipeline Storage & Telemetry Suite

## Core Responsibilities

1. **Uniform Forward Slash (`/`) Formatting**:
   - In all pipeline commands (`gitmap pipeline errors`, `gitmap pipeline logs`, `gitmap pipeline status`, `gitmap pipeline history`, `gitmap pipeline db status`), format all folder and file paths (`Combined Report`, `Latest Run Log`, `Saved Log`, `Pipeline DB`, etc.) with forward slashes (`/`), eliminating OS-dependent mixed backslashes on Windows.
   - Enforce `filepath.ToSlash(...)` across path resolution, serialization, and terminal rendering layers.

2. **Repository-Scoped Pipeline Folder Architecture**:
   - In the pipeline data root (`<dataDir>/pipeline/`), create an isolated directory for each repository named by its sanitized repository slug (e.g. `<dataDir>/pipeline/<repo-slug>/` like `alimtvnetwork-gitmap-v28/`).
   - Store all repository pipeline artifacts within this folder:
     - SQLite Database: `<dataDir>/pipeline/<repo-slug>/pipeline.db` (with backward-compatible detection/migration of legacy `pipeline_<slug>.db`).
     - Raw & Parsed Logs: `<dataDir>/pipeline/<repo-slug>/<run-id>.log`
     - Job Metadata: `<dataDir>/pipeline/<repo-slug>/<run-id>.jobs.json`
     - Run Cache JSON: `<dataDir>/pipeline/<repo-slug>/<run-id>.json`
     - Error Report: `<dataDir>/pipeline/<repo-slug>/pipeline_errors.log`
     - Last Error Log: `<dataDir>/pipeline/<repo-slug>/last_error.log`

3. **Repository Pipeline Clear Commands**:
   - Inside a repository without repo args: `gitmap pipeline clear` and `gitmap pipeline errors clear` automatically resolve the current repository slug, purge all cached logs, reports, and reset/clear the SQLite DB in `<dataDir>/pipeline/<repo-slug>/`.
   - From anywhere with target repo arg: `gitmap pipeline clear <repo>` and `gitmap pipeline errors clear <repo>` clear that specific repository's pipeline folder and DB.
   - Support `-y` / `--yes` flag to bypass confirmation prompt, and display clear summary with reclaimed space and cleaned directory in forward slashes.
   - Retain backward compatibility for `gitmap pipeline clear-db` and `gitmap pipeline db clear`.

4. **Accurate Pipeline Duration Baseline & ETA**:
   - Filter out early-aborted, cancelled, and skipped workflow runs from duration averaging. Only analyze full, valid successful runs (`conclusion == "success"` and duration exceeding minimal thresholds).
   - Target workflow-specific durations rather than averaging completely disjoint workflows together.
   - Use robust duration statistics (e.g. median / 75th percentile / unskewed mean) and fallback baselines to ensure realistic rerun ETAs.
   - Display informative baseline context in terminal output (e.g., `Based on historical successful pipeline runs baseline`).
