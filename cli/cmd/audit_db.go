package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func completePendingCommandTask(db *store.DB, exitCode int, summary string) {
	pending, err := db.ListPendingTasks()
	if err != nil || len(pending) == 0 {
		return
	}

	latest := pending[len(pending)-1]
	if exitCode == 0 {
		_ = db.CompleteTask(latest.ID)
	} else {
		_ = db.FailTask(latest.ID, summary)
	}
}

func openAuditDB() (*store.DB, error) {
	prevQuiet := os.Getenv(constants.EnvGitMapQuiet)
	os.Setenv(constants.EnvGitMapQuiet, constants.EnvGitMapQuietTrue)
	defer os.Setenv(constants.EnvGitMapQuiet, prevQuiet)

	return initAndMigrateAuditDB()
}

func initAndMigrateAuditDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return nil, err
	}

	if err := db.Migrate(); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Audit DB migration failed: %v\n", err)
	}

	return db, nil
}

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

func resolveAlias(command string) string {
	return command
}
