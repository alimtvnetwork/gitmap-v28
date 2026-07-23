// Package cmd — scanprojectoutput.go writes project-specific JSON files.
package cmd

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"

	"github.com/alimtvnetwork/gitmap-v27/gitmap/constants"
	"github.com/alimtvnetwork/gitmap-v27/gitmap/detector"
)

// projectTypeJSONMap maps project type keys to JSON filenames.
var projectTypeJSONMap = map[string]string{
	constants.ProjectKeyGo:     constants.JSONFileGoProjects,
	constants.ProjectKeyNode:   constants.JSONFileNodeProjects,
	constants.ProjectKeyReact:  constants.JSONFileReactProjects,
	constants.ProjectKeyCpp:    constants.JSONFileCppProjects,
	constants.ProjectKeyCsharp: constants.JSONFileCsharpProjects,
}

// writeProjectJSONFiles writes per-type JSON files to the output directory.
func writeProjectJSONFiles(results []detector.DetectionResult, outputDir string) {
	grouped := groupByType(results)
	for typeKey, items := range grouped {
		filename, exists := projectTypeJSONMap[typeKey]
		if exists {
			sortResults(items)
			writeProjectJSON(items, outputDir, filename)
		}
	}
}

// groupByType groups detection results by project type key.
func groupByType(results []detector.DetectionResult) map[string][]detector.DetectionResult {
	grouped := make(map[string][]detector.DetectionResult)
	for _, r := range results {
		key := r.Project.ProjectType
		grouped[key] = append(grouped[key], r)
	}

	return grouped
}

// sortResults sorts results by repo name then relative path.
func sortResults(results []detector.DetectionResult) {
	sort.Slice(results, func(i, j int) bool {
		if results[i].Project.RepoName == results[j].Project.RepoName {
			return results[i].Project.RelativePath < results[j].Project.RelativePath
		}

		return results[i].Project.RepoName < results[j].Project.RepoName
	})
}

// writeProjectJSON writes a single project type JSON file.
func writeProjectJSON(results []detector.DetectionResult, outputDir, filename string) {
	records := buildJSONRecords(results)
	path := filepath.Join(outputDir, filename)
	file, err := createOutputFile(path)
	if err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProjectJSONWrite, filename, err)

		return
	}
	defer file.Close()
	encoder := json.NewEncoder(file)
	encoder.SetIndent("", "  ")
	if err := encoder.Encode(records); err != nil {
		fmt.Fprintf(os.Stderr, constants.ErrProjectJSONWrite, filename, err)

		return
	}
	fmt.Printf(constants.MsgProjectJSONWritten, filename, len(records))
}

// buildJSONRecords converts DetectionResults to JSON-ready structures.
func buildJSONRecords(results []detector.DetectionResult) []interface{} {
	records := make([]interface{}, 0, len(results))
	for _, r := range results {
		records = append(records, buildSingleRecord(r))
	}

	return records
}

// buildSingleRecord creates the appropriate record with optional metadata.
func buildSingleRecord(r detector.DetectionResult) interface{} {
	if r.GoMeta != nil {
		return struct {
			detector.DetectionResult
		}{r}
	}
	if r.Csharp != nil {
		return struct {
			detector.DetectionResult
		}{r}
	}

	return r
}
