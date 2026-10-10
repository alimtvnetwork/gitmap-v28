package cmdwatch

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
	"os"
	"path/filepath"
	"strings"
)

// loadRecordsJSONFallback loads records from .gitmap/output/gitmap.json.
func loadRecordsJSONFallback() []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile)
	if _, statErr := os.Stat(jsonPath); os.IsNotExist(statErr) {
		return loadAllRecordsDBOrEmpty()
	}

	records, err := loadStatusRecords(jsonPath)
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.status.loadJSON",
			"E1085",
			"failed to read status records from JSON",
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": jsonPath},
		)
		cliexit.HandleGeneralError(appErr)

		return nil
	}

	return records
}

func openDB() (*store.DB, error) {
	db, err := store.OpenDefault()
	if err != nil {
		return store.OpenGlobalDefault()
	}
	_ = db.Migrate()
	if countRegisteredSSHHosts(db) == 0 {
		return resolveGlobalSSHDBFallback(db), nil
	}
	return db, nil
}

// truncate shortens a string to max length with ellipsis.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-1] + "…"
}

func countRegisteredSSHHosts(db *store.DB) int {
	var count int
	row := db.SQL().QueryRow("SELECT count(*) FROM ssh_hosts")
	if scanErr := row.Scan(&count); scanErr != nil {
		return 0
	}
	return count
}
func resolveGlobalSSHDBFallback(localDB *store.DB) *store.DB {
	globalDB, err := store.OpenGlobalDefault()
	if err != nil {
		return localDB
	}
	_ = globalDB.Migrate()
	if countRegisteredSSHHosts(globalDB) > 0 {
		localDB.Close()
		return globalDB
	}
	globalDB.Close()
	return localDB
}

// loadAllRecordsDBOrEmpty returns DB records, or exits with a friendly
// "run gitmap scan first" message when the DB has no repos yet.
func loadAllRecordsDBOrEmpty() []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		cliexit.HandleGeneralError(newStatusNoDataError("openDB", "E1086"))

		return nil
	}

	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		handleStatusDBError(err)
	}

	if len(records) == 0 {
		cliexit.HandleGeneralError(newStatusNoDataError("noRepos", "E1087"))

		return nil
	}

	return records
}

func newStatusNoDataError(op, code string) *apperror.AppError {
	return apperror.NewWithDetails(
		"cmd.status."+op,
		code,
		constants.MsgStatusNoData,
		"cmd.status",
		apperror.ErrorTypePrecondition,
		apperror.SeverityError,
		nil,
	)
}

// loadStatusRecords reads ScanRecords from gitmap.json.
func loadStatusRecords(path string) ([]model.ScanRecord, error) {
	return model.LoadStatusRecords(path)
}

// isLegacyDataError checks if an error indicates legacy UUID-format data.
func isLegacyDataError(err error) bool {
	return strings.Contains(err.Error(), "Scan error") ||
		strings.Contains(err.Error(), "converting driver.Value type string")
}

func handleStatusDBError(err error) {
	if isLegacyDataError(err) {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.status.legacyData",
			"E1088",
			constants.MsgLegacyProjectData,
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleGeneralError(appErr)

		return
	}

	appErr := apperror.WrapWithDetails(
		err,
		"cmd.status.dbError",
		"E1089",
		"database operation failed during status lookup",
		"cmd.status",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
	cliexit.HandleGeneralError(appErr)
}
