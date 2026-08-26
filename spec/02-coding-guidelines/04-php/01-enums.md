# PHP Enums — Complete Reference

> **Version:** 7.1.0  
> **Updated:** 2026-03-09  
> **Applies to:** WordPress companion plugins (PHP 8.1+)

---

## Overview

All enum-like constants MUST use **PHP 8.1+ native backed enums** with proper namespaces.
The old pattern of `class FooEnum { public const BAR = '...'; }` and `define()` constants
is **deprecated** and must be migrated.

### Naming Convention: `Type` Suffix

All enums MUST use the **`Type` suffix** in their name. This clearly distinguishes enums from classes and makes the type nature explicit at every usage site.

| ❌ Forbidden Name | ✅ Required Name |
|------------------|-----------------|
| `UploadSource` | `UploadSourceType` |
| `Capability` | `CapabilityType` |
| `HttpMethod` | `HttpMethodType` |
| `Hook` | `HookType` |

> **Non-enum constant classes** (`ErrorType`) keep their existing names — they are `final class`, not `enum`. The former `PathConst` class has been decomposed into 4 domain-specific enums (see below).

### Architectural Rules

1. **All enums live in `includes/Enums/`** — one file per enum.
2. **File name = Definition name** — e.g., `UploadSourceType.php` → contains `enum UploadSourceType: string`.
3. **Namespace:** `RiseupAsia\\Enums` — every enum file declares this namespace.
4. **`Type` suffix required** — use `UploadSourceType`, not `UploadSource`.
5. **String-backed** (`enum Foo: string`) for all enums whose values are strings.
6. **Case names use PascalCase** — `case RestApi`, not `case REST_API`.
7. **No `RISEUP_` prefix** on anything — namespace provides scoping.
8. **`define()` constants are prohibited** for values that belong in an enum.
9. **Access pattern:** `UploadSourceType::Script` (the enum case) or `UploadSourceType::Script->value` (the raw string).
10. **Validation helpers** go as `static` methods on the enum itself (camelCase: `validValues()`, `isValid()`).
11. **Non-enum constants classes** (ErrorType) use the same namespace and folder but remain `final class` with `public const`.
12. **`isEqual()` method required** — every backed enum MUST include the `isEqual(self $other): bool` instance method (see below).

### File Loading

Enum files are loaded via `require_once` before the dependency loader:

```php
// In riseup-asia-uploader.php (bootstrap)
require_once __DIR__ . '/includes/Enums/UploadSourceType.php';
require_once __DIR__ . '/includes/Enums/CapabilityType.php';
require_once __DIR__ . '/includes/Enums/HttpMethodType.php';
require_once __DIR__ . '/includes/Enums/HookType.php';
require_once __DIR__ . '/includes/Enums/EndpointType.php';
require_once __DIR__ . '/includes/Enums/PathSubdirType.php';
require_once __DIR__ . '/includes/Enums/PathDatabaseType.php';
require_once __DIR__ . '/includes/Enums/PathLogFileType.php';
require_once __DIR__ . '/includes/Enums/PathConfigType.php';
require_once __DIR__ . '/includes/Enums/WpErrorCodeType.php';
require_once __DIR__ . '/includes/Enums/ErrorType.php';
```

At call sites, use the `use` import:

```php
use RiseupAsia\\Enums\\UploadSourceType;
use RiseupAsia\\Enums\\CapabilityType;
```

---

## isEqual() — Universal Enum Comparison

### Rule: Every backed enum MUST have `isEqual(self $other): bool`

This method provides a clean, fluent API for comparing enum cases. It replaces raw `===` comparisons at call sites and inside helper methods, making conditionals more readable and eliminating bare operator usage.

### Implementation (identical in every enum)

```php
/** Check if this enum case equals the given case. */
public function isEqual(self $other): bool
{
    return $this === $other;
}
```

### Why isEqual() Instead of Raw `===`

| Aspect | `===` (forbidden) | `isEqual()` (required) |
|--------|-------------------|----------------------|
| Readability | `$status === StatusType::Success` | `$status->isEqual(StatusType::Success)` |
| Fluency | Operator-based, breaks chain | Method-based, reads like English |
| Internal helpers | `$this === self::X` | `$this->isEqual(self::X)` |
| Consistency | Mixed styles across codebase | Single pattern everywhere |

### Usage — Call Sites

