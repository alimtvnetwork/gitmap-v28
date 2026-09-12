package searcher

import (
	"context"
	"database/sql"
	"testing"

	_ "modernc.org/sqlite"
)

func initSearchDB(db *sql.DB) error {
	schema := `CREATE TABLE RepoFile (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		RelativePath TEXT NOT NULL,
		AbsolutePath TEXT NOT NULL,
		Content TEXT,
		IsBig INTEGER NOT NULL,
		WriteTime INTEGER NOT NULL,
		CreatedAt INTEGER NOT NULL,
		UpdatedAt INTEGER NOT NULL
	);
	CREATE TABLE SearchCache (
		Id INTEGER PRIMARY KEY AUTOINCREMENT,
		Query TEXT NOT NULL UNIQUE,
		Hits INTEGER NOT NULL,
		ResultJson TEXT NOT NULL,
		CreatedAt INTEGER NOT NULL,
		UpdatedAt INTEGER NOT NULL
	);`
	_, err := db.Exec(schema)

	return err
}

func setupSearchTestDB(t *testing.T) *sql.DB {
	db, err := sql.Open("sqlite", ":memory:")
	if err != nil {
		t.Fatalf("open sqlite: %v", err)
	}

	db.SetMaxOpenConns(1)
	if err := initSearchDB(db); err != nil {
		t.Fatalf("init schema: %v", err)
	}

	return db
}

func insertTestFiles(t *testing.T, db *sql.DB) {
	insertQuery := `INSERT INTO RepoFile (RelativePath, AbsolutePath, Content, IsBig, WriteTime, CreatedAt, UpdatedAt)
		VALUES ('main.go', '/app/main.go', 'func main() {\n\tprintln("hello gitmap")\n}', 0, 100, 100, 100),
		       ('binary.bin', '/app/binary.bin', '', 1, 100, 100, 100);`
	if _, err := db.Exec(insertQuery); err != nil {
		t.Fatalf("insert test files: %v", err)
	}
}

func verifyCachedResult(t *testing.T, ctx context.Context, db *sql.DB) {
	cached, err := SearchRepoDB(ctx, db, "gitmap", 10, true)
	if err != nil || len(cached) != 1 {
		t.Fatalf("cached search failed: %v", err)
	}
}

func TestSearchRepoDBExact(t *testing.T) {
	ctx := context.Background()
	db := setupSearchTestDB(t)
	defer db.Close()
	insertTestFiles(t, db)

	results, err := SearchRepoDB(ctx, db, "gitmap", 10, true)
	if err != nil || len(results) != 1 {
		t.Fatalf("search exact failed: %v", err)
	}

	if results[0].MatchedText != "gitmap" {
		t.Errorf("expected matched text 'gitmap', got %q", results[0].MatchedText)
	}

	verifyCachedResult(t, ctx, db)
}

func TestSearchRepoDBRegex(t *testing.T) {
	ctx := context.Background()
	db := setupSearchTestDB(t)
	defer db.Close()
	insertTestFiles(t, db)

	results, err := SearchRepoDBRegex(ctx, db, "println\\(.*\\)", 10, true)
	if err != nil {
		t.Fatalf("search regex failed: %v", err)
	}

	if len(results) != 1 {
		t.Fatalf("expected 1 regex result, got %d", len(results))
	}
}

func TestApplyLimit(t *testing.T) {
	items := []SearchResult{
		{MatchedText: "1"},
		{MatchedText: "2"},
		{MatchedText: "3"},
	}

	limited := applyLimit(items, 2)
	if len(limited) != 2 {
		t.Errorf("expected 2 items, got %d", len(limited))
	}

	all := applyLimit(items, 0)
	if len(all) != 3 {
		t.Errorf("expected 3 items, got %d", len(all))
	}
}
