// Package cmdsummary implements CI/CD pipeline error extraction and integration for full summary trees.
package cmdsummary

import (
	"os"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// attachPipelineErrorTelemetry interrogates the repository's pipeline DB and populates failure traces.
func attachPipelineErrorTelemetry(record *RepoSummaryRecord) {
	if record == nil || record.CanonicalSlug == "" {
		return
	}

	slug := pipelinedb.SanitizeRepoSlug(record.CanonicalSlug)
	dbPath := pipelinedb.ResolvePipelineDbPath(slug)
	if _, statErr := os.Stat(dbPath); statErr != nil {
		return
	}

	db, errOpen := pipelinedb.OpenPipelineSplitDb(slug)
	if errOpen != nil {
		return
	}
	defer db.Close()

	run, runErr := db.QueryRunByNegativeOffset(-1)
	if runErr != nil || run == nil {
		return
	}

	if run.IsSuccess || strings.EqualFold(run.Conclusion, "success") {
		record.IsPipelineClean = true
		return
	}

	errSummary := "Pipeline workflow run failed"
	stackTrace := ""
	stepName := "unknown"
	jobName := "unknown"
	exitCode := 1

	compactRes := db.QueryCompactErrorLogsByRunId(run.RunId)
	if compactRes.IsSuccess() && len(compactRes.Data) > 0 {
		first := compactRes.Data[0]
		errSummary = first.ErrorText
		stepName = first.StepName
		jobName = "workflow-job"
		exitCode = 1

		var traceLines []string
		for _, item := range compactRes.Data {
			if item.CompactLogs != "" {
				for _, l := range strings.Split(item.CompactLogs, "\n") {
					trimmed := strings.TrimRight(l, "\r")
					if trimmed != "" {
						traceLines = append(traceLines, trimmed)
					}
				}
			} else if item.ErrorText != "" {
				traceLines = append(traceLines, strings.TrimRight(item.ErrorText, "\r"))
			}
		}
		if len(traceLines) > 25 {
			traceLines = traceLines[len(traceLines)-25:]
		}
		stackTrace = strings.Join(traceLines, "\n")
	}

	record.PipelineError = &RepoPipelineError{
		RepoSlug:     slug,
		WorkflowName: run.WorkflowName,
		JobName:      jobName,
		StepName:     stepName,
		ExitCode:     exitCode,
		ErrorSummary: errSummary,
		StackTrace:   stackTrace,
		RunURL:       run.RunUrl,
	}
}
