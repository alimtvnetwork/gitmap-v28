package constants

import (
	"strings"
	"testing"
)

func TestCommandConstants(t *testing.T) {
	if CmdLogin != "login" {
		t.Errorf("expected CmdLogin == 'login', got %q", CmdLogin)
	}

	if CmdLogout != "logout" {
		t.Errorf("expected CmdLogout == 'logout', got %q", CmdLogout)
	}

	subcmds := []string{
		SubCmdLoginSSH,
		SubCmdLoginBrowser,
		SubCmdLoginBrowserFlow,
		SubCmdLoginToken,
		SubCmdLoginPAT,
		SubCmdLoginStatus,
	}

	for _, sub := range subcmds {
		if len(sub) == 0 {
			t.Errorf("expected subcommand constant to be non-empty")
		}
	}
}

func TestFlagConstants(t *testing.T) {
	flags := []string{
		FlagHelp,
		FlagHelpArg,
		FlagHelpShortArg,
		FlagStatus,
		FlagStatusArg,
		FlagWeb,
		FlagWebArg,
		FlagBrowser,
		FlagBrowserArg,
		FlagNoVerify,
		FlagNoVerifyArg,
		FlagToken,
		FlagTokenArg,
		FlagTokenPrefix,
		FlagYes,
		FlagYesArg,
		FlagYesShortArg,
		FlagForce,
		FlagForceArg,
		FlagForceShortArg,
		FlagDir,
		FlagDirArg,
		FlagDirShortArg,
		FlagDirPrefix,
	}

	for _, f := range flags {
		if len(f) == 0 {
			t.Errorf("expected flag constant to be non-empty")
		}
	}

	if !strings.HasPrefix(FlagTokenPrefix, "--token=") {
		t.Errorf("expected FlagTokenPrefix to start with '--token='")
	}

	if !strings.HasPrefix(FlagDirPrefix, "--dir=") {
		t.Errorf("expected FlagDirPrefix to start with '--dir='")
	}
}

func TestGitConfigConstants(t *testing.T) {
	if GitConfigGithubToken != "github.token" {
		t.Errorf("expected GitConfigGithubToken == 'github.token', got %q", GitConfigGithubToken)
	}

	if GitConfigCoreSSHCommand != "core.sshCommand" {
		t.Errorf("expected GitConfigCoreSSHCommand == 'core.sshCommand', got %q", GitConfigCoreSSHCommand)
	}

	if len(GitCredentialRejectPayload) == 0 {
		t.Errorf("expected GitCredentialRejectPayload to be non-empty")
	}
}

func TestEndpointAndURLConstants(t *testing.T) {
	endpoints := []string{
		GithubUserAPI,
		GithubNewTokenURL,
		GithubNewSSHKeyURL,
		GitlabNewSSHKeyURL,
	}

	for _, ep := range endpoints {
		if !strings.HasPrefix(ep, "https://") {
			t.Errorf("expected HTTPS URL, got %q", ep)
		}
	}
}

func TestToolAndClipboardConstants(t *testing.T) {
	tools := []string{
		BinGh,
		BinSshKeygen,
		ClipWindows,
		ClipDarwin,
		ClipWlCopy,
		ClipXclip,
		ClipXsel,
	}

	for _, tool := range tools {
		if len(tool) == 0 {
			t.Errorf("expected tool name to be non-empty")
		}
	}

	if len(BannerDivider) != 80 {
		t.Errorf("expected 80-char banner divider, got len %d", len(BannerDivider))
	}
}
