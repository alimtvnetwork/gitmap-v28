# Database Schema Design

**Version:** 3.3.0  
**Updated:** 2026-04-19

---

## Overview

Rules for designing database schemas that are efficient, maintainable, and correctly normalized. Covers key sizing, primary key strategy, normalization, and the Split DB pattern.

---

## 1. Primary Key Strategy

### 1.1 Always Use Integer Primary Keys

Primary keys MUST be integer-based. Choose the **smallest type** that fits the expected data volume:

| Expected Rows (10-year horizon) | Key Type | Range | Storage |
|--------------------------------|----------|-------|---------|
| < 32,000 | `SMALLINT` | ±32K | 2 bytes |
| < 2 billion | `INTEGER` (default) | ±2.1B | 4 bytes |
| > 2 billion | `BIGINT` | ±9.2 quintillion | 8 bytes |

### 1.2 Decision Flow

```
How many rows in 10 years?
├── < 32,000 → SMALLINT
├── < 2,000,000,000 → INTEGER (this is the default)
└── > 2,000,000,000 → BIGINT
```

### 1.3 UUID/GUID — Avoid

| Aspect | INTEGER | UUID |
|--------|---------|------|
| Storage | 4 bytes | 16 bytes (4x larger) |
| Index performance | Fast (sequential) | Slow (random distribution) |
| Readability | Easy to debug | Hard to read |
| Fragmentation | None | High (random inserts) |

> **Rule:** ❌ Do NOT use UUID/GUID as primary key unless there is an **explicit requirement** (e.g., distributed systems with no central authority, public-facing IDs that must not be guessable).

```sql
-- ❌ AVOID
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE User (
    UserId TEXT PRIMARY KEY DEFAULT (lower(hex(randomblob(16))))
);

-- ✅ PREFERRED
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE User (
    UserId INTEGER PRIMARY KEY AUTOINCREMENT
);
```

### 1.4 When UUID Is Acceptable

Only use UUID when **all** of these are true:
1. Records are created across multiple disconnected systems
2. There is no central ID authority
3. IDs must be publicly exposed and non-guessable
4. The team has explicitly approved UUID for this table

If UUID is used, store it as `BLOB(16)` (not `TEXT(36)`) for storage efficiency.

---

## 2. Key Sizing — Smallest Possible Type

Apply the smallest-type principle to ALL columns, not just primary keys:

| Data | ❌ Oversized | ✅ Right-Sized |
|------|-------------|---------------|
| Status (5 values) | `TEXT` | `TINYINT` + lookup table |
| Age | `INTEGER` | `TINYINT` (0-255) |
| Year | `INTEGER` | `SMALLINT` (0-65535) |
| Boolean | `INTEGER` | `TINYINT(1)` or `BOOLEAN` |
| Country code | `TEXT` | `CHAR(2)` |
| Currency amount | `REAL` | `DECIMAL(10,2)` |

---

## 3. Normalization — Repeated Values Become Tables

### 3.1 The Rule

> **Any column that contains a repeated set of values MUST be extracted into a separate lookup table with a foreign key relationship.**

This applies to: status types, file types, category types, role types, priority levels, etc.

### 3.2 ❌ Wrong — Repeated Strings

```sql
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE Transaction (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    Status        TEXT,     -- 'Pending', 'Complete', 'Failed' repeated thousands of times
    FileType      TEXT      -- 'Plugin', 'Theme', 'MuPlugin' repeated thousands of times
);
```

### 3.3 ✅ Correct — Normalized with Lookup Tables

```sql
-- Lookup table for statuses
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE StatusType (
    StatusTypeId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name         TEXT NOT NULL UNIQUE   -- 'Pending', 'Complete', 'Failed'
);

-- Lookup table for file types
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE FileType (
    FileTypeId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name       TEXT NOT NULL UNIQUE   -- 'Plugin', 'Theme', 'MuPlugin'
);

-- Main table references lookup tables via FK
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE Transaction (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    StatusTypeId  INTEGER NOT NULL,
    FileTypeId    INTEGER NOT NULL,
    FOREIGN KEY (StatusTypeId) REFERENCES StatusType(StatusTypeId),
    FOREIGN KEY (FileTypeId)   REFERENCES FileType(FileTypeId)
);
```

### 3.4 Many-to-Many (N-to-M) Relationships

