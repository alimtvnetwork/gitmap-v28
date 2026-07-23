package cmd

import (
	"bufio"
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// finalizePush handles the final phase: --no-push prints the manual
// command, --yes pushes immediately, otherwise prompts.
func finalizePush(sandbox, originURL string, opts historyOpts) {
	if opts.noPush {
		fmt.Fprintf(os.Stdout, constants.HistoryMsgManualPush, sandbox, originURL)
		return
	}
	if !opts.yes && !confirmHistoryPush(sandbox, originURL, opts) {
		fmt.Fprintf(os.Stderr, constants.HistoryMsgUserAborted, sandbox)
		fmt.Fprintf(os.Stdout, constants.HistoryMsgManualPush, sandbox, originURL)
		return
	}
	pushSandbox(sandbox, originURL, opts)
}

// confirmHistoryPush prints the verification banner and blocks on
// stdin until the user types the literal token "yes". Any other
// input (including empty, "y", "Y") aborts. Default is abort.
func confirmHistoryPush(sandbox, originURL string, opts historyOpts) bool {
	fmt.Fprintf(os.Stderr, constants.HistoryMsgConfirmBanner,
		opts.modeLabel, opts.pathCount, sandbox, originURL)
	fmt.Fprintf(os.Stderr, constants.HistoryMsgConfirmPush, originURL)
	reader := bufio.NewReader(os.Stdin)
	line, err := reader.ReadString('\n')
	if err != nil {
		return false
	}
	return strings.TrimSpace(line) == "yes"
}

// pushSandbox runs `git -C <sandbox> push --mirror --force-with-lease
// <origin>`. Exit 7 on failure.
func pushSandbox(sandbox, originURL string, opts historyOpts) {
	if !opts.quiet {
		fmt.Fprintf(os.Stderr, constants.HistoryMsgPhasePush, originURL)
	}
	cmd := exec.Command(constants.HistoryGitBin, "-C", sandbox, "push",
		constants.HistoryPushRefSpec, constants.HistoryForceWithLease, originURL)
	cmd.Stdout, cmd.Stderr = os.Stderr, os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrPush, err)
		os.Exit(constants.HistoryExitPushFailed)
	}
	fmt.Fprint(os.Stderr, constants.HistoryMsgPushOk)
}
