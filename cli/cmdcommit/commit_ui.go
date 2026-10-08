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

// execGitPaddedFiltered runs a git command, filters noise lines, and indents lines by 4 spaces.
func execGitPaddedFiltered(gitArgs ...string) error {
	cmd := exec.Command("git", gitArgs...)
	pipe, err := cmd.StdoutPipe()
	if err != nil {
		return cmd.Run()
	}

	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return err
	}

	scanPaddedFiltered(pipe)

	return cmd.Wait()
}

func scanPaddedFiltered(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := scanner.Text()
		if isGitNoiseLine(line) {
			continue
		}

		fmt.Printf("    %s\n", line)
	}
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
