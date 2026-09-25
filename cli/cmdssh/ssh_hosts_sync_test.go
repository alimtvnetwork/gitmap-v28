// Package cmdssh — ssh_hosts_sync_test.go tests bidirectional sync and dual-table node persistence.
package cmdssh

import (
	"context"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/cli/db"
	"github.com/alimtvnetwork/gitmap-v28/cli/store"
)

func TestSyncSSHHostsFromConnections(t *testing.T) {
	dbConn, err := openSSHDBFunc()
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer dbConn.Close()
	ctx := context.Background()

	// Clear existing hosts
	_, _ = store.DeleteAllSSHHosts(ctx, dbConn.SQL())
	_ = db.DeleteAllSSHConnections(ctx, dbConn.SQL())

	// Insert into SSHConnection table only
	conn := db.SSHConnection{
		Alias:             "node-alpha",
		IPAddress:         "10.200.1.10",
		Username:          "root",
		EncryptedPassword: "enc",
		CreatedAt:         time.Now().UTC(),
	}
	appErr := db.InsertOrUpdateSSHConnection(ctx, dbConn.SQL(), conn)
	if appErr != nil {
		t.Fatalf("failed to insert SSHConnection: %v", appErr)
	}

	// Verify ssh_hosts is initially empty
	hostsBefore, err := store.ListHosts(ctx, dbConn.SQL())
	if err != nil {
		t.Fatalf("ListHosts failed: %v", err)
	}
	if len(hostsBefore) != 0 {
		t.Fatalf("expected 0 hosts in ssh_hosts before sync, got %d", len(hostsBefore))
	}

	// Run sync
	syncSSHHostsFromConnections(ctx, dbConn.SQL())

	// Verify ssh_hosts now contains the synced node
	hostsAfter, err := store.ListHosts(ctx, dbConn.SQL())
	if err != nil {
		t.Fatalf("ListHosts after sync failed: %v", err)
	}
	if len(hostsAfter) != 1 {
		t.Fatalf("expected 1 host in ssh_hosts after sync, got %d", len(hostsAfter))
	}
	if hostsAfter[0].Alias != "node-alpha" || hostsAfter[0].IP != "10.200.1.10" {
		t.Fatalf("unexpected host data: %+v", hostsAfter[0])
	}
}

func TestImportConnectionsLocally_DualTablePersistence(t *testing.T) {
	conns := []db.SSHConnection{
		{
			Alias:     "node-dual-1",
			IPAddress: "10.200.1.20",
			Username:  "ubuntu",
		},
		{
			Alias:     "node-dual-2",
			IPAddress: "10.200.1.21",
			Username:  "admin",
		},
	}
	imported := importConnectionsLocally(conns)
	if imported != 2 {
		t.Fatalf("expected 2 imported, got %d", imported)
	}

	dbConn, err := openSSHDBFunc()
	if err != nil {
		t.Fatalf("failed to open test db: %v", err)
	}
	defer dbConn.Close()

	ctx := context.Background()
	hosts, err := store.ListHosts(ctx, dbConn.SQL())
	if err != nil {
		t.Fatalf("ListHosts failed: %v", err)
	}

	hasDual1 := false
	hasDual2 := false
	for _, h := range hosts {
		if h.Alias == "node-dual-1" && h.IP == "10.200.1.20" {
			hasDual1 = true
		}
		if h.Alias == "node-dual-2" && h.IP == "10.200.1.21" {
			hasDual2 = true
		}
	}
	if !hasDual1 || !hasDual2 {
		t.Fatalf("expected dual nodes in ssh_hosts table, found: %+v", hosts)
	}
}

func TestCleanBase64Payload(t *testing.T) {
	input := "  \"eyJzY2hlbWFfdmVyc2lvbiI6IjEuMCIs\r\nImV4cG9ydGVkX2F0IjoiMjAy\"  "
	expected := "eyJzY2hlbWFfdmVyc2lvbiI6IjEuMCIsImV4cG9ydGVkX2F0IjoiMjAy"
	got := cleanBase64Payload(input)
	if got != expected {
		t.Fatalf("cleanBase64Payload failed: got %q, expected %q", got, expected)
	}
}

func TestBuildCompactSSHNodesExportEnvelope(t *testing.T) {
	env, err := BuildCompactSSHNodesExportEnvelope()
	if err != nil {
		t.Fatalf("BuildCompactSSHNodesExportEnvelope failed: %v", err)
	}
	if len(env.Connections) != 0 {
		t.Fatalf("expected empty Connections in compact envelope, got %d", len(env.Connections))
	}
}
