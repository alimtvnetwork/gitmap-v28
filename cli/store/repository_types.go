// Package store — repository_types.go defines generic repository contracts, models, and error envelopes.
package store

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/appfault"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

type (
	// ScanRecordWrap is the canonical single reusable result container for ScanRecord.
	ScanRecordWrap = result.Wrap[model.ScanRecord]

	// ScanRecordSliceWrap is the canonical single reusable collection container for ScanRecord.
	ScanRecordSliceWrap = result.WrapSlice[model.ScanRecord]

	// ScanRecordMapWrap is the canonical single reusable keyed lookup container for ScanRecord.
	ScanRecordMapWrap = result.WrapMap[string, model.ScanRecord]

	// RepoPathWrap is the canonical single reusable result container for repository paths.
	RepoPathWrap = result.Wrap[string]

	// SQLiteDBWrap is the canonical single reusable result container for an opened SQLite connection.
	SQLiteDBWrap = result.Wrap[*sql.DB]

	// ScanRecordsResult is an alias for ScanRecordSliceWrap using canonical WrapSlice.
	ScanRecordsResult = result.WrapSlice[model.ScanRecord]
)

// IRepository defines the generic domain repository interface for CRUD persistence operations.
type IRepository[T any, ID any] interface {
	FindByID(ctx context.Context, id ID) (*T, *appfault.AppError)
	FindAll(ctx context.Context) ([]T, *appfault.AppError)
	Save(ctx context.Context, entity *T) (*T, *appfault.AppError)
	Delete(ctx context.Context, id ID) *appfault.AppError
}

// IPromptBackupRepository defines domain operations for prompt backup batches and items.
type IPromptBackupRepository interface {
	CreateBatch(ctx context.Context, batch *PromptBackupBatch) (*PromptBackupBatch, *appfault.AppError)
	ListBatches(ctx context.Context) ([]PromptBackupBatch, *appfault.AppError)
	GetBatch(ctx context.Context, batchID string) (*PromptBackupBatch, *appfault.AppError)
	GetBatchItems(ctx context.Context, batchID string) ([]PromptBackupItem, *appfault.AppError)
	RecordRestore(ctx context.Context, ledger *PromptRestoreLedger) *appfault.AppError
}

// PromptBackupBatch models a snapshot batch of prompts stored in SQLite.
type PromptBackupBatch struct {
	BatchID       string             `json:"batchId"`
	SourcePath    string             `json:"sourcePath"`
	TotalPrompts  int                `json:"totalPrompts"`
	RunningCount  int                `json:"runningCount"`
	EnqueuedCount int                `json:"enqueuedCount"`
	CreatedAt     string             `json:"createdAt"`
	Note          string             `json:"note"`
	Status        string             `json:"status,omitempty"`
	Node          string             `json:"node,omitempty"`
	Host          string             `json:"host,omitempty"`
	Items         []PromptBackupItem `json:"items,omitempty"`
}

// PromptBackupItem models a single prompt record in a backup batch.
type PromptBackupItem struct {
	ItemID         string   `json:"itemId"`
	BatchID        string   `json:"batchId"`
	ProjectName    string   `json:"projectName"`
	ProjectPath    string   `json:"projectPath"`
	ProjectID      string   `json:"projectId"`
	ConversationID string   `json:"conversationId"`
	SequenceID     string   `json:"sequenceId"`
	PromptText     string   `json:"promptText"`
	PromptStatus   string   `json:"promptStatus"`
	WordCount      int      `json:"wordCount"`
	MediaPaths     []string `json:"mediaPaths,omitempty"`
	CreatedAt      string   `json:"createdAt"`
}

// PromptRestoreLedger models a restore operation record.
type PromptRestoreLedger struct {
	RestoreID     string `json:"restoreId"`
	BatchID       string `json:"batchId"`
	RestoredAt    string `json:"restoredAt"`
	RestoredCount int    `json:"restoredCount"`
	IsKept        bool   `json:"isKept"`
	PrunedAt      string `json:"prunedAt"`
}

// RepoResult encapsulates a single entity result with *appfault.AppError.
type RepoResult[T any] struct {
	Data    *T
	AppErr  *appfault.AppError
	IsFound bool
}

// RepoListResult encapsulates a slice entity result with *appfault.AppError.
type RepoListResult[T any] struct {
	Items  []T
	AppErr *appfault.AppError
	Count  int
}

// IsSuccess reports whether the single result represents a successful operation.
func (r RepoResult[T]) IsSuccess() bool {
	return r.AppErr == nil || r.AppErr.HasNoError()
}

// HasError reports whether the single result contains an error.
func (r RepoResult[T]) HasError() bool {
	return r.AppErr != nil && r.AppErr.HasError()
}

// IsSuccess reports whether the list result represents a successful operation.
func (r RepoListResult[T]) IsSuccess() bool {
	return r.AppErr == nil || r.AppErr.HasNoError()
}

// HasError reports whether the list result contains an error.
func (r RepoListResult[T]) HasError() bool {
	return r.AppErr != nil && r.AppErr.HasError()
}

// NewRepoSuccess creates a successful RepoResult.
func NewRepoSuccess[T any](data *T) RepoResult[T] {
	return RepoResult[T]{
		Data:    data,
		IsFound: data != nil,
	}
}

// NewRepoFailure creates a failed RepoResult.
func NewRepoFailure[T any](err *appfault.AppError) RepoResult[T] {
	return RepoResult[T]{
		AppErr:  err,
		IsFound: false,
	}
}

// NewRepoListSuccess creates a successful RepoListResult.
func NewRepoListSuccess[T any](items []T) RepoListResult[T] {
	return RepoListResult[T]{
		Items: items,
		Count: len(items),
	}
}

// NewRepoListFailure creates a failed RepoListResult.
func NewRepoListFailure[T any](err *appfault.AppError) RepoListResult[T] {
	return RepoListResult[T]{
		Items:  nil,
		AppErr: err,
		Count:  0,
	}
}

// NewRepoError wraps an underlying error into an *appfault.AppError with execution error type.
func NewRepoError(op string, err error, metadata map[string]any) *appfault.AppError {
	if err == nil {
		return nil
	}

	appErr := apperror.Wrap(err, op, metadata)
	appErr.Code = "E_DATABASE_ERROR"
	appErr.Type = apperror.ErrorTypeExecution

	return appErr
}

// NewRepoNotFoundError creates an *appfault.AppError specialized for missing records.
func NewRepoNotFoundError(op string, entity string, id any) *appfault.AppError {
	msg := fmt.Sprintf("%s with id %v not found", entity, id)
	appErr := apperror.NewNotFound(op, "E_NOT_FOUND", msg)
	if appErr.Ctx == nil {
		appErr.Ctx = make(map[string]any)
	}
	appErr.Ctx["entity"] = entity
	appErr.Ctx["id"] = id

	return appErr
}

// NewRepoTxError creates an *appfault.AppError specialized for transaction failures.
func NewRepoTxError(op string, err error) *appfault.AppError {
	if err == nil {
		return nil
	}

	appErr := apperror.WrapSimple(err, op)
	appErr.Code = "E_TRANSACTION_ERROR"
	appErr.Type = apperror.ErrorTypeExecution

	return appErr
}

// NewRepoValidationError creates an *appfault.AppError specialized for validation errors.
func NewRepoValidationError(op string, msg string) *appfault.AppError {
	return apperror.NewValidation(op, "E_VALIDATION", msg)
}
