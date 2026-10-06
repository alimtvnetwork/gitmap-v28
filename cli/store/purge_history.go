package store

import (
	"database/sql"
	"encoding/json"
	"os"
	"path/filepath"
	"sync"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// PurgeHistoryLog represents a legacy record of a git purge operation.
type PurgeHistoryLog struct {
	PurgeHistoryLogId int64  `json:"purgeHistoryLogId"`
	RepoPath          string `json:"repoPath"`
	Pattern           string `json:"pattern"`
	BackupBranch      string `json:"backupBranch"`
	TempDir           string `json:"tempDir"`
	Files             string `json:"files"`
	Timestamp         int64  `json:"timestamp"`
	IsRestored        bool   `json:"isRestored"`
	Notes             string `json:"notes,omitempty"`
	Comments          string `json:"comments,omitempty"`
}

const (
	sqlCreatePurgeHistory = `CREATE TABLE IF NOT EXISTS PurgeHistoryLog (
	PurgeHistoryLogId INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoPath TEXT NOT NULL,
	Pattern TEXT NOT NULL,
	BackupBranch TEXT NOT NULL,
	TempDir TEXT NOT NULL,
	Files TEXT NOT NULL,
	Timestamp INTEGER NOT NULL,
	IsRestored INTEGER NOT NULL DEFAULT 0,
	Notes TEXT NULL,
	Comments TEXT NULL
)`

	sqlInsertPurgeHistory = `INSERT INTO PurgeHistoryLog (RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, IsRestored, Notes, Comments) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlCreateHistoryPurgeOp = `CREATE TABLE IF NOT EXISTS HistoryPurgeOperation (
	OperationId INTEGER PRIMARY KEY AUTOINCREMENT,
	RepoSlug TEXT NOT NULL DEFAULT '',
	TargetType TEXT NOT NULL DEFAULT '',
	TargetPath TEXT NOT NULL DEFAULT '',
	CommitListJson TEXT NOT NULL DEFAULT '[]',
	TotalCommitsScanned INTEGER NOT NULL DEFAULT 0,
	AffectedCommitsCount INTEGER NOT NULL DEFAULT 0,
	BackedUpFilesCount INTEGER NOT NULL DEFAULT 0,
	BackupVaultPath TEXT NOT NULL DEFAULT '',
	IsDryRun INTEGER NOT NULL DEFAULT 0,
	IsVerified INTEGER NOT NULL DEFAULT 0,
	HasPushed INTEGER NOT NULL DEFAULT 0,
	IsUndone INTEGER NOT NULL DEFAULT 0,
	IsSuccess INTEGER NOT NULL DEFAULT 0,
	CreatedAt INTEGER NOT NULL DEFAULT 0,
	CompletedAt INTEGER NOT NULL DEFAULT 0,
	Notes TEXT NOT NULL DEFAULT ''
);
CREATE INDEX IF NOT EXISTS idx_history_purge_op_repo ON HistoryPurgeOperation(RepoSlug);
CREATE INDEX IF NOT EXISTS idx_history_purge_op_created ON HistoryPurgeOperation(CreatedAt);`

	sqlCreateHistoryPurgeCommitMap = `CREATE TABLE IF NOT EXISTS HistoryPurgeCommitMap (
	MapId INTEGER PRIMARY KEY AUTOINCREMENT,
	OperationId INTEGER NOT NULL,
	OriginalCommitSha TEXT NOT NULL,
	RewrittenCommitSha TEXT NOT NULL,
	ParentOriginalSha TEXT NOT NULL DEFAULT '',
	ParentRewrittenSha TEXT NOT NULL DEFAULT '',
	AuthorName TEXT NOT NULL DEFAULT '',
	AuthorEmail TEXT NOT NULL DEFAULT '',
	CommitTimestamp INTEGER NOT NULL DEFAULT 0,
	CommitMessage TEXT NOT NULL DEFAULT '',
	IsPurged INTEGER NOT NULL DEFAULT 0,
	CreatedAt INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY(OperationId) REFERENCES HistoryPurgeOperation(OperationId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_history_purge_cmap_op ON HistoryPurgeCommitMap(OperationId);
CREATE INDEX IF NOT EXISTS idx_history_purge_cmap_orig ON HistoryPurgeCommitMap(OriginalCommitSha);
CREATE INDEX IF NOT EXISTS idx_history_purge_cmap_rewr ON HistoryPurgeCommitMap(RewrittenCommitSha);`

	sqlCreateHistoryPurgeFile = `CREATE TABLE IF NOT EXISTS HistoryPurgeFile (
	FileId INTEGER PRIMARY KEY AUTOINCREMENT,
	OperationId INTEGER NOT NULL,
	OriginalCommitSha TEXT NOT NULL,
	RelativePath TEXT NOT NULL,
	BlobSha TEXT NOT NULL DEFAULT '',
	FileSizeBytes INTEGER NOT NULL DEFAULT 0,
	FileMode INTEGER NOT NULL DEFAULT 420,
	BackupRelPath TEXT NOT NULL,
	IsPurged INTEGER NOT NULL DEFAULT 1,
	CreatedAt INTEGER NOT NULL DEFAULT 0,
	FOREIGN KEY(OperationId) REFERENCES HistoryPurgeOperation(OperationId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_history_purge_file_op ON HistoryPurgeFile(OperationId);
CREATE INDEX IF NOT EXISTS idx_history_purge_file_path ON HistoryPurgeFile(RelativePath);`

	sqlCreateHistoryUndoOp = `CREATE TABLE IF NOT EXISTS HistoryUndoOperation (
	UndoId INTEGER PRIMARY KEY AUTOINCREMENT,
	OperationId INTEGER NOT NULL,
	RepoSlug TEXT NOT NULL DEFAULT '',
	RestoredCommitCount INTEGER NOT NULL DEFAULT 0,
	RestoredFileCount INTEGER NOT NULL DEFAULT 0,
	BackupVaultPath TEXT NOT NULL DEFAULT '',
	IsVerified INTEGER NOT NULL DEFAULT 0,
	IsSuccess INTEGER NOT NULL DEFAULT 0,
	CreatedAt INTEGER NOT NULL DEFAULT 0,
	CompletedAt INTEGER NOT NULL DEFAULT 0,
	ErrorMessage TEXT NOT NULL DEFAULT '',
	FOREIGN KEY(OperationId) REFERENCES HistoryPurgeOperation(OperationId) ON DELETE CASCADE
);
CREATE INDEX IF NOT EXISTS idx_history_undo_op ON HistoryUndoOperation(OperationId);`

	sqlInsertHistoryPurgeOp = `INSERT INTO HistoryPurgeOperation (
	RepoSlug, TargetType, TargetPath, CommitListJson, TotalCommitsScanned,
	AffectedCommitsCount, BackedUpFilesCount, BackupVaultPath, IsDryRun,
	IsVerified, HasPushed, IsUndone, IsSuccess, CreatedAt, CompletedAt, Notes
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlInsertCommitMap = `INSERT INTO HistoryPurgeCommitMap (
	OperationId, OriginalCommitSha, RewrittenCommitSha, ParentOriginalSha,
	ParentRewrittenSha, AuthorName, AuthorEmail, CommitTimestamp,
	CommitMessage, IsPurged, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlInsertPurgeFile = `INSERT INTO HistoryPurgeFile (
	OperationId, OriginalCommitSha, RelativePath, BlobSha,
	FileSizeBytes, FileMode, BackupRelPath, IsPurged, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`

	sqlInsertUndoOp = `INSERT INTO HistoryUndoOperation (
	OperationId, RepoSlug, RestoredCommitCount, RestoredFileCount,
	BackupVaultPath, IsVerified, IsSuccess, CreatedAt, CompletedAt, ErrorMessage
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
)

// PurgeHistoryDB wraps dedicated SQLite database operations for history purging.
type PurgeHistoryDB struct {
	*DB
	mu   sync.Mutex
	path string
}

// OpenPurgeHistoryDB opens the dedicated SQLite database with WAL mode.
func OpenPurgeHistoryDB(customPath string) (*PurgeHistoryDB, error) {
	dbPath := resolvePurgeDbPath(customPath)
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		return nil, apperror.WrapSimple(err, "open purge history db: mkdir")
	}

	innerDB, err := OpenAt(dbPath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "open purge history db: open")
	}

	pdb := &PurgeHistoryDB{DB: innerDB, path: dbPath}
	if err := pdb.InitSchema(); err != nil {
		_ = pdb.Close()
		return nil, err
	}

	return pdb, nil
}

func resolvePurgeDbPath(customPath string) string {
	if customPath != "" {
		return customPath
	}

	return filepath.Join(BinaryDataDir(), "repodb", "history_purge.db")
}

// InitSchema creates the tables and indexes if they do not exist.
func (pdb *PurgeHistoryDB) InitSchema() error {
	pdb.mu.Lock()
	defer pdb.mu.Unlock()

	return pdb.DB.EnsureHistoryPurgeTables()
}

// EnsureHistoryPurgeTables creates tables and configures WAL pragmas.
func (db *DB) EnsureHistoryPurgeTables() error {
	pragmas := []string{
		"PRAGMA journal_mode=WAL;",
		"PRAGMA busy_timeout=5000;",
		"PRAGMA synchronous=NORMAL;",
		"PRAGMA foreign_keys=ON;",
	}
	for _, p := range pragmas {
		_, _ = ExecWrapper(db.conn, p).Destruct()
	}

	return db.execPurgeDDLStatements()
}

func (db *DB) execPurgeDDLStatements() error {
	statements := []string{
		sqlCreateHistoryPurgeOp,
		sqlCreateHistoryPurgeCommitMap,
		sqlCreateHistoryPurgeFile,
		sqlCreateHistoryUndoOp,
	}
	for _, stmt := range statements {
		if _, err := ExecWrapper(db.conn, stmt).Destruct(); err != nil {
			return apperror.WrapSimple(err, "exec purge ddl statement")
		}
	}

	return nil
}

// InsertPurgeHistoryLog inserts a legacy purge history entry.
func (db *DB) InsertPurgeHistoryLog(log *PurgeHistoryLog) error {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return err
	}

	res, err := ExecWrapper(db.conn, sqlInsertPurgeHistory,
		log.RepoPath, log.Pattern, log.BackupBranch, log.TempDir, log.Files, log.Timestamp, boolToInt(log.IsRestored), log.Notes, log.Comments,
	).Destruct()
	if err != nil {
		return err
	}

	id, err := res.LastInsertId()
	if err == nil {
		log.PurgeHistoryLogId = id
	}

	return err
}

// GetLastPurgeHistoryLog retrieves the most recent unrestored purge history log for a repo.
func (db *DB) GetLastPurgeHistoryLog(repoPath string) (*PurgeHistoryLog, error) {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return nil, err
	}

	row := QueryRowWrapper(
		db.conn,
		`SELECT PurgeHistoryLogId, RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, IsRestored, COALESCE(Notes, ''), COALESCE(Comments, '') FROM PurgeHistoryLog WHERE RepoPath = ? AND IsRestored = 0 ORDER BY PurgeHistoryLogId DESC LIMIT 1`,
		repoPath,
	)

	return scanPurgeHistoryRow(row)
}

// GetPurgeHistoryLogById retrieves a purge history log by its ID.
func (db *DB) GetPurgeHistoryLogById(id int64) (*PurgeHistoryLog, error) {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return nil, err
	}

	row := QueryRowWrapper(
		db.conn,
		`SELECT PurgeHistoryLogId, RepoPath, Pattern, BackupBranch, TempDir, Files, Timestamp, IsRestored, COALESCE(Notes, ''), COALESCE(Comments, '') FROM PurgeHistoryLog WHERE PurgeHistoryLogId = ?`,
		id,
	)

	return scanPurgeHistoryRow(row)
}

func scanPurgeHistoryRow(row *sql.Row) (*PurgeHistoryLog, error) {
	var log PurgeHistoryLog
	var isRestoredInt int
	err := row.Scan(
		&log.PurgeHistoryLogId, &log.RepoPath, &log.Pattern,
		&log.BackupBranch, &log.TempDir, &log.Files,
		&log.Timestamp, &isRestoredInt, &log.Notes, &log.Comments,
	)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	log.IsRestored = isRestoredInt == 1

	return &log, nil
}

// MarkPurgeHistoryRestored marks a purge history log as restored.
func (db *DB) MarkPurgeHistoryRestored(id int64) error {
	if err := db.EnsurePurgeHistoryTable(); err != nil {
		return err
	}

	_, err := ExecWrapper(db.conn, `UPDATE PurgeHistoryLog SET IsRestored = 1 WHERE PurgeHistoryLogId = ?`, id).Destruct()

	return err
}

// InsertHistoryPurgeOperation inserts a new HistoryPurgeOperation record.
func (db *DB) InsertHistoryPurgeOperation(op *HistoryPurgeOperation) (int64, error) {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return 0, err
	}
	if op.CreatedAt == 0 {
		op.CreatedAt = time.Now().Unix()
	}

	res, err := ExecWrapper(db.conn, sqlInsertHistoryPurgeOp,
		op.RepoSlug, op.TargetType, op.TargetPath, op.CommitListJson,
		op.TotalCommitsScanned, op.AffectedCommitsCount, op.BackedUpFilesCount,
		op.BackupVaultPath, boolToInt(op.IsDryRun), boolToInt(op.IsVerified),
		boolToInt(op.HasPushed), boolToInt(op.IsUndone), boolToInt(op.IsSuccess),
		op.CreatedAt, op.CompletedAt, op.Notes,
	).Destruct()
	if err != nil {
		return 0, apperror.WrapSimple(err, "insert history purge op")
	}

	id, err := res.LastInsertId()
	if err == nil {
		op.OperationId = id
	}

	return id, err
}

// CreatePurgeOperation wraps parameters and inserts a new HistoryPurgeOperation record.
func (db *DB) CreatePurgeOperation(opts CreatePurgeOptions) (int64, error) {
	rawJson, _ := json.Marshal(opts.CommitList)
	op := &HistoryPurgeOperation{
		RepoSlug:       opts.RepoSlug,
		TargetType:     opts.TargetType,
		TargetPath:     opts.TargetPath,
		CommitListJson: string(rawJson),
		IsDryRun:       opts.IsDryRun,
		HasPushed:      opts.HasPushed,
	}

	return db.InsertHistoryPurgeOperation(op)
}

// UpdateHistoryPurgeOperationStatus updates completion status of an operation.
func (db *DB) UpdateHistoryPurgeOperationStatus(opts UpdatePurgeStatusOptions) error {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return err
	}

	now := time.Now().Unix()
	query := `UPDATE HistoryPurgeOperation SET IsSuccess = ?, IsVerified = ?, Notes = ?, CompletedAt = ? WHERE OperationId = ?`
	_, err := ExecWrapper(db.conn, query, boolToInt(opts.IsSuccess), boolToInt(opts.IsVerified), opts.Notes, now, opts.OperationId).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "update purge op status")
	}

	return nil
}

