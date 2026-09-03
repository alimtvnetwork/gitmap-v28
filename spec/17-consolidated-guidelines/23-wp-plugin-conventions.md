# Consolidated: WordPress Plugin Development — Complete Reference

**Version:** 3.3.0
**Updated:** 2026-04-16
**Source Module:** [`02-spec/18-wp-plugin-how-to/`](../18-wp-plugin-how-to/01-index.md)

---

## Purpose

This is the **standalone consolidated reference** for WordPress plugin development following the Gold Standard architecture. An AI reading only this file must be able to build a compliant plugin.

---

## Prerequisites

| Requirement | Value |
|-------------|-------|
| PHP | 8.1+ (backed enums, readonly, fibers) |
| WordPress | 5.6+ |
| Array syntax | `[]` only — never `array()` |
| Exception catching | Always `Throwable`, never `Exception` |
| Magic strings | **Forbidden** — use backed enums |

---

## Canonical Folder Structure

```
plugin-slug/
├── plugin-slug.php              ← Bootstrap (non-namespaced, 3 tasks only)
├── includes/
│   ├── Autoloader.php           ← PSR-4 autoloader (non-namespaced)
│   ├── Core/
│   │   └── Plugin.php           ← Singleton, composes traits, wires hooks
│   ├── Enums/                   ← All backed enums
│   │   ├── EndpointType.php
│   │   ├── HttpMethodType.php
│   │   ├── PluginConfigType.php
│   │   ├── ResponseKeyType.php
│   │   ├── FilterKeyType.php
│   │   ├── RequestFieldType.php
│   │   ├── PaginationConfigType.php
│   │   ├── PhpNativeType.php
│   │   └── [Domain]Type.php
│   ├── Helpers/
│   │   ├── EnvelopeBuilder.php  ← Fluent response builder
│   │   ├── PathHelper.php       ← File system path resolution
│   │   ├── BooleanHelpers.php   ← Semantic boolean guards
│   │   ├── DateHelper.php       ← Timestamp formatting
│   │   └── InitHelpers.php      ← Bootstrap initialization
│   ├── Logging/
│   │   └── FileLogger.php       ← Structured file logger (singleton)
│   └── Traits/
│       ├── Auth/                ← Permission checks
│       ├── Core/                ← ResponseTrait, StatusHandlerTrait, PingHandlerTrait
│       ├── Route/               ← RouteRegistrationTrait
│       ├── Upload/              ← FileUploadTrait
│       └── [Feature]/           ← One subfolder per feature domain
├── templates/                   ← Admin page templates (200-line max)
├── assets/                      ← CSS/JS/images
├── data/                        ← JSON data files (colors, endpoints, openapi)
│   └── seeds/                   ← Seed JSON files + manifest
└── spec/                        ← Documentation
```

---

## Bootstrap File (3 Tasks Only)

The main `plugin-slug.php` file does exactly 3 things:

```php
<?php
/**
 * Plugin Name: My Plugin
 * Version: 1.0.0
 * Author: Md. Alim Ul Karim
 */

if (!defined('ABSPATH')) { exit; }

define('MY_PLUGIN_DEBUG', false);  // Debug mode gate

require_once __DIR__ . '/includes/Autoloader.php';

PluginSlug\Core\Plugin::getInstance();
```

Nothing else. No function definitions, no hook registrations, no constants beyond debug mode.

---

## Singleton Plugin Core

`Plugin.php` uses `getInstance()` pattern, composes traits via `use`, and wires all WordPress hooks in `__construct()`.

```php
class Plugin {
    use ResponseTrait;
    use RouteRegistrationTrait;
    use AuthTrait;
    use StatusHandlerTrait;
    use PingHandlerTrait;
    use TypeCheckerTrait;
    use FileUploadTrait;
    // ... feature traits

    private static ?self $instance = null;
    private FileLogger $fileLogger;

    public static function getInstance(): self {
        if (self::$instance === null) {
            self::$instance = new self();
        }
        return self::$instance;
    }

    private function __construct() {
        $this->fileLogger = FileLogger::getInstance();
        add_action('rest_api_init', [$this, 'registerRoutes']);
        add_action('admin_menu', [$this, 'registerAdminPages']);
        // ... other hooks
    }
}
```

---

## Traits Over Inheritance

All feature logic lives in traits, organized by domain folder. Composition over inheritance — always.

### Trait vs Helper vs Enum Decision Table

