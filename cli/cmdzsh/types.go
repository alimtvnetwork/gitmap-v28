// Package cmdzsh provides native ZSH and Oh-My-Zsh management and parity engine.
package cmdzsh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdzsh/enums"

// ZshThemeType represents supported ZSH prompt theme names.
// ZshThemeType aliases the strongly-typed enum from the enums package.
type ZshThemeType = enums.ZshThemeType

const (
	ZshThemeRobbyrussell = enums.ZshThemeRobbyrussell
	ZshThemeAgnoster     = enums.ZshThemeAgnoster
	ZshThemeFletcherm    = enums.ZshThemeFletcherm
	ZshThemeAfMagic      = enums.ZshThemeAfMagic
	ZshThemeAfowler      = enums.ZshThemeAfowler
	ZshThemeBira         = enums.ZshThemeBira
	ZshThemeCandy        = enums.ZshThemeCandy
	ZshThemeClean        = enums.ZshThemeClean
	ZshThemeCloud        = enums.ZshThemeCloud
	ZshThemeGnzh         = enums.ZshThemeGnzh
	ZshThemeHalfLife     = enums.ZshThemeHalfLife
)

// PackageManagerType identifies supported system package managers.
// PackageManagerType aliases the strongly-typed enum from the enums package.
type PackageManagerType = enums.PackageManagerType

const (
	PkgMgrApt     = enums.PkgMgrApt
	PkgMgrDnf     = enums.PkgMgrDnf
	PkgMgrBrew    = enums.PkgMgrBrew
	PkgMgrPacman  = enums.PkgMgrPacman
	PkgMgrUnknown = enums.PkgMgrUnknown
)

// ZshOptions encapsulates configuration parameters for ZSH installation and suite setup.
type ZshOptions struct {
	Theme          string
	TargetHome     string
	TargetUser     string
	CustomZshrc    string
	AuthorizedKeys string
	IsAppendZshrc  bool
	IsInstallZsh   bool
	IsChangeShell  bool
	IsDryRun       bool
	IsVerbose      bool
	IsUnattended   bool
}

// ZshCleanOptions encapsulates parameters for cleaning and reinstalling Oh-My-Zsh.
type ZshCleanOptions struct {
	Theme       string
	TargetHome  string
	IsBackup    bool
	IsReinstall bool
	IsDryRun    bool
}

// ZshSwitchOptions encapsulates parameters for switching the default shell.
type ZshSwitchOptions struct {
	TargetUser string
	IsDryRun   bool
}

// ZshProfileOptions encapsulates parameters for setting up user workspace directories.
type ZshProfileOptions struct {
	TargetHome        string
	AuthorizedKeyFile string
	IsDryRun          bool
}

// ZshStatus captures inspection details of the local ZSH environment.
type ZshStatus struct {
	ZshVersion         string
	CurrentTheme       string
	DefaultShell       string
	ZshPath            string
	OhMyZshPath        string
	ActivePlugins      []string
	IsZshInstalled     bool
	IsOhMyZshInstalled bool
	IsDefaultShell     bool
}

type (
	// ZshStatusResult encapsulates a ZshStatus outcome.
	ZshStatusResult = result.Wrap[ZshStatus]

	// ZshStringResult encapsulates a string outcome.
	ZshStringResult = result.StringResult

	// ZshStringsResult encapsulates a slice of strings outcome.
	ZshStringsResult = result.StringSliceResult

	// ZshBoolResult encapsulates a boolean outcome.
	ZshBoolResult = result.BoolResult
)
