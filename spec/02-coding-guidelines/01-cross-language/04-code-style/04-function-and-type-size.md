# Function & Type Size Limits

> **Version:** 4.0.0  
> **Updated:** 2026-03-31  
> **Applies to:** PHP, TypeScript, Go  
> **Rules covered:** 6, 17

---

## Rule 6: Maximum 15 Lines Per Function — Extract Small Helpers

Every function/method body must be **15 lines or fewer** (excluding blank lines, comments, and the signature). If a function exceeds this limit, extract logic into small, well-named helper functions.

### 6a — Error Handling Lines Are Exempt

Error-wrapping lines, error-guard blocks (`if err != nil`), and `apperror.Wrap()` / `FailWrap()` chains **do not count** toward the 15-line limit. These are structural safety code, not business logic. Counting them would incentivize combining error handling into fewer lines, which **harms readability**.

```go
// ✅ This function is within the 15-line limit
// Error handling lines (marked with //*) are EXEMPT from count
func ProcessUpload(ctx context.Context, req UploadRequest) error {
    if err := validateUpload(req); err != nil {                           //* exempt

        return err                                                        //* exempt
    }                                                                     //* exempt

    result, err := executeUpload(ctx, req)

    if err != nil {                                                       //* exempt

        return apperror.Wrap(                                               //* exempt
            err, apperror.ErrUpload, "upload failed",                        //* exempt
        ).                                                                   //* exempt
            WithSiteId(req.SiteId).                                          //* exempt
            WithPath(req.Path)                                               //* exempt
    }                                                                     //* exempt

    return logAndRespond(ctx, result)
}
```

**Rationale:** Error wrapping must always be explicit — each `.WithX()` context value on its own line (see [Rule 11 — Method Chaining](./05-multi-line-formatting.md)). Never compress error handling to save line count.

### 6b — Why

- Short functions are easier to read, test, and debug
- Named helpers act as documentation — the function name describes intent
- Reduces cognitive load — each function does exactly one thing
- Makes code review faster — reviewers can understand each piece in isolation

### 6c — How to Flatten

| Problem | Solution |
|---------|----------|
| Long setup + logic + cleanup | Extract each phase into a helper |
| Multiple validation checks | Extract `validateRequest()` helper |
| Complex data transformation | Extract `transformPayload()` helper |
| Repeated patterns | Extract shared utility function |

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: 25+ line function
public function handleUpload($request) {
    $file = $request->get_param('file');
    $source = $request->get_param('source');
    // ... validation ...
    // ... processing ...
    // ... logging ...
    // ... response building ...
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

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Long function
const handleSubmit = async (data: FormData) => {
    // 20+ lines of validation, API call, state updates, toasts...
};

// ✅ REQUIRED: Decomposed
const handleSubmit = async (data: FormData) => {
    const validated = validateFormData(data);
    const result = await submitToApi(validated);
    updateLocalState(result);
    showSuccessToast(result.message);
};
```

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: Long function
func ProcessUpload(ctx context.Context, req UploadRequest) error {
    // 20+ lines...
}

// ✅ REQUIRED: Decomposed
func ProcessUpload(ctx context.Context, req UploadRequest) error {
    if err := validateUpload(req); err != nil {
        return err
    }

    result, err := executeUpload(ctx, req)

    if err != nil {
        return apperror.Wrap(
            err, apperror.ErrSyncCheck, "upload failed",
        )
    }

    return logAndRespond(ctx, result)
}
```

---

## Rule 17: Struct / Class Size — Maximum 120 Lines

Every struct, class, or interface definition (including its methods in the same file) should target **120 lines or fewer**. If a type exceeds this limit, split using file suffixes (e.g., `_crud.go`, `_helpers.go`, `_validation.go` in Go; separate concern files in TypeScript/PHP).

```go
// ── Go ───────────────────────────────────────────────────────

// ❌ FORBIDDEN: 200+ line struct file with all methods
// plugin.go — struct + 15 methods = 250 lines

// ✅ REQUIRED: Split by concern
// plugin.go           — struct + constructor (30 lines)
// plugin_crud.go      — database CRUD methods (80 lines)
// plugin_helpers.go   — private utility methods (60 lines)
```

```typescript
// ── TypeScript ───────────────────────────────────────────────

// ❌ FORBIDDEN: Single component file with 200+ lines of logic
// UserProfile.tsx — component + hooks + helpers + types

// ✅ REQUIRED: Extract concerns
// UserProfile.tsx           — component JSX only
// useUserProfile.ts         — custom hook with logic
// userProfileHelpers.ts     — utility functions
// userProfile.types.ts      — type definitions
```

```php
// ── PHP ──────────────────────────────────────────────────────

// ❌ FORBIDDEN: God class with 300+ lines
// PluginController.php — validation + CRUD + formatting

// ✅ REQUIRED: Split by responsibility
// PluginController.php      — route handlers only
// PluginValidator.php       — validation logic
// PluginFormatter.php       — response formatting
```

---

*Part of [Code Style](./00-overview.md) — Rules 6, 17*
