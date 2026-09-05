package pipelinedb

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/dbengine"
)

// PipelineDbStats encapsulates health and sizing metrics for the pipeline split database.
type PipelineDbStats struct {
	Path          string `json:"path"`
	Size          uint64 `json:"size"`
	TotalRuns     int    `json:"totalRuns"`
	SuccessRuns   int    `json:"successRuns"`
	FailedRuns    int    `json:"failedRuns"`
	ErrorLogCount int    `json:"errorLogCount"`
	SegmentCount  int    `json:"segmentCount"`
	LastUpdated   string `json:"lastUpdated"`
}

// PipelineDBStats is an alias to PipelineDbStats for backward compatibility.
type PipelineDBStats = PipelineDbStats

// PipelineDbStatsFieldType represents column name enums for PipelineDbStats.
type PipelineDbStatsFieldType string

// Name returns the identifier name of the field enum.
func (e PipelineDbStatsFieldType) Name() string {
	return string(e)
}

// String returns the string representation of the field enum.
func (e PipelineDbStatsFieldType) String() string {
	return string(e)
}

// Value returns the raw string value of the field enum.
func (e PipelineDbStatsFieldType) Value() string {
	return string(e)
}

// IsCompare checks equality against another field enum object.
func (e PipelineDbStatsFieldType) IsCompare(target PipelineDbStatsFieldType) bool {
	return e == target
}

// IsEnum checks whether this field enum exists in the valid enum map.
func (e PipelineDbStatsFieldType) IsEnum() bool {
	return pipelineDbStatsValidMap[e]
}

// MarshalJSON implements json.Marshaler.
func (e PipelineDbStatsFieldType) MarshalJSON() ([]byte, error) {
	return json.Marshal(string(e))
}

// UnmarshalJSON implements json.Unmarshaler with strict map validation.
func (e *PipelineDbStatsFieldType) UnmarshalJSON(data []byte) error {
	var s string
	if err := json.Unmarshal(data, &s); err != nil {
		return err
	}
	target := PipelineDbStatsFieldType(s)
	if !pipelineDbStatsValidMap[target] {
		return fmt.Errorf("invalid %s enum: %s", "PipelineDbStatsFieldType", s)
	}
	*e = target
	return nil
}

// ToJSON converts the field enum to a JSON string representation, returning an AppError on failure.
func (e PipelineDbStatsFieldType) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(string(e))
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize field to json")
	}
	return string(b), nil
}

// FromJSON parses a field enum from a JSON string representation, returning an AppError on failure.
func (e *PipelineDbStatsFieldType) FromJSON(s string) *apperror.AppError {
	var str string
	if err := json.Unmarshal([]byte(s), &str); err != nil {
		return apperror.WrapSimple(err, "deserialize field from json")
	}
	target := PipelineDbStatsFieldType(str)
	if !pipelineDbStatsValidMap[target] {
		return apperror.WrapSimple(fmt.Errorf("invalid %s enum: %s", "PipelineDbStatsFieldType", str), "validate field enum from json")
	}
	*e = target
	return nil
}

// IsPath checks whether this field enum instance is Path.
func (e PipelineDbStatsFieldType) IsPath() bool {
	return e == PipelineDbStatsDb.Path
}

// IsSize checks whether this field enum instance is Size.
func (e PipelineDbStatsFieldType) IsSize() bool {
	return e == PipelineDbStatsDb.Size
}

// IsTotalRuns checks whether this field enum instance is TotalRuns.
func (e PipelineDbStatsFieldType) IsTotalRuns() bool {
	return e == PipelineDbStatsDb.TotalRuns
}

// IsSuccessRuns checks whether this field enum instance is SuccessRuns.
func (e PipelineDbStatsFieldType) IsSuccessRuns() bool {
	return e == PipelineDbStatsDb.SuccessRuns
}

// IsFailedRuns checks whether this field enum instance is FailedRuns.
func (e PipelineDbStatsFieldType) IsFailedRuns() bool {
	return e == PipelineDbStatsDb.FailedRuns
}

// IsErrorLogCount checks whether this field enum instance is ErrorLogCount.
func (e PipelineDbStatsFieldType) IsErrorLogCount() bool {
	return e == PipelineDbStatsDb.ErrorLogCount
}

