# PHP Forbidden Patterns — Quick Reference Checklist

> **Version:** 4.0.0  
> **Updated:** 2026-03-09  
> **Consolidates:** [readme.md](./07-php-standards-reference/00-overview.md), [enums.md](./01-enums.md), WP Error Handling <!-- external: spec/28-wp-plugin-development/07-error-handling.md -->

---

## How to Use

Every pattern below is **forbidden** in production code. The ✅ column shows the required replacement. Use this as a pre-commit or code-review checklist.

---

## 1. Error Handling

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 1.1 | `catch (Exception $e)` | `catch (Throwable $e)` | Misses PHP 7+ `Error`, `TypeError`, `ParseError` |
| 1.2 | `$error && in_array($error['type'], [E_ERROR, ...])` | `ErrorChecker::isFatalError($error)` | Duplicated logic; central list in `ErrorType::FATAL_TYPES` |
| 1.3 | Inline `E_*` → string mapping arrays | `ErrorChecker::getTypeLabel($type)` | Uses `ErrorType::TYPE_LABELS`; one place to update |
| 1.4 | `wp_die()` in REST handlers | `wp_send_json_error()` or `$this->envelope->error()` | `wp_die()` breaks JSON response format |
| 1.5 | `error_log()` for diagnostics | `FileLogger` / `$this->fileLogger` | No structure, no stack trace, no audit trail |
| 1.6 | `!class_exists('PDO') \|\| !extension_loaded(...)` inline | `ErrorChecker::isInvalidPdoExtension()` | Centralized; self-documenting |
| 1.7 | Unchecked `new PDO()` without any guard | `ErrorChecker::isInvalidPdoExtension()` check first | Fatal error if extension missing |
| 1.7 | REST handler without `safeExecute` wrapper | Wrap in `$this->safeExecute(fn() => ...)` | Unhandled exceptions crash the endpoint |
| 1.8 | `->error($e->getMessage(), __FILE__, __LINE__)` — message as first param | `->error($e, __FILE__, __LINE__)` — Throwable as first param | Stack trace lost; `error()` requires `Throwable` as first parameter |
| 1.9 | `error_log('Context: ' . $e->getMessage() . "\n" . $e->getTraceAsString())` | `ErrorLog($e, 'Context:')` | Manual concatenation is error-prone; `ErrorLog()` enforces Throwable-first and consistent formatting |

---

## 2. Magic Strings — Hooks

| # | ❌ Forbidden | ✅ Required | Source |
|---|-------------|------------|--------|
| 2.1 | `add_action('init', ...)` | `add_action(HookType::Init->value, ...)` | `HookType::Init` |
| 2.2 | `add_action('plugins_loaded', ...)` | `add_action(HookType::PluginsLoaded->value, ...)` | `HookType::PluginsLoaded` |
| 2.3 | `add_action('rest_api_init', ...)` | `add_action(HookType::RestApiInit->value, ...)` | `HookType::RestApiInit` |
| 2.4 | `add_action('admin_init', ...)` | `add_action(HookType::AdminInit->value, ...)` | `HookType::AdminInit` |
| 2.5 | `add_action('admin_menu', ...)` | `add_action(HookType::AdminMenu->value, ...)` | `HookType::AdminMenu` |
| 2.6 | `add_action('admin_notices', ...)` | `add_action(HookType::AdminNotices->value, ...)` | `HookType::AdminNotices` |
| 2.7 | `add_action('admin_enqueue_scripts', ...)` | `add_action(HookType::AdminEnqueue->value, ...)` | `HookType::AdminEnqueue` |
| 2.8 | `add_action('activated_plugin', ...)` | `add_action(HookType::ActivatedPlugin->value, ...)` | `HookType::ActivatedPlugin` |
| 2.9 | `add_action('deactivated_plugin', ...)` | `add_action(HookType::DeactivatedPlugin->value, ...)` | `HookType::DeactivatedPlugin` |
| 2.10 | `add_action('deleted_plugin', ...)` | `add_action(HookType::DeletedPlugin->value, ...)` | `HookType::DeletedPlugin` |
| 2.11 | `add_filter('rest_post_dispatch', ...)` | `add_filter(HookType::RestPostDispatch->value, ...)` | `HookType::RestPostDispatch` |
| 2.12 | `add_filter('cron_schedules', ...)` | `add_filter(HookType::CronSchedules->value, ...)` | `HookType::CronSchedules` |
| 2.13 | `add_action('wp_ajax_my_action', ...)` | `define('HOOK_AJAX_MY_ACTION', HookType::ajax(ACTION_MY_ACTION));` then `add_action(HOOK_AJAX_MY_ACTION, ...)` | Named composed constant |
| 2.14 | `add_action(HookType::ajax(ACTION_X), ...)` inline | Compose a named constant first, then use it | No inline concatenation at call site |
| 2.15 | `rest_url(REST_NAMESPACE . '/' . ACTION_X)` | `define('REST_URL_X', REST_NAMESPACE . '/' . ACTION_X);` then `rest_url(REST_URL_X)` | No inline concatenation at call site |
| 2.16 | `current_user_can('manage_options')` | `current_user_can(CapabilityType::ManageOptions->value)` | `CapabilityType` enum |
| 2.17 | `'POST'` or `WP_REST_Server::CREATABLE` in routes | `HttpMethodType::Post->value` | `HttpMethodType` enum |

