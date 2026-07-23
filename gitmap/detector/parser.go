// Package detector — parser.go parses package.json for Node/React classification.
package detector

import (
	"encoding/json"
	"os"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

// packageJSON represents the relevant fields of a package.json file.
type packageJSON struct {
	Name            string            `json:"name"`
	Dependencies    map[string]string `json:"dependencies"`
	DevDependencies map[string]string `json:"devDependencies"`
}

// detectNodeOrReact reads package.json and classifies as node or react.
func detectNodeOrReact(dir, pkgPath, repoPath string, repoID int64, repoName string, results *[]DetectionResult) {
	pkg, err := parsePackageJSON(pkgPath)
	if err != nil {
		return
	}
	if isDuplicate(dir, constants.ProjectKeyReact, results) {
		return
	}
	if isDuplicate(dir, constants.ProjectKeyNode, results) {
		return
	}

	projName := pkg.Name
	if len(projName) == 0 {
		projName = dir
	}

	if isReactProject(pkg) {
		addResult(dir, repoPath, repoID, repoName, constants.ProjectTypeReactID,
			constants.ProjectKeyReact, projName, constants.IndicatorPackageJSON, results)

		return
	}
	addResult(dir, repoPath, repoID, repoName, constants.ProjectTypeNodeID,
		constants.ProjectKeyNode, projName, constants.IndicatorPackageJSON, results)
}

// parsePackageJSON reads and decodes a package.json file.
func parsePackageJSON(path string) (*packageJSON, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	var pkg packageJSON
	err = json.Unmarshal(data, &pkg)

	return &pkg, err
}

// isReactProject checks if any React indicator dependency is present.
func isReactProject(pkg *packageJSON) bool {
	for _, dep := range constants.ReactIndicatorDeps {
		if hasDepKey(pkg.Dependencies, dep) {
			return true
		}
		if hasDepKey(pkg.DevDependencies, dep) {
			return true
		}
	}

	return false
}

// hasDepKey checks if a key exists in a dependency map.
func hasDepKey(deps map[string]string, key string) bool {
	_, exists := deps[key]

	return exists
}
