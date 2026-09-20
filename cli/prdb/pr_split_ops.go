package prdb

import (
	"database/sql"
	"errors"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/dbengine"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func nullableString(val string) any {
	if val == "" {
		return nil
	}

	return val
}

func nullableInt64(val int64) any {
	if val <= 0 {
		return nil
	}

	return val
}

func boolToInt(b bool) int {
	if b {
		return 1
	}

	return 0
}

func normalizePullRequestRecord(rec *PullRequestRecord) {
	now := time.Now().UTC().Unix()
	if rec.CreatedAt == 0 {
		rec.CreatedAt = now
	}
	if rec.UpdatedAt == 0 {
		rec.UpdatedAt = rec.CreatedAt
	}
	if rec.Status == "" {
		rec.Status = "open"
	}
	if rec.TargetBranch == "" {
		rec.TargetBranch = "main"
	}
}

// CreatePullRequest inserts a new pull request record and returns its generated ID.
func (p *PrSplitDb) CreatePullRequest(rec PullRequestRecord) result.Result[int64] {
	normalizePullRequestRecord(&rec)
	res, err := p.conn.Exec(
		constants.SQLInsertPullRequest,
		rec.PrNumber, rec.Title, rec.Description, rec.SourceBranch, rec.TargetBranch,
		rec.Status, nullableString(rec.MergeCommitSha), rec.CreatedAt,
		nullableInt64(rec.MergedAt), nullableInt64(rec.ClosedAt),
		nullableString(rec.Notes), nullableString(rec.Comments), rec.UpdatedAt,
	)
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "prdb.createPullRequest"))
	}

	id, _ := res.LastInsertId()

	return result.Ok(id)
}

func scanPullRequestRecord(row *sql.Row) result.Result[*PullRequestRecord] {
	var (
		rec                      PullRequestRecord
		rawMergeSha, rawNotes    any
		rawComments              any
		rawMergedAt, rawClosedAt any
	)
	err := row.Scan(
		&rec.PullRequestId, &rec.PrNumber, &rec.Title, &rec.Description,
		&rec.SourceBranch, &rec.TargetBranch, &rec.Status, &rawMergeSha,
		&rec.CreatedAt, &rawMergedAt, &rawClosedAt, &rawNotes, &rawComments, &rec.UpdatedAt,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return result.Ok[*PullRequestRecord](nil)
		}

		return result.Fail[*PullRequestRecord](apperror.WrapSimple(err, "prdb.scanPullRequestRecord"))
	}
	rec.MergeCommitSha = dbengine.ScanString(rawMergeSha)
	rec.MergedAt = dbengine.ScanInt64(rawMergedAt)
	rec.ClosedAt = dbengine.ScanInt64(rawClosedAt)
	rec.Notes = dbengine.ScanString(rawNotes)
	rec.Comments = dbengine.ScanString(rawComments)

	return result.Ok(&rec)
}

// GetPullRequestByNumber retrieves a pull request record by its PR number.
func (p *PrSplitDb) GetPullRequestByNumber(prNum int) result.Result[*PullRequestRecord] {
	row := p.conn.QueryRow(constants.SQLSelectPullRequestByNumber, prNum)

	return scanPullRequestRecord(row)
}

// UpdatePullRequestStatus updates status and merge SHA for a PR number.
func (p *PrSplitDb) UpdatePullRequestStatus(prNum int, status, mergeSha string) result.Result[bool] {
	now := time.Now().UTC().Unix()
	res, err := p.conn.Exec(
		constants.SQLUpdatePullRequestStatus,
		status, mergeSha, mergeSha, status, now, status, now, now, prNum,
	)
	if err != nil {
		return result.Fail[bool](apperror.WrapSimple(err, "prdb.updatePullRequestStatus"))
	}

	rowsAffected, _ := res.RowsAffected()

	return result.Ok(rowsAffected > 0)
}

// AddPrRelease inserts a new PR release record and returns its generated ID.
func (p *PrSplitDb) AddPrRelease(rec PrReleaseRecord) result.Result[int64] {
	if rec.CreatedAt == 0 {
		rec.CreatedAt = time.Now().UTC().Unix()
	}
	res, err := p.conn.Exec(
		constants.SQLInsertPrRelease,
		rec.PullRequestId, rec.ReleaseTag, rec.CommitSha,
		nullableString(rec.Notes), nullableString(rec.Comments), rec.CreatedAt,
	)
	if err != nil {
		return result.Fail[int64](apperror.WrapSimple(err, "prdb.addPrRelease"))
	}

	id, _ := res.LastInsertId()

	return result.Ok(id)
}

