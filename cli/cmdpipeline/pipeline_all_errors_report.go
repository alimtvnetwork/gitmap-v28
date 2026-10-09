package cmdpipeline

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/pipelinedb"
)

// allErrorsLogExcerptMaxLines caps the RawLogs window kept per error block.
const allErrorsLogExcerptMaxLines = 200

// writeAllPipelineErrorsReportFile writes the combined report to flags.FilePath.
// With --json the JSON report is written; otherwise the AI-readable markdown.
// The bytes are written raw: no glyph or emoji filtering, ever.
func writeAllPipelineErrorsReportFile(flags PipelineErrorFlags, summary AllPipelineSummary) error {
	var content []byte

	if flags.IsJSON {
		jsonBytes, marshalErr := json.MarshalIndent(summary, "", "  ")

		if marshalErr != nil {
			return apperror.WrapSimple(marshalErr, "marshal all pipeline errors json report")
		}

		content = append(jsonBytes, '\n')
	} else {
		content = []byte(buildAllPipelineErrorsMarkdown(summary))
	}

	dir := filepath.Dir(flags.FilePath)

	if len(dir) > 0 {
		if mkErr := os.MkdirAll(dir, 0755); mkErr != nil {
			return apperror.WrapSimple(mkErr, "create report directory")
		}
	}

	if writeErr := os.WriteFile(flags.FilePath, content, 0644); writeErr != nil {
		return apperror.WrapSimple(writeErr, "write all pipeline errors report")
	}

	fmt.Printf("  %s✓%s Output written to %s\n", constants.ColorGreen, constants.ColorReset, flags.FilePath)

	return nil
}

// buildAllPipelineErrorsMarkdown renders the AI-fix dossier: header, one
// section per failing repo with full per-error blocks, the trailing
// repos-without-errors table, and the pull-errors appendix when present.
func buildAllPipelineErrorsMarkdown(summary AllPipelineSummary) string {
	var sb strings.Builder

	writeAllErrorsMarkdownHeader(&sb, summary)

	for _, item := range summary.FailedItems {
		writeFailedRepoMarkdownSection(&sb, summary, item)
	}

	writeNoErrorsMarkdownTable(&sb, summary)
	writePullErrorsMarkdownAppendix(&sb, summary)

	return sb.String()
}

func writeAllErrorsMarkdownHeader(sb *strings.Builder, summary AllPipelineSummary) {
	sb.WriteString("# Pipeline Errors — All Repositories (AI Fix Dossier)\n\n")
	sb.WriteString("- **Generated At:** " + time.Now().UTC().Format(time.RFC3339) + "\n")
	sb.WriteString("- **Command:** `gitmap pipeline errors all --file`\n")
	sb.WriteString(fmt.Sprintf("- **Repos Scanned:** %d\n", summary.ReposScanned))
	sb.WriteString(fmt.Sprintf("- **Repos With Errors:** %d\n", summary.FailedCount))
	sb.WriteString(fmt.Sprintf("- **Total Errors:** %d\n\n", summary.TotalErrorEntries))
	sb.WriteString("---\n\n")
}

func writeFailedRepoMarkdownSection(sb *strings.Builder, summary AllPipelineSummary, item AllPipelineFailedItem) {
	errorRecords := summary.errorDetails[item.RepoSlug]
	jobNames := summary.stepJobs[item.RepoSlug]

	sb.WriteString("## Repo: " + item.RepoSlug + "\n\n")
	sb.WriteString("- **Repo Slug:** `" + item.RepoSlug + "`\n")
	sb.WriteString("- **Repo Path:** `" + summary.repoAbsPaths[item.RepoSlug] + "`\n")
	sb.WriteString("- **Branch:** `" + item.Branch + "`\n")
	sb.WriteString("- **Latest Run:** `" + fmt.Sprintf("%d", item.RunId) + "` (" + item.WorkflowName +
		", SHA `" + item.Sha + "`, " + normalizeUtcTimestamp(item.UpdatedAt) + ")\n")
	sb.WriteString(fmt.Sprintf("- **Status:** FAILING (%d error(s))\n\n", len(errorRecords)))

	for i, rec := range errorRecords {
		writeErrorBlockMarkdown(sb, item, rec, i+1, jobNames[rec.StepName])
	}
}

func writeErrorBlockMarkdown(sb *strings.Builder, item AllPipelineFailedItem, rec pipelinedb.PipelineErrorRecord, index int, jobName string) {
	stepLine := rec.StepName

	if jobName != "" {
		stepLine = rec.StepName + " (job: " + jobName + ")"
	}

	timestamp := rec.CreatedAt

	if timestamp == "" {
		timestamp = item.UpdatedAt
	}

	sb.WriteString(fmt.Sprintf("### Error %d\n\n", index))
	sb.WriteString("- **Workflow:** `" + rec.WorkflowName + "`\n")
	sb.WriteString("- **Run ID:** `" + fmt.Sprintf("%d", rec.RunId) + "`\n")
	sb.WriteString("- **Branch:** `" + item.Branch + "`\n")
	sb.WriteString("- **SHA:** `" + item.Sha + "`\n")
	sb.WriteString("- **Failed Step:** `" + stepLine + "`\n")
	sb.WriteString("- **Timestamp:** `" + normalizeUtcTimestamp(timestamp) + "`\n")
	sb.WriteString("- **Error Summary:** `" + firstLineOr(item.ErrorSummary, rec.ErrorText) + "`\n")

	if rec.ErrorText == "" && !hasTracebackMarkers(rec.RawLogs) {
		sb.WriteString("- **Traceback:** not available — raw log excerpt only\n")
	}

	sb.WriteString("\n**Full Error Text:**\n\n```text\n")

	if rec.ErrorText == "" {
		sb.WriteString("(no extracted error text — see log excerpt below)\n")
	} else {
		sb.WriteString(rec.ErrorText)

		if !strings.HasSuffix(rec.ErrorText, "\n") {
			sb.WriteString("\n")
		}
	}

	sb.WriteString("```\n\n**Log Excerpt:**\n\n```text\n")
	sb.WriteString(boundLogExcerpt(rec.RawLogs, allErrorsLogExcerptMaxLines))
	sb.WriteString("\n```\n\n---\n\n")
}

