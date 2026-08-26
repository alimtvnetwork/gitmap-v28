# Database Testing Strategy

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

All database schemas, migrations, and queries MUST be tested. Testing follows a **two-tier approach**: unit tests for schema logic and integration tests with in-memory databases for data operations.

---

## 1. Testing Tiers

| Tier | What It Tests | Database | Speed |
|------|--------------|----------|-------|
| **Unit tests** | Schema creation, migrations, constraints, column types | In-memory SQLite (`:memory:`) | Fast (ms) |
| **Integration tests** | Full CRUD operations, views, relationships, ORM queries | In-memory SQLite (`:memory:`) | Fast (ms) |

> **Key insight:** SQLite's in-memory mode (`:memory:`) means both tiers run without touching disk — no test database setup, no cleanup, no Docker containers.

---

## 2. Unit Tests — Schema Validation

### 2.1 What to Test

| Test | What It Validates |
|------|-------------------|
| Table creation | Migration creates expected tables |
| Column existence | All columns present with correct types |
| Constraints | NOT NULL, UNIQUE, DEFAULT values work |
| Foreign keys | FK relationships are enforced |
| Indexes | Expected indexes exist |
| Views | Views return expected columns |

### 2.2 Go Example

```go
func TestTransactionSchema(t *testing.T) {
    db, err := sql.Open("sqlite3", ":memory:")
    require.NoError(t, err)
    defer db.Close()

    // Run migrations
    err = RunMigrations(db)
    require.NoError(t, err)

    // Verify table exists and columns are correct
    rows, err := db.Query("PRAGMA table_info(Transactions)")
    require.NoError(t, err)
    defer rows.Close()

    columns := make(map[string]string)
    for rows.Next() {
        var cid int
        var name, colType string
        var notNull, pk int
        var dfltValue sql.NullString
        rows.Scan(&cid, &name, &colType, &notNull, &dfltValue, &pk)
        columns[name] = colType
    }

    assert.Equal(t, "INTEGER", columns["TransactionId"])
    assert.Equal(t, "TEXT", columns["PluginSlug"])
    assert.Equal(t, "INTEGER", columns["StatusTypeId"])
    assert.Contains(t, columns, "CreatedAt")
}

func TestForeignKeyEnforcement(t *testing.T) {
    db, _ := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
    defer db.Close()

    RunMigrations(db)

    // Insert with invalid FK should fail
    _, err := db.Exec(
        "INSERT INTO Transactions (TransactionId, StatusTypeId) VALUES (1, 999)",
    )

    assert.Error(t, err, "FK constraint should reject invalid StatusTypeId")
}
```

### 2.3 PHP Example

```php
public function testTransactionTableHasCorrectColumns(): void
{
    $pdo = new PDO('sqlite::memory:');
    $pdo->setAttribute(PDO::ATTR_ERRMODE, PDO::ERRMODE_EXCEPTION);

    // Run migrations
    (new MigrationRunner($pdo))->runAll();

    $stmt = $pdo->query("PRAGMA table_info(Transactions)");
    $columns = array_column($stmt->fetchAll(PDO::FETCH_ASSOC), 'type', 'name');

    $this->assertEquals('INTEGER', $columns['TransactionId']);
    $this->assertEquals('TEXT', $columns['PluginSlug']);
    $this->assertEquals('INTEGER', $columns['StatusTypeId']);
}
```

---

## 3. Integration Tests — Data Operations

### 3.1 What to Test

| Test | What It Validates |
|------|-------------------|
| CRUD operations | Insert, select, update, delete through ORM |
| View queries | Views return correctly joined data |
| Cascade behavior | ON DELETE CASCADE works as expected |
| Lookup table data | Seed data is correctly populated |
| Edge cases | NULL handling, empty results, boundary values |
| Migration idempotency | Running migrations twice doesn't break |

### 3.2 Go Example — ORM Integration Test

