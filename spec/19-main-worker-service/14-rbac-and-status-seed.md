# 14 — RBAC + Status Seed (Concrete Row Sets)

**Spec:** `19-main-worker-service`
**Version:** 2.0.0

> **v2.0.0 rename (Phase 1):** `EnumPage` → `AccessItem`, `RolePageAccess` → `RoleAccessItem`. New column `AccessItem.PageUrlSuffix` (route matcher). Old table names accepted as aliases for one release per `03-main-db-schema.md` §2.6.1 deprecation notice.
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** audit findings F-B-09, F-B-10, F-X-06 (top-10 fix #6). Closes AC-5.
**Mechanism:** `spec/06-seedable-config-architecture/02-features/07-reference-table-seeding.md` (`Tables` block in `config.seed.json`).
**Authority:** This file is the canonical row set for Main-tier RBAC and worker-status reference tables. Implementations MUST copy the JSON in §3 verbatim into `config.seed.json`.

---

## 1. What this seeds

| Table | Tier | Rows | Purpose |
|---|---|---:|---|
| `Role` | Settings | 3 | Top-level roles |
| `AccessItem` | Settings | 9 | Capability catalog (renamed from `EnumPage` in v2.0.0) |
| `RoleAccessItem` | Settings | 19 | Role↔AccessItem grants (renamed from `RolePageAccess`) |
| `WorkerNodeStatus` | Root | 7 | Worker lifecycle states (5 primary + 2 backup, Phase 6) |
| `AuthMechanism` | Settings | 4 | Endpoint auth toggles |

All five seeds ship together at SemVer `1.4.0` of `config.seed.json`.

---

## 2. Schema recap (PascalCase, Code Red Schema Rules 10/11/12)

```sql
CREATE TABLE Role (
    RoleId       INTEGER PRIMARY KEY AUTOINCREMENT,
    RoleCode     TEXT NOT NULL UNIQUE,
    RoleLabel    TEXT NOT NULL,
    Description  TEXT NULL
);

CREATE TABLE AccessItem (
    AccessItemId   INTEGER PRIMARY KEY AUTOINCREMENT,
    Code           TEXT NOT NULL UNIQUE,
    Label          TEXT NOT NULL,
    PageUrlSuffix  TEXT NULL,
    Description    TEXT NULL
);
CREATE INDEX IX_AccessItem_PageUrlSuffix ON AccessItem(PageUrlSuffix);

CREATE TABLE RoleAccessItem (   -- join table, exempt from Description per Rule 12 carve-out
    RoleAccessItemId INTEGER PRIMARY KEY AUTOINCREMENT,
    RoleId           INTEGER NOT NULL REFERENCES Role(RoleId),
    AccessItemId     INTEGER NOT NULL REFERENCES AccessItem(AccessItemId),
    CanRead          INTEGER NOT NULL,
    CanWrite         INTEGER NOT NULL,
    UNIQUE (RoleId, AccessItemId)
);

CREATE TABLE WorkerNodeStatus (
    WorkerNodeStatusId    INTEGER PRIMARY KEY AUTOINCREMENT,
    WorkerNodeStatusCode  TEXT NOT NULL UNIQUE,
    WorkerNodeStatusLabel TEXT NOT NULL,
    Description           TEXT NULL
);

CREATE TABLE AuthMechanism (
    AuthMechanismId    INTEGER PRIMARY KEY AUTOINCREMENT,
    AuthMechanismCode  TEXT NOT NULL UNIQUE,
    AuthMechanismLabel TEXT NOT NULL,
    Description        TEXT NULL
);
```

> `RoleAccessItem` is a join table — Description column intentionally absent (Code Red Rule 12 join-table carve-out).

---

## 3. `Tables` block (paste into `config.seed.json`)

```jsonc
"Tables": {

  "Role": {
    "AddedIn": "1.3.0",
    "Version": "1.3.0",
    "PrimaryKey": ["RoleCode"],
    "MergeStrategy": "UpsertByLogicalKey",
    "Description": "Top-level roles for Main-tier RBAC.",
    "Rows": [
      { "RoleCode": "PowerAdmin", "RoleLabel": "Power Administrator", "Description": "Full cross-tenant control. Issued only to platform operators." },
      { "RoleCode": "AdminUser",  "RoleLabel": "Company Administrator", "Description": "Manages a single company: users, billing, settings." },
      { "RoleCode": "Member",     "RoleLabel": "Member",                "Description": "Standard end-user. Read-only outside personal scope." }
    ]
  },

  "AccessItem": {
    "AddedIn": "1.4.0",
    "Version": "1.4.0",
    "PrimaryKey": ["Code"],
    "MergeStrategy": "UpsertByLogicalKey",
    "Description": "Capability catalog enforced by RoleAccessItem. PageUrlSuffix is the route matcher (suffix match against normalized request path).",
    "Rows": [
      { "Code": "PowerAdmin",      "Label": "Power Admin",       "PageUrlSuffix": "/poweradmin",        "Description": "Cross-tenant ops console." },
      { "Code": "Admin",           "Label": "Admin",             "PageUrlSuffix": "/admin",             "Description": "Per-company admin home." },
      { "Code": "Billing",         "Label": "Billing",           "PageUrlSuffix": "/billing",           "Description": "Invoices, plan, payment methods." },
      { "Code": "CompanySettings", "Label": "Company Settings",  "PageUrlSuffix": "/settings/company",  "Description": "Company profile + branding." },
      { "Code": "UserManagement",  "Label": "User Management",   "PageUrlSuffix": "/settings/users",    "Description": "Add/remove users, assign roles." },
      { "Code": "WorkerRegistry",  "Label": "Region Registry",   "PageUrlSuffix": "/regions",           "Description": "List + status of worker nodes (UI label: Regions)." },
      { "Code": "PushUpdate",      "Label": "Push Update",       "PageUrlSuffix": "/regions/update",    "Description": "Trigger worker self-updates." },
      { "Code": "AuditLog",        "Label": "Audit Log",         "PageUrlSuffix": "/audit",             "Description": "Read-only audit history." },
      { "Code": "Dashboard",       "Label": "Dashboard",         "PageUrlSuffix": "/dashboard",         "Description": "Default landing page." }
    ]
  },

  "RoleAccessItem": {
    "AddedIn": "1.4.0",
    "Version": "1.4.0",
    "PrimaryKey": ["RoleId", "AccessItemId"],
    "MergeStrategy": "UpsertByLogicalKey",
    "Description": "Default grants. Operators may add/remove via UI; seed only enforces base baseline. CanRead=1, CanWrite=1 unless noted.",
    "Rows": [
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.PowerAdmin",      "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.Admin",           "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.Billing",         "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.CompanySettings", "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.UserManagement",  "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.WorkerRegistry",  "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.PushUpdate",      "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.AuditLog",        "CanRead": 1, "CanWrite": 0 },
      { "RoleId": "@Role.PowerAdmin", "AccessItemId": "@AccessItem.Dashboard",       "CanRead": 1, "CanWrite": 1 },

      { "RoleId": "@Role.AdminUser",  "AccessItemId": "@AccessItem.Admin",           "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.AdminUser",  "AccessItemId": "@AccessItem.Billing",         "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.AdminUser",  "AccessItemId": "@AccessItem.CompanySettings", "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.AdminUser",  "AccessItemId": "@AccessItem.UserManagement",  "CanRead": 1, "CanWrite": 1 },
      { "RoleId": "@Role.AdminUser",  "AccessItemId": "@AccessItem.AuditLog",        "CanRead": 1, "CanWrite": 0 },
      { "RoleId": "@Role.AdminUser",  "AccessItemId": "@AccessItem.Dashboard",       "CanRead": 1, "CanWrite": 1 },

      { "RoleId": "@Role.Member",     "AccessItemId": "@AccessItem.Dashboard",       "CanRead": 1, "CanWrite": 0 }
    ]
  },

  "WorkerNodeStatus": {
    "AddedIn": "1.3.0",
    "Version": "1.5.0",
    "PrimaryKey": ["WorkerNodeStatusCode"],
    "MergeStrategy": "UpsertByLogicalKey",
    "Description": "Worker lifecycle states observed by Main. Codes Active/Draining/Offline/Quarantined/Retired apply to primaries (IsBackup=0); BackupAttached/BackupLagging apply to backups (IsBackup=1) per 18-backup-nodes.md §4.3.",
    "Rows": [
      { "WorkerNodeStatusCode": "Registering",    "WorkerNodeStatusLabel": "Registering",     "Description": "Initial state — handshake in progress." },
      { "WorkerNodeStatusCode": "Provisioning",   "WorkerNodeStatusLabel": "Provisioning",    "Description": "Backup registered; awaiting first diff. Phase 6." },
      { "WorkerNodeStatusCode": "Active",         "WorkerNodeStatusLabel": "Active",          "Description": "Primary heartbeating; eligible for routing. Never assigned to IsBackup=1 rows." },
      { "WorkerNodeStatusCode": "Quarantined",    "WorkerNodeStatusLabel": "Quarantined",     "Description": "Primary missed >=3 heartbeats; new tenants not routed here." },
      { "WorkerNodeStatusCode": "Retired",        "WorkerNodeStatusLabel": "Retired",         "Description": "Permanently removed; preserved for audit." },
      { "WorkerNodeStatusCode": "BackupAttached", "WorkerNodeStatusLabel": "Backup Attached", "Description": "Backup healthy; last diff applied within MainWorker.Backup.LagWarningSeconds. Phase 6." },
      { "WorkerNodeStatusCode": "BackupLagging",  "WorkerNodeStatusLabel": "Backup Lagging",  "Description": "Backup attached but diff lag exceeds MainWorker.Backup.LagWarningSeconds. No auto-failover. Phase 6." }
    ]
  },

  "AuthMechanism": {
    "AddedIn": "1.3.0",
    "Version": "1.3.0",
    "PrimaryKey": ["AuthMechanismCode"],
    "MergeStrategy": "UpsertByLogicalKey",
    "Description": "Per-endpoint auth toggle catalog (used by EndpointAuthSetting).",
    "Rows": [
      { "AuthMechanismCode": "Session", "AuthMechanismLabel": "Main session cookie",         "Description": "HTTPOnly cookie set by Main." },
      { "AuthMechanismCode": "Jwt",     "AuthMechanismLabel": "Worker JWT (RS256)",          "Description": "Per spec/19/12-jwt-delivery-contract.md." },
      { "AuthMechanismCode": "OAuth",   "AuthMechanismLabel": "OAuth client-credentials",    "Description": "Worker↔Main service auth." },
      { "AuthMechanismCode": "None",    "AuthMechanismLabel": "Public",                      "Description": "No auth required — explicitly opt-in." }
    ]
  }

}
```

### 3.1 Logical-key references (`@Role.PowerAdmin`)

The `@<Table>.<Code>` syntax tells the seeder to substitute the AUTOINCREMENT id of the row whose unique-code column equals `Code`. Resolved at apply time. If the referenced row is missing, fail with `MWS-21002 SplitDbTierMissing` (closest fit) and abort the merge for that table.

Implementer note: this is a seeder feature (FU-13 — extend `01-fundamentals.md` `mergeSeed` to honor `@`-references). Until shipped, use a two-pass approach: seed `Role` first, then `RoleAccessItem` with literal IDs derived from a SELECT.

---

## 4. Verification (post-boot self-check)

After the seeder runs, the worker MUST verify:

```sql
SELECT COUNT(*) FROM Role;             -- expect ≥ 3
SELECT COUNT(*) FROM AccessItem;       -- expect ≥ 9
SELECT COUNT(*) FROM RoleAccessItem
  WHERE RoleId = (SELECT RoleId FROM Role WHERE RoleCode='PowerAdmin');
-- expect ≥ 9 (PowerAdmin has all access items)
SELECT COUNT(*) FROM WorkerNodeStatus; -- expect ≥ 4
SELECT COUNT(*) FROM AuthMechanism;    -- expect ≥ 4
```

Any count below expected → exit `MWS-21002 SplitDbTierMissing`.

---

## 5. Re-seeding rules (when authoring `1.4.0`)

- **Adding a new page:** add to `AccessItem.Rows`, add corresponding `RoleAccessItem` rows, bump both blocks' `Version` to `1.4.0`. Keep `AddedIn` as the original.
- **Revoking a default grant:** Per spec/06 §SeedWithVersionCheck, `UpsertByLogicalKey` does NOT delete. Operators clear via UI; seed cannot rescind. (Documented as expected behavior — defaults are floors, not ceilings.)
- **Renaming a page:** add new code, keep the old (deprecate via `Description` text), migrate in code over one release.

---

## 6. Cross-references

- `spec/06-seedable-config-architecture/02-features/07-reference-table-seeding.md` — the `Tables` block mechanism.
- `spec/19-main-worker-service/07-role-based-dashboards.md` — consumer of `Role`/`AccessItem`/`RoleAccessItem`.
- `spec/19-main-worker-service/10-worker-bootstrap-protocol.md` §8 — consumer of `WorkerNodeStatus`.
- `spec/19-main-worker-service/06-core-api-endpoints.md` §5 — consumer of `AuthMechanism`.
- `spec/19-main-worker-service/11-split-db-tier-reconciliation.md` §4 — tier placement.

---

*RBAC + status seed v2.0.0 — 2026-05-06 (Phase 1: AccessItem rename)*
