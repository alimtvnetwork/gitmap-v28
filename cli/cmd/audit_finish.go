package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func finishCommandAudit(
	shouldAudit bool,
	id int64,
	start time.Time,
	exitCode int,
	summary string,
	repoCount int,
) {
	isNonAudit := !shouldAudit
	if isNonAudit {
		return
	}

	recordAuditEnd(id, start, exitCode, summary, repoCount)
}

func recordAuditEnd(id int64, start time.Time, exitCode int, summary string, repoCount int) {
	record := buildAuditCompletionRecord(id, start, exitCode, summary, repoCount)

	db, err := openAuditDB()
	if err != nil {
		return
	}

	defer db.Close()

	updateAuditRecordAndTasks(db, record, exitCode, summary)
	maybeRecordAiExecutionHistory(record, exitCode)
}

func maybeRecordAiExecutionHistory(record model.CommandHistoryRecord, exitCode int) {
	if os.Getenv("GITMAP_AI_TRACKING") != "1" {
		return
	}

	cwd, _ := os.Getwd()
	cmdStr := resolveAiCommandString()
	isSuccess := exitCode == 0
	_ = store.RecordAiExecution(
		"cli",
		cmdStr,
		"[]",
		cwd,
		"127.0.0.1",
		int(record.DurationMs),
		exitCode,
		record.Summary,
		"",
		isSuccess,
	)
}

func resolveAiCommandString() string {
	if len(os.Args) < 2 {
		return "gitmap"
	}
	return "gitmap " + strings.Join(os.Args[1:], " ")
}

func buildAuditCompletionRecord(id int64, start time.Time, exitCode int, summary string, repoCount int) model.CommandHistoryRecord {
	end := time.Now()

	return model.CommandHistoryRecord{
		ID:         id,
		FinishedAt: end.Format(time.RFC3339),
		DurationMs: end.Sub(start).Milliseconds(),
		ExitCode:   exitCode,
		Summary:    summary,
		RepoCount:  repoCount,
	}
}

func updateAuditRecordAndTasks(db *store.DB, record model.CommandHistoryRecord, exitCode int, summary string) {
	if err := db.UpdateHistory(record); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not update command history: %v\n", err)
	}

	completePendingCommandTask(db, exitCode, summary)
}
