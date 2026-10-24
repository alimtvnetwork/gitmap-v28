package gitutil

import (
	"strings"
	"testing"
	"time"
)

func makeTestBinding(path, keyPath string, priority int, isActive bool) FolderSSHBinding {
	return FolderSSHBinding{
		FolderPath: path,
		KeyPath:    keyPath,
		Priority:   priority,
		IsActive:   isActive,
		UpdatedAt:  time.Now(),
	}
}

func TestGitSSH_LongestPrefixResolution(t *testing.T) {
	bindings := []FolderSSHBinding{
		makeTestBinding("D:/work", "key_work", 100, true),
		makeTestBinding("D:/work/clientA", "key_client_a", 100, true),
	}

	best, hasMatch := FindLongestPrefixBinding("D:/work/clientA/repo1", bindings)
	if !hasMatch {
		t.Fatalf("expected match, got none")
	}

	if best.KeyPath != "key_client_a" {
		t.Errorf("expected key_client_a, got %s", best.KeyPath)
	}
}

func TestGitSSH_DeepPrefixResolution(t *testing.T) {
	bindings := []FolderSSHBinding{
		makeTestBinding("D:/work", "key_work", 100, true),
		makeTestBinding("D:/work/clientA", "key_client_a", 100, true),
		makeTestBinding("D:/work/clientA/special", "key_special", 100, true),
	}

	best, hasMatch := FindLongestPrefixBinding("D:/work/clientA/special/subrepo", bindings)
	if !hasMatch || best.KeyPath != "key_special" {
		t.Errorf("expected key_special, got %s", best.KeyPath)
	}
}

func TestGitSSH_BoundarySafety(t *testing.T) {
	bindings := []FolderSSHBinding{
		makeTestBinding("D:/work", "key_work", 100, true),
		makeTestBinding("D:/work/clientA", "key_client_a", 100, true),
	}

	best, hasMatch := FindLongestPrefixBinding("D:/work/clientA-extra/repo1", bindings)
	if !hasMatch {
		t.Fatalf("expected match with work, got none")
	}

	if best.KeyPath != "key_work" {
		t.Errorf("expected key_work boundary fallback, got %s", best.KeyPath)
	}
}

func TestGitSSH_CaseInsensitiveMatching(t *testing.T) {
	bindings := []FolderSSHBinding{
		makeTestBinding("D:/work/clienta", "key_client_a", 100, true),
	}

	best, hasMatch := FindLongestPrefixBinding("d:\\Work\\ClientA\\Repo", bindings)
	if !hasMatch {
		t.Fatalf("expected case-insensitive match on Windows drive, got none")
	}

	if best.KeyPath != "key_client_a" {
		t.Errorf("expected key_client_a, got %s", best.KeyPath)
	}
}

func TestGitSSH_InactiveBindingsIgnored(t *testing.T) {
	bindings := []FolderSSHBinding{
		makeTestBinding("D:/work", "key_work", 100, true),
		makeTestBinding("D:/work/clientA", "key_inactive", 100, false),
	}

	best, hasMatch := FindLongestPrefixBinding("D:/work/clientA/repo1", bindings)
	if !hasMatch {
		t.Fatalf("expected match with D:/work, got none")
	}

	if best.KeyPath != "key_work" {
		t.Errorf("expected key_work because clientA is inactive, got %s", best.KeyPath)
	}
}

func TestGitSSH_PriorityTieBreaking(t *testing.T) {
	bindings := []FolderSSHBinding{
		makeTestBinding("D:/work/shared", "key_low", 50, true),
		makeTestBinding("D:/work/shared", "key_high", 200, true),
	}

	best, hasMatch := FindLongestPrefixBinding("D:/work/shared/repo", bindings)
	if !hasMatch {
		t.Fatalf("expected match, got none")
	}

	if best.KeyPath != "key_high" {
		t.Errorf("expected higher priority key_high, got %s", best.KeyPath)
	}
}

func TestGitSSH_TimestampTieBreaking(t *testing.T) {
	now := time.Now()
	bOlder := FolderSSHBinding{FolderPath: "D:/work/tie", KeyPath: "key_old", Priority: 100, IsActive: true, UpdatedAt: now.Add(-time.Hour)}
	bNewer := FolderSSHBinding{FolderPath: "D:/work/tie", KeyPath: "key_new", Priority: 100, IsActive: true, UpdatedAt: now}

	best, hasMatch := FindLongestPrefixBinding("D:/work/tie/repo", []FolderSSHBinding{bOlder, bNewer})
	if !hasMatch {
		t.Fatalf("expected match, got none")
	}

	if best.KeyPath != "key_new" {
		t.Errorf("expected newer timestamp key_new, got %s", best.KeyPath)
	}
}

