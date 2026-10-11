// Package cmdresolver implements the Universal 'Repo Feature' destination target resolution engine.
package cmdresolver

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// RepoFeatureTargetType enumerates classification modes for destination targets.
type RepoFeatureTargetType string

const (
	TargetTypeRemoteURL  RepoFeatureTargetType = "remote-url"
	TargetTypeLocalGit   RepoFeatureTargetType = "local-git"
	TargetTypeLocalNoGit RepoFeatureTargetType = "local-no-git"
	TargetTypeBareSlug   RepoFeatureTargetType = "bare-slug"
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

	if info, errStat := os.Stat(trimmed); errStat == nil && info.IsDir() {
		gitPath := filepath.Join(trimmed, ".git")
		if _, errGit := os.Stat(gitPath); errGit == nil {
			return resolveExistingLocalGit(trimmed)
		}
		return resolveLocalNoGitFeature(trimmed)
	}

	return resolveBareSlugFeature(trimmed)
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

	if db, errDB := store.OpenDefault(); errDB == nil {
		defer db.Close()
		if records, errList := db.ListRepos(); errList == nil {
			for _, rec := range records {
				if strings.EqualFold(rec.HTTPSUrl, url) || strings.EqualFold(rec.SSHUrl, url) || strings.EqualFold(rec.DiscoveredURL, url) || strings.EqualFold(rec.Slug, slug) {
					absPath := rec.AbsolutePath
					if absPath == "" {
						absPath = rec.RelativePath
					}
					if _, statErr := os.Stat(filepath.Join(absPath, ".git")); statErr == nil {
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
				}
			}
		}
	}

	workspaceRoot := resolveDefaultWorkspaceRoot()
	destPath := filepath.Join(workspaceRoot, slug)
	_ = os.MkdirAll(workspaceRoot, 0755)

	cmd := exec.Command("git", "clone", url, destPath)
	if out, errClone := cmd.CombinedOutput(); errClone != nil {
		return nil, fmt.Errorf("failed to clone remote repository %s: %s (%w)", url, string(out), errClone)
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

	cmd := exec.Command("git", "-C", abs, "remote", "get-url", "origin")
	remoteURL, _ := cmd.CombinedOutput()

	return &ResolvedRepoFeature{
		OriginalInput: dir,
		TargetType:    TargetTypeLocalGit,
		CanonicalSlug: slug,
		LocalPath:     abs,
		RemoteURL:     strings.TrimSpace(string(remoteURL)),
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
		return nil, fmt.Errorf("failed to initialize git in %s: %s (%w)", abs, string(out), errInit)
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

func resolveBareSlugFeature(slug string) (*ResolvedRepoFeature, error) {
	cleanSlug := SanitizeSlug(slug)
	if cleanSlug == "" {
		cleanSlug = "repo"
	}

	workspaceRoot := resolveDefaultWorkspaceRoot()
	destPath := filepath.Join(workspaceRoot, cleanSlug)

	if _, statErr := os.Stat(destPath); os.IsNotExist(statErr) {
		if errMk := os.MkdirAll(destPath, 0755); errMk != nil {
			return nil, fmt.Errorf("failed to create directory %s: %w", destPath, errMk)
		}
	}

	gitDir := filepath.Join(destPath, ".git")
	isNew := false
	if _, statGit := os.Stat(gitDir); os.IsNotExist(statGit) {
		initCmd := exec.Command("git", "-C", destPath, "init", "-b", "main")
		_ = initCmd.Run()
		isNew = true
	}

	cmd := exec.Command("git", "-C", destPath, "remote", "get-url", "origin")
	remoteURL, _ := cmd.CombinedOutput()

	return &ResolvedRepoFeature{
		OriginalInput: slug,
		TargetType:    TargetTypeBareSlug,
		CanonicalSlug: cleanSlug,
		LocalPath:     destPath,
		RemoteURL:     strings.TrimSpace(string(remoteURL)),
		IsNewRepo:     isNew,
		IsCloned:      false,
	}, nil
}

func resolveDefaultWorkspaceRoot() string {
	if defDB, err := store.OpenDefault(); err == nil {
		defer defDB.Close()
		if wd, err := defDB.GetDefaultWorkDir(); err == nil && wd != nil && len(wd.AbsolutePath) > 0 {
			if info, errStat := os.Stat(wd.AbsolutePath); errStat == nil && info.IsDir() {
				return wd.AbsolutePath
			}
		}
	}

	if info, err := os.Stat(`D:\work`); err == nil && info.IsDir() {
		return `D:\work`
	}

	cwd, err := os.Getwd()
	if err == nil && cwd != "" {
		parent := filepath.Dir(cwd)
		if parent != "" && parent != cwd {
			return parent
		}
		return cwd
	}

	return "."
}
