package cmdssh

import (
	"context"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestParseDeployConfigSSHFlags(t *testing.T) {
	cases := []struct {
		name     string
		args     []string
		expected DeployConfigSSHOptions
	}{
		{
			name: "default args",
			args: []string{},
			expected: DeployConfigSSHOptions{
				Target: "all",
			},
		},
		{
			name: "target specific node",
			args: []string{"w1"},
			expected: DeployConfigSSHOptions{
				Target: "w1",
			},
		},
		{
			name: "dry run with except and file",
			args: []string{"all", "--dry-run", "--file", "my-nodes.json", "--except", "w2,w3"},
			expected: DeployConfigSSHOptions{
				Target:   "all",
				IsDryRun: true,
				FilePath: "my-nodes.json",
				Except:   "w2,w3",
			},
		},
		{
			name: "json and force flags",
			args: []string{"--json", "--force", "192.168.1.5"},
			expected: DeployConfigSSHOptions{
				Target:  "192.168.1.5",
				IsJSON:  true,
				IsForce: true,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			opts := parseDeployConfigSSHFlags(tc.args)
			if opts.Target != tc.expected.Target {
				t.Errorf("expected Target %q, got %q", tc.expected.Target, opts.Target)
			}
			if opts.IsDryRun != tc.expected.IsDryRun {
				t.Errorf("expected IsDryRun %v, got %v", tc.expected.IsDryRun, opts.IsDryRun)
			}
			if opts.IsJSON != tc.expected.IsJSON {
				t.Errorf("expected IsJSON %v, got %v", tc.expected.IsJSON, opts.IsJSON)
			}
			if opts.IsForce != tc.expected.IsForce {
				t.Errorf("expected IsForce %v, got %v", tc.expected.IsForce, opts.IsForce)
			}
			if opts.FilePath != tc.expected.FilePath {
				t.Errorf("expected FilePath %q, got %q", tc.expected.FilePath, opts.FilePath)
			}
			if opts.Except != tc.expected.Except {
				t.Errorf("expected Except %q, got %q", tc.expected.Except, opts.Except)
			}
		})
	}
}

func TestIsHelpDeployConfigSSHRequest(t *testing.T) {
	if !isHelpDeployConfigSSHRequest([]string{"help"}) {
		t.Errorf("expected true for 'help'")
	}
	if !isHelpDeployConfigSSHRequest([]string{"--help"}) {
		t.Errorf("expected true for '--help'")
	}
	if !isHelpDeployConfigSSHRequest([]string{"-h"}) {
		t.Errorf("expected true for '-h'")
	}
	if !isHelpDeployConfigSSHRequest([]string{"w1", "--help"}) {
		t.Errorf("expected true for 'w1 --help'")
	}
	if isHelpDeployConfigSSHRequest([]string{"w1"}) {
		t.Errorf("expected false for 'w1'")
	}
	if isHelpDeployConfigSSHRequest([]string{}) {
		t.Errorf("expected false for empty args")
	}
}

func TestStripLeadingSSHArg(t *testing.T) {
	res1 := stripLeadingSSHArg([]string{"ssh", "w1"})
	if len(res1) != 1 || res1[0] != "w1" {
		t.Errorf("expected ['w1'], got %v", res1)
	}

	res2 := stripLeadingSSHArg([]string{"SSH", "all"})
	if len(res2) != 1 || res2[0] != "all" {
		t.Errorf("expected ['all'], got %v", res2)
	}

	res3 := stripLeadingSSHArg([]string{"w1"})
	if len(res3) != 1 || res3[0] != "w1" {
		t.Errorf("expected ['w1'], got %v", res3)
	}
}

func TestPortableEncryptConnections(t *testing.T) {
	conns := []db.SSHConnection{
		{
			Alias:     "node-test",
			IPAddress: "192.168.1.99",
			Username:  "Administrator",
			OS:        "windows",
		},
	}
	encrypted := portableEncryptConnections(conns)
	if len(encrypted) != 1 {
		t.Fatalf("expected 1 connection")
	}
	if encrypted[0].EncryptedPassword == "" {
		t.Errorf("expected non-empty encrypted password")
	}
	if !isAESCiphertext(encrypted[0].EncryptedPassword) {
		t.Errorf("expected aes: prefix for portable ciphertext, got: %s", encrypted[0].EncryptedPassword)
	}
	decrypted, err := DecryptSSHPassword(encrypted[0].EncryptedPassword)
	if err != nil {
		t.Fatalf("failed to decrypt portable password: %v", err)
	}
	if decrypted == "" {
		t.Errorf("expected non-empty decrypted password")
	}
}

func TestMergeConnectionChanges_PasswordPreservation(t *testing.T) {
	validPass, err := encryptWithFallbackAES("existing-valid-password")
	if err != nil {
		t.Fatalf("failed to encrypt test password: %v", err)
	}

	existing := db.SSHConnection{
		Alias:             "w1",
		IPAddress:         "192.168.1.3",
		Username:          "Administrator",
		EncryptedPassword: validPass,
		OS:                "windows",
	}

	// Incoming with empty password: must preserve existing password
	incomingEmptyPass := db.SSHConnection{
		Alias:     "w1",
		IPAddress: "192.168.1.30", // changed IP
		Username:  "Administrator",
		OS:        "windows",
	}

	isChanged, merged := mergeConnectionChanges(existing, incomingEmptyPass)
	if !isChanged {
		t.Errorf("expected isChanged = true due to IP change")
	}
	if merged.IPAddress != "192.168.1.30" {
		t.Errorf("expected updated IP '192.168.1.30', got %s", merged.IPAddress)
	}
	if merged.EncryptedPassword != validPass {
		t.Errorf("expected preserved password %q, got %q", validPass, merged.EncryptedPassword)
	}
}

func TestSyncSSHConnectionsLocally_MatchAndUpdate(t *testing.T) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		t.Skip("local database unavailable in test environment")
		return
	}
	defer dbConn.Close()
	_ = store.EnsureSSHTables(dbConn.SQL())

	ctx := context.Background()
	testAlias := "test-sync-node-99"
	testIP := "192.168.99.99"

	// Cleanup test node if exists
	_, _ = store.DeleteHostByAliasOrIP(ctx, testAlias, dbConn.SQL())
	_, _ = dbConn.SQL().ExecContext(ctx, "DELETE FROM SSHConnection WHERE Alias = ?", testAlias)

	defer func() {
		_, _ = store.DeleteHostByAliasOrIP(ctx, testAlias, dbConn.SQL())
		_, _ = dbConn.SQL().ExecContext(ctx, "DELETE FROM SSHConnection WHERE Alias = ?", testAlias)
	}()

	now := time.Now().UTC()
	initConn := db.SSHConnection{
		Alias:             testAlias,
		IPAddress:         testIP,
		Username:          "Administrator",
		EncryptedPassword: "aes:sample-pass",
		OS:                "windows",
		CreatedAt:         now,
	}

	// 1. Initial insert
	stats1, err1 := SyncSSHConnectionsLocally([]db.SSHConnection{initConn})
	if err1 != nil {
		t.Fatalf("initial sync failed: %v", err1)
	}
	if stats1.Inserted != 1 {
		t.Errorf("expected 1 inserted, got %d", stats1.Inserted)
	}

	// 2. Sync again with same data: should be Matched and Unchanged
	stats2, err2 := SyncSSHConnectionsLocally([]db.SSHConnection{initConn})
	if err2 != nil {
		t.Fatalf("second sync failed: %v", err2)
	}
	if stats2.Matched != 1 || stats2.Unchanged != 1 {
		t.Errorf("expected 1 matched and 1 unchanged, got matched=%d unchanged=%d", stats2.Matched, stats2.Unchanged)
	}

	// 3. Sync with updated IP: should be Matched and Updated
	updateConn := initConn
	updateConn.IPAddress = "192.168.99.100"
	stats3, err3 := SyncSSHConnectionsLocally([]db.SSHConnection{updateConn})
	if err3 != nil {
		t.Fatalf("third sync failed: %v", err3)
	}
	if stats3.Matched != 1 || stats3.Updated != 1 {
		t.Errorf("expected 1 matched and 1 updated, got matched=%d updated=%d", stats3.Matched, stats3.Updated)
	}

	// Verify database record has updated IP
	fetched, fetchErr := db.GetSSHConnectionByAlias(ctx, dbConn.SQL(), testAlias)
	if fetchErr != nil || fetched == nil {
		t.Fatalf("failed to fetch updated connection: %v", fetchErr)
	}
	if fetched.IPAddress != "192.168.99.100" {
		t.Errorf("expected updated IP '192.168.99.100', got %s", fetched.IPAddress)
	}
}

func TestRenderDeployConfigSSHHelp(t *testing.T) {
	// Verify help renders without panicking
	RenderDeployConfigSSHHelp()
}
