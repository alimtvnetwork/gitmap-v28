package cmdfoldertree

import (
	"encoding/json"
	"fmt"
	"strings"

	"gopkg.in/yaml.v3"
)

func parseImportPayload(data []byte, filename string) (*FolderTreeNode, error) {
	if isYAMLFile(filename) {
		return parseYAMLPayload(data)
	}
	if isJSONFile(filename) {
		return parseJSONPayload(data)
	}
	node, err := parseJSONPayload(data)
	if err == nil {
		return node, nil
	}
	node, err = parseYAMLPayload(data)
	if err == nil {
		return node, nil
	}
	return parseTextPathsPayload(data)
}

func isYAMLFile(filename string) bool {
	low := strings.ToLower(filename)
	return strings.HasSuffix(low, ".yaml") || strings.HasSuffix(low, ".yml")
}

func isJSONFile(filename string) bool {
	low := strings.ToLower(filename)
	return strings.HasSuffix(low, ".json")
}

func parseJSONPayload(data []byte) (*FolderTreeNode, error) {
	var doc FolderTreeExportDoc
	if err := json.Unmarshal(data, &doc); err == nil && doc.Tree != nil {
		return doc.Tree, nil
	}
	var node FolderTreeNode
	if err := json.Unmarshal(data, &node); err == nil && node.Name != "" {
		return &node, nil
	}
	var pathList []string
	if err := json.Unmarshal(data, &pathList); err == nil && len(pathList) > 0 {
		return buildTreeFromPathList(pathList), nil
	}
	return nil, fmt.Errorf("parse json import: invalid format")
}

func parseYAMLPayload(data []byte) (*FolderTreeNode, error) {
	var doc FolderTreeExportDoc
	if err := yaml.Unmarshal(data, &doc); err == nil && doc.Tree != nil {
		return doc.Tree, nil
	}
	var node FolderTreeNode
	if err := yaml.Unmarshal(data, &node); err == nil && node.Name != "" {
		return &node, nil
	}
	var pathList []string
	if err := yaml.Unmarshal(data, &pathList); err == nil && len(pathList) > 0 {
		return buildTreeFromPathList(pathList), nil
	}
	return nil, fmt.Errorf("parse yaml import: invalid format")
}