When an entity can have multiple values of a type (e.g., a user has multiple roles):

```sql
-- Entity tables
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE User (
    UserId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name   TEXT NOT NULL
);

-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE Role (
    RoleId INTEGER PRIMARY KEY AUTOINCREMENT,
    Name   TEXT NOT NULL UNIQUE   -- 'Admin', 'Editor', 'Viewer'
);

-- Junction table for N-to-M
-- linter-waive: MISSING-DESC-001 reason="PK / normalization / FK example; Rules 10/11/12 demonstrated separately in §6.4"
CREATE TABLE UserRole (
    UserRoleId INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId     INTEGER NOT NULL,
    RoleId     INTEGER NOT NULL,
    UNIQUE (UserId, RoleId),
    FOREIGN KEY (UserId) REFERENCES User(UserId),
    FOREIGN KEY (RoleId) REFERENCES Role(RoleId)
);
```

---

## 4. Database Engine — SQLite First (Split DB)

### 4.1 Default: SQLite with Split DB Pattern

The project follows the **Split DB** pattern: multiple small SQLite databases per domain concern rather than one monolithic database.

| Advantage | Description |
|-----------|-------------|
| Zero config | No server process, no credentials |
| Portable | Single file per database, easy to backup/copy |
| Isolation | Domain failures don't cascade |
| Performance | Each DB has its own WAL, no lock contention across domains |
| Testable | In-memory mode for fast tests |

> See [07-split-db-pattern.md](./07-split-db-pattern.md) for the full Split DB specification including directory layout, DB registry, cross-domain rules, and migration strategy.

### 4.2 Fallback: MySQL

Use MySQL only when:
- High concurrent write volume exceeds SQLite's single-writer model
- Multi-server access to the same database is required
- The application requires server-level replication

All naming conventions (PascalCase) apply equally to MySQL.

---

## 5. Schema Documentation

Every database schema MUST be documented with:

1. **Table purpose** — one-line description
2. **Column definitions** — name, type, constraints, description
3. **Relationships** — FK references with cardinality
4. **Indexes** — which columns and why
5. **Expected volume** — estimated rows in 10 years (drives key sizing)

### Template

```markdown

### TableName

**Purpose:** [What this table stores]  
**Expected volume:** [N rows in 10 years]

| Column | Type | Constraints | Description |
|--------|------|-------------|-------------|
| TableNameId | INTEGER | PK, AUTOINCREMENT | Primary key |
| ForeignTableId | INTEGER | FK, NOT NULL | References ForeignTable |
| Name | TEXT | NOT NULL | Human-readable name |
| CreatedAt | INTEGER | NOT NULL, DEFAULT (unixepoch()) | Epoch seconds, UTC (Naming Rule 7.1) |
```

---

## 6. Mandatory Descriptive Free-Text Columns

> **Aligned with:** Naming Conventions **v3.5.0** — Rules 10, 11, and 12.

