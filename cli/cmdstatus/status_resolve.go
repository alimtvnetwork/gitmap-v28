package cmdstatus

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// loadStatusByScope returns records filtered by alias, group, all DB repos, or JSON fallback.
func loadStatusByScope(groupName string, all bool) ([]model.ScanRecord, error) {
	if HasAlias() {
		return []model.ScanRecord{{
			RepoName:     GetAliasSlug(),
			Slug:         GetAliasSlug(),
			AbsolutePath: GetAliasPath(),
		}}, nil
	}

	hasGroup := len(groupName) > 0
	if hasGroup {
		return loadRecordsByGroup(groupName)
	}

	if all {
		return loadAllRecordsDBWithErr()
	}

	return loadRecordsJSONFallbackWithErr()
}

func loadRecordsByGroup(groupName string) ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		return nil, buildStatusDBOpenError("loadGroup", groupName, err)
	}

	defer db.Close()

	records, err := db.ShowGroup(groupName)
	if err != nil {
		return nil, handleStatusDBError(err)
	}

	return records, nil
}

func buildStatusDBOpenError(op, group string, err error) *apperror.AppError {
	ctx := map[string]any(nil)
	hasGroup := group != ""
	if hasGroup {
		ctx = map[string]any{"group": group}
	}

	return apperror.WrapWithDetails(
		err,
		"cmd.status."+op+".openDB",
		"E1083",
		"failed to open database for status group load",
		"cmd.status",
		apperror.ErrorTypeExecution,
		apperror.SeverityFatal,
		ctx,
	)
}

// loadAllRecordsDB loads all repos from the database, matching hook signature.
func loadAllRecordsDB() []model.ScanRecord {
	records, _ := loadAllRecordsDBWithErr()

	return records
}

func loadAllRecordsDBWithErr() ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		return nil, apperror.WrapWithDetails(
			err,
			"cmd.status.loadAll.openDB",
			"E1084",
			"failed to open database for status load",
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
	}

	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		return nil, handleStatusDBError(err)
	}

	return records, nil
}

// LoadRecordsJSONFallback loads records from .gitmap/output/gitmap.json for hook wiring.
func LoadRecordsJSONFallback() []model.ScanRecord {
	records, _ := loadRecordsJSONFallbackWithErr()

	return records
}

func loadRecordsJSONFallbackWithErr() ([]model.ScanRecord, error) {
	jsonPath := filepath.Join(constants.DefaultOutputDir, constants.DefaultJSONFile)
	_, statErr := os.Stat(jsonPath)
	isMissing := os.IsNotExist(statErr)
	if isMissing {
		return loadAllRecordsDBOrEmpty()
	}

	records, err := loadStatusRecords(jsonPath)
	if err != nil {
		return nil, apperror.WrapWithDetails(
			err,
			"cmd.status.loadJSON",
			"E1085",
			"failed to read status records from JSON",
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			map[string]any{"path": jsonPath},
		)
	}

	return records, nil
}

func loadAllRecordsDBOrEmpty() ([]model.ScanRecord, error) {
	db, err := openDB()
	if err != nil {
		return nil, newStatusNoDataError("openDB", "E1086")
	}

	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		return nil, handleStatusDBError(err)
	}

	hasEmpty := len(records) == 0
	if hasEmpty {
		return nil, newStatusNoDataError("noRepos", "E1087")
	}

	return records, nil
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

func loadStatusRecords(path string) ([]model.ScanRecord, error) {
	return model.LoadStatusRecords(path)
}

func handleStatusDBError(err error) error {
	isLegacy := isLegacyDataError(err)
	if isLegacy {
		return apperror.WrapWithDetails(
			err,
			"cmd.status.legacyData",
			"E1088",
			constants.MsgLegacyProjectData,
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityError,
			nil,
		)
	}

	return apperror.WrapWithDetails(
		err,
		"cmd.status.dbError",
		"E1089",
		"database operation failed during status lookup",
		"cmd.status",
		apperror.ErrorTypeExecution,
		apperror.SeverityError,
		nil,
	)
}
