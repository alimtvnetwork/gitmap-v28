# Issues & Fixes Log — Historical Reference

> **Version:** 1.0.0
> **Updated:** 2026-03-09
> **Purpose:** Comprehensive log of all coding standard violations found and fixed, with root cause analysis and prevention rules

---

## Issue Categories

| Category | Issues |
|----------|--------|
| [Naming Convention Violations](#naming-convention-violations) | #01–#04 |
| [Database Casing Violations](#database-casing-violations) | #05–#08 |
| [Boolean & Negation Violations](#boolean--negation-violations) | #09–#11 |
| [Enum Standard Violations](#enum-standard-violations) | #12–#14 |
| [Formatting Violations](#formatting-violations) | #15–#17 |
| [Type Safety Violations](#type-safety-violations) | #18–#19 |

---

## Naming Convention Violations

### Issue #01 — snake_case Log Context Keys (Batch G)

**Scope:** 8 PHP files, 17 keys
**Root Cause:** Log context arrays used `snake_case` keys (`'post_id'`, `'master_dir'`) instead of the required camelCase.
**Impact:** Inconsistent log output, harder to grep/parse logs across the system.

**Before (❌):**
```php
$this->fileLogger->info('Post created', array('post_id' => $postId));
$this->fileLogger->debug('Backup started', array('master_dir' => $dir));
$this->fileLogger->warn('Upload conflict', array('duplicate_dir' => $path));
```

**After (✅):**
```php
$this->fileLogger->info('Post created', array('postId' => $postId));
$this->fileLogger->debug('Backup started', array('masterDir' => $dir));
$this->fileLogger->warn('Upload conflict', array('duplicateDir' => $path));
```

**Files Fixed:**
- `PostCrudTrait.php`: `post_id` → `postId`
- `CategoryTrait.php`: `term_id` → `termId`
- `UploadZipTrait.php`: `duplicate_dir` → `duplicateDir`, `target_slug` → `targetSlug`, `old_version` → `oldVersion`, `new_version` → `newVersion`, `file_size` → `fileSize`
- `SnapshotExportHandlerTrait.php`: `snapshot_id` → `snapshotId`
- `AgentRemoteCoreTrait.php`: `agent_id` → `agentId`
- `RootDbSchemaTrait.php`: `mysql_version` → `mysqlVersion`, `wp_version` → `wpVersion`
- `DetectorSettingsTrait.php`: `changed_keys` → `changedKeys`
- `IncrementalBackup.php`: `master_dir` → `masterDir`
- `UploadInstallExtractTrait.php`: `target_dir` → `targetDir`

**Prevention:** All log context array keys must use camelCase. See [naming-conventions.md § Array Key Conventions](../04-php/03-naming-conventions.md#array-key-conventions).

---

### Issue #02 — Legacy Class Naming (`class Riseup_*`)

**Scope:** All PHP class files
**Root Cause:** Legacy WordPress naming convention used `Riseup_` prefix with underscores (e.g., `class Riseup_Upload_Manager`).
**Impact:** Violates PSR-4 autoloading, creates inconsistent naming with modern codebase.

**Before (❌):**
```php
class Riseup_Upload_Manager { ... }
class Riseup_Snapshot_Factory { ... }
```

**After (✅):**
```php
namespace RiseupAsia\Upload;
class UploadManager { ... }

namespace RiseupAsia\Snapshot;
class SnapshotFactory { ... }
```

**Prevention:** PSR-4 autoloading enforces PascalCase class names = file names. No `Riseup_` prefix.

---

### Issue #03 — Leading Backslash on Global Types

**Scope:** 16 PHP files
**Root Cause:** PHP global types like `Throwable`, `PDO`, `Exception` were used with leading backslash (`\Throwable`) instead of being imported via `use`.

**Before (❌):**
```php
catch (\Throwable $e) { ... }
$db = new \PDO($dsn);
```

**After (✅):**
```php
use Throwable;
use PDO;
// ...
catch (Throwable $e) { ... }
$db = new PDO($dsn);
```

**Prevention:** All global types must be imported via `use` statement. `Autoloader.php` is the sole exemption.

---

### Issue #04 — snake_case Method Names in PHP

**Scope:** All PHP trait and class methods
**Root Cause:** Some methods used `snake_case` (WordPress convention) instead of `camelCase`.

**Before (❌):**
```php
public function get_plugin_file() { ... }
private function handle_upload_error() { ... }
```

**After (✅):**
```php
public function getPluginFile() { ... }
private function handleUploadError() { ... }
```

**Prevention:** Zero-underscore policy. All methods use camelCase. WordPress hook callbacks are the only exemption.

---

## Database Casing Violations

### Issue #05 — snake_case Table Names (V13 Migration)

**Scope:** 12 SQLite tables
**Root Cause:** Original schema used snake_case table names (`agent_sites`, `snapshot_progress`).
**Impact:** Inconsistent with PascalCase standard, cross-language enum mismatch.

**Before (❌):**
```sql
CREATE TABLE transactions (...);
CREATE TABLE agent_sites (...);
CREATE TABLE snapshot_progress (...);
```

**After (✅):**
```sql
ALTER TABLE transactions RENAME TO Transactions;
ALTER TABLE agent_sites RENAME TO AgentSites;
ALTER TABLE snapshot_progress RENAME TO SnapshotProgress;
```

**Prevention:** `TableType` enum values must match PascalCase schema. Migration v13 handles legacy renames.

---

### Issue #06 — snake_case Column Names (V13 Migration)

**Scope:** ~50 columns across 12 tables
**Root Cause:** Columns used `snake_case` (`plugin_slug`, `created_at`, `agent_site_id`).

**Before (❌):**
```sql
SELECT plugin_slug, created_at FROM transactions WHERE agent_site_id = ?
```

**After (✅):**
```sql
SELECT PluginSlug, CreatedAt FROM Transactions WHERE AgentSiteId = ?
```

**Prevention:** All SQL queries must use PascalCase column names. `LogColumnType` enum provides type-safe column access.

---

### Issue #07 — PHP Consumer Code Using Old Column Names

**Scope:** 9 PHP files
**Root Cause:** After V13 migration renamed columns to PascalCase, consumer code still used old snake_case keys in arrays and SQL.

**Before (❌):**
```php
$record = array(
    'plugin_slug' => $slug,
    'created_at'  => $now,
    'status'      => 'success',
);
```

**After (✅):**
```php
$record = array(
    'PluginSlug' => $slug,
    'CreatedAt'  => $now,
    'Status'     => StatusType::Success->value,
);
```

**Prevention:** DB insert/update arrays must use PascalCase keys matching the schema. This caused null-value regressions when code used old keys.

---

### Issue #08 — ImportExecutionTrait buildSnapshotRecord Inconsistency

**Scope:** 1 file, 11 keys
**Root Cause:** `buildSnapshotRecord()` used a mix of camelCase and lowercase keys for DB column names.
**Impact:** Silent data loss — PDO/SQLite ignores unrecognized column names, resulting in NULL values.

**Before (❌):**
```php
return array(
    'sequence'      => $this->manager->getNextSequence(),
    'filename'      => basename($destDir),
    'totalRows'     => $metadata['total_rows'] ?? 0,
    'triggerSource'  => SnapshotTriggerType::Api->value,
    'importSource'   => json_encode($meta),
);
```

**After (✅):**
```php

return array(
    'Sequence'      => $this->manager->getNextSequence(),
    'Filename'      => basename($destDir),
    'TotalRows'     => $metadata['total_rows'] ?? 0,
    'TriggerSource' => SnapshotTriggerType::Api->value,
    'ImportSource'  => json_encode($meta),
);
```

**Prevention:** Every array key used in `$this->db->insert()` or `$this->db->update()` must be PascalCase.

---

## Boolean & Negation Violations

### Issue #09 — `!file_exists()` / `!is_dir()` Raw Negation

**Scope:** ~30 PHP call sites
**Root Cause:** Code used raw `!` operator on PHP filesystem functions.

**Before (❌):**
```php
if (!file_exists($path)) { return false; }
if (!is_dir($dir)) { mkdir($dir, 0755, true); }
if (!self::makeDirectory($dir)) { throw new Exception('...'); }
```

**After (✅):**
```php

if (PathHelper::isFileMissing($path)) { return false; }

if (PathHelper::isDirMissing($dir)) { mkdir($dir, 0755, true); }
$isBaseDirFailed = !self::makeDirectory($dir);

if ($isBaseDirFailed) { throw new Exception('...'); }
```

**Prevention:** Use `PathHelper` and `BooleanHelpers` guard functions. See [no-negatives.md](./12-no-negatives.md).

---

### Issue #10 — Go Enum `String()` Using `!v.IsValid()`

**Scope:** 21 Go enum packages
**Root Cause:** Every enum's `String()` method used `if !v.IsValid()` instead of the positive `if v.IsInvalid()`.

**Before (❌):**
```go
func (v Variant) String() string {
    if !v.IsValid() {
        return variantLabels[Invalid]
    }

    return variantLabels[v]
}
```

**After (✅):**
```go
func (v Variant) String() string {
    if v.IsInvalid() {
        return variantLabels[Invalid]
    }

    return variantLabels[v]
}
```

**Prevention:** Use `IsInvalid()` (already defined on all enums) instead of `!IsValid()`.

---

### Issue #11 — Go `!pathutil.IsDir()` Negation

**Scope:** 4 call sites in `git/service.go`
**Root Cause:** No `IsDirMissing()` counterpart existed in `pathutil` package.

**Before (❌):**
```go
if !pathutil.IsDir(gitDir) {
    return apperror.FailNew[StatusResult](apperror.ErrGitNotRepo, "directory is not a git repository")
}
```

**After (✅):**
```go
// Added to pathutil package:
func IsDirMissing(path string) bool { return !IsDir(path) }

// Usage:
if pathutil.IsDirMissing(gitDir) {
    return apperror.FailNew[StatusResult](apperror.ErrGitNotRepo, "directory is not a git repository")
}
```

**Prevention:** Every boolean function must have a positive counterpart. When `IsX()` exists and `!IsX()` is used 2+ times, add `IsXMissing()` or `IsXAbsent()`.

---

## Enum Standard Violations

### Issue #12 — Magic Strings Instead of Enum Values

**Scope:** ~40 PHP call sites across 15 files
**Root Cause:** Hardcoded strings used where enum values should be referenced.

**Before (❌):**
```php
if ($status === 'success') { ... }
$scope = 'all';
register_rest_route($ns, '/upload', ['methods' => 'POST', ...]);
if (current_user_can('manage_options')) { ... }
```

**After (✅):**
```php

if ($status->isEqual(StatusType::Success)) { ... }
$scope = SnapshotScopeType::All->value;
register_rest_route($ns, '/upload', ['methods' => HttpMethodType::Post->value, ...]);

if (current_user_can(CapabilityType::ManageOptions->value)) { ... }
```

**Prevention:** Zero magic strings policy. Every domain string must reference an enum case.

---

### Issue #13 — Go Enum Labels Using snake_case

**Scope:** Go enum packages (pre-standard)
**Root Cause:** Early Go enums used snake_case or lowercase strings in `variantLabels`.

**Before (❌):**
```go
var variantLabels = [...]string{
    Invalid:  "invalid",
    PerTable: "per_table",
    SingleDb: "single_db",
}
```

**After (✅):**
```go
var variantLabels = [...]string{
    Invalid:  "Invalid",
    PerTable: "PerTable",
    SingleDb: "SingleDb",
}
```

**Prevention:** `variantLabels` must use PascalCase strings matching the constant names. Protocol-driven enums (`content_type`, `endpoint`, `header`, `response_key`, `response_message`) are exempt.

**v2.1.0 Update (2026-02-25):** The dual-table pattern (`variantStrings` + `variantLabels`) is now fully deprecated. All enums must use a single `variantLabels` table with PascalCase values. `Label()` delegates to `String()`. `Parse()` must use `strings.EqualFold()` for case-insensitive matching. Migration notices added to `spec/08`, `spec/10`, `spec/13` enum architecture files.

---

### Issue #14 — PHP Raw `===` for Enum Comparison

**Scope:** All PHP enum usage sites
**Root Cause:** Code used `$status === StatusType::Success` instead of `$status->isEqual(StatusType::Success)`.

**Before (❌):**
```php
if ($status === StatusType::Success) { ... }
if ($level === LogLevelType::Error || $level === LogLevelType::Warn) { ... }
```

**After (✅):**
```php

if ($status->isEqual(StatusType::Success)) { ... }

if ($level->isErrorOrWarn()) { ... }
```

**Prevention:** `isEqual()` is mandatory on all backed enums. Raw `===` is forbidden for enum comparison.

---

## Formatting Violations

### Issue #15 — Missing Blank Line Before `return`/`throw`

**Scope:** 22 violations across 12 PHP files
**Root Cause:** Multi-statement blocks ended with `return`/`throw` without a preceding blank line.

**Before (❌):**
```php
$result = $this->compute($data);
return $result;
```

**After (✅):**
```php
$result = $this->compute($data);

return $result;
```

**Prevention:** Rule 4 in [code-style.md](./04-code-style/00-overview.md). Exception: if `return`/`throw` is the only statement, no blank line needed.

---

### Issue #16 — Functions Exceeding Parameter Limits

**Scope:** 51 signatures across 33 PHP files
**Root Cause:** Functions with >2 parameters were on a single line instead of one-per-line.

**Before (❌):**
```php
public function processUpload(string $path, string $slug, int $postId, bool $isActive): array {
```

**After (✅):**
```php
public function processUpload(
    string $path,
    string $slug,
    int $postId,
    bool $isActive,
): array {
```

**Prevention:** Rule 9a — signatures with >2 parameters must be one-per-line with trailing comma.

---

### Issue #17 — Nested `if` Blocks

**Scope:** Multiple files across PHP and Go
**Root Cause:** Developer habit of nesting conditions instead of using early returns.

**Before (❌):**
```php
if ($request !== null) {
    if ($request->hasParam('file')) {
        if ($this->isValidFile($request->getParam('file'))) {
            $this->process($request);
        }
    }
}
```

**After (✅):**
```php

if ($request === null) {
    return;
}

$hasValidFile = $request->hasParam('file')
    && $this->isValidFile($request->getParam('file'));

if ($hasValidFile) {
    $this->process($request);
}
```

**Prevention:** Rule 2 / Rule 7 — zero nested `if`, absolute ban.

---

## Type Safety Violations

### Issue #18 — Go `interface{}` / `any` in Business Logic

**Scope:** Multiple Go service files
**Root Cause:** Using type erasure (`map[string]any`, `interface{}`) instead of concrete structs.

**Before (❌):**
```go
func ProcessData(data interface{}) interface{} { ... }
func FetchResults() (any, error) { ... }
```

**After (✅):**
```go
func ProcessData(data PluginDetails) apperror.Result[PluginSummary] { ... }
func FetchResults[T any]() apperror.Result[T] { ... }
```

**Prevention:** Zero `any`/`interface{}` in exported APIs. Only permitted in SQL args, logger fields, and third-party library boundaries.

---

### Issue #19 — Go Raw `(T, error)` Returns from Services

**Scope:** All Go service methods
**Root Cause:** Service methods returned raw tuples instead of typed result wrappers.

**Before (❌):**
```go
func (s *PluginService) GetById(ctx context.Context, id int64) (*Plugin, error) { ... }
```

**After (✅):**
```go
func (s *PluginService) GetById(ctx context.Context, id int64) apperror.Result[Plugin] { ... }
```

**Prevention:** All service methods must return `apperror.Result[T]`, `apperror.ResultSlice[T]`, or `apperror.ResultMap[K,V]`.

---

## Summary Statistics

| Category | Issues Found | Issues Fixed | Files Affected |
|----------|-------------|-------------|----------------|
| Naming | 4 | 4 | ~80 |
| Database | 4 | 4 | ~25 |
| Boolean | 3 | 3 | ~55 |
| Enum | 3 | 3 | ~40 |
| Formatting | 3 | 3 | ~45 |
| Type Safety | 2 | 2 | ~20 |
| **Total** | **19** | **19** | **~265** |

---

*Issues and fixes log v1.0.0 — 2026-02-23*
