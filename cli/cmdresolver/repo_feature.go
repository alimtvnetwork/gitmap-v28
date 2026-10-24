// Package cmdresolver implements the Universal 'Repo Feature' destination target resolution engine.
package cmdresolver

import (
	"fmt"
	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/model"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)
import "github.com/alimtvnetwork/gitmap-v28/cli/cmdresolver/enums"

// RepoFeatureTargetType enumerates classification modes for destination targets.
// RepoFeatureTargetType aliases the strongly-typed enum from the enums package.
type RepoFeatureTargetType = enums.RepoFeatureTargetType

const (
	TargetTypeRemoteURL  = enums.TargetTypeRemoteURL
	TargetTypeLocalGit   = enums.TargetTypeLocalGit
	TargetTypeLocalNoGit = enums.TargetTypeLocalNoGit
	TargetTypeBareSlug   = enums.TargetTypeBareSlug
)

// ResolvedRepoFeature represents the resolved workspace destination for a 'repo feature' target.
type ResolvedRepoFeature struct {
	OriginalInput string                `json:"originalInput"`
	TargetType    RepoFeatureTargetType `json:"targetType"`
	CanonicalSlug string                `json:"canonicalSlug"`
	LocalPath     string                `json:"localPath"`
	RemoteURL     string                `json:"remoteUrl,omitempty"`
	IsNewRepo     bool                  `json:"isNewRepo"`
	IsCloned      bool                  `json:"isCloned"`
}

// ResolveRepoFeature resolves any destination identifier per the Repo Feature specification.
func ResolveRepoFeature(target string) (*ResolvedRepoFeature, error) {
	trimmed := strings.TrimSpace(target)
	if trimmed == "" {
		trimmed = "."
	}

	if isRemoteGitURL(trimmed) {
		return resolveRemoteURLFeature(trimmed)
	}

	if isExistingLocalDir(trimmed) {
		return resolveExistingDirectoryFeature(trimmed)
	}

	return resolveBareSlugFeature(trimmed)
}

func isExistingLocalDir(target string) bool {
	info, errStat := os.Stat(target)
	return errStat == nil && info.IsDir()
}

func resolveExistingDirectoryFeature(dir string) (*ResolvedRepoFeature, error) {
	gitPath := filepath.Join(dir, ".git")
	if _, errGit := os.Stat(gitPath); errGit == nil {
		return resolveExistingLocalGit(dir)
	}

	return resolveLocalNoGitFeature(dir)
}

func isRemoteGitURL(target string) bool {
	lower := strings.ToLower(target)
	return strings.HasPrefix(lower, "https://") ||
		strings.HasPrefix(lower, "http://") ||
		strings.HasPrefix(lower, "git@") ||
		strings.HasPrefix(lower, "ssh://")
}

func extractSlugFromURL(url string) string {
	clean := strings.TrimSuffix(url, ".git")
	clean = strings.TrimRight(clean, "/")
	parts := strings.Split(clean, "/")
	if len(parts) > 0 {
		return parts[len(parts)-1]
	}
	return "repo"
}

func resolveRemoteURLFeature(url string) (*ResolvedRepoFeature, error) {
	slug := extractSlugFromURL(url)

	if absPath := findMatchingRepoPath(url, slug); len(absPath) > 0 {
		return &ResolvedRepoFeature{
			OriginalInput: url,
			TargetType:    TargetTypeRemoteURL,
			CanonicalSlug: slug,
			LocalPath:     absPath,
			RemoteURL:     url,
			IsNewRepo:     false,
			IsCloned:      false,
		}, nil
	}

	workspaceRoot := resolveDefaultWorkspaceRoot()
	destPath := filepath.Join(workspaceRoot, slug)
	_ = os.MkdirAll(workspaceRoot, 0755)

	cmd := exec.Command("git", "clone", url, destPath)
	if out, errClone := cmd.CombinedOutput(); errClone != nil {
		return nil, apperror.WrapSimple(errClone, fmt.Sprintf("failed to clone remote repository %s: %s", url, string(out)))
	}

	return &ResolvedRepoFeature{
		OriginalInput: url,
		TargetType:    TargetTypeRemoteURL,
		CanonicalSlug: slug,
		LocalPath:     destPath,
		RemoteURL:     url,
		IsNewRepo:     false,
		IsCloned:      true,
	}, nil
}

func resolveExistingLocalGit(dir string) (*ResolvedRepoFeature, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	slug := filepath.Base(abs)

	remoteURL := ""
	if out, errRemote := exec.Command("git", "-C", abs, "remote", "get-url", "origin").Output(); errRemote == nil {
		remoteURL = strings.TrimSpace(string(out))
	}

	return &ResolvedRepoFeature{
		OriginalInput: dir,
		TargetType:    TargetTypeLocalGit,
		CanonicalSlug: slug,
		LocalPath:     abs,
		RemoteURL:     remoteURL,
		IsNewRepo:     false,
		IsCloned:      false,
	}, nil
}

func resolveLocalNoGitFeature(dir string) (*ResolvedRepoFeature, error) {
	abs, err := filepath.Abs(dir)
	if err != nil {
		abs = dir
	}
	slug := filepath.Base(abs)

	initCmd := exec.Command("git", "-C", abs, "init", "-b", "main")
	if out, errInit := initCmd.CombinedOutput(); errInit != nil {
		return nil, apperror.WrapSimple(errInit, fmt.Sprintf("failed to initialize git in %s: %s", abs, string(out)))
	}

	return &ResolvedRepoFeature{
		OriginalInput: dir,
		TargetType:    TargetTypeLocalNoGit,
		CanonicalSlug: slug,
		LocalPath:     abs,
		IsNewRepo:     true,
		IsCloned:      false,
	}, nil
}

// SanitizeSlug normalizes any target slug into a clean kebab-case identifier.
func SanitizeSlug(slug string) string {
	clean := strings.ToLower(strings.TrimSpace(slug))
	clean = strings.ReplaceAll(clean, " ", "-")
	clean = strings.ReplaceAll(clean, "_", "-")

	var b strings.Builder
	for _, r := range clean {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') || r == '-' {
			b.WriteRune(r)
		}
	}
	clean = b.String()

	for strings.Contains(clean, "--") {
		clean = strings.ReplaceAll(clean, "--", "-")
	}

	return strings.Trim(clean, "-")
}

