# Consolidated: App Database — Complete Reference

**Version:** 3.2.0
**Updated:** 2026-04-16

---

## Purpose

This is the **standalone consolidated reference** for app-specific database design. An AI reading only this file must be able to design, implement, and validate the application's data model — including table schemas, migration strategies, query patterns, and data integrity rules — without consulting source specs.

**Source:** `02-spec/23-app-db/` + `02-spec/04-database-conventions/` (naming rules) + `02-spec/05-split-db-architecture/` (SQLite partitioning)

---

## 1. Relationship to Core Database Specs

This file covers **app-specific** database decisions. For **cross-project** rules, see the companion consolidated files:

| Concern | Consolidated File | What It Covers |
|---------|-------------------|----------------|
| Naming, PK/FK, types, views | [21-database-conventions.md](./21-database-conventions.md) | PascalCase everything, singular tables, integer PKs, column right-sizing, ORM patterns |
| Split DB architecture | [08-split-db-architecture.md](./08-split-db-architecture.md) | Root/App/Session hierarchy, WAL, dynamic creation, RBAC, reset API |
| Seedable config | [09-seedable-config.md](./09-seedable-config.md) | `config.seed.json` merge strategy, schema versioning |

**Rule:** All conventions from `21-database-conventions.md` apply unconditionally to app tables. This file does NOT override core naming rules.

---

## 2. Core Naming Rules (Inherited — Quick Reference)

These are inherited from the core database conventions. Repeated here for standalone completeness:

| Object | Convention | Example |
|--------|-----------|---------|
| Table names | PascalCase, **singular** | `User`, `AgentSite`, `Transaction` |
| Column names | PascalCase | `PluginSlug`, `CreatedAt`, `DisplayName` |
| Primary key | `{TableName}Id` | `UserId`, `TransactionId` |
| Foreign key column | Same name as referenced PK | `AgentSiteId` in both parent and child tables |
| Boolean columns | `Is`/`Has` prefix, **positive only** | `IsActive`, `HasLicense` — never `IsDisabled` |
| Index names | `Idx{Table}_{Column}` | `IdxTransaction_CreatedAt` |
| View names | `Vw` prefix + PascalCase | `VwTransactionDetail` |
| Abbreviations | First letter only capitalized | `Id`, `Url`, `Api` — never `ID`, `URL`, `API` |
| PK type | `INTEGER PRIMARY KEY AUTOINCREMENT` | Default for entity tables |
| UUID | ❌ Forbidden unless all 3 conditions met | Multi-system, no central authority, public non-guessable |

---

## 3. App Database Placement

### 3.1 Which Database Layer?

The split-DB architecture defines three layers. App-specific tables belong in the **App DB** layer:

```
data/
├── root.db                    # Root registry — global settings, app list
├── {project-slug}/
│   ├── app.db                 # ← APP-SPECIFIC TABLES GO HERE
│   ├── config/
│   │   └── settings.db        # Project config
│   ├── cache/
│   │   └── search-cache.db    # Cache with TTL
│   └── chat/
│       └── 001-{id}.db        # Session-scoped data
```

### 3.2 Decision Guide: Which Layer?

| Question | Answer → Layer |
|----------|----------------|
| Is this shared across all apps? | Root DB (`root.db`) |
| Is this specific to one app/project? | **App DB** (`{project}/app.db`) |
| Is this tied to a user session or conversation? | Session DB (`{project}/chat/001-{id}.db`) |
| Is this cached/ephemeral with TTL? | Cache DB (`{project}/cache/`) |
| Is this a document chunk or embedding? | Document DB (`{project}/rag/documents/`) |

---

## 4. App Table Design Template

### 4.1 Entity Table Template

