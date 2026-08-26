package cmd_test

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/cmd"
)

func TestPromptPlatformContract(t *testing.T) {
	unixCmd := cmd.BuildUnixPromptInstallCmd()
	if !strings.Contains(unixCmd, "curl -sL") || !strings.Contains(unixCmd, "install.sh") {
		t.Fatalf("unexpected unix command: %s", unixCmd)
	}

	winCmd := cmd.BuildWindowsPromptInstallCmd()
	if !strings.Contains(winCmd, "Invoke-Expression") || !strings.Contains(winCmd, "install.ps1") {
		t.Fatalf("unexpected win command: %s", winCmd)
	}
}
