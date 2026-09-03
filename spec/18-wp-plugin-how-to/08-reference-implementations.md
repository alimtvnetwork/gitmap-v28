# Phase 7 — Complete Reference Implementations

> **Purpose:** Provide full, copy-paste-ready implementations of every foundational file so no AI has to guess. Every file shown here is a **complete, working implementation** — not a fragment.
> **Rule:** When generating a new plugin, copy these files verbatim and search-replace `PluginName` / `plugin-name` / `PLUGIN_NAME` with the actual plugin identity.

---

## 7.1 Bootstrap File — `plugin-name.php`

This is the only file WordPress loads directly. It must be non-namespaced.

```php
<?php
/**
 * Plugin Name:       Plugin Name
 * Plugin URI:        https://example.com/plugin-name
 * Description:       One-line description of what this plugin does.
 * Version:           1.0.0
 * Requires at least: 5.6
 * Requires PHP:      8.1
 * Author:            Your Name
 * Author URI:        https://example.com
 * License:           GPL-2.0-or-later
 * License URI:       https://www.gnu.org/licenses/gpl-2.0.html
 * Text Domain:       plugin-name
 */

if (!defined('ABSPATH')) {
    exit;
}

/** Enable debug mode — exposes stack traces in API responses. Set to true only in development. */
define('PLUGIN_NAME_DEBUG', false);

// ── Load autoloader ──
$autoloaderPath = __DIR__ . '/includes/Autoloader.php';
$hasAutoloader = file_exists($autoloaderPath);

if (!$hasAutoloader) {
    error_log('[PluginName] FATAL: Autoloader.php not found at ' . $autoloaderPath);

    return;
}

require_once $autoloaderPath;

// ── Boot plugin on rest_api_init (for REST-only plugins) ──
// For plugins that also need admin UI, use 'plugins_loaded' instead
add_action('rest_api_init', function (): void {
    try {
        \PluginName\Core\Plugin::getInstance();
    } catch (\Throwable $e) {
        error_log('[PluginName] Boot failed: ' . $e->getMessage() . "\n" . $e->getTraceAsString());
    }
});

// ── Register activation hook ──
register_activation_hook(__FILE__, function (): void {
    try {
        \PluginName\Core\Activator::activate();
    } catch (\Throwable $e) {
        error_log('[PluginName] Activation failed: ' . $e->getMessage() . "\n" . $e->getTraceAsString());
    }
});

// ── Register deactivation hook ──
register_deactivation_hook(__FILE__, function (): void {
    try {
        \PluginName\Core\Deactivator::deactivate();
    } catch (\Throwable $e) {
        error_log('[PluginName] Deactivation failed: ' . $e->getMessage() . "\n" . $e->getTraceAsString());
    }
});
```

### Key observations

| Pattern | Reason |
|---------|--------|
| `define('PLUGIN_NAME_DEBUG', false)` | Gates stack traces in API responses (see Phase 4, §4.2) |
| Anonymous closures for hooks | Prevents global function name collisions |
| Every hook body wrapped in try-catch | Bootstrap errors must never crash WordPress |
| Autoloader existence check before `require_once` | Graceful failure instead of PHP fatal |
| `return` (not `exit`) on autoloader failure | Allows WordPress to continue loading other plugins |

---

## 7.2 Autoloader — `includes/Autoloader.php`