```sql
-- ✅ CORRECT: App entity table following all conventions (incl. Rule 10)
CREATE TABLE {EntityName} (
    {EntityName}Id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- Required columns
    CreatedAt     DATETIME NOT NULL DEFAULT (datetime('now')),
    UpdatedAt     DATETIME NOT NULL DEFAULT (datetime('now')),
    -- Business columns
    {Column1}     {TYPE} NOT NULL,
    {Column2}     {TYPE},
    -- Boolean: positive naming only
    IsActive      BOOLEAN NOT NULL DEFAULT 1,
    -- Foreign key
    {ParentTable}Id INTEGER NOT NULL,
    -- Rule 10: nullable free-text column for future-proofing / hints
    Description   TEXT NULL,
    FOREIGN KEY ({ParentTable}Id) REFERENCES {ParentTable}({ParentTable}Id)
);

-- Index: FK columns and frequently queried columns
CREATE INDEX Idx{EntityName}_{ParentTable}Id ON {EntityName}({ParentTable}Id);
CREATE INDEX Idx{EntityName}_CreatedAt ON {EntityName}(CreatedAt);
```

### 4.2 Lookup/Reference Table Template

```sql
-- ✅ CORRECT: Lookup table — no AUTOINCREMENT, manually seeded (incl. Rule 10)
CREATE TABLE {TypeName} (
    {TypeName}Id  SMALLINT PRIMARY KEY,
    Name          TEXT NOT NULL UNIQUE,
    DisplayOrder  SMALLINT NOT NULL DEFAULT 0,
    -- Rule 10: nullable hint text — "What this row means"
    Description   TEXT NULL
);

-- Seed values
INSERT INTO {TypeName} ({TypeName}Id, Name, DisplayOrder, Description) VALUES
    (1, 'Active',   1, 'Currently in use'),
    (2, 'Inactive', 2, 'Temporarily disabled'),
    (3, 'Archived', 3, 'No longer used; retained for history');
```

### 4.3 Junction/Bridge Table Template

```sql
-- ✅ CORRECT: Many-to-many junction table — exempt from Rules 10/11/12
-- linter-waive: MISSING-DESC-001 reason="Pure junction table; relational glue only"
CREATE TABLE {TableA}{TableB} (
    {TableA}{TableB}Id INTEGER PRIMARY KEY AUTOINCREMENT,
    {TableA}Id         INTEGER NOT NULL,
    {TableB}Id         INTEGER NOT NULL,
    CreatedAt          DATETIME NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY ({TableA}Id) REFERENCES {TableA}({TableA}Id),
    FOREIGN KEY ({TableB}Id) REFERENCES {TableB}({TableB}Id),
    UNIQUE ({TableA}Id, {TableB}Id)
);
```

### 4.4 Transactional Table Template

```sql
-- ✅ CORRECT: Transactional table following Rule 11 (Notes + Comments)
CREATE TABLE {TransactionName} (
    {TransactionName}Id INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId              INTEGER NOT NULL,
    Amount              DECIMAL(10,2) NOT NULL,
    {TransactionName}StatusId INTEGER NOT NULL,
    -- Rule 11: nullable free-text columns for operational + human context
    Notes               TEXT NULL,   -- internal/operational context
    Comments            TEXT NULL,   -- human-facing/discussion context
    CreatedAt           DATETIME NOT NULL DEFAULT (datetime('now')),
    FOREIGN KEY (UserId) REFERENCES User(UserId),
    FOREIGN KEY ({TransactionName}StatusId) REFERENCES {TransactionName}Status({TransactionName}StatusId)
);
```

---

## 5. Migration Strategy

### 5.1 Migration File Naming

Migrations follow numbered-prefix convention:

```
migrations/
├── 001-create-user-table.sql
├── 002-create-agent-site-table.sql
├── 003-add-user-email-index.sql
└── 004-create-transaction-table.sql
```

### 5.2 Migration Rules

| Rule | Detail |
|------|--------|
| Sequential numbering | 3-digit zero-padded prefix (`001-`, `002-`, ...) |
| Descriptive suffix | Verb-noun kebab-case (`create-user-table`, `add-email-index`) |
| Forward-only | No down migrations — design additive changes |
| Idempotent | Use `IF NOT EXISTS` for CREATE, `IF EXISTS` for DROP |
| One concern per file | Don't mix table creation with data seeding |

