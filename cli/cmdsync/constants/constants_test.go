package constants

import (
	"testing"
)

func TestSubcommandConstants(t *testing.T) {
	if CmdSync != "sync" {
		t.Errorf("expected sync, got %s", CmdSync)
	}
	if CmdSyncAlias != "sy" {
		t.Errorf("expected sy, got %s", CmdSyncAlias)
	}
	if CmdCommon != "common" {
		t.Errorf("expected common, got %s", CmdCommon)
	}
	if CmdCommonAlias != "co" {
		t.Errorf("expected co, got %s", CmdCommonAlias)
	}
	if CmdAddLFSInstall != "add-lfs-install" {
		t.Errorf("expected add-lfs-install, got %s", CmdAddLFSInstall)
	}
}

func TestFlagConstants(t *testing.T) {
	if FlagProjects != "projects" || FlagProjectsShort != "p" {
		t.Errorf("unexpected projects flag values")
	}
	if FlagRepo != "repo" || FlagRepoShort != "r" {
		t.Errorf("unexpected repo flag values")
	}
	if FlagWorkers != "workers" || FlagWorkersShort != "w" {
		t.Errorf("unexpected workers flag values")
	}
	if FlagDryRun != "dry-run" || FlagDryRunShort != "n" {
		t.Errorf("unexpected dry-run flag values")
	}
	if FlagNoPush != "no-push" || FlagNoRelease != "no-release" {
		t.Errorf("unexpected push/release flag values")
	}
	if FlagVerbose != "verbose" {
		t.Errorf("unexpected verbose flag value")
	}
	if FlagForce != "force" || FlagForceShort != "f" {
		t.Errorf("unexpected force flag values")
	}
	if FlagHelp != "help" || FlagHelpShort != "h" {
		t.Errorf("unexpected help flag values")
	}
}

func TestTargetAndFileConstants(t *testing.T) {
	if TargetAll != "all" || TargetIgnore != "ignore" || TargetAttributes != "attributes" {
		t.Errorf("unexpected target constants")
	}
	if TargetLFSInstall != "lfs-install" || TargetPrettierIgnore != "prettier-ignore" || TargetPrettierRC != "prettier-rc" {
		t.Errorf("unexpected target prettier/lfs constants")
	}
	if FileGitignore != ".gitignore" || FileGitattributes != ".gitattributes" {
		t.Errorf("unexpected git ignore/attributes file names")
	}
	if FilePrettierignore != ".prettierignore" || FilePrettierRC != ".prettierrc" {
		t.Errorf("unexpected prettier file names")
	}
	if DirGit != ".git" || DirPrompts != "01-prompts" || DirSpec != "02-spec" {
		t.Errorf("unexpected dir tokens")
	}
}

func TestStatusAndActionConstants(t *testing.T) {
	if StatusOk != "OK" || StatusFail != "FAIL" || StatusSkipped != "SKIPPED" {
		t.Errorf("unexpected status constants")
	}
	if ActionNoChanges != "No Changes" || ActionAlreadySynced != "Already Synced" {
		t.Errorf("unexpected action constants")
	}
	if ActionDryRun != "Dry Run" || ActionMirrorFailed != "Mirror Failed" {
		t.Errorf("unexpected mirror/dry-run action constants")
	}
	if DefaultBranch != "main" || DefaultInitialTag != "v0.1.0" {
		t.Errorf("unexpected default git values")
	}
	if DefaultWorkers != 8 || MinWorkers != 4 || MaxWorkers != 16 {
		t.Errorf("unexpected worker thresholds")
	}
}