---

## 3. Magic Strings — File Paths

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 3.1 | `WP_CONTENT_DIR . '/uploads/.../file.db'` | `PathHelper::getRootDb()` | Manual concatenation; magic string |
| 3.2 | `PathHelper::getDataDir() . '/file.db'` | `PathHelper::getRootDb()` | Partial accessor; magic filename at call site |
| 3.3 | `PathHelper::getDataDir() . PathDatabaseType::Root->value` | `PathHelper::getRootDb()` | Leaks internal composition to caller |
| 3.4 | Any path without a typed accessor | Create accessor in `PathHelper` first | Every path must have a single-call accessor |

---

## 4. Boolean Logic & No Raw Negations

> **Canonical source:** [No Raw Negations](../01-cross-language/12-no-negatives.md)

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 4.1 | `BooleanHelpers::isFalsy(...)` | `$plugin->isDisabled()` | Generic helper obscures intent |
| 4.2 | `BooleanHelpers::isTruthy(...)` | `$isValue` | Unnecessary indirection |
| 4.3 | `!$plugin->isActive()` | `$plugin->isDisabled()` | Negation is easy to miss; use semantic inverse |
| 4.4 | `$value` for boolean variables | `$isValue`, `$hasPermission` | Ambiguous naming; must use `$is*` / `$has*` prefix |
| 4.5 | `!file_exists($path)` | `PathHelper::isFileMissing($path)` | Raw negation; use positive guard |
| 4.6 | `!is_dir($path)` | `PathHelper::isDirMissing($path)` | Raw negation; use positive guard |
| 4.7 | `!class_exists('X')` | `BooleanHelpers::isClassMissing('X')` | Raw negation; use positive guard |
| 4.8 | `!function_exists('f')` | `BooleanHelpers::isFuncMissing('f')` | Raw negation; use positive guard |
| 4.9 | `!extension_loaded('e')` | `BooleanHelpers::isExtensionMissing('e')` | Raw negation; use positive guard |
| 4.10 | `$isX && !$isY` in `if` | `$isConflict = $isX && !$isY` (extracted) | Mixed polarity; extract to named boolean |

---

## 5. Initialization & Architecture

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 5.1 | WordPress calls in `__construct()` | Lazy `initialize()` method with guard | Load order issues; WP may not be ready |
| 5.2 | Raw `require_once` for non-foundation files | `OnboardIncludeFiles` loader utility | Loader logs failures with stack trace |
| 5.3 | `define()` for categorized constants | Native backed enum (`HookType`, `CapabilityType`, `HttpMethodType`, etc.) | Enums group related constants with PHPDoc |

> **Exception for 5.3:** Plugin-specific custom hooks (e.g., `CRON_SNAPSHOT_*`) may use `define()` in `constants.php` when they are plugin-scoped and not WordPress core hooks.

