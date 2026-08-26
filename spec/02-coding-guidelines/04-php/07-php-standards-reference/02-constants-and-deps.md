# PHP Coding Standards — Constants, enums, dependency checks, file paths

> **Parent:** [PHP Coding Standards](./00-overview.md)  
> **Version:** 5.1.0  
> **Updated:** 2026-03-31

---

## Constants & Enums — No Magic Strings

### Rule: All identifiers in `constants.php` or native backed enums

Every endpoint path, action name, capability string, option key, **hook name**, **file path segment**, **HTTP method**, and **WordPress capability** must be defined centrally. Use PHP `constants.php` for simple values and **PHP 8.1+ native backed enums** in `includes/Enums/` for categorized groups.

> **See [enums.md](../01-enums.md)** for the full enum specification (v4.0.0), including file naming rules, namespace conventions, and all enum/const class definitions.

### Hook Names — HookType enum

```php
// ❌ FORBIDDEN: Magic hook strings
add_action('init', [$this, 'setup']);
add_action('rest_api_init', [$this, 'registerRoutes']);
add_action('plugins_loaded', [$this, 'onPluginsLoaded']);

// ✅ REQUIRED: Hook names from HookType enum
use RiseupAsia\Enums\HookType;

add_action(HookType::Init->value, [$this, 'setup']);
add_action(HookType::RestApiInit->value, [$this, 'registerRoutes']);
add_action(HookType::PluginsLoaded->value, [$this, 'onPluginsLoaded']);
```

### Action Names — Named Composed Constants

Inline concatenation at call sites is **forbidden** — even when using centralized base constants. Instead, compose descriptively named constants from base constants, then use those named constants directly.

```php
// ❌ FORBIDDEN: Magic strings
add_action('wp_ajax_my_action', [$this, 'handle']);
$url = rest_url('riseup-asia-uploader/v1/upload');

// ❌ FORBIDDEN: Inline concatenation at call site (even with constants)
// In constants.php:
define('REST_NAMESPACE', 'riseup-asia-uploader/v1');
define('ACTION_UPLOAD', 'upload');
// In handlers:
add_action(Hook::ajax('upload') , [$this, 'handle']);  // ← still inline concat
$url = rest_url(REST_NAMESPACE . '/' . ACTION_UPLOAD);

// ✅ REQUIRED: Compose named constants, then use them directly
// In constants.php:
define('REST_NAMESPACE', 'riseup-asia-uploader/v1');
define('ACTION_UPLOAD', 'upload');
define('REST_URL_UPLOAD', REST_NAMESPACE . '/' . ACTION_UPLOAD);
define('HOOK_AJAX_UPLOAD', HookType::ajax(ACTION_UPLOAD));

// In handlers — clean, readable, no concatenation:
add_action(HOOK_AJAX_UPLOAD, [$this, 'handle']);
$url = rest_url(REST_URL_UPLOAD);
```

> **Naming:** Constants must NOT use the `RISEUP_` prefix. Use descriptive names that convey purpose: `HOOK_AJAX_UPLOAD`, `REST_URL_UPLOAD`, `REST_NAMESPACE`.

---


---

## Dependency Checks

### Rule: Delegate to ErrorChecker — no inline extension checks

Before using external dependencies (PDO, extensions), verify availability via `ErrorChecker`. Never write inline `class_exists()` / `extension_loaded()` checks in business logic.

```php
// ❌ FORBIDDEN: Inline extension checks in business logic
if (!class_exists('PDO') || !extension_loaded('pdo_sqlite')) {
    $this->logger->error('PDO/SQLite not available');
    return $this->envelope->error('SQLite support not available', 500);
}

// ✅ REQUIRED: Centralized check via ErrorChecker
if (ErrorChecker::isInvalidPdoExtension()) {
    $this->logger->error('PDO/SQLite not available');

    return $this->envelope->error('SQLite support not available', 500);
}
```

Throttle repeated initialization errors to prevent log bloat.

---


---

## File Path Resolution

### Rule: Use fully-typed path accessors backed by PathConst constants

Never construct file paths with string concatenation or partial accessors. Every path must resolve to a **single typed accessor method** that internally composes a directory method + a `PathConst` constant.

### How It Works (Internal Architecture)

```
Caller code          →  PathHelper::getRootDb()
                              ↓
Accessor internals   →  self::getDataDir() + PathConst::ROOT_DB
                              ↓                    ↓
Directory method     →  WP_CONTENT_DIR + ...   '/a-root.db'
                              ↓
Final path           →  '/var/www/.../uploads/riseup-asia-uploader/a-root.db'
```

**The caller only ever sees the accessor.** The composition of directory + constant is an internal implementation detail.

### Forbidden vs Required

```php
// ❌ FORBIDDEN: Manual path construction with string literals
$path = WP_CONTENT_DIR . '/uploads/riseup-asia-uploader/data.db';

// ❌ FORBIDDEN: Partial accessor — caller still concatenates a magic string
$path = PathHelper::getDataDir() . '/data.db';

// ❌ FORBIDDEN: Using PathConst directly in business logic (leaks internals)
$path = PathHelper::getDataDir() . PathConst::ROOT_DB;

// ✅ REQUIRED: Single typed accessor — no path fragments visible to caller
$path = PathHelper::getRootDb();
```

### Why This Matters

| If you use... | Problem |
|---------------|---------|
| Manual concatenation | Filename is a magic string; renaming requires find-and-replace |
| `getDataDir() . '/file.db'` | Partial accessor; magic string still exists at the call site |
| `getDataDir() . PathConst::X` | Leaks the composition pattern; callers shouldn't know how paths are built |
| `getRootDb()` ✅ | Filename lives in `PathConst`, directory in `getDataDir()`, both hidden from caller |

> **Rule:** If a path does not have a typed accessor in `PathHelper`, create one before using it. See [PHP Enum Spec](../01-enums.md) for full `PathEnum` and `PathHelper` listings.

---
