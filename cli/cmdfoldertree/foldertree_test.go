package cmdfoldertree

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createMockTree(t *testing.T) string {
	t.Helper()
	tmp := t.TempDir()

	subA := filepath.Join(tmp, "subA")
	subB := filepath.Join(tmp, "subB")
	repoGit := filepath.Join(tmp, "repoGit")
	if err := os.MkdirAll(subA, 0755); err != nil {
		t.Fatalf("mkdir subA: %v", err)
	}
	if err := os.MkdirAll(subB, 0755); err != nil {
		t.Fatalf("mkdir subB: %v", err)
	}
	if err := os.MkdirAll(repoGit, 0755); err != nil {
		t.Fatalf("mkdir repoGit: %v", err)
	}

	if err := os.WriteFile(filepath.Join(subA, "fileA.txt"), []byte("hello"), 0644); err != nil {
		t.Fatalf("write fileA: %v", err)
	}
	if err := os.WriteFile(filepath.Join(tmp, "root.md"), []byte("readme"), 0644); err != nil {
		t.Fatalf("write root.md: %v", err)
	}

	dotGit := filepath.Join(repoGit, ".git")
	if err := os.MkdirAll(dotGit, 0755); err != nil {
		t.Fatalf("mkdir dotGit: %v", err)
	}
	if err := os.WriteFile(filepath.Join(dotGit, "HEAD"), []byte("ref: refs/heads/feature-x\n"), 0644); err != nil {
		t.Fatalf("write HEAD: %v", err)
	}

	return tmp
}

func TestScanFolderTree(t *testing.T) {
	rootPath := createMockTree(t)
	opts := FolderTreeOptions{
		TargetDir: rootPath,
		MaxDepth:  0,
		DirsOnly:  false,
	}

	rootNode, err := ScanFolderTree(rootPath, opts)
	if err != nil {
		t.Fatalf("ScanFolderTree error: %v", err)
	}
	if rootNode == nil || len(rootNode.Children) == 0 {
		t.Fatalf("expected children in root node")
	}

	var gitNode *FolderTreeNode
	for _, child := range rootNode.Children {
		if child.Name == "repoGit" {
			gitNode = child
			break
		}
	}

	if gitNode == nil {
		t.Fatalf("repoGit node not found in children")
	}
	if !gitNode.IsGit {
		t.Errorf("expected repoGit to have IsGit = true")
	}
	if gitNode.GitBranch != "feature-x" {
		t.Errorf("expected git branch feature-x, got %q", gitNode.GitBranch)
	}
}

func TestScanFolderTreeDirsOnly(t *testing.T) {
	rootPath := createMockTree(t)
	opts := FolderTreeOptions{
		TargetDir: rootPath,
		DirsOnly:  true,
	}

	rootNode, err := ScanFolderTree(rootPath, opts)
	if err != nil {
		t.Fatalf("ScanFolderTree error: %v", err)
	}

	for _, child := range rootNode.Children {
		if !child.IsDir {
			t.Errorf("found file node %s when DirsOnly is true", child.Name)
		}
	}
}

func TestRenderTree(t *testing.T) {
	rootPath := createMockTree(t)
	opts := FolderTreeOptions{TargetDir: rootPath}

	rootNode, err := ScanFolderTree(rootPath, opts)
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	output := RenderTree(rootNode, opts)
	if !strings.Contains(output, "📁") && !strings.Contains(output, "📦") {
		t.Errorf("expected emoji in RenderTree output, got:\n%s", output)
	}
	if !strings.Contains(output, "[git: feature-x]") {
		t.Errorf("expected git branch tag in output, got:\n%s", output)
	}
}

func TestRenderPreview(t *testing.T) {
	rootPath := createMockTree(t)
	opts := FolderTreeOptions{TargetDir: rootPath}

	rootNode, err := ScanFolderTree(rootPath, opts)
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	output := RenderPreview(rootNode, opts)
	if !strings.Contains(output, "\n\n") {
		t.Errorf("expected blank line gap between entries in RenderPreview")
	}
	if !strings.Contains(output, "repoGit") {
		t.Errorf("expected repoGit in RenderPreview output")
	}
}

