package cmddb

import (
	"database/sql"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type QueryPayload struct {
	Sql   string `json:"Sql"`
	Query string `json:"Query"`
}

// handleQuery executes arbitrary SQL queries and mutations.
func handleQuery(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias := r.PathValue("alias")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)

			return
		}

		sqlStr, parseErr := extractSqlQuery(r)
		if parseErr != nil || strings.TrimSpace(sqlStr) == "" {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.NewValidationError("missing or empty sql query"))

			return
		}

		executeAndRespondQuery(w, mgr, conn, alias, strings.TrimSpace(sqlStr))
	}
}

func extractSqlQuery(r *http.Request) (string, error) {
	if r.Method == http.MethodGet {
		q := r.URL.Query().Get("sql")
		if q == "" {
			q = r.URL.Query().Get("q")
		}

		return q, nil
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		return "", err
	}
	var payload QueryPayload
	if jsonErr := json.Unmarshal(body, &payload); jsonErr == nil {
		if payload.Sql != "" {
			return payload.Sql, nil
		}
		if payload.Query != "" {
			return payload.Query, nil
		}
	}

	return string(body), nil
}

func executeAndRespondQuery(w http.ResponseWriter, mgr *MultiDbManager, conn *sql.DB, alias, sqlStr string) {
	if isReadOnlyQuery(sqlStr) {
		rows, err := conn.Query(sqlStr)
		if err != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, err)

			return
		}
		defer rows.Close()

		results, scanErr := scanAllRows(rows)
		if scanErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, scanErr)

			return
		}
		writeRestEnvelope(w, http.StatusOK, results, nil, nil)

		return
	}

	executeMutationQuery(w, mgr, conn, alias, sqlStr)
}

func executeMutationQuery(w http.ResponseWriter, mgr *MultiDbManager, conn *sql.DB, alias, sqlStr string) {
	res, err := conn.Exec(sqlStr)
	if err != nil {
		writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, err)

		return
	}
	rowsAff, _ := res.RowsAffected()
	lastId, _ := res.LastInsertId()
	if isSchemaMutation(sqlStr) {
		mgr.InvalidateCache(alias)
	}

	resPayload := map[string]any{"RowsAffected": rowsAff, "LastInsertId": lastId}
	writeRestEnvelope(w, http.StatusOK, resPayload, map[string]any{"IsSingle": true}, nil)
}

func isReadOnlyQuery(sqlStr string) bool {
	upper := strings.ToUpper(strings.TrimSpace(sqlStr))

	return strings.HasPrefix(upper, "SELECT") ||
		strings.HasPrefix(upper, "PRAGMA") ||
		strings.HasPrefix(upper, "EXPLAIN") ||
		strings.HasPrefix(upper, "WITH")
}

func isSchemaMutation(sqlStr string) bool {
	upper := strings.ToUpper(strings.TrimSpace(sqlStr))

	return strings.HasPrefix(upper, "CREATE") ||
		strings.HasPrefix(upper, "DROP") ||
		strings.HasPrefix(upper, "ALTER")
}