```php
<?php
/**
 * Autoloader — PSR-4 class loader for the plugin namespace.
 *
 * This file is non-namespaced because it must load before
 * the namespace system is available.
 *
 * @package PluginName
 * @since   1.0.0
 */

if (!defined('ABSPATH')) {
    exit;
}

final class PluginNameAutoloader
{
    /** @var string The root namespace this autoloader handles. */
    private const NAMESPACE_PREFIX = 'PluginName\\';

    /** @var string The base directory for class files. */
    private string $baseDir;

    /** @var string Path to the autoloader diagnostic log file. */
    private string $logPath;

    public function __construct()
    {
        $this->baseDir = __DIR__ . '/';

        $uploadDir = wp_upload_dir();
        $hasUploadDir = (!empty($uploadDir['basedir']));
        $uploadsBase = $hasUploadDir ? $uploadDir['basedir'] : WP_CONTENT_DIR . '/uploads';

        $logDir = $uploadsBase . '/plugin-name/logs';
        $hasDirExists = is_dir($logDir);

        if (!$hasDirExists) {
            wp_mkdir_p($logDir);
        }

        $this->logPath = $logDir . '/autoloader.log';
    }

    /**
     * Attempt to load a class file for the given fully-qualified class name.
     *
     * @param string $className The fully-qualified class name (e.g., 'PluginName\Enums\HttpStatusType')
     */
    public function loadClass(string $className): void
    {
        $prefixLength = strlen(self::NAMESPACE_PREFIX);
        $isOurNamespace = (strncmp($className, self::NAMESPACE_PREFIX, $prefixLength) === 0);

        if (!$isOurNamespace) {
            return;
        }

        $relativeClass = substr($className, $prefixLength);
        $filePath = $this->baseDir . str_replace('\\', '/', $relativeClass) . '.php';

        $hasFile = file_exists($filePath);

        if (!$hasFile) {
            $this->log("Class file not found: {$className} → {$filePath}");

            return;
        }

        try {
            require_once $filePath;
        } catch (\Throwable $e) {
            $this->log(
                "Failed to load {$className}: {$e->getMessage()}\n{$e->getTraceAsString()}"
            );

            throw $e;
        }
    }

    /**
     * Write a diagnostic entry to the autoloader log.
     *
     * @param string $message The log message
     */
    private function log(string $message): void
    {
        $timestamp = gmdate('d-M-y g:i A');
        $entry = "[{$timestamp}] [Autoloader] {$message}\n";

        error_log(trim($entry));
        @file_put_contents($this->logPath, $entry, FILE_APPEND | LOCK_EX);
    }

    /**
     * Register this autoloader with PHP's SPL autoload stack.
     */
    public function register(): void
    {
        spl_autoload_register([$this, 'loadClass']);
    }
}

// Self-register on include
(new PluginNameAutoloader())->register();
```

### Why this design

| Decision | Reason |
|----------|--------|
| `final` class | No inheritance expected or wanted |
| Non-namespaced | Loaded before namespaces are available |
| `strncmp` for prefix check | Faster than `str_starts_with()` and available in PHP 8.0 fallback scenarios |
| Dual logging (error_log + file) | `error_log` goes to WP_DEBUG log; file is plugin-specific and always readable |
| `@file_put_contents` with `@` | Suppresses warnings if log directory was deleted mid-request |
| Re-throw on require failure | Caller (bootstrap) needs to know the load failed |

---

## 7.3 Plugin.php — `includes/Core/Plugin.php`

