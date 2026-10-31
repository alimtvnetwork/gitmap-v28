package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func RunPowerShell(args []string) error {
	psPath, hasPS := findPowerShellExecutable()
	if !hasPS {
		printPowerShellMissingAdvice()

		return apperror.NewSimple("powershell not found", "E_POWERSHELL_NOT_FOUND")
	}

	cleanArgs := stripForceFlag(args)
	if len(cleanArgs) == 0 {
		return executeInteractivePowerShell(psPath)
	}

	if err := inspectPowerShellPreflight(cleanArgs); err != nil {
		return err
	}

	return executePowerShellWithInterceptor(psPath, cleanArgs)
}

func executeInteractivePowerShell(psPath string) error {
	cmd := exec.Command(psPath)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return apperror.WrapSimple(err, "execute powershell")
	}

	return nil
}

func executePowerShellProcess(psPath string, args []string) error {
	return executePowerShellWithInterceptor(psPath, args)
}

func findPowerShellExecutable() (string, bool) {
	if p, err := exec.LookPath("pwsh"); err == nil && p != "" {
		return p, true
	}

	if !constants.IsWindowsPlatform() {
		return "", false
	}

	p, err := exec.LookPath("powershell")

	return p, err == nil && p != ""
}

func buildPowerShellCmd(psPath string, args []string) *exec.Cmd {
	if len(args) == 0 {
		return exec.Command(psPath)
	}

	cleanArgs := stripForceFlag(args)
	if len(cleanArgs) == 0 {
		return exec.Command(psPath)
	}

	return assemblePowerShellCmd(psPath, cleanArgs)
}

func assemblePowerShellCmd(psPath string, cleanArgs []string) *exec.Cmd {
	isWin := constants.IsWindowsPlatform()
	scriptText := extractPowerShellScript(cleanArgs)

	if isWin {
		return exec.Command(psPath, "-NoProfile", "-Command", scriptText)
	}

	return exec.Command(psPath, "-Command", scriptText)
}

func extractPowerShellScript(cleanArgs []string) string {
	if cleanArgs[0] == "-Command" || cleanArgs[0] == "-c" {
		return strings.Join(cleanArgs[1:], " ")
	}

	return strings.Join(cleanArgs, " ")
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