// detailedErrorRecordsForRun fetches the full error records for a run
// (ErrorText + RawLogs), used for the per-error markdown blocks.
func detailedErrorRecordsForRun(db *pipelinedb.PipelineSplitDb, runId uint64) []pipelinedb.PipelineErrorRecord {
	res := db.QueryErrorLogsByRunId(runId)

	if res.IsFailure() {
		return nil
	}

	return res.Data
}

// stepJobNamesForRun maps step names to their parent job names via the run segments.
func stepJobNamesForRun(db *pipelinedb.PipelineSplitDb, runId uint64) map[string]string {
	names := map[string]string{}
	segments, err := db.QueryPipelineSegments(runId)

	if err != nil {
		return names
	}

	for _, seg := range segments {
		if _, exists := names[seg.StepName]; !exists {
			names[seg.StepName] = seg.JobName
		}
	}

	return names
}

// boundLogExcerpt keeps a bounded RawLogs window, cutting from the head so the
// failure stays at the tail; truncation is always marked explicitly.
func boundLogExcerpt(rawLogs string, maxLines int) string {
	if rawLogs == "" {
		return "(no log excerpt available)"
	}

	lines := strings.Split(rawLogs, "\n")

	if len(lines) <= maxLines {
		return rawLogs
	}

	truncated := fmt.Sprintf("[log excerpt truncated: showing last %d lines of %d]", maxLines, len(lines))

	return truncated + "\n" + strings.Join(lines[len(lines)-maxLines:], "\n")
}

func writeNoErrorsMarkdownTable(sb *strings.Builder, summary AllPipelineSummary) {
	sb.WriteString("## Repos Without Errors\n\n")
	sb.WriteString("| Repo | Path | Status |\n")
	sb.WriteString("| ---- | ---- | ------ |\n")

	for _, entry := range sortedNonFailedEntries(summary) {
		sb.WriteString("| `" + entry.slug + "` | `" + entry.absPath + "` | " +
			displayDataStatus(summary.repoDataStatus[entry.slug]) + " |\n")
	}

	sb.WriteString("\n")
}

func sortedNonFailedEntries(summary AllPipelineSummary) []repoCatalogEntry {
	failed := map[string]bool{}

	for _, item := range summary.FailedItems {
		failed[item.RepoSlug] = true
	}

	entries := []repoCatalogEntry{}

	for _, entry := range summary.repoList {
		if !failed[entry.slug] {
			entries = append(entries, entry)
		}
	}

	sort.Slice(entries, func(i, j int) bool {
		if entries[i].slug != entries[j].slug {
			return entries[i].slug < entries[j].slug
		}

		return entries[i].absPath < entries[j].absPath
	})

	return entries
}

func displayDataStatus(status string) string {
	switch status {
	case allRepoStatusClean:
		return "clean"
	case allRepoStatusDbUnreachable:
		return "pipeline DB unavailable"
	default:
		return "no pipeline data"
	}
}

func writePullErrorsMarkdownAppendix(sb *strings.Builder, summary AllPipelineSummary) {
	if len(summary.PullErrors) == 0 {
		return
	}

	sb.WriteString(fmt.Sprintf("## Recorded Pull Errors (%d)\n\n", len(summary.PullErrors)))

	for _, rec := range summary.PullErrors {
		sb.WriteString("- `" + rec.RepoSlug + "` — " + rec.ErrorType + ": " + oneLineSummary(rec.ErrorText) +
			" (" + rec.CreatedAt.UTC().Format(time.RFC3339) + ")\n")
	}

	sb.WriteString("\n")
}

// normalizeUtcTimestamp re-formats a stored timestamp as UTC RFC3339,
// returning the raw value unchanged when it cannot be parsed.
func normalizeUtcTimestamp(raw string) string {
	trimmed := strings.TrimSpace(raw)

	for _, layout := range []string{time.RFC3339, "2006-01-02T15:04:05", "2006-01-02 15:04:05", "2006-01-02"} {
		if parsed, parseErr := time.Parse(layout, trimmed); parseErr == nil {
			return parsed.UTC().Format(time.RFC3339)
		}
	}

	return trimmed
}

func hasTracebackMarkers(logs string) bool {
	lower := strings.ToLower(logs)

	for _, marker := range []string{"traceback", "--- fail", "assertionerror", "panic:", "fatal", "error:"} {
		if strings.Contains(lower, marker) {
			return true
		}
	}

	return false
}

func firstLineOr(primary, fallback string) string {
	if strings.TrimSpace(primary) != "" {
		return strings.TrimSpace(primary)
	}

	return oneLineSummary(fallback)
}

func oneLineSummary(text string) string {
	line := strings.TrimSpace(strings.SplitN(text, "\n", 2)[0])

	if line == "" {
		return "(no summary)"
	}

	return line
}
