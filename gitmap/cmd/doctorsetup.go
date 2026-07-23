package cmd

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

// checkSetupConfig verifies git-setup.json can be resolved from the binary.
func checkSetupConfig() int {
	configPath := resolveSetupConfigPath(constants.DefaultSetupConfigPath, false)
	_, err := os.Stat(configPath)

	if err != nil {
		printWarn(constants.DoctorSetupConfigMissing)

		return 0
	}

	absPath, absErr := filepath.Abs(configPath)
	if absErr != nil {
		absPath = configPath
	}

	printOK(constants.DoctorSetupConfigOKFmt, absPath)

	return 0
}

// checkShellWrapper verifies the GITMAP_WRAPPER env var is set.
func checkShellWrapper() int {
	if isWrapperActive() {
		printOK(constants.DoctorWrapperOK)

		return 0
	}

	dataDir := store.BinaryDataDir()
	profileDir := filepath.Dir(dataDir)
	printWarn(constants.DoctorWrapperNotLoaded)
	printFix(constants.DoctorWrapperFix)

	_ = profileDir

	return 0
}
