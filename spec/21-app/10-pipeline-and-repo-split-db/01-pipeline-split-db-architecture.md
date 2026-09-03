# Pipeline Split Database Architecture

## 1. Storage Location & Path Formulation

Pipeline split databases are isolated per repository slug to ensure zero lock contention between concurrent CI/CD pipeline monitoring routines:

- **Root Directory:** `<BinaryDataDir>/pipeline_db/`
- **File Naming Pattern:** `pipeline_<sanitized_slug>.db`
- **Sanitization Rule:** Replaces forward slashes (`/`), colons (`:`), and backslashes (`\`) with hyphens (`-`).
  - Example: Repository `alimtvnetwork/gitmap-v28` resolves to:
    `data/pipeline_db/pipeline_alimtvnetwork-gitmap-v28.db`

## 2. Table Schemas

All schemas adhere strictly to `spec/04-database-conventions/` and `spec/05-split-db-architecture/`:

### Table: `PipelineRun`

Tracks historical workflow runs and completion duration baselines.

```sql
CREATE TABLE IF NOT EXISTS PipelineRun (
    PipelineRunId INTEGER PRIMARY KEY AUTOINCREMENT,
    RunId INTEGER NOT NULL UNIQUE,
    RepoSlug TEXT NOT NULL,
    WorkflowName TEXT NOT NULL,
    Status TEXT NOT NULL,
    Conclusion TEXT NOT NULL,
    Branch TEXT NOT NULL,
    Sha TEXT NOT NULL,
    EtaSeconds INTEGER DEFAULT 0,
    DurationSeconds INTEGER DEFAULT 0,
    RunUrl TEXT NOT NULL,
    IsSuccess INTEGER DEFAULT 0,
    Notes TEXT NULL,
    Comments TEXT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);
```

### Table: `PipelineErrorLog`

Isolates and records targeted CI/CD failure errors, stack traces, and failing steps.

```sql
CREATE TABLE IF NOT EXISTS PipelineErrorLog (
    PipelineErrorLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    RunId INTEGER NOT NULL,
    RepoSlug TEXT NOT NULL,
    WorkflowName TEXT NOT NULL,
    StepName TEXT NOT NULL,
    ErrorText TEXT NOT NULL,
    RawLogs TEXT NULL,
    Notes TEXT NULL,
    Comments TEXT NULL,
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);
```

### Table: `PipelineSegment`

Tracks individual step and job execution runtimes to power segment-level ETA estimations.

```sql
CREATE TABLE IF NOT EXISTS PipelineSegment (
    PipelineSegmentId INTEGER PRIMARY KEY AUTOINCREMENT,
    RunId INTEGER NOT NULL,
    JobName TEXT NOT NULL,
    StepName TEXT NOT NULL,
    StepNumber INTEGER DEFAULT 0,
    Status TEXT NOT NULL,
    Conclusion TEXT NOT NULL,
    DurationSeconds INTEGER DEFAULT 0,
    Notes TEXT NULL,
    Comments TEXT NULL,
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);
```

## 3. CLI Subcommands: `gitmap pipeline db`

| Subcommand | Description |
|---|---|
| `status` (default) | Shows pipeline split DB location, path, size, run counts, error log counts, and SQLite metadata |
| `clear` | Truncates run and error log records with confirmation prompt (`-y` to skip) |
| `reset` | Drops all tables and re-executes clean schema initialization |
| `optimize` | Runs `PRAGMA wal_checkpoint(TRUNCATE)`, `VACUUM`, `PRAGMA optimize`, and `ANALYZE`, returning bytes reclaimed |
| `errorlogs` / `error-logs` | Queries and lists stored failure logs for the repository |
| `help` | Prints detailed syntax and example usage |
