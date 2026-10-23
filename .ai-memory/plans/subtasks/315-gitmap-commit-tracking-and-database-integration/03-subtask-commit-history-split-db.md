# Subtask 03: Commit History Split-DB Engine & Storage Registry Integration

> **Subtask ID:** `SUBTASK-315-03`  
> **Parent Task:** `315-gitmap-commit-tracking-and-database-integration`  
> **Status:** Pending  
> **Target Files:**  
> - `cli/store/split_db_path.go` (modify)  
> - `cli/store/commit_history_types.go` (new file)  
> - `cli/store/commit_history_schema.go` (new file)  
> - `cli/store/commit_history_split_db.go` (new file)  
> - `cli/store/split_database_registry_sync.go` (modify)  
> - `cli/store/storage_inventory.go` (modify)  
> **Related Architecture Spec:** `02-spec/21-app/315-gitmap-commit-tracking-and-database-integration/02-component-spec.md`  

---

## 1. Objective

Implement the dedicated external Commit History Split-DB engine (`repodb/commit_history.db`). The engine provides sub-millisecond (<0.5ms) lookup of commit records, file touch histories, diff line metrics, and AI agent provenance. The database file is resolved strictly outside the repository root in `<EffectiveDataDir>/repodb/commit_history.db` to prevent dirtying the active git workspace. It incorporates PascalCase relational schemas (`CommitHistoryRecord`, `CommitFileTouch`), covering indexes, positive boolean attributes (`IsAIMode`, `HasPushed`, `IsActive`), connection mutex serialization, WAL mode optimizations, and registration with `SplitDatabaseRegistry` and `storage_inventory.go`.

---

## 2. Action Checklist

- [ ] 1. Update `cli/store/split_db_path.go`:
  - [ ] Add `ResolveCommitHistoryDbPath(repoRoot string) string` helper function.
  - [ ] Target directory: `<EffectiveDataDir>/repodb/commit_history.db`.
  - [ ] Guarantee auto-creation of parent directories with `os.MkdirAll(dir, 0755)`.
  - [ ] Adhere to the 15-line function limit (preferred <= 8 lines).

- [ ] 2. Author `cli/store/commit_history_types.go`:
  - [ ] Package declaration `package store`.
  - [ ] Import `github.com/alimtvnetwork/gitmap-v28/cli/result`.
  - [ ] Define `CommitHistoryRecord` domain struct:
    - [ ] `RecordID int64`
    - [ ] `CommitSha string`
    - [ ] `ShortSha string`
    - [ ] `RepoRoot string`
    - [ ] `RepoName string`
    - [ ] `BranchName string`
    - [ ] `AuthorName string`
    - [ ] `AuthorEmail string`
    - [ ] `CommitMessage string`
    - [ ] `TaskSlug string`
    - [ ] `AgentID string`
    - [ ] `HeatmapID string`
    - [ ] `AIName string`
    - [ ] `SpecPath string`
    - [ ] `FilesChanged int`
    - [ ] `Insertions int`
    - [ ] `Deletions int`
    - [ ] `IsAIMode bool`
    - [ ] `HasPushed bool`
    - [ ] `CommittedAt int64`
    - [ ] `CreatedAt int64`
  - [ ] Define `CommitFileTouch` domain struct:
    - [ ] `TouchID int64`
    - [ ] `CommitSha string`
    - [ ] `RepoRoot string`
    - [ ] `FilePath string`
    - [ ] `ChangeType string`
    - [ ] `Insertions int`
    - [ ] `Deletions int`
    - [ ] `CreatedAt int64`
  - [ ] Define Result wrapper type aliases:
    - [ ] `CommitHistoryRecordWrap = result.Wrap[CommitHistoryRecord]`
    - [ ] `CommitHistorySliceWrap = result.WrapSlice[CommitHistoryRecord]`
    - [ ] `CommitFileTouchSliceWrap = result.WrapSlice[CommitFileTouch]`
  - [ ] Define parameter struct `CommitHistoryQueryParam` for filtering and pagination.

