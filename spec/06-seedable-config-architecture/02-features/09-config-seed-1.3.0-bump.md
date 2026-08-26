# 09 — `config.seed.json` Bump 1.2.x → 1.3.0 (FU-12)

**Spec:** `06-seedable-config-architecture` (extension)
**Version:** 1.0.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** FU-12 from `spec/19-main-worker-service/14-rbac-and-status-seed.md` §1.
**Authority:** Canonical changelog + delivery instructions for the SemVer bump that introduces the `Tables` block. On any wording conflict with `01-fundamentals.md` §SeedWithVersionCheck, **`01-…` wins**.

---

## 1. Scope

`config.seed.json` moves from `1.2.x` → **`1.3.0`** (MINOR bump per SemVer — adds capabilities, no breaking removals).

Two additive changes:

1. New top-level `Tables` block (per `07-reference-table-seeding.md` + `08-config-schema-tables-block.md`).
2. New `Categories.MainWorker` category carrying the 27 tunables from `spec/19-main-worker-service/15-tunable-constants.md`.

No existing keys are renamed or removed.

---

## 2. Delta summary

| Change | Direction | Source | Risk |
|--------|-----------|--------|------|
| `Version`: `1.2.x` → `1.3.0` | bump | this file | none — monotone |
| `Categories.MainWorker.*` | add (27 keys) | `spec/19/15-tunable-constants.md` §3 | none — new category |
| `Tables.Role` | add (3 rows) | `spec/19/14-rbac-and-status-seed.md` §3 | none — new table seed |
| `Tables.EnumPage` | add (9 rows) | same | none |
| `Tables.RolePageAccess` | add (19 rows, uses `@-ref`) | same | needs FU-13 resolver |
| `Tables.WorkerNodeStatus` | add (4 rows) | same | none |
| `Tables.AuthMechanism` | add (4 rows) | same | none |
| `$schema` reference | add | `08-config-schema-tables-block.md` | none |

---

## 3. Top-level shape after bump

```jsonc
{
  "$schema":   "./config.schema.json",
  "Version":   "1.3.0",
  "Categories": {
    "Rag":        { /* unchanged from 1.2.x */ },
    "Validation": { /* unchanged */ },
    "Update":     { /* unchanged */ },
    "MainWorker": { /* new — see §4 */ }
  },
  "Tables": {
    "Role":             { /* see spec/19/14 §3 */ },
    "EnumPage":         { /* see spec/19/14 §3 */ },
    "RolePageAccess":   { /* see spec/19/14 §3 */ },
    "WorkerNodeStatus": { /* see spec/19/14 §3 */ },
    "AuthMechanism":    { /* see spec/19/14 §3 */ }
  }
}
```

Concrete row sets are not duplicated here — `spec/19/14-rbac-and-status-seed.md` §3 is the single source of truth.

---

## 4. New `Categories.MainWorker` block

Verbatim from `spec/19-main-worker-service/15-tunable-constants.md` §3:

```jsonc
"MainWorker": {
  "Routing.MaxAttempts":            { "Value": 3,        "Description": "Per-call retry cap." },
  "Routing.BackoffSeconds":         { "Value": [2,8,30], "Description": "Per-attempt sleep." },
  "Routing.HeartbeatWindowSeconds": { "Value": 60,       "Description": "Eligibility freshness window." },
  "Routing.MaxCompaniesPerWorker":  { "Value": 0,        "Description": "0 = unlimited." },
  "Heartbeat.IntervalSeconds":      { "Value": 30,       "Description": "Worker → Main ping cadence." },
  "Heartbeat.MissThreshold":        { "Value": 3,        "Description": "Misses before quarantine." },
  "Idempotency.KeyTtlSeconds":      { "Value": 86400,    "Description": "Replay-cache window (24h)." },
  "Auth.WorkerJwtTtlSeconds":       { "Value": 900,      "Description": "Worker JWT lifetime (15m)." },
  "Auth.ClockSkewToleranceSeconds": { "Value": 60,       "Description": "Allowed clock drift on `exp`." }
  /* … balance of the 27 keys per spec/19/15 §3 … */
}
```

Implementer note: the 18 omitted keys live in `spec/19/15` §3; copy them verbatim. They are not re-listed here to keep this file the **delta** rather than a duplicate.

---

## 5. Migration from `1.2.x`

Per `01-fundamentals.md` §SeedWithVersionCheck, the boot-time sequence on first 1.3.0 run is:

```
1. Read installed seed.Version from `config_meta`.
2. If installed < "1.3.0":
     a. Apply NEW Categories keys (additive merge — never overwrite operator-set values).
     b. Apply NEW Tables blocks via 07-reference-table-seeding.md §6 algorithm.
     c. Resolve @-refs via 10-at-ref-resolver.md.
     d. Update config_meta.SeedVersion → "1.3.0".
     e. Insert TableSeedChangelog rows (one per table).
3. If installed >= "1.3.0": skip — already migrated.
```

Per Code Red error-handling: any failure aborts the boot, logs the offending key/row/path, and is surfaced via the registered MWS error code (`MWS-21002 SplitDbTierMissing` for missing-row failures during `@-ref` resolution).

---

## 6. Verification (post-merge)

Mirrors `spec/19/14-rbac-and-status-seed.md` §4 — must all pass:

```sql
SELECT SeedVersion FROM config_meta;            -- expect '1.3.0'
SELECT COUNT(*) FROM Role;                      -- expect >= 3
SELECT COUNT(*) FROM EnumPage;                  -- expect >= 9
SELECT COUNT(*) FROM RolePageAccess;            -- expect >= 19
SELECT COUNT(*) FROM WorkerNodeStatus;          -- expect >= 4
SELECT COUNT(*) FROM AuthMechanism;             -- expect >= 4
SELECT COUNT(*) FROM TableSeedMeta;             -- expect >= 5
SELECT COUNT(*) FROM TableSeedChangelog
  WHERE ToVersion = '1.3.0';                    -- expect >= 5 (one per new Tables entry)
```

---

## 7. Rollback policy

Seed downgrade is **not supported** (per `01-fundamentals.md`). If a deploy of 1.3.0 must be reverted:

1. Roll back the application binary to a 1.2.x-aware build.
2. Leave the seeded rows in place. The 1.2.x build never reads `Tables.*` and ignores them safely (closed-additive-properties contract).
3. Operators MAY manually `DELETE FROM Role/EnumPage/...` if a clean state is required, but this is operational tooling and out of scope for the seeder.

---

## 8. Cross-references

- `spec/06-seedable-config-architecture/01-fundamentals.md` — `SeedWithVersionCheck` algorithm.
- `spec/06-seedable-config-architecture/02-features/07-reference-table-seeding.md` — `Tables` block mechanism.
- `spec/06-seedable-config-architecture/02-features/08-config-schema-tables-block.md` — JSON-Schema fragment used at lint time (FU-11).
- `spec/06-seedable-config-architecture/02-features/10-at-ref-resolver.md` — `@-ref` resolver (FU-13).
- `spec/19-main-worker-service/14-rbac-and-status-seed.md` — concrete row sets.
- `spec/19-main-worker-service/15-tunable-constants.md` §3 — concrete `MainWorker` category values.

---

## 9. Open Questions (logged, non-blocking)

- **OQ-09-1** Should the bump also gain a `Version` history field for forensic traceability beyond `config_meta`? Inferred: no — `TableSeedChangelog.FromVersion` already captures it.

---

*config.seed.json bump 1.3.0 v1.0.0 — 2026-05-04 (FU-12 closed)*
