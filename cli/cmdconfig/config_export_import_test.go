package cmdconfig

import (
	"os"
	"path/filepath"
	"testing"

	cfgconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/constants"
	cfgenum "github.com/alimtvnetwork/gitmap-v28/cli/cmdconfig/enums"
)

func TestNormalizeConfigTool(t *testing.T) {
	cases := []struct {
		input string
		want  string
	}{
		{cfgconst.ToolVSCode, cfgconst.ToolVSCode},
		{cfgconst.ToolCode, cfgconst.ToolVSCode},
		{cfgconst.ToolVCode, cfgconst.ToolVSCode},
		{cfgconst.ToolQTorrent, cfgconst.ToolQTorrent},
		{cfgconst.ToolQBittorrent, cfgconst.ToolQTorrent},
		{cfgconst.ToolQBit, cfgconst.ToolQTorrent},
		{cfgconst.ToolUTorrent, cfgconst.ToolUTorrent},
		{cfgconst.ToolUTTorrent, cfgconst.ToolUTorrent},
		{cfgconst.ToolUTorrentHyphen, cfgconst.ToolUTorrent},
		{cfgconst.ToolAll, cfgconst.ToolAll},
	}

	for _, tc := range cases {
		got := normalizeConfigTool(tc.input)
		if got != tc.want {
			t.Errorf("normalizeConfigTool(%q) = %q; want %q", tc.input, got, tc.want)
		}
	}
}

func TestResolveDefaultConfigFileName(t *testing.T) {
	if got := resolveDefaultConfigFileName(cfgconst.ToolQTorrent, cfgconst.ToolQTorrent); got != cfgconst.ToolQTorrent+cfgconst.ExtJSON {
		t.Errorf("expected qtorrent.json, got %s", got)
	}

	if got := resolveDefaultConfigFileName(cfgconst.ToolUTTorrent, cfgconst.ToolUTorrent); got != cfgconst.ToolUTTorrent+cfgconst.ExtJSON {
		t.Errorf("expected uttorrent.json, got %s", got)
	}

	if got := resolveDefaultConfigFileName(cfgconst.ToolVSCode, cfgconst.ToolVSCode); got != cfgconst.ToolVSCode+cfgconst.ExtJSON {
		t.Errorf("expected vscode.json, got %s", got)
	}
}

func TestExportAndImportRoundtrip(t *testing.T) {
	tmpDir := t.TempDir()
	exportPath := filepath.Join(tmpDir, cfgconst.ToolQTorrent+cfgconst.ExtJSON)
	files := collectDefaultToolFiles(cfgconst.ToolQTorrent)
	bundle := buildConfigBundle(cfgconst.ToolQTorrent, tmpDir, files, nil)
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

	expectedFiles := []string{
		cfgconst.ToolVSCode + cfgconst.ExtJSON,
		cfgconst.ToolQTorrent + cfgconst.ExtJSON,
		cfgconst.ToolUTorrent + cfgconst.ExtJSON,
	}
	for _, f := range expectedFiles {
		fullPath := filepath.Join(tmpDir, f)
		if _, err := os.Stat(fullPath); err != nil {
			t.Errorf("expected exported file %s does not exist", fullPath)
		}
	}
}

func TestConfigCLICalls(t *testing.T) {
	if err := runExportConfig([]string{cfgconst.FlagArgHelp}); err != nil {
		t.Errorf("export-config help returned error: %v", err)
	}

	if err := runImportConfig([]string{cfgconst.FlagArgHelp}); err != nil {
		t.Errorf("import-config help returned error: %v", err)
	}
}

func TestResolveConfigFilePathDirectoryHandling(t *testing.T) {
	tmpDir := t.TempDir()
	defaultFile := cfgconst.ToolQTorrent + cfgconst.ExtJSON
	got := resolveConfigFilePath(tmpDir, defaultFile)
	want := filepath.Join(tmpDir, defaultFile)
	if got != want {
		t.Errorf("resolveConfigFilePath(%q) = %q; want %q", tmpDir, got, want)
	}
}

func TestBase64BinaryDecodeRoundtrip(t *testing.T) {
	payload := ConfigFilePayload{
		Name:     cfgconst.FileSettingsDat,
		Encoding: cfgconst.EncodingBase64,
		Content:  "ZDRpZDFleQ==",
	}

	bytes, err := decodeConfigFileBytes(payload)
	if err != nil || len(bytes) == 0 {
		t.Fatalf("failed to decode base64 payload: %v", err)
	}
}

func TestResolveTargetFileName(t *testing.T) {
	orig := cfgconst.FileQBittorrentConf
	target := resolveTargetFileName(cfgconst.ToolQTorrent, orig)
	if target != cfgconst.FileQBittorrentIni && target != cfgconst.FileQBittorrentConf {
		t.Errorf("unexpected target file name %q", target)
	}
}

func TestImportAllToolConfigs(t *testing.T) {
	tmpDir := t.TempDir()
	if err := exportAllToolConfigs(tmpDir); err != nil {
		t.Fatalf("exportAllToolConfigs failed: %v", err)
	}

	if err := importAllToolConfigs(tmpDir); err != nil {
		t.Fatalf("importAllToolConfigs failed: %v", err)
	}
}

func TestConfigEnumsAndConstantsCoverage(t *testing.T) {
	action := cfgenum.ConfigActionExport
	if !action.IsExport() {
		t.Errorf("expected action to be export")
	}

	format := cfgenum.ConfigFileFormatJSON
	if !format.IsJSON() {
		t.Errorf("expected format to be JSON")
	}
}
