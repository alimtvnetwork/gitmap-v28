package cmdcommit

import (
	"bufio"
	"fmt"
	"io"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

var gitNoisePatterns = []string{
	"will be replaced by CRLF",
	"will be replaced by LF",
	"warning: in the working copy of",
	"The file will have its original line endings",
}

// isGitNoiseLine checks whether a line contains git CRLF/LF warning noise.
func isGitNoiseLine(line string) bool {
	for _, pattern := range gitNoisePatterns {
		if strings.Contains(line, pattern) {
			return true
		}
	}

	return false
}

func printPaddedInfo(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorCyan, "INFO", constants.ColorReset, msg)
}

func printPaddedSuccess(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorGreen, "SUCCESS", constants.ColorReset, msg)
}

func printPaddedWarning(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorYellow, "WARN", constants.ColorReset, msg)
}

func printPaddedError(format string, a ...any) {
	msg := fmt.Sprintf(format, a...)
	fmt.Printf("  %s%s%s %s\n", constants.ColorRed, "ERROR", constants.ColorReset, msg)
}

func printCommandVersionFooter() {
	fmt.Println()
	fmt.Printf("  %sgitmap v%s%s\n", constants.ColorDim, constants.Version, constants.ColorReset)
}

// execGitPaddedFiltered runs a git command, filters noise lines, and indents lines by 4 spaces.
func execGitPaddedFiltered(gitArgs ...string) error {
	_, err := execGitPaddedFilteredCapture(gitArgs...)
	return err
}

// execGitPaddedFilteredCapture behaves like execGitPaddedFiltered but also
// returns the filtered output, so a caller can attach git's own explanation
// to the error it reports instead of a bare exit code.
func execGitPaddedFilteredCapture(gitArgs ...string) (string, error) {
	cmd := exec.Command("git", gitArgs...)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return "", cmd.Run()
	}

	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return "", err
	}

	captured := scanPaddedFilteredCollect(pipe)

	return captured, cmd.Wait()
}

// gitErrorWithOutput wraps err with the last lines of captured git output so
// the failure names git's actual reason (e.g. a rejected push or a hook
// failure) rather than only "exit status 1".
func gitErrorWithOutput(err error, captured string) error {
	if err == nil {
		return nil
	}
	var tail []string
	for _, line := range strings.Split(captured, "\n") {
		if trimmed := strings.TrimSpace(line); trimmed != "" {
			tail = append(tail, trimmed)
		}
	}
	if len(tail) == 0 {
		return err
	}
	if len(tail) > 3 {
		tail = tail[len(tail)-3:]
	}

	return fmt.Errorf("%w — git said: %s", err, strings.Join(tail, " | "))
}

func scanPaddedFilteredCollect(r io.Reader) string {
	scanner := bufio.NewScanner(r)
	var kept []string
	for scanner.Scan() {
		line := scanner.Text()
		if isGitNoiseLine(line) {
			continue
		}

		fmt.Printf("    %s\n", line)
		kept = append(kept, line)
	}

	return strings.Join(kept, "\n")
}

// renderCommitPushSummaryCard renders a boxed completion summary card in the terminal.
func renderCommitPushSummaryCard(branch, sha, message, status string) {
	fmt.Println()
	fmt.Printf("  %s┌── Commit & Push Summary ──────────────────────────────────────────┐%s\n", constants.ColorCyan, constants.ColorReset)
	fmt.Printf("  %s│%s  Branch  : %s%s%s\n", constants.ColorCyan, constants.ColorReset, constants.ColorWhite, branch, constants.ColorReset)
	fmt.Printf("  %s│%s  Commit  : %s%s%s\n", constants.ColorCyan, constants.ColorReset, constants.ColorYellow, sha, constants.ColorReset)
	fmt.Printf("  %s│%s  Message : %s\n", constants.ColorCyan, constants.ColorReset, message)
	fmt.Printf("  %s│%s  Status  : %s%s%s\n", constants.ColorCyan, constants.ColorReset, constants.ColorGreen, status, constants.ColorReset)
	fmt.Printf("  %s└─────────────────────────────────────────────────────────────────────┘%s\n", constants.ColorCyan, constants.ColorReset)
}

// newNonGitRepoAbortError creates an AppError that cleanly aborts without dumping a stack trace.
func newNonGitRepoAbortError(msg string) *apperror.AppError {
	return apperror.NewWithDetails(
		"commit",
		"E9001",
		msg,
		"cli",
		apperror.ErrorTypeAbort,
		apperror.SeverityWarn,
		map[string]any{"reported": true},
	)
}

// handleNonGitRepoCommitPush handles commit-push execution outside a git repository.
func handleNonGitRepoCommitPush() *apperror.AppError {
	printPaddedError("Not a git repository (or any of the parent directories).")
	printPaddedInfo("To commit across child repositories, run: gitmap commit all")

	return newNonGitRepoAbortError("not a git repository")
}
