package cmd

import (
	"os"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/constants"
)

type mockFileInfo struct {
	os.FileInfo
	mode os.FileMode
}

func (m mockFileInfo) Mode() os.FileMode {

	return m.mode
}

func TestProbeToolZsh_Registration(t *testing.T) {
	bins, args := resolveToolCandidates(constants.ToolZsh)
	isBinMatch := len(bins) == 1 && bins[0] == "zsh"
	isArgMatch := len(args) == 1 && args[0] == "--version"
	isValidConfig := isBinMatch && isArgMatch
	if isValidConfig {

		return
	}
	t.Fatalf("unexpected probe config for ToolZsh: bins=%v, args=%v", bins, args)
}

func TestProbeZshBinary_LookPathFound(t *testing.T) {
	oldLookPath := lookPathFunc
	defer func() { lookPathFunc = oldLookPath }()
	lookPathFunc = func(file string) (string, error) {

		return "/custom/bin/zsh", nil
	}
	path, isFound := findZshBinary()
	isMatch := isFound && path == "/custom/bin/zsh"
	if isMatch {

		return
	}
	t.Fatalf("expected /custom/bin/zsh, got path=%q, isFound=%v", path, isFound)
}

func TestProbeSkipZshEnv(t *testing.T) {
	oldGetenv := getenvFunc
	defer func() { getenvFunc = oldGetenv }()
	getenvFunc = func(key string) string {

		return "1"
	}
	isSkip := isSkipZshEnv()
	if isSkip {

		return
	}
	t.Fatalf("expected isSkipZshEnv() to be true when GITMAP_SKIP_ZSH=1")
}

func TestProbeStdinTerminal(t *testing.T) {
	oldStdinStat := stdinStatFunc
	defer func() { stdinStatFunc = oldStdinStat }()
	stdinStatFunc = func() (os.FileInfo, error) {

		return mockFileInfo{mode: os.ModeCharDevice}, nil
	}
	isTerminal := isStdinTerminal()
	if isTerminal {

		return
	}
	t.Fatalf("expected isStdinTerminal() to be true for char device")
}
