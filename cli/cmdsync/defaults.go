package cmdsync

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdagent"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdsync/constants"
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
	"seo-writing/alim-seo-writing",
	"seo-packages",
}

// ResolveSourceRoot detects the canonical coding-guidelines repository root.
func ResolveSourceRoot() (string, error) {
	if root, ok := checkCurrentAndRepoRoot(); ok {
		return root, nil
	}
	if root, ok := checkParentAndSiblingCwd(); ok {
		return root, nil
	}

	return "", apperror.NewNotFoundError("canonical source repository (coding-guidelines) not found")
}

func checkCurrentAndRepoRoot() (string, bool) {
	cwd, err := os.Getwd()
	if err == nil && isCanonicalSourceRepo(cwd) && strings.Contains(strings.ToLower(cwd), "coding-guidelines") {
		return filepath.Clean(cwd), true
	}
	repoRoot := cmdagent.FindRepoRoot()
	if isCanonicalSourceRepo(repoRoot) && strings.Contains(strings.ToLower(repoRoot), "coding-guidelines") {
		return filepath.Clean(repoRoot), true
	}

	return "", false
}

func checkParentAndSiblingCwd() (string, bool) {
	repoRoot := cmdagent.FindRepoRoot()
	parent := filepath.Dir(repoRoot)
	cgPath := filepath.Join(parent, "coding-guidelines")
	if isCanonicalSourceRepo(cgPath) {
		return filepath.Clean(cgPath), true
	}
	cwd, err := os.Getwd()
	if siblingCg := checkSiblingCodingGuidelines(cwd, err); siblingCg != "" {
		return siblingCg, true
	}

	return "", false
}

func checkSiblingCodingGuidelines(cwd string, err error) string {
	if err != nil {
		return ""
	}
	sibling := filepath.Join(filepath.Dir(cwd), "coding-guidelines")
	if isCanonicalSourceRepo(sibling) {
		return filepath.Clean(sibling)
	}

	return ""
}

func isCanonicalSourceRepo(dir string) bool {
	promptsDir := filepath.Join(dir, constants.DirPrompts)
	skillsDir := filepath.Join(dir, filepath.FromSlash(constants.DirAgentsSkills))
	pInfo, pErr := os.Stat(promptsDir)
	sInfo, sErr := os.Stat(skillsDir)

	return pErr == nil && pInfo.IsDir() && sErr == nil && sInfo.IsDir()
}

// ResolveTargetProjects resolves target repositories based on flags and defaults.
func ResolveTargetProjects(opts SyncOptions, sourceRoot string) ([]ProjectConfig, error) {
	workRoot := filepath.Dir(sourceRoot)
	targets, err := resolveRawTargetProjects(opts, workRoot)
	if err != nil {
		return nil, err
	}

	return filterTargetProjects(targets, sourceRoot, opts.Repos), nil
}

func filterTargetProjects(targets []ProjectConfig, sourceRoot string, filters []string) []ProjectConfig {
	var filtered []ProjectConfig
	sourceClean := filepath.Clean(sourceRoot)
	for _, t := range targets {
		if isTargetProjectIncluded(t, sourceClean, filters) {
			filtered = append(filtered, t)
		}
	}

	return filtered
}

func isTargetProjectIncluded(t ProjectConfig, sourceClean string, filters []string) bool {
	if filepath.Clean(t.Path) == sourceClean {
		return false
	}
	if IsNonOwnedRepo(t.Name) || IsNonOwnedRepo(t.Path) {
		return false
	}

	return matchesAnyRepoFilter(t, filters)
}

func resolveRawTargetProjects(opts SyncOptions, workRoot string) ([]ProjectConfig, error) {
	if strings.TrimSpace(opts.Projects) != "" {
		return parseProjectsInput(strings.TrimSpace(opts.Projects), workRoot)
	}

	return defaultTargetProjects(workRoot), nil
}

func defaultTargetProjects(workRoot string) []ProjectConfig {
	var targets []ProjectConfig
	for _, rel := range DefaultTargetRepos {
		fullPath := filepath.Join(workRoot, filepath.FromSlash(rel))
		targets = append(targets, ProjectConfig{
			Name:   filepath.Base(rel),
			Path:   fullPath,
			Folder: rel,
			Mode:   constants.TargetAll,
		})
	}

	return targets
}

func matchesAnyRepoFilter(t ProjectConfig, filters []string) bool {
	if len(filters) == 0 {
		return true
	}
	tNameLower := strings.ToLower(t.Name)
	tFolderLower := strings.ToLower(t.Folder)
	for _, r := range filters {
		rLower := strings.ToLower(r)
		if strings.Contains(tNameLower, rLower) || strings.Contains(tFolderLower, rLower) {
			return true
		}
	}

	return false
}