```php
use RiseupAsia\Enums\StatusType;
use RiseupAsia\Enums\LogLevelType;

// ❌ FORBIDDEN: Raw === comparison
if ($status === StatusType::Success) { ... }

if ($level === LogLevelType::Error || $level === LogLevelType::Warn) { ... }

// ✅ REQUIRED: isEqual() method
if ($status->isEqual(StatusType::Success)) { ... }

if ($level->isEqual(LogLevelType::Error) || $level->isEqual(LogLevelType::Warn)) { ... }
```

### Usage — Internal Helper Methods

Existing domain-specific helpers (e.g., `isSuccess()`, `isError()`) MUST delegate to `isEqual()` internally:

```php
enum StatusType: string
{
    case Success = 'success';
    case Failed  = 'failed';

    public function isEqual(self $other): bool
    {
        return $this === $other;
    }

    // ❌ FORBIDDEN: Direct === in helpers
    public function isSuccess(): bool
    {
        return $this === self::Success;
    }

    // ✅ REQUIRED: Delegate to isEqual()
    public function isSuccess(): bool
    {
        return $this->isEqual(self::Success);
    }
}
```

### Usage — Compound Checks

For helpers that check multiple cases, each comparison uses `isEqual()`:

```php
// ✅ Compound check with isEqual()
public function isLifecycle(): bool
{
    return $this->isEqual(self::Enable)
        || $this->isEqual(self::Disable)
        || $this->isEqual(self::Delete);
}

public function isErrorOrWarn(): bool
{
    return $this->isEqual(self::Error) || $this->isEqual(self::Warn);
}
```

### When NOT to Use isEqual()

- **Domain checks using `str_starts_with()`** — These are prefix-based, not case-based. Keep as-is:
  ```php
  // ✅ Correct: prefix check, not enum comparison
  public function isSnapshot(): bool
  {
      return str_starts_with($this->value, 'snapshot_');
  }
  ```
- **Static validation** — `tryFrom()` and `validValues()` are not comparisons.

---

## Complete Enum Inventory

### UploadSourceType — Upload Origin

Identifies how a plugin upload was initiated.

```php
enum UploadSourceType: string
{
    case Script  = 'upload_script';
    case RestApi = 'rest_api';
    case AdminUi = 'admin_ui';
    case WpCli   = 'wp_cli';

    public function isEqual(self $other): bool { return $this === $other; }

    public static function validValues(): array
    {
        return array_column(self::cases(), 'value');
    }

    public static function isValid(string $source): bool
    {
        return self::tryFrom($source) !== null;
    }
}
```

#### Usage

```php
use RiseupAsia\\Enums\\UploadSourceType;

// ❌ FORBIDDEN
define('UPLOAD_SOURCE_SCRIPT', 'upload_script');

if ($source === UploadSourceType::Script) { ... }

// ✅ REQUIRED
$source  = UploadSourceType::Script;
$value   = UploadSourceType::Script->value;
$parsed  = UploadSourceType::tryFrom('rest_api');
$isValid = UploadSourceType::isValid($input);

if ($source->isEqual(UploadSourceType::RestApi)) { ... }
```

---

### CapabilityType — WordPress Capabilities

```php
enum CapabilityType: string
{
    case ManageOptions   = 'manage_options';
    case ActivatePlugins = 'activate_plugins';
    case PublishPosts    = 'publish_posts';
    case UploadFiles     = 'upload_files';
    case EditPosts       = 'edit_posts';
    case DeletePlugins   = 'delete_plugins';
    case InstallPlugins  = 'install_plugins';
    case UpdatePlugins   = 'update_plugins';
    case SwitchThemes    = 'switch_themes';
    case ManageUsers     = 'manage_users';
    case ManageNetwork   = 'manage_network';

    public function isEqual(self $other): bool { return $this === $other; }
}
```

#### Usage

```php
use RiseupAsia\\Enums\\CapabilityType;

// ❌ FORBIDDEN
if (current_user_can('manage_options')) { ... }

// ✅ REQUIRED
if (current_user_can(CapabilityType::ManageOptions->value)) { ... }
```

---

### HttpMethodType — REST API Methods

```php
enum HttpMethodType: string
{
    case Get    = 'GET';
    case Post   = 'POST';
    case Put    = 'PUT';
    case Patch  = 'PATCH';
    case Delete = 'DELETE';

    public function isEqual(self $other): bool { return $this === $other; }

    public static function editable(): string
    {
        return 'PUT, PATCH';
    }
}
```

#### Usage

