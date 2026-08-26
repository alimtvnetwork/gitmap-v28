# Split DB Pattern

**Version:** 3.2.0  
**Updated:** 2026-04-16

---

## Overview

The **Split DB** pattern uses multiple small SQLite database files — one per bounded context — instead of a single monolithic database. Each domain concern owns its data in isolation.

---

## 1. Core Concept

```
Monolithic (❌ avoid)          Split DB (✅ preferred)
┌─────────────────────┐        ┌──────────────┐  ┌──────────────┐
│   app.db            │        │ transactions │  │  snapshots   │
│                     │        │     .db      │  │     .db      │
│  Transactions       │        │              │  │              │
│  Snapshots          │   →    │ Transactions │  │ Snapshots    │
│  AgentSites         │        │ StatusTypes  │  │ SnapshotMeta │
│  StatusTypes        │        └──────────────┘  └──────────────┘
│  SnapshotMeta       │
│  Plugins            │        ┌──────────────┐  ┌──────────────┐
│  Themes             │        │ agent-sites  │  │   plugins    │
│  ...                │        │     .db      │  │     .db      │
└─────────────────────┘        │              │  │              │
                               │ AgentSites   │  │ Plugins      │
                               │ SiteConfigs  │  │ Themes       │
                               └──────────────┘  └──────────────┘
```

### 1.1 One Database Per Bounded Context

A **bounded context** is a domain area with cohesive data that changes together. Each gets its own `.db` file.

| Bounded Context | Database File | Tables (examples) |
|----------------|---------------|-------------------|
| Transactions | `transactions.db` | `Transactions`, `StatusTypes`, `TransactionItems` |
| Snapshots | `snapshots.db` | `Snapshots`, `SnapshotMeta`, `SnapshotFiles` |
| Agent Sites | `agent-sites.db` | `AgentSites`, `SiteConfigs`, `SiteDomains` |
| Plugins | `plugins.db` | `Plugins`, `Themes`, `PluginVersions` |
| Auth | `auth.db` | `Users`, `Sessions`, `UserRoles`, `Roles` |
| Audit | `audit.db` | `AuditLog`, `EventTypes` |

### 1.2 Naming Convention

Database files use **kebab-case** with `.db` extension:

```
✅ transactions.db
✅ agent-sites.db
✅ snapshot-meta.db
❌ AgentSites.db        (PascalCase is for tables/columns, not files)
❌ agent_sites.db       (snake_case not allowed for files)
❌ app.db               (too generic — name by domain)
```

> **Rule:** File names use kebab-case per [../02-coding-guidelines/08-file-folder-naming/00-overview.md](../02-coding-guidelines/08-file-folder-naming/00-overview.md). Table/column names inside the DB use PascalCase per [./01-naming-conventions.md](./01-naming-conventions.md).

---

## 2. Why Split DB

| Benefit | Explanation |
|---------|-------------|
| **Isolation** | A corrupt or locked DB affects only one domain |
| **Performance** | Each DB has its own WAL — no cross-domain lock contention |
| **Portability** | Copy/backup a single domain independently |
| **Testability** | Spin up one in-memory DB per domain in tests |
| **Deployment** | Deploy/migrate domains independently |
| **Schema clarity** | Each DB schema is small and easy to reason about |

| Risk | Mitigation |
|------|------------|
| Cross-domain queries | Use application-level joins (see Section 4) |
| Transaction spanning domains | Use eventual consistency or saga pattern |
| Connection management | Use a DB registry (see Section 3) |
| More files to manage | Automate with migration tooling |

---

## 3. Architecture

### 3.1 Directory Layout

```
data/
├── transactions.db
├── snapshots.db
├── agent-sites.db
├── plugins.db
├── auth.db
└── audit.db
```

All `.db` files live in a single `data/` directory at the project root.

### 3.2 DB Registry Pattern

A central registry maps domain names to database connections:

```
┌─────────────────────────────────────────┐
│              DB Registry                │
│                                         │
│  "transactions" → transactions.db conn  │
│  "snapshots"    → snapshots.db conn     │
│  "agent-sites"  → agent-sites.db conn   │
│  "plugins"      → plugins.db conn       │
│  "auth"         → auth.db conn          │
│  "audit"        → audit.db conn         │
└─────────────────────────────────────────┘
         │
         ▼
   ORM layer uses registry
   to get correct connection
```

### 3.3 Implementation (Go Example)

```go
// DbRegistry manages connections to split databases
type DbRegistry struct {
    connections map[string]*sql.DB
    dataDir     string
}

// Get returns the connection for a bounded context
func (r *DbRegistry) Get(context string) (*sql.DB, error) {
    if conn, ok := r.connections[context]; ok {
        return conn, nil
    }
    return nil, fmt.Errorf("unknown context: %s", context)
}

// Usage
db, _ := registry.Get("transactions")
```

### 3.4 Implementation (PHP Example)

```php
class DbRegistry
{
    /** @var array<string, PDO> */
    private array $connections = [];

    public function get(string $context): PDO
    {
        if (!isset($this->connections[$context])) {
            throw new RuntimeException("Unknown context: {$context}");
        }
        return $this->connections[$context];
    }
}

// Usage
$db = $registry->get('transactions');
```

