package enums

import (
	"testing"
)

func TestConfigActionType(t *testing.T) {
	all := AllConfigActionTypes()
	if len(all) != 8 {
		t.Fatalf("expected 8 action types, got %d", len(all))
	}

	for _, action := range all {
		if !action.IsValid() {
			t.Errorf("action %s should be valid", action)
		}

		if action.String() == "" {
			t.Errorf("action %s String() should not be empty", action)
		}
	}

	if (ConfigActionType("invalid")).IsValid() {
		t.Errorf("invalid action should not be valid")
	}

	if !ConfigActionExport.IsExport() {
		t.Errorf("expected IsExport true")
	}

	if !ConfigActionImport.IsImport() {
		t.Errorf("expected IsImport true")
	}

	if !ConfigActionWhatConfigs.IsWhatConfigs() {
		t.Errorf("expected IsWhatConfigs true")
	}

	if !ConfigActionWhichFormat.IsWhichFormat() {
		t.Errorf("expected IsWhichFormat true")
	}

	if !ConfigActionGet.IsGet() {
		t.Errorf("expected IsGet true")
	}

	if !ConfigActionSet.IsSet() {
		t.Errorf("expected IsSet true")
	}

	if !ConfigActionList.IsList() {
		t.Errorf("expected IsList true")
	}

	if !ConfigActionReset.IsReset() {
		t.Errorf("expected IsReset true")
	}

	parsed, err := ParseConfigActionType("export")
	if err != nil || parsed != ConfigActionExport {
		t.Errorf("ParseConfigActionType(export) failed: %v", err)
	}

	_, parseErr := ParseConfigActionType("nonexistent")
	if parseErr == nil {
		t.Errorf("ParseConfigActionType(nonexistent) expected error")
	}
}

func TestConfigOutputFormatType(t *testing.T) {
	all := AllConfigOutputFormatTypes()
	if len(all) != 4 {
		t.Fatalf("expected 4 output format types, got %d", len(all))
	}

	for _, format := range all {
		if !format.IsValid() {
			t.Errorf("format %s should be valid", format)
		}
	}

	if (ConfigOutputFormatType("xml")).IsValid() {
		t.Errorf("xml output format should not be valid")
	}

	if !ConfigOutputFormatTable.IsTable() {
		t.Errorf("expected IsTable true")
	}

	if !ConfigOutputFormatJSON.IsJSON() {
		t.Errorf("expected IsJSON true")
	}

	if !ConfigOutputFormatYAML.IsYAML() {
		t.Errorf("expected IsYAML true")
	}

	if !ConfigOutputFormatCompact.IsCompact() {
		t.Errorf("expected IsCompact true")
	}

	parsed, err := ParseConfigOutputFormatType("table")
	if err != nil || parsed != ConfigOutputFormatTable {
		t.Errorf("ParseConfigOutputFormatType(table) failed: %v", err)
	}

	_, parseErr := ParseConfigOutputFormatType("unknown")
	if parseErr == nil {
		t.Errorf("ParseConfigOutputFormatType(unknown) expected error")
	}
}

func TestConfigScopeType(t *testing.T) {
	all := AllConfigScopeTypes()
	if len(all) != 4 {
		t.Fatalf("expected 4 scope types, got %d", len(all))
	}

	for _, scope := range all {
		if !scope.IsValid() {
			t.Errorf("scope %s should be valid", scope)
		}
	}

	if (ConfigScopeType("cluster")).IsValid() {
		t.Errorf("cluster scope should not be valid")
	}

	if !ConfigScopeLocal.IsLocal() {
		t.Errorf("expected IsLocal true")
	}

	if !ConfigScopeGlobal.IsGlobal() {
		t.Errorf("expected IsGlobal true")
	}

	if !ConfigScopeSystem.IsSystem() {
		t.Errorf("expected IsSystem true")
	}

	if !ConfigScopeProject.IsProject() {
		t.Errorf("expected IsProject true")
	}

	parsed, err := ParseConfigScopeType("global")
	if err != nil || parsed != ConfigScopeGlobal {
		t.Errorf("ParseConfigScopeType(global) failed: %v", err)
	}

	_, parseErr := ParseConfigScopeType("invalid")
	if parseErr == nil {
		t.Errorf("ParseConfigScopeType(invalid) expected error")
	}
}

