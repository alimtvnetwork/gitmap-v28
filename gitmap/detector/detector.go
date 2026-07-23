// Package detector walks repository trees and classifies project types.
package detector

import (
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v28/gitmap/model"
)

// DetectProjects scans a repo directory for all supported project types.
//
// Performance contract: the repo tree is walked exactly ONCE. The previous
// implementation did a full pre-walk to collect .sln directories followed
// by a second full walk for detection — doubling I/O on every repo. We now
// collect raw hits in a single pass and resolve .sln precedence in memory.
func DetectProjects(repoPath string, repoID int64, repoName string) []DetectionResult {
	hits := walkRepoOnce(repoPath)

	return classifyHits(repoPath, repoID, repoName, hits)
}

// DetectionResult holds a detected project and optional metadata.
type DetectionResult struct {
	Project model.DetectedProject
	GoMeta  *model.GoProjectMetadata
	Csharp  *model.CsharpProjectMetadata
}

// fileHit captures a single interesting file found during the walk.
type fileHit struct {
	path string // absolute file path
	dir  string // absolute parent directory
	name string // base name
}

// repoHits is the raw output of a single tree walk.
type repoHits struct {
	files   []fileHit
	slnDirs map[string]bool
}

// walkRepoOnce walks the repo tree exactly once, collecting every file
// that any classifier might care about plus the set of directories that
// directly contain a .sln file (used for C# precedence rules).
func walkRepoOnce(repoPath string) repoHits {
	out := repoHits{slnDirs: map[string]bool{}}
	_ = filepath.Walk(repoPath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() {
			if shouldExcludeDir(info.Name()) {
				return filepath.SkipDir
			}

			return nil
		}
		name := info.Name()
		if !isInterestingFile(name) {
			return nil
		}
		out.files = append(out.files, fileHit{path: path, dir: filepath.Dir(path), name: name})
		if strings.HasSuffix(name, constants.ExtSln) {
			out.slnDirs[filepath.Dir(path)] = true
		}

		return nil
	})

	return out
}

// isInterestingFile returns true when a file name matters to any classifier.
// Cheap O(1) prefix/suffix checks keep the hot path tight.
func isInterestingFile(name string) bool {
	if name == constants.IndicatorGoMod || name == constants.IndicatorPackageJSON {
		return true
	}
	if name == constants.IndicatorCMakeLists || name == constants.IndicatorMesonBuild {
		return true
	}
	if strings.HasSuffix(name, constants.ExtVcxproj) {
		return true
	}
	if strings.HasSuffix(name, constants.ExtSln) || strings.HasSuffix(name, constants.ExtCsproj) {
		return true
	}

	return false
}

// classifyHits feeds collected file hits through the per-language detectors.
func classifyHits(repoPath string, repoID int64, repoName string, hits repoHits) []DetectionResult {
	var results []DetectionResult
	for _, h := range hits.files {
		detectFile(h.path, repoPath, repoID, repoName, hits.slnDirs, &results)
	}

	return results
}

// detectFile checks a single file against all detection rules.
func detectFile(path, repoPath string, repoID int64, repoName string, slnDirs map[string]bool, results *[]DetectionResult) {
	name := filepath.Base(path)
	dir := filepath.Dir(path)

	if name == constants.IndicatorGoMod {
		detectGo(dir, repoPath, repoID, repoName, results)
	}
	if name == constants.IndicatorPackageJSON {
		detectNodeOrReact(dir, path, repoPath, repoID, repoName, results)
	}
	detectCpp(name, dir, repoPath, repoID, repoName, results)
	detectCsharpFile(name, dir, repoPath, repoID, repoName, slnDirs, results)
}

// shouldExcludeDir checks if a directory name should be skipped.
func shouldExcludeDir(name string) bool {
	if strings.HasPrefix(name, constants.CMakeBuildPfx) {
		return true
	}
	for _, excluded := range constants.ProjectExcludeDirs {
		if name == excluded {
			return true
		}
	}

	return false
}

// buildRelativePath returns the relative path from repo root.
func buildRelativePath(dir, repoPath string) string {
	rel, err := filepath.Rel(repoPath, dir)
	if err != nil {
		return "."
	}

	return rel
}
