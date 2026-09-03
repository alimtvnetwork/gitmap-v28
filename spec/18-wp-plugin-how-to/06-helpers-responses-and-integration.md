# Phase 5 — Helpers, Response Envelope, and Integration

> **Purpose:** Define helper class patterns, the standard API response format, and how all pieces integrate.

---

## 5.1 Helper Classes

Helpers are **stateless static utility classes**. They never hold instance state, never depend on `$this`, and never access WordPress hooks.

### Standard helpers

| Helper | Responsibility |
|--------|---------------|
| `DateHelper` | All timestamp formatting and timezone conversion (see Phase 4, §4.13) |
| `PathHelper` | File path resolution, directory creation, uploads dir resolution (§5.1) |
| `BooleanHelpers` | Semantic boolean guards — readable negation wrappers (§5.1.1) |
| `InitHelpers` | Bootstrap initialization, SQLite connection, component tracking (§5.1.2) |
| `ErrorLogHelper` | Native `error_log()` wrapper with stack traces (see Phase 4, §4.11) |
| `EnvelopeBuilder` | Constructs the standard API response envelope (§5.3) |
| `HttpConfigType` | HTTP request option factories for `wp_remote_*()` (§5.1.3) — note: this is an enum with static factory methods |
| `TypeCheckerTrait` | Syntax-validator-safe type checking (see Phase 3, §3.8) — note: this is a trait, not a helper, because it needs `$this` context |

### PathHelper specification

PathHelper centralises all file system path resolution. It uses trait decomposition for large helpers:

```
Helpers/
├── PathHelper.php              ← Shell class, composes 3 traits
└── Traits/
    ├── PathHelperCoreTrait.php  ← Base dir, join(), typed accessors
    ├── PathHelperDirTrait.php   ← Directory guards, creation, security files
    └── PathHelperFileTrait.php  ← File guards, delete, relative paths
```

#### Core methods (PathHelperCoreTrait)

| Method | Returns |
|--------|---------|
| `join(string ...$segments)` | Normalised path from segments (forward slashes, no doubles) |
| `getBaseDir()` | `wp-content/uploads/{plugin-slug}` via `wp_upload_dir()` |
| `getLogsDir()` | `{baseDir}/logs` — uses `PathSubdirType::Logs` enum |
| `getTempDir()` | `{baseDir}/temp` — uses `PathSubdirType::Temp` enum |
| `getDbPath()` | `{baseDir}/riseup-asia-uploader.db` — uses `PathDatabaseType::Plugin` enum |
| `getPluginDir()` | `WP_PLUGIN_DIR/{slug}` |
| `getPluginMainFile()` | `{pluginDir}/{slug}.php` |
| `getConstantsFile()` | `{pluginDir}/includes/constants.php` |
| `getEndpointsJsonPath()` | `{pluginDir}/data/endpoints.json` |

#### Directory methods (PathHelperDirTrait)

| Method | Purpose |
|--------|---------|
| `isDirExists($path)` | Guard: non-empty and `is_dir()` |
| `isDirMissing($path)` | Negation of `isDirExists()` |
| `isDirWritable($path)` | Guard: exists and writable |
| `isDirReadonly($path)` | Negation of `isDirWritable()` |
| `makeDirectory($path, $secure)` | Creates dir via `wp_mkdir_p()`, optionally adds security files |
| `addSecurityFiles($path)` | Creates `.htaccess` (deny all) and `index.php` (silence) |
| `ensureDirectory($dir)` | Recursively creates directory if missing; returns `bool` |
| `ensureFileParentDirectory($filePath)` | Ensures parent dir exists for a file path |
| `isSafePath($path, $basePath)` | Path traversal guard — prevents `../` attacks |
| `isDirEmpty($path)` | True when directory has no entries |

#### File methods (PathHelperFileTrait)

| Method | Purpose |
|--------|---------|
| `isFileExists($path)` | Guard: non-empty and `file_exists()` |
| `isFileMissing($path)` | Negation wrapper |
| `isFileUnreadable($path)` | Missing or not readable |
| `deleteFile($path)` | Safe unlink with chmod retry and clearstatcache |
| `deleteDir($path)` | Recursive directory removal |
| `getRelativePath($fullPath)` | Strips base dir prefix |
| `formatBytes($bytes)` | Human-readable file size |

#### Typed path accessors via enums

All subdirectory names and database file names come from enums — no magic strings:

