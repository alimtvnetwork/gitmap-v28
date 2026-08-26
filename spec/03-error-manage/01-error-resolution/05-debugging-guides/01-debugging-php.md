# PHP Debugging Guide

> **Version:** 1.0.0  
> **Created:** 2026-02-04  
> **Applies To:** WordPress Plugins, PHP Backends

---

## Overview

This guide covers debugging patterns for PHP applications, particularly WordPress plugins. It includes initialization order, logging patterns, and common troubleshooting scenarios.

---

## Initialization Order (CRITICAL)

The plugin follows a strict initialization order to prevent errors:

1. **STEP 1: Directories** - All directories are created and verified FIRST
2. **STEP 2: Database** - Database is initialized ONLY AFTER directories exist
3. **STEP 3+: Components** - All other components initialize after database is ready

This order is enforced by **reusable helper functions**:
- `ensure_directories_exist()` - Creates/verifies all directories (only runs once)
- `ensure_database_ready()` - Initializes database (only runs once, requires directories)

---

## Log Files Location

All log files are stored in:
```
wp-content/uploads/{plugin-slug}/logs/
```

### Available Log Files

| File | Purpose |
|------|---------|
| `debug.log` | Detailed execution trace |
| `error.log` | Errors and exceptions with stack traces |

---

## Enabling/Disabling Logging

Logging is controlled by constants in `includes/constants.php`:

```php
// Enable debug logging (trace every step)
define('PLUGIN_DEBUG_LOGGING', true);  // Set to false to disable

// Enable error logging (exceptions and errors)
define('PLUGIN_ERROR_LOGGING', true);  // Set to false to disable
```

---

## What Gets Logged

### Debug Log (debug.log)

- Plugin initialization start
- PHP and WordPress versions
- Class loading (each file loaded)
- Constructor calls
- Hook registrations
- Database connection steps
- Directory creation
- Component initialization
- Success/failure for each step

### Error Log (error.log)

- Critical errors with full stack traces
- Database connection failures
- Missing classes or dependencies
- File loading errors
- Exception details (message, code, file, line, trace)

---

## How to Debug Plugin Activation

1. **Clear the logs** (optional):
   ```
   Delete: wp-content/uploads/{plugin-slug}/logs/debug.log
   Delete: wp-content/uploads/{plugin-slug}/logs/error.log
   ```

2. **Try to activate the plugin** in WordPress admin

3. **Check the debug log**:
   ```
   Open: wp-content/uploads/{plugin-slug}/logs/debug.log
   ```
   - Look for the last successful step
   - The log stops where the error occurred

4. **Check the error log**:
   ```
   Open: wp-content/uploads/{plugin-slug}/logs/error.log
   ```
   - Find exception details
   - Review stack traces
   - Identify the exact file and line number causing issues

---

## Common Issues and Solutions

### Issue: Plugin won't activate

**Check:**
1. PDO SQLite extension installed?
   - Look for: "pdo_sqlite extension is loaded" in debug.log
   - If not found, install PHP SQLite extension

2. Directory permissions
   - Look for: "Database directory is writable" in debug.log
   - If not, check permissions on wp-content/uploads/

3. Missing classes
   - Look for: "✓ Loaded successfully" for each file
   - If you see "✗ File not found", check file exists

### Issue: Database connection fails

**Check debug.log for:**
- "Database directory: [path]" - Is path correct?
- "Database directory created successfully" - Did creation work?
- "Database directory is writable" - Check permissions
- "Main database PDO connection established" - Did connection succeed?

**Check error.log for:**
- PDO exceptions with detailed error messages
- File permission errors
- SQLite errors

### Issue: Components not initializing

**Check debug.log for:**
- "Initializing Component..." - Each component logs initialization
- Look for "✓ Component initialized" vs "✗ Component class not found"
- Missing dependencies will be logged

---

## Log Format

### Debug Log Entry Format:

```
[2026-02-04 15:30:45] [DEBUG] [Memory: 12.5 MB] Message here
--------------------------------------------------------------------------------
```

