package cmdupdate

import (
	"fmt"

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
