package cmd_test

import (
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdprompt"
)

func TestPromptPlatformContract(t *testing.T) {
	unixCmd := cmdprompt.BuildUnixPromptInstallCmd()
	if !strings.Contains(unixCmd, "curl -sL") || !strings.Contains(unixCmd, "install.sh") {
		t.Fatalf("unexpected unix command: %s", unixCmd)
	}

	winCmd := cmdprompt.BuildWindowsPromptInstallCmd()
	if !strings.Contains(winCmd, "Invoke-Expression") || !strings.Contains(winCmd, "install.ps1") {
		t.Fatalf("unexpected win command: %s", winCmd)
	}
}
