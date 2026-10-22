package cmddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// RestViewInfo represents metadata and DDL for a SQLite view.
type RestViewInfo struct {
	Name string `json:"Name"`
	Sql  string `json:"Sql"`
}

// CreateViewPayload holds fields for defining a view.
type CreateViewPayload struct {
	Name     string `json:"Name"`
	ViewName string `json:"ViewName"`
	Query    string `json:"Query"`
	Sql      string `json:"Sql"`
}

// handleCreateView creates a view on the fly.
func handleCreateView(mgr *MultiDbManager, cache *MigrationCache) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		var payload CreateViewPayload
		if decErr := json.NewDecoder(r.Body).Decode(&payload); decErr != nil {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.WrapValidation(decErr, "invalid json payload"))

			return
		}

		executeCreateView(w, mgr, cache, conn, alias, payload)
	}
}

func executeCreateView(w http.ResponseWriter, mgr *MultiDbManager, cache *MigrationCache, conn *sql.DB, alias string, payload CreateViewPayload) {
	name, query, err := extractViewParams(payload)
	if err != nil {
		writeRestEnvelope(w, http.StatusBadRequest, nil, nil, err)

		return
	}

	q := fmt.Sprintf(`CREATE VIEW IF NOT EXISTS "%s" AS %s`, name, query)
	if _, execErr := conn.Exec(q); execErr != nil {
		writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, execErr)

		return
	}

	mgr.InvalidateCache(alias)
	if cache != nil {
		_ = cache.LogViewDefinition(alias, name, query)
	}
	res := map[string]any{"ViewName": name, "IsCreated": true}
	writeRestEnvelope(w, http.StatusCreated, res, map[string]any{"IsSingle": true}, nil)
}

func extractViewParams(p CreateViewPayload) (string, string, error) {
	name := strings.TrimSpace(p.Name)
	if name == "" {
		name = strings.TrimSpace(p.ViewName)
	}
	query := strings.TrimSpace(p.Query)
	if query == "" {
		query = strings.TrimSpace(p.Sql)
	}
	if name == "" || query == "" {
		return "", "", apperror.NewValidationError("both view name and query sql are required")
	}

	return name, query, nil
}

// handleListViews lists all views in the database.
func handleListViews(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		q := `SELECT name, COALESCE(sql, '') FROM sqlite_master WHERE type = 'view' AND name NOT LIKE 'sqlite_%' ORDER BY name`
		rows, qErr := conn.Query(q)
		if qErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, qErr)

			return
		}
		defer rows.Close()

		views, scanErr := scanViewRows(rows)
		if scanErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, scanErr)

			return
		}
		writeRestEnvelope(w, http.StatusOK, views, nil, nil)
	}
}

func scanViewRows(rows *sql.Rows) ([]RestViewInfo, error) {
	var views []RestViewInfo
	for rows.Next() {
		var v RestViewInfo
		if err := rows.Scan(&v.Name, &v.Sql); err != nil {
			return nil, err
		}
		views = append(views, v)
	}

	return views, rows.Err()
}

// handleDeleteView drops an existing view.
func handleDeleteView(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		name := r.PathValue("name")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		q := fmt.Sprintf(`DROP VIEW IF EXISTS "%s"`, name)
		if _, execErr := conn.Exec(q); execErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, execErr)

			return
		}

		mgr.InvalidateCache(alias)
		attrs := map[string]any{"DroppedView": name, "TotalRecords": 0}
		writeRestEnvelope(w, http.StatusOK, []any{}, attrs, nil)
	}
}
