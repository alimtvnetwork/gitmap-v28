package cmdsupabase

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const createTableQuery = `
CREATE TABLE IF NOT EXISTS supabase_databases (
    alias TEXT PRIMARY KEY,
    project_ref TEXT,
    api_url TEXT NOT NULL,
    anon_key_enc TEXT NOT NULL,
    service_key_enc TEXT NOT NULL,
    db_url_enc TEXT,
    status TEXT NOT NULL DEFAULT 'active',
    description TEXT,
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL
);
CREATE INDEX IF NOT EXISTS idx_supabase_status ON supabase_databases(status);
`

// OpenSupabaseDB opens or initializes the dedicated SQLite split database for Supabase.
func OpenSupabaseDB() (*sql.DB, error) {
	dbPath := store.ResolveSplitDbPath(store.SectionInstallation, "supabase", "")
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "OpenSupabaseDB_Open")
	}

	if err := initSupabaseSchema(db); err != nil {
		_ = db.Close()
		return nil, err
	}

	return db, nil
}

func initSupabaseSchema(db *sql.DB) error {
	_, err := db.Exec(createTableQuery)
	if err != nil {
		return apperror.WrapSimple(err, "initSupabaseSchema_Exec")
	}

	return nil
}

// InsertProject persists a new or updated encrypted Supabase project record.
func InsertProject(record SupabaseProjectRecord) error {
	db, err := OpenSupabaseDB()
	if err != nil {
		return err
	}
	defer db.Close()

	return executeInsertProject(db, record)
}

func executeInsertProject(db *sql.DB, record SupabaseProjectRecord) error {
	query := `INSERT OR REPLACE INTO supabase_databases 
		(alias, project_ref, api_url, anon_key_enc, service_key_enc, db_url_enc, status, description, created_at, updated_at) 
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	now := time.Now().UTC()
	createdAt := resolveCreatedAt(record.CreatedAt, now)

	_, err := db.Exec(query, record.Alias, record.ProjectRef, record.ApiUrl,
		record.AnonKeyEnc, record.ServiceKeyEnc, record.DbUrlEnc,
		record.Status, record.Description, createdAt, now)
	if err != nil {
		return apperror.WrapSimple(err, "executeInsertProject")
	}

	return nil
}

func resolveCreatedAt(t, fallback time.Time) time.Time {
	if t.IsZero() {
		return fallback
	}

	return t
}

// GetProjectByAlias queries an encrypted project record by its alias.
func GetProjectByAlias(alias string) (*SupabaseProjectRecord, error) {
	db, err := OpenSupabaseDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return queryProjectByAlias(db, alias)
}

func queryProjectByAlias(db *sql.DB, alias string) (*SupabaseProjectRecord, error) {
	query := `SELECT alias, project_ref, api_url, anon_key_enc, service_key_enc, db_url_enc, status, description, created_at, updated_at 
		FROM supabase_databases WHERE alias = ?`

	var rec SupabaseProjectRecord
	row := db.QueryRow(query, alias)
	err := row.Scan(&rec.Alias, &rec.ProjectRef, &rec.ApiUrl,
		&rec.AnonKeyEnc, &rec.ServiceKeyEnc, &rec.DbUrlEnc,
		&rec.Status, &rec.Description, &rec.CreatedAt, &rec.UpdatedAt)
	if err == sql.ErrNoRows {
		return nil, apperror.NewNotFoundError("supabase project alias not found: " + alias)
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "queryProjectByAlias_Scan")
	}

	return &rec, nil
}

// ListAllProjects returns all stored Supabase project records ordered by alias.
func ListAllProjects() ([]SupabaseProjectRecord, error) {
	db, err := OpenSupabaseDB()
	if err != nil {
		return nil, err
	}
	defer db.Close()

	return queryAllProjects(db)
}

func queryAllProjects(db *sql.DB) ([]SupabaseProjectRecord, error) {
	query := `SELECT alias, project_ref, api_url, anon_key_enc, service_key_enc, db_url_enc, status, description, created_at, updated_at 
		FROM supabase_databases ORDER BY alias ASC`

	rows, err := db.Query(query)
	if err != nil {
		return nil, apperror.WrapSimple(err, "queryAllProjects_Query")
	}
	defer rows.Close()

	return scanProjectRows(rows)
}

func scanProjectRows(rows *sql.Rows) ([]SupabaseProjectRecord, error) {
	var projects []SupabaseProjectRecord
	for rows.Next() {
		var rec SupabaseProjectRecord
		if err := rows.Scan(&rec.Alias, &rec.ProjectRef, &rec.ApiUrl,
			&rec.AnonKeyEnc, &rec.ServiceKeyEnc, &rec.DbUrlEnc,
			&rec.Status, &rec.Description, &rec.CreatedAt, &rec.UpdatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "scanProjectRows_Scan")
		}
		projects = append(projects, rec)
	}

	return projects, nil
}

// UpdateProjectStatus updates the status and timestamp of a Supabase project.
func UpdateProjectStatus(alias string, status string) error {
	db, err := OpenSupabaseDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `UPDATE supabase_databases SET status = ?, updated_at = ? WHERE alias = ?`
	res, err := db.Exec(query, status, time.Now().UTC(), alias)
	if err != nil {
		return apperror.WrapSimple(err, "UpdateProjectStatus_Exec")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return apperror.NewNotFoundError("supabase project alias not found: " + alias)
	}

	return nil
}

// DeleteProjectByAlias removes a project entry from the vault.
func DeleteProjectByAlias(alias string) error {
	db, err := OpenSupabaseDB()
	if err != nil {
		return err
	}
	defer db.Close()

	query := `DELETE FROM supabase_databases WHERE alias = ?`
	res, err := db.Exec(query, alias)
	if err != nil {
		return apperror.WrapSimple(err, "DeleteProjectByAlias_Exec")
	}

	rows, _ := res.RowsAffected()
	if rows == 0 {
		return apperror.NewNotFoundError("supabase project alias not found: " + alias)
	}

	return nil
}
