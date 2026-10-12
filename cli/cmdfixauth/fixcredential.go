package cmdfixauth

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// runFixCredential resolves the 'cache' credential store error on Windows
// by configuring the git credential helper to 'manager' (Git Credential Manager).
func RunFixCredential(args []string) *apperror.AppError {
	if runtime.GOOS != constants.PlatformWindows {
		fmt.Println("gitmap fix-credential is only required on Windows. Your OS does not have this issue.")
		return nil
	}

	fmt.Println("Executing Git credential store fix for Windows...")

	// 1. Unset the problematic 'cache' helper if it exists.
	exec.Command("git", "config", "--global", "--unset-all", "credential.helper").Run()

	// 2. Set it to 'manager' (Git Credential Manager)
	cmd := exec.Command("git", "config", "--global", "credential.helper", "manager")
	if out, err := cmd.CombinedOutput(); err != nil {
		fmt.Printf("Failed to set credential helper: %s\n%v\n", strings.TrimSpace(string(out)), err)
		return apperror.WrapSimple(err, "git config set manager")
	}

	// 3. Set credential store to 'wincred' to avoid headless/SSH wincredman persistence failures.
	storeType := "wincred"
	for _, arg := range args {
		if arg == "--dpapi" {
			storeType = "dpapi"
		}
	}
	exec.Command("git", "config", "--global", "credential.credentialStore", storeType).Run()

	fmt.Printf(constants.ColorGreen+"✔ Git credential helper successfully updated to 'manager' (credentialStore: %s).\n"+constants.ColorReset, storeType)
	fmt.Println("  The 'cache' and 'wincredman' credential store issues on Windows are now resolved.")


	return nil
}
