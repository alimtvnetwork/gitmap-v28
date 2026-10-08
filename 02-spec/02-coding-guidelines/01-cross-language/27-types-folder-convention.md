# Types Folder Convention & Common Type Aliases (AI Execution Prompt)

> **/goal** Standardize application type architecture by placing shared types, enums, and type aliases into dedicated single-type files within a centralized `types/` folder and defining reusable aliases for repeated generics.
> **/learn** Enforce the "one definition per file" rule matching the PascalCase type name, eliminate bundled multi-type junk drawers, establish domain-specific type aliases for generic Result specializations used 3+ times, and prohibit raw string primitives where typed enums exist.

## 🎯 Actionable CI/CD & Agent Checklist

- [ ] `/goal` Create a dedicated `types/` folder (or language equivalent) containing shared enums, type aliases, and contract models.
- [ ] `/learn` Never bundle unrelated types into a generic `types/common.go` file; enforce exactly one definition per file matching PascalCase naming.
- [ ] `/goal` Define reusable type aliases (e.g., `BoolResult`, `StringResult`) whenever a generic type (e.g., `Result[T]`) is specialized 3 or more times.
- [ ] `/learn` Use strongly-typed enums instead of raw magic strings or integers for HTTP methods, statuses, content types, and log levels.

. **CRITICAL AI INSTRUCTION:** This specification is an active AI execution directive. All code generated or modified must strictly follow the rules below.

