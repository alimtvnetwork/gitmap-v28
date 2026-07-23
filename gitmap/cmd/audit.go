package cmd

import (
	"fmt"
	"os"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/store"
)

func beginCommandAudit(command string, args []string) (int64, time.Time, bool) {
	start := time.Now()
	if !shouldAuditCommand(command) {
		return 0, start, false
	}

	return recordAuditStart(command, args)
}

func finishCommandAudit(shouldAudit bool, id int64, start time.Time, exitCode int, summary string, repoCount int) {
	if !shouldAudit {
		return
	}

	recordAuditEnd(id, start, exitCode, summary, repoCount)
}

func shouldAuditCommand(command string) bool {
	return command != constants.CmdVersion && command != constants.CmdVersionAlias
}

// recordAuditStart inserts a new history record at command start.
func recordAuditStart(command string, args []string) (int64, time.Time, bool) {
	start := time.Now()
	alias, flags, positional := classifyArgs(command, args)

	record := model.CommandHistoryRecord{
		Command:   command,
		Alias:     alias,
		Args:      positional,
		Flags:     flags,
		StartedAt: start.Format(time.RFC3339),
	}

	db, err := openAuditDB()
	if err != nil {
		return 0, start, false
	}
	defer db.Close()

	id, insertErr := db.InsertHistory(record)
	if insertErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record command history: %v\n", insertErr)
	}

	return id, start, true
}

// recordAuditEnd updates a history record with completion details.
func recordAuditEnd(id int64, start time.Time, exitCode int, summary string, repoCount int) {
	end := time.Now()
	duration := end.Sub(start).Milliseconds()

	record := model.CommandHistoryRecord{
		ID:         id,
		FinishedAt: end.Format(time.RFC3339),
		DurationMs: duration,
		ExitCode:   exitCode,
		Summary:    summary,
		RepoCount:  repoCount,
	}

	db, err := openAuditDB()
	if err != nil {
		return
	}
	defer db.Close()

	if err := db.UpdateHistory(record); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not update command history: %v\n", err)
	}
}

// openAuditDB opens the database silently (no error output).
func openAuditDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Audit DB migration failed: %v\n", err)
	}

	return db, nil
}

// classifyArgs separates flags from positional arguments.
func classifyArgs(command string, args []string) (string, string, string) {
	alias := resolveAlias(command)
	var flags, positional []string

	for _, arg := range args {
		if fmt.Sprintf("%c", arg[0]) == "-" {
			flags = append(flags, arg)
		} else {
			positional = append(positional, arg)
		}
	}

	return alias, joinStrings(flags), joinStrings(positional)
}

// joinStrings joins string slices with spaces.
func joinStrings(s []string) string {
	result := ""
	for i, v := range s {
		if i > 0 {
			result += " "
		}
		result += v
	}

	return result
}

// resolveAlias returns the alias if the command was invoked by alias.
func resolveAlias(command string) string {
	return command
}
