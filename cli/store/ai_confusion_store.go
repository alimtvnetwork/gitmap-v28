package store

import (
	"database/sql"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/uuidutil"
)

// AiConfusionStore handles persistence operations for AI confusion records.
type AiConfusionStore struct {
	db     *sql.DB
	dbPath string
}

// DBPath returns the underlying file path of the SQLite database.
func (s *AiConfusionStore) DBPath() string {
	return s.dbPath
}

// DB returns the underlying *sql.DB handle.
func (s *AiConfusionStore) DB() *sql.DB {
	return s.db
}

// Close closes the underlying SQLite database connection.
func (s *AiConfusionStore) Close() error {
	hasDB := s.db != nil
	if hasDB {
		return s.db.Close()
	}

	return nil
}

// OpenAiConfusionStore opens or initializes the SQLite database for AI confusion episodes.
func OpenAiConfusionStore(repoCacheDir string) (*AiConfusionStore, error) {
	dbPath := ResolveAiConfusionDbPath(repoCacheDir)
	dirErr := ensureConfusionParentDir(dbPath)
	if dirErr != nil {
		return nil, dirErr
	}

	openRes := OpenSQLiteDB(dbPath)
	if openRes.IsFailure() {
		return nil, openRes.AppError()
	}

	db := openRes.Value
	schemaErr := InitAiConfusionSchema(db)
	if schemaErr != nil {
		_ = db.Close()

		return nil, schemaErr
	}

	return &AiConfusionStore{db: db, dbPath: dbPath}, nil
}

func ensureConfusionParentDir(dbPath string) error {
	dir := filepath.Dir(dbPath)
	err := os.MkdirAll(dir, 0755)
	if err != nil {
		return apperror.WrapSimple(err, "ensure ai confusion parent dir")
	}

	return nil
}

const sqlInsertAiConfusion = `
INSERT INTO AiConfusion (
    ConfusionUuid, AgentSlug, FileNames, MdFiles, Spec, Task, Subtask,
    HeatmapInfo, LineNumbers, CodeSnippet, Reason, InitialIntention,
    ReasoningThinking, ConfusionQuestion, IsCompleted, CreatedAt, UpdatedAt
) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

// InsertAiConfusion inserts a new AI confusion episode into the database.
func (s *AiConfusionStore) InsertAiConfusion(rec *AiConfusionRecord) error {
	hasNilRec := rec == nil
	if hasNilRec {
		return apperror.NewValidationError("ai confusion record is nil")
	}

	s.normalizeTimestampsAndUuid(rec)
	args := s.buildInsertArgs(rec)

	return s.executeInsert(rec, args)
}

func (s *AiConfusionStore) executeInsert(rec *AiConfusionRecord, args []any) error {
	res, err := s.db.Exec(sqlInsertAiConfusion, args...)
	if err != nil {
		return apperror.WrapSimple(err, "insert ai confusion record")
	}

	lastId, idErr := res.LastInsertId()
	hasId := idErr == nil
	if hasId {
		rec.AiConfusionId = lastId
	}

	return nil
}

func (s *AiConfusionStore) normalizeTimestampsAndUuid(rec *AiConfusionRecord) {
	hasEmptyUuid := len(rec.ConfusionUuid) == 0
	if hasEmptyUuid {
		rec.ConfusionUuid = uuidutil.GenerateUUID()
	}

	now := time.Now().UTC()
	hasZeroCreated := rec.CreatedAt.IsZero()
	if hasZeroCreated {
		rec.CreatedAt = now
	}

	hasZeroUpdated := rec.UpdatedAt.IsZero()
	if hasZeroUpdated {
		rec.UpdatedAt = now
	}
}

func (s *AiConfusionStore) buildInsertArgs(rec *AiConfusionRecord) []any {
	completedInt := 0
	if rec.IsCompleted {
		completedInt = 1
	}

	return []any{
		rec.ConfusionUuid, rec.AgentSlug, rec.FileNames, rec.MdFiles,
		rec.Spec, rec.Task, rec.Subtask, rec.HeatmapInfo,
		rec.LineNumbers, rec.CodeSnippet, rec.Reason, rec.InitialIntention,
		rec.ReasoningThinking, rec.ConfusionQuestion, completedInt,
		rec.CreatedAt.Format(time.RFC3339), rec.UpdatedAt.Format(time.RFC3339),
	}
}

const sqlSelectConfusionBase = `
SELECT AiConfusionId, ConfusionUuid, AgentSlug, FileNames, MdFiles, Spec, Task, Subtask,
       HeatmapInfo, LineNumbers, CodeSnippet, Reason, InitialIntention,
       ReasoningThinking, ConfusionQuestion, IsCompleted, CreatedAt, UpdatedAt
