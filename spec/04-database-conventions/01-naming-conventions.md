# Database Naming Conventions

**Version:** 3.5.0  
**Updated:** 2026-04-19

---

## Overview

All database objects use **PascalCase**. This document summarizes the rules and references the full cross-language database naming spec.

> **Full specification:** [../02-coding-guidelines/01-cross-language/07-database-naming.md](../02-coding-guidelines/01-cross-language/07-database-naming.md)

---

## Summary of Rules

| Object | Convention | Example |
|--------|-----------|---------|
| Table names | PascalCase, **singular** | `AgentSite`, `Transaction` |
| Column names | PascalCase | `PluginSlug`, `CreatedAt` |
| Primary key | `{TableName}Id` | `TransactionId`, `AgentSiteId` |
| Foreign key column | Same name as referenced PK | `AgentSiteId` references `AgentSite.AgentSiteId` |
| Boolean columns | `Is` or `Has` prefix, **positive only** | `IsActive`, `HasLicense` |
| Descriptive text on entity/reference tables | `Description TEXT NULL` | `AgentSite.Description`, `Currency.Description` |
| Transactional free-text columns | `Notes TEXT NULL`, `Comments TEXT NULL` | `Invoice.Notes`, `Transaction.Comments` |
| Index names | `Idx{Table}_{Column}` | `IdxTransactions_CreatedAt` |
| View names | PascalCase with `Vw` prefix | `VwTransactionDetail`, `VwActiveAgentSite` |
| Abbreviations | First letter only capitalized | `Id`, `Url`, `Api` — never `ID`, `URL`, `API` |

---

## Boolean Column Rules

Boolean columns follow the same principles as the [cross-language boolean conventions](../02-coding-guidelines/01-cross-language/02-boolean-principles/00-overview.md) and [no-negatives rule](../02-coding-guidelines/01-cross-language/12-no-negatives.md).

### Rule 1: Always Use `Is` or `Has` Prefix

Every boolean column MUST start with `Is` or `Has`:

```sql
-- ✅ CORRECT
IsActive      BOOLEAN NOT NULL DEFAULT 1
IsVerified    BOOLEAN NOT NULL DEFAULT 0
IsPublished   BOOLEAN NOT NULL DEFAULT 0
HasLicense    BOOLEAN NOT NULL DEFAULT 0
HasChildren   BOOLEAN NOT NULL DEFAULT 0
IsEnabled     BOOLEAN NOT NULL DEFAULT 1

-- ❌ WRONG — no prefix
Active        BOOLEAN
Verified      BOOLEAN
Published     BOOLEAN
Licensed      BOOLEAN
```

### Rule 2: Never Use Double-Negative or "Not" Boolean Column Names

Boolean columns MUST express a **single, self-contained state** — never an explicit negation that begins with `Not` or `No`. The forbidden pattern is the **double negative** (e.g., `WHERE IsNotActive = 0` reads "is not not active").

> **Important clarification (v3.3.0):** Names like `IsDisabled`, `IsInvalid`, `IsUnverified`, and `IsUnpublished` are **NOT considered negatives** here. They are legitimate positive states of a domain concept (a row genuinely *is* in the disabled / invalid / unverified / unpublished state). The forbidden form is the **explicit `Not`/`No` prefix**, which produces double negatives in WHERE clauses.

```sql
-- ❌ FORBIDDEN — explicit Not/No prefix produces double negatives
IsNotActive          BOOLEAN    -- WHERE IsNotActive = 0  → "is not not active"
IsNotPublished       BOOLEAN    -- explicit negation
IsNotLocked          BOOLEAN    -- explicit negation
HasNoLicense         BOOLEAN    -- WHERE HasNoLicense = 0 → "has no no license"
HasNoAccess          BOOLEAN    -- use HasAccess, or IsUnauthorized if domain-specific
HasNoChildren        BOOLEAN    -- use HasChildren, or IsSingle/IsLeaf if domain-specific
DoesNotExist         BOOLEAN    -- explicit negation

-- ✅ ALLOWED — "Approved Inverse of Positive": legitimate domain states
-- (single negative root, no Not/No prefix, no double-negative risk)
IsDisabled           BOOLEAN    -- the row IS in the disabled state
IsInvalid            BOOLEAN    -- the row IS in the invalid state
IsIncomplete         BOOLEAN    -- the row IS in the incomplete state
IsUnavailable        BOOLEAN    -- the row IS in the unavailable state
IsUnread             BOOLEAN    -- the row IS unread (default 0)
IsHidden             BOOLEAN    -- the row IS in the hidden state
IsBroken             BOOLEAN    -- the row IS in the broken state
IsLocked             BOOLEAN    -- the row IS in the locked state
IsUnpublished        BOOLEAN    -- the row IS in the unpublished state (default 0)
IsUnverified         BOOLEAN    -- the row IS in the unverified state
HasInvalidLicense    BOOLEAN    -- a license exists, and it IS invalid

-- ✅ ALSO CORRECT — the natural positive form
IsActive             BOOLEAN
IsValid              BOOLEAN
IsVerified           BOOLEAN
IsPublished          BOOLEAN
HasLicense           BOOLEAN
```

