package store

import (
	"database/sql"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const commitSelectColumns = `SELECT CommitHistoryId, CommitSha, ShortSha, RepoName, RepoSlug, RepoRoot,
	BranchName, AuthorName, AuthorEmail, CommitMessage, CommittedAtUnix,
	AiModelName, AiToolName, AgentId, HeatmapSessionId, HeatmapActionId,
	TaskSlug, SpecPath, PlansCount, SubtasksCount,
	IsAiCommitted, IsCurrentMachineAi, Notes,
	FilesChanged, LinesAdded, LinesDeleted, HasPushed, CreatedAtUnix
FROM CommitHistoryRecord`

// RecordCommit inserts a commit record and all associated file touches atomically.
func (db *CommitHistorySplitDB) RecordCommit(record CommitHistoryRecord, touches []CommitFileTouch) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	tx, err := db.conn.Begin()
	if err != nil {
		return apperror.WrapSimple(err, "store.RecordCommit.begin")
	}
	defer func() { _ = tx.Rollback() }()

	if errInsert := insertCommitRecordTx(tx, record); errInsert != nil {
		return errInsert
	}
	if errTouches := insertTouchesTx(tx, record.CommitSha, touches); errTouches != nil {
		return errTouches
	}

	return commitTxOrError(tx)
}

func commitTxOrError(tx *sql.Tx) error {
	if err := tx.Commit(); err != nil {
		return apperror.WrapSimple(err, "store.RecordCommit.commit")
	}

	return nil
}

