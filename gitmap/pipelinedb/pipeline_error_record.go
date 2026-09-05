package pipelinedb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

// PipelineErrorRecord represents an isolated error diagnostic record.
type PipelineErrorRecord struct {
	RunId        uint64 `json:"runId"`
	RepoSlug     string `json:"repoSlug"`
	WorkflowName string `json:"workflowName"`
	StepName     string `json:"stepName"`
	ErrorText    string `json:"errorText"`
	RawLogs      string `json:"rawLogs,omitempty"`
	Notes        string `json:"notes,omitempty"`
	Comments     string `json:"comments,omitempty"`
	CreatedAt    string `json:"createdAt"`
}

// PipelineErrorRecordFieldType represents column name enums for PipelineErrorRecord.
type PipelineErrorRecordFieldType string

// Name returns the identifier name of the field enum.
func (e PipelineErrorRecordFieldType) Name() string {
	return string(e)
}

// String returns the string representation of the field enum.
func (e PipelineErrorRecordFieldType) String() string {
	return string(e)
}

// Value returns the raw string value of the field enum.
func (e PipelineErrorRecordFieldType) Value() string {
	return string(e)
}

// IsCompare checks equality against another field enum object.
func (e PipelineErrorRecordFieldType) IsCompare(target PipelineErrorRecordFieldType) bool {
	return e == target
}

// IsEnum checks whether this field enum exists in the valid enum map.
func (e PipelineErrorRecordFieldType) IsEnum() bool {
	return pipelineErrorRecordValidMap[e]
}

// MarshalJSON implements json.Marshaler.
func (e PipelineErrorRecordFieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// UnmarshalJSON implements json.Unmarshaler with strict map validation.
func (e *PipelineErrorRecordFieldType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	target := PipelineErrorRecordFieldType(s)
	if !pipelineErrorRecordValidMap[target] {
		return fmt.Errorf("invalid %s enum: %s", "PipelineErrorRecordFieldType", s)
	}
	*e = target
	return nil
}

// ToJSON converts the field enum to a JSON string representation, returning an AppError on failure.
func (e PipelineErrorRecordFieldType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(e))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize field to json")
	}
	return string(b), nil
}

// FromJSON parses a field enum from a JSON string representation, returning an AppError on failure.
func (e *PipelineErrorRecordFieldType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize field from json")
	}
	target := PipelineErrorRecordFieldType(str)
	if !pipelineErrorRecordValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid %s enum: %s", "PipelineErrorRecordFieldType", str), "validate field enum from json")
	}
	*e = target
	return nil
}

// IsRunId checks whether this field enum instance is RunId.
func (e PipelineErrorRecordFieldType) IsRunId() bool {
	return e == PipelineErrorRecordDb.RunId
}

// IsRepoSlug checks whether this field enum instance is RepoSlug.
func (e PipelineErrorRecordFieldType) IsRepoSlug() bool {
	return e == PipelineErrorRecordDb.RepoSlug
}

// IsWorkflowName checks whether this field enum instance is WorkflowName.
func (e PipelineErrorRecordFieldType) IsWorkflowName() bool {
	return e == PipelineErrorRecordDb.WorkflowName
}

// IsStepName checks whether this field enum instance is StepName.
func (e PipelineErrorRecordFieldType) IsStepName() bool {
	return e == PipelineErrorRecordDb.StepName
}

// IsErrorText checks whether this field enum instance is ErrorText.
func (e PipelineErrorRecordFieldType) IsErrorText() bool {
	return e == PipelineErrorRecordDb.ErrorText
}

// IsRawLogs checks whether this field enum instance is RawLogs.
func (e PipelineErrorRecordFieldType) IsRawLogs() bool {
	return e == PipelineErrorRecordDb.RawLogs
}

// IsNotes checks whether this field enum instance is Notes.
func (e PipelineErrorRecordFieldType) IsNotes() bool {
	return e == PipelineErrorRecordDb.Notes
}

// IsComments checks whether this field enum instance is Comments.
func (e PipelineErrorRecordFieldType) IsComments() bool {
	return e == PipelineErrorRecordDb.Comments
}

// IsCreatedAt checks whether this field enum instance is CreatedAt.
func (e PipelineErrorRecordFieldType) IsCreatedAt() bool {
	return e == PipelineErrorRecordDb.CreatedAt
}

