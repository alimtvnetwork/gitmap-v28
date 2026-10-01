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

	sqlCreateIgnoreBinding = `CREATE TABLE IF NOT EXISTS IgnoreRepoBinding (
    RepoPath   TEXT NOT NULL,
    GroupName  TEXT NOT NULL,
    PRIMARY KEY (RepoPath, GroupName)
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

// IgnoreRepoBindingRecord represents a repo-to-group mapping in the database.
type IgnoreRepoBindingRecord struct {
	RepoPath  string `json:"repoPath"`
	GroupName string `json:"groupName"`
}

// OpenIgnoreSplitDB opens the split DB for ignore groups.
func OpenIgnoreSplitDB() (*IgnoreSplitDB, *apperror.AppError) {
	dbPath := filepath.Join(BinaryDataDir(), IgnoreDBFileName)
	db, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open ignore db")
	}
	initErr := initIgnoreTables(db)
	if initErr != nil {
		_ = db.Close()
		return nil, initErr
	}
	return &IgnoreSplitDB{dbPath: dbPath, db: db}, nil
}

func initIgnoreTables(db *sql.DB) *apperror.AppError {
	if _, err := db.Exec(sqlCreateIgnoreGroup); err != nil {
		return apperror.WrapSimple(err, "create ignore group table")
	}
	if _, err := db.Exec(sqlCreateIgnoreBinding); err != nil {
		return apperror.WrapSimple(err, "create ignore binding table")
	}
	return nil
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
	isDef := (isDefault == 1)
	return IgnoreGroupRecord{
		Name:      name,
		IsDefault: isDef,
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
	if _, err := s.db.Exec(query, group.Name, isDefaultInt, string(patternsData)); err != nil {
		return apperror.WrapSimple(err, "save ignore group")
	}
	return nil
}

// DeleteGroup removes an ignore group by name.
func (s *IgnoreSplitDB) DeleteGroup(name string) *apperror.AppError {
	if _, err := s.db.Exec("DELETE FROM IgnoreGroup WHERE GroupName = ?", name); err != nil {
		return apperror.WrapSimple(err, "delete ignore group")
	}
	return nil
}

// ClearGroups deletes all ignore groups.
func (s *IgnoreSplitDB) ClearGroups() *apperror.AppError {
	if _, err := s.db.Exec("DELETE FROM IgnoreGroup"); err != nil {
		return apperror.WrapSimple(err, "clear ignore groups")
	}
	return nil
}

// SaveRepoBinding persists a repo to group association.
func (s *IgnoreSplitDB) SaveRepoBinding(repoPath, groupName string) *apperror.AppError {
	query := `INSERT INTO IgnoreRepoBinding (RepoPath, GroupName) VALUES (?, ?)
              ON CONFLICT(RepoPath, GroupName) DO NOTHING`
	if _, err := s.db.Exec(query, repoPath, groupName); err != nil {
		return apperror.WrapSimple(err, "save repo binding")
	}
	return nil
}

// DeleteRepoBinding removes a repo to group association.
func (s *IgnoreSplitDB) DeleteRepoBinding(repoPath, groupName string) *apperror.AppError {
	query := `DELETE FROM IgnoreRepoBinding WHERE RepoPath = ? AND GroupName = ?`
	if _, err := s.db.Exec(query, repoPath, groupName); err != nil {
		return apperror.WrapSimple(err, "delete repo binding")
	}
	return nil
}

// DeleteBindingsForGroup removes all bindings associated with a group.
func (s *IgnoreSplitDB) DeleteBindingsForGroup(groupName string) *apperror.AppError {
	query := `DELETE FROM IgnoreRepoBinding WHERE GroupName = ?`
	if _, err := s.db.Exec(query, groupName); err != nil {
		return apperror.WrapSimple(err, "delete bindings for group")
	}
	return nil
}

// LoadRepoBindings returns all configured repo to group bindings.
func (s *IgnoreSplitDB) LoadRepoBindings() ([]IgnoreRepoBindingRecord, *apperror.AppError) {
	rows, err := s.db.Query("SELECT RepoPath, GroupName FROM IgnoreRepoBinding ORDER BY GroupName, RepoPath")
	if err != nil {
		return nil, apperror.WrapSimple(err, "query ignore repo bindings")
	}
	defer rows.Close()
	return parseBindingRows(rows), nil
}

func parseBindingRows(rows *sql.Rows) []IgnoreRepoBindingRecord {
	var bindings []IgnoreRepoBindingRecord
	for rows.Next() {
		var repoPath, groupName string
		if err := rows.Scan(&repoPath, &groupName); err == nil {
			bindings = append(bindings, IgnoreRepoBindingRecord{
				RepoPath:  repoPath,
				GroupName: groupName,
			})
		}
	}
	return bindings
}

// LoadBindingsForGroup returns repo paths bound to a specific group.
func (s *IgnoreSplitDB) LoadBindingsForGroup(groupName string) ([]string, *apperror.AppError) {
	query := `SELECT RepoPath FROM IgnoreRepoBinding WHERE GroupName = ? ORDER BY RepoPath`
	rows, err := s.db.Query(query, groupName)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query bindings for group")
	}
	defer rows.Close()
	return scanSingleColumnStrings(rows), nil
}

// LoadBindingsForRepo returns group names bound to a specific repo.
func (s *IgnoreSplitDB) LoadBindingsForRepo(repoPath string) ([]string, *apperror.AppError) {
	query := `SELECT GroupName FROM IgnoreRepoBinding WHERE RepoPath = ? ORDER BY GroupName`
	rows, err := s.db.Query(query, repoPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query bindings for repo")
	}
	defer rows.Close()
	return scanSingleColumnStrings(rows), nil
}

func scanSingleColumnStrings(rows *sql.Rows) []string {
	var list []string
	for rows.Next() {
		var val string
		if err := rows.Scan(&val); err == nil {
			list = append(list, val)
		}
	}
	return list
}

// ClearRepoBindings removes all repo bindings.
func (s *IgnoreSplitDB) ClearRepoBindings() *apperror.AppError {
	if _, err := s.db.Exec("DELETE FROM IgnoreRepoBinding"); err != nil {
		return apperror.WrapSimple(err, "clear repo bindings")
	}
	return nil
}