---

## 5A. Namespace Imports in Traits and Classes — CRITICAL

> **Added:** 2026-02-20 — Root cause of the v1.58.0 → v1.59.0 critical site-down incident.

### Incident Summary

In v1.58.0, two missing `use` import statements caused a **full site crash** ("critical error on this website"):

1. **`PluginRoutesTrait`** referenced `SnapshotRouteRegistrationTrait` (namespace `RiseupAsia\Traits\Snapshot`) without a `use` import. PHP resolved it within the current namespace (`RiseupAsia\Traits\Plugin\SnapshotRouteRegistrationTrait`), which does not exist → **fatal error at class compile time**.
2. **`ActivationHandler`** referenced `PluginConfigType` (namespace `RiseupAsia\Enums`) without a `use` import. On cold activations (no OPcache), this caused a **fatal "class not found" error**.

### Why Traits Are Especially Dangerous

- Trait `use` statements are resolved **at compile time**, before any error handler, shutdown function, or `try/catch` can intervene.
- A missing import in a trait used by the main `Plugin` class **kills the entire WordPress site** — not just the plugin.
- Unlike runtime errors, there is **no graceful degradation** — the PHP parser cannot even finish reading the file.

### Forbidden Patterns

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 5A.1 | `use SomeTraitName;` inside a trait body without a file-level `use` import for cross-namespace traits | Add `use RiseupAsia\Full\Namespace\SomeTraitName;` at the top of the file | PHP resolves unimported names in the current namespace → fatal error |
| 5A.2 | Referencing any enum, class, or interface from another namespace without a `use` import | Add explicit `use` import for every cross-namespace reference | Cold activations and OPcache flushes expose missing imports |
| 5A.3 | Assuming a class/enum "works because it was loaded elsewhere" | Every file must be independently self-sufficient with its own imports | Autoloader order is not guaranteed; OPcache invalidation resets everything |
| 5A.4 | Adding a new `use SomeTrait;` in a trait/class body without verifying the namespace | Check the trait's actual namespace in its source file; add the corresponding `use` import | Copy-paste errors silently point to wrong namespace |

### Dos and Don'ts

#### ✅ DO:

- **Import every cross-namespace symbol** at the top of every PHP file, even if it "seems to work without it"
- **Verify the full namespace path** of any trait, class, or enum before adding a `use` statement
- **Run a cold activation test** (deactivate → delete OPcache → activate) after adding new trait compositions
- **Grep for unimported cross-namespace references** before every release: any `use SomeTrait;` in a trait body must have a matching file-level `use RiseupAsia\...\SomeTrait;`
- **Treat trait files identically to class files** — they do NOT inherit imports from the class that mixes them in

#### ❌ DON'T:

- **Don't rely on OPcache** to mask missing imports — it works until the cache is flushed, then the site crashes
- **Don't assume "it works on my machine" means it's correct** — local environments often have warm OPcache that hides missing imports
- **Don't copy trait `use` statements without checking the namespace** — the trait name alone is not enough; the full namespace must match
- **Don't skip the `use` import because "the autoloader will find it"** — the autoloader resolves the class, but PHP still needs to know which namespace you mean
- **Don't mix up trait `use` (inside class body) with namespace `use` (at file top)** — both are required for cross-namespace traits

### Pre-Release Checklist Addition

```
[ ] Every trait `use` inside a class/trait body has a matching file-level `use` import
[ ] Every enum, class, and interface reference has a file-level `use` import
[ ] No file relies on "ambient" loading from other files
[ ] Cold activation tested (OPcache cleared, plugin deactivated and reactivated)
```

---

## 6. Condition Complexity & Function Size (All Languages)

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 6.1 | Inline `if` with 2+ operators (`&&`, `\|\|`, `!`) | Extract to named `$is*`/`$has*` variable or method | Reads as intent, not implementation |
| 6.2 | `$error && in_array($error['type'], [...])` | `ErrorChecker::isFatalError($error)` | Reusable, self-documenting |
| 6.3 | `!class_exists('PDO') \|\| !extension_loaded(...)` | `ErrorChecker::isInvalidPdoExtension()` | Centralized check |
| 6.4 | Nested `if` (any depth) | **Zero tolerance** — flatten with early returns or combined conditions | Absolute ban |
| 6.5 | Functions > 15 lines | Extract helpers; each function does one thing | Max 15 lines per function body |