**Decision rule for AI implementers:**

1. If the candidate name starts with `Not` or `No` → **rename it** (use the opposite, e.g., `IsNotActive` → `IsActive`).
2. If the candidate name uses an `Un-`, `In-`, `Dis-`, or domain-specific negative *root* (e.g., `Disabled`, `Invalid`, `Unverified`, `Unpublished`) → **allowed**, as long as it represents the canonical state the business actually stores.
3. If both the positive and the inverted form are needed in code, **store only the positive form in the database** and derive the inverse as a computed field in the application layer (see Rule 9 below).

### Rule 3: Query Readability Test

A well-named boolean column reads naturally in both true and false checks:

| Column | True Check | False Check | Reads Naturally? |
|--------|-----------|-------------|-----------------|
| `IsActive` | `WHERE IsActive = 1` → "is active" | `WHERE IsActive = 0` → "is not active" | ✅ Yes |
| `HasLicense` | `WHERE HasLicense = 1` → "has license" | `WHERE HasLicense = 0` → "has no license" | ✅ Yes |
| `IsDisabled` | `WHERE IsDisabled = 1` → "is disabled" | `WHERE IsDisabled = 0` → "is not disabled" (??) | ❌ Confusing |
| `IsNotActive` | `WHERE IsNotActive = 1` → "is not active" | `WHERE IsNotActive = 0` → "is not not active" (??) | ❌ Double negative |

### Rule 4: ORM Mapping

Boolean columns with `Is`/`Has` prefix map cleanly to code:

```go
// Go struct — boolean fields match column names
type User struct {
    UserId     int64 `db:"UserId"`
    IsActive   bool  `db:"IsActive"`
    IsVerified bool  `db:"IsVerified"`
    HasLicense bool  `db:"HasLicense"`
}

// Clean, readable business logic
if user.IsActive && user.HasLicense {
    // grant access
}
```

```php
// PHP — clean property access
if ($user->IsActive && $user->HasLicense) {
    // grant access
}
```

```typescript
// TypeScript — clean conditionals
if (user.IsActive && user.HasLicense) {
    // grant access
}
```

### Rule 5: Only `Is` and `Has` — No Other Prefixes

Only two prefixes are allowed. Do NOT use `Can`, `Should`, `Was`, `Will`, `Did`, or bare adjectives:

```sql
-- ❌ WRONG — other prefixes
CanEdit       BOOLEAN    -- use IsEditable
ShouldSync    BOOLEAN    -- use IsSyncRequired
WasProcessed  BOOLEAN    -- use IsProcessed
WillExpire    BOOLEAN    -- use IsExpiring or use ExpiresAt timestamp
DidComplete   BOOLEAN    -- use IsComplete

-- ✅ CORRECT
IsEditable       BOOLEAN NOT NULL DEFAULT 0
IsSyncRequired   BOOLEAN NOT NULL DEFAULT 0
IsProcessed      BOOLEAN NOT NULL DEFAULT 0
IsComplete       BOOLEAN NOT NULL DEFAULT 0
```

### Rule 6: Always `NOT NULL DEFAULT`

Boolean columns MUST never be nullable. A three-state boolean (`true`/`false`/`NULL`) is a logic bug waiting to happen.

```sql
-- ❌ WRONG — nullable boolean
IsActive BOOLEAN    -- NULL = unknown state, breaks WHERE IsActive = 0

-- ✅ CORRECT — always NOT NULL with explicit default
IsActive BOOLEAN NOT NULL DEFAULT 1
IsVerified BOOLEAN NOT NULL DEFAULT 0
```

### Rule 7: Prefer Timestamp Over Boolean When Applicable

If you need to know **when** something happened (not just whether it happened), use a nullable timestamp column instead of a boolean:

