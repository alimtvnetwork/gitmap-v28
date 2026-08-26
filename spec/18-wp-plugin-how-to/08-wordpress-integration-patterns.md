# Phase 8 — WordPress Integration Patterns

> **Purpose:** Define patterns for WordPress-specific features that go beyond REST API endpoints: admin pages, WP-Cron, AJAX handlers, file uploads, database migrations, and transient caching. Each pattern includes a complete implementation with edge cases.

---

## 8.1 Admin Pages & Settings

### When to use

Any plugin that needs a settings screen visible in the WordPress admin sidebar.

### Registration pattern

Create a trait: `Traits/Admin/AdminPageTrait.php`

```php
namespace PluginName\Traits\Admin;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\PluginConfigType;
use PluginName\Enums\CapabilityType;

trait AdminPageTrait
{
    /**
     * Register admin menu pages. Call from Plugin::__construct() via
     * add_action('admin_menu', [$this, 'registerAdminPages']);
     */
    public function registerAdminPages(): void
    {
        add_menu_page(
            PluginConfigType::Name->value,                       // Page title
            PluginConfigType::ShortName->value,                  // Menu title
            CapabilityType::ManageOptions->value,                // Capability
            PluginConfigType::Slug->value,                       // Menu slug
            [$this, 'renderSettingsPage'],                       // Callback
            'dashicons-admin-generic',                           // Icon
            80,                                                  // Position
        );
    }

    /**
     * Render the main settings page.
     */
    public function renderSettingsPage(): void
    {
        $isAuthorized = current_user_can(CapabilityType::ManageOptions->value);

        if (!$isAuthorized) {
            wp_die('Unauthorized access');
        }

        // Use a template file instead of inline HTML
        $templatePath = plugin_dir_path(dirname(__DIR__, 2)) . 'templates/settings.php';
        $hasTemplate = file_exists($templatePath);

        if ($hasTemplate) {
            include $templatePath;
        }
    }
}
```

### Settings API pattern

```php
/**
 * Register settings fields. Call from Plugin::__construct() via
 * add_action('admin_init', [$this, 'registerSettings']);
 */
public function registerSettings(): void
{
    register_setting(
        PluginConfigType::SettingsGroup->value,
        'plugin_name_settings',
        [
            'type'              => 'object',
            'sanitize_callback' => [$this, 'sanitizeSettings'],
        ],
    );

    add_settings_section(
        'plugin_name_general',
        'General Settings',
        null,
        PluginConfigType::Slug->value,
    );

    add_settings_field(
        'plugin_name_api_key',
        'API Key',
        [$this, 'renderApiKeyField'],
        PluginConfigType::Slug->value,
        'plugin_name_general',
    );
}

/**
 * Sanitize settings before saving.
 *
 * @param array<string, mixed> $input Raw form input
 *
 * @return array<string, mixed> Sanitized settings
 */
public function sanitizeSettings(array $input): array
{
    $sanitized = [];

    $apiKey = $input['api_key'] ?? '';
    $hasApiKey = ($apiKey !== '');

    if ($hasApiKey) {
        $sanitized['api_key'] = sanitize_text_field($apiKey);
    }

    return $sanitized;
}
```

### Nonce verification for admin forms

```php
// In the form template:
wp_nonce_field('plugin_name_save_settings', 'plugin_name_nonce');

// In the handler:
$nonce = $_POST['plugin_name_nonce'] ?? '';
$isValidNonce = wp_verify_nonce($nonce, 'plugin_name_save_settings');

if (!$isValidNonce) {
    wp_die('Security check failed');
}
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| Plugin activated on multisite | Use `is_multisite()` check; register network admin pages with `network_admin_menu` |
| User lacks capability | `wp_die()` with friendly message — never show partial page |
| Settings page renders with PHP warnings | Wrap all rendering in try-catch, log errors, show fallback message |

### Admin Menu Error Count Badge

When the plugin has unseen errors, display a count badge on the admin menu item. This uses WordPress's built-in `<span class="update-plugins">` pattern.

```php
trait AdminPageTrait
{
    public function registerAdminPages(): void
    {
        // Get unseen error count for badge
        $unseenCount = $this->getUnseenErrorCount();
        $menuTitle = PluginConfigType::ShortName->value;

        $hasBadge = ($unseenCount > 0);

        if ($hasBadge) {
            $menuTitle .= sprintf(
                ' <span class="update-plugins count-%d"><span class="plugin-count">%d</span></span>',
                $unseenCount,
                $unseenCount,
            );
        }

        add_menu_page(
            PluginConfigType::Name->value,
            $menuTitle,                                            // Menu title with badge
            CapabilityType::ManageOptions->value,
            PluginConfigType::Slug->value,
            [$this, 'renderSettingsPage'],
            'dashicons-admin-generic',
            80,
        );
    }

