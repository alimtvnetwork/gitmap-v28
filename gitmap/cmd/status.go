package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// runStatus handles the "status" subcommand.
func runStatus(args []string) {
	checkHelp("status", args)
	groupName, all := parseStatusFlags(args)
	records := loadStatusByScope(groupName, all)

	printStatusBanner(len(records))
	prog := cloner.NewBatchProgress(len(records), "Status", true)
	summary := printStatusTableTracked(records, prog)
	printStatusSummary(summary)
}

// parseStatusFlags parses --group and --all flags.
func parseStatusFlags(args []string) (groupName string, all bool) {
	fs := flag.NewFlagSet(constants.CmdStatus, flag.ExitOnError)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	fs.Parse(args)

	return *gFlag, *aFlag
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
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	records, err := db.ShowGroup(groupName)
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrGenericFmt, err)
		os.Exit(1)
	}

	return records
}

// loadAllRecordsDB loads all repos from the database.
func loadAllRecordsDB() []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrListDBFailed, err)
		os.Exit(1)
	}
	defer db.Close()
	records, err := db.ListRepos()
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrGenericFmt, err)
		os.Exit(1)
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
		fmt.Fprintf(os.Stderr, constants.ErrStatusLoadFailed, jsonPath, err)
		os.Exit(1)
	}

	return records
}

// loadAllRecordsDBOrEmpty returns DB records, or exits with a friendly
// "run gitmap scan first" message when the DB has no repos yet.
func loadAllRecordsDBOrEmpty() []model.ScanRecord {
	db, err := openDB()
	if err != nil {
		fmt.Fprint(os.Stderr, constants.MsgStatusNoData)
		os.Exit(1)
	}
	defer db.Close()
	records, err := db.ListRepos()
	if err != nil {
		if isLegacyDataError(err) {
			fmt.Fprint(os.Stderr, constants.MsgLegacyProjectData)
			os.Exit(1)
		}
		fmt.Fprintf(os.Stderr, constants.ErrGenericFmt, err)
		os.Exit(1)
	}
	if len(records) == 0 {
		fmt.Fprint(os.Stderr, constants.MsgStatusNoData)
		os.Exit(1)
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
