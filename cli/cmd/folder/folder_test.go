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

func TestParseArgs_WildcardPositional(t *testing.T) {
	opts, err := ParseArgs([]string{"*.go"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.TargetDir != "." {
		t.Errorf("expected TargetDir '.', got '%s'", opts.TargetDir)
	}

	if len(opts.Filter.IncludeGlobs) != 1 || opts.Filter.IncludeGlobs[0] != "*.go" {
		t.Errorf("expected IncludeGlobs ['*.go'], got %v", opts.Filter.IncludeGlobs)
	}

	if len(opts.Filter.Extensions) != 1 || opts.Filter.Extensions[0] != "go" {
		t.Errorf("expected Extensions ['go'], got %v", opts.Filter.Extensions)
	}

	if opts.OutFile != "" {
		t.Errorf("expected empty OutFile, got '%s'", opts.OutFile)
	}
}

func TestParseArgs_FileExportFlags(t *testing.T) {
	cases := []struct {
		name           string
		args           []string
		expectedOut    string
		expectedFormat OutputFormat
	}{
		{
			name:           "flag --file txt",
			args:           []string{"--file", "a.txt"},
			expectedOut:    "a.txt",
			expectedFormat: FormatTree,
		},
		{
			name:           "flag -f json",
			args:           []string{"-f", "a.json"},
			expectedOut:    "a.json",
			expectedFormat: FormatJson,
		},
		{
			name:           "flag -f yaml",
			args:           []string{"-f", "a.yaml"},
			expectedOut:    "a.yaml",
			expectedFormat: FormatYaml,
		},
		{
			name:           "flag -f yml",
			args:           []string{"-f", "a.yml"},
			expectedOut:    "a.yml",
			expectedFormat: FormatYaml,
		},
		{
			name:           "flag --file md",
			args:           []string{"--file", "a.md"},
			expectedOut:    "a.md",
			expectedFormat: FormatMd,
		},
		{
			name:           "flag -f tree",
			args:           []string{"-f", "a.tree"},
			expectedOut:    "a.tree",
			expectedFormat: FormatTree,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts, err := ParseArgs(tc.args)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if opts.OutFile != tc.expectedOut {
				t.Errorf("expected OutFile '%s', got '%s'", tc.expectedOut, opts.OutFile)
			}

			if opts.Format != tc.expectedFormat {
				t.Errorf("expected Format '%s', got '%s'", tc.expectedFormat, opts.Format)
			}
		})
	}
}

func TestParseArgs_CombinedTargetAndGlob(t *testing.T) {
	opts, err := ParseArgs([]string{"src", "*.ts"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if opts.TargetDir != "src" {
		t.Errorf("expected TargetDir 'src', got '%s'", opts.TargetDir)
	}

	if len(opts.Filter.IncludeGlobs) != 1 || opts.Filter.IncludeGlobs[0] != "*.ts" {
		t.Errorf("expected IncludeGlobs ['*.ts'], got %v", opts.Filter.IncludeGlobs)
	}

	if len(opts.Filter.Extensions) != 1 || opts.Filter.Extensions[0] != "ts" {
		t.Errorf("expected Extensions ['ts'], got %v", opts.Filter.Extensions)
	}

	if opts.OutFile != "" {
		t.Errorf("expected empty OutFile, got '%s'", opts.OutFile)
	}
}

func TestParseArgs_ExplicitFormatPrecedence(t *testing.T) {
	optsFlat, err := ParseArgs([]string{"--flat", "-f", "out.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if optsFlat.Format != FormatFlat {
		t.Errorf("expected FormatFlat, got '%s'", optsFlat.Format)
	}

	optsJson, err := ParseArgs([]string{"--json", "-f", "out.txt"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if optsJson.Format != FormatJson {
		t.Errorf("expected FormatJson, got '%s'", optsJson.Format)
	}
}

func TestDeduceFormatFromExt(t *testing.T) {
	if f := deduceFormatFromExt("out.json", FormatTree); f != FormatJson {
		t.Errorf("expected FormatJson, got %s", f)
	}

	if f := deduceFormatFromExt("out.yaml", FormatTree); f != FormatYaml {
		t.Errorf("expected FormatYaml, got %s", f)
	}

	if f := deduceFormatFromExt("out.yml", FormatTree); f != FormatYaml {
		t.Errorf("expected FormatYaml, got %s", f)
	}

	if f := deduceFormatFromExt("out.txt", FormatTree); f != FormatTree {
		t.Errorf("expected FormatTree, got %s", f)
	}

	if f := deduceFormatFromExt("out.tree", FormatTree); f != FormatTree {
		t.Errorf("expected FormatTree, got %s", f)
	}

	if f := deduceFormatFromExt("out.md", FormatTree); f != FormatMd {
		t.Errorf("expected FormatMd, got %s", f)
	}

	if f := deduceFormatFromExt("out.unknown", FormatTree); f != FormatTree {
		t.Errorf("expected FormatTree fallback, got %s", f)
	}

	if f := deduceFormatFromExt("out.txt", FormatFlat); f != FormatFlat {
		t.Errorf("expected FormatFlat when fallback is flat, got %s", f)
	}
}

func TestScanDirectory_IncludeGlobs(t *testing.T) {
	ws := createTestWorkspace(t)
	filter := FilterConfig{
		IncludeGlobs: []string{"*.go"},
	}

	files, err := ScanDirectory(ws, filter)
	if err != nil {
		t.Fatalf("ScanDirectory failed: %v", err)
	}

	if len(files) != 1 {
		t.Fatalf("expected 1 file matching *.go, got %d", len(files))
	}

	if files[0].Filename != "util.go" {
		t.Errorf("expected util.go, got %s", files[0].Filename)
	}
}

func TestWriteOrPrintOutput_NestedDirectories(t *testing.T) {
	tempDir := t.TempDir()
	targetPath := filepath.Join(tempDir, "sub", "dir", "tree.txt")
	testContent := "test tree content\nline 2\n"

	err := writeOrPrintOutput(testContent, targetPath)
	if err != nil {
		t.Fatalf("writeOrPrintOutput failed: %v", err)
	}

	readBytes, errRead := os.ReadFile(targetPath)
	if errRead != nil {
		t.Fatalf("failed to read created file: %v", errRead)
	}

	if string(readBytes) != testContent {
		t.Errorf("expected content '%s', got '%s'", testContent, string(readBytes))
	}
}
