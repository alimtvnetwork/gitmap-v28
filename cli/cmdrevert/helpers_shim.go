package cmdrevert

import (
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

var openDBFn = store.OpenDefault

func openDB() (*store.DB, error) {
	return openDBFn()
}

func getRepoPath() string {
	if RepoPathFn != nil {
		return RepoPathFn()
	}

	return ""
}
