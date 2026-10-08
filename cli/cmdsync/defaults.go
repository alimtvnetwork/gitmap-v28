package cmdsync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagent"
)

// DefaultTargetRepos defines the 43 canonical fleet repositories.
var DefaultTargetRepos = []string{
	"02-prompts/ai-empathy-prompt-tuner",
	"alim-cv",
	"alim-karim-profile",
	"aukgit/alim.karim.profile",
	"antigravity-manager",
	"cat-my",
	"03-aukgo/core",
	"digital-name-card",
	"presentations-repos/flat-slide-show",
	"gitlogger-new",
	"gitmap",
	"presentations-repos/global-ppt-v1",
	"presentations-repos/hiltrax",
	"icon-coding-guidelines",
	"img-pdf",
	"presentations-repos/ki-health-ppt",
	"aukgit/kubernetes-training",
	"lara-licensing",
	"lara-publishing",
	"laravel-automation",
	"letsmarknow-ui",
	"letsmarknow",
	"macro-ahk",
	"presentations-repos/maid-app-spec-presentation",
	"movie-cli",
	"03-aukgo/pathhelper",
	"presentations-repos/presentation-aug-2026-plans-alim",
	"02-prompts/prompts-connect",
	"punam-case-studies-v1",
	"presentations-repos/rasia-logo",
	"scripts-fixer",
	"presentations-repos/slides-spec",
	"spec-builder",
	"web-system/sweet-digs-finder",
	"ui-prompts-cat",
	"presentations-repos/white-presentation-v1",
	"workflowy-ui",
	"workflowy",
	"wp-exam",
	"wp-git-log",
	"wp-html-automate",
	"wp-link-manager",
	"wp-onboarding",
}

// ResolveSourceRoot detects the canonical coding-guidelines repository root.
func ResolveSourceRoot() (string, error) {
	cwd, err := os.Getwd()
	if err == nil && isCanonicalSourceRepo(cwd) && strings.Contains(strings.ToLower(cwd), "coding-guidelines") {
		return filepath.Clean(cwd), nil
	}

	repoRoot := cmdagent.FindRepoRoot()
	if isCanonicalSourceRepo(repoRoot) && strings.Contains(strings.ToLower(repoRoot), "coding-guidelines") {
		return filepath.Clean(repoRoot), nil
	}

	parent := filepath.Dir(repoRoot)
	cgPath := filepath.Join(parent, "coding-guidelines")
	if isCanonicalSourceRepo(cgPath) {
		return filepath.Clean(cgPath), nil
	}

	// Sibling check from cwd
	if err == nil {
		sibling := filepath.Join(filepath.Dir(cwd), "coding-guidelines")
		if isCanonicalSourceRepo(sibling) {
			return filepath.Clean(sibling), nil
		}
	}

	return "", fmt.Errorf("canonical source repository (coding-guidelines) not found")
}

func isCanonicalSourceRepo(dir string) bool {
	promptsDir := filepath.Join(dir, "01-prompts")
	skillsDir := filepath.Join(dir, ".agents", "skills")
	pInfo, pErr := os.Stat(promptsDir)
	sInfo, sErr := os.Stat(skillsDir)
	return pErr == nil && pInfo.IsDir() && sErr == nil && sInfo.IsDir()
}

// ResolveTargetProjects resolves target repositories based on flags and defaults.
func ResolveTargetProjects(opts SyncOptions, sourceRoot string) ([]ProjectConfig, error) {
	workRoot := filepath.Dir(sourceRoot)
	var targets []ProjectConfig

	if strings.TrimSpace(opts.Projects) != "" {
		parsed, parseErr := parseProjectsInput(strings.TrimSpace(opts.Projects), workRoot)
		if parseErr != nil {
			return nil, parseErr
		}
		targets = parsed
	} else {
		for _, rel := range DefaultTargetRepos {
			fullPath := filepath.Join(workRoot, filepath.FromSlash(rel))
			name := filepath.Base(rel)
			targets = append(targets, ProjectConfig{
				Name:   name,
				Path:   fullPath,
				Folder: rel,
				Mode:   "all",
			})
		}
	}

	var filtered []ProjectConfig
	sourceClean := filepath.Clean(sourceRoot)

	for _, t := range targets {
		if filepath.Clean(t.Path) == sourceClean {
			continue
		}
		if IsNonOwnedRepo(t.Name) || IsNonOwnedRepo(t.Path) {
			continue
		}

		if len(opts.Repos) > 0 {
			matches := false
			for _, r := range opts.Repos {
				rLower := strings.ToLower(r)
				if strings.Contains(strings.ToLower(t.Name), rLower) ||
					strings.Contains(strings.ToLower(t.Folder), rLower) {
					matches = true
					break
				}
			}
			if !matches {
				continue
			}
		}

		filtered = append(filtered, t)
	}

	return filtered, nil
}

func parseProjectsInput(input, workRoot string) ([]ProjectConfig, error) {
	raw := input
	if _, err := os.Stat(input); err == nil {
		bytes, readErr := os.ReadFile(input)
		if readErr != nil {
			return nil, fmt.Errorf("failed to read projects file %s: %w", input, readErr)
		}
		raw = string(bytes)
	}

	rawTrim := strings.TrimSpace(raw)
	if strings.HasPrefix(rawTrim, "[") {
		// Try parsing as []ProjectConfig
		var list []ProjectConfig
		if err := json.Unmarshal([]byte(rawTrim), &list); err == nil {
			for i := range list {
				if list[i].Path == "" && list[i].Folder != "" {
					list[i].Path = filepath.Join(workRoot, filepath.FromSlash(list[i].Folder))
				}
				if list[i].Name == "" {
					list[i].Name = filepath.Base(list[i].Path)
				}
			}
			return list, nil
		}

		// Try parsing as []string
		var strList []string
		if err := json.Unmarshal([]byte(rawTrim), &strList); err == nil {
			var res []ProjectConfig
			for _, s := range strList {
				p := s
				if !filepath.IsAbs(p) {
					p = filepath.Join(workRoot, filepath.FromSlash(s))
				}
				res = append(res, ProjectConfig{
					Name:   filepath.Base(p),
					Path:   p,
					Folder: s,
					Mode:   "all",
				})
			}
			return res, nil
		}
	}

	// If directory path provided, scan for immediate child git repos
	if info, err := os.Stat(input); err == nil && info.IsDir() {
		entries, _ := os.ReadDir(input)
		var res []ProjectConfig
		for _, e := range entries {
			if !e.IsDir() {
				continue
			}
			child := filepath.Join(input, e.Name())
			gitDir := filepath.Join(child, ".git")
			if _, gErr := os.Stat(gitDir); gErr == nil {
				res = append(res, ProjectConfig{
					Name:   e.Name(),
					Path:   child,
					Folder: e.Name(),
					Mode:   "all",
				})
			}
		}
		return res, nil
	}

	return nil, fmt.Errorf("invalid projects input: expected JSON array or existing directory")
}
