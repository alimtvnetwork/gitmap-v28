package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// runHasAnyUpdates checks if the current repo is behind remote.
func runHasAnyUpdates(args []string) {
	checkHelp("has-any-updates", args)

	if !isInsideGitRepo() {
		fmt.Fprint(os.Stderr, constants.ErrHAUNotRepo)
		os.Exit(1)
	}

	fmt.Fprint(os.Stderr, constants.MsgHAUChecking)
	fetchRemote()

	ahead, behind, ok := hauAheadBehind()
	if !ok {
		fmt.Fprint(os.Stdout, constants.MsgHAUNoUpstream)

		return
	}

	hauPrintResult(ahead, behind)
}

// isInsideGitRepo returns true if cwd is inside a git work tree.
func isInsideGitRepo() bool {
	cmd := exec.Command(constants.GitBin, constants.GitRevParse, "--is-inside-work-tree")
	out, err := cmd.Output()

	return err == nil && strings.TrimSpace(string(out)) == "true"
}

// fetchRemote runs git fetch to update remote refs.
func fetchRemote() {
	cmd := exec.Command(constants.GitBin, "fetch")
	cmd.Stdout = nil
	cmd.Stderr = nil

	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrHAUFetchFailed, err)
	}
}

// hauAheadBehind returns ahead/behind counts relative to upstream.
func hauAheadBehind() (int, int, bool) {
	cmd := exec.Command(constants.GitBin, "rev-list", "--left-right", "--count", "HEAD...@{upstream}")
	out, err := cmd.Output()
	if err != nil {
		return 0, 0, false
	}

	parts := strings.Fields(strings.TrimSpace(string(out)))
	if len(parts) != 2 {
		return 0, 0, false
	}

	ahead, _ := strconv.Atoi(parts[0])
	behind, _ := strconv.Atoi(parts[1])

	return ahead, behind, true
}

// hauPrintResult prints the appropriate message based on ahead/behind.
func hauPrintResult(ahead, behind int) {
	switch {
	case behind > 0 && ahead > 0:
		fmt.Printf(constants.MsgHAUDiverged, ahead, behind)
	case behind > 0:
		fmt.Printf(constants.MsgHAUYes, behind)
	case ahead > 0:
		fmt.Printf(constants.MsgHAUAhead, ahead)
	default:
		fmt.Print(constants.MsgHAUNo)
	}
}