| Boolean | Timestamp Alternative | When to Use Timestamp |
|---------|----------------------|----------------------|
| `IsDeleted` | `DeletedAt INTEGER NULL` | Soft deletes — need to know when |
| `IsExpired` | `ExpiresAt INTEGER NULL` | Expiration tracking |
| `IsBanned`  | `BannedAt  INTEGER NULL` | Audit trail needed |
| `IsCompleted` | `CompletedAt INTEGER NULL` | Duration tracking |

```sql
-- Boolean is fine when you only need true/false
IsActive   INTEGER NOT NULL DEFAULT 1   -- SQLite boolean = INTEGER 0/1

-- Timestamp is better when "when" matters (epoch seconds, UTC)
DeletedAt  INTEGER NULL    -- NULL = not deleted, non-NULL = epoch seconds
ExpiresAt  INTEGER NULL    -- NULL = never expires, non-NULL = epoch seconds
```

> **Rule of thumb:** If the business logic ever asks "when did this happen?", use a timestamp. If it only asks "is this the case?", use a boolean.

### Rule 7.1: Epoch-INTEGER Timestamp (canonical storage format) — v2.0

**Effective v5.13.0**, all timestamp columns MUST be stored as `INTEGER` representing **Unix epoch seconds in UTC** (signed 64-bit). TEXT/ISO-8601 storage is **deprecated and forbidden** for new schemas.

| Property | Rule | Rationale |
|----------|------|-----------|
| Type | `INTEGER` (SQLite) / `BIGINT` (Postgres/MySQL) | Compact, no parsing, native arithmetic, range-query friendly |
| Unit | **Seconds** since 1970-01-01T00:00:00Z | Matches Unix `time()`, language stdlibs, JWT `iat`/`exp` |
| Range | Signed 64-bit (effectively unbounded) | Y2038-safe |
| Timezone | Always UTC; no offset stored | Single source of truth |
| NULL semantics | NULL = "never / not yet" | Same as before |
| Sub-second precision | **Not stored at column level**. If sub-second ordering matters, add a separate `SeqNo INTEGER` companion column. | Avoids float drift, keeps integer comparisons exact |

**Naming suffix unchanged:** the `*At` suffix still indicates a point-in-time column.

**Defaults (SQLite):**

```sql
CreatedAt INTEGER NOT NULL DEFAULT (unixepoch())
UpdatedAt INTEGER NOT NULL DEFAULT (unixepoch())
DeletedAt INTEGER NULL
```

**Defaults (Postgres):**

```sql
CreatedAt BIGINT NOT NULL DEFAULT EXTRACT(EPOCH FROM NOW())::BIGINT
```

**Wire format (JSON / API):** transmit as JSON `number` (integer epoch seconds). Frontend converts to `Date` using `new Date(epoch * 1000)`. Do **not** emit ISO strings on the wire for `*At` fields.

**Examples:**

```
✅ "CreatedAt": 1735689600
❌ "CreatedAt": "2025-01-01T00:00:00.000Z"   (TEXT format — deprecated)
❌ "CreatedAt": 1735689600000                 (milliseconds — wrong unit)
```

**Migration note (v5.12.0 → v5.13.0):** existing TEXT timestamp columns must be migrated via `ALTER TABLE ... ADD COLUMN <Name>Epoch INTEGER` + backfill `unixepoch(<Name>)` + drop old column in a follow-up migration. Track in `98-changelog.md`.

> **Rule of thumb:** Every `*At` column on the wire, in the DB, and in logs is **epoch seconds (INTEGER, UTC)**. No exceptions, no per-table overrides.

### Rule 8: Negative-to-Positive Conversion + Approved Inverses

This is the canonical reference table for boolean column naming. Every candidate name lands in **one of three buckets**:

| Bucket | Meaning | Example |
|--------|---------|---------|
| ❌ **Forbidden** | Uses `Not`/`No` prefix → produces double negatives in WHERE clauses | `IsNotActive`, `HasNoLicense` |
| ✅ **Positive (preferred)** | Natural positive form of the concept | `IsActive`, `HasLicense` |
| 🟦 **Approved Inverse of Positive** | Single-root negative (`Un-`, `In-`, `Dis-`, or domain word) that names a legitimate stored state | `IsDisabled`, `IsLocked`, `IsUnpublished` |

> **Rule of thumb:** Either the **positive** column or its **approved inverse** may be stored — but never both, and never the forbidden form. If both views are needed in code, store one and derive the other via Rule 9.