| Question | Answer → Use |
|----------|-------------|
| Does it handle a REST endpoint? | Trait |
| Does it need `$this->fileLogger`? | Trait |
| Is it stateless utility logic? | Helper (static class) |
| Is it a set of named constants? | Enum |
| Does it represent configuration? | Enum |

### Trait Anatomy

Every handler trait follows this structure:

```php
trait SomeHandlerTrait
{
    // Public: route handler → wraps in safeExecute
    public function handleSomething(WP_REST_Request $request): WP_REST_Response
    {
        return $this->safeExecute(
            fn() => $this->executeSomething($request),
            'something',
        );
    }

    // Private: business logic
    private function executeSomething(WP_REST_Request $request): WP_REST_Response
    {
        // 1. Extract and validate input
        // 2. Execute business logic
        // 3. Log result
        // 4. Return envelope response

        return EnvelopeBuilder::success('Operation complete')
            ->setRequestedAt($request->get_route())
            ->setSingleResult($data)
            ->toResponse();
    }
}
```

---

## Backed Enums (PHP 8.1+)

All string constants use backed enums with:

- `Type` suffix (e.g., `ResponseKeyType`, `StatusType`)
- `isEqual()` method for comparison
- PascalCase values

```php
enum ResponseKeyType: string {
    case Status = 'Status';
    case IsSuccess = 'IsSuccess';
    case Results = 'Results';
    case ErrorCode = 'ErrorCode';
    case Author = 'Author';
    case Company = 'Company';
    case Version = 'Version';

    public function isEqual(self $other): bool {
        return $this->value === $other->value;
    }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Standard Enum Categories

| Category | Examples | Purpose |
|----------|----------|---------|
| Config | `PluginConfigType` | Plugin slug, version, API namespace, debug constants |
| Route | `EndpointType`, `HttpMethodType` | Route paths and HTTP methods |
| Response | `ResponseKeyType` | All PascalCase response JSON keys |
| Validation | `PhpNativeType` | Type checking via `TypeCheckerTrait` |
| Domain | `ActionType`, `StatusType` | Business-specific enums |

---

## Error Handling

- Always catch `Throwable`, never `Exception`
- Use `FileLogger` for structured logging (singleton)
- Never use `error_log()` for diagnostics (only as Tier 1 fallback before autoloader)
- Include context (file path, operation) in every error message
- `safeExecute()` wraps all endpoint handlers — catches Throwable, logs, returns structured envelope

### Debug Mode Gate

| Feature | Debug ON | Debug OFF |
|---------|----------|-----------|
| Stack traces in API responses | ✅ Full frames | ❌ Omitted |
| Verbose log entries | ✅ `debug()` writes to `info.log` | ❌ Silently skipped |
| Error response detail | ✅ Full exception message | ⚠️ Generic message |

---

## Input Validation

### Philosophy

| Principle | Rule |
|-----------|------|
| Fail fast | Validate all inputs at top of handler, before business logic |
| Guard clauses | Each validation is a standalone `if` → `return error` |
| Positive booleans | Extract every check into `$has…` / `$is…` variable |
| No silent defaults | Never substitute a default for a missing required field |

### Standard Flow

```php
private function executeCreateWidget(WP_REST_Request $request): WP_REST_Response
{
    $body = $request->get_json_params();
    $hasBody = ($body !== null && $this->isArray($body));
    if (!$hasBody) {
        return $this->validationError('Request body must be a JSON object', $request);
    }

    $name = $body['name'] ?? null;
    $hasName = ($name !== null && $this->isString($name));
    if (!$hasName) {
        return $this->validationError('Missing required field: name', $request);
    }

    $sanitisedName = sanitize_text_field($name);
    // ... business logic
}
```

---

## File Upload Flow

### FileUploadTrait

Complete upload handling via `Traits/Upload/FileUploadTrait.php`:

```php
trait FileUploadTrait
{
    public function handleFileUpload(WP_REST_Request $request): WP_REST_Response
    {
        return $this->safeExecute(
            fn() => $this->executeFileUpload($request),
            'file-upload',
        );
    }

