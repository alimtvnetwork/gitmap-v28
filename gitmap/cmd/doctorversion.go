package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// checkVersionMismatch compares PATH vs deployed vs source versions.
func checkVersionMismatch() int {
	sourceVersion := fmt.Sprintf(constants.MsgVersionFmt[:len(constants.MsgVersionFmt)-1], constants.Version)
	activeVersion, activePath := getActiveVersion()
	deployedVersion, deployedPath := getDeployedVersion()
	issues := 0

	issues += checkActiveVsSource(activeVersion, sourceVersion)
	issues += checkDeployedVsSource(deployedVersion, sourceVersion)
	issues += checkActiveVsDeployed(activeVersion, deployedVersion, activePath, deployedPath)

	if issues == 0 {
		printOK(constants.DoctorSourceOKFmt, sourceVersion)
	}

	return issues
}

// getActiveVersion returns version and path of the active PATH binary.
func getActiveVersion() (string, string) {
	path, err := exec.LookPath(constants.GitMapBin)
	if err != nil {
		return "", ""
	}

	absPath, absErr := filepath.Abs(path)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", path, absErr)
		absPath = path
	}

	return getBinaryVersion(absPath), absPath
}

// getDeployedVersion returns version and path of the deployed binary.
func getDeployedVersion() (string, string) {
	if len(constants.RepoPath) == 0 {
		return "", ""
	}

	data, err := readPowershellJSON()
	if err != nil {
		return "", ""
	}

	return resolveDeployedVersionPath(data)
}

// resolveDeployedVersionPath extracts version and path from config data.
func resolveDeployedVersionPath(data []byte) (string, string) {
	dp := extractJSONString(data, constants.JSONKeyDeployPath)
	bn := extractJSONString(data, constants.JSONKeyBinaryName)
	if len(bn) == 0 {
		bn = constants.DoctorDefaultBinary
	}

	if len(dp) == 0 {
		return "", ""
	}

	deployedPath := filepath.Join(dp, constants.GitMapCliSubdir, bn)

	return getBinaryVersion(deployedPath), deployedPath
}

// checkActiveVsSource reports if PATH binary differs from source.
func checkActiveVsSource(activeVersion, sourceVersion string) int {
	if len(activeVersion) > 0 && activeVersion != sourceVersion {
		printIssue(constants.DoctorVersionMismatch,
			fmt.Sprintf(constants.DoctorVMismatchFmt, activeVersion, sourceVersion))
		printFix(constants.DoctorVMismatchFix)

		return 1
	}

	return 0
}

// checkDeployedVsSource reports if deployed binary differs from source.
func checkDeployedVsSource(deployedVersion, sourceVersion string) int {
	if len(deployedVersion) > 0 && deployedVersion != sourceVersion {
		printIssue(constants.DoctorDeployMismatch,
			fmt.Sprintf(constants.DoctorDMismatchFmt, deployedVersion, sourceVersion))
		printFix(constants.DoctorDMismatchFix)

		return 1
	}

	return 0
}

// checkActiveVsDeployed reports if PATH and deployed binaries differ.
func checkActiveVsDeployed(activeVersion, deployedVersion, activePath, deployedPath string) int {
	if len(activeVersion) == 0 || len(deployedVersion) == 0 {
		return 0
	}

	if activeVersion == deployedVersion {
		return 0
	}

	absActive, err1 := filepath.Abs(activePath)
	if err1 != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", activePath, err1)
		absActive = activePath
	}
	absDeployed, err2 := filepath.Abs(deployedPath)
	if err2 != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", deployedPath, err2)
		absDeployed = deployedPath
	}
	if absActive == absDeployed {
		return 0
	}

	printIssue(constants.DoctorBinariesDiffer,
		fmt.Sprintf(constants.DoctorBDifferFmt, absActive, activeVersion, absDeployed, deployedVersion))
	printFix(constants.DoctorBDifferFix)

	return 1
}
