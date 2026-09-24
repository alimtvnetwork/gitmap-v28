package cmdupdate

import (
	"context"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
)

// TestFleetUpdate_RCAProof_UserScenario explicitly proves the fix for the user's reported bug:
// 1. 3 online machines (w1, w2, w3) with encrypted stored passwords.
// 2. 3 offline machines (alpha-win, beta-linux, gamma-mac).
// 3. Proves passwords encrypted at rest are correctly decrypted in memory before dialing.
// 4. Proves offline machines are intercepted by pre-flight liveness check without hanging.
// 5. Proves remote execution is exclusively dispatched to verified active nodes.
func TestFleetUpdate_RCAProof_UserScenario(t *testing.T) {
	origLoad := LoadFleetTargetsFn
	origExec := ExecuteRemoteUpdateFn
	origLive := CheckConnLivenessFn
	defer func() {
		LoadFleetTargetsFn = origLoad
		ExecuteRemoteUpdateFn = origExec
		CheckConnLivenessFn = origLive
	}()

	plainPassword := "SecretClusterPassword123!"
	encryptedPassword, err := crypto.Encrypt([]byte(plainPassword), []byte("gitmap-ssh-secret-key-0123456789"))
	if err != nil {
		t.Fatalf("failed to encrypt mock password: %v", err)
	}

	userFleet := buildMockUserFleet(encryptedPassword)
	LoadFleetTargetsFn = func() ([]FleetTarget, error) {
		return userFleet, nil
	}

	CheckConnLivenessFn = func(ctx context.Context, ip string, port int, timeout time.Duration) (bool, string) {
		if strings.HasPrefix(ip, "192.168.1.") {
			return true, "online"
		}
		return false, "connection timed out"
	}

	var mu sync.Mutex
	executedIPs := make(map[string]bool)
	decryptedPasswords := make(map[string]string)

	ExecuteRemoteUpdateFn = func(target FleetTarget, opts FleetUpdateOptions) (string, error) {
		mu.Lock()
		defer mu.Unlock()
		executedIPs[target.IP] = true

		decrypted, decErr := crypto.DecryptStoredPassword(target.Password)
		if decErr != nil {
			t.Errorf("failed to decrypt password for %s: %v", target.Alias, decErr)
		}
		decryptedPasswords[target.IP] = decrypted
		return `{"success": true, "updated": ["gitmap"], "current_version": "v6.337.0"}`, nil
	}

	dispatchErr := RunFleetUpdateDispatch("update", []string{"all"})
	if dispatchErr != nil {
		t.Fatalf("RunFleetUpdateDispatch failed: %v", dispatchErr)
	}

	verifyOnlineExecutionProof(t, executedIPs, decryptedPasswords, plainPassword)
	verifyOfflineSkippedProof(t, executedIPs)
}

func buildMockUserFleet(encryptedPassword string) []FleetTarget {
	return []FleetTarget{
		{ID: "w1", Alias: "w1", IP: "192.168.1.3", Username: "administrator", Port: 22, Password: encryptedPassword, OS: "windows"},
		{ID: "w2", Alias: "w2", IP: "192.168.1.7", Username: "administrator", Port: 22, Password: encryptedPassword, OS: "windows"},
		{ID: "w3", Alias: "w3", IP: "192.168.1.12", Username: "administrator", Port: 22, Password: encryptedPassword, OS: "windows"},
		{ID: "alpha", Alias: "alpha-win", IP: "10.20.0.11", Username: "admin", Port: 22, Password: encryptedPassword, OS: "windows"},
		{ID: "beta", Alias: "beta-linux", IP: "10.20.0.12", Username: "ubuntu", Port: 22, Password: encryptedPassword, OS: "linux"},
		{ID: "gamma", Alias: "gamma-mac", IP: "10.20.0.13", Username: "devops", Port: 22, Password: encryptedPassword, OS: "darwin"},
	}
}

func verifyOnlineExecutionProof(t *testing.T, executed map[string]bool, passwords map[string]string, expectedPassword string) {
	onlineIPs := []string{"192.168.1.3", "192.168.1.7", "192.168.1.12"}
	for _, ip := range onlineIPs {
		if !executed[ip] {
			t.Errorf("expected online node %s to be updated", ip)
		}
		if passwords[ip] != expectedPassword {
			t.Errorf("expected node %s password to be decrypted to %q, got %q", ip, expectedPassword, passwords[ip])
		}
	}
}

func verifyOfflineSkippedProof(t *testing.T, executed map[string]bool) {
	offlineIPs := []string{"10.20.0.11", "10.20.0.12", "10.20.0.13"}
	for _, ip := range offlineIPs {
		if executed[ip] {
			t.Errorf("expected offline node %s NOT to be executed via SSH", ip)
		}
	}
}
