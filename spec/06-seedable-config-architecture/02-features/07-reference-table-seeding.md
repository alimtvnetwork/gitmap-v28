# 07 — Reference-Table Seeding (Top-Level `Tables` Block)

**Spec:** `06-seedable-config-architecture` (new feature)
**Version:** 1.0.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** spec/19 audit findings F-B-09, F-B-10, F-X-06 (top-10 fix #6). Unblocks AC-5.
**Authority:** Extends the `config.seed.json` schema defined in `spec/06-seedable-config-architecture/01-fundamentals.md` with a sibling top-level `Tables` block for seeding small reference / enum / join tables. Existing `Categories` block is unchanged.

---

## 1. Why this feature exists

`config.seed.json` (per `01-fundamentals.md`) seeds **key/value settings** stored in the `config` table. Several specs — including `spec/19-main-worker-service/07-role-based-dashboards.md` — also need to seed **reference tables** (e.g. `Role`, `EnumPage`, `RolePageAccess`, `WorkerNodeStatus`) at install time and merge new rows on SemVer bumps.

Before this feature, those specs assumed an unspecified mechanism. The audit (F-B-09/10) flagged this as a BLOCKER for AC-5 (RBAC). This file defines the mechanism.

---

## 2. Schema extension

`config.seed.json` gains a new sibling to `Categories`:

```jsonc
{
  "$schema": "./config.schema.json",
  "Version": "1.3.0",
  "Categories": { /* unchanged, per 01-fundamentals */ },
  "Tables": {
    "<TableName>": {
      "AddedIn": "1.3.0",                    // SemVer when this table seed first appeared
      "Version": "1.3.0",                    // SemVer of THIS table's row set
      "PrimaryKey": ["<ColumnName>"],        // logical key for upsert (NOT the AUTOINCREMENT PK)
      "MergeStrategy": "UpsertByLogicalKey", // see §3
      "Description": "What this table holds and why it's seeded",
      "Rows": [
        { "<ColumnName>": "<Value>", "Description": "..." }
      ]
    }
  }
}
```

### 2.1 Field-by-field

| Field | Type | Nullable | Notes |
|---|---|---|---|
| `AddedIn` | SemVer | No | First SemVer that introduced this table seed. |
| `Version` | SemVer | No | Bump when row set changes. Skipped if equal/lower than DB-recorded version. |
| `PrimaryKey` | string[] | No | Columns that uniquely identify a row for merge. Length ≥ 1. |
| `MergeStrategy` | enum | No | `UpsertByLogicalKey` \| `ReplaceAll` \| `AppendOnly`. v1.0 only `UpsertByLogicalKey` and `AppendOnly` implemented. |
| `Description` | TEXT | No | Free text, copied to `config_changelog`. |
| `Rows` | object[] | No | Each row's keys must match table columns. |

---

## 3. Merge strategies

| Strategy | Behavior |
|---|---|
| `UpsertByLogicalKey` | For each row: SELECT by `PrimaryKey`. If exists, UPDATE non-key columns (skip `*Id`). Else INSERT. |
| `AppendOnly` | INSERT if row with same `PrimaryKey` does not exist. Never UPDATE. Used for audit-style enums. |
| `ReplaceAll` *(reserved)* | Truncate table then INSERT all rows. **Not implemented in v1.0** — destructive. |

`AUTOINCREMENT` PKs (per Code Red Schema rule: `{TableName}Id INTEGER PRIMARY KEY AUTOINCREMENT`) are **never** specified in `Rows`. The seeder lets the engine assign them.

---

## 4. Worked example — `EnumPage` seed

```jsonc
"Tables": {
  "EnumPage": {
    "AddedIn": "1.3.0",
    "Version": "1.3.0",
    "PrimaryKey": ["EnumPageCode"],
    "MergeStrategy": "UpsertByLogicalKey",
    "Description": "Capability catalog used by RolePageAccess.",
    "Rows": [
      { "EnumPageCode": "PowerAdminPage",     "EnumPageLabel": "Power Admin",       "Description": null },
      { "EnumPageCode": "AdminPage",          "EnumPageLabel": "Admin",             "Description": null },
      { "EnumPageCode": "BillingPage",        "EnumPageLabel": "Billing",           "Description": null }
    ]
  }
}
```

Row count is the source of truth — adding a 4th row + bumping `Version` triggers an additive merge on next boot.

---

## 5. Bookkeeping tables

The seeder records progress in two tables (created if missing):

```sql
CREATE TABLE TableSeedMeta (
    TableSeedMetaId  INTEGER PRIMARY KEY AUTOINCREMENT,
    TableName        TEXT NOT NULL UNIQUE,
    SeedVersion      TEXT NOT NULL,
    LastSeededAtUtc  TEXT NOT NULL,
    Description      TEXT NULL
);

CREATE TABLE TableSeedChangelog (
    TableSeedChangelogId  INTEGER PRIMARY KEY AUTOINCREMENT,
    TableName             TEXT NOT NULL,
    FromVersion           TEXT NULL,
    ToVersion             TEXT NOT NULL,
    RowsInserted          INTEGER NOT NULL,
    RowsUpdated           INTEGER NOT NULL,
    AppliedAtUtc          TEXT NOT NULL,
    Notes                 TEXT NULL,
    Comments              TEXT NULL,
    Description           TEXT NULL
);
```

Compliant with Code Red Schema Rules 10/11/12: entity table has `Description`; transactional has `Notes`+`Comments`+`Description`; all NULL-able with no DEFAULT.

---

## 6. Boot algorithm (mirrors `01-fundamentals.md` §SeedWithVersionCheck)

```
for each (TableName, TableBlock) in seed.Tables:
    metaRow = TableSeedMeta.get(TableName)              // may be null
    if metaRow exists AND semver(TableBlock.Version) <= semver(metaRow.SeedVersion):
        continue                                          // skip; up-to-date
    apply(TableBlock.Rows, TableBlock.MergeStrategy, TableBlock.PrimaryKey)
    upsert TableSeedMeta(TableName, TableBlock.Version, now)
    insert TableSeedChangelog(...)
```

Per spec/06 conventions: failure = error, never silent skip. Per Code Red error-handling: log `TableName + offending row index + cause path`.

---

## 7. JSON-Schema additions

`config.schema.json` MUST gain (follow-up FU-11):

```jsonc
"Tables": {
  "type": "object",
  "additionalProperties": {
    "type": "object",
    "required": ["AddedIn", "Version", "PrimaryKey", "MergeStrategy", "Rows"],
    "properties": {
      "AddedIn":       { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
      "Version":       { "type": "string", "pattern": "^\\d+\\.\\d+\\.\\d+$" },
      "PrimaryKey":    { "type": "array", "items": { "type": "string" }, "minItems": 1 },
      "MergeStrategy": { "enum": ["UpsertByLogicalKey", "AppendOnly"] },
      "Rows":          { "type": "array", "items": { "type": "object" } },
      "Description":   { "type": ["string", "null"] }
    }
  }
}
```

---

## 8. Spec/19 binding — Main-tier reference tables to seed

Per `spec/19-main-worker-service/07-role-based-dashboards.md` and `10-worker-bootstrap-protocol.md`, the following Main-tier tables MUST be seeded via this mechanism. Concrete row sets live in `spec/19-main-worker-service/14-rbac-and-status-seed.md`.

| Table | Tier (per `19/11`) | Strategy | Purpose |
|---|---|---|---|
| `Role` | Settings | UpsertByLogicalKey | Roles: PowerAdmin, AdminUser, Member |
| `EnumPage` | Settings | UpsertByLogicalKey | Capability catalog |
| `RolePageAccess` | Settings | UpsertByLogicalKey | Role↔Page grants |
| `WorkerNodeStatus` | Root | UpsertByLogicalKey | Registering / Active / Quarantined / Retired |
| `AuthMechanism` | Settings | UpsertByLogicalKey | Session / Jwt / OAuth / None |

---

## 9. Open Questions (logged, non-blocking)

- **OQ-07-1** Should `AppendOnly` track which rows were skipped for audit? Inferred: no — bookkeeping stays in `TableSeedChangelog` aggregate counts; row-level audit is logger's job.
- **OQ-07-2** Cross-table FK ordering (e.g. `RolePageAccess` references `Role` + `EnumPage`)? Inferred: seed in declaration order; spec/19/14 declares parents first. Seeder MUST NOT reorder.

---

## 10. Cross-references

- `spec/06-seedable-config-architecture/01-fundamentals.md` — sibling `Categories` block; merge-on-version-bump pattern.
- `spec/06-seedable-config-architecture/02-features/05-validation-data-seeding.md` — earlier validation-only seeding (different scope).
- `spec/19-main-worker-service/07-role-based-dashboards.md` — RBAC schema that consumes this feature.
- `spec/19-main-worker-service/10-worker-bootstrap-protocol.md` §8 — `WorkerNodeStatus` rows.
- `spec/19-main-worker-service/14-rbac-and-status-seed.md` — concrete row-set binding.

---

*Reference-table seeding v1.0.0 — 2026-05-04*
