# PHP–Go Cross-Language Consistency Audit (AI Execution Prompt)

> **/goal** Enforce strict structural, architectural, and naming parity between PHP (companion plugins) and Go (backend services) across database schemas, enums, identifier casing, and API responses.
> **/learn** Master the cross-language architectural mapping: PascalCase SQLite tables/columns, mirrored enum contracts (`StatusType` vs `status.Variant`), identical casing conventions (`Id`, `Url`, `Md5`), and documented framework exemptions.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Verify database schema parity: PascalCase table names, column names, and index names across PHP and Go.
- [ ] `/learn` Enforce enum pattern parity: PascalCase cases, label values, and symmetric comparison helpers.
- [ ] `/learn` Audit identifier casing: camelCase variables/methods, PascalCase structs/classes, and camelCase log context keys.
- [ ] `/learn` Ensure API response contracts match across PHP (`ResponseKeyType`) and Go (`responsekey.Variant`).
- [ ] `/goal` Verify compliance using `python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only`.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16
**Status:** Active
**AI Confidence:** Production-Ready
**Ambiguity:** None

---

## Purpose

Documents the cross-language alignment between PHP (WordPress plugin) and Go (backend services) for naming conventions, enum patterns, database schemas, and API contracts.

---

## 1. Database Naming — PascalCase

| Aspect | PHP | Go | Status |
|--------|-----|-----|--------|
| Table names | `TableType::Transactions = 'Transactions'` | `CREATE TABLE Transactions` | ✅ Aligned |
| Column names | `'PluginSlug'`, `'CreatedAt'` | `db:"PluginSlug"` | ✅ Aligned |
| Index names | `IdxTransactions_CreatedAt` | `IdxTransactions_CreatedAt` | ✅ Aligned |
| Abbreviations | `Id`, `Url`, `Md5` | `Id`, `Url`, `Md5` | ✅ Aligned |
| WP core tables | snake_case (exempt) | N/A | ✅ Exempt |

**Reference:** [database-naming.md](../01-cross-language/07-database-naming.md)

---

## 2. Enum Patterns

| Aspect | PHP | Go | Status |
|--------|-----|-----|--------|
| Type suffix | `StatusType`, `ActionType` | `status.Variant`, `action.Variant` | ✅ Aligned (different idiom, same concept) |
| Case naming | PascalCase (`Success`, `Failed`) | PascalCase constants (`Success`, `Failed`) | ✅ Aligned |
| Label values | PascalCase strings | PascalCase `variantLabels` | ✅ Aligned |
| Comparison | `isEqual()`, `isOtherThan()`, `isAnyOf()` | `Is{Value}()`, `IsOther()`, `IsAnyOf()` | ✅ Aligned |
| Parsing | `tryFrom()` | `Parse()` (case-insensitive) | ✅ Aligned |
| JSON serialization | `->value` (string-backed) | `MarshalJSON()` / `UnmarshalJSON()` | ✅ Aligned |
| Zero value | N/A (PHP enums have no zero) | `Invalid = iota` | ✅ By design |
| Protocol-driven exemptions | N/A | Preserve functional values (`application/json`) | ✅ Documented |

**Reference:** [PHP enums.md](./02-enums.md), [Go 02-required-methods.md](../03-golang/01-enum-specification/03-required-methods.md)

---

## 3. Identifier Casing

| Identifier Type | PHP | Go | Status |
|----------------|-----|-----|--------|
| Class/struct names | PascalCase | PascalCase | ✅ Aligned |
| Method names | camelCase | PascalCase (exported) | ✅ Language idiom |
| Variables | camelCase | camelCase | ✅ Aligned |
| Log context keys | camelCase (`'postId'`) | camelCase (struct fields) | ✅ Aligned |
| DB column keys | PascalCase (`'PluginSlug'`) | PascalCase (`db:"PluginSlug"`) | ✅ Aligned |
| JSON struct tags | N/A | PascalCase (default, redundant tags removed) | ✅ Aligned |
| Abbreviations | `Id`, `Url`, `Md5` | `Id`, `Url`, `Md5` | ✅ Aligned |

---

## 4. API Response Keys

| Aspect | PHP | Go | Status |
|--------|-----|-----|--------|
| Key source | `ResponseKeyType` enum | `responsekey.Variant` | ✅ Aligned |
| Key casing | camelCase values (`'pluginSlug'`, `'isUpdate'`) | camelCase (JSON output) | ✅ Aligned |
| Envelope keys | `Success`, `Error`, `Results` | `Success`, `Error`, `Results` | ✅ Aligned |

**Reference:** [response-key-type-inventory.md](./08-response-key-type-inventory.md)

---

## 5. HTTP Status Codes

| Aspect | PHP | Go | Status |
|--------|-----|-----|--------|
| Type | `HttpStatusType: int` | `HttpStatusType` (int, exempt from byte pattern) | ✅ Aligned |
| Helpers | `isSuccess()`, `isRetryable()`, `isRedirect()` | `IsSuccess()`, `IsRetryable()`, `IsRedirect()` | ✅ Aligned |
| Retryable codes | 408, 429, 502, 503, 504 | 408, 429, 502, 503, 504 | ✅ Aligned |

---

## 6. Migrations Completed

| Phase | Scope | Status |
|-------|-------|--------|
| Phase 1 | Specs & standards | ✅ Complete |
| Phase 2A | Go SplitDB (3 tables) | ✅ Complete |
| Phase 2B | Go E2E Service (4 tables) | ✅ Complete |
| Phase 3 | PHP Plugin SQLite v13 (12 tables) | ✅ Complete |
| Phase 4 | PHP Root DB (5 tables + backward compat) | ✅ Complete |
| Phase 5 | Validation sweep | ✅ Complete |
| Batch G | camelCase log context keys (8 files) | ✅ Complete |

---

## 7. Known Exemptions

| Exemption | Reason |
|-----------|--------|
| WordPress core tables (`wp_posts`, `wp_options`) | Managed by WordPress core |
| `schema_version` table | Internal migration tracking |
| PHP `HookType`, `CapabilityType`, `NonceType` backed values | WordPress API requires snake_case |
| `WpErrorCodeType::RestForbidden`, `RestDisabled` | WordPress REST API convention |
| Go protocol-driven enums (`content_type`, `endpoint`, `header`) | Preserve functional values |
| SQLite `_snapshot_meta` keys | Per-snapshot internal persistence |
| `wp_options` setting keys | WordPress persistence layer |
| V1–V12 migration DDL | Historical, immutable |

---

## Cross-References

- [Database Naming Convention](../01-cross-language/07-database-naming.md)
- [PHP Enum Specification](./02-enums.md)
- [PHP Naming Conventions](../../01-spec-authoring-guide/03-naming-conventions.md)
- [Go Enum Specification](../03-golang/01-enum-specification/readme.md)
- [Go Required Methods](../03-golang/01-enum-specification/03-required-methods.md)

*Cross-language consistency audit v1.0.0 — 2026-02-23*

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PHP-009: PHP-Go Cross-Language Architecture Consistency

**Given** Polyglot implementations spanning PHP WordPress companion plugins and Go backend services.
**When** Cross-language architectures are audited for naming, enum patterns, and schema contracts.
**Then** All tables, columns, enums, response keys, and HTTP status handling remain 100% aligned with zero architectural drift and exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.
