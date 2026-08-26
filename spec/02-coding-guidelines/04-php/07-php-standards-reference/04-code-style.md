# PHP Coding Standards — Braces, nesting, spacing, function size

> **Parent:** [PHP Coding Standards](./00-overview.md)  
> **Version:** 5.1.0  
> **Updated:** 2026-03-31

---

## Code Style — Braces, Nesting, Spacing & Function Size

> These rules apply across **all languages** (PHP, TypeScript, Go).  
> **Canonical source:** [Cross-Language Code Style](../../01-cross-language/04-code-style/00-overview.md) — this section repeats key rules with PHP-specific examples.

### Rule 1: Always use braces — no single-line returns

Every `if`, `for`, `foreach`, `while` block must use curly braces, even for single-statement bodies.

```php
// ❌ FORBIDDEN: Single-line return without braces
if ($this->initialized) return;
if ($error === null) return false;

// ✅ REQUIRED: Always use braces
if ($this->initialized) {
    return;
}

if ($error === null) {
    return false;
}
```

### Rule 2: Zero nested `if` — absolute ban

Nested `if` blocks are **absolutely forbidden** — zero tolerance, no exceptions. Flatten using early returns, combined conditions, or extracted helper functions. If a helper function already handles the null/empty check internally (e.g., `ErrorChecker::isFatalError()` already returns `false` for `null`), rely on it — don't wrap it in a redundant outer guard.

```php
// ❌ FORBIDDEN: Nested if — redundant null guard
if ($error !== null) {
    if (ErrorChecker::isFatalError($error)) {
        $this->logger->fatal($error);
    }
}

// ✅ REQUIRED: Flat — isFatalError() handles null internally
if (ErrorChecker::isFatalError($error)) {
    $this->logger->fatal($error);
}

// ✅ ALSO OK: Early return to flatten
if ($request === null) {
    return;
}

if ($request->has_param('file')) {
    $this->process($request);
}
```

### Rule 3: Extract complex conditions — no inline multi-part checks

When an `if` condition contains **two or more operators** (`&&`, `||`, `!`), it must be extracted into one of:

1. **A named boolean variable** (`$is_*` / `$has_*`) — for local, one-off checks
2. **A dedicated method/function** — for reusable or domain-meaningful checks
3. **A named constant** — for static flag combinations

The goal: every `if` reads as a **single intent**, not as implementation logic.

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: Inline multi-part condition
if ($error && in_array($error['type'], [E_ERROR, E_PARSE, E_CORE_ERROR], true)) {
    $this->logger->fatal($error);
}

// ✅ REQUIRED: Extracted into a dedicated method
if (ErrorChecker::isFatalError($error)) {
    $this->logger->fatal($error);
}

// ❌ FORBIDDEN: Inline extension check
if (!class_exists('PDO') || !extension_loaded('pdo_sqlite')) {
    return $this->envelope->error('SQLite not available', 500);
}

// ✅ REQUIRED: Extracted into ErrorChecker
if (ErrorChecker::isInvalidPdoExtension()) {
    return $this->envelope->error('SQLite not available', 500);
}

// ❌ FORBIDDEN: Combinable nested conditions left inline
if ($request !== null && $request->hasParam('file') && $request->getParam('file') !== '') {
    $this->process($request);
}

// ✅ REQUIRED: Named boolean for clarity
$hasFileParam = $request !== null
    && $request->hasParam('file')
    && $request->getParam('file') !== '';

if ($hasFileParam) {
    $this->process($request);
}
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Inline multi-part condition
if (response && response.status >= 400 && response.data?.code?.startsWith('E8')) {
    showDelegatedError(response);
}

// ✅ REQUIRED: Named boolean
const isDelegatedError = response != null
    && response.status >= 400
    && response.data?.code?.startsWith('E8');

if (isDelegatedError) {
    showDelegatedError(response);
}

// ✅ ALSO OK: Dedicated function for reusable checks
function isDelegatedError(res: ApiResponse | null): res is DelegatedErrorResponse {
    return res != null && res.status >= 400 && res.data?.code?.startsWith('E8');
}

if (isDelegatedError(response)) {
    showDelegatedError(response);
}
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: Inline multi-part condition
if err != nil && resp != nil && resp.StatusCode >= 400 {
    handleUpstreamError(resp)
}

// ✅ REQUIRED: Separate error guard, then named boolean for domain check
if err != nil {
    return apperror.Wrap[UpstreamResult](err, errors.ErrUpstream, "upstream call failed")
}

isUpstreamError := resp != nil && resp.StatusCode >= 400

if isUpstreamError {
    handleUpstreamError(resp)
}
```

#### When to use which extraction:

| Complexity | Extraction | Example |
|------------|-----------|---------|
| 2 conditions, used once | Named `$is*` / `$has*` variable | `$hasFileParam = $req !== null && $req->hasParam('file');` |
| 2+ conditions, used in multiple places | Dedicated method/function | `ErrorChecker::isFatalError($error)` |
| Static flag combination | Named constant | `const EDITABLE = 'PUT, PATCH';` |

### Rule 4: Blank line before `return` when preceded by other statements

If a block contains statements before `return`, insert **one blank line** before the `return`. If `return` is the **only statement**, no blank line is needed.

```php
// ❌ FORBIDDEN: No blank line before return
if (ErrorChecker::isInvalidPdoExtension()) {
    $this->logger->error('PDO/SQLite not available');
    return $this->envelope->error('SQLite support not available', 500);
}

// ✅ REQUIRED: Blank line separates logic from exit
if (ErrorChecker::isInvalidPdoExtension()) {
    $this->logger->error('PDO/SQLite not available');

    return $this->envelope->error('SQLite support not available', 500);
}

// ✅ OK: Return is the only statement — no blank line needed
if ($error === null) {
    return false;
}
```

### Rule 5: Blank line after closing `}` when followed by more code

If code continues after a closing `}` (i.e., not followed by another `}` or end of function), insert **one blank line** after it.

```php
// ❌ FORBIDDEN: No blank line after block when code follows
if ($this->initialized) {
    return;
}
$this->initialized = true;
add_action(Hook::Init->value, [$this, 'setup']);

// ✅ REQUIRED: Blank line after block when code follows
if ($this->initialized) {
    return;
}

$this->initialized = true;
add_action(Hook::Init->value, [$this, 'setup']);
```

---

### Rule 6: Maximum 15 lines per function

> **Canonical source:** [Cross-Language Code Style](../../01-cross-language/04-code-style/00-overview.md) — Rule 6

Every function/method body must be **15 lines or fewer** (excluding blank lines, comments, and the signature). Extract logic into small, well-named helper functions.

```php
// ❌ FORBIDDEN: 25+ line function
public function handleUpload($request) {
    // validation, processing, logging, response... all inline
}

// ✅ REQUIRED: Short top-level, helpers do the work
public function handleUpload(WP_REST_Request $request): WP_REST_Response {
    $params = $this->extractUploadParams($request);
    $this->validateUpload($params);
    $result = $this->processUpload($params);
    $this->logUpload($result);

    return $this->envelope->success($result);
}
```

---
