package cmd

import (
	"os"
	"path/filepath"
)

func processDirectoryBrokenLinks(dir string, opts FixLinkOptions) ([]LinkResult, error) {
	broken, err := scanBrokenLinks(dir, opts.IsRecursive)
	if err != nil {
		return nil, err
	}

	if len(broken) == 0 {
		return []LinkResult{{Path: dir, IsHealthy: true, Message: "no broken symlinks found in directory"}}, nil
	}

	var results []LinkResult
	for _, p := range broken {
		res := repairSingleLink(p, opts)
		results = append(results, res)
	}

	return results, nil
}

func scanBrokenLinks(dir string, isRecursive bool) ([]string, error) {
	var broken []string
	err := inspectEntriesInDir(dir, isRecursive, &broken)
	if err != nil {
		return nil, err
	}

	return broken, nil
}

func inspectEntriesInDir(dir string, isRecursive bool, broken *[]string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return err
	}

	for _, entry := range entries {
		fullPath := filepath.Join(dir, entry.Name())
		checkAndCollectEntry(fullPath, entry, isRecursive, broken)
	}

	return nil
}

func checkAndCollectEntry(fullPath string, entry os.DirEntry, isRecursive bool, broken *[]string) {
	if entry.Type()&os.ModeSymlink != 0 {
		collectIfBroken(fullPath, broken)

		return
	}

	if isRecursive && entry.IsDir() {
		_ = inspectEntriesInDir(fullPath, isRecursive, broken)
	}
}

func collectIfBroken(path string, broken *[]string) {
	_, err := os.Stat(path)
	if err != nil {
		*broken = append(*broken, path)
	}
}
