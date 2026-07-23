package store

import (
	"database/sql"
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// InsertTransactionAction appends one typed action row and returns its id.
// Seq must be strictly the previous max + 1; callers fetch it via
// NextActionSeq under the same lock to keep the invariant tight.
func (db *DB) InsertTransactionAction(r model.TransactionActionRecord) (int64, error) {
	res, err := db.conn.Exec(constants.SQLInsertTransactionAction,
		r.TransactionID, r.Seq, r.Kind, r.ForwardJSON, r.ReverseJSON,
		r.BackupRef, nowUnix())
	if err != nil {
		return 0, fmt.Errorf(constants.ErrActionInsert, err)
	}

	return res.LastInsertId()
}

// NextActionSeq returns the next 1-based Seq for the given transaction.
func (db *DB) NextActionSeq(txnID int64) (int64, error) {
	row := db.conn.QueryRow(constants.SQLSelectMaxActionSeq, txnID)
	var maxSeq int64
	if err := row.Scan(&maxSeq); err != nil {
		return 0, fmt.Errorf(constants.ErrActionList, err)
	}

	return maxSeq + 1, nil
}

// ListTransactionActions returns every action row for one transaction in
// forward (Seq ASC) order. Use ListTransactionActionsReverse for revert.
func (db *DB) ListTransactionActions(txnID int64) ([]model.TransactionActionRecord, error) {
	return db.queryActions(constants.SQLSelectTransactionActions, txnID)
}

// ListTransactionActionsReverse returns rows newest-first for revert.
func (db *DB) ListTransactionActionsReverse(txnID int64) ([]model.TransactionActionRecord, error) {
	return db.queryActions(constants.SQLSelectTransactionActionsReverse, txnID)
}

// queryActions is the shared scan path for both list orderings.
func (db *DB) queryActions(query string, txnID int64) ([]model.TransactionActionRecord, error) {
	rows, err := db.conn.Query(query, txnID)
	if err != nil {
		return nil, fmt.Errorf(constants.ErrActionList, err)
	}
	defer rows.Close()

	return scanActionRows(rows)
}

// MarkTransactionActionReverted stamps RevertedAt for one action row.
func (db *DB) MarkTransactionActionReverted(actionID int64) error {
	_, err := db.conn.Exec(constants.SQLMarkTransactionActionReverted,
		nowUnix(), actionID)
	if err != nil {
		return fmt.Errorf(constants.ErrActionMarkReverted, err)
	}

	return nil
}

// scanActionRows reads every action row into a typed slice.
func scanActionRows(rows *sql.Rows) ([]model.TransactionActionRecord, error) {
	var out []model.TransactionActionRecord
	for rows.Next() {
		r, err := scanOneAction(rows)
		if err != nil {
			return nil, err
		}
		out = append(out, r)
	}

	return out, nil
}

// scanOneAction reads a single row into a TransactionActionRecord.
func scanOneAction(rows *sql.Rows) (model.TransactionActionRecord, error) {
	var r model.TransactionActionRecord
	var reverted sql.NullInt64
	err := rows.Scan(&r.ID, &r.TransactionID, &r.Seq, &r.Kind,
		&r.ForwardJSON, &r.ReverseJSON, &r.BackupRef,
		&r.AppliedAt, &reverted)
	if err != nil {
		return r, fmt.Errorf(constants.ErrActionList, err)
	}
	r.RevertedAt = nullInt(reverted)

	return r, nil
}
