package cmdupdate

import (
	"fmt"
	"net/http"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// resolveTargetSlug returns the repo slug to install from.
func resolveTargetSlug() (string, string, error) {
	if hasFlag(constants.FlagNoProbe) {
		fmt.Printf(constants.MsgUpdateProbeSkipped, constants.UpdateCurrentRepoSlug)

		return constants.UpdateCurrentRepoSlug, constants.UpdateProbeSourceMain, nil
	}

	return resolveLatestRepoSlug(newProbeClient())
}

// installerURLFor builds the raw.githubusercontent installer URL for slug.
func installerURLFor(slug string) string {
	name := constants.UpdateInstallerNameBash
	if hasWindowsRuntime() {
		name = constants.UpdateInstallerNamePwsh
	}

	return fmt.Sprintf(constants.UpdateRawInstallerTmpl,
		constants.UpdateRepoOwner, slug, name)
}

// downloadRemoteInstaller fetches url into a platform-appropriate temp file.
func downloadRemoteInstaller(url string) (string, error) {
	resp, err := http.Get(url) //nolint:gosec
	if err != nil {
		return "", err
	}

	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP %d from %s", resp.StatusCode, url)
	}

	return writeInstallerTempFile(resp.Body)
}
