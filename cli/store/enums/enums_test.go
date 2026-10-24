package enums

import "testing"

func TestStoreTableType_Methods(t *testing.T) {
	tables := AllStoreTables()
	if len(tables) != 11 {
		t.Fatalf("expected 11 store tables, got %d", len(tables))
	}

	for _, tbl := range tables {
		if !tbl.IsValid() {
			t.Errorf("expected %s to be valid", tbl)
		}

		if tbl.String() != tbl.TableName() {
			t.Errorf("expected String() == TableName() for %s", tbl)
		}
	}

	invalid := StoreTableType("InvalidTable")
	if invalid.IsValid() {
		t.Errorf("expected invalid table to return false")
	}
}

func TestSQLitePragmaType_Methods(t *testing.T) {
	pragmas := DefaultPragmas()
	if len(pragmas) != 4 {
		t.Fatalf("expected 4 default pragmas, got %d", len(pragmas))
	}

	for _, p := range pragmas {
		if !p.IsValid() {
			t.Errorf("expected %s to be valid", p)
		}

		if p.String() != p.Query() {
			t.Errorf("expected String() == Query() for %s", p)
		}
	}

	invalid := SQLitePragmaType("PRAGMA invalid;")
	if invalid.IsValid() {
		t.Errorf("expected invalid pragma to return false")
	}
}

func TestSplitDatabaseType_Methods(t *testing.T) {
	dbs := AllSplitDatabaseTypes()
	if len(dbs) != 6 {
		t.Fatalf("expected 6 split database types, got %d", len(dbs))
	}

	for _, db := range dbs {
		if !db.IsValid() {
			t.Errorf("expected %s to be valid", db)
		}

		if db.String() == "" {
			t.Errorf("expected non-empty string for %s", db)
		}
	}

	invalid := SplitDatabaseType("invalid_db")
	if invalid.IsValid() {
		t.Errorf("expected invalid db to return false")
	}
}
