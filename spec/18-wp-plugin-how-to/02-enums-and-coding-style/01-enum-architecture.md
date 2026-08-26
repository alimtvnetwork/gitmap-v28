# Phase 2 — Enums and Coding Style

> **Purpose:** Define how enums eliminate magic strings and how code must be written for consistency.
> **Cross-reference:** [`../../02-coding-guidelines/01-cross-language/04-code-style/00-overview.md`](../../02-coding-guidelines/01-cross-language/04-code-style/00-overview.md) for all formatting rules (R1–R13).

---

## 2.1 Enum Architecture

Enums are the **primary tool for eliminating magic strings**. Every string or integer literal that appears more than once, or that represents a domain concept, must be captured in a backed enum.

### Standard enum categories

Every plugin should define at minimum these enum types:

| Enum | Backing type | Purpose |
|------|-------------|---------|
| `PluginConfigType` | `string` | Plugin identity: slug, name, version, API namespace, log prefix |
| `EndpointType` | `string` | REST API endpoint path fragments |
| `HttpMethodType` | `string` | HTTP verbs (GET, POST, PUT, DELETE) |
| `HttpStatusType` | `int` | HTTP status codes (200, 400, 401, 403, 404, 500) |
| `HookType` | `string` | WordPress hook names (`rest_api_init`, `admin_menu`, etc.) |
| `LogLevelType` | `string` | Log severity levels (Debug, Info, Warn, Error) |
| `ResponseKeyType` | `string` | Standardised keys in API response envelopes (PascalCase) |
| `CapabilityType` | `string` | WordPress capabilities (`activate_plugins`, `manage_options`) |
| `WpErrorCodeType` | `string` | WP_Error code constants |
| `PathLogFileType` | `string` | Log file name fragments (`/info.log`, `/error.log`) |
| `PhpNativeType` | `string` | PHP type names returned by `gettype()` — see §2.3 |

### Enum file template

```
namespace PluginName\Enums;

if (!defined('ABSPATH')) {
    exit;
}

enum ExampleType: string
{
    case SomeName = 'some_value';

    public function isEqual(self $other): bool { return $this === $other; }
    public function isOtherThan(self $other): bool { return $this !== $other; }
    public function isAnyOf(self ...$others): bool { return in_array($this, $others, true); }
}
```

### Three standard comparison methods

Every enum must include these three methods. They provide a consistent, readable API for all comparisons:

| Method | Usage |
|--------|-------|
| `isEqual($other)` | Replace `$enum === SomeType::Value` |
| `isOtherThan($other)` | Replace `$enum !== SomeType::Value` |
| `isAnyOf(...$others)` | Replace `in_array($enum, [...], true)` |

### Per-case `is*()` methods

Every enum should provide an individual `is*()` method for **each case**. This makes calling code more readable:

```php
// ✅ Readable
if ($action->isUpload()) { ... }

// ❌ Verbose — avoid in calling code
if ($action->isEqual(ActionType::Upload)) { ... }
```

Implementation: each `is*()` delegates to `isEqual()`:

```php
public function isUpload(): bool { return $this->isEqual(self::Upload); }
public function isEnable(): bool { return $this->isEqual(self::Enable); }
```

### Group helper methods

When multiple cases share a domain concept, add a group helper:

| Pattern | Example |
|---------|---------|
| Prefix-based groups | `isSnapshot()` → `str_starts_with($this->value, 'Snapshot')` |
| Explicit groups | `isLifecycle()` → `isAnyOf(Enable, Disable, Delete)` |
| Computed values | `EndpointType::route()` → prepends `/` to the value |
| Builder patterns | `HookType::ajax($action)` → builds `wp_ajax_` prefix |

---

## 2.2 PluginConfigType — The Identity Enum

This is the most important enum. It centralises every identity value the plugin uses:

| Case | Example value | Used for |
|------|--------------|----------|
| `Slug` | `'my-plugin'` | Directory names, option keys, CSS classes |
| `ShortName` | `'MyPlugin'` | Log prefixes, class references |
| `Name` | `'My Plugin'` | Display name in admin UI |
| `Version` | `'1.0.0'` | Log lines, API responses, headers |
| `MinWpVersion` | `'5.6'` | Compatibility checks |
| `MinPhpVersion` | `'8.1'` | Compatibility checks |
| `ApiNamespace` | `'my-plugin-api'` | REST route namespace |
| `ApiVersion` | `'v1'` | REST route version segment |
| `LogPrefix` | `'[MyPlugin]'` | Fallback log prefix |
| `SettingsGroup` | `'my_plugin_settings'` | WordPress settings API group |

---

## 2.3 PhpNativeType — Syntax Validator Compatibility

Some deployment pipelines validate PHP files using `token_get_all()` before activation. The following native PHP functions tokenize as language constructs (`T_ARRAY`, `T_STRING_CAST`, etc.) and **will cause validation failures**:

