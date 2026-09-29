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
		OSGroup:           "unix",
		OSVersion:         "Ubuntu 22.04.4 LTS",
		BuildVersion:      "5.15.0-generic",
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

	isMatchGroup := fetched.OSGroup == "unix"
	if isMatchGroup == false {
		t.Fatalf("expected OSGroup 'unix', got: %s", fetched.OSGroup)
	}

	isMatchBuild := fetched.BuildVersion == "5.15.0-generic"
	if isMatchBuild == false {
		t.Fatalf("expected BuildVersion '5.15.0-generic', got: %s", fetched.BuildVersion)
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

func TestSSHConnection_MatchAndUpdate(t *testing.T) {
	conn := setupTestSSHDB(t)
	defer conn.Close()
	ctx := context.Background()

	now := time.Now().UTC().Truncate(time.Second)
	initial := SSHConnection{
		Alias:             "w1",
		IPAddress:         "192.168.1.3",
		Username:          "Administrator",
		EncryptedPassword: "rsa:secret-password",
		OS:                "windows",
		CreatedAt:         now,
	}
	if err := InsertOrUpdateSSHConnection(ctx, conn, initial); err != nil {
		t.Fatalf("initial insert failed: %v", err)
	}

	// 1. Update by Alias with empty password: should update OS and preserve password
	updateAlias := SSHConnection{
		Alias:     "w1",
		IPAddress: "192.168.1.3",
		Username:  "Administrator",
		OS:        "linux",
	}
	if err := InsertOrUpdateSSHConnection(ctx, conn, updateAlias); err != nil {
		t.Fatalf("update by alias failed: %v", err)
	}
	fetched1, _ := GetSSHConnectionByAlias(ctx, conn, "w1")
	if fetched1.EncryptedPassword != "rsa:secret-password" {
		t.Fatalf("expected preserved password, got: %s", fetched1.EncryptedPassword)
	}
	if fetched1.OS != "linux" {
		t.Fatalf("expected updated OS 'linux', got: %s", fetched1.OS)
	}

	// 2. Match by IP address with new Alias: should update Alias and preserve password
	updateIP := SSHConnection{
		Alias:     "w1-renamed",
		IPAddress: "192.168.1.3",
		Username:  "Administrator",
	}
	if err := InsertOrUpdateSSHConnection(ctx, conn, updateIP); err != nil {
		t.Fatalf("update by IP failed: %v", err)
	}
	fetched2, _ := GetSSHConnectionByIP(ctx, conn, "192.168.1.3")
	if fetched2.Alias != "w1-renamed" {
		t.Fatalf("expected updated alias 'w1-renamed', got: %s", fetched2.Alias)
	}
	if fetched2.EncryptedPassword != "rsa:secret-password" {
		t.Fatalf("expected preserved password after IP match, got: %s", fetched2.EncryptedPassword)
	}
}
