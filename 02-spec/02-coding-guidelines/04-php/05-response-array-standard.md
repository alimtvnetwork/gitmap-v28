# PHP Response Array Standard (AI Execution Prompt)

> **/goal** Enforce structured PHP response arrays across internal services using `ResultHelper` factory methods and `ResponseKeyType` enum keys.
> **/learn** Distinguish internal `ResultHelper` arrays from HTTP `EnvelopeBuilder` envelopes, eliminating bare string keys, inline error payloads, and unstructured arrays.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Construct all service result arrays using `ResultHelper::ok()`, `ResultHelper::error()`, or `ResultHelper::errorWithCode()`.
- [ ] `/learn` Never use bare string keys (`'rows'`, `'total'`); access keys via `ResponseKeyType::*->value`.
- [ ] `/goal` Format all multi-item result arrays with one key-value pair per line for clean diffing.
- [ ] `/learn` Import `RiseupAsia\Helpers\ResultHelper` and `RiseupAsia\Enums\ResponseKeyType` at file level in all consumer files.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Last Updated:** 2026-04-16

> **Version:** 1.0.0
> **Since:** 2.2.0
> **Scope:** All internal service result arrays in the `RiseupAsia` namespace.

---

## 1. Overview

Internal PHP service methods return structured arrays to communicate success/failure, data, and error context between traits, services, and REST handlers. This document standardizes:

1. **Factory methods** via `ResultHelper` — eliminates boilerplate.
2. **Key standardization** via `ResponseKeyType` enum — eliminates magic strings.
3. **Multi-line formatting** — improves readability and diff quality.

> **Important distinction:** `ResultHelper` builds **internal service result arrays** (small, flat arrays passed between methods). `EnvelopeBuilder` builds **REST API response envelopes** (structured, multi-section JSON for HTTP responses). Do not confuse the two.

---

## 2. ResultHelper — Factory Methods

**Location:** `includes/Helpers/ResultHelper.php`
**Namespace:** `RiseupAsia\Helpers`

### 2.1 Available Methods

| Method | Returns | Use Case |
|--------|---------|----------|
| `ResultHelper::ok(array $extra = [])` | `{success: true, ...}` | Successful operation |
| `ResultHelper::failed(array $extra = [])` | `{success: false, ...}` | Bare failure (no message) |
| `ResultHelper::error(string $msg, array $extra = [])` | `{success: false, error: $msg, ...}` | Failure with message |
| `ResultHelper::errorWithCode(string $msg, string $code, array $extra = [])` | `{success: false, error: $msg, code: $code, ...}` | Failure with message and code |
| `ResultHelper::errorFromException(Throwable $e, array $extra = [])` | `{success: false, error: $e->getMessage(), ...}` | Failure from caught exception |

### 2.2 Usage Examples

```php
// ✅ REQUIRED — Simple success
return ResultHelper::ok();

// ✅ REQUIRED — Success with extra data
return ResultHelper::ok(array(
    ResponseKeyType::Rows->value => $totalRows,
));

// ✅ REQUIRED — Error with message
return ResultHelper::error('Table not found in snapshot');

// ✅ REQUIRED — Error with extra context
return ResultHelper::error(
    'Table not found in snapshot',
    array(ResponseKeyType::Rows->value => 0),
);

// ✅ REQUIRED — Error with code (typically from an error-domain enum)
return ResultHelper::errorWithCode(
    'Database unavailable',
    SnapshotErrorType::ExportBuildFailed->value,
);

// ✅ REQUIRED — Error from exception in catch block
return ResultHelper::errorFromException(
    $e,
    array(ResponseKeyType::Rows->value => 0),
);

// ✅ REQUIRED — Bare failure (no message needed)
return ResultHelper::failed();
```

### 2.3 Anti-Patterns

```php
// ❌ FORBIDDEN — Manual success/error boilerplate
return array(
    ResponseKeyType::Success->value => false,
    ResponseKeyType::Error->value   => $e->getMessage(),
    ResponseKeyType::Rows->value    => 0,
);

// ❌ FORBIDDEN — Inline success key
return array(ResponseKeyType::Success->value => true, ResponseKeyType::Rows->value => $totalRows);
```

---

## 3. ResponseKeyType — Zero Magic Strings

**Location:** `includes/Enums/ResponseKeyType.php`
**Namespace:** `RiseupAsia\Enums`

Every key in a service result array or structured response **must** use a `ResponseKeyType` case. No bare string keys are permitted in response arrays.

### 3.1 Key Categories

| Category | Cases |
|----------|-------|
| **Envelope** | `Success`, `Error`, `Message`, `Data`, `Code`, `Valid`, `Errors`, `Cached`, `Phase`, `Reason` |
| **Domain Collections** | `Total`, `Agents`, `Actions`, `Logs`, `Snapshots`, `Sql`, `Params`, `Sets`, `Plugins`, `Tables` |
| **Pagination** | `Limit`, `Offset` |
| **Domain Entities** | `Posts`, `Categories`, `Category`, `Export`, `Incrementals`, `TotalSize`, `Applied`, `Folder` |
| **File/Size** | `Rows`, `Bytes`, `Size`, `FileSize`, `Path`, `Filename`, `Checksum`, `Duration`, `Count`, `Files`, `Directory`, `Scope`, `Exported`, `Entry`, `Computed`, `Removed` |
| **Snapshot Domain** | `SnapshotId`, `Sequence`, `FolderName`, `TablesChanged`, `TotalRows`, `TotalNewRows`, `ZipSize`, `BackupId`, `ZipFailed`, `SkipAudit`, `TablesRestored` |

