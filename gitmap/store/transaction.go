package store

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// ErrTransactionNotFound is returned by FindTransactionByID when no row
// matches the supplied id. Callers check via errors.Is so the existing
// zero-swallow error policy still applies for unexpected SQL failures.
var ErrTransactionNotFound = errors.New("transaction not found")

// InsertTransaction creates a pending transaction row and returns its id.
func (db *DB) InsertTransaction(r model.TransactionRecord) (int64, error) {
	res, err := db.conn.Exec(constants.SQLInsertTransaction,
		r.Kind, constants.TxnStatusPending, r.Argv, r.Cwd,
		nowUnix(), r.ReverseSummary, r.RepoSlug, r.GitSha)
	if err != nil {
		return 0, fmt.Errorf("transaction insert: %w", err)
	}

	return res.LastInsertId()
}

// MarkTransactionCommitted flips a pending row to committed.
func (db *DB) MarkTransactionCommitted(id int64) error {
	return db.updateTxnStatus(id, constants.TxnStatusCommitted, true, false)
}

// MarkTransactionAborted flips a pending row to aborted (commit failed).
func (db *DB) MarkTransactionAborted(id int64) error {
	return db.updateTxnStatus(id, constants.TxnStatusAborted, false, false)
}

// MarkTransactionReverted flips a committed row to reverted.
func (db *DB) MarkTransactionReverted(id int64) error {
	return db.updateTxnStatus(id, constants.TxnStatusReverted, false, true)
}

// updateTxnStatus is the shared status-write path for the three flips above.
func (db *DB) updateTxnStatus(id int64, status string, setCommit, setRevert bool) error {
	committed := nullableUnix(setCommit)
	reverted := nullableUnix(setRevert)
	_, err := db.conn.Exec(constants.SQLUpdateTransactionStatus,
		status, committed, reverted, id)
	if err != nil {
		return fmt.Errorf("transaction status update: %w", err)
	}

	return nil
}

// InsertTransactionFile records one snapshotted file under a transaction.
func (db *DB) InsertTransactionFile(r model.TransactionFileRecord) error {
	_, err := db.conn.Exec(constants.SQLInsertTransactionFile,
		r.TransactionID, r.RelPath, r.AbsPath, r.BackupPath,
		r.ByteSize, r.Sha256, r.Action)
	if err != nil {
		return fmt.Errorf("transaction file insert: %w", err)
	}

	return nil
}

// FindTransactionByID loads one row or returns ErrTransactionNotFound.
func (db *DB) FindTransactionByID(id int64) (model.TransactionRecord, error) {
	row := db.conn.QueryRow(constants.SQLSelectTransactionByID, id)

	return scanTransactionRow(row)
}

// ListTransactions returns the most recent rows, newest first, capped at limit.
func (db *DB) ListTransactions(limit int) ([]model.TransactionRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectTransactionsRecent, limit)
	if err != nil {
		return nil, fmt.Errorf("transaction list: %w", err)
	}
	defer rows.Close()

	return scanTransactionRows(rows)
}

// ListTransactionFiles returns every snapshotted file for one transaction.
func (db *DB) ListTransactionFiles(txnID int64) ([]model.TransactionFileRecord, error) {
	rows, err := db.conn.Query(constants.SQLSelectTransactionFiles, txnID)
	if err != nil {
		return nil, fmt.Errorf("transaction files: %w", err)
	}
	defer rows.Close()

	return scanTransactionFileRows(rows)
}

// LastCommittedTransactionID returns the newest committed id, or 0 if none.
func (db *DB) LastCommittedTransactionID() (int64, error) {
	row := db.conn.QueryRow(constants.SQLSelectLastCommittedTransaction)
	var id int64
	err := row.Scan(&id)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, nil
	}
	if err != nil {
		return 0, fmt.Errorf("transaction last-committed: %w", err)
	}

	return id, nil
}

