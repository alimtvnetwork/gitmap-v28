package cmd

// Entrypoints for `gitmap history-purge` (`hp`) and `gitmap history-pin`
// (`hpin`). Both commands wrap `git filter-repo` in a mirror-clone
// sandbox so the user's working repository is never rewritten in
// place. Spec: spec/04-generic-cli/16-history-rewrite.md.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
)

// historyMode picks which filter-repo behavior to wrap.
type historyMode int

const (
	historyModePurge historyMode = iota
	historyModePin
)

// runHistoryPurge is the dispatch entry for `history-purge` / `hp`.
func runHistoryPurge(args []string) {
	checkHelp(constants.CmdHistoryPurge, args)
	runHistoryRewrite(historyModePurge, args)
}

// runHistoryPin is the dispatch entry for `history-pin` / `hpin`.
func runHistoryPin(args []string) {
	checkHelp(constants.CmdHistoryPin, args)
	runHistoryRewrite(historyModePin, args)
}

// runHistoryRewrite is the shared phase pipeline for both commands.
// Each phase is its own function so this stays under 15 lines.
func runHistoryRewrite(mode historyMode, args []string) {
	opts, paths := parseHistoryArgs(args)
	opts.modeLabel = historyModeLabel(mode)
	opts.pathCount = len(paths)
	ensureFilterRepoInstalled()
	originURL := readOriginURL()
	sandbox := mirrorClone(originURL, opts)
	defer cleanupSandbox(sandbox, opts)

	pinPayloads := loadPinPayloads(mode, paths)
	runFilterRepo(mode, sandbox, paths, pinPayloads, opts)
	verifyHistoryRewrite(mode, sandbox, paths)

	if opts.dryRun {
		fmt.Fprintf(os.Stdout, constants.HistoryMsgDryRunDone, sandbox)
		return
	}
	finalizePush(sandbox, originURL, opts)
}

// historyModeLabel returns a short human label used in the confirm
// banner.
func historyModeLabel(mode historyMode) string {
	if mode == historyModePin {
		return "history-pin"
	}
	return "history-purge"
}

// loadPinPayloads reads current bytes for each path when in pin mode.
// Purge mode returns nil (no payloads needed). Errors exit 4.
func loadPinPayloads(mode historyMode, paths []string) map[string][]byte {
	if mode != historyModePin {
		return nil
	}
	out := make(map[string][]byte, len(paths))
	for _, p := range paths {
		data, err := os.ReadFile(p)
		if err != nil {
			fmt.Fprintf(os.Stderr, constants.HistoryErrBadArgs,
				fmt.Sprintf(constants.HistoryErrPathNotReadable, p, err))
			os.Exit(constants.HistoryExitBadArgs)
		}
		out[p] = data
	}
	return out
}

// cleanupSandbox honors --keep-sandbox; otherwise removes the temp
// mirror-clone. Always called from a deferred wrapper.
func cleanupSandbox(sandbox string, opts historyOpts) {
	if opts.keepSandbox {
		fmt.Fprintf(os.Stderr, constants.HistoryMsgKeepSandbox, sandbox)
		return
	}
	_ = os.RemoveAll(sandbox)
}