- [ ] 3. Author `cli/store/commit_history_schema.go`:
  - [ ] Package declaration `package store`.
  - [ ] Define constant `sqlCreateCommitHistoryTables` containing DDL for `CommitHistoryRecord` and `CommitFileTouch`.
  - [ ] Define constant `sqlCreateCommitHistoryIndexes` containing covering and lookup indexes:
    - [ ] `IdxCommitHistory_Covering` ON `CommitHistoryRecord(RepoRoot, CommittedAt DESC, CommitSha)`
    - [ ] `IdxCommitHistory_ShortSha` ON `CommitHistoryRecord(ShortSha)`
    - [ ] `IdxCommitHistory_TaskSlug` ON `CommitHistoryRecord(TaskSlug)`
    - [ ] `IdxCommitHistory_AgentId` ON `CommitHistoryRecord(AgentId)`
    - [ ] `IdxCommitFileTouch_CommitSha` ON `CommitFileTouch(CommitSha)`
    - [ ] `IdxCommitFileTouch_FilePath` ON `CommitFileTouch(FilePath)`

- [ ] 4. Author `cli/store/commit_history_split_db.go`:
  - [ ] Define `CommitHistoryDB` struct wrapping `sync.Mutex`, `*sql.DB`, and database file path.
  - [ ] Implement `OpenCommitHistoryDB(repoRoot string) (*CommitHistoryDB, error)`.
  - [ ] Implement `OpenCommitHistoryDBAt(dbPath string) (*CommitHistoryDB, error)`.
  - [ ] Configure SQLite PRAGMAs:
    - [ ] `PRAGMA journal_mode = WAL;`
    - [ ] `PRAGMA synchronous = NORMAL;`
    - [ ] `PRAGMA busy_timeout = 5000;`
    - [ ] `PRAGMA mmap_size = 268435456;` (256MB)
    - [ ] `PRAGMA cache_size = -64000;` (64MB)
    - [ ] `PRAGMA foreign_keys = ON;`
  - [ ] Implement schema migration helper `initCommitHistorySchema(conn *sql.DB) error`.
  - [ ] Implement transactional operations:
    - [ ] `RecordCommit(rec CommitHistoryRecord, touches []CommitFileTouch) error`
    - [ ] `GetCommitBySha(commitSha string) (CommitHistoryRecord, bool, error)`
    - [ ] `ListRecentCommits(limit, offset int) ([]CommitHistoryRecord, error)`
    - [ ] `ListCommitsByTask(taskSlug string) ([]CommitHistoryRecord, error)`
    - [ ] `ListTouchesByCommitSha(commitSha string) ([]CommitFileTouch, error)`
    - [ ] `Close() error`
  - [ ] Wrap errors with `apperror.WrapSimple(err, "store.<funcName>")`.
  - [ ] Ensure all functions <= 15 lines, guard clauses only, zero nested ifs.

- [ ] 5. Update `cli/store/split_database_registry_sync.go`:
  - [ ] Add `syncCommitHistoryDB()` method on `*DB`.
  - [ ] Invoke `db.syncCommitHistoryDB()` within `SyncKnownSplitDatabases()`.
  - [ ] Register with `DatabaseType: "repodb"` and `DatabaseKey: "commit_history"`.

- [ ] 6. Update `cli/store/storage_inventory.go`:
  - [ ] Verify `collectDataSubdirs` inspects `repodb/commit_history.db`.
  - [ ] Ensure `CollectAllDatabaseEntries` includes `commit_history.db` under the `"repodb"` section.

- [ ] 7. Quality & Invariant Verification:
  - [ ] Strictly relative paths (zero drive letters or URI schemes).
  - [ ] Strictly lowercase filenames.
  - [ ] Positive booleans only (`IsAIMode`, `HasPushed`, `IsActive`, `HasFound`).
  - [ ] Zero git commands, zero builds, zero test executions.

---

## 3. Concrete Code Specifications

### Specification: `cli/store/split_db_path.go`

```go
// ResolveCommitHistoryDbPath resolves the external commit history SQLite database path.
func ResolveCommitHistoryDbPath(repoRoot string) string {
	dataDir := EffectiveDataDir()
	dir := filepath.Join(dataDir, "repodb")
	_ = os.MkdirAll(dir, 0755)

	return filepath.ToSlash(filepath.Join(dir, "commit_history.db"))
}
```

