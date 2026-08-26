package installer

import (
	"testing"
)

func TestPromptRunnersSuite(t *testing.T) {
	tempDir := t.TempDir()

	hasPS := HasPowerShell()
	if !hasPS {
		t.Log("PowerShell not found in test environment")
	}

	hasBash := HasBashAndCurl()
	if !hasBash {
		t.Log("Bash/curl not found in test environment")
	}

	if HasExistingPrompts(tempDir) {
		t.Fatal("expected tempDir to have no existing prompts")
	}
}
