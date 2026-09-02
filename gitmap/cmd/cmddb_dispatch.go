package cmd

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runDB routes the gitmap db command to appropriate sub-handlers.
func runDB(args []string) error {
	if len(args) == 0 {
		return runDBLs(nil)
	}
	sub := strings.ToLower(strings.TrimSpace(args[0]))
	return routeDBSubcommand(sub, args[1:])
}

func routeDBSubcommand(sub string, tail []string) error {
	switch sub {
	case "ls", "list":
		return runDBLs(tail)
	case "help", "-h", "--help":
		return runDBHelp()
	case "repo-db", "repodb":
		return runDBRepoDB(tail)
	case "sizes", "size":
		return runDBSizes(tail)
	case "reset", "clear":
		return runDBResetAction(tail)
	default:
		return handleUnknownDBSub(sub)
	}
}

func handleUnknownDBSub(sub string) error {
	fmt.Printf(constants.ColorRed+"Unknown db subcommand '%s'"+constants.ColorReset+"\n\n", sub)
	_ = runDBHelp()
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
