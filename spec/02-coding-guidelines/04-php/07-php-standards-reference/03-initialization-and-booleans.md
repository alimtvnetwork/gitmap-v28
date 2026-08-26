# PHP Coding Standards — Constructor rules, boolean logic, isDefined guards

> **Parent:** [PHP Coding Standards](./00-overview.md)  
> **Version:** 5.1.0  
> **Updated:** 2026-03-31

---

## Initialization — No WordPress Calls in Constructors

### Rule: Lazy initialization with HookType enum

Never call WordPress functions (`add_action`, `register_rest_route`, etc.) in class constructors. All hook registrations must use `HookType` enum cases:

```php
// ❌ FORBIDDEN: WordPress call in constructor + magic string
class MyPlugin {
    public function __construct() {
        add_action('init', [$this, 'setup']); // May fail if WP not loaded
    }
}

// ✅ REQUIRED: Lazy initialization with HookType enum
use RiseupAsia\Enums\HookType;

class MyPlugin {
    private bool $isInitialized = false;
    
    public function initialize() {
        if ($this->isInitialized) {
            return;
        }

        $this->isInitialized = true;
        add_action(HookType::Init->value, [$this, 'setup']);
    }
}
```

---


---

## Boolean Logic

### Rule: No raw negations — use positive guard functions

> **Canonical source:** [No Raw Negations](../../01-cross-language/12-no-negatives.md)

**Never use `!` on a function call in a condition.** Every negative check must be wrapped in a positively named guard function that reads as a single intent. See the canonical spec for the full cross-language rule and all guard function tables.

### Rule: Use semantic method names — no trivial wrapper helpers

Boolean checks must be self-documenting through **semantic method names** on the object itself. Trivial wrappers that merely restate native PHP operators are **prohibited** — they add indirection without clarity.

### Prohibited Trivial Wrappers (deprecated since 1.19.0)

The following methods from the legacy `BooleanHelpers` class are **deprecated and must not be used**. Use native PHP instead:

| ❌ Deprecated method | ✅ Native replacement |
|----------------------|----------------------|
| `BooleanHelpers::isFalsy($x)` | `!$x` |
| `BooleanHelpers::isTruthy($x)` | `(bool) $x` |
| `BooleanHelpers::isNull($x)` | `$x === null` |
| `BooleanHelpers::isSet($x)` | `$x !== null` |
| `BooleanHelpers::isEmpty($x)` | `empty($x)` |
| `BooleanHelpers::hasContent($x)` | `!empty($x)` |

```php
// ❌ FORBIDDEN: Trivial wrappers — use native PHP
if (BooleanHelpers::isFalsy($value)) { ... }
if (BooleanHelpers::isNull($config)) { ... }
if (BooleanHelpers::hasContent($name)) { ... }

// ✅ REQUIRED: Native PHP operators
if (!$value) { ... }
if ($config === null) { ... }
if (!empty($name)) { ... }
```

### Allowed Domain-Specific Helpers

The following `BooleanHelpers` methods **are allowed** because they encapsulate multi-step checks with safety guards (e.g., `empty()` + native function) that would be error-prone inline:

| Method | Semantics | Internal logic |
|--------|-----------|----------------|
| `isFuncExists($name)` | Function is available | `function_exists($name)` |
| `isFuncMissing($name)` | Function is not available | `!function_exists($name)` |
| `isClassExists($name)` | Class is available | `class_exists($name)` |
| `isClassMissing($name)` | Class is not available | `!class_exists($name)` |
| `isExtensionLoaded($name)` | PHP extension is loaded | `extension_loaded($name)` |
| `isExtensionMissing($name)` | PHP extension is not loaded | `!extension_loaded($name)` |
| `isDirExists($path)` | Directory exists | `!empty($path) && is_dir($path)` |
| `isDirMissing($path)` | Directory does not exist | `empty($path) \|\| !is_dir($path)` |
| `isDirWritable($path)` | Directory exists and is writable | `!empty($path) && is_dir($path) && is_writable($path)` |
| `isDirReadonly($path)` | Directory missing or not writable | `empty($path) \|\| !is_dir($path) \|\| !is_writable($path)` |
| `isFileExists($path)` | File exists | `!empty($path) && file_exists($path)` |
| `isFileMissing($path)` | File does not exist | `empty($path) \|\| !file_exists($path)` |
| `isDbConnected($db)` | DB object is connected | `$db !== null && $db->isConnected()` |
| `isDbDisconnected($db)` | DB object is not connected | `$db === null \|\| !$db->isConnected()` |

