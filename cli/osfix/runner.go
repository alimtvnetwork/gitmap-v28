package osfix

import (
	"fmt"
	"os"
	"os/exec"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunFixByName finds and executes a registered OS fix.
func RunFixByName(name string) FixBoolResult {
	res := GetFix(name)
	if res.IsFailure() {
		return result.Fail[bool](res.AppError())
	}
	fmt.Printf("▶ Running registered fix: %s\n", res.Value.Name)
	return RunFixCommand(res.Value.Command)
}

// RunFixCommand executes a fix command line string directly.
func RunFixCommand(command string) FixBoolResult {
	cmd := buildShellCommand(command)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	if err := cmd.Run(); err != nil {
		msg := fmt.Sprintf("execution failed for fix command: %s", command)
		return result.Fail[bool](apperror.WrapWithDetails(
			err, "RunFixCommand", "E_FIX_EXEC_FAIL", msg,
			"osfix", apperror.ErrorTypeExecution, apperror.SeverityError,
			map[string]any{"command": command},
		))
	}
	return result.Ok(true)
}

func buildShellCommand(command string) *exec.Cmd {
	if runtime.GOOS == constants.OSWindows {
		return buildWindowsShellCmd(command)
	}
	return exec.Command("sh", "-c", command)
}

func buildWindowsShellCmd(command string) *exec.Cmd {
	if strings.Contains(command, "|") || strings.Contains(command, ";") {
		return exec.Command("powershell.exe", "-NoProfile", "-Command", command)
	}
	return exec.Command("cmd.exe", "/c", command)
}