    /**
     * Get count of unseen error sessions.
     * Uses wp_options for fast retrieval without DB query on every admin page load.
     */
    private function getUnseenErrorCount(): int
    {
        $optionKey = PluginConfigType::Slug->value . '_unseen_error_count';
        $count = get_option($optionKey, 0);

        return (int) $count;
    }
}
```

#### Badge update flow

| Event | Action |
|-------|--------|
| New error logged | `update_option($optionKey, $currentCount + 1)` in error handler |
| Admin views error page | `update_option($optionKey, 0)` — clears badge |
| Flash banner dismissed | `update_option($optionKey, 0)` — clears badge via AJAX |
| Plugin deactivated | `delete_option($optionKey)` in `Deactivator` |

#### Rules

1. Badge count stored in `wp_options` — NOT computed via DB query on every page load
2. The `update-plugins` class is a WordPress convention that styles the red bubble automatically
3. Badge is cleared when admin navigates to the error page OR dismisses the flash banner
4. The `count-{N}` class is required for WordPress core CSS to render correctly

---

## 8.2 AJAX Handlers (Non-REST)

### When to use

WordPress admin AJAX calls via `wp_ajax_{action}` — typically for admin UI interactions that don't need REST API formality.

### Pattern

```php
namespace PluginName\Traits\Admin;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\AjaxActionType;
use PluginName\Enums\CapabilityType;

trait AdminAjaxTrait
{
    /**
     * Register AJAX handlers. Call from Plugin::__construct() via
     * add_action('admin_init', [$this, 'registerAjaxHandlers']);
     */
    public function registerAjaxHandlers(): void
    {
        // Authenticated admin AJAX only — no wp_ajax_nopriv_ (never allow public AJAX)
        add_action(
            'wp_ajax_' . AjaxActionType::ClearCache->value,
            [$this, 'handleClearCache'],
        );
    }

    /**
     * Handle the clear-cache AJAX request.
     */
    public function handleClearCache(): void
    {
        try {
            // ── 1. Verify nonce ──
            $nonce = $_POST['nonce'] ?? '';
            $isValidNonce = wp_verify_nonce($nonce, AjaxActionType::ClearCache->value);

            if (!$isValidNonce) {
                wp_send_json_error(['message' => 'Invalid security token'], 403);

                return;
            }

            // ── 2. Verify capability ──
            $isAuthorized = current_user_can(CapabilityType::ManageOptions->value);

            if (!$isAuthorized) {
                wp_send_json_error(['message' => 'Unauthorized'], 403);

                return;
            }

            // ── 3. Business logic ──
            delete_transient('plugin_name_cache');

            // ── 4. Success response ──
            wp_send_json_success(['message' => 'Cache cleared']);
        } catch (\Throwable $e) {
            $this->fileLogger->logException($e, 'AJAX:clear-cache');

            wp_send_json_error(['message' => 'An error occurred'], 500);
        }
    }
}
```

### AjaxActionType enum

```php
enum AjaxActionType: string
{
    case ClearCache    = 'plugin_name_clear_cache';
    case ExportData    = 'plugin_name_export_data';
    case RunDiagnostic = 'plugin_name_run_diagnostic';

    /** Build the nonce action string. */
    public function nonceAction(): string
    {
        return $this->value;
    }

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### JavaScript side (enqueue in admin page)

```javascript
jQuery.post(ajaxurl, {
    action: 'plugin_name_clear_cache',
    nonce: pluginNameData.nonce,   // Localized via wp_localize_script()
}, function(response) {
    if (response.success) {
        alert(response.data.message);
    }
});
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| Missing `ajaxurl` in JS | Always use `wp_localize_script()` to pass URL and nonce |
| Expired nonce (user left tab open) | Return 403 with "Security token expired. Please refresh the page." |
| AJAX called without login | Don't register `wp_ajax_nopriv_` for admin actions — WordPress returns 0 automatically |

---

## 8.3 WP-Cron — Scheduled Tasks

### When to use

Recurring background tasks: log rotation, cache cleanup, data sync, health checks.

### Registration pattern

```php
namespace PluginName\Traits\Cron;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\CronScheduleType;

trait CronSchedulerTrait
{
    /**
     * Register cron hooks. Call from Plugin::__construct().
     */
    public function registerCronHooks(): void
    {
        add_action(CronScheduleType::LogRotation->hookName(), [$this, 'executeLogRotation']);
        add_action(CronScheduleType::CacheCleanup->hookName(), [$this, 'executeCacheCleanup']);
    }

    /**
     * Schedule cron events. Called from Activator::activate().
     */
    public static function scheduleCronEvents(): void
    {
        foreach (CronScheduleType::cases() as $schedule) {
            $isAlreadyScheduled = (wp_next_scheduled($schedule->hookName()) !== false);

            if ($isAlreadyScheduled) {
                continue;
            }

            wp_schedule_event(time(), $schedule->recurrence(), $schedule->hookName());
        }
    }