// PruneOldestTransactions deletes everything beyond the cap (newest cap rows
// kept). Returns the deleted ids so the caller can clean their backup dirs.
func (db *DB) PruneOldestTransactions(cap int) ([]int64, error) {
	rows, err := db.conn.Query(constants.SQLSelectExcessTransactionIDs, cap)
	if err != nil {
		return nil, fmt.Errorf("transaction prune list: %w", err)
	}
	ids, err := collectInt64Column(rows)
	if err != nil {
		return nil, err
	}

	return ids, deleteTransactionRows(db, ids)
}

// deleteTransactionRows removes each id; the FK cascade drops file rows.
func deleteTransactionRows(db *DB, ids []int64) error {
	for _, id := range ids {
		if _, err := db.conn.Exec(constants.SQLDeleteTransaction, id); err != nil {
			return fmt.Errorf("transaction prune delete: %w", err)
		}
	}

	return nil
}

// scanTransactionRow reads one *sql.Row into a TransactionRecord.
func scanTransactionRow(row *sql.Row) (model.TransactionRecord, error) {
	var r model.TransactionRecord
	var committed, reverted sql.NullInt64
	err := row.Scan(&r.ID, &r.Kind, &r.Status, &r.Argv, &r.Cwd,
		&r.CreatedAt, &committed, &reverted, &r.ReverseSummary,
		&r.RepoSlug, &r.GitSha)
	if errors.Is(err, sql.ErrNoRows) {
		return r, ErrTransactionNotFound
	}
	if err != nil {
		return r, fmt.Errorf("transaction scan: %w", err)
	}
	r.CommittedAt = nullInt(committed)
	r.RevertedAt = nullInt(reverted)

	return r, nil
}

// scanTransactionRows reads every *sql.Rows entry into a slice.
func scanTransactionRows(rows *sql.Rows) ([]model.TransactionRecord, error) {
	var out []model.TransactionRecord
	for rows.Next() {
		r, err := scanOneTxnFromRows(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, nil
}

// scanOneTxnFromRows is split out to keep scanTransactionRows under the line cap.
func scanOneTxnFromRows(rows *sql.Rows) (model.TransactionRecord, error) {
	var r model.TransactionRecord
	var committed, reverted sql.NullInt64
	err := rows.Scan(&r.ID, &r.Kind, &r.Status, &r.Argv, &r.Cwd,
		&r.CreatedAt, &committed, &reverted, &r.ReverseSummary,
		&r.RepoSlug, &r.GitSha)
	if err != nil {
		return r, fmt.Errorf("transaction scan: %w", err)
	}
	r.CommittedAt = nullInt(committed)
	r.RevertedAt = nullInt(reverted)

	return r, nil
}

// scanTransactionFileRows reads every TransactionFile row.
func scanTransactionFileRows(rows *sql.Rows) ([]model.TransactionFileRecord, error) {
	var out []model.TransactionFileRecord
	for rows.Next() {
		var r model.TransactionFileRecord
		err := rows.Scan(&r.ID, &r.TransactionID, &r.RelPath, &r.AbsPath,
			&r.BackupPath, &r.ByteSize, &r.Sha256, &r.Action)
		if err != nil {
			return nil, fmt.Errorf("transaction file scan: %w", err)
		}
		out = append(out, r)
	}

	return out, nil
}

// collectInt64Column reads a single-column INTEGER result set into a slice.
func collectInt64Column(rows *sql.Rows) ([]int64, error) {
	defer rows.Close()
	var ids []int64
	for rows.Next() {
		var id int64
		if err := rows.Scan(&id); err != nil {
			return nil, fmt.Errorf("transaction id scan: %w", err)
		}
		ids = append(ids, id)
	}

	return ids, nil
}

// nowUnix returns the current unix-seconds timestamp.
func nowUnix() int64 { return time.Now().Unix() }

// nullableUnix returns a NullInt64 set when the boolean is true.
func nullableUnix(set bool) sql.NullInt64 {
	if set {
		return sql.NullInt64{Int64: nowUnix(), Valid: true}
	}

	return sql.NullInt64{}
}

// nullInt unwraps a sql.NullInt64 into a plain int64 (0 when null).
func nullInt(n sql.NullInt64) int64 {
	if n.Valid {
		return n.Int64
	}

	return 0
}
