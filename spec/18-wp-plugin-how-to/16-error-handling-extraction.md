# Phase 16 — Error Handling & Diagnostics Extraction

> **Purpose:** Define the complete error handling, error log viewing, error session management, and PHP error classification patterns extracted from the RiseUpAsia codebase. Supplements Phase 4 (Logging) with admin-facing error management UI and API patterns.
> **Audience:** AI code generators and human developers.
> **Prerequisite:** Phases 1–4 must be read first.

---

## 16.1 Error Type Classification

PHP errors are classified into three severity groups using a dedicated `ErrorType` class (NOT a backed enum — it holds arrays of constants):

```php
final class ErrorType
{
    public const FATAL_TYPES = [
        E_ERROR, E_PARSE, E_CORE_ERROR, E_COMPILE_ERROR, E_USER_ERROR,
    ];

    public const WARNING_TYPES = [
        E_WARNING, E_CORE_WARNING, E_USER_WARNING,
        E_NOTICE, E_USER_NOTICE, E_DEPRECATED, E_USER_DEPRECATED,
    ];

    public const RECOVERABLE_TYPES = [
        E_RECOVERABLE_ERROR, E_STRICT,
    ];

    public const TYPE_LABELS = [
        E_ERROR => 'E_ERROR',
        // ... one entry per constant
    ];
}
```

### Why a final class instead of an enum

- Each case would need to hold an **array** of PHP `E_*` constants — backed enums only support `string|int`.
- The class groups related constants; individual error codes are not discrete enum cases.
- `TYPE_LABELS` provides human-readable names for display in admin UI tables.

### Rules

| Rule | Detail |
|------|--------|
| Fatal detection | `in_array($errno, ErrorType::FATAL_TYPES, true)` |
| Label lookup | `ErrorType::TYPE_LABELS[$errno] ?? 'UNKNOWN'` |
| No instantiation | Class is `final` with only `public const` members |
| Namespace | `RiseupAsia\Enums` (lives alongside real enums for discoverability) |

---

## 16.2 Two-Tier Error Capture

### Tier 1 — Bootstrap Errors (before autoloader)

```php
// In InitHelpers (available from the main plugin file)
public static function errorLogWithPrefix(string $message): void {
    error_log(PluginConfigType::LogPrefix->value . ' ' . $message);
}

public static function errorLog(Throwable $e, string $context): void {
    error_log($context . ' ' . $e->getMessage() . "\n" . $e->getTraceAsString());
}
```

**When to use:** Only during bootstrap, activation hooks, or when `FileLogger` is not yet available.

### Tier 2 — FileLogger Errors (after initialization)

All post-bootstrap errors go through `FileLogger` which writes to structured log files with rotation and deduplication (see Phase 4).

---

## 16.3 Error Log Retrieval API

The plugin exposes two diagnostic endpoints for error management:

### Endpoints

| Endpoint | Method | Purpose |
|----------|--------|---------|
| `error-logs` | GET | Query PHP error log entries with configurable inclusion |
| `error-sessions` | GET | Query grouped error sessions |

### Configurable Settings

Error log retrieval is controlled by `OptionNameType::LogRetrieval` settings:

```php
$defaults = [
    'include_error_log'  => true,   // Include error.log content
    'include_full_log'   => false,  // Include info.log content
    'include_stacktrace' => true,   // Include stack trace data
    'max_lines'          => 200,    // Maximum log lines to return
];
```

### Resolution Order

Settings are resolved in this priority:

1. **Request parameters** — Query params override stored settings per-request
2. **Stored settings** — `OptionNameType::LogRetrieval` from `wp_options`
3. **Defaults** — Hardcoded fallbacks above

```php
private function resolveSettings(WP_REST_Request $request): array {
    $logSettings = get_option(OptionNameType::LogRetrieval->value, []);

    $resolved = [
        'include_error_log'  => isset($logSettings['include_error_log'])
            ? (bool) $logSettings['include_error_log'] : true,
        // ... repeat for each key
    ];

    // Per-request overrides
    foreach (['include_error_log', 'include_full_log', 'include_stacktrace'] as $key) {
        if ($request->get_param($key) !== null) {
            $resolved[$key] = (bool) $request->get_param($key);
        }
    }

    return $resolved;
}
```