    private function executeFileUpload(WP_REST_Request $request): WP_REST_Response
    {
        // ── 1. Check for uploaded files ──
        $files = $request->get_file_params();
        $hasFiles = (!empty($files) && isset($files['file']));

        if (!$hasFiles) {
            return $this->validationError('No file uploaded', $request);
        }

        $file = $files['file'];

        // ── 2. Check for upload errors ──
        $uploadError = $file['error'] ?? UPLOAD_ERR_NO_FILE;
        $hasUploadError = ($uploadError !== UPLOAD_ERR_OK);

        if ($hasUploadError) {
            $errorMessage = $this->resolveUploadError($uploadError);
            return $this->validationError($errorMessage, $request);
        }

        // ── 3. Validate file type ──
        $fileName = sanitize_file_name($file['name']);
        $extension = strtolower(pathinfo($fileName, PATHINFO_EXTENSION));
        $allowedExtensions = ['zip', 'json', 'csv'];
        $isAllowedType = in_array($extension, $allowedExtensions, true);

        if (!$isAllowedType) {
            $allowed = implode(', ', $allowedExtensions);
            return $this->validationError(
                "File type '.{$extension}' is not allowed. Allowed: {$allowed}",
                $request,
            );
        }

        // ── 4. Validate file size (max 50MB) ──
        $maxSizeBytes = 50 * 1024 * 1024;
        $fileSize = $file['size'] ?? 0;
        $isTooBig = ($fileSize > $maxSizeBytes);

        if ($isTooBig) {
            $maxMb = $maxSizeBytes / 1048576;
            return $this->validationError(
                "File exceeds maximum size of {$maxMb}MB",
                $request,
            );
        }

        // ── 5. Move file to plugin directory ──
        $tempPath = $file['tmp_name'];
        $targetDir = PathHelper::getTempDir();
        PathHelper::ensureDirectory($targetDir);

        $targetPath = $targetDir . '/' . $fileName;
        $isMoved = move_uploaded_file($tempPath, $targetPath);

        if (!$isMoved) {
            return EnvelopeBuilder::error('Failed to save uploaded file', 500)
                ->setRequestedAt($request->get_route())
                ->toResponse();
        }

        $this->fileLogger->info('File uploaded', [
            'fileName' => $fileName,
            'size'     => $fileSize,
        ]);

        return EnvelopeBuilder::success('File uploaded successfully')
            ->setRequestedAt($request->get_route())
            ->setSingleResult([
                'fileName' => $fileName,
                'size'     => $fileSize,
                'path'     => $targetPath,
            ])
            ->toResponse();
    }

    private function resolveUploadError(int $errorCode): string
    {
        return match ($errorCode) {
            UPLOAD_ERR_INI_SIZE   => 'File exceeds server upload limit',
            UPLOAD_ERR_FORM_SIZE  => 'File exceeds form upload limit',
            UPLOAD_ERR_PARTIAL    => 'File was only partially uploaded',
            UPLOAD_ERR_NO_FILE    => 'No file was uploaded',
            UPLOAD_ERR_NO_TMP_DIR => 'Server misconfiguration: missing temp directory',
            UPLOAD_ERR_CANT_WRITE => 'Server failed to write file to disk',
            UPLOAD_ERR_EXTENSION  => 'Upload blocked by server extension',
            default               => "Unknown upload error (code: {$errorCode})",
        };
    }
}
```

### Upload Route Registration

```php
$safeRegister(
    EndpointType::FileUpload->route(),
    [
        'methods'             => HttpMethodType::Post->value,
        'callback'            => [$this, 'handleFileUpload'],
        'permission_callback' => [$this, 'checkPluginPermission'],
    ],
);
```

### Upload Edge Cases

| Scenario | Handling |
|----------|----------|
| PHP `upload_max_filesize` too low | Check `ini_get('upload_max_filesize')` and warn in health check |
| `post_max_size` smaller than file | Returns empty `$_FILES` — check for this before validation |
| Symlink attack on temp file | Use `move_uploaded_file()` (not `rename()`) — it validates the temp path |
| File name collision | Prepend timestamp or UUID: `time() . '_' . $fileName` |
| Binary file disguised as allowed type | Check MIME via `wp_check_filetype()` in addition to extension |

### Plugin ZIP Upload Flow

For plugins that manage other plugin installations (e.g., a plugin manager):

1. **Receive** — `POST /upload-active` with ZIP file via `multipart/form-data`
2. **Validate** — Extension must be `.zip`, size ≤ 50MB, MIME check
3. **Extract** — Use `unzip_file()` (WP built-in) to temp directory
4. **Verify** — Check extracted folder contains a valid plugin header file
5. **Install** — Move to `WP_PLUGIN_DIR/{plugin-slug}/`
6. **Activate** — Call `activate_plugin($pluginPath)` and check for errors
7. **Log** — Record action in `FileLogger` with plugin slug, version, source
8. **Respond** — Envelope response with installed plugin metadata

---

## REST API Conventions

### Namespace Convention

```
{plugin-slug}/v{major}
```

Constructed via `PluginConfigType::apiFullNamespace()`. Version is major only (`v1`, `v2`). Never hardcode.

### Route Naming Rules

| Rule | ✅ Correct | ❌ Wrong |
|------|-----------|---------|
| Plural nouns for collections | `plugins`, `agents` | `plugin`, `agent` |
| Kebab-case for multi-word paths | `upload-active` | `uploadActive` |
| Action verbs as path segments | `plugins/enable` | `enablePlugin` |
| No trailing slashes | `/plugins` | `/plugins/` |

### EndpointType Enum

Every route path is defined as an enum case — no string literals in registration:

```php
enum EndpointType: string {
    case Status  = 'status';
    case Ping    = 'ping';
    case Plugins = 'plugins';
    case Logs    = 'logs';
    case FileUpload = 'upload';