type pipelineErrorRecordDbRegistry struct {
	RunId PipelineErrorRecordFieldType
	RepoSlug PipelineErrorRecordFieldType
	WorkflowName PipelineErrorRecordFieldType
	StepName PipelineErrorRecordFieldType
	ErrorText PipelineErrorRecordFieldType
	RawLogs PipelineErrorRecordFieldType
	Notes PipelineErrorRecordFieldType
	Comments PipelineErrorRecordFieldType
	CreatedAt PipelineErrorRecordFieldType
}

// All returns a slice of all field enums in PipelineErrorRecord.
func (r pipelineErrorRecordDbRegistry) All() []PipelineErrorRecordFieldType {
	return []PipelineErrorRecordFieldType{
		r.RunId,
		r.RepoSlug,
		r.WorkflowName,
		r.StepName,
		r.ErrorText,
		r.RawLogs,
		r.Notes,
		r.Comments,
		r.CreatedAt,
	}
}

// Names returns a slice of string names for all fields in PipelineErrorRecord.
func (r pipelineErrorRecordDbRegistry) Names() []string {
	return []string{
		"RunId",
		"RepoSlug",
		"WorkflowName",
		"StepName",
		"ErrorText",
		"RawLogs",
		"Notes",
		"Comments",
		"CreatedAt",
	}
}

// IsEnum checks whether the target object matches any registered field enum in PipelineErrorRecord.
func (r pipelineErrorRecordDbRegistry) IsEnum(target PipelineErrorRecordFieldType) bool {
	return pipelineErrorRecordValidMap[target]
}

// IsRunId checks whether the target object is RunId.
func (r pipelineErrorRecordDbRegistry) IsRunId(target PipelineErrorRecordFieldType) bool {
	return target == r.RunId
}

// IsRepoSlug checks whether the target object is RepoSlug.
func (r pipelineErrorRecordDbRegistry) IsRepoSlug(target PipelineErrorRecordFieldType) bool {
	return target == r.RepoSlug
}

// IsWorkflowName checks whether the target object is WorkflowName.
func (r pipelineErrorRecordDbRegistry) IsWorkflowName(target PipelineErrorRecordFieldType) bool {
	return target == r.WorkflowName
}

// IsStepName checks whether the target object is StepName.
func (r pipelineErrorRecordDbRegistry) IsStepName(target PipelineErrorRecordFieldType) bool {
	return target == r.StepName
}

// IsErrorText checks whether the target object is ErrorText.
func (r pipelineErrorRecordDbRegistry) IsErrorText(target PipelineErrorRecordFieldType) bool {
	return target == r.ErrorText
}

// IsRawLogs checks whether the target object is RawLogs.
func (r pipelineErrorRecordDbRegistry) IsRawLogs(target PipelineErrorRecordFieldType) bool {
	return target == r.RawLogs
}

// IsNotes checks whether the target object is Notes.
func (r pipelineErrorRecordDbRegistry) IsNotes(target PipelineErrorRecordFieldType) bool {
	return target == r.Notes
}

// IsComments checks whether the target object is Comments.
func (r pipelineErrorRecordDbRegistry) IsComments(target PipelineErrorRecordFieldType) bool {
	return target == r.Comments
}

// IsCreatedAt checks whether the target object is CreatedAt.
func (r pipelineErrorRecordDbRegistry) IsCreatedAt(target PipelineErrorRecordFieldType) bool {
	return target == r.CreatedAt
}

// ToJSON converts the field registry to a JSON string representation, returning an AppError on failure.
func (r pipelineErrorRecordDbRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize registry to json")
	}
	return string(b), nil
}

// PipelineErrorRecordDb provides scoped access to field enums: PipelineErrorRecordDb.<Field>.
var PipelineErrorRecordDb = pipelineErrorRecordDbRegistry{
	RunId: "RunId",
	RepoSlug: "RepoSlug",
	WorkflowName: "WorkflowName",
	StepName: "StepName",
	ErrorText: "ErrorText",
	RawLogs: "RawLogs",
	Notes: "Notes",
	Comments: "Comments",
	CreatedAt: "CreatedAt",
}

// pipelineErrorRecordValidMap provides O(1) map validation for field enums.
var pipelineErrorRecordValidMap = map[PipelineErrorRecordFieldType]bool{
	PipelineErrorRecordDb.RunId: true,
	PipelineErrorRecordDb.RepoSlug: true,
	PipelineErrorRecordDb.WorkflowName: true,
	PipelineErrorRecordDb.StepName: true,
	PipelineErrorRecordDb.ErrorText: true,
	PipelineErrorRecordDb.RawLogs: true,
	PipelineErrorRecordDb.Notes: true,
	PipelineErrorRecordDb.Comments: true,
	PipelineErrorRecordDb.CreatedAt: true,
}