---

## 7. Error Type Constants

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 7.1 | `[E_ERROR, E_PARSE, E_CORE_ERROR, ...]` inline | `ErrorType::FATAL_TYPES` | Centralized; update one place for new PHP versions |
| 7.2 | `[E_WARNING, E_NOTICE, ...]` inline | `ErrorType::WARNING_TYPES` | Same principle |
| 7.3 | Custom `errorTypeToString()` functions | `ErrorChecker::getTypeLabel($type)` | Uses `ErrorType::TYPE_LABELS` map |

---

## 8. Magic Strings — Structured Array Keys (`ResponseKeyType`)

> **Added:** 2026-02-20 — Eliminates magic string fragmentation across API responses, log context arrays, and inter-service data transfer.

### Why This Matters

Structured response arrays (`['success' => true, 'error' => '...']`) appear in **every** REST handler, logger call, and service return value. When these keys are raw strings, typos silently break consumers and grep-based audits miss variants (`'errors'` vs `'error'`).

`ResponseKeyType` centralizes all envelope and domain keys so that:
- A single rename propagates everywhere.
- IDE autocompletion prevents typos.
- The Go proxy and TypeScript frontend can mirror the same enum for end-to-end type safety.

### Forbidden Patterns

| # | ❌ Forbidden | ✅ Required | Enum Case |
|---|-------------|------------|-----------|
| 8.1  | `['success']` or `=> 'success'` | `[ResponseKeyType::Success->value]` | `Success` |
| 8.2  | `['error']` | `[ResponseKeyType::Error->value]` | `Error` |
| 8.3  | `['message']` | `[ResponseKeyType::Message->value]` | `Message` |
| 8.4  | `['data']` | `[ResponseKeyType::Data->value]` | `Data` |
| 8.5  | `['code']` | `[ResponseKeyType::Code->value]` | `Code` |
| 8.6  | `['valid']` | `[ResponseKeyType::Valid->value]` | `Valid` |
| 8.7  | `['errors']` | `[ResponseKeyType::Errors->value]` | `Errors` |
| 8.8  | `['cached']` | `[ResponseKeyType::Cached->value]` | `Cached` |
| 8.9  | `['phase']` | `[ResponseKeyType::Phase->value]` | `Phase` |
| 8.10 | `['reason']` | `[ResponseKeyType::Reason->value]` | `Reason` |
| 8.11 | `['total']` | `[ResponseKeyType::Total->value]` | `Total` |
| 8.12 | `['agents']` | `[ResponseKeyType::Agents->value]` | `Agents` |
| 8.13 | `['actions']` | `[ResponseKeyType::Actions->value]` | `Actions` |
| 8.14 | `['logs']` | `[ResponseKeyType::Logs->value]` | `Logs` |
| 8.15 | `['snapshots']` | `[ResponseKeyType::Snapshots->value]` | `Snapshots` |
| 8.16 | `['sql']` | `[ResponseKeyType::Sql->value]` | `Sql` |
| 8.17 | `['params']` | `[ResponseKeyType::Params->value]` | `Params` |
| 8.18 | `['sets']` | `[ResponseKeyType::Sets->value]` | `Sets` |
| 8.19 | `['plugins']` | `[ResponseKeyType::Plugins->value]` | `Plugins` |
| 8.20 | `['tables']` | `[ResponseKeyType::Tables->value]` | `Tables` |
| 8.21 | `['rows']` | `[ResponseKeyType::Rows->value]` | `Rows` |
| 8.22 | `['bytes']` | `[ResponseKeyType::Bytes->value]` | `Bytes` |
| 8.23 | `['size']` | `[ResponseKeyType::Size->value]` | `Size` |
| 8.24 | `['file_size']` | `[ResponseKeyType::FileSize->value]` | `FileSize` |
| 8.25 | `['path']` | `[ResponseKeyType::Path->value]` | `Path` |
| 8.26 | `['filename']` | `[ResponseKeyType::Filename->value]` | `Filename` |
| 8.27 | `['checksum']` | `[ResponseKeyType::Checksum->value]` | `Checksum` |
| 8.28 | `['duration']` | `[ResponseKeyType::Duration->value]` | `Duration` |
| 8.29 | `['count']` | `[ResponseKeyType::Count->value]` | `Count` |
| 8.30 | `['files']` | `[ResponseKeyType::Files->value]` | `Files` |
| 8.31 | `['directory']` | `[ResponseKeyType::Directory->value]` | `Directory` |
| 8.32 | `['scope']` | `[ResponseKeyType::Scope->value]` | `Scope` |
| 8.33 | `['exported']` | `[ResponseKeyType::Exported->value]` | `Exported` |
| 8.34 | `['entry']` | `[ResponseKeyType::Entry->value]` | `Entry` |
| 8.35 | `['snapshot_id']` | `[ResponseKeyType::SnapshotId->value]` | `SnapshotId` |
| 8.36 | `['sequence']` | `[ResponseKeyType::Sequence->value]` | `Sequence` |
| 8.37 | `['folder_name']` | `[ResponseKeyType::FolderName->value]` | `FolderName` |
| 8.38 | `['tables_changed']` | `[ResponseKeyType::TablesChanged->value]` | `TablesChanged` |
| 8.39 | `['total_rows']` | `[ResponseKeyType::TotalRows->value]` | `TotalRows` |
| 8.40 | `['total_new_rows']` | `[ResponseKeyType::TotalNewRows->value]` | `TotalNewRows` |
| 8.41 | `['zip_size']` | `[ResponseKeyType::ZipSize->value]` | `ZipSize` |
| 8.42 | `['backup_id']` | `[ResponseKeyType::BackupId->value]` | `BackupId` |
| 8.43 | `['zip_failed']` | `[ResponseKeyType::ZipFailed->value]` | `ZipFailed` |
| 8.44 | `['skip_audit']` | `[ResponseKeyType::SkipAudit->value]` | `SkipAudit` |
| 8.45 | `['tables_restored']` | `[ResponseKeyType::TablesRestored->value]` | `TablesRestored` |

