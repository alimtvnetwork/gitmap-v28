# Phase 6 — Input Validation Patterns

> **Purpose:** Define how REST request parameters are validated, required fields enforced, and type-safe parsing performed using the `TypeCheckerTrait` and `PhpNativeType` enum. Every endpoint must reject invalid input early with clear, structured error responses.

---

## 6.1 Validation Philosophy

| Principle | Rule |
|-----------|------|
| **Fail fast** | Validate all inputs at the top of the private handler method, before any business logic |
| **Guard clauses** | Each validation is a standalone `if` → `return error`. No nested else chains |
| **Positive booleans** | Extract every check into a named `$has…` / `$is…` variable |
| **Structured errors** | All validation failures return a 400 envelope with a human-readable message |
| **No silent defaults** | Never silently substitute a default for a missing required field — reject it |

---

## 6.2 Standard Validation Flow

Every handler that accepts input follows this pattern:

```
private function executeCreateWidget(WP_REST_Request $request): WP_REST_Response
{
    // ── Step 1: Extract raw input ──
    $body = $request->get_json_params();

    // ── Step 2: Validate body exists and is an array ──
    $hasBody = ($body !== null && $this->isArray($body));

    if (!$hasBody) {
        return $this->validationError('Request body must be a JSON object', $request);
    }

    // ── Step 3: Validate required fields ──
    $name = $body['name'] ?? null;
    $hasName = ($name !== null && $this->isString($name));

    if (!$hasName) {
        return $this->validationError('Missing required field: name', $request);
    }

    $nameLength = mb_strlen($name);
    $isNameTooLong = ($nameLength > 200);

    if ($isNameTooLong) {
        return $this->validationError('Field "name" must not exceed 200 characters', $request);
    }

    // ── Step 4: Validate optional fields with type checks ──
    $priority = $body['priority'] ?? null;
    $hasPriority = ($priority !== null);

    if ($hasPriority) {
        $isPriorityValid = $this->isInteger($priority);

        if (!$isPriorityValid) {
            return $this->validationError('Field "priority" must be an integer', $request);
        }
    }

    // ── Step 5: Sanitise ──
    $sanitisedName = sanitize_text_field($name);
    $resolvedPriority = $hasPriority ? absint($priority) : 0;

    // ── Step 6: Business logic (all inputs are now safe) ──
    // ...
}
```

---

## 6.3 validationError() — Helper Method

Add this to `ResponseTrait` alongside `safeExecute()`:

```
/**
 * Return a 400 validation error with a descriptive message.
 * Validation messages are never gated by debug mode — they contain
 * no internal implementation details and help the caller fix their request.
 *
 * @param string          $message  Human-readable description of what is wrong
 * @param WP_REST_Request $request  The incoming request (for RequestedAt path)
 *
 * @return WP_REST_Response 400 envelope response
 */
protected function validationError(
    string $message,
    WP_REST_Request $request,
): WP_REST_Response {
    $this->fileLogger->warn('Validation failed: ' . $message, [
        'endpoint' => $request->get_route(),
        'method'   => $request->get_method(),
    ]);

    return EnvelopeBuilder::error($message, 400)
        ->setRequestedAt($request->get_route())
        ->toResponse();
}
```

### Response example

```json
{
  "Status": {
    "IsSuccess": false,
    "IsFailed": true,
    "Code": 400,
    "Message": "Missing required field: name",
    "Timestamp": "2026-04-08T10:00:00Z"
  },
  "Attributes": {
    "RequestedAt": "/my-plugin-api/v1/widgets",
    "TotalRecords": 0
  },
  "Results": []
}
```

No `Errors` key — validation errors never include stack traces, regardless of debug mode.

---

## 6.4 Field Extraction Patterns

### Required string field

```
$slug = $body['slug'] ?? null;
$hasSlug = ($slug !== null && $this->isString($slug));

if (!$hasSlug) {
    return $this->validationError('Missing required field: slug', $request);
}

$slug = sanitize_text_field($slug);
```

### Required integer field

```
$count = $body['count'] ?? null;
$hasCount = ($count !== null && $this->isInteger($count));

if (!$hasCount) {
    return $this->validationError('Missing required field: count (integer)', $request);
}

$count = absint($count);
```

### Required boolean field

```
$isActive = $body['is_active'] ?? null;
$hasIsActive = ($isActive !== null && $this->isBoolean($isActive));

if (!$hasIsActive) {
    return $this->validationError('Missing required field: is_active (boolean)', $request);
}
```

### Required array field

```
$siteKeys = $body['site_keys'] ?? null;
$hasSiteKeys = ($siteKeys !== null && $this->isArray($siteKeys));

if (!$hasSiteKeys) {
    return $this->validationError('Missing required field: site_keys (array)', $request);
}

$isEmpty = (count($siteKeys) === 0);

if ($isEmpty) {
    return $this->validationError('Field "site_keys" must not be empty', $request);
}
```

### Optional field with default

