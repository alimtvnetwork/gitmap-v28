# 08 — `config.schema.json` Tables-Block Schema (FU-11)

**Spec:** `06-seedable-config-architecture` (extension)
**Version:** 1.0.0
**Created:** 2026-05-04
**Status:** Authoritative
**Resolves:** FU-11 from `spec/19-main-worker-service/14-rbac-and-status-seed.md` §3.1.
**Authority:** Canonical JSON-Schema fragment for the new top-level `Tables` block introduced by `07-reference-table-seeding.md`. Implementations MUST embed §2 verbatim into the project's `config.schema.json`. On any conflict between this file and `07-reference-table-seeding.md`, **`07-…` wins** (this file is its machine-readable mirror).

---

## 1. Why this file exists

`07-reference-table-seeding.md` §7 sketched the schema fragment but did not deliver a paste-ready, draft-2020-12-compliant block. This file delivers it, plus the diagnostics, examples, and CI guidance the implementer needs to wire it into `config.schema.json`.

---

## 2. Schema fragment (paste into `config.schema.json`)

```jsonc
{
  "$schema": "https://json-schema.org/draft/2020-12/schema",
  "$id": "https://riseup.asia/schemas/config.schema.json",
  "title": "Seedable Config",
  "type": "object",
  "required": ["Version", "Categories"],
  "properties": {
    "$schema": { "type": "string" },
    "Version": { "$ref": "#/$defs/SemVer" },
    "Categories": { "type": "object" },
    "Tables": { "$ref": "#/$defs/TablesBlock" }
  },
  "additionalProperties": false,

  "$defs": {
    "SemVer": {
      "type": "string",
      "pattern": "^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)$",
      "description": "Strict X.Y.Z. No pre-release / build metadata in v1.0."
    },

    "AtRefString": {
      "type": "string",
      "pattern": "^@[A-Z][A-Za-z0-9]*\\.[A-Za-z0-9_]+$",
      "description": "Logical reference; resolved at apply time per spec/06/02-features/10-at-ref-resolver.md."
    },

    "RowValue": {
      "oneOf": [
        { "type": "string" },
        { "type": "number" },
        { "type": "boolean" },
        { "type": "null" },
        { "$ref": "#/$defs/AtRefString" }
      ]
    },

    "MergeStrategy": {
      "enum": ["UpsertByLogicalKey", "AppendOnly"],
      "description": "ReplaceAll is reserved (destructive) and NOT implemented in v1.0."
    },

    "TableBlock": {
      "type": "object",
      "required": ["AddedIn", "Version", "PrimaryKey", "MergeStrategy", "Rows"],
      "additionalProperties": false,
      "properties": {
        "AddedIn":       { "$ref": "#/$defs/SemVer" },
        "Version":       { "$ref": "#/$defs/SemVer" },
        "PrimaryKey":    {
          "type": "array",
          "items": { "type": "string", "pattern": "^[A-Z][A-Za-z0-9]*$" },
          "minItems": 1,
          "uniqueItems": true
        },
        "MergeStrategy": { "$ref": "#/$defs/MergeStrategy" },
        "Description":   { "type": ["string", "null"] },
        "Notes":         { "type": ["string", "null"] },
        "Comments":      { "type": ["string", "null"] },
        "Rows": {
          "type": "array",
          "minItems": 0,
          "items": {
            "type": "object",
            "additionalProperties": { "$ref": "#/$defs/RowValue" }
          }
        }
      }
    },

    "TablesBlock": {
      "type": "object",
      "propertyNames": {
        "pattern": "^[A-Z][A-Za-z0-9]*$",
        "description": "PascalCase table name."
      },
      "additionalProperties": { "$ref": "#/$defs/TableBlock" }
    }
  }
}
```

---

## 3. Validation rules surfaced by the schema

