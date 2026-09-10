package cmd

import (
	"os"
	"path/filepath"
	"strings"
)

func isInsidePath(current, target string) bool {
	cClean := filepath.Clean(current)
	tClean := filepath.Clean(target)
	if cClean == tClean {
		return true
	}
	return strings.HasPrefix(cClean, tClean+string(filepath.Separator))
}

func escapeFromDir(corruptedPath, safeDir string) (string, string, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return "", "", err
	}
	if !isInsidePath(cwd, corruptedPath) {
		return "", "", nil
	}
	if err := os.Chdir(safeDir); err != nil {
		return "", "", err
	}
	return cwd, safeDir, nil
}

func handleDirEscape(dir, safeDir string, res *CleanResult) {
	from, to, err := escapeFromDir(dir, safeDir)
	if err == nil && from != "" {
		res.EscapedFrom = from
		res.EscapedTo = to
	}
}

func cleanSingleDir(info CorruptedDirInfo, opts CleanOptions, targetDir, safeDir string, res *CleanResult) error {
	handleDirEscape(info.Path, safeDir, res)
	rec, _ := RecoverCorruptedDirAssets(info, targetDir)
	res.RecoveredFiles = append(res.RecoveredFiles, rec...)
	res.RemovedDirs = append(res.RemovedDirs, info.Path)
	if opts.IsDryRun {
		return nil
	}
	return os.RemoveAll(info.Path)
}

// CleanCorruptedDirs detects, recovers assets from, and removes corrupted directories.
func CleanCorruptedDirs(opts CleanOptions) (CleanResult, error) {
	detected, err := DetectCorruptedDirs()
	if err != nil {
		return CleanResult{}, err
	}
	targetDir := resolveCanonicalInstallDir()
	homeDir, _ := os.UserHomeDir()
	safeDir := homeDir
	if safeDir == "" {
		safeDir = targetDir
	}
	res := CleanResult{DetectedDirs: detected}
	for _, info := range detected {
		if err := cleanSingleDir(info, opts, targetDir, safeDir, &res); err != nil {
			return res, err
		}
	}
	return res, nil
}