| ❌ Negative Name (Forbidden) | ✅ Positive Replacement | 🟦 Approved Inverse of Positive | Notes |
|---|---|---|---|
| `IsNotActive` | `IsActive` | — | Flip the default; `Not` prefix forbidden |
| `IsNotEnabled` | `IsEnabled` | `IsDisabled` | `Not` prefix forbidden; `IsDisabled` stores the legitimate inverse state |
| `IsNotValid` | `IsValid` | `IsInvalid` | `Not` prefix forbidden; `IsInvalid` is an approved inverse |
| `IsNotComplete` | `IsComplete` | `IsIncomplete` | `Not` prefix forbidden; `IsIncomplete` is an approved inverse |
| `IsNotAvailable` | `IsAvailable` | `IsUnavailable` | `Not` prefix forbidden; `IsUnavailable` is an approved inverse |
| `IsNotRead` | `IsRead` | `IsUnread` | `Not` prefix forbidden; `IsUnread` is an approved inverse (default 0) |
| `IsNotVisible` | `IsVisible` | `IsHidden` | `Not` prefix forbidden; `IsHidden` is an approved inverse |
| `IsNotWorking` | `IsWorking` | `IsBroken` | `Not` prefix forbidden; `IsBroken` is an approved inverse |
| `IsNotLocked` | `IsUnlocked` *(reframe)* | `IsLocked` | `Not` prefix forbidden; `IsLocked` is an approved inverse |
| `IsNotPublished` | `IsPublished` | `IsUnpublished` | `Not` prefix forbidden; `IsUnpublished` is an approved inverse (default 0) |
| `IsNotVerified` | `IsVerified` | `IsUnverified` | `Not` prefix forbidden; `IsUnverified` is an approved inverse |
| `HasNoAccess` | `HasAccess` *or* `IsUnauthorized` | — | `No` prefix forbidden; pick `HasAccess` for resource-level checks, `IsUnauthorized` for session/role-level state. **No approved inverse column** — derive in code. |
| `HasNoChildren` | `HasChildren` *or* `IsSingle` / `IsLeaf` | — | `No` prefix forbidden; pick the domain term (`IsSingle` for users, `IsLeaf` for tree nodes). **No approved inverse column** — derive in code. |
| `HasNoLicense` | `HasLicense` *or* `HasInvalidLicense` | — | `No` prefix forbidden; use `HasInvalidLicense` only if the domain truly tracks license validity. **No approved inverse column** — derive in code. |

**Reading the table:**

- The **left** column shows forbidden names that MUST be renamed.
- The **middle** column is the canonical positive form — always safe to use.
- The **right (🟦) Approved Inverse of Positive** column lists names that *are also acceptable as stored columns* because they describe a legitimate domain state, not a syntactic negation. Choose **either** the positive **or** its approved inverse based on which one matches the canonical state of your business domain — never both.
- A dash (—) in the inverse column means **no approved inverse exists** for that concept; the derived inverse must come from Rule 9 codegen, not from a second column.

> **Approved inverse names retained from earlier versions:** `IsDisabled`, `IsInvalid`, `IsIncomplete`, `IsUnavailable`, `IsUnread`, `IsHidden`, `IsBroken`, `IsLocked`, `IsUnpublished`, `IsUnverified`. See Rule 2 clarification.

---

### Rule 9: Auto-Generated Inverted (Computed) Fields in Code

The database stores **only the positive (canonical) form** of a boolean. When code needs the inverse, it MUST be exposed as an **auto-generated computed/virtual field** at the language layer — never as a second column.

**Programming term:** This pattern is known as a **computed property**, **derived field**, or **virtual getter** (depending on the language). It is a read-only accessor whose value is mechanically derived from a stored field by negation.

**Naming convention for the derived field:**

| Stored DB Column | Auto-Derived Code Field | Derivation |
|------------------|-------------------------|-----------|
| `IsActive`       | `IsInactive`            | `!IsActive` |
| `IsEnabled`      | `IsDisabled`            | `!IsEnabled` |
| `IsValid`        | `IsInvalid`             | `!IsValid` |
| `IsComplete`     | `IsIncomplete`          | `!IsComplete` |
| `IsAvailable`    | `IsUnavailable`         | `!IsAvailable` |
| `IsRead`         | `IsUnread`              | `!IsRead` |
| `IsVisible`      | `IsHidden`              | `!IsVisible` |
| `IsWorking`      | `IsBroken`              | `!IsWorking` |
| `IsUnlocked`     | `IsLocked`              | `!IsUnlocked` |
| `IsVerified`     | `IsUnverified`          | `!IsVerified` |
| `IsPublished`    | `IsUnpublished`         | `!IsPublished` |
| `HasAccess`      | `IsUnauthorized`        | `!HasAccess` (domain-specific inverse — preferred over `HasNoAccess`) |
| `HasChildren`    | `IsSingle` *or* `IsLeaf` | `!HasChildren` (pick the domain term — preferred over `HasNoChildren`) |
| `HasLicense`     | `HasNoLicense`          | `!HasLicense` (code-only fallback when no domain term exists; never a column) |

