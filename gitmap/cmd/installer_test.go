package cmd

import (
	"bytes"
	"testing"

	"github.com/spf13/cobra"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func TestInstallerCmd(t *testing.T) {
	t.Parallel()

	if installerCmd == nil {
		t.Fatal("expected installerCmd to be initialized")
	}

	if installerCmd.Use != "installer" {
		t.Fatalf("expected Use to be 'installer', got %q", installerCmd.Use)
	}

	// Test success boundary with buffered output
	buf := new(bytes.Buffer)
	installerCmd.SetOut(buf)
	installerCmd.SetErr(buf)

	err := runInstaller(installerCmd, []string{})
	if err != nil {
		t.Fatalf("runInstaller failed: %v", err)
	}

	if buf.Len() == 0 {
		t.Fatal("expected help output to be written to buffer")
	}

	// Test RunE hook
	runEErr := installerCmd.RunE(installerCmd, []string{})
	if runEErr != nil {
		t.Fatalf("installerCmd.RunE failed: %v", runEErr)
	}

	// Test with a fresh dummy command
	dummyCmd := &cobra.Command{Use: "test"}
	dummyBuf := new(bytes.Buffer)
	dummyCmd.SetOut(dummyBuf)
	if errDummy := runInstaller(dummyCmd, []string{}); errDummy != nil {
		t.Fatalf("runInstaller on dummyCmd failed: %v", errDummy)
	}

	// Test failure boundary with nil command
	errNil := runInstaller(nil, []string{})
	if errNil == nil {
		t.Fatal("expected error when passing nil command")
	}

	appErr, ok := errNil.(*apperror.AppError)
	if !ok {
		t.Fatalf("expected *apperror.AppError, got %T", errNil)
	}

	if appErr.Code != "E_INSTALLER_NIL_COMMAND" {
		t.Fatalf("expected error code E_INSTALLER_NIL_COMMAND, got %q", appErr.Code)
	}
}
