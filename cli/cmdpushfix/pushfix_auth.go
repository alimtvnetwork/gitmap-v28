package cmdpushfix

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/alimtvnetwork/gitmap-v28/cli/cmdclone"
	"github.com/alimtvnetwork/gitmap-v28/cli/cmdssh"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

// ProbeSSHConnection tests SSH public key authentication against the git host.
func ProbeSSHConnection(host string) (bool, string) {
	if len(host) == 0 {
		host = "github.com"
	}
	cmd := exec.Command("ssh", "-Tv", "-o", "BatchMode=yes", "-o", "ConnectTimeout=5", "git@"+host)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	_ = cmd.Run()
	output := buf.String()
	isSuccess := strings.Contains(output, "successfully authenticated")
	return isSuccess, output
}

// RemediateSSHKeyAndConfig verifies local keys and auto-synchronizes SSH config.
func RemediateSSHKeyAndConfig(host string) []string {
	var remediations []string
	isSuccess, _ := ProbeSSHConnection(host)
	if isSuccess {
		return remediations
	}

	remediations = append(remediations, attemptCandidateKeySync()...)
	remediations = append(remediations, purgeBogusDashedKeys()...)
	remediations = append(remediations, syncAllSSHConfigs()...)
	return remediations
}

func attemptCandidateKeySync() []string {
	var actions []string
	altPath := filepath.Join(getSystemSSHDir(), "id_rsa_-y")
	if _, err := os.Stat(altPath); err != nil {
		return actions
	}

	destPriv := filepath.Join(getSystemSSHDir(), "id_rsa")
	destPub := filepath.Join(getSystemSSHDir(), "id_rsa.pub")
	if copyFileSafe(altPath, destPriv) && copyFileSafe(altPath+".pub", destPub) {
		actions = append(actions, "Promoted working key id_rsa_-y to default id_rsa")
		PrintRemediationApplied("Synchronized authenticated key pair to standard id_rsa")
	}
	return actions
}

func purgeBogusDashedKeys() []string {
	var actions []string
	db, err := store.OpenDefault()
	if err != nil {
		return actions
	}
	defer db.Close()

	keys, _ := db.ListSSHKeys()
	for _, k := range keys {
		if strings.HasPrefix(k.Name, "-") {
			_ = db.DeleteSSHKey(k.Name)
			actions = append(actions, fmt.Sprintf("Purged invalid dashed key %q from SQLite store", k.Name))
			PrintRemediationApplied(fmt.Sprintf("Cleaned up transient dashed key %q from DB", k.Name))
		}
	}
	return actions
}

func syncAllSSHConfigs() []string {
	_ = cmdssh.RunSSHConfig(nil)
	PrintRemediationApplied("Regenerated canonical Host github.com across SSH profiles")
	return []string{"Regenerated ~/.ssh/config and system profile config"}
}

func getSystemSSHDir() string {
	drive := os.Getenv("SystemDrive")
	user := os.Getenv("USERNAME")
	if len(drive) > 0 && len(user) > 0 && runtime.GOOS == "windows" {
		return filepath.Join(drive, "\\Users", user, ".ssh")
	}
	home, _ := os.UserHomeDir()
	return filepath.Join(home, ".ssh")
}

func copyFileSafe(src, dst string) bool {
	data, err := os.ReadFile(src)
	if err != nil {
		return false
	}
	return os.WriteFile(dst, data, 0o600) == nil
}

// ConvertRemoteToSSH flips a repository's remote URL from HTTPS to SSH.
func ConvertRemoteToSSH(repoDir, remoteName, currentURL string) bool {
	sshURL, ok := cmdclone.ConvertURLToSSH(currentURL)
	if !ok || sshURL == currentURL {
		return false
	}
	cmd := exec.Command("git", "-C", repoDir, "remote", "set-url", remoteName, sshURL)
	if err := cmd.Run(); err != nil {
		return false
	}
	PrintRemediationApplied(fmt.Sprintf("Converted %s from HTTPS to SSH (%s)", remoteName, sshURL))
	return true
}
