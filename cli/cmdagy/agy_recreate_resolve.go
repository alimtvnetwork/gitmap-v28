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
	if len(args) == 0 || (len(args) == 1 && args[0] == ".") {
		return resolveCurrentDirTarget(projects)
	}
	return resolveSpecifiedTargets(args, projects)
}

func resolveCurrentDirTarget(projects []AgyProject) ([]AgyProject, error) {
	cwd, err := os.Getwd()
	if err != nil {
		return nil, apperror.WrapSimple(err, "get current working directory")
	}
	targetDir := resolveWorkingRepoRoot(cwd)
	for _, p := range projects {
		if isMatchingProjectPath(p.GetPath(), targetDir) {
			return []AgyProject{p}, nil
		}
	}
	return []AgyProject{buildAdHocProject(targetDir)}, nil
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

func resolveSpecifiedTargets(args []string, projects []AgyProject) ([]AgyProject, error) {
	tokens := expandCommaSeparatedTokens(args)
	var results []AgyProject
	seen := make(map[string]bool)

	for _, token := range tokens {
		matched, err := resolveSingleRecreateToken(token, projects)
		if err != nil {
			return nil, err
		}
		results = appendUniqueProjects(results, matched, seen)
	}
	return results, nil
}

func resolveSingleRecreateToken(token string, projects []AgyProject) ([]AgyProject, error) {
	if seqMatch, isSeq := matchBySequence(token, projects); isSeq {
		return []AgyProject{seqMatch}, nil
	}
	if idMatches := matchByID(token, projects); len(idMatches) > 0 {
		return idMatches, nil
	}
	if isLikelyPath(token) {
		if pathMatches := matchByExactPath(token, projects); len(pathMatches) > 0 {
			return pathMatches, nil
		}
		if folderMatches := collectProjectsUnderFolder(token, projects); len(folderMatches) > 0 {
			return folderMatches, nil
		}
		if info, err := os.Stat(token); err == nil && info.IsDir() {
			targetDir := resolveWorkingRepoRoot(token)
			return []AgyProject{buildAdHocProject(targetDir)}, nil
		}
	}
	if slugMatches := matchBySlug(token, projects); len(slugMatches) > 0 {
		return slugMatches, nil
	}
	return resolveDirectoryOrFolderTarget(token, projects)
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

func resolveDirectoryOrFolderTarget(token string, projects []AgyProject) ([]AgyProject, error) {
	if folderMatches := collectProjectsUnderFolder(token, projects); len(folderMatches) > 0 {
		return folderMatches, nil
	}
	if info, err := os.Stat(token); err == nil && info.IsDir() {
		targetDir := resolveWorkingRepoRoot(token)
		return []AgyProject{buildAdHocProject(targetDir)}, nil
	}
	return nil, fmt.Errorf("target %q did not match any project sequence, ID, alias, or path", token)
}