func findMatchingRepoPath(url, slug string) string {
	db, errDB := store.OpenDefault()
	if errDB != nil {
		return ""
	}
	defer db.Close()

	records, errList := db.ListRepos()
	if errList != nil {
		return ""
	}

	for _, rec := range records {
		if path := matchRecordPath(rec, url, slug); len(path) > 0 {
			return path
		}
	}

	return ""
}

func matchRecordPath(rec model.ScanRecord, url, slug string) string {
	if !isRecordMatchingURLOrSlug(rec, url, slug) {
		return ""
	}

	absPath := rec.AbsolutePath
	if len(absPath) == 0 {
		absPath = rec.RelativePath
	}

	if _, statErr := os.Stat(filepath.Join(absPath, ".git")); statErr == nil {
		return absPath
	}

	return ""
}

func isRecordMatchingURLOrSlug(rec model.ScanRecord, url, slug string) bool {
	return strings.EqualFold(rec.HTTPSUrl, url) ||
		strings.EqualFold(rec.SSHUrl, url) ||
		strings.EqualFold(rec.DiscoveredURL, url) ||
		strings.EqualFold(rec.Slug, slug)
}

func resolveBareSlugFeature(slug string) (*ResolvedRepoFeature, error) {
	cleanSlug := SanitizeSlug(slug)
	if cleanSlug == "" {
		cleanSlug = "repo"
	}

	workspaceRoot := resolveDefaultWorkspaceRoot()
	destPath := filepath.Join(workspaceRoot, cleanSlug)

	if errMk := ensureDirExists(destPath); errMk != nil {
		return nil, apperror.WrapSimple(errMk, fmt.Sprintf("failed to create directory %s", destPath))
	}

	gitDir := filepath.Join(destPath, ".git")
	isNew := false
	if _, statGit := os.Stat(gitDir); os.IsNotExist(statGit) {
		initCmd := exec.Command("git", "-C", destPath, "init", "-b", "main")
		_ = initCmd.Run()
		isNew = true
	}

	remoteURL := ""
	if out, errRemote := exec.Command("git", "-C", destPath, "remote", "get-url", "origin").Output(); errRemote == nil {
		remoteURL = strings.TrimSpace(string(out))
	}

	return &ResolvedRepoFeature{
		OriginalInput: slug,
		TargetType:    TargetTypeBareSlug,
		CanonicalSlug: cleanSlug,
		LocalPath:     destPath,
		RemoteURL:     remoteURL,
		IsNewRepo:     isNew,
		IsCloned:      false,
	}, nil
}
