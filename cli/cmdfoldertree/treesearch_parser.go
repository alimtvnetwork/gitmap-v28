package cmdfoldertree

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

// TreeItem represents a normalized file or directory entry extracted from any tree format.
type TreeItem struct {
	Id        int64  `json:"id,omitempty"`
	RootPath  string `json:"rootPath,omitempty"`
	RelPath   string `json:"relPath"`
	FileName  string `json:"fileName"`
	DirPath   string `json:"dirPath"`
	Extension string `json:"extension"`
	SizeBytes int64  `json:"sizeBytes"`
	IsDir     bool   `json:"isDir"`
	Depth     int    `json:"depth"`
	UpdatedAt int64  `json:"updatedAt"`
}

// TreeFileRecord is an alias for TreeItem for backward compatibility with component specifications.
type TreeFileRecord = TreeItem

// ParseTreeFile reads and parses a tree file (.json, .yaml, .yml, .txt, or auto-detected).
func ParseTreeFile(filePath string) ([]TreeItem, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, fmt.Errorf("input file path is required")
	}

	info, err := os.Stat(filePath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, fmt.Errorf("tree file not found: %s", filePath)
		}
		return nil, fmt.Errorf("failed to access tree file: %w", err)
	}
	if info.IsDir() {
		return nil, fmt.Errorf("specified path is a directory, not a tree file: %s", filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read tree file %s: %w", filePath, err)
	}

	ext := strings.ToLower(filepath.Ext(filePath))
	switch ext {
	case ".json":
		return parseJSONTreeContent(content)
	case ".yaml", ".yml":
		return parseYAMLTreeContent(content)
	case ".txt", ".text", ".log":
		return parseTextTreeContent(content)
	default:
		// Sniff content format
		trimmed := bytes.TrimSpace(content)
		if len(trimmed) > 0 && (trimmed[0] == '{' || trimmed[0] == '[') {
			items, jsonErr := parseJSONTreeContent(content)
			if jsonErr == nil && len(items) > 0 {
				return items, nil
			}
		}
		if bytes.Contains(content, []byte("rootPath:")) || bytes.Contains(content, []byte("tree:")) || bytes.Contains(content, []byte("files:")) {
			items, yamlErr := parseYAMLTreeContent(content)
			if yamlErr == nil && len(items) > 0 {
				return items, nil
			}
		}
		return parseTextTreeContent(content)
	}
}

