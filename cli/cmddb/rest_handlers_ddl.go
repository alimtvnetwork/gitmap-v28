package cmddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// CreateTablePayload holds table creation DDL parameters.
type CreateTablePayload struct {
	TableName string `json:"TableName"`
	Ddl       string `json:"Ddl"`
	Sql       string `json:"Sql"`
}

// AlterTablePayload holds table alteration parameters.
type AlterTablePayload struct {
	Table  string `json:"Table"`
	Action string `json:"Action"`
	Column string `json:"Column"`
	Type   string `json:"Type"`
	Sql    string `json:"Sql"`
}

// MigratePayload holds schema migration script details.
type MigratePayload struct {
	Version string `json:"Version"`
	Name    string `json:"Name"`
	Sql     string `json:"Sql"`
}

// ClonePayload holds parameters for cloning an active database.
type ClonePayload struct {
	TargetAlias string `json:"TargetAlias"`
	TargetPath  string `json:"TargetPath"`
}

// handleCreateTable creates a new table via DDL.
func handleCreateTable(mgr *MultiDbManager, cache *MigrationCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		var payload CreateTablePayload
		if decErr := json.NewDecoder(r.Body).Decode(&payload); decErr != nil {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, decErr)

			return
		}

		executeCreateTable(w, mgr, cache, conn, alias, payload)
	}
}

func executeCreateTable(w http.ResponseWriter, mgr *MultiDbManager, cache *MigrationCache, conn *sql.DB, alias string, p CreateTablePayload) {
	ddl := strings.TrimSpace(p.Ddl)
	if ddl == "" {
		ddl = strings.TrimSpace(p.Sql)
	}
	if ddl == "" {
		writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.NewValidationError("ddl sql is required"))

		return
	}

	if _, err := conn.Exec(ddl); err != nil {
		writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, err)

		return
	}

	mgr.InvalidateCache(alias)
	if cache != nil {
		_ = cache.LogAlteration(alias, p.TableName, "CREATE_TABLE", ddl, true)
	}
	writeRestEnvelope(w, http.StatusCreated, map[string]any{"TableName": p.TableName, "IsCreated": true}, map[string]any{"IsSingle": true}, nil)
}

// handleAlter applies schema alterations such as adding columns.
func handleAlter(mgr *MultiDbManager, cache *MigrationCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		var payload AlterTablePayload
		if decErr := json.NewDecoder(r.Body).Decode(&payload); decErr != nil {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, decErr)

			return
		}

		executeAlter(w, mgr, cache, conn, alias, payload)
	}
}

func executeAlter(w http.ResponseWriter, mgr *MultiDbManager, cache *MigrationCache, conn *sql.DB, alias string, p AlterTablePayload) {
	sqlStr := resolveAlterSql(p)
	if sqlStr == "" {
		writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.NewValidationError("incomplete alter payload"))

		return
	}

	if _, err := conn.Exec(sqlStr); err != nil {
		writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, err)

		return
	}

	mgr.InvalidateCache(alias)
	if cache != nil {
		_ = cache.LogAlteration(alias, p.Table, "ALTER_TABLE", sqlStr, true)
	}
	writeRestEnvelope(w, http.StatusOK, map[string]any{"Table": p.Table, "IsAltered": true}, map[string]any{"IsSingle": true}, nil)
}

func resolveAlterSql(p AlterTablePayload) string {
	if strings.TrimSpace(p.Sql) != "" {
		return strings.TrimSpace(p.Sql)
	}
	if p.Table == "" || p.Column == "" {
		return ""
	}
	action := strings.TrimSpace(p.Action)
	if action == "" {
		action = "ADD COLUMN"
	}

	return fmt.Sprintf(`ALTER TABLE "%s" %s "%s" %s`, p.Table, action, p.Column, p.Type)
}

// handleMigrate runs a migration SQL script and records it.
func handleMigrate(mgr *MultiDbManager, cache *MigrationCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		var payload MigratePayload
		if decErr := json.NewDecoder(r.Body).Decode(&payload); decErr != nil {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, decErr)

			return
		}

		executeMigrate(w, mgr, cache, conn, alias, payload)
	}
}

func logMigrationResult(cache *MigrationCache, alias string, p MigratePayload, applied bool, errMsg string) {
	if cache == nil {
		return
	}
	_ = cache.LogMigration(alias, p.Version, p.Name, p.Sql, applied, errMsg)
}

func executeMigrate(w http.ResponseWriter, mgr *MultiDbManager, cache *MigrationCache, conn *sql.DB, alias string, p MigratePayload) {
	if strings.TrimSpace(p.Sql) == "" {
		writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.NewValidationError("migration sql is empty"))

		return
	}

	if _, err := conn.Exec(p.Sql); err != nil {
		logMigrationResult(cache, alias, p, false, err.Error())
		writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, err)

		return
	}

	mgr.InvalidateCache(alias)
	logMigrationResult(cache, alias, p, true, "")
	res := map[string]any{"Version": p.Version, "Name": p.Name, "HasApplied": true}
	writeRestEnvelope(w, http.StatusOK, res, map[string]any{"IsSingle": true}, nil)
}

// handleClone clones a database to a target path and dynamically registers it.
func handleClone(mgr *MultiDbManager, cache *MigrationCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		srcPath, err := mgr.GetPath(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		var payload ClonePayload
		if decErr := json.NewDecoder(r.Body).Decode(&payload); decErr != nil {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, decErr)

			return
		}

		executeClone(w, mgr, cache, alias, srcPath, payload)
	}
}

func executeClone(w http.ResponseWriter, mgr *MultiDbManager, cache *MigrationCache, alias, srcPath string, p ClonePayload) {
	targetAlias, targetPath := strings.TrimSpace(p.TargetAlias), strings.TrimSpace(p.TargetPath)
	if targetAlias == "" || targetPath == "" {
		writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.NewValidationError("targetAlias and targetPath are required"))

		return
	}

	size, copyErr := copyFileSafe(srcPath, targetPath)
	if copyErr != nil {
		writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, copyErr)

		return
	}

	if verifyErr := verifyClonedSqlite(targetPath); verifyErr != nil {
		writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, verifyErr)

		return
	}

	_ = mgr.Mount(targetAlias, targetPath)
	if cache != nil {
		_ = cache.LogClonedDatabase(alias, targetAlias, srcPath, targetPath, size)
	}
	res := map[string]any{
		"SourceAlias":   alias,
		"TargetAlias":   targetAlias,
		"TargetPath":    targetPath,
		"FileSizeBytes": size,
		"IsMounted":     true,
	}
	writeRestEnvelope(w, http.StatusOK, res, map[string]any{"IsSingle": true}, nil)
}

func copyFileSafe(src, dst string) (int64, error) {
	if err := os.MkdirAll(filepath.Dir(dst), 0755); err != nil {
		return 0, err
	}
	in, err := os.Open(src)
	if err != nil {
		return 0, err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return 0, err
	}
	defer out.Close()

	return io.Copy(out, in)
}

func verifyClonedSqlite(path string) error {
	info, err := os.Stat(path)
	if err != nil || info.Size() == 0 {
		return apperror.NewExecutionError("cloned file verification failed: empty or missing")
	}

	openRes := store.OpenSQLiteDB(path)
	if openRes.IsFailure() {
		return apperror.WrapSimple(openRes.AppError(), "cloned database verification failed")
	}
	_ = openRes.Value.Close()

	return nil
}
