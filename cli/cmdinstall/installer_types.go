// Package cmdinstall — installer_types.go provides OS targeting and sub-installer chaining data structures.
package cmdinstall

import (
	"sort"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// InstallerOSType defines supported operating systems for modular recipes.
type InstallerOSType string

const (
	InstallerOSWindows InstallerOSType = "windows"
	InstallerOSUnix    InstallerOSType = "unix"
	InstallerOSUbuntu  InstallerOSType = "ubuntu"
	InstallerOSCentOS  InstallerOSType = "centos"
	InstallerOSAll     InstallerOSType = "all"
)

// ChainedSubInstaller represents a downstream installer dependency executed sequentially.
type ChainedSubInstaller struct {
	SubInstallerId string          `json:"subInstallerId"`
	ParentId       string          `json:"parentId"`
	Name           string          `json:"name"`
	TargetOS       InstallerOSType `json:"targetOS"`
	Command        string          `json:"command"`
	SequenceOrder  int             `json:"sequenceOrder"`
	IsRequired     bool            `json:"isRequired"`
}

// ModularInstallerRecipe defines an installer package with target OS filtering and child sub-installers.
type ModularInstallerRecipe struct {
	RecipeId      string                `json:"recipeId"`
	Name          string                `json:"name"`
	Description   string                `json:"description"`
	TargetOS      InstallerOSType       `json:"targetOS"`
	MainCommand   string                `json:"mainCommand"`
	SubInstallers []ChainedSubInstaller `json:"subInstallers"`
}

// MatchesTargetOS verifies if the recipe or sub-installer is applicable to the current host OS.
func MatchesTargetOS(target InstallerOSType, currentOS string) bool {
	lowTarget := strings.ToLower(string(target))
	lowCurrent := strings.ToLower(currentOS)

	if lowTarget == "" || lowTarget == string(InstallerOSAll) {
		return true
	}

	if lowTarget == string(InstallerOSWindows) && lowCurrent == constants.PlatformWindows {
		return true
	}

	if lowTarget == string(InstallerOSUnix) && (lowCurrent == constants.PlatformLinux || lowCurrent == constants.PlatformDarwin) {
		return true
	}

	if lowTarget == string(InstallerOSUbuntu) && lowCurrent == constants.PlatformLinux {
		return true
	}

	if lowTarget == string(InstallerOSCentOS) && lowCurrent == constants.PlatformLinux {
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
