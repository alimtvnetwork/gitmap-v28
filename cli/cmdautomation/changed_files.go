package cmdautomation

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/result"
)

// RunChangedFiles discovers modified, staged, and untracked files.
func RunChangedFiles(opts ChangedFilesOptions) ChangedFilesResultMonad {
	start := time.Now()
	baseDir := resolveArtifactBaseDir(opts.Dir)
	headHash := fetchHeadHash(baseDir)
	rawItems := collectGitChanges(opts, baseDir)
	deduped := deduplicateChangedFiles(rawItems, opts.IsVerify, baseDir)
	res := aggregateChangedFiles(opts.BaseRef, headHash, deduped, start)
	return result.Ok(res)
}

func fetchGitCommandOutput(dir string, args ...string) string {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir
	out, err := cmd.Output()
	if err != nil {
		return ""
	}
	return string(out)
}

func fetchHeadHash(baseDir string) string {
	out := fetchGitCommandOutput(baseDir, "rev-parse", "HEAD")
	return strings.TrimSpace(out)
}

func mapStatusCode(code string) string {
	if strings.HasPrefix(code, "A") {
		return "added"
	}
	if strings.HasPrefix(code, "M") {
		return "modified"
	}
	if strings.HasPrefix(code, "D") {
		return "deleted"
	}
	if strings.HasPrefix(code, "R") {
		return "renamed"
	}
	return "changed"
}

func parseStatusLine(line string, isStaged bool) ChangedFileItem {
	parts := strings.Fields(line)
	if len(parts) < 2 {
		return ChangedFileItem{}
	}
	statusCode := parts[0]
	path := filepath.ToSlash(parts[len(parts)-1])
	status := mapStatusCode(statusCode)
	return ChangedFileItem{
		Path:      path,
		Status:    status,
		Extension: filepath.Ext(path),
		IsStaged:  isStaged,
		IsExists:  true,
	}
}

func parseGitNameStatus(output string, isStaged bool) []ChangedFileItem {
	lines := strings.Split(output, "\n")
	items := make([]ChangedFileItem, 0, len(lines))
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if len(trimmed) == 0 {
			continue
		}
		item := parseStatusLine(trimmed, isStaged)
		if len(item.Path) > 0 {
			items = append(items, item)
		}
	}
	return items
}

func fetchUntrackedFiles(baseDir string) []ChangedFileItem {
	out := fetchGitCommandOutput(baseDir, "ls-files", "--others", "--exclude-standard")
	lines := strings.Split(out, "\n")
	items := make([]ChangedFileItem, 0, len(lines))
	for _, l := range lines {
		trimmed := strings.TrimSpace(l)
		if len(trimmed) == 0 {
			continue
		}
		p := filepath.ToSlash(trimmed)
		items = append(items, ChangedFileItem{
			Path:      p,
			Status:    "untracked",
			Extension: filepath.Ext(p),
			IsStaged:  false,
			IsExists:  true,
		})
	}
	return items
}

func fetchStagedFiles(baseDir string) []ChangedFileItem {
	out := fetchGitCommandOutput(baseDir, "diff", "--cached", "--name-status")
	return parseGitNameStatus(out, true)
}

func fetchUnstagedFiles(baseDir string) []ChangedFileItem {
	out := fetchGitCommandOutput(baseDir, "diff", "--name-status")
	return parseGitNameStatus(out, false)
}

func fetchBaseDiffFiles(baseDir string, baseRef string) []ChangedFileItem {
	rangeSpec := fmt.Sprintf("%s...HEAD", baseRef)
	out := fetchGitCommandOutput(baseDir, "diff", "--name-status", rangeSpec)
	return parseGitNameStatus(out, false)
}

func collectGitChanges(opts ChangedFilesOptions, baseDir string) []ChangedFileItem {
	var all []ChangedFileItem
	hasBase := len(opts.BaseRef) > 0
	if hasBase {
		baseItems := fetchBaseDiffFiles(baseDir, opts.BaseRef)
		all = append(all, baseItems...)
	}
	stagedItems := fetchStagedFiles(baseDir)
	all = append(all, stagedItems...)
	if opts.IsStagedOnly {
		return all
	}
	unstagedItems := fetchUnstagedFiles(baseDir)
	all = append(all, unstagedItems...)
	untrackedItems := fetchUntrackedFiles(baseDir)
	all = append(all, untrackedItems...)
	return all
}

func verifyFileExistence(baseDir string, relPath string) bool {
	full := filepath.Join(baseDir, relPath)
	info, err := os.Stat(full)
	if err != nil {
		return false
	}
	if info.IsDir() {
		return false
	}
	return true
}

func deduplicateChangedFiles(items []ChangedFileItem, isVerify bool, baseDir string) []ChangedFileItem {
	seen := make(map[string]bool)
	var deduped []ChangedFileItem
	for _, it := range items {
		if seen[it.Path] {
			continue
		}
		seen[it.Path] = true
		if isVerify {
			it.IsExists = verifyFileExistence(baseDir, it.Path)
		}
		deduped = append(deduped, it)
	}
	return deduped
}

func aggregateChangedFiles(baseRef string, headHash string, items []ChangedFileItem, start time.Time) ChangedFilesResult {
	return ChangedFilesResult{
		BaseRef:    baseRef,
		HeadHash:   headHash,
		TotalFiles: len(items),
		Files:      items,
		Duration:   time.Since(start),
		IsSuccess:  true,
	}
}