```php
enum PathSubdirType: string {
    case Logs      = '/logs';
    case Temp      = '/temp';
    case Snapshots = '/snapshots';
    case Exports   = '/exports';
    case Backups   = '/backups';
}

enum PathDatabaseType: string {
    case Root     = '/a-root.db';
    case Activity = '/activity.db';
    case Snapshot = '/snapshots.db';
    case Plugin   = '/riseup-asia-uploader.db';
}
```

Usage: `self::join(self::getBaseDir(), PathSubdirType::Logs->value)` — never hardcode `/logs`.

### Key design principle

Path resolution uses `wp_upload_dir()` for the base and caches the result. If `wp_upload_dir()` returns an invalid result, it falls back to `WP_CONTENT_DIR . '/uploads'`.

---

## 5.1.1 BooleanHelpers — Semantic Guard Class

BooleanHelpers is a **static utility class** providing readable boolean checks for common conditions. It eliminates negation-heavy code like `!class_exists()` and `!function_exists()` in favour of intention-revealing method names.

### Structure

```
Helpers/
├── BooleanHelpers.php           ← Shell class, composes trait
└── Traits/
    └── BooleanDomainTrait.php   ← All boolean methods
```

### Method categories

#### Environment guards

| Method | Replaces |
|--------|----------|
| `isClassExists($name)` | `class_exists($name)` |
| `isClassMissing($name)` | `!class_exists($name)` |
| `isClassUnregistered($name)` | `!class_exists($name, false)` — no autoload |
| `isFuncExists($name)` | `function_exists($name)` |
| `isFuncMissing($name)` | `!function_exists($name)` |
| `isExtensionLoaded($name)` | `extension_loaded($name)` |
| `isExtensionMissing($name)` | `!extension_loaded($name)` |
| `isConstantMissing($name)` | `!defined($name)` |

#### Data guards

| Method | Replaces |
|--------|----------|
| `isKeySet($data, $key)` | `isset($data[$key])` |
| `isKeyMissing($data, $key)` | `!isset($data[$key])` |
| `hasValue($value)` | `!empty($value)` |
| `isValueEmpty($value)` | `empty($value)` |
| `isNull($value)` | `$value === null` |
| `isAbsentFromList($needle, $haystack)` | `!in_array($needle, $haystack)` |

#### String inspection

| Method | Replaces |
|--------|----------|
| `hasSubstring($haystack, $needle)` | `str_contains()` |
| `lacksSubstring($haystack, $needle)` | `!str_contains()` |
| `hasPrefix($haystack, $prefix)` | `str_starts_with()` |
| `hasSuffix($haystack, $suffix)` | `str_ends_with()` |
| `isStringPopulated($value)` | `$value !== ''` |
| `isStringEmpty($value)` | `$value === ''` |

#### WordPress guards

| Method | Replaces |
|--------|----------|
| `isWpScheduleMissing($hook)` | `!wp_next_scheduled($hook)` |
| `isCapabilityMissing($cap)` | `!current_user_can($cap)` |
| `isDbConnected($db)` | Checks `$db !== null && $db->isReady()` |
| `isDbDisconnected($db)` | Negation of connected check |

### Usage pattern

```php
// Before (negation-heavy, easy to misread)
if (!class_exists('PDO')) { return null; }
if (!extension_loaded('pdo_sqlite')) { return null; }

// After (intention-revealing)
if (BooleanHelpers::isClassMissing('PDO')) { return null; }
if (BooleanHelpers::isExtensionMissing('pdo_sqlite')) { return null; }
```

---

## 5.1.2 InitHelpers — Bootstrap Initialization

InitHelpers handles early-boot operations where `FileLogger` may not yet be available. It uses two traits:

```
Helpers/
├── InitHelpers.php           ← Shell class, composes 2 traits
└── Traits/
    ├── InitDirTrait.php      ← Directory creation with fallback chain
    └── InitStartupTrait.php  ← Component startup tracking and diagnostics
```

### Core methods (InitHelpers shell)

| Method | Purpose |
|--------|---------|
| `initSqliteConnection($dbPath, $logger)` | Creates PDO connection with WAL mode and pragmas |
| `errorLogWithPrefix($message)` | Native `error_log()` with plugin prefix — use before FileLogger is available |
| `errorLog($e, $context)` | Log exception to native `error_log()` |
| `errorLogAndThrow($e, $context)` | Log then re-throw — for boot failures that must propagate |