// UpdatePurgeStatus is an alias for UpdateHistoryPurgeOperationStatus.
func (db *DB) UpdatePurgeStatus(opts UpdatePurgeStatusOptions) error {
	return db.UpdateHistoryPurgeOperationStatus(opts)
}

// InsertHistoryPurgeCommitMapBatch batch inserts rewritten commit mappings in a transaction.
func (db *DB) InsertHistoryPurgeCommitMapBatch(opId int64, mappings []HistoryPurgeCommitMap) error {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return err
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "begin commit map batch tx")
	}
	defer func() { _ = tx.Rollback() }()

	if err := execBatchCommitMap(tx, opId, mappings); err != nil {
		return err
	}

	return tx.Commit()
}

func execBatchCommitMap(tx *sql.Tx, opId int64, mappings []HistoryPurgeCommitMap) error {
	stmt, err := tx.Prepare(sqlInsertCommitMap)
	if err != nil {
		return apperror.WrapSimple(err, "prepare insert commit map")
	}
	defer stmt.Close()

	now := time.Now().Unix()
	for _, m := range mappings {
		if m.CreatedAt == 0 {
			m.CreatedAt = now
		}
		_, err := stmt.Exec(opId, m.OriginalCommitSha, m.RewrittenCommitSha, m.ParentOriginalSha,
			m.ParentRewrittenSha, m.AuthorName, m.AuthorEmail, m.CommitTimestamp,
			m.CommitMessage, boolToInt(m.IsPurged), m.CreatedAt)
		if err != nil {
			return apperror.WrapSimple(err, "exec insert commit map")
		}
	}

	return nil
}

