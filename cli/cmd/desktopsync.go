package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/desktop"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// runDesktopSync handles the "desktop-sync" subcommand.
func runDesktopSync(args []string) error {
	if err := ensureGHDesktopInstalled(args); err != nil {
		return err
	}

	return executeDesktopSync()
}

func executeDesktopSync() error {
	outputDir := constants.DefaultOutputFolder
	jsonPath := filepath.Join(outputDir, constants.DefaultJSONFile)
	if appErr := validateDesktopSyncPaths(outputDir, jsonPath); appErr != nil {
		return appErr
	}

	records, appErr := loadDesktopRecords(jsonPath)
	if appErr != nil {
		return appErr
	}

	return syncToDesktop(records, jsonPath)
}

// validateDesktopSyncPaths checks that the output dir and JSON file exist.
func validateDesktopSyncPaths(outputDir, jsonPath string) *apperror.AppError {
	info, err := os.Stat(outputDir)
	if err != nil || !info.IsDir() {
		return apperror.NewSimple(constants.MsgNoOutputDir, "E9000")
	}

	_, jsonErr := os.Stat(jsonPath)
	if jsonErr != nil {
		return apperror.NewSimple(constants.MsgNoJSONFile, "E9000")
	}

	return nil
}

// loadDesktopRecords reads and parses the JSON file into ScanRecords.
func loadDesktopRecords(path string) ([]model.ScanRecord, *apperror.AppError) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, apperror.WrapSimple(err, constants.ErrDesktopReadFailed)
	}

	var records []model.ScanRecord
	err = json.Unmarshal(data, &records)
	if err != nil {
		return nil, apperror.WrapSimple(err, constants.ErrDesktopParseFailed)
	}

	return records, nil
}

// syncToDesktop registers each repo with GitHub Desktop.
func syncToDesktop(records []model.ScanRecord, source string) *apperror.AppError {
	cli := desktop.ResolveCLI()
	if cli == "" {
		desktop.PrintInstallSuggestions()
		return desktop.NewMissingCLIError()
	}

	fmt.Printf(constants.MsgDesktopSyncStart, source)
	added, skipped, failed := syncAll(records, cli)
	fmt.Printf(constants.MsgDesktopSyncDone, added, skipped, failed)

	return nil
}

// syncAll iterates records and syncs each to GitHub Desktop.
func syncAll(records []model.ScanRecord, cli string) (added, skipped, failed int) {
	for _, r := range records {
		result := syncOne(r, cli)
		added, skipped, failed = tallyResult(result, added, skipped, failed)
	}

	return added, skipped, failed
}
