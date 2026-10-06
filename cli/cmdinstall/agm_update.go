package cmdinstall

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var (
	agmUpdateRemote string
	agmUpdateSSH    bool
	agmUpdateExcept string
)

var agmUpdateCmd = &cobra.Command{
	Use:     "update [flags] [ssh]",
	Aliases: []string{"up", "u"},
	Short:   "Update Antigravity Manager GUI / Tools to latest",
	RunE:    runAgmUpdateCmd,
}

var agmUpdateAllCmd = &cobra.Command{
	Use:     "update-all [flags] [ssh]",
	Aliases: []string{"update-all-nodes", "updateallnodes", "update-nodes", "updateall", "up-all"},
	Short:   "Update Antigravity Manager across all SSH nodes or local",
	RunE:    runAgmUpdateAllCmd,
}

func init() {
	bindAgmUpdateFlags()
	AgmCmd.AddCommand(agmUpdateCmd)
	AgmCmd.AddCommand(agmUpdateAllCmd)
}

func bindAgmUpdateFlags() {
	bindCommonUpdateFlags(agmUpdateCmd)
	bindCommonUpdateFlags(agmUpdateAllCmd)
}

func bindCommonUpdateFlags(cmd *cobra.Command) {
	cmd.Flags().BoolVarP(&agmInstallDryRun, "dry-run", "n", false, "Simulate update without downloading")
	cmd.Flags().BoolVarP(&agmInstallYes, "yes", "y", false, "Automatic yes to prompts")
	cmd.Flags().BoolVarP(&agmInstallVerbose, "verbose", "v", false, "Enable verbose output")
	cmd.Flags().StringVar(&agmInstallVersion, "version", "", "Specific release version to update to (e.g. 4.7.1)")
	cmd.Flags().StringVarP(&agmUpdateRemote, "remote", "r", "", "Remote node alias or IP to update")
	cmd.Flags().BoolVarP(&agmUpdateSSH, "ssh", "s", false, "Update remote fleet machines via SSH")
	cmd.Flags().StringVarP(&agmUpdateExcept, "except", "e", "", "Exclude machines by ID, alias, or IP")
	cmd.Flags().StringVar(&agmUpdateExcept, "excep", "", "Exclude machines by ID, alias, or IP")
}

func runAgmUpdateCmd(cmd *cobra.Command, args []string) error {
	hasSSH := agmUpdateSSH || hasSSHArg(args)
	hasAll := hasAllArg(args)
	if (hasSSH || agmUpdateRemote != "" || hasAll) && RemoteAgmUpdateFleetFn != nil {
		target := resolveAgmUpdateTarget(agmUpdateRemote, hasAll)
		return RemoteAgmUpdateFleetFn(target, agmUpdateExcept)
	}
	opts := buildAgmInstallOptions()
	for _, a := range args {
		if !strings.HasPrefix(a, "-") && a != "ssh" && a != "all" && opts.Version == "" {
			opts.Version = strings.TrimPrefix(a, "v")
		}
	}
	return runAgmUpdateWithQuietOutput(opts)
}

func resolveAgmUpdateTarget(target string, hasAll bool) string {
	if target == "" || hasAll {
		return "all-nodes"
	}

	return target
}

func runAgmUpdateAllCmd(cmd *cobra.Command, args []string) error {
	if RemoteAgmUpdateFleetFn != nil {
		return RemoteAgmUpdateFleetFn("all-nodes", agmUpdateExcept)
	}
	opts := buildAgmInstallOptions()
	return runAgmUpdateWithQuietOutput(opts)
}

func runAgmUpdateWithQuietOutput(opts installOptions) error {
	if opts.DryRun {
		fmt.Printf("  [dry-run] Would run: %s (version: %s)\n", resolveAgManagerCommandForOS(), opts.Version)
		return nil
	}
	if runtime.GOOS == "windows" {
		return runUpdateAgManagerWithOpts(opts)
	}

	return runAgmUpdateLinuxQuiet(opts)
}

func runAgmUpdateLinuxQuiet(opts installOptions) error {
	version := opts.Version
	if version == "" {
		version = resolveLatestAgManagerReleaseVersion()
	}
	installCmd := constants.AgManagerUnixInstallCmd
	if version != "" {
		clean := strings.TrimPrefix(version, "v")
		installCmd = fmt.Sprintf(`curl -fsSL https://raw.githubusercontent.com/alimtvnetwork/Antigravity-Manager/main/install.sh | bash -s -- --version '%s'`, clean)
	}
	cmd := exec.Command("bash", "-c", installCmd)
	out, err := cmd.CombinedOutput()
	if err != nil {
		printUpdateFailureDetails(err, out)
		return apperror.WrapSimple(err, "Antigravity Manager update failed")
	}

	recordAgManagerInstalled(resolveInstalledVerName(version))
	verLabel := formatAgManagerVerLabel(version)
	fmt.Printf("%s✓%s Antigravity Manager%s updated successfully.\n", constants.ColorGreen, constants.ColorReset, verLabel)
	return nil
}

func printUpdateFailureDetails(err error, out []byte) {
	fmt.Fprintf(os.Stderr, "  %s✗%s Antigravity Manager update failed: %v\n", constants.ColorRed, constants.ColorReset, err)
	if len(out) > 0 {
		fmt.Fprintf(os.Stderr, "    Output:\n%s\n", strings.TrimSpace(string(out)))
	}
	trace := apperror.CaptureStackTrace(2)
	if trace != "" {
		fmt.Fprintf(os.Stderr, "    Stack trace:\n%s\n", trace)
	}
	verifyAgManagerOnFailure()
}

// PromptAgmBatchFailureResolution prompts the user when multiple failures occur during updates.
func PromptAgmBatchFailureResolution(failedTargets []string, retryFunc func(string) error) error {
	hasMultipleFailures := len(failedTargets) > 1
	if !hasMultipleFailures || !isInteractiveStdin() {
		return nil
	}
	fmt.Printf("\n  %s[?]%s %s%d targets failed during update. Do you want to resolve these failures?%s [Y/n]: ",
		constants.ColorCyan, constants.ColorReset,
		constants.ColorBold, len(failedTargets), constants.ColorReset)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil && len(line) == 0 {
		return nil
	}
	ans := strings.ToLower(strings.TrimSpace(line))
	if ans == "n" || ans == "no" || ans == "q" {
		return nil
	}

	fmt.Printf("\n  %s→ Resolving failures across %d target(s)...%s\n\n",
		constants.ColorCyan, len(failedTargets), constants.ColorReset)
	for _, target := range failedTargets {
		if err := retryFunc(target); err != nil {
			fmt.Printf("  %s✗%s Resolution for %s failed: %v\n", constants.ColorRed, constants.ColorReset, target, err)
		} else {
			fmt.Printf("  %s✓%s Resolution for %s succeeded.\n", constants.ColorGreen, constants.ColorReset, target)
		}
	}
	return nil
}

func isInteractiveStdin() bool {
	stat, err := os.Stdin.Stat()
	if err != nil {
		return false
	}
	return (stat.Mode() & os.ModeCharDevice) != 0
}

func hasSSHArg(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "ssh" || low == "--ssh" || low == "-s" {
			return true
		}
	}
	return false
}

func hasAllArg(args []string) bool {
	for _, a := range args {
		low := strings.ToLower(a)
		if low == "all" || low == "all-nodes" || low == "update-all" {
			return true
		}
	}
	return false
}
