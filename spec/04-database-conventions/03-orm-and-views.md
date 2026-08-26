# ORM Usage and Database Views

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

All database interactions in business logic MUST go through an **ORM or query builder**. Raw SQL is only permitted in migrations and database view definitions. Complex joins MUST be pre-defined as **database views** so the business layer queries a flat result.

---

## 1. ORM-First Rule

### 1.1 The Rule

> **Never write raw SQL in the business/service layer.** Use the best ORM or query builder available for the language.

| Layer | Raw SQL Allowed? | What to Use |
|-------|-----------------|-------------|
| Business logic / services | ❌ No | ORM methods |
| Repository / data access | ❌ No | ORM / query builder |
| Migrations | ✅ Yes | Raw DDL statements |
| View definitions | ✅ Yes | `CREATE VIEW` statements |
| One-off scripts | ✅ Yes | With approval |

### 1.2 Recommended ORMs by Language

| Language | ORM / Query Builder | Why |
|----------|-------------------|-----|
| **Go** | `sqlc` or `GORM` | Type-safe generated code (sqlc) or full ORM (GORM) |
| **PHP** | Custom `Orm` class or Eloquent | Project uses custom PascalCase-aware Orm |
| **TypeScript** | Prisma or Drizzle | Type-safe, schema-first |
| **Rust** | Diesel or SeaORM | Compile-time query validation |
| **C#** | Entity Framework Core | Industry standard, LINQ queries |

### 1.3 Examples

#### ❌ Wrong — Raw SQL in Business Logic

```go
// Service layer — FORBIDDEN
func (s *TransactionService) GetPending() ([]Transaction, error) {
    rows, err := s.db.Query(
        "SELECT TransactionId, PluginSlug, Amount FROM Transactions WHERE StatusTypeId = 1 ORDER BY CreatedAt DESC",
    )
    // manual scanning...
}
```

#### ✅ Correct — ORM in Business Logic

```go
// Service layer — uses repository with ORM
func (s *TransactionService) GetPending() ([]Transaction, error) {
    return s.repo.FindAll(TransactionFilter{
        StatusTypeId: statustype.Pending,
        OrderBy:      "CreatedAt DESC",
    })
}
```

```php
// PHP — uses Orm class
$orm = new Orm(TableType::Transactions->value);
$pending = $orm->findAll(
    ['StatusTypeId' => StatusType::Pending->value],
    'CreatedAt DESC'
);
```

```typescript
// TypeScript — uses Prisma
const pending = await prisma.transactions.findMany({
    where: { StatusTypeId: StatusType.Pending },
    orderBy: { CreatedAt: "desc" },
});
```

---

## 2. Database Views — No On-the-Fly Joins

### 2.1 The Rule

> **When business logic needs data from multiple tables, create a database VIEW.** Do not write JOIN queries in the application code.

### 2.2 Why Views?

| Aspect | On-the-fly JOIN | Database VIEW |
|--------|----------------|---------------|
| Maintainability | JOIN logic scattered across codebase | Defined once in migration |
| Performance | No query plan caching guarantee | DB can optimize the view |
| Testability | Must test JOIN logic in every caller | Test the view once |
| Readability | Complex SQL embedded in business code | Simple `SELECT * FROM VwName` |
| Reusability | Copy-paste JOINs | Single source of truth |

### 2.3 View Naming Convention

Views use PascalCase with a `Vw` prefix:

```
VwTransactionDetails
VwActiveAgentSites
VwUserRoleSummary
VwPendingSnapshotJobs
```

### 2.4 Example — Creating and Using a View

#### Step 1: Define the View in a Migration

```sql
-- Migration: create VwTransactionDetails
CREATE VIEW VwTransactionDetails AS
SELECT
    t.TransactionId,
    t.PluginSlug,
    t.Amount,
    t.CreatedAt,
    st.Name       AS StatusName,
    ft.Name       AS FileTypeName,
    a.SiteName    AS AgentSiteName,
    a.SiteUrl     AS AgentSiteUrl
FROM Transactions t
INNER JOIN StatusTypes st ON t.StatusTypeId = st.StatusTypeId
INNER JOIN FileTypes ft   ON t.FileTypeId = ft.FileTypeId
LEFT JOIN AgentSites a    ON t.AgentSiteId = a.AgentSiteId;
```

#### Step 2: Query the View via ORM (Flat Result)

```go
// Go — query the view like a simple table
type TransactionDetail struct {
    TransactionId int64  `db:"TransactionId"`
    PluginSlug    string `db:"PluginSlug"`
    Amount        float64 `db:"Amount"`
    StatusName    string `db:"StatusName"`
    FileTypeName  string `db:"FileTypeName"`
    AgentSiteName string `db:"AgentSiteName"`
}

details, err := repo.FindAllFromView("VwTransactionDetails", filter)
```

```php
// PHP — query view through Orm
$orm = new Orm('VwTransactionDetails');
$details = $orm->findAll(['StatusName' => 'Pending']);
```

### 2.5 When to Create a View

Create a database view when:

- A query JOINs **2 or more tables**
- The same JOIN is needed in **more than one place**
- The business layer needs a **flattened result** from related tables
- A report or dashboard aggregates data across tables

### 2.6 View Maintenance

- Views are created/updated via **migrations** (never manually)
- Views are versioned alongside table migrations
- When underlying tables change, update dependent views in the same migration
- Document views in the schema documentation alongside tables

---

## 3. Raw SQL Boundaries

### 3.1 Where Raw SQL IS Allowed

```go
// ✅ Migration — raw DDL
func migrateV5(db *sql.DB) error {
    _, err := db.Exec(`
        CREATE VIEW VwTransactionDetails AS
        SELECT t.TransactionId, t.PluginSlug, st.Name AS StatusName
        FROM Transactions t
        INNER JOIN StatusTypes st ON t.StatusTypeId = st.StatusTypeId
    `)

    return err
}
```

### 3.2 Where Raw SQL is FORBIDDEN

```go
// ❌ Service layer — FORBIDDEN
func (s *Service) GetReport() {
    rows, _ := s.db.Query(`
        SELECT t.*, st.Name
        FROM Transactions t
        JOIN StatusTypes st ON t.StatusTypeId = st.StatusTypeId
        WHERE t.CreatedAt > ?
    `, cutoff)
}

// ✅ Instead — use the view via ORM
func (s *Service) GetReport() {
    return s.repo.FindFromView("VwTransactionDetails", filter)
}
```

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Schema design | [./02-schema-design.md](./02-schema-design.md) |
| Testing strategy | [./04-testing-strategy.md](./04-testing-strategy.md) |
| Naming conventions | [./01-naming-conventions.md](./01-naming-conventions.md) |
