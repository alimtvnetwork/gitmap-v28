// Package store — prompt_backup_repo.go implements the domain repository for prompt backups.
package store

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	_ "modernc.org/sqlite"
)

var (
	_ IPromptBackupRepository                = (*PromptBackupRepository)(nil)
	_ IRepository[PromptBackupBatch, string] = (*PromptBackupRepository)(nil)
)

// PromptBackupRepository provides SQLite-backed repository operations for prompt snapshots.
type PromptBackupRepository struct {
	db   *sql.DB
	path string
}

// NewPromptBackupRepository initializes a PromptBackupRepository with default database path.
func NewPromptBackupRepository() (*PromptBackupRepository, *appfault.AppError) {
	return NewPromptBackupRepositoryWithPath("")
}

// NewPromptBackupRepositoryWithPath initializes a repository at a specific path or default if empty.
func NewPromptBackupRepositoryWithPath(customPath string) (*PromptBackupRepository, *appfault.AppError) {
	dbPath := resolveBackupPromptsPath(customPath)
	dirPath := filepath.Dir(dbPath)
	if err := os.MkdirAll(dirPath, 0755); err != nil {
		return nil, NewRepoError("NewPromptBackupRepository.Mkdir", err, map[string]any{"path": dirPath})
	}
	conn, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return nil, NewRepoError("NewPromptBackupRepository.Open", err, map[string]any{"path": dbPath})
	}
	return initPromptBackupRepo(conn, dbPath)
}

// NewPromptBackupRepositoryWithDB wraps an existing *sql.DB connection.
func NewPromptBackupRepositoryWithDB(db *sql.DB, path string) (*PromptBackupRepository, *appfault.AppError) {
	return initPromptBackupRepo(db, path)
}

func initPromptBackupRepo(conn *sql.DB, dbPath string) (*PromptBackupRepository, *appfault.AppError) {
	if conn == nil {
		return nil, NewRepoValidationError("initPromptBackupRepo", "database connection cannot be nil")
	}
	if err := ConfigureSQLiteConn(conn); err != nil {
		_ = conn.Close()
		return nil, NewRepoError("initPromptBackupRepo.Configure", err, map[string]any{"path": dbPath})
	}
	repo := &PromptBackupRepository{db: conn, path: dbPath}
	if appErr := repo.ensureSchema(); appErr != nil {
		_ = conn.Close()
		return nil, appErr
	}
	return repo, nil
}

func (r *PromptBackupRepository) ensureSchema() *appfault.AppError {
	for _, stmt := range []string{sqlCreatePromptBackupBatch, sqlCreatePromptBackupItem, sqlCreatePromptRestoreLedger} {
		if _, err := r.db.Exec(stmt); err != nil {
			return NewRepoError("ensureSchema.Exec", err, nil)
		}
	}
	_, _ = r.db.Exec("ALTER TABLE PromptBackupItem ADD COLUMN MediaPathsJson TEXT NOT NULL DEFAULT '[]'")
	return nil
}

// Close closes the underlying SQLite database connection.
func (r *PromptBackupRepository) Close() *appfault.AppError {
	if r.db == nil {
		return nil
	}
	if err := r.db.Close(); err != nil {
		return NewRepoError("Close", err, map[string]any{"path": r.path})
	}
	return nil
}

// DB returns the underlying sql.DB connection.
func (r *PromptBackupRepository) DB() *sql.DB { return r.db }

// Path returns the SQLite database file path.
func (r *PromptBackupRepository) Path() string { return r.path }

// CreateBatch inserts a backup batch and its associated items into the database.
func (r *PromptBackupRepository) CreateBatch(ctx context.Context, batch *PromptBackupBatch) (*PromptBackupBatch, *appfault.AppError) {
	if batch == nil {
		return nil, NewRepoValidationError("CreateBatch", "batch cannot be nil")
	}
	if batch.BatchID == "" {
		batch.BatchID = fmt.Sprintf("batch-%d", time.Now().UnixNano())
	}
	if batch.CreatedAt == "" {
		batch.CreatedAt = time.Now().UTC().Format(time.RFC3339)
	}
	hasItems := len(batch.Items) > 0
	if hasItems && batch.TotalPrompts == 0 {
		batch.TotalPrompts = len(batch.Items)
	}
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, NewRepoTxError("CreateBatch.BeginTx", err)
	}
	defer tx.Rollback()

	query := `INSERT INTO PromptBackupBatch (BatchId, SourcePath, TotalPrompts, RunningCount, EnqueuedCount, CreatedAt, Note) VALUES (?, ?, ?, ?, ?, ?, ?)`
	_, err = tx.ExecContext(ctx, query, batch.BatchID, batch.SourcePath, batch.TotalPrompts, batch.RunningCount, batch.EnqueuedCount, batch.CreatedAt, batch.Note)
	if err != nil {
		return nil, NewRepoError("CreateBatch.InsertBatch", err, map[string]any{"batchId": batch.BatchID})
	}
	if appErr := r.insertBatchItemsTx(ctx, tx, batch.BatchID, batch.Items); appErr != nil {
		return nil, appErr
	}
	if err := tx.Commit(); err != nil {
		return nil, NewRepoTxError("CreateBatch.Commit", err)
	}
	return batch, nil
}