### InitDirTrait — Directory creation

| Method | Purpose |
|--------|---------|
| `makeDirectory($path, $secure)` | Creates directory with deduplication cache; delegates to PathHelper when available, falls back to native `mkdir()` |
| `makeSubDirectory($baseDir, $subDir, $secure)` | Creates parent + child directory in one call |
| `resolveBaseDir()` | Resolves uploads base dir with `wp_upload_dir()` fallback |
| `addSecurityFiles($path)` | Creates `.htaccess` and `index.php` security files |

### InitStartupTrait — Component tracking

| Method | Purpose |
|--------|---------|
| `initComponent($name, $initFn)` | Wraps a boot callable with timing, error capture, and verbose logging |
| `getStartupResults()` | Returns all tracked component results |
| `getFailedStartups()` | Returns only failed components |
| `allStartupsSucceeded()` | Boolean: true if no failures |
| `logStartupSummary($logger)` | Logs total count, failures, and elapsed time |

### Boot verbose mode

Add `define('PLUGINNAME_DEBUG_BOOT', true)` to `wp-config.php` to enable per-component init logging for troubleshooting startup issues. The `isBootVerbose()` method gates these detailed log entries.

### SQLite connection pattern

```php
$pdo = InitHelpers::initSqliteConnection($dbPath, $logger);
// Internally:
// 1. Checks PDO class exists (BooleanHelpers::isClassMissing)
// 2. Checks pdo_sqlite extension loaded (BooleanHelpers::isExtensionMissing)
// 3. Creates PDO with 'sqlite:' DSN
// 4. Sets ERRMODE_EXCEPTION and FETCH_ASSOC
// 5. Applies PRAGMA journal_mode = WAL and auto_vacuum = INCREMENTAL
// 6. Returns PDO or null on failure
```

---

## 5.1.3 HttpConfigType — HTTP Request Factories

`HttpConfigType` is a backed enum that centralises HTTP timeout values and provides static factory methods for `wp_remote_get()` / `wp_remote_request()` option arrays.

```php
enum HttpConfigType: int {
    case TimeoutDefault = 30;
    case TimeoutShort   = 15;

    public static function headRedirectOptions(): array {
        return [
            'timeout'     => self::TimeoutShort->value,
            'redirection' => 0,
            'sslverify'   => true,
        ];
    }

    public static function defaultGetOptions(): array {
        return [
            'timeout'   => self::TimeoutDefault->value,
            'sslverify' => true,
        ];
    }

    public static function authenticatedOptions(string $method, string $authHeader): array {
        return [
            'method'    => strtoupper($method),
            'timeout'   => self::TimeoutDefault->value,
            'headers'   => [
                'Authorization' => $authHeader,
                'Content-Type'  => ContentTypeValueType::Json->value,
            ],
            'sslverify' => true,
        ];
    }
}
```

### Usage

```php
// HEAD request (no redirects, short timeout)
$response = wp_remote_head($url, HttpConfigType::headRedirectOptions());

// Standard GET
$response = wp_remote_get($url, HttpConfigType::defaultGetOptions());

// Authenticated API call
$response = wp_remote_request($url, HttpConfigType::authenticatedOptions('POST', $authHeader));
```

### Supporting enum — ContentTypeValueType

```php
enum ContentTypeValueType: string {
    case Json     = 'application/json';
    case JsonUtf8 = 'application/json; charset=utf-8';
}
```

All MIME types used in HTTP headers come from this enum — never hardcode `'application/json'`.

---

## 5.2 Response Envelope — The Standard API Format

Every REST endpoint returns responses in this exact envelope structure:

```json
{
  "Status": {
    "IsSuccess": true,
    "IsFailed": false,
    "Code": 200,
    "Message": "OK",
    "Timestamp": "2026-04-07T14:30:00Z"
  },
  "Attributes": {
    "RequestedAt": "/api-namespace/v1/endpoint",
    "TotalRecords": 1
  },
  "Results": [
    { "key": "value" }
  ]
}
```

### Error envelope — Debug mode ON (stack trace included)

