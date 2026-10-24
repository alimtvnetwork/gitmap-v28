package constants

import (
	"strings"
	"testing"
)

func TestCommandConstants(t *testing.T) {
	if CmdRevert != "revert" {
		t.Errorf("expected 'revert', got %q", CmdRevert)
	}

	if CmdRevertRunner != "revert-runner" {
		t.Errorf("expected 'revert-runner', got %q", CmdRevertRunner)
	}
}

func TestFlagConstants(t *testing.T) {
	flags := []string{
		FlagListTxn, FlagShowTxn, FlagPruneTxn, FlagTxn,
		FlagLastTxn, FlagLastNTxn, FlagForce, FlagVerbose,
	}

	for _, f := range flags {
		if len(f) == 0 {
			t.Error("flag name should not be empty")
		}
	}

	args := []string{
		FlagArgListTxn, FlagArgShowTxn, FlagArgPruneTxn, FlagArgTxn,
		FlagArgLastTxn, FlagArgLastNTxn, FlagArgForce, FlagArgVerbose,
	}

	for _, a := range args {
		if !strings.HasPrefix(a, "--") {
			t.Errorf("flag arg %q must start with '--'", a)
		}
	}
}

func TestRetentionAndBinaryConstants(t *testing.T) {
	if TxnRetentionCap != 50 {
		t.Errorf("expected retention cap 50, got %d", TxnRetentionCap)
	}

	if TxnStatusCommitted != "committed" || TxnStatusReverted != "reverted" {
		t.Error("unexpected transaction status constants")
	}

	if ConfirmTokenYes != "yes" {
		t.Errorf("expected confirm token 'yes', got %q", ConfirmTokenYes)
	}

	if BinGit != "git" || BinPowershell != "powershell" {
		t.Error("unexpected binary names")
	}

	if ScriptRunPS1 != "run.ps1" {
		t.Errorf("expected 'run.ps1', got %q", ScriptRunPS1)
	}
}

func TestTemplatesAndMessages(t *testing.T) {
	if !strings.Contains(RevertPSHeader, "Set-Location") {
		t.Error("RevertPSHeader missing Set-Location")
	}

	if !strings.Contains(RevertPSBuild, "Building from checked-out version") {
		t.Error("RevertPSBuild missing build message")
	}

	if !strings.Contains(RevertPSPostActions, "Active version") {
		t.Error("RevertPSPostActions missing version check")
	}

	if len(MsgRevertUsage) == 0 || len(MsgRevertCheckout) == 0 || len(MsgRevertStarting) == 0 {
		t.Error("revert messages should not be empty")
	}
}

func TestErrorCodesAndDomains(t *testing.T) {
	if ErrCodeE9000 != "E9000" || ErrCodeE2012 != "E2012" {
		t.Error("unexpected error codes")
	}

	if DomainCmdRevert != "cmd.revert" || DomainCmdRevertTxn != "cmd.reverttxn" {
		t.Error("unexpected error domains")
	}

	if OpRevertCheckout != "revert.checkout" || OpTxnRevert != "txn.revert" {
		t.Error("unexpected operation tags")
	}
}
