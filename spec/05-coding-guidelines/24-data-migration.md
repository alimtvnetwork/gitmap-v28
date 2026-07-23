# Data Migration

Universal guidelines for schema versioning, rollback strategies, zero-downtime migrations, and data validation across all projects.

---

## 1. Schema Versioning

### Sequential Migration Files

Every schema change is a numbered, immutable migration file:

```
migrations/
├── 001_create_repos.sql
├── 002_add_groups.sql
├── 003_add_releases.sql
├── 004_add_command_history.sql
├── 005_add_detected_projects.sql
└── 006_convert_uuid_to_int.sql
```

### Rules

| Rule | Detail |
|---|---|
| Sequential numbering | `001_`, `002_`, … — never reorder or reuse |
| One concern per file | Each migration does one logical change |
| Immutable once applied | Never edit a migration after it runs in any environment |
| Forward-only by default | New migration to fix a previous one — don't modify the original |
| Descriptive names | `003_add_releases.sql` not `003_update.sql` |

### Version Tracking

Maintain a `schema_versions` table to track applied migrations:

```sql
CREATE TABLE IF NOT EXISTS schema_versions (
    version   INTEGER PRIMARY KEY,
    name      TEXT    NOT NULL,
    applied   TEXT    NOT NULL DEFAULT (datetime('now'))
);
```

```go
func MigrateUp(db *sql.DB, migrations []Migration) error {
    current := getCurrentVersion(db)

    for _, m := range migrations {
        if m.Version <= current {
            continue
        }
        if err := executeMigration(db, m); err != nil {
            return fmt.Errorf("migration %d (%s) failed: %w", m.Version, m.Name, err)
        }
        recordVersion(db, m.Version, m.Name)
    }
    return nil
}
```

### TypeScript

```typescript
interface Migration {
  version: number;
  name: string;
  up: string;   // SQL to apply
  down: string;  // SQL to rollback
}

async function migrateUp(db: Database, migrations: Migration[]): Promise<void> {
  const current = await getCurrentVersion(db);

  for (const m of migrations) {
    if (m.version <= current) continue;
    await db.exec(m.up);
    await recordVersion(db, m.version, m.name);
  }
}
```

---

## 2. Rollback Strategies

### Rollback Types

| Type | When | Method |
|---|---|---|
| Schema rollback | Migration broke the schema | Run `down` migration |
| Data rollback | Bad data transformation | Restore from backup or reverse transform |
| Application rollback | New code incompatible with old schema | Deploy previous app version |
| Full rollback | Catastrophic failure | Restore database snapshot + previous app |

### Down Migrations

Every `up` migration should have a corresponding `down`:

```sql
-- 006_add_status_column.up.sql
ALTER TABLE repos ADD COLUMN status TEXT NOT NULL DEFAULT 'active';

-- 006_add_status_column.down.sql
ALTER TABLE repos DROP COLUMN status;
```

### Rollback Rules

| Rule | Detail |
|---|---|
| Every migration has a `down` | Even if it's "not expected to rollback" |
| Test rollbacks | CI runs `up` then `down` then `up` again |
| Backup before destructive changes | `DROP`, `ALTER`, data transforms |
| Time-box rollback decisions | If not resolved in 30 min, rollback |
| Rollback is not failure | It's a valid recovery strategy |

### When Rollback Is Not Possible

Some migrations are inherently irreversible:

| Migration | Why | Mitigation |
|---|---|---|
| Drop column | Data is lost | Backup column data to a side table first |
| Change column type | Precision loss | Keep old column, add new, migrate gradually |
| Merge tables | Original boundaries lost | Maintain mapping table during transition |

For irreversible migrations, document the risk and require explicit confirmation:

```go
if !flags.Confirm {
    fmt.Println("WARNING: This migration is irreversible. Use --confirm to proceed.")
    return nil
}
```

---

## 3. Zero-Downtime Migrations

### Expand-Contract Pattern

Never modify or remove a column/table in a single step. Use a two-phase approach:

```
Phase 1 (Expand): Add new structure alongside old
Phase 2 (Contract): Remove old structure after all code uses new
```

### Example: Rename Column

**Wrong** (causes downtime):
```sql
ALTER TABLE repos RENAME COLUMN url TO clone_url;  -- breaks old code immediately
```

**Right** (zero-downtime):
```sql
-- Migration 1: Expand
ALTER TABLE repos ADD COLUMN clone_url TEXT;
UPDATE repos SET clone_url = url;

-- Deploy code that reads from clone_url, writes to both

-- Migration 2: Contract (after all code migrated)
ALTER TABLE repos DROP COLUMN url;
```

### Example: Change Column Type (UUID → Integer)

```sql
-- Migration 1: Add new column
ALTER TABLE repos ADD COLUMN new_id INTEGER;

-- Migration 2: Backfill
UPDATE repos SET new_id = rowid;

-- Migration 3: Swap (after code updated)
-- Drop dependent tables, rebuild with new FK type
-- Re-insert data with integer references
```

### Rules

| Rule | Detail |
|---|---|
| Additive changes only in Phase 1 | Add columns, add tables, add indexes |
| Dual-write during transition | Write to both old and new columns |
| Read from new, fallback to old | Gradually shift reads |
| Remove old only after full migration | All code, all environments |
| No long-running transactions | Lock contention kills availability |

### SQLite Considerations

SQLite has limited `ALTER TABLE` support:

