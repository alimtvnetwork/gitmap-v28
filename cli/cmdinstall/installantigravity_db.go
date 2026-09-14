package cmdinstall

import (
	"database/sql"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RecordAntigravityDualDatabase records Antigravity and agy in both installation.db and gitmap.db.
func RecordAntigravityDualDatabase(installPath string, durationMs int64, exitCode int) error {
	if errSplit := recordInstallationSplitDB(installPath, durationMs, exitCode); errSplit != nil {
		return errSplit
	}

	return syncGitmapRootDB(installPath)
}

func recordInstallationSplitDB(installPath string, durationMs int64, exitCode int) error {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "installation_split.open")
	}
	defer splitDB.Close()

	if errTools := saveSplitInstalledTools(splitDB, installPath); errTools != nil {
		return errTools
	}
	if errLogs := recordSplitInstallTelemetry(splitDB, exitCode); errLogs != nil {
		return errLogs
	}
	return saveSplitInstallationLog(splitDB, durationMs, exitCode)
}

func recordSplitInstallTelemetry(splitDB *store.InstallationSplitDB, exitCode int) error {
	if exitCode == 0 {
		return recordSplitTelemetrySuccess(splitDB)
	}

	return recordSplitTelemetryFailure(splitDB, exitCode)
}

func recordSplitTelemetrySuccess(splitDB *store.InstallationSplitDB) error {
	if err := splitDB.RecordInstallSuccess("tool", constants.ToolAntigravity); err != nil {
		return err
	}

	return splitDB.RecordInstallSuccess("tool", constants.ToolAgy)
}

func recordSplitTelemetryFailure(splitDB *store.InstallationSplitDB, exitCode int) error {
	errMsg := "Antigravity installation failed"
	if err := splitDB.RecordInstallFailure("tool", constants.ToolAntigravity, exitCode, errMsg); err != nil {
		return err
	}

	return splitDB.RecordInstallFailure("tool", constants.ToolAgy, exitCode, errMsg)
}

// RecordAntigravityInstallStart logs the beginning of an Antigravity installation.
func RecordAntigravityInstallStart() (string, error) {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return "", apperror.WrapSimple(err, "antigravity_install.openStart")
	}
	defer splitDB.Close()

	return splitDB.RecordInstallStart("tool", constants.ToolAntigravity, "install")
}

// RecordAntigravityInstallSuccess logs the successful completion of an Antigravity installation.
func RecordAntigravityInstallSuccess() error {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "antigravity_install.openSuccess")
	}
	defer splitDB.Close()

	return recordSplitTelemetrySuccess(splitDB)
}

// RecordAntigravityInstallFailure logs a failed Antigravity installation.
func RecordAntigravityInstallFailure(exitCode int, errorMsg string) error {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "antigravity_install.openFailure")
	}
	defer splitDB.Close()

	return splitDB.RecordInstallFailure("tool", constants.ToolAntigravity, exitCode, errorMsg)
}

// RecordAntigravityInstallSkipped logs a skipped Antigravity installation.
func RecordAntigravityInstallSkipped(reason string) error {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "antigravity_install.openSkipped")
	}
	defer splitDB.Close()

	return splitDB.RecordInstallSkipped("tool", constants.ToolAntigravity, reason)
}

func saveSplitInstalledTools(splitDB *store.InstallationSplitDB, installPath string) error {
	if err := splitDB.SaveInstalledTool(constants.ToolAntigravity, AntigravityDefaultVersion, "installer"); err != nil {
		return err
	}
	if err := splitDB.SaveInstalledTool(constants.ToolAgy, AntigravityDefaultVersion, "installer"); err != nil {
		return err
	}
	if err := updateSplitInstallPath(splitDB.Conn(), constants.ToolAntigravity, installPath); err != nil {
		return err
	}
	return updateSplitInstallPath(splitDB.Conn(), constants.ToolAgy, installPath)
}

func updateSplitInstallPath(conn *sql.DB, tool, installPath string) error {
	if conn == nil || installPath == "" {
		return nil
	}
	if _, err := conn.Exec("UPDATE InstalledTool SET InstallPath = ? WHERE Tool = ?", installPath, tool); err != nil {
		return apperror.WrapSimple(err, "split_db.update_install_path")
	}
	return nil
}

func saveSplitInstallationLog(splitDB *store.InstallationSplitDB, durationMs int64, exitCode int) error {
	isSuccess := exitCode == 0
	rec := buildInstallationLogRecord(durationMs, exitCode, isSuccess)
	return splitDB.RecordLog(rec)
}

func buildInstallationLogRecord(durationMs int64, exitCode int, isSuccess bool) store.InstallationLogRecord {
	return store.InstallationLogRecord{
		Tool:           constants.ToolAntigravity,
		Action:         "install",
		Version:        AntigravityDefaultVersion,
		PackageManager: "installer",
		DurationMs:     durationMs,
		IsSuccess:      isSuccess,
		ExitCode:       exitCode,
		Notes:          "Antigravity desktop IDE and CLI installation",
	}
}

func syncGitmapRootDB(installPath string) error {
	rootDB, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "gitmap_root.open")
	}
	defer rootDB.Close()

	if errSync := syncRootInstalledTools(rootDB.Conn(), installPath); errSync != nil {
		return errSync
	}
	if errReg := rootDB.SyncKnownSplitDatabases(); errReg != nil {
		return apperror.WrapSimple(errReg, "gitmap_root.syncRegistry")
	}
	return nil
}

func syncRootInstalledTools(conn *sql.DB, installPath string) error {
	if conn == nil {
		return nil
	}
	if err := syncSingleRootTool(conn, constants.ToolAntigravity, installPath); err != nil {
		return err
	}
	return syncSingleRootTool(conn, constants.ToolAgy, installPath)
}

func syncSingleRootTool(conn *sql.DB, tool, installPath string) error {
	_, err := conn.Exec(constants.SQLInsertInstalledTool,
		tool, 2, 13, 0, 0, AntigravityDefaultVersion, "installer", installPath)
	if err != nil {
		return apperror.WrapSimple(err, "root_db.insert_tool")
	}
	return nil
}

// PurgeAntigravityDualDatabase removes Antigravity and agy records from installation.db and gitmap.db.
func PurgeAntigravityDualDatabase() error {
	if errSplit := purgeInstallationSplitDB(); errSplit != nil {
		return errSplit
	}
	return purgeGitmapRootDB()
}

func purgeInstallationSplitDB() error {
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return apperror.WrapSimple(err, "split_db.open_purge")
	}
	defer splitDB.Close()

	if errTool := splitDB.RemoveInstalledTool(constants.ToolAntigravity); errTool != nil {
		return errTool
	}
	return splitDB.RemoveInstalledTool(constants.ToolAgy)
}

func purgeGitmapRootDB() error {
	rootDB, err := store.OpenDefault()
	if err != nil {
		return apperror.WrapSimple(err, "root_db.open_purge")
	}
	defer rootDB.Close()

	if errTools := purgeRootTools(rootDB.Conn()); errTools != nil {
		return errTools
	}
	return rootDB.SyncKnownSplitDatabases()
}

func purgeRootTools(conn *sql.DB) error {
	if conn == nil {
		return nil
	}
	if _, err := conn.Exec(constants.SQLDeleteInstalledTool, constants.ToolAntigravity); err != nil {
		return apperror.WrapSimple(err, "root_db.delete_antigravity")
	}
	if _, err := conn.Exec(constants.SQLDeleteInstalledTool, constants.ToolAgy); err != nil {
		return apperror.WrapSimple(err, "root_db.delete_agy")
	}
	return nil
}
