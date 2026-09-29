package cmdfoldertree

import (
	"fmt"
	"os"
	"path/filepath"
)

// ImportFolderTree reads JSON/YAML/text and recreates directory and placeholder files on disk.
func ImportFolderTree(inputFile, targetDir string, opts FolderTreeOptions) (*ImportSummary, error) {
	data, err := os.ReadFile(inputFile)
	if err != nil {
		return nil, fmt.Errorf("read input file: %w", err)
	}
	rootNode, err := parseImportPayload(data, inputFile)
	if err != nil {
		return nil, err
	}
	absTarget, err := filepath.Abs(targetDir)
	if err != nil {
		return nil, fmt.Errorf("resolve target directory: %w", err)
	}
	summary := &ImportSummary{TargetDir: absTarget}
	rootNode.RelPath = "."
	err = materializeNode(rootNode, absTarget, opts, summary)
	return summary, err
}

func materializeNode(node *FolderTreeNode, currentDest string, opts FolderTreeOptions, summary *ImportSummary) error {
	destPath := filepath.Join(currentDest, node.Name)
	if node.RelPath == "." || node.Name == "." {
		destPath = currentDest
	}
	if node.IsDir {
		return materializeDirNode(node, destPath, opts, summary)
	}
	return materializeFileNode(destPath, opts, summary)
}

func materializeDirNode(node *FolderTreeNode, destPath string, opts FolderTreeOptions, summary *ImportSummary) error {
	if err := createDirIfNotDryRun(destPath, opts.DryRun); err != nil {
		return err
	}
	summary.DirsCreated++
	summary.CreatedPaths = append(summary.CreatedPaths, destPath)
	for _, child := range node.Children {
		if err := materializeNode(child, destPath, opts, summary); err != nil {
			return err
		}
	}
	return nil
}

func createDirIfNotDryRun(path string, dryRun bool) error {
	if dryRun {
		return nil
	}
	if err := os.MkdirAll(path, 0755); err != nil {
		return fmt.Errorf("create dir %s: %w", path, err)
	}
	return nil
}

func materializeFileNode(destPath string, opts FolderTreeOptions, summary *ImportSummary) error {
	if opts.DirsOnly {
		return nil
	}
	if err := createFileIfNotDryRun(destPath, opts.DryRun); err != nil {
		return err
	}
	summary.FilesCreated++
	summary.CreatedPaths = append(summary.CreatedPaths, destPath)
	return nil
}

func createFileIfNotDryRun(path string, dryRun bool) error {
	if dryRun {
		return nil
	}
	parent := filepath.Dir(path)
	_ = os.MkdirAll(parent, 0755)
	if err := os.WriteFile(path, []byte(""), 0644); err != nil {
		return fmt.Errorf("create file %s: %w", path, err)
	}
	return nil
}
