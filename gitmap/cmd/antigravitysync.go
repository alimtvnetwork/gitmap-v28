package cmd

import (
	"github.com/alimtvnetwork/gitmap-v28/gitmap/workspacesync"
)

func syncAntigravityHelper(absPath, repoName string) {
	workspacesync.SyncAntigravity(absPath, repoName)
}
