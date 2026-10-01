package store

import (
	"database/sql"
	"encoding/json"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

const (
	IgnoreDBFileName = "gitmap-ignore.db"

	sqlCreateIgnoreGroup = `CREATE TABLE IF NOT EXISTS IgnoreGroup (
    GroupName  TEXT PRIMARY KEY,
    IsDefault  INTEGER NOT NULL DEFAULT 0,
    Patterns   TEXT NOT NULL DEFAULT '[]'
);`
)

// IgnoreSplitDB represents an isolated SQLite database connection for ignores.
type IgnoreSplitDB struct {
	dbPath string
	db     *sql.DB
}

// IgnoreGroupRecord represents an ignore group in the database.
type IgnoreGroupRecord struct {
	Name      string
	IsDefault bool
	Patterns  []string
}

// OpenIgnoreSplitDB opens the split DB for ignore groups.
func OpenIgnoreSplitDB() (*IgnoreSplitDB, *apperror.AppError) {
	dbPath := filepath.Join(BinaryDataDir(), IgnoreDBFileName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open ignore db")
	}

	if _, err := db.Exec(sqlCreateIgnoreGroup); err != nil {
		_ = db.Close()
		return nil, apperror.WrapSimple(err, "create ignore group table")
	}

	return &IgnoreSplitDB{dbPath: dbPath, db: db}, nil
}

// Close closes the connection.
func (s *IgnoreSplitDB) Close() *apperror.AppError {
	if err := s.db.Close(); err != nil {
		return apperror.WrapSimple(err, "close ignore db")
	}
	return nil
}

// LoadGroups returns all configured ignore groups.
func (s *IgnoreSplitDB) LoadGroups() (map[string]IgnoreGroupRecord, *apperror.AppError) {
	rows, err := s.db.Query("SELECT GroupName, IsDefault, Patterns FROM IgnoreGroup")
	if err != nil {
		return nil, apperror.WrapSimple(err, "query ignore groups")
	}
	defer rows.Close()
	return parseGroupRows(rows), nil
}

func parseGroupRows(rows *sql.Rows) map[string]IgnoreGroupRecord {
	res := make(map[string]IgnoreGroupRecord)
	for rows.Next() {
		var name, patternsStr string
		var isDefaultInt int
		if err := rows.Scan(&name, &isDefaultInt, &patternsStr); err == nil {
			res[name] = buildGroupRecord(name, isDefaultInt, patternsStr)
		}
	}
	return res
}

func buildGroupRecord(name string, isDefault int, patternsStr string) IgnoreGroupRecord {
	var patterns []string
	_ = json.Unmarshal([]byte(patternsStr), &patterns)
	return IgnoreGroupRecord{
		Name:      name,
		IsDefault: isDefault == 1,
		Patterns:  patterns,
	}
}

// SaveGroup persists a single ignore group.
func (s *IgnoreSplitDB) SaveGroup(group IgnoreGroupRecord) *apperror.AppError {
	patternsData, _ := json.Marshal(group.Patterns)
	isDefaultInt := 0
	if group.IsDefault {
		isDefaultInt = 1
	}

	query := `INSERT INTO IgnoreGroup (GroupName, IsDefault, Patterns) VALUES (?, ?, ?)
              ON CONFLICT(GroupName) DO UPDATE SET IsDefault = excluded.IsDefault, Patterns = excluded.Patterns`
	_, err := s.db.Exec(query, group.Name, isDefaultInt, string(patternsData))
	if err != nil {
		return apperror.WrapSimple(err, "save ignore group")
	}
	return nil
}

// DeleteGroup removes an ignore group by name.
func (s *IgnoreSplitDB) DeleteGroup(name string) *apperror.AppError {
	_, err := s.db.Exec("DELETE FROM IgnoreGroup WHERE GroupName = ?", name)
	if err != nil {
		return apperror.WrapSimple(err, "delete ignore group")
	}
	return nil
}

// ClearGroups deletes all ignore groups.
func (s *IgnoreSplitDB) ClearGroups() *apperror.AppError {
	_, err := s.db.Exec("DELETE FROM IgnoreGroup")
	if err != nil {
		return apperror.WrapSimple(err, "clear ignore groups")
	}
	return nil
}