### Scope

This rule applies to **all** array key access where the key matches a `ResponseKeyType` case value:

- REST handler response arrays (`new WP_REST_Response(array(...))`)
- Service return arrays (`return array(ResponseKeyType::Success->value => true, ...)`)
- Log context arrays (`$this->fileLogger->info('...', array(ResponseKeyType::Phase->value => '...'))`)
- Internal data transfer between traits/classes

### Exceptions

Keys that are **not** in `ResponseKeyType` remain as literal strings (e.g., domain-specific keys like `'retention'`, `'orphans'`, `'stuck'`, `'settings'`, `'providers'`, `'content'`). If a key appears in 3+ files, consider adding it to the enum.

Keys inside `$_FILES` superglobal access (e.g., `$files['file']['size']`) and WordPress hook/filter names are exempt — these are PHP/WordPress API contracts.

### Example

```php
// ❌ FORBIDDEN — magic string keys
return array(
    'success' => true,
    'snapshot_id' => $id,
    'total_rows' => $rows,
    'duration' => $elapsed,
    'errors' => $errors,
);

// ✅ REQUIRED — ResponseKeyType enum
return array(
    ResponseKeyType::Success->value => true,
    ResponseKeyType::SnapshotId->value => $id,
    ResponseKeyType::TotalRows->value => $rows,
    ResponseKeyType::Duration->value => $elapsed,
    ResponseKeyType::Errors->value => $errors,
);
```

---

## 9. Magic Strings — Plugin Identity (`PluginConfigType`)

> **Added:** 2026-02-23 — Eliminates hardcoded plugin name and log prefix strings scattered across the codebase.

### Why This Matters

