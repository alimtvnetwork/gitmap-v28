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

	// 3. (Optional) set credential store explicitly to avoid conflicts.
	exec.Command("git", "config", "--global", "--unset-all", "credential.credentialStore").Run()

	fmt.Println(constants.ColorGreen + "✔ Git credential helper successfully updated to 'manager'." + constants.ColorReset)
	fmt.Println("  The 'cache' credential store issue on Windows is now resolved.")

	return nil
}
