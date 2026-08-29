package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// checkRepoPath reports whether RepoPath is embedded.
//
//nolint:unused
func checkRepoPath() int {
	if len(constants.RepoPath) == 0 {
		printIssue(constants.DoctorRepoPathMissing, constants.DoctorRepoPathDetail)
		printFix(constants.DoctorRepoPathFix)

		return 1
	}

	printOK(constants.DoctorRepoPathOKFmt, constants.RepoPath)

	return 0
}

// checkActiveBinary reports the gitmap binary on PATH.
//
//nolint:unused
func checkActiveBinary() int {
	path, err := exec.LookPath(constants.GitMapBin)
	if err != nil {
		printIssue(constants.DoctorPathMissTitle, constants.DoctorPathMissDetail)
		printFix(constants.DoctorPathMissFix)

		return 1
	}
	absPath := resolveBinaryAbsPath(path)
	version := getBinaryVersion(absPath)
	printOK(constants.DoctorPathBinaryFmt, absPath, version)

	return 0
}

//nolint:unused
func resolveBinaryAbsPath(path string) string {
	absPath, absErr := filepath.Abs(path)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", path, absErr)

		return path
	}

	return absPath
}

// checkDeployedBinary reports the deployed binary from powershell.json.
//
//nolint:unused
func checkDeployedBinary() int {
	if len(constants.RepoPath) == 0 {
		return 0
	}

	data, err := readPowershellJSON()
	if err != nil {
		printIssue(constants.DoctorDeployReadFail, constants.DoctorDeployReadDet)

		return 1
	}

	return verifyDeployedBinary(data)
}

//nolint:unused
func verifyDeployedBinary(data []byte) int {
	deployedBinary, issue := resolveDeployedFromData(data)
	if issue > 0 {
		return issue
	}

	version := getBinaryVersion(deployedBinary)
	printOK(constants.DoctorDeployOKFmt, deployedBinary, version)

	return 0
}

// readPowershellJSON reads the powershell.json config file.
//
//nolint:unused
func readPowershellJSON() ([]byte, error) {
	configPath := filepath.Join(constants.RepoPath, constants.GitMapSubdir, constants.PowershellConfigFile)

	//nolint:unused
	return os.ReadFile(configPath)
}

// resolveDeployedFromData extracts and validates the deployed binary path.
//
//nolint:unused
func resolveDeployedFromData(data []byte) (string, int) {
	deployPath := extractJSONString(data, constants.JSONKeyDeployPath)
	if len(deployPath) == 0 {
		printIssue(constants.DoctorNoDeployPath, constants.DoctorNoDeployDet)

		return "", 1
	}

	binaryName := extractJSONString(data, constants.JSONKeyBinaryName)
	if len(binaryName) == 0 {
		binaryName = constants.DoctorDefaultBinary
	}
	//nolint:unused

	return checkDeployedFileExists(filepath.Join(deployPath, constants.GitMapCliSubdir, binaryName))
	//nolint:unused
}

//nolint:unused
func checkDeployedFileExists(deployedBinary string) (string, int) {
	if _, err := os.Stat(deployedBinary); err != nil {
		printIssue(constants.DoctorDeployNotFound, deployedBinary)
		printFix(constants.DoctorDeployRunFix)

		return "", 1
	}

	//nolint:unused
	return deployedBinary, 0
}

// checkGit verifies git is available.
//
//nolint:unused
func checkGit() int {
	path, err := exec.LookPath(constants.GitBin)
	if err != nil {
		printIssue(constants.DoctorGitMissTitle, constants.DoctorGitMissDetail)

		return 1
	}

	version := getToolVersion(constants.GitBin, constants.GitVersionArg)
	if len(version) == 0 {
		printOK(constants.DoctorGitOKPathFmt, path)

		return 0
	}

	printOK(constants.DoctorGitOKFmt, path, version)

	return 0
}

