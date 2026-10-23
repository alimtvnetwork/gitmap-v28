package cmdfoldertree

import (
	"bytes"
	"encoding/json"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

type jsonFolderReportFile struct {
	Path          string `json:"path"`
	Filename      string `json:"filename"`
	Directory     string `json:"directory"`
	Extension     string `json:"extension"`
	SizeBytes     int64  `json:"sizeBytes"`
	SizeFormatted string `json:"sizeFormatted"`
	LinesOfCode   int    `json:"linesOfCode"`
	IsBinary      bool   `json:"isBinary"`
	Sequence      int    `json:"sequence"`
}

type jsonFolderReport struct {
	Root               string                 `json:"root"`
	TotalFiles         int                    `json:"totalFiles"`
	TotalLines         int                    `json:"totalLines"`
	TotalSizeBytes     int64                  `json:"totalSizeBytes"`
	TotalSizeFormatted string                 `json:"totalSizeFormatted"`
	Files              []jsonFolderReportFile `json:"files"`
}

type jsonExportDoc struct {
	RootPath   string          `json:"rootPath"`
	RootName   string          `json:"rootName"`
	TotalNodes int             `json:"totalNodes"`
	TotalDirs  int             `json:"totalDirs"`
	TotalFiles int             `json:"totalFiles"`
	Tree       *FolderTreeNode `json:"tree"`
}

func parseJSONTreeContent(content []byte) ([]TreeItem, error) {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return []TreeItem{}, nil
	}

	nowUnix := time.Now().Unix()

	if items, ok := tryParseDirectTreeItems(content, nowUnix); ok {
		return items, nil
	}
	if items, ok := tryParseRawPaths(content, nowUnix); ok {
		return items, nil
	}
	if items, ok := tryParseFolderReport(content, nowUnix); ok {
		return items, nil
	}
	if items, ok := tryParseExportDoc(content, nowUnix); ok {
		return items, nil
	}
	if items, ok := tryParseGenericMap(content, nowUnix); ok {
		return items, nil
	}

	return nil, apperror.NewValidationError("unable to recognize JSON tree format (expected FolderReport, FolderTreeExportDoc, or []TreeItem)")
}

func tryParseDirectTreeItems(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var directItems []TreeItem
	if err := json.Unmarshal(content, &directItems); err != nil || len(directItems) == 0 {
		return nil, false
	}
	for i := range directItems {
		normalizeTreeItem(&directItems[i], nowUnix)
	}
	return directItems, true
}

func tryParseRawPaths(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var rawPaths []string
	if err := json.Unmarshal(content, &rawPaths); err != nil || len(rawPaths) == 0 {
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

func tryParseFolderReport(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var report jsonFolderReport
	if err := json.Unmarshal(content, &report); err != nil || len(report.Files) == 0 {
		return nil, false
	}
	items := make([]TreeItem, 0, len(report.Files))
	for _, f := range report.Files {
		items = append(items, buildTreeItemFromFolderReportFile(f, report.Root, nowUnix))
	}
	return items, true
}

func buildTreeItemFromFolderReportFile(f jsonFolderReportFile, root string, nowUnix int64) TreeItem {
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
		IsDir:     false,
		Depth:     depth,
		UpdatedAt: nowUnix,
	}
}

func resolveFileName(fn, rel string) string {
	if fn != "" {
		return fn
	}
	return filepath.Base(rel)
}

func resolveDirPath(dir, rel string) string {
	clean := normalizeSlash(dir)
	if clean != "" && clean != "." {
		return clean
	}
	parent := normalizeSlash(filepath.Dir(rel))
	if parent == "." {
		return ""
	}
	return parent
}

func resolveExt(ext, fileName string) string {
	if ext != "" {
		return ext
	}
	return filepath.Ext(fileName)
}

func tryParseExportDoc(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var exportDoc jsonExportDoc
	if err := json.Unmarshal(content, &exportDoc); err != nil || exportDoc.Tree == nil {
		return nil, false
	}
	items := make([]TreeItem, 0, 64)
	walkExportTreeNode(exportDoc.Tree, exportDoc.RootPath, &items, nowUnix)
	return items, true
}

func buildTreeItemFromExportNode(node *FolderTreeNode, rootPath string, nowUnix int64) TreeItem {
	relPath := normalizeSlash(node.RelPath)
	fileName := resolveFileName(node.Name, relPath)
	dirPath := cleanDirPath(normalizeSlash(filepath.Dir(relPath)))
	ext := filepath.Ext(fileName)
	depth := computePathDepth(relPath)

	return TreeItem{
		RootPath:  rootPath,
		RelPath:   relPath,
		FileName:  fileName,
		DirPath:   dirPath,
		Extension: ext,
		SizeBytes: 0,
		IsDir:     node.IsDir,
		Depth:     depth,
		UpdatedAt: nowUnix,
	}
}

// walkExportTreeNode recursively visits FolderTreeNode structures.
func walkExportTreeNode(node *FolderTreeNode, rootPath string, items *[]TreeItem, nowUnix int64) {
	if node == nil {
		return
	}
	if normalizeSlash(node.RelPath) != "" {
		*items = append(*items, buildTreeItemFromExportNode(node, rootPath, nowUnix))
	}
	for _, child := range node.Children {
		walkExportTreeNode(child, rootPath, items, nowUnix)
	}
}

func tryParseGenericMap(content []byte, nowUnix int64) ([]TreeItem, bool) {
	var genericMap map[string]interface{}
	if err := json.Unmarshal(content, &genericMap); err != nil {
		return nil, false
	}
	rawFiles, ok := genericMap["files"].([]interface{})
	if !ok || len(rawFiles) == 0 {
		return nil, false
	}
	rootPath := resolveRootPath(genericMap)
	items := make([]TreeItem, 0, len(rawFiles))
	for _, rf := range rawFiles {
		if item, ok := buildGenericMapItem(rf, rootPath, nowUnix); ok {
			items = append(items, item)
		}
	}
	if len(items) == 0 {
		return nil, false
	}
	return items, true
}

func resolveRootPath(m map[string]interface{}) string {
	if rootPath, _ := m["rootPath"].(string); rootPath != "" {
		return rootPath
	}
	root, _ := m["root"].(string)
	return root
}

func buildGenericMapItem(rf interface{}, rootPath string, nowUnix int64) (TreeItem, bool) {
	switch val := rf.(type) {
	case string:
		return buildTreeItemFromPath(val, rootPath, 0, nowUnix), true
	case map[string]interface{}:
		rel, _ := val["relPath"].(string)
		if rel == "" {
			rel, _ = val["path"].(string)
		}
		if rel == "" {
			return TreeItem{}, false
		}
		size := int64(0)
		if sz, ok := val["sizeBytes"].(float64); ok {
			size = int64(sz)
		}
		isDir, _ := val["isDir"].(bool)
		item := buildTreeItemFromPath(rel, rootPath, size, nowUnix)
		item.IsDir = isDir
		return item, true
	default:
		return TreeItem{}, false
	}
}
