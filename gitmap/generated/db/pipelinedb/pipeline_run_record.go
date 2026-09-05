package pipelinedb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

// PipelineRunRecord represents a recorded workflow run in the pipeline database.
type PipelineRunRecord struct {
	RunId           uint64 `json:"runId"`
	RepoSlug        string `json:"repoSlug"`
	WorkflowName    string `json:"workflowName"`
	Status          string `json:"status"`
	Conclusion      string `json:"conclusion"`
	Branch          string `json:"branch"`
	Sha             string `json:"sha"`
	EtaSeconds      int    `json:"etaSeconds"`
	DurationSeconds int    `json:"durationSeconds"`
	RunUrl          string `json:"runUrl"`
	IsSuccess       bool   `json:"isSuccess"`
	Notes           string `json:"notes,omitempty"`
	Comments        string `json:"comments,omitempty"`
	CreatedAt       string `json:"createdAt"`
	UpdatedAt       string `json:"updatedAt"`
}

// PipelineRunRecordFieldType represents column name enums for PipelineRunRecord.
type PipelineRunRecordFieldType string

// Name returns the identifier name of the field enum.
func (e PipelineRunRecordFieldType) Name() string {
	return string(e)
}

// String returns the string representation of the field enum.
func (e PipelineRunRecordFieldType) String() string {
	return string(e)
}

// Value returns the raw string value of the field enum.
func (e PipelineRunRecordFieldType) Value() string {
	return string(e)
}

// IsCompare checks equality against another field enum object.
func (e PipelineRunRecordFieldType) IsCompare(target PipelineRunRecordFieldType) bool {
	return e == target
}

// IsEnum checks whether this field enum exists in the valid enum map.
func (e PipelineRunRecordFieldType) IsEnum() bool {
	return pipelineRunRecordValidMap[e]
}

// MarshalJSON implements json.Marshaler.
func (e PipelineRunRecordFieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// UnmarshalJSON implements json.Unmarshaler with strict map validation.
func (e *PipelineRunRecordFieldType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	target := PipelineRunRecordFieldType(s)
	if !pipelineRunRecordValidMap[target] {
		return fmt.Errorf("invalid %s enum: %s", "PipelineRunRecordFieldType", s)
	}
	*e = target
	return nil
}

// ToJSON converts the field enum to a JSON string representation, returning an AppError on failure.
func (e PipelineRunRecordFieldType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(e))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize field to json")
	}
	return string(b), nil
}

// FromJSON parses a field enum from a JSON string representation, returning an AppError on failure.
func (e *PipelineRunRecordFieldType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize field from json")
	}
	target := PipelineRunRecordFieldType(str)
	if !pipelineRunRecordValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid %s enum: %s", "PipelineRunRecordFieldType", str), "validate field enum from json")
	}
	*e = target
	return nil
}

// IsRunId checks whether this field enum instance is RunId.
func (e PipelineRunRecordFieldType) IsRunId() bool {
	return e == PipelineRunRecordDb.RunId
}

// IsRepoSlug checks whether this field enum instance is RepoSlug.
func (e PipelineRunRecordFieldType) IsRepoSlug() bool {
	return e == PipelineRunRecordDb.RepoSlug
}

// IsWorkflowName checks whether this field enum instance is WorkflowName.
func (e PipelineRunRecordFieldType) IsWorkflowName() bool {
	return e == PipelineRunRecordDb.WorkflowName
}

// IsStatus checks whether this field enum instance is Status.
func (e PipelineRunRecordFieldType) IsStatus() bool {
	return e == PipelineRunRecordDb.Status
}

// IsConclusion checks whether this field enum instance is Conclusion.
func (e PipelineRunRecordFieldType) IsConclusion() bool {
	return e == PipelineRunRecordDb.Conclusion
}

// IsBranch checks whether this field enum instance is Branch.
func (e PipelineRunRecordFieldType) IsBranch() bool {
	return e == PipelineRunRecordDb.Branch
}

// IsSha checks whether this field enum instance is Sha.
func (e PipelineRunRecordFieldType) IsSha() bool {
	return e == PipelineRunRecordDb.Sha
}

// IsEtaSeconds checks whether this field enum instance is EtaSeconds.
func (e PipelineRunRecordFieldType) IsEtaSeconds() bool {
	return e == PipelineRunRecordDb.EtaSeconds
}

