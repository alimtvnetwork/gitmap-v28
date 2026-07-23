package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// runReleaseAliasPull runs `git pull --ff-only` in target. Aborts on failure.
func runReleaseAliasPull(target string) {
	fmt.Printf(constants.MsgRAPullingFmt, target)

	cmd := exec.Command(constants.GitBin, "pull", "--ff-only")
	cmd.Dir = target
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRAPullFailedFmt, target, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}
}

// autoStashIfDirty stashes the working tree when dirty and returns the label
// used; returns "" when the tree was already clean.
func autoStashIfDirty(target, alias, version string) string {
	if !isWorkingTreeDirty(target) {
		return ""
	}

	label := fmt.Sprintf(constants.RAStashMessageFmt,
		fmt.Sprintf("%s-%s-%d", alias, version, time.Now().Unix()))

	cmd := exec.Command(constants.GitBin, "stash", "push", "--include-untracked", "-m", label)
	cmd.Dir = target
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrRAStashFailedFmt, target, err)
		fmt.Fprintln(os.Stderr)
		os.Exit(1)
	}

	fmt.Printf(constants.MsgRAStashCreatedFmt, label)

	return label
}

// popAutoStash attempts to pop the named stash entry. Failures only warn.
func popAutoStash(target, label string) {
	idx := findStashIndex(target, label)
	if idx == "" {
		return
	}

	cmd := exec.Command(constants.GitBin, "stash", "pop", idx)
	cmd.Dir = target
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		fmt.Fprintln(os.Stderr, constants.WarnRAStashPopFailed)

		return
	}

	fmt.Printf(constants.MsgRAStashPoppedFmt, label)
}

// isWorkingTreeDirty reports whether `git status --porcelain` returns content.
func isWorkingTreeDirty(target string) bool {
	cmd := exec.Command(constants.GitBin, "status", "--porcelain")
	cmd.Dir = target

	out, err := cmd.Output()
	if err != nil {
		return false
	}

	return len(strings.TrimSpace(string(out))) > 0
}

// findStashIndex locates the stash@{N} reference for a given label, "" if not found.
func findStashIndex(target, label string) string {
	cmd := exec.Command(constants.GitBin, "stash", "list")
	cmd.Dir = target

	out, err := cmd.Output()
	if err != nil {
		return ""
	}

	for _, line := range strings.Split(string(out), "\n") {
		if strings.Contains(line, label) {
			parts := strings.SplitN(line, ":", 2)

			return parts[0]
		}
	}

	return ""
}
