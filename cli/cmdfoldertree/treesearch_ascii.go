package cmdfoldertree

import (
	"path/filepath"
	"strings"
)

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
	indentLoop:
		for i < len(runes) {
			r := runes[i]
			switch r {
			case ' ':
				indent++
				i++
			case '\t':
				indent += 4
				i++
			case '│':
				// Tree vertical bar: usually counts as 4 chars (e.g. "│   ")
				indent += 4
				i++
				// consume trailing spaces if present
				for i < len(runes) && runes[i] == ' ' {
					i++
				}
			case '├', '└':
				// Branch connector: consume "├── " or "└── "
				indent += 4
				i++
				for i < len(runes) && (runes[i] == '─' || runes[i] == ' ') {
					i++
				}
				break indentLoop
			default:
				break indentLoop
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