    public function route(): string { return '/' . $this->value; }
}
```

### Route Registration Pattern

```php
$safeRegister(EndpointType::Ping->route(), [
    'methods'             => HttpMethodType::Get->value,
    'callback'            => [$this, 'handlePing'],
    'permission_callback' => PluginConfigType::pingPermissionCallback($this),
], 'system');
```

Fault-tolerant `$safeRegister` closure — logs failures without stopping other routes.

### Pagination

- Query params: `limit` (default 50, max 500), `offset`, `page`
- `PaginationConfigType` enum for defaults
- Response includes `TotalRecords`, `Limit`, `Offset`, `Page`, `TotalPages` in `Attributes`

### Filtering

- `FilterKeyType` enum for all filter parameter names
- camelCase for query filter keys: `triggeredBy`, `uploadSource`
- ISO 8601 for date filters

### Key Naming by Layer

| Layer | Convention | Example |
|-------|-----------|---------|
| Request body fields | snake_case | `plugin_zip`, `upload_source` |
| Query parameters | camelCase | `triggeredBy`, `uploadSource` |
| Response keys | PascalCase | `PluginSlug`, `TotalRecords` |
| URL path segments | kebab-case | `upload-active`, `sync-manifest` |

---

## Ping Endpoint (Mandatory)

Every plugin **must** expose a `GET /ping` endpoint that returns:

| Field | Type | Source |
|-------|------|--------|
| `Author` | string | `PluginConfigType::Author->value` |
| `Company` | string | `PluginConfigType::Company->value` |
| `Version` | string | `PluginConfigType::Version->value` |

**Authorization mode** configurable via `PluginConfigType::IsPingAuthorized`:

- `'true'` → requires authentication (`checkPingPermission`)
- `'false'` → publicly accessible (`__return_true`)

Response message is `"pong"`.

---

## Deployment & Versioning

### Semantic Versioning

All plugins follow SemVer 2.0.0. Version is declared in exactly two places (must always match):

| Location | Read by |
|----------|---------|
| Main plugin file header (`* Version: X.Y.Z`) | WordPress core |
| `PluginConfigType::Version` enum case | All plugin code |

**Rule:** Every code change requires at least a minor version bump. The `.release` folder is exempt.

### Version Bump Checklist

1. Update main plugin file header → `Version: X.Y.Z`
2. Update `PluginConfigType::Version` → `case Version = 'X.Y.Z'`
3. Update `changelog.md` → add entry under `## [X.Y.Z] - YYYY-MM-DD`
4. If `composer.json` has a version field → update it too

### ZIP Packaging

**Distribution structure** (what ships):

```
plugin-slug/                      ← ZIP root
├── plugin-slug.php               ← Main plugin file (required)
├── uninstall.php                 ← Clean removal hook
├── settings.json                 ← Default configuration
├── assets/                       ← Admin CSS/JS/images
├── data/seeds/                   ← Seed JSON files + manifest
├── includes/                     ← All PHP source (PSR-4)
├── templates/                    ← PHP view templates
└── vendor/                       ← Composer autoloader (production only)
```

**Excluded from ZIP:** `.git/`, `.github/`, `.ai-instructions`, `tests/`, `phpunit.xml`, `phpstan.neon`, `composer.lock`, `spec/`, `*.log`, `node_modules/`

### ZIP Integrity Requirements

| Check | Rule |
|-------|------|
| Main plugin file exists | `plugin-slug/plugin-slug.php` must be at ZIP root level |
| No nested folders | ZIP must not contain `plugin-slug/plugin-slug/` (double-nesting) |
| vendor/ present | Autoloader must be included — plugin won't boot without it |
| No dev dependencies | `vendor/phpunit/` must NOT appear in the ZIP |
| File permissions | PHP files: 644, directories: 755 |

