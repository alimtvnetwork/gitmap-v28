package cmdfoldertree

import (
	"path/filepath"
	"strings"
)

func resolveNodeIcon(node *FolderTreeNode) string {
	if node.IsGit {
		return "📦"
	}
	if node.IsDir {
		return "📁"
	}
	return resolveFileIcon(node.Name)
}

func resolveFileIcon(name string) string {
	ext := strings.ToLower(filepath.Ext(name))
	switch ext {
	case ".go", ".ts", ".js", ".py", ".rs", ".cs", ".php", ".c", ".cpp":
		return "💻"
	case ".json", ".yaml", ".yml", ".toml", ".ini", ".env", ".xml":
		return "📜"
	case ".md", ".txt", ".rst", ".doc", ".pdf":
		return "📝"
	case ".png", ".jpg", ".jpeg", ".svg", ".ico", ".webp", ".gif":
		return "🖼️"
	default:
		return "📄"
	}
}
