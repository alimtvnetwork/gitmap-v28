package cmd

import (
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func beginCommandAudit(command string, args []string) (int64, time.Time, bool) {
	start := time.Now()
	isNonAuditCommand := !isAuditableCommand(command)
	if isNonAuditCommand {
		return 0, start, false
	}

	return recordAuditStart(command, args)
}

func isAuditableCommand(command string) bool {
	if command == constants.CmdVersion || command == constants.CmdVersionAlias {
		return false
	}

	if command == constants.CmdReset || command == constants.CmdDBReset || command == "db-reset" {
		return false
	}

	return true
}

func recordAuditStart(command string, args []string) (int64, time.Time, bool) {
	start := time.Now()
	id, isRecorded := executeAuditDBInsert(command, args, start)
	if !isRecorded {
		return 0, start, false
	}

	recordPendingTaskAudit(command, args)

	return id, start, true
}

func executeAuditDBInsert(command string, args []string, start time.Time) (int64, bool) {
	record := buildCommandAuditRecord(command, args, start)
	db, err := openAuditDB()
	if err != nil {
		return 0, false
	}

	defer db.Close()

	return insertAuditRecord(db, record), true
}

func buildCommandAuditRecord(command string, args []string, start time.Time) model.CommandHistoryRecord {
	alias, flags, positional := classifyArgs(command, args)

	return model.CommandHistoryRecord{
		Command:   command,
		Alias:     alias,
		Args:      positional,
		Flags:     flags,
		StartedAt: start.Format(time.RFC3339),
	}
}

func insertAuditRecord(db *store.DB, record model.CommandHistoryRecord) int64 {
	id, insertErr := db.InsertHistory(record)
	if insertErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not record command history: %v\n", insertErr)
	}

	return id
}

func recordPendingTaskAudit(command string, args []string) {
	if !isPendingTaskCommand(command) {
		return
	}
	cwd, _ := os.Getwd()
	cmdArgs := strings.Join(args, " ")
	_, _ = createPendingTask(command, cwd, cwd, command, cmdArgs)
}

func isPendingTaskCommand(command string) bool {
	switch command {
	case "pending", "status", "st", "version", "--version", "-v", "help", "--help", "-h", "docs":
		return false
	case "storage", "disk", "df", "ps", "ls":
		return false
	}

	return true
}