This section is the **schema-design counterpart** of [Naming Conventions Rules 10, 11, and 12](./01-naming-conventions.md#descriptive-free-text-column-naming). Naming conventions define **what the columns are called**; this section defines **which tables must include them and how they are placed in the schema**. The two documents MUST stay aligned — if Rule 10/11/12 changes, this section changes.

Every table (except pure join/pivot tables) MUST reserve at least one nullable free-text column to absorb future context that the schema did not anticipate. The exact column depends on the **table category**:

- **Entity / reference / master-data tables** → `Description TEXT NULL` (Rule 10)
- **Transactional / invoice / billing / payment / order tables** → `Notes TEXT NULL` **and** `Comments TEXT NULL` (Rule 11)
- **All three columns** → MUST be nullable; they are optional context fields, not required business data (Rule 12)

> **Why:** Schemas evolve slower than business needs. A nullable text column lets operators, support staff, and downstream tools attach context (explanations, hints, audit notes, dispute reasons) **without a migration**. Cost is near-zero (NULL rows store no data); benefit is operational flexibility.

### 6.1 Table Categories and Required Columns

| Category | Examples | Required Nullable Columns | Rule | Purpose |
|----------|----------|---------------------------|------|---------|
| **Lookup / reference** | `Country`, `Currency`, `Status`, `Role` | `Description TEXT NULL` | Rule 10 | Explain what each row means; hint for UI tooltips |
| **Entity / master data** | `AgentSite`, `User`, `Plugin`, `Product` | `Description TEXT NULL` | Rule 10 | Free-form context, help text, future-proofing |
| **Transactional / event** | `Transaction`, `Invoice`, `Order`, `Payment`, `UserBill` | `Notes TEXT NULL` **and** `Comments TEXT NULL` | Rule 11 | Capture exceptions, dispute reasons, manual overrides |
| **Audit / log** | `CommandHistory`, `AuditLog` | `Notes TEXT NULL` | Rule 11 (notes-only) | Reviewer remarks |
| **Join / pivot** | `GroupItems`, `UserRole` | None required | — | Pure relational glue |

### 6.2 Column Semantics

The naming guideline (Rule 11) distinguishes the two transactional columns by **audience**:

| Column | Cardinality | Who Writes It | Typical Content | Per Naming Rule |
|--------|-------------|---------------|-----------------|-----------------|
| `Description` | One per row, mostly static | Schema designer / admin | "What this row represents" — set once, rarely changes | Rule 10 |
| `Notes` | One per row, mutable | System or operator | Internal or operational context — back-office, system-generated, processing notes | Rule 11 (`Notes` → internal/operational) |
| `Comments` | One per row, mutable | End user or support staff | Human-facing or discussion-oriented context — customer remarks, dispute text, manual overrides | Rule 11 (`Comments` → human-facing) |

> Both `Notes` and `Comments` exist on transactional tables because they serve **different audiences**: `Notes` is for back-office/system context; `Comments` is for human-to-human conversation about the record. This is the wording used in Rule 11 of the naming spec — do not paraphrase it elsewhere.

### 6.3 Schema-Design Rules (alignment with Naming Rules 10–12)

1. **Always nullable** (Rule 12) — these columns MUST be `TEXT NULL` with no `NOT NULL` constraint and no `DEFAULT` value. A NULL value means "nothing to add", which is the common case.
2. **Never indexed** — they are read-on-detail, not searched. No `Idx{Table}_Description` / `Idx{Table}_Notes` / `Idx{Table}_Comments` indexes.
3. **No length cap in SQLite**; in MySQL use `TEXT` (not `VARCHAR`) so length is effectively unbounded.
4. **Never required by the application** — code MUST treat them as informational only. No business logic may branch on their content.
5. **Excluded from list/index views by default** — fetch only when the user opens a detail view, to keep row payloads small.
6. **PascalCase, exact spelling** — `Description`, `Notes`, `Comments`. No `desc`, `note`, `comment`, `Remarks`, `Memo`, etc. The names are fixed by Rules 10/11.

### 6.4 Examples

```sql
-- ✅ Entity table — gains Description (Rule 10), epoch CreatedAt (Rule 7.1)
CREATE TABLE AgentSite (
    AgentSiteId  INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName     TEXT NOT NULL,
    Description  TEXT NULL,                                  -- future-proofing / help / hint
    CreatedAt    INTEGER NOT NULL DEFAULT (unixepoch())      -- epoch seconds, UTC
);

-- ✅ Lookup / enum table — canonical (Id, Code, Label) shape (Rule 13)
--    See §6.5 below. Replaces older `{Table}Code` style.
CREATE TABLE TransactionStatus (
    TransactionStatusId  INTEGER PRIMARY KEY AUTOINCREMENT,
    Code                 TEXT NOT NULL UNIQUE,               -- machine identifier, e.g. 'PendingSettlement'
    Label                TEXT NOT NULL,                      -- human-readable, e.g. 'Pending settlement'
    Description          TEXT NULL                           -- "Pending settlement with bank"
);

-- ✅ Transactional table — gains Notes AND Comments (Rule 11), epoch CreatedAt
CREATE TABLE Transaction (
    TransactionId        INTEGER PRIMARY KEY AUTOINCREMENT,
    AgentSiteId          INTEGER NOT NULL,
    Amount               REAL    NOT NULL,
    TransactionStatusId  INTEGER NOT NULL,
    Notes                TEXT NULL,                          -- internal or operational context
    Comments             TEXT NULL,                          -- human-facing or discussion-oriented context
    CreatedAt            INTEGER NOT NULL DEFAULT (unixepoch()),
    FOREIGN KEY (AgentSiteId)         REFERENCES AgentSite(AgentSiteId),
    FOREIGN KEY (TransactionStatusId) REFERENCES TransactionStatus(TransactionStatusId)
);

-- ✅ Invoice/bill table — same pattern (Rule 11)
CREATE TABLE UserBill (
    UserBillId   INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId       INTEGER NOT NULL,
    AmountDue    REAL    NOT NULL,
    DueDate      INTEGER NOT NULL,                           -- epoch seconds (date-of-day = midnight UTC)
    Notes        TEXT NULL,                                  -- "Adjusted for credit memo #4421"
    Comments     TEXT NULL,                                  -- "Customer disputes line 3"
    FOREIGN KEY (UserId) REFERENCES User(UserId)
);
```

### 6.5 Rule 13 — Enum / Lookup Table Canonical Shape

> **Effective v5.13.0.** Every enum-like reference table (status, type, role, kind, category, channel, etc.) MUST use exactly the column set below.

**The primary-key rule (Rule 1) is universal and is NOT relaxed by Rule 13.** The PK column on every table — entity, transactional, ref, join — is `{TableName}Id`. The "simplification" introduced by Rule 13 applies **only** to the descriptive columns `Code`, `Label`, and `Description` (no `{Table}` prefix on those three).

| Column | Type | Constraint | Purpose |
|--------|------|------------|---------|
| `{TableName}Id` | INTEGER | PK, AUTOINCREMENT | Surrogate key. **Always `{TableName}Id`** — e.g. `WorkerNodeStatusId`, `RoleId`, `TransactionStatusId`. Never bare `Id`. |
| `Code` | TEXT | NOT NULL, UNIQUE | Stable machine identifier (PascalCase). Application code references rows by `Code`, not by `Id`. **Bare `Code` — no `{Table}` prefix.** |
| `Label` | TEXT | NOT NULL | Default human-readable label (English). I18n is handled by a separate `{Table}Translation` table when needed. **Bare `Label` — no `{Table}` prefix.** |
| `Description` | TEXT | NULL | Optional explanatory hint (Rule 10). **Bare `Description` — no `{Table}` prefix.** |

**Forbidden alternatives** (deprecated and removed in v5.13.0):

| ❌ Forbidden | ✅ Required | Why |
|-------------|------------|-----|
| `Id` (bare PK) | `{TableName}Id` | Universal Rule 1; ambiguous in JOINs |
| `{Table}Code` (e.g. `RoleCode`, `WorkerNodeStatusCode`) | `Code` | Rule 13 simplification |
| `{Table}Label` (e.g. `RoleLabel`) | `Label` | Rule 13 simplification |
| `Name`, `DisplayName`, `Title` (instead of `Label`) | `Label` | Single canonical name |
| `{Table}Description` (e.g. `RoleDescription`) | `Description` | Rule 13 simplification |
| Mixing `Code` + `Slug` on the same enum table | Pick one (`Code`) | One stable identifier per row |

**Rationale:** PK keeps the `{TableName}Id` prefix because PKs travel into other tables as FKs, where the `{Table}Id` prefix is what makes the FK self-documenting. The descriptive columns (`Code`, `Label`, `Description`) never travel — they live only on the enum row itself — so the prefix adds noise without information.

**Canonical example:**

```sql
CREATE TABLE WorkerNodeStatus (
    WorkerNodeStatusId  INTEGER PRIMARY KEY AUTOINCREMENT,  -- Rule 1: {TableName}Id
    Code                TEXT NOT NULL UNIQUE,               -- Rule 13: bare 'Code'
    Label               TEXT NOT NULL,                      -- Rule 13: bare 'Label'
    Description         TEXT NULL                           -- Rule 13: bare 'Description'
);
```

**Lookup pattern:**
```sql
SELECT WorkerNodeStatusId FROM WorkerNodeStatus WHERE Code = 'Active';
```

**Seeding pattern** — see [`14-rbac-and-status-seed.md`](../19-main-worker-service/14-rbac-and-status-seed.md) for the canonical seed format used by `19-main-worker-service`.

```sql
-- ❌ WRONG — entity table without Description (violates Rule 10)
-- linter-waive: MISSING-DESC-001 reason="Intentional anti-example for Rule 10 — teaches the violation"
CREATE TABLE AgentSite (
    AgentSiteId INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName    TEXT NOT NULL
);

-- ❌ WRONG — transactional table missing Notes and Comments (violates Rule 11)
-- linter-waive: MISSING-DESC-001 reason="Intentional anti-example for Rule 11 — teaches the violation"
CREATE TABLE Transaction (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    AgentSiteId   INTEGER NOT NULL,
    Amount        REAL    NOT NULL
);

-- ❌ WRONG — Description marked NOT NULL (violates Rule 12)
-- linter-waive: MISSING-DESC-001 reason="Intentional anti-example for Rule 12 — teaches the violation"
CREATE TABLE AgentSite (
    AgentSiteId  INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName     TEXT NOT NULL,
    Description  TEXT NOT NULL DEFAULT ''      -- forbidden: must be nullable
);
```

### 6.5 Codegen / AI Implementer Contract

When generating or reviewing a `CREATE TABLE` statement:

1. **Classify the table** using §6.1 (lookup, entity, transactional, audit, or join).
2. **Inject the required nullable columns** if missing — `Description` for lookup/entity (Rule 10), `Notes` + `Comments` for transactional (Rule 11), `Notes` only for audit/log (Rule 11 notes-only variant).
3. **Place them last** in the column list, just before timestamps and foreign-key clauses, so they read as "metadata".
4. **Never** mark them `NOT NULL`, never give them a `DEFAULT` value, never index them (Rule 12 + §6.3).
5. **Skip** join/pivot tables — they exist solely to relate two PKs and gain nothing from text columns.
6. **Use the exact column names** `Description`, `Notes`, `Comments` — no synonyms (`Memo`, `Remarks`, etc.). Rules 10/11 fix the names.

> **Naming-only summary (mirrors Rule 12 callout):**
> - Entity / reference / master tables → `Description`
> - Transaction / invoice / bill / payment / order tables → `Notes`, `Comments`
> - All three use PascalCase and remain nullable

### 6.6 Linter Waivers (MISSING-DESC-001 v1.1)

The `MISSING-DESC-001` SQL linter enforces Rules 10/11/12 across every
`CREATE TABLE` it sees in `.sql` files and inside ` ```sql ` markdown
fences. Some examples in this spec tree exist to teach a different
concept (PK choice, normalization, FK syntax, junction-table mechanics).
Adding `Description`/`Notes`/`Comments` to those examples would obscure
the lesson, so the linter accepts a **waiver** comment with a mandatory
`reason="..."` clause.

**Block-level waiver** (skip the next `CREATE TABLE` only):

```sql
-- linter-waive: MISSING-DESC-001 reason="PK-sizing example; Description omitted for clarity"
CREATE TABLE User (
    UserId INTEGER PRIMARY KEY AUTOINCREMENT
);
```

The waiver MUST appear within the **5 non-blank lines** immediately
preceding `CREATE TABLE` and the chain must consist only of `--`
comments — any other SQL statement breaks the chain.

**File-level waiver** (skip every block in the file):

```sql
-- linter-waive-file: MISSING-DESC-001 reason="Migration history; legacy schema frozen"
```

**Rules for waivers:**

1. The `reason="..."` clause is **mandatory**. Bare waivers are ignored
   so silent suppressions can't sneak through review.
2. Waivers are **per-rule** — `MISSING-DESC-001` only. Other linters
   ignore them.
3. **When to waive:** pedagogical mini-examples, ASCII/Mermaid blocks
   that look like SQL, frozen legacy schemas, runtime-generated DDL.
4. **When NOT to waive:** real production schemas, migration files, or
   any spec example showing a table the app will actually create.

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Naming conventions (Rules 10/11/12 — canonical wording) | [./01-naming-conventions.md#descriptive-free-text-column-naming](./01-naming-conventions.md#descriptive-free-text-column-naming) |
| ORM and views | [./03-orm-and-views.md](./03-orm-and-views.md) |
| Testing strategy | [./04-testing-strategy.md](./04-testing-strategy.md) |
| Cross-language DB naming | [../02-coding-guidelines/01-cross-language/07-database-naming.md](../02-coding-guidelines/01-cross-language/07-database-naming.md) |

---

*Schema design — v3.4.0 — 2026-04-19. §6 wording aligned with naming-conventions v3.5.0 Rules 10/11/12.*
