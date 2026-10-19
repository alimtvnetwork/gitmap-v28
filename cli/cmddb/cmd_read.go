package cmddb

import (
	"database/sql"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

type readFlags struct {
	aliasName   string
	limit       int
	isJson      bool
	showColumns bool
	isHelp      bool
}

// RunDbRead executes the unified gitmap db-read command.
func RunDbRead(args []string) *appfault.AppError {
	flags, posArgs := parseReadFlags(args)
	if flags.isHelp || (len(posArgs) == 0 && flags.aliasName == "") {
		printDbReadHelp(os.Stdout)
		return nil
	}

	return dispatchReadOperation(flags, posArgs)
}

func parseReadFlags(args []string) (readFlags, []string) {
	var flags readFlags
	var posArgs []string

	i := 0
	for i < len(args) {
		consumed, handled := parseFlagItem(args, i, &flags)
		if handled {
			i += consumed
			continue
		}

		posArgs = append(posArgs, args[i])
		i++
	}

	return flags, posArgs
}

func parseFlagItem(args []string, idx int, flags *readFlags) (int, bool) {
	arg := args[idx]
	switch {
	case arg == "--help" || arg == "-h":
		flags.isHelp = true
		return 1, true
	case arg == "--json" || arg == "-j":
		flags.isJson = true
		return 1, true
	case arg == "--columns" || arg == "-c":
		flags.showColumns = true
		return 1, true
	case (arg == "-as" || arg == "--alias") && idx+1 < len(args):
		flags.aliasName = args[idx+1]
		return 2, true
	case strings.HasPrefix(arg, "-as="):
		flags.aliasName = strings.TrimPrefix(arg, "-as=")
		return 1, true
	case strings.HasPrefix(arg, "--alias="):
		flags.aliasName = strings.TrimPrefix(arg, "--alias=")
		return 1, true
	case (arg == "-l" || arg == "--limit" || arg == "-n") && idx+1 < len(args):
		flags.limit, _ = strconv.Atoi(args[idx+1])
		return 2, true
	case strings.HasPrefix(arg, "--limit="):
		flags.limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "--limit="))
		return 1, true
	case strings.HasPrefix(arg, "-l="):
		flags.limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "-l="))
		return 1, true
	}
	return 0, false
}

func dispatchReadOperation(flags readFlags, posArgs []string) *appfault.AppError {
	if len(posArgs) == 1 && flags.aliasName != "" && !strings.EqualFold(posArgs[0], "ls") {
		return handleAliasRegistrationOnly(flags, posArgs[0])
	}
	if len(posArgs) > 0 && strings.EqualFold(posArgs[0], "ls") {
		return handleReadLs(flags, resolveEffectiveTargetDb(flags.aliasName, ""))
	}
	if len(posArgs) > 1 && strings.EqualFold(posArgs[1], "ls") {
		maybeAutoRegisterAlias(flags, posArgs[0])
		return handleReadLs(flags, posArgs[0])
	}
	if len(posArgs) > 0 && strings.EqualFold(posArgs[0], "table") {
		return handleReadTable(flags, resolveEffectiveTargetDb(flags.aliasName, ""), posArgs[1:])
	}
	if len(posArgs) > 1 && strings.EqualFold(posArgs[1], "table") {
		maybeAutoRegisterAlias(flags, posArgs[0])
		return handleReadTable(flags, posArgs[0], posArgs[2:])
	}
	if len(posArgs) >= 2 {
		maybeAutoRegisterAlias(flags, posArgs[0])
		return handleReadSql(flags, posArgs[0], strings.Join(posArgs[1:], " "))
	}
	if len(posArgs) == 1 && isQuerySafe(posArgs[0]) {
		return handleReadSql(flags, resolveEffectiveTargetDb(flags.aliasName, ""), posArgs[0])
	}
	if len(posArgs) == 1 {
		maybeAutoRegisterAlias(flags, posArgs[0])
		return handleReadLs(flags, posArgs[0])
	}
	return handleReadLs(flags, "master")
}