// parseJSONTreeContent handles FolderReport, FolderTreeExportDoc, []TreeItem, or []string.
func parseJSONTreeContent(content []byte) ([]TreeItem, error) {
	trimmed := bytes.TrimSpace(content)
	if len(trimmed) == 0 {
		return []TreeItem{}, nil
	}

	nowUnix := time.Now().Unix()

	// Try 1: Direct []TreeItem array
	var directItems []TreeItem
	if err := json.Unmarshal(content, &directItems); err == nil && len(directItems) > 0 {
		for i := range directItems {
			normalizeTreeItem(&directItems[i], nowUnix)
		}
		return directItems, nil
	}

	// Try 2: Raw []string array
	var rawPaths []string
	if err := json.Unmarshal(content, &rawPaths); err == nil && len(rawPaths) > 0 {
		items := make([]TreeItem, 0, len(rawPaths))
		for _, p := range rawPaths {
			p = strings.TrimSpace(p)
			if p == "" {
				continue
			}
			item := buildTreeItemFromPath(p, "", 0, nowUnix)
			items = append(items, item)
		}
		return items, nil
	}

	// Try 3: FolderReport schema (root, files: []*FileMeta)
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

	var report jsonFolderReport
	if err := json.Unmarshal(content, &report); err == nil && len(report.Files) > 0 {
		items := make([]TreeItem, 0, len(report.Files))
		for _, f := range report.Files {
			relPath := normalizeSlash(f.Path)
			fileName := f.Filename
			if fileName == "" {
				fileName = filepath.Base(relPath)
			}
			dirPath := normalizeSlash(f.Directory)
			if dirPath == "" || dirPath == "." {
				dirPath = normalizeSlash(filepath.Dir(relPath))
				if dirPath == "." {
					dirPath = ""
				}
			}
			ext := f.Extension
			if ext == "" {
				ext = filepath.Ext(fileName)
			}
			depth := computePathDepth(relPath)
			items = append(items, TreeItem{
				RootPath:  report.Root,
				RelPath:   relPath,
				FileName:  fileName,
				DirPath:   dirPath,
				Extension: ext,
				SizeBytes: f.SizeBytes,
				IsDir:     false,
				Depth:     depth,
				UpdatedAt: nowUnix,
			})
		}
		return items, nil
	}

	// Try 4: FolderTreeExportDoc schema (rootPath, rootName, tree: *FolderTreeNode)
	type jsonExportDoc struct {
		RootPath   string          `json:"rootPath"`
		RootName   string          `json:"rootName"`
		TotalNodes int             `json:"totalNodes"`
		TotalDirs  int             `json:"totalDirs"`
		TotalFiles int             `json:"totalFiles"`
		Tree       *FolderTreeNode `json:"tree"`
	}

	var exportDoc jsonExportDoc
	if err := json.Unmarshal(content, &exportDoc); err == nil && exportDoc.Tree != nil {
		items := make([]TreeItem, 0, 64)
		walkExportTreeNode(exportDoc.Tree, exportDoc.RootPath, &items, nowUnix)
		return items, nil
	}

	// Try 5: Generic object with "files" array of string or objects
	var genericMap map[string]interface{}
	if err := json.Unmarshal(content, &genericMap); err == nil {
		rootPath, _ := genericMap["rootPath"].(string)
		if rootPath == "" {
			rootPath, _ = genericMap["root"].(string)
		}
		if rawFiles, ok := genericMap["files"].([]interface{}); ok && len(rawFiles) > 0 {
			items := make([]TreeItem, 0, len(rawFiles))
			for _, rf := range rawFiles {
				switch val := rf.(type) {
				case string:
					items = append(items, buildTreeItemFromPath(val, rootPath, 0, nowUnix))
				case map[string]interface{}:
					rel, _ := val["relPath"].(string)
					if rel == "" {
						rel, _ = val["path"].(string)
					}
					if rel != "" {
						size := int64(0)
						if sz, ok := val["sizeBytes"].(float64); ok {
							size = int64(sz)
						}
						isDir, _ := val["isDir"].(bool)
						item := buildTreeItemFromPath(rel, rootPath, size, nowUnix)
						if isDir {
							item.IsDir = true
						}
						items = append(items, item)
					}
				}
			}
			if len(items) > 0 {
				return items, nil
			}
		}
	}

	return nil, fmt.Errorf("unable to recognize JSON tree format (expected FolderReport, FolderTreeExportDoc, or []TreeItem)")
}

// walkExportTreeNode recursively visits FolderTreeNode structures.
func walkExportTreeNode(node *FolderTreeNode, rootPath string, items *[]TreeItem, nowUnix int64) {
	if node == nil {
		return
	}
	relPath := normalizeSlash(node.RelPath)
	if relPath != "" {
		fileName := node.Name
		if fileName == "" {
			fileName = filepath.Base(relPath)
		}
		dirPath := normalizeSlash(filepath.Dir(relPath))
		if dirPath == "." {
			dirPath = ""
		}
		ext := filepath.Ext(fileName)
		depth := computePathDepth(relPath)
		*items = append(*items, TreeItem{
			RootPath:  rootPath,
			RelPath:   relPath,
			FileName:  fileName,
			DirPath:   dirPath,
			Extension: ext,
			SizeBytes: 0,
			IsDir:     node.IsDir,
			Depth:     depth,
			UpdatedAt: nowUnix,
		})
	}
	for _, child := range node.Children {
		walkExportTreeNode(child, rootPath, items, nowUnix)
	}
}

