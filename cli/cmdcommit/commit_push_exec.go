package cmdcommit

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func resolveShortHeadSha() string {
	headSha, _ := execGitOutputCP("rev-parse", "--short", "HEAD")
	sha := strings.TrimSpace(headSha)
	if sha == "" {
		return "-"
	}

	return sha
}

func handleCleanWorkingTreeCP() *apperror.AppError {
	printPaddedInfo("Working tree clean, nothing to commit.")

	unpushed := countUnpushedCommitsCP()
	if unpushed == 0 {
		printPaddedSuccess("Everything is up to date.")

		return nil
	}

	return pushUnpushedCommitsCP(unpushed)
}

func pushUnpushedCommitsCP(unpushed int) *apperror.AppError {
	printPaddedInfo("Pushing %d unpushed commit(s) to remote...", unpushed)

	if errPush := performPushStep(); errPush != nil {
		return errPush
	}

	printPaddedSuccess("Pushed %d commit(s) to remote.", unpushed)

	return nil
}

// ExecGitInheritCP runs a git command with inherited stdio.
func ExecGitInheritCP(gitArgs ...string) error {
	cmd := exec.Command("git", gitArgs...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	return cmd.Run()
}

// execGitOutputCP runs a git command and captures its stdout.
func execGitOutputCP(gitArgs ...string) (string, error) {
	cmd := exec.Command("git", gitArgs...)
	out, err := cmd.Output()

	return string(out), err
}

// execGitPadded runs a git command and indents all output lines by 2 spaces.
func execGitPadded(gitArgs ...string) error {
	cmd := exec.Command("git", gitArgs...)
	stdoutPipe, err := cmd.StdoutPipe()
	if err != nil {
		return cmd.Run()
	}

	cmd.Stderr = cmd.Stdout
	if err := cmd.Start(); err != nil {
		return err
	}

	scanPaddedLines(stdoutPipe)

	return cmd.Wait()
}

func scanPaddedLines(r io.Reader) {
	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		fmt.Printf("  %s\n", scanner.Text())
	}
}
