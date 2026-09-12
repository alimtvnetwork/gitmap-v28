package cmd

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/gitutil"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runStatus handles the "status" subcommand.
func runStatus(args []string) error {
	checkHelp("status", args)
	groupName, all, onlyDirty := parseStatusFlags(args)
	records := loadStatusByScope(groupName, all)

	if onlyDirty {
		records = filterDirtyRecords(records)
	}

	if onlyDirty && len(records) == 0 {
		fmt.Println("✨ All repositories are !clean")

		return nil
	}

	return renderStatusView(records)
}

func renderStatusView(records []model.ScanRecord) error {
	printStatusBanner(len(records))
	prog := cloner.NewBatchProgress(len(records), "Status", true)
	summary := printStatusTableTracked(records, prog)
	printStatusSummary(summary)

	return nil
}

func filterDirtyRecords(records []model.ScanRecord) []model.ScanRecord {
	var dirty []model.ScanRecord
	for _, rec := range records {
		rs := gitutil.Status(rec.AbsolutePath)
		if rs.Dirty || rs.Ahead > 0 || rs.Behind > 0 || rs.StashCount > 0 {
			dirty = append(dirty, rec)
		}
	}

	return dirty
}

// parseStatusFlags parses --group, --all, and --dirty flags.
func parseStatusFlags(args []string) (groupName string, all bool, onlyDirty bool) {
	fs := flag.NewFlagSet(constants.CmdStatus, flag.ExitOnError)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	dFlag := fs.Bool("dirty", false, "Display only repositories with uncommitted, unstaged, or unpushed changes")
	fs.BoolVar(dFlag, "only-dirty", false, "Display only repositories with uncommitted, unstaged, or unpushed changes")
	fs.Parse(args)

	return *gFlag, *aFlag, *dFlag
}

// loadStatusByScope returns records filtered by alias, group, all DB repos, or JSON fallback.
func loadStatusByScope(groupName string, all bool) []model.ScanRecord {
	if HasAlias() {
		return []model.ScanRecord{{
			RepoName:     GetAliasSlug(),
			Slug:         GetAliasSlug(),
			AbsolutePath: GetAliasPath(),
		}}
	}

	if len(groupName) > 0 {
		return loadRecordsByGroup(groupName)
	}

	if all {
		return loadAllRecordsDB()
	}

	return loadRecordsJSONFallback()
}

// loadRecordsByGroup loads repos from a specific group in the database.
func loadRecordsByGroup(groupName string) []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		appErr := buildStatusDBOpenError("loadGroup", groupName, err)
		cliexit.HandleGeneralError(appErr)

		return nil
	}

	defer db.Close()

	records, err := db.ShowGroup(groupName)
	if err != nil {
		handleStatusDBError(err)
	}

	return records
}

func buildStatusDBOpenError(op, group string, err error) *apperror.AppError {
	ctx := map[string]any(nil)
	if group != "" {
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

// loadAllRecordsDB loads all repos from the database.
func loadAllRecordsDB() []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.status.loadAll.openDB",
			"E1084",
			"failed to open database for status load",
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			nil,
		)
		cliexit.HandleGeneralError(appErr)

		return nil
	}

	defer db.Close()

	records, err := db.ListRepos()
	if err != nil {
		handleStatusDBError(err)
	}

	return records
}

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

// statusSummary aggregates counts across all repos.
type statusSummary struct {
	Total   int
	Clean   int
	Dirty   int
	Ahead   int
	Behind  int
	Stashed int
	Missing int
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
