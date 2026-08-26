# PHP Naming Conventions

> **Version:** 1.3.0
> **Updated:** 2026-03-09
> **Baseline:** PSR-12 / PSR-1  
> **Applies to:** All PHP code unless overridden by project-specific or framework-specific conventions

---

> **Override Rule:** These are the default conventions for PHP. Any language-specific, framework-specific, or project-specific convention file **takes precedence** when it conflicts with this baseline. For example, a WordPress plugin project may mandate `snake_case` methods — that override wins.

---

## Classes, Interfaces, Traits, Enums

Use **PascalCase**. Enums MUST use the **`Type` suffix**.

```php
class UploadManager {}
interface CacheDriver {}
trait HasTimestamps {}
enum UploadSourceType: string {}  // ✅ Type suffix required
enum CapabilityType: string {}    // ✅ Type suffix required
```

**Rules:**

- One class per file
- File name matches class name
- Names should be nouns
- Avoid abbreviations unless universal (`HTTP`, `API`, `CLI`)

---

## Methods and Functions

Use **camelCase**.

```php
function processUpload() {}
function getRetryCount() {}
```

**Rules:**

- Verb or verb phrase
- Describe behavior, not implementation
- Avoid prefixes like `do`, `handle`, `run`

```php
// ✅ Good
calculateChecksum()

// ❌ Bad
doChecksumThing()
```

---

## Variables

Use **camelCase**.

```php
$maxRetries = 3;
$uploadSource = UploadSource::Script;
```

**Rules:**

- Clear intent over short names
- Avoid Hungarian notation
- Avoid `snake_case` for variables in modern PHP
- Boolean variables must use `$is` or `$has` prefix in **camelCase** (e.g., `$isActive`, `$hasErrors`) — not `$is_active` or `$has_errors`. The prefix follows the same camelCase rule as all other variables

---

## Constants

Use **UPPER_SNAKE_CASE**.

```php
const MAX_RETRIES = 3;
const DEFAULT_TIMEOUT = 30;
```

For class constants:

```php
class Limits
{
    public const MAX_RETRIES = 3;
}
```

Enums are the exception: enum cases follow **PascalCase**, not uppercase. Enum names use `Type` suffix.

```php
UploadSourceType::RestApi
CapabilityType::ManageOptions
```

---

## Enum Cases

Use **PascalCase**. Enum names MUST have `Type` suffix.

```php
enum UploadSourceType: string
{
    case Script;
    case RestApi;
    case AdminUi;
    case WpCli;
}
```

**Rules:**

- Match domain naming
- Avoid screaming uppercase
- Treat them like class names

---

## Namespaces

Use **PascalCase**, structured by domain.

```php
namespace App\Domain\Upload;
namespace App\Infrastructure\Http;
```

**Rules:**

- Reflect architecture, not folders alone
- Avoid generic buckets like `Utils` or `Helpers`

---

## Files

File names follow **PSR-4 autoloading** — the file name must match the class, enum, or interface inside it exactly.

### Class / Enum / Interface / Trait files

Use **PascalCase** and match the symbol name 1:1.

```
UploadManager.php
HttpClient.php
CacheDriver.php
SnapshotFactory.php
```

**Rules:**

- One class or enum per file
- File name equals class name exactly
- Case-sensitive on Linux servers
- No underscores, no `snake_case`, no `class-kebab-case` prefixes

```php
// ❌ Bad
class-upload-manager.php
upload_source.php
Upload_Source.php
uploadsource.php

// ✅ Good
UploadManager.php
UploadSource.php
```

### Domain-based directory structure

Files are organized into **domain folders** within `includes/`:

```
includes/
  Admin/           — Admin UI and settings
  Agent/           — Agent management
  Database/        — Database, ORM, caching
  Enums/           — Backed enums (PSR-4 namespace)
  Helpers/         — Utility classes (path, envelope, error checking)
  Logging/         — Logger implementations
  Post/            — Post/content management
  Snapshot/        — Snapshot system (backup, restore, providers)
  Update/          — Auto-update resolver
  Upload/          — Upload ignore rules
  constants.php    — Global constants (unprefixed)
```

### Namespaced directory structure

Directory structure mirrors the namespace:

```php
namespace RiseupAsia\Enums;
```

```
includes/
  Enums/
    CapabilityType.php
    ErrorType.php
    HookType.php
    HttpMethodType.php
    UploadSourceType.php
    PathSubdirType.php
    PathDatabaseType.php
    PathLogFileType.php
    PathConfigType.php
```

Autoloaders depend on this mapping.

### Files without classes (rare)

For procedural or config files, use **lowercase with underscores**:

```
constants.php
constants-compat.php
```

These are exceptions and should be minimal in modern architecture.

---

## Array Key Conventions

### Log Context Keys — camelCase

Internal log context array keys (passed to logger calls) MUST use **camelCase**:

```php
// ✅ REQUIRED
$this->fileLogger->info('Post created', array('postId' => $postId));
$this->fileLogger->warn('Duplicate detected', array('duplicateDir' => $dir, 'targetSlug' => $slug));
$this->fileLogger->debug('Agent API request', array('agentId' => $id, 'method' => $method));
```

### Database Column Keys — PascalCase

Array keys referencing database columns (inserts, updates, WHERE conditions) MUST use **PascalCase** to match the schema:

```php
// ✅ REQUIRED
$this->db->insert(TableType::Transactions->value, array('PluginSlug' => $slug, 'CreatedAt' => $now));
```

### API Response Keys — PascalCase