### 5.3 Migration Template

```sql
-- Migration: 001-create-user-table.sql
-- Created: 2026-04-16
-- Purpose: Create the User entity table

CREATE TABLE IF NOT EXISTS User (
    UserId      INTEGER PRIMARY KEY AUTOINCREMENT,
    Email       TEXT NOT NULL UNIQUE,
    Name        TEXT NOT NULL,
    IsActive    BOOLEAN NOT NULL DEFAULT 1,
    -- Rule 10: nullable free-text column for future-proofing
    Description TEXT NULL,
    CreatedAt   DATETIME NOT NULL DEFAULT (datetime('now')),
    UpdatedAt   DATETIME NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS IdxUser_Email ON User(Email);
CREATE INDEX IF NOT EXISTS IdxUser_CreatedAt ON User(CreatedAt);
```

---

## 6. Query Patterns

### 6.1 Standard CRUD

```sql
-- INSERT: Always specify columns explicitly
INSERT INTO User (Email, Name, IsActive, CreatedAt, UpdatedAt)
VALUES (?, ?, 1, datetime('now'), datetime('now'));

-- SELECT: Never use SELECT * — list columns
SELECT UserId, Email, Name, IsActive, CreatedAt
FROM User
WHERE IsActive = 1
ORDER BY CreatedAt DESC;

-- UPDATE: Always update UpdatedAt
UPDATE User
SET Name = ?, UpdatedAt = datetime('now')
WHERE UserId = ?;

-- DELETE: Prefer soft delete (IsActive = 0) over hard delete
UPDATE User
SET IsActive = 0, UpdatedAt = datetime('now')
WHERE UserId = ?;
```

### 6.2 View Pattern

```sql
-- Views use Vw prefix
CREATE VIEW IF NOT EXISTS VwActiveUser AS
SELECT UserId, Email, Name, CreatedAt
FROM User
WHERE IsActive = 1;
```

### 6.3 Pagination Pattern

```sql
-- Cursor-based pagination (preferred over OFFSET)
SELECT UserId, Email, Name, CreatedAt
FROM User
WHERE IsActive = 1
  AND UserId > ?  -- cursor: last seen UserId
ORDER BY UserId ASC
LIMIT 20;
```

---

## 7. Data Integrity Rules

### 7.1 Constraint Checklist

| Constraint | When to Apply |
|------------|---------------|
| `NOT NULL` | Every column unless business logic requires nullable |
| `UNIQUE` | Email, slug, external IDs, natural keys |
| `FOREIGN KEY` | Every relationship — no orphaned records |
| `DEFAULT` | `CreatedAt`, `UpdatedAt`, `IsActive`, boolean columns |
| `CHECK` | Range validation (`CHECK(Age >= 0 AND Age <= 150)`) |

### 7.2 Soft Delete Convention

```sql
-- Prefer soft delete over hard delete
-- Column: IsActive BOOLEAN NOT NULL DEFAULT 1
-- Active records: IsActive = 1
-- Deleted records: IsActive = 0

-- All queries must filter by IsActive unless explicitly querying archived data
SELECT * FROM User WHERE IsActive = 1;

-- Views should pre-filter
CREATE VIEW VwActiveUser AS SELECT ... FROM User WHERE IsActive = 1;
```

### 7.3 Timestamp Convention

| Column | Type | Default | Nullable |
|--------|------|---------|----------|
| `CreatedAt` | `DATETIME` | `datetime('now')` | `NOT NULL` |
| `UpdatedAt` | `DATETIME` | `datetime('now')` | `NOT NULL` |
| `DeletedAt` | `DATETIME` | `NULL` | Nullable (only if using soft-delete-with-timestamp pattern) |

---

## 8. ORM Integration (GORM)