```php
<?php
/**
 * Plugin — Composition root and singleton entry point.
 *
 * @package PluginName\Core
 * @since   1.0.0
 */

namespace PluginName\Core;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;
use PluginName\Logging\FileLogger;
use PluginName\Enums\PluginConfigType;
use PluginName\Enums\HookType;

// Trait imports — grouped by domain
use PluginName\Traits\Auth\AuthTrait;
use PluginName\Traits\Route\RouteRegistrationTrait;
use PluginName\Traits\Core\ResponseTrait;
use PluginName\Traits\Core\TypeCheckerTrait;
use PluginName\Traits\Core\StatusHandlerTrait;
use PluginName\Traits\Core\PluginInventoryTrait;
// Add feature-domain trait imports here

final class Plugin
{
    // ── Auth ──
    use AuthTrait;

    // ── Routes ──
    use RouteRegistrationTrait;

    // ── Core infrastructure ──
    use ResponseTrait;
    use TypeCheckerTrait;
    use StatusHandlerTrait;
    use PluginInventoryTrait;

    // ── Feature domains (add new traits here) ──
    // use UploadHandlerTrait;
    // use ActivateHandlerTrait;

    /** @var self|null Singleton instance. */
    private static ?self $instance = null;

    /** @var FileLogger Structured file logger. */
    private FileLogger $fileLogger;

    /**
     * Get the singleton instance. Creates it on first call.
     */
    public static function getInstance(): self
    {
        $hasInstance = (self::$instance !== null);

        if (!$hasInstance) {
            self::$instance = new self();
        }

        return self::$instance;
    }

    /**
     * Private constructor — called once via getInstance().
     * Wires hooks and initialises logging. No business logic here.
     */
    private function __construct()
    {
        $startTime = microtime(true);

        // Obtain logger
        $this->fileLogger = FileLogger::getInstance();

        // Register REST routes
        $this->registerRoutes();

        // Register shutdown handler for fatal errors
        $this->registerShutdownHandler();

        // Log boot summary
        $elapsedMs = round((microtime(true) - $startTime) * 1000, 2);
        $this->fileLogger->info('Plugin initialized', [
            'version'  => PluginConfigType::Version->value,
            'timeMs'   => $elapsedMs,
            'isDebug'  => PluginConfigType::isDebugMode(),
        ]);
    }

    /**
     * Return the plugin slug for logging and identification.
     */
    public function getPluginSlug(): string
    {
        return PluginConfigType::Slug->value;
    }

    /**
     * Register the shutdown handler for uncaught fatal errors.
     */
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

    /** Prevent cloning. */
    private function __clone() {}

    /** Prevent unserialization. */
    public function __wakeup(): void
    {
        throw new \RuntimeException('Cannot unserialize singleton');
    }
}
```

---

## 7.4 Activator & Deactivator — `includes/Core/Activator.php`

```php
<?php
/**
 * Activator — Runs on plugin activation.
 *
 * @package PluginName\Core
 * @since   1.0.0
 */

namespace PluginName\Core;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Helpers\PathHelper;

final class Activator
{
    /**
     * Run activation tasks. Called from register_activation_hook().
     */
    public static function activate(): void
    {
        // ── 1. Create required directories ──
        PathHelper::ensureDirectory(PathHelper::getLogsDir());
        PathHelper::ensureDirectory(PathHelper::getTempDir());

        // ── 2. Create or migrate database (if using SQLite) ──
        // DatabaseMigrator::runPending();

        // ── 3. Set default options ──
        $hasExistingSettings = (get_option('plugin_name_settings') !== false);

        if (!$hasExistingSettings) {
            add_option('plugin_name_settings', [
                'version' => '1.0.0',
                'activated_at' => gmdate('c'),
            ]);
        }

        // ── 4. Schedule cron events (if needed) ──
        // See Phase 8, §8.3 for WP-Cron patterns

        // ── 5. Flush rewrite rules (if adding custom REST routes via rewrite) ──
        flush_rewrite_rules();

        error_log('[PluginName] Plugin activated successfully');
    }
}
```

### Deactivator — `includes/Core/Deactivator.php`

```php
<?php
/**
 * Deactivator — Runs on plugin deactivation.
 *
 * @package PluginName\Core
 * @since   1.0.0
 */

namespace PluginName\Core;

if (!defined('ABSPATH')) {
    exit;
}

final class Deactivator
{
    /**
     * Run deactivation tasks. Called from register_deactivation_hook().
     */
    public static function deactivate(): void
    {
        // ── 1. Unschedule cron events ──
        $nextRun = wp_next_scheduled('plugin_name_cron_hook');
        $hasScheduledEvent = ($nextRun !== false);

        if ($hasScheduledEvent) {
            wp_unschedule_event($nextRun, 'plugin_name_cron_hook');
        }

        // ── 2. Clean up transients ──
        delete_transient('plugin_name_cache');

        // ── 3. Flush rewrite rules ──
        flush_rewrite_rules();

        error_log('[PluginName] Plugin deactivated');
    }
}
```

