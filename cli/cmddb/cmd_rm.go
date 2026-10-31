package cmddb

import (
	"database/sql"
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const dbRmHelp = `gitmap db-rm (dbrm, db-delete) — Remove a database alias and optionally purge physical SQLite file.

Usage:
  gitmap db-rm <alias> [--purge] [-y]

Flags:
  --purge, --purge-file   Permanently delete the underlying database file from disk
  -y, --yes, --force      Bypass confirmation prompt
  -h, --help              Show help text

Examples:
  gitmap db-rm customdb
  gitmap db-rm tempdb --purge -y
  gitmap db-delete olddb --purge
`

type dbRmParams struct {
	aliasName   string
	isPurge     bool
	isConfirmed bool
}

// RunDbRm executes database alias unbinding and optional file purge.
func RunDbRm(args []string) *appfault.AppError {
	if isDbHelpRequested(args) {
		fmt.Print(dbRmHelp)

		return nil
	}

	params, err := parseDbRmArgs(args)
	if err != nil {
		return err
	}

	return processAliasRemoval(params)
}

func processAliasRemoval(params dbRmParams) *appfault.AppError {
	if isProtectedSystemAlias(params.aliasName) {
		return createForbiddenAliasError(params.aliasName)
	}

	resolvedPath, findErr := findAliasPathInRegistry(params.aliasName)
	if findErr != nil {
		return findErr
	}

	if params.isPurge && !params.isConfirmed && !askPurgeConfirmation(resolvedPath) {
		fmt.Println("[ABORTED] Database file purge canceled.")

		return nil
	}

	return executeAliasRemoval(params, resolvedPath)
}

func isProtectedSystemAlias(alias string) bool {
	lower := strings.ToLower(strings.TrimSpace(alias))
	switch lower {
	case "master", "gitmap", "pipeline", "pipelines", "heatmap", "tasks", "agents", "agy", "antigravity", "convs", "conversation-summaries":
		return true
	default:
		return false
	}
}

func createForbiddenAliasError(alias string) *appfault.AppError {
	return &appfault.AppError{
		Op:       "db.rm",
		Code:     "E4030",
		Type:     apperror.ErrorTypePrecondition,
		Severity: apperror.SeverityError,
		Message:  fmt.Sprintf("system database alias '%s' cannot be removed or deleted", alias),
		Caller:   apperror.CaptureCaller(apperror.DefaultCallerSkip),
	}
}

func parseDbRmArgs(args []string) (dbRmParams, *appfault.AppError) {
	var p dbRmParams
	for _, arg := range args {
		consumeDbRmArg(arg, &p)
	}

	if p.aliasName == "" {
		return p, appfault.NewValidationError("usage: gitmap db-rm <alias> [--purge] [-y]")
	}

	return p, nil
}

func consumeDbRmArg(arg string, p *dbRmParams) {
	if slices.Contains([]string{"--purge", "--purge-file"}, arg) {
		p.isPurge = true
		return
	}

	if isConfirmFlag(arg) {
		p.isConfirmed = true
		return
	}

	if !strings.HasPrefix(arg, "-") && p.aliasName == "" {
		p.aliasName = strings.TrimSpace(arg)
	}
}

func findAliasPathInRegistry(alias string) (string, *appfault.AppError) {
	openRes := store.OpenSQLiteDB(store.DefaultDBPath())
	if openRes.IsFailure() {
		return "", appfault.NewNotFoundError(fmt.Sprintf("open master db: %v", openRes.AppError()))
	}
	conn := openRes.Value
	defer conn.Close()

	return queryAliasPath(conn, alias)
}

func queryAliasPath(conn *sql.DB, alias string) (string, *appfault.AppError) {
	var path string
	query := `SELECT DatabasePath FROM DatabaseAlias WHERE LOWER(AliasName) = LOWER(?) LIMIT 1`
	err := conn.QueryRow(query, alias).Scan(&path)
	if err == sql.ErrNoRows {
		return "", appfault.NewNotFoundError(fmt.Sprintf("database alias '%s' is not registered", alias))
	}

	if err != nil {
		return "", appfault.NewValidationError(fmt.Sprintf("query alias '%s': %v", alias, err))
	}

	return path, nil
}

func askPurgeConfirmation(path string) bool {
	fmt.Printf("[WARNING] Permanently delete database file '%s'? [y/N]: ", path)
	var response string
	_, err := fmt.Scanln(&response)
	if err != nil {
		return false
	}

	trimmed := strings.ToLower(strings.TrimSpace(response))

	return trimmed == "y" || trimmed == "yes"
}

func executeAliasRemoval(p dbRmParams, resolvedPath string) *appfault.AppError {
	delErr := deleteAliasFromDatabase(p.aliasName)
	if delErr != nil {
		return delErr
	}

	InvalidateTableMetadata(resolvedPath)

	return finishAliasRemoval(p, resolvedPath)
}

func deleteAliasFromDatabase(aliasName string) *appfault.AppError {
	openRes := store.OpenSQLiteDB(store.DefaultDBPath())
	if openRes.IsFailure() {
		return appfault.NewValidationError(fmt.Sprintf("open master db: %v", openRes.AppError()))
	}
	conn := openRes.Value
	defer conn.Close()

	delQuery := `DELETE FROM DatabaseAlias WHERE LOWER(AliasName) = LOWER(?)`
	_, delErr := conn.Exec(delQuery, aliasName)
	if delErr != nil {
		return appfault.NewValidationError(fmt.Sprintf("delete alias '%s': %v", aliasName, delErr))
	}

	return nil
}

func finishAliasRemoval(p dbRmParams, resolvedPath string) *appfault.AppError {
	if p.isPurge {
		removeDatabasePhysicalFiles(resolvedPath)
		fmt.Printf("[SUCCESS] Purged database file '%s' and removed alias '%s' from registry.\n", resolvedPath, p.aliasName)

		return nil
	}

	fmt.Printf("[SUCCESS] Removed alias '%s' from registry (file preserved: %s).\n", p.aliasName, resolvedPath)

	return nil
}

func removeDatabasePhysicalFiles(dbPath string) {
	_ = os.Remove(dbPath)
	_ = os.Remove(dbPath + "-wal")
	_ = os.Remove(dbPath + "-shm")
	_ = os.Remove(dbPath + "-journal")
}
