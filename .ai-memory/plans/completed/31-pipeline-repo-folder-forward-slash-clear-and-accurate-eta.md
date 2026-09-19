# Plan 31: Pipeline Repository Folder Architecture, Forward Slash Paths, Clear Commands & Accurate ETA

## User Request (Verbatim)

```text
The pipeline errors in Git map make sure that the folder, all the folder and file paths in the same, um, folder slash. That would be the, the, the forward slash because it is the same, uh, regardless of the OS. And also there should be one or two commands that would help the user to clear the pipeline for this repo. So there should be like, uh, user can run the command inside the repo, then it would clean the pipe, I mean, uh, logs and everything for this, uh, this repository automatically. Uh, that means it would also include the repository logs. Uh, remember when you keep the repository logs, that is also lo- wrong because it needs to be in the repository folder. So inside the pipeline folder, I think the first thing you should do is create the repository folder, and inside this you will have like the, the, um, pipeline DB or DB, uh, or it's SQL.db. That would be the pipeline DB. And the logs and everything should be inside that folder so that it can be removed easily. Okay? And the commands should actually help, like both running inside the repo without repo args, or running with the repo name from outside. Also ETA duration calculation should be accurate... Currently, ETA (~68s) is skewed because it improperly calculates averages (e.g., mixing failed runs that abort in seconds with full successful runs) and provides inaccurate estimates. Store and compute accurate run duration baselines from historical successful runs.
```

## Extracted Actionable Task List & Verified Delivery

1. **Uniform Forward Slash Path Normalization**:
   - Normalized all file and directory paths in terminal rendering, payload metadata, error logs, and SQLite DB references via `filepath.ToSlash(...)`.
   - Updated `renderSavedLocationsTerminal`, `renderSavedDbAndUrl`, `renderRunCardBranchAndLog`, `renderCleanSuccessDbAndHistory`, `renderFailedJobSection`, `FormatRelativeDbPath`, and `GetDatabaseInfo` to guarantee pure forward slashes (`/`) across all platforms.
2. **Repository-Scoped Pipeline Folder Hierarchy**:
   - Reorganized pipeline filesystem storage under `<dataDir>/pipeline/<repo-slug>/` (e.g. `<dataDir>/pipeline/alimtvnetwork-gitmap-v28/`).
   - All repository-specific pipeline data—the SQLite DB (`pipeline.db` / `sql.db`), run logs (`<run-id>.log`), job metadata (`<run-id>.jobs.json`), run cache (`<run-id>.json`), and error reports (`pipeline_errors.log`, `last_error.log`)—now reside strictly within the dedicated repository directory.
   - Built automatic multi-path legacy migration in `migrateOrFallbackPipelineDb`, moving legacy databases and error logs from root `pipeline/` or `pipeline_db/` directly into the repository folder.
3. **Repository Pipeline Clear Subcommands**:
   - Implemented `gitmap pipeline clear` and `gitmap pipeline errors clear` to clear pipeline logs, cache, and database for the current repository when executed inside a repo directory.
   - Implemented `gitmap pipeline clear <repo>` and `gitmap pipeline errors clear <repo>` to clear a targeted repository's pipeline data from anywhere.
   - Added `extractTargetRepo` to parse target repository arguments and `-y` / `--yes` flag to bypass confirmation prompts.
   - Implemented `purgeRepoPipelineFolder` and `printPipelineClearSummary` displaying cleaned directory paths in forward slashes with reclaimed space and purged file counts.
4. **Accurate Pipeline Rerun Duration & ETA Engine**:
   - Redesigned `calculateAverageDuration` to filter out early-aborted, cancelled, and short skipped runs below a 45s threshold that skewed averages down to ~68s.
   - Implemented `QuerySuccessfulRunDurations` in `PipelineSplitDb` to query historical successful runs from SQLite whenever recent GitHub API runs contain only failures.
   - Implemented `computeBaselineDuration` taking robust upper-median durations across valid runs.
   - Updated `fallbackWorkflowDuration` with realistic baselines (180s for CI/smoke/test, 95s for release, 120s default).
5. **Documentation and Help Text Parity**:
   - Updated CLI help documentation in `cli/helptext/pipeline.md` and terminal help in `cli/cmdpipeline/pipeline.go` and `pipeline_db_dispatch.go`.
   - Added command examples and descriptions for `gitmap pipeline clear`, `gitmap pipeline clear <repo>`, and `gitmap pipeline errors clear`.

## Consolidated Subtasks

- Subtask 1: `01-forward-slash-normalization` - Verified all terminal renders and payloads use `filepath.ToSlash(...)`.
- Subtask 2: `02-repo-scoped-pipeline-folder-architecture` - Implemented `RepoPipelineDir`, `<dataDir>/pipeline/<slug>/pipeline.db`, and automatic legacy file migration.
- Subtask 3: `03-pipeline-clear-commands` - Registered `clear` and `errors clear`, supported repo args and `-y` flag, added `purgeRepoPipelineFolder`.
- Subtask 4: `04-accurate-duration-eta-engine` - Added outlier filtering, DB duration queries, upper-median baseline calculation, and enhanced ETA displays.
- Subtask 5: `05-cli-help-and-documentation` - Documented new clear commands and examples in markdown and terminal help.