### Uninstall — `uninstall.php` (plugin root)

```php
<?php
/**
 * Uninstall — Runs when the plugin is deleted via WordPress admin.
 *
 * This file is called by WordPress directly. It must not load the plugin.
 * It performs destructive cleanup: delete options, drop tables, remove files.
 */

if (!defined('WP_UNINSTALL_PLUGIN')) {
    exit;
}

// ── 1. Delete options ──
delete_option('plugin_name_settings');

// ── 2. Delete upload directory and all contents ──
$uploadDir = wp_upload_dir();
$pluginDir = $uploadDir['basedir'] . '/plugin-name';
$hasDir = is_dir($pluginDir);

if ($hasDir) {
    // Recursive delete
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

// ── 3. Delete scheduled cron events ──
wp_clear_scheduled_hook('plugin_name_cron_hook');
```

### Activation vs Deactivation vs Uninstall — When to use what

| Hook | Fires when | Do | Don't |
|------|-----------|-----|-------|
| `register_activation_hook` | Plugin activated | Create dirs, set defaults, schedule cron | Delete data, show admin notices |
| `register_deactivation_hook` | Plugin deactivated | Unschedule cron, flush rewrite rules | Delete data (user may reactivate) |
| `uninstall.php` | Plugin deleted from admin | Delete ALL plugin data: options, files, tables | Reference plugin classes (not loaded) |

---

## 7.5 EnvelopeBuilder — `includes/Helpers/EnvelopeBuilder.php`

