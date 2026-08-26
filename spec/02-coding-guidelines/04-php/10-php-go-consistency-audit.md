# PHP–Go Cross-Language Consistency Audit

> **Version:** 1.0.0  
> **Updated:** 2026-02-23  
> **Status:** All phases complete

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

**Reference:** [PHP enums.md](./01-enums.md), [Go 02-required-methods.md](../03-golang/01-enum-specification/02-required-methods.md)

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

**Reference:** [response-key-type-inventory.md](./09-response-key-type-inventory.md)

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
- [PHP Enum Specification](./01-enums.md)
- [PHP Naming Conventions](./03-naming-conventions.md)
- [Go Enum Specification](../03-golang/01-enum-specification/00-overview.md)
- [Go Required Methods](../03-golang/01-enum-specification/02-required-methods.md)

---

*Cross-language consistency audit v1.0.0 — 2026-02-23*