### Error Log Entry Format:

```
[2026-02-04 15:30:45] [ERROR] [Memory: 12.5 MB] Error message
Context: Array
(
    [exception] => Array
        (
            [message] => Error details
            [code] => 0
            [file] => /path/to/file.php
            [line] => 123
            [trace] => Full stack trace here
        )
)
--------------------------------------------------------------------------------
```

---

## Reading the Logs

### Step-by-Step Execution Trace

The debug log shows exactly how far the plugin got before failing:

```
=== PLUGIN INITIALIZATION STARTED ===
Plugin Version: 1.0.6
WordPress Version: 6.4.2
PHP Version: 8.1.0
Loading Paths class...
Paths class loaded successfully
Loading Config class...
Config class loaded successfully
Loading plugin dependencies...
Loading file: includes/class-database.php
  ✓ Loaded successfully: includes/class-database.php
...
[Last successful line before error]
```

### Finding the Exact Failure Point

1. Scroll to the bottom of debug.log
2. The last logged message shows where execution stopped
3. Check error.log for the corresponding exception
4. The exception shows:
   - Exact file path
   - Line number
   - Full stack trace

---

## Logger Class Pattern

```php
<?php
namespace PluginNamespace\Utils;

class Logger
{
    /**
     * Log debug message
     */
    public static function debug(string $message, array $context = []): void
    {
        if (!defined('PLUGIN_DEBUG_LOGGING') || !PLUGIN_DEBUG_LOGGING) {
            return;
        }
        
        self::write('DEBUG', $message, $context);
    }
    
    /**
     * Log error with stack trace
     *
     * CRITICAL: First parameter is always Throwable.
     * The message is extracted from $e->getMessage() automatically.
     * The stack trace is always appended via $e->getTraceAsString().
     * See: spec/02-coding-guidelines/04-php/forbidden-patterns.md rule 1.8
     */
    public static function error(Throwable $exception, array $context = []): void
    {
        if (!defined('PLUGIN_ERROR_LOGGING') || !PLUGIN_ERROR_LOGGING) {
            return;
        }
        
        $context['exception'] = [
            'message' => $exception->getMessage(),
            'code' => $exception->getCode(),
            'file' => $exception->getFile(),
            'line' => $exception->getLine(),
            'trace' => $exception->getTraceAsString(),
        ];
        
        self::write('ERROR', $exception->getMessage(), $context, 'error.log');
    }
    
    /**
     * Write to log file
     */
    private static function write(
        string $level, 
        string $message, 
        array $context = [], 
        string $file = 'debug.log'
    ): void {
        $logDir = WP_CONTENT_DIR . '/uploads/' . PLUGIN_SLUG . '/logs/';
        
        if (!is_dir($logDir)) {
            wp_mkdir_p($logDir);
        }
        
        $memory = round(memory_get_usage() / 1024 / 1024, 2);
        $timestamp = date('Y-m-d H:i:s');
        
        $entry = sprintf(
            "[%s] [%s] [Memory: %s MB] %s\n",
            $timestamp,
            $level,
            $memory,
            $message
        );
        
        if (!empty($context)) {
            $entry .= "Context: " . print_r($context, true) . "\n";
        }
        
        $entry .= str_repeat('-', 80) . "\n";
        
        file_put_contents($logDir . $file, $entry, FILE_APPEND | LOCK_EX);
    }
}
```

---

## Disabling Logs in Production

Once debugging is complete, disable logs by editing constants:

```php
define('PLUGIN_DEBUG_LOGGING', false);
define('PLUGIN_ERROR_LOGGING', false);
```

Or set environment variables:
```php
putenv('PLUGIN_DEBUG_LOGGING=false');
putenv('PLUGIN_ERROR_LOGGING=false');
```

---

## Cross-Reference

- [Error Resolution Overview](../../00-overview.md)
- PHP Coding Guidelines *(external spec)*
- WordPress Guidelines *(external spec)*
