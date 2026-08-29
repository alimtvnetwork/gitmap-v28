package cmd

import (
	"encoding/json"
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
		if len(records) == 0 {
			fmt.Println("✨ All repositories are clean!")
			return nil
		}
	}

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
		appErr := apperror.WrapWithDetails(
			err,
			"cmd.status.loadGroup.openDB",
			"E1083",
			"failed to open database for status group load",
			"cmd.status",
			apperror.ErrorTypeExecution,
			apperror.SeverityFatal,
			map[string]any{"group": groupName},
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	defer db.Close()
	records, err := db.ShowGroup(groupName)
	if err != nil {
		handleStatusDBError(err)
	}

	return records
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
		cliexit.HandleError(appErr, 1)
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
// If the JSON file is missing (e.g. user has not run `gitmap scan` from this
// exact directory), fall through to the database — the DB is the source of
// truth post-v2 and usually has every repo the user has ever scanned.
//
// Bug fix (v3.32.0): previously this looked at the legacy bare "output/"
// path AND exited with an error when the file was missing, even though the
// DB had perfectly good data. Users hit this whenever they ran `gitmap status`
// from a directory they had never scanned (e.g. a parent shell prompt).
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
		cliexit.HandleError(appErr, 1)
		return nil
	}

	return records
}

// loadAllRecordsDBOrEmpty returns DB records, or exits with a friendly
// "run gitmap scan first" message when the DB has no repos yet.
func loadAllRecordsDBOrEmpty() []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		appErr := apperror.NewWithDetails(
			"cmd.status.openDB",
			"E1086",
			constants.MsgStatusNoData,
			"cmd.status",
			apperror.ErrorTypePrecondition,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}
	defer db.Close()
	records, err := db.ListRepos()
	if err != nil {
		handleStatusDBError(err)
	}
	if len(records) == 0 {
		appErr := apperror.NewWithDetails(
			"cmd.status.noRepos",
			"E1087",
			constants.MsgStatusNoData,
			"cmd.status",
			apperror.ErrorTypePrecondition,
			apperror.SeverityError,
			nil,
		)
		cliexit.HandleError(appErr, 1)
		return nil
	}

	return records
}

// loadStatusRecords reads ScanRecords from gitmap.json.
func loadStatusRecords(path string) ([]model.ScanRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var records []model.ScanRecord
	err = json.Unmarshal(data, &records)

	return records, err
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
		cliexit.HandleError(appErr, 1)
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
	cliexit.HandleError(appErr, 1)
}