func insertCommitRecordTx(tx *sql.Tx, rec CommitHistoryRecord) error {
	query := `INSERT OR REPLACE INTO CommitHistoryRecord (
		CommitSha, ShortSha, RepoName, RepoSlug, RepoRoot, BranchName,
		AuthorName, AuthorEmail, CommitMessage, CommittedAtUnix,
		AiModelName, AiToolName, AgentId, HeatmapSessionId, HeatmapActionId,
		TaskSlug, SpecPath, PlansCount, SubtasksCount,
		IsAiCommitted, IsCurrentMachineAi, Notes,
		FilesChanged, LinesAdded, LinesDeleted, HasPushed, CreatedAtUnix
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	nowUnix := resolveNowUnix(rec.CreatedAtUnix)
	_, err := tx.Exec(query,
		rec.CommitSha, rec.ShortSha, rec.RepoName, rec.RepoSlug, rec.RepoRoot, rec.BranchName,
		rec.AuthorName, rec.AuthorEmail, rec.CommitMessage, rec.CommittedAtUnix,
		rec.AiModelName, rec.AiToolName, rec.AgentId, rec.HeatmapSessionId, rec.HeatmapActionId,
		rec.TaskSlug, rec.SpecPath, rec.PlansCount, rec.SubtasksCount,
		boolToInt(rec.IsAiCommitted), boolToInt(rec.IsCurrentMachineAi), rec.Notes,
		rec.FilesChanged, rec.LinesAdded, rec.LinesDeleted, boolToInt(rec.HasPushed), nowUnix,
	)
	if err != nil {
		return apperror.WrapSimple(err, "store.insertCommitRecordTx")
	}

	return nil
}

func insertTouchesTx(tx *sql.Tx, sha string, touches []CommitFileTouch) error {
	query := `INSERT INTO CommitFileTouch (
		CommitSha, FilePath, ChangeType, LinesAdded, LinesDeleted, CreatedAtUnix
	) VALUES (?, ?, ?, ?, ?, ?);`

	for _, t := range touches {
		nowUnix := resolveNowUnix(t.CreatedAtUnix)
		_, err := tx.Exec(query, sha, t.FilePath, t.ChangeType, t.LinesAdded, t.LinesDeleted, nowUnix)
		if err != nil {
			return apperror.WrapSimple(err, "store.insertTouchesTx")
		}
	}

	return nil
}

func resolveNowUnix(ts int64) int64 {
	if ts > 0 {
		return ts
	}

	return time.Now().Unix()
}

// GetBySha retrieves a commit record by full SHA or ShortSha.
func (db *CommitHistorySplitDB) GetBySha(sha string) (*CommitHistoryRecord, bool, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	cleanSha := strings.TrimSpace(sha)
	query := commitSelectColumns + " WHERE CommitSha = ? OR ShortSha = ? LIMIT 1;"

	row := db.conn.QueryRow(query, cleanSha, cleanSha)
	rec, hasFound, err := scanCommitRecordRow(row)
	if err != nil {
		return nil, false, apperror.WrapSimple(err, "store.GetBySha")
	}

	return rec, hasFound, nil
}

// GetCommitBySha is a compatibility alias for GetBySha.
func (db *CommitHistorySplitDB) GetCommitBySha(sha string) (*CommitHistoryRecord, bool, error) {
	return db.GetBySha(sha)
}

// SearchCommits searches commit history matching SHA prefix, suffix, substring, or metadata.
func (db *CommitHistorySplitDB) SearchCommits(queryStr string, limit int) ([]CommitHistoryRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	effLimit := normalizeLimit(limit)
	pattern := "%" + strings.TrimSpace(queryStr) + "%"
	sqlQuery := commitSelectColumns + ` WHERE CommitSha LIKE ? OR ShortSha LIKE ? 
		OR CommitMessage LIKE ? OR TaskSlug LIKE ? OR AgentId LIKE ?
		ORDER BY CommittedAtUnix DESC LIMIT ?;`

	rows, err := db.conn.Query(sqlQuery, pattern, pattern, pattern, pattern, pattern, effLimit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.SearchCommits.query")
	}
	defer rows.Close()

	return scanCommitRecords(rows)
}

// ListRecentCommits retrieves the latest commit records ordered by commit time descending.
func (db *CommitHistorySplitDB) ListRecentCommits(limit int) ([]CommitHistoryRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	effLimit := normalizeLimit(limit)
	query := commitSelectColumns + " ORDER BY CommittedAtUnix DESC LIMIT ?;"

	rows, err := db.conn.Query(query, effLimit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListRecentCommits.query")
	}
	defer rows.Close()

	return scanCommitRecords(rows)
}

// ListFileTouches retrieves all touched files associated with a given commit SHA.
func (db *CommitHistorySplitDB) ListFileTouches(sha string) ([]CommitFileTouch, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	query := `SELECT CommitFileTouchId, CommitSha, FilePath, ChangeType, LinesAdded, LinesDeleted, CreatedAtUnix
		FROM CommitFileTouch WHERE CommitSha = ? ORDER BY FilePath ASC;`

	rows, err := db.conn.Query(query, strings.TrimSpace(sha))
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListFileTouches.query")
	}
	defer rows.Close()

	return scanFileTouches(rows)
}

// ListCommits retrieves commit records matching filter options.
func (db *CommitHistorySplitDB) ListCommits(filter CommitHistoryFilter) ([]CommitHistoryRecord, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	query, args := buildFilterQuery(filter)
	rows, err := db.conn.Query(query, args...)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.ListCommits.query")
	}
	defer rows.Close()

	return scanCommitRecords(rows)
}

func buildFilterQuery(filter CommitHistoryFilter) (string, []any) {
	whereClauses, args := collectFilterClauses(filter)
	query := commitSelectColumns
	if len(whereClauses) > 0 {
		query += " WHERE " + strings.Join(whereClauses, " AND ")
	}

	query += " ORDER BY CommittedAtUnix DESC LIMIT ? OFFSET ?;"
	args = append(args, normalizeLimit(filter.Limit), normalizeOffset(filter.Offset))

	return query, args
}

func collectFilterClauses(f CommitHistoryFilter) ([]string, []any) {
	var clauses []string
	var args []any
	if f.RepoSlug != "" {
		clauses = append(clauses, "RepoSlug = ?")
		args = append(args, f.RepoSlug)
	}
	if f.TaskSlug != "" {
		clauses = append(clauses, "TaskSlug = ?")
		args = append(args, f.TaskSlug)
	}
	if f.AgentId != "" {
		clauses = append(clauses, "AgentId = ?")
		args = append(args, f.AgentId)
	}
	if f.IsAiOnly {
		clauses = append(clauses, "IsAiCommitted = 1")
	}

	return clauses, args
}

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return 50
	}
	if limit > 500 {
		return 500
	}

	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}

	return offset
}

func scanCommitRecordRow(s interface{ Scan(...any) error }) (*CommitHistoryRecord, bool, error) {
	var r CommitHistoryRecord
	var isAi, isLocal, hasPush int

	err := s.Scan(
		&r.CommitHistoryId, &r.CommitSha, &r.ShortSha, &r.RepoName, &r.RepoSlug, &r.RepoRoot,
		&r.BranchName, &r.AuthorName, &r.AuthorEmail, &r.CommitMessage, &r.CommittedAtUnix,
		&r.AiModelName, &r.AiToolName, &r.AgentId, &r.HeatmapSessionId, &r.HeatmapActionId,
		&r.TaskSlug, &r.SpecPath, &r.PlansCount, &r.SubtasksCount,
		&isAi, &isLocal, &r.Notes,
		&r.FilesChanged, &r.LinesAdded, &r.LinesDeleted, &hasPush, &r.CreatedAtUnix,
	)
	if err == sql.ErrNoRows {
		return nil, false, nil
	}
	if err != nil {
		return nil, false, err
	}

	r.IsAiCommitted = isAi == 1
	r.IsCurrentMachineAi = isLocal == 1
	r.HasPushed = hasPush == 1

	return &r, true, nil
}

func scanCommitRecords(rows *sql.Rows) ([]CommitHistoryRecord, error) {
	var list []CommitHistoryRecord
	for rows.Next() {
		rec, hasFound, err := scanCommitRecordRow(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "store.scanCommitRecords")
		}
		if hasFound && rec != nil {
			list = append(list, *rec)
		}
	}

	return list, nil
}

func scanFileTouches(rows *sql.Rows) ([]CommitFileTouch, error) {
	var list []CommitFileTouch
	for rows.Next() {
		var t CommitFileTouch
		err := rows.Scan(
			&t.CommitFileTouchId, &t.CommitSha, &t.FilePath, &t.ChangeType,
			&t.LinesAdded, &t.LinesDeleted, &t.CreatedAtUnix,
		)
		if err != nil {
			return nil, apperror.WrapSimple(err, "store.scanFileTouches")
		}
		list = append(list, t)
	}

	return list, nil
}
