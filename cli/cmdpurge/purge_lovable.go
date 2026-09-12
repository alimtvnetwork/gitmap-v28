package cmdpurge

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

func doPurgeLovable(repoPath string) error {
	trackedRes := getTrackedLovableFiles(repoPath)
	if trackedRes.IsFailure() {
		return trackedRes.AppError()
	}

	purged, err := removeUntrackedLovable(repoPath, trackedRes.Data)
	if err != nil {
		return err
	}

	fmt.Printf("Purged %d untracked files from .lovable directory.\n", purged)

	return nil
}

func getTrackedLovableFiles(repoPath string) result.ResultMap[string, bool] {
	out, err := runPurgeCmd("git", "-C", repoPath, "ls-files", ".lovable")
	if err != nil {
		appErr := apperror.WrapSimple(err, "git ls-files .lovable")

		return result.FailMap[string, bool](appErr)
	}

	lines := strings.Split(out, "\n")
	res := make(map[string]bool)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			res[l] = true
		}
	}

	return result.OkMap(res)
}

func removeUntrackedLovable(repoPath string, tracked map[string]bool) (int, error) {
	lovableDir := filepath.Join(repoPath, ".lovable")
	if _, err := os.Stat(lovableDir); os.IsNotExist(err) {
		return 0, nil
	}

	purged := 0
	err := filepath.Walk(lovableDir, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return nil
		}

		rel, _ := filepath.Rel(repoPath, path)
		relSlash := filepath.ToSlash(rel)
		if !tracked[relSlash] {
			_ = os.Remove(path)
			purged++
		}

		return nil
	})

	return purged, err
}
