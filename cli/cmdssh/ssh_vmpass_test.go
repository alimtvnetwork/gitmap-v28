package cmdssh

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

func TestResolveFallbackCredentials_SaltedAndPlain(t *testing.T) {
	tempDir := t.TempDir()
	origWd, err := os.Getwd()
	if err != nil {
		t.Fatalf("getwd: %v", err)
	}
	defer func() {
		_ = os.Chdir(origWd)
	}()

	if chErr := os.Chdir(tempDir); chErr != nil {
		t.Fatalf("chdir: %v", chErr)
	}

	winPass := "rtyrty123@"
	linuxPass := "ubuntu-secret-42"
	saltWin := "9f8b2c4e"
	encWin := crypto.EncryptSalted(winPass, saltWin)
	encLinux := crypto.EncryptCaesar(linuxPass, 15, "salttest")

	jsonContent := `{
  "attributes": {
    "type": "vm-credentials",
    "version": "2.0",
    "workDirectory": {
      "path": "${workDir}",
      "isApplied": true
    }
  },
  "variables": {
    "workDir": "D:\\work",
    "adminUser": "Administrator",
    "winPass": "` + encWin + `",
    "linuxUser": "a",
    "linuxPass": "` + encLinux + `"
  },
  "data": {
    "windows": {
      "user": "${adminUser}",
      "pass": "${winPass}"
    },
    "ubuntu": {
      "user": "${linuxUser}",
      "pass": "${linuxPass}"
    }
  }
}`

	if writeErr := os.WriteFile(filepath.Join(tempDir, "vmpass.json"), []byte(jsonContent), 0644); writeErr != nil {
		t.Fatalf("write temp vmpass.json: %v", writeErr)
	}

	resolvedWin := ResolveFallbackCredentials("Administrator", "windows")
	if resolvedWin != winPass {
		t.Errorf("expected decrypted windows pass %q, got %q", winPass, resolvedWin)
	}

	resolvedUbuntu := ResolveFallbackCredentials("a", "linux")
	if resolvedUbuntu != linuxPass {
		t.Errorf("expected decrypted ubuntu pass %q, got %q", linuxPass, resolvedUbuntu)
	}
}
