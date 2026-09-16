package cmd

// Entrypoints for `gitmap history-purge` (`hp`) and `gitmap history-pin`
// (`hpin`). Both commands wrap `git filter-repo` in a mirror-clone
// sandbox so the user's working repository is never rewritten in
// place. Spec: spec/04-generic-cli/16-history-rewrite.md.

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// HistoryModeType picks which filter-repo behavior to wrap.
type HistoryModeType int

type historyMode = HistoryModeType

const (
	HistoryModeTypePurge HistoryModeType = iota
	HistoryModeTypePin
)

const (
	historyModePurge = HistoryModeTypePurge
	historyModePin   = HistoryModeTypePin
)

// runHistoryPurge is the dispatch entry for `history-purge` / `hp`.
func runHistoryPurge(args []string) error {
	checkHelp(constants.CmdHistoryPurge, args)
	runHistoryRewrite(HistoryModeTypePurge, args)

	return nil
}

// runHistoryPin is the dispatch entry for `history-pin` / `hpin`.
func runHistoryPin(args []string) error {
	checkHelp(constants.CmdHistoryPin, args)
	runHistoryRewrite(HistoryModeTypePin, args)

	return nil
}

// runHistoryRewrite is the shared phase pipeline for both commands.
func runHistoryRewrite(mode HistoryModeType, args []string) error {
	opts, paths := parseHistoryArgs(args)
	opts.modeLabel = historyModeLabel(mode)
	opts.pathCount = len(paths)
	ensureFilterRepoInstalled()
	originURL := readOriginURL()
	sandbox := mirrorClone(originURL, opts)
	defer cleanupSandbox(sandbox, opts)

	executeRewriteSandbox(mode, sandbox, paths, opts)

	return finalizeRewritePush(sandbox, originURL, opts)
}

func executeRewriteSandbox(mode HistoryModeType, sandbox string, paths []string, opts historyOpts) {
	pinPayloads := loadPinPayloads(mode, paths)
	runFilterRepo(mode, sandbox, paths, pinPayloads, opts)
	verifyHistoryRewrite(mode, sandbox, paths)
}

func finalizeRewritePush(sandbox, originURL string, opts historyOpts) error {
	if opts.dryRun {
		fmt.Fprintf(os.Stdout, constants.HistoryMsgDryRunDone, sandbox)

		return nil
	}

	finalizePush(sandbox, originURL, opts)

	return nil
}