```php
<?php
/**
 * EnvelopeBuilder — Fluent builder for the standard API response envelope.
 *
 * @package PluginName\Helpers
 * @since   1.0.0
 */

namespace PluginName\Helpers;

if (!defined('ABSPATH')) {
    exit;
}

use WP_REST_Response;
use PluginName\Enums\PluginConfigType;
use PluginName\Enums\ResponseKeyType;

final class EnvelopeBuilder
{
    private bool $isSuccess;
    private int $code;
    private string $message;
    private string $requestedAt = '';

    /** @var array<int, array<string, mixed>> */
    private array $results = [];

    private ?string $backendMessage = null;
    private ?string $exceptionType = null;

    /** @var array<int, string>|null */
    private ?array $stackFrames = null;

    private function __construct(bool $isSuccess, string $message, int $code)
    {
        $this->isSuccess = $isSuccess;
        $this->message = $message;
        $this->code = $code;
    }

    /**
     * Create a success envelope.
     *
     * @param string $message Status message (default: 'OK')
     * @param int    $code    HTTP status code (default: 200)
     */
    public static function success(string $message = 'OK', int $code = 200): self
    {
        return new self(true, $message, $code);
    }

    /**
     * Create an error envelope. Automatically extracts exception info if provided.
     *
     * @param string          $message   Error description
     * @param int             $code      HTTP status code (default: 500)
     * @param \Throwable|null $exception Optional exception for trace extraction
     */
    public static function error(
        string $message,
        int $code = 500,
        ?\Throwable $exception = null,
    ): self {
        $isDebug = PluginConfigType::isDebugMode();
        $isServerError = ($code >= 500);

        // Gate the error message for 5xx in production
        $resolvedMessage = ($isServerError && !$isDebug)
            ? 'An internal error occurred'
            : $message;

        $builder = new self(false, $resolvedMessage, $code);
        $builder->backendMessage = $message;

        $hasException = ($exception !== null);

        if ($hasException && $isDebug) {
            $builder->exceptionType = get_class($exception);
            $builder->stackFrames = self::extractFrames($exception);
        }

        return $builder;
    }

    public function setRequestedAt(string $path): self
    {
        $this->requestedAt = $path;

        return $this;
    }

    /**
     * @param array<string, mixed> $item Single result item
     */
    public function setSingleResult(array $item): self
    {
        $this->results = [$item];

        return $this;
    }

    /**
     * @param array<int, array<string, mixed>> $items List of result items
     */
    public function setListResult(array $items): self
    {
        $this->results = $items;

        return $this;
    }

    /**
     * @param array<int, string> $frames Stack trace frames
     */
    public function setStackTrace(array $frames): self
    {
        $isDebug = PluginConfigType::isDebugMode();

        if ($isDebug) {
            $this->stackFrames = $frames;
        }

        return $this;
    }

    /**
     * Build and return the final WP_REST_Response.
     */
    public function toResponse(): WP_REST_Response
    {
        $envelope = [
            ResponseKeyType::Status->value => [
                ResponseKeyType::IsSuccess->value => $this->isSuccess,
                ResponseKeyType::IsFailed->value  => !$this->isSuccess,
                ResponseKeyType::Code->value      => $this->code,
                ResponseKeyType::Message->value   => $this->message,
                ResponseKeyType::Timestamp->value => DateHelper::nowUtc(),
            ],
            ResponseKeyType::Attributes->value => [
                ResponseKeyType::RequestedAt->value  => $this->requestedAt,
                ResponseKeyType::TotalRecords->value => count($this->results),
            ],
            ResponseKeyType::Results->value => $this->results,
        ];

        // Only add Errors key on failure AND when there's debug info to show
        $hasErrors = (!$this->isSuccess && $this->stackFrames !== null);

        if ($hasErrors) {
            $envelope[ResponseKeyType::Errors->value] = [
                'BackendMessage' => $this->backendMessage,
                'ExceptionType'  => $this->exceptionType,
                'Backend'        => $this->stackFrames,
            ];
        }

        return new WP_REST_Response($envelope, $this->code);
    }

    /**
     * @param \Throwable $e
     *
     * @return array<int, string>
     */
    private static function extractFrames(\Throwable $e): array
    {
        $rawTrace = $e->getTraceAsString();
        $lines = explode("\n", $rawTrace);
        $frames = [];

        foreach ($lines as $line) {
            $trimmedLine = trim($line);
            $hasContent = ($trimmedLine !== '');

            if ($hasContent) {
                $frames[] = $trimmedLine;
            }
        }

        return $frames;
    }
}
```

---

## 7.6 ResponseKeyType Enum — `includes/Enums/ResponseKeyType.php`

```php
<?php
/**
 * ResponseKeyType — Standard keys used in the API response envelope.
 *
 * @package PluginName\Enums
 * @since   1.0.0
 */

namespace PluginName\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum ResponseKeyType: string
{
    // Status block
    case Status      = 'Status';
    case IsSuccess   = 'IsSuccess';
    case IsFailed    = 'IsFailed';
    case Code        = 'Code';
    case Message     = 'Message';
    case Timestamp   = 'Timestamp';

    // Attributes block
    case Attributes   = 'Attributes';
    case RequestedAt  = 'RequestedAt';
    case TotalRecords = 'TotalRecords';

    // Data blocks
    case Results = 'Results';
    case Errors  = 'Errors';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

---

## 7.7 `.ai-instructions` File Template

Place this file at the plugin root. AI code generators read it before generating code.

```markdown

# AI Instructions for {Plugin Name}

## Architecture

- This plugin follows the spec at `02-spec/18-wp-plugin-how-to/`
- Read ALL 11 phases before writing any code

## Critical Rules

