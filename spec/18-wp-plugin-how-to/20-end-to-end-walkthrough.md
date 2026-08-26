# Phase 20 — End-to-End Walkthrough: Building "Task Tracker" from Scratch

> **Purpose:** Walk through building a complete, minimal WordPress plugin — **Task Tracker** — from an empty folder to a deployable plugin, applying every phase of the Gold Standard specification. Each step references the authoritative phase so you can cross-check rules as you go.  
> **Audience:** AI code generators receiving this spec for the first time and human developers seeking a concrete example.  
> **Plugin identity:** `task-tracker` / `TaskTracker` / `TASK_TRACKER`

---

## 20.1 What We're Building

A WordPress REST API plugin that lets an authenticated user:

1. **Create** a task (title, optional priority)
2. **List** all tasks (with pagination)
3. **Complete** a task (set status to `done`)
4. **View admin settings** (enable/disable task notifications)

This covers: bootstrap, autoloader, enums, traits, logging, validation, response envelope, REST endpoints, SQLite database, admin settings page, and testing.

### Plugin identity table

| Key | Value | Phase |
|-----|-------|-------|
| Slug | `task-tracker` | [Phase 1, §1.8](01-foundation-and-architecture.md) |
| ShortName | `TaskTracker` | — |
| Namespace | `TaskTracker\` | — |
| API namespace | `task-tracker-api/v1` | [Phase 14, §14.1](14-rest-api-conventions.md) |
| Debug constant | `TASK_TRACKER_DEBUG` | [Phase 4, §4.2](04-logging-and-error-handling.md) |
| Log prefix | `[TaskTracker]` | — |

---

## 20.2 Step 1 — Folder Structure

> **Phase 1, §1.2** — Canonical folder structure

```
task-tracker/
├── task-tracker.php                    ← Bootstrap (Step 2)
├── uninstall.php                       ← Cleanup on delete (Step 14)
├── includes/
│   ├── Autoloader.php                  ← PSR-4 loader (Step 3)
│   ├── Core/
│   │   ├── Plugin.php                  ← Singleton composition root (Step 7)
│   │   ├── Activator.php              ← Activation hook handler
│   │   └── Deactivator.php            ← Deactivation hook handler
│   ├── Enums/
│   │   ├── PluginConfigType.php        ← Identity enum (Step 4)
│   │   ├── EndpointType.php            ← Route paths (Step 4)
│   │   ├── HttpMethodType.php          ← GET/POST/PUT/DELETE
│   │   ├── HttpStatusType.php          ← Status codes
│   │   ├── HookType.php                ← WordPress hook names
│   │   ├── TaskStatusType.php          ← pending/done (Step 4)
│   │   ├── OptionNameType.php          ← wp_options keys
│   │   ├── ResponseKeyType.php         ← Envelope keys
│   │   ├── PhpNativeType.php           ← Type-checker backing
│   │   ├── CapabilityType.php          ← WP capabilities
│   │   └── LogLevelType.php            ← Log levels
│   ├── Helpers/
│   │   ├── EnvelopeBuilder.php         ← Response envelope (Step 6)
│   │   ├── DateHelper.php              ← Timestamp formatting
│   │   ├── PathHelper.php              ← File path resolution
│   │   └── ErrorLogHelper.php          ← Tier 1 logging helper
│   ├── Logging/
│   │   └── FileLogger.php              ← Tier 2 logger (Step 5)
│   ├── Database/
│   │   └── DatabaseMigrationsTrait.php ← SQLite schema (Step 9)
│   └── Traits/
│       ├── Auth/
│       │   └── AuthTrait.php           ← Permission checks (Step 6)
│       ├── Core/
│       │   ├── ResponseTrait.php       ← safeExecute + envelope (Step 6)
│       │   └── TypeCheckerTrait.php    ← Safe type checking
│       ├── Route/
│       │   └── RouteRegistrationTrait.php  ← Route wiring (Step 8)
│       └── Task/
│           ├── TaskCreateTrait.php     ← POST /tasks (Step 10)
│           ├── TaskListTrait.php       ← GET /tasks (Step 10)
│           └── TaskCompleteTrait.php   ← POST /tasks/complete (Step 10)
├── templates/
│   └── settings.php                    ← Admin settings page (Step 12)
├── data/
│   └── seeds/
│       ├── manifest.json               ← Seed registry
│       └── default-settings.json       ← Default settings seed
├── assets/
│   └── css/
│       └── admin.css                   ← Admin page styles
└── tests/
    ├── bootstrap.php
    └── Unit/
        └── Enums/
            └── TaskStatusTypeTest.php  ← Example test (Step 13)