// InsertCommitMappings is an alias for InsertHistoryPurgeCommitMapBatch.
func (db *DB) InsertCommitMappings(opId int64, mappings []HistoryPurgeCommitMap) error {
	return db.InsertHistoryPurgeCommitMapBatch(opId, mappings)
}

// InsertHistoryPurgeFileBatch batch inserts purged file records in a transaction.
func (db *DB) InsertHistoryPurgeFileBatch(opId int64, files []HistoryPurgeFile) error {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return err
	}

	tx, err := db.conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "begin purge file batch tx")
	}
	defer func() { _ = tx.Rollback() }()

	if err := execBatchPurgedFiles(tx, opId, files); err != nil {
		return err
	}

	return tx.Commit()
}

func execBatchPurgedFiles(tx *sql.Tx, opId int64, files []HistoryPurgeFile) error {
	stmt, err := tx.Prepare(sqlInsertPurgeFile)
	if err != nil {
		return apperror.WrapSimple(err, "prepare insert purge file")
	}
	defer stmt.Close()

	now := time.Now().Unix()
	for _, f := range files {
		if f.CreatedAt == 0 {
			f.CreatedAt = now
		}
		_, err := stmt.Exec(opId, f.OriginalCommitSha, f.RelativePath, f.BlobSha,
			f.FileSizeBytes, f.FileMode, f.BackupRelPath, boolToInt(f.IsPurged), f.CreatedAt)
		if err != nil {
			return apperror.WrapSimple(err, "exec insert purge file")
		}
	}

	return nil
}

