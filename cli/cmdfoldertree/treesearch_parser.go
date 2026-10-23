package cmdfoldertree

import (
	"bytes"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

// TreeItem represents an indexed file or directory record from an imported tree.
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

// TreeFileRecord is an alias for TreeItem used by the import and query engine.
type TreeFileRecord = TreeItem

// ParseTreeFile detects the format of a tree export file and returns unified TreeItem records.
func ParseTreeFile(filePath string) ([]TreeItem, error) {
	if strings.TrimSpace(filePath) == "" {
		return nil, apperror.NewValidationError("input file path is required")
	}

	info, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return nil, apperror.NewNotFoundError("tree file not found: " + filePath)
	}
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to access tree file")
	}
	if info.IsDir() {
		return nil, apperror.NewValidationError("specified path is a directory, not a tree file: " + filePath)
	}

	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, apperror.WrapSimple(err, "failed to read tree file "+filePath)
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
		trimmed := bytes.TrimSpace(content)
		if items, ok := tryParseJSONTree(content, trimmed); ok {
			return items, nil
		}
		if items, ok := tryParseYAMLTree(content); ok {
			return items, nil
		}
		return parseTextTreeContent(content)
	}
}

func tryParseJSONTree(content, trimmed []byte) ([]TreeItem, bool) {
	if len(trimmed) == 0 || (trimmed[0] != '{' && trimmed[0] != '[') {
		return nil, false
	}
	items, err := parseJSONTreeContent(content)
	if err != nil || len(items) == 0 {
		return nil, false
	}
	return items, true
}

func tryParseYAMLTree(content []byte) ([]TreeItem, bool) {
	hasMarker := bytes.Contains(content, []byte("rootPath:")) ||
		bytes.Contains(content, []byte("tree:")) ||
		bytes.Contains(content, []byte("files:"))
	if !hasMarker {
		return nil, false
	}
	items, err := parseYAMLTreeContent(content)
	if err != nil || len(items) == 0 {
		return nil, false
	}
	return items, true
}

// parseTextTreeContent parses plain text: either flat path lists or indented ASCII tree.
func parseTextTreeContent(content []byte) ([]TreeItem, error) {
	lines := strings.Split(string(content), "\n")
	nowUnix := time.Now().Unix()

	isAsciiTree := false
	for _, line := range lines {
		if strings.Contains(line, "├──") || strings.Contains(line, "└──") || strings.Contains(line, "│") {
			isAsciiTree = true
			break
		}
	}

	if isAsciiTree {
		return parseIndentedAsciiTree(lines, nowUnix)
	}

	items := make([]TreeItem, 0, len(lines))
	for _, rawLine := range lines {
		line := strings.TrimSpace(rawLine)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		items = append(items, buildTreeItemFromPath(line, "", 0, nowUnix))
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
	dirPath := cleanDirPath(normalizeSlash(filepath.Dir(norm)))
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

func cleanDirPath(dir string) string {
	if dir == "." {
		return ""
	}
	return dir
}

func normalizeTreeItem(item *TreeItem, nowUnix int64) {
	item.RelPath = normalizeSlash(item.RelPath)
	if item.FileName == "" {
		item.FileName = filepath.Base(item.RelPath)
	}
	if item.DirPath == "" {
		item.DirPath = cleanDirPath(normalizeSlash(filepath.Dir(item.RelPath)))
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