Array keys in REST API responses (success payloads, error enrichment, health checks) MUST use **PascalCase**:

```php
// ❌ WRONG — snake_case in API response
return new WP_REST_Response([
    'plugin_version' => PluginConfigType::Version->value,  // Wrong
    'timestamp' => gmdate('c'),                            // Wrong
    'log_hint' => $this->getLogHint($status),              // Wrong
]);

// ✅ CORRECT — PascalCase in API response
return new WP_REST_Response([
    'PluginVersion' => PluginConfigType::Version->value,
    'Timestamp' => gmdate(DateFormatType::Iso8601->value),
    'LogHint' => $this->getLogHint($status),
]);
```

### Persistence-Level Keys — Exempt

WordPress `wp_options` keys, SQLite `_snapshot_meta` keys, and external API response keys retain their native casing (typically snake_case). These are **not** subject to the PascalCase rule.

---

## Summary Table

| Element                    | Convention                     | Example                          |
|----------------------------|--------------------------------|----------------------------------|
| Class / Interface          | PascalCase                     | `SnapshotFactory`                |
| Enum                       | PascalCase + `Type` suffix     | `UploadSourceType`               |
| Trait                      | PascalCase                     | `HasTimestamps`                  |
| Method / Function          | camelCase                      | `processUpload()`                |
| Variable                   | camelCase                      | `$maxRetries`                    |
| Boolean variable           | `$is` / `$has` + camelCase     | `$isActive`, `$hasErrors`        |
| Constant                   | UPPER_SNAKE_CASE               | `MAX_RETRIES`                    |
| Enum case                  | PascalCase                     | `RestApi`                        |
| Namespace                  | PascalCase                     | `RiseupAsia\Enums`               |
| File (class/trait)         | PascalCase.php                 | `SnapshotFactory.php`            |
| File (enum)                | PascalCase + Type.php          | `UploadSourceType.php`           |
| File (config/procedural)   | lowercase_with_underscores.php | `constants.php`                  |
| Directory (domain folder)  | PascalCase                     | `Snapshot/`, `Database/`         |
| Log context array key      | camelCase                      | `'postId'`, `'masterDir'`        |
| DB column array key        | PascalCase                     | `'PluginSlug'`, `'CreatedAt'`    |
| API response array key     | PascalCase                     | `'PluginVersion'`, `'Timestamp'` |

---

## Common Mistakes — PHP Naming

These are real violations found and fixed across the codebase. Use as a checklist to avoid repeating them.

### Mistake 1: snake_case Log Context Keys

```php
// ❌ WRONG — snake_case in log context
$this->fileLogger->info('Post created', array(
    'post_id'    => $postId,     // Wrong
    'master_dir' => $dir,        // Wrong
    'agent_id'   => $agentId,    // Wrong
));

// ✅ CORRECT — camelCase
$this->fileLogger->info('Post created', array(
    'postId'    => $postId,
    'masterDir' => $dir,
    'agentId'   => $agentId,
));
```

### Mistake 2: camelCase or snake_case for DB Column / API Response Keys

```php
// ❌ WRONG — camelCase/snake_case for DB columns
$this->db->insert(TableType::Snapshots->value, array(
    'totalRows'     => $rows,      // Wrong — camelCase
    'trigger_source' => $trigger,  // Wrong — snake_case
    'status'        => $status,    // Wrong — lowercase
));

// ✅ CORRECT — PascalCase matching schema
$this->db->insert(TableType::Snapshots->value, array(
    'TotalRows'     => $rows,
    'TriggerSource' => $trigger,
    'Status'        => $status,
));

// ❌ WRONG — snake_case in API response
$data['plugin_version'] = PluginConfigType::Version->value;
$data['log_hint'] = $this->getLogHint($status);

// ✅ CORRECT — PascalCase in API response
$data['PluginVersion'] = PluginConfigType::Version->value;
$data['LogHint'] = $this->getLogHint($status);
```

### Mistake 3: snake_case Method Names

```php
// ❌ WRONG — WordPress-style snake_case
public function get_plugin_file() { ... }
private function handle_upload_error() { ... }

// ✅ CORRECT — camelCase
public function getPluginFile() { ... }
private function handleUploadError() { ... }
```

### Mistake 4: Missing `Type` Suffix on Enums

```php
// ❌ WRONG — no Type suffix
enum UploadSource: string { ... }
enum Capability: string { ... }

// ✅ CORRECT — Type suffix required
enum UploadSourceType: string { ... }
enum CapabilityType: string { ... }
```

### Mistake 5: Leading Backslash on Global Types

```php
// ❌ WRONG — backslash-qualified
catch (\Throwable $e) { ... }
$db = new \PDO($dsn);

// ✅ CORRECT — import via use
use Throwable;
use PDO;
catch (Throwable $e) { ... }
$db = new PDO($dsn);
```

### Mistake 6: Raw `===` for Enum Comparison

```php
// ❌ WRONG — raw operator
if ($status === StatusType::Success) { ... }

// ✅ CORRECT — isEqual() method
if ($status->isEqual(StatusType::Success)) { ... }
```

### Mistake 7: Uppercase Abbreviations

```php
// ❌ WRONG
$postId = 5;
$fileUrl = '/path';
$hashMD5 = md5($data);

// ✅ CORRECT
$postId = 5;
$fileUrl = '/path';
$hashMd5 = md5($data);
```

---

> **Reminder:** This is the PSR-12 baseline. Project-level specs (e.g., WordPress plugin conventions in this same folder's [readme.md](./07-php-standards-reference/00-overview.md)) may override specific rules — those overrides take precedence.
