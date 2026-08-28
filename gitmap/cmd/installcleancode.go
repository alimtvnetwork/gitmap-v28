package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// cleanCodeAliases is the closed set of CLI tokens that all dispatch to the
// same coding-guidelines installer. Kept in one place so install.go and the
// validator stay in sync.
//
//	gitmap install clean-code
//	gitmap install code-guide
//	gitmap i cg
//	gitmap i cc
var cleanCodeAliases = map[string]struct{}{
	"clean-code": {},
	"code-guide": {},
	"cg":         {},
	"cc":         {},
}

// isCleanCodeAlias reports whether the tool argument matches any of the four
// coding-guidelines installer aliases.
func isCleanCodeAlias(tool string) bool {
	_, ok := cleanCodeAliases[tool]

	return ok
}

// runInstallCleanCode pipes the published install.ps1 through PowerShell.
// On Windows it prefers `powershell` (pre-installed); falls back to `pwsh`
// (PowerShell 7+) on every platform.
func runInstallCleanCode() error {
	pwsh := resolvePowerShellBinary()
	if pwsh == "" {
		return apperror.New(constants.MsgCleanCodeNoPwsh, "E9000", nil)
	}

	if runtime.GOOS != "windows" {
		fmt.Print(constants.MsgCleanCodeNonWin)
	}

	fmt.Printf(constants.MsgCleanCodeRunning, constants.DefaultCleanCodeURL)

	// PowerShell equivalent of: irm <url> | iex
	script := fmt.Sprintf("irm %s | iex", constants.DefaultCleanCodeURL)
	cmd := exec.Command(pwsh, "-NoProfile", "-ExecutionPolicy", "Bypass", "-Command", script)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		return apperror.Wrap(err, constants.ErrCleanCodeFailed, nil)
	}

	fmt.Print(constants.MsgCleanCodeDone)
	return nil
}

// resolvePowerShellBinary returns the first available PowerShell on PATH,
// preferring Windows PowerShell on Windows and pwsh elsewhere.
func resolvePowerShellBinary() string {
	return resolvePowerShellBinaryWithLookPath(exec.LookPath)
}

// resolvePowerShellBinaryWithLookPath resolves PowerShell with a custom LookPath func.
func resolvePowerShellBinaryWithLookPath(lookPath func(string) (string, error)) string {
	if lookPath == nil {
		lookPath = exec.LookPath
	}
	candidates := []string{"pwsh", "powershell"}
	if runtime.GOOS == "windows" {
		candidates = []string{"powershell", "pwsh"}
	}
	for _, name := range candidates {
		if path, err := lookPath(name); err == nil {
			return path
		}
	}

	return ""
}