    /**
     * Unschedule all cron events. Called from Deactivator::deactivate().
     */
    public static function unscheduleCronEvents(): void
    {
        foreach (CronScheduleType::cases() as $schedule) {
            $nextRun = wp_next_scheduled($schedule->hookName());
            $hasEvent = ($nextRun !== false);

            if ($hasEvent) {
                wp_unschedule_event($nextRun, $schedule->hookName());
            }
        }
    }

    /**
     * Execute log rotation task.
     */
    public function executeLogRotation(): void
    {
        try {
            $this->fileLogger->info('Cron: log rotation started');
            // ... rotation logic
            $this->fileLogger->info('Cron: log rotation complete');
        } catch (\Throwable $e) {
            $this->fileLogger->logException($e, 'Cron:log-rotation');
        }
    }

    /**
     * Execute cache cleanup task.
     */
    public function executeCacheCleanup(): void
    {
        try {
            $this->fileLogger->info('Cron: cache cleanup started');
            // ... cleanup logic
            $this->fileLogger->info('Cron: cache cleanup complete');
        } catch (\Throwable $e) {
            $this->fileLogger->logException($e, 'Cron:cache-cleanup');
        }
    }
}
```

### CronScheduleType enum

```php
enum CronScheduleType: string
{
    case LogRotation  = 'log_rotation';
    case CacheCleanup = 'cache_cleanup';

    /** WordPress hook name for this scheduled task. */
    public function hookName(): string
    {
        return 'plugin_name_cron_' . $this->value;
    }

    /** WordPress cron recurrence interval. */
    public function recurrence(): string
    {
        return match ($this) {
            self::LogRotation  => 'daily',
            self::CacheCleanup => 'hourly',
        };
    }

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Custom cron intervals

If WordPress's built-in intervals (`hourly`, `twicedaily`, `daily`, `weekly`) aren't enough:

```php
// In Plugin::__construct()
add_filter('cron_schedules', [$this, 'addCronIntervals']);

public function addCronIntervals(array $schedules): array
{
    $schedules['every_five_minutes'] = [
        'interval' => 300,
        'display'  => 'Every 5 Minutes',
    ];

    return $schedules;
}
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| WP-Cron disabled (`DISABLE_WP_CRON`) | Document that server-level cron must call `wp-cron.php` |
| Cron runs overlapping (long task) | Use a transient lock: `set_transient('lock', true, 300)` at start, check before running |
| Cron fires after plugin deactivated | Check `function_exists()` or class existence at hook callback start |
| Multiple cron events for same hook | `wp_next_scheduled()` only returns the next one — use `_get_cron_array()` to check all |

### Transient lock pattern (prevent overlapping runs)

```php
public function executeLogRotation(): void
{
    $lockKey = 'plugin_name_lock_log_rotation';
    $isLocked = (get_transient($lockKey) !== false);

    if ($isLocked) {
        $this->fileLogger->debug('Cron: log rotation skipped — already running');

        return;
    }

    // Acquire lock (5-minute TTL)
    set_transient($lockKey, true, 300);

    try {
        // ... rotation logic
    } catch (\Throwable $e) {
        $this->fileLogger->logException($e, 'Cron:log-rotation');
    } finally {
        delete_transient($lockKey);
    }
}
```

---

## 8.4 File Upload Handling

### REST endpoint for file uploads

```php
namespace PluginName\Traits\Upload;

if (!defined('ABSPATH')) {
    exit;
}

use WP_REST_Request;
use WP_REST_Response;
use PluginName\Helpers\EnvelopeBuilder;
use PluginName\Helpers\PathHelper;

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

    /**
     * Map PHP upload error codes to human-readable messages.
     */
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

### Route registration for file upload

```php
// In RouteRegistrationTrait — file uploads need special handling
$safeRegister(
    EndpointType::FileUpload->route(),
    [
        'methods'             => HttpMethodType::Post->value,
        'callback'            => [$this, 'handleFileUpload'],
        'permission_callback' => [$this, 'checkPluginPermission'],
    ],
);
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| PHP `upload_max_filesize` too low | Check `ini_get('upload_max_filesize')` and warn in health check |
| `post_max_size` smaller than file | Returns empty `$_FILES` — check for this before validation |
| Symlink attack on temp file | Use `move_uploaded_file()` (not `rename()`) — it validates the temp path |
| File name collision | Prepend timestamp or UUID: `time() . '_' . $fileName` |
| Binary file disguised as allowed type | Check MIME via `wp_check_filetype()` in addition to extension |

---

## 8.5 Database Migrations (SQLite)

### Migration runner pattern

```php
namespace PluginName\Database;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;
use PluginName\Helpers\PathHelper;
use PluginName\Helpers\ErrorLogHelper;

final class DatabaseMigrator
{
    private \SQLite3 $db;