func TestExportAndImportRoundTrip(t *testing.T) {
	rootPath := createMockTree(t)
	opts := FolderTreeOptions{
		TargetDir: rootPath,
		Format:    "json",
	}

	rootNode, err := ScanFolderTree(rootPath, opts)
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	exportedJSON, err := ExportFolderTree(rootNode, opts)
	if err != nil {
		t.Fatalf("ExportFolderTree error: %v", err)
	}

	jsonFile := filepath.Join(t.TempDir(), "export.json")
	if err := os.WriteFile(jsonFile, []byte(exportedJSON), 0644); err != nil {
		t.Fatalf("write jsonFile: %v", err)
	}

	// Dry Run
	destDir := t.TempDir()
	optsDry := FolderTreeOptions{
		DryRun: true,
	}
	summaryDry, err := ImportFolderTree(jsonFile, destDir, optsDry)
	if err != nil {
		t.Fatalf("ImportFolderTree dry run error: %v", err)
	}
	if summaryDry.DirsCreated == 0 {
		t.Errorf("expected dirs count > 0 in dry run summary")
	}

	// Check dry run did not touch disk
	entries, _ := os.ReadDir(destDir)
	if len(entries) > 0 {
		t.Errorf("dry run created files on disk")
	}

	// Live Run
	optsLive := FolderTreeOptions{
		DryRun: false,
	}
	summaryLive, err := ImportFolderTree(jsonFile, destDir, optsLive)
	if err != nil {
		t.Fatalf("ImportFolderTree live run error: %v", err)
	}
	if summaryLive.DirsCreated == 0 {
		t.Errorf("expected live run to create directories")
	}

	liveEntries, _ := os.ReadDir(destDir)
	if len(liveEntries) == 0 {
		t.Errorf("expected files/directories in destDir after live import")
	}
}

func TestExportYAML(t *testing.T) {
	rootPath := createMockTree(t)
	opts := FolderTreeOptions{
		TargetDir: rootPath,
		Format:    "yaml",
	}

	rootNode, err := ScanFolderTree(rootPath, opts)
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	out, err := ExportFolderTree(rootNode, opts)
	if err != nil {
		t.Fatalf("ExportFolderTree yaml error: %v", err)
	}
	if !strings.Contains(out, "rootPath:") && !strings.Contains(out, "totalNodes:") {
		t.Errorf("expected yaml structure in output: %s", out)
	}
}

func TestExportFolderPaths(t *testing.T) {
	rootPath := createMockTree(t)
	opts := FolderTreeOptions{
		TargetDir: rootPath,
		Format:    "folder",
	}

	rootNode, err := ScanFolderTree(rootPath, opts)
	if err != nil {
		t.Fatalf("scan error: %v", err)
	}

	out, err := ExportFolderTree(rootNode, opts)
	if err != nil {
		t.Fatalf("ExportFolderTree folder error: %v", err)
	}
	if !strings.Contains(out, "subA/") && !strings.Contains(out, "subB/") {
		t.Errorf("expected folder paths in output: %s", out)
	}
}

func TestImportFolderOnly(t *testing.T) {
	jsonContent := `{
		"rootName": "sample",
		"tree": {
			"name": "sample",
			"isDir": true,
			"children": [
				{
					"name": "config",
					"isDir": true,
					"children": [
						{"name": "app.json", "isDir": false}
					]
				}
			]
		}
	}`
	jsonFile := filepath.Join(t.TempDir(), "structure.json")
	if err := os.WriteFile(jsonFile, []byte(jsonContent), 0644); err != nil {
		t.Fatalf("write jsonFile: %v", err)
	}

	destDir := t.TempDir()
	opts := FolderTreeOptions{DirsOnly: true}
	summary, err := ImportFolderTree(jsonFile, destDir, opts)
	if err != nil {
		t.Fatalf("ImportFolderTree error: %v", err)
	}
	if summary.DirsCreated == 0 {
		t.Errorf("expected dirs created > 0")
	}
	if summary.FilesCreated != 0 {
		t.Errorf("expected 0 files created when DirsOnly is true, got %d", summary.FilesCreated)
	}
	if _, err := os.Stat(filepath.Join(destDir, "config")); os.IsNotExist(err) {
		t.Errorf("expected config directory to exist")
	}
	if _, err := os.Stat(filepath.Join(destDir, "config", "app.json")); !os.IsNotExist(err) {
		t.Errorf("expected app.json to NOT exist when DirsOnly is true")
	}
}

func TestImportTextPaths(t *testing.T) {
	textContent := "models/\nmodels/user.go\nservices/auth.go\n"
	txtFile := filepath.Join(t.TempDir(), "paths.txt")
	if err := os.WriteFile(txtFile, []byte(textContent), 0644); err != nil {
		t.Fatalf("write txtFile: %v", err)
	}

	destDir := t.TempDir()
	opts := FolderTreeOptions{DirsOnly: false}
	summary, err := ImportFolderTree(txtFile, destDir, opts)
	if err != nil {
		t.Fatalf("ImportFolderTree text error: %v", err)
	}
	if summary.FilesCreated != 2 {
		t.Errorf("expected 2 files created, got %d", summary.FilesCreated)
	}
	userFile := filepath.Join(destDir, "models", "user.go")
	if _, err := os.Stat(userFile); os.IsNotExist(err) {
		t.Errorf("expected placeholder user.go to exist")
	}
}
