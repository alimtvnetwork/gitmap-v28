// Package cmdagy — agy_recreate_resolve.go resolves target Antigravity projects for recreate.
package cmdagy

import (
	"fmt"
	"net/url"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/gitutil"
)

// ResolveAgyRecreateTargets resolves raw arguments or current working directory into AgyProject targets.
func ResolveAgyRecreateTargets(args []string, projects []AgyProject) ([]AgyProject, error) {
	return ResolveAgyRecreateTargetsWithConfirm(args, projects, recreateConfirmFlag)
}

// ResolveAgyRecreateTargetsWithConfirm resolves raw arguments or cwd into AgyProject targets with explicit confirm flag.
func ResolveAgyRecreateTargetsWithConfirm(args []string, projects []AgyProject, isConfirm bool) ([]AgyProject, error) {
	if len(args) == 0 || (len(args) == 1 && args[0] == ".") {
		return resolveCurrentDirTarget(projects, isConfirm)
	}
	return resolveSpecifiedTargets(args, projects, isConfirm)
}

func resolveCurrentDirTarget(projects []AgyProject, isConfirm bool) ([]AgyProject, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, apperror.WrapSimple(err, "get current working directory")
	}
	if IsRestrictedSystemOrHomeDir(cwd) {
		return nil, apperror.NewSimple(fmt.Sprintf("restricted system or home directory cannot be recreated as an Antigravity project: %s", cwd), "E9101")
	}
	return resolveWorkingDirAdHoc(projects, cwd, isConfirm)
}

func resolveWorkingDirAdHoc(projects []AgyProject, cwd string, isConfirm bool) ([]AgyProject, error) {
	targetDir := resolveWorkingRepoRoot(cwd)
	if matched := findProjectByPathOrCwd(projects, targetDir, cwd); matched != nil {
		return []AgyProject{*matched}, nil
	}
	if !isGitRepo(targetDir) {
		return nil, apperror.NewSimple(fmt.Sprintf("directory %q is neither an existing Antigravity project nor a valid Git repository; refusing to create ad-hoc project", targetDir), "E9102")
	}
	if !isConfirm {
		return nil, apperror.NewSimple(fmt.Sprintf("directory %q is not an existing registered Antigravity project; pass --confirm (-y) to confirm creating and recreating this workspace", targetDir), "E9103")
	}
	return []AgyProject{buildAdHocProject(targetDir)}, nil
}

func findProjectByPathOrCwd(projects []AgyProject, targetDir, cwd string) *AgyProject {
	if matched := findProjectByPath(projects, targetDir); matched != nil {
		return matched
	}
	return findProjectByPath(projects, cwd)
}

func findProjectByPath(projects []AgyProject, targetDir string) *AgyProject {
	for _, p := range projects {
		if isMatchingProjectPath(p.GetPath(), targetDir) {
			return &p
		}
	}
	return nil
}

func resolveWorkingRepoRoot(dirPath string) string {
	root, err := gitutil.RepoRoot(dirPath)
	if err == nil && len(root) > 0 {
		return root
	}
	return dirPath
}

func isMatchingProjectPath(pathA, pathB string) bool {
	if pathA == "" || pathB == "" {
		return false
	}
	cleanA := strings.ToLower(filepath.Clean(pathA))
	cleanB := strings.ToLower(filepath.Clean(pathB))
	return cleanA == cleanB
}

func buildAdHocProject(targetDir string) AgyProject {
	absPath, err := filepath.Abs(targetDir)
	if err == nil {
		targetDir = absPath
	}
	name := filepath.Base(targetDir)
	return AgyProject{
		ID:   "",
		Name: name,
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{
					GitFolder: &AgyGitFolder{
						FolderURI:     buildFolderURI(targetDir),
						DefaultBranch: "main",
					},
				},
			},
		},
	}
}

func buildFolderURI(dirPath string) string {
	slashPath := filepath.ToSlash(dirPath)
	if len(slashPath) > 1 && slashPath[1] == ':' {
		return "file:///" + string(slashPath[0]) + "%3A" + url.PathEscape(slashPath[2:])
	}
	if len(slashPath) > 0 && slashPath[0] == '/' {
		return "file://" + url.PathEscape(slashPath)
	}
	return "file:///" + url.PathEscape(slashPath)
}

func resolveSpecifiedTargets(args []string, projects []AgyProject, isConfirm bool) ([]AgyProject, error) {
	tokens := expandCommaSeparatedTokens(args)
	var results []AgyProject
	seen := make(map[string]bool)

	for _, token := range tokens {
		matched, err := resolveSingleRecreateToken(token, projects, isConfirm)
		if err != nil {
			return nil, err
		}
		results = appendUniqueProjects(results, matched, seen)
	}
	return results, nil
}

func resolveSingleRecreateToken(token string, projects []AgyProject, isConfirm bool) ([]AgyProject, error) {
	if seqMatch, isSeq := matchBySequence(token, projects); isSeq {
		return []AgyProject{seqMatch}, nil
	}
	if idMatches := matchByID(token, projects); len(idMatches) > 0 {
		return idMatches, nil
	}
	if pathProjects := resolveLikelyPathToken(token, projects); len(pathProjects) > 0 {
		return pathProjects, nil
	}
	if slugMatches := matchBySlug(token, projects); len(slugMatches) > 0 {
		return slugMatches, nil
	}
	return resolveDirectoryOrFolderTarget(token, projects, isConfirm)
}

