package cmdpy

import (
	"fmt"
	"os"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// PrintPyHelp renders help information for gitmap py / gitmap python.
func PrintPyHelp() {
	fmt.Println()
	fmt.Printf("  %s╔══════════════════════════════════════════════════════════════════╗%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s║             gitmap py / python - Python Execution Runner         ║%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s╚══════════════════════════════════════════════════════════════════╝%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Println()
	fmt.Println("  Usage:")
	fmt.Println("    gitmap py -c \"<code>\"             Execute inline Python code")
	fmt.Println("    gitmap py <script.py> [args...]   Execute Python script file")
	fmt.Println("    gitmap py [args...]               Pass arguments directly to Python")
	fmt.Println("    gitmap python [args...]           Canonical alias for gitmap py")
	fmt.Println()
}

var pyHelpFlagTokens = []string{"help", "--help", "-h"}

func isPyHelpRequest(args []string) bool {
	if len(args) == 0 {
		return false
	}
	arg := strings.ToLower(args[0])
	return slices.Contains(pyHelpFlagTokens, arg)
}

func handleExitCode(exitCode int, runErr error) error {
	if exitCode != 0 && runErr == nil {
		cliexit.HandleError(nil, exitCode)
	}
	if runErr != nil {
		return apperror.WrapSimple(runErr, "python execution failed")
	}
	return nil
}

func warnAdHocFileRead(args []string) {
	if !hasAdHocFileRead(args) {
		return
	}

	fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python file reading (`open(...)`). Use 'gitmap cat <filepath> -L <start>:<end>', 'gitmap cat <filepath> -q <query>', or native tool calls.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func hasAdHocFileRead(args []string) bool {
	for idx, arg := range args {
		if (arg == "-c" || arg == "--command") && idx+1 < len(args) {
			code := args[idx+1]
			if strings.Contains(code, "open(") && (strings.Contains(code, ".read") || strings.Contains(code, "for line in") || strings.Contains(code, "readlines")) {
				return true
			}
		}
	}

	return false
}

// RunPy executes the Python runner CLI command.
func RunPy(args []string) error {
	if isPyHelpRequest(args) {
		PrintPyHelp()
		return nil
	}
	warnAdHocFileRead(args)
	binary, err := resolvePythonBinary()
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	exitCode, runErr := executePythonStreaming(binary, args, cwd)
	return handleExitCode(exitCode, runErr)
}
