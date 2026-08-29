package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/desktop"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runDesktopSync handles the "desktop-sync" subcommand.
func runDesktopSync() error {
	outputDir := constants.DefaultOutputFolder
	jsonPath := filepath.Join(outputDir, constants.DefaultJSONFile)
	validateDesktopSyncPaths(outputDir, jsonPath)
	records := loadDesktopRecords(jsonPath)
	syncToDesktop(records, jsonPath)
	return nil
}

// validateDesktopSyncPaths checks that the output dir and JSON file exist.
func validateDesktopSyncPaths(outputDir, jsonPath string) {
	info, err := os.Stat(outputDir)
	if err != nil || !info.IsDir() {
		apperror.NewSimple(constants.MsgNoOutputDir, "E9000")
		return
	}
	_, jsonErr := os.Stat(jsonPath)
	if jsonErr != nil {
		apperror.NewSimple(constants.MsgNoJSONFile, "E9000")
		return
	}
}

// loadDesktopRecords reads and parses the JSON file into ScanRecords.
func loadDesktopRecords(path string) []model.ScanRecord {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintln(os.Stderr, apperror.NewSimple(constants.ErrDesktopReadFailed, "E9000").Error())
		cliexit.HandleError(nil, 1)
	}
	var records []model.ScanRecord
	err = json.Unmarshal(data, &records)
	if err != nil {
		fmt.Fprintln(os.Stderr, apperror.NewSimple(constants.ErrDesktopParseFailed, "E9000").Error())
		cliexit.HandleError(nil, 1)
	}

	return records
}

// syncToDesktop registers each repo with GitHub Desktop.
func syncToDesktop(records []model.ScanRecord, source string) {
	cli := desktop.ResolveCLI()
	if cli == "" {
		apperror.NewSimple(constants.MsgDesktopNotFound, "E9000")
		return
	}
	fmt.Printf(constants.MsgDesktopSyncStart, source)
	added, skipped, failed := syncAll(records, cli)
	fmt.Printf(constants.MsgDesktopSyncDone, added, skipped, failed)
}

// syncAll iterates records and syncs each to GitHub Desktop.
func syncAll(records []model.ScanRecord, cli string) (added, skipped, failed int) {
	for _, r := range records {
		result := syncOne(r, cli)
		added, skipped, failed = tallyResult(result, added, skipped, failed)
	}

	return added, skipped, failed
}

// syncResult represents the outcome of syncing one repo.
type syncResult int

const (
	syncAdded syncResult = iota
	syncSkipped
	syncFailed
)

// syncOne attempts to register a single repo with GitHub Desktop.
func syncOne(r model.ScanRecord, cli string) syncResult {
	if len(r.AbsolutePath) == 0 {
		fmt.Printf(constants.MsgDesktopSyncFailed, r.RepoName, constants.ErrNoAbsPath)

		return syncFailed
	}

	return syncExistingPath(r, cli)
}

// syncExistingPath checks path existence and registers with Desktop.
func syncExistingPath(r model.ScanRecord, cli string) syncResult {
	_, err := os.Stat(r.AbsolutePath)
	if err == nil {
		return registerOne(r.RepoName, r.AbsolutePath, cli)
	}
	fmt.Printf(constants.MsgDesktopSyncSkipped, r.RepoName)

	return syncSkipped
}

// registerOne calls the GitHub Desktop CLI for a single repo.
func registerOne(name, repoPath, cli string) syncResult {
	cmd := exec.Command(cli, repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(constants.MsgDesktopSyncFailed, name, fmt.Sprintf("%v: %s", err, output))

		return syncFailed
	}
	fmt.Printf(constants.MsgDesktopSyncAdded, name)

	return syncAdded
}

// tallyResult increments the appropriate counter.
func tallyResult(r syncResult, added, skipped, failed int) (int, int, int) {
	if r == syncAdded {
		return added + 1, skipped, failed
	}
	if r == syncSkipped {
		return added, skipped + 1, failed
	}

	return added, skipped, failed + 1
}
