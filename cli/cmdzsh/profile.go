// Package cmdzsh provides user profile and workspace initialization logic.
package cmdzsh

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// SetupUserProfile provisions standard workspace directories and .ssh permissions.
func SetupUserProfile(opts ZshProfileOptions) *apperror.AppError {
	home := ResolveTargetHome(opts.TargetHome)
	if opts.IsDryRun {
		fmt.Printf("[dry-run] Would setup user workspace directories in %s%s", home, constants.NewLineUnix)

		return nil
	}

	appErr := createWorkspaceDirectories(home)
	if appErr != nil {
		return appErr
	}

	return setupSshAuthorizedKeys(home, opts.AuthorizedKeyFile)
}

func getStandardWorkspaceDirs(home string) []string {
	return []string{
		filepath.Join(home, "scripts"),
		filepath.Join(home, "gitlab"),
		filepath.Join(home, "github"),
		filepath.Join(home, ".ssh"),
	}
}

func createWorkspaceDirectories(home string) *apperror.AppError {
	for _, d := range getStandardWorkspaceDirs(home) {
		appErr := makeSecuredDirectory(d)
		if appErr != nil {
			return appErr
		}
	}

	return nil
}

func makeSecuredDirectory(path string) *apperror.AppError {
	appErr := EnsureDirectory(path, 0700)
	if appErr != nil {
		return appErr
	}

	_ = os.Chmod(path, 0700)

	return nil
}

func setupSshAuthorizedKeys(home, keyFile string) *apperror.AppError {
	sshDir := filepath.Join(home, ".ssh")
	authFile := filepath.Join(sshDir, "authorized_keys")
	appErr := touchAuthorizedKeysFile(authFile)
	if appErr != nil {
		return appErr
	}

	if keyFile != "" && HasFile(keyFile) {
		return copyAuthorizedKeyContent(keyFile, authFile)
	}

	return nil
}

func touchAuthorizedKeysFile(authFile string) *apperror.AppError {
	if HasFile(authFile) {
		_ = os.Chmod(authFile, 0600)

		return nil
	}

	return createEmptyAuthorizedKeys(authFile)
}

func createEmptyAuthorizedKeys(authFile string) *apperror.AppError {
	f, err := os.OpenFile(authFile, os.O_CREATE|os.O_WRONLY, 0600)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.createEmptyAuthorizedKeys")
	}

	defer f.Close()

	_ = os.Chmod(authFile, 0600)

	return nil
}

func copyAuthorizedKeyContent(srcFile, dstFile string) *apperror.AppError {
	data, err := os.ReadFile(srcFile)
	if err != nil {
		return apperror.WrapSimple(err, "cmdzsh.copyAuthorizedKeyContent.read")
	}

	writeErr := os.WriteFile(dstFile, data, 0600)
	if writeErr != nil {
		return apperror.WrapSimple(writeErr, "cmdzsh.copyAuthorizedKeyContent.write")
	}

	_ = os.Chmod(dstFile, 0600)

	return nil
}
