package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func runPowerShell(args []string) error {
	psPath, hasPS := findPowerShellExecutable()
	if !hasPS {
		printPowerShellMissingAdvice()
		return apperror.NewSimple("powershell not found", "E_POWERSHELL_NOT_FOUND")
	}

	return executePowerShellProcess(psPath, args)
}

func findPowerShellExecutable() (string, bool) {
	if p, err := exec.LookPath("pwsh"); err == nil && p != "" {
		return p, true
	}
	if runtime.GOOS != constants.OSWindows {
		return "", false
	}
	p, err := exec.LookPath("powershell")
	return p, err == nil && p != ""
}

func executePowerShellProcess(psPath string, args []string) error {
	cmd := buildPowerShellCmd(psPath, args)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "execute powershell")
	}
	return nil
}

func buildPowerShellCmd(psPath string, args []string) *exec.Cmd {
	if len(args) == 0 {
		return exec.Command(psPath)
	}
	cleanArgs := stripForceFlag(args)
	if len(cleanArgs) == 0 {
		return exec.Command(psPath)
	}
	if cleanArgs[0] == "-Command" || cleanArgs[0] == "-c" {
		return exec.Command(psPath, cleanArgs...)
	}
	isWin := runtime.GOOS == constants.OSWindows
	if isWin {
		return exec.Command(psPath, "-NoProfile", "-Command", strings.Join(cleanArgs, " "))
	}
	return exec.Command(psPath, "-Command", strings.Join(cleanArgs, " "))
}

func stripForceFlag(args []string) []string {
	var clean []string
	for _, a := range args {
		if a != "--force-all" && a != "-f" && a != "force-all" {
			clean = append(clean, a)
		}
	}
	return clean
}

func printPowerShellMissingAdvice() {
	printPowerShellMissingHeader()
	printPowerShellInstallOptions()
}

func printPowerShellMissingHeader() {
	fmt.Println()
	fmt.Println("✗ PowerShell (pwsh) is not installed on this system.")
	fmt.Println()
	fmt.Println("Do you like to install the PowerShell using GitMap? (gitmap install powershell)")
	fmt.Println()
}

func printPowerShellInstallOptions() {
	fmt.Println("  Local Installation:")
	fmt.Println("    ● gitmap install powershell")
	fmt.Println("    ● winget install --id Microsoft.Powershell [Windows]")
	fmt.Println()
	fmt.Println("  Remote SSH / Unix Installation:")
	fmt.Println("    ● gitmap ssh exec <node-alias> \"sudo snap install powershell --classic\"")
	fmt.Println("    ● gitmap ssh exec <node-alias> \"sudo apt-get install -y powershell\"")
	fmt.Println()
}
