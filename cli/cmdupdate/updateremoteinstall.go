package cmdupdate

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runUpdateRemoteInstall is the v5.52.0+ remote-installer flow.
func runUpdateRemoteInstall() bool {
	slug, source, err := resolveTargetSlug()
	if err != nil {
		return false
	}

	if hasFlag(constants.FlagProbeOnly) {
		fmt.Printf(constants.MsgUpdateProbeOnly, slug, source)

		return true
	}

	return startRemoteUpdateWorkflow(slug)
}

func startRemoteUpdateWorkflow(slug string) bool {
	currentVersion := constants.Version
	targetVersion := fetchRemoteTargetVersion(slug)
	url := installerURLFor(slug)
	fmt.Printf(constants.MsgUpdateRemoteFetch, url)

	return executeRemoteUpdateWorkflow(url, currentVersion, targetVersion)
}

func executeRemoteUpdateWorkflow(url, currentVersion, targetVersion string) bool {
	scriptPath, err := downloadRemoteInstaller(url)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrUpdateRemoteDownload, err)

		return false
	}

	defer os.Remove(scriptPath)
	announceRemoteUpdate(currentVersion, targetVersion, scriptPath)

	return runAndFinishRemoteUpdate(scriptPath, currentVersion, targetVersion, url)
}

func runAndFinishRemoteUpdate(scriptPath, currentVersion, targetVersion, url string) bool {
	if errRun := runRemoteInstaller(scriptPath); errRun != nil {
		handleRemoteInstallerError(errRun)

		return false
	}

	finishRemoteUpdate(currentVersion, targetVersion, url)

	return true
}

func announceRemoteUpdate(currentVersion, targetVersion, scriptPath string) {
	fmt.Printf(constants.MsgUpdateVersionCompare, currentVersion, targetVersion)
	fmt.Printf(constants.MsgUpdateRemoteRun, scriptPath)
}

func finishRemoteUpdate(currentVersion, targetVersion, url string) {
	fmt.Printf(constants.MsgUpdateSummaryDetail, currentVersion, targetVersion, url)
	printPostUpdateIdentity()
}
