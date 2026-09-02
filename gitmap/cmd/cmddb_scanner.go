package cmd

import (
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/store"
)

func collectSplitDBs() []DBFileInfo {
	dirs := findSplitDBDirs()
	var out []DBFileInfo
	for _, dir := range dirs {
		files, err := os.ReadDir(dir)
		if err != nil {
			continue
		}
		for _, f := range files {
			if f.IsDir() || !strings.HasSuffix(f.Name(), ".db") {
				continue
			}
			fullPath := filepath.Join(dir, f.Name())
			if item, hasInfo := inspectSplitDBFile(fullPath, f.Name()); hasInfo {
				out = append(out, item)
			}
		}
	}
	return out
}

func inspectSplitDBFile(fullPath, filename string) (DBFileInfo, bool) {
	st, err := os.Stat(fullPath)
	if err != nil {
		return DBFileInfo{}, false
	}
	slug, repoID := parseSplitDBFilename(filename)
	purpose := "Isolated repository index storing RepoFile index, SearchCache, and FileSequence."
	return DBFileInfo{
		Name:     filename,
		Path:     fullPath,
		Size:     st.Size(),
		Category: "Split Repository DB",
		Purpose:  purpose,
		RepoID:   repoID,
		RepoSlug: slug,
	}, true
}

func parseSplitDBFilename(filename string) (string, int64) {
	base := strings.TrimSuffix(filename, ".db")
	lastDash := strings.LastIndex(base, "-")
	if lastDash == -1 || lastDash == len(base)-1 {
		return base, 0
	}
	idPart := base[lastDash+1:]
	slugPart := base[:lastDash]
	id, err := strconv.ParseInt(idPart, 10, 64)
	if err != nil {
		return base, 0
	}
	return slugPart, id
}

func collectProfileDBs() []DBFileInfo {
	binDir := store.BinaryDataDir()
	files, err := os.ReadDir(binDir)
	if err != nil {
		return nil
	}
	var out []DBFileInfo
	for _, f := range files {
		if f.IsDir() || !strings.HasSuffix(f.Name(), ".db") || f.Name() == "gitmap.db" {
			continue
		}
		fullPath := filepath.Join(binDir, f.Name())
		if st, sErr := os.Stat(fullPath); sErr == nil {
			out = append(out, DBFileInfo{
				Name:     f.Name(),
				Path:     fullPath,
				Size:     st.Size(),
				Category: "Profile DB",
				Purpose:  "Workspace-isolated SQLite database for alternate profiles.",
			})
		}
	}
	return out
}