### 3.2 When to Add a New Case

Add a new `ResponseKeyType` case when:

- A string key appears in **two or more** response arrays across different traits/classes.
- A string key is part of a **client-facing API response** (even if used only once).

Do **not** add cases for:

- Internal-only metadata keys used within a single method (e.g., database record fields).
- WordPress core keys (`type` in `error_get_last()`, array config keys).
- Log context arrays (these are for debugging, not structured API contracts).

### 3.3 Exemptions

The following bare string keys remain permitted per project conventions:

- **WordPress API contracts:** Hook names, option names, table column names (snake_case).
- **Database record fields:** Column names read from `PDO::FETCH_ASSOC` results.
- **Log context arrays:** Keys in the second argument of `$this->log()` calls.
- **i18n translation calls:** The plugin slug in `__()` and `_e()`.

---

## 4. Multi-Line Array Formatting

All response arrays **must** use one key-value pair per line. This improves:

- **Readability** — Each field is immediately visible.
- **Diff quality** — Changes to a single field produce a one-line diff.
- **Code review** — Reviewers can scan vertically rather than parsing long horizontal lines.

### 4.1 Formatting Rules

```php
// ✅ REQUIRED — One pair per line, trailing comma
return ResultHelper::ok(array(
    ResponseKeyType::Total->value  => $query->found_posts,
    ResponseKeyType::Limit->value  => $args['posts_per_page'],
    ResponseKeyType::Offset->value => $args['offset'],
    ResponseKeyType::Posts->value   => $posts,
));

// ✅ REQUIRED — Non-ResultHelper arrays also use multi-line
return array(
    ResponseKeyType::Count->value    => $count,
    ResponseKeyType::TotalSize->value => $total_size,
    ResponseKeyType::Plugins->value  => $plugin_list,
);

// ❌ FORBIDDEN — Horizontal packing
return array(ResponseKeyType::Count->value => $count, 'total_size' => $total_size, ResponseKeyType::Plugins->value => $plugin_list);
```

### 4.2 Exception: Single-Key Arrays

Arrays with a single key-value pair **may** remain on one line:

```php

return ResultHelper::ok(array(ResponseKeyType::Rows->value => $totalRows));
```

---

## 5. Namespace Import Requirements

Every trait file that uses `ResultHelper` **must** include the file-level import:

```php
use RiseupAsia\Helpers\ResultHelper;
```

Per the namespace import integrity lesson, PHP traits do not inherit `use` imports from the classes they are mixed into. Missing imports cause fatal errors that only surface after OPcache is cleared.

---

## 6. Decision Matrix: ResultHelper vs EnvelopeBuilder

| Scenario | Use |
|----------|-----|
| Internal service method returning success/failure to a caller | `ResultHelper` |
| REST endpoint building a full JSON response for the client | `EnvelopeBuilder` |
| Admin AJAX handler returning via `wp_send_json_success/error` | `ResultHelper` for building the data array |
| Trait returning partial data to an orchestrating method | `ResultHelper` or plain array with `ResponseKeyType` keys |

---

## 7. Migration Checklist

When refactoring existing code to this standard:

1. **Replace inline `array(ResponseKeyType::Success->value => ...)` with `ResultHelper::ok/error/...`.**
2. **Replace all bare string keys** (`'total'`, `'rows'`, `'export'`, etc.) with `ResponseKeyType::X->value`.
3. **Format arrays** with one key-value pair per line.
4. **Add `use RiseupAsia\Helpers\ResultHelper;`** to the file-level imports.
5. **Update callers** that read the result to use `ResponseKeyType::X->value` for key access.
6. **Cold-activate test** — clear OPcache and verify no fatal errors.

---

## 8. Files Affected (Reference)

The following files have been refactored as reference implementations:

- `Post/Traits/PostQueryTrait.php` — Pagination with `Limit`, `Offset`, `Posts`
- `Post/Traits/CategoryTrait.php` — `Category`, `Categories` entity keys
- `Snapshot/Traits/ManagerTableRestoreTrait.php` — `ResultHelper::error/errorFromException/ok`
- `Snapshot/Traits/ExporterPublicApiTrait.php` — `ResultHelper::errorWithCode`, `Export` key
- `Snapshot/Traits/ExporterBuildTrait.php` — `ResultHelper::errorWithCode`, `Incrementals`, `Files` keys
- `Snapshot/Traits/OrchestratorPluginTrait.php` — `ResultHelper::failed/error/ok`, `TotalSize` key
- `Snapshot/Traits/RestoreIncrementalTrait.php` — `Applied` key, multi-line formatting
- `Snapshot/Traits/ImportExecutionTrait.php` — `Folder`, `Incrementals` keys
- `Traits/Snapshot/SnapshotExportHandlerTrait.php` — `Export` key read-side update

> **Remaining files:** ~30+ additional trait files contain `return array(ResponseKeyType::...` patterns that should be migrated in subsequent passes.

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-PHP-005: PHP REST API PascalCase Response Keys Conformance

**Given** PHP service methods, traits, and REST response handlers.
**When** Guideline linters audit response arrays and return structures.
**Then** All response arrays use `ResultHelper` factory methods, `ResponseKeyType` enum cases with PascalCase values, and multi-line formatting with trailing commas with zero bare string keys and zero violations detected.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/04-php --check-only
```
**Expected:** exit 0. Zero violations.
