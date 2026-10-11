package cmdscan

import (
	"os"
	"path/filepath"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func init() {
	store.OnRepoAddedHook = TryIncludeRepoInRepoCacheJSON
}

// TryIncludeRepoInRepoCacheJSON checks whether repo-cache exists in the default work directory,
// and if present, synchronizes the provided repository records into the repo-cache gitmap.json.
func TryIncludeRepoInRepoCacheJSON(records []model.ScanRecord) {
	if len(records) == 0 {
		return
	}

	root, _, hasRepoCache := probeRepoCacheDir()
	if !hasRepoCache || len(root) == 0 {
		return
	}

	_ = MergeScanRecordsWithExistingManifest(root, records, true)
	maybeSyncRootRepoCacheJSON(root, records)
}

func maybeSyncRootRepoCacheJSON(root string, records []model.ScanRecord) {
	manifestPath := filepath.Join(root, "gitmap.json")
	if isDir(root) && isRegularFileOnDisk(manifestPath) {
		_ = mergeExistingManifest(manifestPath, records, true)
	}
}

func isRegularFileOnDisk(p string) bool {
	info, err := os.Stat(p)
	return err == nil && !info.IsDir()
}
