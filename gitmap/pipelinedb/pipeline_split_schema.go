package pipelinedb

const (
	sqlCreatePipelineRun = `
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
);`

	sqlCreatePipelineErrorLog = `
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
);`

	sqlCreatePipelineSegment = `
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
);`
)
