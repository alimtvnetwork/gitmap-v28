package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDb routes the gitmap db command to appropriate sub-handlers.
func runDb(args []string) error {
	if len(args) == 0 {
		return runDBLs(nil)
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))

	return routeDbSubcommand(sub, args[1:])
}

// Backwards-compatible alias for runDb
var runDB = runDb

func routeDbSubcommand(sub string, tail []string) error {
	switch sub {
	case "status", "st", "info":
		return runDBStatus(tail)
	case "optimize", "opt":
		return runDBOptimize(tail)
	case "ls", "list":
		return runDBLs(tail)
	case "help", "-h", "--help":
		return runDbHelp()
	case "repo-db", "repodb":
		return runDBRepoDB(tail)
	case "sizes", "size":
		return runDBSizes(tail)
	case "reset":
		return runDbResetAction(tail)
	case "clear":
		return runDBClearAction(tail)
	default:
		return handleUnknownDbSub(sub)
	}
}

// Backwards-compatible alias for routeDbSubcommand

func handleUnknownDbSub(sub string) error {
	fmt.Printf(constants.ColorRed+"Unknown db subcommand '%s'"+constants.ColorReset+"\n\n", sub)
	_ = runDbHelp()

	return apperror.NewWithDetails(
		"cmd.db.dispatch",
		"E1050",
		fmt.Sprintf("unknown db subcommand '%s'", sub),
		"cmd.db",
		apperror.ErrorTypeValidation,
		apperror.SeverityError,
		map[string]any{"subcommand": sub},
	)
}

// Backwards-compatible alias for handleUnknownDbSub
