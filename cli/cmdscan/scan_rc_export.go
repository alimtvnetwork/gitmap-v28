package cmdscan

import (
	"fmt"
	"os"
	"path/filepath"
	"strconv"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ResolveRepoCacheRootFn is an optional test hook override.
var ResolveRepoCacheRootFn func() string

// ResolveRepoCacheRoot locates the companion repository directory for repo-cache.
func ResolveRepoCacheRoot() string {
	if ResolveRepoCacheRootFn != nil {
		root := ResolveRepoCacheRootFn()
		_ = os.MkdirAll(root, 0755)
		return root
	}

	root := resolveRepoCacheRootDefault()
	_ = os.MkdirAll(root, 0755)
	return root
}

func resolveRepoCacheRootDefault() string {
	if splitDB, err := store.OpenSpecialReposSplitDB(); err == nil {
		defer splitDB.Close()
		if rec, err := splitDB.GetSpecialRepo("repo-cache"); err == nil && rec != nil {
			if len(rec.LocalPath) > 0 {
				if info, err := os.Stat(rec.LocalPath); err == nil && info.IsDir() {
					return rec.LocalPath
				}
			}
		}
	}

	if defDB, err := store.OpenDefault(); err == nil {
		defer defDB.Close()
		if wd, err := defDB.GetDefaultWorkDir(); err == nil && wd != nil && len(wd.AbsolutePath) > 0 {
			return filepath.Join(wd.AbsolutePath, "repo-cache")
		}
	}

	if info, err := os.Stat(`D:\work`); err == nil && info.IsDir() {
		return filepath.Join(`D:\work`, "repo-cache")
	}

	if cwd, err := os.Getwd(); err == nil {
		return filepath.Join(filepath.Dir(cwd), "repo-cache")
	}

	return "repo-cache"
}

// ExportScanToRepoCache routes scan results to either the master merged manifest
// or a newly allocated sequential manifest in repo-cache/.
func ExportScanToRepoCache(records []model.ScanRecord, isSeparate, isQuiet bool) error {
	repoCacheRoot := ResolveRepoCacheRoot()
	if isSeparate {
		return AllocateNextSequentialFile(repoCacheRoot, records, isQuiet)
	}

	return MergeScanRecordsWithExistingManifest(repoCacheRoot, records, isQuiet)
}

// AllocateNextSequentialFile finds the next sequential prefix (XX-gitmap.json)
// and persists scan records as an isolated manifest.
func AllocateNextSequentialFile(root string, records []model.ScanRecord, isQuiet bool) error {
	maxSeq := resolveMaxSequentialPrefix(root)
	nextSeq := maxSeq + 1
	fileName := fmt.Sprintf("%02d-gitmap.json", nextSeq)
	targetPath := filepath.Join(root, fileName)

	exports := make([]exportRecord, 0, len(records))
	for _, r := range records {
		exports = append(exports, toExportRecord(r))
	}

	if err := writeExportRecordsJSON(targetPath, exports); err != nil {
		return apperror.WrapSimple(err, "allocate sequential repo-cache manifest")
	}

	if !isQuiet {
		fmt.Printf("  ✓ Exported %d repo(s) to sequential manifest: %s\n", len(exports), targetPath)
	}

	return nil
}

func resolveMaxSequentialPrefix(root string) int {
	maxSeq := 0
	if _, err := os.Stat(filepath.Join(root, "01-gitmap", "gitmap.json")); err == nil {
		maxSeq = 1
	}

	entries, err := os.ReadDir(root)
	if err != nil {
		return maxSeq
	}

	for _, entry := range entries {
		name := entry.Name()
		if len(name) < 2 {
			continue
		}
		if name[0] < '0' || name[0] > '9' || name[1] < '0' || name[1] > '9' {
			continue
		}
		if len(name) > 2 && name[2] != '-' {
			continue
		}
		seq, err := strconv.Atoi(name[:2])
		if err == nil && seq > maxSeq {
			maxSeq = seq
		}
	}

	return maxSeq
}
