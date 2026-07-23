// Package detector — csharpparser.go parses C# project metadata.
package detector

import (
	"encoding/json"
	"encoding/xml"
	"os"
	"path/filepath"
	"strings"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/model"
)

// parseCsharpProject collects C# metadata from a solution directory.
func parseCsharpProject(dir, repoPath string) *model.CsharpProjectMetadata {
	meta := &model.CsharpProjectMetadata{}
	findSlnFile(dir, meta)
	findGlobalJSON(dir, meta)
	meta.ProjectFiles = findCsprojFiles(dir, repoPath)
	meta.KeyFiles = findKeyFiles(dir, repoPath)

	return meta
}

// findSlnFile locates the .sln file in the directory.
func findSlnFile(dir string, meta *model.CsharpProjectMetadata) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return
	}
	for _, entry := range entries {
		if strings.HasSuffix(entry.Name(), constants.ExtSln) {
			meta.SlnPath = filepath.Join(dir, entry.Name())
			meta.SlnName = entry.Name()

			return
		}
	}
}

// findGlobalJSON looks for global.json and parses the SDK version.
func findGlobalJSON(dir string, meta *model.CsharpProjectMetadata) {
	path := filepath.Join(dir, "global.json")
	if fileExists(path) {
		meta.GlobalJsonPath = path
		meta.SdkVersion = parseGlobalJSONSdk(path)
	}
}

// globalJSON represents the relevant fields of global.json.
type globalJSON struct {
	Sdk struct {
		Version string `json:"version"`
	} `json:"sdk"`
}

// parseGlobalJSONSdk extracts the SDK version from global.json.
func parseGlobalJSONSdk(path string) string {
	data, err := os.ReadFile(path)
	if err != nil {
		return ""
	}
	var g globalJSON
	err = json.Unmarshal(data, &g)
	if err != nil {
		return ""
	}

	return g.Sdk.Version
}

// csprojXML represents the relevant XML structure of a .csproj file.
type csprojXML struct {
	XMLName xml.Name `xml:"Project"`
	Sdk     string   `xml:"Sdk,attr"`
	Groups  []struct {
		TargetFramework string `xml:"TargetFramework"`
		OutputType      string `xml:"OutputType"`
	} `xml:"PropertyGroup"`
}

// findCsprojFiles walks the tree for .csproj and .fsproj files.
func findCsprojFiles(dir, _ string) []model.CsharpProjectFile {
	var files []model.CsharpProjectFile
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && shouldExcludeDir(info.Name()) {
			return filepath.SkipDir
		}
		if isCsprojFile(info.Name()) {
			f := parseCsprojFile(path, dir)
			files = append(files, f)
		}

		return nil
	})

	return files
}

// isCsprojFile checks if a filename is a C# project file.
func isCsprojFile(name string) bool {
	if strings.HasSuffix(name, constants.ExtCsproj) {
		return true
	}

	return strings.HasSuffix(name, constants.ExtFsproj)
}

// parseCsprojFile extracts metadata from a .csproj XML file.
func parseCsprojFile(path, baseDir string) model.CsharpProjectFile {
	rel := buildRelativePath(filepath.Dir(path), baseDir)
	name := filepath.Base(path)
	relPath := filepath.Join(rel, name)
	projName := strings.TrimSuffix(name, filepath.Ext(name))
	f := model.CsharpProjectFile{
		FilePath:     path,
		RelativePath: relPath,
		FileName:     name,
		ProjectName:  projName,
	}
	parseCsprojXML(path, &f)

	return f
}

// parseCsprojXML reads the XML and populates framework/output/SDK fields.
func parseCsprojXML(path string, f *model.CsharpProjectFile) {
	data, err := os.ReadFile(path)
	if err != nil {
		return
	}
	var proj csprojXML
	err = xml.Unmarshal(data, &proj)
	if err != nil {
		return
	}
	f.Sdk = proj.Sdk
	for _, pg := range proj.Groups {
		if len(pg.TargetFramework) > 0 {
			f.TargetFramework = pg.TargetFramework
		}
		if len(pg.OutputType) > 0 {
			f.OutputType = pg.OutputType
		}
	}
}

// findKeyFiles collects known C# key files in the project tree.
func findKeyFiles(dir, _ string) []model.CsharpKeyFile {
	var files []model.CsharpKeyFile
	_ = filepath.Walk(dir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return nil
		}
		if info.IsDir() && shouldExcludeDir(info.Name()) {
			return filepath.SkipDir
		}
		if isKeyFile(info.Name()) {
			rel := buildRelativePath(filepath.Dir(path), dir)
			relPath := filepath.Join(rel, info.Name())
			files = append(files, model.CsharpKeyFile{
				FileType:     info.Name(),
				FilePath:     path,
				RelativePath: relPath,
			})
		}

		return nil
	})

	return files
}

// isKeyFile checks if a filename matches any C# key file pattern.
func isKeyFile(name string) bool {
	for _, pattern := range constants.CsharpKeyFilePatterns {
		if name == pattern {
			return true
		}
	}
	if strings.HasSuffix(name, ".props") {
		return true
	}

	return strings.HasSuffix(name, ".targets")
}
