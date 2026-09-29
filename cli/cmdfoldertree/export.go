package cmdfoldertree

import (
	"encoding/json"
	"os"
	"strings"

	"gopkg.in/yaml.v3"
)

// ExportFolderTree exports the scanned hierarchy in the requested format.
func ExportFolderTree(root *FolderTreeNode, opts FolderTreeOptions) (string, error) {
	doc := buildExportDoc(root)
	content, err := renderFormatContent(doc, root, opts)
	if err != nil {
		return "", err
	}
	if opts.OutputFile != "" {
		err := os.WriteFile(opts.OutputFile, []byte(content), 0644)
		return content, err
	}
	return content, nil
}

func renderFormatContent(doc FolderTreeExportDoc, root *FolderTreeNode, opts FolderTreeOptions) (string, error) {
	format := strings.ToLower(opts.Format)
	switch format {
	case "yaml", "yml":
		return exportYAML(doc)
	case "tree", "text":
		return RenderTree(root, opts), nil
	case "preview", "gap":
		return RenderPreview(root, opts), nil
	default:
		return exportJSON(doc)
	}
}

func buildExportDoc(root *FolderTreeNode) FolderTreeExportDoc {
	total, dirs, files := 0, 0, 0
	countNodes(root, &total, &dirs, &files)
	return FolderTreeExportDoc{
		RootPath:   root.Path,
		RootName:   root.Name,
		TotalNodes: total,
		TotalDirs:  dirs,
		TotalFiles: files,
		Tree:       root,
	}
}

func countNodes(node *FolderTreeNode, total, dirs, files *int) {
	*total++
	if node.IsDir {
		*dirs++
	} else {
		*files++
	}
	for _, child := range node.Children {
		countNodes(child, total, dirs, files)
	}
}

func exportJSON(doc FolderTreeExportDoc) (string, error) {
	b, err := json.MarshalIndent(doc, "", "  ")
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func exportYAML(doc FolderTreeExportDoc) (string, error) {
	b, err := yaml.Marshal(doc)
	if err != nil {
		return "", err
	}
	return string(b), nil
}