---

### Specification: `cli/store/commit_history_types.go`

```go
package store

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// CommitHistoryRecord models a historical commit entry with AI provenance.
type CommitHistoryRecord struct {
	RecordID      int64  `json:"recordId"`
	CommitSha     string `json:"commitSha"`
	ShortSha      string `json:"shortSha"`
	RepoRoot      string `json:"repoRoot"`
	RepoName      string `json:"repoName"`
	BranchName    string `json:"branchName"`
	AuthorName    string `json:"authorName"`
	AuthorEmail   string `json:"authorEmail"`
	CommitMessage string `json:"commitMessage"`
	TaskSlug      string `json:"taskSlug"`
	AgentID       string `json:"agentId"`
	HeatmapID     string `json:"heatmapId"`
	AIName        string `json:"aiName"`
	SpecPath      string `json:"specPath"`
	FilesChanged  int    `json:"filesChanged"`
	Insertions    int    `json:"insertions"`
	Deletions     int    `json:"deletions"`
	IsAIMode      bool   `json:"isAiMode"`
	HasPushed     bool   `json:"hasPushed"`
	CommittedAt   int64  `json:"committedAt"`
	CreatedAt     int64  `json:"createdAt"`
}

// CommitFileTouch models an individual file modified within a commit.
type CommitFileTouch struct {
	TouchID    int64  `json:"touchId"`
	CommitSha  string `json:"commitSha"`
	RepoRoot   string `json:"repoRoot"`
	FilePath   string `json:"filePath"`
	ChangeType string `json:"changeType"`
	Insertions int    `json:"insertions"`
	Deletions  int    `json:"deletions"`
	CreatedAt  int64  `json:"createdAt"`
}

// CommitHistoryQueryParam encapsulates filtering criteria for commit queries.
type CommitHistoryQueryParam struct {
	RepoRoot string
	TaskSlug string
	AgentID  string
	Limit    int
	Offset   int
}

// Result wrapper type aliases.
type CommitHistoryRecordWrap = result.Wrap[CommitHistoryRecord]
type CommitHistorySliceWrap = result.WrapSlice[CommitHistoryRecord]
type CommitFileTouchSliceWrap = result.WrapSlice[CommitFileTouch]
```

---

### Specification: `cli/store/commit_history_schema.go`

```go
package store

const sqlCreateCommitHistoryTables = `
CREATE TABLE IF NOT EXISTS CommitHistoryRecord (
    RecordId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommitSha TEXT NOT NULL,
    ShortSha TEXT NOT NULL,
    RepoRoot TEXT NOT NULL,
    RepoName TEXT NOT NULL,
    BranchName TEXT NOT NULL,
    AuthorName TEXT NOT NULL,
    AuthorEmail TEXT NOT NULL,
    CommitMessage TEXT NOT NULL,
    TaskSlug TEXT NOT NULL DEFAULT '',
    AgentId TEXT NOT NULL DEFAULT '',
    HeatmapId TEXT NOT NULL DEFAULT '',
    AiName TEXT NOT NULL DEFAULT '',
    SpecPath TEXT NOT NULL DEFAULT '',
    FilesChanged INTEGER NOT NULL DEFAULT 0,
    Insertions INTEGER NOT NULL DEFAULT 0,
    Deletions INTEGER NOT NULL DEFAULT 0,
    IsAiMode INTEGER NOT NULL DEFAULT 0,
    HasPushed INTEGER NOT NULL DEFAULT 0,
    CommittedAt INTEGER NOT NULL,
    CreatedAt INTEGER NOT NULL,
    CONSTRAINT Uq_CommitHistory_Sha UNIQUE (CommitSha)
);