// parseYAMLTreeContent handles YAML formatted tree documents.
func parseYAMLTreeContent(content []byte) ([]TreeItem, error) {
	nowUnix := time.Now().Unix()

	// Try 1: FolderTreeExportDoc representation
	var exportDoc struct {
		RootPath   string          `yaml:"rootPath"`
		RootName   string          `yaml:"rootName"`
		TotalNodes int             `yaml:"totalNodes"`
		TotalDirs  int             `yaml:"totalDirs"`
		TotalFiles int             `yaml:"totalFiles"`
		Tree       *FolderTreeNode `yaml:"tree"`
	}
	if err := yaml.Unmarshal(content, &exportDoc); err == nil && exportDoc.Tree != nil {
		items := make([]TreeItem, 0, 64)
		walkExportTreeNode(exportDoc.Tree, exportDoc.RootPath, &items, nowUnix)
		return items, nil
	}

	// Try 2: FolderReport representation
	var report struct {
		Root  string `yaml:"root"`
		Files []struct {
			Path      string `yaml:"path"`
			Filename  string `yaml:"filename"`
			Directory string `yaml:"directory"`
			Extension string `yaml:"extension"`
			SizeBytes int64  `yaml:"sizeBytes"`
			IsDir     bool   `yaml:"isDir"`
		} `yaml:"files"`
	}
	if err := yaml.Unmarshal(content, &report); err == nil && len(report.Files) > 0 {
		items := make([]TreeItem, 0, len(report.Files))
		for _, f := range report.Files {
			relPath := normalizeSlash(f.Path)
			fileName := f.Filename
			if fileName == "" {
				fileName = filepath.Base(relPath)
			}
			dirPath := normalizeSlash(f.Directory)
			if dirPath == "" || dirPath == "." {
				dirPath = normalizeSlash(filepath.Dir(relPath))
				if dirPath == "." {
					dirPath = ""
				}
			}
			ext := f.Extension
			if ext == "" {
				ext = filepath.Ext(fileName)
			}
			depth := computePathDepth(relPath)
			items = append(items, TreeItem{
				RootPath:  report.Root,
				RelPath:   relPath,
				FileName:  fileName,
				DirPath:   dirPath,
				Extension: ext,
				SizeBytes: f.SizeBytes,
				IsDir:     f.IsDir,
				Depth:     depth,
				UpdatedAt: nowUnix,
			})
		}
		return items, nil
	}

	// Try 3: List of path strings in YAML
	var rawPaths []string
	if err := yaml.Unmarshal(content, &rawPaths); err == nil && len(rawPaths) > 0 {
		items := make([]TreeItem, 0, len(rawPaths))
		for _, p := range rawPaths {
			p = strings.TrimSpace(p)
			if p != "" {
				items = append(items, buildTreeItemFromPath(p, "", 0, nowUnix))
			}
		}
		return items, nil
	}

	// Try 4: Generic YAML map
	var genMap map[string]interface{}
	if err := yaml.Unmarshal(content, &genMap); err == nil {
		if rawFiles, ok := genMap["files"].([]interface{}); ok && len(rawFiles) > 0 {
			items := make([]TreeItem, 0, len(rawFiles))
			for _, rf := range rawFiles {
				if pathStr, ok := rf.(string); ok {
					items = append(items, buildTreeItemFromPath(pathStr, "", 0, nowUnix))
				}
			}
			if len(items) > 0 {
				return items, nil
			}
		}
	}

	// Fallback to text parsing
	return parseTextTreeContent(content)
}

// parseTextTreeContent parses plain text: either flat path lists or indented ASCII tree.
func parseTextTreeContent(content []byte) ([]TreeItem, error) {
	scanner := bufio.NewScanner(bytes.NewReader(content))
	lines := make([]string, 0, 128)
	hasTreeGlyphs := false

	for scanner.Scan() {
		raw := scanner.Text()
		trimmed := strings.TrimRight(raw, "\r\n")
		if strings.TrimSpace(trimmed) == "" {
			continue
		}
		// Skip comment banners like "=== tree ===" or "# tree"
		s := strings.TrimSpace(trimmed)
		if strings.HasPrefix(s, "===") || strings.HasPrefix(s, "---") || strings.HasPrefix(s, "###") {
			continue
		}
		if strings.ContainsAny(raw, "├└│") {
			hasTreeGlyphs = true
		}
		lines = append(lines, trimmed)
	}

	if len(lines) == 0 {
		return []TreeItem{}, nil
	}

	nowUnix := time.Now().Unix()

	if hasTreeGlyphs {
		return parseIndentedAsciiTree(lines, nowUnix)
	}

	// Check if lines are indented by spaces/tabs (hierarchy without box chars)
	hasIndentation := false
	for _, l := range lines {
		if strings.HasPrefix(l, "  ") || strings.HasPrefix(l, "\t") {
			hasIndentation = true
			break
		}
	}
	if hasIndentation {
		return parseIndentedAsciiTree(lines, nowUnix)
	}

	// Plain flat list
	items := make([]TreeItem, 0, len(lines))
	for _, line := range lines {
		line = strings.TrimSpace(line)
		if line == "" {
			continue
		}
		items = append(items, buildTreeItemFromPath(line, "", 0, nowUnix))
	}
	return items, nil
}

type dirStackEntry struct {
	depth int
	path  string
}