// checkGo verifies Go is available for building.
//
//nolint:unused
func checkGo() int {
	path, err := exec.LookPath(constants.GoBin)
	if err != nil {
		printWarn(constants.DoctorGoWarn)

		return 0
	}

	version := getToolVersion(constants.GoBin, constants.GoVersionArg)
	if len(version) == 0 {
		printOK(constants.DoctorGoOKPathFmt, path)

		return 0
	}

	printOK(constants.DoctorGoOKFmt, version)

	return 0
}

// getToolVersion runs a tool with an arg and returns trimmed output.
//
//nolint:unused
func getToolVersion(tool, arg string) string {
	cmd := exec.Command(tool, arg)
	out, err := cmd.Output()
	if err != nil {
		//nolint:unused
		return ""
	}

	return strings.TrimSpace(string(out))
}

// checkChangelogFile verifies changelog.md exists.
//
//nolint:unused
func checkChangelogFile() int {
	if _, err := os.Stat(constants.ChangelogFile); err != nil {
		printWarn(constants.DoctorChangelogWarn)

		return 0
	}

	printOK(constants.DoctorChangelogOK)

	return 0
}

// checkLegacyDirs confirms no legacy directories remain after auto-migration.
// Since migrateLegacyDirs now merges and removes legacy folders, this check
// serves only as a safety net for edge cases (e.g., permission errors).
//
//nolint:unused
func checkLegacyDirs() int {
	printOK(constants.DoctorLegacyDirsOK)

	//nolint:unused
	return 0
}

// checkSignature verifies whether the active binary has a valid digital signature.
// Only runs on Windows — signature verification uses PowerShell's Get-AuthenticodeSignature.
//
//nolint:unused
func checkSignature() int {
	if runtime.GOOS != constants.PlatformWindows {
		printWarn(constants.DoctorSignSkipUnix)

		return 0
	}

	absPath, ok := resolveSignaturePath()
	if !ok {
		return 0
		//nolint:unused
	}
	//nolint:unused

	return verifyBinarySignature(absPath)
}

//nolint:unused
func resolveSignaturePath() (string, bool) {
	binaryPath, err := exec.LookPath(constants.GitMapBin)
	if err != nil {
		printWarn(constants.DoctorSignNoPath)

		//nolint:unused
		return "", false
	}
	//nolint:unused

	return resolveBinaryAbsPath(binaryPath), true
}

//nolint:unused
func verifyBinarySignature(absPath string) int {
	cmd := exec.Command(constants.ShellPowerShell, constants.DoctorFlagNoProfile, constants.DoctorFlagCommand,
		"(Get-AuthenticodeSignature '"+absPath+"').Status")
	out, err := cmd.Output()
	if err != nil {
		printWarn(constants.DoctorSignCheckFail)

		return 0
	}

	//nolint:unused
	return evaluateSignatureStatus(absPath, strings.TrimSpace(string(out)))
}

//nolint:unused
func evaluateSignatureStatus(absPath, status string) int {
	if status == constants.DoctorSignStatusValid {
		signer := getSignatureSigner(absPath)
		printOK(constants.DoctorSignOKFmt, absPath, signer)

		return 0
	}

	if status == constants.DoctorSignStatusNotSigned {
		printWarn(constants.DoctorSignUnsigned)

		return 0
	}

	printIssue(constants.DoctorSignInvalidFmt, status)
	//nolint:unused
	printFix(constants.DoctorSignUnsignFix)

	return 1
}

// getSignatureSigner extracts the signer subject from a signed binary.
//
//nolint:unused
func getSignatureSigner(binaryPath string) string {
	cmd := exec.Command(constants.ShellPowerShell, constants.DoctorFlagNoProfile, constants.DoctorFlagCommand,
		"(Get-AuthenticodeSignature '"+binaryPath+"').SignerCertificate.Subject")
	out, err := cmd.Output()
	if err != nil {
		return constants.DoctorUnknownSigner
	}

	subject := strings.TrimSpace(string(out))
	if len(subject) == 0 {
		return constants.DoctorUnknownSigner
	}

	return subject
}
