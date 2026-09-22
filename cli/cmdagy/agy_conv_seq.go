package cmdagy

import (
	"database/sql"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// AgyProjectSequenceRecord models a persistent sequence mapping for AGY projects.
type AgyProjectSequenceRecord struct {
	ProjectID     string `json:"projectId"`
	ProjectName   string `json:"projectName"`
	WorkspacePath string `json:"workspacePath"`
	SequenceNum   int    `json:"sequenceNum"`
}

func getAgySequenceDBPath() string {
	return store.DefaultDBPath()
}

func openAgySequenceDB() (*sql.DB, error) {
	dbPath := getAgySequenceDBPath()
	_ = os.MkdirAll(filepath.Dir(dbPath), 0755)
	conn, err := store.OpenSQLiteDB(dbPath)
	if err != nil {
		return nil, err
	}
	initErr := initAgySequenceTable(conn)
	if initErr != nil {
		_ = conn.Close()
		return nil, initErr
	}

	return conn, nil
}

func initAgySequenceTable(conn *sql.DB) error {
	schema := `CREATE TABLE IF NOT EXISTS AgyProjectSequence (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		project_id TEXT UNIQUE NOT NULL,
		project_name TEXT NOT NULL,
		workspace_path TEXT NOT NULL,
		sequence_num INTEGER UNIQUE NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);`
	_, err := conn.Exec(schema)
	if err != nil {
		return apperror.WrapSimple(err, "init AgyProjectSequence table")
	}

	return nil
}

// GetOrAssignProjectSequence retrieves or generates a stable global sequence number.
func GetOrAssignProjectSequence(projectID, projectName, workspacePath string) (int, error) {
	conn, err := openAgySequenceDB()
	if err != nil {
		return 0, err
	}
	defer conn.Close()

	seq, isFound, findErr := queryExistingProjectSequence(conn, projectID, workspacePath)
	if findErr != nil {
		return 0, findErr
	}
	if isFound {
		return seq, nil
	}

	return insertNextProjectSequence(conn, projectID, projectName, workspacePath)
}

func queryExistingProjectSequence(conn *sql.DB, projectID, workspacePath string) (int, bool, error) {
	query := "SELECT sequence_num FROM AgyProjectSequence WHERE project_id = ? OR (workspace_path != '' AND workspace_path = ?) LIMIT 1"
	row := conn.QueryRow(query, projectID, workspacePath)
	var seq int
	scanErr := row.Scan(&seq)
	if scanErr == sql.ErrNoRows {
		return 0, false, nil
	}
	if scanErr != nil {
		return 0, false, apperror.WrapSimple(scanErr, "scan project sequence")
	}

	return seq, true, nil
}

func insertNextProjectSequence(conn *sql.DB, projectID, projectName, workspacePath string) (int, error) {
	nextSeq, seqErr := queryNextMaxSequence(conn)
	if seqErr != nil {
		return 0, seqErr
	}
	insertSQL := "INSERT INTO AgyProjectSequence (project_id, project_name, workspace_path, sequence_num) VALUES (?, ?, ?, ?)"
	_, execErr := conn.Exec(insertSQL, projectID, projectName, workspacePath, nextSeq)
	if execErr != nil {
		return 0, apperror.WrapSimple(execErr, "insert project sequence")
	}

	return nextSeq, nil
}

func queryNextMaxSequence(conn *sql.DB) (int, error) {
	query := "SELECT COALESCE(MAX(sequence_num), 0) + 1 FROM AgyProjectSequence"
	var nextSeq int
	err := conn.QueryRow(query).Scan(&nextSeq)
	if err != nil {
		return 0, apperror.WrapSimple(err, "query max sequence")
	}

	return nextSeq, nil
}

// GetProjectBySequence resolves a project record using its unique sequence number.
func GetProjectBySequence(seq int) (*AgyProjectSequenceRecord, error) {
	conn, err := openAgySequenceDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	query := "SELECT project_id, project_name, workspace_path, sequence_num FROM AgyProjectSequence WHERE sequence_num = ? LIMIT 1"
	var rec AgyProjectSequenceRecord
	scanErr := conn.QueryRow(query, seq).Scan(&rec.ProjectID, &rec.ProjectName, &rec.WorkspacePath, &rec.SequenceNum)
	if scanErr == sql.ErrNoRows {
		return nil, nil
	}
	if scanErr != nil {
		return nil, apperror.WrapSimple(scanErr, "scan project by sequence")
	}

	return &rec, nil
}

// GetAllProjectSequences loads all persisted sequences keyed by project ID.
func GetAllProjectSequences() (map[string]int, error) {
	conn, err := openAgySequenceDB()
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	rows, qErr := conn.Query("SELECT project_id, sequence_num FROM AgyProjectSequence")
	if qErr != nil {
		return nil, apperror.WrapSimple(qErr, "query all project sequences")
	}
	defer rows.Close()

	return scanSequenceRows(rows)
}

func scanSequenceRows(rows *sql.Rows) (map[string]int, error) {
	m := make(map[string]int)
	for rows.Next() {
		var id string
		var seq int
		if err := rows.Scan(&id, &seq); err == nil {
			m[id] = seq
		}
	}

	return m, nil
}