1. **No `is_array()`, `is_string()`, `is_int()`, etc.** — Use `TypeCheckerTrait` methods (`$this->isArray()`) or `PhpNativeType::matches()` in static contexts. The syntax validator blocks T_ARRAY tokens.
2. **No `array()` syntax** — Use `[]` exclusively.
3. **No `Exception` catches** — Always catch `Throwable`.
4. **No negative conditions** — Extract to `$is…`/`$has…`/`$should…` booleans.
5. **No magic strings** — Every string literal that represents a domain concept must be an enum case.
6. **No direct error responses** — Always use `EnvelopeBuilder` via `safeExecute()` or `validationError()`.
7. **Stack traces are debug-mode gated** — Never expose in production (see Phase 4, §4.2).
8. **All REST handlers wrapped in `safeExecute()`** — No bare try-catch in endpoints.
9. **Every file starts with ABSPATH guard** (after namespace, if namespaced).
10. **Version lives in `PluginConfigType::Version` only** — never hardcoded elsewhere.
11. **Seed data lives in `data/seeds/`** — JSON files + `manifest.json` define initial DB state (see Phase 8, §8.5.1). Never hardcode INSERT statements in migrations for reference data.
12. **File size limit: 200 lines max** — Templates, classes, traits, partials, JS, CSS. Extract into partials or sub-components when exceeded (see Phase 11, §11.1).
13. **Templates are orchestrators** — Page templates set variables and include partials. No business logic in templates (see Phase 11, §11.3).
14. **React requires developer confirmation** — Never default to React for admin UI without explicit developer approval (see Phase 11, §11.10).
15. **Source maps: dev only** — Production builds must NOT include `.map` files (see Phase 11, §11.8).

## File Generation Order

When creating a new plugin from scratch:
1. `plugin-name.php` (bootstrap)
2. `includes/Autoloader.php`
3. `includes/Enums/` (all enums first — they have no dependencies)
4. `includes/Helpers/` (DateHelper, PathHelper, ErrorLogHelper, EnvelopeBuilder)
5. `includes/Logging/FileLogger.php`
6. `includes/Traits/Core/` (ResponseTrait, TypeCheckerTrait)
7. `includes/Traits/Auth/AuthTrait.php`
8. `includes/Traits/Route/RouteRegistrationTrait.php`
9. `includes/Core/Plugin.php`
10. `includes/Core/Activator.php` + `Deactivator.php`
11. `includes/Database/DatabaseMigrator.php` + `DatabaseSeeder.php`
12. `data/seeds/manifest.json` + seed JSON files (initial reference data)
13. `uninstall.php`
14. `templates/` — page templates + `partials/shared/` (page-header, pagination, etc.)
15. `assets/css/` + `assets/js/` — per-page styles and scripts
16. Feature-domain traits (one per endpoint)
17. *(Optional, if React confirmed)* `frontend/` — React source → builds to `assets/dist/`

## Formatting Rules

- R1: Always use braces, even for single-line if/foreach/while
- R4: Blank line before return/throw only when preceded by other statements
- R9a: 3+ parameters → one per line, trailing comma
- R9b: 3+ arguments → one per line, trailing comma
- R9c: 3+ array items → one per line, trailing comma
- R12: No empty line after opening brace
- R13: No empty first line after <?php
```

---

## 7.8 Complete PluginConfigType Example

```php
<?php
/**
 * PluginConfigType — The identity enum. Single source of truth for all plugin metadata.
 *
 * @package PluginName\Enums
 * @since   1.0.0
 */

namespace PluginName\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum PluginConfigType: string
{
    case Slug          = 'plugin-name';
    case ShortName     = 'PluginName';
    case Name          = 'Plugin Name';
    case Version       = '1.0.0';
    case MinWpVersion  = '5.6';
    case MinPhpVersion = '8.1';
    case ApiNamespace  = 'plugin-name-api';
    case ApiVersion    = 'v1';
    case LogPrefix     = '[PluginName]';
    case SettingsGroup = 'plugin_name_settings';
    case DebugConstant = 'PLUGIN_NAME_DEBUG';

    /** Build the full REST API namespace: 'plugin-name-api/v1' */
    public static function apiFullNamespace(): string
    {
        return self::ApiNamespace->value . '/' . self::ApiVersion->value;
    }

    /** Check if the plugin is running in debug mode. */
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
