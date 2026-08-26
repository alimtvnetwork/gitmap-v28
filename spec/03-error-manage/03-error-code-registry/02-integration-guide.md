# Error Code Registry - Integration Guide


**Last Updated:** 2026-04-16  

> **Version:** 1.0.0  
> **Purpose:** How to integrate standardized error codes into your project

---

## Quick Start

### 1. Claim Your Project Prefix

Add your project to `01-registry.md`:

```markdown
| `XX` | Your Project | X000-X999 | Your Team | ✅ Active |
```

### 2. Define Your Error Codes

Create error definitions in your project:

```go
// Go: errors/codes.go
package errors

const (
    // SM-000: Initialization
    ErrConfigMissing = "SM-000-01"
    ErrConfigInvalid = "SM-000-02"
    
    // SM-400: Business Logic
    ErrSpecParseError = "SM-400-01"
)
```

```typescript
// TypeScript: src/errors/codes.ts
export const ErrorCodes = {
  // SM-000: Initialization
  CONFIG_MISSING: "SM-000-01",
  CONFIG_INVALID: "SM-000-02",
  
  // SM-400: Business Logic
  SPEC_PARSE_ERROR: "SM-400-01",
} as const;
```

---

## Implementation Patterns

### Go Backend

```go
package errors

import "fmt"

type AppError struct {
    Code    string
    Message string
    Cause   error  // EXEMPTED: AppError internal cause (I-2)
}

func (e *AppError) Error() string {
    if e.Cause != nil {
        return fmt.Sprintf("%s: %s (%v)", e.Code, e.Message, e.Cause)
    }

    return fmt.Sprintf("%s: %s", e.Code, e.Message)
}

func New(code, message string) *AppError {
    return &AppError{Code: code, Message: message}
}

func Wrap(code, message string, cause error) *AppError {
    return &AppError{Code: code, Message: message, Cause: cause}
}

// Usage:
// return errors.New(ErrConfigMissing, "Configuration file not found")
// return errors.Wrap(ErrDbQuery, "Failed to fetch user", err)
```

### TypeScript Frontend

```typescript
// src/errors/AppError.ts
export class AppError extends Error {
  constructor(
    public readonly code: string,
    message: string,
    public readonly cause?: Error
  ) {
    super(`${code}: ${message}`);
    this.name = 'AppError';
  }
  
  toJson() {
    return {
      code: this.code,
      message: this.message,
      cause: this.cause?.message,
    };
  }
}

// Usage:
// throw new AppError(ErrorCodes.CONFIG_MISSING, "Configuration file not found");
```

### API Response Format

```typescript
interface ErrorResponse {
  success: false;
  error: {
    code: string;      // "SM-400-01"
    message: string;   // Human-readable message
    details?: unknown; // Optional additional context
  };
}

// Example response:
{
  "success": false,
  "error": {
    "code": "SM-400-01",
    "message": "Failed to parse specification file",
    "details": {
      "file": "spec.md",
      "line": 42
    }
  }
}
```

---

## Error Handling Best Practices

### 1. Log with Context

```go
log.Error().
    Str("code", err.Code).
    Str("file", filename).
    Int("line", lineNum).
    Err(err.Cause).
    Msg(err.Message)
```

### 2. User-Facing vs Internal Errors

```typescript
// Internal: Full details for logging
console.error(`[${error.code}] ${error.message}`, error.cause);

// User-facing: Safe message only
toast.error(getHumanMessage(error.code));
```

### 3. Error Code Lookup

```typescript
const errorMessages: Record<string, string> = {
  "GEN-100-01": "Please log in to continue",
  "GEN-100-02": "Your session has expired. Please log in again",
  "SM-400-01": "Unable to read the specification file",
};

function getHumanMessage(code: string): string {
  return errorMessages[code] ?? "An unexpected error occurred";
}
```

---

## PowerShell Integration

For PowerShell scripts, use exit codes:

```powershell

# Define exit codes

$ErrorCodes = @{
    Success = 0
    ConfigMissing = 9510
    BuildFailed = 9520
    BackendFailed = 9530
}

# Usage

if (-not (Test-Path $configPath)) {
    Write-Error "PS-9510-01: Configuration file not found"
    exit $ErrorCodes.ConfigMissing
}
```

---

## Validation

Use the JSON schema to validate error code format:

```bash

# Validate error code format

echo '{"code": "SM-400-01"}' | \
  npx ajv validate -s schemas/error-code.schema.json -d -
```

---

## Checklist

- [ ] Claimed project prefix in registry
- [ ] Created error constants file
- [ ] Implemented AppError class/struct
- [ ] Added to API response format
- [ ] Created user-facing message lookup
- [ ] Updated registry with new codes