func resolveEffectiveTargetDb(flagAlias, fallback string) string {
	if flagAlias != "" {
		return flagAlias
	}
	if fallback != "" {
		return fallback
	}
	return "master"
}

func maybeAutoRegisterAlias(flags readFlags, dbArg string) {
	if flags.aliasName == "" || dbArg == "" {
		return
	}
	if _, statErr := os.Stat(dbArg); statErr == nil {
		_, _ = registerDbPathAlias(dbArg, flags.aliasName)
	}
}

func handleAliasRegistrationOnly(flags readFlags, rawPath string) *appfault.AppError {
	start := time.Now()
	absPath, err := registerDbPathAlias(rawPath, flags.aliasName)
	elapsed := time.Since(start)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:alias", rawPath, flags.aliasName, "ERR_DB_NOT_FOUND", 20, err.Error(), "Ensure target SQLite file exists before registering.")
	}

	msg := fmt.Sprintf("Registered alias '%s' -> '%s' in database registry.", flags.aliasName, absPath)
	if flags.isJson {
		data := map[string]any{"alias": flags.aliasName, "path": absPath, "message": msg}
		outputReadSuccess(os.Stdout, "db-read:alias", absPath, flags.aliasName, data, elapsed)
		return nil
	}

	fmt.Println(msg)
	fmt.Println(formatCompletionTime(elapsed))
	return nil
}

func handleReadLs(flags readFlags, dbArg string) *appfault.AppError {
	start := time.Now()
	conn, dbPath, err := openReadOnlyResolvedDb(dbArg)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:ls", dbArg, dbArg, "ERR_DB_NOT_FOUND", 20, err.Error(), "Verify database alias or run 'gitmap scan' to initialize.")
	}
	defer conn.Close()

	tables, listErr := listTablesAndViews(conn)
	elapsed := time.Since(start)
	if listErr != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:ls", dbPath, dbArg, "ERR_EXECUTION", 23, listErr.Error(), "Check database schema integrity.")
	}

	if flags.isJson {
		data := map[string]any{"tables": tables, "totalTables": len(tables)}
		outputReadSuccess(os.Stdout, "db-read:ls", dbPath, dbArg, data, elapsed)
		return nil
	}

	renderTablesTerminal(os.Stdout, dbPath, tables)
	fmt.Println(formatCompletionTime(elapsed))
	return nil
}

func handleReadTable(flags readFlags, dbArg string, tables []string) *appfault.AppError {
	if len(tables) == 0 {
		return failRead(os.Stdout, flags.isJson, "db-read:table", dbArg, dbArg, "ERR_TABLE_NOT_FOUND", 22, "table name is required (e.g. gitmap db-read <db> table <t1>)", "Specify at least one table name.")
	}

	start := time.Now()
	conn, dbPath, err := openReadOnlyResolvedDb(dbArg)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read:table", dbArg, dbArg, "ERR_DB_NOT_FOUND", 20, err.Error(), "Verify database alias or path.")
	}
	defer conn.Close()

	if flags.showColumns {
		return handleTableColumnsOutput(flags, conn, dbPath, dbArg, tables, start)
	}

	return handleTableOverviewOutput(flags, conn, dbPath, dbArg, tables, start)
}

func handleTableColumnsOutput(flags readFlags, conn *sql.DB, dbPath, dbArg string, tables []string, start time.Time) *appfault.AppError {
	var columnOverviews []TableColumnOverview
	for _, t := range tables {
		cols, err := fetchColumnInfos(conn, t)
		if err != nil || len(cols) == 0 {
			msg := fmt.Sprintf("table '%s' not found or has no columns in %s", t, dbPath)
			return failRead(os.Stdout, flags.isJson, "db-read:table", dbPath, dbArg, "ERR_TABLE_NOT_FOUND", 22, msg, "Verify table name via 'gitmap db-read ls'.")
		}
		columnOverviews = append(columnOverviews, TableColumnOverview{Table: t, Columns: cols})
	}

	elapsed := time.Since(start)
	if flags.isJson {
		outputReadSuccess(os.Stdout, "db-read:table", dbPath, dbArg, map[string]any{"tables": columnOverviews}, elapsed)
		return nil
	}

	for _, co := range columnOverviews {
		renderColumnsTerminal(os.Stdout, dbPath, co.Table, co.Columns)
		fmt.Println()
	}
	fmt.Println(formatCompletionTime(elapsed))
	return nil
}

