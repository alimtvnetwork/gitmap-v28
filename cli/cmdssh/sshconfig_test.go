package cmdssh

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSanitizeSSHConfigContent_AuthorizedKeysFile(t *testing.T) {
	input := `Host github.com
    HostName github.com
    User git
    IdentityFile ~/.ssh/id_rsa
    IdentitiesOnly yes

Host server
    HostName 192.168.1.3
    User git
    IdentityFile ~/.ssh/id_rsa
    IdentitiesOnly yes
    PasswordAuthentication no
    pubkeyauthentication yes
    authorizedkeysfile ~/.ssh/authorized_keys
`

	sanitized, count := SanitizeSSHConfigContent(input)
	if count != 1 {
		t.Fatalf("expected 1 removed directive, got %d", count)
	}

	if strings.Contains(sanitized, "\n    authorizedkeysfile") {
		t.Fatalf("expected authorizedkeysfile to be removed or commented, got:\n%s", sanitized)
	}

	if !strings.Contains(sanitized, "[gitmap removed server-only directive]: authorizedkeysfile") {
		t.Fatalf("expected gitmap removal comment, got:\n%s", sanitized)
	}

	if !strings.Contains(sanitized, "Host github.com") || !strings.Contains(sanitized, "Host server") {
		t.Fatalf("expected host blocks to be preserved")
	}
}

func TestSanitizeSSHConfigContent_MultipleDaemonDirectives(t *testing.T) {
	input := `PermitRootLogin no
Subsystem sftp /usr/lib/openssh/sftp-server
ClientAliveInterval 300
ClientAliveCountMax 3
StrictModes yes

Host devbox
    HostName 10.0.0.5
    User dev
`

	sanitized, count := SanitizeSSHConfigContent(input)
	if count != 5 {
		t.Fatalf("expected 5 removed directives, got %d", count)
	}

	if !strings.Contains(sanitized, "Host devbox") {
		t.Fatalf("expected Host devbox to be preserved")
	}
}

func TestSanitizeSSHConfigFile_Hermetic(t *testing.T) {
	tmpDir := t.TempDir()
	configPath := filepath.Join(tmpDir, "config")

	content := `Host my-server
    HostName 192.168.1.50
    authorizedkeysfile /etc/ssh/authorized_keys
`
	if err := os.WriteFile(configPath, []byte(content), 0o600); err != nil {
		t.Fatalf("failed to write test config: %v", err)
	}

	changed, appErr := SanitizeSSHConfigFile(configPath)
	if appErr != nil {
		t.Fatalf("unexpected error: %v", appErr)
	}
	if !changed {
		t.Fatalf("expected changed=true")
	}

	data, err := os.ReadFile(configPath)
	if err != nil {
		t.Fatalf("failed to read back config: %v", err)
	}

	if strings.Contains(string(data), "\n    authorizedkeysfile") {
		t.Fatalf("authorizedkeysfile was not cleaned from file")
	}

	// Second run should be a no-op
	changedAgain, appErrAgain := SanitizeSSHConfigFile(configPath)
	if appErrAgain != nil {
		t.Fatalf("unexpected error on second run: %v", appErrAgain)
	}
	if changedAgain {
		t.Fatalf("expected changed=false on second run")
	}
}
