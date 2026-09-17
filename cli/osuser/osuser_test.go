package osuser

import (
	"os/exec"
	"strings"
	"testing"
)

func TestBuildWindowsNetUserArgs(t *testing.T) {
	argsWithPass := buildWindowsNetUserArgs("alice", "secret")
	if len(argsWithPass) != 4 || argsWithPass[2] != "secret" {
		t.Fatalf("expected password args, got: %v", argsWithPass)
	}
	argsNoPass := buildWindowsNetUserArgs("bob", "")
	if len(argsNoPass) != 3 {
		t.Fatalf("expected 3 args without password, got: %v", argsNoPass)
	}
}

func TestBuildLinuxUserAddArgs(t *testing.T) {
	opts := UserCreateOptions{
		Username: "alice",
		HomeDir:  "/home/custom",
		Shell:    "/bin/bash",
	}
	args := buildLinuxUserAddArgs(opts)
	if len(args) != 6 {
		t.Fatalf("unexpected args length: %v", args)
	}
}

func TestFilterRetainedSudoersLines(t *testing.T) {
	content := "root ALL=(ALL:ALL) ALL\nalice ALL=(ALL) NOPASSWD:ALL\nbob ALL=(ALL) ALL\n"
	lines := filterRetainedSudoersLines(content, "alice")
	output := strings.Join(lines, "\n")
	if strings.Contains(output, "alice") {
		t.Fatalf("expected alice to be removed, got: %s", output)
	}
	if !strings.Contains(output, "root") || !strings.Contains(output, "bob") {
		t.Fatalf("expected root and bob retained, got: %s", output)
	}
}

func TestKeyDeduplication(t *testing.T) {
	lines := []string{"ssh-rsa AAAA...", "ssh-ed25519 BBBB..."}
	if !hasExistingKey(lines, "ssh-rsa AAAA...") {
		t.Fatal("expected existing key to be found")
	}
	if hasExistingKey(lines, "ssh-rsa CCCC...") {
		t.Fatal("unexpected key found")
	}
}

func TestIsProtectedUser(t *testing.T) {
	if !isProtectedUser("root") || !isProtectedUser("0") || !isProtectedUser("Administrator") {
		t.Fatal("expected root and administrator to be protected")
	}
	if isProtectedUser("regularuser") {
		t.Fatal("regularuser should not be protected")
	}
}

func TestBuildDeluserAndUserdelArgs(t *testing.T) {
	deluserArgs := buildDeluserArgs("testuser", true)
	if len(deluserArgs) != 2 || deluserArgs[0] != "--remove-home" {
		t.Fatalf("expected --remove-home, got: %v", deluserArgs)
	}
	userdelArgs := buildUserdelArgs("testuser", true)
	if len(userdelArgs) != 2 || userdelArgs[0] != "-r" {
		t.Fatalf("expected -r, got: %v", userdelArgs)
	}
}

func TestValidateCreateOptions(t *testing.T) {
	emptyOpts := UserCreateOptions{}
	if err := validateCreateOptions(emptyOpts); err == nil {
		t.Fatal("expected error on empty username")
	}
	validOpts := UserCreateOptions{Username: "valid"}
	if err := validateCreateOptions(validOpts); err != nil {
		t.Fatalf("unexpected error on valid options: %v", err)
	}
}

type mockCmdRunner struct {
	calls []*exec.Cmd
}

func setupMockRunner() (*mockCmdRunner, func()) {
	orig := defaultOSCommandRunner
	mock := &mockCmdRunner{}
	defaultOSCommandRunner = func(cmd *exec.Cmd) ([]byte, error) {
		mock.calls = append(mock.calls, cmd)
		return []byte("success"), nil
	}
	return mock, func() {
		defaultOSCommandRunner = orig
	}
}

func TestCreateRootUser_Mocked(t *testing.T) {
	mock, cleanup := setupMockRunner()
	defer cleanup()

	opts := UserCreateOptions{Username: "alice", Password: "pwd"}
	if err := CreateRootUser(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected mock runner calls, got 0")
	}
}

func TestRemoveEnhancedUser_Mocked(t *testing.T) {
	mock, cleanup := setupMockRunner()
	defer cleanup()

	opts := UserRemoveOptions{Username: "bob"}
	if err := RemoveEnhancedUser(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected mock runner calls for remove, got 0")
	}
}

func TestKillUserProcesses_Mocked(t *testing.T) {
	mock, cleanup := setupMockRunner()
	defer cleanup()

	opts := UserKillOptions{Username: "charlie", IsForce: true}
	if err := KillUserProcesses(opts); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(mock.calls) == 0 {
		t.Fatal("expected mock runner calls for kill, got 0")
	}
}
