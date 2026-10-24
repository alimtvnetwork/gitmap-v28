package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
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
func (db *DB) GetLatestIPSnapshot(ifaceName string) IPSnapshotRecordWrap {
	if db == nil || db.conn == nil {
		return LoadLatestIPSnapshotJSON()
	}

	if err := db.EnsureIPSnapshotTable(); err != nil {
		return LoadLatestIPSnapshotJSON()
	}

	return db.queryLatestSnapshot(ifaceName)
}

func (db *DB) queryLatestSnapshot(ifaceName string) IPSnapshotRecordWrap {
	row := db.conn.QueryRow(sqlSelectLatestIPSnapshot, ifaceName)
	scanRes := scanIPSnapshotRow(row)
	if scanRes.IsFailure() || scanRes.Value == nil {
		return LoadLatestIPSnapshotJSON()
	}

	return scanRes
}

func scanIPSnapshotRow(row *sql.Row) IPSnapshotRecordWrap {
	var rec IPSnapshotRecord
	var isDHCP int
	if err := scanFields(row.Scan, &rec, &isDHCP); err != nil {
		return handleRowScanError(err)
	}

	rec.IsDHCP = isDHCP == 1

	return result.Ok(&rec)
}

func handleRowScanError(err error) IPSnapshotRecordWrap {
	if err == sql.ErrNoRows {
		return result.Ok((*IPSnapshotRecord)(nil))
	}

	return result.Fail[*IPSnapshotRecord](apperror.WrapSimple(err, "store.handleRowScanError"))
}

// ListIPSnapshots returns historical snapshots for an interface.
func (db *DB) ListIPSnapshots(
	ifaceName string,
	limit int,
) IPSnapshotRecordSliceWrap {
	rowsRes := db.querySnapshotRows(ifaceName, limit)
	if rowsRes.IsFailure() {
		return result.FailSlice[IPSnapshotRecord](rowsRes.AppError())
	}
	rows := rowsRes.Value
	defer rows.Close()

	return scanIPSnapshotRows(rows)
}

func (db *DB) querySnapshotRows(ifaceName string, limit int) SQLRowsWrap {
	if err := db.checkSnapshotTable(); err != nil {
		return result.Fail[*sql.Rows](err)
	}

	rows, err := db.conn.Query(sqlSelectListIPSnapshots, ifaceName, limit)
	if err != nil {
		return result.Fail[*sql.Rows](apperror.WrapSimple(err, "store.querySnapshotRows.Query"))
	}

	return result.Ok(rows)
}

func (db *DB) checkSnapshotTable() *apperror.AppError {
	if db == nil || db.conn == nil {
		return apperror.NewSimple("store.checkSnapshotTable", "E_NIL_CONN")
	}

	return db.EnsureIPSnapshotTable()
}

func scanIPSnapshotRows(rows *sql.Rows) IPSnapshotRecordSliceWrap {
	records := make([]IPSnapshotRecord, 0)
	for rows.Next() {
		recRes := scanOneIPSnapshotRow(rows)
		if recRes.IsFailure() {
			return result.FailSlice[IPSnapshotRecord](recRes.AppError())
		}
		records = append(records, *recRes.Value)
	}

	return result.OkSlice(records)
}

func scanOneIPSnapshotRow(rows *sql.Rows) IPSnapshotRecordWrap {
	var rec IPSnapshotRecord
	var isDHCP int
	if err := scanFields(rows.Scan, &rec, &isDHCP); err != nil {
		return result.Fail[*IPSnapshotRecord](apperror.WrapSimple(err, "store.scanOneIPSnapshotRow"))
	}

	rec.IsDHCP = isDHCP == 1

	return result.Ok(&rec)
}
