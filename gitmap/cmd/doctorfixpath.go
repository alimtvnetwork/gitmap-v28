package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runFixPath syncs the active PATH binary from the deployed binary.
func runFixPath() {
	fmt.Println()
	fmt.Printf(constants.DoctorFixBannerFmt, constants.Version)
	fmt.Println(constants.DoctorBannerRule)

	absActive, activeVersion := resolveActiveBinary()
	if len(absActive) == 0 {
		return
	}

	absDeployed, deployedVersion := resolveDeployedForSync()
	if len(absDeployed) == 0 {
		return
	}

	printFixPathInfo(absActive, activeVersion, absDeployed, deployedVersion)
	syncBinaries(absActive, activeVersion, absDeployed, deployedVersion)
}

// resolveActiveBinary finds and validates the active PATH binary.
func resolveActiveBinary() (string, string) {
	activePath, activeErr := exec.LookPath(constants.GitMapBin)
	if activeErr != nil {
		printIssue(constants.DoctorNotOnPath, constants.DoctorNoSync)
		printFix(constants.DoctorAddPathFix)

		return "", ""
	}

	absActive, absErr := filepath.Abs(activePath)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", activePath, absErr)
		absActive = activePath
	}

	return absActive, getBinaryVersion(absActive)
}

// resolveDeployedForSync finds the deployed binary for syncing.
func resolveDeployedForSync() (string, string) {
	deployedPath, deployedErr := resolveDeployedBinary()
	if deployedErr != nil {
		printIssue(constants.DoctorCannotResolve, deployedErr.Error())

		return "", ""
	}

	absDeployed, absErr := filepath.Abs(deployedPath)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", deployedPath, absErr)
		absDeployed = deployedPath
	}

	return absDeployed, getBinaryVersion(absDeployed)
}

// printFixPathInfo displays the active and deployed binary paths.
func printFixPathInfo(absActive, activeVersion, absDeployed, deployedVersion string) {
	fmt.Printf(constants.DoctorActivePathFmt, absActive, activeVersion)
	fmt.Printf(constants.DoctorDeployedFmt, absDeployed, deployedVersion)
}

// syncBinaries orchestrates the 3-layer sync strategy.
func syncBinaries(absActive, activeVersion, absDeployed, deployedVersion string) {
	if absActive == absDeployed {
		printOK(constants.DoctorAlreadySynced)

		return
	}

	if activeVersion == deployedVersion {
		printOK(constants.DoctorVersionsMatch, activeVersion)

		return
	}

	fmt.Println()
	fmt.Printf(constants.DoctorSyncingFmt, absDeployed, absActive)
	attemptSync(absDeployed, absActive, deployedVersion)
}

// attemptSync tries copy, rename fallback, and kill strategies.
func attemptSync(absDeployed, absActive, deployedVersion string) {
	if tryCopyWithRetry(absDeployed, absActive, 20, 500*timeMillisecond) {
		verifySync(absActive, deployedVersion)

		return
	}

	if tryRenameFallback(absDeployed, absActive) {
		verifySync(absActive, deployedVersion)

		return
	}

	if tryKillAndCopy(absDeployed, absActive) {
		verifySync(absActive, deployedVersion)

		return
	}

	printSyncFailure(absDeployed, absActive)
}

// printSyncFailure reports that all sync strategies failed.
func printSyncFailure(absDeployed, absActive string) {
	fmt.Println()
	printIssue(constants.DoctorSyncFailTitle, constants.DoctorSyncFailDetail)
	printFix(constants.DoctorSyncFailFix1)
	printFix(fmt.Sprintf(constants.DoctorSyncFailFix2Fmt, absDeployed, absActive))
}

// resolveDeployedBinary finds the deployed binary path from powershell.json.
func resolveDeployedBinary() (string, error) {
	if len(constants.RepoPath) == 0 {
		return "", fmt.Errorf(constants.DoctorResolveNoRepo)
	}

	configPath := filepath.Join(constants.RepoPath, constants.GitMapSubdir, constants.PowershellConfigFile)
	data, err := os.ReadFile(configPath)
	if err != nil {
		return "", fmt.Errorf(constants.DoctorResolveNoRead, err)
	}

	return resolveDeployedPath(data)
}

// resolveDeployedPath extracts and validates the deployed path from config data.
func resolveDeployedPath(data []byte) (string, error) {
	deployPath := extractJSONString(data, constants.JSONKeyDeployPath)
	if len(deployPath) == 0 {
		return "", fmt.Errorf(constants.DoctorResolveNoDeploy)
	}

	binaryName := extractJSONString(data, constants.JSONKeyBinaryName)
	if len(binaryName) == 0 {
		binaryName = constants.DoctorDefaultBinary
	}

	deployed := filepath.Join(deployPath, constants.GitMapCliSubdir, binaryName)
	if _, err := os.Stat(deployed); err != nil {
		return "", fmt.Errorf(constants.DoctorResolveNotFound, deployed)
	}

	return deployed, nil
}

// verifySync checks that the synced binary reports the expected version.
func verifySync(path, expectedVersion string) {
	fmt.Println()
	actualVersion := getBinaryVersion(path)
	if actualVersion == expectedVersion {
		printOK(constants.DoctorOKPathFmt, actualVersion)

		return
	}

	printWarn(fmt.Sprintf(constants.DoctorWarnSyncFmt, actualVersion, expectedVersion))
}