func TestGitSSH_BuildDeterministicSSHCommand(t *testing.T) {
	cmd := BuildDeterministicSSHCommand("C:\\Users\\dev\\.ssh\\id_ed25519")
	expected := `ssh -F /dev/null -o IdentitiesOnly=yes -i "C:/Users/dev/.ssh/id_ed25519" -o StrictHostKeyChecking=accept-new -o BatchMode=yes -o ConnectTimeout=15`

	if cmd != expected {
		t.Errorf("command mismatch:\nexpected: %s\ngot:      %s", expected, cmd)
	}
}

func TestGitSSH_BuildDeterministicSSHCommand_SpacesInPath(t *testing.T) {
	cmd := BuildDeterministicSSHCommand("C:\\Users\\John Doe\\.ssh\\id_rsa")
	expected := `ssh -F /dev/null -o IdentitiesOnly=yes -i "C:/Users/John Doe/.ssh/id_rsa" -o StrictHostKeyChecking=accept-new -o BatchMode=yes -o ConnectTimeout=15`

	if cmd != expected {
		t.Errorf("spaces in path command mismatch:\nexpected: %s\ngot:      %s", expected, cmd)
	}
}

func testGitSSH_AssertAuthorEnv(t *testing.T, env []string) {
	assertEnvContains(t, env, "GIT_TERMINAL_PROMPT=0")
	assertEnvContains(t, env, "GIT_AUTHOR_NAME=Alice Engineer")
	assertEnvContains(t, env, "GIT_AUTHOR_EMAIL=alice@clienta.com")
	assertEnvContains(t, env, "CUSTOM_VAR=hello")
}

func TestGitSSH_BuildSafeGitEnvForRepo(t *testing.T) {
	ctx := ResolvedGitContext{
		RepoPath:       "D:/work/clientA/repo1",
		GitSSHCommand:  BuildDeterministicSSHCommand("C:/keys/id_ed25519"),
		AuthorName:     "Alice Engineer",
		AuthorEmail:    "alice@clienta.com",
		HasInjectedSSH: true,
	}

	env := BuildSafeGitEnvForRepo(ctx, "CUSTOM_VAR=hello")
	sshCmd, hasExtracted := ExtractGitSSHCommand(env)
	if !hasExtracted || !strings.Contains(sshCmd, "id_ed25519") {
		t.Errorf("failed extracting valid ssh command from env: %s", sshCmd)
	}

	testGitSSH_AssertAuthorEnv(t, env)
}

func TestGitSSH_BuildSafeGitEnvForRepo_NoSSH(t *testing.T) {
	ctx := ResolvedGitContext{RepoPath: "D:/work/repo"}
	env := BuildSafeGitEnvForRepo(ctx)

	if HasGitSSHCommand(env) {
		t.Errorf("expected no GIT_SSH_COMMAND when HasInjectedSSH is false")
	}

	assertEnvContains(t, env, "GIT_TERMINAL_PROMPT=0")
	assertEnvContains(t, env, "GCM_INTERACTIVE=never")
}

func TestGitSSH_ResolveGitContextFromBindings(t *testing.T) {
	b := FolderSSHBinding{
		FolderPath: "D:/work/corp", KeyPath: "C:/keys/corp_key",
		AuthorName: "Bob Corp", AuthorEmail: "bob@corp.com", Priority: 100, IsActive: true,
	}

	ctx := ResolveGitContextFromBindings("D:/work/corp/service", []FolderSSHBinding{b})
	if !ctx.HasInjectedSSH {
		t.Fatalf("expected HasInjectedSSH to be true")
	}

	if ctx.AuthorName != "Bob Corp" || ctx.AuthorEmail != "bob@corp.com" {
		t.Errorf("unexpected author data: %s <%s>", ctx.AuthorName, ctx.AuthorEmail)
	}
}

func TestGitSSH_CreateGitCmdInDirWithContext(t *testing.T) {
	ctx := ResolvedGitContext{
		GitSSHCommand:  BuildDeterministicSSHCommand("C:/key"),
		HasInjectedSSH: true,
	}

	cmd := CreateGitCmdInDirWithContext("D:/repo", ctx, "status")
	if len(cmd.Args) < 4 || cmd.Args[1] != "-C" || cmd.Args[2] != "D:/repo" {
		t.Errorf("unexpected cmd args: %v", cmd.Args)
	}

	if !HasGitSSHCommand(cmd.Env) {
		t.Errorf("cmd.Env missing GIT_SSH_COMMAND")
	}
}

func assertEnvContains(t *testing.T, env []string, target string) {
	for _, v := range env {
		if v == target {
			return
		}
	}

	t.Errorf("expected env to contain %q", target)
}