CREATE TABLE IF NOT EXISTS CommitFileTouch (
    TouchId INTEGER PRIMARY KEY AUTOINCREMENT,
    CommitSha TEXT NOT NULL,
    RepoRoot TEXT NOT NULL,
    FilePath TEXT NOT NULL,
    ChangeType TEXT NOT NULL,
    Insertions INTEGER NOT NULL DEFAULT 0,
    Deletions INTEGER NOT NULL DEFAULT 0,
    CreatedAt INTEGER NOT NULL,
    FOREIGN KEY (CommitSha) REFERENCES CommitHistoryRecord(CommitSha) ON DELETE CASCADE
);
`

const sqlCreateCommitHistoryIndexes = `
CREATE INDEX IF NOT EXISTS IdxCommitHistory_Covering 
    ON CommitHistoryRecord(RepoRoot, CommittedAt DESC, CommitSha);

CREATE INDEX IF NOT EXISTS IdxCommitHistory_ShortSha 
    ON CommitHistoryRecord(ShortSha);

CREATE INDEX IF NOT EXISTS IdxCommitHistory_TaskSlug 
    ON CommitHistoryRecord(TaskSlug);

CREATE INDEX IF NOT EXISTS IdxCommitHistory_AgentId 
    ON CommitHistoryRecord(AgentId);

CREATE INDEX IF NOT EXISTS IdxCommitFileTouch_CommitSha 
    ON CommitFileTouch(CommitSha);

CREATE INDEX IF NOT EXISTS IdxCommitFileTouch_FilePath 
    ON CommitFileTouch(FilePath);
`
```

---

### Specification: `cli/store/commit_history_split_db.go`

```go
package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	_ "modernc.org/sqlite"
)

// CommitHistoryDB manages high-performance storage for commit records and file touches.
type CommitHistoryDB struct {
	mu   sync.Mutex
	conn *sql.DB
	path string
}

// OpenCommitHistoryDB opens or initializes the commit history split SQLite database.
func OpenCommitHistoryDB(repoRoot string) (*CommitHistoryDB, error) {
	dbPath := ResolveCommitHistoryDbPath(repoRoot)

	return OpenCommitHistoryDBAt(dbPath)
}

// OpenCommitHistoryDBAt opens the database at a specific filesystem path.
func OpenCommitHistoryDBAt(dbPath string) (*CommitHistoryDB, error) {
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenCommitHistoryDBAt.mkdir")
	}

	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.OpenCommitHistoryDBAt.open")
	}

	if errPragma := applyCommitHistoryPragmas(conn); errPragma != nil {
		_ = conn.Close()

		return nil, errPragma
	}

	if errSchema := initCommitHistorySchema(conn); errSchema != nil {
		_ = conn.Close()

		return nil, errSchema
	}

	return &CommitHistoryDB{conn: conn, path: dbPath}, nil
}

func applyCommitHistoryPragmas(conn *sql.DB) error {
	pragmas := []string{
		"PRAGMA journal_mode = WAL;",
		"PRAGMA synchronous = NORMAL;",
		"PRAGMA busy_timeout = 5000;",
		"PRAGMA mmap_size = 268435456;",
		"PRAGMA cache_size = -64000;",
		"PRAGMA foreign_keys = ON;",
	}
	for _, p := range pragmas {
		if _, err := conn.Exec(p); err != nil {
			return apperror.WrapSimple(err, "store.applyCommitHistoryPragmas")
		}
	}

	return nil
}

func initCommitHistorySchema(conn *sql.DB) error {
	if _, err := conn.Exec(sqlCreateCommitHistoryTables); err != nil {
		return apperror.WrapSimple(err, "store.initCommitHistorySchema.tables")
	}
	if _, err := conn.Exec(sqlCreateCommitHistoryIndexes); err != nil {
		return apperror.WrapSimple(err, "store.initCommitHistorySchema.indexes")
	}

	return nil
}

// RecordCommit inserts a commit record and all associated file touches atomically.
func (db *CommitHistoryDB) RecordCommit(rec CommitHistoryRecord, touches []CommitFileTouch) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "store.RecordCommit.begin")
	}
	defer func() { _ = tx.Rollback() }()

	if errInsert := insertCommitRecordTx(tx, rec); errInsert != nil {
		return errInsert
	}

	if errTouches := insertTouchesTx(tx, rec.CommitSha, rec.RepoRoot, touches); errTouches != nil {
		return errTouches
	}

	if errCommit := tx.Commit(); errCommit != nil {
		return apperror.WrapSimple(errCommit, "store.RecordCommit.commit")
	}

	return nil
}

func insertCommitRecordTx(tx *sql.Tx, rec CommitHistoryRecord) error {
	query := `INSERT OR REPLACE INTO CommitHistoryRecord (
		CommitSha, ShortSha, RepoRoot, RepoName, BranchName,
		AuthorName, AuthorEmail, CommitMessage, TaskSlug, AgentId,
		HeatmapId, AiName, SpecPath, FilesChanged, Insertions,
		Deletions, IsAiMode, HasPushed, CommittedAt, CreatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	nowUnix := time.Now().Unix()
	isAiInt := 0
	if rec.IsAIMode {
		isAiInt = 1
	}
	hasPushInt := 0
	if rec.HasPushed {
		hasPushInt = 1
	}

	_, err := tx.Exec(query,
		rec.CommitSha, rec.ShortSha, rec.RepoRoot, rec.RepoName, rec.BranchName,
		rec.AuthorName, rec.AuthorEmail, rec.CommitMessage, rec.TaskSlug, rec.AgentID,
		rec.HeatmapID, rec.AIName, rec.SpecPath, rec.FilesChanged, rec.Insertions,
		rec.Deletions, isAiInt, hasPushInt, rec.CommittedAt, nowUnix,
	)
	if err != nil {
		return apperror.WrapSimple(err, "store.insertCommitRecordTx")
	}

	return nil
}

