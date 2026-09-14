package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	sqlInsertIPSnapshot = `INSERT INTO IPSnapshot (
		InterfaceName, IP, Netmask, Gateway, DNS, IsDHCP, Timestamp, Notes, Comments
	) VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlSelectLatestIPSnapshot = `SELECT
		IPSnapshotId, InterfaceName, IP, Netmask, Gateway, DNS, IsDHCP, Timestamp,
		COALESCE(Notes, ''), COALESCE(Comments, '')
	FROM IPSnapshot
	WHERE InterfaceName = ?
	ORDER BY IPSnapshotId DESC
	LIMIT 1;`

	sqlSelectListIPSnapshots = `SELECT
		IPSnapshotId, InterfaceName, IP, Netmask, Gateway, DNS, IsDHCP, Timestamp,
		COALESCE(Notes, ''), COALESCE(Comments, '')
	FROM IPSnapshot
	WHERE InterfaceName = ?
	ORDER BY IPSnapshotId DESC
	LIMIT ?;`
)

type scanFn func(dest ...any) error

func scanFields(scan scanFn, rec *IPSnapshotRecord, isDHCP *int) error {
	return scan(
		&rec.IPSnapshotId, &rec.InterfaceName, &rec.IP, &rec.Netmask,
		&rec.Gateway, &rec.DNS, isDHCP, &rec.Timestamp,
		&rec.Notes, &rec.Comments,
	)
}

// InsertIPSnapshot saves a snapshot into SQLite with automatic JSON backup.
func (db *DB) InsertIPSnapshot(record *IPSnapshotRecord) *apperror.AppError {
	if record == nil {
		return apperror.NewSimple("store.InsertIPSnapshot", "E_NIL_RECORD")
	}

	_ = SaveIPSnapshotJSON(record)
	if db == nil || db.conn == nil {
		return nil
	}

	return db.execInsertIPSnapshot(record)
}

func (db *DB) execInsertIPSnapshot(rec *IPSnapshotRecord) *apperror.AppError {
	if err := db.EnsureIPSnapshotTable(); err != nil {
		return err
	}

	res, execErr := db.conn.Exec(sqlInsertIPSnapshot,
		rec.InterfaceName, rec.IP, rec.Netmask, rec.Gateway, rec.DNS,
		boolToInt(rec.IsDHCP), rec.Timestamp, rec.Notes, rec.Comments,
	)
	if execErr != nil {
		return apperror.WrapSimple(execErr, "store.execInsertIPSnapshot")
	}

	return assignInsertID(rec, res)
}

func assignInsertID(rec *IPSnapshotRecord, res sql.Result) *apperror.AppError {
	id, err := res.LastInsertId()
	if err != nil {
		return apperror.WrapSimple(err, "store.assignInsertID")
	}

	rec.IPSnapshotId = id

	return nil
}

// GetLatestIPSnapshot fetches the latest snapshot for an interface or falls back to JSON.
func (db *DB) GetLatestIPSnapshot(ifaceName string) (*IPSnapshotRecord, *apperror.AppError) {
	if db == nil || db.conn == nil {
		return LoadLatestIPSnapshotJSON()
	}

	if err := db.EnsureIPSnapshotTable(); err != nil {
		return LoadLatestIPSnapshotJSON()
	}

	return db.queryLatestSnapshot(ifaceName)
}

func (db *DB) queryLatestSnapshot(ifaceName string) (*IPSnapshotRecord, *apperror.AppError) {
	row := db.conn.QueryRow(sqlSelectLatestIPSnapshot, ifaceName)
	rec, err := scanIPSnapshotRow(row)
	if err != nil || rec == nil {
		return LoadLatestIPSnapshotJSON()
	}

	return rec, nil
}

func scanIPSnapshotRow(row *sql.Row) (*IPSnapshotRecord, *apperror.AppError) {
	var rec IPSnapshotRecord
	var isDHCP int
	if err := scanFields(row.Scan, &rec, &isDHCP); err != nil {
		return handleRowScanError(err)
	}

	rec.IsDHCP = isDHCP == 1

	return &rec, nil
}

func handleRowScanError(err error) (*IPSnapshotRecord, *apperror.AppError) {
	if err == sql.ErrNoRows {
		return nil, nil
	}

	return nil, apperror.WrapSimple(err, "store.handleRowScanError")
}

// ListIPSnapshots returns historical snapshots for an interface.
func (db *DB) ListIPSnapshots(
	ifaceName string,
	limit int,
) ([]IPSnapshotRecord, *apperror.AppError) {
	rows, err := db.querySnapshotRows(ifaceName, limit)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	return scanIPSnapshotRows(rows)
}

func (db *DB) querySnapshotRows(ifaceName string, limit int) (*sql.Rows, *apperror.AppError) {
	if err := db.checkSnapshotTable(); err != nil {
		return nil, err
	}

	rows, err := db.conn.Query(sqlSelectListIPSnapshots, ifaceName, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.querySnapshotRows.Query")
	}

	return rows, nil
}

func (db *DB) checkSnapshotTable() *apperror.AppError {
	if db == nil || db.conn == nil {
		return apperror.NewSimple("store.checkSnapshotTable", "E_NIL_CONN")
	}

	return db.EnsureIPSnapshotTable()
}

func scanIPSnapshotRows(rows *sql.Rows) ([]IPSnapshotRecord, *apperror.AppError) {
	records := make([]IPSnapshotRecord, 0)
	for rows.Next() {
		rec, err := scanOneIPSnapshotRow(rows)
		if err != nil {
			return nil, err
		}
		records = append(records, *rec)
	}

	return records, nil
}

func scanOneIPSnapshotRow(rows *sql.Rows) (*IPSnapshotRecord, *apperror.AppError) {
	var rec IPSnapshotRecord
	var isDHCP int
	if err := scanFields(rows.Scan, &rec, &isDHCP); err != nil {
		return nil, apperror.WrapSimple(err, "store.scanOneIPSnapshotRow")
	}

	rec.IsDHCP = isDHCP == 1

	return &rec, nil
}