### 8.1 Struct Mapping

```go
// ✅ CORRECT: GORM struct following PascalCase convention
type User struct {
    UserId    int       `gorm:"column:UserId;primaryKey;autoIncrement"`
    Email     string    `gorm:"column:Email;uniqueIndex;not null"`
    Name      string    `gorm:"column:Name;not null"`
    IsActive  bool      `gorm:"column:IsActive;not null;default:true"`
    CreatedAt time.Time `gorm:"column:CreatedAt;autoCreateTime"`
    UpdatedAt time.Time `gorm:"column:UpdatedAt;autoUpdateTime"`
}

// TableName overrides GORM's default pluralization
func (User) TableName() string {
    return "User"
}
```

### 8.2 GORM Anti-Patterns

| ❌ Wrong | ✅ Correct | Why |
|----------|-----------|-----|
| Let GORM auto-name tables (`users`) | Override `TableName()` → `User` | Singular PascalCase required |
| Use `gorm:"primaryKey"` on `ID` field | Name field `UserId` with explicit column tag | PK must be `{TableName}Id` |
| Use `gorm.Model` (embeds `ID`, `CreatedAt`, etc.) | Define fields explicitly | GORM defaults use wrong naming |

---

## 9. Schema Documentation Requirements

Every app table MUST be documented in `02-spec/23-app-db/` with:

| Section | Content |
|---------|---------|
| Table name | PascalCase singular name |
| Purpose | One-sentence description |
| Columns | Full column listing with types, constraints, defaults |
| Indexes | All indexes with rationale |
| Foreign keys | All relationships with referenced tables |
| Sample queries | Common CRUD operations |
| Migration file | Reference to migration file that creates it |

### Documentation File Template

```markdown

# {TableName} Table

**Version:** 1.0.0
**Migration:** `001-create-{table-name}-table.sql`

## Purpose

{One sentence describing what this table stores}

## Schema

| Column | Type | Nullable | Default | Constraint | Description |
|--------|------|----------|---------|------------|-------------|
| {TableName}Id | INTEGER | NOT NULL | AUTOINCREMENT | PK | Primary key |
| ... | ... | ... | ... | ... | ... |

## Indexes

| Index | Columns | Rationale |
|-------|---------|-----------|
| ... | ... | ... |

## Relationships

| FK Column | References | On Delete |
|-----------|------------|-----------|
| ... | ... | CASCADE / RESTRICT |
```

---

## 10. Anti-Patterns

| ❌ Anti-Pattern | ✅ Correct Approach |
|-----------------|---------------------|
| Plural table names (`Users`) | Singular (`User`) |
| `snake_case` columns (`created_at`) | PascalCase (`CreatedAt`) |
| UUID primary keys | INTEGER AUTOINCREMENT |
| `SELECT *` in queries | Explicit column list |
| Hard delete (`DELETE FROM`) | Soft delete (`IsActive = 0`) |
| OFFSET-based pagination | Cursor-based pagination |
| Missing `UpdatedAt` on mutations | Always set `UpdatedAt = datetime('now')` |
| Storing enums as TEXT | SMALLINT + lookup table |
| Missing indexes on FK columns | Index every FK column |
| No foreign key constraints | Always define FK with ON DELETE behavior |

---

## 11. Cross-References

| Topic | File |
|-------|------|
| Core database naming conventions | [21-database-conventions.md](./21-database-conventions.md) |
| Split DB architecture | [08-split-db-architecture.md](./08-split-db-architecture.md) |
| Seedable configuration | [09-seedable-config.md](./09-seedable-config.md) |
| App-specific features | [16-app.md](./16-app.md) |
| Enum standards (for lookup tables) | [07-enum-standards.md](./07-enum-standards.md) |
| App design system (related app spec) | [19-app-design-system-and-ui.md](./19-app-design-system-and-ui.md) |
| Source spec folder | `02-spec/23-app-db/` |

---

*App database consolidated — created 2026-04-16*
