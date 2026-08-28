package cmd

import (
	"fmt"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/workspacesync"
)

func runCloneSync() error {
	args := argsTail()
	if len(args) == 0 {
		fmt.Println("Usage: gitmap clone-sync <url1> [url2] ...")
		return apperror.NewSimple("fatal error", "E9000")
	}

	for _, arg := range args {
		// Use internal direct clone one which handles retry, db upsert, etc
		err := executeDirectCloneOne(arg, "", false, false)
		if err != nil {
			continue
		}

		name := extractRepoNameFromURL(arg)
		// Assuming we just clone into current directory for now
		absPath := resolveCloneFolder(name, "")

		workspacesync.SyncAll(absPath, name)
	}
	return nil
}

func extractRepoNameFromURL(u string) string {
	return repoNameFromURL(u)
}