// InsertPurgedFiles is an alias for InsertHistoryPurgeFileBatch.
func (db *DB) InsertPurgedFiles(opId int64, files []HistoryPurgeFile) error {
	return db.InsertHistoryPurgeFileBatch(opId, files)
}

// GetHistoryPurgeOperationById retrieves an operation record by ID.
func (db *DB) GetHistoryPurgeOperationById(opId int64) (*HistoryPurgeOperation, error) {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return nil, err
	}

	row := QueryRowWrapper(db.conn, `SELECT OperationId, RepoSlug, TargetType, TargetPath,
		CommitListJson, TotalCommitsScanned, AffectedCommitsCount, BackedUpFilesCount,
		BackupVaultPath, IsDryRun, IsVerified, HasPushed, IsUndone, IsSuccess,
		CreatedAt, CompletedAt, Notes FROM HistoryPurgeOperation WHERE OperationId = ?`, opId)

	return scanHistoryPurgeOpRow(row)
}

// GetPurgeOperation is an alias for GetHistoryPurgeOperationById.
func (db *DB) GetPurgeOperation(opId int64) (*HistoryPurgeOperation, error) {
	return db.GetHistoryPurgeOperationById(opId)
}

// GetLastHistoryPurgeOperation retrieves the most recent purge operation for a repo slug.
func (db *DB) GetLastHistoryPurgeOperation(repoSlug string) (*HistoryPurgeOperation, error) {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return nil, err
	}

	row := QueryRowWrapper(db.conn, `SELECT OperationId, RepoSlug, TargetType, TargetPath,
		CommitListJson, TotalCommitsScanned, AffectedCommitsCount, BackedUpFilesCount,
		BackupVaultPath, IsDryRun, IsVerified, HasPushed, IsUndone, IsSuccess,
		CreatedAt, CompletedAt, Notes FROM HistoryPurgeOperation WHERE RepoSlug = ? ORDER BY OperationId DESC LIMIT 1`, repoSlug)

	return scanHistoryPurgeOpRow(row)
}