> **Conflict resolution:** If the database stores the *approved inverse* form as the canonical state (e.g., `IsDisabled` because that is the legitimate domain state), then the derived code field is the positive (`IsEnabled = !IsDisabled`). The rule is symmetric: **store one, derive the other**.

#### Implementation Examples

```go
// Go — computed method on the struct
type User struct {
    UserId   int64 `db:"UserId"`
    IsActive bool  `db:"IsActive"`
}

func (u User) IsInactive() bool { return !u.IsActive }
```

```php
// PHP — readonly virtual getter
class User {
    public bool $IsActive;

    public function getIsInactive(): bool {
        return !$this->IsActive;
    }
}
```

```typescript
// TypeScript — getter on the model class
class User {
    constructor(public readonly IsActive: boolean) {}

    get IsInactive(): boolean {
        return !this.IsActive;
    }
}
```

#### Code-Generation Contract (for AI implementers and codegen tools)

When generating models from a database schema, an AI or code generator MUST:

1. **Detect** every boolean column matching `^(Is|Has)[A-Z][A-Za-z]+$`.
2. **Emit** the stored field exactly as named in the database (PascalCase, with the `db` tag).
3. **Emit a derived sibling** using the inversion table above. For arbitrary names, prepend `Not` after the prefix only if no canonical inverse exists (e.g., `HasLicense` → `HasNoLicense`); otherwise use the linguistic inverse from the table.
4. **Never persist** the derived field — it must be a getter/method/computed property only, with no `db` tag and no migration.
5. **Document** the derivation in a comment so reviewers can trace it back to this rule.

This guarantees: one source of truth in the database, ergonomic positive-and-negative readability in code, and zero chance of two columns drifting out of sync.

### Complete Example

```sql
-- linter-waive: MISSING-DESC-001 reason="Column-naming example; free-text columns covered in 02-schema-design.md §6.4"
CREATE TABLE User (
    UserId        INTEGER PRIMARY KEY AUTOINCREMENT,
    Name          TEXT NOT NULL,
    Email         TEXT NOT NULL UNIQUE,
    IsActive      BOOLEAN NOT NULL DEFAULT 1,
    IsVerified    BOOLEAN NOT NULL DEFAULT 0,
    HasLicense    BOOLEAN NOT NULL DEFAULT 0,
    IsAdmin       BOOLEAN NOT NULL DEFAULT 0,
    CreatedAt     INTEGER NOT NULL DEFAULT (unixepoch()),  -- epoch seconds, UTC (Rule 7.1)
    DeletedAt     INTEGER NULL                              -- soft delete (epoch seconds; NULL = not deleted)
);

-- Clean queries
SELECT * FROM User WHERE IsActive = 1 AND IsVerified = 1;
SELECT * FROM User WHERE HasLicense = 1 AND IsAdmin = 0;
SELECT * FROM User WHERE DeletedAt IS NULL;              -- not deleted
```

---

## Descriptive Free-Text Column Naming

Certain table categories MUST reserve standard nullable text columns so future context can be stored without a schema rewrite.

### Rule 10: Entity and Reference Tables Must Include `Description`

Entity/reference/master-data tables such as `AgentSite`, `Country`, `Currency`, `Status`, `Role`, and similar lookup-style tables MUST include:

```sql
Description TEXT NULL
```

This field is for explanation, help text, hints, future annotations, or any lightweight descriptive context that may be needed later.

```sql
-- ✅ CORRECT — entity/reference table with Description
CREATE TABLE AgentSite (
    AgentSiteId   INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName      TEXT NOT NULL,
    Description   TEXT NULL
);

-- ❌ WRONG — missing Description on a table that should have it
-- linter-waive: MISSING-DESC-001 reason="Column-naming example; free-text columns covered in 02-schema-design.md §6.4"
CREATE TABLE AgentSite (
    AgentSiteId INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName    TEXT NOT NULL
);
```

### Rule 11: Transactional Tables MUST Include `Notes` and `Comments`