func parseProjectsInput(input, workRoot string) ([]ProjectConfig, error) {
	raw, err := readProjectsInputRaw(input)
	if err != nil {
		return nil, err
	}

	if jsonProjects, ok := tryParseJsonProjects(strings.TrimSpace(raw), workRoot); ok {
		return jsonProjects, nil
	}

	if scanned, ok := tryScanDirectoryForGitRepos(input); ok {
		return scanned, nil
	}

	return nil, apperror.NewValidationError("invalid projects input: expected JSON array or existing directory")
}

func readProjectsInputRaw(input string) (string, error) {
	fi, err := os.Stat(input)
	if err != nil {
		return input, nil
	}
	if fi.IsDir() {
		return "", nil
	}
	bytes, readErr := os.ReadFile(input)
	if readErr != nil {
		return "", apperror.WrapSimple(readErr, fmt.Sprintf("failed to read projects file %s", input))
	}

	return string(bytes), nil
}

func tryParseJsonProjects(rawTrim string, workRoot string) ([]ProjectConfig, bool) {
	if !strings.HasPrefix(rawTrim, "[") {
		return nil, false
	}
	list, isConfigList := tryParseProjectConfigList(rawTrim, workRoot)
	if isConfigList {
		return list, true
	}

	return tryParseStringPathList(rawTrim, workRoot)
}

func tryParseProjectConfigList(raw string, workRoot string) ([]ProjectConfig, bool) {
	var list []ProjectConfig
	if err := json.Unmarshal([]byte(raw), &list); err != nil {
		return nil, false
	}
	for i := range list {
		normalizeProjectConfigEntry(&list[i], workRoot)
	}

	return list, true
}

func normalizeProjectConfigEntry(entry *ProjectConfig, workRoot string) {
	if entry.Path == "" && entry.Folder != "" {
		entry.Path = filepath.Join(workRoot, filepath.FromSlash(entry.Folder))
	}
	if entry.Name == "" {
		entry.Name = filepath.Base(entry.Path)
	}
}

func tryParseStringPathList(raw string, workRoot string) ([]ProjectConfig, bool) {
	var strList []string
	if err := json.Unmarshal([]byte(raw), &strList); err != nil {
		return nil, false
	}
	var res []ProjectConfig
	for _, s := range strList {
		res = append(res, buildProjectConfigFromPath(s, workRoot))
	}

	return res, true
}

func buildProjectConfigFromPath(pathOrRel string, workRoot string) ProjectConfig {
	p := pathOrRel
	if !filepath.IsAbs(p) {
		p = filepath.Join(workRoot, filepath.FromSlash(pathOrRel))
	}

	return ProjectConfig{
		Name:   filepath.Base(p),
		Path:   p,
		Folder: pathOrRel,
		Mode:   constants.TargetAll,
	}
}

func tryScanDirectoryForGitRepos(dirPath string) ([]ProjectConfig, bool) {
	info, err := os.Stat(dirPath)
	if err != nil || !info.IsDir() {
		return nil, false
	}
	if isGitDirectory(dirPath) {
		return []ProjectConfig{newSingleGitProjectConfig(dirPath)}, true
	}

	return scanChildGitRepos(dirPath)
}

func isGitDirectory(dirPath string) bool {
	_, err := os.Stat(filepath.Join(dirPath, constants.DirGit))

	return err == nil
}

func newSingleGitProjectConfig(dirPath string) ProjectConfig {
	baseName := filepath.Base(dirPath)

	return ProjectConfig{
		Name:   baseName,
		Path:   dirPath,
		Folder: baseName,
		Mode:   constants.TargetAll,
	}
}

func scanChildGitRepos(dirPath string) ([]ProjectConfig, bool) {
	entries, _ := os.ReadDir(dirPath)
	var res []ProjectConfig
	for _, e := range entries {
		cfg, isGit := inspectChildGitRepo(dirPath, e)
		if isGit {
			res = append(res, cfg)
		}
	}

	return res, len(res) > 0
}

func inspectChildGitRepo(parentDir string, e os.DirEntry) (ProjectConfig, bool) {
	if !e.IsDir() {
		return ProjectConfig{}, false
	}
	child := filepath.Join(parentDir, e.Name())
	if !isGitDirectory(child) {
		return ProjectConfig{}, false
	}

	return ProjectConfig{
		Name:   e.Name(),
		Path:   child,
		Folder: e.Name(),
		Mode:   constants.TargetAll,
	}, true
}
