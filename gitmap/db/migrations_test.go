package db

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func TestApplyMigrations(t *testing.T) {
	conn, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("failed to open memory db: %v", err)
	}
	defer conn.Close()

	ctx := context.Background()
	err = ApplyMigrations(ctx, conn)
	if err != nil {
		t.Fatalf("ApplyMigrations failed: %v", err)
	}
}