Transactional / invoice / billing / payment / order-like tables MUST include these nullable free-text columns (aligns with project memory rule 11; absence triggers linter `MISSING-DESC-001` per the SQL linter spec):

```sql
Notes    TEXT NULL,
Comments TEXT NULL
```

Use them for exceptional cases, manual remarks, dispute context, billing comments, support explanations, or other non-structural information that should not force an immediate schema change.

- `Notes` → internal or operational context
- `Comments` → human-facing or discussion-oriented context

```sql
-- ✅ CORRECT — transactional table with Notes and Comments
CREATE TABLE Transaction (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    AgentSiteId   INTEGER NOT NULL,
    Amount        REAL NOT NULL,
    Notes         TEXT NULL,
    Comments      TEXT NULL,
    FOREIGN KEY (AgentSiteId) REFERENCES AgentSite(AgentSiteId)
);

-- ✅ CORRECT — billing/invoice style table
CREATE TABLE UserBill (
    UserBillId INTEGER PRIMARY KEY AUTOINCREMENT,
    UserId     INTEGER NOT NULL,
    AmountDue  REAL NOT NULL,
    Notes      TEXT NULL,
    Comments   TEXT NULL,
    FOREIGN KEY (UserId) REFERENCES User(UserId)
);
```

### Rule 12: These Columns Must Be Nullable

`Description`, `Notes`, and `Comments` MUST be nullable (`TEXT NULL`). They are optional context fields, not required business data.

> **Naming-only summary:**
> - Entity/reference/master tables → `Description`
> - Transaction / invoice / bill / payment / order tables → `Notes`, `Comments`
> - All three use PascalCase and remain nullable

---

## Primary Key Naming Pattern

The primary key MUST be named `{TableName}Id` — not just `Id`:

```sql
-- ❌ WRONG — generic Id
CREATE TABLE Transaction (
    Id INTEGER PRIMARY KEY AUTOINCREMENT
);

-- ✅ CORRECT — TableNameId
-- linter-waive: MISSING-DESC-001 reason="Column-naming example; free-text columns covered in 02-schema-design.md §6.4"
CREATE TABLE Transaction (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT
);
```

**Why:** When this column appears as a foreign key in another table, the name is self-documenting:

```sql
-- linter-waive: MISSING-DESC-001 reason="Column-naming example; free-text columns covered in 02-schema-design.md §6.4"
CREATE TABLE TransactionLog (
    TransactionLogId INTEGER PRIMARY KEY AUTOINCREMENT,
    TransactionId    INTEGER NOT NULL,  -- clearly references Transactions
    LogMessage       TEXT,
    FOREIGN KEY (TransactionId) REFERENCES Transaction(TransactionId)
);
```

---

## Foreign Key Column Naming

Foreign key columns MUST use the **exact same name** as the primary key they reference:

```sql
-- Source table
CREATE TABLE AgentSite (
    AgentSiteId  INTEGER PRIMARY KEY AUTOINCREMENT,
    SiteName     TEXT NOT NULL,
    Description  TEXT NULL
);

-- Referencing table — FK column matches PK name exactly
CREATE TABLE Transaction (
    TransactionId INTEGER PRIMARY KEY AUTOINCREMENT,
    AgentSiteId   INTEGER NOT NULL,  -- same name as AgentSite.AgentSiteId
    Amount        REAL,
    Notes         TEXT NULL,
    Comments      TEXT NULL,
    FOREIGN KEY (AgentSiteId) REFERENCES AgentSite(AgentSiteId)
);
```

---

## WordPress Exception

WordPress core tables (`wp_posts`, `wp_options`) retain their native `snake_case` naming. Only custom tables follow PascalCase.

See full details: [../02-coding-guidelines/01-cross-language/07-database-naming.md](../02-coding-guidelines/01-cross-language/07-database-naming.md)

---

## Cross-References

| Reference | Location |
|-----------|----------|
| Full naming spec | [../02-coding-guidelines/01-cross-language/07-database-naming.md](../02-coding-guidelines/01-cross-language/07-database-naming.md) |
| Key naming PascalCase | [../02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md](../02-coding-guidelines/01-cross-language/11-key-naming-pascalcase.md) |
| Boolean principles | [../02-coding-guidelines/01-cross-language/02-boolean-principles/00-overview.md](../02-coding-guidelines/01-cross-language/02-boolean-principles/00-overview.md) |
| No-negatives rule | [../02-coding-guidelines/01-cross-language/12-no-negatives.md](../02-coding-guidelines/01-cross-language/12-no-negatives.md) |
