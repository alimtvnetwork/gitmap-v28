package store

import (
	"database/sql"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

const (
	sqlUpsertProfileStart = `INSERT INTO ProfileInstallation
(ProfileName, ProfileAlias, Status, IsSuccess, TotalTools, Description, InstalledAt, UpdatedAt)
VALUES (?, ?, 'installing', 0, ?, ?, ?, ?)
ON CONFLICT(ProfileName) DO UPDATE SET
ProfileAlias=excluded.ProfileAlias,
Status='installing',
IsSuccess=0,
TotalTools=excluded.TotalTools,
Description=excluded.Description,
UpdatedAt=excluded.UpdatedAt;`

	sqlSelectProfileIdByName = `SELECT ProfileInstallationId FROM ProfileInstallation WHERE ProfileName = ?;`

	sqlUpdateProfileCompletion = `UPDATE ProfileInstallation
SET Status = ?, IsSuccess = ?, DurationMs = ?, InstalledCount = ?, FailedCount = ?,
StackTrace = ?, ErrorLog = ?, Notes = ?, UpdatedAt = ?
WHERE ProfileInstallationId = ?;`

	sqlSelectProfileByName = `SELECT ProfileInstallationId, ProfileName, ProfileAlias, Status, IsSuccess,
DurationMs, TotalTools, InstalledCount, FailedCount, COALESCE(StackTrace, ''),
COALESCE(ErrorLog, ''), COALESCE(Description, ''), COALESCE(Notes, ''), COALESCE(Comments, ''),
InstalledAt, UpdatedAt
FROM ProfileInstallation WHERE ProfileName = ?;`

	sqlCountInstalledProfile = `SELECT COUNT(*) FROM ProfileInstallation WHERE ProfileName = ? AND Status = 'installed' AND IsSuccess = 1;`

	sqlSelectAllProfiles = `SELECT ProfileInstallationId, ProfileName, ProfileAlias, Status, IsSuccess,
DurationMs, TotalTools, InstalledCount, FailedCount, COALESCE(StackTrace, ''),
COALESCE(ErrorLog, ''), COALESCE(Description, ''), COALESCE(Notes, ''), COALESCE(Comments, ''),
InstalledAt, UpdatedAt
FROM ProfileInstallation ORDER BY ProfileInstallationId ASC;`

	sqlDeleteProfileByName = `DELETE FROM ProfileInstallation WHERE ProfileName = ?;`

	sqlInsertPackageInstallation = `INSERT INTO PackageInstallation
(ProfileInstallationId, PackageName, Version, PackageManager, InstallPath, Status,
IsSuccess, ExitCode, DurationMs, Stdout, Stderr, StackTrace, CommandLine, Description,
Notes, Comments, InstalledAt, UpdatedAt)
VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?);`

	sqlSelectPackagesByProfileId = `SELECT PackageInstallationId, COALESCE(ProfileInstallationId, 0),
PackageName, Version, PackageManager, InstallPath, Status, IsSuccess, ExitCode,
DurationMs, COALESCE(Stdout, ''), COALESCE(Stderr, ''), COALESCE(StackTrace, ''),
COALESCE(CommandLine, ''), COALESCE(Description, ''), COALESCE(Notes, ''),
COALESCE(Comments, ''), InstalledAt, UpdatedAt
FROM PackageInstallation WHERE ProfileInstallationId = ? ORDER BY PackageInstallationId ASC;`
)

// ProfileInstallationRecord represents a tracked profile installation event.
type ProfileInstallationRecord struct {
	ProfileInstallationId int64  `json:"profileInstallationId"`
	ProfileName           string `json:"profileName"`
	ProfileAlias          string `json:"profileAlias"`
	Status                string `json:"status"`
	IsSuccess             bool   `json:"isSuccess"`
	DurationMs            int64  `json:"durationMs"`
	TotalTools            int    `json:"totalTools"`
	InstalledCount        int    `json:"installedCount"`
	FailedCount           int    `json:"failedCount"`
	StackTrace            string `json:"stackTrace,omitempty"`
	ErrorLog              string `json:"errorLog,omitempty"`
	Description           string `json:"description,omitempty"`
	Notes                 string `json:"notes,omitempty"`
	Comments              string `json:"comments,omitempty"`
	InstalledAt           string `json:"installedAt"`
	UpdatedAt             string `json:"updatedAt"`
}

// PackageInstallationRecord represents execution telemetry for a profile package step.
type PackageInstallationRecord struct {
	PackageInstallationId int64  `json:"packageInstallationId"`
	ProfileInstallationId int64  `json:"profileInstallationId,omitempty"`
	PackageName           string `json:"packageName"`
	Version               string `json:"version"`
	PackageManager        string `json:"packageManager"`
	InstallPath           string `json:"installPath"`
	Status                string `json:"status"`
	IsSuccess             bool   `json:"isSuccess"`
	ExitCode              int    `json:"exitCode"`
	DurationMs            int64  `json:"durationMs"`
	Stdout                string `json:"stdout,omitempty"`
	Stderr                string `json:"stderr,omitempty"`
	StackTrace            string `json:"stackTrace,omitempty"`
	CommandLine           string `json:"commandLine,omitempty"`
	Description           string `json:"description,omitempty"`
	Notes                 string `json:"notes,omitempty"`
	Comments              string `json:"comments,omitempty"`
	InstalledAt           string `json:"installedAt"`
	UpdatedAt             string `json:"updatedAt"`
}

// RecordProfileStart inserts or initializes an in-progress profile installation record.
func (s *InstallationSplitDB) RecordProfileStart(profileName, alias, description string, totalTools int) (int64, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	_, err := s.conn.Exec(sqlUpsertProfileStart, profileName, alias, totalTools, description, now, now)
	if err != nil {
		return 0, apperror.WrapSimple(err, "installation_split.recordProfileStart.exec")
	}

	return s.queryProfileId(profileName)
}

func (s *InstallationSplitDB) queryProfileId(profileName string) (int64, error) {
	var id int64
	err := s.conn.QueryRow(sqlSelectProfileIdByName, profileName).Scan(&id)
	if err != nil {
		return 0, apperror.WrapSimple(err, "installation_split.queryProfileId")
	}

	return id, nil
}

// RecordProfileCompletion updates the profile installation status, duration, counts, and error metadata.
func (s *InstallationSplitDB) RecordProfileCompletion(
	profileId int64, durationMs int64, isSuccess bool,
	installedCount, failedCount int, stackTrace, errorLog, notes string,
) error {
	status := resolveCompletionStatus(isSuccess)
	now := time.Now().UTC().Format(time.RFC3339)
	successInt := boolToInt(isSuccess)
	sanitizedStack := SanitizeLogOutput(stackTrace)
	sanitizedLog := SanitizeLogOutput(errorLog)

	_, err := s.conn.Exec(sqlUpdateProfileCompletion, status, successInt, durationMs,
		installedCount, failedCount, sanitizedStack, sanitizedLog, notes, now, profileId)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.recordProfileCompletion")
	}

	return nil
}

