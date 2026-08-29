package cmd

import (
	"fmt"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cliexit"
)

// runSetSourceRepo handles the hidden "set-source-repo" command.
// Called by run.ps1 after deploy to persist the current repo root
// so future "gitmap update" uses the correct source location.
func runSetSourceRepo() error {
	args := os.Args[2:]
	if len(args) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrSetSourceRepoNoPath)
		cliexit.HandleError(nil, 1)
	}

	path := args[0]
	normalized := normalizeRepoPath(path)
	if len(normalized) == 0 {
		fmt.Fprintf(os.Stderr, constants.ErrSetSourceRepoInvalid, path)
		cliexit.HandleError(nil, 1)
	}

	saveRepoPathToDB(normalized)
	fmt.Printf(constants.MsgSetSourceRepoDone, normalized)
	return nil
}
