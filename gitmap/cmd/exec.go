package cmd

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cloner"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// runExec handles the "exec" subcommand.
func runExec(args []string) {
	checkHelp("exec", args)
	groupName, all, stopOnFail, gitArgs := parseExecFlags(args)
	if len(gitArgs) == 0 {
		fmt.Fprintln(os.Stderr, constants.ErrExecUsage)
		os.Exit(1)
	}

	records := loadExecByScope(groupName, all)

	workDir, wdErr := os.Getwd()
	if wdErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not determine working directory: %v\n", wdErr)
	}
	cmdArgs := buildCommandArgs(append([]string{"exec"}, os.Args[2:]...))
	taskID, taskDB := createPendingTask(constants.TaskTypeExec, workDir, workDir, "exec", cmdArgs)

	printExecBanner(gitArgs, len(records))

	prog := cloner.NewBatchProgress(len(records), "Exec", false)
	prog.SetStopOnFail(stopOnFail)
	succeeded, failed, missing := execAllReposTracked(records, gitArgs, prog)
	prog.PrintSummary()
	prog.PrintFailureReport()
	printExecSummary(succeeded, failed, missing, len(records))

	if code := prog.ExitCodeForBatch(); code != 0 {
		failPendingTask(taskDB, taskID, fmt.Sprintf("exec batch failed with exit code %d", code))
		closeTaskDB(taskDB)
		os.Exit(code)
	}

	completePendingTask(taskDB, taskID)
	closeTaskDB(taskDB)
}

// execAllReposTracked runs a git command across all repos with progress.
func execAllReposTracked(records []model.ScanRecord, gitArgs []string, prog *cloner.BatchProgress) (int, int, int) {
	var succeeded, failed, missing int
	for _, rec := range records {
		if prog.Stopped() {
			break
		}
		prog.BeginItem(rec.RepoName)
		s, f, m := execOneRepoTracked(rec, gitArgs, prog)
		succeeded += s
		failed += f
		missing += m
	}

	return succeeded, failed, missing
}

// execOneRepoTracked runs a git command in one repo with progress tracking.
func execOneRepoTracked(rec model.ScanRecord, gitArgs []string, prog *cloner.BatchProgress) (int, int, int) {
	_, err := os.Stat(rec.AbsolutePath)
	if err != nil {
		prog.Skip()

		return 0, 0, 1
	}

	if execInRepo(rec, gitArgs) {
		prog.Succeed()

		return 1, 0, 0
	}

	prog.FailWithError(rec.RepoName, "command exited with non-zero status")

	return 0, 1, 0
}

// execAllRepos runs a git command across all repos and returns totals.
func execAllRepos(records []model.ScanRecord, gitArgs []string) (int, int, int) {
	var succeeded, failed, missing int
	for _, rec := range records {
		s, f, m := execOneRepo(rec, gitArgs)
		succeeded += s
		failed += f
		missing += m
	}

	return succeeded, failed, missing
}

// execOneRepo runs a git command in one repo, returning increment counts.
func execOneRepo(rec model.ScanRecord, gitArgs []string) (int, int, int) {
	_, err := os.Stat(rec.AbsolutePath)
	if err == nil && execInRepo(rec, gitArgs) {
		return 1, 0, 0
	}
	if err == nil {
		return 0, 1, 0
	}

	fmt.Printf(constants.ExecMissingFmt,
		constants.ColorDim, truncate(rec.RepoName, 22),
		constants.ColorYellow, constants.ColorReset)

	return 0, 0, 1
}

// parseExecFlags parses --group, --all, and --stop-on-fail flags, returning remaining args as git args.
func parseExecFlags(args []string) (groupName string, all, stopOnFail bool, gitArgs []string) {
	fs := flag.NewFlagSet(constants.CmdExec, flag.ExitOnError)
	gFlag := fs.String("group", "", constants.FlagDescGroup)
	fs.StringVar(gFlag, "g", "", constants.FlagDescGroup)
	aFlag := fs.Bool("all", false, constants.FlagDescAll)
	sFlag := fs.Bool(constants.FlagStopOnFail, false, constants.FlagDescStopOnFail)
	fs.Parse(args)

	return *gFlag, *aFlag, *sFlag, fs.Args()
}

// loadExecByScope returns records filtered by alias, group, all DB repos, or JSON fallback.
func loadExecByScope(groupName string, all bool) []model.ScanRecord {
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

	return loadExecRecordsJSON()
}

// loadExecRecordsJSON reads ScanRecords from gitmap.json.
func loadExecRecordsJSON() []model.ScanRecord {
	jsonPath := filepath.Join(constants.DefaultOutputFolder, constants.DefaultJSONFile)
	records, err := loadExecRecords(jsonPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrExecLoadFailed, jsonPath, err)
		os.Exit(1)
	}

	return records
}

// loadExecRecords reads ScanRecords from a JSON file.
func loadExecRecords(path string) ([]model.ScanRecord, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var records []model.ScanRecord
	err = json.Unmarshal(data, &records)

	return records, err
}

// truncate shortens a string to max length with ellipsis.
func truncate(s string, max int) string {
	if len(s) <= max {
		return s
	}

	return s[:max-1] + "…"
}
