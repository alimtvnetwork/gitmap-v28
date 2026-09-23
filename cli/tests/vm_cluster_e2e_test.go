//go:build e2e

package tests_test

import (
	"context"
	"encoding/json"
	"net"
	"os"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

type VMClusterCredentials struct {
	Windows VMUserPass `json:"windows"`
	Ubuntu  VMUserPass `json:"ubuntu"`
}

type VMUserPass struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

type VMNodeInfo struct {
	Alias string
	IP    string
	OS    string
}

func loadVMCredentials(t *testing.T) *VMClusterCredentials {
	t.Helper()
	rootPaths := []string{
		"../../vmpass.json",
		"../vmpass.json",
		"vmpass.json",
	}
	for _, p := range rootPaths {
		if creds := parseCredentialsFile(p); creds != nil {
			return creds
		}
	}
	t.Skip("vmpass.json not found; skipping local-only VM cluster e2e test")
	return nil
}

func parseCredentialsFile(path string) *VMClusterCredentials {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}
	var creds VMClusterCredentials
	if jsonErr := json.Unmarshal(data, &creds); jsonErr != nil {
		return nil
	}
	return &creds
}

func isTCPPortOpen(ip string, port int, timeout time.Duration) bool {
	addr := net.JoinHostPort(ip, "22")
	conn, err := net.DialTimeout("tcp", addr, timeout)
	if err != nil {
		return false
	}
	_ = conn.Close()
	return true
}

func resolveNodeCredentials(creds *VMClusterCredentials, osType string) (string, string) {
	if osType == "windows" {
		return creds.Windows.User, creds.Windows.Pass
	}
	return creds.Ubuntu.User, creds.Ubuntu.Pass
}

func testSingleNodeSSH(t *testing.T, node VMNodeInfo, user, pass string) {
	t.Run(node.Alias, func(t *testing.T) {
		if !isTCPPortOpen(node.IP, 22, 1*time.Second) {
			t.Logf("Node %s (%s) is offline or unreachable; skipping", node.Alias, node.IP)
			return
		}
		client, err := crypto.ConnectWithPassword(node.IP, user, pass)
		if err != nil {
			t.Fatalf("SSH connect failed for node %s: %v", node.Alias, err)
		}
		defer client.Close()

		shell := ""
		if node.OS == "windows" {
			shell = "ps"
		}
		out, runErr := crypto.RunCommand(client, "hostname", shell)
		if runErr != nil {
			t.Fatalf("SSH command execution failed on node %s: %v", node.Alias, runErr)
		}
		if len(out) == 0 {
			t.Fatalf("expected non-empty output from node %s", node.Alias)
		}
		t.Logf("Node %s (%s) response: %s", node.Alias, node.IP, out)
	})
}

func clusterTestNodes() []VMNodeInfo {
	return []VMNodeInfo{
		{Alias: "w1", IP: "192.168.1.3", OS: "windows"},
		{Alias: "w2", IP: "192.168.1.7", OS: "windows"},
		{Alias: "w3", IP: "192.168.1.12", OS: "windows"},
		{Alias: "w4", IP: "192.168.1.13", OS: "windows"},
		{Alias: "u1", IP: "192.168.1.22", OS: "linux"},
	}
}

func TestVMClusterE2EConnectivity(t *testing.T) {
	creds := loadVMCredentials(t)
	nodes := clusterTestNodes()
	for _, node := range nodes {
		user, pass := resolveNodeCredentials(creds, node.OS)
		testSingleNodeSSH(t, node, user, pass)
	}
}

func TestVMClusterE2EGitMapVersions(t *testing.T) {
	creds := loadVMCredentials(t)
	nodes := clusterTestNodes()
	for _, node := range nodes {
		user, pass := resolveNodeCredentials(creds, node.OS)
		testSingleNodeVersion(t, node, user, pass)
	}
}

func testSingleNodeVersion(t *testing.T, node VMNodeInfo, user, pass string) {
	t.Run(node.Alias+"_version", func(t *testing.T) {
		if !isTCPPortOpen(node.IP, 22, 1*time.Second) {
			t.Logf("Node %s is offline; skipping version check", node.Alias)
			return
		}
		client, err := crypto.ConnectWithPassword(node.IP, user, pass)
		if err != nil {
			t.Fatalf("SSH connect failed: %v", err)
		}
		defer client.Close()
		shell := ""
		if node.OS == "windows" {
			shell = "ps"
		}
		cmd := "gitmap version"
		if node.OS == "linux" {
			cmd = "export PATH=\"$HOME/.local/bin:$HOME/.local/bin/gitmap-cli:$PATH\"; gitmap version"
		}
		out, runErr := crypto.RunCommand(client, cmd, shell)
		if runErr != nil {
			t.Fatalf("gitmap version failed on %s: %v", node.Alias, runErr)
		}
		if !strings.Contains(out, "gitmap") {
			t.Fatalf("unexpected version output on %s: %s", node.Alias, out)
		}
		t.Logf("Node %s running GitMap version: %s", node.Alias, strings.TrimSpace(out))
	})
}

func TestVMClusterE2ENodeLifecycle(t *testing.T) {
	dbConn, err := store.OpenDefault()
	if err != nil {
		t.Skipf("cannot open store: %v", err)
		return
	}
	defer dbConn.Close()

	ctx := context.Background()
	testConn := db.SSHConnection{
		Alias:             "e2e-test-node",
		IPAddress:         "192.168.1.250",
		Username:          "e2e-user",
		EncryptedPassword: "e2e-enc-pass",
		OS:                "linux",
	}

	delErr := db.DeleteSSHConnection(ctx, dbConn.SQL(), testConn.Alias)
	if delErr != nil {
		t.Logf("pre-clean notice: %v", delErr)
	}

	insErr := db.InsertOrUpdateSSHConnection(ctx, dbConn.SQL(), testConn)
	if insErr != nil {
		t.Fatalf("failed to insert test connection: %v", insErr)
	}

	connsRes := db.GetSSHConnections(ctx, dbConn.SQL())
	if connsRes.IsFailure() {
		t.Fatalf("failed to get connections: %v", connsRes.AppError())
	}
	found := isNodeAliasPresent(connsRes.Data, testConn.Alias)
	if !found {
		t.Fatalf("expected node %s to be present in database", testConn.Alias)
	}

	cleanErr := db.DeleteSSHConnection(ctx, dbConn.SQL(), testConn.Alias)
	if cleanErr != nil {
		t.Fatalf("failed to cleanup test node: %v", cleanErr)
	}
}

func isNodeAliasPresent(conns []db.SSHConnection, targetAlias string) bool {
	for _, c := range conns {
		if c.Alias == targetAlias {
			return true
		}
	}
	return false
}
