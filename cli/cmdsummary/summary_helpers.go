package cmdsummary

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/fsutil"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// resolveWorkspaceRepositories resolves all repo scan records in cwd or parent workspace.
func resolveWorkspaceRepositories(cwd string) []model.ScanRecord {
	targets := cmdpull.ResolvePullDirectoryTargets(cwd)
	var filtered []model.ScanRecord
	for _, rec := range targets {
		if fsutil.IsExcludedScanPath(rec.AbsolutePath) {
			continue
		}
		if isScratchOrCorruptRepo(rec) {
			continue
		}
		filtered = append(filtered, rec)
	}
	if len(filtered) > 0 {
		return filtered
	}

	gitPath := filepath.Join(cwd, ".git")
	if _, err := os.Stat(gitPath); err == nil && !fsutil.IsExcludedScanPath(cwd) && fsutil.IsValidGitRepo(cwd) {
		base := filepath.Base(cwd)
		rec := model.ScanRecord{
			RepoName:     base,
			Slug:         base,
			AbsolutePath: cwd,
			RelativePath: ".",
		}
		if !isScratchOrCorruptRepo(rec) {
			return []model.ScanRecord{rec}
		}
	}

	return nil
}

func isScratchOrCorruptRepo(rec model.ScanRecord) bool {
	if !fsutil.IsValidGitRepo(rec.AbsolutePath) {
		return true
	}
	name := strings.ToLower(rec.RepoName)
	if name == "dst" || name == "src" || name == "temp" || name == "test" {
		remoteURL, _ := gitExec(rec.AbsolutePath, "remote", "get-url", "origin")
		if strings.TrimSpace(remoteURL) == "" || strings.HasPrefix(remoteURL, "error:") || strings.HasPrefix(remoteURL, "fatal:") {
			return true
		}
	}
	return false
}

// parsePorcelainStatusLines parses git status --porcelain output into uncommitted counts and pending file paths.
func parsePorcelainStatusLines(output string) (int, int, int, []string) {
	var untracked, modified, staged int
	var pendingFiles []string

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		lower := strings.ToLower(trimmed)
		if strings.HasPrefix(lower, "fatal:") || strings.HasPrefix(lower, "error:") || strings.HasPrefix(lower, "warning:") {
			continue
		}
		pendingFiles = append(pendingFiles, trimmed)
		if strings.HasPrefix(line, "??") {
			untracked++
			continue
		}
		stagedDelta, modifiedDelta := classifyStatusChars(line)
		staged += stagedDelta
		modified += modifiedDelta
	}

	return untracked, modified, staged, pendingFiles
}

func classifyStatusChars(line string) (int, int) {
	if len(line) < 2 {
		return 0, 0
	}
	stagedDelta := 0
	modifiedDelta := 0
	idxChar := line[0]
	wtChar := line[1]
	if idxChar != ' ' && idxChar != '?' {
		stagedDelta = 1
	}
	if wtChar != ' ' && wtChar != '?' {
		modifiedDelta = 1
	}
	return stagedDelta, modifiedDelta
}
