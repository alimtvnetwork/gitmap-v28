// Package cmdzsh provides native ZSH and Oh-My-Zsh management and parity engine.
package cmdzsh

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// ZshThemeType represents supported ZSH prompt theme names.
type ZshThemeType string

const (
	ZshThemeRobbyrussell ZshThemeType = "robbyrussell"
	ZshThemeAgnoster     ZshThemeType = "agnoster"
	ZshThemeFletcherm    ZshThemeType = "fletcherm"
	ZshThemeAfMagic      ZshThemeType = "af-magic"
	ZshThemeAfowler      ZshThemeType = "afowler"
	ZshThemeBira         ZshThemeType = "bira"
	ZshThemeCandy        ZshThemeType = "candy"
	ZshThemeClean        ZshThemeType = "clean"
	ZshThemeCloud        ZshThemeType = "cloud"
	ZshThemeGnzh         ZshThemeType = "gnzh"
	ZshThemeHalfLife     ZshThemeType = "half-life"
)

// PackageManagerType identifies supported system package managers.
type PackageManagerType string

const (
	PkgMgrApt     PackageManagerType = "apt"
	PkgMgrDnf     PackageManagerType = "dnf"
	PkgMgrBrew    PackageManagerType = "brew"
	PkgMgrPacman  PackageManagerType = "pacman"
	PkgMgrUnknown PackageManagerType = "unknown"
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
	ZshStatusResult = result.Result[ZshStatus]

	// ZshStringResult encapsulates a string outcome.
	ZshStringResult = result.Result[string]

	// ZshStringsResult encapsulates a slice of strings outcome.
	ZshStringsResult = result.ResultSlice[string]

	// ZshBoolResult encapsulates a boolean outcome.
	ZshBoolResult = result.Result[bool]
)
