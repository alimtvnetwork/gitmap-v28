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

func resolveTargetCredentialStore(args []string) string {
	for _, arg := range args {
		if arg == "--dpapi" {
			return "dpapi"
		}

		if arg == "--wincredman" {
			return "wincredman"
		}
	}

	return gitutil.ResolveWindowsCredentialStore()
}

func unsetLegacyCredentialHelper() {
	_ = exec.Command("git", "config", "--global", "--unset-all", "credential.helper").Run()
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

func applyCredentialStoreType(storeType string) {
	_ = exec.Command("git", "config", "--global", "credential.credentialStore", storeType).Run()
}

func printFixCredentialSuccess(storeType string) {
	fmt.Printf(constants.ColorGreen+"✔ Git credential helper successfully updated to 'manager' (credentialStore: %s).\n"+constants.ColorReset, storeType)
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
	unsetLegacyCredentialHelper()
	if appErr := setCredentialManagerHelper(); appErr != nil {
		return appErr
	}

	storeType := resolveTargetCredentialStore(args)
	applyCredentialStoreType(storeType)
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