// IsSegmentCount checks whether this field enum instance is SegmentCount.
func (e PipelineDbStatsFieldType) IsSegmentCount() bool {
	return e == PipelineDbStatsDb.SegmentCount
}

// IsLastUpdated checks whether this field enum instance is LastUpdated.
func (e PipelineDbStatsFieldType) IsLastUpdated() bool {
	return e == PipelineDbStatsDb.LastUpdated
}

type pipelineDbStatsDbRegistry struct {
	Path PipelineDbStatsFieldType
	Size PipelineDbStatsFieldType
	TotalRuns PipelineDbStatsFieldType
	SuccessRuns PipelineDbStatsFieldType
	FailedRuns PipelineDbStatsFieldType
	ErrorLogCount PipelineDbStatsFieldType
	SegmentCount PipelineDbStatsFieldType
	LastUpdated PipelineDbStatsFieldType
}

// All returns a slice of all field enums in PipelineDbStats.
func (r pipelineDbStatsDbRegistry) All() []PipelineDbStatsFieldType {
	return []PipelineDbStatsFieldType{
		r.Path,
		r.Size,
		r.TotalRuns,
		r.SuccessRuns,
		r.FailedRuns,
		r.ErrorLogCount,
		r.SegmentCount,
		r.LastUpdated,
	}
}

// Names returns a slice of string names for all fields in PipelineDbStats.
func (r pipelineDbStatsDbRegistry) Names() []string {
	return []string{
		"Path",
		"Size",
		"TotalRuns",
		"SuccessRuns",
		"FailedRuns",
		"ErrorLogCount",
		"SegmentCount",
		"LastUpdated",
	}
}

// IsEnum checks whether the target object matches any registered field enum in PipelineDbStats.
func (r pipelineDbStatsDbRegistry) IsEnum(target PipelineDbStatsFieldType) bool {
	return pipelineDbStatsValidMap[target]
}

// IsPath checks whether the target object is Path.
func (r pipelineDbStatsDbRegistry) IsPath(target PipelineDbStatsFieldType) bool {
	return target == r.Path
}

// IsSize checks whether the target object is Size.
func (r pipelineDbStatsDbRegistry) IsSize(target PipelineDbStatsFieldType) bool {
	return target == r.Size
}

// IsTotalRuns checks whether the target object is TotalRuns.
func (r pipelineDbStatsDbRegistry) IsTotalRuns(target PipelineDbStatsFieldType) bool {
	return target == r.TotalRuns
}

// IsSuccessRuns checks whether the target object is SuccessRuns.
func (r pipelineDbStatsDbRegistry) IsSuccessRuns(target PipelineDbStatsFieldType) bool {
	return target == r.SuccessRuns
}

// IsFailedRuns checks whether the target object is FailedRuns.
func (r pipelineDbStatsDbRegistry) IsFailedRuns(target PipelineDbStatsFieldType) bool {
	return target == r.FailedRuns
}

// IsErrorLogCount checks whether the target object is ErrorLogCount.
func (r pipelineDbStatsDbRegistry) IsErrorLogCount(target PipelineDbStatsFieldType) bool {
	return target == r.ErrorLogCount
}

// IsSegmentCount checks whether the target object is SegmentCount.
func (r pipelineDbStatsDbRegistry) IsSegmentCount(target PipelineDbStatsFieldType) bool {
	return target == r.SegmentCount
}

// IsLastUpdated checks whether the target object is LastUpdated.
func (r pipelineDbStatsDbRegistry) IsLastUpdated(target PipelineDbStatsFieldType) bool {
	return target == r.LastUpdated
}

// ToJSON converts the field registry to a JSON string representation, returning an AppError on failure.
func (r pipelineDbStatsDbRegistry) ToJSON() (string, *apperror.AppError) {
	b, err := json.Marshal(r)
	if err != nil {
		return "", apperror.WrapSimple(err, "serialize registry to json")
	}
	return string(b), nil
}

// PipelineDbStatsDb provides scoped access to field enums: PipelineDbStatsDb.<Field>.
var PipelineDbStatsDb = pipelineDbStatsDbRegistry{
	Path: "Path",
	Size: "Size",
	TotalRuns: "TotalRuns",
	SuccessRuns: "SuccessRuns",
	FailedRuns: "FailedRuns",
	ErrorLogCount: "ErrorLogCount",
	SegmentCount: "SegmentCount",
	LastUpdated: "LastUpdated",
}