```

---

## 20.3 Step 2 — Bootstrap File

> **Phase 1, §1.3** + **Phase 7, §7.1** — Three tasks only: header, autoloader, singleton.

**File: `task-tracker.php`**

```php
<?php
/**
 * Plugin Name:       Task Tracker
 * Plugin URI:        https://example.com/task-tracker
 * Description:       A minimal task management plugin built with the Gold Standard architecture.
 * Version:           1.0.0
 * Requires at least: 5.6
 * Requires PHP:      8.1
 * Author:            Your Name
 * Author URI:        https://example.com
 * License:           GPL-2.0-or-later
 * License URI:       https://www.gnu.org/licenses/gpl-2.0.html
 * Text Domain:       task-tracker
 */

if (!defined('ABSPATH')) {
    exit;
}

/** Enable debug mode — exposes stack traces in API responses. Set false in production. */
define('TASK_TRACKER_DEBUG', false);

// ── Load autoloader ──
$autoloaderPath = __DIR__ . '/includes/Autoloader.php';
$hasAutoloader = file_exists($autoloaderPath);

if (!$hasAutoloader) {
    error_log('[TaskTracker] FATAL: Autoloader.php not found at ' . $autoloaderPath);

    return;
}

require_once $autoloaderPath;

// ── Boot plugin ──
// Uses 'plugins_loaded' because we need admin UI + REST API
add_action('plugins_loaded', function (): void {
    try {
        \TaskTracker\Core\Plugin::getInstance();
    } catch (\Throwable $e) {
        error_log('[TaskTracker] Boot failed: ' . $e->getMessage() . "\n" . $e->getTraceAsString());
    }
});

// ── Activation hook ──
register_activation_hook(__FILE__, function (): void {
    try {
        \TaskTracker\Core\Activator::activate();
    } catch (\Throwable $e) {
        error_log('[TaskTracker] Activation failed: ' . $e->getMessage() . "\n" . $e->getTraceAsString());
    }
});

// ── Deactivation hook ──
register_deactivation_hook(__FILE__, function (): void {
    try {
        \TaskTracker\Core\Deactivator::deactivate();
    } catch (\Throwable $e) {
        error_log('[TaskTracker] Deactivation failed: ' . $e->getMessage() . "\n" . $e->getTraceAsString());
    }
});
```

### Checklist — What to verify

| ✅ Rule | Source |
|---------|--------|
| ABSPATH guard at top | Phase 1, §1.6 |
| `return` not `exit` on autoloader failure | Phase 7, §7.1 |
| Every hook body in try-catch `Throwable` | Phase 4, §4.8 Rule 1 |
| Stack trace appended in every catch | Phase 4, §4.8 Rule 2 |
| Debug constant uses `PLUGIN_NAME_DEBUG` pattern | Phase 4, §4.2 |

---

## 20.4 Step 3 — Autoloader

> **Phase 1, §1.4** + **Phase 7, §7.2** — PSR-4 mapping, diagnostic logging, self-register.

**File: `includes/Autoloader.php`**

Copy the reference autoloader from Phase 7, §7.2 verbatim, then search-replace:

| Find | Replace |
|------|---------|
| `PluginNameAutoloader` | `TaskTrackerAutoloader` |
| `PluginName\\` | `TaskTracker\\` |
| `plugin-name` | `task-tracker` |
| `[PluginName]` | `[TaskTracker]` |

The autoloader is **non-namespaced** and registers itself at the bottom via `spl_autoload_register()`.

---

## 20.5 Step 4 — Enums (No Dependencies)

> **Phase 2** — Backed enums for all constants. Create these first because every other file depends on them.

### 4a. PluginConfigType — The identity enum

> **Phase 7, §7.8** — Single source of truth for all metadata.

**File: `includes/Enums/PluginConfigType.php`**

```php
<?php
namespace TaskTracker\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum PluginConfigType: string
{
    case Slug          = 'task-tracker';
    case ShortName     = 'TaskTracker';
    case Name          = 'Task Tracker';
    case Version       = '1.0.0';
    case MinWpVersion  = '5.6';
    case MinPhpVersion = '8.1';
    case ApiNamespace  = 'task-tracker-api';
    case ApiVersion    = 'v1';
    case LogPrefix     = '[TaskTracker]';
    case SettingsGroup = 'task_tracker_settings';
    case DebugConstant = 'TASK_TRACKER_DEBUG';

    public static function apiFullNamespace(): string
    {
        return self::ApiNamespace->value . '/' . self::ApiVersion->value;
    }

    public static function isDebugMode(): bool
    {
        $constantName = self::DebugConstant->value;
        $isDefined = defined($constantName);

        return $isDefined && constant($constantName) === true;
    }

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### 4b. TaskStatusType — Domain enum

> **Phase 2, §2.1** — Every domain concept gets a backed enum.

**File: `includes/Enums/TaskStatusType.php`**

```php
<?php
namespace TaskTracker\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum TaskStatusType: string
{
    case Pending = 'pending';
    case Done    = 'done';

    /** Label for display in admin UI. */
    public function label(): string
    {
        return match ($this) {
            self::Pending => 'Pending',
            self::Done    => 'Done',
        };
    }

    /** CSS class for badge styling. */
    public function cssClass(): string
    {
        return match ($this) {
            self::Pending => 'badge--warning',
            self::Done    => 'badge--success',
        };
    }

    public function isPending(): bool { return $this === self::Pending; }
    public function isDone(): bool { return $this === self::Done; }

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### 4c. EndpointType — Route paths

> **Phase 14, §14.2** — Resource-based, kebab-case route naming.

**File: `includes/Enums/EndpointType.php`**

```php
<?php
namespace TaskTracker\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum EndpointType: string
{
    case Tasks         = 'tasks';
    case TaskComplete  = 'tasks/complete';
    case Status        = 'status';