// IsDurationSeconds checks whether this field enum instance is DurationSeconds.
func (e PipelineRunRecordFieldType) IsDurationSeconds() bool {
	return e == PipelineRunRecordDb.DurationSeconds
}

// IsRunUrl checks whether this field enum instance is RunUrl.
func (e PipelineRunRecordFieldType) IsRunUrl() bool {
	return e == PipelineRunRecordDb.RunUrl
}

// IsIsSuccess checks whether this field enum instance is IsSuccess.
func (e PipelineRunRecordFieldType) IsIsSuccess() bool {
	return e == PipelineRunRecordDb.IsSuccess
}

// IsNotes checks whether this field enum instance is Notes.
func (e PipelineRunRecordFieldType) IsNotes() bool {
	return e == PipelineRunRecordDb.Notes
}

// IsComments checks whether this field enum instance is Comments.
func (e PipelineRunRecordFieldType) IsComments() bool {
	return e == PipelineRunRecordDb.Comments
}

// IsCreatedAt checks whether this field enum instance is CreatedAt.
func (e PipelineRunRecordFieldType) IsCreatedAt() bool {
	return e == PipelineRunRecordDb.CreatedAt
}

// IsUpdatedAt checks whether this field enum instance is UpdatedAt.
func (e PipelineRunRecordFieldType) IsUpdatedAt() bool {
	return e == PipelineRunRecordDb.UpdatedAt
}

type pipelineRunRecordDbRegistry struct {
	RunId PipelineRunRecordFieldType
	RepoSlug PipelineRunRecordFieldType
	WorkflowName PipelineRunRecordFieldType
	Status PipelineRunRecordFieldType
	Conclusion PipelineRunRecordFieldType
	Branch PipelineRunRecordFieldType
	Sha PipelineRunRecordFieldType
	EtaSeconds PipelineRunRecordFieldType
	DurationSeconds PipelineRunRecordFieldType
	RunUrl PipelineRunRecordFieldType
	IsSuccess PipelineRunRecordFieldType
	Notes PipelineRunRecordFieldType
	Comments PipelineRunRecordFieldType
	CreatedAt PipelineRunRecordFieldType
	UpdatedAt PipelineRunRecordFieldType
}

// All returns a slice of all field enums in PipelineRunRecord.
func (r pipelineRunRecordDbRegistry) All() []PipelineRunRecordFieldType {
	return []PipelineRunRecordFieldType{
		r.RunId,
		r.RepoSlug,
		r.WorkflowName,
		r.Status,
		r.Conclusion,
		r.Branch,
		r.Sha,
		r.EtaSeconds,
		r.DurationSeconds,
		r.RunUrl,
		r.IsSuccess,
		r.Notes,
		r.Comments,
		r.CreatedAt,
		r.UpdatedAt,
	}
}

// Names returns a slice of string names for all fields in PipelineRunRecord.
func (r pipelineRunRecordDbRegistry) Names() []string {
	return []string{
		"RunId",
		"RepoSlug",
		"WorkflowName",
		"Status",
		"Conclusion",
		"Branch",
		"Sha",
		"EtaSeconds",
		"DurationSeconds",
		"RunUrl",
		"IsSuccess",
		"Notes",
		"Comments",
		"CreatedAt",
		"UpdatedAt",
	}
}

// IsEnum checks whether the target object matches any registered field enum in PipelineRunRecord.
func (r pipelineRunRecordDbRegistry) IsEnum(target PipelineRunRecordFieldType) bool {
	return pipelineRunRecordValidMap[target]
}

// IsRunId checks whether the target object is RunId.
func (r pipelineRunRecordDbRegistry) IsRunId(target PipelineRunRecordFieldType) bool {
	return target == r.RunId
}

// IsRepoSlug checks whether the target object is RepoSlug.
func (r pipelineRunRecordDbRegistry) IsRepoSlug(target PipelineRunRecordFieldType) bool {
	return target == r.RepoSlug
}

// IsWorkflowName checks whether the target object is WorkflowName.
func (r pipelineRunRecordDbRegistry) IsWorkflowName(target PipelineRunRecordFieldType) bool {
	return target == r.WorkflowName
}

// IsStatus checks whether the target object is Status.
func (r pipelineRunRecordDbRegistry) IsStatus(target PipelineRunRecordFieldType) bool {
	return target == r.Status
}

// IsConclusion checks whether the target object is Conclusion.
func (r pipelineRunRecordDbRegistry) IsConclusion(target PipelineRunRecordFieldType) bool {
	return target == r.Conclusion
}