---

## 16.4 Error Session Model

Errors are grouped into **sessions** — a session represents a single request that produced one or more errors. This enables the admin UI to show errors in context rather than as isolated log lines.

### Session Data Structure

| Field | Type | Description |
|-------|------|-------------|
| `session_id` | string | Unique identifier (UUID or timestamp-based) |
| `created_at` | datetime | When the session started |
| `error_count` | int | Number of errors in this session |
| `is_seen` | boolean | Whether admin has dismissed/acknowledged |
| `errors` | array | Individual error entries within the session |

### Admin Actions

| Action | AJAX Handler | Description |
|--------|-------------|-------------|
| Dismiss flash | `dismissFlash` | Mark all unseen errors as seen (badge clears) |
| Clear all | `clearSessions` | Delete all error sessions |
| View details | Modal | Show full error context in modal overlay |

---

## 16.5 Error Admin Page JavaScript Pattern

The error management page follows the standard localized-object pattern:

```javascript
jQuery(document).ready(function($) {
    var C = window.RiseupErrors;
    var ajaxNonce = C.nonce;
    var activeTab = C.activeTab;
    var autoRefreshTimer = null;

    // Flash banner dismiss
    $('#riseup-dismiss-flash').on('click', function() {
        $.post(ajaxurl, {
            action: C.actions.dismissFlash,
            nonce: ajaxNonce
        }, function(response) {
            if (response.success) {
                $('#riseup-flash-banner').slideUp(300);
                $('.tab-badge, .error-count-badge').fadeOut(200);
            }
        });
    });
});
```

### Localized Object Shape (`RiseupErrors`)

```javascript
window.RiseupErrors = {
    nonce: '...',
    activeTab: 'errors',           // Current active tab
    actions: {
        dismissFlash: 'riseup_dismiss_error_flash',
        clearSessions: 'riseup_clear_error_sessions',
    },
    i18n: {
        dismissing: 'Dismissing...',
        markAsSeen: 'Mark as Seen',
        confirmClearAll: 'Are you sure you want to clear all error sessions?',
    }
};
```

---

## 16.6 Flash Banner Pattern

Unseen errors trigger a **flash banner** at the top of the error admin page:

### Requirements

| Requirement | Implementation |
|-------------|----------------|
| Visibility | Show only when `unseen_count > 0` |
| Badge | Display count in tab badge AND menu badge |
| Dismiss | Single click marks all as seen via AJAX |
| Animation | `slideUp(300)` on dismiss, `fadeOut(200)` on badges |
| Persistence | State stored in database, not session/cookie |

---

## 16.7 Auto-Refresh for Error Pages

Error pages support automatic polling for new errors:

```javascript
var autoRefreshTimer = null;

function startAutoRefresh(intervalMs) {
    stopAutoRefresh();
    autoRefreshTimer = setInterval(function() {
        loadErrors();  // Re-fetch and re-render table
    }, intervalMs);
}

function stopAutoRefresh() {
    if (autoRefreshTimer) {
        clearInterval(autoRefreshTimer);
        autoRefreshTimer = null;
    }
}
```

### Rules

- Default interval: **30 seconds**
- Auto-refresh MUST stop when modal is open
- Auto-refresh MUST stop when user is interacting (e.g., selecting text)
- Toggle state persists via a UI switch (see Phase 13 — Toggle Switch)

---

## 16.8 safeExecute Wrapper

All REST endpoint handlers MUST use the `safeExecute` wrapper to catch exceptions and return standardized error envelopes:

```php
public function handleErrorLogs(WP_REST_Request $request): WP_REST_Response {
    return $this->safeExecute(function() use ($request) {
        // ... handler logic
        return EnvelopeBuilder::success()
            ->autoDetectRequestedAt()
            ->setSingleResult($result)
            ->toResponse();
    }, 'error_logs');  // Context label for logging
}
```

### What safeExecute provides

| Feature | Detail |
|---------|--------|
| Exception catch | Wraps callback in try/catch, returns error envelope on failure |
| Context label | Second argument used in log messages for traceability |
| Debug gating | Stack traces included in response only when debug mode is ON |
| Consistent shape | All responses use `EnvelopeBuilder` regardless of success/failure |

