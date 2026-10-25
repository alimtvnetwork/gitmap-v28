package gitutil

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// FileStatusCode represents porcelain status indicators.
type FileStatusCode string

// Standard porcelain status codes.
const (
	StatusModified  FileStatusCode = "[M]"
	StatusAdded     FileStatusCode = "[A]"
	StatusDeleted   FileStatusCode = "[D]"
	StatusUntracked FileStatusCode = "[?]"
	StatusRenamed   FileStatusCode = "[R]"
)

// DirtyFileInfo details status for a single dirty file.
type DirtyFileInfo struct {
	Path       string         `json:"path"`
	StatusCode FileStatusCode `json:"statusCode"`
	IsStaged   bool           `json:"isStaged"`
}

// RepoDirtyManifest contains file-level uncommitted state details.
type RepoDirtyManifest struct {
	RepoPath       string          `json:"repoPath"`
	RepoName       string          `json:"repoName"`
	HasDirtyWork   bool            `json:"hasDirtyWork"`
	Files          []DirtyFileInfo `json:"files"`
	ModifiedCount  int             `json:"modifiedCount"`
	UntrackedCount int             `json:"untrackedCount"`
	DeletedCount   int             `json:"deletedCount"`
}

func hasGitMarker(cleanPath string) bool {
	gitMarker := filepath.Join(cleanPath, ".git")
	markerInfo, markerErr := os.Stat(gitMarker)
	if markerErr != nil {
		return false
	}

	return markerInfo.IsDir() || !markerInfo.IsDir()
}

func isInsideWorkTree(cleanPath string) bool {
	cmd := exec.Command("git", "-C", cleanPath, "rev-parse", "--is-inside-work-tree")
	out, cmdErr := cmd.Output()

	return cmdErr == nil && strings.TrimSpace(string(out)) == "true"
}

// IsValidGitRepository verifies that the target path contains a verified .git directory or worktree.
func IsValidGitRepository(targetPath string) bool {
	cleanPath := filepath.Clean(targetPath)
	info, err := os.Stat(cleanPath)
	if err != nil || !info.IsDir() {
		return false
	}

	if hasGitMarker(cleanPath) {
		return true
	}

	return isInsideWorkTree(cleanPath)
}

// IsWorkspaceContainerDirectory checks whether a path represents a parent workspace folder rather than a leaf repository.
func IsWorkspaceContainerDirectory(targetPath string) bool {
	clean := filepath.Clean(targetPath)
	base := strings.ToLower(filepath.Base(clean))
	if base == "work" || base == "git-work" || base == "workspace" || base == "projects" {
		gitMarker := filepath.Join(clean, ".git")
		_, err := os.Stat(gitMarker)

		return os.IsNotExist(err)
	}

	return false
}

// InspectDirtyManifest gathers detailed file status flags using 'git status --porcelain=v1'.
func InspectDirtyManifest(repoPath string) (RepoDirtyManifest, error) {
	manifest := RepoDirtyManifest{
		RepoPath: repoPath,
		RepoName: filepath.Base(repoPath),
	}
	if !IsValidGitRepository(repoPath) {
		return manifest, nil
	}

	out, err := exec.Command("git", "-C", repoPath, "status", "--porcelain=v1").Output()
	if err != nil {
		return manifest, err
	}

	populateManifestFiles(&manifest, string(out))

	return manifest, nil
}

func populateManifestFiles(manifest *RepoDirtyManifest, raw string) {
	lines := strings.Split(raw, "\n")
	for _, l := range lines {
		trimmed := strings.TrimRight(l, "\r")
		info := parsePorcelainLine(trimmed)
		if info != nil {
			manifest.Files = append(manifest.Files, *info)
			updateManifestCounts(manifest, info.StatusCode)
		}
	}

	manifest.HasDirtyWork = len(manifest.Files) > 0
}

func updateManifestCounts(manifest *RepoDirtyManifest, code FileStatusCode) {
	if code == StatusModified || code == StatusRenamed {
		manifest.ModifiedCount++
	}

	if code == StatusUntracked {
		manifest.UntrackedCount++
	}

	if code == StatusDeleted {
		manifest.DeletedCount++
	}
}

func parsePorcelainLine(line string) *DirtyFileInfo {
	if len(line) < 3 {
		return nil
	}

	x := line[0]
	y := line[1]
	relPath := strings.TrimSpace(line[3:])
	if strings.Contains(relPath, " -> ") {
		parts := strings.Split(relPath, " -> ")
		relPath = parts[len(parts)-1]
	}

	code, isStaged := resolveStatusCode(x, y)

	return &DirtyFileInfo{
		Path:       relPath,
		StatusCode: code,
		IsStaged:   isStaged,
	}
}

func resolveStatusCode(x, y byte) (FileStatusCode, bool) {
	if x == '?' && y == '?' {
		return StatusUntracked, false
	}

	if x == 'R' || y == 'R' {
		return StatusRenamed, x != ' '
	}

	if x == 'D' || y == 'D' {
		return StatusDeleted, x == 'D'
	}

	if x == 'A' || y == 'A' {
		return StatusAdded, x == 'A'
	}

	return StatusModified, x != ' ' && x != '?'
}