// IsBranch checks whether the target object is Branch.
func (r pipelineRunRecordDbRegistry) IsBranch(target PipelineRunRecordFieldType) bool {
	return target == r.Branch
}

// IsSha checks whether the target object is Sha.
func (r pipelineRunRecordDbRegistry) IsSha(target PipelineRunRecordFieldType) bool {
	return target == r.Sha
}

// IsEtaSeconds checks whether the target object is EtaSeconds.
func (r pipelineRunRecordDbRegistry) IsEtaSeconds(target PipelineRunRecordFieldType) bool {
	return target == r.EtaSeconds
}

// IsDurationSeconds checks whether the target object is DurationSeconds.
func (r pipelineRunRecordDbRegistry) IsDurationSeconds(target PipelineRunRecordFieldType) bool {
	return target == r.DurationSeconds
}

// IsRunUrl checks whether the target object is RunUrl.
func (r pipelineRunRecordDbRegistry) IsRunUrl(target PipelineRunRecordFieldType) bool {
	return target == r.RunUrl
}

// IsIsSuccess checks whether the target object is IsSuccess.
func (r pipelineRunRecordDbRegistry) IsIsSuccess(target PipelineRunRecordFieldType) bool {
	return target == r.IsSuccess
}

// IsNotes checks whether the target object is Notes.
func (r pipelineRunRecordDbRegistry) IsNotes(target PipelineRunRecordFieldType) bool {
	return target == r.Notes
}

// IsComments checks whether the target object is Comments.
func (r pipelineRunRecordDbRegistry) IsComments(target PipelineRunRecordFieldType) bool {
	return target == r.Comments
}

// IsCreatedAt checks whether the target object is CreatedAt.
func (r pipelineRunRecordDbRegistry) IsCreatedAt(target PipelineRunRecordFieldType) bool {
	return target == r.CreatedAt
}

// IsUpdatedAt checks whether the target object is UpdatedAt.
func (r pipelineRunRecordDbRegistry) IsUpdatedAt(target PipelineRunRecordFieldType) bool {
	return target == r.UpdatedAt
}

// ToJSON converts the field registry to a JSON string representation, returning an AppError on failure.
func (r pipelineRunRecordDbRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize registry to json")
	}
	return string(b), nil
}

// PipelineRunRecordDb provides scoped access to field enums: PipelineRunRecordDb.<Field>.
var PipelineRunRecordDb = pipelineRunRecordDbRegistry{
	RunId: "RunId",
	RepoSlug: "RepoSlug",
	WorkflowName: "WorkflowName",
	Status: "Status",
	Conclusion: "Conclusion",
	Branch: "Branch",
	Sha: "Sha",
	EtaSeconds: "EtaSeconds",
	DurationSeconds: "DurationSeconds",
	RunUrl: "RunUrl",
	IsSuccess: "IsSuccess",
	Notes: "Notes",
	Comments: "Comments",
	CreatedAt: "CreatedAt",
	UpdatedAt: "UpdatedAt",
}

// pipelineRunRecordValidMap provides O(1) map validation for field enums.
var pipelineRunRecordValidMap = map[PipelineRunRecordFieldType]bool{
	PipelineRunRecordDb.RunId: true,
	PipelineRunRecordDb.RepoSlug: true,
	PipelineRunRecordDb.WorkflowName: true,
	PipelineRunRecordDb.Status: true,
	PipelineRunRecordDb.Conclusion: true,
	PipelineRunRecordDb.Branch: true,
	PipelineRunRecordDb.Sha: true,
	PipelineRunRecordDb.EtaSeconds: true,
	PipelineRunRecordDb.DurationSeconds: true,
	PipelineRunRecordDb.RunUrl: true,
	PipelineRunRecordDb.IsSuccess: true,
	PipelineRunRecordDb.Notes: true,
	PipelineRunRecordDb.Comments: true,
	PipelineRunRecordDb.CreatedAt: true,
	PipelineRunRecordDb.UpdatedAt: true,
}

// PipelineRunRecordField is an alias to PipelineRunRecordDb.
var PipelineRunRecordField = PipelineRunRecordDb

