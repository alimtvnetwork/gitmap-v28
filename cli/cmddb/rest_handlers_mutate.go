package cmddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// handleInsertRecord inserts a record from JSON body.
func handleInsertRecord(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, table := r.PathValue("alias"), r.PathValue("table")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)
			return
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.WrapValidation(err, "invalid json payload"))
			return
		}
		res, insErr := executeRestInsert(conn, table, payload)
		if insErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, insErr)
			return
		}
		writeRestEnvelope(w, http.StatusCreated, res, map[string]any{"IsSingle": true}, nil)
	}
}

func executeRestInsert(conn *sql.DB, table string, payload map[string]any) (map[string]any, error) {
	if len(payload) == 0 {
		return nil, apperror.NewValidationError("empty insert payload")
	}
	var cols, placeholders []string
	var vals []any
	for k, v := range payload {
		cols = append(cols, fmt.Sprintf(`"%s"`, k))
		placeholders = append(placeholders, "?")
		vals = append(vals, v)
	}
	q := fmt.Sprintf(`INSERT INTO "%s" (%s) VALUES (%s)`, table, strings.Join(cols, ", "), strings.Join(placeholders, ", "))
	res, err := conn.Exec(q, vals...)
	if err != nil {
		return nil, err
	}
	lastId, _ := res.LastInsertId()
	return map[string]any{"Id": lastId, "RowsAffected": int64(1)}, nil
}

// handleUpdateRecord updates a record by primary key.
func handleUpdateRecord(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, table, id := r.PathValue("alias"), r.PathValue("table"), r.PathValue("id")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)
			return
		}
		var payload map[string]any
		if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, apperror.WrapValidation(err, "invalid json payload"))
			return
		}
		pk := findRestPrimaryKeyCol(mgr, alias, table)
		res, updErr := executeRestUpdate(conn, table, pk, id, payload)
		if updErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, updErr)
			return
		}
		writeRestEnvelope(w, http.StatusOK, res, map[string]any{"IsSingle": true}, nil)
	}
}

func executeRestUpdate(conn *sql.DB, table, pk, id string, payload map[string]any) (map[string]any, error) {
	if len(payload) == 0 {
		return nil, apperror.NewValidationError("empty update payload")
	}
	var sets []string
	var vals []any
	for k, v := range payload {
		sets = append(sets, fmt.Sprintf(`"%s" = ?`, k))
		vals = append(vals, v)
	}
	vals = append(vals, id)
	q := fmt.Sprintf(`UPDATE "%s" SET %s WHERE "%s" = ?`, table, strings.Join(sets, ", "), pk)
	res, err := conn.Exec(q, vals...)
	if err != nil {
		return nil, err
	}
	rowsAff, _ := res.RowsAffected()
	return map[string]any{"Id": id, "RowsAffected": rowsAff}, nil
}

// handleDeleteRecord deletes a record by primary key.
func handleDeleteRecord(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, table, id := r.PathValue("alias"), r.PathValue("table"), r.PathValue("id")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)
			return
		}
		pk := findRestPrimaryKeyCol(mgr, alias, table)
		q := fmt.Sprintf(`DELETE FROM "%s" WHERE "%s" = ?`, table, pk)
		res, err := conn.Exec(q, id)
		if err != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, err)
			return
		}
		rowsAff, _ := res.RowsAffected()
		attrs := map[string]any{"RowsAffected": rowsAff, "TotalRecords": 0}
		writeRestEnvelope(w, http.StatusOK, []any{}, attrs, nil)
	}
}
