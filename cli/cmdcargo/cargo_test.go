package cmdcargo

import (
	"errors"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	cargoconst "github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/constants"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdcargo/enums"
	"github.com/alimtvnetwork/gitmap-v28/cli/constants"
)

type flagTestCase struct {
	args        []string
	wantHas     bool
	wantCleaned []string
}

func installFlagCases() []flagTestCase {
	return []flagTestCase{
		{
			args:        []string{"build", cargoconst.FlagLongRelease},
			wantHas:     false,
			wantCleaned: []string{"build", cargoconst.FlagLongRelease},
		},
		{
			args:        []string{cargoconst.FlagLongInstall, "fmt", cargoconst.FlagDoubleDash, cargoconst.FlagLongCheck},
			wantHas:     true,
			wantCleaned: []string{"fmt", cargoconst.FlagDoubleDash, cargoconst.FlagLongCheck},
		},
		{
			args:        []string{cargoconst.FlagShortInstall, "check"},
			wantHas:     true,
			wantCleaned: []string{"check"},
		},
	}
}

func assertInstallFlagResult(t *testing.T, args []string, has, wantHas bool, gotLen, wantLen int) {
	if has != wantHas {
		t.Errorf("extractInstallFlag(%v) has = %v, want %v", args, has, wantHas)
	}
	if gotLen != wantLen {
		t.Errorf("extractInstallFlag(%v) cleaned len = %d, want %d", args, gotLen, wantLen)
	}
}

func TestExtractInstallFlag(t *testing.T) {
	for _, tc := range installFlagCases() {
		cleaned, has := extractInstallFlag(tc.args)
		assertInstallFlagResult(t, tc.args, has, tc.wantHas, len(cleaned), len(tc.wantCleaned))
	}
}

func TestFormatCargoInstallSuggestions(t *testing.T) {
	for _, targetOS := range []string{constants.PlatformWindows, constants.PlatformDarwin, constants.PlatformLinux} {
		out := FormatCargoInstallSuggestions(targetOS)
		if !strings.Contains(out, cargoconst.SuggestGitMapInstall) {
			t.Errorf("expected %q to contain %q for OS %s", out, cargoconst.SuggestGitMapInstall, targetOS)
		}
		if !strings.Contains(out, "rustup.rs") {
			t.Errorf("expected %q to contain rustup link for OS %s", out, targetOS)
		}
	}
}

func TestNewMissingCargoError(t *testing.T) {
	err := NewMissingCargoError()
	if err.Code != cargoconst.ErrCodeMissingCargo {
		t.Errorf("expected error code %s, got %s", cargoconst.ErrCodeMissingCargo, err.Code)
	}
}

func TestSubcommandMatchers(t *testing.T) {
	if !isStatusSubcommand(string(enums.CargoSubcommandStatus)) ||
		!isStatusSubcommand(string(enums.CargoSubcommandSt)) ||
		!isStatusSubcommand(string(enums.CargoSubcommandInfo)) {
		t.Errorf("expected status subcommands to match")
	}
	if !isInstallSubcommand(string(enums.CargoSubcommandInstall)) ||
		!isInstallSubcommand(string(enums.CargoSubcommandIn)) {
		t.Errorf("expected install subcommands to match")
	}
	if isStatusSubcommand("invalid") || isInstallSubcommand("invalid") {
		t.Errorf("expected invalid subcommands to return false")
	}
}

func TestHandleCargoExitError(t *testing.T) {
	genericErr := errors.New("command not found")
	returned := handleCargoExitError(genericErr)
	if returned == nil || returned.Cause != genericErr {
		t.Errorf("expected wrapped cause %v, got %v", genericErr, returned)
	}

	exitErr := &exec.ExitError{}
	returnedExit := handleCargoExitError(exitErr)
	if returnedExit == nil || returnedExit.Cause != exitErr {
		t.Errorf("expected wrapped exitErr cause, got %v", returnedExit)
	}
	var asExit *exec.ExitError
	if !errors.As(returnedExit, &asExit) || asExit != exitErr {
		t.Errorf("expected errors.As to recover exitErr, got %v", asExit)
	}
	if returnedExit.Ctx["exit_code"] != exitErr.ExitCode() {
		t.Errorf("expected exit_code context %d, got %v", exitErr.ExitCode(), returnedExit.Ctx["exit_code"])
	}
}

func TestEnrichCargoEnv(t *testing.T) {
	testBin := filepath.Join("test", "bin", "cargo")
	expectedDir := filepath.Dir(testBin)
	env := enrichCargoEnv(testBin)
	hasPath := false
	for _, e := range env {
		if !strings.HasPrefix(e, "PATH=") {
			continue
		}

		hasPath = true
		if !strings.Contains(e, expectedDir) {
			t.Errorf("expected PATH to contain %s, got %s", expectedDir, e)
		}
	}
	if !hasPath {
		t.Errorf("expected enriched env to contain PATH")
	}
}
