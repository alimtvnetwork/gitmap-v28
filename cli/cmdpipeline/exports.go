package cmdpipeline

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// RunPipeline executes the pipeline command group.
func RunPipeline(args []string) error {
	return runPipeline(args)
}

// RunPipelineAI executes the pipeline AI command with delays.
func RunPipelineAI(args []string) error {
	return runPipelineAI(args)
}

// HandlePipelineStatus handles pipeline status output.
func HandlePipelineStatus(args []string) error {
	return handlePipelineStatus(args)
}

// HandlePipelineErrorLogs handles pipeline error log inspection.
func HandlePipelineErrorLogs(args []string) error {
	return handlePipelineErrorLogs(args)
}

// HandlePipelineDB handles pipeline DB operations.
func HandlePipelineDB(args []string) error {
	return handlePipelineDB(args)
}

// HandlePipelineWaitTime calculates the ETA wait time.
func HandlePipelineWaitTime(args []string) error {
	return handlePipelineWaitTime(args)
}

// ResolveTempDir returns the configured or default temp dir.
func ResolveTempDir() string {
	return resolveTempDir()
}

// GhRunItem aliases the unexported ghRunItem.
type GhRunItem = ghRunItem

// BuildErrorLogsPayload builds the pipeline error logs payload.
func BuildErrorLogsPayload(repo string, runs []GhRunItem) PipelineErrorLogsPayload {
	return buildErrorLogsPayload(repo, runs)
}

// RecordRunInSplitDb writes run telemetry into the pipeline split DB.
func RecordRunInSplitDb(db *pipelinedb.PipelineSplitDb, repo string, run GhRunItem) error {
	return recordRunInSplitDb(db, repo, run)
}

// SaveParsedFailedJobs records parsed failure structures into the split DB.
func SaveParsedFailedJobs(db *pipelinedb.PipelineSplitDb, repo string, run GhRunItem, jobs []FailedJobItem, rawLogs string) {
	saveParsedFailedJobs(db, repo, run, jobs, rawLogs)
}

// FetchLatestPipelineErrorReport returns the formatted error report string and whether failures exist.
func FetchLatestPipelineErrorReport(repo string, isDetailed bool) (string, bool) {
	_, report, hasFailures := FetchPipelineErrorReportWithMeta(repo, isDetailed)

	return report, hasFailures
}

// FetchPipelineErrorReportWithMeta returns payload metadata, formatted error report, and failure status.
func FetchPipelineErrorReportWithMeta(repo string, isDetailed bool) (PipelineErrorLogsPayload, string, bool) {
	targetRepo := repo
	if len(targetRepo) == 0 {
		targetRepo = resolveCurrentRepoSlug()
	}

	runs := queryWorkflowRuns(targetRepo)
	payload := buildErrorLogsPayload(targetRepo, runs)
	if !isDetailed {
		compactErrorPayload(&payload)
	}

	hasFailure := payload.Conclusion == "failure" || len(payload.FailedRuns) > 0
	if hasFailure {
		return payload, buildClipboardErrorReport(payload), true
	}

	return payload, buildClipboardCleanReport(payload), false
}