func insertTouchesTx(tx *sql.Tx, sha, repoRoot string, touches []CommitFileTouch) error {
	query := `INSERT INTO CommitFileTouch (
		CommitSha, RepoRoot, FilePath, ChangeType, Insertions, Deletions, CreatedAt
	) VALUES (?, ?, ?, ?, ?, ?, ?);`

	nowUnix := time.Now().Unix()
	for _, t := range touches {
		_, err := tx.Exec(query, sha, repoRoot, t.FilePath, t.ChangeType, t.Insertions, t.Deletions, nowUnix)
		if err != nil {
			return apperror.WrapSimple(err, "store.insertTouchesTx")
		}
	}

	return nil
}

// GetCommitBySha fetches a commit record by its commit hash.
func (db *CommitHistoryDB) GetCommitBySha(sha string) (CommitHistoryRecord, bool, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `SELECT RecordId, CommitSha, ShortSha, RepoRoot, RepoName, BranchName,
		AuthorName, AuthorEmail, CommitMessage, TaskSlug, AgentId, HeatmapId,
		AiName, SpecPath, FilesChanged, Insertions, Deletions, IsAiMode, HasPushed,
		CommittedAt, CreatedAt
	FROM CommitHistoryRecord WHERE CommitSha = ? OR ShortSha = ? LIMIT 1;`

	row := db.conn.QueryRow(query, sha, sha)
	rec, hasFound, err := scanCommitRecord(row)
	if err != nil {
		return rec, false, apperror.WrapSimple(err, "store.GetCommitBySha")
	}

	return rec, hasFound, nil
}

func scanCommitRecord(scanner interface{ Scan(dest ...any) error }) (CommitHistoryRecord, bool, error) {
	var rec CommitHistoryRecord
	var isAiInt, hasPushInt int

	err := scanner.Scan(
		&rec.RecordID, &rec.CommitSha, &rec.ShortSha, &rec.RepoRoot, &rec.RepoName,
		&rec.BranchName, &rec.AuthorName, &rec.AuthorEmail, &rec.CommitMessage,
		&rec.TaskSlug, &rec.AgentID, &rec.HeatmapID, &rec.AIName, &rec.SpecPath,
		&rec.FilesChanged, &rec.Insertions, &rec.Deletions, &isAiInt, &hasPushInt,
		&rec.CommittedAt, &rec.CreatedAt,
	)
	if err == sql.ErrNoRows {
		return rec, false, nil
	}
	if err != nil {
		return rec, false, err
	}

	rec.IsAIMode = isAiInt == 1
	rec.HasPushed = hasPushInt == 1

	return rec, true, nil
}

