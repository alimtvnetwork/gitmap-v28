package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/cli/cliexit"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

// historyModeLabel returns a short human label used in the confirm banner.
func historyModeLabel(mode HistoryModeType) string {
	if mode == HistoryModeTypePin {
		return "history-pin"
	}

	return "history-purge"
}

// loadPinPayloads reads current bytes for each path when in pin mode.
func loadPinPayloads(mode HistoryModeType, paths []string) map[string][]byte {
	if mode != HistoryModeTypePin {
		return nil
	}

	out := make(map[string][]byte, len(paths))
	for _, p := range paths {
		out[p] = readPinPathPayload(p)
	}

	return out
}

func readPinPathPayload(p string) []byte {
	data, err := os.ReadFile(p)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.HistoryErrBadArgs,
			fmt.Sprintf(constants.HistoryErrPathNotReadable, p, err))
		cliexit.HandleError(nil, constants.HistoryExitBadArgs)
	}

	return data
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