| Rule | Schema mechanism | Failure example |
|------|------------------|-----------------|
| PascalCase table names | `TablesBlock.propertyNames.pattern` | `"role": {…}` rejected (must be `"Role"`). |
| PascalCase column names in `PrimaryKey` | `TableBlock.PrimaryKey.items.pattern` | `["role_id"]` rejected. |
| Strict SemVer | `SemVer.pattern` | `"1.3"`, `"1.3.0-rc1"` rejected. |
| Closed strategy enum | `MergeStrategy.enum` | `"ReplaceAll"` rejected (reserved). |
| `@-ref` shape | `AtRefString.pattern` | `"@role.PowerAdmin"` rejected (lowercase table). |
| No stray top-level keys | `additionalProperties: false` on root | Catches typos like `"Tabels"`. |
| `PrimaryKey` non-empty + unique | `minItems: 1`, `uniqueItems: true` | `[]` and `["X","X"]` rejected. |

Semantic rules **not** expressible in JSON-Schema (enforced by the seeder, see §5):

- Every `PrimaryKey` column MUST appear as a key in every row.
- Every `@<Table>.<Code>` reference MUST resolve to a real row (per resolver spec FU-13).
- `Version` MUST be `>= AddedIn`.

---

## 4. Worked validation example

Valid — `EnumPage` block from `spec/19/14-rbac-and-status-seed.md` §3:

```jsonc
"Tables": {
  "EnumPage": {
    "AddedIn": "1.3.0",
    "Version": "1.3.0",
    "PrimaryKey": ["EnumPageCode"],
    "MergeStrategy": "UpsertByLogicalKey",
    "Description": "Capability catalog enforced by RolePageAccess.",
    "Rows": [
      { "EnumPageCode": "PowerAdminPage", "EnumPageLabel": "Power Admin", "Description": "Cross-tenant ops console." }
    ]
  }
}
```

Invalid — illustrates each guard:

```jsonc
"Tables": {
  "role":            {                              /* fails propertyNames */ },
  "Role": {
    "AddedIn":       "1.3",                         /* fails SemVer */
    "Version":       "1.3.0",
    "PrimaryKey":    ["role_code"],                 /* fails PascalCase */
    "MergeStrategy": "ReplaceAll",                  /* fails enum */
    "Rows":          [ { "RoleCode": "@role.X" } ]  /* fails AtRefString */
  }
}
```

---

## 5. Pre-apply checks the seeder MUST run (beyond JSON-Schema)

Each is a positive guard per Code Red:

```
isVersionMonotone(block)        → block.Version >= block.AddedIn
hasAllPrimaryKeyColumns(block)  → every row contains every PrimaryKey column
hasResolvableAtRefs(block, db)  → per FU-13 resolver
isStrategyImplemented(block)    → MergeStrategy in {UpsertByLogicalKey, AppendOnly}
```

Failure of any guard → abort merge for that table (NEVER silent skip), log `TableName + offending row index + offending column + cause path` per `mem://architecture/error-handling`.

---

## 6. CI integration

A shared validator runs in CI before any seed is applied to a development DB:

```
ajv validate \
  -s spec/06-seedable-config-architecture/02-features/08-config-schema-tables-block.md#section-2 \
  -d spec/19-main-worker-service/14-rbac-and-status-seed.md#section-3
```

The build script extracts §2 (this file) and §3 (`14-…md`) at lint time. A future linter `check-config-seed.py` may absorb this — tracked, not in scope here.

---

## 7. Cross-references

- `spec/06-seedable-config-architecture/01-fundamentals.md` — sibling `Categories` block.
- `spec/06-seedable-config-architecture/02-features/07-reference-table-seeding.md` — feature definition this file mechanizes.
- `spec/06-seedable-config-architecture/02-features/09-config-seed-1.3.0-bump.md` — first seed file that uses this schema (FU-12).
- `spec/06-seedable-config-architecture/02-features/10-at-ref-resolver.md` — `@<Table>.<Code>` resolver (FU-13).
- `spec/19-main-worker-service/14-rbac-and-status-seed.md` — first concrete consumer.

---

## 8. Open Questions (logged, non-blocking)

- **OQ-08-1** Add `MaxRows` per table to bound seed bloat? Inferred: defer; row counts already gated by code review.
- **OQ-08-2** Allow `Version` downgrades (rollback)? Inferred: no — `01-fundamentals.md` §SeedWithVersionCheck mandates monotone increase.

---

*Tables-block schema v1.0.0 — 2026-05-04 (FU-11 closed)*
