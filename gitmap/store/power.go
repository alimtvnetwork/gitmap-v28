package store

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/power"
)

const (
	sqlInsertPowerProfile = `INSERT OR REPLACE INTO PowerSetting
(ProfileName, Platform, DisplayTimeoutMinutes, SleepTimeoutMinutes, DiskTimeoutMinutes, IsNeverSleep, IsLockDisabled, IsActive, UpdatedAt, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, COALESCE((SELECT CreatedAt FROM PowerSetting WHERE ProfileName = ?), ?));`

	sqlSelectPowerProfile = `SELECT Platform, DisplayTimeoutMinutes, SleepTimeoutMinutes, DiskTimeoutMinutes, IsNeverSleep, IsLockDisabled
FROM PowerSetting WHERE ProfileName = ?;`

	sqlSelectActivePower = `SELECT Platform, DisplayTimeoutMinutes, SleepTimeoutMinutes, DiskTimeoutMinutes, IsNeverSleep, IsLockDisabled
FROM PowerSetting WHERE IsActive = 1 LIMIT 1;`

	sqlInsertPowerHistory = `INSERT INTO PowerSettingHistory
(Action, Platform, DisplayTimeoutMinutes, SleepTimeoutMinutes, IsNeverSleep, Notes, Comments, CreatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?);`

	sqlSelectListPowerHistory = `SELECT PowerSettingHistoryId, Action, Platform, DisplayTimeoutMinutes, SleepTimeoutMinutes, IsNeverSleep, Notes, Comments, CreatedAt
FROM PowerSettingHistory ORDER BY PowerSettingHistoryId DESC LIMIT ?;`
)

// PowerHistoryRecord encapsulates a logged power configuration change.
type PowerHistoryRecord struct {
	ID                    int64
	Action                string
	Platform              string
	DisplayTimeoutMinutes int
	SleepTimeoutMinutes   int
	IsNeverSleep          bool
	Notes                 string
	Comments              string
	CreatedAt             string
}

// SavePowerProfile stores or updates a named power configuration profile.
func (db *DB) SavePowerProfile(profile string, s power.Settings, isActive bool) error {
	if err := EnsurePowerTables(db.conn); err != nil {
		return err
	}

	if err := db.maybeDeactivateProfiles(isActive); err != nil {
		return err
	}

	return db.insertPowerProfile(profile, s, isActive)
}

func (db *DB) maybeDeactivateProfiles(isActive bool) error {
	if !isActive {
		return nil
	}

	return db.deactivateAllProfiles()
}

func (db *DB) deactivateAllProfiles() error {
	_, err := db.conn.Exec("UPDATE PowerSetting SET IsActive = 0;")
	if err != nil {
		return apperror.WrapSimple(err, "store.savePowerProfile.DeactivatePrior")
	}

	return nil
}

func (db *DB) insertPowerProfile(profile string, s power.Settings, isActive bool) error {
	activeInt := 0
	if isActive {
		activeInt = 1
	}

	now := time.Now().UTC().Format(time.RFC3339)
	neverInt, lockInt := boolToInt(s.IsNeverSleep), boolToInt(s.IsLockDisabled)
	_, err := db.conn.Exec(sqlInsertPowerProfile, profile, s.Platform,
		s.DisplayTimeoutMinutes, s.SleepTimeoutMinutes, s.DiskTimeoutMinutes,
		neverInt, lockInt, activeInt, now, profile, now)
	if err != nil {
		return apperror.WrapSimple(err, "store.savePowerProfile")
	}

	return nil
}

// GetPowerProfile loads a power profile by name.
func (db *DB) GetPowerProfile(profile string) (power.Settings, error) {
	if err := EnsurePowerTables(db.conn); err != nil {
		return power.Settings{}, err
	}

	var s power.Settings
	var neverInt, lockInt int

	err := db.conn.QueryRow(sqlSelectPowerProfile, profile).Scan(
		&s.Platform, &s.DisplayTimeoutMinutes, &s.SleepTimeoutMinutes,
		&s.DiskTimeoutMinutes, &neverInt, &lockInt,
	)
	if err != nil {
		return s, apperror.WrapSimple(err, "store.getPowerProfile")
	}

	s.IsNeverSleep, s.IsLockDisabled = neverInt == 1, lockInt == 1
	s.Source = "profile:" + profile

	return s, nil
}

// GetActivePowerSetting loads the currently active power profile.
func (db *DB) GetActivePowerSetting() (power.Settings, error) {
	if err := EnsurePowerTables(db.conn); err != nil {
		return power.Settings{}, err
	}

	var s power.Settings
	var neverInt, lockInt int

	err := db.conn.QueryRow(sqlSelectActivePower).Scan(
		&s.Platform, &s.DisplayTimeoutMinutes, &s.SleepTimeoutMinutes,
		&s.DiskTimeoutMinutes, &neverInt, &lockInt,
	)
	if err != nil {
		return s, apperror.WrapSimple(err, "store.getActivePowerSetting")
	}

	s.IsNeverSleep, s.IsLockDisabled = neverInt == 1, lockInt == 1

	return s, nil
}

// RecordPowerHistory adds an audit log entry for a power change action.
func (db *DB) RecordPowerHistory(action string, s power.Settings, note string) error {
	if err := EnsurePowerTables(db.conn); err != nil {
		return err
	}

	now := time.Now().UTC().Format(time.RFC3339)
	neverInt := boolToInt(s.IsNeverSleep)

	_, err := db.conn.Exec(sqlInsertPowerHistory, action, s.Platform,
		s.DisplayTimeoutMinutes, s.SleepTimeoutMinutes, neverInt, note, "", now)
	if err != nil {
		return apperror.WrapSimple(err, "store.recordPowerHistory")
	}

	return nil
}

// ListPowerHistory queries the most recent power history records.
func (db *DB) ListPowerHistory(limit int) ([]PowerHistoryRecord, error) {
	if err := EnsurePowerTables(db.conn); err != nil {
		return nil, err
	}

	if limit <= 0 {
		limit = 50
	}

	rows, err := db.conn.Query(sqlSelectListPowerHistory, limit)
	if err != nil {
		return nil, apperror.WrapSimple(err, "store.listPowerHistory")
	}

	defer rows.Close()

	return scanPowerHistoryRows(rows)
}

func scanPowerHistoryRows(rows *sql.Rows) ([]PowerHistoryRecord, error) {
	var records []PowerHistoryRecord
	for rows.Next() {
		var r PowerHistoryRecord
		var neverInt int
		err := rows.Scan(&r.ID, &r.Action, &r.Platform, &r.DisplayTimeoutMinutes,
			&r.SleepTimeoutMinutes, &neverInt, &r.Notes, &r.Comments, &r.CreatedAt)
		if err != nil {
			return nil, apperror.WrapSimple(err, "store.scanPowerHistory")
		}

		r.IsNeverSleep = neverInt == 1
		records = append(records, r)
	}

	return records, nil
}