> **Why these are allowed:** Each combines a null/empty guard with a native function call — a pattern that is easy to get wrong inline. The semantic method name (`isDirMissing`) reads as a single intent.

### Semantic Object Methods

```php
// ❌ FORBIDDEN: Raw negation — easy to miss the "!"
if (!$plugin->isActive()) { ... }

// ✅ REQUIRED: Semantic inverse methods on the object
if ($plugin->isDisabled()) { ... }

// ✅ REQUIRED: Descriptive boolean variable names (Is/Has prefix)
if ($isValue) { ... }
if ($hasPermission) { ... }
```

### Guidelines

1. **Every `isX()` method should have a semantic inverse** (e.g., `isActive()` ↔ `isDisabled()`) rather than relying on `!isActive()`.
2. **Boolean variables must use `$is*` or `$has*` prefix** — never store a boolean in `$value` or `$result`.
3. **Never create new trivial wrapper helpers** — if the check is a single native operator (`!`, `empty()`, `=== null`), use PHP directly. Only create helpers for multi-step checks with safety guards.

---


---

## `isDefined()` and `isDefinedAndValid()` — Positive Null/Existence Guards

> **Cross-language parity:** [Go Standards §IsDefined](../../03-golang/04-golang-standards-reference/00-overview.md)

Raw `!== null` combined with negation or nested validity checks creates cognitive overhead. Use positive guard methods that express intent clearly.

### `isDefined()` — Value Existence Check

Returns `true` when the value has been set (is not null). Replaces `!== null` checks.

```php
// ❌ FORBIDDEN: Raw null check
if ($config !== null) {
    $this->applyConfig($config);
}

// ❌ FORBIDDEN: Double negation
if (!is_null($config)) {
    $this->applyConfig($config);
}

// ✅ REQUIRED: Positive existence check
if ($config->isDefined()) {
    $this->applyConfig($config);
}
```

### `isDefinedAndValid()` — Existence + Validity Combined

Returns `true` when the value exists AND passes its own validation rules. Replaces nested null+validity checks.

```php
// ❌ FORBIDDEN: Nested null + validity check
if ($config !== null) {
    if ($config->isValid()) {
        $this->applyConfig($config);
    }
}

// ❌ FORBIDDEN: Compound with null check
if ($config !== null && $config->isValid()) {
    $this->applyConfig($config);
}

// ✅ REQUIRED: Single positive guard
if ($config->isDefinedAndValid()) {
    $this->applyConfig($config);
}
```

### Implementation Pattern

Every class that can be null or absent should implement both methods:

```php
class Config
{
    public function isDefined(): bool
    {
        return true; // Non-null instance always returns true
    }

    public function isDefinedAndValid(): bool
    {
        return $this->validate()->isSafe();
    }
}

// For nullable wrappers / optional values:
class Optional
{
    private mixed $value;

    public function isDefined(): bool
    {
        return $this->value !== null;
    }

    public function isDefinedAndValid(): bool
    {
        return $this->value !== null && $this->value->isValid();
    }
}
```

### Guard Method Table

| Guard Method | Replaces | Description |
|-------------|----------|-------------|
| `isDefined()` | `!== null`, `isset()` | Value exists (not null) |
| `isDefinedAndValid()` | `!== null && isValid()` | Value exists AND passes validation |
| `isEmpty()` | `=== null` | No value set (absent) |
| `isInvalid()` | `!isValid()` | Value fails validation |

### On `DbResult<T>` (Already Built-in)

| Method | Meaning |
|--------|---------|
| `isDefined()` | Row was found (regardless of error state) |
| `isSafe()` | Row exists AND no error — equivalent to `isDefinedAndValid()` |
| `isEmpty()` | No row found |
| `hasError()` | Query failed |

### Real-World Example

```php
// Service layer — checking optional input
public function update(UpdateSiteInput $input): AppResult
{
    if ($input->config->isDefined()) {
        if ($input->config->isDefinedAndValid()) {
            $this->applyConfig($input->config);
        } else {
            return AppResult::fail('E3010', 'invalid site config');
        }
    }
    // ...
}
```

---