FROM AiConfusion `

// GetAiConfusionByUuidOrId retrieves a single confusion record by integer ID or UUID string.
func (s *AiConfusionStore) GetAiConfusionByUuidOrId(identifier string) (*AiConfusionRecord, error) {
	cleanId := strings.TrimSpace(identifier)
	hasEmptyId := len(cleanId) == 0
	if hasEmptyId {
		return nil, apperror.NewValidationError("identifier is empty")
	}

	query, args := s.buildGetQuery(cleanId)
	row := s.db.QueryRow(query, args...)

	return s.scanSingleRecord(row)
}

func (s *AiConfusionStore) buildGetQuery(identifier string) (string, []any) {
	idNum, numErr := strconv.ParseInt(identifier, 10, 64)
	isNumeric := numErr == nil
	if isNumeric {
		return sqlSelectConfusionBase + "WHERE AiConfusionId = ? OR ConfusionUuid = ? LIMIT 1;", []any{idNum, identifier}
	}

	return sqlSelectConfusionBase + "WHERE ConfusionUuid = ? LIMIT 1;", []any{identifier}
}

type rowScanner interface {
	Scan(dest ...any) error
}

func (s *AiConfusionStore) scanSingleRecord(scanner rowScanner) (*AiConfusionRecord, error) {
	var rec AiConfusionRecord
	var completedInt int
	var createdStr, updatedStr string

	hasRow, err := s.executeScan(scanner, &rec, &completedInt, &createdStr, &updatedStr)
	if err != nil {
		return nil, err
	}

	if hasRow {
		s.populateScannedMetadata(&rec, completedInt, createdStr, updatedStr)

		return &rec, nil
	}

	return nil, nil
}

func (s *AiConfusionStore) executeScan(sc rowScanner, r *AiConfusionRecord, comp *int, cStr, uStr *string) (bool, error) {
	err := sc.Scan(
		&r.AiConfusionId, &r.ConfusionUuid, &r.AgentSlug, &r.FileNames,
		&r.MdFiles, &r.Spec, &r.Task, &r.Subtask,
		&r.HeatmapInfo, &r.LineNumbers, &r.CodeSnippet, &r.Reason,
		&r.InitialIntention, &r.ReasoningThinking, &r.ConfusionQuestion,
		comp, cStr, uStr,
	)
	isNoRows := err == sql.ErrNoRows
	if isNoRows {
		return false, nil
	}

	if err != nil {
		return false, apperror.WrapSimple(err, "scan ai confusion record")
	}

	return true, nil
}

func (s *AiConfusionStore) populateScannedMetadata(rec *AiConfusionRecord, comp int, cStr, uStr string) {
	rec.IsCompleted = comp == 1
	rec.CreatedAt, _ = time.Parse(time.RFC3339, cStr)
	rec.UpdatedAt, _ = time.Parse(time.RFC3339, uStr)
}

// ListAiConfusions lists confusion episodes filtered by task, agent, status, or id.
func (s *AiConfusionStore) ListAiConfusions(filter AiConfusionFilter) ([]AiConfusionRecord, error) {
	query, args := s.buildListQuery(filter)
	rows, err := s.db.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "list ai confusion records")
	}
	defer rows.Close()

	return s.scanConfusionRows(rows)
}

func (s *AiConfusionStore) buildListQuery(filter AiConfusionFilter) (string, []any) {
	var sb strings.Builder
	var args []any
	sb.WriteString(sqlSelectConfusionBase + "WHERE 1=1")
	s.appendFilterConditions(&sb, &args, filter)
	sb.WriteString(" ORDER BY AiConfusionId DESC")
	s.appendLimitCondition(&sb, &args, filter.Limit)

	return sb.String(), args
}

func (s *AiConfusionStore) appendFilterConditions(sb *strings.Builder, args *[]any, filter AiConfusionFilter) {
	s.appendIdFilter(sb, args, strings.TrimSpace(filter.Id))
	s.appendTaskFilter(sb, args, strings.TrimSpace(filter.Task))
	s.appendAgentFilter(sb, args, strings.TrimSpace(filter.AgentSlug))
	s.appendCompletedFilter(sb, args, filter.IsCompleted)
}

func (s *AiConfusionStore) appendIdFilter(sb *strings.Builder, args *[]any, id string) {
	hasEmptyId := len(id) == 0
	if hasEmptyId {
		return
	}

	idNum, numErr := strconv.ParseInt(id, 10, 64)
	isNumeric := numErr == nil
	if isNumeric {
		sb.WriteString(" AND (AiConfusionId = ? OR ConfusionUuid = ?)")
		*args = append(*args, idNum, id)

		return
	}

	sb.WriteString(" AND ConfusionUuid = ?")
	*args = append(*args, id)
}

func (s *AiConfusionStore) appendTaskFilter(sb *strings.Builder, args *[]any, task string) {
	hasEmptyTask := len(task) == 0
	if hasEmptyTask {
		return
	}

	sb.WriteString(" AND Task = ?")
	*args = append(*args, task)
}

func (s *AiConfusionStore) appendAgentFilter(sb *strings.Builder, args *[]any, agent string) {
	hasEmptyAgent := len(agent) == 0
	if hasEmptyAgent {
		return
	}

	sb.WriteString(" AND AgentSlug = ?")
	*args = append(*args, agent)
}

func (s *AiConfusionStore) appendCompletedFilter(sb *strings.Builder, args *[]any, isComp *bool) {
	hasNilFilter := isComp == nil
	if hasNilFilter {
		return
	}

	val := 0
	if *isComp {
		val = 1
	}

	sb.WriteString(" AND IsCompleted = ?")
	*args = append(*args, val)
}

func (s *AiConfusionStore) appendLimitCondition(sb *strings.Builder, args *[]any, limit int) {
	hasPositiveLimit := limit > 0
	if hasPositiveLimit {
		sb.WriteString(" LIMIT ?")
		*args = append(*args, limit)
	}
}

func (s *AiConfusionStore) scanConfusionRows(rows *sql.Rows) ([]AiConfusionRecord, error) {
	var records []AiConfusionRecord
	for rows.Next() {
		rec, err := s.scanSingleRecord(rows)
		if err != nil {
			return nil, err
		}
		hasRecord := rec != nil
		if hasRecord {
			records = append(records, *rec)
		}
	}

	return records, rows.Err()
}

const sqlCompleteConfusion = `
UPDATE AiConfusion
SET IsCompleted = 1,
    Reason = CASE WHEN ? != '' THEN ? ELSE Reason END,
    UpdatedAt = ?