The plugin name (`'Riseup Asia Uploader'`) and log prefix (`'[Riseup Asia]'`) are already centralized in `PluginConfigType::Name` and `PluginConfigType::LogPrefix`. Hardcoding these strings in class constants, email subjects, admin notices, or generated file comments creates maintenance debt — a rebrand or rename requires a full codebase grep instead of a single enum update.

### Forbidden Patterns

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 9.1 | `private const LOG_PREFIX = '[Riseup Asia] ClassName: '` | Derive from `PluginConfigType::LogPrefix->value` at runtime | Hardcoded prefix duplicates enum value |
| 9.2 | `'[Riseup Asia] Plugin Boot Errors on ' . $site` | `PluginConfigType::LogPrefix->value . ' Plugin Boot Errors on ' . $site` | Email subject uses hardcoded prefix |
| 9.3 | `'Riseup Asia Uploader — Boot Error Report'` | `PluginConfigType::Name->value . ' — Boot Error Report'` | User-facing text uses hardcoded name |
| 9.4 | `'⚠️ Riseup Asia Uploader:'` in admin notices | `'⚠️ ' . PluginConfigType::Name->value . ':'` | Admin HTML uses hardcoded name |
| 9.5 | `'# Riseup Asia Uploader - Security'` in generated files | `'# ' . PluginConfigType::Name->value . ' - Security'` | Generated file comments use hardcoded name |
| 9.6 | `'from the Riseup Asia Uploader plugin.'` | `'from the ' . PluginConfigType::Name->value . ' plugin.'` | Email body uses hardcoded name |
| 9.7 | `'riseup-asia-uploader'` as inline string in templates | `PluginConfigType::Slug->value` passed via controller or used directly | Template uses hardcoded slug |
| 9.8 | `'Riseup Asia Uploader'` as inline string in template HTML | `<?= esc_html(PluginConfigType::Name->value) ?>` | Template uses hardcoded plugin name |

### PHP Const Limitation

PHP `const` expressions cannot reference enum cases (`PluginConfigType::LogPrefix->value`). For classes that previously used `private const LOG_PREFIX = '[Riseup Asia] ClassName: '`, replace with a private static method:

```php
// ❌ FORBIDDEN — hardcoded prefix in const
private const LOG_PREFIX = '[Riseup Asia] AdminMailer: ';

// ✅ REQUIRED — derived from enum at runtime
private static function logPrefix(): string {
    return PluginConfigType::LogPrefix->value . ' AdminMailer: ';
}
```

### Exception

The `Autoloader` class is **exempt** — it loads before enums are available and explicitly cannot depend on `PluginConfigType`. Its hardcoded `LOG_PREFIX` is acceptable.

PHPDoc `@package` headers and file-level doc block comments (e.g., `* Riseup Asia Uploader - File Logger`) are documentation, not logic — they are exempt from this rule.

---

## 10. Hardcoded Date Format Strings (`DateFormatType`)

> **Added:** 2026-02-25 — All date format strings must come from `DateFormatType` enum, not hardcoded in `gmdate()`, `date()`, or class constants.

### Forbidden Patterns

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 10.1 | `gmdate('c')` | `gmdate(DateFormatType::Iso8601->value)` | Hardcoded format string |
| 10.2 | `gmdate('Y-m-d H:i:s')` | `gmdate(DateFormatType::DateTime->value)` | Hardcoded format string |
| 10.3 | `gmdate('Y-m-d')` | `gmdate(DateFormatType::DateOnly->value)` | Hardcoded format string |
| 10.4 | `private const TIMESTAMP_FORMAT = 'Y-m-d\TH:i:s'` | Use `DateFormatType::LogTimestamp->value` | Class const duplicates enum |

### Exception

The `Autoloader` class is **exempt** — it loads before enums are available.

---

## 11. snake_case API Response Keys

> **Added:** 2026-02-25 — API response array keys must use PascalCase, matching the DB column key convention.

### Forbidden Patterns

