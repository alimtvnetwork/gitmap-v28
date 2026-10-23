package cmdfoldertree

import (
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

type yamlReportFile struct {
	Path      string `yaml:"path"`
	Filename  string `yaml:"filename"`
	Directory string `yaml:"directory"`
	Extension string `yaml:"extension"`
	SizeBytes int64  `yaml:"sizeBytes"`
	IsDir     bool   `yaml:"isDir"`
}

type yamlFolderReport struct {
	Root  string           `yaml:"root"`
	Files []yamlReportFile `yaml:"files"`
}

func parseYAMLTreeContent(content []byte) ([]TreeItem, error) {
	nowUnix := time.Now().Unix()

	if items, ok := tryParseYAMLExportDoc(content, nowUnix); ok {
		return items, nil
	}
	if items, ok := tryParseYAMLReport(content, nowUnix); ok {
		return items, nil
	}
	if items, ok := tryParseYAMLPaths(content, nowUnix); ok {
		return items, nil
	}
	if items, ok := tryParseYAMLGenericMap(content, nowUnix); ok {
		return items, nil
	}

	return parseTextTreeContent(content)
}

func tryParseYAMLExportDoc(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var exportDoc struct {
		RootPath   string          `yaml:"rootPath"`
		RootName   string          `yaml:"rootName"`
		TotalNodes int             `yaml:"totalNodes"`
		TotalDirs  int             `yaml:"totalDirs"`
		TotalFiles int             `yaml:"totalFiles"`
		Tree       *FolderTreeNode `yaml:"tree"`
	}
	if err := yaml.Unmarshal(content, &exportDoc); err != nil || exportDoc.Tree == nil {
		return nil, false
	}
	items := make([]TreeItem, 0, 64)
	walkExportTreeNode(exportDoc.Tree, exportDoc.RootPath, &items, nowUnix)
	return items, true
}

func tryParseYAMLReport(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var report yamlFolderReport
	if err := yaml.Unmarshal(content, &report); err != nil || len(report.Files) == 0 {
		return nil, false
	}
	items := make([]TreeItem, 0, len(report.Files))
	for _, f := range report.Files {
		items = append(items, buildTreeItemFromYamlReportFile(f, report.Root, nowUnix))
	}
	return items, true
}

func buildTreeItemFromYamlReportFile(f yamlReportFile, root string, nowUnix int64) TreeItem {
	relPath := normalizeSlash(f.Path)
	fileName := resolveFileName(f.Filename, relPath)
	dirPath := resolveDirPath(f.Directory, relPath)
	ext := resolveExt(f.Extension, fileName)
	depth := computePathDepth(relPath)

	return TreeItem{
		RootPath:  root,
		RelPath:   relPath,
		FileName:  fileName,
		DirPath:   dirPath,
		Extension: ext,
		SizeBytes: f.SizeBytes,
		IsDir:     f.IsDir,
		Depth:     depth,
		UpdatedAt: nowUnix,
	}
}

func tryParseYAMLPaths(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var rawPaths []string
	if err := yaml.Unmarshal(content, &rawPaths); err != nil || len(rawPaths) == 0 {
		return nil, false
	}
	items := make([]TreeItem, 0, len(rawPaths))
	for _, p := range rawPaths {
		clean := strings.TrimSpace(p)
		if clean == "" {
			continue
		}
		items = append(items, buildTreeItemFromPath(clean, "", 0, nowUnix))
	}
	return items, true
}

func tryParseYAMLGenericMap(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var genMap map[string]interface{}
	if err := yaml.Unmarshal(content, &genMap); err != nil {
		return nil, false
	}
	rawFiles, ok := genMap["files"].([]interface{})
	if !ok || len(rawFiles) == 0 {
		return nil, false
	}
	items := make([]TreeItem, 0, len(rawFiles))
	for _, rf := range rawFiles {
		if pathStr, ok := rf.(string); ok {
			items = append(items, buildTreeItemFromPath(pathStr, "", 0, nowUnix))
		}
	}
	if len(items) == 0 {
		return nil, false
	}
	return items, true
}