func handleTableOverviewOutput(flags readFlags, conn *sql.DB, dbPath, dbArg string, tables []string, start time.Time) *appfault.AppError {
	var overviews []TableOverview
	for _, t := range tables {
		overview, err := fetchTableOverview(conn, t, dbPath, flags.limit)
		if err != nil {
			return failRead(os.Stdout, flags.isJson, "db-read:table", dbPath, dbArg, "ERR_TABLE_NOT_FOUND", 22, err.Error(), "Verify table name via 'gitmap db-read ls'.")
		}
		overviews = append(overviews, overview)
	}

	elapsed := time.Since(start)
	if flags.isJson {
		outputReadSuccess(os.Stdout, "db-read:table", dbPath, dbArg, map[string]any{"tables": overviews}, elapsed)
		return nil
	}

	for _, ov := range overviews {
		renderTableOverviewTerminal(os.Stdout, ov)
		fmt.Println()
	}
	fmt.Println(formatCompletionTime(elapsed))
	return nil
}

func handleReadSql(flags readFlags, dbArg, sqlStr string) *appfault.AppError {
	cleanSql := strings.TrimSpace(sqlStr)
	if cleanSql == "" {
		return failRead(os.Stdout, flags.isJson, "db-read:sql", dbArg, dbArg, "ERR_INVALID_SQL", 23, "SQL query string is required", "Provide a non-empty SQL query string.")
	}

	if !isQuerySafe(cleanSql) {
		msg := "mutation attempted during db-read: attempt to write a readonly database"
		return failRead(os.Stdout, flags.isJson, "db-read", dbArg, dbArg, "ERR_READONLY_VIOLATION", 24, msg, "Use 'gitmap db-add' or 'gitmap db-edit' for write operations.")
	}

	start := time.Now()
	conn, dbPath, err := openReadOnlyResolvedDb(dbArg)
	if err != nil {
		return failRead(os.Stdout, flags.isJson, "db-read", dbArg, dbArg, "ERR_DB_NOT_FOUND", 20, err.Error(), "Verify database alias or path.")
	}
	defer conn.Close()

	effectiveSql := applyQueryLimit(cleanSql, flags.limit)
	res, queryErr := executeDynamicQuery(conn, effectiveSql)
	elapsed := time.Since(start)
	if queryErr != nil {
		return failRead(os.Stdout, flags.isJson, "db-read", dbPath, dbArg, "ERR_INVALID_SQL", 23, queryErr.Error(), "Verify SQL syntax against standard SQLite grammar.")
	}

	if flags.isJson {
		data := map[string]any{"columns": res.Columns, "rows": res.Rows, "count": res.Count}
		outputReadSuccess(os.Stdout, "db-read", dbPath, dbArg, data, elapsed)
		return nil
	}

	renderQueryTerminal(os.Stdout, res)
	fmt.Println(formatCompletionTime(elapsed))
	return nil
}

const dbReadHelp = `gitmap db-read — Read and inspect SQLite databases across Split-DB tiers

Usage:
  gitmap db-read <path-or-alias> "sql" [flags]
  gitmap db-read <path-or-alias> ls [flags]
  gitmap db-read <path-or-alias> table <t1> [t2...] [flags]
  gitmap db-read <path> -as <aliasName> [flags]

Flags:
  -as, --alias <name>   Database alias or path (default: master)
  -l,  --limit <N>      Maximum rows to return for query (default: 100)
  -c,  --columns        List columns, types, notnull, default, and PK
  -j,  --json           Output results in structured JSON envelope
  -h,  --help           Show this help text
`

func printDbReadHelp(out io.Writer) {
	fmt.Fprint(out, dbReadHelp)
}