func resolveCompletionStatus(isSuccess bool) string {
	if isSuccess {
		return "installed"
	}

	return "failed"
}

// GetProfileInstallation retrieves a single profile installation record by name.
func (s *InstallationSplitDB) GetProfileInstallation(profileName string) (*ProfileInstallationRecord, error) {
	row := s.conn.QueryRow(sqlSelectProfileByName, profileName)
	rec, err := scanProfileRow(row)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.getProfileInstallation")
	}

	return rec, nil
}

func scanProfileRow(scanner interface{ Scan(dest ...any) error }) (*ProfileInstallationRecord, error) {
	var r ProfileInstallationRecord
	var successInt int
	err := scanner.Scan(
		&r.ProfileInstallationId, &r.ProfileName, &r.ProfileAlias, &r.Status,
		&successInt, &r.DurationMs, &r.TotalTools, &r.InstalledCount, &r.FailedCount,
		&r.StackTrace, &r.ErrorLog, &r.Description, &r.Notes, &r.Comments,
		&r.InstalledAt, &r.UpdatedAt,
	)
	if err != nil {
		return nil, err
	}
	r.IsSuccess = successInt == 1

	return &r, nil
}

// IsProfileInstalled reports whether a profile is successfully installed.
func (s *InstallationSplitDB) IsProfileInstalled(profileName string) bool {
	var count int
	err := s.conn.QueryRow(sqlCountInstalledProfile, profileName).Scan(&count)
	if err != nil {
		return false
	}

	return count > 0
}