// ListRecentCommits retrieves the latest commit records ordered by commit time.
func (db *CommitHistoryDB) ListRecentCommits(limit, offset int) ([]CommitHistoryRecord, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `SELECT RecordId, CommitSha, ShortSha, RepoRoot, RepoName, BranchName,
		AuthorName, AuthorEmail, CommitMessage, TaskSlug, AgentId, HeatmapId,
		AiName, SpecPath, FilesChanged, Insertions, Deletions, IsAiMode, HasPushed,
		CommittedAt, CreatedAt
	FROM CommitHistoryRecord ORDER BY CommittedAt DESC LIMIT ? OFFSET ?;`

	rows, err := db.conn.Query(query, limit, offset)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListRecentCommits.query")
	}
	defer rows.Close()

	return scanCommitRecordRows(rows)
}

func scanCommitRecordRows(rows *sql.Rows) ([]CommitHistoryRecord, error) {
	var list []CommitHistoryRecord
	for rows.Next() {
		rec, hasFound, err := scanCommitRecord(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "store.scanCommitRecordRows")
		}
		if hasFound {
			list = append(list, rec)
		}
	}

	return list, nil
}

// ListCommitsByTask returns all commits linked to a specific task slug.
func (db *CommitHistoryDB) ListCommitsByTask(taskSlug string) ([]CommitHistoryRecord, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `SELECT RecordId, CommitSha, ShortSha, RepoRoot, RepoName, BranchName,
		AuthorName, AuthorEmail, CommitMessage, TaskSlug, AgentId, HeatmapId,
		AiName, SpecPath, FilesChanged, Insertions, Deletions, IsAiMode, HasPushed,
		CommittedAt, CreatedAt
	FROM CommitHistoryRecord WHERE TaskSlug = ? ORDER BY CommittedAt DESC;`

	rows, err := db.conn.Query(query, taskSlug)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListCommitsByTask.query")
	}
	defer rows.Close()

	return scanCommitRecordRows(rows)
}

// ListTouchesByCommitSha retrieves all touched files for a given commit SHA.
func (db *CommitHistoryDB) ListTouchesByCommitSha(sha string) ([]CommitFileTouch, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	query := `SELECT TouchId, CommitSha, RepoRoot, FilePath, ChangeType, Insertions, Deletions, CreatedAt
	FROM CommitFileTouch WHERE CommitSha = ? ORDER BY FilePath ASC;`

	rows, err := db.conn.Query(query, sha)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListTouchesByCommitSha.query")
	}
	defer rows.Close()

	var touches []CommitFileTouch
	for rows.Next() {
		var t CommitFileTouch
		if err := rows.Scan(&t.TouchID, &t.CommitSha, &t.RepoRoot, &t.FilePath, &t.ChangeType, &t.Insertions, &t.Deletions, &t.CreatedAt); err != nil {
			return nil, apperror.WrapSimple(err, "store.ListTouchesByCommitSha.scan")
		}
		touches = append(touches, t)
	}

	return touches, nil
}

// Close gracefully releases the database connection.
func (db *CommitHistoryDB) Close() error {
	if db == nil || db.conn == nil {
		return nil
	}

	return db.conn.Close()
}
```

---

### Specification: `cli/store/split_database_registry_sync.go`

```go
func (db *DB) syncCommitHistoryDB() {
	path := ResolveCommitHistoryDbPath("")
	desc := "Repository commit history, file touches, and AI agent provenance split database"
	entry := inspectSplitDBFile("repodb", "commit_history", path, desc)
	_ = db.RegisterSplitDB(entry)
}
```

Wiring into `SyncKnownSplitDatabases()`:

```go
// SyncKnownSplitDatabases discovers and updates metadata for all known split databases.
func (db *DB) SyncKnownSplitDatabases() error {
	if err := db.syncInstallationDB(); err != nil {
		return apperror.WrapSimple(err, "store.SyncKnownSplitDatabases.installation")
	}

	db.syncTasksDB()
	db.syncScheduleDBs()
	db.syncPipelineDBs()
	db.syncAiInstructionDB()
	db.syncSearchDB()
	db.syncSectionTasksDBs()
	db.syncCommitHistoryDB()

	return nil
}
```
