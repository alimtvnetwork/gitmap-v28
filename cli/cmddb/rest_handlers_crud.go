package cmddb

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

// handleGetRecords queries paginated, filtered table records.
func handleGetRecords(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, table := r.PathValue("alias"), r.PathValue("table")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)
			return
		}
		records, attrs, qErr := queryTableRecords(mgr, conn, alias, table, r)
		if qErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, qErr)
			return
		}
		writeRestEnvelope(w, http.StatusOK, records, attrs, nil)
	}
}

func queryTableRecords(mgr *MultiDbManager, conn *sql.DB, alias, table string, r *http.Request) ([]map[string]any, map[string]any, error) {
	page, limit := parsePagination(r)
	cols, _ := mgr.GetColumns(alias, table)
	where, args := buildWhereFilters(r, cols)
	total := countMatchingRecords(conn, table, where, args)

	order := buildOrderClause(r, cols)
	offset := (page - 1) * limit
	query := fmt.Sprintf(`SELECT * FROM "%s"%s%s LIMIT ? OFFSET ?`, table, where, order)
	queryArgs := append(args, limit, offset)

	rows, err := conn.Query(query, queryArgs...)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	records, err := scanAllRows(rows)
	attrs := map[string]any{"TotalRecords": total, "PerPage": limit, "CurrentPage": page, "IsSingle": false}
	return records, attrs, err
}

func parsePagination(r *http.Request) (int, int) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	if page < 1 {
		page = 1
	}
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	if limit < 1 || limit > 1000 {
		limit = 20
	}
	return page, limit
}

func buildWhereFilters(r *http.Request, cols []RestColumnInfo) (string, []any) {
	colMap := make(map[string]bool, len(cols))
	for _, c := range cols {
		colMap[strings.ToLower(c.Name)] = true
	}
	var clauses []string
	var args []any
	for k, v := range r.URL.Query() {
		lk := strings.ToLower(k)
		if lk == "page" || lk == "limit" || lk == "sort" || lk == "order" {
			continue
		}
		if colMap[lk] && len(v) > 0 {
			clauses = append(clauses, fmt.Sprintf(`"%s" = ?`, k))
			args = append(args, v[0])
		}
	}
	if len(clauses) == 0 {
		return "", nil
	}
	return " WHERE " + strings.Join(clauses, " AND "), args
}

func buildOrderClause(r *http.Request, cols []RestColumnInfo) string {
	sortCol := r.URL.Query().Get("sort")
	if sortCol == "" {
		return ""
	}
	for _, c := range cols {
		if strings.EqualFold(c.Name, sortCol) {
			order := "ASC"
			if strings.EqualFold(r.URL.Query().Get("order"), "desc") {
				order = "DESC"
			}
			return fmt.Sprintf(` ORDER BY "%s" %s`, c.Name, order)
		}
	}
	return ""
}

func countMatchingRecords(conn *sql.DB, table, where string, args []any) int {
	var count int
	q := fmt.Sprintf(`SELECT COUNT(*) FROM "%s"%s`, table, where)
	_ = conn.QueryRow(q, args...).Scan(&count)
	return count
}

func scanAllRows(rows *sql.Rows) ([]map[string]any, error) {
	cols, err := rows.Columns()
	if err != nil {
		return nil, err
	}
	results := make([]map[string]any, 0)
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if err := rows.Scan(ptrs...); err != nil {
			return nil, err
		}
		row := make(map[string]any, len(cols))
		for i, c := range cols {
			if b, ok := vals[i].([]byte); ok {
				row[c] = string(b)
			} else {
				row[c] = vals[i]
			}
		}
		results = append(results, row)
	}
	return results, rows.Err()
}

// handleGetRecordById queries a single record by primary key.
func handleGetRecordById(mgr *MultiDbManager) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		alias, table, id := r.PathValue("alias"), r.PathValue("table"), r.PathValue("id")
		conn, _, err := mgr.GetDB(alias)
		if err != nil {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, err)
			return
		}
		pk := findRestPrimaryKeyCol(mgr, alias, table)
		q := fmt.Sprintf(`SELECT * FROM "%s" WHERE "%s" = ? LIMIT 1`, table, pk)
		rows, qErr := conn.Query(q, id)
		if qErr != nil {
			writeRestEnvelope(w, http.StatusInternalServerError, nil, nil, qErr)
			return
		}
		defer rows.Close()

		records, _ := scanAllRows(rows)
		if len(records) == 0 {
			writeRestEnvelope(w, http.StatusNotFound, nil, nil, fmt.Errorf("record not found: %s", id))
			return
		}
		writeRestEnvelope(w, http.StatusOK, records, map[string]any{"IsSingle": true, "TotalRecords": 1}, nil)
	}
}

func findRestPrimaryKeyCol(mgr *MultiDbManager, alias, table string) string {
	cols, err := mgr.GetColumns(alias, table)
	if err == nil {
		for _, c := range cols {
			if c.IsPrimaryKey {
				return c.Name
			}
		}
	}
	return "rowid"
}

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
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, fmt.Errorf("invalid json payload: %w", err))
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
		return nil, fmt.Errorf("empty insert payload")
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
			writeRestEnvelope(w, http.StatusBadRequest, nil, nil, fmt.Errorf("invalid json payload: %w", err))
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
		return nil, fmt.Errorf("empty update payload")
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
