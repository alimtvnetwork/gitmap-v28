// Package detector — goparser.go parses Go metadata and finds runnables.
package detector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// detectGo creates a Go project entry with metadata from go.mod.
func detectGo(dir, repoPath string, repoID int64, repoName string, results *[]DetectionResult) {
	if isDuplicate(dir, constants.ProjectKeyGo, results) {
		return
	}
	modPath := filepath.Join(dir, constants.IndicatorGoMod)
	meta := parseGoMetadata(dir, modPath)
	projName := meta.ModuleName
	if len(projName) == 0 {
		projName = filepath.Base(dir)
	}

	result := buildBaseResult(dir, repoPath, repoID, repoName,
		constants.ProjectTypeGoID, constants.ProjectKeyGo, projName, constants.IndicatorGoMod)
	result.GoMeta = meta
	*results = append(*results, result)
}

// parseGoMetadata extracts module name, version, and paths from go.mod.
func parseGoMetadata(dir, modPath string) *model.GoProjectMetadata {
	meta := &model.GoProjectMetadata{
		GoModPath: modPath,
	}
	sumPath := filepath.Join(dir, constants.GoSumFile)
	if fileExists(sumPath) {
		meta.GoSumPath = sumPath
	}
	parseGoModContent(modPath, meta)
	meta.Runnables = findGoRunnables(dir)

	return meta
}

// parseGoModContent reads go.mod and extracts module and go directives.
func parseGoModContent(modPath string, meta *model.GoProjectMetadata) {
	data, err := os.ReadFile(modPath)
	if err != nil {
		return
	}
	for _, line := range strings.Split(string(data), "\n") {
		trimmed := strings.TrimSpace(line)
		if strings.HasPrefix(trimmed, "module ") {
			meta.ModuleName = strings.TrimSpace(strings.TrimPrefix(trimmed, "module "))
		}
		if strings.HasPrefix(trimmed, "go ") {
			meta.GoVersion = strings.TrimSpace(strings.TrimPrefix(trimmed, "go "))
		}
	}
}

// findGoRunnables scans cmd/ subdirectories and root for main.go files.
func findGoRunnables(projectDir string) []model.GoRunnableFile {
	var runnables []model.GoRunnableFile
	runnables = findCmdRunnables(projectDir, runnables)
	runnables = findRootRunnable(projectDir, runnables)

	return runnables
}

// findCmdRunnables scans the cmd/ directory for main.go files.
func findCmdRunnables(projectDir string, runnables []model.GoRunnableFile) []model.GoRunnableFile {
	cmdDir := filepath.Join(projectDir, constants.GoCmdDir)
	entries, err := os.ReadDir(cmdDir)
	if err != nil {
		return runnables
	}
	for _, entry := range entries {
		if entry.IsDir() {
			runnables = checkCmdSubdir(cmdDir, entry.Name(), projectDir, runnables)
		}
	}

	return runnables
}

// checkCmdSubdir looks for main.go in a cmd subdirectory.
func checkCmdSubdir(cmdDir, subName, projectDir string, runnables []model.GoRunnableFile) []model.GoRunnableFile {
	mainPath := filepath.Join(cmdDir, subName, constants.GoMainFile)
	if fileExists(mainPath) {
		rel := buildRelativePath(filepath.Dir(mainPath), projectDir)

		return append(runnables, model.GoRunnableFile{
			RunnableName: subName,
			FilePath:     mainPath,
			RelativePath: filepath.Join(rel, constants.GoMainFile),
		})
	}
	nestedPath := filepath.Join(cmdDir, subName, "main", constants.GoMainFile)
	if fileExists(nestedPath) {
		rel := buildRelativePath(filepath.Dir(nestedPath), projectDir)

		return append(runnables, model.GoRunnableFile{
			RunnableName: subName,
			FilePath:     nestedPath,
			RelativePath: filepath.Join(rel, constants.GoMainFile),
		})
	}

	return runnables
}

// findRootRunnable checks for main.go at the project root.
func findRootRunnable(projectDir string, runnables []model.GoRunnableFile) []model.GoRunnableFile {
	rootMain := filepath.Join(projectDir, constants.GoMainFile)
	if fileExists(rootMain) {
		return append(runnables, model.GoRunnableFile{
			RunnableName: filepath.Base(projectDir),
			FilePath:     rootMain,
			RelativePath: constants.GoMainFile,
		})
	}

	return runnables
}

// fileExists returns true if the path exists and is a regular file.
func fileExists(path string) bool {
	info, err := os.Stat(path)
	if err != nil {
		return false
	}

	return info.Mode().IsRegular()
}
