# Update.* Configuration Keys

**Version:** 1.0.0  
**Created:** 2026-04-20  
**Status:** Active  
**Purpose:** Register the `Update.*` and `Storage.Backend` keys consumed by the update-check subsystem ([spec/14-update/24-update-check-mechanism](../../14-update/24-update-check-mechanism/00-overview.md))

---

## Overview

The CLI update-check subsystem reads four keys from CW Config on every
invocation. These keys MUST be present in `config.seed.json` so that
first-run seeding produces a working configuration without code-side
defaults leaking in.

---

## Key Inventory

| Key | Type | Default | Purpose |
|-----|------|---------|---------|
| `Update.BackgroundUpdateCheckEnabled` | `Boolean` | `true`  | Master switch for the pre-command hook spawn |
| `Update.CheckIntervalHours`           | `Integer` | `12`    | Interval gate; pre-hook only spawns if `Now ≥ LastCheckedAt + Interval` |
| `Update.PendingUpdateWarningEnabled`  | `Boolean` | `true`  | Master switch for the post-command warning line |
| `Storage.Backend`                     | `Enum`    | `Sqlite`| `Sqlite` or `JsonFile` — selects the `UpdateStore` backend |

PascalCase applies to keys AND to enum values
(`Sqlite`, `JsonFile`).

---

## `config.seed.json` Snippet

```json
{
  "Update": {
    "BackgroundUpdateCheckEnabled": true,
    "CheckIntervalHours": 12,
    "PendingUpdateWarningEnabled": true
  },
  "Storage": {
    "Backend": "Sqlite"
  }
}
```

The `Update` and `Storage` objects are **independent top-level
sections** — they MUST NOT be nested under each other or under any
existing section.

---

## Validation Rules

| Key | Rule | On violation |
|-----|------|--------------|
| `Update.BackgroundUpdateCheckEnabled` | strict `Boolean` (no truthy strings) | Reject seed; fall back to `true` and log |
| `Update.CheckIntervalHours`           | Integer in `[1, 168]` (1 hour to 7 days) | Reject seed; fall back to `12` and log |
| `Update.PendingUpdateWarningEnabled`  | strict `Boolean` | Reject seed; fall back to `true` and log |
| `Storage.Backend`                     | `ParseStorageBackend()` strict enum | Reject seed; fall back to `Sqlite` and log |

Validation lives in the same file as other CW Config validators
(see [02-rag-validation-helpers.md](./02-rag-validation-helpers.md) for
the helper pattern). Functions follow CODE RED metrics: 8–15 lines,
positively named, max 2 operands.

---

## Independence of the Two Switches

`BackgroundUpdateCheckEnabled` and `PendingUpdateWarningEnabled` are
**independent** by design:

| Scenario | Background spawn | Trailing warning |
|----------|------------------|------------------|
| Default | ✅ | ✅ |
| User wants no network calls but wants to be reminded if a check happened earlier | ❌ | ✅ |
| User wants automatic checks but no noise in piped output | ✅ | ❌ |
| User wants the subsystem fully silent | ❌ | ❌ |

Do not collapse them into a single `UpdateChecksEnabled` flag.

---

## `Storage.Backend` Selection

`Sqlite` (default) routes `UpdateCheckerService` to the host CLI's
existing SQLite database. `JsonFile` routes it to
`~/.<CliName>/data/UpdateChecker.json` per
[09-json-fallback-store.md](../../14-update/24-update-check-mechanism/09-json-fallback-store.md).

The host CLI MAY override this at boot (e.g. a CLI with no DB at all
hard-codes `JsonFile`). The seeded value is the user-visible default.

---

## Migration

When upgrading a project that does not yet have these keys:

1. The seedable-config merge strategy (see
   [01-fundamentals.md](../01-fundamentals.md)) detects the missing
   keys and inserts them with the defaults above.
2. `changelog.md` records the seed-version bump.
3. No data migration is needed — the keys are read-only inputs to the
   update subsystem; they do not own any persisted state.

---

## Cross-References

| Reference | Description |
|-----------|-------------|
| [Update Check Mechanism — Overview](../../14-update/24-update-check-mechanism/00-overview.md) | Subsystem that consumes these keys |
| [Pre-Command Hook](../../14-update/24-update-check-mechanism/07-pre-command-hook.md) | Reads both `Update.*` boolean switches |
| [JSON Fallback Store](../../14-update/24-update-check-mechanism/09-json-fallback-store.md) | Activated when `Storage.Backend = "JsonFile"` |
| [Validation Helpers](./02-rag-validation-helpers.md) | Helper pattern these validators follow |

---

*Update.* configuration keys — v1.0.0 — 2026-04-20*