// parseIndentedAsciiTree reconstructs directory hierarchy from indented text.
func parseIndentedAsciiTree(lines []string, nowUnix int64) ([]TreeItem, error) {
	items := make([]TreeItem, 0, len(lines))
	dirStack := make([]dirStackEntry, 0, 16)

	for _, line := range lines {
		// Calculate indentation level
		indent := 0
		i := 0
		runes := []rune(line)
		for i < len(runes) {
			r := runes[i]
			if r == ' ' {
				indent++
				i++
			} else if r == '\t' {
				indent += 4
				i++
			} else if r == '│' {
				// Tree vertical bar: usually counts as 4 chars (e.g. "│   ")
				indent += 4
				i++
				// consume trailing spaces if present
				for i < len(runes) && runes[i] == ' ' {
					i++
				}
			} else if r == '├' || r == '└' {
				// Branch connector: consume "├── " or "└── "
				indent += 4
				i++
				for i < len(runes) && (runes[i] == '─' || runes[i] == ' ') {
					i++
				}
				break
			} else {
				break
			}
		}

		cleanName := strings.TrimSpace(string(runes[i:]))
		// Strip any remaining branch decorations
		cleanName = strings.TrimLeft(cleanName, "├─└│ \t")
		cleanName = strings.TrimSpace(cleanName)
		if cleanName == "" {
			continue
		}

		depth := indent / 4

		isDir := strings.HasSuffix(cleanName, "/") || strings.HasSuffix(cleanName, "\\")
		cleanName = strings.TrimRight(cleanName, "/\\")

		// Pop stack items that are deeper than or equal to current depth
		for len(dirStack) > 0 && dirStack[len(dirStack)-1].depth >= depth {
			dirStack = dirStack[:len(dirStack)-1]
		}

		var relPath string
		var dirPath string
		if len(dirStack) > 0 {
			dirPath = dirStack[len(dirStack)-1].path
			relPath = dirPath + "/" + cleanName
		} else {
			dirPath = ""
			relPath = cleanName
		}

		ext := ""
		if !isDir {
			ext = filepath.Ext(cleanName)
		}

		item := TreeItem{
			RelPath:   normalizeSlash(relPath),
			FileName:  cleanName,
			DirPath:   normalizeSlash(dirPath),
			Extension: ext,
			IsDir:     isDir,
			Depth:     computePathDepth(relPath),
			UpdatedAt: nowUnix,
		}
		items = append(items, item)

		if isDir {
			dirStack = append(dirStack, dirStackEntry{
				depth: depth,
				path:  relPath,
			})
		}
	}

	return items, nil
}

// buildTreeItemFromPath creates a TreeItem from a flat path string.
func buildTreeItemFromPath(rawPath string, rootPath string, sizeBytes int64, nowUnix int64) TreeItem {
	norm := normalizeSlash(rawPath)
	norm = strings.TrimPrefix(norm, "./")

	isDir := strings.HasSuffix(norm, "/")
	norm = strings.TrimRight(norm, "/")

	fileName := filepath.Base(norm)
	dirPath := normalizeSlash(filepath.Dir(norm))
	if dirPath == "." {
		dirPath = ""
	}
	ext := ""
	if !isDir {
		ext = filepath.Ext(fileName)
	}

	return TreeItem{
		RootPath:  rootPath,
		RelPath:   norm,
		FileName:  fileName,
		DirPath:   dirPath,
		Extension: ext,
		SizeBytes: sizeBytes,
		IsDir:     isDir,
		Depth:     computePathDepth(norm),
		UpdatedAt: nowUnix,
	}
}

func normalizeTreeItem(item *TreeItem, nowUnix int64) {
	item.RelPath = normalizeSlash(item.RelPath)
	if item.FileName == "" {
		item.FileName = filepath.Base(item.RelPath)
	}
	if item.DirPath == "" {
		dir := normalizeSlash(filepath.Dir(item.RelPath))
		if dir == "." {
			dir = ""
		}
		item.DirPath = dir
	}
	if item.Extension == "" && !item.IsDir {
		item.Extension = filepath.Ext(item.FileName)
	}
	if item.Depth == 0 {
		item.Depth = computePathDepth(item.RelPath)
	}
	if item.UpdatedAt == 0 {
		item.UpdatedAt = nowUnix
	}
}

func normalizeSlash(p string) string {
	return strings.ReplaceAll(filepath.ToSlash(strings.TrimSpace(p)), "\\", "/")
}

func computePathDepth(p string) int {
	norm := normalizeSlash(p)
	norm = strings.Trim(norm, "/")
	if norm == "" {
		return 0
	}
	return strings.Count(norm, "/")
}
