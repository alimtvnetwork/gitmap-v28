package cmdpy

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"slices"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
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
	if runErr == nil && exitCode != 0 {
		return buildPythonAbortError(exitCode)
	}

	if runErr != nil {
		var exitErr *exec.ExitError
		if errors.As(runErr, &exitErr) {
			return buildPythonAbortError(exitErr.ExitCode())
		}

		return apperror.WrapSimple(runErr, "python execution failed")
	}

	return nil
}

func buildPythonAbortError(code int) error {
	return apperror.NewWithDetails(
		"python",
		"E9001",
		fmt.Sprintf("python process exited with code %d", code),
		"cmdpy",
		apperror.ErrorTypeAbort,
		apperror.SeverityWarn,
		map[string]any{"reported": true},
	)
}

func warnAdHocAntiPatterns(args []string) {
	code := extractCodeOrScriptSnippet(args)
	if code == "" {
		return
	}

	checkFileReadTip(code)
	checkOsWalkTip(code)
	checkSleepTip(code)
	checkTranscriptTip(code)
	checkSocketCDPTip(code)
	checkAgySummariesTip(code)
	checkSqliteTaskTip(code)
}

func extractInlineCode(args []string) string {
	for idx, arg := range args {
		if (arg == "-c" || arg == "--command") && idx+1 < len(args) {
			return args[idx+1]
		}
	}

	return ""
}

func extractCodeOrScriptSnippet(args []string) string {
	inline := extractInlineCode(args)
	if inline != "" {
		return inline
	}

	for _, arg := range args {
		if strings.HasSuffix(strings.ToLower(arg), ".py") {
			if data, err := os.ReadFile(arg); err == nil {
				if len(data) > 8192 {
					return string(data[:8192])
				}
				return string(data)
			}
		}
	}

	return ""
}

func checkFileReadTip(code string) {
	hasOpen := strings.Contains(code, "open(")
	hasRead := strings.Contains(code, ".read") || strings.Contains(code, "readlines") || strings.Contains(code, "splitlines") || strings.Contains(code, "for line in")
	if !hasOpen || !hasRead {
		return
	}

	fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python file reading (`open(...)`). Use 'gitmap cat <filepath> [-L start:end] [-q query] [-l limit]'.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func checkOsWalkTip(code string) {
	if !strings.Contains(code, "os.walk") {
		return
	}

	fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python filesystem crawling (`os.walk(...)`). Use 'gitmap find \"<pattern>\" [dir]' or 'gitmap lf [dir]'.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func checkSleepTip(code string) {
	if !strings.Contains(code, "time.sleep") {
		return
	}

	fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python sleep (`time.sleep(...)`). Use Antigravity 'schedule', 'gitmap watch', or background tasks.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func checkTranscriptTip(code string) {
	hasTranscript := strings.Contains(code, "transcript.jsonl") || strings.Contains(code, "transcript_full.jsonl")
	if !hasTranscript {
		return
	}

	fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python transcript parsing. Use 'gitmap cat <log_path> -q <pattern>' or read compact transcript directly.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func checkSocketCDPTip(code string) {
	hasSocket := strings.Contains(code, "socket.socket") || strings.Contains(code, "send_frame") || strings.Contains(code, "recv_frame")
	hasDevTools := strings.Contains(code, "socket") && strings.Contains(code, "devtools")
	if !hasSocket && !hasDevTools {
		return
	}

	fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python socket DevTools scripts. Use native 'gitmap cdp probe', 'gitmap cdp eval <expr>', or 'gitmap cdp target'.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func checkAgySummariesTip(code string) {
	hasAgyDB := strings.Contains(code, "conversation_summaries")
	if !hasAgyDB {
		return
	}

	fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python queries on conversation_summaries.db. Use native 'gitmap agy summaries [N]' or 'gitmap db-read agy \"<sql>\"'.%s\n\n", constants.ColorYellow, constants.ColorReset)
}

func checkSqliteTaskTip(code string) {
	hasSqlite := strings.Contains(code, "sqlite3")
	if !hasSqlite {
		return
	}

	if strings.Contains(code, "agent-task.db") || strings.Contains(code, "subtask") {
		fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python SQLite queries on agent tasks. Use native 'gitmap agent task ls' or 'gitmap agent subtask ls'.%s\n\n", constants.ColorYellow, constants.ColorReset)
		return
	}

	if strings.Contains(code, "connect(") || strings.Contains(code, ".db") {
		fmt.Fprintf(os.Stderr, "\n  %s💡 GitMap Tip: Avoid ad-hoc Python SQLite scripts. Use native 'gitmap db-read <alias|path> \"<sql>\"' or 'gitmap db-find'.%s\n\n", constants.ColorYellow, constants.ColorReset)
	}
}

// RunPy executes the Python runner CLI command.
func RunPy(args []string) error {
	if isPyHelpRequest(args) {
		PrintPyHelp()
		return nil
	}
	warnAdHocAntiPatterns(args)
	binary, err := resolvePythonBinary()
	if err != nil {
		return err
	}
	cwd, _ := os.Getwd()
	exitCode, runErr := executePythonStreaming(binary, args, cwd)
	return handleExitCode(exitCode, runErr)
}