---

## 4. Cross-Domain Data Access

### 4.1 The Rule

> **Cross-domain queries are NOT allowed at the database level.** If Domain A needs data from Domain B, it goes through Domain B's service/repository layer.

### 4.2 Pattern: Application-Level Join

```
❌ WRONG — SQL join across databases
SELECT t.*, a.SiteName
FROM transactions.Transactions t
JOIN agent_sites.AgentSites a ON t.AgentSiteId = a.AgentSiteId

✅ CORRECT — Application-level composition
1. TransactionService.GetTransaction(id)     → Transaction{AgentSiteId: 5}
2. AgentSiteService.GetAgentSite(5)          → AgentSite{SiteName: "example.com"}
3. Compose in application layer              → TransactionWithSite{...}
```

### 4.3 Shared Identifiers

Domains reference each other via **ID only**. The foreign key is stored but NOT enforced across databases:

```sql
-- In transactions.db
-- linter-waive: MISSING-DESC-001 reason="Split-DB pattern example; focus on cross-DB references"
CREATE TABLE Transaction (
    TransactionId   INTEGER PRIMARY KEY AUTOINCREMENT,
    AgentSiteId     INTEGER NOT NULL,  -- References agent-sites.db but NO FK constraint
    StatusTypeId    INTEGER NOT NULL,
    FOREIGN KEY (StatusTypeId) REFERENCES StatusTypes(StatusTypeId)  -- Same DB = FK OK
    -- NO FOREIGN KEY for AgentSiteId (different DB)
);
```

> **Rule:** Foreign keys are enforced ONLY within the same database file. Cross-database references store the ID but rely on application-layer validation.

### 4.4 Decision Flow: Same DB or Different DB?

```
Do these tables change together in the same transaction?
├── YES → Same database file
└── NO
    ├── Does Table A need FK-enforced integrity with Table B?
    │   ├── YES → Same database file
    │   └── NO → Different database files
    └── Are they in the same bounded context?
        ├── YES → Same database file
        └── NO → Different database files
```

---

## 5. Migrations

Each database has its own migration folder:

```
migrations/
├── transactions/
│   ├── 001_create_status_types.sql
│   ├── 002_create_transactions.sql
│   └── 003_add_transaction_items.sql
├── snapshots/
│   ├── 001_create_snapshots.sql
│   └── 002_add_snapshot_meta.sql
├── agent-sites/
│   ├── 001_create_agent_sites.sql
│   └── 002_add_site_configs.sql
└── auth/
    ├── 001_create_users.sql
    ├── 002_create_roles.sql
    └── 003_create_user_roles.sql
```

### 5.1 Migration Rules

1. Migrations are **numbered sequentially** per database (not globally)
2. Each migration targets **one database only**
3. Migration file names use **snake_case** with numeric prefix
4. Migrations are **forward-only** — no down migrations

---

## 6. Testing with Split DB

### 6.1 In-Memory Per Domain

```go
// Test setup — each domain gets its own in-memory DB
func setupTestDBs() *DbRegistry {
    registry := NewDbRegistry()
    registry.Register("transactions", ":memory:")
    registry.Register("snapshots", ":memory:")
    registry.Register("auth", ":memory:")
    return registry
}
```

### 6.2 Test Isolation

| Strategy | Description |
|----------|-------------|
| Per-test DB | Create fresh in-memory DB for each test |
| Per-domain DB | Each domain test suite uses its own DB |
| Shared fixtures | Lookup tables (StatusTypes, Roles) seeded once per suite |

> Cross-reference: [./04-testing-strategy.md](./04-testing-strategy.md) for full testing conventions.

---

## 7. When NOT to Use Split DB

Split DB is the **default**, but consolidate into fewer databases when:

| Scenario | Action |
|----------|--------|
| Two domains always transact together | Merge into one DB |
| Domain has only 1-2 tables | Consider merging with a related domain |
| Team explicitly needs cross-domain FK integrity | Merge or use MySQL |
| Write concurrency exceeds SQLite limits | Move affected domain to MySQL |

---

## 8. Checklist — New Domain Setup

When adding a new bounded context:

- [ ] Create `{domain-name}.db` file in `data/`
- [ ] Register in DB registry
- [ ] Create `migrations/{domain-name}/` folder
- [ ] Write initial migration(s)
- [ ] Define ORM models scoped to this DB
- [ ] Create service/repository layer for the domain
- [ ] Add in-memory test setup
- [ ] Document in this file (Section 1.1 table)

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Schema design & key sizing | [./02-schema-design.md](./02-schema-design.md) |
| Naming conventions | [./01-naming-conventions.md](./01-naming-conventions.md) |
| ORM and views | [./03-orm-and-views.md](./03-orm-and-views.md) |
| Testing strategy | [./04-testing-strategy.md](./04-testing-strategy.md) |
| File & folder naming | [../02-coding-guidelines/08-file-folder-naming/00-overview.md](../02-coding-guidelines/08-file-folder-naming/00-overview.md) |
| Database conventions overview | [./00-overview.md](./00-overview.md) |