func scanHistoryPurgeOpRow(row *sql.Row) (*HistoryPurgeOperation, error) {
	var op HistoryPurgeOperation
	var isDryRun, isVerified, hasPushed, isUndone, isSuccess int
	err := row.Scan(&op.OperationId, &op.RepoSlug, &op.TargetType, &op.TargetPath,
		&op.CommitListJson, &op.TotalCommitsScanned, &op.AffectedCommitsCount, &op.BackedUpFilesCount,
		&op.BackupVaultPath, &isDryRun, &isVerified, &hasPushed, &isUndone, &isSuccess,
		&op.CreatedAt, &op.CompletedAt, &op.Notes)
	if err == sql.ErrNoRows {
		return nil, nil
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "scan history purge op row")
	}
	op.IsDryRun = isDryRun == 1
	op.IsVerified = isVerified == 1
	op.HasPushed = hasPushed == 1
	op.IsUndone = isUndone == 1
	op.IsSuccess = isSuccess == 1

	return &op, nil
}

// GetCommitMappings retrieves all commit mappings for an operation.
func (db *DB) GetCommitMappings(opId int64) ([]HistoryPurgeCommitMap, error) {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return nil, err
	}

	query := `SELECT MapId, OperationId, OriginalCommitSha, RewrittenCommitSha, ParentOriginalSha,
		ParentRewrittenSha, AuthorName, AuthorEmail, CommitTimestamp, CommitMessage, IsPurged, CreatedAt
		FROM HistoryPurgeCommitMap WHERE OperationId = ? ORDER BY MapId ASC`
	rows, err := db.conn.Query(query, opId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query commit mappings")
	}
	defer rows.Close()

	return scanCommitMappingsList(rows)
}

