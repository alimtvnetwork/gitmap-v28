package cmddb

import (
	"fmt"
	"io"
	"os"
	"slices"
	"strconv"
	"strings"
	"time"

	appfault "github.com/alimtvnetwork/gitmap-v28/cli/appfault"
)

type readFlags struct {
	aliasName     string
	limit         int
	isJson        bool
	isShowColumns bool
	isHelp        bool
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
		consumed, isHandled := parseFlagItem(args, i, &flags)
		if isHandled {
			i += consumed
			continue
		}

		posArgs = append(posArgs, args[i])
		i++
	}

	return flags, posArgs
}

func parseFlagItem(args []string, idx int, flags *readFlags) (int, bool) {
	if n, isHandled := parseSimpleReadFlag(args[idx], flags); isHandled {
		return n, true
	}

	if n, isHandled := parseAliasReadFlag(args, idx, flags); isHandled {
		return n, true
	}

	return parseLimitReadFlag(args, idx, flags)
}

func parseSimpleReadFlag(arg string, flags *readFlags) (int, bool) {
	switch arg {
	case "--help", "-h":
		flags.isHelp = true
		return 1, true
	case "--json", "-j":
		flags.isJson = true
		return 1, true
	case "--columns", "-c":
		flags.isShowColumns = true
		return 1, true
	default:
		return 0, false
	}
}

func parseAliasReadFlag(args []string, idx int, flags *readFlags) (int, bool) {
	arg := args[idx]
	if slices.Contains([]string{"-as", "--alias"}, arg) && idx+1 < len(args) {
		flags.aliasName = args[idx+1]
		return 2, true
	}

	if strings.HasPrefix(arg, "-as=") {
		flags.aliasName = strings.TrimPrefix(arg, "-as=")
		return 1, true
	}

	if strings.HasPrefix(arg, "--alias=") {
		flags.aliasName = strings.TrimPrefix(arg, "--alias=")
		return 1, true
	}

	return 0, false
}

func parseLimitReadFlag(args []string, idx int, flags *readFlags) (int, bool) {
	arg := args[idx]
	if slices.Contains([]string{"-l", "--limit", "-n"}, arg) && idx+1 < len(args) {
		flags.limit, _ = strconv.Atoi(args[idx+1])
		return 2, true
	}

	if strings.HasPrefix(arg, "--limit=") {
		flags.limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "--limit="))
		return 1, true
	}

	if strings.HasPrefix(arg, "-l=") {
		flags.limit, _ = strconv.Atoi(strings.TrimPrefix(arg, "-l="))
		return 1, true
	}

	return 0, false
}

func dispatchReadOperation(flags readFlags, posArgs []string) *appfault.AppError {
	if len(posArgs) == 0 {
		return handleReadInspect(flags, resolveEffectiveTargetDb(flags.aliasName, ""))
	}

	if isExplicitLsDispatch(posArgs) {
		return dispatchLsOperation(flags, posArgs)
	}

	if isExplicitTableDispatch(posArgs) {
		return dispatchTableOperation(flags, posArgs)
	}

	if len(posArgs) >= 2 {
		return dispatchMultiArgOperation(flags, posArgs)
	}

	return dispatchSingleArgOperation(flags, posArgs[0])
}

func isExplicitLsDispatch(posArgs []string) bool {
	if len(posArgs) > 0 && strings.EqualFold(posArgs[0], "ls") {
		return true
	}

	return len(posArgs) > 1 && strings.EqualFold(posArgs[1], "ls")
}

func dispatchLsOperation(flags readFlags, posArgs []string) *appfault.AppError {
	if strings.EqualFold(posArgs[0], "ls") {
		return handleReadLs(flags, resolveEffectiveTargetDb(flags.aliasName, ""))
	}

	maybeAutoRegisterAlias(flags, posArgs[0])

	return handleReadLs(flags, posArgs[0])
}

func isExplicitTableDispatch(posArgs []string) bool {
	if len(posArgs) > 0 && strings.EqualFold(posArgs[0], "table") {
		return true
	}

	return len(posArgs) > 1 && strings.EqualFold(posArgs[1], "table")
}

func dispatchTableOperation(flags readFlags, posArgs []string) *appfault.AppError {
	if strings.EqualFold(posArgs[0], "table") {
		return handleReadTable(flags, resolveEffectiveTargetDb(flags.aliasName, ""), posArgs[1:])
	}

	maybeAutoRegisterAlias(flags, posArgs[0])

	return handleReadTable(flags, posArgs[0], posArgs[2:])
}

func dispatchMultiArgOperation(flags readFlags, posArgs []string) *appfault.AppError {
	dbArg := posArgs[0]
	maybeAutoRegisterAlias(flags, dbArg)
	queryCandidate := strings.Join(posArgs[1:], " ")
	if isQuerySafe(queryCandidate) || isSqlStatement(queryCandidate) {
		return handleReadSql(flags, dbArg, queryCandidate)
	}

	return handleReadTable(flags, dbArg, posArgs[1:])
}

func dispatchSingleArgOperation(flags readFlags, arg string) *appfault.AppError {
	if isQuerySafe(arg) {
		return handleReadSql(flags, resolveEffectiveTargetDb(flags.aliasName, ""), arg)
	}

	if flags.aliasName != "" {
		return handleSinglePosArgWithAlias(flags, arg)
	}

	if isDatabaseTarget(arg) {
		maybeAutoRegisterAlias(flags, arg)

		return handleReadInspect(flags, arg)
	}

	return handleReadTable(flags, resolveEffectiveTargetDb(flags.aliasName, ""), []string{arg})
}

func handleSinglePosArgWithAlias(flags readFlags, target string) *appfault.AppError {
	if isExistingFilePath(target) {
		return handleAliasRegistrationOnly(flags, target)
	}

	return handleReadTable(flags, flags.aliasName, []string{target})
}

func isExistingFilePath(target string) bool {
	fi, err := os.Stat(target)
	if err != nil {
		return false
	}

	return !fi.IsDir()
}

func isDatabaseTarget(target string) bool {
	if isExistingFilePath(target) {
		return true
	}

	_, hasSys := resolveSystemDbAlias(target)
	if hasSys {
		return true
	}

	_, hasReg := lookupAliasInRegistry(target)

	return hasReg
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

	return renderAliasRegistrationResult(flags, absPath, elapsed)
}

func renderAliasRegistrationResult(flags readFlags, absPath string, elapsed time.Duration) *appfault.AppError {
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

	return outputLsResult(flags, dbPath, dbArg, tables, elapsed)
}

func outputLsResult(flags readFlags, dbPath, dbArg string, tables []TableInfo, elapsed time.Duration) *appfault.AppError {
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

	if flags.isShowColumns {
		return handleTableColumnsOutput(flags, conn, dbPath, dbArg, tables, start)
	}

	return handleTableOverviewOutput(flags, conn, dbPath, dbArg, tables, start)
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

	return executeReadOnlySql(flags, dbArg, cleanSql)
}

func executeReadOnlySql(flags readFlags, dbArg, cleanSql string) *appfault.AppError {
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

	return outputSqlResult(flags, dbPath, dbArg, res, elapsed)
}

func outputSqlResult(flags readFlags, dbPath, dbArg string, res QueryResult, elapsed time.Duration) *appfault.AppError {
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
  gitmap db-read <path-or-alias> [table] <t1> [t2...] [flags]
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