```go
func TestTransactionCRUD(t *testing.T) {
    db, _ := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
    defer db.Close()

    RunMigrations(db)
    SeedLookupTables(db) // Insert StatusTypes, FileTypes, etc.
    repo := NewTransactionRepo(db)

    // Create
    tx := Transaction{
        PluginSlug:   "my-plugin",
        StatusTypeId: 1, // Pending
        FileTypeId:   1, // Plugin
        Amount:       29.99,
    }
    id, err := repo.Insert(tx)
    require.NoError(t, err)
    assert.Greater(t, id, int64(0))

    // Read
    found, err := repo.FindById(id)
    require.NoError(t, err)
    assert.Equal(t, "my-plugin", found.PluginSlug)

    // Update
    err = repo.UpdateStatus(id, 2) // Complete
    require.NoError(t, err)

    // Verify via view
    details, err := repo.FindDetailById(id) // Queries VwTransactionDetails
    require.NoError(t, err)
    assert.Equal(t, "Complete", details.StatusName)

    // Delete
    err = repo.Delete(id)
    require.NoError(t, err)

    _, err = repo.FindById(id)
    assert.Error(t, err)
}
```

### 3.3 Test Helper — In-Memory Database Setup

```go
// testutil/db.go — reusable test helper
func NewTestDB(t *testing.T) *sql.DB {
    t.Helper()
    db, err := sql.Open("sqlite3", ":memory:?_foreign_keys=on")
    require.NoError(t, err)

    t.Cleanup(func() { db.Close() })

    err = RunMigrations(db)
    require.NoError(t, err)

    err = SeedLookupTables(db)
    require.NoError(t, err)

    return db
}

// Usage in any test
func TestSomething(t *testing.T) {
    db := testutil.NewTestDB(t)
    repo := NewMyRepo(db)
    // ... test with fully migrated + seeded in-memory DB
}
```

---

## 4. View Testing

### 4.1 Test That Views Return Expected Columns

```go
func TestVwTransactionDetailsView(t *testing.T) {
    db := testutil.NewTestDB(t)

    // Insert test data
    db.Exec("INSERT INTO StatusTypes (StatusTypeId, Name) VALUES (1, 'Pending')")
    db.Exec("INSERT INTO FileTypes (FileTypeId, Name) VALUES (1, 'Plugin')")
    db.Exec(`INSERT INTO Transactions
        (TransactionId, PluginSlug, StatusTypeId, FileTypeId, Amount)
        VALUES (1, 'test-plugin', 1, 1, 9.99)`)

    // Query the view
    row := db.QueryRow("SELECT StatusName, FileTypeName FROM VwTransactionDetails WHERE TransactionId = 1")

    var statusName, fileTypeName string
    err := row.Scan(&statusName, &fileTypeName)

    require.NoError(t, err)
    assert.Equal(t, "Pending", statusName)
    assert.Equal(t, "Plugin", fileTypeName)
}
```

---

## 5. Migration Testing

### 5.1 Test Migrations Are Idempotent

```go
func TestMigrationsIdempotent(t *testing.T) {
    db := testutil.NewTestDB(t) // Already ran migrations once

    // Running again should not error
    err := RunMigrations(db)

    assert.NoError(t, err, "Migrations must be safe to run multiple times")
}
```

### 5.2 Test Migration Versioning

```go
func TestMigrationVersionTracking(t *testing.T) {
    db := testutil.NewTestDB(t)

    version, err := GetCurrentMigrationVersion(db)
    require.NoError(t, err)
    assert.Greater(t, version, 0, "Migration version should be tracked")
}
```

---

## 6. Test File Organization

```
project/
├── internal/
│   ├── db/
│   │   ├── migrations.go
│   │   ├── migrations_test.go      ← Schema unit tests
│   │   └── seed.go
│   ├── repos/
│   │   ├── transaction_repo.go
│   │   └── transaction_repo_test.go ← Integration tests
│   └── testutil/
│       └── db.go                    ← Shared test helpers
└── ...
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Schema design | [./02-schema-design.md](./02-schema-design.md) |
| ORM and views | [./03-orm-and-views.md](./03-orm-and-views.md) |
| Test naming conventions | [../02-coding-guidelines/01-cross-language/14-test-naming-and-structure.md](../02-coding-guidelines/01-cross-language/14-test-naming-and-structure.md) |