```json
{
  "Status": { "IsSuccess": false, "IsFailed": true, "Code": 500, "Message": "Connection refused", "Timestamp": "2026-04-07T14:30:00Z" },
  "Attributes": { "RequestedAt": "/my-plugin-api/v1/activate", "TotalRecords": 0 },
  "Results": [],
  "Errors": {
    "BackendMessage": "Connection refused",
    "ExceptionType": "RuntimeException",
    "Backend": [
      "#0 ActivateHandlerTrait.php(78): PluginName\\Traits\\Activate\\ActivateHandlerTrait->executeActivation()",
      "#1 ResponseTrait.php(35): PluginName\\Traits\\Core\\ResponseTrait->safeExecute()",
      "#2 WP_REST_Server.php(1181): WP_REST_Server->dispatch()"
    ]
  }
}
```

### Error envelope — Debug mode OFF (production-safe)

```json
{
  "Status": { "IsSuccess": false, "IsFailed": true, "Code": 500, "Message": "An internal error occurred", "Timestamp": "2026-04-07T14:30:00Z" },
  "Attributes": { "RequestedAt": "/my-plugin-api/v1/activate", "TotalRecords": 0 },
  "Results": []
}
```

> **Note:** The `Errors` key is completely omitted in production to prevent leaking internal file paths, class names, or PHP version details. See Phase 4, §4.10 for full examples.

### Envelope rules

| Rule | Detail |
|------|--------|
| `IsSuccess` and `IsFailed` are always both present | They are logical inverses |
| `Timestamp` is always UTC ISO 8601 | From `DateHelper::nowUtc()` |
| `Results` is always an array | Even for single results, wrap in array |
| `Errors` key only present on failure **and** debug mode | Never include in production or on success |
| All keys are PascalCase | Defined in `ResponseKeyType` enum |
| 400-level errors always include descriptive message | Validation errors are not sensitive — always show real message |
| 500-level errors are gated by debug mode | Generic message in production, real message in debug |

---

## 5.3 EnvelopeBuilder — Fluent API

The EnvelopeBuilder uses the builder pattern with static factory methods:

### Success flow

```
EnvelopeBuilder::success('OK', 200)
    ->setRequestedAt('/namespace/v1/endpoint')
    ->setSingleResult(['key' => 'value'])
    ->toResponse();
```

### Error flow — with debug-mode gating

```
// The builder handles debug-mode gating internally
EnvelopeBuilder::error('Something failed', 500, $exception)
    ->setRequestedAt('/namespace/v1/endpoint')
    ->toResponse();
// If debug mode ON:  includes Errors.Backend with trace frames
// If debug mode OFF: omits Errors entirely, uses generic message for 500s
```

### Full method signatures

```
class EnvelopeBuilder
{
    /**
     * Create a success envelope.
     *
     * @param string $message  Status message (e.g., 'OK', 'Plugin activated')
     * @param int    $code     HTTP status code (200, 201, etc.)
     *
     * @return self Fluent builder instance
     */
    public static function success(string $message, int $code = 200): self;

    /**
     * Create an error envelope. Automatically gates stack trace by debug mode.
     *
     * @param string         $message    Error message (used as-is in debug, genericised in production for 5xx)
     * @param int            $code       HTTP status code (400, 401, 403, 404, 500)
     * @param Throwable|null $exception  Optional exception for stack trace extraction
     *
     * @return self Fluent builder instance
     */
    public static function error(string $message, int $code, ?Throwable $exception = null): self;

    /**
     * Set the requested endpoint path in Attributes.
     *
     * @param string $path  The REST route path (e.g., '/my-plugin-api/v1/status')
     *
     * @return self
     */
    public function setRequestedAt(string $path): self;

    /**
     * Wrap a single associative array in Results: [$item].
     *
     * @param array<string, mixed> $item  The single result item
     *
     * @return self
     */
    public function setSingleResult(array $item): self;

    /**
     * Set Results to the provided array directly (for list endpoints).
     *
     * @param array<int, array<string, mixed>> $items  The result items
     *
     * @return self
     */
    public function setListResult(array $items): self;

    /**
     * Manually set stack trace frames (used by safeExecute).
     * Only included in the response if debug mode is enabled.
     *
     * @param array<int, string> $frames  Formatted trace lines
     *
     * @return self
     */
    public function setStackTrace(array $frames): self;

    /**
     * Build and return the final WP_REST_Response.
     *
     * @return WP_REST_Response
     */
    public function toResponse(): WP_REST_Response;
}
```

### Fallback safety

If `EnvelopeBuilder` cannot be loaded (autoloader failure), `ResponseTrait` has an inline fallback that builds the same envelope structure manually. This ensures the plugin never returns a bare PHP error.

---

## 5.4 Integration Checklist — Adding a New Feature

