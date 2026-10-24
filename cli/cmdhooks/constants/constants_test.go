package constants

import (
	"strings"
	"testing"
)

func TestConstantsInvariants(t *testing.T) {
	testCommandConstants(t)
	testFlagConstants(t)
	testPathAndPermConstants(t)
	testTemplateAndGuardConstants(t)
	testHelpAndMessageConstants(t)
	testErrorConstants(t)
}

func testCommandConstants(t *testing.T) {
	if CmdHooks != "hooks" || CmdHooksAlias != "hkin" {
		t.Error("cmdhooks command and alias mismatch")
	}
	if SubCmdInstall != "install" || SubCmdUninstall != "uninstall" {
		t.Error("subcommands install or uninstall mismatch")
	}
	if SubCmdRun != "run" || SubCmdCheck != "check" || SubCmdList != "list" {
		t.Error("subcommands run/check/list mismatch")
	}
	if SubCmdStatus != "status" || SubCmdHelp != "help" {
		t.Error("subcommands status/help mismatch")
	}
}

func testFlagConstants(t *testing.T) {
	if FlagAll != "all" || FlagWorkdir != "workdir" || FlagDryRun != "dry-run" {
		t.Error("flags all/workdir/dry-run mismatch")
	}
	if FlagYes != "yes" || FlagHelp != "help" {
		t.Error("flags yes/help mismatch")
	}
	if FlagArgAll != "--all" || FlagArgWorkdir != "--workdir" {
		t.Error("flag args all/workdir mismatch")
	}
	if FlagArgDryRun != "--dry-run" || FlagArgYes != "--yes" {
		t.Error("flag args dry-run/yes mismatch")
	}
	if FlagArgHelp != "--help" || FlagArgHelpShort != "-h" {
		t.Error("flag args help mismatch")
	}
	if FlagDescAll == "" || FlagDescWorkdir == "" || FlagSetHooksInstall == "" {
		t.Error("descriptions or flagset must not be empty")
	}
}

func testPathAndPermConstants(t *testing.T) {
	if DirGit != ".git" || DirHooks != "hooks" || FileCommitMsg != "commit-msg" {
		t.Error("directory or file constants mismatch")
	}
	if DefaultReposDir != "repos" || DefaultWorkspaceDir != "workspace" {
		t.Error("default dir constants mismatch")
	}
	if SentinelCurrentDir != "." || PathHookCommitMsg != ".git/hooks/commit-msg" {
		t.Error("sentinel or path mismatch")
	}
	if EnvGitmapWorkdir != "GITMAP_WORKDIR" || HookFilePerm != 0755 {
		t.Error("env or file perm mismatch")
	}
}

func testTemplateAndGuardConstants(t *testing.T) {
	if GuardMarker != "gitmap commit-msg guard" {
		t.Error("guard marker constant mismatch")
	}
	if !strings.Contains(CommitMsgGuardTemplate, GuardMarker) {
		t.Error("template must contain guard marker")
	}
	if !strings.Contains(CommitMsgGuardTemplate, "#!/bin/sh") {
		t.Error("template must contain shebang")
	}
}

func testHelpAndMessageConstants(t *testing.T) {
	if !strings.Contains(HelpTextHooksInstall, "gitmap hooks install") {
		t.Error("help text must mention hooks install")
	}
	if DetailAlreadyGuarded != "guard already installed" {
		t.Error("detail already guarded mismatch")
	}
	if DetailHookRelPath != ".git/hooks/commit-msg" {
		t.Error("detail hook rel path mismatch")
	}
	if MsgDryRunNotice == "" || MsgNoReposFoundFmt == "" {
		t.Error("messages must not be empty")
	}
	if MsgPlanHeaderFmt == "" || MsgPlanRowFmt == "" || MsgInstallSummaryFmt == "" {
		t.Error("formatted messages must not be empty")
	}
}

func testErrorConstants(t *testing.T) {
	if ErrMissingSubcommand == "" || ErrUnknownSubcommandFmt == "" {
		t.Error("error messages must not be empty")
	}
	if ErrMissingAllFlag != "missing required flag: --all" {
		t.Error("missing all flag message mismatch")
	}
	if ErrWrapFlags != "hooks.install.flags" || ErrWrapWritePrefix != "hooks.install.write: " {
		t.Error("wrap context prefixes mismatch")
	}
}
