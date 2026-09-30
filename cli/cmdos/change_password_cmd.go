// Package cmdos — change_password_cmd.go manages changing OS user account passwords across Windows, Linux, and macOS.
package cmdos

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"os/user"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// ChangePasswordOptions holds configuration for modifying OS account passwords.
type ChangePasswordOptions struct {
	Username string
	Password string
	IsYes    bool
	IsDryRun bool
}

// RunChangePasswordCLI dispatches the OS password modification across supported operating systems.
func RunChangePasswordCLI(args []string) error {
	opts, shouldExit, err := parseChangePasswordArgs(args)
	if shouldExit || err != nil {
		return err
	}
	if err := resolvePasswordTarget(&opts); err != nil {
		return err
	}
	if !opts.IsYes && !opts.IsDryRun {
		if !confirmPasswordChange(opts.Username) {
			fmt.Println("Canceled by user.")
			return nil
		}
	}
	if opts.Password == "" {
		opts.Password = promptPassword(fmt.Sprintf("Enter new password for '%s': ", opts.Username))
		if opts.Password == "" {
			return apperror.NewValidationError("password cannot be empty")
		}
	}
	if opts.IsDryRun {
		fmt.Printf("[dry-run] Would change password for user '%s' on %s\n", opts.Username, runtime.GOOS)
		return nil
	}
	return executePlatformPasswordChange(opts.Username, opts.Password)
}

func parseChangePasswordArgs(args []string) (ChangePasswordOptions, bool, error) {
	opts := ChangePasswordOptions{}
	var posArgs []string
	for i := 0; i < len(args); i++ {
		if isChangePasswordHelpArg(args[i]) {
			printChangePasswordUsage()
			return opts, true, nil
		}
		if nextIdx, isHandled := consumeChangePasswordFlag(args, i, &opts); isHandled {
			i = nextIdx
			continue
		}
		if !strings.HasPrefix(args[i], "-") {
			posArgs = append(posArgs, args[i])
		}
	}
	assignPositionalArgs(&opts, posArgs)
	return opts, false, nil
}

func isChangePasswordHelpArg(arg string) bool {
	return arg == "-h" || arg == "--help" || arg == "help"
}

func consumeChangePasswordFlag(args []string, i int, opts *ChangePasswordOptions) (int, bool) {
	switch args[i] {
	case "-y", "--yes":
		opts.IsYes = true
		return i, true
	case "--dry-run":
		opts.IsDryRun = true
		return i, true
	case "-u", "--user":
		if i+1 < len(args) {
			opts.Username = args[i+1]
			return i + 1, true
		}
	case "-p", "--password":
		if i+1 < len(args) {
			opts.Password = args[i+1]
			return i + 1, true
		}
	}
	return i, false
}

func assignPositionalArgs(opts *ChangePasswordOptions, pos []string) {
	if len(pos) == 1 {
		opts.Password = pos[0]
		return
	}
	if len(pos) >= 2 {
		opts.Username = pos[0]
		opts.Password = pos[1]
	}
}

func resolvePasswordTarget(opts *ChangePasswordOptions) error {
	if opts.Username != "" {
		return nil
	}
	currentUser, _ := user.Current()
	opts.Username = resolveDefaultUsername(currentUser)
	if opts.Username == "" {
		return apperror.NewValidationError("could not detect current OS user; specify username explicitly")
	}
	return nil
}

func confirmPasswordChange(targetUser string) bool {
	fmt.Printf("Proceed with changing OS password for user '%s'? (yes/no): ", targetUser)
	reader := bufio.NewReader(os.Stdin)
	input, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	ans := strings.ToLower(strings.TrimSpace(input))
	return ans == "y" || ans == "yes"
}

func executePlatformPasswordChange(username, password string) error {
	switch runtime.GOOS {
	case "windows":
		return applyWindowsPassword(username, password)
	case "linux":
		return applyLinuxPassword(username, password)
	case "darwin":
		return applyDarwinPassword(username, password)
	default:
		return apperror.NewExecutionError("unsupported OS for change-password: " + runtime.GOOS)
	}
}

func applyWindowsPassword(username, password string) error {
	cmd := exec.Command("net", "user", username, password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("net user failed: %s (%v)", strings.TrimSpace(string(output)), err))
	}
	fmt.Printf("✓ %sSuccessfully updated Windows OS password for user '%s'%s\n", constants.ColorGreen, username, constants.ColorReset)
	return nil
}

func applyLinuxPassword(username, password string) error {
	cmd := exec.Command("chpasswd")
	cmd.Stdin = strings.NewReader(fmt.Sprintf("%s:%s\n", username, password))
	output, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("chpasswd failed: %s (%v)", strings.TrimSpace(string(output)), err))
	}
	fmt.Printf("✓ %sSuccessfully updated Linux OS password for user '%s'%s\n", constants.ColorGreen, username, constants.ColorReset)
	return nil
}

func applyDarwinPassword(username, password string) error {
	cmd := exec.Command("dscl", ".", "-passwd", "/Users/"+username, password)
	output, err := cmd.CombinedOutput()
	if err != nil {
		return apperror.NewExecutionError(fmt.Sprintf("dscl passwd failed: %s (%v)", strings.TrimSpace(string(output)), err))
	}
	fmt.Printf("✓ %sSuccessfully updated macOS password for user '%s'%s\n", constants.ColorGreen, username, constants.ColorReset)
	return nil
}

func printChangePasswordUsage() {
	fmt.Printf(`Change OS User Account Password (gitmap os change-password)

Usage:
  gitmap os change-password <new-password> [flags]
  gitmap os change-password <user> <new-password> [flags]
  gitmap os passwd [user] [password] [flags]
  gitmap change-password [user] [password] [flags]

Flags:
  -u, --user <user>      Target OS username (defaults to current user)
  -p, --password <pass>  New password to set
  -y, --yes              Confirm password change immediately without prompting
      --dry-run          Simulate password change without executing OS modification
  -h, --help             Show this help message

Examples:
  gitmap os change-password MyNewPass123!             # Change current user password after confirmation
  gitmap os change-password                           # Prompt for confirmation and password for current user
  gitmap os change-password Administrator MyNewPass   # Update Windows Administrator password (admin mode)
  gitmap os passwd ubuntu Secret123 -y                # Non-interactive password change
  gitmap change-password                              # Top-level shortcut
`)
}
