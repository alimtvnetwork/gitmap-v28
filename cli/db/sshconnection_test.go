package db

import (
	"context"
	"database/sql"
	"testing"
	"time"

	_ "modernc.org/sqlite"
)

func setupTestSSHDB(t *testing.T) *sql.DB {
	conn, err := sql.Open("sqlite", ":memory:")
	isNil := err == nil
	if isNil == false {
		t.Fatalf("failed to open memory db: %v", err)
	}

	return conn
}

func TestSSHConnection_InsertAndRetrieveWithOSMetadata(t *testing.T) {
	conn := setupTestSSHDB(t)
	defer conn.Close()
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	testConn := SSHConnection{
		Alias:             "node-u2",
		IPAddress:         "192.168.1.5",
		Username:          "alim",
		EncryptedPassword: "rsa:test-encrypted-pass",
		KeyPath:           "/home/user/.ssh/id_rsa",
		OS:                "linux",
		OSVersion:         "Ubuntu 22.04.4 LTS",
		FirstRunAt:        now,
		CreatedAt:         now,
	}

	appErr := InsertOrUpdateSSHConnection(ctx, conn, testConn)
	isSuccess := appErr == nil
	if isSuccess == false {
		t.Fatalf("InsertOrUpdateSSHConnection failed: %v", appErr)
	}

	fetched, err := GetSSHConnectionByAlias(ctx, conn, "node-u2")
	hasErr := err != nil
	if hasErr {
		t.Fatalf("GetSSHConnectionByAlias failed: %v", err)
	}

	isMatchOS := fetched.OSVersion == "Ubuntu 22.04.4 LTS"
	if isMatchOS == false {
		t.Fatalf("expected OSVersion 'Ubuntu 22.04.4 LTS', got: %s", fetched.OSVersion)
	}

	isMatchFirstRun := fetched.FirstRunAt.Equal(now)
	if isMatchFirstRun == false {
		t.Fatalf("expected FirstRunAt %v, got: %v", now, fetched.FirstRunAt)
	}
}

func TestSSHConnection_GetByIP(t *testing.T) {
	conn := setupTestSSHDB(t)
	defer conn.Close()
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	testConn := SSHConnection{
		Alias:     "win-w2",
		IPAddress: "192.168.1.8",
		Username:  "Administrator",
		OS:        "windows",
		OSVersion: "Microsoft Windows 10 Pro",
		CreatedAt: now,
	}
	_ = InsertOrUpdateSSHConnection(ctx, conn, testConn)

	fetched, err := GetSSHConnectionByIP(ctx, conn, "192.168.1.8")
	hasErr := err != nil
	if hasErr {
		t.Fatalf("GetSSHConnectionByIP failed: %v", err)
	}

	isWindows := fetched.OS == "windows"
	if isWindows == false {
		t.Fatalf("expected OS 'windows', got: %s", fetched.OS)
	}
}
