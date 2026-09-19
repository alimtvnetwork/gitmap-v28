package cmdssh

import (
	"bytes"
	"context"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

const sampleKnownHostsLine = "192.168.1.5 ssh-ed25519 AAAAC3NzaC1lZDI1NTE5AAAAIGitMapTestKeyHostKeyEntry test-comment"

func TestParseKnownHostLine(t *testing.T) {
	entry, ok := parseKnownHostLine(sampleKnownHostsLine)
	if !ok {
		t.Fatal("expected parseKnownHostLine to succeed")
	}
	if entry.Host != "192.168.1.5" || entry.KeyType != "ssh-ed25519" {
		t.Fatalf("unexpected entry: %+v", entry)
	}
	if !strings.HasPrefix(entry.Fingerprint, "SHA256:") {
		t.Fatalf("expected SHA256: prefix in fingerprint, got %s", entry.Fingerprint)
	}
	if entry.Comment != "test-comment" {
		t.Fatalf("expected comment 'test-comment', got %s", entry.Comment)
	}
}

func TestParseKnownHostLine_Ignored(t *testing.T) {
	_, ok1 := parseKnownHostLine("# this is a comment")
	_, ok2 := parseKnownHostLine("")
	_, ok3 := parseKnownHostLine("   ")
	if ok1 || ok2 || ok3 {
		t.Fatal("expected comments and blank lines to be ignored")
	}
}

func TestAppendAndRemoveKnownHostsFile(t *testing.T) {
	tmpFile := filepath.Join(t.TempDir(), "known_hosts")
	err := AppendKnownHostFile(tmpFile, "192.168.1.5", "ssh-ed25519", "AAAAC3NzaC1lZDI1NTE5AAAAI...")
	if err != nil {
		t.Fatalf("failed to append: %v", err)
	}
	content, _ := os.ReadFile(tmpFile)
	if !strings.Contains(string(content), "192.168.1.5") {
		t.Fatalf("expected file to contain host, got: %s", string(content))
	}

	removed, err := RemoveFromKnownHostsFile(tmpFile, "192.168.1.5")
	if err != nil || removed != 1 {
		t.Fatalf("failed to remove: %v, count=%d", err, removed)
	}
	contentAfter, _ := os.ReadFile(tmpFile)
	if strings.Contains(string(contentAfter), "192.168.1.5") {
		t.Fatalf("expected host to be removed, got: %s", string(contentAfter))
	}
}

func TestSSHKnownHostsRepo_CRUD(t *testing.T) {
	testDB := setupTestDB(t)
	ctx := context.Background()

	kh := store.SSHKnownHost{
		ID:          "kh-192.168.1.5-ssh-ed25519",
		Host:        "192.168.1.5",
		KeyType:     "ssh-ed25519",
		PublicKey:   "AAAAC3NzaC1lZDI1NTE5AAAAI...",
		Fingerprint: "SHA256:ROdbMLgtCZ47lGGUryHnDCpik9/G5V8H/ie6tuqai2k",
		Comment:     "test",
		CreatedAt:   time.Now().UTC(),
		UpdatedAt:   time.Now().UTC(),
	}

	if err := store.UpsertSSHKnownHost(ctx, kh, testDB.Conn()); err != nil {
		t.Fatalf("upsert failed: %v", err)
	}

	hosts, err := store.ListSSHKnownHosts(ctx, testDB.Conn())
	if err != nil || len(hosts) != 1 {
		t.Fatalf("list failed: %v, count=%d", err, len(hosts))
	}

	got, err := store.GetSSHKnownHost(ctx, "192.168.1.5", testDB.Conn())
	if err != nil || got.Host != "192.168.1.5" {
		t.Fatalf("get failed: %v, got=%+v", err, got)
	}

	rows, err := store.DeleteSSHKnownHostByTarget(ctx, "192.168.1.5", testDB.Conn())
	if err != nil || rows != 1 {
		t.Fatalf("delete failed: %v, rows=%d", err, rows)
	}
}

func TestAppendHostKeyCheckingDefault(t *testing.T) {
	res1 := appendHostKeyCheckingDefault([]string{}, []string{})
	if len(res1) != 2 || res1[0] != "-o" || res1[1] != "StrictHostKeyChecking=accept-new" {
		t.Fatalf("expected StrictHostKeyChecking=accept-new, got: %v", res1)
	}

	customArgs := []string{"-o", "StrictHostKeyChecking=no"}
	res2 := appendHostKeyCheckingDefault([]string{}, customArgs)
	if len(res2) != 0 {
		t.Fatalf("expected no default appended when custom present, got: %v", res2)
	}
}

func TestRenderKnownHostsTable(t *testing.T) {
	var buf bytes.Buffer
	hosts := []store.SSHKnownHost{
		{
			Host:        "192.168.1.5",
			KeyType:     "ssh-ed25519",
			Fingerprint: "SHA256:ROdbMLgtCZ47lGGUryHnDCpik9/G5V8H/ie6tuqai2k",
			UpdatedAt:   time.Date(2026, 9, 20, 1, 21, 0, 0, time.UTC),
		},
	}
	renderKnownHostsTable(&buf, hosts)
	out := buf.String()
	if !strings.Contains(out, "192.168.1.5") || !strings.Contains(out, "ssh-ed25519") {
		t.Fatalf("unexpected table output: %s", out)
	}
}
