package cmddb

import (
	"fmt"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runDb routes the gitmap db command to appropriate sub-handlers.
func runDb(args []string) error {
	if len(args) == 0 {
		return runDBLs(nil)
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))
	return result.AsError(routeDbSubcommand(sub, args[1:]))
}

// Backwards-compatible alias for runDb
var runDB = runDb

func routeDbSubcommand(sub string, tail []string) result.ErrorWrapper {
	switch sub {
	case "status", "st", "info":
		return result.MatchWrapper(runDBStatus(tail))
	case "optimize", "opt":
		return result.MatchWrapper(runDBOptimize(tail))
	case "ls", "list":
		return result.MatchWrapper(runDBLs(tail))
	case "help", "-h", "--help":
		return result.MatchWrapper(runDbHelp())
	case "repo-db", "repodb":
		return result.MatchWrapper(runDBRepoDB(tail))
	case "sizes", "size":
		return result.MatchWrapper(runDBSizes(tail))
	case "reset":
		return result.MatchWrapper(runDbResetAction(tail))
	case "clear":
		return result.MatchWrapper(runDBClearAction(tail))
	default:
		return result.FailureWrapperErr(handleUnknownDbSub(sub))
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
