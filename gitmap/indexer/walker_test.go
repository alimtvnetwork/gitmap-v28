package indexer

import (
	"context"
	"database/sql"
	"os"
	"path/filepath"
	"testing"

	_ "modernc.org/sqlite"
)

func createSchema(db *sql.DB) error {
	schema := `CREATE TABLE RepoFile (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		RelativePath TEXT NOT NULL UNIQUE,
		AbsolutePath TEXT NOT NULL,
		Content TEXT,
		IsBig INTEGER NOT NULL,
		WriteTime INTEGER NOT NULL,
		CreatedAt INTEGER NOT NULL,
		UpdatedAt INTEGER NOT NULL
	);`
	_, err := db.Exec(schema)

	return err
}

func setupTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open in-memory db: %v", err)
	}

	db.SetMaxOpenConns(1)
	if err := createSchema(db); err != nil {
		t.Fatalf("create table: %v", err)
	}

	return db
}

func TestBinarySniffer(t *testing.T) {
	textData := []byte("Hello, world! This is plain text.")
	if isBinaryContent(textData) {
		t.Errorf("expected textData to not be binary")
	}

	binData := []byte("Hello\x00World")
	if !isBinaryContent(binData) {
		t.Errorf("expected binData to be binary")
	}
}

func TestExcludedDirs(t *testing.T) {
	targets := []string{
		".git", "node_modules", ".venv", "dist", "build", "bin",
		"vendor", ".gemini", "coverage", "tmp", "__pycache__", ".turbo",
	}

	for _, dir := range targets {
		if !isExcludedDir(dir) {
			t.Errorf("expected directory %s to be excluded", dir)
		}
	}

	if isExcludedDir("src") {
		t.Errorf("src directory should not be excluded")
	}
}

func writeTestFiles(t *testing.T, tmpDir string) {
	if err := os.WriteFile(filepath.Join(tmpDir, "text.txt"), []byte("normal text file"), 0644); err != nil {
		t.Fatalf("write text: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tmpDir, "binary.bin"), []byte("null\x00byte\x00binary"), 0644); err != nil {
		t.Fatalf("write bin: %v", err)
	}

	venvDir := filepath.Join(tmpDir, ".venv")
	if err := os.MkdirAll(venvDir, 0755); err != nil {
		t.Fatalf("mkdir venv: %v", err)
	}

	if err := os.WriteFile(filepath.Join(venvDir, "skipped.py"), []byte("skip me"), 0644); err != nil {
		t.Fatalf("write skipped: %v", err)
	}
}

func createTestTree(t *testing.T) string {
	tmpDir, err := os.MkdirTemp("", "walker-test-*")
	if err != nil {
		t.Fatalf("mkdirtemp: %v", err)
	}

	t.Cleanup(func() {
		_ = os.RemoveAll(tmpDir)
	})
	writeTestFiles(t, tmpDir)

	return tmpDir
}

func verifyBinaryFile(t *testing.T, db *sql.DB) {
	var isBig int
	var content string
	query := "SELECT IsBig, Content FROM RepoFile WHERE RelativePath = 'binary.bin'"
	if err := db.QueryRow(query).Scan(&isBig, &content); err != nil {
		t.Fatalf("query binary file: %v", err)
	}

	if isBig != 1 {
		t.Errorf("expected binary.bin to have IsBig=1, got %d", isBig)
	}

	if content != "" {
		t.Errorf("expected binary.bin content to be empty, got %q", content)
	}
}

func verifyIndexedFiles(t *testing.T, db *sql.DB) {
	var count int
	if err := db.QueryRow("SELECT COUNT(*) FROM RepoFile").Scan(&count); err != nil {
		t.Fatalf("count files: %v", err)
	}

	if count != 2 {
		t.Fatalf("expected 2 indexed files (text.txt, binary.bin), got %d", count)
	}

	verifyBinaryFile(t, db)
}

func TestWalkerUpfrontBatchAndSniffer(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()
	tmpDir := createTestTree(t)

	w := NewWalker(tmpDir, db, false)
	if err := w.Walk(ctx, 2); err != nil {
		t.Fatalf("walk failed: %v", err)
	}

	verifyIndexedFiles(t, db)
}

func TestWalkerDeltaSkip(t *testing.T) {
	ctx := context.Background()
	db := setupTestDB(t)
	defer db.Close()
	tmpDir := createTestTree(t)

	w := NewWalker(tmpDir, db, false)
	if err := w.Walk(ctx, 2); err != nil {
		t.Fatalf("first walk failed: %v", err)
	}

	if err := w.Walk(ctx, 2); err != nil {
		t.Fatalf("second walk failed: %v", err)
	}

	verifyIndexedFiles(t, db)
}
