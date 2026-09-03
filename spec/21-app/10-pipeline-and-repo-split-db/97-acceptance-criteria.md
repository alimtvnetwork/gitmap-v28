# Acceptance Criteria: Pipeline and Repo Split Databases

## Scenario 1: Pipeline Split Database Initialization
- **Given** an active repository with remote origin `alimtvnetwork/gitmap-v28`
- **When** `gitmap pipeline db status` is executed
- **Then** the command outputs the pipeline split database path under `data/pipeline_db/pipeline_alimtvnetwork-gitmap-v28.db`
- **And** displays total file size, run count, and error count.

## Scenario 2: Pipeline Split Database Clear and Reset
- **Given** a pipeline split database containing recorded runs
- **When** `gitmap pipeline db clear -y` is executed
- **Then** all rows from `PipelineRun`, `PipelineErrorLog`, and `PipelineSegment` are truncated
- **And** subsequent `status` reports 0 runs and 0 error logs.

## Scenario 3: Pipeline Split Database Optimization
- **Given** a pipeline split database file
- **When** `gitmap pipeline db optimize` is executed
- **Then** SQLite `VACUUM` and `PRAGMA optimize` execute successfully
- **And** disk space reclaimed is reported to stdout.

## Scenario 4: Repository Split Database Status and Logging
- **Given** an initialized repository in GitMap
- **When** `gitmap repo db status` is executed
- **Then** the command displays the split DB path under `repo_search/`
- **And** displays table counts for `RepoFile`, `SearchCache`, and `RepoScanLog`.

## Scenario 5: Unified Database Status
- **Given** initialized Master, Repo, and Pipeline databases
- **When** `gitmap db status` is executed
- **Then** a consolidated summary of all three tiers with file paths and sizes is printed.
