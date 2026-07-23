package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// gitStage runs git add for a file.
func gitStage(file string) {
	cmd := exec.Command("git", "add", file)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOGitStage, err)
	}
}

// gitCommitWithAuthor creates a commit with optional author override.
func gitCommitWithAuthor(title, description, authorName, authorEmail string) {
	msg := title + "\n\n" + description

	if authorName != "" || authorEmail != "" {
		author := resolveAuthorFlag(authorName, authorEmail)
		cmd := exec.Command("git", "commit", "-m", msg, "--author", author)
		if err := cmd.Run(); err != nil {
			fmt.Fprintf(os.Stderr, constants.ErrSEOGitCommit, err)
		}

		return
	}

	cmd := exec.Command("git", "commit", "-m", msg)
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOGitCommit, err)
	}
}

// resolveAuthorFlag builds the --author "Name <email>" string.
func resolveAuthorFlag(name, email string) string {
	if name == "" {
		out, gitErr := exec.Command("git", "config", "user.name").Output()
		if gitErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not read git user.name: %v\n", gitErr)
		}
		name = strings.TrimSpace(string(out))
	}

	if email == "" {
		out, gitErr := exec.Command("git", "config", "user.email").Output()
		if gitErr != nil {
			fmt.Fprintf(os.Stderr, "  ⚠ Could not read git user.email: %v\n", gitErr)
		}
		email = strings.TrimSpace(string(out))
	}

	return name + " <" + email + ">"
}

// gitPush pushes to the remote.
func gitPush() {
	cmd := exec.Command("git", "push")
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrSEOGitPush, err)
	}
}

// appendToFile appends text to a file for rotation mode.
func appendToFile(path, text string) {
	f, err := os.OpenFile(path, os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to open file %s for append: %v\n", path, err)

		return
	}
	defer f.Close()

	if _, writeErr := f.WriteString("\n" + text); writeErr != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not append to %s: %v\n", path, writeErr)
	}
}

// revertFile removes the appended text from the file.
func revertFile(path, text string) {
	data, err := os.ReadFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, "  Error: failed to read file %s for revert: %v\n", path, err)

		return
	}

	cleaned := strings.Replace(string(data), "\n"+text, "", 1)
	if err := os.WriteFile(path, []byte(cleaned), 0o644); err != nil {
		fmt.Fprintf(os.Stderr, "  ⚠ Could not revert file %s: %v\n", path, err)
	}
}

// printHeader outputs the commit plan header.
func printHeader(max, minSec, maxSec int) {
	if max > 0 {
		fmt.Printf(constants.MsgSEOHeader, max, minSec, maxSec)

		return
	}

	fmt.Printf(constants.MsgSEOHeaderUnlimited, minSec, maxSec)
}

// printCommitLine outputs a single commit progress line.
func printCommitLine(max, current, total int, title, file string) {
	if max > 0 {
		fmt.Printf(constants.MsgSEOCommit, current, max, title, file)

		return
	}

	fmt.Printf(constants.MsgSEOCommitOpen, current, title, file)
}

// printRotationLine outputs a rotation progress line.
func printRotationLine(max, current int, file string) {
	if max > 0 {
		fmt.Printf(constants.MsgSEORotation, current, max, file)

		return
	}

	fmt.Printf(constants.MsgSEORotationOpen, current, file)
}

// printDone outputs the final summary line.
func printDone(count int, elapsed time.Duration) {
	fmt.Printf(constants.MsgSEODone, count, formatDuration(elapsed))
}

// formatDuration formats a duration into a human-readable string.
func formatDuration(d time.Duration) string {
	h := int(d.Hours())
	m := int(d.Minutes()) % 60

	if h > 0 {
		return fmt.Sprintf("%dh %dm", h, m)
	}

	return fmt.Sprintf("%dm", m)
}

// shouldStop checks if the loop should terminate.
func shouldStop(stop <-chan bool, maxCommits, count int) bool {
	select {
	case <-stop:
		return true
	default:
	}

	if maxCommits > 0 && count >= maxCommits {
		return true
	}

	return false
}
