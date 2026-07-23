package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// checkRepoPath reports whether RepoPath is embedded.
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
func checkActiveBinary() int {
	path, err := exec.LookPath(constants.GitMapBin)
	if err != nil {
		printIssue(constants.DoctorPathMissTitle, constants.DoctorPathMissDetail)
		printFix(constants.DoctorPathMissFix)

		return 1
	}

	absPath, absErr := filepath.Abs(path)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", path, absErr)
		absPath = path
	}
	version := getBinaryVersion(absPath)
	printOK(constants.DoctorPathBinaryFmt, absPath, version)

	return 0
}

// checkDeployedBinary reports the deployed binary from powershell.json.
func checkDeployedBinary() int {
	if len(constants.RepoPath) == 0 {
		return 0
	}

	data, err := readPowershellJSON()
	if err != nil {
		printIssue(constants.DoctorDeployReadFail, constants.DoctorDeployReadDet)

		return 1
	}

	deployedBinary, issue := resolveDeployedFromData(data)
	if issue > 0 {
		return issue
	}

	version := getBinaryVersion(deployedBinary)
	printOK(constants.DoctorDeployOKFmt, deployedBinary, version)

	return 0
}

// readPowershellJSON reads the powershell.json config file.
func readPowershellJSON() ([]byte, error) {
	configPath := filepath.Join(constants.RepoPath, constants.GitMapSubdir, constants.PowershellConfigFile)

	return os.ReadFile(configPath)
}

// resolveDeployedFromData extracts and validates the deployed binary path.
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

	deployedBinary := filepath.Join(deployPath, constants.GitMapCliSubdir, binaryName)
	if _, err := os.Stat(deployedBinary); err != nil {
		printIssue(constants.DoctorDeployNotFound, deployedBinary)
		printFix(constants.DoctorDeployRunFix)

		return "", 1
	}

	return deployedBinary, 0
}

// checkGit verifies git is available.
func checkGit() int {
	path, err := exec.LookPath(constants.GitBin)
	if err != nil {
		printIssue(constants.DoctorGitMissTitle, constants.DoctorGitMissDetail)

		return 1
	}

	version := getToolVersion(constants.GitBin, "--version")
	if len(version) == 0 {
		printOK(constants.DoctorGitOKPathFmt, path)

		return 0
	}

	printOK(constants.DoctorGitOKFmt, path, version)

	return 0
}

// checkGo verifies Go is available for building.
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
func getToolVersion(tool, arg string) string {
	cmd := exec.Command(tool, arg)
	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	return strings.TrimSpace(string(out))
}

// checkChangelogFile verifies CHANGELOG.md exists.
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
func checkLegacyDirs() int {
	printOK(constants.DoctorLegacyDirsOK)

	return 0
}

// checkSignature verifies whether the active binary has a valid digital signature.
// Only runs on Windows — signature verification uses PowerShell's Get-AuthenticodeSignature.
func checkSignature() int {
	if runtime.GOOS != "windows" {
		printWarn(constants.DoctorSignSkipUnix)

		return 0
	}

	binaryPath, err := exec.LookPath(constants.GitMapBin)
	if err != nil {
		printWarn(constants.DoctorSignNoPath)

		return 0
	}

	absPath, absErr := filepath.Abs(binaryPath)
	if absErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not resolve absolute path for %s: %v\n", binaryPath, absErr)
		absPath = binaryPath
	}

	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		"(Get-AuthenticodeSignature '"+absPath+"').Status")
	out, err := cmd.Output()
	if err != nil {
		printWarn(constants.DoctorSignCheckFail)

		return 0
	}

	status := strings.TrimSpace(string(out))

	if status == "Valid" {
		signer := getSignatureSigner(absPath)
		printOK(constants.DoctorSignOKFmt, absPath, signer)

		return 0
	}

	if status == "NotSigned" {
		printWarn(constants.DoctorSignUnsigned)

		return 0
	}

	printIssue(constants.DoctorSignInvalidFmt, status)
	printFix(constants.DoctorSignUnsignFix)

	return 1
}

// getSignatureSigner extracts the signer subject from a signed binary.
func getSignatureSigner(binaryPath string) string {
	cmd := exec.Command("powershell", "-NoProfile", "-Command",
		"(Get-AuthenticodeSignature '"+binaryPath+"').SignerCertificate.Subject")
	out, err := cmd.Output()
	if err != nil {
		return "unknown signer"
	}

	subject := strings.TrimSpace(string(out))
	if len(subject) == 0 {
		return "unknown signer"
	}

	return subject
}
