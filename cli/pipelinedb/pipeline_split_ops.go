package pipelinedb

const sqlRecordRun = `
INSERT INTO PipelineRun (
    RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
    EtaSeconds, DurationSeconds, RunUrl, IsSuccess, Notes, Comments, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId) DO UPDATE SET
    Status = excluded.Status,
    Conclusion = excluded.Conclusion,
    EtaSeconds = excluded.EtaSeconds,
    DurationSeconds = excluded.DurationSeconds,
    IsSuccess = excluded.IsSuccess,
    UpdatedAt = excluded.UpdatedAt;`

const sqlRecordErrorLog = `
INSERT INTO PipelineErrorLog (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, RawLogs, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId, StepName) DO UPDATE SET
    RepoSlug = excluded.RepoSlug,
    WorkflowName = excluded.WorkflowName,
    ErrorText = excluded.ErrorText,
    RawLogs = excluded.RawLogs,
    Notes = excluded.Notes,
    Comments = excluded.Comments,
    CreatedAt = excluded.CreatedAt;`

const sqlRecordDetailErrorLog = `
INSERT INTO PipelineDetailErrorLog (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, RawLogs, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId, StepName) DO UPDATE SET
    RepoSlug = excluded.RepoSlug,
    WorkflowName = excluded.WorkflowName,
    ErrorText = excluded.ErrorText,
    RawLogs = excluded.RawLogs,
    Notes = excluded.Notes,
    Comments = excluded.Comments,
    CreatedAt = excluded.CreatedAt;`

const sqlRecordCompactErrorLog = `
INSERT INTO PipelineCompactErrorLog (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, CompactLogs, FilteredOkCount, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)
ON CONFLICT(RunId, StepName) DO UPDATE SET
    RepoSlug = excluded.RepoSlug,
    WorkflowName = excluded.WorkflowName,
    ErrorText = excluded.ErrorText,
    CompactLogs = excluded.CompactLogs,
    FilteredOkCount = excluded.FilteredOkCount,
    Notes = excluded.Notes,
    Comments = excluded.Comments,
    CreatedAt = excluded.CreatedAt;`

const sqlQueryRecentRuns = `
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun ORDER BY CreatedAt DESC, RunId DESC LIMIT ?;`

const sqlQueryRecentErrors = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineErrorLog ORDER BY PipelineErrorLogId DESC LIMIT ?;`

const sqlQueryRecentDetailErrors = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineDetailErrorLog ORDER BY PipelineDetailErrorLogId DESC LIMIT ?;`

const sqlQueryRecentCompactErrors = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(CompactLogs, ''), FilteredOkCount, CreatedAt
FROM PipelineCompactErrorLog ORDER BY PipelineCompactErrorLogId DESC LIMIT ?;`

const sqlQueryCachedErrorRunIds = `
SELECT DISTINCT RunId FROM (
	SELECT RunId FROM PipelineErrorLog
	UNION
	SELECT RunId FROM PipelineDetailErrorLog
	UNION
	SELECT RunId FROM PipelineCompactErrorLog
) ORDER BY RunId DESC;`

const sqlQueryRunByOffset = `
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun ORDER BY CreatedAt DESC, RunId DESC LIMIT 1 OFFSET ?;`

const sqlQueryLastFailedRuns = `
SELECT RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
       EtaSeconds, DurationSeconds, RunUrl, IsSuccess, CreatedAt, UpdatedAt
FROM PipelineRun WHERE IsSuccess = 0 ORDER BY CreatedAt DESC, RunId DESC LIMIT ?;`

const sqlQueryErrorLogsByRunId = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineErrorLog WHERE RunId = ? ORDER BY PipelineErrorLogId ASC;`

const sqlQueryDetailErrorLogsByRunId = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(RawLogs, ''), CreatedAt
FROM PipelineDetailErrorLog WHERE RunId = ? ORDER BY PipelineDetailErrorLogId ASC;`

const sqlQueryCompactErrorLogsByRunId = `
SELECT RunId, RepoSlug, WorkflowName, StepName, ErrorText, COALESCE(CompactLogs, ''), FilteredOkCount, CreatedAt
FROM PipelineCompactErrorLog WHERE RunId = ? ORDER BY PipelineCompactErrorLogId ASC;`