func TestConfigKeyType(t *testing.T) {
	all := AllConfigKeyTypes()
	if len(all) != 6 {
		t.Fatalf("expected 6 config key types, got %d", len(all))
	}

	for _, key := range all {
		if !key.IsValid() {
			t.Errorf("key %s should be valid", key)
		}
	}

	if (ConfigKeyType("invalid_key")).IsValid() {
		t.Errorf("invalid_key should not be valid")
	}

	if !ConfigKeyTargetDir.IsTargetDir() {
		t.Errorf("expected IsTargetDir true")
	}

	if !ConfigKeyDefaultTool.IsDefaultTool() {
		t.Errorf("expected IsDefaultTool true")
	}

	if !ConfigKeyExportDir.IsExportDir() {
		t.Errorf("expected IsExportDir true")
	}

	if !ConfigKeySchemaVersion.IsSchemaVersion() {
		t.Errorf("expected IsSchemaVersion true")
	}

	if !ConfigKeyEncoding.IsEncoding() {
		t.Errorf("expected IsEncoding true")
	}

	if !ConfigKeyAutoImport.IsAutoImport() {
		t.Errorf("expected IsAutoImport true")
	}

	parsed, err := ParseConfigKeyType("target_dir")
	if err != nil || parsed != ConfigKeyTargetDir {
		t.Errorf("ParseConfigKeyType(target_dir) failed: %v", err)
	}

	_, parseErr := ParseConfigKeyType("foo")
	if parseErr == nil {
		t.Errorf("ParseConfigKeyType(foo) expected error")
	}
}

func TestConfigToolType(t *testing.T) {
	all := AllConfigToolTypes()
	if len(all) != 4 {
		t.Fatalf("expected 4 tool types, got %d", len(all))
	}

	for _, tool := range all {
		if !tool.IsValid() {
			t.Errorf("tool %s should be valid", tool)
		}
	}

	if (ConfigToolType("notepad")).IsValid() {
		t.Errorf("notepad should not be valid")
	}

	if !ConfigToolVSCode.IsVSCode() {
		t.Errorf("expected IsVSCode true")
	}

	if !ConfigToolQTorrent.IsQTorrent() {
		t.Errorf("expected IsQTorrent true")
	}

	if !ConfigToolUTorrent.IsUTorrent() {
		t.Errorf("expected IsUTorrent true")
	}

	if !ConfigToolAll.IsAll() {
		t.Errorf("expected IsAll true")
	}

	cases := []struct {
		input string
		want  ConfigToolType
	}{
		{"vscode", ConfigToolVSCode},
		{"code", ConfigToolVSCode},
		{"vcode", ConfigToolVSCode},
		{"qtorrent", ConfigToolQTorrent},
		{"qbittorrent", ConfigToolQTorrent},
		{"qbit", ConfigToolQTorrent},
		{"utorrent", ConfigToolUTorrent},
		{"uttorrent", ConfigToolUTorrent},
		{"u-torrent", ConfigToolUTorrent},
		{"all", ConfigToolAll},
		{"*", ConfigToolAll},
	}

	for _, tc := range cases {
		got, err := ParseConfigToolType(tc.input)
		if err != nil || got != tc.want {
			t.Errorf("ParseConfigToolType(%q) = %q, want %q (err: %v)", tc.input, got, tc.want, err)
		}
	}

	_, parseErr := ParseConfigToolType("invalid_tool")
	if parseErr == nil {
		t.Errorf("ParseConfigToolType(invalid_tool) expected error")
	}
}

