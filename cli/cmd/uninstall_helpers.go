package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func isToolTrackedInDB(db *store.DB, tool, canonical string) bool {
	if db != nil && (db.IsToolInstalled(tool) || db.IsToolInstalled(canonical)) {
		return true
	}
	splitDB, err := store.OpenInstallationSplitDB()
	if err != nil {
		return false
	}
	defer splitDB.Close()
	return splitDB.IsToolInstalled(tool) || splitDB.IsToolInstalled(canonical)
}

func isBinaryPresent(tool, canonical string) bool {
	return cmdinstall.GetInstalledVersion(tool) != "" || cmdinstall.GetInstalledVersion(canonical) != ""
}

func checkUninstallEligibility(db *store.DB, tool, canonical string, force bool) error {
	if force {
		return nil
	}
	if isToolTrackedInDB(db, tool, canonical) || isBinaryPresent(tool, canonical) {
		return nil
	}
	return apperror.NewSimple(fmt.Sprintf(constants.ErrUninstallNotFound, tool), "E9000")
}

func dispatchUninstallExecution(db *store.DB, tool, canonical string, dryRun, purge bool) error {
	if isCustomStandaloneTool(tool) || isCustomStandaloneTool(canonical) {
		return executeCustomUninstall(db, tool, canonical, dryRun, purge)
	}
	return executeStandardUninstall(db, tool, dryRun, purge)
}

func executeCustomUninstall(db *store.DB, tool, canonical string, dryRun, purge bool) error {
	if dryRun {
		fmt.Printf("  [dry-run] Would uninstall custom tool %s\n", tool)
		return nil
	}
	if err := uninstallCustomTool(canonical, purge); err != nil {
		return err
	}
	cleanLegacyDBRecords(db, tool, canonical)
	return nil
}

func cleanLegacyDBRecords(db *store.DB, tool, canonical string) {
	if db == nil {
		return
	}
	_ = db.RemoveInstalledTool(tool)
	_ = db.RemoveInstalledTool(canonical)
}
