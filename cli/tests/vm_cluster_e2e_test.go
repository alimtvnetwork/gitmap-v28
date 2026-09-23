//go:build e2e

package tests_test

import (
	"encoding/json"
	"net"
	"os"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/crypto"
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

func TestVMClusterE2EConnectivity(t *testing.T) {
	creds := loadVMCredentials(t)
	nodes := []VMNodeInfo{
		{Alias: "w1", IP: "192.168.1.3", OS: "windows"},
		{Alias: "w2", IP: "192.168.1.7", OS: "windows"},
		{Alias: "w3", IP: "192.168.1.12", OS: "windows"},
		{Alias: "w4", IP: "192.168.1.13", OS: "windows"},
		{Alias: "u1", IP: "192.168.1.22", OS: "linux"},
	}
	for _, node := range nodes {
		user, pass := resolveNodeCredentials(creds, node.OS)
		testSingleNodeSSH(t, node, user, pass)
	}
}
