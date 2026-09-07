package store

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// MigrateInstalledToolsFromRoot copies any existing InstalledTool rows from root DB to split DB.
func MigrateInstalledToolsFromRoot(rootConn *sql.DB, splitDB *InstallationSplitDB) error {
	if rootConn == nil || splitDB == nil {
		return nil
	}

	var tableExists int
	checkErr := rootConn.QueryRow("SELECT COUNT(*) FROM sqlite_master WHERE type='table' AND name='InstalledTool';").Scan(&tableExists)
	if checkErr != nil || tableExists == 0 {
		return nil
	}

	rows, err := rootConn.Query(constants.SQLSelectAllInstalled)
	if err != nil {
		return nil
	}
	defer rows.Close()

	return copyRowsToSplit(rows, splitDB)
}

func copyRowsToSplit(rows *sql.Rows, splitDB *InstallationSplitDB) error {
	for rows.Next() {
		var t InstalledTool
		err := rows.Scan(&t.ID, &t.Tool, &t.VersionMajor, &t.VersionMinor,
			&t.VersionPatch, &t.VersionBuild, &t.VersionString, &t.PackageManager,
			&t.InstallPath, &t.InstalledAt, &t.UpdatedAt)
		if err != nil {
			return apperror.WrapSimple(err, "installation_split.migrateScan")
		}

		_ = splitDB.SaveInstalledTool(t.Tool, t.VersionString, t.PackageManager)
	}

	return nil
}
