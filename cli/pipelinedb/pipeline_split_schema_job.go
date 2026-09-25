package pipelinedb

const sqlCreatePipelineJob = `
CREATE TABLE IF NOT EXISTS PipelineJob (
    PipelineJobId INTEGER PRIMARY KEY AUTOINCREMENT,
    JobId INTEGER NOT NULL UNIQUE,
    RunId INTEGER NOT NULL,
    RepoSlug TEXT NOT NULL,
    JobName TEXT NOT NULL,
    Status TEXT NOT NULL,
    Conclusion TEXT NOT NULL,
    StartedAt TEXT NOT NULL,
    CompletedAt TEXT NOT NULL,
    DurationSeconds INTEGER DEFAULT 0,
    JobUrl TEXT NOT NULL,
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);
CREATE INDEX IF NOT EXISTS IdxPipelineJob_RunId ON PipelineJob (RunId);
CREATE INDEX IF NOT EXISTS IdxPipelineJob_RepoSlug ON PipelineJob (RepoSlug);
CREATE INDEX IF NOT EXISTS IdxPipelineJob_RepoSlug_RunId ON PipelineJob (RepoSlug, RunId);
CREATE INDEX IF NOT EXISTS IdxPipelineJob_JobName ON PipelineJob (JobName);`