    /** @var array<int, callable> Migrations keyed by version number. */
    private array $migrations = [];

    public function __construct()
    {
        $dbPath = PathHelper::getBaseDir() . '/plugin-name.db';
        PathHelper::ensureFileParentDirectory($dbPath);

        $this->db = new \SQLite3($dbPath);
        $this->db->enableExceptions(true);

        // Ensure migrations table exists
        $this->db->exec('
            CREATE TABLE IF NOT EXISTS migrations (
                version INTEGER PRIMARY KEY,
                applied_at TEXT NOT NULL
            )
        ');

        $this->registerMigrations();
    }

    /**
     * Register all migrations. Each migration is a callable that receives the SQLite3 instance.
     */
    private function registerMigrations(): void
    {
        $this->migrations[1] = function (\SQLite3 $db): void {
            $db->exec('
                CREATE TABLE IF NOT EXISTS widgets (
                    id TEXT PRIMARY KEY,
                    name TEXT NOT NULL,
                    slug TEXT NOT NULL UNIQUE,
                    created_at TEXT NOT NULL,
                    updated_at TEXT NOT NULL
                )
            ');
        };

        $this->migrations[2] = function (\SQLite3 $db): void {
            $db->exec('ALTER TABLE widgets ADD COLUMN is_active INTEGER NOT NULL DEFAULT 1');
        };

        // Add new migrations here with incrementing version numbers
    }

    /**
     * Run all pending migrations in order.
     */
    public function runPending(): void
    {
        $currentVersion = $this->getCurrentVersion();

        foreach ($this->migrations as $version => $migration) {
            $isAlreadyApplied = ($version <= $currentVersion);

            if ($isAlreadyApplied) {
                continue;
            }

            try {
                $this->db->exec('BEGIN TRANSACTION');

                $migration($this->db);

                $stmt = $this->db->prepare(
                    'INSERT INTO migrations (version, applied_at) VALUES (:version, :applied_at)'
                );
                $stmt->bindValue(':version', $version, SQLITE3_INTEGER);
                $stmt->bindValue(':applied_at', gmdate('c'), SQLITE3_TEXT);
                $stmt->execute();

                $this->db->exec('COMMIT');
            } catch (Throwable $e) {
                $this->db->exec('ROLLBACK');

                ErrorLogHelper::logAndThrow($e, "Migration v{$version}:");
            }
        }
    }

    /**
     * Get the highest applied migration version.
     */
    private function getCurrentVersion(): int
    {
        $result = $this->db->querySingle('SELECT MAX(version) FROM migrations');
        $hasVersion = ($result !== null && $result !== false);

        return $hasVersion ? (int) $result : 0;
    }

    /**
     * Static entry point for Activator.
     */
    public static function runAllPending(): void
    {
        $migrator = new self();
        $migrator->runPending();
    }
}
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| SQLite3 extension not installed | Check `extension_loaded('sqlite3')` in Activator, fail gracefully |
| Migration fails midway | Transaction + ROLLBACK ensures no partial schema changes |
| Database file locked (concurrent access) | SQLite handles this with WAL mode; enable via `$db->exec('PRAGMA journal_mode=WAL')` |
| Database file deleted while plugin active | Recreate on next access — check file exists before opening |
| Adding column to existing table | SQLite doesn't support `DROP COLUMN` — use `ALTER TABLE ADD COLUMN` only |

## 8.5.1 Database Seeding — Initial Data from JSON Files

### Concept

Seeding is the process of populating database tables with **default data** on first activation and **updated data** on version upgrades. Seed data lives in JSON files inside the plugin source — not hardcoded in PHP. This allows non-developers to review and modify defaults without touching code.

### Folder structure

```
plugin-slug/
├── data/
│   └── seeds/
│       ├── manifest.json          ← Declares which seeds exist and their strategy
│       ├── settings.json          ← Default plugin settings
│       ├── templates.json         ← Default templates / presets
│       └── permissions.json       ← Default role-capability mappings
```

### manifest.json — Seed Registry

The manifest declares every seed file, which table it targets, and the strategy to apply. The optional `version` field is **metadata only** — it documents which plugin version introduced the seed but is **not used by the seeder logic**. Version tracking is handled per-file via the `seed_history` table (see `getLastSeededVersion()`).

```json
{
  "seeds": [
    {
      "file": "settings.json",
      "table": "settings",
      "version": "1.0.0",
      "strategy": "insert_if_empty"
    },
    {
      "file": "templates.json",
      "table": "templates",
      "version": "1.0.0",
      "strategy": "insert_if_empty"
    },
    {
      "file": "permissions.json",
      "table": "permissions",
      "version": "1.2.0",
      "strategy": "upsert_by_key"
    }
  ]
}
```

> **Note:** The `version` field is kept as human-readable metadata so developers can see when each seed was introduced. The seeder ignores it — re-seeding is triggered solely by comparing the plugin's current version against `seed_history.last_seeded_ver` per file.

### Seeding strategies

| Strategy | Behaviour | Use when |
|----------|-----------|----------|
| `insert_if_empty` | Only seeds if the target table has zero rows | First activation — don't overwrite user customisations |
| `upsert_by_key` | Inserts new rows, updates existing rows by primary key | Version upgrade adds new defaults without wiping user data |
| `replace_all` | Truncates table and re-inserts all rows | Non-user-editable reference data (e.g., error codes, system constants) |

### Seed file format

Each JSON file is an array of row objects. Keys match column names:

```json
[
  {
    "key": "log_level",
    "value": "info",
    "description": "Minimum log severity to persist",
    "is_user_editable": true
  },
  {
    "key": "max_upload_size_mb",
    "value": "50",
    "description": "Maximum file upload size in megabytes",
    "is_user_editable": true
  },
  {
    "key": "api_version",
    "value": "v1",
    "description": "Current REST API version",
    "is_user_editable": false
  }
]
```

### DatabaseSeeder — Complete Implementation

```php
namespace PluginName\Database;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;
use PluginName\Enums\PluginConfigType;
use PluginName\Enums\SeedStrategyType;
use PluginName\Helpers\PathHelper;
use PluginName\Helpers\ErrorLogHelper;

final class DatabaseSeeder
{
    private \SQLite3 $db;
    private string $seedsDir;

    /**
     * @param \SQLite3    $db       Database connection
     * @param string|null $seedsDir Path to seeds directory (default: plugin's data/seeds/)
     */
    public function __construct(\SQLite3 $db, ?string $seedsDir = null)
    {
        $this->db = $db;
        $this->seedsDir = $seedsDir ?? plugin_dir_path(dirname(__DIR__)) . 'data/seeds';
    }

    /**
     * Run all pending seeds based on the manifest and current version.
     * Called after migrations complete (schema must exist before data).
     *
     * @param string|null $version Override version (default: reads from PluginConfigType::Version)
     */
    public function seedAll(?string $version = null): void
    {
        $manifestPath = $this->seedsDir . '/manifest.json';
        $hasManifest = file_exists($manifestPath);

        if (!$hasManifest) {
            return;
        }

        $manifestRaw = file_get_contents($manifestPath);
        $manifest = json_decode($manifestRaw, true);
        $hasSeeds = (gettype($manifest) === 'array' && isset($manifest['seeds']));

        if (!$hasSeeds) {
            return;
        }

        $currentVersion = $version ?? PluginConfigType::Version->value;

        $this->ensureSeedHistoryTable();

        foreach ($manifest['seeds'] as $seedEntry) {
            $this->processSeedEntry($seedEntry, $currentVersion);
        }
    }

    /**
     * Process a single seed entry from the manifest.
     *
     * @param array<string, string> $entry          Manifest entry
     * @param string                $currentVersion Current plugin version
     */
    private function processSeedEntry(array $entry, string $currentVersion): void
    {
        $seedFile = $entry['file'] ?? '';
        $table = $entry['table'] ?? '';
        $strategyValue = $entry['strategy'] ?? 'insert_if_empty';

        // Per-file version tracking — skip if already seeded at this version or newer
        $lastSeeded = $this->getLastSeededVersion($seedFile);
        $isAlreadySeeded = (
            $lastSeeded !== null
            && version_compare($currentVersion, $lastSeeded, '<=')
        );

        if ($isAlreadySeeded) {
            return;
        }

        $filePath = $this->seedsDir . '/' . $seedFile;
        $hasFile = file_exists($filePath);

        if (!$hasFile) {
            error_log("[Seeder] Seed file not found: {$filePath}");

            return;
        }

        $strategy = SeedStrategyType::tryFrom($strategyValue);
        $hasStrategy = ($strategy !== null);

        if (!$hasStrategy) {
            error_log("[Seeder] Unknown strategy: {$strategyValue}");

            return;
        }

        try {
            $rawData = file_get_contents($filePath);
            $rows = json_decode($rawData, true);
            $hasRows = (gettype($rows) === 'array' && count($rows) > 0);

            if (!$hasRows) {
                return;
            }

            $this->db->exec('BEGIN TRANSACTION');

            match ($strategy) {
                SeedStrategyType::InsertIfEmpty => $this->seedInsertIfEmpty($table, $rows),
                SeedStrategyType::UpsertByKey   => $this->seedUpsertByKey($table, $rows),
                SeedStrategyType::ReplaceAll    => $this->seedReplaceAll($table, $rows),
            };

            $this->db->exec('COMMIT');

            // Record per-file version after successful seed
            $this->setLastSeededVersion($seedFile, $currentVersion);
        } catch (Throwable $e) {
            $this->db->exec('ROLLBACK');

            ErrorLogHelper::log($e, "Seeder:{$seedFile}:");
        }
    }

    /**
     * Insert rows only if the table is completely empty.
     */
    private function seedInsertIfEmpty(string $table, array $rows): void
    {
        $count = $this->db->querySingle("SELECT COUNT(*) FROM {$table}");
        $hasExistingData = ($count > 0);

        if ($hasExistingData) {
            return;
        }

        $this->insertRows($table, $rows);
    }

    /**
     * Insert new rows, update existing rows by primary key.
     * Uses SQLite's INSERT OR REPLACE.
     */
    private function seedUpsertByKey(string $table, array $rows): void
    {
        $this->insertRows($table, $rows, 'INSERT OR REPLACE');
    }

    /**
     * Delete all existing rows and re-insert from seed file.
     */
    private function seedReplaceAll(string $table, array $rows): void
    {
        $this->db->exec("DELETE FROM {$table}");
        $this->insertRows($table, $rows);
    }

    /**
     * Insert an array of rows into a table.
     *
     * @param string                          $table  Target table name
     * @param array<int, array<string, mixed>> $rows   Row data
     * @param string                          $verb   SQL verb (INSERT or INSERT OR REPLACE)
     */
    private function insertRows(string $table, array $rows, string $verb = 'INSERT'): void
    {
        $columns = array_keys($rows[0]);
        $columnList = implode(', ', $columns);
        $placeholders = implode(', ', array_map(fn($c) => ":{$c}", $columns));

        $sql = "{$verb} INTO {$table} ({$columnList}) VALUES ({$placeholders})";
        $stmt = $this->db->prepare($sql);

        foreach ($rows as $row) {
            foreach ($columns as $column) {
                $value = $row[$column] ?? null;
                $stmt->bindValue(":{$column}", $value);
            }

            $stmt->execute();
            $stmt->reset();
        }
    }

    /**
     * Ensure the seed_history table exists.
     */
    private function ensureSeedHistoryTable(): void
    {
        $this->db->exec('
            CREATE TABLE IF NOT EXISTS seed_history (
                id              INTEGER PRIMARY KEY AUTOINCREMENT,
                seed_file       TEXT    NOT NULL UNIQUE,
                last_seeded_ver TEXT    NOT NULL,
                seeded_at       TEXT    NOT NULL DEFAULT CURRENT_TIMESTAMP
            )
        ');
    }

    /**
     * Get the last version a specific seed file was applied at.
     */
    private function getLastSeededVersion(string $seedFile): ?string
    {
        $stmt = $this->db->prepare(
            "SELECT last_seeded_ver FROM seed_history WHERE seed_file = :file"
        );
        $stmt->bindValue(':file', $seedFile, SQLITE3_TEXT);
        $result = $stmt->execute()->fetchArray(SQLITE3_ASSOC);

        $hasResult = (gettype($result) === 'array' && isset($result['last_seeded_ver']));

        return $hasResult ? $result['last_seeded_ver'] : null;
    }

    /**
     * Record the version a specific seed file was applied at.
     */
    private function setLastSeededVersion(string $seedFile, string $version): void
    {
        $stmt = $this->db->prepare(
            "INSERT OR REPLACE INTO seed_history (seed_file, last_seeded_ver)
             VALUES (:file, :version)"
        );
        $stmt->bindValue(':file', $seedFile, SQLITE3_TEXT);
        $stmt->bindValue(':version', $version, SQLITE3_TEXT);
        $stmt->execute();
    }
}
```

### SeedStrategyType Enum

```php
namespace PluginName\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum SeedStrategyType: string
{
    case InsertIfEmpty = 'insert_if_empty';
    case UpsertByKey   = 'upsert_by_key';
    case ReplaceAll    = 'replace_all';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Integration with DatabaseMigrator

Seeding runs **after** migrations (schema must exist before data is inserted):

```php
// In DatabaseMigrator::runAllPending() — add seeding after migrations
public static function runAllPending(): void
{
    $migrator = new self();
    $migrator->runPending();

    // Seed data after schema is up to date
    $seeder = new DatabaseSeeder($migrator->db);
    $seeder->seedAll();  // Uses PluginConfigType::Version by default
}
```

### Version upgrade flow

```
1. User updates plugin from v1.0.0 to v1.2.0
2. WordPress calls register_activation_hook → Activator::activate()
3. Activator calls DatabaseMigrator::runAllPending()
4. Migrator runs migration v2 (adds is_active column)
5. Migrator calls DatabaseSeeder::seedAll()
6. Seeder reads manifest.json and checks per-file seed_history:
   - settings.json → seed_history shows last_seeded_ver "1.0.0" < "1.2.0" → RE-SEED
   - permissions.json → no seed_history entry → SEED (first run for this file)
7. Seeder records last_seeded_ver = "1.2.0" per file in seed_history table
8. Next activation at v1.2.0: seedAll() sees per-file versions unchanged → skips all
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| Seed file has invalid JSON | `json_decode` returns null → logged, skipped, no crash |
| Seed file references non-existent table | SQLite throws → caught by transaction rollback |
| User modified seeded data, then plugin upgrades | `insert_if_empty` preserves user data; `upsert_by_key` adds new rows but updates existing keys; `replace_all` only for system data the user should never edit |
| Downgrade (v1.2.0 → v1.0.0) | `lastSeededVersion` is higher than manifest entries → all seeds skipped → safe |
| Manifest missing | `seedAll()` returns immediately — no error |
| Seed file added in new version but missing from manifest | Not processed — all seeds must be declared in manifest |
| Empty seed file (empty JSON array) | Detected by `count($rows) === 0` → skipped |

---

## 8.6 Transient Caching

### When to use

Cache expensive operations (remote API calls, complex queries) that don't need real-time accuracy.

### Pattern

```php
/**
 * Get data with transient caching.
 *
 * @param string   $key      Transient key (auto-prefixed)
 * @param int      $ttl      Cache duration in seconds
 * @param callable $callback Function to compute the value if cache misses
 *
 * @return mixed The cached or freshly computed value
 */
protected function getCached(string $key, int $ttl, callable $callback): mixed
{
    $fullKey = 'plugin_name_' . $key;
    $cached = get_transient($fullKey);
    $hasCached = ($cached !== false);

    if ($hasCached) {
        $this->fileLogger->debug('Cache hit', ['key' => $key]);

        return $cached;
    }

    $this->fileLogger->debug('Cache miss — computing', ['key' => $key]);

    $value = $callback();
    set_transient($fullKey, $value, $ttl);

    return $value;
}

/**
 * Invalidate a specific cache key.
 */
protected function invalidateCache(string $key): void
{
    $fullKey = 'plugin_name_' . $key;
    delete_transient($fullKey);

    $this->fileLogger->debug('Cache invalidated', ['key' => $key]);
}
```

### Usage in a handler

```php
private function executeGetStats(WP_REST_Request $request): WP_REST_Response
{
    $stats = $this->getCached('dashboard_stats', 300, function () {
        // Expensive computation — cached for 5 minutes
        return [
            'totalWidgets'  => $this->countWidgets(),
            'activeWidgets' => $this->countActiveWidgets(),
            'computedAt'    => DateHelper::nowUtc(),
        ];
    });

    return EnvelopeBuilder::success()
        ->setRequestedAt($request->get_route())
        ->setSingleResult($stats)
        ->toResponse();
}
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| Object caching plugin installed | `get_transient` automatically uses object cache — no changes needed |
| Transient stores `false` as value | Wrap value in array: `['data' => $value]` to distinguish from cache miss |
| Cache stampede (many requests hit miss simultaneously) | Use transient lock (same pattern as cron lock in §8.3) |
| Multisite | Use `get_site_transient()` / `set_site_transient()` for network-wide caching |

---

## 8.7 HTTP Requests to External APIs

### Pattern with error handling

```php
/**
 * Make a GET request to an external API with structured error handling.
 *
 * @param string               $url     The full URL to request
 * @param array<string, string> $headers Additional headers
 *
 * @return array{success: bool, data: mixed, error: string|null}
 */
protected function externalGet(string $url, array $headers = []): array
{
    $response = wp_remote_get($url, [
        'timeout' => 15,
        'headers' => $headers,
    ]);

    $isWpError = is_wp_error($response);

    if ($isWpError) {
        $errorMessage = $response->get_error_message();
        $this->fileLogger->error('External API request failed', [
            'url'   => $url,
            'error' => $errorMessage,
        ]);

        return ['success' => false, 'data' => null, 'error' => $errorMessage];
    }

    $statusCode = wp_remote_retrieve_response_code($response);
    $body = wp_remote_retrieve_body($response);
    $isSuccess = ($statusCode >= 200 && $statusCode < 300);

    if (!$isSuccess) {
        $this->fileLogger->warn('External API returned non-2xx', [
            'url'    => $url,
            'status' => $statusCode,
            'body'   => mb_substr($body, 0, 500),
        ]);

        return ['success' => false, 'data' => $body, 'error' => "HTTP {$statusCode}"];
    }

    $decoded = json_decode($body, true);
    $isJson = ($decoded !== null);

    return [
        'success' => true,
        'data'    => $isJson ? $decoded : $body,
        'error'   => null,
    ];
}
```

### Edge cases

| Scenario | Handling |
|----------|----------|
| SSL certificate issues | Don't disable SSL verification — fix the server config |
| Timeout | Set explicit `timeout` in `wp_remote_get` args (default 5s is often too low) |
| WordPress HTTP API blocked | Some hosts block `wp_remote_get` — check `WP_HTTP_BLOCK_EXTERNAL` constant |
| Response is not JSON | Check content type or use `json_decode` return value to detect |
| Rate limiting (429) | Respect `Retry-After` header; log and return structured error |

---

## 8.8 Integration into Plugin.php

When using admin pages, cron, and AJAX alongside REST endpoints, the Plugin constructor changes:

```php
private function __construct()
{
    $startTime = microtime(true);
    $this->fileLogger = FileLogger::getInstance();

    // ── REST routes (only on REST requests) ──
    $this->registerRoutes();

    // ── Admin pages (only in admin context) ──
    $isAdmin = is_admin();

    if ($isAdmin) {
        add_action('admin_menu', [$this, 'registerAdminPages']);
        add_action('admin_init', [$this, 'registerSettings']);
        add_action('admin_init', [$this, 'registerAjaxHandlers']);
    }

    // ── Cron hooks (always — cron runs outside admin) ──
    $this->registerCronHooks();

    // ── Shutdown handler ──
    $this->registerShutdownHandler();

    // ── Boot log ──
    $elapsedMs = round((microtime(true) - $startTime) * 1000, 2);
    $this->fileLogger->info('Plugin initialized', [
        'version' => PluginConfigType::Version->value,
        'timeMs'  => $elapsedMs,
        'isAdmin' => $isAdmin,
    ]);
}
```

### Bootstrap file change

When the plugin needs admin features, change the hook from `rest_api_init` to `plugins_loaded`:

```php
// In plugin-name.php — use 'plugins_loaded' for full-featured plugins
add_action('plugins_loaded', function (): void {
    try {
        \PluginName\Core\Plugin::getInstance();
    } catch (\Throwable $e) {
        error_log('[PluginName] Boot failed: ' . $e->getMessage() . "\n" . $e->getTraceAsString());
    }
});
```

### Why `plugins_loaded` instead of `rest_api_init`

| Hook | When it fires | Use for |
|------|--------------|---------|
| `rest_api_init` | Only on REST API requests | REST-only plugins (no admin UI) |
| `plugins_loaded` | Every WordPress request | Plugins with admin pages, cron, AJAX, and REST |
| `init` | After `plugins_loaded` | Only if you need custom post types or taxonomies registered |

---

## 8.9 Complete Enum Inventory for a Full-Featured Plugin

| Enum | File | Cases (minimum) |
|------|------|-----------------|
| `PluginConfigType` | `Enums/PluginConfigType.php` | Slug, ShortName, Name, Version, MinWpVersion, MinPhpVersion, ApiNamespace, ApiVersion, LogPrefix, SettingsGroup, DebugConstant |
| `EndpointType` | `Enums/EndpointType.php` | Status, Plugins, Activate, Deactivate, Upload, FileUpload, Logs, Settings |
| `HttpMethodType` | `Enums/HttpMethodType.php` | Get, Post, Put, Delete |
| `HttpStatusType` | `Enums/HttpStatusType.php` | Ok, Created, BadRequest, Unauthorized, Forbidden, NotFound, ServerError |
| `HookType` | `Enums/HookType.php` | RestApiInit, PluginsLoaded, AdminMenu, AdminInit |
| `LogLevelType` | `Enums/LogLevelType.php` | Debug, Info, Warn, Error |
| `ResponseKeyType` | `Enums/ResponseKeyType.php` | Status, IsSuccess, IsFailed, Code, Message, Timestamp, Attributes, RequestedAt, TotalRecords, Results, Errors |
| `CapabilityType` | `Enums/CapabilityType.php` | ActivatePlugins, ManageOptions |
| `WpErrorCodeType` | `Enums/WpErrorCodeType.php` | Unauthorized, Forbidden, InvalidCredentials |
| `PathLogFileType` | `Enums/PathLogFileType.php` | Info, Error, Stacktrace, Autoloader, Fatal |
| `PhpNativeType` | `Enums/PhpNativeType.php` | PhpArray, PhpString, PhpInteger, PhpDouble, PhpBoolean, PhpObject, PhpNull |
| `AjaxActionType` | `Enums/AjaxActionType.php` | ClearCache, ExportData, RunDiagnostic |
| `CronScheduleType` | `Enums/CronScheduleType.php` | LogRotation, CacheCleanup |
| `SeedStrategyType` | `Enums/SeedStrategyType.php` | InsertIfEmpty, UpsertByKey, ReplaceAll |