func resolveLikelyPathToken(token string, projects []AgyProject) []AgyProject {
	if !isLikelyPath(token) {
		return nil
	}
	if pathMatches := matchByExactPath(token, projects); len(pathMatches) > 0 {
		return pathMatches
	}
	return collectProjectsUnderFolder(token, projects)
}

func isLikelyPath(token string) bool {
	return strings.Contains(token, "\\") ||
		strings.Contains(token, "/") ||
		strings.HasPrefix(token, ".") ||
		(len(token) > 1 && token[1] == ':')
}

func matchByExactPath(token string, projects []AgyProject) []AgyProject {
	var matches []AgyProject
	absToken, err := filepath.Abs(token)
	if err != nil {
		absToken = token
	}
	for _, p := range projects {
		if isMatchingProjectPath(p.GetPath(), absToken) {
			matches = append(matches, p)
		}
	}
	return matches
}

func resolveDirectoryOrFolderTarget(token string, projects []AgyProject, isConfirm bool) ([]AgyProject, error) {
	if folderMatches := collectProjectsUnderFolder(token, projects); len(folderMatches) > 0 {
		return folderMatches, nil
	}
	info, err := os.Stat(token)
	if err != nil || !info.IsDir() {
		return nil, fmt.Errorf("target %q did not match any project sequence, ID, alias, or path", token)
	}
	return resolveDirectoryAdHocTarget(token, isConfirm)
}

func resolveDirectoryAdHocTarget(token string, isConfirm bool) ([]AgyProject, error) {
	absDir, err := filepath.Abs(token)
	if err == nil {
		token = absDir
	}
	if IsRestrictedSystemOrHomeDir(token) {
		return nil, apperror.NewSimple(fmt.Sprintf("restricted system or home directory cannot be recreated as an Antigravity project: %s", token), "E9101")
	}
	return validateAndBuildAdHoc(token, isConfirm)
}

func validateAndBuildAdHoc(token string, isConfirm bool) ([]AgyProject, error) {
	targetDir := resolveWorkingRepoRoot(token)
	if !isGitRepo(targetDir) {
		return nil, apperror.NewSimple(fmt.Sprintf("directory %q is neither an existing Antigravity project nor a valid Git repository; refusing to create ad-hoc project", targetDir), "E9102")
	}
	if !isConfirm {
		return nil, apperror.NewSimple(fmt.Sprintf("directory %q is not an existing registered Antigravity project; pass --confirm (-y) to confirm creating and recreating this workspace", targetDir), "E9103")
	}
	return []AgyProject{buildAdHocProject(targetDir)}, nil
}

// IsRestrictedSystemOrHomeDir reports whether dirPath is a drive root, user home, or OS system directory.
func IsRestrictedSystemOrHomeDir(dirPath string) bool {
	if strings.TrimSpace(dirPath) == "" {
		return true
	}
	clean := filepath.Clean(dirPath)
	if abs, err := filepath.Abs(clean); err == nil {
		clean = abs
	}
	if isVolumeOrRootDrive(clean) || isUserHomeOrParent(clean) {
		return true
	}
	return isOperatingSystemDir(clean)
}

func isVolumeOrRootDrive(clean string) bool {
	if clean == "/" || clean == "\\" {
		return true
	}
	if filepath.Dir(clean) == clean {
		return true
	}
	vol := filepath.VolumeName(clean)
	return vol != "" && (clean == vol || clean == vol+"\\" || clean == vol+"/")
}

func isUserHomeOrParent(clean string) bool {
	home, err := os.UserHomeDir()
	if err != nil || home == "" {
		return false
	}
	cleanHome := filepath.Clean(home)
	if strings.EqualFold(clean, cleanHome) {
		return true
	}
	parent := filepath.Dir(cleanHome)
	return parent != cleanHome && strings.EqualFold(clean, parent)
}

func isOperatingSystemDir(clean string) bool {
	candidates := collectSystemDirCandidates()
	for _, d := range candidates {
		if d != "" && strings.EqualFold(clean, filepath.Clean(d)) {
			return true
		}
	}
	return false
}

func collectSystemDirCandidates() []string {
	return []string{
		os.Getenv("WINDIR"),
		os.Getenv("SystemRoot"),
		os.Getenv("ProgramFiles"),
		os.Getenv("ProgramFiles(x86)"),
		os.Getenv("ProgramData"),
		"/etc", "/usr", "/bin", "/sbin", "/var", "/root", "/System", "/Library",
	}
}

func isGitRepo(dirPath string) bool {
	if strings.TrimSpace(dirPath) == "" {
		return false
	}
	dotGit := filepath.Join(dirPath, ".git")
	if _, err := os.Stat(dotGit); err == nil {
		return true
	}
	root, err := gitutil.RepoRoot(dirPath)
	return err == nil && len(root) > 0
}
