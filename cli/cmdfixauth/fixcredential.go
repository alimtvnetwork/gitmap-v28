package cmdfixauth

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

func resolveTargetCredentialStore(args []string) constants.CredentialStoreType {
	for _, arg := range args {
		if arg == "--dpapi" {
			return constants.CredentialStoreDPAPI
		}

		if arg == "--wincredman" {
			return constants.CredentialStoreWinCredMan
		}
	}

	return gitutil.ResolveWindowsCredentialStore()
}

func unsetLegacyCredentialHelper() *apperror.AppError {
	cmd := exec.Command("git", "config", "--global", "--unset-all", "credential.helper")
	if err := cmd.Run(); err != nil && !isBenignUnsetExit(err) {
		return apperror.WrapSimple(err, "git config unset credential.helper")
	}

	return nil
}

func isBenignUnsetExit(err error) bool {
	if exitErr, ok := err.(*exec.ExitError); ok {
		return exitErr.ExitCode() == 5
	}

	return false
}

func setCredentialManagerHelper() *apperror.AppError {
	cmd := exec.Command("git", "config", "--global", "credential.helper", "manager")
	out, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf("Failed to set credential helper: %s\n%v\n", strings.TrimSpace(string(out)), err)

		return apperror.WrapSimple(err, "git config set manager")
	}

	return nil
}

func applyCredentialStoreType(storeType constants.CredentialStoreType) *apperror.AppError {
	cmd := exec.Command("git", "config", "--global", "credential.credentialStore", string(storeType))
	if out, err := cmd.CombinedOutput(); err != nil {
		return apperror.WrapSimple(fmt.Errorf("%w: %s", err, strings.TrimSpace(string(out))), "git config set credentialStore")
	}

	return nil
}

func printFixCredentialSuccess(storeType constants.CredentialStoreType) {
	fmt.Printf(constants.ColorGreen+"✔ Git credential helper successfully updated to 'manager' (credentialStore: %s).\n"+constants.ColorReset, string(storeType))
	fmt.Println("  The 'cache' and 'wincredman' credential store issues on Windows are now resolved.")
}

func isWindowsPlatform() bool {
	if runtime.GOOS != constants.PlatformWindows {
		fmt.Println("gitmap fix-credential is only required on Windows. Your OS does not have this issue.")

		return false
	}

	return true
}

func applyWindowsCredentialFix(args []string) *apperror.AppError {
	if appErr := unsetLegacyCredentialHelper(); appErr != nil {
		return appErr
	}

	if appErr := setCredentialManagerHelper(); appErr != nil {
		return appErr
	}

	storeType := resolveTargetCredentialStore(args)
	if appErr := applyCredentialStoreType(storeType); appErr != nil {
		return appErr
	}

	printFixCredentialSuccess(storeType)

	return nil
}

// RunFixCredential resolves Git credential store issues on Windows.
func RunFixCredential(args []string) *apperror.AppError {
	if !isWindowsPlatform() {
		return nil
	}

	fmt.Println("Executing Git credential store fix for Windows...")

	return applyWindowsCredentialFix(args)
}