    /** Route path with leading slash for register_rest_route(). */
    public function route(): string
    {
        return '/' . $this->value;
    }

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### 4d. Remaining enums

Create these by copying patterns from Phase 7 and Phase 2:

| Enum | Cases | Source |
|------|-------|--------|
| `HttpMethodType` | GET, POST, PUT, DELETE | Phase 14, §14.3 |
| `HttpStatusType` | Ok=200, Created=201, BadRequest=400, Unauthorized=401, InternalError=500 | Phase 5, §5.2 |
| `ResponseKeyType` | Status, IsSuccess, IsFailed, Code, Message, Timestamp, Attributes, RequestedAt, TotalRecords, Results, Errors | Phase 7, §7.6 |
| `PhpNativeType` | PhpArray='array', PhpString='string', PhpInteger='integer', etc. | Phase 3, §3.8 |
| `CapabilityType` | ManageOptions='manage_options', ActivatePlugins='activate_plugins' | Phase 3, §3.6 |
| `LogLevelType` | Debug, Info, Warn, Error | Phase 4, §4.3 |
| `OptionNameType` | NotificationsEnabled='task_tracker_notifications' | Phase 15 |
| `HookType` | RestApiInit='rest_api_init', AdminMenu='admin_menu', AdminInit='admin_init' | Phase 8, §8.1 |

**Each enum file follows the same template:**

1. Namespace declaration
2. ABSPATH guard
3. Enum with backed values
4. `match`-based metadata methods where useful
5. `isEqual()`, `isOtherThan()`, `isAnyOf()` comparison trio

---

## 20.6 Step 5 — FileLogger

> **Phase 4, §4.3–§4.7** — Singleton, three log files, rotation, dedup.

**File: `includes/Logging/FileLogger.php`**

Copy the FileLogger reference from Phase 4, replacing `PluginName` → `TaskTracker` and `plugin-name` → `task-tracker`.

Key implementation checkpoints:

| Feature | Verify |
|---------|--------|
| Singleton via `getInstance()` | Phase 4, §4.3 |
| Writes to `info.log`, `error.log`, `stacktrace.log` | Phase 4, §4.3 |
| Log format: `[{timestamp} v{version}] [{Level}] {message} ({file}:{line}) {json}` | Phase 4, §4.4 |
| `debug()` skipped when not in debug mode | Phase 4, §4.3 |
| Stack trace written with `=` separator blocks | Phase 4, §4.5 |
| Rotation at 512KB, max 10 archives | Phase 4, §4.6 |
| In-memory + persistent dedup | Phase 4, §4.7 |

---

## 20.7 Step 6 — Core Infrastructure Traits + Helpers

> **Phase 3, §3.4** — ResponseTrait (safeExecute)  
> **Phase 3, §3.6** — AuthTrait  
> **Phase 3, §3.8** — TypeCheckerTrait  
> **Phase 5, §5.3** — EnvelopeBuilder  
> **Phase 6, §6.3** — validationError()

### 6a. ResponseTrait

**File: `includes/Traits/Core/ResponseTrait.php`**

Contains:
- `safeExecute(callable $callback, string $endpointName): WP_REST_Response` — universal error boundary
- `buildErrorResponse(Throwable $e, string $endpointName): WP_REST_Response` — debug-gated error envelope
- `formatStackFrames(Throwable $e): array` — trace extraction
- `validationError(string $message, WP_REST_Request $request): WP_REST_Response` — 400 errors

Every pattern comes from Phase 4, §4.9 and Phase 6, §6.3.

### 6b. TypeCheckerTrait

**File: `includes/Traits/Core/TypeCheckerTrait.php`**

Copy verbatim from Phase 3, §3.8. Replace namespace `PluginName` → `TaskTracker`.

### 6c. AuthTrait

**File: `includes/Traits/Auth/AuthTrait.php`**

Implements `checkPluginPermission(WP_REST_Request $request)` per Phase 3, §3.6.

### 6d. EnvelopeBuilder

**File: `includes/Helpers/EnvelopeBuilder.php`**

Copy from Phase 7, §7.5. Replace `PluginName` → `TaskTracker`.

### 6e. Other helpers

| Helper | File | Source |
|--------|------|--------|
| `DateHelper` | `Helpers/DateHelper.php` | Phase 4, §4.13 |
| `PathHelper` | `Helpers/PathHelper.php` | Phase 5, §5.1 |
| `ErrorLogHelper` | `Helpers/ErrorLogHelper.php` | Phase 4, §4.11 |

---

## 20.8 Step 7 — Plugin.php (Composition Root)

> **Phase 1, §1.5** + **Phase 7, §7.3** — Singleton, compose all traits, wire hooks, no business logic.

**File: `includes/Core/Plugin.php`**

```php
<?php
namespace TaskTracker\Core;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;
use TaskTracker\Logging\FileLogger;
use TaskTracker\Enums\PluginConfigType;
use TaskTracker\Enums\HookType;

// Trait imports — grouped by domain
use TaskTracker\Traits\Auth\AuthTrait;
use TaskTracker\Traits\Route\RouteRegistrationTrait;
use TaskTracker\Traits\Core\ResponseTrait;
use TaskTracker\Traits\Core\TypeCheckerTrait;

// Feature-domain traits
use TaskTracker\Traits\Task\TaskCreateTrait;
use TaskTracker\Traits\Task\TaskListTrait;
use TaskTracker\Traits\Task\TaskCompleteTrait;

final class Plugin
{
    // ── Auth ──
    use AuthTrait;

    // ── Routes ──
    use RouteRegistrationTrait;

    // ── Core infrastructure ──
    use ResponseTrait;
    use TypeCheckerTrait;

    // ── Feature domains ──
    use TaskCreateTrait;
    use TaskListTrait;
    use TaskCompleteTrait;

    private static ?self $instance = null;
    private FileLogger $fileLogger;

    public static function getInstance(): self
    {
        $hasInstance = (self::$instance !== null);

        if (!$hasInstance) {
            self::$instance = new self();
        }

        return self::$instance;
    }

    private function __construct()
    {
        $startTime = microtime(true);

        $this->fileLogger = FileLogger::getInstance();

        // Register REST routes (fires on rest_api_init)
        add_action('rest_api_init', [$this, 'registerRoutes']);

        // Register admin menu
        add_action('admin_menu', [$this, 'registerAdminPages']);

        // Register shutdown handler
        $this->registerShutdownHandler();

        // Log boot summary
        $elapsedMs = round((microtime(true) - $startTime) * 1000, 2);
        $this->fileLogger->info('Plugin initialized', [
            'version' => PluginConfigType::Version->value,
            'timeMs'  => $elapsedMs,
            'isDebug' => PluginConfigType::isDebugMode(),
        ]);
    }

    public function getPluginSlug(): string
    {
        return PluginConfigType::Slug->value;
    }

    private function registerShutdownHandler(): void
    {
        register_shutdown_function(function (): void {
            $lastError = error_get_last();
            $hasError = ($lastError !== null);

            if (!$hasError) {
                return;
            }

            $fatalTypes = [E_ERROR, E_CORE_ERROR, E_COMPILE_ERROR, E_PARSE];
            $isFatal = in_array($lastError['type'], $fatalTypes, true);

            if (!$isFatal) {
                return;
            }

            $memoryMb = round(memory_get_peak_usage(true) / 1048576, 2);

            error_log(sprintf(
                '[%s] FATAL: %s in %s:%d | Memory: %sMB',
                PluginConfigType::ShortName->value,
                $lastError['message'],
                $lastError['file'],
                $lastError['line'],
                $memoryMb,
            ));
        });
    }

    private function __clone() {}

    public function __wakeup(): void
    {
        throw new \RuntimeException('Cannot unserialize singleton');
    }
}
```

### Observations

| Pattern | Phase |
|---------|-------|
| No business logic in constructor | Phase 1, §1.5 |
| Traits composed by domain groups | Phase 3, §3.7 |
| Shutdown handler for fatals | Phase 7, §7.3 |
| Boot timing logged | Phase 4, §4.4 |

---

## 20.9 Step 8 — Route Registration

> **Phase 3, §3.5** + **Phase 14** — Grouped, fault-tolerant route registration.

**File: `includes/Traits/Route/RouteRegistrationTrait.php`**

```php
<?php
namespace TaskTracker\Traits\Route;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;
use TaskTracker\Enums\PluginConfigType;
use TaskTracker\Enums\EndpointType;
use TaskTracker\Enums\HttpMethodType;

trait RouteRegistrationTrait
{
    public function registerRoutes(): void
    {
        $namespace = PluginConfigType::apiFullNamespace();
        $registered = 0;
        $failed = 0;

        $safeRegister = function (
            string $route,
            string $method,
            callable $callback,
            callable $permission,
        ) use ($namespace, &$registered, &$failed): void {
            try {
                register_rest_route($namespace, $route, [
                    'methods'             => $method,
                    'callback'            => $callback,
                    'permission_callback' => $permission,
                ]);
                $registered++;
            } catch (Throwable $e) {
                $failed++;
                error_log(
                    "[TaskTracker] Route registration failed for {$route}: "
                    . $e->getMessage() . "\n" . $e->getTraceAsString()
                );
            }
        };

        // ── Task endpoints ──
        $safeRegister(
            EndpointType::Tasks->route(),
            HttpMethodType::POST->value,
            [$this, 'handleCreateTask'],
            [$this, 'checkPluginPermission'],
        );

        $safeRegister(
            EndpointType::Tasks->route(),
            HttpMethodType::GET->value,
            [$this, 'handleListTasks'],
            [$this, 'checkPluginPermission'],
        );

        $safeRegister(
            EndpointType::TaskComplete->route(),
            HttpMethodType::POST->value,
            [$this, 'handleCompleteTask'],
            [$this, 'checkPluginPermission'],
        );

        // ── Log summary ──
        $this->fileLogger->info('Routes registered', [
            'registered' => $registered,
            'failed'     => $failed,
        ]);
    }
}
```

---

## 20.10 Step 9 — Database Migration (SQLite)

> **Phase 8, §8.5** — SQLite migrations via traits.

**File: `includes/Database/DatabaseMigrationsTrait.php`**

```php
<?php
namespace TaskTracker\Database;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;

trait DatabaseMigrationsTrait
{
    /**
     * Run all pending migrations. Called from Activator::activate().
     *
     * @param \PDO $pdo The SQLite connection
     */
    private function runMigrations(\PDO $pdo): void
    {
        // ── Migration 001: Tasks table ──
        $pdo->exec("
            CREATE TABLE IF NOT EXISTS tasks (
                Id         INTEGER PRIMARY KEY AUTOINCREMENT,
                Title      TEXT    NOT NULL,
                Priority   INTEGER NOT NULL DEFAULT 0,
                Status     TEXT    NOT NULL DEFAULT 'pending',
                CreatedAt  TEXT    NOT NULL DEFAULT (datetime('now')),
                UpdatedAt  TEXT    NOT NULL DEFAULT (datetime('now'))
            )
        ");

        $this->fileLogger->info('Database migrations complete');
    }
}
```

### Rules applied

| Rule | Detail | Phase |
|------|--------|-------|
| `CREATE TABLE IF NOT EXISTS` | Idempotent migrations | Phase 8, §8.5 |
| Status column uses enum value | `'pending'` = `TaskStatusType::Pending->value` | Phase 2 |
| Datetime defaults use SQLite `datetime('now')` | UTC timestamps | Phase 4, §4.13 |

---

## 20.11 Step 10 — Feature Handler Traits

> **Phase 3, §3.3** — Trait anatomy: public handler wraps `safeExecute`, private method has logic.  
> **Phase 6** — Validation guard clauses at top of every handler.

### 10a. TaskCreateTrait — POST /tasks

**File: `includes/Traits/Task/TaskCreateTrait.php`**

```php
<?php
namespace TaskTracker\Traits\Task;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;
use WP_REST_Request;
use WP_REST_Response;
use TaskTracker\Enums\TaskStatusType;
use TaskTracker\Helpers\EnvelopeBuilder;

trait TaskCreateTrait
{
    /**
     * Handle POST /tasks — Create a new task.
     */
    public function handleCreateTask(WP_REST_Request $request): WP_REST_Response
    {
        return $this->safeExecute(
            fn() => $this->executeCreateTask($request),
            'create-task',
        );
    }

    private function executeCreateTask(WP_REST_Request $request): WP_REST_Response
    {
        // ── Validate body ──
        $body = $request->get_json_params();
        $hasBody = ($body !== null && $this->isArray($body));

        if (!$hasBody) {
            return $this->validationError('Request body must be a JSON object', $request);
        }

        // ── Validate required: title ──
        $title = $body['title'] ?? null;
        $hasTitle = ($title !== null && $this->isString($title));

        if (!$hasTitle) {
            return $this->validationError('Missing required field: title', $request);
        }

        $titleLength = mb_strlen($title);
        $isTitleTooLong = ($titleLength > 200);

        if ($isTitleTooLong) {
            return $this->validationError(
                'Field "title" must not exceed 200 characters',
                $request,
            );
        }

        // ── Validate optional: priority ──
        $priority = $body['priority'] ?? null;
        $hasPriority = ($priority !== null);

        if ($hasPriority) {
            $isPriorityValid = $this->isInteger($priority);

            if (!$isPriorityValid) {
                return $this->validationError(
                    'Field "priority" must be an integer',
                    $request,
                );
            }
        }

        // ── Sanitise ──
        $sanitisedTitle = sanitize_text_field($title);
        $resolvedPriority = $hasPriority ? absint($priority) : 0;

        // ── Insert ──
        $pdo = $this->getDatabase();
        $stmt = $pdo->prepare(
            'INSERT INTO tasks (Title, Priority, Status) VALUES (:title, :priority, :status)'
        );
        $stmt->execute([
            ':title'    => $sanitisedTitle,
            ':priority' => $resolvedPriority,
            ':status'   => TaskStatusType::Pending->value,
        ]);

        $taskId = (int) $pdo->lastInsertId();

        $this->fileLogger->info('Task created', [
            'taskId'   => $taskId,
            'title'    => $sanitisedTitle,
            'priority' => $resolvedPriority,
        ]);

        return EnvelopeBuilder::success('Task created', 201)
            ->setRequestedAt($request->get_route())
            ->setSingleResult([
                'Id'       => $taskId,
                'Title'    => $sanitisedTitle,
                'Priority' => $resolvedPriority,
                'Status'   => TaskStatusType::Pending->value,
            ])
            ->toResponse();
    }
}
```

### Validation checklist applied (Phase 6, §6.8)

- [x] Body validated as array
- [x] Required `title` has `$hasTitle` guard
- [x] String length enforced (200 chars)
- [x] Optional `priority` type-checked with `TypeCheckerTrait`
- [x] Sanitisation after validation (`sanitize_text_field`)
- [x] Uses `$this->validationError()`, not `EnvelopeBuilder::error()` directly
- [x] Enum value for status, no magic string

### 10b. TaskListTrait — GET /tasks

**File: `includes/Traits/Task/TaskListTrait.php`**

```php
<?php
namespace TaskTracker\Traits\Task;

if (!defined('ABSPATH')) {
    exit;
}

use WP_REST_Request;
use WP_REST_Response;
use TaskTracker\Helpers\EnvelopeBuilder;

trait TaskListTrait
{
    public function handleListTasks(WP_REST_Request $request): WP_REST_Response
    {
        return $this->safeExecute(
            fn() => $this->executeListTasks($request),
            'list-tasks',
        );
    }

    private function executeListTasks(WP_REST_Request $request): WP_REST_Response
    {
        // ── Pagination (silent defaults allowed per Phase 6, §6.5) ──
        $page = $request->get_param('page');
        $perPage = $request->get_param('per_page');

        $resolvedPage = ($page !== null && $this->isNumeric($page))
            ? max(1, absint($page))
            : 1;
        $resolvedPerPage = ($perPage !== null && $this->isNumeric($perPage))
            ? min(100, max(1, absint($perPage)))
            : 20;

        $offset = ($resolvedPage - 1) * $resolvedPerPage;

        // ── Query ──
        $pdo = $this->getDatabase();

        $countStmt = $pdo->query('SELECT COUNT(*) FROM tasks');
        $totalRecords = (int) $countStmt->fetchColumn();

        $stmt = $pdo->prepare(
            'SELECT Id, Title, Priority, Status, CreatedAt, UpdatedAt
             FROM tasks
             ORDER BY CreatedAt DESC
             LIMIT :limit OFFSET :offset'
        );
        $stmt->bindValue(':limit', $resolvedPerPage, \PDO::PARAM_INT);
        $stmt->bindValue(':offset', $offset, \PDO::PARAM_INT);
        $stmt->execute();

        $tasks = $stmt->fetchAll(\PDO::FETCH_ASSOC);

        return EnvelopeBuilder::success()
            ->setRequestedAt($request->get_route())
            ->setListResult($tasks)
            ->toResponse();
    }
}
```

### 10c. TaskCompleteTrait — POST /tasks/complete

**File: `includes/Traits/Task/TaskCompleteTrait.php`**

```php
<?php
namespace TaskTracker\Traits\Task;

if (!defined('ABSPATH')) {
    exit;
}

use WP_REST_Request;
use WP_REST_Response;
use TaskTracker\Enums\TaskStatusType;
use TaskTracker\Helpers\EnvelopeBuilder;

trait TaskCompleteTrait
{
    public function handleCompleteTask(WP_REST_Request $request): WP_REST_Response
    {
        return $this->safeExecute(
            fn() => $this->executeCompleteTask($request),
            'complete-task',
        );
    }

    private function executeCompleteTask(WP_REST_Request $request): WP_REST_Response
    {
        // ── Validate ──
        $body = $request->get_json_params();
        $hasBody = ($body !== null && $this->isArray($body));

        if (!$hasBody) {
            return $this->validationError('Request body must be a JSON object', $request);
        }

        $taskId = $body['task_id'] ?? null;
        $hasTaskId = ($taskId !== null && $this->isInteger($taskId));

        if (!$hasTaskId) {
            return $this->validationError('Missing required field: task_id (integer)', $request);
        }

        // ── Update ──
        $pdo = $this->getDatabase();
        $stmt = $pdo->prepare(
            "UPDATE tasks SET Status = :status, UpdatedAt = datetime('now') WHERE Id = :id"
        );
        $stmt->execute([
            ':status' => TaskStatusType::Done->value,
            ':id'     => absint($taskId),
        ]);

        $affectedRows = $stmt->rowCount();
        $hasMatch = ($affectedRows > 0);

        if (!$hasMatch) {
            return EnvelopeBuilder::error('Task not found', 404)
                ->setRequestedAt($request->get_route())
                ->toResponse();
        }

        $this->fileLogger->info('Task completed', ['taskId' => $taskId]);

        return EnvelopeBuilder::success('Task marked as done')
            ->setRequestedAt($request->get_route())
            ->setSingleResult(['Id' => $taskId, 'Status' => TaskStatusType::Done->value])
            ->toResponse();
    }
}
```

---

## 20.12 Step 11 — Admin Settings Page

> **Phase 8, §8.1** — Admin menu registration  
> **Phase 11, §11.1–§11.4** — Template ≤200 lines, orchestrator pattern  
> **Phase 15** — Settings architecture  
> **Phase 12** — Design system tokens

### 11a. Admin menu registration (in Plugin.php)

Add this method to Plugin.php (or a separate `AdminPageTrait`):

```php
public function registerAdminPages(): void
{
    add_menu_page(
        PluginConfigType::Name->value,
        PluginConfigType::ShortName->value,
        CapabilityType::ManageOptions->value,
        PluginConfigType::Slug->value,
        [$this, 'renderSettingsPage'],
        'dashicons-list-view',
        80,
    );
}

public function renderSettingsPage(): void
{
    $isAuthorized = current_user_can(CapabilityType::ManageOptions->value);

    if (!$isAuthorized) {
        wp_die('Unauthorized access');
    }

    $templatePath = plugin_dir_path(dirname(__DIR__, 2)) . 'templates/settings.php';
    $hasTemplate = file_exists($templatePath);

    if ($hasTemplate) {
        include $templatePath;
    }
}
```

### 11b. Settings template

**File: `templates/settings.php`**

```php
<?php
/**
 * Settings page template — Task Tracker.
 *
 * @var none — all data fetched inline (simple page)
 * @package TaskTracker
 */

if (!defined('ABSPATH')) {
    exit;
}

$isNotificationsEnabled = get_option('task_tracker_notifications', 'no');
?>
<div class="wrap">
    <h1><?php echo esc_html(get_admin_page_title()); ?></h1>

    <form method="post" action="options.php">
        <?php settings_fields('task_tracker_settings'); ?>

        <table class="form-table">
            <tr>
                <th scope="row">
                    <label for="task_tracker_notifications">Enable Notifications</label>
                </th>
                <td>
                    <select name="task_tracker_notifications" id="task_tracker_notifications">
                        <option value="yes" <?php selected($isNotificationsEnabled, 'yes'); ?>>
                            Yes
                        </option>
                        <option value="no" <?php selected($isNotificationsEnabled, 'no'); ?>>
                            No
                        </option>
                    </select>
                    <p class="description">
                        Send email notifications when a task is completed.
                    </p>
                </td>
            </tr>
        </table>

        <?php submit_button('Save Settings'); ?>
    </form>
</div>
```

**Template line count:** ~42 lines — well within the 200-line limit (Phase 11, §11.1).

---

## 20.13 Step 12 — Admin Settings Registration

> **Phase 15** — `OptionNameType` enum for wp_options keys.

Register settings in `admin_init`:

```php
public function registerSettings(): void
{
    register_setting(
        PluginConfigType::SettingsGroup->value,
        OptionNameType::NotificationsEnabled->value,
        [
            'type'              => 'string',
            'sanitize_callback' => 'sanitize_text_field',
            'default'           => 'no',
        ],
    );
}
```

Wire in Plugin constructor: `add_action('admin_init', [$this, 'registerSettings']);`

---

## 20.14 Step 13 — Tests

> **Phase 9, §9.1–§9.3** — Unit tests run without WordPress.

### 13a. Test bootstrap

**File: `tests/bootstrap.php`**

```php
<?php
// Mock ABSPATH so includes load
define('ABSPATH', '/tmp/wordpress/');

// Load autoloader
require_once __DIR__ . '/../includes/Autoloader.php';
```

### 13b. Enum test

**File: `tests/Unit/Enums/TaskStatusTypeTest.php`**

```php
<?php
namespace TaskTracker\Tests\Unit\Enums;

use PHPUnit\Framework\TestCase;
use TaskTracker\Enums\TaskStatusType;

final class TaskStatusTypeTest extends TestCase
{
    public function testPendingValueIsCorrect(): void
    {
        $this->assertSame('pending', TaskStatusType::Pending->value);
    }

    public function testDoneValueIsCorrect(): void
    {
        $this->assertSame('done', TaskStatusType::Done->value);
    }

    public function testLabelReturnsHumanReadable(): void
    {
        $this->assertSame('Pending', TaskStatusType::Pending->label());
        $this->assertSame('Done', TaskStatusType::Done->label());
    }

    public function testCssClassReturnsBadgeClass(): void
    {
        $this->assertSame('badge--warning', TaskStatusType::Pending->cssClass());
        $this->assertSame('badge--success', TaskStatusType::Done->cssClass());
    }

    public function testIsPendingReturnsCorrectly(): void
    {
        $this->assertTrue(TaskStatusType::Pending->isPending());
        $this->assertFalse(TaskStatusType::Done->isPending());
    }

    public function testTryFromReturnsNullForInvalidValue(): void
    {
        $result = TaskStatusType::tryFrom('invalid');

        $this->assertNull($result);
    }

    public function testIsEqualComparison(): void
    {
        $this->assertTrue(TaskStatusType::Pending->isEqual(TaskStatusType::Pending));
        $this->assertFalse(TaskStatusType::Pending->isEqual(TaskStatusType::Done));
    }

    public function testIsAnyOfComparison(): void
    {
        $result = TaskStatusType::Pending->isAnyOf(
            TaskStatusType::Pending,
            TaskStatusType::Done,
        );

        $this->assertTrue($result);
    }
}
```

---

## 20.15 Step 14 — Uninstall Cleanup

> **Phase 7, §7.4** — Destructive cleanup when plugin is deleted.

**File: `uninstall.php`**

```php
<?php
if (!defined('WP_UNINSTALL_PLUGIN')) {
    exit;
}

// Delete options
delete_option('task_tracker_notifications');
delete_option('task_tracker_settings');

// Delete upload directory (logs, database)
$uploadDir = wp_upload_dir();
$pluginDir = $uploadDir['basedir'] . '/task-tracker';
$hasDir = is_dir($pluginDir);

if ($hasDir) {
    $iterator = new RecursiveIteratorIterator(
        new RecursiveDirectoryIterator($pluginDir, RecursiveDirectoryIterator::SKIP_DOTS),
        RecursiveIteratorIterator::CHILD_FIRST,
    );

    foreach ($iterator as $item) {
        $isDir = $item->isDir();

        if ($isDir) {
            rmdir($item->getRealPath());
        } else {
            unlink($item->getRealPath());
        }
    }

    rmdir($pluginDir);
}
```

---

## 20.16 Phase Coverage Matrix

Every phase of the Gold Standard applied in this walkthrough:

| Phase | Where Applied | Step |
|-------|---------------|------|
| **0 — Quick Start** | Followed the 5-file creation order | Steps 2–7 |
| **1 — Foundation** | Folder structure, ABSPATH guards, namespace, bootstrap | Steps 1–3 |
| **2 — Enums** | All constants as backed enums, `match` metadata, comparison trio | Step 4 |
| **3 — Traits** | ResponseTrait, AuthTrait, TypeCheckerTrait, feature traits | Steps 6–10 |
| **4 — Logging** | FileLogger, two-tier architecture, debug-mode gating, shutdown handler | Steps 5, 7 |
| **5 — Helpers** | EnvelopeBuilder, PathHelper, DateHelper, ErrorLogHelper | Step 6 |
| **6 — Validation** | Guard clauses, `validationError()`, type checking, sanitisation | Step 10 |
| **7 — Reference Impl** | All files follow Phase 7 templates with search-replace | Steps 2–7 |
| **8 — WP Integration** | Admin menu, settings page, SQLite migration | Steps 9, 11, 12 |
| **9 — Testing** | PHPUnit test structure, enum tests | Step 13 |
| **10 — Deployment** | Uninstall cleanup | Step 14 |
| **11 — Frontend** | Template ≤200 lines, output escaping | Step 11 |
| **12 — Design System** | CSS variable tokens in admin styles, slug substitution | Step 11 |
| **13 — Admin UI** | Settings page layout, form controls | Step 11 |
| **14 — REST API** | Namespace, route naming, grouped registration | Steps 4, 8 |
| **15 — Settings** | `OptionNameType`, settings registration, defaults | Step 12 |
| **16 — Error Handling** | `ErrorResponse`, forbidden patterns avoided | Steps 2, 10 |
| **17 — Data Files** | Seed manifest (referenced) | Step 1 folder structure |
| **18 — Frontend JS** | (Not needed — no JS interactivity in this minimal example) | — |
| **19 — Micro-ORM** | (Not used — direct PDO for minimal plugin; use ORM for complex queries) | — |

---

## 20.17 Final Checklist — "Is My Plugin Gold Standard?"

Run through this checklist before shipping:

- [ ] Every PHP file has ABSPATH guard after namespace
- [ ] Zero magic strings — all constants are enum cases
- [ ] Every REST handler wrapped in `safeExecute()`
- [ ] Every catch block catches `Throwable` with stack trace logged
- [ ] No forbidden error patterns (Phase 4, §4.8)
- [ ] All validation uses `$this->validationError()` with guard clauses
- [ ] All strings sanitised after validation, not instead of
- [ ] FileLogger with rotation and dedup configured
- [ ] Version in `PluginConfigType::Version` only
- [ ] Template files ≤200 lines
- [ ] Admin pages use capability checks
- [ ] `uninstall.php` removes ALL plugin data
- [ ] Tests cover every enum and every public helper method
- [ ] No `array()` syntax — only `[]`
- [ ] No `is_array()` / `is_string()` — use `TypeCheckerTrait`

---

*This walkthrough demonstrates a complete plugin. For production plugins with more features, repeat Steps 4 (enums), 9 (migration), and 10 (handler traits) for each new domain.*