When adding a new feature endpoint to the plugin, follow this exact sequence:

### Step 1: Define enums

| What to add | Where |
|-------------|-------|
| Endpoint path | New case in `EndpointType` |
| New response keys | New cases in `ResponseKeyType` |
| New capabilities (if any) | New case in `CapabilityType` |

### Step 2: Create the handler trait

1. Create a new file: `Traits/{FeatureDomain}/{FeatureName}Trait.php`
2. Follow the trait anatomy from Phase 3, §3.3
3. The public handler method wraps logic in `$this->safeExecute()`
4. Use `EnvelopeBuilder` for all responses
5. Use `$this->fileLogger` for all logging
6. Use enum values for all string literals

### Step 3: Register the route

1. Add a new registration method in `RouteRegistrationTrait` (or add to an existing group)
2. Wire it using the `$safeRegister` closure pattern
3. Use `EndpointType::NewEndpoint->route()` for the path
4. Use `HttpMethodType` for the method
5. Point to the correct permission callback

### Step 4: Compose in Plugin.php

1. Add `use PluginName\Traits\{FeatureDomain}\{FeatureName}Trait;` import
2. Add `use {FeatureName}Trait;` inside the class body
3. If a new route group was created, add it to the `$groups` array in `registerRoutes()`

### Step 5: Bump version

Update `PluginConfigType::Version` case value.

---

## 5.5 Database — Split DB Concept

Plugins that need data persistence use SQLite stored in `wp-content/uploads/{plugin-slug}/` rather than WordPress's MySQL database. This provides:

| Benefit | Detail |
|---------|--------|
| Isolation | Plugin data is completely separate from WordPress tables |
| Portability | Database file can be backed up, moved, or deleted independently |
| No migration conflicts | No interference with WordPress core or other plugin migrations |
| Schema versioning | Track migration versions inside the SQLite database itself |

### Database location

```
wp-content/uploads/{plugin-slug}/{plugin-slug}.db
```

### Multi-DB expansion (Split DB)

Larger plugins may split data across multiple SQLite files. Each database is identified by a `PathDatabaseType` enum case:

```
wp-content/uploads/{plugin-slug}/
├── {plugin-slug}.db   ← Main plugin data (PathDatabaseType::Plugin)
├── a-root.db          ← Cross-plugin registration (PathDatabaseType::Root)
├── activity.db        ← Transaction/activity logs (PathDatabaseType::Activity)
└── snapshots.db       ← Snapshot exports (PathDatabaseType::Snapshot)
```

Each database path is resolved via `PathHelper::join(PathHelper::getBaseDir(), PathDatabaseType::Plugin->value)` — never hardcoded.

> **Reference:** For the Go backend implementation of Split DB with child database connection pooling, see `backend/internal/database/Database.go`.

### Schema versioning

Store a `schema_version` value in the database (either a dedicated table or SQLite `user_version` pragma). On plugin activation or init, compare the stored version to the expected version and run any pending migrations.

---

## 5.5.1 Database Class — Trait Decomposition

The Database class follows the same shell + trait pattern as helpers. The class holds shared state (`$pdo`, `$dbPath`, `$fileLogger`) while traits provide domain-specific logic:

```
Database/
├── Database.php                        ← Shell class (singleton, shared state)
└── Traits/
    ├── DatabaseConnectionTrait.php     ← init(), getPdo(), isReady(), migration orchestrator
    ├── DatabaseConvenienceTrait.php    ← queryAll(), querySingle(), insert(), update(), delete()
    ├── DatabaseQueryTrait.php          ← Shell: composes query sub-traits
    │   ├── DatabaseQueryLogTrait.php   ← Transaction log queries
    │   └── DatabaseQuerySearchTrait.php ← Search/filter queries
    ├── DatabaseMigrationsEarlyTrait.php ← Shell: composes v1-v5 migration traits
    │   ├── DatabaseMigrationsV1V3Trait.php
    │   └── DatabaseMigrationsV4V5Trait.php
    ├── DatabaseMigrationsLateTrait.php  ← Shell: composes v6+ migration traits
    │   ├── DatabaseMigrationsV6V8Trait.php
    │   ├── DatabaseMigrationsV9V11Trait.php
    │   └── DatabaseMigrationsV12Trait.php ... V22Trait.php
    └── OrmQueryTrait.php, OrmMutationTrait.php, OrmWhereTrait.php  ← Micro-ORM
```

### Shell class pattern