// ListProfileInstallations returns all tracked profile installations from installation.db.
func (s *InstallationSplitDB) ListProfileInstallations() ([]ProfileInstallationRecord, error) {
	rows, err := s.conn.Query(sqlSelectAllProfiles)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.listProfileInstallations")
	}

	defer rows.Close()

	return scanProfileList(rows)
}

func scanProfileList(rows *sql.Rows) ([]ProfileInstallationRecord, error) {
	var list []ProfileInstallationRecord
	for rows.Next() {
		rec, err := scanProfileRow(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "installation_split.scanProfileList")
		}

		list = append(list, *rec)
	}

	return list, nil
}

// DeleteProfileInstallation deletes a profile installation record by name.
func (s *InstallationSplitDB) DeleteProfileInstallation(profileName string) error {
	_, err := s.conn.Exec(sqlDeleteProfileByName, profileName)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.deleteProfileInstallation")
	}

	return nil
}

// RecordPackageInstallation records execution telemetry for a profile package step.
func (s *InstallationSplitDB) RecordPackageInstallation(rec PackageInstallationRecord) error {
	now := resolveCreatedAt(rec.InstalledAt)
	updatedAt := resolveCreatedAt(rec.UpdatedAt)
	successInt := boolToInt(rec.IsSuccess)
	profID := resolveNullableProfileId(rec.ProfileInstallationId)

	_, err := s.conn.Exec(sqlInsertPackageInstallation, profID, rec.PackageName,
		rec.Version, rec.PackageManager, rec.InstallPath, rec.Status, successInt,
		rec.ExitCode, rec.DurationMs, SanitizeLogOutput(rec.Stdout), SanitizeLogOutput(rec.Stderr),
		SanitizeLogOutput(rec.StackTrace), rec.CommandLine, rec.Description, rec.Notes,
		rec.Comments, now, updatedAt)
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.recordPackageInstallation")
	}

	return nil
}

func resolveNullableProfileId(id int64) *int64 {
	if id > 0 {
		return &id
	}

	return nil
}

// ListPackagesForProfile returns all tracked packages associated with a profile installation.
func (s *InstallationSplitDB) ListPackagesForProfile(profileId int64) ([]PackageInstallationRecord, error) {
	rows, err := s.conn.Query(sqlSelectPackagesByProfileId, profileId)
	if err != nil {
		return nil, apperror.WrapSimple(err, "installation_split.listPackagesForProfile")
	}

	defer rows.Close()

	return scanPackageList(rows)
}

func scanPackageList(rows *sql.Rows) ([]PackageInstallationRecord, error) {
	var list []PackageInstallationRecord
	for rows.Next() {
		rec, err := scanPackageRow(rows)
		if err != nil {
			return nil, apperror.WrapSimple(err, "installation_split.scanPackageList")
		}

		list = append(list, rec)
	}

	return list, nil
}

func scanPackageRow(rows *sql.Rows) (PackageInstallationRecord, error) {
	var r PackageInstallationRecord
	var successInt int
	var profID sql.NullInt64
	err := rows.Scan(&r.PackageInstallationId, &profID, &r.PackageName, &r.Version,
		&r.PackageManager, &r.InstallPath, &r.Status, &successInt, &r.ExitCode,
		&r.DurationMs, &r.Stdout, &r.Stderr, &r.StackTrace, &r.CommandLine,
		&r.Description, &r.Notes, &r.Comments, &r.InstalledAt, &r.UpdatedAt)
	if err != nil {
		return r, err
	}
	if profID.Valid {
		r.ProfileInstallationId = profID.Int64
	}
	r.IsSuccess = successInt == 1

	return r, nil
}
