package desktop

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

func TestFormatInstallSuggestions_Windows(t *testing.T) {
	out := FormatInstallSuggestions(constants.PlatformWindows)
	assertContains(t, out, "gitmap install github-desktop")
	assertContains(t, out, "winget install")
	assertContains(t, out, URLGitHubDesktopOfficial)
}

func TestFormatInstallSuggestions_Linux(t *testing.T) {
	out := FormatInstallSuggestions(constants.PlatformLinux)
	assertContains(t, out, "gitmap install github-desktop")
	assertContains(t, out, "snap install")
	assertContains(t, out, URLGitHubDesktopOfficial)
}

func TestFormatInstallSuggestions_Darwin(t *testing.T) {
	out := FormatInstallSuggestions(constants.PlatformDarwin)
	assertContains(t, out, "gitmap install github-desktop")
	assertContains(t, out, "brew install")
	assertContains(t, out, URLGitHubDesktopOfficial)
}

func TestFormatInstallSuggestions_Unknown(t *testing.T) {
	out := FormatInstallSuggestions("unknown-os")
	assertContains(t, out, "gitmap install github-desktop")
	assertContains(t, out, URLGitHubDesktopOfficial)
}

func TestResolveNativeCommand(t *testing.T) {
	winCmd := resolveNativeCommand(constants.PlatformWindows)
	linuxCmd := resolveNativeCommand(constants.PlatformLinux)
	darwinCmd := resolveNativeCommand(constants.PlatformDarwin)
	otherCmd := resolveNativeCommand("plan9")
	assertCommandMatch(t, winCmd, "winget install")
	assertCommandMatch(t, linuxCmd, "snap install")
	assertCommandMatch(t, darwinCmd, "brew install")
	assertCommandMatch(t, otherCmd, "")
}

func TestNewMissingCLIError(t *testing.T) {
	appErr := NewMissingCLIError()
	if appErr == nil {
		t.Fatal("expected non-nil AppError")
	}
	assertErrorProps(t, appErr)
}

func TestPrintInstallSuggestions(t *testing.T) {
	PrintInstallSuggestions()
}

func assertErrorProps(t *testing.T, appErr *apperror.AppError) {
	if appErr.Type != apperror.ErrorTypeValidation {
		t.Errorf("got type %v, want validation", appErr.Type)
	}
	if appErr.Message != constants.MsgDesktopNotFound {
		t.Errorf("got message %q, want %q", appErr.Message, constants.MsgDesktopNotFound)
	}
}

func assertContains(t *testing.T, content, substr string) {
	hasSubstr := strings.Contains(content, substr)
	if !hasSubstr {
		t.Errorf("expected output to contain %q, got: %s", substr, content)
	}
}

func assertCommandMatch(t *testing.T, got, wantSubstr string) {
	isEmptyExpected := wantSubstr == ""
	if isEmptyExpected {
		assertEmptyCommand(t, got)
		return
	}
	hasSubstr := strings.Contains(got, wantSubstr)
	if !hasSubstr {
		t.Errorf("command %q does not contain %q", got, wantSubstr)
	}
}

func assertEmptyCommand(t *testing.T, got string) {
	if got != "" {
		t.Errorf("command = %q, want empty", got)
	}
}
