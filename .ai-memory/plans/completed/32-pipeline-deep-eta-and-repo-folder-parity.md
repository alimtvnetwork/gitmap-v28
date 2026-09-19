# Plan 32: Deep Pipeline ETA Decision Engine, Database Telemetry Persistence, and Repository Folder Parity

## Overview & Execution Context

- **Workflow:** Parent Task N-Step Continuous Loop & Multi-Agent Orchestration (Prompt v2.2.0, N=250)
- **Start Reference:** User requested rigorous audit of all missing elements, a clear Done vs. Pending table, forward-slash normalization across all OSes, repository-scoped pipeline folder architecture, dual-mode clear commands (`inside repo` and `with repo name`), two annotated examples in the pipeline section, and a deep overhaul of how ETA decisions are computed, how duration/success data is stored in SQLite, and how ETA information is presented to the user.
- **Execution Loops:** Completed in 4 consolidated loops across planning, implementation, verification, and documentation.

## User Request (Verbatim)

```text
is it done properly, please check carefully all the missing elements and tasks from the below, can you please showcase a list or table of what is done and what is pending?



The pipeline errors in Git map make sure that the folder, all the folder and file paths in the same, um, folder slash. That would be the, the, the forward slash because it is the same, uh, regardless of the OS. And also there should be one or two commands that would help the user to clear the pipeline for this repo. So there should be like, uh, user can run the command inside the repo, then it would clean the pipe, I mean, uh, logs and everything for this, uh, this repository automatically. Uh, that means it would also include the repository logs. Uh, remember when you keep the repository logs, that is also lo- wrong because it needs to be in the repository folder. So inside the pipeline folder, I think the first thing you should do is create the repository folder, and inside this you will have like the, the, um, pipeline DB or DB, uh, or it's SQL.db. That would be the pipeline DB. And the logs and everything should be inside that folder so that it can be removed easily. Okay? And the commands should actually appear, uh, that it could run, uh, inside the repo without the repo name, and also can be run from anywhere with the repo name to clear all the pipeline. It could be like pipeline space clear, something like this, um, but also the name. So two examples needs to be there in the pipeline section. And, uh, yeah, I think a-another, uh, important information that is required, uh, currently you do have it, the ETA, uh, the total ETA that is required to run the whole pipeline. I think that needs to be a little bit more accurate. Currently, how you are making the decision, how you're keeping the data in the database and, uh, providing the information, these are not accurate, not accurate. So I want you to dig deep, understand this and improve these parts. Do you understand? Can you please help me in this order?
```

## Extracted Actionable Task List & Verified Delivery

1. **Uniform Forward Slash Path Normalization Across All OSes**:
   - Normalized all file and directory paths in terminal rendering, payload metadata, error logs, and SQLite DB references via `filepath.ToSlash(...)`.
   - Updated `renderSavedLocationsTerminal`, `renderSavedDbAndUrl`, `renderRunCardBranchAndLog`, `renderCleanSuccessDbAndHistory`, `renderFailedJobSection`, `FormatRelativeDbPath`, and `GetDatabaseInfo` to guarantee pure forward slashes (`/`) across all platforms.
2. **Repository-Scoped Pipeline Folder Hierarchy**:
   - Reorganized pipeline filesystem storage under `<dataDir>/pipeline/<repo-slug>/` (e.g. `<dataDir>/pipeline/alimtvnetwork-gitmap-v28/`).
   - All repository-specific pipeline data—the SQLite DB (`pipeline.db` / `sql.db`), run logs (`<run-id>.log`), job metadata (`<run-id>.jobs.json`), run cache (`<run-id>.json`), and error reports (`pipeline_errors.log`, `last_error.log`)—now reside strictly within the dedicated repository directory.
   - Built automatic multi-path legacy migration in `migrateOrFallbackPipelineDb`, moving legacy databases and error logs from root `pipeline/` or `pipeline_db/` directly into the repository folder.
3. **Repository Pipeline Clear Subcommands & Two Annotated Examples**:
   - Implemented `gitmap pipeline clear` and `gitmap pipeline errors clear` to clear pipeline logs, cache, and database for the current repository when executed inside a repo directory.
   - Implemented `gitmap pipeline clear <repo>` and `gitmap pipeline errors clear <repo>` to clear a targeted repository's pipeline data from anywhere.
   - Added `extractTargetRepo` to parse target repository arguments and `-y` / `--yes` flag to bypass confirmation prompts.
   - Implemented `purgeRepoPipelineFolder` and `printPipelineClearSummary` displaying cleaned directory paths in forward slashes with reclaimed space and purged file counts.
   - Annotated two explicit clear examples in `cli/cmdpipeline/pipeline.go`, `cli/cmdpipeline/pipeline_db_dispatch.go`, and `cli/helptext/pipeline.md`:
     - `gitmap pipeline clear -y                          # Clear pipeline for current repo (run inside repo)`
     - `gitmap pipeline clear alimtvnetwork/gitmap-v28 -y # Clear pipeline for target repo (run from anywhere)`
4. **Deep Pipeline ETA Decision Engine, Database Persistence & Human-Friendly Presentation**:
   - **Database Telemetry Storage**: Fixed `buildPipelineRunRecord` and `insertSingleMasterRun` in `cli/cmdpipeline/pipeline_recorder.go` to compute and record true `DurationSeconds` (`calculateRunDuration`), `IsSuccess` boolean flag (`r.Conclusion == "success"`), timestamps (`CreatedAt`, `UpdatedAt`), and zero out remaining `EtaSeconds` for completed runs.
   - **Historical Database Querying**: Implemented `QuerySuccessfulRunDurations` in `PipelineSplitDb` to pull historical successful run durations directly from SQLite when GitHub Actions runs contain only failures or recent aborts.
   - **Outlier & Abort Filtering**: Redesigned `calculateAverageDuration` to filter out early-aborted, canceled, and short skipped runs below a 45s threshold that skewed averages down to ~68s.
   - **Upper-Median Baseline Decision**: Implemented `computeBaselineDuration` taking robust upper-median durations across valid runs rather than skewed arithmetic means.
   - **Human-Friendly Duration Presentation**: Added `formatEtaDisplay` in `cli/cmdpipeline/pipeline_query.go` providing human-readable duration strings (e.g. `~180s (3m)`, `~125s (2m 5s)`) in terminal banners, dynamic timeline, and status cards.

## Consolidated Subtasks

- Subtask 1: `01-forward-slash-and-repo-folder-audit` - Verified uniform forward slash paths and `<dataDir>/pipeline/<repo-slug>/` folder hierarchy.
- Subtask 2: `02-clear-commands-and-help-examples-audit` - Verified dual clear commands (inside repo & from anywhere) and annotated help examples.
- Subtask 3: `03-deep-eta-decision-and-db-storage-engine` - Implemented database telemetry persistence, outlier filtering, historical DB querying, upper-median decision baseline, and human-friendly duration formatting.