### Self-Update

- Self-hosted update server via two WordPress filters (`site_transient_update_plugins`, `plugins_api`)
- Update JSON endpoint serves version, download URL, tested/requires info
- `SelfUpdateStatusType` enum tracks update stages
- Rollback capability via version pinning

---

## Data File Patterns

### endpoints.json

Every plugin maintains `data/endpoints.json` documenting all REST routes:

```json
{
  "namespace": "plugin-slug/v1",
  "version": "2.0.0",
  "endpoints": [
    {
      "path": "ping",
      "methods": ["GET"],
      "category": "system",
      "description": "Health check — returns author, company, version",
      "auth": true
    }
  ]
}
```

### colors.json

Design tokens in `data/colors.json`, accessed via `ColorGroupType` enum with static caching.

---

## Admin UI Patterns

- Templates in `templates/` folder, max 200 lines each
- Orchestrator pattern: main template composes partials
- Use WordPress admin styles, extend with scoped custom CSS
- Data files in `data/` folder (JSON format)
- Register via `add_action('admin_menu', [$this, 'registerAdminPages'])`
- Always check `current_user_can(CapabilityType::ManageOptions->value)` before rendering

---

## Micro-ORM & Root Database

### TypedQuery

Fluent query builder with Go-style result wrappers:

| Method | Returns | Purpose |
|--------|---------|---------|
| `queryOne()` | `DbResult<T>` | Single row with typed mapper |
| `queryMany()` | `DbResultSet<T>` | Multiple rows with typed mapper |
| `execute()` | `DbExecResult` | INSERT/UPDATE/DELETE |

### Root Database (RootDb)

Cross-plugin shared SQLite database for:

- Plugin manifest (installed plugins, versions)
- Shared configuration
- Cross-plugin communication

---

## Testing Patterns

- Unit tests for all helpers and utility classes
- Integration tests for WordPress hook wiring
- Mock WordPress functions where needed
- Test naming: `testMethodName_Scenario_ExpectedResult`

---

## Forbidden Patterns

| Pattern | Required Alternative |
|---------|---------------------|
| `array()` syntax | `[]` |
| `catch (Exception $e)` | `catch (Throwable $e)` |
| `error_log()` for diagnostics | `FileLogger` |
| Magic string constants | Backed enums |
| String literals in route registration | `EndpointType` enum |
| String literals in HTTP methods | `HttpMethodType` enum |
| Hardcoded namespace string | `PluginConfigType::apiFullNamespace()` |
| `__return_true` on authenticated endpoints | Proper `permission_callback` |
| Nested `if` blocks | Guard clauses with early returns |
| Boolean flags as parameters | Separate named methods |
| `rename()` for uploaded files | `move_uploaded_file()` |
| Inline HTML in PHP classes | Template files in `templates/` |

---

## Cross-References

| Topic | Source Location |
|-------|---------------|
| Quick Start | `02-spec/18-wp-plugin-how-to/02-quick-start.md` |
| Foundation & Architecture | `02-spec/18-wp-plugin-how-to/03-foundation-and-architecture.md` |
| Enums & Coding Style | `02-spec/18-wp-plugin-how-to/02-enums-and-coding-style/` |
| Traits & Composition | `02-spec/18-wp-plugin-how-to/04-traits-and-composition.md` |
| Logging & Error Handling | `02-spec/18-wp-plugin-how-to/05-logging-and-error-handling.md` |
| Input Validation | `02-spec/18-wp-plugin-how-to/07-input-validation-patterns.md` |
| WordPress Integration | `02-spec/18-wp-plugin-how-to/09-wordpress-integration-patterns.md` |
| Deployment Patterns | `02-spec/18-wp-plugin-how-to/11-deployment-patterns.md` |
| REST API Conventions | `02-spec/18-wp-plugin-how-to/15-rest-api-conventions.md` |
| Settings Architecture | `02-spec/18-wp-plugin-how-to/16-settings-architecture.md` |
| Micro ORM & Root DB | `02-spec/18-wp-plugin-how-to/20-micro-orm-and-root-db.md` |
| Ping Endpoint | `02-spec/18-wp-plugin-how-to/22-ping-endpoint.md` |
| End-to-End Walkthrough | `02-spec/18-wp-plugin-how-to/21-end-to-end-walkthrough.md` |

---

*Consolidated WordPress plugin conventions — v3.3.0 — 2026-04-16*