func (r *PromptBackupRepository) insertBatchItemsTx(ctx context.Context, tx *sql.Tx, batchID string, items []PromptBackupItem) *appfault.AppError {
	query := `INSERT INTO PromptBackupItem (ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`
	for idx, item := range items {
		itemID := item.ItemID
		if itemID == "" {
			itemID = fmt.Sprintf("item-%s-%d", batchID, idx+1)
		}
		itemCreatedAt := item.CreatedAt
		if itemCreatedAt == "" {
			itemCreatedAt = time.Now().UTC().Format(time.RFC3339)
		}
		mediaJSON := marshalMediaPaths(item.MediaPaths)
		_, err := tx.ExecContext(ctx, query, itemID, batchID, item.ProjectName, item.ProjectPath, item.ProjectID, item.ConversationID, item.SequenceID, item.PromptText, item.PromptStatus, item.WordCount, mediaJSON, itemCreatedAt)
		if err != nil {
			return NewRepoError("insertBatchItemsTx.Exec", err, map[string]any{"itemId": itemID, "batchId": batchID})
		}
	}
	return nil
}

// ListBatches returns all recorded backup batches ordered by newest first.
func (r *PromptBackupRepository) ListBatches(ctx context.Context) ([]PromptBackupBatch, *appfault.AppError) {
	query := `SELECT b.BatchId, b.SourcePath, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, b.CreatedAt, b.Note,
		CASE WHEN r.RestoreId IS NOT NULL THEN 'restored' ELSE 'active' END as Status
	FROM PromptBackupBatch b
	LEFT JOIN (SELECT BatchId, RestoreId FROM PromptRestoreLedger GROUP BY BatchId) r ON b.BatchId = r.BatchId
	ORDER BY b.CreatedAt DESC`

	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, NewRepoError("ListBatches.Query", err, nil)
	}
	defer rows.Close()

	var batches []PromptBackupBatch
	for rows.Next() {
		var b PromptBackupBatch
		err := rows.Scan(&b.BatchID, &b.SourcePath, &b.TotalPrompts, &b.RunningCount, &b.EnqueuedCount, &b.CreatedAt, &b.Note, &b.Status)
		if err != nil {
			return nil, NewRepoError("ListBatches.Scan", err, nil)
		}
		batches = append(batches, b)
	}
	if err := rows.Err(); err != nil {
		return nil, NewRepoError("ListBatches.RowsErr", err, nil)
	}
	if batches == nil {
		return []PromptBackupBatch{}, nil
	}
	return batches, nil
}

// GetBatch retrieves a single backup batch by its unique batch ID.
func (r *PromptBackupRepository) GetBatch(ctx context.Context, batchID string) (*PromptBackupBatch, *appfault.AppError) {
	if batchID == "" {
		return nil, NewRepoValidationError("GetBatch", "batchId cannot be empty")
	}
	query := `SELECT b.BatchId, b.SourcePath, b.TotalPrompts, b.RunningCount, b.EnqueuedCount, b.CreatedAt, b.Note,
		CASE WHEN r.RestoreId IS NOT NULL THEN 'restored' ELSE 'active' END as Status
	FROM PromptBackupBatch b
	LEFT JOIN (SELECT BatchId, RestoreId FROM PromptRestoreLedger GROUP BY BatchId) r ON b.BatchId = r.BatchId
	WHERE b.BatchId = ?
	LIMIT 1`

	var b PromptBackupBatch
	err := r.db.QueryRowContext(ctx, query, batchID).Scan(&b.BatchID, &b.SourcePath, &b.TotalPrompts, &b.RunningCount, &b.EnqueuedCount, &b.CreatedAt, &b.Note, &b.Status)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, NewRepoNotFoundError("GetBatch", "PromptBackupBatch", batchID)
	}
	if err != nil {
		return nil, NewRepoError("GetBatch.Scan", err, map[string]any{"batchId": batchID})
	}
	items, appErr := r.GetBatchItems(ctx, batchID)
	if appErr != nil {
		return nil, appErr
	}
	b.Items = items
	return &b, nil
}