WHERE AiConfusionId = ? OR ConfusionUuid = ?;`

// MarkConfusionCompleted marks a confusion episode as resolved.
func (s *AiConfusionStore) MarkConfusionCompleted(identifier string, notes string) error {
	cleanId, idErr := s.validateIdentifier(identifier)
	if idErr != nil {
		return idErr
	}

	targetId := s.resolveNumericId(cleanId)
	nowStr := time.Now().UTC().Format(time.RFC3339)

	return s.executeCompleteUpdate(targetId, cleanId, notes, nowStr)
}

func (s *AiConfusionStore) validateIdentifier(identifier string) (string, error) {
	cleanId := strings.TrimSpace(identifier)
	hasEmptyId := len(cleanId) == 0
	if hasEmptyId {
		return "", apperror.NewValidationError("identifier is empty")
	}

	return cleanId, nil
}

func (s *AiConfusionStore) resolveNumericId(cleanId string) int64 {
	idNum, numErr := strconv.ParseInt(cleanId, 10, 64)
	isNumeric := numErr == nil
	if isNumeric {
		return idNum
	}

	return int64(-1)
}

func (s *AiConfusionStore) executeCompleteUpdate(targetId int64, cleanId, notes, nowStr string) error {
	res, err := s.db.Exec(sqlCompleteConfusion, notes, notes, nowStr, targetId, cleanId)
	if err != nil {
		return apperror.WrapSimple(err, "mark confusion completed")
	}

	affected, _ := res.RowsAffected()
	hasZeroAffected := affected == 0
	if hasZeroAffected {
		return apperror.NewValidationError("confusion record not found: " + cleanId)
	}

	return nil
}
