package cmdsummary

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdpull"
	"github.com/alimtvnetwork/gitmap-v28/cli/model"
)

// resolveWorkspaceRepositories resolves all repo scan records in cwd or parent workspace.
func resolveWorkspaceRepositories(cwd string) []model.ScanRecord {
	targets := cmdpull.ResolvePullDirectoryTargets(cwd)
	if len(targets) > 0 {
		return targets
	}

	gitPath := filepath.Join(cwd, ".git")
	if _, err := os.Stat(gitPath); err == nil {
		base := filepath.Base(cwd)
		return []model.ScanRecord{
			{
				RepoName:     base,
				Slug:         base,
				AbsolutePath: cwd,
				RelativePath: ".",
			},
		}
	}

	return nil
}

// parsePorcelainStatusLines parses git status --porcelain output into uncommitted counts and pending file paths.
func parsePorcelainStatusLines(output string) (int, int, int, []string) {
	var untracked, modified, staged int
	var pendingFiles []string

	lines := strings.Split(output, "\n")
	for _, line := range lines {
		line = strings.TrimRight(line, "\r")
		if len(line) == 0 {
			continue
		}
		pendingFiles = append(pendingFiles, strings.TrimSpace(line))
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