```php
use RiseupAsia\\Enums\\HttpMethodType;

// ❌ FORBIDDEN
register_rest_route($ns, '/upload', ['methods' => 'POST', ...]);

// ✅ REQUIRED
register_rest_route($ns, '/upload', ['methods' => HttpMethodType::Post->value, ...]);
```

---

### HookType — WordPress Hook Names

```php
enum HookType: string
{
    // ── Core Lifecycle ──────────────────────────────────────────
    case Init           = 'init';
    case PluginsLoaded  = 'plugins_loaded';
    case RestApiInit    = 'rest_api_init';
    case AdminInit      = 'admin_init';
    case Shutdown       = 'shutdown';

    // ── Plugin Lifecycle ────────────────────────────────────────
    case ActivatedPlugin   = 'activated_plugin';
    case DeactivatedPlugin = 'deactivated_plugin';
    case DeletedPlugin     = 'deleted_plugin';

    // ── Admin UI ────────────────────────────────────────────────
    case AdminNotices   = 'admin_notices';
    case AdminEnqueue   = 'admin_enqueue_scripts';
    case AdminMenu      = 'admin_menu';

    // ── Filters ─────────────────────────────────────────────────
    case RestPostDispatch                  = 'rest_post_dispatch';
    case PluginActionLinks                 = 'plugin_action_links';
    case PreSetSiteTransientUpdatePlugins  = 'pre_set_site_transient_update_plugins';
    case PluginsApi                        = 'plugins_api';
    case CronSchedules                     = 'cron_schedules';

    public function isEqual(self $other): bool { return $this === $other; }

    public static function ajax(string $action): string
    {
        return 'wp_ajax_' . $action;
    }

    public static function ajaxNopriv(string $action): string
    {
        return 'wp_ajax_nopriv_' . $action;
    }
}
```

#### Usage

```php
use RiseupAsia\\Enums\\HookType;

// ❌ FORBIDDEN
add_action('rest_api_init', [$this, 'registerRoutes']);

// ✅ REQUIRED
add_action(HookType::RestApiInit->value, [$this, 'registerRoutes']);
add_action(HookType::ajax('riseup_test'), [$this, 'ajaxTest']);
```

---

### LogLevelType — Log Severity Levels

```php
enum LogLevelType: string
{
    case Debug = 'DEBUG';
    case Info  = 'INFO';
    case Warn  = 'WARN';
    case Error = 'ERROR';

    public function isEqual(self $other): bool { return $this === $other; }

    public function isError(): bool { return $this->isEqual(self::Error); }
    public function isWarn(): bool  { return $this->isEqual(self::Warn); }
    public function isInfo(): bool  { return $this->isEqual(self::Info); }
    public function isDebug(): bool { return $this->isEqual(self::Debug); }

    public function isErrorOrWarn(): bool
    {
        return $this->isEqual(self::Error) || $this->isEqual(self::Warn);
    }
}
```

#### Usage

```php
use RiseupAsia\Enums\LogLevelType;

// ❌ FORBIDDEN
if ($level === LogLevelType::Error) { ... }

// ✅ REQUIRED — use domain helper or isEqual()
if ($level->isError()) { ... }

if ($level->isEqual(LogLevelType::Error)) { ... }
```

---

### StatusType — Transaction Result Status

```php
enum StatusType: string
{
    case Success = 'success';
    case Failed  = 'failed';

    public function isEqual(self $other): bool { return $this === $other; }

    public function isSuccess(): bool { return $this->isEqual(self::Success); }
    public function isFailed(): bool  { return $this->isEqual(self::Failed); }
}
```

---

### PostStatusType — WordPress Post Statuses

```php
enum PostStatusType: string
{
    case Publish = 'publish';
    case Draft   = 'draft';
    case Pending = 'pending';

    public function isEqual(self $other): bool { return $this === $other; }

    public function isPublic(): bool { return $this->isEqual(self::Publish); }

    public static function validValues(): array
    {
        return array_map(fn(self $case) => $case->value, self::cases());
    }
}
```

---

### ActionType — Transaction Logging Actions

42 cases across Core, Post, Auth, Export, Update, Agent, Snapshot domains.

