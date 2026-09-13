package osuser

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// AddUser creates a new OS user across Windows and Linux (backward compatibility).
func AddUser(username, password string) error {
	opts := buildAddUserOptions(username, password)
	appErr := CreateRootUser(opts)
	if appErr != nil {
		return appErr
	}
	return nil
}

func buildAddUserOptions(username, password string) UserCreateOptions {
	return UserCreateOptions{
		Username: username,
		Password: password,
	}
}

// RemoveUser deletes an OS user and their profile/home directory (backward compatibility).
func RemoveUser(username string) error {
	opts := buildRemoveUserOptions(username)
	appErr := RemoveEnhancedUser(opts)
	if appErr != nil {
		return appErr
	}
	return nil
}

func buildRemoveUserOptions(username string) UserRemoveOptions {
	return UserRemoveOptions{
		Username:       username,
		IsRemoveHome:   true,
		IsCleanSudoers: true,
	}
}

// CreateRoot provisions a root/administrative user with sudoers and optional ZSH.
func CreateRoot(opts UserCreateOptions) *apperror.AppError {
	return CreateRootUser(opts)
}

// Remove deletes an OS user with full cleanup of home directory and sudoers.
func Remove(opts UserRemoveOptions) *apperror.AppError {
	return RemoveEnhancedUser(opts)
}

// Kill terminates processes owned by the given user.
func Kill(opts UserKillOptions) *apperror.AppError {
	return KillUserProcesses(opts)
}

// InstallKey installs an SSH public key for the given user.
func InstallKey(opts SSHKeyOptions) *apperror.AppError {
	return InstallSSHKey(opts)
}
