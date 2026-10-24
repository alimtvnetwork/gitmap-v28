// Package cmdinstall — installer_types.go provides OS targeting and sub-installer chaining data structures.
package cmdinstall

import (
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdinstall/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// InstallerOSType is an alias for the canonical enum in the package's
// enums subpackage, kept so existing callers keep compiling.
type InstallerOSType = enums.InstallerOSType

// Installer OS constants, aliased from the enums subpackage.
const (
	InstallerOSWindows = enums.InstallerOSWindows
	InstallerOSUnix    = enums.InstallerOSUnix
	InstallerOSUbuntu  = enums.InstallerOSUbuntu
	InstallerOSCentOS  = enums.InstallerOSCentOS
	InstallerOSAll     = enums.InstallerOSAll
)

// ChainedSubInstaller represents a downstream installer dependency executed sequentially.
type ChainedSubInstaller struct {
	SubInstallerId string                `json:"subInstallerId"`
	ParentId       string                `json:"parentId"`
	Name           string                `json:"name"`
	TargetOS       enums.InstallerOSType `json:"targetOS"`
	Command        string                `json:"command"`
	SequenceOrder  int                   `json:"sequenceOrder"`
	IsRequired     bool                  `json:"isRequired"`
}

// ModularInstallerRecipe defines an installer package with target OS filtering and child sub-installers.
type ModularInstallerRecipe struct {
	RecipeId      string                `json:"recipeId"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	TargetOS      enums.InstallerOSType `json:"targetOS"`
	MainCommand   string                `json:"mainCommand"`
	SubInstallers []ChainedSubInstaller `json:"subInstallers"`
}

// MatchesTargetOS verifies if the recipe or sub-installer is applicable to the current host OS.
func MatchesTargetOS(target enums.InstallerOSType, currentOS string) bool {
	lowTarget := strings.ToLower(target.String())
	lowCurrent := strings.ToLower(currentOS)

	if lowTarget == "" || lowTarget == enums.InstallerOSAll.String() {
		return true
	}

	if lowTarget == enums.InstallerOSWindows.String() && lowCurrent == constants.PlatformWindows {
		return true
	}

	if lowTarget == enums.InstallerOSUnix.String() && (lowCurrent == constants.PlatformLinux || lowCurrent == constants.PlatformDarwin) {
		return true
	}

	if lowTarget == enums.InstallerOSUbuntu.String() && lowCurrent == constants.PlatformLinux {
		return true
	}

	if lowTarget == enums.InstallerOSCentOS.String() && lowCurrent == constants.PlatformLinux {
		return true
	}

	return lowTarget == lowCurrent
}

// OrderSubInstallers returns sub-installers sorted deterministically by SequenceOrder.
func OrderSubInstallers(chain []ChainedSubInstaller) []ChainedSubInstaller {
	sorted := make([]ChainedSubInstaller, len(chain))
	copy(sorted, chain)

	sort.SliceStable(sorted, func(i, j int) bool {
		return sorted[i].SequenceOrder < sorted[j].SequenceOrder
	})

	return sorted
}

// ValidateSubInstallerChain ensures unique IDs and non-empty commands.
func ValidateSubInstallerChain(chain []ChainedSubInstaller) error {
	seen := make(map[string]bool)
	for _, sub := range chain {
		if sub.SubInstallerId == "" {
			return apperror.NewValidationError("subInstallerId cannot be empty")
		}
		if seen[sub.SubInstallerId] {
			return apperror.NewValidationError("duplicate subInstallerId: " + sub.SubInstallerId)
		}
		seen[sub.SubInstallerId] = true

		if strings.TrimSpace(sub.Command) == "" {
			return apperror.NewValidationError("subInstaller command cannot be empty for " + sub.SubInstallerId)
		}
	}

	return nil
}
