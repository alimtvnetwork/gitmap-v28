# PHP Error Management — Known Issues & Cases

> **Version:** 1.0.0
> **Updated:** 2026-03-12
> **Status:** Active

---

## QUpload Syntax Validator False Positives

The QUpload plugin runs `token_get_all($content, TOKEN_PARSE)` on all non-template PHP files before activating an uploaded plugin. This can produce false-positive syntax errors.

### Template Detection Logic

Files are skipped (not validated) when they:
- Start with `<?php`
- **AND** contain `?>`

Pure PHP files (no closing `?>` tag) are **always validated**.

### Known Blocked Patterns

| Pattern | Error Message | Severity | Fix |
|---------|--------------|----------|-----|
| `is_array($var)` | `syntax error, unexpected token "array"` | Critical — blocks activation | Replace with `gettype($var) === PhpNativeType::PhpArray->value` |
| `array()` long syntax | `syntax error, unexpected token "array"` | Critical — blocks activation | Replace with `[]` short array syntax |
| `= array()` default param | `syntax error, unexpected token "array"` | Critical — blocks activation | Replace with `= []` |

### Diagnostic Indicators

When QUpload blocks activation due to syntax:
- **error.txt**: `Plugin syntax validation failed for {slug} file: {path}: syntax error, unexpected token "array" (:{line})`
- **stacktrace.txt**: Full trace from `UploadExtractTrait::checkPhpFileSyntax()`
- **API response**: `Plugin uploaded but activation was blocked due to PHP syntax error in {path}`

The plugin IS uploaded and extracted — it just isn't activated. Manual activation via WP admin will also fail if the server PHP version triggers the same error.

---

## Log File Management

### Log Files (QUpload)

| File | Path | Content |
|------|------|---------|
| `log.txt` | `wp-content/uploads/qupload/logs/log.txt` | General operational log |
| `error.txt` | `wp-content/uploads/qupload/logs/error.txt` | Errors only |
| `stacktrace.txt` | `wp-content/uploads/qupload/logs/stacktrace.txt` | Full stack traces |

### Lifecycle

- **Activation**: All log files are cleared (`FileLogger::clearAllLogFiles()`) for a fresh start
- **Runtime**: Logs are appended as errors/events occur
- **Deactivation**: Temp directory (`qupload/temp/`) is cleared; logs are preserved
- **Uninstall**: Entire `wp-content/uploads/qupload/` directory is recursively deleted

---

## PhpNativeType Enum Reference

Located at `includes/Enums/PhpNativeType.php` in the Riseup Asia Uploader plugin.

| Case | Value | Replaces |
|------|-------|----------|
| `PhpArray` | `'array'` | `is_array()` checks |
| `PhpString` | `'string'` | `is_string()` checks |
| `PhpInteger` | `'integer'` | `is_int()` checks |
| `PhpDouble` | `'double'` | `is_float()` checks |
| `PhpBoolean` | `'boolean'` | `is_bool()` checks |
| `PhpObject` | `'object'` | `is_object()` checks |
| `PhpNull` | `'NULL'` | `is_null()` checks |

### Usage Pattern

```php
// WRONG — blocked by QUpload validator
$isBackupError = is_array($backupId);

// WRONG — magic string
$isBackupError = gettype($backupId) === 'array';

// CORRECT — typed constant
$isBackupError = gettype($backupId) === PhpNativeType::PhpArray->value;
```