---

## 16.9 Error Notification Settings

Stored under `OptionNameType::ErrorNotification`:

```php
case ErrorNotification = 'RiseupErrorNotificationSettings';
```

### Configurable Fields

| Field | Type | Default | Description |
|-------|------|---------|-------------|
| `email_enabled` | bool | false | Send email on fatal errors |
| `email_address` | string | admin_email | Recipient address |
| `threshold` | int | 5 | Minimum errors before notification |
| `cooldown_minutes` | int | 60 | Minimum time between notifications |

---

## 16.10 ErrorResponse — Consolidated Error Return Helper

Non-endpoint catch blocks (helpers, services, traits that don't use `safeExecute`) need a consistent way to log + return an error value. The `ErrorResponse` class provides static methods for every return-type scenario.

### Location

`includes/ErrorHandling/ErrorResponse.php` — namespace `PluginName\ErrorHandling`.

### Full implementation pattern

```php
namespace PluginName\ErrorHandling;

use Throwable;
use WP_Error;
use WP_REST_Response;

use PluginName\Enums\HttpStatusType;
use PluginName\Helpers\ResultHelper;
use PluginName\Logging\FileLogger;

class ErrorResponse
{
    /** Log exception and return standardized error array. */
    public static function logAndReturn(
        FileLogger $logger,
        Throwable $e,
        string $context = '',
    ): array {
        $logger->logException($e, $context);

        return ResultHelper::errorFromException($e);
    }

    /** Log exception and return a WP_REST_Response error envelope. */
    public static function logAndReturnEnvelope(
        FileLogger $logger,
        Throwable $e,
        string $context = '',
        int $status = HttpStatusType::ServerError->value,
    ): WP_REST_Response {
        $logger->logException($e, $context);

        return new WP_REST_Response(
            ResultHelper::errorFromException($e),
            $status,
        );
    }

    /** Log exception and return a WP_Error object. */
    public static function logAndReturnWpError(
        FileLogger $logger,
        Throwable $e,
        string $context = '',
        string $code = 'InternalError',
        int $status = HttpStatusType::ServerError->value,
    ): WP_Error {
        $logger->logException($e, $context);

        return new WP_Error(
            $code,
            $e->getMessage(),
            ['status' => $status],
        );
    }

    /** Log exception and return false. */
    public static function logAndReturnFalse(
        FileLogger $logger,
        Throwable $e,
        string $context = '',
    ): false {
        $logger->logException($e, $context);

        return false;
    }
}
```

### When to use each method

| Method | Return type | Use case |
|--------|------------|----------|
| `logAndReturn()` | `array` | Internal service methods returning result arrays |
| `logAndReturnEnvelope()` | `WP_REST_Response` | Non-safeExecute REST responses (rare) |
| `logAndReturnWpError()` | `WP_Error` | WordPress hooks that expect WP_Error (e.g., `pre_update_option`) |
| `logAndReturnFalse()` | `false` | Boolean-return methods (e.g., `isValid()`, `canProceed()`) |

### Rules

1. **Every** non-endpoint catch block MUST use an `ErrorResponse` method — no bare `return false` after catch
2. The `$context` parameter should identify the method: `'MyService::processItem'`
3. `logAndReturn` always calls `$logger->logException()` which writes to all three log files (info, error, stacktrace)
4. For catch blocks where FileLogger is unavailable, use `ErrorLogHelper::log()` instead (see Phase 4, §4.11)

---

## 16.11 AdminErrorAjaxTrait — Log File Operations

Admin pages that display log files need AJAX handlers for reading, clearing, and bulk-clearing log files. This trait provides the standard pattern.

### Location

`includes/Admin/Traits/AdminErrorAjaxTrait.php`

### Full implementation pattern

```php
namespace PluginName\Admin\Traits;

use PluginName\Enums\AjaxActionType;
use PluginName\Enums\CapabilityType;
use PluginName\Helpers\PathHelper;

trait AdminErrorAjaxTrait
{
    /**
     * AJAX handler: Read a log file's content.
     * Action: wp_ajax_{plugin}_read_log_file
     */
    public function ajaxReadLogFile(): void
    {
        check_ajax_referer(AjaxActionType::ReadLogFile->nonceAction(), 'nonce');

        $hasPermission = current_user_can(CapabilityType::ManageOptions->value);
        if (!$hasPermission) {
            wp_send_json_error(['message' => 'Unauthorized'], 403);
        }

        $filename = sanitize_file_name($_POST['filename'] ?? '');
        $allowedFiles = ['error.log', 'info.log', 'stacktrace.log', 'fatal-errors.log'];
        $isAllowed = in_array($filename, $allowedFiles, true);

        if (!$isAllowed) {
            wp_send_json_error(['message' => 'Invalid log file'], 400);
        }

        $filePath = PathHelper::getLogsDir() . '/' . $filename;
        $fileExists = file_exists($filePath);
        $content = $fileExists ? file_get_contents($filePath) : '';

        wp_send_json_success([
            'filename' => $filename,
            'content'  => $content,
            'size'     => $fileExists ? filesize($filePath) : 0,
            'isEmpty'  => ($content === '' || $content === false),
        ]);
    }

    /**
     * AJAX handler: Clear a single log file.
     * Action: wp_ajax_{plugin}_clear_log_file
     */
    public function ajaxClearLogFile(): void
    {
        check_ajax_referer(AjaxActionType::ClearLogFile->nonceAction(), 'nonce');

        $hasPermission = current_user_can(CapabilityType::ManageOptions->value);
        if (!$hasPermission) {
            wp_send_json_error(['message' => 'Unauthorized'], 403);
        }

        $filename = sanitize_file_name($_POST['filename'] ?? '');
        $filePath = PathHelper::getLogsDir() . '/' . $filename;
        $fileExists = file_exists($filePath);

        if ($fileExists) {
            file_put_contents($filePath, '');
        }

        wp_send_json_success(['filename' => $filename, 'cleared' => true]);
    }

    /**
     * AJAX handler: Clear all log files.
     * Action: wp_ajax_{plugin}_clear_all_logs
     */
    public function ajaxClearAllLogs(): void
    {
        check_ajax_referer(AjaxActionType::ClearAllLogs->nonceAction(), 'nonce');

        $hasPermission = current_user_can(CapabilityType::ManageOptions->value);
        if (!$hasPermission) {
            wp_send_json_error(['message' => 'Unauthorized'], 403);
        }

        $logsDir = PathHelper::getLogsDir();
        $logFiles = glob($logsDir . '/*.log');
        $clearedCount = 0;

        foreach ($logFiles as $file) {
            file_put_contents($file, '');
            $clearedCount++;
        }

        wp_send_json_success(['cleared' => $clearedCount]);
    }
}
```

### AJAX Action Registration

In the `Admin` class constructor:

```php
add_action(HookType::ajax(AjaxActionType::ReadLogFile->value), [$this, 'ajaxReadLogFile']);
add_action(HookType::ajax(AjaxActionType::ClearLogFile->value), [$this, 'ajaxClearLogFile']);
add_action(HookType::ajax(AjaxActionType::ClearAllLogs->value), [$this, 'ajaxClearAllLogs']);
```

### Required Enum Cases

```php
enum AjaxActionType: string
{
    case ReadLogFile = 'pluginname_read_log_file';
    case ClearLogFile = 'pluginname_clear_log_file';
    case ClearAllLogs = 'pluginname_clear_all_logs';

    public function nonceAction(): string
    {
        return $this->value . '_nonce';
    }
}
```

### REST vs AJAX Decision Matrix

| Operation | Pattern | Reason |
|-----------|---------|--------|
| Log file read/clear (admin-only) | **AJAX** (`wp_ajax_*`) | No external consumers, admin-page JS only |
| Error log query (structured data) | **REST** | May be consumed by Go backend or external tools |
| Error session management | **REST** | Structured data with filtering, pagination |
| Flash banner dismiss | **AJAX** | Simple toggle, admin-page JS only |
| Plugin settings update | **REST** | Standard CRUD, may have external consumers |

---

## 16.12 Admin Errors Page — Complete Template

Every plugin MUST include an error management admin page. This template orchestrates the error viewing UI from partials.

### Orchestrator Template: `templates/admin-errors.php`

```php
<?php
/**
 * Admin Errors Page — Error log viewer and session manager.
 *
 * Orchestrates error management UI via partials.
 *
 * @package PluginName
 * @since   1.0.0
 */

use PluginName\Enums\PluginConfigType;
use PluginName\Enums\AdminPageType;

if (!defined('ABSPATH')) {
    exit;
}

$pluginName = PluginConfigType::Name->value;
$pluginSlug = PluginConfigType::Slug->value;
?>
<div class="wrap pluginname-admin pluginname-errors">
    <?php
    // ── Page Header ──────────────────────────────────
    $pageIcon = 'dashicons-warning';
    $pageTitle = __('Error Management', $pluginSlug);
    $pageDescription = __('View, inspect, and clear plugin error logs.', $pluginSlug);
    include __DIR__ . '/partials/shared/page-header.php';

    // ── Flash Banner (unseen errors) ─────────────────
    include __DIR__ . '/partials/errors/flash-banner.php';
    ?>

    <!-- ── Actions Bar ──────────────────────────────── -->
    <div class="pluginname-card">
        <div class="pluginname-actions-row">
            <select id="log-file-selector">
                <option value="error.log"><?php esc_html_e('Error Log', $pluginSlug); ?></option>
                <option value="info.log"><?php esc_html_e('Info Log', $pluginSlug); ?></option>
                <option value="stacktrace.log"><?php esc_html_e('Stack Trace Log', $pluginSlug); ?></option>
                <option value="fatal-errors.log"><?php esc_html_e('Fatal Errors', $pluginSlug); ?></option>
            </select>
            <button id="btn-read-log" class="button button-primary">
                <span class="dashicons dashicons-visibility"></span>
                <?php esc_html_e('View Log', $pluginSlug); ?>
            </button>
            <button id="btn-clear-log" class="button button-secondary">
                <span class="dashicons dashicons-trash"></span>
                <?php esc_html_e('Clear File', $pluginSlug); ?>
            </button>
            <button id="btn-clear-all" class="button button-secondary">
                <span class="dashicons dashicons-dismiss"></span>
                <?php esc_html_e('Clear All Logs', $pluginSlug); ?>
            </button>
            <span id="log-action-status" class="pluginname-inline-status"></span>
        </div>
    </div>

    <!-- ── Log Content Viewer ───────────────────────── -->
    <?php include __DIR__ . '/partials/errors/log-file-viewer.php'; ?>

    <!-- ── Error Sessions Table ─────────────────────── -->
    <?php include __DIR__ . '/partials/errors/error-session-table.php'; ?>

    <!-- ── Error Detail Modal ───────────────────────── -->
    <?php include __DIR__ . '/partials/errors/error-detail-modal.php'; ?>
</div>
```

### Partial: `partials/errors/flash-banner.php`

```php
<?php
/**
 * Flash Banner — Shows unseen error count.
 * Receives: $pluginSlug (from parent template)
 */

if (!defined('ABSPATH')) { exit; }

$unseenCount = get_option(PluginConfigType::Slug->value . '_unseen_error_count', 0);
$hasUnseen = ($unseenCount > 0);

if (!$hasUnseen) { return; }
?>
<div id="pluginname-flash-banner" class="pluginname-flash-banner">
    <span class="flash-icon">⚠️</span>
    <div class="flash-content">
        <strong>
            <?php printf(
                esc_html__('%d new error(s) since your last visit.', $pluginSlug),
                $unseenCount
            ); ?>
        </strong>
    </div>
    <button id="btn-dismiss-flash" class="button flash-dismiss">
        <?php esc_html_e('Mark as Seen', $pluginSlug); ?>
    </button>
</div>
```

### Partial: `partials/errors/log-file-viewer.php`

```php
<?php
/**
 * Log File Viewer — Dark terminal-style content display.
 * Receives: $pluginSlug (from parent template)
 */

if (!defined('ABSPATH')) { exit; }
?>
<div class="pluginname-card" id="log-viewer-card">
    <h2>
        <span class="dashicons dashicons-media-text"></span>
        <?php esc_html_e('Log File Content', $pluginSlug); ?>
        <span id="log-file-name" class="pluginname-version-badge">error.log</span>
        <span id="log-file-size" class="pluginname-version-badge"></span>
    </h2>
    <div id="log-content-wrapper">
        <div id="log-empty" class="file-empty" style="display: none;">
            <span class="dashicons dashicons-yes-alt"></span>
            <p><?php esc_html_e('No errors found — looking good!', $pluginSlug); ?></p>
        </div>
        <pre id="log-content" class="code-pre" style="display: none;"></pre>
    </div>
</div>
```

### Partial: `partials/errors/error-session-table.php`

```php
<?php
/**
 * Error Session Table — Groups errors by request session.
 * Receives: $pluginSlug, $sessions (from parent template or AJAX)
 */

if (!defined('ABSPATH')) { exit; }
?>
<div class="pluginname-card">
    <h2>
        <span class="dashicons dashicons-database"></span>
        <?php esc_html_e('Error Sessions', $pluginSlug); ?>
    </h2>
    <table class="wp-list-table widefat fixed striped">
        <thead>
            <tr>
                <th class="column-time"><?php esc_html_e('Time', $pluginSlug); ?></th>
                <th class="column-count"><?php esc_html_e('Errors', $pluginSlug); ?></th>
                <th class="column-level"><?php esc_html_e('Severity', $pluginSlug); ?></th>
                <th class="column-message"><?php esc_html_e('First Error', $pluginSlug); ?></th>
                <th class="column-status"><?php esc_html_e('Status', $pluginSlug); ?></th>
                <th class="column-actions"><?php esc_html_e('Actions', $pluginSlug); ?></th>
            </tr>
        </thead>
        <tbody id="error-sessions-body">
            <tr>
                <td colspan="6" class="no-items">
                    <?php esc_html_e('Loading sessions...', $pluginSlug); ?>
                </td>
            </tr>
        </tbody>
    </table>
</div>
```

### Partial: `partials/errors/error-detail-modal.php`

```php
<?php
/**
 * Error Detail Modal — Shows full error context with stack trace.
 * Receives: $pluginSlug (from parent template)
 */

if (!defined('ABSPATH')) { exit; }

$modalId = 'error-detail-modal';
$modalTitle = __('Error Details', $pluginSlug);
$modalIcon = 'dashicons-warning';
$modalIconColor = '#dc2626';
$modalMaxWidth = '1000px';
$modalBody = '
    <div class="error-detail-content">
        <div class="error-meta" id="error-meta"></div>
        <h4>' . esc_html__('Stack Trace', $pluginSlug) . '</h4>
        <pre class="stack-trace" id="error-stack-trace"></pre>
    </div>';
$modalFooter = '<button class="button" onclick="document.getElementById(\'' . $modalId . '\').style.display=\'none\'">'
    . esc_html__('Close', $pluginSlug) . '</button>';

include __DIR__ . '/../shared/modal-wrapper.php';
?>
```

### Required partials directory structure

```
templates/
├── admin-errors.php                    ← Orchestrator (under 80 lines)
└── partials/
    ├── shared/
    │   ├── page-header.php
    │   └── modal-wrapper.php
    └── errors/
        ├── flash-banner.php            ← Unseen error alert
        ├── log-file-viewer.php         ← Dark terminal log display
        ├── error-session-table.php     ← Session-grouped error table
        └── error-detail-modal.php      ← Full error context modal
```

---

## 16.13 ErrorSessions Table — SQLite Migration

Error sessions require a dedicated SQLite table. This migration follows the pattern from Phase 8, §8.5.

### Migration trait: `DatabaseMigrationsErrorSessionsTrait.php`

```php
namespace PluginName\Database\Traits;

if (!defined('ABSPATH')) {
    exit;
}

use PDOException;
use PluginName\Enums\TableType;

trait DatabaseMigrationsErrorSessionsTrait
{
    /**
     * Create ErrorSessions table for grouped error tracking.
     * Version: Assign the next sequential migration number in your plugin.
     */
    private function migrateErrorSessions(int $current, int $version): void
    {
        if ($current >= $version) {
            return;
        }

        $this->fileLogger->info("Applying migration v{$version}: ErrorSessions table");
        $table = TableType::ErrorSessions->value;

        $sql = <<<SQL
            CREATE TABLE IF NOT EXISTS {$table} (
                Id              INTEGER PRIMARY KEY AUTOINCREMENT,
                SessionId       TEXT    NOT NULL UNIQUE,
                ErrorCount      INTEGER NOT NULL DEFAULT 0,
                FirstErrorType  TEXT    DEFAULT '',
                FirstMessage    TEXT    DEFAULT '',
                Severity        TEXT    NOT NULL DEFAULT 'error',
                IsSeen          INTEGER NOT NULL DEFAULT 0,
                PluginVersion   TEXT    DEFAULT '',
                CreatedAt       TEXT    NOT NULL DEFAULT (datetime('now')),
                UpdatedAt       TEXT    NOT NULL DEFAULT (datetime('now'))
            )
        SQL;

        $this->pdo->exec($sql);

        $this->pdo->exec("CREATE INDEX IF NOT EXISTS idx_es_session_id ON {$table}(SessionId)");
        $this->pdo->exec("CREATE INDEX IF NOT EXISTS idx_es_is_seen ON {$table}(IsSeen)");
        $this->pdo->exec("CREATE INDEX IF NOT EXISTS idx_es_created ON {$table}(CreatedAt)");
        $this->pdo->exec("CREATE INDEX IF NOT EXISTS idx_es_severity ON {$table}(Severity)");

        $this->recordMigration($version);
    }
}
```

### Required TableType enum case

```php
enum TableType: string
{
    // ... existing cases
    case ErrorSessions = 'ErrorSessions';
}
```

### Column reference

| Column | Type | Description |
|--------|------|-------------|
| `Id` | INTEGER PK | Auto-increment primary key |
| `SessionId` | TEXT UNIQUE | UUID or timestamp-based session identifier |
| `ErrorCount` | INTEGER | Number of errors captured in this session |
| `FirstErrorType` | TEXT | PHP error type label (e.g., `E_WARNING`) |
| `FirstMessage` | TEXT | First error message in the session |
| `Severity` | TEXT | Highest severity: `fatal`, `error`, `warning` |
| `IsSeen` | INTEGER | `0` = unseen (shows badge), `1` = dismissed |
| `PluginVersion` | TEXT | Plugin version that generated the errors |
| `CreatedAt` | TEXT | Session creation timestamp |
| `UpdatedAt` | TEXT | Last error added timestamp |

### Integration points

| Component | How it uses ErrorSessions |
|-----------|--------------------------|
| Shutdown handler (Phase 4, §4.13) | Creates new session or increments `ErrorCount` |
| Admin menu badge (Phase 8, §8.1) | Counts rows where `IsSeen = 0` |
| Flash banner (§16.6) | Queries unseen count |
| Dismiss AJAX (§16.5) | Sets `IsSeen = 1` for all rows |
| Clear all (§16.11) | Deletes all rows from table |
| Error detail modal (§16.12) | Queries single session by `SessionId` |

---

## 16.14 Checklist

- [ ] `ErrorType` class with `FATAL_TYPES`, `WARNING_TYPES`, `RECOVERABLE_TYPES`, `TYPE_LABELS`
- [ ] `InitHelpers::errorLogWithPrefix()` and `errorLog()` for Tier 1 logging
- [ ] `FileLogger` for Tier 2 structured logging (Phase 4)
- [ ] `ErrorResponse` class with `logAndReturn()`, `logAndReturnFalse()`, `logAndReturnEnvelope()`, `logAndReturnWpError()`
- [ ] `error-logs` and `error-sessions` REST endpoints
- [ ] Error log retrieval settings with 3-level resolution (request → stored → defaults)
- [ ] Error session model with `is_seen` tracking
- [ ] `ErrorSessions` SQLite table migration with indexes (§16.13)
- [ ] `AdminErrorAjaxTrait` with read/clear/clear-all handlers
- [ ] `admin-errors.php` orchestrator template with 4 partials
- [ ] Flash banner with AJAX dismiss
- [ ] Auto-refresh with stop-on-modal behavior
- [ ] `safeExecute` wrapper on all REST handlers
- [ ] Error notification settings in `OptionNameType`
- [ ] Admin menu error count badge (Phase 8, §8.1)
- [ ] No forbidden error patterns in codebase (see Phase 4, §4.8, Rule 6)

---

*Last Updated: 2026-04-14*