| Supported | Not Supported |
|---|---|
| `ADD COLUMN` | `DROP COLUMN` (before 3.35) |
| `RENAME TABLE` | `ALTER COLUMN` type |
| `RENAME COLUMN` (3.25+) | `ADD CONSTRAINT` |

For unsupported operations, use the rebuild pattern:

```sql
-- 1. Create new table with desired schema
CREATE TABLE repos_new (
    id    INTEGER PRIMARY KEY AUTOINCREMENT,
    name  TEXT NOT NULL,
    url   TEXT NOT NULL
);

-- 2. Copy data
INSERT INTO repos_new (name, url) SELECT name, url FROM repos;

-- 3. Drop old table
DROP TABLE repos;

-- 4. Rename new table
ALTER TABLE repos_new RENAME TO repos;
```

---

## 4. Data Validation

### Pre-Migration Validation

Before running a migration, verify preconditions:

```go
func validatePreMigration(db *sql.DB) error {
    // Check row counts
    count := getRowCount(db, "repos")
    if count == 0 {
        log.Println("WARN: repos table is empty — migration is a no-op")
    }

    // Check for orphaned references
    orphans := countOrphans(db, "group_repos", "repo_id", "repos", "id")
    if orphans > 0 {
        return fmt.Errorf("found %d orphaned group_repos records — clean up first", orphans)
    }

    return nil
}
```

### Post-Migration Validation

After every migration, verify the result:

```go
func validatePostMigration(db *sql.DB, expected int) error {
    // Verify row count preserved
    actual := getRowCount(db, "repos")
    if actual != expected {
        return fmt.Errorf("row count mismatch: expected %d, got %d", expected, actual)
    }

    // Verify no NULL values in required columns
    nulls := countNulls(db, "repos", "name")
    if nulls > 0 {
        return fmt.Errorf("found %d NULL name values after migration", nulls)
    }

    // Verify foreign key integrity
    if err := checkForeignKeys(db); err != nil {
        return fmt.Errorf("foreign key violation: %w", err)
    }

    return nil
}
```

### Validation Checklist

| Check | When | Example |
|---|---|---|
| Row count preservation | After data copy/transform | Before count == after count |
| NOT NULL constraints | After adding required columns | No unexpected NULLs |
| Foreign key integrity | After any FK-related change | All references resolve |
| Unique constraint | After merges or deduplication | No duplicate keys |
| Type correctness | After type conversion | All values parse correctly |
| Index existence | After rebuild operations | Required indexes present |

### Graceful Recovery

When validation fails, provide clear user-facing instructions:

```go
const MsgMigrationFailed = `
Migration failed: %s

To recover:
  1. Run 'gitmap-v27 db-reset --confirm' to rebuild the database
  2. Run 'gitmap-v27 rescan' to re-import repository data

Your Git repositories are not affected — only cached metadata is reset.
`
```

```go
if err := validatePostMigration(db, expectedCount); err != nil {
    log.Printf(MsgMigrationFailed, err)
    return err
}
```

---

## 5. Migration Testing

### CI Pipeline

```
Run migrations up → Validate schema → Run migrations down → Run migrations up again → Run application tests
```

### Test Requirements

| Test | Detail |
|---|---|
| Fresh install | All migrations run on an empty database |
| Upgrade path | Migrations run on a database with existing data |
| Rollback path | Down migrations execute without errors |
| Idempotency | Running the same migration twice doesn't break |
| Data integrity | Row counts and constraints hold after migration |

### Go Test Example

```go
func TestMigrations_UpDown(t *testing.T) {
    db := openTestDB(t)

    // Run all migrations up
    err := MigrateUp(db, allMigrations)
    require.NoError(t, err)

    // Verify schema
    assertTableExists(t, db, "repos")
    assertTableExists(t, db, "schema_versions")

    // Run all migrations down
    err = MigrateDown(db, allMigrations)
    require.NoError(t, err)

    // Run up again — must succeed
    err = MigrateUp(db, allMigrations)
    require.NoError(t, err)
}
```

---

## Constraints

| Constraint | Detail |
|---|---|
| No manual schema changes | All changes via numbered migrations |
| No edited migrations | Immutable once applied — create a new one |
| No downtime for additive changes | `ADD COLUMN` and `CREATE TABLE` are always safe |
| No `DROP` without backup | Always preserve data before destructive operations |
| No long transactions | Batch large data moves to avoid locking |
| Require `--confirm` for destructive ops | User must explicitly opt in |
| Post-migration validation is mandatory | Every migration verifies its result |

---

## References

- [Database Patterns](./11-database-patterns.md)
- [Error Handling Patterns](./04-error-handling.md)
- [Testing Patterns](./06-testing-patterns.md)
- [Database Spec](../01-app/16-database.md)
- [Legacy ID Migration](../02-app-issues/12-legacy-id-migration.md)

---

## Contributors

- [**Md. Alim Ul Karim**](https://www.linkedin.com/in/alimkarim) — Creator & Lead Architect. System architect with 20+ years of professional software engineering experience across enterprise, fintech, and distributed systems. Recognized as one of the top software architects globally. Alim's architectural philosophy — consistency over cleverness, convention over configuration — is the driving force behind every design decision in this framework.
  - [Google Profile](https://www.google.com/search?q=Alim+Ul+Karim)
- [Riseup Asia LLC (Top Leading Software Company in WY)](https://riseup-asia.com) (2026)
  - [Facebook](https://www.facebook.com/riseupasia.talent/)
  - [LinkedIn](https://www.linkedin.com/company/105304484/)
  - [YouTube](https://www.youtube.com/@riseup-asia)
