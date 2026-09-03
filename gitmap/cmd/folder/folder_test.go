package folder

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func createTestWorkspace(t *testing.T) string {
	tempDir := t.TempDir()

	_ = os.MkdirAll(filepath.Join(tempDir, "src", "core"), 0755)
	_ = os.MkdirAll(filepath.Join(tempDir, "vendor", "lib"), 0755)
	_ = os.MkdirAll(filepath.Join(tempDir, "docs"), 0755)

	_ = os.WriteFile(filepath.Join(tempDir, "src", "01-app.ts"), []byte("line1\nline2\nline3\n"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "src", "core", "util.go"), []byte("package core\nfunc Run() {}\n"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "docs", "01-index.md"), []byte("# Index\nDocs here\n"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "vendor", "lib", "dep.js"), []byte("console.log('dep');"), 0644)
	_ = os.WriteFile(filepath.Join(tempDir, "src", "logo.png"), []byte{0x89, 0x50, 0x4E, 0x47, 0x00}, 0644)

	return tempDir
}

func TestScanDirectoryWithMultiGlobExcept(t *testing.T) {
	ws := createTestWorkspace(t)

	filter := FilterConfig{
		ExceptGlobs: []string{"vendor/**", "*.png"},
	}

	files, err := ScanDirectory(ws, filter)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(files) != 3 {
		t.Fatalf("expected 3 files (excluding vendor and png), got %d", len(files))
	}

	for _, f := range files {
		if strings.HasPrefix(f.Path, "vendor") {
			t.Errorf("vendor file not excluded: %s", f.Path)
		}

		if f.Extension == ".png" {
			t.Errorf("png file not excluded: %s", f.Path)
		}
	}
}

func TestRenderTreeDetailed(t *testing.T) {
	ws := createTestWorkspace(t)
	filter := FilterConfig{
		ExceptGlobs: []string{"vendor/**", "*.png"},
	}

	files, _ := ScanDirectory(ws, filter)
	rootNode := BuildTree("my-project", files)
	rendered := RenderTree(rootNode, true)

	if !strings.Contains(rendered, "01-app.ts (seq: 01, 3 lines") {
		t.Errorf("expected detailed metadata for 01-app.ts, got:\n%s", rendered)
	}

	if !strings.Contains(rendered, "my-project/") {
		t.Errorf("expected root header in tree, got:\n%s", rendered)
	}
}

func TestRenderMarkdownNestedList(t *testing.T) {
	ws := createTestWorkspace(t)
	filter := FilterConfig{
		ExceptGlobs: []string{"vendor/**", "*.png"},
	}

	files, _ := ScanDirectory(ws, filter)
	rootNode := BuildTree("root", files)
	md := RenderMarkdown(rootNode, true)

	if !strings.Contains(md, "- src/") {
		t.Errorf("expected markdown folder bullet, got:\n%s", md)
	}

	if !strings.Contains(md, "01-app.ts (seq: 01, 3 lines") {
		t.Errorf("expected markdown file bullet with details, got:\n%s", md)
	}
}

func TestRenderJson(t *testing.T) {
	ws := createTestWorkspace(t)
	filter := FilterConfig{
		ExceptGlobs: []string{"vendor/**"},
	}

	files, _ := ScanDirectory(ws, filter)
	report := BuildReport(ws, files)
	jsonStr, err := RenderJson(report)
	if err != nil {
		t.Fatalf("RenderJson failed: %v", err)
	}

	if !strings.Contains(jsonStr, `"totalFiles": 4`) {
		t.Errorf("expected totalFiles 4 in JSON, got:\n%s", jsonStr)
	}

	if !strings.Contains(jsonStr, `"isBinary": true`) {
		t.Errorf("expected isBinary true for png in JSON, got:\n%s", jsonStr)
	}
}