func scanCommitMappingsList(rows *sql.Rows) ([]HistoryPurgeCommitMap, error) {
	var list []HistoryPurgeCommitMap
	for rows.Next() {
		var m HistoryPurgeCommitMap
		var isPurged int
		err := rows.Scan(&m.MapId, &m.OperationId, &m.OriginalCommitSha, &m.RewrittenCommitSha,
			&m.ParentOriginalSha, &m.ParentRewrittenSha, &m.AuthorName, &m.AuthorEmail,
			&m.CommitTimestamp, &m.CommitMessage, &isPurged, &m.CreatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scan commit map item")
		}
		m.IsPurged = isPurged == 1
		list = append(list, m)
	}

	return list, rows.Err()
}

// GetPurgedFiles retrieves all purged files for an operation.
func (db *DB) GetPurgedFiles(opId int64) ([]HistoryPurgeFile, error) {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return nil, err
	}

	query := `SELECT FileId, OperationId, OriginalCommitSha, RelativePath, BlobSha,
		FileSizeBytes, FileMode, BackupRelPath, IsPurged, CreatedAt
		FROM HistoryPurgeFile WHERE OperationId = ? ORDER BY FileId ASC`
	rows, err := db.conn.Query(query, opId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "query purged files")
	}
	defer rows.Close()

	return scanPurgedFilesList(rows)
}

func scanPurgedFilesList(rows *sql.Rows) ([]HistoryPurgeFile, error) {
	var list []HistoryPurgeFile
	for rows.Next() {
		var f HistoryPurgeFile
		var isPurged int
		err := rows.Scan(&f.FileId, &f.OperationId, &f.OriginalCommitSha, &f.RelativePath,
			&f.BlobSha, &f.FileSizeBytes, &f.FileMode, &f.BackupRelPath, &isPurged, &f.CreatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "scan purged file item")
		}
		f.IsPurged = isPurged == 1
		list = append(list, f)
	}

	return list, rows.Err()
}

// InsertHistoryUndoOperation records an execution of the undo engine.
func (db *DB) InsertHistoryUndoOperation(undo *HistoryUndoOperation) (int64, error) {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return 0, err
	}
	if undo.CreatedAt == 0 {
		undo.CreatedAt = time.Now().Unix()
	}

	res, err := ExecWrapper(db.conn, sqlInsertUndoOp,
		undo.OperationId, undo.RepoSlug, undo.RestoredCommitCount,
		undo.RestoredFileCount, undo.BackupVaultPath, boolToInt(undo.IsVerified),
		boolToInt(undo.IsSuccess), undo.CreatedAt, undo.CompletedAt, undo.ErrorMessage,
	).Destruct()
	if err != nil {
		return 0, apperror.WrapSimple(err, "insert undo op")
	}

	id, err := res.LastInsertId()
	if err == nil {
		undo.UndoId = id
	}

	return id, err
}

// CreateUndoOperation records an undo attempt for an operation.
func (db *DB) CreateUndoOperation(opId int64, repoSlug string) (int64, error) {
	undo := &HistoryUndoOperation{
		OperationId: opId,
		RepoSlug:    repoSlug,
	}

	return db.InsertHistoryUndoOperation(undo)
}

// UpdateUndoStatus updates the completion status of an undo attempt.
func (db *DB) UpdateUndoStatus(undoId int64, isSuccess bool, errMsg string) error {
	if err := db.EnsureHistoryPurgeTables(); err != nil {
		return err
	}

	now := time.Now().Unix()
	query := `UPDATE HistoryUndoOperation SET IsSuccess = ?, ErrorMessage = ?, CompletedAt = ? WHERE UndoId = ?`
	_, err := ExecWrapper(db.conn, query, boolToInt(isSuccess), errMsg, now, undoId).Destruct()
	if err != nil {
		return apperror.WrapSimple(err, "update undo status")
	}

	return nil
}
