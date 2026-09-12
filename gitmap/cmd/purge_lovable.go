package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

func doPurgeLovable(repoPath string) error {
	tracked, err := getTrackedLovableFiles(repoPath)
	if err != nil {
		return err
	}

	purged, err := removeUntrackedLovable(repoPath, tracked)
	if err != nil {
		return err
	}

	fmt.Printf("Purged %d untracked files from .lovable directory.\n", purged)

	return nil
}

func getTrackedLovableFiles(repoPath string) (map[string]bool, error) {
	out, err := runPurgeCmd("git", "-C", repoPath, "ls-files", ".lovable")
	if err != nil {
		return nil, err
	}

	lines := strings.Split(out, "\n")
	res := make(map[string]bool)
	for _, l := range lines {
		l = strings.TrimSpace(l)
		if l != "" {
			res[l] = true
		}
	}

	return res, nil
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
