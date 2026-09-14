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
	assertContains(t, out, "choco install")
	assertContains(t, out, URLGitHubDesktopOfficial)
}

func TestFormatInstallSuggestions_Linux(t *testing.T) {
	out := FormatInstallSuggestions(constants.PlatformLinux)
	assertContains(t, out, "gitmap install github-desktop")
	assertContains(t, out, "shiftkey")
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

func TestResolveNativeCommands(t *testing.T) {
	winCmds := resolveNativeCommands(constants.PlatformWindows)
	linuxCmds := resolveNativeCommands(constants.PlatformLinux)
	darwinCmds := resolveNativeCommands(constants.PlatformDarwin)
	otherCmds := resolveNativeCommands("plan9")
	assertCommandCount(t, len(winCmds), 2)
	assertCommandCount(t, len(linuxCmds), 2)
	assertCommandCount(t, len(darwinCmds), 1)
	assertCommandCount(t, len(otherCmds), 0)
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

func assertCommandCount(t *testing.T, got, want int) {
	if got != want {
		t.Errorf("command count = %d, want %d", got, want)
	}
}
