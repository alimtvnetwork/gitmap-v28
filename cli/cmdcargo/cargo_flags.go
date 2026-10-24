package cmdcargo

import (
	"slices"

	cargoconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/constants"
)

func extractInstallFlag(args []string) ([]string, bool) {
	var cleaned []string
	hasInstall := false
	for _, a := range args {
		if isInstallFlag(a) {
			hasInstall = true

			continue
		}
		cleaned = append(cleaned, a)
	}

	return cleaned, hasInstall
}

func isInstallFlag(arg string) bool {
	return slices.Contains([]string{cargoconst.FlagLongInstall, cargoconst.FlagShortInstall}, arg)
}