// ScanPipelineRunRecord maps a database row scanner to a PipelineRunRecord entity.
func ScanPipelineRunRecord(row dbengine.RowScanner) (*PipelineRunRecord, error) {
	var item PipelineRunRecord
	var (
		raw_RunId any
		raw_RepoSlug any
		raw_WorkflowName any
		raw_Status any
		raw_Conclusion any
		raw_Branch any
		raw_Sha any
		raw_EtaSeconds any
		raw_DurationSeconds any
		raw_RunUrl any
		raw_IsSuccess any
		raw_Notes any
		raw_Comments any
		raw_CreatedAt any
		raw_UpdatedAt any
	)
	err := row.Scan(
		&raw_RunId,
		&raw_RepoSlug,
		&raw_WorkflowName,
		&raw_Status,
		&raw_Conclusion,
		&raw_Branch,
		&raw_Sha,
		&raw_EtaSeconds,
		&raw_DurationSeconds,
		&raw_RunUrl,
		&raw_IsSuccess,
		&raw_Notes,
		&raw_Comments,
		&raw_CreatedAt,
		&raw_UpdatedAt,
	)
	if err != nil {
		return nil, err
	}

	item.RunId = dbengine.ScanUint64(raw_RunId)
	item.RepoSlug = dbengine.ScanString(raw_RepoSlug)
	item.WorkflowName = dbengine.ScanString(raw_WorkflowName)
	item.Status = dbengine.ScanString(raw_Status)
	item.Conclusion = dbengine.ScanString(raw_Conclusion)
	item.Branch = dbengine.ScanString(raw_Branch)
	item.Sha = dbengine.ScanString(raw_Sha)
	item.EtaSeconds = dbengine.ScanInt(raw_EtaSeconds)
	item.DurationSeconds = dbengine.ScanInt(raw_DurationSeconds)
	item.RunUrl = dbengine.ScanString(raw_RunUrl)
	item.IsSuccess = dbengine.ScanBool(raw_IsSuccess)
	item.Notes = dbengine.ScanString(raw_Notes)
	item.Comments = dbengine.ScanString(raw_Comments)
	item.CreatedAt = dbengine.ScanString(raw_CreatedAt)
	item.UpdatedAt = dbengine.ScanString(raw_UpdatedAt)
	return &item, nil
}

// PipelineRunRecordDbRepo provides typed database repository access for PipelineRunRecord.
type PipelineRunRecordDbRepo struct {
	db   *dbengine.DbWrapper
	repo *PipelineRunRecordRepository
}

// NewPipelineRunRecordDbRepo initializes a typed repository for PipelineRunRecord.
func NewPipelineRunRecordDbRepo(db *dbengine.DbWrapper) *PipelineRunRecordDbRepo {
	repo := dbengine.NewRepository[PipelineRunRecord, PipelineRunRecordFieldType](
		db,
		PipelineRunRecordTable,
		ScanPipelineRunRecord,
	)
	return &PipelineRunRecordDbRepo{
		db:   db,
		repo: repo,
	}
}

// NewPipelineRunDbRepo is an alias constructor for PipelineRunRecordDbRepo.
func NewPipelineRunDbRepo(db *dbengine.DbWrapper) *PipelineRunDbRepo {
	return NewPipelineRunRecordDbRepo(db)
}

// Db returns the underlying DbWrapper.
func (r *PipelineRunRecordDbRepo) Db() *dbengine.DbWrapper {
	return r.db
}

// Repo returns the underlying generic Repository.
func (r *PipelineRunRecordDbRepo) Repo() *PipelineRunRecordRepository {
	return r.repo
}

// Query returns a fluent QueryBuilder initialized with all standard fields projected.
func (r *PipelineRunRecordDbRepo) Query() *PipelineRunRecordQueryBuilder {
	return r.repo.Query().Select(PipelineRunRecordDb.All()...)
}

// QueryBare returns a fluent QueryBuilder without any pre-selected fields.
func (r *PipelineRunRecordDbRepo) QueryBare() *PipelineRunRecordQueryBuilder {
	return r.repo.Query()
}

// FindAll executes the query selecting all fields and returns a ListResult envelope.
func (r *PipelineRunRecordDbRepo) FindAll(ctx context.Context) dbengine.ListResult[PipelineRunRecord] {
	return r.Query().FindAll(ctx)
}

// First executes the query selecting all fields and returns the first record in an EntityResult envelope.
func (r *PipelineRunRecordDbRepo) First(ctx context.Context) dbengine.EntityResult[PipelineRunRecord] {
	return r.Query().First(ctx)
}

// Count returns the total number of records matching the query.
func (r *PipelineRunRecordDbRepo) Count(ctx context.Context) dbengine.Int64Result {
	return r.Query().Count(ctx)
}
