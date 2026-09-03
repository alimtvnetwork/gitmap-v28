// Package cmd — agy_pin_projects_resolve.go resolves target strings to Antigravity projects.
package cmd

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func resolveTargetAgyProject(target string) (*AgyProject, *apperror.AppError) {
	if target == "" || target == "." {
		return resolveCurrentDirProject()
	}

	dirPath, pathErr := getProjectsDirPath()

	if pathErr != nil {
		return nil, apperror.WrapSimple(pathErr, "resolveTargetAgyProject.path")
	}

	projects, loadErr := loadAllAgyProjects(dirPath)

	if loadErr != nil {
		return nil, apperror.WrapSimple(loadErr, "resolveTargetAgyProject.load")
	}

	if match := matchAgyProject(projects, target); match != nil {
		return match, nil
	}

	return resolveFallbackLocalDir(target)
}

func resolveCurrentDirProject() (*AgyProject, *apperror.AppError) {
	cwd, cwdErr := os.Getwd()

	if cwdErr != nil {
		return nil, apperror.WrapSimple(cwdErr, "resolveCurrentDirProject.cwd")
	}

	return resolveTargetAgyProject(cwd)
}

func matchAgyProject(projects []AgyProject, target string) *AgyProject {
	low := strings.ToLower(target)
	cleanPath := filepath.Clean(target)

	for i := range projects {
		p := &projects[i]
		if p.ID == target || strings.HasPrefix(p.ID, target) {
			return p
		}

		if strings.ToLower(p.Name) == low || filepath.Clean(p.GetPath()) == cleanPath {
			return p
		}
	}

	return nil
}

func resolveFallbackLocalDir(target string) (*AgyProject, *apperror.AppError) {
	abs, absErr := filepath.Abs(target)

	if absErr != nil {
		return nil, apperror.WrapSimple(absErr, "resolveFallbackLocalDir.abs")
	}

	stat, statErr := os.Stat(abs)

	if statErr != nil || !stat.IsDir() {
		return nil, apperror.NewSimple(fmt.Sprintf("antigravity project or directory %q not found", target), "E9000")
	}

	return &AgyProject{
		ID:   filepath.Base(abs),
		Name: filepath.Base(abs),
		ProjectResources: &AgyProjectResources{
			Resources: []AgyResource{
				{
					GitFolder: &AgyGitFolder{
						FolderURI:     "file:///" + filepath.ToSlash(abs),
						DefaultBranch: "main",
					},
				},
			},
		},
	}, nil
}