```
$limit = $body['limit'] ?? null;
$hasLimit = ($limit !== null);

if ($hasLimit) {
    $isLimitValid = $this->isInteger($limit);

    if (!$isLimitValid) {
        return $this->validationError('Field "limit" must be an integer', $request);
    }
}

$resolvedLimit = $hasLimit ? absint($limit) : 10;
```

### Enum-constrained field

```
$status = $body['status'] ?? null;
$hasStatus = ($status !== null && $this->isString($status));

if (!$hasStatus) {
    return $this->validationError('Missing required field: status', $request);
}

$matchedStatus = StatusType::tryFrom($status);
$isValidStatus = ($matchedStatus !== null);

if (!$isValidStatus) {
    $allowed = implode(', ', array_column(StatusType::cases(), 'value'));

    return $this->validationError(
        "Invalid status '{$status}'. Allowed values: {$allowed}",
        $request,
    );
}
```

---

## 6.5 URL Parameter Validation

For GET endpoints that use URL path or query parameters:

### Path parameter (from route regex)

```
// Route: /widgets/(?P<id>[a-zA-Z0-9-]+)
$id = $request->get_param('id');
$hasId = ($id !== null && $this->isString($id));

if (!$hasId) {
    return $this->validationError('Missing path parameter: id', $request);
}

$id = sanitize_text_field($id);
```

### Query parameter (optional filter)

```
$statusFilter = $request->get_param('status');
$hasStatusFilter = ($statusFilter !== null && $this->isString($statusFilter));

if ($hasStatusFilter) {
    $matchedFilter = StatusType::tryFrom($statusFilter);
    $isValidFilter = ($matchedFilter !== null);

    if (!$isValidFilter) {
        return $this->validationError("Invalid filter status: {$statusFilter}", $request);
    }
}
```

### Pagination parameters

```
$page = $request->get_param('page');
$perPage = $request->get_param('per_page');

$resolvedPage = ($page !== null && $this->isNumeric($page)) ? max(1, absint($page)) : 1;
$resolvedPerPage = ($perPage !== null && $this->isNumeric($perPage)) ? min(100, max(1, absint($perPage))) : 20;
```

Pagination parameters are an exception to the "no silent defaults" rule — they should fall back to sensible defaults because every list endpoint needs them.

---

## 6.6 Nested Object Validation

When the request body contains nested objects:

```
$config = $body['config'] ?? null;
$hasConfig = ($config !== null && $this->isArray($config));

if (!$hasConfig) {
    return $this->validationError('Missing required field: config (object)', $request);
}

// Validate nested fields with dotted context in error messages
$retryCount = $config['retry_count'] ?? null;
$hasRetryCount = ($retryCount !== null && $this->isInteger($retryCount));

if (!$hasRetryCount) {
    return $this->validationError('Missing required field: config.retry_count (integer)', $request);
}

$isRetryCountInRange = ($retryCount >= 0 && $retryCount <= 10);

if (!$isRetryCountInRange) {
    return $this->validationError('Field "config.retry_count" must be between 0 and 10', $request);
}
```

Use dotted field paths (`config.retry_count`) in error messages so the caller knows exactly which nested field is invalid.

---

## 6.7 Sanitisation Reference

After validation passes, sanitise all string values before use:

| Input type | Sanitisation function | Notes |
|------------|----------------------|-------|
| Plain text | `sanitize_text_field($value)` | Strips tags, removes extra whitespace |
| File name | `sanitize_file_name($value)` | Removes special chars unsafe for filenames |
| Slug | `sanitize_title($value)` | Lowercase, hyphens only |
| URL | `esc_url_raw($value)` | For storage; use `esc_url()` for display |
| Integer | `absint($value)` | Absolute integer (non-negative) |
| Integer (allow negative) | `intval($value)` | Preserves sign |
| HTML content | `wp_kses($value, $allowedTags)` | Whitelist-based HTML sanitisation |

### Rule: Sanitise after validation, not instead of validation

```
// ❌ WRONG — sanitise without validating type first
$name = sanitize_text_field($body['name'] ?? '');

// ✅ CORRECT — validate type, reject if invalid, then sanitise
$name = $body['name'] ?? null;
$hasName = ($name !== null && $this->isString($name));

if (!$hasName) {
    return $this->validationError('Missing required field: name', $request);
}

$name = sanitize_text_field($name);
```

---

## 6.8 Validation Checklist for New Endpoints

When adding a new endpoint, verify:

- [ ] Body is validated as array (for POST/PUT)
- [ ] Every required field has a `$has…` guard with early return
- [ ] Every field type is checked with `TypeCheckerTrait` methods
- [ ] String length limits are enforced where applicable
- [ ] Integer ranges are validated (min/max)
- [ ] Enum-constrained fields use `tryFrom()` with allowed-values error message
- [ ] Nested objects validate each level with dotted field paths in errors
- [ ] All strings are sanitised after validation
- [ ] Optional fields fall back to explicit defaults (not `null`)
- [ ] Validation errors use `$this->validationError()` (not `EnvelopeBuilder::error()` directly)
- [ ] No business logic executes before all validation passes
