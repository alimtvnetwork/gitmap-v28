// Package osuser provides cross-platform OS user management, root provisioning,
// SSH credential configuration, and process termination.
package osuser

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// Default constants for user provisioning.
const (
	DefaultLinuxShell = "/bin/zsh"
	DefaultZshTheme   = "fletcherm"
	SudoersDirPath    = "/etc/sudoers.d"
	SudoersFilePath   = "/etc/sudoers"
)

// UserCreateOptions defines parameters for provisioning a new OS user.
type UserCreateOptions struct {
	Username       string
	Password       string
	HomeDir        string
	Shell          string
	Theme          string
	IsSudoer       bool
	IsConfigureZsh bool
	IsDryRun       bool
	IsVerbose      bool
}

// UserRemoveOptions defines parameters for deleting an OS user.
type UserRemoveOptions struct {
	Username       string
	IsRemoveHome   bool
	IsForceKill    bool
	IsCleanSudoers bool
	IsDryRun       bool
	IsVerbose      bool
}

// UserKillOptions defines parameters for terminating user processes.
type UserKillOptions struct {
	Username  string
	IsForce   bool
	IsDryRun  bool
	IsVerbose bool
}

// SSHKeyOptions defines parameters for installing an SSH public key.
type SSHKeyOptions struct {
	Username    string
	PublicKey   string
	KeyFilePath string
	TargetHome  string
	IsDryRun    bool
	IsVerbose   bool
}

// UserInfo encapsulates summary information for an OS user.
type UserInfo struct {
	Username string
	HomeDir  string
	Shell    string
	IsSudoer bool
	IsActive bool
}

type (
	// UserResult encapsulates a UserInfo outcome.
	UserResult = result.Result[UserInfo]

	// UserBoolResult encapsulates a boolean outcome.
	UserBoolResult = result.Result[bool]

	// UserStringResult encapsulates a string outcome.
	UserStringResult = result.Result[string]

	// UserStringsResult encapsulates a slice of strings outcome.
	UserStringsResult = result.ResultSlice[string]
)
