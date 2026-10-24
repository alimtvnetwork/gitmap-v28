package cmddb

import (
	"context"
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// runDb routes the gitmap db command to appropriate sub-handlers.
func runDb(args []string) error {
	if len(args) == 0 {
		return runDBLs(nil)
	}

	sub := strings.ToLower(strings.TrimSpace(args[0]))
	switch sub {
	case "tables", "table", "t", "schema", "ddl", "s", "columns", "cols", "column", "c", "query", "q", "sql", "records", "rows":
		return DispatchDbSubsystem(context.Background(), args, os.Stdout)
	case "status", "st", "info":
		return result.AsError(result.MatchWrapper(runDBStatus(args[1:])))
	case "optimize", "opt":
		return result.AsError(result.MatchWrapper(runDBOptimize(args[1:])))
	case "vacuum", "vac":
		return result.AsError(result.MatchWrapper(runDBVacuum(args[1:])))
	case "ls", "list":
		if len(args) > 1 {
			return DispatchDbSubsystem(context.Background(), args, os.Stdout)
		}
		return result.AsError(result.MatchWrapper(runDBLs(args[1:])))
	case "help", "-h", "--help":
		return printDbHelp(os.Stdout)
	case "repo-db", "repodb":
		return result.AsError(result.MatchWrapper(runDBRepoDB(args[1:])))
	case "sizes", "size":
		return result.AsError(result.MatchWrapper(runDBSizes(args[1:])))
	case "doctor", "doc":
		return RunDbDoctor(args[1:])
	case "listing", "dbl":
		return RunDbListing(args[1:])
	case "read", "dbr", "db-read":
		return RunDbRead(args[1:])
	case "add", "db-add":
		return RunDbAdd(args[1:])
	case "edit", "db-edit":
		return RunDbEdit(args[1:])
	case "find", "db-find":
		return RunDbFind(args[1:])
	case "view", "db-view":
		return RunDbView(args[1:])
	case "reset":
		return result.AsError(result.MatchWrapper(runDbResetAction(args[1:])))
	case "clear":
		return result.AsError(result.MatchWrapper(runDBClearAction(args[1:])))
	default:
		return DispatchDbSubsystem(context.Background(), args, os.Stdout)
	}
}

// DispatchDbSubsystem routes SQLite database introspection commands.
func DispatchDbSubsystem(ctx context.Context, args []string, out io.Writer) error {
	subcmd, pos, opts := parseDbSubsystemArgs(args)
	switch strings.ToLower(subcmd) {
	case "tables", "table", "t":
		dbArg := opts.DbPath
		if dbArg == "" && len(pos) > 0 {
			dbArg = pos[0]
		}
		return RunDbTables(dbArg, opts.IsJson, out)
	case "schema", "ddl", "s":
		tableArg := ""
		if len(pos) > 0 {
			tableArg = pos[0]
		}
		return RunDbSchema(tableArg, opts.DbPath, opts.IsJson, out)
	case "columns", "cols", "column", "c":
		tableArg := ""
		if len(pos) > 0 {
			tableArg = pos[0]
		}
		return RunDbColumns(tableArg, opts.DbPath, opts.IsJson, out)
	case "query", "q", "sql":
		queryArg := ""
		if len(pos) > 0 {
			queryArg = strings.Join(pos, " ")
		}
		return RunDbQuery(queryArg, opts.DbPath, opts.IsJson, out)
	case "list", "records", "rows", "ls", "l":
		tableArg := ""
		if len(pos) > 0 {
			tableArg = pos[0]
		}
		return RunDbList(tableArg, opts.DbPath, opts.Limit, opts.IsJson, out)
	case "add":
		return RunDbAdd(pos)
	case "edit":
		return RunDbEdit(pos)
	case "find":
		return RunDbFind(pos)
	case "view":
		return RunDbView(pos)
	case "help", "--help", "-h", "":
		return printDbHelp(out)
	default:
		return apperror.NewValidationError(fmt.Sprintf("unknown db subcommand '%s'. Run 'gitmap db help' for usage", subcmd))
	}
}

func parseDbSubsystemArgs(args []string) (string, []string, CommonDbOptions) {
	if len(args) == 0 {
		return "help", nil, CommonDbOptions{}
	}

	var (
		subcmd  string
		posArgs []string
		opts    CommonDbOptions
	)

	i := 0
	for i < len(args) {
		arg := args[i]
		if subcmd == "" && !strings.HasPrefix(arg, "-") {
			subcmd = arg
			i++
			continue
		}

		consumed, handled := parseFlagToken(args, i, &opts)
		if handled {
			i += consumed
			continue
		}

		posArgs = append(posArgs, arg)
		i++
	}

	if subcmd == "" {
		subcmd = "help"
	}

	return subcmd, posArgs, opts
}

func parseFlagToken(args []string, idx int, opts *CommonDbOptions) (int, bool) {
	arg := args[idx]
	if slices.Contains([]string{"--json", "-j"}, arg) {
		opts.IsJson = true
		return 1, true
	}

	if slices.Contains([]string{"--db", "-d"}, arg) && idx+1 < len(args) {
		opts.DbPath = args[idx+1]
		return 2, true
	}

	if strings.HasPrefix(arg, "--db=") {
		opts.DbPath = strings.TrimPrefix(arg, "--db=")
		return 1, true
	}

	if slices.Contains([]string{"--limit", "-n", "-l"}, arg) && idx+1 < len(args) {
		opts.Limit, _ = strconv.Atoi(args[idx+1])
		return 2, true
	}

	if strings.HasPrefix(arg, "--limit=") {
		opts.Limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "--limit="))
		return 1, true
	}

	return 0, false
}

func printDbHelp(out io.Writer) error {
	help := `GitMap SQLite Database Introspection Subsystem

Usage:
  gitmap db <command> [arguments...] [flags]

Aliases:
  database

Available Commands:
  tables  [db]            List tables, views, column counts, and record counts
  schema  [table]         Dump CREATE TABLE and index DDL
  columns <table>         Inspect column definitions (PRAGMA table_info)
  query   "<SQL>"         Execute read-only SQL queries (SELECT, PRAGMA, EXPLAIN, WITH)
  list    <table>         Display table rows with limit paging
  listing [flags]         List registered and discovered SQLite databases (dbl)
  help                    Show this help text

Flags:
  -d, --db <path|alias>   Target database: 'master'/'gitmap', 'heatmap', 'pipeline', 'tasks', or path
  -j, --json              Output results in JSON format
  -n, --limit <N>         Maximum rows to return for 'list' (default 20)

Examples:
  gitmap db tables
  gitmap db tables heatmap
  gitmap db schema Subtask --db tasks
  gitmap db columns Subtask --json
  gitmap db query "SELECT SubtaskId, Status FROM Subtask" --db tasks
  gitmap db list Subtask --limit 10
`
	_, err := fmt.Fprint(out, help)
	return err
}
