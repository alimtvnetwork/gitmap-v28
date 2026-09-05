package pipelinedb

import (
	"context"
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

// PipelineRepository provides domain business logic methods for pipeline runs and error diagnostics.
type PipelineRepository struct {
	db   *dbengine.DbWrapper
	repo *dbengine.Repository[PipelineRunRecord, PipelineRunRecordFieldType]
}

// NewPipelineRepository creates a new domain repository wrapping a DbWrapper.
func NewPipelineRepository(db *dbengine.DbWrapper) *PipelineRepository {
	repo := dbengine.NewRepository[PipelineRunRecord, PipelineRunRecordFieldType](
		db,
		PipelineRunRecordTable,
		ScanPipelineRunRecord,
	)
	return &PipelineRepository{
		db:   db,
		repo: repo,
	}
}

// Db returns the underlying DbWrapper.
func (r *PipelineRepository) Db() *dbengine.DbWrapper {
	return r.db
}

// Repo returns the underlying generic Repository.
func (r *PipelineRepository) Repo() *dbengine.Repository[PipelineRunRecord, PipelineRunRecordFieldType] {
	return r.repo
}

// Query returns a fluent QueryBuilder initialized with all standard fields projected.
func (r *PipelineRepository) Query() *dbengine.QueryBuilder[PipelineRunRecord, PipelineRunRecordFieldType] {
	return r.repo.Query().Select(PipelineRunRecordDb.All()...)
}

// QueryBare returns a fluent QueryBuilder without any pre-selected fields.
func (r *PipelineRepository) QueryBare() *dbengine.QueryBuilder[PipelineRunRecord, PipelineRunRecordFieldType] {
	return r.repo.Query()
}

// ScanPipelineRunRecord maps a database row scanner to a PipelineRunRecord entity.
func ScanPipelineRunRecord(row dbengine.RowScanner) (*PipelineRunRecord, error) {
	var r PipelineRunRecord
	var isSuccessInt int
	var notes, comments sql.NullString
	err := row.Scan(
		&r.RunId,
		&r.RepoSlug,
		&r.WorkflowName,
		&r.Status,
		&r.Conclusion,
		&r.Branch,
		&r.Sha,
		&r.EtaSeconds,
		&r.DurationSeconds,
		&r.RunUrl,
		&isSuccessInt,
		&notes,
		&comments,
		&r.CreatedAt,
		&r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	r.IsSuccess = isSuccessInt == 1
	if notes.Valid {
		r.Notes = notes.String
	}
	if comments.Valid {
		r.Comments = comments.String
	}
	return &r, nil
}

// ScanPipelineErrorRecord maps a database row scanner to a PipelineErrorRecord entity.
func ScanPipelineErrorRecord(row dbengine.RowScanner) (*PipelineErrorRecord, error) {
	var e PipelineErrorRecord
	var rawLogs, notes, comments sql.NullString
	err := row.Scan(
		&e.RunId,
		&e.RepoSlug,
		&e.WorkflowName,
		&e.StepName,
		&e.ErrorText,
		&rawLogs,
		&notes,
		&comments,
		&e.CreatedAt,
	)
	if err != nil {
		return nil, err
	}
	if rawLogs.Valid {
		e.RawLogs = rawLogs.String
	}
	if notes.Valid {
		e.Notes = notes.String
	}
	if comments.Valid {
		e.Comments = comments.String
	}
	return &e, nil
}

// InitSchema ensures the PipelineRunRecord and PipelineErrorRecord tables exist.
func (r *PipelineRepository) InitSchema(ctx context.Context) *apperror.AppError {
	createRunQuery := `
CREATE TABLE IF NOT EXISTS PipelineRunRecord (
    RunId INTEGER NOT NULL UNIQUE,
    RepoSlug TEXT NOT NULL,
    WorkflowName TEXT NOT NULL,
    Status TEXT NOT NULL,
    Conclusion TEXT NOT NULL,
    Branch TEXT NOT NULL,
    Sha TEXT NOT NULL,
    EtaSeconds INTEGER DEFAULT 0,
    DurationSeconds INTEGER DEFAULT 0,
    RunUrl TEXT NOT NULL,
    IsSuccess INTEGER DEFAULT 0,
    Notes TEXT NULL,
    Comments TEXT NULL,
    CreatedAt TEXT NOT NULL,
    UpdatedAt TEXT NOT NULL
);`
	_, runErr := r.db.Exec(ctx, createRunQuery)
	if runErr != nil {
		return runErr
	}

	createErrQuery := `
CREATE TABLE IF NOT EXISTS PipelineErrorRecord (
    RunId INTEGER NOT NULL,
    RepoSlug TEXT NOT NULL,
    WorkflowName TEXT NOT NULL,
    StepName TEXT NOT NULL,
    ErrorText TEXT NOT NULL,
    RawLogs TEXT NULL,
    Notes TEXT NULL,
    Comments TEXT NULL,
    CreatedAt TEXT DEFAULT CURRENT_TIMESTAMP
);`
	_, errLogErr := r.db.Exec(ctx, createErrQuery)
	if errLogErr != nil {
		return errLogErr
	}

	return nil
}

// InsertRun inserts a new PipelineRunRecord into the database.
func (r *PipelineRepository) InsertRun(ctx context.Context, run PipelineRunRecord) dbengine.RowsAffectedResult {
	query := `
INSERT INTO PipelineRunRecord (
    RunId, RepoSlug, WorkflowName, Status, Conclusion, Branch, Sha,
    EtaSeconds, DurationSeconds, RunUrl, IsSuccess, Notes, Comments, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`
	isSuccessInt := 0
	if run.IsSuccess {
		isSuccessInt = 1
	}
	return r.db.ExecRowsAffected(ctx, query,
		run.RunId, run.RepoSlug, run.WorkflowName, run.Status, run.Conclusion,
		run.Branch, run.Sha, run.EtaSeconds, run.DurationSeconds, run.RunUrl,
		isSuccessInt, run.Notes, run.Comments, run.CreatedAt, run.UpdatedAt,
	)
}

// InsertErrorRecord inserts an isolated diagnostic record into PipelineErrorRecord.
func (r *PipelineRepository) InsertErrorRecord(ctx context.Context, errLog PipelineErrorRecord) dbengine.RowsAffectedResult {
	query := `
INSERT INTO PipelineErrorRecord (
    RunId, RepoSlug, WorkflowName, StepName, ErrorText, RawLogs, Notes, Comments, CreatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`
	return r.db.ExecRowsAffected(ctx, query,
		errLog.RunId, errLog.RepoSlug, errLog.WorkflowName, errLog.StepName,
		errLog.ErrorText, errLog.RawLogs, errLog.Notes, errLog.Comments, errLog.CreatedAt,
	)
}

// GetRunById retrieves a run record by its unique RunId.
func (r *PipelineRepository) GetRunById(ctx context.Context, runId uint64) dbengine.EntityResult[PipelineRunRecord] {
	return r.Query().
		WhereOp(PipelineRunRecordDb.RunId, dbengine.SqlOperators.Equal, runId).
		First(ctx)
}

// GetRecentRuns retrieves recent runs for a repository slug ordered descending by RunId.
func (r *PipelineRepository) GetRecentRuns(ctx context.Context, repoSlug string, limit int) dbengine.ListResult[PipelineRunRecord] {
	return r.Query().
		WhereOp(PipelineRunRecordDb.RepoSlug, dbengine.SqlOperators.Equal, repoSlug).
		OrderByDesc(PipelineRunRecordDb.RunId).
		Limit(limit).
		FindAll(ctx)
}

// EnsureActiveErrorsView creates or reuses the ActiveCiErrors database view using type-safe joins and automated query hashing.
func (r *PipelineRepository) EnsureActiveErrorsView(ctx context.Context) dbengine.BoolResult {
	return r.repo.Query().
		Select(PipelineRunRecordDb.RunId, PipelineRunRecordDb.RepoSlug, PipelineRunRecordDb.WorkflowName).
		InnerJoin(PipelineErrorRecordTable).
		Select(PipelineErrorRecordDb.StepName, PipelineErrorRecordDb.ErrorText).
		OnField(PipelineRunRecordDb.RunId, dbengine.SqlOperators.Equal, PipelineErrorRecordDb.RunId).
		WhereOp(PipelineRunRecordDb.Status, dbengine.SqlOperators.Equal, "completed").
		CreateViewOrUseView(ctx, "ActiveCiErrors")
}
