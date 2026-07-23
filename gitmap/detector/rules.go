// Package detector — rules.go handles detection classification for C++ and C#.
package detector

import (
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// detectCpp checks if a file indicates a C++ project.
func detectCpp(name, dir, repoPath string, repoID int64, repoName string, results *[]DetectionResult) {
	indicator := ""
	if name == constants.IndicatorCMakeLists {
		indicator = constants.IndicatorCMakeLists
	}
	if name == constants.IndicatorMesonBuild {
		indicator = constants.IndicatorMesonBuild
	}
	if strings.HasSuffix(name, constants.ExtVcxproj) {
		indicator = name
	}
	if len(indicator) == 0 {
		return
	}
	if isDuplicate(dir, constants.ProjectKeyCpp, results) {
		return
	}
	addResult(dir, repoPath, repoID, repoName, constants.ProjectTypeCppID,
		constants.ProjectKeyCpp, filepath.Base(dir), indicator, results)
}

// detectCsharpFile handles .sln and standalone .csproj detection.
func detectCsharpFile(name, dir, repoPath string, repoID int64, repoName string, slnDirs map[string]bool, results *[]DetectionResult) {
	if strings.HasSuffix(name, constants.ExtSln) {
		detectCsharpSln(name, dir, repoPath, repoID, repoName, results)

		return
	}
	if strings.HasSuffix(name, constants.ExtCsproj) {
		detectCsharpStandalone(name, dir, repoPath, repoID, repoName, slnDirs, results)
	}
}

// detectCsharpSln creates a project entry for a .sln file.
func detectCsharpSln(name, dir, repoPath string, repoID int64, repoName string, results *[]DetectionResult) {
	if isDuplicate(dir, constants.ProjectKeyCsharp, results) {
		return
	}
	projName := strings.TrimSuffix(name, constants.ExtSln)
	result := buildBaseResult(dir, repoPath, repoID, repoName,
		constants.ProjectTypeCsharpID, constants.ProjectKeyCsharp, projName, name)
	meta := parseCsharpProject(dir, repoPath)
	result.Csharp = meta
	*results = append(*results, result)
}

// detectCsharpStandalone creates an entry for .csproj without parent .sln.
func detectCsharpStandalone(name, dir, repoPath string, repoID int64, repoName string, slnDirs map[string]bool, results *[]DetectionResult) {
	if isUnderSlnDir(dir, slnDirs) {
		return
	}
	if isDuplicate(dir, constants.ProjectKeyCsharp, results) {
		return
	}
	projName := strings.TrimSuffix(name, constants.ExtCsproj)
	addResult(dir, repoPath, repoID, repoName, constants.ProjectTypeCsharpID,
		constants.ProjectKeyCsharp, projName, name, results)
}

// isUnderSlnDir checks if dir is at or under any .sln directory.
func isUnderSlnDir(dir string, slnDirs map[string]bool) bool {
	for slnDir := range slnDirs {
		if strings.HasPrefix(dir, slnDir) {
			return true
		}
	}

	return false
}

// isDuplicate checks if a project of the given type already exists at dir.
func isDuplicate(dir, typeKey string, results *[]DetectionResult) bool {
	for _, r := range *results {
		if r.Project.AbsolutePath == dir && r.Project.ProjectType == typeKey {
			return true
		}
	}

	return false
}

// addResult creates and appends a basic DetectionResult.
func addResult(dir, repoPath string, repoID int64, repoName string, typeID int64, typeKey, projName, indicator string, results *[]DetectionResult) {
	result := buildBaseResult(dir, repoPath, repoID, repoName, typeID, typeKey, projName, indicator)
	*results = append(*results, result)
}

// buildBaseResult creates a DetectionResult with the project fields populated.
func buildBaseResult(dir, repoPath string, repoID int64, repoName string, typeID int64, typeKey, projName, indicator string) DetectionResult {
	relPath := buildRelativePath(dir, repoPath)

	return DetectionResult{
		Project: model.DetectedProject{
			RepoID:           repoID,
			RepoName:         repoName,
			ProjectTypeID:    typeID,
			ProjectType:      typeKey,
			ProjectName:      projName,
			AbsolutePath:     dir,
			RepoPath:         repoPath,
			RelativePath:     relPath,
			PrimaryIndicator: indicator,
		},
	}
}
