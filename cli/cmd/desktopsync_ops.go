package cmd

import (
	"fmt"
	"os"
	"os/exec"

	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// SyncResultType represents the outcome of syncing one repo.
type SyncResultType int

type syncResult = SyncResultType

const (
	SyncResultTypeAdded SyncResultType = iota
	SyncResultTypeSkipped
	SyncResultTypeFailed
)

const (
	syncAdded   = SyncResultTypeAdded
	syncSkipped = SyncResultTypeSkipped
	syncFailed  = SyncResultTypeFailed
)

// syncOne attempts to register a single repo with GitHub Desktop.
func syncOne(r model.ScanRecord, cli string) SyncResultType {
	if len(r.AbsolutePath) == 0 {
		fmt.Printf(constants.MsgDesktopSyncFailed, r.RepoName, constants.ErrNoAbsPath)

		return SyncResultTypeFailed
	}

	return syncExistingPath(r, cli)
}

// syncExistingPath checks path existence and registers with Desktop.
func syncExistingPath(r model.ScanRecord, cli string) SyncResultType {
	_, err := os.Stat(r.AbsolutePath)
	if err == nil {
		return registerOne(r.RepoName, r.AbsolutePath, cli)
	}

	fmt.Printf(constants.MsgDesktopSyncSkipped, r.RepoName)

	return SyncResultTypeSkipped
}

// registerOne calls the GitHub Desktop CLI for a single repo.
func registerOne(name, repoPath, cli string) SyncResultType {
	cmd := exec.Command(cli, repoPath)
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Printf(constants.MsgDesktopSyncFailed, name, fmt.Sprintf("%v: %s", err, output))

		return SyncResultTypeFailed
	}

	fmt.Printf(constants.MsgDesktopSyncAdded, name)

	return SyncResultTypeAdded
}

// tallyResult increments the appropriate counter.
func tallyResult(r SyncResultType, added, skipped, failed int) (int, int, int) {
	if r == SyncResultTypeAdded {
		return added + 1, skipped, failed
	}

	if r == SyncResultTypeSkipped {
		return added, skipped + 1, failed
	}

	return added, skipped, failed + 1
}