```php
class Database {
    use DatabaseConnectionTrait;
    use DatabaseConvenienceTrait;
    use DatabaseMigrationsEarlyTrait;
    use DatabaseMigrationsLateTrait;
    use DatabaseQueryTrait;

    private ?PDO $pdo = null;
    private string $dbPath = '';
    private FileLogger $fileLogger;
    private static ?self $instance = null;
    private bool $isInitAttempted = false;

    public static function getInstance(): self { /* singleton */ }
    private function __construct() { $this->fileLogger = FileLogger::getInstance(); }
}
```

### ConnectionTrait responsibilities

| Method | Purpose |
|--------|---------|
| `init()` | Lazy initialization — only runs once via `$isInitAttempted` guard |
| `getDatabasePath()` | Resolves path via `PathHelper::getDbPath()`, ensures directory exists |
| `initDatabase()` | Creates PDO, configures ORM, runs migrations, logs timing |
| `createTables()` | Migration orchestrator — calls all `migrateV*()` methods in sequence |
| `ensureSchemaVersionTable()` | Creates `schema_version` table if missing |
| `getCurrentSchemaVersion()` | `SELECT MAX(version)` from schema_version |
| `recordMigration($version)` | Inserts version + UTC timestamp after successful migration |
| `getPdo()` | Returns PDO (triggers lazy init if needed) |
| `isReady()` | Boolean: is PDO connected? |

### ConvenienceTrait — Thin PDO wrappers

| Method | Purpose |
|--------|---------|
| `queryAll($sql, $params)` | SELECT → all rows as assoc arrays |
| `querySingle($sql, $params)` | SELECT → single row |
| `insert($table, $data)` | INSERT from associative array |
| `update($table, $data, $where, $whereParams)` | UPDATE with WHERE clause |
| `delete($table, $where, $whereParams)` | DELETE with WHERE clause |
| `execute($sql, $params)` | Generic exec (DDL, DML) |
| `execIfColumnMissing($table, $column, $sql)` | Safe ADD COLUMN for idempotent migrations |

### Migration trait pattern

Each migration method follows this guard pattern:

```php
private function migrateV4SourceMachine(int $current): void {
    if ($current >= 4) { return; }

    $this->pdo->exec("ALTER TABLE transactions ADD COLUMN source_machine TEXT DEFAULT ''");
    $this->recordMigration(4);
}
```

When migrations grow beyond 200 lines, split into sub-traits grouped by version range and compose them via shell traits (`DatabaseMigrationsEarlyTrait`, `DatabaseMigrationsLateTrait`).

---

## 5.6 Notification Patterns

When the plugin needs to send notifications (email, admin notices), follow these patterns:

| Pattern | Implementation |
|---------|---------------|
| Email sending | Delegate to `wp_mail()` with structured HTML templates |
| Log email endpoint | A dedicated REST endpoint that emails log contents to a configured recipient |
| Admin notices | Only show on plugin's own admin pages, never globally |
| Error notifications | Log the error; optionally email if severity is critical |

---

## 5.7 Security Checklist

| Requirement | Implementation |
|-------------|---------------|
| Authentication | WordPress Application Passwords via Basic Auth |
| Capability checks | Every endpoint has a `permission_callback` |
| Input sanitisation | Use `sanitize_text_field()`, `absint()`, `wp_kses()` |
| Output escaping | Use `esc_html()`, `esc_attr()`, `esc_url()` in admin pages |
| Nonce verification | All admin AJAX actions verify nonces via enum-defined values |
| Rate limiting | Track request counts per IP in transients or custom table |
| ABSPATH guard | Every PHP file checks `defined('ABSPATH')` |

---

## 5.8 Summary — The Complete Pattern

```
Request → WordPress REST API
  → RouteRegistrationTrait (resolves endpoint)
  → AuthTrait (validates credentials + capabilities)
  → Handler Trait (public method)
    → safeExecute() (error boundary from ResponseTrait)
      → Private method (business logic)
        → FileLogger (structured logging)
        → EnvelopeBuilder (response construction)
        → Enums (all string/int constants)
        → Helpers (stateless utilities)
      ← WP_REST_Response (envelope format)
    ← Throwable caught → errorResponse() → WP_REST_Response
  ← JSON response to client
```

Every layer has a single responsibility. Every string is an enum value. Every error is caught, logged with a stack trace, and returned in a structured format. The pattern is identical for every endpoint.