| # | ❌ Forbidden | ✅ Required | Why |
|---|-------------|------------|-----|
| 11.1 | `$data['plugin_version']` | `$data['PluginVersion']` | snake_case in response |
| 11.2 | `$data['timestamp']` | `$data['Timestamp']` | lowercase in response |
| 11.3 | `$data['log_hint']` | `$data['LogHint']` | snake_case in response |
| 11.4 | `'success' => true` | `'Success' => true` | lowercase in response |
| 11.5 | `'error' => [...]` | `'Error' => [...]` | lowercase in response |

---

## Checklist Summary (Copy for PRs)

```
[ ] No `catch (Exception $e)` — all use `catch (Throwable $e)`
[ ] Every `catch (Throwable $e)` passes `$e` to the logger — stack trace must never be dropped
[ ] No inline `in_array($error['type'], [...])` — use `ErrorChecker`
[ ] No inline E_* → string maps — use `ErrorChecker::getTypeLabel()`
[ ] No `wp_die()` in REST handlers
[ ] No `error_log()` — use structured logger
[ ] No `error_log('...' . $e->getMessage() . ...)` — use `ErrorLog($e, 'context')` for pre-FileLogger exception logging
[ ] No string literals in add_action/add_filter — use `HookType::*->value`
[ ] No inline concatenation at call sites — compose named constants first
[ ] No manual path concatenation — use `PathHelper` accessors
[ ] No `BooleanHelpers` trivial wrappers — use semantic methods
[ ] No `!$obj->isActive()` — use `$obj->isDisabled()`
[ ] No `!file_exists()` — use `PathHelper::isFileMissing()`
[ ] No `!is_dir()` — use `PathHelper::isDirMissing()`
[ ] No `!class_exists()` — use `BooleanHelpers::isClassMissing()`
[ ] No `!function_exists()` — use `BooleanHelpers::isFuncMissing()`
[ ] No `!extension_loaded()` — use `BooleanHelpers::isExtensionMissing()`
[ ] No raw `!` on any function call — use positive guard function
[ ] No boolean vars without `$is*` / `$has*` prefix
[ ] No WordPress calls in constructors
[ ] No inline `!class_exists('PDO')` — use `ErrorChecker::isInvalidPdoExtension()`
[ ] Blank line before `return` or `throw` when preceded by other statements
[ ] No single-line `if (...) return;` — always use braces
[ ] Blank line after closing `}` when followed by more code
[ ] No nested `if` — ZERO TOLERANCE — absolute ban
[ ] No inline multi-part `if` (2+ operators) — extract to `$is*` variable or method
[ ] Functions max 15 lines — extract helpers for longer logic
[ ] No leading backslash on `Throwable` — use `catch (Throwable $e)`
[ ] Functions with >2 params — one param per line with trailing comma
[ ] Every trait `use` in class/trait body has matching file-level `use` import
[ ] Every cross-namespace enum/class/interface has file-level `use` import
[ ] No file relies on ambient loading — each file is self-sufficient
[ ] Cold activation tested after adding new trait compositions
[ ] No magic string array keys matching ResponseKeyType cases — use enum->value
[ ] Every new repeated key (3+ files) added to ResponseKeyType enum
[ ] No hardcoded date format strings — use `DateFormatType::*->value`
[ ] No snake_case API response keys — use PascalCase
```

---

## Cross-References

- [PHP Coding Standards](./07-php-standards-reference/00-overview.md) — Full spec with examples
- [PHP Enum Classes](./01-enums.md) — `HookType`, `CapabilityType`, `HttpMethodType`, Path enums, `ErrorType`, `ErrorChecker`
- [Cross-Language Code Style](../01-cross-language/04-code-style/00-overview.md) — Rules 1-9 (braces, nesting, spacing, function size, Throwable, multi-line params)
- WordPress Error Handling <!-- external: spec/28-wp-plugin-development/07-error-handling.md -->
- WordPress Initialization <!-- external: spec/28-wp-plugin-development/01-initialization-patterns.md -->
- WordPress API Design <!-- external: spec/28-wp-plugin-development/04-api-design.md -->

---

*Forbidden patterns checklist v5.0.0 — 2026-02-25*
