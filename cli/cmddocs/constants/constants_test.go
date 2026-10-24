package constants

import (
	"strings"
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testCommandConstants(t)
	testURLConstants(t)
	testFlagConstants(t)
	testFlagArgConstants(t)
	testOSTokenConstants(t)
	testMessageConstants(t)
	testErrorCodeConstants(t)
	testExitCodeConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdDocs != "docs" || CmdDocsAlias != "d" {
		t.Error("docs command and alias mismatch")
	}
	if HelpDocs == "" || HelpUsage == "" {
		t.Error("help strings must not be empty")
	}
}

func testURLConstants(t *testing.T) {
	if !strings.HasPrefix(DocsURL, "https://") {
		t.Errorf("expected DocsURL to begin with https://, got %s", DocsURL)
	}
	if DefaultServerHost != "127.0.0.1" || DefaultServerPort != 8080 {
		t.Error("default server host/port mismatch")
	}
}

func testFlagConstants(t *testing.T) {
	if FlagBrowser != "browser" || FlagTerminal != "terminal" {
		t.Error("browser or terminal flag mismatch")
	}
	if FlagServer != "server" || FlagCategory != "category" {
		t.Error("server or category flag mismatch")
	}
	if FlagFormat != "format" || FlagPort != "port" || FlagHelp != "help" {
		t.Error("format, port, or help flag mismatch")
	}
	if len(AllDocsFlags()) != 7 {
		t.Errorf("expected 7 docs flags, got %d", len(AllDocsFlags()))
	}
}

func testFlagArgConstants(t *testing.T) {
	if FlagArgBrowser != "--browser" || FlagArgTerminal != "--terminal" {
		t.Error("browser or terminal flag argument mismatch")
	}
	if FlagArgServer != "--server" || FlagArgCategory != "--category" {
		t.Error("server or category flag argument mismatch")
	}
	if FlagArgFormat != "--format" || FlagArgPort != "--port" {
		t.Error("format or port flag argument mismatch")
	}
	if FlagArgHelp != "--help" || FlagArgHelpShort != "-h" {
		t.Error("help flag argument mismatch")
	}
}

func testOSTokenConstants(t *testing.T) {
	if OSWindows != "windows" || OSDarwin != "darwin" || OSLinux != "linux" {
		t.Error("OS identifier constants mismatch")
	}
	if CmdWindowsShell != "cmd" || CmdArgSlashC != "/c" || CmdArgStart != "start" {
		t.Error("Windows execution tokens mismatch")
	}
	if CmdOpen != "open" || CmdXdgOpen != "xdg-open" {
		t.Error("Unix open command tokens mismatch")
	}
}

func testMessageConstants(t *testing.T) {
	if !strings.Contains(MsgDocsOpened, "Opened") {
		t.Error("MsgDocsOpened mismatch")
	}
	if !strings.Contains(MsgDocsServerStarting, "Starting") {
		t.Error("MsgDocsServerStarting mismatch")
	}
	if !strings.Contains(ErrDocsOpen, "Failed to open") {
		t.Error("ErrDocsOpen mismatch")
	}
	if ErrDocsOpenFailed != "failed to open docs browser" {
		t.Error("ErrDocsOpenFailed mismatch")
	}
}

func testErrorCodeConstants(t *testing.T) {
	if ErrCodeDocsLaunchFailed != "EDOCS01" {
		t.Error("ErrCodeDocsLaunchFailed mismatch")
	}
	if ErrCodeDocsInvalidArgs != "EDOCS02" {
		t.Error("ErrCodeDocsInvalidArgs mismatch")
	}
	if ErrCodeDocsServerFailed != "EDOCS03" {
		t.Error("ErrCodeDocsServerFailed mismatch")
	}
}

func testExitCodeConstants(t *testing.T) {
	if ExitSuccess != 0 || ExitGeneralError != 1 {
		t.Error("exit code constants mismatch")
	}
}