```php
enum ActionType: string
{
    // Core: Upload, UploadActive, UploadInitiated, Enable, Disable, Delete, ...
    // Post: PostCreate, PostUpdate, CategoryCreate, MediaUpload
    // Auth: AuthFailed
    // Export: ExportSelf, ExportPlugin
    // Update: UpdateCheck, UpdateResolve, UpdateDownload, UpdateInstall
    // Agent: AgentAdd, AgentRemove, AgentTest, AgentSync, ...
    // Snapshot: SnapshotCreate, SnapshotRestore, SnapshotDelete, ...

    public function isEqual(self $other): bool { return $this === $other; }

    // Domain prefix checks (str_starts_with, NOT isEqual)
    public function isSnapshot(): bool { return str_starts_with($this->value, 'snapshot_'); }
    public function isAgent(): bool    { return str_starts_with($this->value, 'agent_'); }
    public function isUpdate(): bool   { return str_starts_with($this->value, 'update_'); }

    // Compound case check (uses isEqual)
    public function isLifecycle(): bool
    {
        return $this->isEqual(self::Enable)
            || $this->isEqual(self::Disable)
            || $this->isEqual(self::Delete);
    }
}
```

---

### TableType — SQLite Table Names (PascalCase)

All custom SQLite table names use **PascalCase** values. This aligns with the [cross-language database naming convention](../01-cross-language/07-database-naming.md).

```php
enum TableType: string
{
    case Transactions     = 'Transactions';
    case AgentSites       = 'AgentSites';
    case AgentActions     = 'AgentActions';
    case Snapshots        = 'Snapshots';
    case SnapshotProgress = 'SnapshotProgress';
    case SnapshotJobs     = 'SnapshotJobs';
    case SnapshotSettings = 'SnapshotSettings';
    case SnapshotExports  = 'SnapshotExports';
    case FileCache        = 'FileCache';
    case RemotePluginsCache = 'RemotePluginsCache';
    case ErrorSessions    = 'ErrorSessions';
    case FlashState       = 'FlashState';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }

    public function isSnapshot(): bool { return str_starts_with($this->value, 'Snapshot'); }
    public function isAgent(): bool    { return str_starts_with($this->value, 'Agent'); }
}
```

### LogColumnType — Log Table Column Names (PascalCase)

Enum for type-safe access to Transactions table columns. Used by `LogValueTrait` for `logValue()` / `logString()` calls.

```php
enum LogColumnType: string
{
    case Id            = 'Id';
    case Action        = 'Action';
    case PluginSlug    = 'PluginSlug';
    case PluginFile    = 'PluginFile';
    case PluginVersion = 'PluginVersion';
    case PostId        = 'PostId';
    case Status        = 'Status';
    case Details       = 'Details';
    case ErrorMsg      = 'ErrorMsg';
    case UserLogin     = 'UserLogin';
    case UserId        = 'UserId';
    case IpAddress     = 'IpAddress';
    case TriggeredBy   = 'TriggeredBy';
    case UploadSource  = 'UploadSource';
    case SourceMachine = 'SourceMachine';
    case CreatedAt     = 'CreatedAt';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

---

### EndpointType — REST API Endpoint Paths

Stores only the path fragment for each REST endpoint. The full WordPress route
is constructed via the `route()` helper, which prepends `/`.

```php
enum EndpointType: string
{
    // Core: Status, Upload, Plugins, ExportSelf, Posts, ...
    // Plugin: PluginFiles, PluginFile, PluginEnable, ...
    // Sync: SyncManifest, Sync
    // Agent: Agents, AgentsAdd, AgentsRemove, ...
    // Snapshot: SnapshotList, SnapshotSchedule, ...

    public function isEqual(self $other): bool { return $this === $other; }

    public function route(): string { return '/' . $this->value; }

    public function isSnapshot(): bool { return str_starts_with($this->value, 'snapshots/'); }
    public function isAgent(): bool    { return str_starts_with($this->value, 'agents'); }
    public function isPlugin(): bool   { return str_starts_with($this->value, 'plugins/'); }
}
```

#### Usage in Route Registration

```php
use RiseupAsia\Enums\EndpointType;
use RiseupAsia\Enums\HttpMethodType;

// ❌ FORBIDDEN: Accessing ->value directly for route construction
$safeRegister(EndpointType::Upload->value, [...]);

