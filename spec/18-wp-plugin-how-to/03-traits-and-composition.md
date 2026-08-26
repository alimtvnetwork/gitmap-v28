# Phase 3 — Traits and Composition

> **Purpose:** Define how plugin functionality is decomposed into traits, how traits interact, and the composition model.

---

## 3.1 Why Traits

The plugin architecture uses traits instead of service classes for a specific reason: all endpoint handlers need access to the shared `$fileLogger` property and the `ResponseTrait` methods. Traits allow flat composition into the singleton Plugin class without constructor injection chains or a DI container.

### Trait vs Helper vs Enum decision table

| Question | Answer → Use |
|----------|-------------|
| Does it handle a REST endpoint? | Trait |
| Does it need `$this->fileLogger`? | Trait |
| Is it stateless utility logic? | Helper (static class) |
| Is it a set of named constants? | Enum |
| Does it represent configuration? | Enum |

---

## 3.2 Trait Folder Structure

Traits are organised by feature domain. Each domain gets its own subfolder under `Traits/`.

```
Traits/
├── Auth/
│   └── AuthTrait.php              ← Permission checks
├── Core/
│   ├── ResponseTrait.php          ← safeExecute, error/success response builders
│   ├── StatusHandlerTrait.php     ← /status endpoint
│   └── PluginInventoryTrait.php   ← /plugins endpoint
├── Route/
│   └── RouteRegistrationTrait.php ← All route registration
├── Upload/
│   └── UploadHandlerTrait.php     ← File upload handling
├── Activate/
│   ├── ActivateHandlerTrait.php   ← Plugin activation logic
│   └── DeactivateEndpointTrait.php
├── Deactivate/
│   └── DeactivateHandlerTrait.php
├── Log/
│   ├── LogStatusTrait.php
│   ├── LogRotationStatusTrait.php
│   ├── LogClearingTrait.php
│   ├── LogEmailTrait.php
│   ├── LogRetrievalTrait.php
│   └── LogDedupRegistryTrait.php
└── Debug/
    └── DebugRoutesTrait.php
```

### Naming convention

- Trait file and class name always end with `Trait`
- Subfolder name matches the feature domain (PascalCase)
- One trait per file, one endpoint handler per trait (preferred)

---

## 3.3 Trait Anatomy

Every handler trait follows this exact structure:

```
namespace PluginName\Traits\FeatureDomain;

if (!defined('ABSPATH')) {
    exit;
}

use Throwable;
use WP_REST_Request;
use WP_REST_Response;

use PluginName\Enums\EndpointType;
use PluginName\Enums\PluginConfigType;
use PluginName\Helpers\EnvelopeBuilder;

trait SomeHandlerTrait
{
    public function handleSomeEndpoint(WP_REST_Request $request): WP_REST_Response {
        return $this->safeExecute(
            fn() => $this->executeSomeLogic($request),
            'some-endpoint',
        );
    }

    private function executeSomeLogic(WP_REST_Request $request): WP_REST_Response {
        // Business logic here
        // Always use $this->fileLogger for logging
        // Always use EnvelopeBuilder for responses
        // Always use enum values, never string literals

        return EnvelopeBuilder::success()
            ->setRequestedAt($requestPath)
            ->setSingleResult($data)
            ->toResponse();
    }
}
```

### Key observations

| Pattern | Reason |
|---------|--------|
| `$this->safeExecute()` wraps every public handler | Guarantees Throwable is caught and logged |
| Private method contains actual logic | Separates error boundary from business logic |
| `$this->fileLogger` is used directly | Available because Plugin.php composes the trait |
| Enum values for all strings | No magic strings in handler code |

---

## 3.4 The ResponseTrait — Error Boundary

The `ResponseTrait` is the most critical infrastructure trait. It provides:

### safeExecute()

Wraps any callable with comprehensive error handling:

1. Calls the callback inside `try`
2. On `Throwable`: logs to PHP `error_log()` with full stack trace, logs via `FileLogger::logException()`, returns a structured error response

This is the **only** place where endpoint-level exceptions are caught. Individual handler traits must **not** have their own try-catch blocks around the entire handler — they delegate to `safeExecute()`.

### errorResponse()

Creates a structured error response with:
- A log entry including backtrace
- An envelope-formatted WP_REST_Response
- Optional exception details (stack trace frames)

### successResponse()

Creates a structured success response with the standard envelope format.

### Fallback safety

If `EnvelopeBuilder` cannot be loaded (autoloader failure), `ResponseTrait` has an inline fallback that builds the same envelope structure manually. This ensures the plugin never returns a bare PHP error.

---

## 3.5 The RouteRegistrationTrait

Route registration follows a grouped, fault-tolerant pattern:

1. Define a `$safeRegister` closure that wraps `register_rest_route()` in a try-catch
2. Group routes by domain (core, logs, machines, debug)
3. Each group is registered in its own try-catch — a failure in one group does not prevent other groups from registering
4. After all groups, log a summary: total registered, total failed, which groups failed

### Route definition pattern

Routes are always defined using enum values:

| Component | Source |
|-----------|--------|
| Namespace | `PluginConfigType::apiFullNamespace()` |
| Path | `EndpointType::SomePath->route()` (prepends `/`) |
| HTTP method | `HttpMethodType::Get->value` |
| Callback | `[$this, 'handleMethodName']` |
| Permission | `[$this, 'checkPluginPermission']` |

