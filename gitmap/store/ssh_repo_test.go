package store

import (
	"context"
	"database/sql"
	"testing"
	"time"

	"github.com/alimtvnetwork/gitmap-v28/gitmap/apperror"
)

func setupSSHTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}

	err = RegisterSSHHostMigration(db, 1, false)
	if err != nil {
		t.Fatalf("failed to run migration: %v", err)
	}
	return db
}

func TestInsertSSHHost(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now().Truncate(time.Second).UTC()

	host := SSHHost{
		ID:        "host-1",
		Alias:     "my-server",
		IP:        "192.168.1.100",
		Username:  "admin",
		CreatedAt: now,
	}

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		t.Fatalf("failed to begin tx: %v", err)
	}

	err = InsertSSHHost(ctx, host, tx)
	if err != nil {
		t.Fatalf("InsertSSHHost failed: %v", err)
	}

	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed to commit: %v", err)
	}

	var count int
	err = db.QueryRow("SELECT COUNT(*) FROM ssh_hosts WHERE id = 'host-1'").Scan(&count)
	if err != nil {
		t.Fatalf("failed to query count: %v", err)
	}
	if count != 1 {
		t.Errorf("expected 1 record, got %d", count)
	}
	
	// Test constraint violation
	tx2, _ := db.BeginTx(ctx, nil)
	err = InsertSSHHost(ctx, host, tx2)
	if err == nil {
		t.Fatal("expected error on duplicate insert, got nil")
	}
	tx2.Rollback()
	
	if appErr, ok := err.(*apperror.AppError); ok {
		if appErr.Code != "E_INTERNAL_ERROR" {
			t.Errorf("expected E_INTERNAL_ERROR, got %s", appErr.Code)
		}
	} else {
		t.Errorf("expected AppError, got %T", err)
	}
}

func TestGetHostByAlias(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	now := time.Now().Truncate(time.Second).UTC()

	host := SSHHost{
		ID:        "host-2",
		Alias:     "web-server",
		IP:        "192.168.1.101",
		Username:  "root",
		CreatedAt: now,
	}

	tx, _ := db.BeginTx(ctx, nil)
	_ = InsertSSHHost(ctx, host, tx)
	tx.Commit()

	// Find existing
	found, err := GetHostByAlias(ctx, "web-server", db)
	if err != nil {
		t.Fatalf("GetHostByAlias failed: %v", err)
	}
	if found.ID != "host-2" {
		t.Errorf("expected ID host-2, got %s", found.ID)
	}

	// Find non-existing
	_, err = GetHostByAlias(ctx, "db-server", db)
	if err == nil {
		t.Fatal("expected error for non-existent alias, got nil")
	}
	
	if appErr, ok := err.(*apperror.AppError); ok {
		if appErr.Code != "E_INTERNAL_ERROR" {
			t.Errorf("expected E_INTERNAL_ERROR, got %s", appErr.Code)
		}
		if appErr.Cause != apperror.ErrNotFound {
			t.Errorf("expected cause ErrNotFound, got %v", appErr.Cause)
		}
	} else {
		t.Errorf("expected AppError, got %T", err)
	}
}

func TestDeleteHostByIP(t *testing.T) {
	db := setupSSHTestDB(t)
	defer db.Close()

	ctx := context.Background()
	
	host := SSHHost{
		ID:        "host-3",
		Alias:     "delete-me",
		IP:        "10.0.0.99",
		Username:  "ubuntu",
		CreatedAt: time.Now().UTC(),
	}

	tx, _ := db.BeginTx(ctx, nil)
	_ = InsertSSHHost(ctx, host, tx)
	tx.Commit()

	// Delete existing
	err := DeleteHostByIP(ctx, "10.0.0.99", db)
	if err != nil {
		t.Fatalf("DeleteHostByIP failed: %v", err)
	}

	var count int
	_ = db.QueryRow("SELECT COUNT(*) FROM ssh_hosts WHERE ip = '10.0.0.99'").Scan(&count)
	if count != 0 {
		t.Errorf("expected 0 records after delete, got %d", count)
	}

	// Idempotency: delete non-existing
	err = DeleteHostByIP(ctx, "10.0.0.99", db)
	if err != nil {
		t.Errorf("expected no error deleting non-existent, got %v", err)
	}
}