// ✅ REQUIRED: Use route() helper
$safeRegister(EndpointType::Upload->route(), array(
    'methods'  => HttpMethodType::Post->value,
    'callback' => array($this, 'handleUpload'),
    ...
));
```

#### When to Use ->value vs ->route()

| Context | Use | Example |
|---------|-----|---------|
| Route registration (`register_rest_route`) | `->route()` | `EndpointType::Upload->route()` → `'/upload'` |
| Building remote API URLs | `->value` | `$baseUrl . '/' . EndpointType::Upload->value` |
| Logging / display | `->value` or `->name` | `'Endpoint: ' . EndpointType::Upload->value` |
| Domain checks | helpers | `$endpoint->isSnapshot()` |

---

## Path Enums — 4 Domain-Specific Enums (replaces PathConst)

The former `PathConst` final class has been decomposed into 4 backed enums. Each answers "which one?" for its domain, qualifying as a proper enum with the `Type` suffix. All include `isEqual()`.

### PathSubdirType — Plugin Subdirectories

```php
enum PathSubdirType: string
{
    case Logs      = '/logs';
    case Temp      = '/temp';
    case Snapshots = '/snapshots';
    case Exports   = '/exports';

    public function isEqual(self $other): bool { return $this === $other; }
}
```

### PathDatabaseType — SQLite Database Files

```php
enum PathDatabaseType: string
{
    case Root     = '/a-root.db';
    case Activity = '/activity.db';
    case Snapshot = '/snapshots.db';
    case Plugin   = '/riseup-asia-uploader.db';

    public function isEqual(self $other): bool { return $this === $other; }
}
```

### PathLogFileType — Log File Names

```php
enum PathLogFileType: string
{
    case Log        = '/log.txt';
    case FatalError = '/fatal-errors.log';
    case Stacktrace = '/stacktrace.txt';
    case Error      = '/error.txt';

    public function isEqual(self $other): bool { return $this === $other; }
}
```

### PathConfigType — Config File Names

```php
enum PathConfigType: string
{
    case Detection = '/wp-plugin-detected.json';

    public function isEqual(self $other): bool { return $this === $other; }
}
```

### Usage in PathHelper

```php
use RiseupAsia\Enums\PathSubdirType;
use RiseupAsia\Enums\PathDatabaseType;

// ❌ FORBIDDEN
$logsDir = self::join(self::getBaseDir(), LOGS_SUBDIR);

// ✅ REQUIRED
$logsDir = self::join(self::getBaseDir(), PathSubdirType::Logs->value);
$dbPath  = self::join(self::getBaseDir(), PathDatabaseType::Plugin->value);
```

---

## WpErrorCodeType — WordPress REST API Error Codes

`WpErrorCodeType` centralizes all `WP_Error` code strings used in REST API permission callbacks, validation, and error responses. Eliminates magic strings in `new WP_Error()` calls.

```php
namespace RiseupAsia\Enums;

enum WpErrorCodeType: string
{
    // Auth & Permission
    case RestForbidden          = 'rest_forbidden';
    case NotAuthenticated       = 'not_authenticated';
    case InsufficientPermissions = 'insufficient_permissions';
    case NoToken                = 'no_token';
    case InvalidToken           = 'invalid_token';
    case RateLimited            = 'rate_limited';

    // Validation
    case ValidationFailed       = 'validation_failed';
    case ValidationError        = 'validation_error';

    // Operations
    case UploadFailed           = 'upload_failed';
    case AuthenticationFailed   = 'authentication_failed';
    case FatalError             = 'fatal_error';

    public function isEqual(self $other): bool {
        return $this->value === $other->value;
    }
}
```

### Usage

```php
use RiseupAsia\Enums\WpErrorCodeType;

// ❌ FORBIDDEN
return new WP_Error('not_authenticated', 'Authentication required', ['status' => 401]);

