package constants

import (
	"strings"
	"testing"
)

func TestSubcommandConstants(t *testing.T) {
	if CmdExportConfig != "export-config" {
		t.Errorf("CmdExportConfig invariant mismatch: %s", CmdExportConfig)
	}

	if CmdImportConfig != "import-config" {
		t.Errorf("CmdImportConfig invariant mismatch: %s", CmdImportConfig)
	}

	if CmdWhatConfigs != "what-configs" {
		t.Errorf("CmdWhatConfigs invariant mismatch: %s", CmdWhatConfigs)
	}

	if CmdWhichFormat != "which-format" {
		t.Errorf("CmdWhichFormat invariant mismatch: %s", CmdWhichFormat)
	}
}

func TestFlagConstants(t *testing.T) {
	if !strings.HasPrefix(FlagArgDir, "--") {
		t.Errorf("FlagArgDir should have double dash prefix")
	}

	if !strings.HasPrefix(FlagArgYes, "--") {
		t.Errorf("FlagArgYes should have double dash prefix")
	}

	if !strings.HasPrefix(FlagArgHelp, "--") {
		t.Errorf("FlagArgHelp should have double dash prefix")
	}

	if FlagArgDirPrefix != "--dir=" {
		t.Errorf("FlagArgDirPrefix mismatch: %s", FlagArgDirPrefix)
	}
}

func TestToolAndPathConstants(t *testing.T) {
	if ToolVSCode != "vscode" {
		t.Errorf("ToolVSCode mismatch: %s", ToolVSCode)
	}

	if ToolQTorrent != "qtorrent" {
		t.Errorf("ToolQTorrent mismatch: %s", ToolQTorrent)
	}

	if ToolUTorrent != "utorrent" {
		t.Errorf("ToolUTorrent mismatch: %s", ToolUTorrent)
	}

	if ToolAll != "all" {
		t.Errorf("ToolAll mismatch: %s", ToolAll)
	}

	if PathVSCodeUser != "Code/User" {
		t.Errorf("PathVSCodeUser mismatch: %s", PathVSCodeUser)
	}
}

func TestCanonicalFileConstants(t *testing.T) {
	files := []string{
		FileSettingsJSON,
		FileKeybindingsJSON,
		FileQBittorrentIni,
		FileQBittorrentConf,
		FileSettingsDat,
		FileExtensionsTxt,
		FilePackageJSON,
		FileVersionJSON,
		FileTSConfigJSON,
		FileCargoTOML,
		FilePyprojectTOML,
		FileInstallationDB,
		FileGitmapDB,
		FileDotEnv,
	}

	for _, f := range files {
		if strings.TrimSpace(f) == "" {
			t.Errorf("empty canonical file constant found")
		}
	}
}

func TestDefaultTemplates(t *testing.T) {
	if !strings.Contains(DefaultVSCodeSettings, "editor.tabSize") {
		t.Errorf("DefaultVSCodeSettings missing tabSize")
	}

	if !strings.Contains(DefaultQBittorrentConfig, "LegalNotice") {
		t.Errorf("DefaultQBittorrentConfig missing LegalNotice")
	}

	if !strings.Contains(DefaultUTorrentSettings, "uTorrent") {
		t.Errorf("DefaultUTorrentSettings missing uTorrent")
	}
}

func TestErrorCodes(t *testing.T) {
	codes := []string{
		ErrCodeConfigExport,
		ErrCodeConfigImport,
		ErrCodeInvalidTool,
		ErrCodeFileNotFound,
	}

	for _, c := range codes {
		if !strings.HasPrefix(c, "E72") {
			t.Errorf("expected error code with E72 prefix, got %s", c)
		}
	}
}
