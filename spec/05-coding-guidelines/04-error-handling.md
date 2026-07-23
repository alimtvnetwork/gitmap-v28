# Error Handling Patterns

## Overview

Consistent error handling across TypeScript and Go to ensure
predictable behavior, clear diagnostics, and clean control flow.

---

## 1. Fail Fast

Check preconditions at the top of a function. Return or throw early
rather than nesting the happy path inside conditions.

### TypeScript

```ts
function processUser(user: User | null): UserResult {
  if (user === null) {
    throw new Error("User is required");
  }

  const isActive = user.status === UserStatus.Active;
  if (isActive === false) {
    throw new Error("User must be active");
  }

  return buildResult(user);
}
```

### Go

```go
func processUser(user *User) (*UserResult, error) {
    if user == nil {
        return nil, fmt.Errorf("user is required")
    }

    if user.Status == StatusInactive {
        return nil, fmt.Errorf("user must be active")
    }

    return buildResult(user), nil
}
```

---

## 2. Go — Check Errors Immediately

Always handle the error on the very next line after a call.

```go
// ✅ Correct
data, err := readFile(path)
if err != nil {
    return fmt.Errorf("reading %s: %w", path, err)
}

// ❌ Wrong — error checked later
data, err := readFile(path)
processed := transform(data) // may panic if data is nil
if err != nil { ... }
```

---

## 3. Go — Wrap Errors with Context

Use `fmt.Errorf` with `%w` to add context while preserving the
original error for `errors.Is` / `errors.As`.

```go
result, err := db.Query(query)
if err != nil {
    return fmt.Errorf("querying repos: %w", err)
}
```

---

## 4. TypeScript — Typed Error Handling

Use discriminated unions or custom error classes instead of
generic `catch (e)`.

```ts
class ValidationError extends Error {
  constructor(
    public readonly field: string,
    message: string,
  ) {
    super(message);
    this.name = "ValidationError";
  }
}

function handleError(error: unknown): void {
  const isValidationError = error instanceof ValidationError;
  if (isValidationError) {
    showFieldError(error.field, error.message);

    return;
  }

  showGenericError("An unexpected error occurred");
}
```

---

## 5. Never Swallow Errors

Every `catch` block must either re-throw, log, or handle meaningfully.
Empty catch blocks are prohibited.

```ts
// ❌ Wrong
try { riskyOp(); } catch (e) {}

// ✅ Correct
try {
  riskyOp();
} catch (error) {
  console.error("Operation failed:", error);
}
```

---

## 6. Go — No `panic` for Expected Conditions

Reserve `panic` for truly unrecoverable programmer errors.
All expected failures use `error` returns.

---

## 7. Mandatory Path Context in File Errors (Code Red Rule)

Every error message involving a file or directory path **must** include:

1. The exact resolved path
2. The operation attempted
3. The specific failure reason

### Format

```
Error: [message] at [path]: [error] (operation: [op], reason: [reason])
```

### Examples

```go
// ✅ Correct
return fmt.Errorf("Error: failed to read config at %s: %v (operation: read, reason: permission denied)", path, err)

// ❌ Wrong — no path, no operation
return fmt.Errorf("config not found")
```

Generic "file not found" messages without paths are **prohibited by convention**.

See: `spec/02-app-issues/27-error-management-file-path-and-missing-file-code-red-rule.md`

---

## 8. Format Verb Compliance

All `fmt.Fprintf`, `fmt.Printf`, and `fmt.Errorf` calls **must** match their format verb count to the argument count. Enforce via `go vet` in CI. Mismatches cause runtime panics.

---

## References

- Code Quality Improvement: `spec/05-coding-guidelines/01-code-quality-improvement.md`
- Go Code Style: `spec/05-coding-guidelines/02-go-code-style.md`