func TestConfigFileFormatType(t *testing.T) {
	all := AllConfigFileFormatTypes()
	if len(all) != 9 {
		t.Fatalf("expected 9 file format types, got %d", len(all))
	}

	for _, format := range all {
		if !format.IsValid() {
			t.Errorf("format %s should be valid", format)
		}
	}

	if (ConfigFileFormatType("exe")).IsValid() {
		t.Errorf("exe format should not be valid")
	}

	if !ConfigFileFormatJSON.IsJSON() {
		t.Errorf("expected IsJSON true")
	}

	if !ConfigFileFormatYAML.IsYAML() {
		t.Errorf("expected IsYAML true")
	}

	if !ConfigFileFormatDB.IsDB() {
		t.Errorf("expected IsDB true")
	}

	if !ConfigFileFormatSQLite.IsSQLite() {
		t.Errorf("expected IsSQLite true")
	}

	if !ConfigFileFormatTOML.IsTOML() {
		t.Errorf("expected IsTOML true")
	}

	if !ConfigFileFormatENV.IsENV() {
		t.Errorf("expected IsENV true")
	}

	if !ConfigFileFormatINI.IsINI() {
		t.Errorf("expected IsINI true")
	}

	if !ConfigFileFormatCFG.IsCFG() {
		t.Errorf("expected IsCFG true")
	}

	if !ConfigFileFormatCONF.IsCONF() {
		t.Errorf("expected IsCONF true")
	}

	cases := []struct {
		input string
		want  ConfigFileFormatType
	}{
		{"json", ConfigFileFormatJSON},
		{".json", ConfigFileFormatJSON},
		{"yaml", ConfigFileFormatYAML},
		{".yml", ConfigFileFormatYAML},
		{"db", ConfigFileFormatDB},
		{".sqlite", ConfigFileFormatSQLite},
		{"sqlite3", ConfigFileFormatSQLite},
		{".toml", ConfigFileFormatTOML},
		{".env", ConfigFileFormatENV},
		{"ini", ConfigFileFormatINI},
		{"cfg", ConfigFileFormatCFG},
		{"conf", ConfigFileFormatCONF},
	}

	for _, tc := range cases {
		got, err := ParseConfigFileFormatType(tc.input)
		if err != nil || got != tc.want {
			t.Errorf("ParseConfigFileFormatType(%q) = %q, want %q (err: %v)", tc.input, got, tc.want, err)
		}
	}

	_, parseErr := ParseConfigFileFormatType("unsupported")
	if parseErr == nil {
		t.Errorf("ParseConfigFileFormatType(unsupported) expected error")
	}
}

func TestConfigFileEncodingType(t *testing.T) {
	all := AllConfigFileEncodingTypes()
	if len(all) != 2 {
		t.Fatalf("expected 2 encoding types, got %d", len(all))
	}

	for _, enc := range all {
		if !enc.IsValid() {
			t.Errorf("encoding %s should be valid", enc)
		}
	}

	if (ConfigFileEncodingType("latin1")).IsValid() {
		t.Errorf("latin1 should not be valid")
	}

	if !ConfigFileEncodingUTF8.IsUTF8() {
		t.Errorf("expected IsUTF8 true")
	}

	if !ConfigFileEncodingBase64.IsBase64() {
		t.Errorf("expected IsBase64 true")
	}

	cases := []struct {
		input string
		want  ConfigFileEncodingType
	}{
		{"utf-8", ConfigFileEncodingUTF8},
		{"utf8", ConfigFileEncodingUTF8},
		{"base64", ConfigFileEncodingBase64},
		{"b64", ConfigFileEncodingBase64},
	}

	for _, tc := range cases {
		got, err := ParseConfigFileEncodingType(tc.input)
		if err != nil || got != tc.want {
			t.Errorf("ParseConfigFileEncodingType(%q) = %q, want %q (err: %v)", tc.input, got, tc.want, err)
		}
	}

	_, parseErr := ParseConfigFileEncodingType("binary")
	if parseErr == nil {
		t.Errorf("ParseConfigFileEncodingType(binary) expected error")
	}
}
