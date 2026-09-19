package cmdautomation

import (
	"database/sql"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunDbMigrate manages and executes SQLite schema migrations.
func RunDbMigrate(opts DbMigrateOptions) DbMigrateResultMonad {
	start := time.Now()
	dbPath := resolveDbPath(opts.DbPath)
	db, errApp := openMigrationDb(dbPath)
	if errApp != nil {
		return result.Fail[DbMigrateResult](errApp)
	}
	defer db.Close()
	res, runErr := executeMigrationPlan(db, opts)
	if runErr != nil {
		return result.Fail[DbMigrateResult](runErr)
	}
	res.Duration = time.Since(start)
	res.IsDryRun = opts.IsDryRun
	res.IsSuccess = true
	return result.Ok(res)
}

func resolveDbPath(dbPath string) string {
	if dbPath != "" {
		return dbPath
	}
	return GetAutomationDbPath()
}

func openMigrationDb(dbPath string) (*sql.DB, *apperror.AppError) {
	dir := filepath.Dir(dbPath)
	_ = os.MkdirAll(dir, 0755)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open migration sqlite database")
	}
	db.SetMaxOpenConns(1)
	if errInit := initMigrationTable(db); errInit != nil {
		_ = db.Close()
		return nil, errInit
	}
	return db, nil
}

func initMigrationTable(db *sql.DB) *apperror.AppError {
	schema := `CREATE TABLE IF NOT EXISTS _migrations (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		migration_name TEXT NOT NULL UNIQUE,
		applied_at DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
	);`
	if _, err := db.Exec(schema); err != nil {
		return apperror.WrapSimple(err, "init _migrations table")
	}
	return nil
}

func executeMigrationPlan(db *sql.DB, opts DbMigrateOptions) (DbMigrateResult, *apperror.AppError) {
	if opts.RollbackStep > 0 {
		return handleRollback(db, opts.RollbackStep, opts.IsDryRun)
	}
	if opts.SqlScript != "" {
		return handleDirectSql(db, opts.SqlScript, opts.IsDryRun)
	}
	if opts.MigrationsDir != "" {
		return handleDirectoryMigrations(db, opts.MigrationsDir, opts.IsDryRun)
	}
	return getMigrationStatus(db)
}

func handleDirectSql(db *sql.DB, sqlScript string, isDryRun bool) (DbMigrateResult, *apperror.AppError) {
	var res DbMigrateResult
	if isDryRun {
		res.AppliedMigrations = append(res.AppliedMigrations, "inline_sql:preview")
		return res, nil
	}
	tx, err := db.Begin()
	if err != nil {
		return res, apperror.WrapSimple(err, "begin sql migration tx")
	}
	if _, execErr := tx.Exec(sqlScript); execErr != nil {
		_ = tx.Rollback()
		return res, apperror.WrapSimple(execErr, "execute inline sql migration")
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return res, apperror.WrapSimple(commitErr, "commit inline sql migration")
	}
	res.AppliedMigrations = append(res.AppliedMigrations, "inline_sql:applied")
	res.TotalApplied = 1
	return res, nil
}

func handleDirectoryMigrations(db *sql.DB, dir string, isDryRun bool) (DbMigrateResult, *apperror.AppError) {
	var res DbMigrateResult
	files, errApp := discoverSqlFiles(dir)
	if errApp != nil {
		return res, errApp
	}
	appliedMap, mapErr := fetchAppliedMap(db)
	if mapErr != nil {
		return res, mapErr
	}
	for _, f := range files {
		name := filepath.Base(f)
		if appliedMap[name] {
			continue
		}
		if applyErr := applySingleMigration(db, f, name, isDryRun); applyErr != nil {
			return res, applyErr
		}
		res.AppliedMigrations = append(res.AppliedMigrations, name)
	}
	res.TotalApplied = len(res.AppliedMigrations)
	return res, nil
}

func discoverSqlFiles(dir string) ([]string, *apperror.AppError) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, apperror.WrapSimple(err, "read migrations directory")
	}
	var files []string
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(strings.ToLower(e.Name()), ".sql") {
			continue
		}
		files = append(files, filepath.Join(dir, e.Name()))
	}
	sort.Strings(files)
	return files, nil
}

func fetchAppliedMap(db *sql.DB) (map[string]bool, *apperror.AppError) {
	rows, err := db.Query(`SELECT migration_name FROM _migrations;`)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query _migrations table")
	}
	defer rows.Close()
	applied := make(map[string]bool)
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			applied[name] = true
		}
	}
	return applied, nil
}

func applySingleMigration(db *sql.DB, filePath, name string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	content, err := os.ReadFile(filePath)
	if err != nil {
		return apperror.WrapSimple(err, "read migration file: "+name)
	}
	tx, txErr := db.Begin()
	if txErr != nil {
		return apperror.WrapSimple(txErr, "begin migration tx: "+name)
	}
	if _, execErr := tx.Exec(string(content)); execErr != nil {
		_ = tx.Rollback()
		return apperror.WrapSimple(execErr, "execute migration: "+name)
	}
	if _, logErr := tx.Exec(`INSERT INTO _migrations (migration_name) VALUES (?);`, name); logErr != nil {
		_ = tx.Rollback()
		return apperror.WrapSimple(logErr, "record applied migration: "+name)
	}
	if commitErr := tx.Commit(); commitErr != nil {
		return apperror.WrapSimple(commitErr, "commit migration tx: "+name)
	}
	return nil
}

func handleRollback(db *sql.DB, steps int, isDryRun bool) (DbMigrateResult, *apperror.AppError) {
	var res DbMigrateResult
	query := `SELECT migration_name FROM _migrations ORDER BY id DESC LIMIT ?;`
	rows, err := db.Query(query, steps)
	if err != nil {
		return res, apperror.WrapSimple(err, "query migrations to rollback")
	}
	defer rows.Close()
	var targets []string
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			targets = append(targets, name)
		}
	}
	return executeRollbackSteps(db, targets, isDryRun)
}

func executeRollbackSteps(db *sql.DB, targets []string, isDryRun bool) (DbMigrateResult, *apperror.AppError) {
	var res DbMigrateResult
	for _, name := range targets {
		if err := deleteMigrationRecord(db, name, isDryRun); err != nil {
			return res, err
		}
		res.AppliedMigrations = append(res.AppliedMigrations, "rollback:"+name)
	}
	res.TotalApplied = len(res.AppliedMigrations)
	return res, nil
}

func deleteMigrationRecord(db *sql.DB, name string, isDryRun bool) *apperror.AppError {
	if isDryRun {
		return nil
	}
	_, err := db.Exec(`DELETE FROM _migrations WHERE migration_name = ?;`, name)
	if err != nil {
		return apperror.WrapSimple(err, "delete migration record: "+name)
	}
	return nil
}

func getMigrationStatus(db *sql.DB) (DbMigrateResult, *apperror.AppError) {
	var res DbMigrateResult
	rows, err := db.Query(`SELECT migration_name FROM _migrations ORDER BY id ASC;`)
	if err != nil {
		return res, apperror.WrapSimple(err, "query migration history")
	}
	defer rows.Close()
	for rows.Next() {
		var name string
		if err := rows.Scan(&name); err == nil {
			res.AppliedMigrations = append(res.AppliedMigrations, name)
			res.CurrentVersion = name
		}
	}
	res.TotalApplied = len(res.AppliedMigrations)
	return res, nil
}