// PipelineErrorRecordField is an alias to PipelineErrorRecordDb.
var PipelineErrorRecordField = PipelineErrorRecordDb

// ScanPipelineErrorRecord maps a database row scanner to a PipelineErrorRecord entity.
func ScanPipelineErrorRecord(row dbengine.RowScanner) (*PipelineErrorRecord, error) {
	var item PipelineErrorRecord
	var (
		raw_RunId any
		raw_RepoSlug any
		raw_WorkflowName any
		raw_StepName any
		raw_ErrorText any
		raw_RawLogs any
		raw_Notes any
		raw_Comments any
		raw_CreatedAt any
	)
	err := row.Scan(
		&raw_RunId,
		&raw_RepoSlug,
		&raw_WorkflowName,
		&raw_StepName,
		&raw_ErrorText,
		&raw_RawLogs,
		&raw_Notes,
		&raw_Comments,
		&raw_CreatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.RunId = dbengine.ScanUint64(raw_RunId)
	item.RepoSlug = dbengine.ScanString(raw_RepoSlug)
	item.WorkflowName = dbengine.ScanString(raw_WorkflowName)
	item.StepName = dbengine.ScanString(raw_StepName)
	item.ErrorText = dbengine.ScanString(raw_ErrorText)
	item.RawLogs = dbengine.ScanString(raw_RawLogs)
	item.Notes = dbengine.ScanString(raw_Notes)
	item.Comments = dbengine.ScanString(raw_Comments)
	item.CreatedAt = dbengine.ScanString(raw_CreatedAt)
	return &item, nil
}

// PipelineErrorRecordDbRepo provides typed database repository access for PipelineErrorRecord.
type PipelineErrorRecordDbRepo struct {
	db   *dbengine.DbWrapper
	repo *PipelineErrorRecordRepository
}

// NewPipelineErrorRecordDbRepo initializes a typed repository for PipelineErrorRecord.
func NewPipelineErrorRecordDbRepo(db *dbengine.DbWrapper) *PipelineErrorRecordDbRepo {
	repo := dbengine.NewRepository[PipelineErrorRecord, PipelineErrorRecordFieldType](
		db,
		PipelineErrorRecordTable,
		ScanPipelineErrorRecord,
	)
	return &PipelineErrorRecordDbRepo{
		db:   db,
		repo: repo,
	}
}

// NewPipelineErrorDbRepo is an alias constructor for PipelineErrorRecordDbRepo.
func NewPipelineErrorDbRepo(db *dbengine.DbWrapper) *PipelineErrorDbRepo {
	return NewPipelineErrorRecordDbRepo(db)
}

// Db returns the underlying DbWrapper.
func (r *PipelineErrorRecordDbRepo) Db() *dbengine.DbWrapper {
	return r.db
}

// Repo returns the underlying generic Repository.
func (r *PipelineErrorRecordDbRepo) Repo() *PipelineErrorRecordRepository {
	return r.repo
}

// Query returns a fluent QueryBuilder initialized with all standard fields projected.
func (r *PipelineErrorRecordDbRepo) Query() *PipelineErrorRecordQueryBuilder {
	return r.repo.Query().Select(PipelineErrorRecordDb.All()...)
}

// QueryBare returns a fluent QueryBuilder without any pre-selected fields.
func (r *PipelineErrorRecordDbRepo) QueryBare() *PipelineErrorRecordQueryBuilder {
	return r.repo.Query()
}

// FindAll executes the query selecting all fields and returns a ListResult envelope.
func (r *PipelineErrorRecordDbRepo) FindAll(ctx context.Context) dbengine.ListResult[PipelineErrorRecord] {
	return r.Query().FindAll(ctx)
}

// First executes the query selecting all fields and returns the first record in an EntityResult envelope.
func (r *PipelineErrorRecordDbRepo) First(ctx context.Context) dbengine.EntityResult[PipelineErrorRecord] {
	return r.Query().First(ctx)
}

// Count returns the total number of records matching the query.
func (r *PipelineErrorRecordDbRepo) Count(ctx context.Context) dbengine.Int64Result {
	return r.Query().Count(ctx)
}