// ✅ REQUIRED
return new WP_Error(WpErrorCodeType::NotAuthenticated->value, 'Authentication required', ['status' => 401]);
```

---

## ErrorType — PHP Error Type Constants (Non-Enum Class)

`ErrorType` holds arrays/maps — not a backed enum. Does NOT get `isEqual()`.

```php
final class ErrorType
{
    public const FATAL_TYPES = [E_ERROR, E_PARSE, E_CORE_ERROR, E_COMPILE_ERROR, E_USER_ERROR];
    public const WARNING_TYPES = [E_WARNING, E_CORE_WARNING, E_USER_WARNING, E_NOTICE, ...];
    public const RECOVERABLE_TYPES = [E_RECOVERABLE_ERROR, E_STRICT];
    public const TYPE_LABELS = [E_ERROR => 'E_ERROR', ...];
}
```

---

## Classification: Enum vs Const Class

| Name               | Type         | Suffix | Has isEqual() | Why                                              |
|--------------------|--------------|--------|---------------|--------------------------------------------------|
| `UploadSourceType` | `enum`       | `Type` | ✅            | Discrete set — "which source?"                   |
| `CapabilityType`   | `enum`       | `Type` | ✅            | Discrete capabilities — "which permission?"      |
| `HttpMethodType`   | `enum`       | `Type` | ✅            | Discrete HTTP verbs — "which method?"            |
| `HookType`         | `enum`       | `Type` | ✅            | Discrete hook names — "which hook?"              |
| `EndpointType`     | `enum`       | `Type` | ✅            | Discrete REST paths — "which endpoint?"          |
| `LogLevelType`     | `enum`       | `Type` | ✅            | Discrete log levels — "which severity?"          |
| `StatusType`       | `enum`       | `Type` | ✅            | Discrete results — "success or failed?"          |
| `PostStatusType`   | `enum`       | `Type` | ✅            | Discrete post states — "which status?"           |
| `ActionType`       | `enum`       | `Type` | ✅            | Discrete actions — "which action?"               |
| `TableType`        | `enum`       | `Type` | ✅            | Discrete tables — "which table?"                 |
| `PathSubdirType`   | `enum`       | `Type` | ✅            | Discrete subdirectories — "which directory?"     |
| `PathDatabaseType` | `enum`       | `Type` | ✅            | Discrete DB files — "which database?"            |
| `PathLogFileType`  | `enum`       | `Type` | ✅            | Discrete log files — "which log?"                |
| `PathConfigType`   | `enum`       | `Type` | ✅            | Discrete config files — "which config?"          |
| `WpErrorCodeType`  | `enum`       | `Type` | ✅            | Discrete WP_Error codes — "which error code?"    |
| `ErrorType`        | `final class`| —      | ❌            | Arrays of E_* constants and label maps           |

### Decision Rule

> If the type answers **"which one of these?"** with a single value → `enum` with `Type` suffix + `isEqual()`.  
> If it holds **arrays, maps, or composable fragments** → `final class` with `public const`.

---

## Adding New Enum Cases — Checklist

1. **Add the case** to the appropriate enum in `includes/Enums/`.
2. **Ensure `isEqual()` exists** — it should already be there; verify.
3. **Add a PHPDoc comment** if the case is non-obvious.
4. **If PathSubdirType:** Add a corresponding typed accessor to `PathHelper`.
5. **If PathDatabaseType/PathLogFileType/PathConfigType:** Add a typed accessor to `PathHelper`.
6. **If HookType:** Update all `add_action`/`add_filter` calls.
7. **If CapabilityType:** Update all `current_user_can()` calls.
8. **If HttpMethodType:** Update all `register_rest_route()` calls.
9. **If EndpointType:** Add the case, then use `->route()` in route registration. Update all callers.
10. **If ErrorType:** Add to the appropriate group array AND to `TYPE_LABELS`.
11. **If WpErrorCodeType:** Update all `new WP_Error()` calls and `$this->envelope->error()` code parameters.
12. **Never skip the enum** — even for "one-time" usage.
13. **Use `isEqual()` for all comparisons** — never raw `===` at call sites.

---

## Cross-References

- [PHP Coding Standards](./07-php-standards-reference/00-overview.md) — Parent spec with forbidden patterns
- [Naming Conventions](./03-naming-conventions.md) — PascalCase for enums, camelCase for methods
- [Golang Standards](../03-golang/04-golang-standards-reference/00-overview.md) — Go equivalent patterns

## Log Context Array Keys — camelCase

Internal log context array keys (passed to `fileLogger` and `logger` calls) MUST use **camelCase**. This matches the zero-underscore policy for all logic-level identifiers.

```php
// ❌ FORBIDDEN: snake_case log context keys
$this->fileLogger->info('Post created', array('post_id' => $postId));
$this->fileLogger->warn('Duplicate detected', array('duplicate_dir' => $dir));

// ✅ REQUIRED: camelCase log context keys
$this->fileLogger->info('Post created', array('postId' => $postId));
$this->fileLogger->warn('Duplicate detected', array('duplicateDir' => $dir));
```

**Exempt:** Persistence-level keys (wp_options, SQLite `_snapshot_meta` table keys, V13 migration mappings), WordPress API return values (e.g., `$result['term_id']`), and internal markers (e.g., `_invocation_chain`).

---

*PHP Enum specification v7.1.0 — 2026-02-23*
