// Package cmdos — os_machine_alias_prompt_test.go tests interactive prompt logic and argument resolution.
package cmdos

import (
	"os"
	"strings"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/apperror"
)

func TestPromptMachineNameFromReaderValid(t *testing.T) {
	input := "dev-win-01\n"
	r := strings.NewReader(input)

	val, err := promptMachineNameFromReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "dev-win-01" {
		t.Errorf("expected 'dev-win-01', got %q", val)
	}
}

func TestPromptMachineNameFromReaderWhitespace(t *testing.T) {
	input := "  ubuntu-node-02  \r\n"
	r := strings.NewReader(input)

	val, err := promptMachineNameFromReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "ubuntu-node-02" {
		t.Errorf("expected 'ubuntu-node-02', got %q", val)
	}
}

func TestPromptMachineNameFromReaderEmpty(t *testing.T) {
	input := "   \n"
	r := strings.NewReader(input)

	val, err := promptMachineNameFromReader(r)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if val != "" {
		t.Errorf("expected empty string, got %q", val)
	}
}

func TestResolveMachineAliasArgsWithArg(t *testing.T) {
	args := []string{"my-node-01"}

	res, err := resolveMachineAliasArgs(args, false)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(res) != 1 || res[0] != "my-node-01" {
		t.Errorf("expected ['my-node-01'], got %v", res)
	}
}

func TestResolveMachineAliasArgsNonInteractiveJSON(t *testing.T) {
	_, err := resolveMachineAliasArgs([]string{}, true)
	if err == nil {
		t.Fatalf("expected error for non-interactive JSON")
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if !isAppErr || appErr.Code != "E_MISSING_MACHINE_NAME" {
		t.Errorf("expected code E_MISSING_MACHINE_NAME, got %v", err)
	}
}

func TestResolveMachineAliasArgsNonInteractiveCI(t *testing.T) {
	os.Setenv("CI", "true")
	defer os.Unsetenv("CI")

	_, err := resolveMachineAliasArgs([]string{}, false)
	if err == nil {
		t.Fatalf("expected error when CI is set")
	}

	appErr, isAppErr := err.(*apperror.AppError)
	if !isAppErr || appErr.Code != "E_MISSING_MACHINE_NAME" {
		t.Errorf("expected code E_MISSING_MACHINE_NAME, got %v", err)
	}
}

func TestIsTerminalInteractive(t *testing.T) {
	os.Setenv("GITMAP_NON_INTERACTIVE", "1")
	defer os.Unsetenv("GITMAP_NON_INTERACTIVE")

	isInteractive := isTerminalInteractive()
	if isInteractive {
		t.Errorf("expected isTerminalInteractive to be false when GITMAP_NON_INTERACTIVE=1")
	}
}