// pipelineDbStatsValidMap provides O(1) map validation for field enums.
var pipelineDbStatsValidMap = map[PipelineDbStatsFieldType]bool{
	PipelineDbStatsDb.Path: true,
	PipelineDbStatsDb.Size: true,
	PipelineDbStatsDb.TotalRuns: true,
	PipelineDbStatsDb.SuccessRuns: true,
	PipelineDbStatsDb.FailedRuns: true,
	PipelineDbStatsDb.ErrorLogCount: true,
	PipelineDbStatsDb.SegmentCount: true,
	PipelineDbStatsDb.LastUpdated: true,
}

// PipelineDbStatsField is an alias to PipelineDbStatsDb.
var PipelineDbStatsField = PipelineDbStatsDb

// ScanPipelineDbStats maps a database row scanner to a PipelineDbStats entity.
func ScanPipelineDbStats(row dbengine.RowScanner) (*PipelineDbStats, error) {
	var item PipelineDbStats
	var (
		raw_Path any
		raw_Size any
		raw_TotalRuns any
		raw_SuccessRuns any
		raw_FailedRuns any
		raw_ErrorLogCount any
		raw_SegmentCount any
		raw_LastUpdated any
	)
	err := row.Scan(
		&raw_Path,
		&raw_Size,
		&raw_TotalRuns,
		&raw_SuccessRuns,
		&raw_FailedRuns,
		&raw_ErrorLogCount,
		&raw_SegmentCount,
		&raw_LastUpdated,
	)
	if err != nil {
		return nil, err
	}

	item.Path = dbengine.ScanString(raw_Path)
	item.Size = dbengine.ScanUint64(raw_Size)
	item.TotalRuns = dbengine.ScanInt(raw_TotalRuns)
	item.SuccessRuns = dbengine.ScanInt(raw_SuccessRuns)
	item.FailedRuns = dbengine.ScanInt(raw_FailedRuns)
	item.ErrorLogCount = dbengine.ScanInt(raw_ErrorLogCount)
	item.SegmentCount = dbengine.ScanInt(raw_SegmentCount)
	item.LastUpdated = dbengine.ScanString(raw_LastUpdated)
	return &item, nil
}

// PipelineDbStatsDbRepo provides typed database repository access for PipelineDbStats.
type PipelineDbStatsDbRepo struct {
	db   *dbengine.DbWrapper
	repo *PipelineDbStatsRepository
}

// NewPipelineDbStatsDbRepo initializes a typed repository for PipelineDbStats.
func NewPipelineDbStatsDbRepo(db *dbengine.DbWrapper) *PipelineDbStatsDbRepo {
	repo := dbengine.NewRepository[PipelineDbStats, PipelineDbStatsFieldType](
		db,
		PipelineDbStatsTable,
		ScanPipelineDbStats,
	)
	return &PipelineDbStatsDbRepo{
		db:   db,
		repo: repo,
	}
}

// Db returns the underlying DbWrapper.
func (r *PipelineDbStatsDbRepo) Db() *dbengine.DbWrapper {
	return r.db
}

// Repo returns the underlying generic Repository.
func (r *PipelineDbStatsDbRepo) Repo() *PipelineDbStatsRepository {
	return r.repo
}

// Query returns a fluent QueryBuilder initialized with all standard fields projected.
func (r *PipelineDbStatsDbRepo) Query() *PipelineDbStatsQueryBuilder {
	return r.repo.Query().Select(PipelineDbStatsDb.All()...)
}

// QueryBare returns a fluent QueryBuilder without any pre-selected fields.
func (r *PipelineDbStatsDbRepo) QueryBare() *PipelineDbStatsQueryBuilder {
	return r.repo.Query()
}

// FindAll executes the query selecting all fields and returns a ListResult envelope.
func (r *PipelineDbStatsDbRepo) FindAll(ctx context.Context) dbengine.ListResult[PipelineDbStats] {
	return r.Query().FindAll(ctx)
}

// First executes the query selecting all fields and returns the first record in an EntityResult envelope.
func (r *PipelineDbStatsDbRepo) First(ctx context.Context) dbengine.EntityResult[PipelineDbStats] {
	return r.Query().First(ctx)
}

// Count returns the total number of records matching the query.
func (r *PipelineDbStatsDbRepo) Count(ctx context.Context) dbengine.Int64Result {
	return r.Query().Count(ctx)
}