---

## 3.6 The AuthTrait

Authentication uses WordPress Application Passwords via Basic Auth:

1. Resolve the Authorization header from the request (check `$request->get_header()`, then `$_SERVER` fallbacks)
2. Decode Base64 credentials
3. Call `wp_authenticate_application_password()`
4. Set the current user with `wp_set_current_user()`
5. Check capabilities using `current_user_can()`

### Permission levels

| Method | Who can access |
|--------|---------------|
| `checkPluginPermission` | Users with `activate_plugins` capability |
| `checkStatusPermission` | Any authenticated user |

All authentication failures return a `WP_Error` with appropriate HTTP status codes from `HttpStatusType` and error codes from `WpErrorCodeType`.

---

## 3.7 Trait Composition in Plugin.php

The Plugin class composes traits via `use` statements at the top of the class body. The order should group related traits:

1. Auth traits
2. Route traits
3. Core infrastructure traits (Response, Status)
4. Feature-domain traits (Upload, Activate, Deactivate, etc.)
5. Logging management traits
6. Debug traits

The Plugin class must **not** override or re-implement any method defined in a trait. If a trait method needs customisation, create a new trait or modify the existing one.

---

## 3.8 TypeCheckerTrait — Safe Type Validation

The `TypeCheckerTrait` wraps the verbose `gettype() === PhpNativeType::...->value` pattern into clean, readable methods. This exists because the syntax validator blocks native PHP type-check functions like `is_array()` (see Phase 2, §2.3).

### Location

```
Traits/
└── Core/
    └── TypeCheckerTrait.php
```

### Full implementation

```
namespace PluginName\Traits\Core;

if (!defined('ABSPATH')) {
    exit;
}

use PluginName\Enums\PhpNativeType;

/**
 * TypeCheckerTrait — Syntax-validator-safe type checking.
 *
 * Replaces blocked native functions (is_array, is_string, etc.)
 * with gettype()-based checks via the PhpNativeType enum.
 *
 * @package PluginName\Traits\Core
 * @since   1.0.0
 */
trait TypeCheckerTrait
{
    protected function isArray(mixed $value): bool
    {
        return gettype($value) === PhpNativeType::PhpArray->value;
    }

    protected function isString(mixed $value): bool
    {
        return gettype($value) === PhpNativeType::PhpString->value;
    }

    protected function isInteger(mixed $value): bool
    {
        return gettype($value) === PhpNativeType::PhpInteger->value;
    }

    protected function isFloat(mixed $value): bool
    {
        return gettype($value) === PhpNativeType::PhpDouble->value;
    }

    protected function isBoolean(mixed $value): bool
    {
        return gettype($value) === PhpNativeType::PhpBoolean->value;
    }

    protected function isObject(mixed $value): bool
    {
        return gettype($value) === PhpNativeType::PhpObject->value;
    }

    protected function isNull(mixed $value): bool
    {
        return gettype($value) === PhpNativeType::PhpNull->value;
    }

    protected function isNumeric(mixed $value): bool
    {
        $type = gettype($value);
        $isInt = ($type === PhpNativeType::PhpInteger->value);
        $isFloat = ($type === PhpNativeType::PhpDouble->value);

        return $isInt || $isFloat;
    }

    protected function isScalar(mixed $value): bool
    {
        $type = gettype($value);
        $isString = ($type === PhpNativeType::PhpString->value);
        $isInt = ($type === PhpNativeType::PhpInteger->value);
        $isFloat = ($type === PhpNativeType::PhpDouble->value);
        $isBool = ($type === PhpNativeType::PhpBoolean->value);

        return $isString || $isInt || $isFloat || $isBool;
    }
}
```

### Usage in a handler trait

```
trait SomeHandlerTrait
{
    use TypeCheckerTrait;

    private function executeSomeLogic(WP_REST_Request $request): WP_REST_Response
    {
        $body = $request->get_json_params();
        $isValidBody = $this->isArray($body);

        if (!$isValidBody) {
            return EnvelopeBuilder::error('Request body must be a JSON object', 400)
                ->setRequestedAt($request->get_route())
                ->toResponse();
        }

        $name = $body['name'] ?? null;
        $hasName = ($name !== null && $this->isString($name));

        if (!$hasName) {
            return EnvelopeBuilder::error('Missing required field: name', 400)
                ->setRequestedAt($request->get_route())
                ->toResponse();
        }

        // ... business logic
    }
}
```

### When to use TypeCheckerTrait vs PhpNativeType::matches()

| Context | Use |
|---------|-----|
| Inside a class that composes traits (Plugin, handlers) | `$this->isArray($var)` via `TypeCheckerTrait` |
| Inside a static helper class | `PhpNativeType::PhpArray->matches($var)` |
| Inside an enum method | `PhpNativeType::PhpArray->matches($var)` |

### Composition in Plugin.php

Add `TypeCheckerTrait` in the Core infrastructure group:

```
use PluginName\Traits\Core\TypeCheckerTrait;

class Plugin
{
    // Auth traits
    use AuthTrait;
    // Route traits
    use RouteRegistrationTrait;
    // Core infrastructure traits
    use ResponseTrait;
    use TypeCheckerTrait;    // ← add here
    use StatusHandlerTrait;
    // ... feature-domain traits
}
```
