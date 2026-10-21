// Package cmdos — os_machine_alias_test.go tests interactive prompt logic and machine alias resolution.
package cmdos

import (
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func TestResolveDefaultMachineName(t *testing.T) {
	name := resolveDefaultMachineName("machine")
	if len(name) == 0 {
		t.Fatalf("expected non-empty machine name")
	}

	aliasName := resolveDefaultMachineName("alias")
	if len(aliasName) == 0 {
		t.Fatalf("expected non-empty alias name")
	}
}

func TestExecuteMachineAliasSetNonInteractiveJSON(t *testing.T) {
	err := executeMachineAliasSet("machine", []string{}, false, true)
	if err == nil {
		t.Fatalf("expected error for non-interactive JSON")
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if !isAppErr || appErr.Code != "E_MISSING_MACHINE_NAME" {
		t.Errorf("expected E_MISSING_MACHINE_NAME, got %v", err)
	}
}

func TestExecuteMachineAliasSetNonInteractiveTerminal(t *testing.T) {
	orig := isTerminalFn
	isTerminalFn = func() bool { return false }
	defer func() { isTerminalFn = orig }()

	err := executeMachineAliasSet("machine", []string{}, false, false)
	if err == nil {
		t.Fatalf("expected error when terminal is non-interactive")
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if !isAppErr || appErr.Code != "E_MISSING_MACHINE_NAME" {
		t.Errorf("expected E_MISSING_MACHINE_NAME, got %v", err)
	}
}

func TestPromptForMachineNameInteractive(t *testing.T) {
	origStdin := osStdinReader
	osStdinReader = strings.NewReader("custom-node-01\n")
	defer func() { osStdinReader = origStdin }()

	val, err := promptForMachineName("machine")
	if err != nil {
		t.Fatalf("unexpected prompt error: %v", err)
	}

	if val != "custom-node-01" {
		t.Errorf("expected 'custom-node-01', got %q", val)
	}
}

func TestPromptForMachineNameDefaultOnEmpty(t *testing.T) {
	origStdin := osStdinReader
	osStdinReader = strings.NewReader("\n")
	defer func() { osStdinReader = origStdin }()

	defaultVal := resolveDefaultMachineName("machine")
	val, err := promptForMachineName("machine")
	if err != nil {
		t.Fatalf("unexpected prompt error: %v", err)
	}

	if val != defaultVal {
		t.Errorf("expected default %q, got %q", defaultVal, val)
	}
}

func TestPromptForMachineNameDefaultOnWhitespace(t *testing.T) {
	origStdin := osStdinReader
	osStdinReader = strings.NewReader("   \r\n")
	defer func() { osStdinReader = origStdin }()

	defaultVal := resolveDefaultMachineName("machine")
	val, err := promptForMachineName("machine")
	if err != nil {
		t.Fatalf("unexpected prompt error: %v", err)
	}

	if val != defaultVal {
		t.Errorf("expected default %q, got %q", defaultVal, val)
	}
}

func TestSanitizeMachineToken(t *testing.T) {
	if sanitizeMachineToken("  dev win 01  ") != "dev-win-01" {
		t.Errorf("expected dev-win-01")
	}

	if sanitizeMachineToken("ubuntu-node-02") != "ubuntu-node-02" {
		t.Errorf("expected ubuntu-node-02")
	}

	if sanitizeMachineToken("alpha   beta") != "alpha---beta" {
		t.Errorf("expected alpha---beta")
	}
}

func TestResolveMachineSetArgsExplicit(t *testing.T) {
	args := []string{"dev-win-01"}
	res, err := resolveMachineSetArgs("machine", args, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 || res[0] != "dev-win-01" {
		t.Errorf("expected ['dev-win-01'], got %v", res)
	}
}

func TestIsInteractiveTerminalDefault(t *testing.T) {
	os.Setenv("GITMAP_NON_INTERACTIVE", "1")
	defer os.Unsetenv("GITMAP_NON_INTERACTIVE")

	if isInteractiveTerminal() {
		t.Errorf("expected isInteractiveTerminal to be false when GITMAP_NON_INTERACTIVE=1")
	}
}