// GetBatchItems retrieves all prompt items belonging to the given batch ID.
func (r *PromptBackupRepository) GetBatchItems(ctx context.Context, batchID string) ([]PromptBackupItem, *appfault.AppError) {
	if batchID == "" {
		return nil, NewRepoValidationError("GetBatchItems", "batchId cannot be empty")
	}
	query := `SELECT ItemId, BatchId, ProjectName, ProjectPath, ProjectId, ConversationId, SequenceId, PromptText, PromptStatus, WordCount, MediaPathsJson, CreatedAt
	FROM PromptBackupItem
	WHERE BatchId = ?
	ORDER BY CreatedAt ASC`

	rows, err := r.db.QueryContext(ctx, query, batchID)
	if err != nil {
		return nil, NewRepoError("GetBatchItems.Query", err, map[string]any{"batchId": batchID})
	}
	defer rows.Close()

	var items []PromptBackupItem
	for rows.Next() {
		var item PromptBackupItem
		var mediaJSON string
		err := rows.Scan(&item.ItemID, &item.BatchID, &item.ProjectName, &item.ProjectPath, &item.ProjectID, &item.ConversationID, &item.SequenceID, &item.PromptText, &item.PromptStatus, &item.WordCount, &mediaJSON, &item.CreatedAt)
		if err != nil {
			return nil, NewRepoError("GetBatchItems.Scan", err, map[string]any{"batchId": batchID})
		}
		item.MediaPaths = unmarshalMediaPaths(mediaJSON)
		items = append(items, item)
	}
	if err := rows.Err(); err != nil {
		return nil, NewRepoError("GetBatchItems.RowsErr", err, map[string]any{"batchId": batchID})
	}
	if items == nil {
		return []PromptBackupItem{}, nil
	}
	return items, nil
}

// RecordRestore logs a prompt restore event into PromptRestoreLedger.
func (r *PromptBackupRepository) RecordRestore(ctx context.Context, ledger *PromptRestoreLedger) *appfault.AppError {
	if ledger == nil {
		return NewRepoValidationError("RecordRestore", "ledger cannot be nil")
	}
	if ledger.RestoreID == "" {
		ledger.RestoreID = fmt.Sprintf("r-%d", time.Now().UnixNano())
	}
	if ledger.RestoredAt == "" {
		ledger.RestoredAt = time.Now().UTC().Format(time.RFC3339)
	}
	isKeptInt := 0
	if ledger.IsKept {
		isKeptInt = 1
	}
	query := `INSERT INTO PromptRestoreLedger (RestoreId, BatchId, RestoredAt, RestoredCount, IsKept, PrunedAt) VALUES (?, ?, ?, ?, ?, ?)`
	_, err := r.db.ExecContext(ctx, query, ledger.RestoreID, ledger.BatchID, ledger.RestoredAt, ledger.RestoredCount, isKeptInt, ledger.PrunedAt)
	if err != nil {
		return NewRepoError("RecordRestore.Exec", err, map[string]any{"restoreId": ledger.RestoreID, "batchId": ledger.BatchID})
	}
	return nil
}

// FindByID retrieves a batch by its primary key ID, implementing IRepository.
func (r *PromptBackupRepository) FindByID(ctx context.Context, id string) (*PromptBackupBatch, *appfault.AppError) {
	return r.GetBatch(ctx, id)
}

// FindAll retrieves all recorded batches, implementing IRepository.
func (r *PromptBackupRepository) FindAll(ctx context.Context) ([]PromptBackupBatch, *appfault.AppError) {
	return r.ListBatches(ctx)
}

// Save persists a batch to the repository, implementing IRepository.
func (r *PromptBackupRepository) Save(ctx context.Context, batch *PromptBackupBatch) (*PromptBackupBatch, *appfault.AppError) {
	return r.CreateBatch(ctx, batch)
}

// Delete removes a batch and all associated records, implementing IRepository.
func (r *PromptBackupRepository) Delete(ctx context.Context, id string) *appfault.AppError {
	if id == "" {
		return NewRepoValidationError("Delete", "id cannot be empty")
	}
	query := `DELETE FROM PromptBackupBatch WHERE BatchId = ?`
	_, err := r.db.ExecContext(ctx, query, id)
	if err != nil {
		return NewRepoError("Delete.Exec", err, map[string]any{"batchId": id})
	}
	return nil
}