func normalizePrBranchRecord(rec *PrBranchRecord) {
	now := time.Now().UTC().Unix()
	if rec.CreatedAt == 0 {
		rec.CreatedAt = now
	}
	if rec.UpdatedAt == 0 {
		rec.UpdatedAt = rec.CreatedAt
	}
	if rec.BranchType == "" {
		rec.BranchType = "feature"
	}
}

// UpsertPrBranch creates or updates a tracked PR branch record.
func (p *PrSplitDb) UpsertPrBranch(rec PrBranchRecord) result.Result[bool] {
	normalizePrBranchRecord(&rec)
	_, err := p.conn.Exec(
		constants.SQLUpsertPrBranch,
		rec.BranchName, rec.BranchType, boolToInt(rec.IsMerged), boolToInt(rec.IsDeleted),
		rec.CreatedAt, nullableInt64(rec.MergedAt), nullableInt64(rec.DeletedAt),
		nullableString(rec.Notes), nullableString(rec.Comments), rec.UpdatedAt,
	)
	if err != nil {
		return result.Fail[bool](apperror.WrapSimple(err, "prdb.upsertPrBranch"))
	}

	return result.Ok(true)
}

func populatePrBranchFields(
	b *PrBranchRecord,
	rawMerged, rawDeleted, rawMergedAt, rawDeletedAt, rawNotes, rawComments any,
) {
	b.IsMerged = dbengine.ScanBool(rawMerged)
	b.IsDeleted = dbengine.ScanBool(rawDeleted)
	b.MergedAt = dbengine.ScanInt64(rawMergedAt)
	b.DeletedAt = dbengine.ScanInt64(rawDeletedAt)
	b.Notes = dbengine.ScanString(rawNotes)
	b.Comments = dbengine.ScanString(rawComments)
}

func scanSinglePrBranch(rows *sql.Rows) result.Result[PrBranchRecord] {
	var (
		b                         PrBranchRecord
		rawMerged, rawDeleted     any
		rawMergedAt, rawDeletedAt any
		rawNotes, rawComments     any
	)
	err := rows.Scan(
		&b.PrBranchId, &b.BranchName, &b.BranchType, &rawMerged, &rawDeleted,
		&b.CreatedAt, &rawMergedAt, &rawDeletedAt, &rawNotes, &rawComments, &b.UpdatedAt,
	)
	if err != nil {
		return result.Fail[PrBranchRecord](apperror.WrapSimple(err, "prdb.scanSinglePrBranch"))
	}
	populatePrBranchFields(&b, rawMerged, rawDeleted, rawMergedAt, rawDeletedAt, rawNotes, rawComments)

	return result.Ok(b)
}

func scanPrBranchRows(rows *sql.Rows) result.Result[[]PrBranchRecord] {
	var list []PrBranchRecord
	for rows.Next() {
		itemRes := scanSinglePrBranch(rows)
		if itemRes.IsFailure() {
			return result.Fail[[]PrBranchRecord](itemRes.Err)
		}
		list = append(list, itemRes.Value)
	}
	if err := rows.Err(); err != nil {
		return result.Fail[[]PrBranchRecord](apperror.WrapSimple(err, "prdb.scanPrBranchRows.rowsErr"))
	}

	return result.Ok(list)
}

// ListActivePrBranches retrieves all active (unmerged, not deleted) PR branches.
func (p *PrSplitDb) ListActivePrBranches() result.Result[[]PrBranchRecord] {
	rows, err := p.conn.Query(constants.SQLSelectListActivePrBranches)
	if err != nil {
		return result.Fail[[]PrBranchRecord](apperror.WrapSimple(err, "prdb.listActivePrBranches"))
	}
	defer rows.Close()

	return scanPrBranchRows(rows)
}

// ListMergedPrBranches retrieves all merged but not yet deleted PR branches.
func (p *PrSplitDb) ListMergedPrBranches() result.Result[[]PrBranchRecord] {
	rows, err := p.conn.Query(constants.SQLSelectListMergedPrBranches)
	if err != nil {
		return result.Fail[[]PrBranchRecord](apperror.WrapSimple(err, "prdb.listMergedPrBranches"))
	}
	defer rows.Close()

	return scanPrBranchRows(rows)
}

// MarkPrBranchDeleted marks a branch record as deleted with timestamp.
func (p *PrSplitDb) MarkPrBranchDeleted(branchName string) result.Result[bool] {
	now := time.Now().UTC().Unix()
	res, err := p.conn.Exec(constants.SQLMarkPrBranchDeleted, now, now, branchName)
	if err != nil {
		return result.Fail[bool](apperror.WrapSimple(err, "prdb.markPrBranchDeleted"))
	}

	rowsAffected, _ := res.RowsAffected()

	return result.Ok(rowsAffected > 0)
}