| Blocked function | Raw replacement | Recommended (via TypeCheckerTrait) |
|-----------------|----------------|-------------------------------------|
| `is_array($var)` | `gettype($var) === PhpNativeType::PhpArray->value` | `$this->isArray($var)` |
| `is_string($var)` | `gettype($var) === PhpNativeType::PhpString->value` | `$this->isString($var)` |
| `is_int($var)` | `gettype($var) === PhpNativeType::PhpInteger->value` | `$this->isInteger($var)` |
| `is_float($var)` | `gettype($var) === PhpNativeType::PhpDouble->value` | `$this->isFloat($var)` |
| `is_bool($var)` | `gettype($var) === PhpNativeType::PhpBoolean->value` | `$this->isBoolean($var)` |
| `is_object($var)` | `gettype($var) === PhpNativeType::PhpObject->value` | `$this->isObject($var)` |
| `is_null($var)` | `gettype($var) === PhpNativeType::PhpNull->value` | `$this->isNull($var)` |

> **Best practice:** Always use the `TypeCheckerTrait` methods (right column) instead of the verbose `gettype()` calls. The raw replacement column is for reference only — showing what the trait does internally. See Phase 3, §3.8 for the full trait specification.

**Safe functions** (tokenize as T_STRING, no issues): `in_array()`, `array_merge()`, `array_filter()`, `array_map()`, `array_key_exists()`, `array_slice()`, `array_pop()`.

### PhpNativeType Enum — isMatches() Method

The enum itself also provides an `isMatches()` method for standalone use outside classes (e.g., in helpers):

```
enum PhpNativeType: string
{
    case PhpArray   = 'array';
    case PhpString  = 'string';
    case PhpInteger = 'integer';
    case PhpDouble  = 'double';
    case PhpBoolean = 'boolean';
    case PhpObject  = 'object';
    case PhpNull    = 'NULL';

    /** Check if a value's type matches this enum case. */
    public function isMatches(mixed $value): bool
    {
        return gettype($value) === $this->value;
    }
}
```

**Usage in helpers (static context where traits are unavailable):**
```
$isValid = PhpNativeType::PhpArray->isMatches($input);
```

---

## 2.4 Coding Style — Conditionals

### Rule: No negative conditions

Every `if` condition must be expressed as a **positive boolean**. Extract the condition into a named boolean variable.

| Forbidden | Required |
|-----------|----------|
| `if (!$result)` | `$isFailed = ($result === false);` then `if ($isFailed)` |
| `if ($x !== null)` | `$hasValue = ($x !== null);` then `if ($hasValue)` |
| `if (!empty($arr))` | `$hasItems = (!empty($arr));` then `if ($hasItems)` |
| `if (!$user->can())` | `$isUnauthorized = (!$user->can());` then `if ($isUnauthorized)` |

### Rule: No complex conditions

Never combine multiple conditions in an `if` statement. Extract each to a named boolean.

**Forbidden:**
```
if ($status === 200 && $body !== null && count($body) > 0) {
```

**Required:**
```
$isSuccessStatus = ($status === 200);
$hasBody = ($body !== null);
$hasItems = ($hasBody && count($body) > 0);

if ($isSuccessStatus && $hasItems) {
```

### Rule: No inline ternaries with conditions

Ternary expressions with non-trivial conditions must extract the condition to a named boolean first.

**Forbidden:**
```
$value = ($count > 0 && $isActive) ? $data : $default;
```

**Required:**
```
$shouldUseData = ($count > 0 && $isActive);
$value = $shouldUseData ? $data : $default;
```

---

## 2.5 Coding Style — Formatting

These rules apply to all PHP code. Full details in [`../../02-coding-guidelines/01-cross-language/04-code-style/00-overview.md`](../../02-coding-guidelines/01-cross-language/04-code-style/00-overview.md).

| Rule | Summary |
|------|---------|
| R1 | Always use braces, even for single-line `if`/`foreach`/`while` |
| R4 | Blank line before `return`/`throw` only when other statements precede it in the block |
| R5 | Blank line after closing `}` when followed by more code (except `else`/`catch`/`finally`) |
| R9a | Function signatures with 3+ parameters: one per line, trailing comma |
| R9b | Function calls with 3+ arguments: one per line, trailing comma |
| R9c | Array literals with 3+ items: one per line, trailing comma |
| R10 | Blank line before control structures when preceded by statements |
| R12 | No empty line after opening brace |
| R13 | No empty first line in source files (after `<?php`) |

---

## 2.6 Coding Style — Naming

| Element | Convention | Example |
|---------|-----------|---------|
| Enum files | `PascalCase.php` | `HttpStatusType.php` |
| Enum cases | `PascalCase` | `ServerError` |
| Trait files | `PascalCaseTrait.php` | `AuthTrait.php` |
| Helper files | `PascalCaseHelper.php` | `DateHelper.php` |
| Boolean variables | Positive name, `$is`/`$has`/`$should`/`$can` prefix | `$isValid`, `$hasData` |
| Private properties | `$camelCase` | `$fileLogger` |
| Constants | `UPPER_SNAKE_CASE` | `MAX_RETRIES` |
| Enum type suffix | Always end with `Type` | `EndpointType`, `HookType` |

---

## 2.7 PHPDoc Standards

Every file must start with a docblock that includes:

```
/**
 * ClassName — One-line description.
 *
 * @package PluginName\Namespace
 * @since   1.0.0
 */
```

Methods require a docblock only when their purpose is not immediately obvious from the name and signature.
