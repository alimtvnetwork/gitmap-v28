package cmd

import (
	"os"
	"path/filepath"
	"testing"
)

func TestNormalizeConfigTool(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{"vscode", "vscode"},
		{"code", "vscode"},
		{"vcode", "vscode"},
		{"qtorrent", "qtorrent"},
		{"qbittorrent", "qtorrent"},
		{"qbit", "qtorrent"},
		{"utorrent", "utorrent"},
		{"uttorrent", "utorrent"},
		{"u-torrent", "utorrent"},
		{"all", "all"},
	}
	for _, tc := range cases {
		got := normalizeConfigTool(tc.input)
		if got != tc.want {
			t.Errorf("normalizeConfigTool(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestResolveDefaultConfigFileName(t *testing.T) {
	if got := resolveDefaultConfigFileName("qtorrent", "qtorrent"); got != "qtorrent.json" {
		t.Errorf("expected qtorrent.json, got %s", got)
	}
	if got := resolveDefaultConfigFileName("uttorrent", "utorrent"); got != "uttorrent.json" {
		t.Errorf("expected uttorrent.json, got %s", got)
	}
	if got := resolveDefaultConfigFileName("vscode", "vscode"); got != "vscode.json" {
		t.Errorf("expected vscode.json, got %s", got)
	}
}

func TestExportAndImportRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	exportPath := filepath.Join(tmpDir, "qtorrent.json")
	files := collectDefaultToolFiles("qtorrent")
	bundle := buildConfigBundle("qtorrent", tmpDir, files, nil)
	if err := writeBundleJSON(bundle, exportPath); err != nil {
		t.Fatalf("failed to write bundle: %v", err)
	}
	loaded, loadErr := loadConfigBundle(exportPath)
	if loadErr != nil {
		t.Fatalf("failed to load bundle: %v", loadErr)
	}
	restoreDir := filepath.Join(tmpDir, "restored")
	count, restoreErr := restoreBundleFiles(loaded, restoreDir)
	if restoreErr != nil {
		t.Fatalf("failed to restore files: %v", restoreErr)
	}
	if count != len(files) {
		t.Errorf("expected %d restored files, got %d", len(files), count)
	}
}

func TestExportAllToolConfigs(t *testing.T) {
	tmpDir := t.TempDir()
	if err := exportAllToolConfigs(tmpDir); err != nil {
		t.Fatalf("exportAllToolConfigs failed: %v", err)
	}
	expectedFiles := []string{"vscode.json", "qtorrent.json", "utorrent.json"}
	for _, f := range expectedFiles {
		fullPath := filepath.Join(tmpDir, f)
		if _, err := os.Stat(fullPath); err != nil {
			t.Errorf("expected exported file %s does not exist", fullPath)
		}
	}
}

func TestConfigCLICalls(t *testing.T) {
	if err := runExportConfig([]string{"--help"}); err != nil {
		t.Errorf("export-config help returned error: %v", err)
	}
	if err := runImportConfig([]string{"--help"}); err != nil {
		t.Errorf("import-config help returned error: %v", err)
	}
}