**Version:** 3.2.0
**Updated:** 2026-04-16
**Applies to:** All languages (Go, TypeScript, PHP, Rust, C#)
**Source:** Derived from `appfault` package patterns, content type conventions, and real-world `riseup-asia-uploader` codebase

---

## 1. Principle

Every project must have a **`types/` folder** (or language equivalent) containing shared type definitions, enums, and type aliases. Each definition gets its **own file** — never bundle unrelated types together.

---

## 2. Folder Structure

```
types/
├── ContentType.go          # MIME content types
├── HttpMethod.go           # HTTP methods (Get, Post, Put, Delete)
├── HttpStatus.go           # HTTP status code groups
├── AppResults.go           # Common Result[T] aliases (BoolResult, StringResult)
├── SortDirection.go        # Asc, Desc
├── Environment.go          # Development, Staging, Production
└── LogLevel.go             # Debug, Info, Warn, Error, Fatal
```

**TypeScript equivalent:**
```
types/
├── ContentType.ts
├── HttpMethod.ts
├── HttpStatus.ts
├── AppResults.ts
├── SortDirection.ts
├── Environment.ts
└── LogLevel.ts
```

**PHP equivalent:**
```
types/
├── ContentType.php
├── HttpMethod.php
├── HttpStatus.php
├── SortDirection.php
├── Environment.php
└── LogLevel.php
```

---

## 3. Rules

### Rule 1: One Definition Per File

Each file contains exactly **one** enum, type alias group, or constant group. The filename matches the type name in PascalCase.

```go
// ❌ FORBIDDEN — multiple unrelated types in one file
// types/Common.go
type ContentType byte
type HttpMethod byte
type LogLevel byte
type SortDirection byte

// ✅ CORRECT — one type per file
// types/ContentType.go
type ContentType byte

// types/HttpMethod.go
type HttpMethod byte
```

### Rule 2: Type Aliases for Repeated Generics

When a generic type is used **3 or more times** with the same parameter, create a type alias.

#### Go

```go
// types/app_results.go (or types.go)
package types

import "coding-guidelines/common/pkg/result"

// Common Result aliases — use these instead of repeating Result[T]
type BoolResult = result.Result[bool]
type StringResult = result.Result[string]
type IntResult = result.Result[int]
type Int64Result = result.Result[int64]

// Domain-specific aliases (add as patterns emerge)
// type PluginResult = result.Result[*Plugin]
// type SiteResult = result.Result[*Site]
```

```go
// Usage — clean, readable, consistent
func (h *PluginHandler) EnablePlugin(siteId string, slug string) types.BoolResult {
    // ...
}

func (h *PluginHandler) GetName(siteId string) types.StringResult {
    // ...
}
```

#### TypeScript

```typescript
// types/AppResults.ts
import type { Result } from "@/lib/result"

export type BoolResult = Result<boolean>
export type StringResult = Result<string>
export type NumberResult = Result<number>
export type VoidResult = Result<void>
```

#### C#

```csharp
// Types/AppResults.cs
namespace MyApp.Types;

// Common Result aliases
using BoolResult = AppError.Result<bool>;
using StringResult = AppError.Result<string>;
using IntResult = AppError.Result<int>;
```

### Rule 3: Enums Over Constants for Finite Sets

If a value belongs to a **known, finite set**, use an enum — not string constants.

---

## 4. Common Type Definitions

### ContentType

```go
// types/ContentType.go
package types

type ContentType byte

const (
    ContentTypeJson ContentType = iota + 1
    ContentTypeXml
    ContentTypeFormData
    ContentTypeTextPlain
    ContentTypeOctetStream
)

var contentTypeLabels = map[ContentType]string{
    ContentTypeJson:        "application/json",
    ContentTypeXml:         "application/xml",
    ContentTypeFormData:    "multipart/form-data",
    ContentTypeTextPlain:   "text/plain",
    ContentTypeOctetStream: "application/octet-stream",
}

func (c ContentType) String() string {
    if label, isFound := contentTypeLabels[c]; isFound {
        return label
    }

    return "application/octet-stream"
}
```

```typescript
// types/ContentType.ts
export enum ContentType {
    Json = "application/json",
    Xml = "application/xml",
    FormData = "multipart/form-data",
    TextPlain = "text/plain",
    OctetStream = "application/octet-stream",
}
```

```php
// types/ContentType.php
enum ContentType: string {
    case Json = 'application/json';
    case Xml = 'application/xml';
    case FormData = 'multipart/form-data';
    case TextPlain = 'text/plain';
    case OctetStream = 'application/octet-stream';
}
```

### HttpMethod

```go
// types/HttpMethod.go
package types

type HttpMethod byte

const (
    HttpMethodGet HttpMethod = iota + 1
    HttpMethodPost
    HttpMethodPut
    HttpMethodPatch
    HttpMethodDelete
    HttpMethodHead
    HttpMethodOptions
)

var httpMethodLabels = map[HttpMethod]string{
    HttpMethodGet:     "GET",
    HttpMethodPost:    "POST",
    HttpMethodPut:     "PUT",
    HttpMethodPatch:   "PATCH",
    HttpMethodDelete:  "DELETE",
    HttpMethodHead:    "HEAD",
    HttpMethodOptions: "OPTIONS",
}

func (m HttpMethod) String() string {
    if label, isFound := httpMethodLabels[m]; isFound {
        return label
    }

    return "GET"
}
```

```typescript
// types/HttpMethod.ts
export enum HttpMethod {
    Get = "GET",
    Post = "POST",
    Put = "PUT",
    Patch = "PATCH",
    Delete = "DELETE",
    Head = "HEAD",
    Options = "OPTIONS",
}
```

```php
// types/HttpMethod.php
enum HttpMethod: string {
    case Get = 'GET';
    case Post = 'POST';
    case Put = 'PUT';
    case Patch = 'PATCH';
    case Delete = 'DELETE';
    case Head = 'HEAD';
    case Options = 'OPTIONS';
}
```

### HttpStatus

```go
// types/HttpStatus.go
package types

type HttpStatus byte

const (
    HttpStatusOk HttpStatus = iota + 1
    HttpStatusCreated
    HttpStatusBadRequest
    HttpStatusUnauthorized
    HttpStatusForbidden
    HttpStatusNotFound
    HttpStatusConflict
    HttpStatusInternalError
    HttpStatusServiceUnavailable
)

var httpStatusCodes = map[HttpStatus]int{
    HttpStatusOk:                 200,
    HttpStatusCreated:            201,
    HttpStatusBadRequest:         400,
    HttpStatusUnauthorized:       401,
    HttpStatusForbidden:          403,
    HttpStatusNotFound:           404,
    HttpStatusConflict:           409,
    HttpStatusInternalError:      500,
    HttpStatusServiceUnavailable: 503,
}

func (s HttpStatus) Code() int {
    if code, isFound := httpStatusCodes[s]; isFound {
        return code
    }

    return 500
}
```

```typescript
// types/HttpStatus.ts
export enum HttpStatus {
    Ok = 200,
    Created = 201,
    BadRequest = 400,
    Unauthorized = 401,
    Forbidden = 403,
    NotFound = 404,
    Conflict = 409,
    InternalError = 500,
    ServiceUnavailable = 503,
}
```

### SortDirection

```go
// types/SortDirection.go
package types

type SortDirection byte

const (
    SortAsc SortDirection = iota + 1
    SortDesc
)
```

```typescript
// types/SortDirection.ts
export enum SortDirection {
    Asc = "asc",
    Desc = "desc",
}
```

### Environment

```go
// types/Environment.go
package types

type Environment byte

const (
    EnvironmentDevelopment Environment = iota + 1
    EnvironmentStaging
    EnvironmentProduction
)
```

```typescript
// types/Environment.ts
export enum Environment {
    Development = "development",
    Staging = "staging",
    Production = "production",
}
```

---

## 5. Anti-Patterns

### ❌ Dumping All Types in One File

```go
// types/types.go — FORBIDDEN
type ContentType byte
type HttpMethod byte
type LogLevel byte
type SortDirection byte
type Environment byte
// ... 200 more lines
```

### ❌ Defining Types Inline Where Used

```go
// handlers/plugin.go — FORBIDDEN
type PluginStatus byte  // Should be in types/PluginStatus.go
```

### ❌ Using Raw Strings Instead of Type Enum

```go
// ❌ FORBIDDEN
req.Header.Set("Content-Type", "application/json")

// ✅ CORRECT
req.Header.Set("Content-Type", ContentTypeJson.String())
```

### ❌ Repeating Generic Specializations

```go
// ❌ FORBIDDEN — Result[bool] repeated 15 times across codebase
func Enable() result.Result[bool] { ... }
func Disable() result.Result[bool] { ... }
func Toggle() result.Result[bool] { ... }

// ✅ CORRECT — alias defined once in types.go or types/
func Enable() types.BoolResult { ... }
func Disable() types.BoolResult { ... }
func Toggle() types.BoolResult { ... }
```

---

## 6. Summary Checklist

```
□ types/ folder exists at project root (or language-appropriate location)
□ One type/enum/alias group per file — filename matches type name
□ Common enums defined: ContentType, HttpMethod, HttpStatus, Environment
□ Result[T] aliases created for any specialization used 3+ times
□ No raw string literals for content types, HTTP methods, or status codes
□ No inline type definitions in handler/service files
□ All enum files follow the language's naming convention (PascalCase.go, PascalCase.ts)
```

---

## 7. Cross-References

- [Magic Values & Immutability](./26-magic-values-and-immutability.md) — No magic strings, use enums
- [Code Mutation Avoidance](./18-code-mutation-avoidance.md) — Immutable type values
- [Strict Typing](./13-strict-typing.md) — No `any`/`interface{}`
- [Generic Return Types](./25-generic-return-types.md) — Result[T] patterns
- [Boolean Principles](./03-boolean-principles.md) — Boolean naming in type definitions

---

## Verification & Acceptance Criteria

_Auto-generated section — see `02-spec/02-coding-guidelines/97-acceptance-criteria.md` for the full criteria index._

### AC-CG-TYPE-027: Centralized Types Folder and Single-Definition File Architecture

**Given** Source code files, type definitions, domain models, and enum declarations across all languages.
**When** Coding guideline linters or architectural analyzers inspect the repository structure and type declarations.
**Then** Shared types reside within dedicated `types/` folders with one definition per file matching the PascalCase type name, generic Result specializations repeated 3+ times use centralized aliases, and no raw primitive strings replace declared enums, achieving exit code 0.

**Verification command:**
```bash
python 03-ai-scripts/05-guideline-autofixer.py 02-spec/02-coding-guidelines/01-cross-language --check-only
```
**Expected:** exit 0. Zero violations.
